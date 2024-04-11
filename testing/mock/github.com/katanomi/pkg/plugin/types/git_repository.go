// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: GitRepositoryCreator,GitRepositoryDeleter,GitRepositoryLister,GitRepositoryGetter)
//
// Generated by this command:
//
//	mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/git_repository.go github.com/katanomi/pkg/plugin/types GitRepositoryCreator,GitRepositoryDeleter,GitRepositoryLister,GitRepositoryGetter
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

// MockGitRepositoryCreator is a mock of GitRepositoryCreator interface.
type MockGitRepositoryCreator struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepositoryCreatorMockRecorder
}

// MockGitRepositoryCreatorMockRecorder is the mock recorder for MockGitRepositoryCreator.
type MockGitRepositoryCreatorMockRecorder struct {
	mock *MockGitRepositoryCreator
}

// NewMockGitRepositoryCreator creates a new mock instance.
func NewMockGitRepositoryCreator(ctrl *gomock.Controller) *MockGitRepositoryCreator {
	mock := &MockGitRepositoryCreator{ctrl: ctrl}
	mock.recorder = &MockGitRepositoryCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepositoryCreator) EXPECT() *MockGitRepositoryCreatorMockRecorder {
	return m.recorder
}

// CreateGitRepository mocks base method.
func (m *MockGitRepositoryCreator) CreateGitRepository(arg0 context.Context, arg1 v1alpha1.CreateGitRepositoryPayload) (v1alpha1.GitRepository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGitRepository", arg0, arg1)
	ret0, _ := ret[0].(v1alpha1.GitRepository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGitRepository indicates an expected call of CreateGitRepository.
func (mr *MockGitRepositoryCreatorMockRecorder) CreateGitRepository(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGitRepository", reflect.TypeOf((*MockGitRepositoryCreator)(nil).CreateGitRepository), arg0, arg1)
}

// Path mocks base method.
func (m *MockGitRepositoryCreator) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitRepositoryCreatorMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitRepositoryCreator)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitRepositoryCreator) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitRepositoryCreatorMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitRepositoryCreator)(nil).Setup), arg0, arg1)
}

// MockGitRepositoryDeleter is a mock of GitRepositoryDeleter interface.
type MockGitRepositoryDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepositoryDeleterMockRecorder
}

// MockGitRepositoryDeleterMockRecorder is the mock recorder for MockGitRepositoryDeleter.
type MockGitRepositoryDeleterMockRecorder struct {
	mock *MockGitRepositoryDeleter
}

// NewMockGitRepositoryDeleter creates a new mock instance.
func NewMockGitRepositoryDeleter(ctrl *gomock.Controller) *MockGitRepositoryDeleter {
	mock := &MockGitRepositoryDeleter{ctrl: ctrl}
	mock.recorder = &MockGitRepositoryDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepositoryDeleter) EXPECT() *MockGitRepositoryDeleterMockRecorder {
	return m.recorder
}

// DeleteGitRepository mocks base method.
func (m *MockGitRepositoryDeleter) DeleteGitRepository(arg0 context.Context, arg1 v1alpha1.GitRepo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGitRepository", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGitRepository indicates an expected call of DeleteGitRepository.
func (mr *MockGitRepositoryDeleterMockRecorder) DeleteGitRepository(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGitRepository", reflect.TypeOf((*MockGitRepositoryDeleter)(nil).DeleteGitRepository), arg0, arg1)
}

// Path mocks base method.
func (m *MockGitRepositoryDeleter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitRepositoryDeleterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitRepositoryDeleter)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitRepositoryDeleter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitRepositoryDeleterMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitRepositoryDeleter)(nil).Setup), arg0, arg1)
}

// MockGitRepositoryLister is a mock of GitRepositoryLister interface.
type MockGitRepositoryLister struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepositoryListerMockRecorder
}

// MockGitRepositoryListerMockRecorder is the mock recorder for MockGitRepositoryLister.
type MockGitRepositoryListerMockRecorder struct {
	mock *MockGitRepositoryLister
}

// NewMockGitRepositoryLister creates a new mock instance.
func NewMockGitRepositoryLister(ctrl *gomock.Controller) *MockGitRepositoryLister {
	mock := &MockGitRepositoryLister{ctrl: ctrl}
	mock.recorder = &MockGitRepositoryListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepositoryLister) EXPECT() *MockGitRepositoryListerMockRecorder {
	return m.recorder
}

// ListGitRepository mocks base method.
func (m *MockGitRepositoryLister) ListGitRepository(arg0 context.Context, arg1, arg2 string, arg3 v1alpha1.ProjectSubType, arg4 v1alpha1.ListOptions) (v1alpha1.GitRepositoryList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGitRepository", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(v1alpha1.GitRepositoryList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGitRepository indicates an expected call of ListGitRepository.
func (mr *MockGitRepositoryListerMockRecorder) ListGitRepository(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGitRepository", reflect.TypeOf((*MockGitRepositoryLister)(nil).ListGitRepository), arg0, arg1, arg2, arg3, arg4)
}

// Path mocks base method.
func (m *MockGitRepositoryLister) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitRepositoryListerMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitRepositoryLister)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitRepositoryLister) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitRepositoryListerMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitRepositoryLister)(nil).Setup), arg0, arg1)
}

// MockGitRepositoryGetter is a mock of GitRepositoryGetter interface.
type MockGitRepositoryGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepositoryGetterMockRecorder
}

// MockGitRepositoryGetterMockRecorder is the mock recorder for MockGitRepositoryGetter.
type MockGitRepositoryGetterMockRecorder struct {
	mock *MockGitRepositoryGetter
}

// NewMockGitRepositoryGetter creates a new mock instance.
func NewMockGitRepositoryGetter(ctrl *gomock.Controller) *MockGitRepositoryGetter {
	mock := &MockGitRepositoryGetter{ctrl: ctrl}
	mock.recorder = &MockGitRepositoryGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepositoryGetter) EXPECT() *MockGitRepositoryGetterMockRecorder {
	return m.recorder
}

// GetGitRepository mocks base method.
func (m *MockGitRepositoryGetter) GetGitRepository(arg0 context.Context, arg1 v1alpha1.GitRepo) (v1alpha1.GitRepository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitRepository", arg0, arg1)
	ret0, _ := ret[0].(v1alpha1.GitRepository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitRepository indicates an expected call of GetGitRepository.
func (mr *MockGitRepositoryGetterMockRecorder) GetGitRepository(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitRepository", reflect.TypeOf((*MockGitRepositoryGetter)(nil).GetGitRepository), arg0, arg1)
}

// Path mocks base method.
func (m *MockGitRepositoryGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockGitRepositoryGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockGitRepositoryGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockGitRepositoryGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockGitRepositoryGetterMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockGitRepositoryGetter)(nil).Setup), arg0, arg1)
}
