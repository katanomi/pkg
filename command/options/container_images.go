/*
Copyright 2022 The Katanomi Authors.

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

package options

import (
	"context"

	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"oras.land/oras-go/v2/registry"
)

// ContainerImagesOption describe container images option
type ContainerImagesOption struct {
	ContainerImages []string

	references    []registry.Reference
	requiredTag   bool
	requiredValue bool
}

// Setup init container images from args
func (m *ContainerImagesOption) Setup(ctx context.Context, _ *cobra.Command, args []string) (err error) {
	m.ContainerImages, _ = pkgargs.GetArrayValues(ctx, args, "container-images")
	return nil
}

// GetReferences returns references of the container-images.
func (m *ContainerImagesOption) GetReferences() []registry.Reference {
	return m.references
}

// SetTagRequired set tag is not an empty tag.
func (m *ContainerImagesOption) SetTagRequired(required bool) *ContainerImagesOption {
	m.requiredTag = required
	return m
}

// SetValueRequired set container-images as a required flag.
func (m *ContainerImagesOption) SetValueRequired(required bool) *ContainerImagesOption {
	m.requiredValue = required
	return m
}

// ValidateReferences check if the container images is valid
func (m *ContainerImagesOption) ValidateReferences(path *field.Path, references []registry.Reference) (errs field.ErrorList) {
	if m.requiredValue && len(references) == 0 {
		errs = append(errs, field.Required(path, "container-images must be set"))
	}

	for idx, reference := range references {
		if err := reference.ValidateReference(); err != nil {
			errs = append(errs, field.Invalid(path.Index(idx), reference, err.Error()))
		}

		if !m.requiredTag {
			continue
		}

		err := reference.ValidateReferenceAsTag()
		if err != nil {
			errs = append(errs, field.Invalid(path.Index(idx), reference, err.Error()))
		}
	}
	return
}

// Validate check if the container images is valid
func (m *ContainerImagesOption) Validate(path *field.Path) (errs field.ErrorList) {
	m.parseContainerImages()
	return m.ValidateReferences(path, m.references)
}

func (m *ContainerImagesOption) parseContainerImages() {
	if len(m.references) == 0 {
		m.references = make([]registry.Reference, 0, len(m.ContainerImages))

		for _, item := range m.ContainerImages {
			reference, _ := registry.ParseReference(item)
			m.references = append(m.references, reference)
		}
	}
}
