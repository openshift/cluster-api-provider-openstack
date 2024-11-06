package scheme

import (
	openshiftconfig "github.com/openshift/api/config/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

func DefaultScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(openshiftconfig.AddToScheme(scheme))
	utilruntime.Must(machinev1beta1.AddToScheme(scheme))
	utilruntime.Must(infrav1alpha1.AddToScheme(scheme))
	utilruntime.Must(infrav1beta1.AddToScheme(scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme))

	return scheme
}
