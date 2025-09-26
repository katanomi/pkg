/*
Copyright 2025 The Katanomi Authors.

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

package tekton

import (
	"context"

	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/logging"
)

// ItemProcessor defines callback function for processing named items
// Returns true if the item was found and processed, false otherwise
type ItemProcessor func(ctx context.Context, itemName string, item NamedItem) bool

// NamedItem defines interface for items that can be identified by name
type NamedItem interface {
	GetName() string
}

// ProcessSelectiveItems processes items selectively, skipping those whose names are in the exclude list.
// This generic function handles the common logic of filtering items and delegates the actual
// processing to the ItemProcessor callback.
func ProcessSelectiveItems[T NamedItem](
	ctx context.Context,
	items []T,
	excludeNames []string,
	processor ItemProcessor,
) {
	if processor == nil || len(items) == 0 {
		return
	}

	logger := logging.FromContext(ctx)

	// Create lookup map for excluded names
	excludeMap := make(map[string]struct{}, len(excludeNames))
	for _, name := range excludeNames {
		excludeMap[name] = struct{}{}
	}

	// Process each item that is not excluded
	for _, item := range items {
		itemName := item.GetName()

		// Skip if item name is in exclude list
		if _, excluded := excludeMap[itemName]; excluded {
			continue
		}

		logger.Debugw("processing item", "item", itemName)

		// Delegate the actual processing to the callback
		processed := processor(ctx, itemName, item)
		if !processed {
			logger.Debugw("item not processed", "item", itemName)
		}
	}
}

// NewSpecMerger creates an ItemProcessor for merging pipev1beta1.ParamSpec arrays
// This is a convenience function for common use cases where you want to merge specs directly
// It uses cached maps to avoid repeated array traversals for better performance
// defaultSource is used for object type parameters to get additional default values
func NewSpecMerger(
	destination *[]pipev1beta1.ParamSpec,
	source []pipev1beta1.ParamSpec,
	defaultSource []pipev1beta1.ParamSpec,
) ItemProcessor {
	// Convert to pointer array to reuse NewSpecMergerForPointers logic
	var pointerDestination *[]*pipev1beta1.ParamSpec
	if destination != nil {
		pointers := make([]*pipev1beta1.ParamSpec, len(*destination))
		for i := range *destination {
			pointers[i] = &(*destination)[i]
		}
		pointerDestination = &pointers
	}

	// Reuse NewSpecMergerForPointers implementation
	return NewSpecMergerForPointers(pointerDestination, source, defaultSource)
}

// NewSpecMergerForPointers creates an ItemProcessor for merging pipev1beta1.ParamSpec pointer arrays.
// This function allows direct modification of the original ParamSpec objects through pointers,
// avoiding the need for data copying and synchronization. If the source and destination
// parameter types do not match, the merger will deliberately skip updating metadata to
// preserve the destination parameter definition.
func NewSpecMergerForPointers(
	destination *[]*pipev1beta1.ParamSpec,
	source []pipev1beta1.ParamSpec,
	defaultSource []pipev1beta1.ParamSpec,
) ItemProcessor {
	// Create cached maps for efficient lookups
	destinationIndexMap := make(map[string]int)
	defaultSpecMap := make(map[string]pipev1beta1.ParamSpec)
	sourceSpecMap := make(map[string]pipev1beta1.ParamSpec)

	// Build destination index map
	if destination != nil {
		for i, spec := range *destination {
			if spec != nil {
				destinationIndexMap[spec.Name] = i
			}
		}
	}

	// Build source spec map
	for _, spec := range source {
		sourceSpecMap[spec.Name] = spec
	}

	// Build default spec map for object type default values
	for _, spec := range defaultSource {
		defaultSpecMap[spec.Name] = spec
	}

	return func(ctx context.Context, itemName string, item NamedItem) bool {
		if destination == nil {
			return false
		}

		// Find destination spec index using cached map
		destinationIndex, destinationExists := destinationIndexMap[itemName]
		if !destinationExists || destinationIndex >= len(*destination) || (*destination)[destinationIndex] == nil {
			return false
		}

		// Find source spec using cached map
		sourceSpec, sourceExists := sourceSpecMap[itemName]
		if !sourceExists {
			return false
		}
		logger := logging.FromContext(ctx)
		logger.Debugw("merging param spec (pointer version)", "param", itemName)

		destinationSpec := (*destination)[destinationIndex]
		destinationType := destinationSpec.Type
		sourceType := sourceSpec.Type
		if destinationType != "" && sourceType != "" && destinationType != sourceType {
			logger.Debugw("skip merging param metadata due to type mismatch",
				"param", itemName,
				"destinationType", destinationType,
				"sourceType", sourceType,
			)
			return true
		}

		// Update basic metadata from source - directly modify the pointed-to object
		destinationSpec.Description = sourceSpec.Description
		if sourceSpec.Default != nil {
			destinationSpec.Default = sourceSpec.Default.DeepCopy()
		}

		// Handle object parameter special processing
		if destinationSpec.Type == pipev1beta1.ParamTypeObject {
			// For object type, merge default values from defaultSource if available
			if defaultSpec, defaultExists := defaultSpecMap[itemName]; defaultExists && defaultSpec.Default != nil {
				// Merge object defaults: source takes precedence, but preserve defaultSource-only keys
				destinationSpec.Default = MergeObjectDefaults(defaultSpec.Default, destinationSpec.Default)
				logger.Debugw("merging param object defaults (pointer version)", "param", itemName)
			}

			// Cleanup default values based on properties schema
			CleanupObjectParamByProperties(destinationSpec)
		}

		(*destination)[destinationIndex] = destinationSpec

		return true
	}
}

// ExtractParamNames extracts parameter names from pipev1beta1.Param slice
// This is a convenience function for creating exclude lists from parameter arrays
func ExtractParamNames(params []pipev1beta1.Param) []string {
	names := make([]string, len(params))
	for i, param := range params {
		names[i] = param.Name
	}
	return names
}
