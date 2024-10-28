#!/bin/bash
# description: download tekton params definition and patches to remove some code related tekton pipeline
# using: ./hack/patch-tekton-params.sh

## how to generate hack/patches/0001-remove-unused-param-logic-based-on-tekton-definition.patch ##
# 1. start it from a clean branch without `param_type.go` and `param_type_test.go`
# 2. download tekton parms code and save to `./apis/meta/v1alpha1/` and commit it
# 3. change `param_type.go` `param_type_test.go` to your desired state and commit it
# 4. now you can get the patch file by executing `git format-patch HEAD^`

set -e

TEKTON_VERSION=v0.64.0

echo "download tekton params code..."
wget https://raw.githubusercontent.com/tektoncd/pipeline/refs/tags/${TEKTON_VERSION}/pkg/apis/pipeline/v1/param_types.go -O ./apis/meta/v1alpha1/param_types.go
wget https://raw.githubusercontent.com/tektoncd/pipeline/refs/tags/${TEKTON_VERSION}/pkg/apis/pipeline/v1/param_types_test.go -O ./apis/meta/v1alpha1/param_types_test.go

echo "apply patches..."
git apply ./hack/patches/0001-remove-unused-param-logic-based-on-tekton-definition.patch
