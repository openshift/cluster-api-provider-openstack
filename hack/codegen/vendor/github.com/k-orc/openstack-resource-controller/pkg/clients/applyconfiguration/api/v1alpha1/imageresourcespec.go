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

package v1alpha1

import (
	v1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"
)

// ImageResourceSpecApplyConfiguration represents an declarative configuration of the ImageResourceSpec type for use
// with apply.
type ImageResourceSpecApplyConfiguration struct {
	Name       *string                            `json:"name,omitempty"`
	Protected  *bool                              `json:"protected,omitempty"`
	Tags       []v1alpha1.ImageTag                `json:"tags,omitempty"`
	Visibility *v1alpha1.ImageVisibility          `json:"visibility,omitempty"`
	Properties *ImagePropertiesApplyConfiguration `json:"properties,omitempty"`
	Content    *ImageContentApplyConfiguration    `json:"content,omitempty"`
}

// ImageResourceSpecApplyConfiguration constructs an declarative configuration of the ImageResourceSpec type for use with
// apply.
func ImageResourceSpec() *ImageResourceSpecApplyConfiguration {
	return &ImageResourceSpecApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ImageResourceSpecApplyConfiguration) WithName(value string) *ImageResourceSpecApplyConfiguration {
	b.Name = &value
	return b
}

// WithProtected sets the Protected field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Protected field is set to the value of the last call.
func (b *ImageResourceSpecApplyConfiguration) WithProtected(value bool) *ImageResourceSpecApplyConfiguration {
	b.Protected = &value
	return b
}

// WithTags adds the given value to the Tags field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tags field.
func (b *ImageResourceSpecApplyConfiguration) WithTags(values ...v1alpha1.ImageTag) *ImageResourceSpecApplyConfiguration {
	for i := range values {
		b.Tags = append(b.Tags, values[i])
	}
	return b
}

// WithVisibility sets the Visibility field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Visibility field is set to the value of the last call.
func (b *ImageResourceSpecApplyConfiguration) WithVisibility(value v1alpha1.ImageVisibility) *ImageResourceSpecApplyConfiguration {
	b.Visibility = &value
	return b
}

// WithProperties sets the Properties field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Properties field is set to the value of the last call.
func (b *ImageResourceSpecApplyConfiguration) WithProperties(value *ImagePropertiesApplyConfiguration) *ImageResourceSpecApplyConfiguration {
	b.Properties = value
	return b
}

// WithContent sets the Content field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Content field is set to the value of the last call.
func (b *ImageResourceSpecApplyConfiguration) WithContent(value *ImageContentApplyConfiguration) *ImageResourceSpecApplyConfiguration {
	b.Content = value
	return b
}