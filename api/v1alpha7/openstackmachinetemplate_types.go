/*
Copyright 2023 The Kubernetes Authors.

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

package v1alpha7

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenStackMachineTemplateSpec defines the desired state of OpenStackMachineTemplate.
type OpenStackMachineTemplateSpec struct {
	Template OpenStackMachineTemplateResource `json:"template"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion:warning="The v1alpha7 version of OpenStackMachineTemplate has been deprecated and will be removed in a future release."
// +kubebuilder:resource:path=openstackmachinetemplates,scope=Namespaced,categories=cluster-api,shortName=osmt

// OpenStackMachineTemplate is the Schema for the openstackmachinetemplates API.
//
// Deprecated: v1alpha7.OpenStackMachineTemplate has been replaced by v1beta1.OpenStackMachineTemplate.
type OpenStackMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OpenStackMachineTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackMachineTemplateList contains a list of OpenStackMachineTemplate.
//
// Deprecated: v1alpha7.OpenStackMachineTemplateList has been replaced by v1beta1.OpenStackMachineTemplateList.
type OpenStackMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackMachineTemplate{}, &OpenStackMachineTemplateList{})
}
