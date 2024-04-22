/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/loadbalancer"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/provider"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	waitForClusterInfrastructureReadyDuration = 15 * time.Second
	TimeoutInstanceCreate                     = 5
	RetryIntervalInstanceStatus               = 10 * time.Second
)

// OpenStackMachineReconciler reconciles a OpenStackMachine object
type OpenStackMachineReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

func (r *OpenStackMachineReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.TODO()
	logger := r.Log.WithValues("namespace", req.Namespace, "openStackMachine", req.Name)

	// Fetch the OpenStackMachine instance.
	openStackMachine := &infrav1.OpenStackMachine{}
	err := r.Get(ctx, req.NamespacedName, openStackMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, openStackMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		logger.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("machine", machine.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		logger.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	if isPaused(cluster, openStackMachine) {
		logger.Info("OpenStackMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("cluster", cluster.Name)

	openStackCluster := &infrav1.OpenStackCluster{}

	openStackClusterName := client.ObjectKey{
		Namespace: openStackMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, openStackClusterName, openStackCluster); err != nil {
		logger.Info("OpenStackCluster is not available yet")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("openStackCluster", openStackCluster.Name)

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(openStackMachine, r)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always patch the openStackMachine when exiting this function so we can persist any AWSMachine changes.
	defer func() {
		if err := patchHelper.Patch(ctx, openStackMachine); err != nil {
			if reterr == nil {
				reterr = err
			}
		}
	}()

	// Handle deleted machines
	if !openStackMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(logger, machine, openStackMachine, cluster, openStackCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, logger, patchHelper, machine, openStackMachine, cluster, openStackCluster)
}

func (r *OpenStackMachineReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.OpenStackMachine{}).
		Watches(
			&source.Kind{Type: &clusterv1.Machine{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("OpenStackMachine")),
			},
		).
		Watches(
			&source.Kind{Type: &infrav1.OpenStackCluster{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: handler.ToRequestsFunc(r.OpenStackClusterToOpenStackMachines)},
		).
		WithEventFilter(pausePredicates).
		Build(r)

	if err != nil {
		return err
	}

	return controller.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(r.requeueOpenStackMachinesForUnpausedCluster),
		},
		predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				oldCluster := e.ObjectOld.(*clusterv1.Cluster)
				newCluster := e.ObjectNew.(*clusterv1.Cluster)
				return oldCluster.Spec.Paused && !newCluster.Spec.Paused
			},
			CreateFunc: func(e event.CreateEvent) bool {
				cluster := e.Object.(*clusterv1.Cluster)
				return !cluster.Spec.Paused
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return false
			},
		},
	)
}

func (r *OpenStackMachineReconciler) reconcileDelete(logger logr.Logger, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (ctrl.Result, error) {
	logger.Info("Handling deleted OpenStackMachine")

	clusterName := fmt.Sprintf("%s-%s", cluster.ObjectMeta.Namespace, cluster.Name)

	osProviderClient, clientOpts, err := provider.NewClientFromMachine(r.Client, openStackMachine)
	if err != nil {
		return ctrl.Result{}, err
	}

	computeService, err := compute.NewService(osProviderClient, clientOpts, logger)
	if err != nil {
		return ctrl.Result{}, err
	}

	loadBalancerService, err := loadbalancer.NewService(osProviderClient, clientOpts, logger, openStackCluster.Spec.UseOctavia)
	if err != nil {
		return ctrl.Result{}, err
	}
	if openStackCluster.Spec.ManagedAPIServerLoadBalancer {
		err = loadBalancerService.DeleteLoadBalancerMember(clusterName, machine, openStackMachine, openStackCluster)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	instance, err := computeService.InstanceExists(openStackMachine)
	if err != nil {
		return ctrl.Result{}, err
	}

	if instance == nil {
		logger.Info("Skipped deleting machine that is already deleted")
		controllerutil.RemoveFinalizer(openStackMachine, infrav1.MachineFinalizer)
		return ctrl.Result{}, nil
	}

	// TODO(sbueringer) wait for instance deleted
	err = computeService.InstanceDelete(machine)
	if err != nil {
		handleUpdateMachineError(logger, openStackMachine, errors.Errorf("error deleting Openstack instance: %v", err))
		return ctrl.Result{}, nil
	}

	logger.Info("OpenStack machine deleted successfully")
	r.Recorder.Eventf(openStackMachine, corev1.EventTypeNormal, "SuccessfulTerminate", "Terminated instance %q", instance.ID)

	// Instance is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(openStackMachine, infrav1.MachineFinalizer)

	return ctrl.Result{}, nil
}

func (r *OpenStackMachineReconciler) reconcileNormal(ctx context.Context, logger logr.Logger, patchHelper *patch.Helper, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (_ ctrl.Result, reterr error) {
	// If the OpenStackMachine is in an error state, return early.
	if openStackMachine.Status.FailureReason != nil || openStackMachine.Status.FailureMessage != nil {
		logger.Info("Error state detected, skipping reconciliation")
		return ctrl.Result{}, nil
	}

	// If the OpenStackMachine doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(openStackMachine, infrav1.ClusterFinalizer)
	// Register the finalizer immediately to avoid orphaning OpenStack resources on delete
	if err := patchHelper.Patch(ctx, openStackMachine); err != nil {
		return ctrl.Result{}, err
	}

	if !cluster.Status.InfrastructureReady {
		logger.Info("Cluster infrastructure is not ready yet, requeuing machine")
		return ctrl.Result{RequeueAfter: waitForClusterInfrastructureReadyDuration}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machine.Spec.Bootstrap.DataSecretName == nil {
		logger.Info("Waiting for bootstrap data to be available")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}
	userData, err := r.getBootstrapData(machine, openStackMachine)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Creating Machine")

	clusterName := fmt.Sprintf("%s-%s", cluster.ObjectMeta.Namespace, cluster.Name)

	osProviderClient, clientOpts, err := provider.NewClientFromMachine(r.Client, openStackMachine)
	if err != nil {
		return ctrl.Result{}, err
	}

	computeService, err := compute.NewService(osProviderClient, clientOpts, logger)
	if err != nil {
		return ctrl.Result{}, err
	}

	networkingService, err := networking.NewService(osProviderClient, clientOpts, logger)
	if err != nil {
		return ctrl.Result{}, err
	}

	instance, err := r.getOrCreate(computeService, machine, openStackMachine, cluster, openStackCluster, userData)
	if err != nil {
		handleUpdateMachineError(logger, openStackMachine, errors.Errorf("OpenStack instance cannot be created: %v", err))
		return ctrl.Result{}, err
	}

	// Set an error message if we couldn't find the instance.
	if instance == nil {
		handleUpdateMachineError(logger, openStackMachine, errors.New("OpenStack instance cannot be found"))
		return ctrl.Result{}, nil
	}

	// TODO(sbueringer) From CAPA: TODO(ncdc): move this validation logic into a validating webhook (for us: create validation logic in webhook)

	openStackMachine.Spec.ProviderID = pointer.StringPtr(fmt.Sprintf("openstack://%s", instance.ID))

	openStackMachine.Status.InstanceState = &instance.State

	// TODO(sbueringer) From CAPA: TODO(vincepri): Remove this annotation when clusterctl is no longer relevant.
	if openStackMachine.Annotations == nil {
		openStackMachine.Annotations = map[string]string{}
	}
	openStackMachine.Annotations["cluster-api-provider-openstack"] = "true"

	switch instance.State {
	case infrav1.InstanceStateActive:
		logger.Info("Machine instance is ACTIVE", "instance-id", instance.ID)
		openStackMachine.Status.Ready = true
	case infrav1.InstanceStateBuilding:
		logger.Info("Machine instance is BUILDING", "instance-id", instance.ID)
	default:
		handleUpdateMachineError(logger, openStackMachine, errors.Errorf("OpenStack instance state %q is unexpected", instance.State))
		return ctrl.Result{}, nil
	}

	if openStackMachine.Spec.FloatingIP != "" {
		err = r.reconcileFloatingIP(computeService, networkingService, instance, openStackMachine, openStackCluster)
		if err != nil {
			handleUpdateMachineError(logger, openStackMachine, errors.Errorf("FloatingIP cannot be reconciled: %v", err))
			return ctrl.Result{}, nil
		}
	}

	if openStackCluster.Spec.ManagedAPIServerLoadBalancer {
		err = r.reconcileLoadBalancerMember(logger, osProviderClient, clientOpts, instance, clusterName, machine, openStackMachine, openStackCluster)
		if err != nil {
			handleUpdateMachineError(logger, openStackMachine, errors.Errorf("LoadBalancerMember cannot be reconciled: %v", err))
			return ctrl.Result{}, nil
		}
	}

	logger.Info("Reconciled Machine create successfully")
	return ctrl.Result{}, nil
}

func (r *OpenStackMachineReconciler) getOrCreate(computeService *compute.Service, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster, userData string) (*compute.Instance, error) {

	instance, err := computeService.InstanceExists(openStackMachine)
	if err != nil {
		return nil, err
	}

	if instance == nil {
		instance, err = computeService.InstanceCreate(cluster.Name, machine, openStackMachine, openStackCluster, userData)
		if err != nil {
			return nil, errors.Errorf("error creating Openstack instance: %v", err)
		}
		instanceCreateTimeout := getTimeout("CLUSTER_API_OPENSTACK_INSTANCE_CREATE_TIMEOUT", TimeoutInstanceCreate)
		instanceCreateTimeout *= time.Minute
		// instance in PollImmediate has to overwrites instance of the outer scope to get an updated instance state,
		// which is then returned at the end of getOrCreate
		err = util.PollImmediate(RetryIntervalInstanceStatus, instanceCreateTimeout, func() (bool, error) {
			instance, err = computeService.GetInstance(instance.ID)
			if err != nil {
				return false, nil
			}
			return instance.Status == "ACTIVE", nil
		})
		if err != nil {
			return nil, errors.Errorf("error creating Openstack instance: %v", err)
		}
	}

	return instance, nil
}

func handleUpdateMachineError(logger logr.Logger, openstackMachine *infrav1.OpenStackMachine, message error) {
	err := capierrors.UpdateMachineError
	openstackMachine.Status.FailureReason = &err
	openstackMachine.Status.FailureMessage = pointer.StringPtr(message.Error())
	// TODO remove if this error is logged redundantly
	logger.Error(fmt.Errorf(string(err)), message.Error())
}

func getTimeout(name string, timeout int) time.Duration {
	if v := os.Getenv(name); v != "" {
		timeout, err := strconv.Atoi(v)
		if err == nil {
			return time.Duration(timeout)
		}
	}
	return time.Duration(timeout)
}

func (r *OpenStackMachineReconciler) reconcileFloatingIP(computeService *compute.Service, networkingService *networking.Service, instance *compute.Instance, openStackMachine *infrav1.OpenStackMachine, openStackCluster *infrav1.OpenStackCluster) error {
	err := networkingService.CreateFloatingIPIfNecessary(openStackCluster, openStackMachine.Spec.FloatingIP)
	if err != nil {
		return fmt.Errorf("error creating floatingIP: %v", err)
	}

	err = computeService.AssociateFloatingIP(instance.ID, openStackMachine.Spec.FloatingIP)
	if err != nil {
		return fmt.Errorf("error associating floatingIP: %v", err)
	}
	return nil
}

func (r *OpenStackMachineReconciler) reconcileLoadBalancerMember(logger logr.Logger, osProviderClient *gophercloud.ProviderClient, clientOpts *clientconfig.ClientOpts, instance *compute.Instance, clusterName string, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, openStackCluster *infrav1.OpenStackCluster) error {
	ip, err := getIPFromInstance(instance)
	if err != nil {
		return err
	}
	loadbalancerService, err := loadbalancer.NewService(osProviderClient, clientOpts, logger, openStackCluster.Spec.UseOctavia)
	if err != nil {
		return err
	}

	if err := loadbalancerService.ReconcileLoadBalancerMember(clusterName, machine, openStackMachine, openStackCluster, ip); err != nil {
		return err
	}
	return nil
}

func getIPFromInstance(instance *compute.Instance) (string, error) {
	if instance.AccessIPv4 != "" && net.ParseIP(instance.AccessIPv4) != nil {
		return instance.AccessIPv4, nil
	}
	type networkInterface struct {
		Address string  `json:"addr"`
		Version float64 `json:"version"`
		Type    string  `json:"OS-EXT-IPS:type"`
	}
	var addrList []string

	for _, b := range instance.Addresses {
		list, err := json.Marshal(b)
		if err != nil {
			return "", fmt.Errorf("extract IP from instance err: %v", err)
		}
		var networks []interface{}
		err = json.Unmarshal(list, &networks)
		if err != nil {
			return "", fmt.Errorf("extract IP from instance err: %v", err)
		}
		for _, network := range networks {
			var netInterface networkInterface
			b, _ := json.Marshal(network)
			err = json.Unmarshal(b, &netInterface)
			if err != nil {
				return "", fmt.Errorf("extract IP from instance err: %v", err)
			}
			if netInterface.Version == 4.0 {
				if netInterface.Type == "floating" {
					return netInterface.Address, nil
				}
				addrList = append(addrList, netInterface.Address)
			}
		}
	}
	if len(addrList) != 0 {
		return addrList[0], nil
	}
	return "", fmt.Errorf("extract IP from instance err")
}

// OpenStackClusterToOpenStackMachine is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// of OpenStackMachines.
func (r *OpenStackMachineReconciler) OpenStackClusterToOpenStackMachines(o handler.MapObject) []ctrl.Request {
	result := []ctrl.Request{}

	c, ok := o.Object.(*infrav1.OpenStackCluster)
	if !ok {
		r.Log.Error(errors.Errorf("expected a OpenStackCluster but got a %T", o.Object), "failed to get OpenStackMachine for OpenStackCluster")
		return nil
	}
	log := r.Log.WithValues("OpenStackCluster", c.Name, "Namespace", c.Namespace)

	cluster, err := util.GetOwnerCluster(context.TODO(), r.Client, c.ObjectMeta)
	switch {
	case apierrors.IsNotFound(err) || cluster == nil:
		return result
	case err != nil:
		log.Error(err, "failed to get owning cluster")
		return result
	}

	labels := map[string]string{clusterv1.ClusterLabelName: cluster.Name}
	machineList := &clusterv1.MachineList{}
	if err := r.List(context.TODO(), machineList, client.InNamespace(c.Namespace), client.MatchingLabels(labels)); err != nil {
		log.Error(err, "failed to list Machines")
		return nil
	}
	for _, m := range machineList.Items {
		if m.Spec.InfrastructureRef.Name == "" {
			continue
		}
		name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}

	return result
}

func (r *OpenStackMachineReconciler) getBootstrapData(machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine) (string, error) {
	if machine.Spec.Bootstrap.DataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: machine.Namespace, Name: *machine.Spec.Bootstrap.DataSecretName}
	if err := r.Client.Get(context.TODO(), key, secret); err != nil {
		return "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for Openstack Machine %s/%s", machine.Namespace, openStackMachine.Name)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return base64.StdEncoding.EncodeToString(value), nil
}

func (r *OpenStackMachineReconciler) requeueOpenStackMachinesForUnpausedCluster(o handler.MapObject) []ctrl.Request {
	c, ok := o.Object.(*clusterv1.Cluster)
	if !ok {
		r.Log.Error(errors.Errorf("expected a Cluster but got a %T", o.Object), "failed to get OpenStackMachines for unpaused Cluster")
		return nil
	}

	// Don't handle deleted clusters
	if !c.ObjectMeta.DeletionTimestamp.IsZero() {
		return nil
	}

	return r.requestsForCluster(c.Namespace, c.Name)
}

func (r *OpenStackMachineReconciler) requestsForCluster(namespace, name string) []ctrl.Request {
	log := r.Log.WithValues("Cluster", name, "Namespace", namespace)
	labels := map[string]string{clusterv1.ClusterLabelName: name}
	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(context.TODO(), machineList, client.InNamespace(namespace), client.MatchingLabels(labels)); err != nil {
		log.Error(err, "failed to get owned Machines")
		return nil
	}

	result := make([]ctrl.Request, 0, len(machineList.Items))
	for _, m := range machineList.Items {
		if m.Spec.InfrastructureRef.Name != "" {
			result = append(result, ctrl.Request{NamespacedName: client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}})
		}
	}
	return result
}
