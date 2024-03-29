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
	"encoding/json"

	artifacts "github.com/katanomi/pkg/apis/artifacts/v1alpha1"
	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/katanomi/pkg/command/io"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ContainerImagesOption describe container images option
type ContainerImagesOption struct {
	ContainerImages []string
	Type            artifacts.ArtifactType
	ResultPath      string

	references []artifacts.URI
	parseErrs  field.ErrorList

	requiredTag   bool
	requiredValue bool
	withoutDigest bool
}

// AddFlags add flags for ContainerImageOption
func (m *ContainerImagesOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.ResultPath, "container-image-result-path", m.ResultPath, `filepath to store container image results`)
}

// Setup init container images from args
func (m *ContainerImagesOption) Setup(ctx context.Context, _ *cobra.Command, args []string) (err error) {
	m.ContainerImages, _ = pkgargs.GetArrayValues(ctx, args, "container-images")
	if m.Type == "" {
		m.Type = artifacts.ArtifactTypeContainerImage
	}
	return nil
}

// GetReferences returns references of the container-images.
func (m *ContainerImagesOption) GetReferences() []artifacts.URI {
	return m.references
}

// GetParseError returns references of the container-images.
func (m *ContainerImagesOption) GetParseError() field.ErrorList {
	return m.parseErrs
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

// SetWithoutDigest set tag is not an empty tag.
func (m *ContainerImagesOption) SetWithoutDigest(required bool) *ContainerImagesOption {
	m.withoutDigest = required
	return m
}

// ValidateReferences check if the container images is valid
func (m *ContainerImagesOption) ValidateReferences(path *field.Path, references []artifacts.URI) (errs field.ErrorList) {
	path = path.Child("container-images")
	if m.requiredValue && len(references) == 0 {
		errs = append(errs, field.Required(path, "container-images must be set"))
	}

	for idx, reference := range references {
		if err := reference.Validate(); err != nil {
			errs = append(errs, field.Invalid(path.Index(idx), reference, err.Error()))
			continue
		}

		if err := reference.ValidateTag(); err != nil && m.requiredTag {
			errs = append(errs, field.Invalid(path.Index(idx), reference, err.Error()))
			continue
		}

		if m.withoutDigest && reference.DigestString() != "" {
			errs = append(errs, field.Forbidden(path, "digest not allowed to be set."))
		}
	}

	errs = append(errs, m.parseErrs...)
	return
}

// Validate check if the container images is valid
func (m *ContainerImagesOption) Validate(path *field.Path) (errs field.ErrorList) {
	m.parseContainerImages()
	return m.ValidateReferences(path, m.references)
}

func (m *ContainerImagesOption) parseContainerImages() (errs field.ErrorList) {
	if len(m.references) == 0 {
		m.parseErrs = field.ErrorList{}
		m.references = make([]artifacts.URI, 0, len(m.ContainerImages))

		path := field.NewPath("container-images")
		for idx, item := range m.ContainerImages {
			reference, err := artifacts.ParseURI(item, m.Type)
			if err != nil {
				m.parseErrs = append(m.parseErrs, field.Invalid(path.Index(idx), item, err.Error()))
				continue
			}
			m.references = append(m.references, reference)
		}
	}

	return m.parseErrs
}

// WriteResult writes a result to the provided path if given
func (m *ContainerImagesOption) WriteResult(artfactList []artifacts.URI) (err error) {
	stringSlice := artifacts.AsDigestStringArray(artfactList...)
	var content []byte
	if content, err = json.Marshal(stringSlice); err != nil {
		return
	}
	err = io.WriteFile(m.ResultPath, content, 0777)
	return
}
