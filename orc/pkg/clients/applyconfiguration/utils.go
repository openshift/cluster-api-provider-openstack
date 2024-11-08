/*
Copyright 2024 The ORC Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"
	apiv1alpha1 "github.com/k-orc/openstack-resource-controller/pkg/clients/applyconfiguration/api/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=openstack.k-orc.cloud, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("CloudCredentialsReference"):
		return &apiv1alpha1.CloudCredentialsReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Image"):
		return &apiv1alpha1.ImageApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageContent"):
		return &apiv1alpha1.ImageContentApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageContentSourceDownload"):
		return &apiv1alpha1.ImageContentSourceDownloadApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageFilter"):
		return &apiv1alpha1.ImageFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageHash"):
		return &apiv1alpha1.ImageHashApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageImport"):
		return &apiv1alpha1.ImageImportApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageProperties"):
		return &apiv1alpha1.ImagePropertiesApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImagePropertiesHardware"):
		return &apiv1alpha1.ImagePropertiesHardwareApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageResourceSpec"):
		return &apiv1alpha1.ImageResourceSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageResourceStatus"):
		return &apiv1alpha1.ImageResourceStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageSpec"):
		return &apiv1alpha1.ImageSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ImageStatus"):
		return &apiv1alpha1.ImageStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ManagedOptions"):
		return &apiv1alpha1.ManagedOptionsApplyConfiguration{}

	}
	return nil
}
