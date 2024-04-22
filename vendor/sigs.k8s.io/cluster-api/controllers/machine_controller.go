/*
Copyright 2019 The Kubernetes Authors.

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
	"sync"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/cluster-api/controllers/remote"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

var (
	errNilNodeRef           = errors.New("noderef is nil")
	errLastControlPlaneNode = errors.New("last control plane member")
	errNoControlPlaneNodes  = errors.New("no control plane members")
)

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io;bootstrap.cluster.x-k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch;create;update;patch;delete

// MachineReconciler reconciles a Machine object
type MachineReconciler struct {
	client.Client
	Log logr.Logger

	controller       controller.Controller
	recorder         record.EventRecorder
	externalWatchers sync.Map
}

func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.Machine{}).
		WithOptions(options).
		Build(r)

	r.controller = c
	r.recorder = mgr.GetEventRecorderFor("machine-controller")
	return err
}

func (r *MachineReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()
	_ = r.Log.WithValues("machine", req.NamespacedName)

	// Fetch the Machine instance
	m := &clusterv1.Machine{}
	if err := r.Client.Get(ctx, req.NamespacedName, m); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(m, r)
	if err != nil {
		return ctrl.Result{}, err
	}

	defer func() {
		// Always reconcile the Status.Phase field.
		r.reconcilePhase(ctx, m)

		// Always attempt to Patch the Machine object and status after each reconciliation.
		if err := patchHelper.Patch(ctx, m); err != nil {
			if reterr == nil {
				reterr = err
			}
		}
	}()

	// Cluster might be nil as some providers might not require a cluster object
	// for machine management.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, m.ObjectMeta)
	if errors.Cause(err) == util.ErrNoCluster {
		klog.V(2).Infof("Machine %q in namespace %q doesn't specify %q label, assuming nil cluster",
			m.Name, m.Namespace, clusterv1.MachineClusterLabelName)
	} else if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to get cluster %q for machine %q in namespace %q",
			m.Labels[clusterv1.MachineClusterLabelName], m.Name, m.Namespace)
	}

	// Handle deletion reconciliation loop.
	if !m.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, cluster, m)
	}

	// Handle normal reconciliation loop.
	return r.reconcile(ctx, cluster, m)
}

func (r *MachineReconciler) reconcile(ctx context.Context, cluster *clusterv1.Cluster, m *clusterv1.Machine) (ctrl.Result, error) {
	// If the Machine belongs to a cluster, add an owner reference.
	if cluster != nil && r.shouldAdopt(m) {
		m.OwnerReferences = util.EnsureOwnerRef(m.OwnerReferences, metav1.OwnerReference{
			APIVersion: cluster.APIVersion,
			Kind:       cluster.Kind,
			Name:       cluster.Name,
			UID:        cluster.UID,
		})
	}

	// If the Machine doesn't have a finalizer, add one.
	if !util.Contains(m.Finalizers, clusterv1.MachineFinalizer) {
		m.Finalizers = append(m.ObjectMeta.Finalizers, clusterv1.MachineFinalizer)
	}

	// Call the inner reconciliation methods.
	reconciliationErrors := []error{
		r.reconcileBootstrap(ctx, m),
		r.reconcileInfrastructure(ctx, m),
		r.reconcileNodeRef(ctx, cluster, m),
		r.reconcileClusterAnnotations(ctx, cluster, m),
	}

	// Parse the errors, making sure we record if there is a RequeueAfterError.
	res := ctrl.Result{}
	errs := []error{}
	for _, err := range reconciliationErrors {
		if requeueErr, ok := errors.Cause(err).(capierrors.HasRequeueAfterError); ok {
			// Only record and log the first RequeueAfterError.
			if !res.Requeue {
				res.Requeue = true
				res.RequeueAfter = requeueErr.GetRequeueAfter()
				klog.Infof("Reconciliation for Machine %q in namespace %q asked to requeue: %v", m.Name, m.Namespace, err)
			}
			continue
		}

		errs = append(errs, err)
	}
	return res, kerrors.NewAggregate(errs)
}

func (r *MachineReconciler) reconcileDelete(ctx context.Context, cluster *clusterv1.Cluster, m *clusterv1.Machine) (ctrl.Result, error) {
	if err := r.isDeleteNodeAllowed(ctx, m); err != nil {
		switch err {
		case errNilNodeRef:
			klog.V(2).Infof("Deleting node is not allowed for machine %q: %v", m.Name, err)
		case errNoControlPlaneNodes, errLastControlPlaneNode:
			klog.V(2).Infof("Deleting node %q is not allowed for machine %q: %v", m.Status.NodeRef.Name, m.Name, err)
		default:
			klog.Errorf("IsDeleteNodeAllowed check failed for machine %q: %v", m.Name, err)
			return ctrl.Result{}, err
		}
	} else {
		klog.Infof("Deleting node %q for machine %q", m.Status.NodeRef.Name, m.Name)
		if err := r.deleteNode(ctx, cluster, m.Status.NodeRef.Name); err != nil && !apierrors.IsNotFound(err) {
			klog.Errorf("Error deleting node %q for machine %q: %v", m.Status.NodeRef.Name, m.Name, err)
			return ctrl.Result{}, err
		}
	}

	if ok, err := r.reconcileDeleteExternal(ctx, m); !ok || err != nil {
		// Return early and don't remove the finalizer if we got an error or
		// the external reconciliation deletion isn't ready.
		return ctrl.Result{}, err
	}

	m.ObjectMeta.Finalizers = util.Filter(m.ObjectMeta.Finalizers, clusterv1.MachineFinalizer)
	return ctrl.Result{}, nil
}

// isDeleteNodeAllowed returns nil only if the Machine's NodeRef is not nil
// and if the Machine is not the last control plane node in the cluster.
func (r *MachineReconciler) isDeleteNodeAllowed(ctx context.Context, machine *clusterv1.Machine) error {
	// Cannot delete something that doesn't exist.
	if machine.Status.NodeRef == nil {
		return errNilNodeRef
	}

	// Get all of the machines that belong to this cluster.
	machines, err := r.getMachinesInCluster(ctx, machine.Namespace, machine.Labels[clusterv1.MachineClusterLabelName])
	if err != nil {
		return err
	}

	// Whether or not it is okay to delete the NodeRef depends on the
	// number of remaining control plane members and whether or not this
	// machine is one of them.
	switch numControlPlaneMachines := len(util.GetControlPlaneMachines(machines)); {
	case numControlPlaneMachines == 0:
		// Do not delete the NodeRef if there are no remaining members of
		// the control plane.
		return errNoControlPlaneNodes
	case numControlPlaneMachines == 1 && util.IsControlPlaneMachine(machine):
		// Do not delete the NodeRef if this is the last member of the
		// control plane.
		return errLastControlPlaneNode
	default:
		// Otherwise it is okay to delete the NodeRef.
		return nil
	}
}

func (r *MachineReconciler) deleteNode(ctx context.Context, cluster *clusterv1.Cluster, name string) error {
	if cluster == nil {
		// Try to retrieve the Node from the local cluster, if no Cluster reference is found.
		var node corev1.Node
		if err := r.Client.Get(ctx, client.ObjectKey{Name: name}, &node); err != nil {
			return err
		}
		return r.Client.Delete(ctx, &node)
	}

	// Otherwise, proceed to get the remote cluster client and get the Node.
	remoteClient, err := remote.NewClusterClient(r.Client, cluster)
	if err != nil {
		klog.Errorf("Error creating a remote client for cluster %q while deleting Machine %q, won't retry: %v",
			cluster.Name, name, err)
		return nil
	}

	corev1Remote, err := remoteClient.CoreV1()
	if err != nil {
		klog.Errorf("Error creating a remote client for cluster %q while deleting Machine %q, won't retry: %v",
			cluster.Name, name, err)
		return nil
	}

	return corev1Remote.Nodes().Delete(name, &metav1.DeleteOptions{})
}

// getMachinesInCluster returns all of the Machine objects that belong to the
// same cluster as the provided Machine
func (r *MachineReconciler) getMachinesInCluster(ctx context.Context, namespace, name string) ([]*clusterv1.Machine, error) {
	if name == "" {
		return nil, nil
	}

	machineList := &clusterv1.MachineList{}
	labels := map[string]string{clusterv1.MachineClusterLabelName: name}

	if err := r.Client.List(ctx, machineList, client.InNamespace(namespace), client.MatchingLabels(labels)); err != nil {
		return nil, errors.Wrap(err, "failed to list machines")
	}

	machines := make([]*clusterv1.Machine, len(machineList.Items))
	for i := range machineList.Items {
		machines[i] = &machineList.Items[i]
	}

	return machines, nil
}

// reconcileDeleteExternal tries to delete external references, returning true if it cannot find any.
func (r *MachineReconciler) reconcileDeleteExternal(ctx context.Context, m *clusterv1.Machine) (bool, error) {
	objects := []*unstructured.Unstructured{}
	references := []*corev1.ObjectReference{
		m.Spec.Bootstrap.ConfigRef,
		&m.Spec.InfrastructureRef,
	}

	// Loop over the references and try to retrieve it with the client.
	for _, ref := range references {
		if ref == nil {
			continue
		}

		obj, err := external.Get(r.Client, ref, m.Namespace)
		if err != nil && !apierrors.IsNotFound(err) {
			return false, errors.Wrapf(err, "failed to get %s %q for Machine %q in namespace %q",
				obj.GroupVersionKind(), obj.GetName(), m.Name, m.Namespace)
		}
		if obj != nil {
			objects = append(objects, obj)
		}
	}

	// Issue a delete request for any object that has been found.
	for _, obj := range objects {
		if err := r.Delete(ctx, obj); err != nil && !apierrors.IsNotFound(err) {
			return false, errors.Wrapf(err,
				"failed to delete %v %q for Machine %q in namespace %q",
				obj.GroupVersionKind(), obj.GetName(), m.Name, m.Namespace)
		}
	}

	// Return true if there are no more external objects.
	return len(objects) == 0, nil
}

func (r *MachineReconciler) shouldAdopt(m *clusterv1.Machine) bool {
	return !util.HasOwner(m.OwnerReferences, clusterv1.GroupVersion.String(), []string{"MachineSet", "Cluster"})
}
