// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source=client.go -destination=../../../testing/mock/github.com/katanomi/pkg/plugin/storage/client/client.go -package=v1alpha1 Interface
//

// Package v1alpha1 is a generated GoMock package.
package v1alpha1

import (
	context "context"
	reflect "reflect"

	resty "github.com/go-resty/resty/v2"
	client "github.com/katanomi/pkg/plugin/client"
	client0 "github.com/katanomi/pkg/plugin/storage/client"
	gomock "go.uber.org/mock/gomock"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// APIVersion mocks base method.
func (m *MockInterface) APIVersion() *schema.GroupVersion {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIVersion")
	ret0, _ := ret[0].(*schema.GroupVersion)
	return ret0
}

// APIVersion indicates an expected call of APIVersion.
func (mr *MockInterfaceMockRecorder) APIVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIVersion", reflect.TypeOf((*MockInterface)(nil).APIVersion))
}

// Delete mocks base method.
func (m *MockInterface) Delete(ctx context.Context, path string, options ...client.OptionFunc) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, path}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockInterfaceMockRecorder) Delete(ctx, path any, options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, path}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockInterface)(nil).Delete), varargs...)
}

// ForGroupVersion mocks base method.
func (m *MockInterface) ForGroupVersion(gv *schema.GroupVersion) client0.Interface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForGroupVersion", gv)
	ret0, _ := ret[0].(client0.Interface)
	return ret0
}

// ForGroupVersion indicates an expected call of ForGroupVersion.
func (mr *MockInterfaceMockRecorder) ForGroupVersion(gv any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForGroupVersion", reflect.TypeOf((*MockInterface)(nil).ForGroupVersion), gv)
}

// Get mocks base method.
func (m *MockInterface) Get(ctx context.Context, path string, options ...client.OptionFunc) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, path}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockInterfaceMockRecorder) Get(ctx, path any, options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, path}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInterface)(nil).Get), varargs...)
}

// GetResponse mocks base method.
func (m *MockInterface) GetResponse(ctx context.Context, path string, options ...client.OptionFunc) (*resty.Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, path}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetResponse", varargs...)
	ret0, _ := ret[0].(*resty.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResponse indicates an expected call of GetResponse.
func (mr *MockInterfaceMockRecorder) GetResponse(ctx, path any, options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, path}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResponse", reflect.TypeOf((*MockInterface)(nil).GetResponse), varargs...)
}

// Post mocks base method.
func (m *MockInterface) Post(ctx context.Context, path string, options ...client.OptionFunc) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, path}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Post", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Post indicates an expected call of Post.
func (mr *MockInterfaceMockRecorder) Post(ctx, path any, options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, path}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockInterface)(nil).Post), varargs...)
}

// Put mocks base method.
func (m *MockInterface) Put(ctx context.Context, path string, options ...client.OptionFunc) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, path}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Put", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockInterfaceMockRecorder) Put(ctx, path any, options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, path}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockInterface)(nil).Put), varargs...)
}
