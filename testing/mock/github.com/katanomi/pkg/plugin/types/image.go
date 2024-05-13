// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: ImageConfigGetter)
//
// Generated by this command:
//
//	mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/image.go github.com/katanomi/pkg/plugin/types ImageConfigGetter
//

// Package types is a generated GoMock package.
package types

import (
	context "context"
	reflect "reflect"

	v1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockImageConfigGetter is a mock of ImageConfigGetter interface.
type MockImageConfigGetter struct {
	ctrl     *gomock.Controller
	recorder *MockImageConfigGetterMockRecorder
}

// MockImageConfigGetterMockRecorder is the mock recorder for MockImageConfigGetter.
type MockImageConfigGetterMockRecorder struct {
	mock *MockImageConfigGetter
}

// NewMockImageConfigGetter creates a new mock instance.
func NewMockImageConfigGetter(ctrl *gomock.Controller) *MockImageConfigGetter {
	mock := &MockImageConfigGetter{ctrl: ctrl}
	mock.recorder = &MockImageConfigGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageConfigGetter) EXPECT() *MockImageConfigGetterMockRecorder {
	return m.recorder
}

// GetImageConfig mocks base method.
func (m *MockImageConfigGetter) GetImageConfig(arg0 context.Context, arg1 v1alpha1.ArtifactOptions) (*v1alpha1.ImageConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageConfig", arg0, arg1)
	ret0, _ := ret[0].(*v1alpha1.ImageConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageConfig indicates an expected call of GetImageConfig.
func (mr *MockImageConfigGetterMockRecorder) GetImageConfig(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageConfig", reflect.TypeOf((*MockImageConfigGetter)(nil).GetImageConfig), arg0, arg1)
}

// Path mocks base method.
func (m *MockImageConfigGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockImageConfigGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockImageConfigGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockImageConfigGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockImageConfigGetterMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockImageConfigGetter)(nil).Setup), arg0, arg1)
}
