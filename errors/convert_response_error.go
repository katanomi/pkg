package errors

import (
	"context"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ConvertResponseError converts a gitlab response and an error into a kubernetes api error
// ctx is the basic context, response is the response object from gitlab sdk, err is the returned error
// gvk is the GroupVersionKind object with type meta for the object
// names supports one optional name to be given and will be attributed as the resource name in the returned error
func ConvertResponseError(ctx context.Context, response *http.Response, err error, gvk schema.GroupVersionKind, names ...string) error {
	if err == nil {
		return err
	}
	statusCode := http.StatusInternalServerError
	method := http.MethodGet
	name := ""
	if response != nil {
		statusCode = response.StatusCode
		method = response.Request.Method
	}
	if len(names) > 0 {
		name = names[0]
	} else if response != nil && response.Request != nil && response.Request.URL != nil {
		name = response.Request.URL.String()
	}

	return errors.NewGenericServerResponse(
		statusCode,
		method,
		schema.GroupResource{Group: gvk.Group, Resource: gvk.Kind},
		name, err.Error(),
		0,
		true)
}
