package infraclustercontroller

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	openshiftconfig "github.com/openshift/api/config/v1"
	mapi "github.com/openshift/api/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
)

const (
	CAPINamespace            = "openshift-cluster-api"
	MAPINamespace            = "openshift-machine-api"
	CloudName                = "openstack" // The name of the clouds.yaml entry to use. It is not configurable in OpenShift.
	CredentialsSecretName    = "openstack-cloud-credentials"
	ControlPlaneEndpointPort = 6443
	OpenStackProviderPrefix  = "openstack:///"
	ClusterOperatorName      = "cluster-api"
	InfrastructureName       = "cluster"

	InfraClusterDegradedCondition = "InfraClusterDegraded"
	InfraClusterReasonUnsupported = "Unsupported"
	InfraClusterReasonAsExpected  = "AsExpected"

	FieldManager = "openstack-infracluster-controller"
)

type OpenShiftClusterReconciler struct {
	Client         client.Client
	Recorder       record.EventRecorder
	ScopeFactory   scope.Factory
	CaCertificates []byte // PEM encoded CA certificates

	Scheme *runtime.Scheme
}

// Cluster-scoped RBAC
// +kubebuilder:rbac:groups=config.openshift.io,resources=clusteroperators,resourceNames=cluster-api,verbs=get;list;watch
// +kubebuilder:rbac:groups=config.openshift.io,resources=clusteroperators/status,resourceNames=cluster-api,verbs=get;list;watch;patch;update
// +kubebuilder:rbac:groups=config.openshift.io,resources=infrastructures;infrastructures/status,verbs=get;list;watch

// Namespace-scoped RBAC
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackclusters,namespace=openshift-cluster-api,verbs=get;list;watch;create
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,namespace=openshift-cluster-api,verbs=get;list;watch;create
// +kubebuilder:rbac:groups="",resources=secrets,namespace=openshift-cluster-api,resourceNames=openstack-cloud-credentials,verbs=get;list;watch
// +kubebuilder:rbac:groups="machine.openshift.io",namespace=openshift-machine-api,resources=machines,verbs=get;list;watch

func (r *OpenShiftClusterReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		Owns(&infrav1.OpenStackCluster{}).
		Owns(&clusterv1.Cluster{}).
		Watches(&openshiftconfig.Infrastructure{}, handler.EnqueueRequestsFromMapFunc(func(context.Context, client.Object) []reconcile.Request {
			return []reconcile.Request{{NamespacedName: types.NamespacedName{Name: ClusterOperatorName}}}
		})).
		For(&openshiftconfig.ClusterOperator{}, builder.WithPredicates(clusterOperatorPredicates())).
		Complete(r)
}

func clusterOperatorPredicates() predicate.Funcs {
	isClusterCAPIOperator := func(obj runtime.Object) bool {
		clusterOperator, ok := obj.(*openshiftconfig.ClusterOperator)
		return ok && clusterOperator.GetName() == ClusterOperatorName
	}

	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return isClusterCAPIOperator(e.Object) },
		UpdateFunc:  func(e event.UpdateEvent) bool { return isClusterCAPIOperator(e.ObjectNew) },
		DeleteFunc:  func(e event.DeleteEvent) bool { return isClusterCAPIOperator(e.Object) },
		GenericFunc: func(e event.GenericEvent) bool { return isClusterCAPIOperator(e.Object) },
	}
}

func (r *OpenShiftClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	co := openshiftconfig.ClusterOperator{}
	if err := r.Client.Get(ctx, req.NamespacedName, &co); err != nil {
		return ctrl.Result{}, err
	}

	infrastructure := openshiftconfig.Infrastructure{}
	if err := r.Client.Get(ctx, types.NamespacedName{Name: InfrastructureName}, &infrastructure); err != nil {
		return reconcile.Result{}, err
	}

	if infrastructure.Status.PlatformStatus.Type != openshiftconfig.OpenStackPlatformType {
		log.V(2).Info("Infrastructure platform type is not OpenStack", "type", infrastructure.Spec.PlatformSpec.Type)
		return ctrl.Result{}, nil
	}

	infraName := infrastructure.Status.InfrastructureName
	log = log.WithValues("infraName", infraName)

	platformStatus := infrastructure.Status.PlatformStatus.OpenStack
	if platformStatus == nil {
		log.V(4).Info("Waiting for OpenStack platform status to be available")
		return ctrl.Result{}, nil
	}

	openStackCluster, err := r.ensureInfraCluster(ctx, log, &infrastructure, &co)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err := r.ensureCAPICluster(ctx, log, &infrastructure, openStackCluster, &co); err != nil {
		return ctrl.Result{}, err
	}

	// Ensure we are not degraded
	if err := r.setClusterOperatorCondition(ctx, &co, openshiftconfig.ConditionFalse, InfraClusterReasonAsExpected, "", log); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *OpenShiftClusterReconciler) setClusterOperatorCondition(ctx context.Context, co *openshiftconfig.ClusterOperator, status openshiftconfig.ConditionStatus, reason, message string, log logr.Logger) error {
	// Check the condition is not already set as requested
	condition := func() *openshiftconfig.ClusterOperatorStatusCondition {
		for i := range co.Status.Conditions {
			if co.Status.Conditions[i].Type == InfraClusterDegradedCondition {
				return &co.Status.Conditions[i]
			}
		}
		return nil
	}()
	if condition != nil && condition.Status == status && condition.Reason == reason && condition.Message == message {
		log.V(4).Info("ClusterOperator status condition already set", "status", status, "reason", reason, "message", message)
		return nil
	}

	/*
		XXX: ClusterOperator Conditions don't have SSA markers so we use
		     strategic merge patch instead. See https://github.com/openshift/api/pull/1631

		apply := openshiftconfig.ClusterOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name: co.Name,
			},
			TypeMeta: co.TypeMeta,
		}
		apply.Status.Conditions = []openshiftconfig.ClusterOperatorStatusCondition{
			{
				Type:               InfraClusterReadyCondition,
				Status:             status,
				LastTransitionTime: metav1.Now(),
				Reason:             reason,
				Message:            message,
			},
		}

		if err := r.Client.Status().Patch(ctx, &apply, client.Apply, client.FieldOwner(FieldManager), client.ForceOwnership); err != nil {
			return err
		}
	*/

	coBefore := co.DeepCopy()
	if condition == nil {
		co.Status.Conditions = append(co.Status.Conditions, openshiftconfig.ClusterOperatorStatusCondition{})
		condition = &co.Status.Conditions[len(co.Status.Conditions)-1]
		condition.Type = InfraClusterDegradedCondition
	}
	condition.Status = status
	condition.Reason = reason
	condition.Message = message
	condition.LastTransitionTime = metav1.Now()

	if err := r.Client.Status().Patch(ctx, co, client.MergeFrom(coBefore)); err != nil {
		return fmt.Errorf("patching ClusterOperator status: %w", err)
	}

	log.V(2).Info("Set ClusterOperator status condition", "status", status, "reason", reason, "message", message)
	return nil
}

// ensureInfraCluster creates a new OpenStackCluster if none exists. It will not modify an existing cluster.
func (r *OpenShiftClusterReconciler) ensureInfraCluster(ctx context.Context, log logr.Logger, infrastructure *openshiftconfig.Infrastructure, co *openshiftconfig.ClusterOperator) (*infrav1.OpenStackCluster, error) {
	infraName := infrastructure.Status.InfrastructureName

	openStackCluster := infrav1.OpenStackCluster{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: infraName, Namespace: CAPINamespace}, &openStackCluster)
	if err == nil {
		// We don't update a cluster that already exists
		log.V(4).Info("OpenStackCluster already exists")
		return &openStackCluster, nil
	}

	if !apierrors.IsNotFound(err) {
		return nil, fmt.Errorf("fetching OpenStackCluster %s: %w", infraName, err)
	}

	// We only get here if the cluster does not exist and we need to create it.
	platformStatus := infrastructure.Status.PlatformStatus.OpenStack

	if platformStatus.LoadBalancer.Type != openshiftconfig.LoadBalancerTypeOpenShiftManagedDefault {
		status := openshiftconfig.ConditionTrue
		reason := InfraClusterReasonUnsupported
		message := fmt.Sprintf("Initialialising infrastructure cluster is not supported for load balancer type %s is not supported", platformStatus.LoadBalancer.Type)

		return nil, r.setClusterOperatorCondition(ctx, co, status, reason, message, log)
	}

	// Set an owner reference to the clusteroperator so we will reconcile on changes
	if err := controllerutil.SetOwnerReference(co, &openStackCluster, r.Client.Scheme()); err != nil {
		return nil, fmt.Errorf("setting owner reference: %w", err)
	}

	/*
	 * Set properties of the OpenStack cluster which are either hard-coded
	 * or can be determined from the infrastructure status.
	 */

	openStackCluster.Name = infraName
	openStackCluster.Namespace = CAPINamespace
	openStackCluster.Labels = map[string]string{
		clusterv1.ClusterNameLabel: infraName,
	}
	openStackCluster.Spec.IdentityRef = infrav1.OpenStackIdentityReference{
		Name:      CredentialsSecretName,
		CloudName: CloudName,
	}

	if len(platformStatus.APIServerInternalIPs) == 0 {
		return nil, fmt.Errorf("no APIServerInternalIPs available")
	}
	openStackCluster.Spec.ControlPlaneEndpoint = &clusterv1.APIEndpoint{
		Host: platformStatus.APIServerInternalIPs[0],
		Port: ControlPlaneEndpointPort,
	}

	openStackCluster.Spec.DisableAPIServerFloatingIP = ptr.To(true)
	openStackCluster.Spec.Tags = []string{
		"openshiftClusterID=" + infraName,
	}

	/*
	 * Set network properties of the OpenStack cluster. We don't have any
	 * reliable record of these, so we have to determine them by examining
	 * the existing control plane machines.
	 */

	scope, err := r.ScopeFactory.NewClientScopeFromObject(ctx, r.Client, r.CaCertificates, log, &openStackCluster)
	if err != nil {
		return nil, fmt.Errorf("creating OpenStack provider scope: %w", err)
	}

	networkClient, err := scope.NewNetworkClient()
	if err != nil {
		return nil, fmt.Errorf("creating OpenStack compute service: %w", err)
	}

	defaultSubnet, err := r.getDefaultSubnetFromMachines(ctx, log, networkClient, platformStatus)
	if err != nil {
		return nil, err
	}
	if defaultSubnet == nil {
		return nil, fmt.Errorf("unable to determine default subnet from control plane machines")
	}
	openStackCluster.Spec.Network = &infrav1.NetworkParam{ID: ptr.To(defaultSubnet.NetworkID)}
	// N.B. Deliberately don't add subnet here: CAPO will use all subnets in network, which should also cover dual stack deployments

	routerID, err := r.getDefaultRouterIDFromSubnet(ctx, networkClient, defaultSubnet)
	if err != nil {
		return nil, err
	}
	openStackCluster.Spec.Router = &infrav1.RouterParam{ID: ptr.To(routerID)}

	router, err := networkClient.GetRouter(routerID)
	if err != nil {
		return nil, fmt.Errorf("getting router %s: %w", routerID, err)
	}
	if router.GatewayInfo.NetworkID == "" {
		return nil, fmt.Errorf("router %s does not have an external gateway", routerID)
	}
	// N.B. The only reason we set ExternalNetworkID in the cluster spec is
	// to avoid an error reconciling the external network if it isn't set.
	// If CAPO ever no longer requires this we can just not set it and
	// remove much of the code above. We don't actually use it.
	openStackCluster.Spec.ExternalNetwork = &infrav1.NetworkParam{ID: ptr.To(router.GatewayInfo.NetworkID)}

	err = r.Client.Create(ctx, &openStackCluster)
	if err != nil {
		return nil, fmt.Errorf("creating OpenStackCluster: %w", err)
	}

	log.V(2).Info("Created OpenStackCluster")
	return &openStackCluster, nil
}

func (r *OpenShiftClusterReconciler) getDefaultRouterIDFromSubnet(_ context.Context, networkClient clients.NetworkClient, subnet *subnets.Subnet) (string, error) {
	// Find the port which owns the subnet's gateway IP
	ports, err := networkClient.ListPort(ports.ListOpts{
		NetworkID: subnet.NetworkID,
		FixedIPs: []ports.FixedIPOpts{
			{
				IPAddress: subnet.GatewayIP,
				// XXX: We should search on both subnet and IP
				// address here, but can't because of
				// https://github.com/gophercloud/gophercloud/issues/2807
				// SubnetID:  subnet.ID,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("listing ports: %w", err)
	}

	if len(ports) == 0 {
		return "", fmt.Errorf("no ports found for subnet %s", subnet.ID)
	}

	if len(ports) > 1 {
		return "", fmt.Errorf("multiple ports found for subnet %s", subnet.ID)
	}

	return ports[0].DeviceID, nil
}

// getDefaultSubnetFromMachines attempts to infer the default cluster subnet by
// directly examining the control plane machines. Specifically it looks for a
// subnet attached to a control plane machine whose CIDR contains the API
// loadbalancer internal VIP.
//
// This heuristic is only valid when the API loadbalancer type is
// LoadBalancerTypeOpenShiftManagedDefault.
func (r *OpenShiftClusterReconciler) getDefaultSubnetFromMachines(ctx context.Context, log logr.Logger, networkClient clients.NetworkClient, platformStatus *openshiftconfig.OpenStackPlatformStatus) (*subnets.Subnet, error) {
	mapiMachines := mapi.MachineList{}
	if err := r.Client.List(ctx, &mapiMachines, client.InNamespace(MAPINamespace), client.MatchingLabels{"machine.openshift.io/cluster-api-machine-role": "master"}); err != nil {
		return nil, fmt.Errorf("listing control plane machines: %w", err)
	}

	if len(mapiMachines.Items) == 0 {
		return nil, fmt.Errorf("no control plane machines found")
	}

	apiServerInternalIPs := make([]net.IP, len(platformStatus.APIServerInternalIPs))
	for i, ipStr := range platformStatus.APIServerInternalIPs {
		apiServerInternalIPs[i] = net.ParseIP(ipStr)
	}

	for _, mapiMachine := range mapiMachines.Items {
		log := log.WithValues("machine", mapiMachine.Name)

		providerID := mapiMachine.Spec.ProviderID
		if providerID == nil {
			log.V(3).Info("Skipping machine: providerID is not set")
			continue
		}

		if !strings.HasPrefix(*providerID, OpenStackProviderPrefix) {
			log.V(2).Info("Skipping machine: providerID has unexpected format", "providerID", *providerID)
			continue
		}

		instanceID := (*providerID)[len(OpenStackProviderPrefix):]

		ports, err := listPortsByInstanceID(networkClient, instanceID)
		if err != nil {
			return nil, fmt.Errorf("listing ports for instance %s: %w", instanceID, err)
		}
		if len(ports) == 0 {
			return nil, fmt.Errorf("no ports found for instance %s", instanceID)
		}
		for _, port := range ports {
			log := log.WithValues("port", port.ID)

			for _, fixedIP := range port.FixedIPs {
				if fixedIP.SubnetID == "" {
					continue
				}

				subnet, err := networkClient.GetSubnet(fixedIP.SubnetID)
				if err != nil {
					return nil, fmt.Errorf("getting subnet %s: %w", fixedIP.SubnetID, err)
				}

				_, cidr, err := net.ParseCIDR(subnet.CIDR)
				if err != nil {
					return nil, fmt.Errorf("parsing subnet CIDR %s: %w", subnet.CIDR, err)
				}
				for i := range apiServerInternalIPs {
					if cidr.Contains(apiServerInternalIPs[i]) {
						return subnet, nil
					}
				}

				log.V(6).Info("subnet does not match any APIServerInternalIPs", "subnet", subnet.CIDR)
			}

			log.V(6).Info("port does not match any APIServerInternalIPs")
		}

		log.V(6).Info("machine does not match any APIServerInternalIPs")
	}

	return nil, nil
}

func listPortsByInstanceID(networkClient clients.NetworkClient, instanceID string) ([]ports.Port, error) {
	portOpts := ports.ListOpts{
		DeviceID: instanceID,
	}
	return networkClient.ListPort(portOpts)
}

// ensureCAPICluster creates a new CAPI cluster if none exists. It will not modify an existing cluster.
func (r *OpenShiftClusterReconciler) ensureCAPICluster(ctx context.Context, log logr.Logger, infrastructure *openshiftconfig.Infrastructure, openStackCluster *infrav1.OpenStackCluster, co *openshiftconfig.ClusterOperator) error {
	infraName := infrastructure.Status.InfrastructureName

	cluster := clusterv1.Cluster{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: infraName, Namespace: CAPINamespace}, &cluster)
	if err == nil {
		// We don't update a cluster that already exists
		log.V(4).Info("Cluster already exists")
		return nil
	}

	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("fetching Cluster %s: %w", infraName, err)
	}

	// We only get here if the cluster does not exist and we need to create it.

	gvk, err := apiutil.GVKForObject(openStackCluster, r.Client.Scheme())
	if err != nil {
		return fmt.Errorf("getting GVK for OpenStackCluster: %w", err)
	}

	// Set an owner reference to the clusteroperator so we will reconcile on changes
	if err := controllerutil.SetOwnerReference(co, &cluster, r.Client.Scheme()); err != nil {
		return fmt.Errorf("setting owner reference: %w", err)
	}

	cluster.Name = infraName
	cluster.Namespace = CAPINamespace
	cluster.Spec.InfrastructureRef = &corev1.ObjectReference{
		Kind:       gvk.Kind,
		Name:       openStackCluster.Name,
		Namespace:  openStackCluster.Namespace,
		APIVersion: gvk.GroupVersion().String(),
	}

	if err := r.Client.Create(ctx, &cluster); err != nil {
		return fmt.Errorf("creating Cluster: %w", err)
	}
	log.V(2).Info("Created Cluster")
	return nil
}
