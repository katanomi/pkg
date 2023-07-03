package plugin

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GenerateTypeMeta(gvk schema.GroupVersionKind) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       gvk.Kind,
		APIVersion: gvk.GroupVersion().String(),
	}
}
