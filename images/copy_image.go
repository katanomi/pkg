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

package images

import (
	"context"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
)

// CopyImage Copy the method of copy image in github.com/containers/image for internal use.
func CopyImage(ctx context.Context, policyContext *signature.PolicyContext, destRef,
	srcRef types.ImageReference, options *copy.Options) (copiedManifest []byte, copyErr error) {
	return copy.Image(ctx, policyContext, destRef, srcRef, options)
}
