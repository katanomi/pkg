// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: WebhookRegister,WebhookCreator,WebhookUpdater,WebhookDeleter,WebhookLister,WebhookResourceDiffer,WebhookReceiver)

// Package types is a generated GoMock package.
package types

import (
	context "context"
	reflect "reflect"

	event "github.com/cloudevents/sdk-go/v2/event"
	restful "github.com/emicklei/go-restful/v3"
	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	zap "go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	apis "knative.dev/pkg/apis"
)

// MockWebhookRegister is a mock of WebhookRegister interface.
type MockWebhookRegister struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookRegisterMockRecorder
}

// MockWebhookRegisterMockRecorder is the mock recorder for MockWebhookRegister.
type MockWebhookRegisterMockRecorder struct {
	mock *MockWebhookRegister
}

// NewMockWebhookRegister creates a new mock instance.
func NewMockWebhookRegister(ctrl *gomock.Controller) *MockWebhookRegister {
	mock := &MockWebhookRegister{ctrl: ctrl}
	mock.recorder = &MockWebhookRegisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebhookRegister) EXPECT() *MockWebhookRegisterMockRecorder {
	return m.recorder
}

// CreateWebhook mocks base method.
func (m *MockWebhookRegister) CreateWebhook(arg0 context.Context, arg1 v1alpha1.WebhookRegisterSpec, arg2 v1.Secret) (v1alpha1.WebhookRegisterStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(v1alpha1.WebhookRegisterStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebhook indicates an expected call of CreateWebhook.
func (mr *MockWebhookRegisterMockRecorder) CreateWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebhook", reflect.TypeOf((*MockWebhookRegister)(nil).CreateWebhook), arg0, arg1, arg2)
}

// DeleteWebhook mocks base method.
func (m *MockWebhookRegister) DeleteWebhook(arg0 context.Context, arg1 v1alpha1.WebhookRegisterSpec, arg2 v1.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebhook indicates an expected call of DeleteWebhook.
func (mr *MockWebhookRegisterMockRecorder) DeleteWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebhook", reflect.TypeOf((*MockWebhookRegister)(nil).DeleteWebhook), arg0, arg1, arg2)
}

// ListWebhooks mocks base method.
func (m *MockWebhookRegister) ListWebhooks(arg0 context.Context, arg1 apis.URL, arg2 v1.Secret) ([]v1alpha1.WebhookRegisterStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWebhooks", arg0, arg1, arg2)
	ret0, _ := ret[0].([]v1alpha1.WebhookRegisterStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWebhooks indicates an expected call of ListWebhooks.
func (mr *MockWebhookRegisterMockRecorder) ListWebhooks(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWebhooks", reflect.TypeOf((*MockWebhookRegister)(nil).ListWebhooks), arg0, arg1, arg2)
}

// UpdateWebhook mocks base method.
func (m *MockWebhookRegister) UpdateWebhook(arg0 context.Context, arg1 v1alpha1.WebhookRegisterSpec, arg2 v1.Secret) (v1alpha1.WebhookRegisterStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(v1alpha1.WebhookRegisterStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWebhook indicates an expected call of UpdateWebhook.
func (mr *MockWebhookRegisterMockRecorder) UpdateWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWebhook", reflect.TypeOf((*MockWebhookRegister)(nil).UpdateWebhook), arg0, arg1, arg2)
}

// MockWebhookCreator is a mock of WebhookCreator interface.
type MockWebhookCreator struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookCreatorMockRecorder
}

// MockWebhookCreatorMockRecorder is the mock recorder for MockWebhookCreator.
type MockWebhookCreatorMockRecorder struct {
	mock *MockWebhookCreator
}

// NewMockWebhookCreator creates a new mock instance.
func NewMockWebhookCreator(ctrl *gomock.Controller) *MockWebhookCreator {
	mock := &MockWebhookCreator{ctrl: ctrl}
	mock.recorder = &MockWebhookCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebhookCreator) EXPECT() *MockWebhookCreatorMockRecorder {
	return m.recorder
}

// CreateWebhook mocks base method.
func (m *MockWebhookCreator) CreateWebhook(arg0 context.Context, arg1 v1alpha1.WebhookRegisterSpec, arg2 v1.Secret) (v1alpha1.WebhookRegisterStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(v1alpha1.WebhookRegisterStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebhook indicates an expected call of CreateWebhook.
func (mr *MockWebhookCreatorMockRecorder) CreateWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebhook", reflect.TypeOf((*MockWebhookCreator)(nil).CreateWebhook), arg0, arg1, arg2)
}

// MockWebhookUpdater is a mock of WebhookUpdater interface.
type MockWebhookUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookUpdaterMockRecorder
}

// MockWebhookUpdaterMockRecorder is the mock recorder for MockWebhookUpdater.
type MockWebhookUpdaterMockRecorder struct {
	mock *MockWebhookUpdater
}

// NewMockWebhookUpdater creates a new mock instance.
func NewMockWebhookUpdater(ctrl *gomock.Controller) *MockWebhookUpdater {
	mock := &MockWebhookUpdater{ctrl: ctrl}
	mock.recorder = &MockWebhookUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebhookUpdater) EXPECT() *MockWebhookUpdaterMockRecorder {
	return m.recorder
}

// UpdateWebhook mocks base method.
func (m *MockWebhookUpdater) UpdateWebhook(arg0 context.Context, arg1 v1alpha1.WebhookRegisterSpec, arg2 v1.Secret) (v1alpha1.WebhookRegisterStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(v1alpha1.WebhookRegisterStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWebhook indicates an expected call of UpdateWebhook.
func (mr *MockWebhookUpdaterMockRecorder) UpdateWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWebhook", reflect.TypeOf((*MockWebhookUpdater)(nil).UpdateWebhook), arg0, arg1, arg2)
}

// MockWebhookDeleter is a mock of WebhookDeleter interface.
type MockWebhookDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookDeleterMockRecorder
}

// MockWebhookDeleterMockRecorder is the mock recorder for MockWebhookDeleter.
type MockWebhookDeleterMockRecorder struct {
	mock *MockWebhookDeleter
}

// NewMockWebhookDeleter creates a new mock instance.
func NewMockWebhookDeleter(ctrl *gomock.Controller) *MockWebhookDeleter {
	mock := &MockWebhookDeleter{ctrl: ctrl}
	mock.recorder = &MockWebhookDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebhookDeleter) EXPECT() *MockWebhookDeleterMockRecorder {
	return m.recorder
}

// DeleteWebhook mocks base method.
func (m *MockWebhookDeleter) DeleteWebhook(arg0 context.Context, arg1 v1alpha1.WebhookRegisterSpec, arg2 v1.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebhook indicates an expected call of DeleteWebhook.
func (mr *MockWebhookDeleterMockRecorder) DeleteWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebhook", reflect.TypeOf((*MockWebhookDeleter)(nil).DeleteWebhook), arg0, arg1, arg2)
}

// MockWebhookLister is a mock of WebhookLister interface.
type MockWebhookLister struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookListerMockRecorder
}

// MockWebhookListerMockRecorder is the mock recorder for MockWebhookLister.
type MockWebhookListerMockRecorder struct {
	mock *MockWebhookLister
}

// NewMockWebhookLister creates a new mock instance.
func NewMockWebhookLister(ctrl *gomock.Controller) *MockWebhookLister {
	mock := &MockWebhookLister{ctrl: ctrl}
	mock.recorder = &MockWebhookListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebhookLister) EXPECT() *MockWebhookListerMockRecorder {
	return m.recorder
}

// ListWebhooks mocks base method.
func (m *MockWebhookLister) ListWebhooks(arg0 context.Context, arg1 apis.URL, arg2 v1.Secret) ([]v1alpha1.WebhookRegisterStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWebhooks", arg0, arg1, arg2)
	ret0, _ := ret[0].([]v1alpha1.WebhookRegisterStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWebhooks indicates an expected call of ListWebhooks.
func (mr *MockWebhookListerMockRecorder) ListWebhooks(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWebhooks", reflect.TypeOf((*MockWebhookLister)(nil).ListWebhooks), arg0, arg1, arg2)
}

// MockWebhookResourceDiffer is a mock of WebhookResourceDiffer interface.
type MockWebhookResourceDiffer struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookResourceDifferMockRecorder
}

// MockWebhookResourceDifferMockRecorder is the mock recorder for MockWebhookResourceDiffer.
type MockWebhookResourceDifferMockRecorder struct {
	mock *MockWebhookResourceDiffer
}

// NewMockWebhookResourceDiffer creates a new mock instance.
func NewMockWebhookResourceDiffer(ctrl *gomock.Controller) *MockWebhookResourceDiffer {
	mock := &MockWebhookResourceDiffer{ctrl: ctrl}
	mock.recorder = &MockWebhookResourceDifferMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebhookResourceDiffer) EXPECT() *MockWebhookResourceDifferMockRecorder {
	return m.recorder
}

// IsSameResource mocks base method.
func (m *MockWebhookResourceDiffer) IsSameResource(arg0 context.Context, arg1, arg2 v1alpha1.ResourceURI) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSameResource", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSameResource indicates an expected call of IsSameResource.
func (mr *MockWebhookResourceDifferMockRecorder) IsSameResource(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSameResource", reflect.TypeOf((*MockWebhookResourceDiffer)(nil).IsSameResource), arg0, arg1, arg2)
}

// MockWebhookReceiver is a mock of WebhookReceiver interface.
type MockWebhookReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookReceiverMockRecorder
}

// MockWebhookReceiverMockRecorder is the mock recorder for MockWebhookReceiver.
type MockWebhookReceiverMockRecorder struct {
	mock *MockWebhookReceiver
}

// NewMockWebhookReceiver creates a new mock instance.
func NewMockWebhookReceiver(ctrl *gomock.Controller) *MockWebhookReceiver {
	mock := &MockWebhookReceiver{ctrl: ctrl}
	mock.recorder = &MockWebhookReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebhookReceiver) EXPECT() *MockWebhookReceiverMockRecorder {
	return m.recorder
}

// Path mocks base method.
func (m *MockWebhookReceiver) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockWebhookReceiverMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockWebhookReceiver)(nil).Path))
}

// ReceiveWebhook mocks base method.
func (m *MockWebhookReceiver) ReceiveWebhook(arg0 context.Context, arg1 *restful.Request, arg2 string) (event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReceiveWebhook indicates an expected call of ReceiveWebhook.
func (mr *MockWebhookReceiverMockRecorder) ReceiveWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveWebhook", reflect.TypeOf((*MockWebhookReceiver)(nil).ReceiveWebhook), arg0, arg1, arg2)
}

// Setup mocks base method.
func (m *MockWebhookReceiver) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockWebhookReceiverMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockWebhookReceiver)(nil).Setup), arg0, arg1)
}
