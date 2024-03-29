// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: GitCommitGetter,GitCommitCreator,GitCommitLister)

// Package types is a generated GoMock package.
package types

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/katanomi/pkg/apis/coderepository/v1alpha1"
	v1alpha10 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	zap "go.uber.org/zap"
)

// MockGitCommitGetter is a mock of GitCommitGetter interface.
type MockGitCommitGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGitCommitGetterMockRecorder
}

// MockGitCommitGetterMockRecorder is the mock recorder for MockGitCommitGetter.
type MockGitCommitGetterMockRecorder struct {
	mock *MockGitCommitGetter
}

// NewMockGitCommitGetter creates a new mock instance.
func NewMockGitCommitGetter(ctrl *gomock.Controller) *MockGitCommitGetter {
	mock := &MockGitCommitGetter{ctrl: ctrl}
	mock.recorder = &MockGitCommitGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitCommitGetter) EXPECT() *MockGitCommitGetterMockRecorder {
	return m.recorder
}

// GetGitCommit mocks base method.
func (m *MockGitCommitGetter) GetGitCommit(arg0 context.Context, arg1 v1alpha10.GitCommitOption) (v1alpha10.GitCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitCommit", arg0, arg1)
	ret0, _ := ret[0].(v1alpha10.GitCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitCommit indicates an expected call of GetGitCommit.
func (mr *MockGitCommitGetterMockRecorder) GetGitCommit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitCommit", reflect.TypeOf((*MockGitCommitGetter)(nil).GetGitCommit), arg0, arg1)
}

// Path mocks base method.
func (m *MockGitCommitGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitCommitGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitCommitGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitCommitGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitCommitGetterMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitCommitGetter)(nil).Setup), arg0, arg1)
}

// MockGitCommitCreator is a mock of GitCommitCreator interface.
type MockGitCommitCreator struct {
	ctrl     *gomock.Controller
	recorder *MockGitCommitCreatorMockRecorder
}

// MockGitCommitCreatorMockRecorder is the mock recorder for MockGitCommitCreator.
type MockGitCommitCreatorMockRecorder struct {
	mock *MockGitCommitCreator
}

// NewMockGitCommitCreator creates a new mock instance.
func NewMockGitCommitCreator(ctrl *gomock.Controller) *MockGitCommitCreator {
	mock := &MockGitCommitCreator{ctrl: ctrl}
	mock.recorder = &MockGitCommitCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitCommitCreator) EXPECT() *MockGitCommitCreatorMockRecorder {
	return m.recorder
}

// CreateGitCommit mocks base method.
func (m *MockGitCommitCreator) CreateGitCommit(arg0 context.Context, arg1 v1alpha1.CreateGitCommitOption) (v1alpha10.GitCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGitCommit", arg0, arg1)
	ret0, _ := ret[0].(v1alpha10.GitCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGitCommit indicates an expected call of CreateGitCommit.
func (mr *MockGitCommitCreatorMockRecorder) CreateGitCommit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGitCommit", reflect.TypeOf((*MockGitCommitCreator)(nil).CreateGitCommit), arg0, arg1)
}

// Path mocks base method.
func (m *MockGitCommitCreator) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitCommitCreatorMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitCommitCreator)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitCommitCreator) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitCommitCreatorMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitCommitCreator)(nil).Setup), arg0, arg1)
}

// MockGitCommitLister is a mock of GitCommitLister interface.
type MockGitCommitLister struct {
	ctrl     *gomock.Controller
	recorder *MockGitCommitListerMockRecorder
}

// MockGitCommitListerMockRecorder is the mock recorder for MockGitCommitLister.
type MockGitCommitListerMockRecorder struct {
	mock *MockGitCommitLister
}

// NewMockGitCommitLister creates a new mock instance.
func NewMockGitCommitLister(ctrl *gomock.Controller) *MockGitCommitLister {
	mock := &MockGitCommitLister{ctrl: ctrl}
	mock.recorder = &MockGitCommitListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitCommitLister) EXPECT() *MockGitCommitListerMockRecorder {
	return m.recorder
}

// ListGitCommit mocks base method.
func (m *MockGitCommitLister) ListGitCommit(arg0 context.Context, arg1 v1alpha10.GitCommitListOption, arg2 v1alpha10.ListOptions) (v1alpha10.GitCommitList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGitCommit", arg0, arg1, arg2)
	ret0, _ := ret[0].(v1alpha10.GitCommitList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGitCommit indicates an expected call of ListGitCommit.
func (mr *MockGitCommitListerMockRecorder) ListGitCommit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGitCommit", reflect.TypeOf((*MockGitCommitLister)(nil).ListGitCommit), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockGitCommitLister) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitCommitListerMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitCommitLister)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitCommitLister) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitCommitListerMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitCommitLister)(nil).Setup), arg0, arg1)
}
