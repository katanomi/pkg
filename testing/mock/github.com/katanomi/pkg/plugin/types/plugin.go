// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/plugin/types (interfaces: Interface,PluginRegister,PluginAddressable,DependentResourceGetter,AdditionalWebhookRegister,ResourcePathFormatter,PluginDisplayColumns,PluginAttributes,PluginVersionAttributes,LivenessChecker,Initializer,ToolMetadataGetter)

// Package types is a generated GoMock package.
package types

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	zap "go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	apis "knative.dev/pkg/apis"
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

// Path mocks base method.
func (m *MockInterface) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockInterfaceMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockInterface)(nil).Path))
}

// Setup mocks base method.
func (m *MockInterface) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockInterfaceMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockInterface)(nil).Setup), arg0, arg1)
}

// MockPluginRegister is a mock of PluginRegister interface.
type MockPluginRegister struct {
	ctrl     *gomock.Controller
	recorder *MockPluginRegisterMockRecorder
}

// MockPluginRegisterMockRecorder is the mock recorder for MockPluginRegister.
type MockPluginRegisterMockRecorder struct {
	mock *MockPluginRegister
}

// NewMockPluginRegister creates a new mock instance.
func NewMockPluginRegister(ctrl *gomock.Controller) *MockPluginRegister {
	mock := &MockPluginRegister{ctrl: ctrl}
	mock.recorder = &MockPluginRegisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPluginRegister) EXPECT() *MockPluginRegisterMockRecorder {
	return m.recorder
}

// GetAddressURL mocks base method.
func (m *MockPluginRegister) GetAddressURL() *apis.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddressURL")
	ret0, _ := ret[0].(*apis.URL)
	return ret0
}

// GetAddressURL indicates an expected call of GetAddressURL.
func (mr *MockPluginRegisterMockRecorder) GetAddressURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddressURL", reflect.TypeOf((*MockPluginRegister)(nil).GetAddressURL))
}

// GetAllowEmptySecret mocks base method.
func (m *MockPluginRegister) GetAllowEmptySecret() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllowEmptySecret")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetAllowEmptySecret indicates an expected call of GetAllowEmptySecret.
func (mr *MockPluginRegisterMockRecorder) GetAllowEmptySecret() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllowEmptySecret", reflect.TypeOf((*MockPluginRegister)(nil).GetAllowEmptySecret))
}

// GetIntegrationClassName mocks base method.
func (m *MockPluginRegister) GetIntegrationClassName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntegrationClassName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetIntegrationClassName indicates an expected call of GetIntegrationClassName.
func (mr *MockPluginRegisterMockRecorder) GetIntegrationClassName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntegrationClassName", reflect.TypeOf((*MockPluginRegister)(nil).GetIntegrationClassName))
}

// GetReplicationPolicyTypes mocks base method.
func (m *MockPluginRegister) GetReplicationPolicyTypes() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReplicationPolicyTypes")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetReplicationPolicyTypes indicates an expected call of GetReplicationPolicyTypes.
func (mr *MockPluginRegisterMockRecorder) GetReplicationPolicyTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReplicationPolicyTypes", reflect.TypeOf((*MockPluginRegister)(nil).GetReplicationPolicyTypes))
}

// GetResourceTypes mocks base method.
func (m *MockPluginRegister) GetResourceTypes() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceTypes")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetResourceTypes indicates an expected call of GetResourceTypes.
func (mr *MockPluginRegisterMockRecorder) GetResourceTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceTypes", reflect.TypeOf((*MockPluginRegister)(nil).GetResourceTypes))
}

// GetSecretTypes mocks base method.
func (m *MockPluginRegister) GetSecretTypes() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretTypes")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetSecretTypes indicates an expected call of GetSecretTypes.
func (mr *MockPluginRegisterMockRecorder) GetSecretTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretTypes", reflect.TypeOf((*MockPluginRegister)(nil).GetSecretTypes))
}

// GetSupportedVersions mocks base method.
func (m *MockPluginRegister) GetSupportedVersions() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupportedVersions")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetSupportedVersions indicates an expected call of GetSupportedVersions.
func (mr *MockPluginRegisterMockRecorder) GetSupportedVersions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupportedVersions", reflect.TypeOf((*MockPluginRegister)(nil).GetSupportedVersions))
}

// GetWebhookURL mocks base method.
func (m *MockPluginRegister) GetWebhookURL() (*apis.URL, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWebhookURL")
	ret0, _ := ret[0].(*apis.URL)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetWebhookURL indicates an expected call of GetWebhookURL.
func (mr *MockPluginRegisterMockRecorder) GetWebhookURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWebhookURL", reflect.TypeOf((*MockPluginRegister)(nil).GetWebhookURL))
}

// Path mocks base method.
func (m *MockPluginRegister) Path() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path.
func (mr *MockPluginRegisterMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockPluginRegister)(nil).Path))
}

// Setup mocks base method.
func (m *MockPluginRegister) Setup(arg0 context.Context, arg1 *zap.SugaredLogger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockPluginRegisterMockRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockPluginRegister)(nil).Setup), arg0, arg1)
}

// MockPluginAddressable is a mock of PluginAddressable interface.
type MockPluginAddressable struct {
	ctrl     *gomock.Controller
	recorder *MockPluginAddressableMockRecorder
}

// MockPluginAddressableMockRecorder is the mock recorder for MockPluginAddressable.
type MockPluginAddressableMockRecorder struct {
	mock *MockPluginAddressable
}

// NewMockPluginAddressable creates a new mock instance.
func NewMockPluginAddressable(ctrl *gomock.Controller) *MockPluginAddressable {
	mock := &MockPluginAddressable{ctrl: ctrl}
	mock.recorder = &MockPluginAddressableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPluginAddressable) EXPECT() *MockPluginAddressableMockRecorder {
	return m.recorder
}

// GetAddressURL mocks base method.
func (m *MockPluginAddressable) GetAddressURL() *apis.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddressURL")
	ret0, _ := ret[0].(*apis.URL)
	return ret0
}

// GetAddressURL indicates an expected call of GetAddressURL.
func (mr *MockPluginAddressableMockRecorder) GetAddressURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddressURL", reflect.TypeOf((*MockPluginAddressable)(nil).GetAddressURL))
}

// MockDependentResourceGetter is a mock of DependentResourceGetter interface.
type MockDependentResourceGetter struct {
	ctrl     *gomock.Controller
	recorder *MockDependentResourceGetterMockRecorder
}

// MockDependentResourceGetterMockRecorder is the mock recorder for MockDependentResourceGetter.
type MockDependentResourceGetterMockRecorder struct {
	mock *MockDependentResourceGetter
}

// NewMockDependentResourceGetter creates a new mock instance.
func NewMockDependentResourceGetter(ctrl *gomock.Controller) *MockDependentResourceGetter {
	mock := &MockDependentResourceGetter{ctrl: ctrl}
	mock.recorder = &MockDependentResourceGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDependentResourceGetter) EXPECT() *MockDependentResourceGetterMockRecorder {
	return m.recorder
}

// GetDependentResources mocks base method.
func (m *MockDependentResourceGetter) GetDependentResources(arg0 context.Context, arg1 []v1alpha1.Param) ([]v1.ObjectReference, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDependentResources", arg0, arg1)
	ret0, _ := ret[0].([]v1.ObjectReference)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDependentResources indicates an expected call of GetDependentResources.
func (mr *MockDependentResourceGetterMockRecorder) GetDependentResources(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDependentResources", reflect.TypeOf((*MockDependentResourceGetter)(nil).GetDependentResources), arg0, arg1)
}

// MockAdditionalWebhookRegister is a mock of AdditionalWebhookRegister interface.
type MockAdditionalWebhookRegister struct {
	ctrl     *gomock.Controller
	recorder *MockAdditionalWebhookRegisterMockRecorder
}

// MockAdditionalWebhookRegisterMockRecorder is the mock recorder for MockAdditionalWebhookRegister.
type MockAdditionalWebhookRegisterMockRecorder struct {
	mock *MockAdditionalWebhookRegister
}

// NewMockAdditionalWebhookRegister creates a new mock instance.
func NewMockAdditionalWebhookRegister(ctrl *gomock.Controller) *MockAdditionalWebhookRegister {
	mock := &MockAdditionalWebhookRegister{ctrl: ctrl}
	mock.recorder = &MockAdditionalWebhookRegisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdditionalWebhookRegister) EXPECT() *MockAdditionalWebhookRegisterMockRecorder {
	return m.recorder
}

// GetWebhookSupport mocks base method.
func (m *MockAdditionalWebhookRegister) GetWebhookSupport() map[v1alpha1.WebhookEventSupportType][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWebhookSupport")
	ret0, _ := ret[0].(map[v1alpha1.WebhookEventSupportType][]string)
	return ret0
}

// GetWebhookSupport indicates an expected call of GetWebhookSupport.
func (mr *MockAdditionalWebhookRegisterMockRecorder) GetWebhookSupport() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWebhookSupport", reflect.TypeOf((*MockAdditionalWebhookRegister)(nil).GetWebhookSupport))
}

// MockResourcePathFormatter is a mock of ResourcePathFormatter interface.
type MockResourcePathFormatter struct {
	ctrl     *gomock.Controller
	recorder *MockResourcePathFormatterMockRecorder
}

// MockResourcePathFormatterMockRecorder is the mock recorder for MockResourcePathFormatter.
type MockResourcePathFormatterMockRecorder struct {
	mock *MockResourcePathFormatter
}

// NewMockResourcePathFormatter creates a new mock instance.
func NewMockResourcePathFormatter(ctrl *gomock.Controller) *MockResourcePathFormatter {
	mock := &MockResourcePathFormatter{ctrl: ctrl}
	mock.recorder = &MockResourcePathFormatterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourcePathFormatter) EXPECT() *MockResourcePathFormatterMockRecorder {
	return m.recorder
}

// GetResourcePathFmt mocks base method.
func (m *MockResourcePathFormatter) GetResourcePathFmt() map[v1alpha1.ResourcePathScene]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourcePathFmt")
	ret0, _ := ret[0].(map[v1alpha1.ResourcePathScene]string)
	return ret0
}

// GetResourcePathFmt indicates an expected call of GetResourcePathFmt.
func (mr *MockResourcePathFormatterMockRecorder) GetResourcePathFmt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourcePathFmt", reflect.TypeOf((*MockResourcePathFormatter)(nil).GetResourcePathFmt))
}

// GetSubResourcePathFmt mocks base method.
func (m *MockResourcePathFormatter) GetSubResourcePathFmt() map[v1alpha1.ResourcePathScene]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubResourcePathFmt")
	ret0, _ := ret[0].(map[v1alpha1.ResourcePathScene]string)
	return ret0
}

// GetSubResourcePathFmt indicates an expected call of GetSubResourcePathFmt.
func (mr *MockResourcePathFormatterMockRecorder) GetSubResourcePathFmt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubResourcePathFmt", reflect.TypeOf((*MockResourcePathFormatter)(nil).GetSubResourcePathFmt))
}

// MockPluginDisplayColumns is a mock of PluginDisplayColumns interface.
type MockPluginDisplayColumns struct {
	ctrl     *gomock.Controller
	recorder *MockPluginDisplayColumnsMockRecorder
}

// MockPluginDisplayColumnsMockRecorder is the mock recorder for MockPluginDisplayColumns.
type MockPluginDisplayColumnsMockRecorder struct {
	mock *MockPluginDisplayColumns
}

// NewMockPluginDisplayColumns creates a new mock instance.
func NewMockPluginDisplayColumns(ctrl *gomock.Controller) *MockPluginDisplayColumns {
	mock := &MockPluginDisplayColumns{ctrl: ctrl}
	mock.recorder = &MockPluginDisplayColumnsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPluginDisplayColumns) EXPECT() *MockPluginDisplayColumnsMockRecorder {
	return m.recorder
}

// GetDisplayColumns mocks base method.
func (m *MockPluginDisplayColumns) GetDisplayColumns() map[string]v1alpha1.DisplayColumns {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDisplayColumns")
	ret0, _ := ret[0].(map[string]v1alpha1.DisplayColumns)
	return ret0
}

// GetDisplayColumns indicates an expected call of GetDisplayColumns.
func (mr *MockPluginDisplayColumnsMockRecorder) GetDisplayColumns() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDisplayColumns", reflect.TypeOf((*MockPluginDisplayColumns)(nil).GetDisplayColumns))
}

// SetDisplayColumns mocks base method.
func (m *MockPluginDisplayColumns) SetDisplayColumns(arg0 string, arg1 ...v1alpha1.DisplayColumn) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SetDisplayColumns", varargs...)
}

// SetDisplayColumns indicates an expected call of SetDisplayColumns.
func (mr *MockPluginDisplayColumnsMockRecorder) SetDisplayColumns(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDisplayColumns", reflect.TypeOf((*MockPluginDisplayColumns)(nil).SetDisplayColumns), varargs...)
}

// MockPluginAttributes is a mock of PluginAttributes interface.
type MockPluginAttributes struct {
	ctrl     *gomock.Controller
	recorder *MockPluginAttributesMockRecorder
}

// MockPluginAttributesMockRecorder is the mock recorder for MockPluginAttributes.
type MockPluginAttributesMockRecorder struct {
	mock *MockPluginAttributes
}

// NewMockPluginAttributes creates a new mock instance.
func NewMockPluginAttributes(ctrl *gomock.Controller) *MockPluginAttributes {
	mock := &MockPluginAttributes{ctrl: ctrl}
	mock.recorder = &MockPluginAttributesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPluginAttributes) EXPECT() *MockPluginAttributesMockRecorder {
	return m.recorder
}

// Attributes mocks base method.
func (m *MockPluginAttributes) Attributes() map[string][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attributes")
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// Attributes indicates an expected call of Attributes.
func (mr *MockPluginAttributesMockRecorder) Attributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attributes", reflect.TypeOf((*MockPluginAttributes)(nil).Attributes))
}

// GetAttribute mocks base method.
func (m *MockPluginAttributes) GetAttribute(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttribute", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetAttribute indicates an expected call of GetAttribute.
func (mr *MockPluginAttributesMockRecorder) GetAttribute(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttribute", reflect.TypeOf((*MockPluginAttributes)(nil).GetAttribute), arg0)
}

// SetAttribute mocks base method.
func (m *MockPluginAttributes) SetAttribute(arg0 string, arg1 ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SetAttribute", varargs...)
}

// SetAttribute indicates an expected call of SetAttribute.
func (mr *MockPluginAttributesMockRecorder) SetAttribute(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAttribute", reflect.TypeOf((*MockPluginAttributes)(nil).SetAttribute), varargs...)
}

// MockPluginVersionAttributes is a mock of PluginVersionAttributes interface.
type MockPluginVersionAttributes struct {
	ctrl     *gomock.Controller
	recorder *MockPluginVersionAttributesMockRecorder
}

// MockPluginVersionAttributesMockRecorder is the mock recorder for MockPluginVersionAttributes.
type MockPluginVersionAttributesMockRecorder struct {
	mock *MockPluginVersionAttributes
}

// NewMockPluginVersionAttributes creates a new mock instance.
func NewMockPluginVersionAttributes(ctrl *gomock.Controller) *MockPluginVersionAttributes {
	mock := &MockPluginVersionAttributes{ctrl: ctrl}
	mock.recorder = &MockPluginVersionAttributesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPluginVersionAttributes) EXPECT() *MockPluginVersionAttributesMockRecorder {
	return m.recorder
}

// GetVersionAttributes mocks base method.
func (m *MockPluginVersionAttributes) GetVersionAttributes(arg0 string) map[string][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersionAttributes", arg0)
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// GetVersionAttributes indicates an expected call of GetVersionAttributes.
func (mr *MockPluginVersionAttributesMockRecorder) GetVersionAttributes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersionAttributes", reflect.TypeOf((*MockPluginVersionAttributes)(nil).GetVersionAttributes), arg0)
}

// SetVersionAttributes mocks base method.
func (m *MockPluginVersionAttributes) SetVersionAttributes(arg0 string, arg1 map[string][]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetVersionAttributes", arg0, arg1)
}

// SetVersionAttributes indicates an expected call of SetVersionAttributes.
func (mr *MockPluginVersionAttributesMockRecorder) SetVersionAttributes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetVersionAttributes", reflect.TypeOf((*MockPluginVersionAttributes)(nil).SetVersionAttributes), arg0, arg1)
}

// MockLivenessChecker is a mock of LivenessChecker interface.
type MockLivenessChecker struct {
	ctrl     *gomock.Controller
	recorder *MockLivenessCheckerMockRecorder
}

// MockLivenessCheckerMockRecorder is the mock recorder for MockLivenessChecker.
type MockLivenessCheckerMockRecorder struct {
	mock *MockLivenessChecker
}

// NewMockLivenessChecker creates a new mock instance.
func NewMockLivenessChecker(ctrl *gomock.Controller) *MockLivenessChecker {
	mock := &MockLivenessChecker{ctrl: ctrl}
	mock.recorder = &MockLivenessCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLivenessChecker) EXPECT() *MockLivenessCheckerMockRecorder {
	return m.recorder
}

// CheckAlive mocks base method.
func (m *MockLivenessChecker) CheckAlive(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAlive", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckAlive indicates an expected call of CheckAlive.
func (mr *MockLivenessCheckerMockRecorder) CheckAlive(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAlive", reflect.TypeOf((*MockLivenessChecker)(nil).CheckAlive), arg0)
}

// MockInitializer is a mock of Initializer interface.
type MockInitializer struct {
	ctrl     *gomock.Controller
	recorder *MockInitializerMockRecorder
}

// MockInitializerMockRecorder is the mock recorder for MockInitializer.
type MockInitializerMockRecorder struct {
	mock *MockInitializer
}

// NewMockInitializer creates a new mock instance.
func NewMockInitializer(ctrl *gomock.Controller) *MockInitializer {
	mock := &MockInitializer{ctrl: ctrl}
	mock.recorder = &MockInitializerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInitializer) EXPECT() *MockInitializerMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockInitializer) Initialize(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockInitializerMockRecorder) Initialize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockInitializer)(nil).Initialize), arg0)
}

// MockToolMetadataGetter is a mock of ToolMetadataGetter interface.
type MockToolMetadataGetter struct {
	ctrl     *gomock.Controller
	recorder *MockToolMetadataGetterMockRecorder
}

// MockToolMetadataGetterMockRecorder is the mock recorder for MockToolMetadataGetter.
type MockToolMetadataGetterMockRecorder struct {
	mock *MockToolMetadataGetter
}

// NewMockToolMetadataGetter creates a new mock instance.
func NewMockToolMetadataGetter(ctrl *gomock.Controller) *MockToolMetadataGetter {
	mock := &MockToolMetadataGetter{ctrl: ctrl}
	mock.recorder = &MockToolMetadataGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockToolMetadataGetter) EXPECT() *MockToolMetadataGetterMockRecorder {
	return m.recorder
}

// GetToolMetadata mocks base method.
func (m *MockToolMetadataGetter) GetToolMetadata(arg0 context.Context) (*v1alpha1.ToolMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToolMetadata", arg0)
	ret0, _ := ret[0].(*v1alpha1.ToolMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetToolMetadata indicates an expected call of GetToolMetadata.
func (mr *MockToolMetadataGetterMockRecorder) GetToolMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToolMetadata", reflect.TypeOf((*MockToolMetadataGetter)(nil).GetToolMetadata), arg0)
}