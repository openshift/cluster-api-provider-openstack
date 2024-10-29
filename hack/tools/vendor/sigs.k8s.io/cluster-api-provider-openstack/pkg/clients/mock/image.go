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
// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/cluster-api-provider-openstack/pkg/clients (interfaces: ImageClient)
//
// Generated by this command:
//
//	mockgen -package mock -destination=image.go sigs.k8s.io/cluster-api-provider-openstack/pkg/clients ImageClient
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	io "io"
	reflect "reflect"

	imageimport "github.com/gophercloud/gophercloud/v2/openstack/image/v2/imageimport"
	images "github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	gomock "go.uber.org/mock/gomock"
)

// MockImageClient is a mock of ImageClient interface.
type MockImageClient struct {
	ctrl     *gomock.Controller
	recorder *MockImageClientMockRecorder
}

// MockImageClientMockRecorder is the mock recorder for MockImageClient.
type MockImageClientMockRecorder struct {
	mock *MockImageClient
}

// NewMockImageClient creates a new mock instance.
func NewMockImageClient(ctrl *gomock.Controller) *MockImageClient {
	mock := &MockImageClient{ctrl: ctrl}
	mock.recorder = &MockImageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageClient) EXPECT() *MockImageClientMockRecorder {
	return m.recorder
}

// CreateImage mocks base method.
func (m *MockImageClient) CreateImage(arg0 context.Context, arg1 images.CreateOptsBuilder) (*images.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImage", arg0, arg1)
	ret0, _ := ret[0].(*images.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateImage indicates an expected call of CreateImage.
func (mr *MockImageClientMockRecorder) CreateImage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImage", reflect.TypeOf((*MockImageClient)(nil).CreateImage), arg0, arg1)
}

// CreateImport mocks base method.
func (m *MockImageClient) CreateImport(arg0 context.Context, arg1 string, arg2 imageimport.CreateOptsBuilder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImport", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateImport indicates an expected call of CreateImport.
func (mr *MockImageClientMockRecorder) CreateImport(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImport", reflect.TypeOf((*MockImageClient)(nil).CreateImport), arg0, arg1, arg2)
}

// DeleteImage mocks base method.
func (m *MockImageClient) DeleteImage(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImage indicates an expected call of DeleteImage.
func (mr *MockImageClientMockRecorder) DeleteImage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockImageClient)(nil).DeleteImage), arg0, arg1)
}

// GetImage mocks base method.
func (m *MockImageClient) GetImage(arg0 string) (*images.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImage", arg0)
	ret0, _ := ret[0].(*images.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImage indicates an expected call of GetImage.
func (mr *MockImageClientMockRecorder) GetImage(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImage", reflect.TypeOf((*MockImageClient)(nil).GetImage), arg0)
}

// GetImportInfo mocks base method.
func (m *MockImageClient) GetImportInfo(arg0 context.Context) (*imageimport.ImportInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImportInfo", arg0)
	ret0, _ := ret[0].(*imageimport.ImportInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImportInfo indicates an expected call of GetImportInfo.
func (mr *MockImageClientMockRecorder) GetImportInfo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImportInfo", reflect.TypeOf((*MockImageClient)(nil).GetImportInfo), arg0)
}

// ListImages mocks base method.
func (m *MockImageClient) ListImages(arg0 images.ListOptsBuilder) ([]images.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListImages", arg0)
	ret0, _ := ret[0].([]images.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListImages indicates an expected call of ListImages.
func (mr *MockImageClientMockRecorder) ListImages(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListImages", reflect.TypeOf((*MockImageClient)(nil).ListImages), arg0)
}

// UploadData mocks base method.
func (m *MockImageClient) UploadData(arg0 context.Context, arg1 string, arg2 io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadData", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadData indicates an expected call of UploadData.
func (mr *MockImageClientMockRecorder) UploadData(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadData", reflect.TypeOf((*MockImageClient)(nil).UploadData), arg0, arg1, arg2)
}
