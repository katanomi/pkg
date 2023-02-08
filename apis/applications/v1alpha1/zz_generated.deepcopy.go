//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeployApplicationResults) DeepCopyInto(out *DeployApplicationResults) {
	*out = *in
	if in.ApplicationRef != nil {
		in, out := &in.ApplicationRef, &out.ApplicationRef
		*out = new(v1.ObjectReference)
		**out = **in
	}
	if in.Before != nil {
		in, out := &in.Before, &out.Before
		*out = make([]DeployApplicationStatus, len(*in))
		copy(*out, *in)
	}
	if in.After != nil {
		in, out := &in.After, &out.After
		*out = make([]DeployApplicationStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeployApplicationResults.
func (in *DeployApplicationResults) DeepCopy() *DeployApplicationResults {
	if in == nil {
		return nil
	}
	out := new(DeployApplicationResults)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeployApplicationStatus) DeepCopyInto(out *DeployApplicationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeployApplicationStatus.
func (in *DeployApplicationStatus) DeepCopy() *DeployApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(DeployApplicationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamedDeployApplicationResult) DeepCopyInto(out *NamedDeployApplicationResult) {
	*out = *in
	in.DeployApplicationResults.DeepCopyInto(&out.DeployApplicationResults)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamedDeployApplicationResult.
func (in *NamedDeployApplicationResult) DeepCopy() *NamedDeployApplicationResult {
	if in == nil {
		return nil
	}
	out := new(NamedDeployApplicationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in NamedDeployApplicationResults) DeepCopyInto(out *NamedDeployApplicationResults) {
	{
		in := &in
		*out = make(NamedDeployApplicationResults, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamedDeployApplicationResults.
func (in NamedDeployApplicationResults) DeepCopy() NamedDeployApplicationResults {
	if in == nil {
		return nil
	}
	out := new(NamedDeployApplicationResults)
	in.DeepCopyInto(out)
	return *out
}
