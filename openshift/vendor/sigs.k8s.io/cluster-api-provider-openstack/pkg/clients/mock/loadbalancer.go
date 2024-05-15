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

// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/cluster-api-provider-openstack/pkg/clients (interfaces: LbClient)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	apiversions "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/apiversions"
	listeners "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	loadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	monitors "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	pools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	providers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/providers"
)

// MockLbClient is a mock of LbClient interface.
type MockLbClient struct {
	ctrl     *gomock.Controller
	recorder *MockLbClientMockRecorder
}

// MockLbClientMockRecorder is the mock recorder for MockLbClient.
type MockLbClientMockRecorder struct {
	mock *MockLbClient
}

// NewMockLbClient creates a new mock instance.
func NewMockLbClient(ctrl *gomock.Controller) *MockLbClient {
	mock := &MockLbClient{ctrl: ctrl}
	mock.recorder = &MockLbClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLbClient) EXPECT() *MockLbClientMockRecorder {
	return m.recorder
}

// CreateListener mocks base method.
func (m *MockLbClient) CreateListener(arg0 listeners.CreateOptsBuilder) (*listeners.Listener, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateListener", arg0)
	ret0, _ := ret[0].(*listeners.Listener)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateListener indicates an expected call of CreateListener.
func (mr *MockLbClientMockRecorder) CreateListener(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateListener", reflect.TypeOf((*MockLbClient)(nil).CreateListener), arg0)
}

// CreateLoadBalancer mocks base method.
func (m *MockLbClient) CreateLoadBalancer(arg0 loadbalancers.CreateOptsBuilder) (*loadbalancers.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLoadBalancer", arg0)
	ret0, _ := ret[0].(*loadbalancers.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLoadBalancer indicates an expected call of CreateLoadBalancer.
func (mr *MockLbClientMockRecorder) CreateLoadBalancer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoadBalancer", reflect.TypeOf((*MockLbClient)(nil).CreateLoadBalancer), arg0)
}

// CreateMonitor mocks base method.
func (m *MockLbClient) CreateMonitor(arg0 monitors.CreateOptsBuilder) (*monitors.Monitor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMonitor", arg0)
	ret0, _ := ret[0].(*monitors.Monitor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMonitor indicates an expected call of CreateMonitor.
func (mr *MockLbClientMockRecorder) CreateMonitor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMonitor", reflect.TypeOf((*MockLbClient)(nil).CreateMonitor), arg0)
}

// CreatePool mocks base method.
func (m *MockLbClient) CreatePool(arg0 pools.CreateOptsBuilder) (*pools.Pool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePool", arg0)
	ret0, _ := ret[0].(*pools.Pool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePool indicates an expected call of CreatePool.
func (mr *MockLbClientMockRecorder) CreatePool(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePool", reflect.TypeOf((*MockLbClient)(nil).CreatePool), arg0)
}

// CreatePoolMember mocks base method.
func (m *MockLbClient) CreatePoolMember(arg0 string, arg1 pools.CreateMemberOptsBuilder) (*pools.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePoolMember", arg0, arg1)
	ret0, _ := ret[0].(*pools.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePoolMember indicates an expected call of CreatePoolMember.
func (mr *MockLbClientMockRecorder) CreatePoolMember(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePoolMember", reflect.TypeOf((*MockLbClient)(nil).CreatePoolMember), arg0, arg1)
}

// DeleteListener mocks base method.
func (m *MockLbClient) DeleteListener(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteListener", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteListener indicates an expected call of DeleteListener.
func (mr *MockLbClientMockRecorder) DeleteListener(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteListener", reflect.TypeOf((*MockLbClient)(nil).DeleteListener), arg0)
}

// DeleteLoadBalancer mocks base method.
func (m *MockLbClient) DeleteLoadBalancer(arg0 string, arg1 loadbalancers.DeleteOptsBuilder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLoadBalancer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLoadBalancer indicates an expected call of DeleteLoadBalancer.
func (mr *MockLbClientMockRecorder) DeleteLoadBalancer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLoadBalancer", reflect.TypeOf((*MockLbClient)(nil).DeleteLoadBalancer), arg0, arg1)
}

// DeleteMonitor mocks base method.
func (m *MockLbClient) DeleteMonitor(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMonitor", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMonitor indicates an expected call of DeleteMonitor.
func (mr *MockLbClientMockRecorder) DeleteMonitor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMonitor", reflect.TypeOf((*MockLbClient)(nil).DeleteMonitor), arg0)
}

// DeletePool mocks base method.
func (m *MockLbClient) DeletePool(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePool", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePool indicates an expected call of DeletePool.
func (mr *MockLbClientMockRecorder) DeletePool(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePool", reflect.TypeOf((*MockLbClient)(nil).DeletePool), arg0)
}

// DeletePoolMember mocks base method.
func (m *MockLbClient) DeletePoolMember(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePoolMember", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePoolMember indicates an expected call of DeletePoolMember.
func (mr *MockLbClientMockRecorder) DeletePoolMember(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePoolMember", reflect.TypeOf((*MockLbClient)(nil).DeletePoolMember), arg0, arg1)
}

// GetListener mocks base method.
func (m *MockLbClient) GetListener(arg0 string) (*listeners.Listener, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListener", arg0)
	ret0, _ := ret[0].(*listeners.Listener)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListener indicates an expected call of GetListener.
func (mr *MockLbClientMockRecorder) GetListener(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListener", reflect.TypeOf((*MockLbClient)(nil).GetListener), arg0)
}

// GetLoadBalancer mocks base method.
func (m *MockLbClient) GetLoadBalancer(arg0 string) (*loadbalancers.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoadBalancer", arg0)
	ret0, _ := ret[0].(*loadbalancers.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoadBalancer indicates an expected call of GetLoadBalancer.
func (mr *MockLbClientMockRecorder) GetLoadBalancer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoadBalancer", reflect.TypeOf((*MockLbClient)(nil).GetLoadBalancer), arg0)
}

// GetPool mocks base method.
func (m *MockLbClient) GetPool(arg0 string) (*pools.Pool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPool", arg0)
	ret0, _ := ret[0].(*pools.Pool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPool indicates an expected call of GetPool.
func (mr *MockLbClientMockRecorder) GetPool(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPool", reflect.TypeOf((*MockLbClient)(nil).GetPool), arg0)
}

// ListListeners mocks base method.
func (m *MockLbClient) ListListeners(arg0 listeners.ListOptsBuilder) ([]listeners.Listener, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListListeners", arg0)
	ret0, _ := ret[0].([]listeners.Listener)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListListeners indicates an expected call of ListListeners.
func (mr *MockLbClientMockRecorder) ListListeners(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListListeners", reflect.TypeOf((*MockLbClient)(nil).ListListeners), arg0)
}

// ListLoadBalancerProviders mocks base method.
func (m *MockLbClient) ListLoadBalancerProviders() ([]providers.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLoadBalancerProviders")
	ret0, _ := ret[0].([]providers.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLoadBalancerProviders indicates an expected call of ListLoadBalancerProviders.
func (mr *MockLbClientMockRecorder) ListLoadBalancerProviders() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLoadBalancerProviders", reflect.TypeOf((*MockLbClient)(nil).ListLoadBalancerProviders))
}

// ListLoadBalancers mocks base method.
func (m *MockLbClient) ListLoadBalancers(arg0 loadbalancers.ListOptsBuilder) ([]loadbalancers.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLoadBalancers", arg0)
	ret0, _ := ret[0].([]loadbalancers.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLoadBalancers indicates an expected call of ListLoadBalancers.
func (mr *MockLbClientMockRecorder) ListLoadBalancers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLoadBalancers", reflect.TypeOf((*MockLbClient)(nil).ListLoadBalancers), arg0)
}

// ListMonitors mocks base method.
func (m *MockLbClient) ListMonitors(arg0 monitors.ListOptsBuilder) ([]monitors.Monitor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMonitors", arg0)
	ret0, _ := ret[0].([]monitors.Monitor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMonitors indicates an expected call of ListMonitors.
func (mr *MockLbClientMockRecorder) ListMonitors(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMonitors", reflect.TypeOf((*MockLbClient)(nil).ListMonitors), arg0)
}

// ListOctaviaVersions mocks base method.
func (m *MockLbClient) ListOctaviaVersions() ([]apiversions.APIVersion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOctaviaVersions")
	ret0, _ := ret[0].([]apiversions.APIVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOctaviaVersions indicates an expected call of ListOctaviaVersions.
func (mr *MockLbClientMockRecorder) ListOctaviaVersions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOctaviaVersions", reflect.TypeOf((*MockLbClient)(nil).ListOctaviaVersions))
}

// ListPoolMember mocks base method.
func (m *MockLbClient) ListPoolMember(arg0 string, arg1 pools.ListMembersOptsBuilder) ([]pools.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPoolMember", arg0, arg1)
	ret0, _ := ret[0].([]pools.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPoolMember indicates an expected call of ListPoolMember.
func (mr *MockLbClientMockRecorder) ListPoolMember(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPoolMember", reflect.TypeOf((*MockLbClient)(nil).ListPoolMember), arg0, arg1)
}

// ListPools mocks base method.
func (m *MockLbClient) ListPools(arg0 pools.ListOptsBuilder) ([]pools.Pool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPools", arg0)
	ret0, _ := ret[0].([]pools.Pool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPools indicates an expected call of ListPools.
func (mr *MockLbClientMockRecorder) ListPools(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPools", reflect.TypeOf((*MockLbClient)(nil).ListPools), arg0)
}

// UpdateListener mocks base method.
func (m *MockLbClient) UpdateListener(arg0 string, arg1 listeners.UpdateOpts) (*listeners.Listener, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateListener", arg0, arg1)
	ret0, _ := ret[0].(*listeners.Listener)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateListener indicates an expected call of UpdateListener.
func (mr *MockLbClientMockRecorder) UpdateListener(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateListener", reflect.TypeOf((*MockLbClient)(nil).UpdateListener), arg0, arg1)
}