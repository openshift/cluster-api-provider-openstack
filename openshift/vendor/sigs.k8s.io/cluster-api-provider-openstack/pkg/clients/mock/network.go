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
// Source: sigs.k8s.io/cluster-api-provider-openstack/pkg/clients (interfaces: NetworkClient)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	extensions "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions"
	attributestags "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	floatingips "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	routers "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	groups "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	rules "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	trunks "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	networks "github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	ports "github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	subnets "github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

// MockNetworkClient is a mock of NetworkClient interface.
type MockNetworkClient struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkClientMockRecorder
}

// MockNetworkClientMockRecorder is the mock recorder for MockNetworkClient.
type MockNetworkClientMockRecorder struct {
	mock *MockNetworkClient
}

// NewMockNetworkClient creates a new mock instance.
func NewMockNetworkClient(ctrl *gomock.Controller) *MockNetworkClient {
	mock := &MockNetworkClient{ctrl: ctrl}
	mock.recorder = &MockNetworkClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworkClient) EXPECT() *MockNetworkClientMockRecorder {
	return m.recorder
}

// AddRouterInterface mocks base method.
func (m *MockNetworkClient) AddRouterInterface(arg0 string, arg1 routers.AddInterfaceOptsBuilder) (*routers.InterfaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRouterInterface", arg0, arg1)
	ret0, _ := ret[0].(*routers.InterfaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRouterInterface indicates an expected call of AddRouterInterface.
func (mr *MockNetworkClientMockRecorder) AddRouterInterface(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRouterInterface", reflect.TypeOf((*MockNetworkClient)(nil).AddRouterInterface), arg0, arg1)
}

// CreateFloatingIP mocks base method.
func (m *MockNetworkClient) CreateFloatingIP(arg0 floatingips.CreateOptsBuilder) (*floatingips.FloatingIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFloatingIP", arg0)
	ret0, _ := ret[0].(*floatingips.FloatingIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFloatingIP indicates an expected call of CreateFloatingIP.
func (mr *MockNetworkClientMockRecorder) CreateFloatingIP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFloatingIP", reflect.TypeOf((*MockNetworkClient)(nil).CreateFloatingIP), arg0)
}

// CreateNetwork mocks base method.
func (m *MockNetworkClient) CreateNetwork(arg0 networks.CreateOptsBuilder) (*networks.Network, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNetwork", arg0)
	ret0, _ := ret[0].(*networks.Network)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNetwork indicates an expected call of CreateNetwork.
func (mr *MockNetworkClientMockRecorder) CreateNetwork(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNetwork", reflect.TypeOf((*MockNetworkClient)(nil).CreateNetwork), arg0)
}

// CreatePort mocks base method.
func (m *MockNetworkClient) CreatePort(arg0 ports.CreateOptsBuilder) (*ports.Port, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePort", arg0)
	ret0, _ := ret[0].(*ports.Port)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePort indicates an expected call of CreatePort.
func (mr *MockNetworkClientMockRecorder) CreatePort(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePort", reflect.TypeOf((*MockNetworkClient)(nil).CreatePort), arg0)
}

// CreateRouter mocks base method.
func (m *MockNetworkClient) CreateRouter(arg0 routers.CreateOptsBuilder) (*routers.Router, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRouter", arg0)
	ret0, _ := ret[0].(*routers.Router)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRouter indicates an expected call of CreateRouter.
func (mr *MockNetworkClientMockRecorder) CreateRouter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRouter", reflect.TypeOf((*MockNetworkClient)(nil).CreateRouter), arg0)
}

// CreateSecGroup mocks base method.
func (m *MockNetworkClient) CreateSecGroup(arg0 groups.CreateOptsBuilder) (*groups.SecGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecGroup", arg0)
	ret0, _ := ret[0].(*groups.SecGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecGroup indicates an expected call of CreateSecGroup.
func (mr *MockNetworkClientMockRecorder) CreateSecGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecGroup", reflect.TypeOf((*MockNetworkClient)(nil).CreateSecGroup), arg0)
}

// CreateSecGroupRule mocks base method.
func (m *MockNetworkClient) CreateSecGroupRule(arg0 rules.CreateOptsBuilder) (*rules.SecGroupRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecGroupRule", arg0)
	ret0, _ := ret[0].(*rules.SecGroupRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecGroupRule indicates an expected call of CreateSecGroupRule.
func (mr *MockNetworkClientMockRecorder) CreateSecGroupRule(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecGroupRule", reflect.TypeOf((*MockNetworkClient)(nil).CreateSecGroupRule), arg0)
}

// CreateSubnet mocks base method.
func (m *MockNetworkClient) CreateSubnet(arg0 subnets.CreateOptsBuilder) (*subnets.Subnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSubnet", arg0)
	ret0, _ := ret[0].(*subnets.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubnet indicates an expected call of CreateSubnet.
func (mr *MockNetworkClientMockRecorder) CreateSubnet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubnet", reflect.TypeOf((*MockNetworkClient)(nil).CreateSubnet), arg0)
}

// CreateTrunk mocks base method.
func (m *MockNetworkClient) CreateTrunk(arg0 trunks.CreateOptsBuilder) (*trunks.Trunk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrunk", arg0)
	ret0, _ := ret[0].(*trunks.Trunk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTrunk indicates an expected call of CreateTrunk.
func (mr *MockNetworkClientMockRecorder) CreateTrunk(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrunk", reflect.TypeOf((*MockNetworkClient)(nil).CreateTrunk), arg0)
}

// DeleteFloatingIP mocks base method.
func (m *MockNetworkClient) DeleteFloatingIP(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFloatingIP", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFloatingIP indicates an expected call of DeleteFloatingIP.
func (mr *MockNetworkClientMockRecorder) DeleteFloatingIP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFloatingIP", reflect.TypeOf((*MockNetworkClient)(nil).DeleteFloatingIP), arg0)
}

// DeleteNetwork mocks base method.
func (m *MockNetworkClient) DeleteNetwork(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNetwork", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNetwork indicates an expected call of DeleteNetwork.
func (mr *MockNetworkClientMockRecorder) DeleteNetwork(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNetwork", reflect.TypeOf((*MockNetworkClient)(nil).DeleteNetwork), arg0)
}

// DeletePort mocks base method.
func (m *MockNetworkClient) DeletePort(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePort", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePort indicates an expected call of DeletePort.
func (mr *MockNetworkClientMockRecorder) DeletePort(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePort", reflect.TypeOf((*MockNetworkClient)(nil).DeletePort), arg0)
}

// DeleteRouter mocks base method.
func (m *MockNetworkClient) DeleteRouter(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRouter", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRouter indicates an expected call of DeleteRouter.
func (mr *MockNetworkClientMockRecorder) DeleteRouter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRouter", reflect.TypeOf((*MockNetworkClient)(nil).DeleteRouter), arg0)
}

// DeleteSecGroup mocks base method.
func (m *MockNetworkClient) DeleteSecGroup(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecGroup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecGroup indicates an expected call of DeleteSecGroup.
func (mr *MockNetworkClientMockRecorder) DeleteSecGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecGroup", reflect.TypeOf((*MockNetworkClient)(nil).DeleteSecGroup), arg0)
}

// DeleteSecGroupRule mocks base method.
func (m *MockNetworkClient) DeleteSecGroupRule(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecGroupRule", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecGroupRule indicates an expected call of DeleteSecGroupRule.
func (mr *MockNetworkClientMockRecorder) DeleteSecGroupRule(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecGroupRule", reflect.TypeOf((*MockNetworkClient)(nil).DeleteSecGroupRule), arg0)
}

// DeleteSubnet mocks base method.
func (m *MockNetworkClient) DeleteSubnet(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubnet", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubnet indicates an expected call of DeleteSubnet.
func (mr *MockNetworkClientMockRecorder) DeleteSubnet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubnet", reflect.TypeOf((*MockNetworkClient)(nil).DeleteSubnet), arg0)
}

// DeleteTrunk mocks base method.
func (m *MockNetworkClient) DeleteTrunk(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrunk", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrunk indicates an expected call of DeleteTrunk.
func (mr *MockNetworkClientMockRecorder) DeleteTrunk(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrunk", reflect.TypeOf((*MockNetworkClient)(nil).DeleteTrunk), arg0)
}

// GetFloatingIP mocks base method.
func (m *MockNetworkClient) GetFloatingIP(arg0 string) (*floatingips.FloatingIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloatingIP", arg0)
	ret0, _ := ret[0].(*floatingips.FloatingIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFloatingIP indicates an expected call of GetFloatingIP.
func (mr *MockNetworkClientMockRecorder) GetFloatingIP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloatingIP", reflect.TypeOf((*MockNetworkClient)(nil).GetFloatingIP), arg0)
}

// GetNetwork mocks base method.
func (m *MockNetworkClient) GetNetwork(arg0 string) (*networks.Network, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetwork", arg0)
	ret0, _ := ret[0].(*networks.Network)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetwork indicates an expected call of GetNetwork.
func (mr *MockNetworkClientMockRecorder) GetNetwork(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetwork", reflect.TypeOf((*MockNetworkClient)(nil).GetNetwork), arg0)
}

// GetPort mocks base method.
func (m *MockNetworkClient) GetPort(arg0 string) (*ports.Port, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPort", arg0)
	ret0, _ := ret[0].(*ports.Port)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPort indicates an expected call of GetPort.
func (mr *MockNetworkClientMockRecorder) GetPort(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPort", reflect.TypeOf((*MockNetworkClient)(nil).GetPort), arg0)
}

// GetRouter mocks base method.
func (m *MockNetworkClient) GetRouter(arg0 string) (*routers.Router, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRouter", arg0)
	ret0, _ := ret[0].(*routers.Router)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRouter indicates an expected call of GetRouter.
func (mr *MockNetworkClientMockRecorder) GetRouter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRouter", reflect.TypeOf((*MockNetworkClient)(nil).GetRouter), arg0)
}

// GetSecGroup mocks base method.
func (m *MockNetworkClient) GetSecGroup(arg0 string) (*groups.SecGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecGroup", arg0)
	ret0, _ := ret[0].(*groups.SecGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecGroup indicates an expected call of GetSecGroup.
func (mr *MockNetworkClientMockRecorder) GetSecGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecGroup", reflect.TypeOf((*MockNetworkClient)(nil).GetSecGroup), arg0)
}

// GetSecGroupRule mocks base method.
func (m *MockNetworkClient) GetSecGroupRule(arg0 string) (*rules.SecGroupRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecGroupRule", arg0)
	ret0, _ := ret[0].(*rules.SecGroupRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecGroupRule indicates an expected call of GetSecGroupRule.
func (mr *MockNetworkClientMockRecorder) GetSecGroupRule(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecGroupRule", reflect.TypeOf((*MockNetworkClient)(nil).GetSecGroupRule), arg0)
}

// GetSubnet mocks base method.
func (m *MockNetworkClient) GetSubnet(arg0 string) (*subnets.Subnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnet", arg0)
	ret0, _ := ret[0].(*subnets.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnet indicates an expected call of GetSubnet.
func (mr *MockNetworkClientMockRecorder) GetSubnet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnet", reflect.TypeOf((*MockNetworkClient)(nil).GetSubnet), arg0)
}

// ListExtensions mocks base method.
func (m *MockNetworkClient) ListExtensions() ([]extensions.Extension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListExtensions")
	ret0, _ := ret[0].([]extensions.Extension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListExtensions indicates an expected call of ListExtensions.
func (mr *MockNetworkClientMockRecorder) ListExtensions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListExtensions", reflect.TypeOf((*MockNetworkClient)(nil).ListExtensions))
}

// ListFloatingIP mocks base method.
func (m *MockNetworkClient) ListFloatingIP(arg0 floatingips.ListOptsBuilder) ([]floatingips.FloatingIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFloatingIP", arg0)
	ret0, _ := ret[0].([]floatingips.FloatingIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFloatingIP indicates an expected call of ListFloatingIP.
func (mr *MockNetworkClientMockRecorder) ListFloatingIP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFloatingIP", reflect.TypeOf((*MockNetworkClient)(nil).ListFloatingIP), arg0)
}

// ListNetwork mocks base method.
func (m *MockNetworkClient) ListNetwork(arg0 networks.ListOptsBuilder) ([]networks.Network, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNetwork", arg0)
	ret0, _ := ret[0].([]networks.Network)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNetwork indicates an expected call of ListNetwork.
func (mr *MockNetworkClientMockRecorder) ListNetwork(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNetwork", reflect.TypeOf((*MockNetworkClient)(nil).ListNetwork), arg0)
}

// ListPort mocks base method.
func (m *MockNetworkClient) ListPort(arg0 ports.ListOptsBuilder) ([]ports.Port, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPort", arg0)
	ret0, _ := ret[0].([]ports.Port)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPort indicates an expected call of ListPort.
func (mr *MockNetworkClientMockRecorder) ListPort(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPort", reflect.TypeOf((*MockNetworkClient)(nil).ListPort), arg0)
}

// ListRouter mocks base method.
func (m *MockNetworkClient) ListRouter(arg0 routers.ListOpts) ([]routers.Router, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRouter", arg0)
	ret0, _ := ret[0].([]routers.Router)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRouter indicates an expected call of ListRouter.
func (mr *MockNetworkClientMockRecorder) ListRouter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRouter", reflect.TypeOf((*MockNetworkClient)(nil).ListRouter), arg0)
}

// ListSecGroup mocks base method.
func (m *MockNetworkClient) ListSecGroup(arg0 groups.ListOpts) ([]groups.SecGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecGroup", arg0)
	ret0, _ := ret[0].([]groups.SecGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecGroup indicates an expected call of ListSecGroup.
func (mr *MockNetworkClientMockRecorder) ListSecGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecGroup", reflect.TypeOf((*MockNetworkClient)(nil).ListSecGroup), arg0)
}

// ListSecGroupRule mocks base method.
func (m *MockNetworkClient) ListSecGroupRule(arg0 rules.ListOpts) ([]rules.SecGroupRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecGroupRule", arg0)
	ret0, _ := ret[0].([]rules.SecGroupRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecGroupRule indicates an expected call of ListSecGroupRule.
func (mr *MockNetworkClientMockRecorder) ListSecGroupRule(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecGroupRule", reflect.TypeOf((*MockNetworkClient)(nil).ListSecGroupRule), arg0)
}

// ListSubnet mocks base method.
func (m *MockNetworkClient) ListSubnet(arg0 subnets.ListOptsBuilder) ([]subnets.Subnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSubnet", arg0)
	ret0, _ := ret[0].([]subnets.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSubnet indicates an expected call of ListSubnet.
func (mr *MockNetworkClientMockRecorder) ListSubnet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSubnet", reflect.TypeOf((*MockNetworkClient)(nil).ListSubnet), arg0)
}

// ListTrunk mocks base method.
func (m *MockNetworkClient) ListTrunk(arg0 trunks.ListOptsBuilder) ([]trunks.Trunk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTrunk", arg0)
	ret0, _ := ret[0].([]trunks.Trunk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTrunk indicates an expected call of ListTrunk.
func (mr *MockNetworkClientMockRecorder) ListTrunk(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTrunk", reflect.TypeOf((*MockNetworkClient)(nil).ListTrunk), arg0)
}

// ListTrunkSubports mocks base method.
func (m *MockNetworkClient) ListTrunkSubports(arg0 string) ([]trunks.Subport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTrunkSubports", arg0)
	ret0, _ := ret[0].([]trunks.Subport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTrunkSubports indicates an expected call of ListTrunkSubports.
func (mr *MockNetworkClientMockRecorder) ListTrunkSubports(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTrunkSubports", reflect.TypeOf((*MockNetworkClient)(nil).ListTrunkSubports), arg0)
}

// RemoveRouterInterface mocks base method.
func (m *MockNetworkClient) RemoveRouterInterface(arg0 string, arg1 routers.RemoveInterfaceOptsBuilder) (*routers.InterfaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRouterInterface", arg0, arg1)
	ret0, _ := ret[0].(*routers.InterfaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveRouterInterface indicates an expected call of RemoveRouterInterface.
func (mr *MockNetworkClientMockRecorder) RemoveRouterInterface(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRouterInterface", reflect.TypeOf((*MockNetworkClient)(nil).RemoveRouterInterface), arg0, arg1)
}

// RemoveSubports mocks base method.
func (m *MockNetworkClient) RemoveSubports(arg0 string, arg1 trunks.RemoveSubportsOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSubports", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSubports indicates an expected call of RemoveSubports.
func (mr *MockNetworkClientMockRecorder) RemoveSubports(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSubports", reflect.TypeOf((*MockNetworkClient)(nil).RemoveSubports), arg0, arg1)
}

// ReplaceAllAttributesTags mocks base method.
func (m *MockNetworkClient) ReplaceAllAttributesTags(arg0, arg1 string, arg2 attributestags.ReplaceAllOptsBuilder) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceAllAttributesTags", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceAllAttributesTags indicates an expected call of ReplaceAllAttributesTags.
func (mr *MockNetworkClientMockRecorder) ReplaceAllAttributesTags(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceAllAttributesTags", reflect.TypeOf((*MockNetworkClient)(nil).ReplaceAllAttributesTags), arg0, arg1, arg2)
}

// UpdateFloatingIP mocks base method.
func (m *MockNetworkClient) UpdateFloatingIP(arg0 string, arg1 floatingips.UpdateOptsBuilder) (*floatingips.FloatingIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFloatingIP", arg0, arg1)
	ret0, _ := ret[0].(*floatingips.FloatingIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFloatingIP indicates an expected call of UpdateFloatingIP.
func (mr *MockNetworkClientMockRecorder) UpdateFloatingIP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFloatingIP", reflect.TypeOf((*MockNetworkClient)(nil).UpdateFloatingIP), arg0, arg1)
}

// UpdateNetwork mocks base method.
func (m *MockNetworkClient) UpdateNetwork(arg0 string, arg1 networks.UpdateOptsBuilder) (*networks.Network, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNetwork", arg0, arg1)
	ret0, _ := ret[0].(*networks.Network)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNetwork indicates an expected call of UpdateNetwork.
func (mr *MockNetworkClientMockRecorder) UpdateNetwork(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNetwork", reflect.TypeOf((*MockNetworkClient)(nil).UpdateNetwork), arg0, arg1)
}

// UpdatePort mocks base method.
func (m *MockNetworkClient) UpdatePort(arg0 string, arg1 ports.UpdateOptsBuilder) (*ports.Port, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePort", arg0, arg1)
	ret0, _ := ret[0].(*ports.Port)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePort indicates an expected call of UpdatePort.
func (mr *MockNetworkClientMockRecorder) UpdatePort(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePort", reflect.TypeOf((*MockNetworkClient)(nil).UpdatePort), arg0, arg1)
}

// UpdateRouter mocks base method.
func (m *MockNetworkClient) UpdateRouter(arg0 string, arg1 routers.UpdateOptsBuilder) (*routers.Router, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRouter", arg0, arg1)
	ret0, _ := ret[0].(*routers.Router)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRouter indicates an expected call of UpdateRouter.
func (mr *MockNetworkClientMockRecorder) UpdateRouter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRouter", reflect.TypeOf((*MockNetworkClient)(nil).UpdateRouter), arg0, arg1)
}

// UpdateSecGroup mocks base method.
func (m *MockNetworkClient) UpdateSecGroup(arg0 string, arg1 groups.UpdateOptsBuilder) (*groups.SecGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecGroup", arg0, arg1)
	ret0, _ := ret[0].(*groups.SecGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSecGroup indicates an expected call of UpdateSecGroup.
func (mr *MockNetworkClientMockRecorder) UpdateSecGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecGroup", reflect.TypeOf((*MockNetworkClient)(nil).UpdateSecGroup), arg0, arg1)
}

// UpdateSubnet mocks base method.
func (m *MockNetworkClient) UpdateSubnet(arg0 string, arg1 subnets.UpdateOptsBuilder) (*subnets.Subnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSubnet", arg0, arg1)
	ret0, _ := ret[0].(*subnets.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSubnet indicates an expected call of UpdateSubnet.
func (mr *MockNetworkClientMockRecorder) UpdateSubnet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSubnet", reflect.TypeOf((*MockNetworkClient)(nil).UpdateSubnet), arg0, arg1)
}
