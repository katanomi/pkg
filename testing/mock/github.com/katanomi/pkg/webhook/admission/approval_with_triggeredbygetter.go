// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/katanomi/pkg/webhook/admission (interfaces: ApprovalWithTriggeredByGetter)

// Package admission is a generated GoMock package.
package admission

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
)

// MockApprovalWithTriggeredByGetter is a mock of ApprovalWithTriggeredByGetter interface.
type MockApprovalWithTriggeredByGetter struct {
	ctrl     *gomock.Controller
	recorder *MockApprovalWithTriggeredByGetterMockRecorder
}

// MockApprovalWithTriggeredByGetterMockRecorder is the mock recorder for MockApprovalWithTriggeredByGetter.
type MockApprovalWithTriggeredByGetterMockRecorder struct {
	mock *MockApprovalWithTriggeredByGetter
}

// NewMockApprovalWithTriggeredByGetter creates a new mock instance.
func NewMockApprovalWithTriggeredByGetter(ctrl *gomock.Controller) *MockApprovalWithTriggeredByGetter {
	mock := &MockApprovalWithTriggeredByGetter{ctrl: ctrl}
	mock.recorder = &MockApprovalWithTriggeredByGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApprovalWithTriggeredByGetter) EXPECT() *MockApprovalWithTriggeredByGetterMockRecorder {
	return m.recorder
}

// DeepCopyObject mocks base method.
func (m *MockApprovalWithTriggeredByGetter) DeepCopyObject() runtime.Object {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeepCopyObject")
	ret0, _ := ret[0].(runtime.Object)
	return ret0
}

// DeepCopyObject indicates an expected call of DeepCopyObject.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) DeepCopyObject() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeepCopyObject", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).DeepCopyObject))
}

// GetAnnotations mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetAnnotations() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnotations")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetAnnotations indicates an expected call of GetAnnotations.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetAnnotations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnotations", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetAnnotations))
}

// GetApprovalSpecs mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetApprovalSpecs(arg0 runtime.Object) []*v1alpha1.ApprovalSpec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApprovalSpecs", arg0)
	ret0, _ := ret[0].([]*v1alpha1.ApprovalSpec)
	return ret0
}

// GetApprovalSpecs indicates an expected call of GetApprovalSpecs.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetApprovalSpecs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApprovalSpecs", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetApprovalSpecs), arg0)
}

// GetChecks mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetChecks(arg0 runtime.Object) []*v1alpha1.Check {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChecks", arg0)
	ret0, _ := ret[0].([]*v1alpha1.Check)
	return ret0
}

// GetChecks indicates an expected call of GetChecks.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetChecks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChecks", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetChecks), arg0)
}

// GetClusterName mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetClusterName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetClusterName indicates an expected call of GetClusterName.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetClusterName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterName", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetClusterName))
}

// GetCreationTimestamp mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetCreationTimestamp() v1.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreationTimestamp")
	ret0, _ := ret[0].(v1.Time)
	return ret0
}

// GetCreationTimestamp indicates an expected call of GetCreationTimestamp.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetCreationTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreationTimestamp", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetCreationTimestamp))
}

// GetDeletionGracePeriodSeconds mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetDeletionGracePeriodSeconds() *int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeletionGracePeriodSeconds")
	ret0, _ := ret[0].(*int64)
	return ret0
}

// GetDeletionGracePeriodSeconds indicates an expected call of GetDeletionGracePeriodSeconds.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetDeletionGracePeriodSeconds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeletionGracePeriodSeconds", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetDeletionGracePeriodSeconds))
}

// GetDeletionTimestamp mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetDeletionTimestamp() *v1.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeletionTimestamp")
	ret0, _ := ret[0].(*v1.Time)
	return ret0
}

// GetDeletionTimestamp indicates an expected call of GetDeletionTimestamp.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetDeletionTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeletionTimestamp", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetDeletionTimestamp))
}

// GetFinalizers mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetFinalizers() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFinalizers")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetFinalizers indicates an expected call of GetFinalizers.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetFinalizers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFinalizers", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetFinalizers))
}

// GetGenerateName mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetGenerateName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenerateName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetGenerateName indicates an expected call of GetGenerateName.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetGenerateName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenerateName", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetGenerateName))
}

// GetGeneration mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetGeneration() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGeneration")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetGeneration indicates an expected call of GetGeneration.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetGeneration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGeneration", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetGeneration))
}

// GetLabels mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetLabels() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLabels")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetLabels indicates an expected call of GetLabels.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetLabels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLabels", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetLabels))
}

// GetManagedFields mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetManagedFields() []v1.ManagedFieldsEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedFields")
	ret0, _ := ret[0].([]v1.ManagedFieldsEntry)
	return ret0
}

// GetManagedFields indicates an expected call of GetManagedFields.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetManagedFields() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedFields", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetManagedFields))
}

// GetName mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetName))
}

// GetNamespace mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetNamespace() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespace")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNamespace indicates an expected call of GetNamespace.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespace", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetNamespace))
}

// GetObjectKind mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetObjectKind() schema.ObjectKind {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectKind")
	ret0, _ := ret[0].(schema.ObjectKind)
	return ret0
}

// GetObjectKind indicates an expected call of GetObjectKind.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetObjectKind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectKind", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetObjectKind))
}

// GetOwnerReferences mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetOwnerReferences() []v1.OwnerReference {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnerReferences")
	ret0, _ := ret[0].([]v1.OwnerReference)
	return ret0
}

// GetOwnerReferences indicates an expected call of GetOwnerReferences.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetOwnerReferences() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnerReferences", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetOwnerReferences))
}

// GetResourceVersion mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetResourceVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetResourceVersion indicates an expected call of GetResourceVersion.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetResourceVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceVersion", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetResourceVersion))
}

// GetSelfLink mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetSelfLink() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelfLink")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetSelfLink indicates an expected call of GetSelfLink.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetSelfLink() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelfLink", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetSelfLink))
}

// GetTriggeredBy mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetTriggeredBy(arg0 runtime.Object) *v1alpha1.TriggeredBy {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTriggeredBy", arg0)
	ret0, _ := ret[0].(*v1alpha1.TriggeredBy)
	return ret0
}

// GetTriggeredBy indicates an expected call of GetTriggeredBy.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetTriggeredBy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTriggeredBy", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetTriggeredBy), arg0)
}

// GetUID mocks base method.
func (m *MockApprovalWithTriggeredByGetter) GetUID() types.UID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUID")
	ret0, _ := ret[0].(types.UID)
	return ret0
}

// GetUID indicates an expected call of GetUID.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) GetUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUID", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).GetUID))
}

// ModifiedOthers mocks base method.
func (m *MockApprovalWithTriggeredByGetter) ModifiedOthers(arg0, arg1 runtime.Object) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModifiedOthers", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ModifiedOthers indicates an expected call of ModifiedOthers.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) ModifiedOthers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModifiedOthers", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).ModifiedOthers), arg0, arg1)
}

// SetAnnotations mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetAnnotations(arg0 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAnnotations", arg0)
}

// SetAnnotations indicates an expected call of SetAnnotations.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetAnnotations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAnnotations", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetAnnotations), arg0)
}

// SetClusterName mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetClusterName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetClusterName", arg0)
}

// SetClusterName indicates an expected call of SetClusterName.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetClusterName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetClusterName", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetClusterName), arg0)
}

// SetCreationTimestamp mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetCreationTimestamp(arg0 v1.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCreationTimestamp", arg0)
}

// SetCreationTimestamp indicates an expected call of SetCreationTimestamp.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetCreationTimestamp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCreationTimestamp", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetCreationTimestamp), arg0)
}

// SetDeletionGracePeriodSeconds mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetDeletionGracePeriodSeconds(arg0 *int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDeletionGracePeriodSeconds", arg0)
}

// SetDeletionGracePeriodSeconds indicates an expected call of SetDeletionGracePeriodSeconds.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetDeletionGracePeriodSeconds(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeletionGracePeriodSeconds", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetDeletionGracePeriodSeconds), arg0)
}

// SetDeletionTimestamp mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetDeletionTimestamp(arg0 *v1.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDeletionTimestamp", arg0)
}

// SetDeletionTimestamp indicates an expected call of SetDeletionTimestamp.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetDeletionTimestamp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeletionTimestamp", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetDeletionTimestamp), arg0)
}

// SetFinalizers mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetFinalizers(arg0 []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFinalizers", arg0)
}

// SetFinalizers indicates an expected call of SetFinalizers.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetFinalizers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFinalizers", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetFinalizers), arg0)
}

// SetGenerateName mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetGenerateName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGenerateName", arg0)
}

// SetGenerateName indicates an expected call of SetGenerateName.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetGenerateName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGenerateName", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetGenerateName), arg0)
}

// SetGeneration mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetGeneration(arg0 int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGeneration", arg0)
}

// SetGeneration indicates an expected call of SetGeneration.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetGeneration(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGeneration", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetGeneration), arg0)
}

// SetLabels mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetLabels(arg0 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLabels", arg0)
}

// SetLabels indicates an expected call of SetLabels.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetLabels(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLabels", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetLabels), arg0)
}

// SetManagedFields mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetManagedFields(arg0 []v1.ManagedFieldsEntry) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetManagedFields", arg0)
}

// SetManagedFields indicates an expected call of SetManagedFields.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetManagedFields(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetManagedFields", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetManagedFields), arg0)
}

// SetName mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetName", arg0)
}

// SetName indicates an expected call of SetName.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetName", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetName), arg0)
}

// SetNamespace mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetNamespace(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNamespace", arg0)
}

// SetNamespace indicates an expected call of SetNamespace.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNamespace", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetNamespace), arg0)
}

// SetOwnerReferences mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetOwnerReferences(arg0 []v1.OwnerReference) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOwnerReferences", arg0)
}

// SetOwnerReferences indicates an expected call of SetOwnerReferences.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetOwnerReferences(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOwnerReferences", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetOwnerReferences), arg0)
}

// SetResourceVersion mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetResourceVersion(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetResourceVersion", arg0)
}

// SetResourceVersion indicates an expected call of SetResourceVersion.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetResourceVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetResourceVersion", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetResourceVersion), arg0)
}

// SetSelfLink mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetSelfLink(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetSelfLink", arg0)
}

// SetSelfLink indicates an expected call of SetSelfLink.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetSelfLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSelfLink", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetSelfLink), arg0)
}

// SetUID mocks base method.
func (m *MockApprovalWithTriggeredByGetter) SetUID(arg0 types.UID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetUID", arg0)
}

// SetUID indicates an expected call of SetUID.
func (mr *MockApprovalWithTriggeredByGetterMockRecorder) SetUID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUID", reflect.TypeOf((*MockApprovalWithTriggeredByGetter)(nil).SetUID), arg0)
}
