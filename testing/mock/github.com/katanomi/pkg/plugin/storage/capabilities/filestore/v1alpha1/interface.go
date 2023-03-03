// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package v1alpha1 is a generated GoMock package.
package v1alpha1

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/katanomi/pkg/apis/storage/v1alpha1"
	v1alpha10 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockFileStoreCapable is a mock of FileStoreCapable interface.
type MockFileStoreCapable struct {
	ctrl     *gomock.Controller
	recorder *MockFileStoreCapableMockRecorder
}

// MockFileStoreCapableMockRecorder is the mock recorder for MockFileStoreCapable.
type MockFileStoreCapableMockRecorder struct {
	mock *MockFileStoreCapable
}

// NewMockFileStoreCapable creates a new mock instance.
func NewMockFileStoreCapable(ctrl *gomock.Controller) *MockFileStoreCapable {
	mock := &MockFileStoreCapable{ctrl: ctrl}
	mock.recorder = &MockFileStoreCapableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStoreCapable) EXPECT() *MockFileStoreCapableMockRecorder {
	return m.recorder
}

// DeleteFileObject mocks base method.
func (m *MockFileStoreCapable) DeleteFileObject(ctx context.Context, objectName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFileObject", ctx, objectName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFileObject indicates an expected call of DeleteFileObject.
func (mr *MockFileStoreCapableMockRecorder) DeleteFileObject(ctx, objectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFileObject", reflect.TypeOf((*MockFileStoreCapable)(nil).DeleteFileObject), ctx, objectName)
}

// GetFileMeta mocks base method.
func (m *MockFileStoreCapable) GetFileMeta(ctx context.Context, objectName string) (*v1alpha1.FileMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileMeta", ctx, objectName)
	ret0, _ := ret[0].(*v1alpha1.FileMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileMeta indicates an expected call of GetFileMeta.
func (mr *MockFileStoreCapableMockRecorder) GetFileMeta(ctx, objectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileMeta", reflect.TypeOf((*MockFileStoreCapable)(nil).GetFileMeta), ctx, objectName)
}

// GetFileObject mocks base method.
func (m *MockFileStoreCapable) GetFileObject(ctx context.Context, objectName string) (*v1alpha10.FileObject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileObject", ctx, objectName)
	ret0, _ := ret[0].(*v1alpha10.FileObject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileObject indicates an expected call of GetFileObject.
func (mr *MockFileStoreCapableMockRecorder) GetFileObject(ctx, objectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileObject", reflect.TypeOf((*MockFileStoreCapable)(nil).GetFileObject), ctx, objectName)
}

// ListFileMetas mocks base method.
func (m *MockFileStoreCapable) ListFileMetas(ctx context.Context, opt *v1.ListOptions) ([]v1alpha1.FileMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFileMetas", ctx, opt)
	ret0, _ := ret[0].([]v1alpha1.FileMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFileMetas indicates an expected call of ListFileMetas.
func (mr *MockFileStoreCapableMockRecorder) ListFileMetas(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFileMetas", reflect.TypeOf((*MockFileStoreCapable)(nil).ListFileMetas), ctx, opt)
}

// PutFileObject mocks base method.
func (m *MockFileStoreCapable) PutFileObject(ctx context.Context, obj *v1alpha10.FileObject) (*v1alpha1.FileMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutFileObject", ctx, obj)
	ret0, _ := ret[0].(*v1alpha1.FileMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutFileObject indicates an expected call of PutFileObject.
func (mr *MockFileStoreCapableMockRecorder) PutFileObject(ctx, obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutFileObject", reflect.TypeOf((*MockFileStoreCapable)(nil).PutFileObject), ctx, obj)
}

// MockFileObjectInterface is a mock of FileObjectInterface interface.
type MockFileObjectInterface struct {
	ctrl     *gomock.Controller
	recorder *MockFileObjectInterfaceMockRecorder
}

// MockFileObjectInterfaceMockRecorder is the mock recorder for MockFileObjectInterface.
type MockFileObjectInterfaceMockRecorder struct {
	mock *MockFileObjectInterface
}

// NewMockFileObjectInterface creates a new mock instance.
func NewMockFileObjectInterface(ctrl *gomock.Controller) *MockFileObjectInterface {
	mock := &MockFileObjectInterface{ctrl: ctrl}
	mock.recorder = &MockFileObjectInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileObjectInterface) EXPECT() *MockFileObjectInterfaceMockRecorder {
	return m.recorder
}

// DeleteFileObject mocks base method.
func (m *MockFileObjectInterface) DeleteFileObject(ctx context.Context, objectName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFileObject", ctx, objectName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFileObject indicates an expected call of DeleteFileObject.
func (mr *MockFileObjectInterfaceMockRecorder) DeleteFileObject(ctx, objectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFileObject", reflect.TypeOf((*MockFileObjectInterface)(nil).DeleteFileObject), ctx, objectName)
}

// GetFileObject mocks base method.
func (m *MockFileObjectInterface) GetFileObject(ctx context.Context, objectName string) (*v1alpha10.FileObject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileObject", ctx, objectName)
	ret0, _ := ret[0].(*v1alpha10.FileObject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileObject indicates an expected call of GetFileObject.
func (mr *MockFileObjectInterfaceMockRecorder) GetFileObject(ctx, objectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileObject", reflect.TypeOf((*MockFileObjectInterface)(nil).GetFileObject), ctx, objectName)
}

// PutFileObject mocks base method.
func (m *MockFileObjectInterface) PutFileObject(ctx context.Context, obj *v1alpha10.FileObject) (*v1alpha1.FileMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutFileObject", ctx, obj)
	ret0, _ := ret[0].(*v1alpha1.FileMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutFileObject indicates an expected call of PutFileObject.
func (mr *MockFileObjectInterfaceMockRecorder) PutFileObject(ctx, obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutFileObject", reflect.TypeOf((*MockFileObjectInterface)(nil).PutFileObject), ctx, obj)
}

// MockFileMetaInterface is a mock of FileMetaInterface interface.
type MockFileMetaInterface struct {
	ctrl     *gomock.Controller
	recorder *MockFileMetaInterfaceMockRecorder
}

// MockFileMetaInterfaceMockRecorder is the mock recorder for MockFileMetaInterface.
type MockFileMetaInterfaceMockRecorder struct {
	mock *MockFileMetaInterface
}

// NewMockFileMetaInterface creates a new mock instance.
func NewMockFileMetaInterface(ctrl *gomock.Controller) *MockFileMetaInterface {
	mock := &MockFileMetaInterface{ctrl: ctrl}
	mock.recorder = &MockFileMetaInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileMetaInterface) EXPECT() *MockFileMetaInterfaceMockRecorder {
	return m.recorder
}

// GetFileMeta mocks base method.
func (m *MockFileMetaInterface) GetFileMeta(ctx context.Context, objectName string) (*v1alpha1.FileMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileMeta", ctx, objectName)
	ret0, _ := ret[0].(*v1alpha1.FileMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileMeta indicates an expected call of GetFileMeta.
func (mr *MockFileMetaInterfaceMockRecorder) GetFileMeta(ctx, objectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileMeta", reflect.TypeOf((*MockFileMetaInterface)(nil).GetFileMeta), ctx, objectName)
}

// ListFileMetas mocks base method.
func (m *MockFileMetaInterface) ListFileMetas(ctx context.Context, opt *v1.ListOptions) ([]v1alpha1.FileMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFileMetas", ctx, opt)
	ret0, _ := ret[0].([]v1alpha1.FileMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFileMetas indicates an expected call of ListFileMetas.
func (mr *MockFileMetaInterfaceMockRecorder) ListFileMetas(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFileMetas", reflect.TypeOf((*MockFileMetaInterface)(nil).ListFileMetas), ctx, opt)
}
