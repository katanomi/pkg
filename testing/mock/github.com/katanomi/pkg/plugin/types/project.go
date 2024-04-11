// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: ProjectLister,ProjectGetter,SubtypeProjectGetter,ProjectCreator,ProjectDeleter)
//
// Generated by this command:
//
//	mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/project.go github.com/katanomi/pkg/plugin/types ProjectLister,ProjectGetter,SubtypeProjectGetter,ProjectCreator,ProjectDeleter
//

// Package types is a generated GoMock package.
package types

import (
	context "context"
	reflect "reflect"

	v1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	types "github.com/katanomi/pkg/plugin/types"
	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockProjectLister is a mock of ProjectLister interface.
type MockProjectLister struct {
	ctrl     *gomock.Controller
	recorder *MockProjectListerMockRecorder
}

// MockProjectListerMockRecorder is the mock recorder for MockProjectLister.
type MockProjectListerMockRecorder struct {
	mock *MockProjectLister
}

// NewMockProjectLister creates a new mock instance.
func NewMockProjectLister(ctrl *gomock.Controller) *MockProjectLister {
	mock := &MockProjectLister{ctrl: ctrl}
	mock.recorder = &MockProjectListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectLister) EXPECT() *MockProjectListerMockRecorder {
	return m.recorder
}

// ListProjects mocks base method.
func (m *MockProjectLister) ListProjects(arg0 context.Context, arg1 v1alpha1.ListOptions) (*v1alpha1.ProjectList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProjects", arg0, arg1)
	ret0, _ := ret[0].(*v1alpha1.ProjectList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProjects indicates an expected call of ListProjects.
func (mr *MockProjectListerMockRecorder) ListProjects(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProjects", reflect.TypeOf((*MockProjectLister)(nil).ListProjects), arg0, arg1)
}

// Path mocks base method.
func (m *MockProjectLister) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockProjectListerMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockProjectLister)(nil).Path))
}

// Setup mocks base method.
func (m *MockProjectLister) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockProjectListerMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockProjectLister)(nil).Setup), arg0, arg1)
}

// MockProjectGetter is a mock of ProjectGetter interface.
type MockProjectGetter struct {
	ctrl     *gomock.Controller
	recorder *MockProjectGetterMockRecorder
}

// MockProjectGetterMockRecorder is the mock recorder for MockProjectGetter.
type MockProjectGetterMockRecorder struct {
	mock *MockProjectGetter
}

// NewMockProjectGetter creates a new mock instance.
func NewMockProjectGetter(ctrl *gomock.Controller) *MockProjectGetter {
	mock := &MockProjectGetter{ctrl: ctrl}
	mock.recorder = &MockProjectGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectGetter) EXPECT() *MockProjectGetterMockRecorder {
	return m.recorder
}

// GetProject mocks base method.
func (m *MockProjectGetter) GetProject(arg0 context.Context, arg1 string) (*v1alpha1.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", arg0, arg1)
	ret0, _ := ret[0].(*v1alpha1.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject.
func (mr *MockProjectGetterMockRecorder) GetProject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockProjectGetter)(nil).GetProject), arg0, arg1)
}

// Path mocks base method.
func (m *MockProjectGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockProjectGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockProjectGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockProjectGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockProjectGetterMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockProjectGetter)(nil).Setup), arg0, arg1)
}

// MockSubtypeProjectGetter is a mock of SubtypeProjectGetter interface.
type MockSubtypeProjectGetter struct {
	ctrl     *gomock.Controller
	recorder *MockSubtypeProjectGetterMockRecorder
}

// MockSubtypeProjectGetterMockRecorder is the mock recorder for MockSubtypeProjectGetter.
type MockSubtypeProjectGetterMockRecorder struct {
	mock *MockSubtypeProjectGetter
}

// NewMockSubtypeProjectGetter creates a new mock instance.
func NewMockSubtypeProjectGetter(ctrl *gomock.Controller) *MockSubtypeProjectGetter {
	mock := &MockSubtypeProjectGetter{ctrl: ctrl}
	mock.recorder = &MockSubtypeProjectGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubtypeProjectGetter) EXPECT() *MockSubtypeProjectGetterMockRecorder {
	return m.recorder
}

// GetSubTypeProject mocks base method.
func (m *MockSubtypeProjectGetter) GetSubTypeProject(arg0 context.Context, arg1 types.GetProjectOption) (*v1alpha1.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubTypeProject", arg0, arg1)
	ret0, _ := ret[0].(*v1alpha1.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubTypeProject indicates an expected call of GetSubTypeProject.
func (mr *MockSubtypeProjectGetterMockRecorder) GetSubTypeProject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubTypeProject", reflect.TypeOf((*MockSubtypeProjectGetter)(nil).GetSubTypeProject), arg0, arg1)
}

// Path mocks base method.
func (m *MockSubtypeProjectGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockSubtypeProjectGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockSubtypeProjectGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockSubtypeProjectGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockSubtypeProjectGetterMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockSubtypeProjectGetter)(nil).Setup), arg0, arg1)
}

// MockProjectCreator is a mock of ProjectCreator interface.
type MockProjectCreator struct {
	ctrl     *gomock.Controller
	recorder *MockProjectCreatorMockRecorder
}

// MockProjectCreatorMockRecorder is the mock recorder for MockProjectCreator.
type MockProjectCreatorMockRecorder struct {
	mock *MockProjectCreator
}

// NewMockProjectCreator creates a new mock instance.
func NewMockProjectCreator(ctrl *gomock.Controller) *MockProjectCreator {
	mock := &MockProjectCreator{ctrl: ctrl}
	mock.recorder = &MockProjectCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectCreator) EXPECT() *MockProjectCreatorMockRecorder {
	return m.recorder
}

// CreateProject mocks base method.
func (m *MockProjectCreator) CreateProject(arg0 context.Context, arg1 *v1alpha1.Project) (*v1alpha1.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", arg0, arg1)
	ret0, _ := ret[0].(*v1alpha1.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockProjectCreatorMockRecorder) CreateProject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectCreator)(nil).CreateProject), arg0, arg1)
}

// Path mocks base method.
func (m *MockProjectCreator) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockProjectCreatorMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockProjectCreator)(nil).Path))
}

// Setup mocks base method.
func (m *MockProjectCreator) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockProjectCreatorMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockProjectCreator)(nil).Setup), arg0, arg1)
}

// MockProjectDeleter is a mock of ProjectDeleter interface.
type MockProjectDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockProjectDeleterMockRecorder
}

// MockProjectDeleterMockRecorder is the mock recorder for MockProjectDeleter.
type MockProjectDeleterMockRecorder struct {
	mock *MockProjectDeleter
}

// NewMockProjectDeleter creates a new mock instance.
func NewMockProjectDeleter(ctrl *gomock.Controller) *MockProjectDeleter {
	mock := &MockProjectDeleter{ctrl: ctrl}
	mock.recorder = &MockProjectDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectDeleter) EXPECT() *MockProjectDeleterMockRecorder {
	return m.recorder
}

// DeleteProject mocks base method.
func (m *MockProjectDeleter) DeleteProject(arg0 context.Context, arg1 *v1alpha1.Project) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject.
func (mr *MockProjectDeleterMockRecorder) DeleteProject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockProjectDeleter)(nil).DeleteProject), arg0, arg1)
}

// Path mocks base method.
func (m *MockProjectDeleter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockProjectDeleterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockProjectDeleter)(nil).Path))
}

// Setup mocks base method.
func (m *MockProjectDeleter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockProjectDeleterMockRecorder) Setup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockProjectDeleter)(nil).Setup), arg0, arg1)
}
