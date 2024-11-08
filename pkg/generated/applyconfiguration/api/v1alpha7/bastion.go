/*
Copyright 2024 The Kubernetes Authors.

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

package v1alpha7

// BastionApplyConfiguration represents an declarative configuration of the Bastion type for use
// with apply.
type BastionApplyConfiguration struct {
	Enabled          *bool                                   `json:"enabled,omitempty"`
	Instance         *OpenStackMachineSpecApplyConfiguration `json:"instance,omitempty"`
	AvailabilityZone *string                                 `json:"availabilityZone,omitempty"`
}

// BastionApplyConfiguration constructs an declarative configuration of the Bastion type for use with
// apply.
func Bastion() *BastionApplyConfiguration {
	return &BastionApplyConfiguration{}
}

// WithEnabled sets the Enabled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Enabled field is set to the value of the last call.
func (b *BastionApplyConfiguration) WithEnabled(value bool) *BastionApplyConfiguration {
	b.Enabled = &value
	return b
}

// WithInstance sets the Instance field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Instance field is set to the value of the last call.
func (b *BastionApplyConfiguration) WithInstance(value *OpenStackMachineSpecApplyConfiguration) *BastionApplyConfiguration {
	b.Instance = value
	return b
}

// WithAvailabilityZone sets the AvailabilityZone field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AvailabilityZone field is set to the value of the last call.
func (b *BastionApplyConfiguration) WithAvailabilityZone(value string) *BastionApplyConfiguration {
	b.AvailabilityZone = &value
	return b
}
