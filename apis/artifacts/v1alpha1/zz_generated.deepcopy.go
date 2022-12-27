//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllFilter) DeepCopyInto(out *AllFilter) {
	*out = *in
	if in.items != nil {
		in, out := &in.items, &out.items
		*out = make([]metav1alpha1.ArtifactFilter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllFilter.
func (in *AllFilter) DeepCopy() *AllFilter {
	if in == nil {
		return nil
	}
	out := new(AllFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnyFilter) DeepCopyInto(out *AnyFilter) {
	*out = *in
	if in.items != nil {
		in, out := &in.items, &out.items
		*out = make([]metav1alpha1.ArtifactFilter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnyFilter.
func (in *AnyFilter) DeepCopy() *AnyFilter {
	if in == nil {
		return nil
	}
	out := new(AnyFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactVersion) DeepCopyInto(out *ArtifactVersion) {
	*out = *in
	if in.Versions != nil {
		in, out := &in.Versions, &out.Versions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactVersion.
func (in *ArtifactVersion) DeepCopy() *ArtifactVersion {
	if in == nil {
		return nil
	}
	out := new(ArtifactVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactVersionCollection) DeepCopyInto(out *ArtifactVersionCollection) {
	*out = *in
	if in.ArtifactVersions != nil {
		in, out := &in.ArtifactVersions, &out.ArtifactVersions
		*out = make([]ArtifactVersion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactVersionCollection.
func (in *ArtifactVersionCollection) DeepCopy() *ArtifactVersionCollection {
	if in == nil {
		return nil
	}
	out := new(ArtifactVersionCollection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerAuthItem) DeepCopyInto(out *DockerAuthItem) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerAuthItem.
func (in *DockerAuthItem) DeepCopy() *DockerAuthItem {
	if in == nil {
		return nil
	}
	out := new(DockerAuthItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerConfigJson) DeepCopyInto(out *DockerConfigJson) {
	*out = *in
	if in.Auths != nil {
		in, out := &in.Auths, &out.Auths
		*out = make(map[string]DockerAuthItem, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerConfigJson.
func (in *DockerConfigJson) DeepCopy() *DockerConfigJson {
	if in == nil {
		return nil
	}
	out := new(DockerConfigJson)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvFilter) DeepCopyInto(out *EnvFilter) {
	*out = *in
	in.ArtifactEnvFilter.DeepCopyInto(&out.ArtifactEnvFilter)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvFilter.
func (in *EnvFilter) DeepCopy() *EnvFilter {
	if in == nil {
		return nil
	}
	out := new(EnvFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in ImageAuths) DeepCopyInto(out *ImageAuths) {
	{
		in := &in
		*out = make(ImageAuths, len(*in))
		for key, val := range *in {
			var outVal []DockerAuthItem
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]DockerAuthItem, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageAuths.
func (in ImageAuths) DeepCopy() ImageAuths {
	if in == nil {
		return nil
	}
	out := new(ImageAuths)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageConfig) DeepCopyInto(out *ImageConfig) {
	*out = *in
	if in.ImageAuths != nil {
		in, out := &in.ImageAuths, &out.ImageAuths
		*out = make(ImageAuths, len(*in))
		for key, val := range *in {
			var outVal []DockerAuthItem
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]DockerAuthItem, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageConfig.
func (in *ImageConfig) DeepCopy() *ImageConfig {
	if in == nil {
		return nil
	}
	out := new(ImageConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelFilter) DeepCopyInto(out *LabelFilter) {
	*out = *in
	in.ArtifactLabelFilter.DeepCopyInto(&out.ArtifactLabelFilter)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelFilter.
func (in *LabelFilter) DeepCopy() *LabelFilter {
	if in == nil {
		return nil
	}
	out := new(LabelFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagFilter) DeepCopyInto(out *TagFilter) {
	*out = *in
	in.ArtifactTagFilter.DeepCopyInto(&out.ArtifactTagFilter)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagFilter.
func (in *TagFilter) DeepCopy() *TagFilter {
	if in == nil {
		return nil
	}
	out := new(TagFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *URI) DeepCopyInto(out *URI) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new URI.
func (in *URI) DeepCopy() *URI {
	if in == nil {
		return nil
	}
	out := new(URI)
	in.DeepCopyInto(out)
	return out
}
