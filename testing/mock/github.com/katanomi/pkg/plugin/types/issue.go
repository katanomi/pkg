// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: IssueLister,IssueGetter,IssueBranchLister,IssueBranchCreator,IssueBranchDeleter,IssueAttributeGetter)

// Package types is a generated GoMock package.
package types

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	zap "go.uber.org/zap"
)

// MockIssueLister is a mock of IssueLister interface.
type MockIssueLister struct {
	ctrl     *gomock.Controller
	recorder *MockIssueListerMockRecorder
}

// MockIssueListerMockRecorder is the mock recorder for MockIssueLister.
type MockIssueListerMockRecorder struct {
	mock *MockIssueLister
}

// NewMockIssueLister creates a new mock instance.
func NewMockIssueLister(ctrl *gomock.Controller) *MockIssueLister {
	mock := &MockIssueLister{ctrl: ctrl}
	mock.recorder = &MockIssueListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIssueLister) EXPECT() *MockIssueListerMockRecorder {
	return m.recorder
}

// ListIssues mocks base method.
func (m *MockIssueLister) ListIssues(arg0 context.Context, arg1 v1alpha1.IssueOptions, arg2 v1alpha1.ListOptions) (*v1alpha1.IssueList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIssues", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1alpha1.IssueList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIssues indicates an expected call of ListIssues.
func (mr *MockIssueListerMockRecorder) ListIssues(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIssues", reflect.TypeOf((*MockIssueLister)(nil).ListIssues), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockIssueLister) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockIssueListerMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockIssueLister)(nil).Path))
}

// Setup mocks base method.
func (m *MockIssueLister) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockIssueListerMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockIssueLister)(nil).Setup), arg0, arg1)
}

// MockIssueGetter is a mock of IssueGetter interface.
type MockIssueGetter struct {
	ctrl     *gomock.Controller
	recorder *MockIssueGetterMockRecorder
}

// MockIssueGetterMockRecorder is the mock recorder for MockIssueGetter.
type MockIssueGetterMockRecorder struct {
	mock *MockIssueGetter
}

// NewMockIssueGetter creates a new mock instance.
func NewMockIssueGetter(ctrl *gomock.Controller) *MockIssueGetter {
	mock := &MockIssueGetter{ctrl: ctrl}
	mock.recorder = &MockIssueGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIssueGetter) EXPECT() *MockIssueGetterMockRecorder {
	return m.recorder
}

// GetIssue mocks base method.
func (m *MockIssueGetter) GetIssue(arg0 context.Context, arg1 v1alpha1.IssueOptions, arg2 v1alpha1.ListOptions) (*v1alpha1.Issue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIssue", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1alpha1.Issue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIssue indicates an expected call of GetIssue.
func (mr *MockIssueGetterMockRecorder) GetIssue(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIssue", reflect.TypeOf((*MockIssueGetter)(nil).GetIssue), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockIssueGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockIssueGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockIssueGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockIssueGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockIssueGetterMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockIssueGetter)(nil).Setup), arg0, arg1)
}

// MockIssueBranchLister is a mock of IssueBranchLister interface.
type MockIssueBranchLister struct {
	ctrl     *gomock.Controller
	recorder *MockIssueBranchListerMockRecorder
}

// MockIssueBranchListerMockRecorder is the mock recorder for MockIssueBranchLister.
type MockIssueBranchListerMockRecorder struct {
	mock *MockIssueBranchLister
}

// NewMockIssueBranchLister creates a new mock instance.
func NewMockIssueBranchLister(ctrl *gomock.Controller) *MockIssueBranchLister {
	mock := &MockIssueBranchLister{ctrl: ctrl}
	mock.recorder = &MockIssueBranchListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIssueBranchLister) EXPECT() *MockIssueBranchListerMockRecorder {
	return m.recorder
}

// ListIssueBranches mocks base method.
func (m *MockIssueBranchLister) ListIssueBranches(arg0 context.Context, arg1 v1alpha1.IssueOptions, arg2 v1alpha1.ListOptions) (*v1alpha1.BranchList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIssueBranches", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1alpha1.BranchList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIssueBranches indicates an expected call of ListIssueBranches.
func (mr *MockIssueBranchListerMockRecorder) ListIssueBranches(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIssueBranches", reflect.TypeOf((*MockIssueBranchLister)(nil).ListIssueBranches), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockIssueBranchLister) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockIssueBranchListerMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockIssueBranchLister)(nil).Path))
}

// Setup mocks base method.
func (m *MockIssueBranchLister) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockIssueBranchListerMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockIssueBranchLister)(nil).Setup), arg0, arg1)
}

// MockIssueBranchCreator is a mock of IssueBranchCreator interface.
type MockIssueBranchCreator struct {
	ctrl     *gomock.Controller
	recorder *MockIssueBranchCreatorMockRecorder
}

// MockIssueBranchCreatorMockRecorder is the mock recorder for MockIssueBranchCreator.
type MockIssueBranchCreatorMockRecorder struct {
	mock *MockIssueBranchCreator
}

// NewMockIssueBranchCreator creates a new mock instance.
func NewMockIssueBranchCreator(ctrl *gomock.Controller) *MockIssueBranchCreator {
	mock := &MockIssueBranchCreator{ctrl: ctrl}
	mock.recorder = &MockIssueBranchCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIssueBranchCreator) EXPECT() *MockIssueBranchCreatorMockRecorder {
	return m.recorder
}

// CreateIssueBranch mocks base method.
func (m *MockIssueBranchCreator) CreateIssueBranch(arg0 context.Context, arg1 v1alpha1.IssueOptions, arg2 v1alpha1.Branch) (*v1alpha1.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIssueBranch", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1alpha1.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIssueBranch indicates an expected call of CreateIssueBranch.
func (mr *MockIssueBranchCreatorMockRecorder) CreateIssueBranch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIssueBranch", reflect.TypeOf((*MockIssueBranchCreator)(nil).CreateIssueBranch), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockIssueBranchCreator) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockIssueBranchCreatorMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockIssueBranchCreator)(nil).Path))
}

// Setup mocks base method.
func (m *MockIssueBranchCreator) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockIssueBranchCreatorMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockIssueBranchCreator)(nil).Setup), arg0, arg1)
}

// MockIssueBranchDeleter is a mock of IssueBranchDeleter interface.
type MockIssueBranchDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockIssueBranchDeleterMockRecorder
}

// MockIssueBranchDeleterMockRecorder is the mock recorder for MockIssueBranchDeleter.
type MockIssueBranchDeleterMockRecorder struct {
	mock *MockIssueBranchDeleter
}

// NewMockIssueBranchDeleter creates a new mock instance.
func NewMockIssueBranchDeleter(ctrl *gomock.Controller) *MockIssueBranchDeleter {
	mock := &MockIssueBranchDeleter{ctrl: ctrl}
	mock.recorder = &MockIssueBranchDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIssueBranchDeleter) EXPECT() *MockIssueBranchDeleterMockRecorder {
	return m.recorder
}

// DeleteIssueBranch mocks base method.
func (m *MockIssueBranchDeleter) DeleteIssueBranch(arg0 context.Context, arg1 v1alpha1.IssueOptions, arg2 v1alpha1.ListOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIssueBranch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIssueBranch indicates an expected call of DeleteIssueBranch.
func (mr *MockIssueBranchDeleterMockRecorder) DeleteIssueBranch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIssueBranch", reflect.TypeOf((*MockIssueBranchDeleter)(nil).DeleteIssueBranch), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockIssueBranchDeleter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockIssueBranchDeleterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockIssueBranchDeleter)(nil).Path))
}

// Setup mocks base method.
func (m *MockIssueBranchDeleter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockIssueBranchDeleterMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockIssueBranchDeleter)(nil).Setup), arg0, arg1)
}

// MockIssueAttributeGetter is a mock of IssueAttributeGetter interface.
type MockIssueAttributeGetter struct {
	ctrl     *gomock.Controller
	recorder *MockIssueAttributeGetterMockRecorder
}

// MockIssueAttributeGetterMockRecorder is the mock recorder for MockIssueAttributeGetter.
type MockIssueAttributeGetterMockRecorder struct {
	mock *MockIssueAttributeGetter
}

// NewMockIssueAttributeGetter creates a new mock instance.
func NewMockIssueAttributeGetter(ctrl *gomock.Controller) *MockIssueAttributeGetter {
	mock := &MockIssueAttributeGetter{ctrl: ctrl}
	mock.recorder = &MockIssueAttributeGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIssueAttributeGetter) EXPECT() *MockIssueAttributeGetterMockRecorder {
	return m.recorder
}

// GetIssueAttribute mocks base method.
func (m *MockIssueAttributeGetter) GetIssueAttribute(arg0 context.Context, arg1 v1alpha1.IssueOptions, arg2 v1alpha1.ListOptions) (*v1alpha1.Attribute, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIssueAttribute", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1alpha1.Attribute)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIssueAttribute indicates an expected call of GetIssueAttribute.
func (mr *MockIssueAttributeGetterMockRecorder) GetIssueAttribute(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIssueAttribute", reflect.TypeOf((*MockIssueAttributeGetter)(nil).GetIssueAttribute), arg0, arg1, arg2)
}

// Path mocks base method.
func (m *MockIssueAttributeGetter) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockIssueAttributeGetterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockIssueAttributeGetter)(nil).Path))
}

// Setup mocks base method.
func (m *MockIssueAttributeGetter) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockIssueAttributeGetterMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockIssueAttributeGetter)(nil).Setup), arg0, arg1)
}
