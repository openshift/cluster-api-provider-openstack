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

package v1alpha3

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

const (
	// ClusterFinalizer allows ReconcileOpenStackCluster to clean up OpenStack resources associated with OpenStackCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "openstackcluster.infrastructure.cluster.x-k8s.io"
)

// OpenStackClusterSpec defines the desired state of OpenStackCluster
type OpenStackClusterSpec struct {

	// The name of the secret containing the openstack credentials
	// +optional
	CloudsSecret *corev1.SecretReference `json:"cloudsSecret"`

	// The name of the cloud to use from the clouds secret
	// +optional
	CloudName string `json:"cloudName"`

	// NodeCIDR is the OpenStack Subnet to be created. Cluster actuator will create a
	// network, a subnet with NodeCIDR, and a router connected to this subnet.
	// If you leave this empty, no network will be created.
	NodeCIDR string `json:"nodeCidr,omitempty"`

	// If NodeCIDR cannot be set this can be used to detect an existing network.
	Network Filter `json:"network,omitempty"`

	// If NodeCIDR cannot be set this can be used to detect an existing subnet.
	Subnet SubnetFilter `json:"subnet,omitempty"`

	// DNSNameservers is the list of nameservers for OpenStack Subnet being created.
	// Set this value when you need create a new network/subnet while the access
	// through DNS is required.
	DNSNameservers []string `json:"dnsNameservers,omitempty"`
	// ExternalRouterIPs is an array of externalIPs on the respective subnets.
	// This is necessary if the router needs a fixed ip in a specific subnet.
	ExternalRouterIPs []ExternalRouterIPParam `json:"externalRouterIPs,omitempty"`
	// ExternalNetworkID is the ID of an external OpenStack Network. This is necessary
	// to get public internet to the VMs.
	// +optional
	ExternalNetworkID string `json:"externalNetworkId,omitempty"`

	// UseOctavia is weather LoadBalancer Service is Octavia or not
	// +optional
	UseOctavia bool `json:"useOctavia,omitempty"`

	// ManagedAPIServerLoadBalancer defines whether a LoadBalancer for the
	// APIServer should be created. If set to true the following properties are
	// mandatory: APIServerLoadBalancerFloatingIP, APIServerLoadBalancerPort
	// +optional
	ManagedAPIServerLoadBalancer bool `json:"managedAPIServerLoadBalancer"`

	// APIServerLoadBalancerFloatingIP is the floatingIP which will be associated
	// to the APIServer loadbalancer. The floatingIP will be created if it not
	// already exists.
	APIServerLoadBalancerFloatingIP string `json:"apiServerLoadBalancerFloatingIP,omitempty"`

	// APIServerLoadBalancerPort is the port on which the listener on the APIServer
	// loadbalancer will be created
	APIServerLoadBalancerPort int `json:"apiServerLoadBalancerPort,omitempty"`

	// APIServerLoadBalancerAdditionalPorts adds additional ports to the APIServerLoadBalancer
	APIServerLoadBalancerAdditionalPorts []int `json:"apiServerLoadBalancerAdditionalPorts,omitempty"`

	// ManagedSecurityGroups defines that kubernetes manages the OpenStack security groups
	// for now, that means that we'll create security group allows traffic to/from
	// machines belonging to that group based on Calico CNI plugin default network
	// requirements: BGP and IP-in-IP for master node(s) and worker node(s) respectively.
	// In the future, we could make this more flexible.
	// +optional
	ManagedSecurityGroups bool `json:"managedSecurityGroups"`

	// DisablePortSecurity disables the port security of the network created for the
	// Kubernetes cluster, which also disables SecurityGroups
	DisablePortSecurity bool `json:"disablePortSecurity,omitempty"`

	// Tags for all resources in cluster
	Tags []string `json:"tags,omitempty"`

	// CAKeyPair is the key pair for ca certs.
	CAKeyPair KeyPair `json:"caKeyPair,omitempty"`

	//EtcdCAKeyPair is the key pair for etcd.
	EtcdCAKeyPair KeyPair `json:"etcdCAKeyPair,omitempty"`

	// FrontProxyCAKeyPair is the key pair for FrontProxyKeyPair.
	FrontProxyCAKeyPair KeyPair `json:"frontProxyCAKeyPair,omitempty"`

	// SAKeyPair is the service account key pair.
	SAKeyPair KeyPair `json:"saKeyPair,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// ControlPlaneAvailabilityZones is the az to deploy control plane to
	ControlPlaneAvailabilityZones []string `json:"controlPlaneAvailabilityZones,omitempty"`
}

// OpenStackClusterStatus defines the observed state of OpenStackCluster
type OpenStackClusterStatus struct {
	Ready bool `json:"ready"`

	// Network contains all information about the created OpenStack Network.
	// It includes Subnets and Router.
	Network *Network `json:"network,omitempty"`

	// External Network contains information about the created OpenStack external network.
	ExternalNetwork *Network `json:"externalNetwork,omitempty"`

	// FailureDomains represent OpenStack availability zones
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	// ControlPlaneSecurityGroups contains all the information about the OpenStack
	// Security Group that needs to be applied to control plane nodes.
	// TODO: Maybe instead of two properties, we add a property to the group?
	ControlPlaneSecurityGroup *SecurityGroup `json:"controlPlaneSecurityGroup,omitempty"`

	// WorkerSecurityGroup contains all the information about the OpenStack Security
	// Group that needs to be applied to worker nodes.
	WorkerSecurityGroup *SecurityGroup `json:"workerSecurityGroup,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=openstackclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this OpenStackCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for OpenStack instances"
// +kubebuilder:printcolumn:name="Network",type="string",JSONPath=".status.network.id",description="Network the cluster is using"
// +kubebuilder:printcolumn:name="Subnet",type="string",JSONPath=".status.network.subnet.id",description="Subnet the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.network.apiServerLoadBalancer.ip",description="API Endpoint",priority=1

// OpenStackCluster is the Schema for the openstackclusters API
type OpenStackCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackClusterSpec   `json:"spec,omitempty"`
	Status OpenStackClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackClusterList contains a list of OpenStackCluster
type OpenStackClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackCluster{}, &OpenStackClusterList{})
}
