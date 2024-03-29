// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: GitRepoFileGetter,GitRepoFileCreator,GitRepositoryFileTreeGetter)

// Package types is a generated GoMock package.
package types

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	zap "go.uber.org/zap"
)

// MockGitRepoFileGetter is a mock of GitRepoFileGetter interface.
type MockGitRepoFileGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepoFileGetterMockRecorder
}

// MockGitRepoFileGetterMockRecorder is the mock recorder for MockGitRepoFileGetter.
type MockGitRepoFileGetterMockRecorder struct {
	mock *MockGitRepoFileGetter
}

// NewMockGitRepoFileGetter creates a new mock instance.
func NewMockGitRepoFileGetter(ctrl *gomock.Controller) *MockGitRepoFileGetter {
	mock := &MockGitRepoFileGetter{ctrl: ctrl}
	mock.recorder = &MockGitRepoFileGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepoFileGetter) EXPECT() *MockGitRepoFileGetterMockRecorder {
	return m.recorder
}

// GetGitRepoFile mocks base method.
func (m *MockGitRepoFileGetter) GetGitRepoFile(arg0 context.Context, arg1 v1alpha1.GitRepoFileOption) (v1alpha1.GitRepoFile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitRepoFile", arg0, arg1)
	ret0, _ := ret[0].(v1alpha1.GitRepoFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitRepoFile indicates an expected call of GetGitRepoFile.
func (mr *MockGitRepoFileGetterMockRecorder) GetGitRepoFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitRepoFile", reflect.TypeOf((*MockGitRepoFileGetter)(nil).GetGitRepoFile), arg0, arg1)
}

// Path mocks base method.
func (m *MockGitRepoFileGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitRepoFileGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitRepoFileGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitRepoFileGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitRepoFileGetterMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitRepoFileGetter)(nil).Setup), arg0, arg1)
}

// MockGitRepoFileCreator is a mock of GitRepoFileCreator interface.
type MockGitRepoFileCreator struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepoFileCreatorMockRecorder
}

// MockGitRepoFileCreatorMockRecorder is the mock recorder for MockGitRepoFileCreator.
type MockGitRepoFileCreatorMockRecorder struct {
	mock *MockGitRepoFileCreator
}

// NewMockGitRepoFileCreator creates a new mock instance.
func NewMockGitRepoFileCreator(ctrl *gomock.Controller) *MockGitRepoFileCreator {
	mock := &MockGitRepoFileCreator{ctrl: ctrl}
	mock.recorder = &MockGitRepoFileCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepoFileCreator) EXPECT() *MockGitRepoFileCreatorMockRecorder {
	return m.recorder
}

// CreateGitRepoFile mocks base method.
func (m *MockGitRepoFileCreator) CreateGitRepoFile(arg0 context.Context, arg1 v1alpha1.CreateRepoFilePayload) (v1alpha1.GitCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGitRepoFile", arg0, arg1)
	ret0, _ := ret[0].(v1alpha1.GitCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGitRepoFile indicates an expected call of CreateGitRepoFile.
func (mr *MockGitRepoFileCreatorMockRecorder) CreateGitRepoFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGitRepoFile", reflect.TypeOf((*MockGitRepoFileCreator)(nil).CreateGitRepoFile), arg0, arg1)
}

// Path mocks base method.
func (m *MockGitRepoFileCreator) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitRepoFileCreatorMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitRepoFileCreator)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitRepoFileCreator) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitRepoFileCreatorMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitRepoFileCreator)(nil).Setup), arg0, arg1)
}

// MockGitRepositoryFileTreeGetter is a mock of GitRepositoryFileTreeGetter interface.
type MockGitRepositoryFileTreeGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepositoryFileTreeGetterMockRecorder
}

// MockGitRepositoryFileTreeGetterMockRecorder is the mock recorder for MockGitRepositoryFileTreeGetter.
type MockGitRepositoryFileTreeGetterMockRecorder struct {
	mock *MockGitRepositoryFileTreeGetter
}

// NewMockGitRepositoryFileTreeGetter creates a new mock instance.
func NewMockGitRepositoryFileTreeGetter(ctrl *gomock.Controller) *MockGitRepositoryFileTreeGetter {
	mock := &MockGitRepositoryFileTreeGetter{ctrl: ctrl}
	mock.recorder = &MockGitRepositoryFileTreeGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepositoryFileTreeGetter) EXPECT() *MockGitRepositoryFileTreeGetterMockRecorder {
	return m.recorder
}

// GetGitRepositoryFileTree mocks base method.
func (m *MockGitRepositoryFileTreeGetter) GetGitRepositoryFileTree(arg0 context.Context, arg1 v1alpha1.GitRepoFileTreeOption, arg2 v1alpha1.ListOptions) (v1alpha1.GitRepositoryFileTree, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitRepositoryFileTree", arg0, arg1, arg2)
	ret0, _ := ret[0].(v1alpha1.GitRepositoryFileTree)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitRepositoryFileTree indicates an expected call of GetGitRepositoryFileTree.
func (mr *MockGitRepositoryFileTreeGetterMockRecorder) GetGitRepositoryFileTree(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitRepositoryFileTree", reflect.TypeOf((*MockGitRepositoryFileTreeGetter)(nil).GetGitRepositoryFileTree), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockGitRepositoryFileTreeGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitRepositoryFileTreeGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitRepositoryFileTreeGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitRepositoryFileTreeGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitRepositoryFileTreeGetterMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitRepositoryFileTreeGetter)(nil).Setup), arg0, arg1)
}
