/*
Copyright 2023 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package variable

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/utils/field"
)

var (
	rbacv1SubjectName         = reflect.TypeOf(rbacv1.Subject{}).Name()
	corev1ObjectReferenceName = reflect.TypeOf(corev1.ObjectReference{}).Name()
	buildGitBranchStatusName  = reflect.TypeOf(v1alpha1.BuildGitBranchStatus{}).Name()

	rbacv1SubjectPkgPath         = reflect.TypeOf(rbacv1.Subject{}).PkgPath()
	corev1ObjectReferencePkgPath = reflect.TypeOf(corev1.ObjectReference{}).PkgPath()
	buildGitBranchStatusPkgPath  = reflect.TypeOf(v1alpha1.BuildGitBranchStatus{}).PkgPath()
)

// DefaultNameMarshalFuncs define a default custom type conversion function
var DefaultNameMarshalFuncs = map[string]ConvertFunc{
	rbacv1SubjectName:         MarshalSubject,
	corev1ObjectReferenceName: MarshalObjectReference,
	buildGitBranchStatusName:  MarshalBuildGitBranchStatus,
}

// MarshalSubject marshal the Subject of k8s.io/api/rbac/v1 to variables.
func MarshalSubject(st reflect.Type, base *field.Path, _ MarshalFuncManager) ([]Variable, error) {
	if err := matchObjectPath(st, rbacv1SubjectName, rbacv1SubjectPkgPath); err != nil {
		return []Variable{}, err
	}

	return []Variable{
		{Name: base.Child("kind").String(), Example: "User"},
		{Name: base.Child("apiGroup").String()},
		{Name: base.Child("name").String(), Example: "joedoe@example.com"},
		{Name: base.Child("namespace").String(), Example: "default"},
	}, nil
}

// MarshalObjectReference marshal the ObjectReference of k8s.io/api/core/v1 to variables.
func MarshalObjectReference(st reflect.Type, base *field.Path, _ MarshalFuncManager) ([]Variable, error) {
	if err := matchObjectPath(st, corev1ObjectReferenceName, corev1ObjectReferencePkgPath); err != nil {
		return []Variable{}, err
	}

	return []Variable{
		{Name: base.Child("kind").String(), Example: "DeliveryRun"},
		{Name: base.Child("namespace").String(), Example: "default"},
		{Name: base.Child("name").String(), Example: "delivery-run-abdexy"},
		{Name: base.Child("uid").String(), Example: "b2fab970-f672-4af0-a9cd-5ad9a8dbcc29"},
		{Name: base.Child("apiVersion").String(), Example: "deliveries.katanomi.dev/v1alpha1"},
	}, nil
}

// MarshalBuildGitBranchStatus marshal the BuildGitBranchStatus of github.com/katanomi/pkg/apis/meta/v1alpha1 to variables
func MarshalBuildGitBranchStatus(st reflect.Type, base *field.Path, convertFuncs MarshalFuncManager) ([]Variable, error) {
	if err := matchObjectPath(st, buildGitBranchStatusName, buildGitBranchStatusPkgPath); err != nil {
		return []Variable{}, err
	}

	// using default func to marshal.
	variables, err := marshalStruct(st, base, convertFuncs)
	if err != nil {
		return []Variable{}, err
	}

	// when parent node is not branch, BuildGitBranchStatus.Name label should be empty.
	baseString := "." + base.String()
	if !strings.HasSuffix(baseString, ".branch") {
		for i, item := range variables {
			// check is name node.
			if strings.HasSuffix(item.Name, ".name") {
				item.Label = ""
				variables[i] = item
				break
			}
		}
	}
	return variables, nil
}

func matchObjectPath(st reflect.Type, name, pkgPath string) error {
	if name != st.Name() || st.PkgPath() != pkgPath {
		return fmt.Errorf(
			"get marshal type[%s/%s] don't match %s/%s",
			st.PkgPath(), st.Name(), pkgPath, name)
	}
	return nil
}
