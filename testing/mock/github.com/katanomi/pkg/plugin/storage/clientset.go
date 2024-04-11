// Code generated by MockGen. DO NOT EDIT.
// Source: clientset.go
//
// Generated by this command:
//
//	mockgen -source=clientset.go -destination=../../testing/mock/github.com/katanomi/pkg/plugin/storage/clientset.go -package=storage Interface
//

// Package storage is a generated GoMock package.
package storage

import (
	reflect "reflect"

	v1alpha1 "github.com/katanomi/pkg/plugin/storage/client/versioned/archive/v1alpha1"
	v1alpha10 "github.com/katanomi/pkg/plugin/storage/client/versioned/core/v1alpha1"
	v1alpha11 "github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1"
	gomock "go.uber.org/mock/gomock"
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

// ArchiveV1alpha1 mocks base method.
func (m *MockInterface) ArchiveV1alpha1() v1alpha1.ArchiveInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveV1alpha1")
	ret0, _ := ret[0].(v1alpha1.ArchiveInterface)
	return ret0
}

// ArchiveV1alpha1 indicates an expected call of ArchiveV1alpha1.
func (mr *MockInterfaceMockRecorder) ArchiveV1alpha1() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveV1alpha1", reflect.TypeOf((*MockInterface)(nil).ArchiveV1alpha1))
}

// CoreV1alpha1 mocks base method.
func (m *MockInterface) CoreV1alpha1() v1alpha10.CoreV1alpha1Interface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CoreV1alpha1")
	ret0, _ := ret[0].(v1alpha10.CoreV1alpha1Interface)
	return ret0
}

// CoreV1alpha1 indicates an expected call of CoreV1alpha1.
func (mr *MockInterfaceMockRecorder) CoreV1alpha1() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CoreV1alpha1", reflect.TypeOf((*MockInterface)(nil).CoreV1alpha1))
}

// FileStoreV1alpha1 mocks base method.
func (m *MockInterface) FileStoreV1alpha1() v1alpha11.FileStoreV1alpha1Interface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileStoreV1alpha1")
	ret0, _ := ret[0].(v1alpha11.FileStoreV1alpha1Interface)
	return ret0
}

// FileStoreV1alpha1 indicates an expected call of FileStoreV1alpha1.
func (mr *MockInterfaceMockRecorder) FileStoreV1alpha1() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileStoreV1alpha1", reflect.TypeOf((*MockInterface)(nil).FileStoreV1alpha1))
}
