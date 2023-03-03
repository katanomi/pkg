// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/apis/meta/v1alpha1 (interfaces: TopLevelConditionObject)

// Package apis is a generated GoMock package.
package apis

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	apis "knative.dev/pkg/apis"
)

// MockTopLevelConditionObject is a mock of TopLevelConditionObject interface.
type MockTopLevelConditionObject struct {
	ctrl     *gomock.Controller
	recorder *MockTopLevelConditionObjectMockRecorder
}

// MockTopLevelConditionObjectMockRecorder is the mock recorder for MockTopLevelConditionObject.
type MockTopLevelConditionObjectMockRecorder struct {
	mock *MockTopLevelConditionObject
}

// NewMockTopLevelConditionObject creates a new mock instance.
func NewMockTopLevelConditionObject(ctrl *gomock.Controller) *MockTopLevelConditionObject {
	mock := &MockTopLevelConditionObject{ctrl: ctrl}
	mock.recorder = &MockTopLevelConditionObjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopLevelConditionObject) EXPECT() *MockTopLevelConditionObjectMockRecorder {
	return m.recorder
}

// GetAnnotations mocks base method.
func (m *MockTopLevelConditionObject) GetAnnotations() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnotations")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetAnnotations indicates an expected call of GetAnnotations.
func (mr *MockTopLevelConditionObjectMockRecorder) GetAnnotations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnotations", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetAnnotations))
}

// GetCreationTimestamp mocks base method.
func (m *MockTopLevelConditionObject) GetCreationTimestamp() v1.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreationTimestamp")
	ret0, _ := ret[0].(v1.Time)
	return ret0
}

// GetCreationTimestamp indicates an expected call of GetCreationTimestamp.
func (mr *MockTopLevelConditionObjectMockRecorder) GetCreationTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreationTimestamp", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetCreationTimestamp))
}

// GetDeletionGracePeriodSeconds mocks base method.
func (m *MockTopLevelConditionObject) GetDeletionGracePeriodSeconds() *int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeletionGracePeriodSeconds")
	ret0, _ := ret[0].(*int64)
	return ret0
}

// GetDeletionGracePeriodSeconds indicates an expected call of GetDeletionGracePeriodSeconds.
func (mr *MockTopLevelConditionObjectMockRecorder) GetDeletionGracePeriodSeconds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeletionGracePeriodSeconds", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetDeletionGracePeriodSeconds))
}

// GetDeletionTimestamp mocks base method.
func (m *MockTopLevelConditionObject) GetDeletionTimestamp() *v1.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeletionTimestamp")
	ret0, _ := ret[0].(*v1.Time)
	return ret0
}

// GetDeletionTimestamp indicates an expected call of GetDeletionTimestamp.
func (mr *MockTopLevelConditionObjectMockRecorder) GetDeletionTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeletionTimestamp", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetDeletionTimestamp))
}

// GetFinalizers mocks base method.
func (m *MockTopLevelConditionObject) GetFinalizers() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFinalizers")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetFinalizers indicates an expected call of GetFinalizers.
func (mr *MockTopLevelConditionObjectMockRecorder) GetFinalizers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFinalizers", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetFinalizers))
}

// GetGenerateName mocks base method.
func (m *MockTopLevelConditionObject) GetGenerateName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenerateName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetGenerateName indicates an expected call of GetGenerateName.
func (mr *MockTopLevelConditionObjectMockRecorder) GetGenerateName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenerateName", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetGenerateName))
}

// GetGeneration mocks base method.
func (m *MockTopLevelConditionObject) GetGeneration() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGeneration")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetGeneration indicates an expected call of GetGeneration.
func (mr *MockTopLevelConditionObjectMockRecorder) GetGeneration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGeneration", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetGeneration))
}

// GetLabels mocks base method.
func (m *MockTopLevelConditionObject) GetLabels() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLabels")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetLabels indicates an expected call of GetLabels.
func (mr *MockTopLevelConditionObjectMockRecorder) GetLabels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLabels", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetLabels))
}

// GetManagedFields mocks base method.
func (m *MockTopLevelConditionObject) GetManagedFields() []v1.ManagedFieldsEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedFields")
	ret0, _ := ret[0].([]v1.ManagedFieldsEntry)
	return ret0
}

// GetManagedFields indicates an expected call of GetManagedFields.
func (mr *MockTopLevelConditionObjectMockRecorder) GetManagedFields() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedFields", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetManagedFields))
}

// GetName mocks base method.
func (m *MockTopLevelConditionObject) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockTopLevelConditionObjectMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetName))
}

// GetNamespace mocks base method.
func (m *MockTopLevelConditionObject) GetNamespace() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespace")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNamespace indicates an expected call of GetNamespace.
func (mr *MockTopLevelConditionObjectMockRecorder) GetNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespace", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetNamespace))
}

// GetOwnerReferences mocks base method.
func (m *MockTopLevelConditionObject) GetOwnerReferences() []v1.OwnerReference {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnerReferences")
	ret0, _ := ret[0].([]v1.OwnerReference)
	return ret0
}

// GetOwnerReferences indicates an expected call of GetOwnerReferences.
func (mr *MockTopLevelConditionObjectMockRecorder) GetOwnerReferences() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnerReferences", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetOwnerReferences))
}

// GetResourceVersion mocks base method.
func (m *MockTopLevelConditionObject) GetResourceVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetResourceVersion indicates an expected call of GetResourceVersion.
func (mr *MockTopLevelConditionObjectMockRecorder) GetResourceVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceVersion", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetResourceVersion))
}

// GetSelfLink mocks base method.
func (m *MockTopLevelConditionObject) GetSelfLink() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelfLink")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetSelfLink indicates an expected call of GetSelfLink.
func (mr *MockTopLevelConditionObjectMockRecorder) GetSelfLink() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelfLink", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetSelfLink))
}

// GetTopLevelCondition mocks base method.
func (m *MockTopLevelConditionObject) GetTopLevelCondition() *apis.Condition {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopLevelCondition")
	ret0, _ := ret[0].(*apis.Condition)
	return ret0
}

// GetTopLevelCondition indicates an expected call of GetTopLevelCondition.
func (mr *MockTopLevelConditionObjectMockRecorder) GetTopLevelCondition() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopLevelCondition", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetTopLevelCondition))
}

// GetUID mocks base method.
func (m *MockTopLevelConditionObject) GetUID() types.UID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUID")
	ret0, _ := ret[0].(types.UID)
	return ret0
}

// GetUID indicates an expected call of GetUID.
func (mr *MockTopLevelConditionObjectMockRecorder) GetUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUID", reflect.TypeOf((*MockTopLevelConditionObject)(nil).GetUID))
}

// SetAnnotations mocks base method.
func (m *MockTopLevelConditionObject) SetAnnotations(arg0 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAnnotations", arg0)
}

// SetAnnotations indicates an expected call of SetAnnotations.
func (mr *MockTopLevelConditionObjectMockRecorder) SetAnnotations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAnnotations", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetAnnotations), arg0)
}

// SetCreationTimestamp mocks base method.
func (m *MockTopLevelConditionObject) SetCreationTimestamp(arg0 v1.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCreationTimestamp", arg0)
}

// SetCreationTimestamp indicates an expected call of SetCreationTimestamp.
func (mr *MockTopLevelConditionObjectMockRecorder) SetCreationTimestamp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCreationTimestamp", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetCreationTimestamp), arg0)
}

// SetDeletionGracePeriodSeconds mocks base method.
func (m *MockTopLevelConditionObject) SetDeletionGracePeriodSeconds(arg0 *int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDeletionGracePeriodSeconds", arg0)
}

// SetDeletionGracePeriodSeconds indicates an expected call of SetDeletionGracePeriodSeconds.
func (mr *MockTopLevelConditionObjectMockRecorder) SetDeletionGracePeriodSeconds(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeletionGracePeriodSeconds", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetDeletionGracePeriodSeconds), arg0)
}

// SetDeletionTimestamp mocks base method.
func (m *MockTopLevelConditionObject) SetDeletionTimestamp(arg0 *v1.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDeletionTimestamp", arg0)
}

// SetDeletionTimestamp indicates an expected call of SetDeletionTimestamp.
func (mr *MockTopLevelConditionObjectMockRecorder) SetDeletionTimestamp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeletionTimestamp", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetDeletionTimestamp), arg0)
}

// SetFinalizers mocks base method.
func (m *MockTopLevelConditionObject) SetFinalizers(arg0 []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFinalizers", arg0)
}

// SetFinalizers indicates an expected call of SetFinalizers.
func (mr *MockTopLevelConditionObjectMockRecorder) SetFinalizers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFinalizers", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetFinalizers), arg0)
}

// SetGenerateName mocks base method.
func (m *MockTopLevelConditionObject) SetGenerateName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGenerateName", arg0)
}

// SetGenerateName indicates an expected call of SetGenerateName.
func (mr *MockTopLevelConditionObjectMockRecorder) SetGenerateName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGenerateName", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetGenerateName), arg0)
}

// SetGeneration mocks base method.
func (m *MockTopLevelConditionObject) SetGeneration(arg0 int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGeneration", arg0)
}

// SetGeneration indicates an expected call of SetGeneration.
func (mr *MockTopLevelConditionObjectMockRecorder) SetGeneration(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGeneration", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetGeneration), arg0)
}

// SetLabels mocks base method.
func (m *MockTopLevelConditionObject) SetLabels(arg0 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLabels", arg0)
}

// SetLabels indicates an expected call of SetLabels.
func (mr *MockTopLevelConditionObjectMockRecorder) SetLabels(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLabels", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetLabels), arg0)
}

// SetManagedFields mocks base method.
func (m *MockTopLevelConditionObject) SetManagedFields(arg0 []v1.ManagedFieldsEntry) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetManagedFields", arg0)
}

// SetManagedFields indicates an expected call of SetManagedFields.
func (mr *MockTopLevelConditionObjectMockRecorder) SetManagedFields(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetManagedFields", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetManagedFields), arg0)
}

// SetName mocks base method.
func (m *MockTopLevelConditionObject) SetName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetName", arg0)
}

// SetName indicates an expected call of SetName.
func (mr *MockTopLevelConditionObjectMockRecorder) SetName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetName", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetName), arg0)
}

// SetNamespace mocks base method.
func (m *MockTopLevelConditionObject) SetNamespace(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNamespace", arg0)
}

// SetNamespace indicates an expected call of SetNamespace.
func (mr *MockTopLevelConditionObjectMockRecorder) SetNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNamespace", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetNamespace), arg0)
}

// SetOwnerReferences mocks base method.
func (m *MockTopLevelConditionObject) SetOwnerReferences(arg0 []v1.OwnerReference) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOwnerReferences", arg0)
}

// SetOwnerReferences indicates an expected call of SetOwnerReferences.
func (mr *MockTopLevelConditionObjectMockRecorder) SetOwnerReferences(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOwnerReferences", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetOwnerReferences), arg0)
}

// SetResourceVersion mocks base method.
func (m *MockTopLevelConditionObject) SetResourceVersion(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetResourceVersion", arg0)
}

// SetResourceVersion indicates an expected call of SetResourceVersion.
func (mr *MockTopLevelConditionObjectMockRecorder) SetResourceVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetResourceVersion", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetResourceVersion), arg0)
}

// SetSelfLink mocks base method.
func (m *MockTopLevelConditionObject) SetSelfLink(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetSelfLink", arg0)
}

// SetSelfLink indicates an expected call of SetSelfLink.
func (mr *MockTopLevelConditionObjectMockRecorder) SetSelfLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSelfLink", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetSelfLink), arg0)
}

// SetUID mocks base method.
func (m *MockTopLevelConditionObject) SetUID(arg0 types.UID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetUID", arg0)
}

// SetUID indicates an expected call of SetUID.
func (mr *MockTopLevelConditionObjectMockRecorder) SetUID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUID", reflect.TypeOf((*MockTopLevelConditionObject)(nil).SetUID), arg0)
}
