// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/yolo-sh/aws-cloud-provider/service (interfaces: UserConfigLoader)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	gomock "github.com/golang/mock/gomock"
	userconfig "github.com/yolo-sh/aws-cloud-provider/userconfig"
)

// MockUserConfigLoader is a mock of UserConfigLoader interface.
type MockUserConfigLoader struct {
	ctrl     *gomock.Controller
	recorder *MockUserConfigLoaderMockRecorder
}

// MockUserConfigLoaderMockRecorder is the mock recorder for MockUserConfigLoader.
type MockUserConfigLoaderMockRecorder struct {
	mock *MockUserConfigLoader
}

// NewMockUserConfigLoader creates a new mock instance.
func NewMockUserConfigLoader(ctrl *gomock.Controller) *MockUserConfigLoader {
	mock := &MockUserConfigLoader{ctrl: ctrl}
	mock.recorder = &MockUserConfigLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserConfigLoader) EXPECT() *MockUserConfigLoaderMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockUserConfigLoader) Load(arg0 *userconfig.Config) (aws.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0)
	ret0, _ := ret[0].(aws.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockUserConfigLoaderMockRecorder) Load(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockUserConfigLoader)(nil).Load), arg0)
}