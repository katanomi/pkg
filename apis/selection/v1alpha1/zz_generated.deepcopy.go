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
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BaseFilterRule) DeepCopyInto(out *BaseFilterRule) {
	*out = *in
	if in.Exact != nil {
		in, out := &in.Exact, &out.Exact
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BaseFilterRule.
func (in *BaseFilterRule) DeepCopy() *BaseFilterRule {
	if in == nil {
		return nil
	}
	out := new(BaseFilterRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFilter) DeepCopyInto(out *ClusterFilter) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Filter != nil {
		in, out := &in.Filter, &out.Filter
		*out = new(ClusterFilterRule)
		(*in).DeepCopyInto(*out)
	}
	if in.Refs != nil {
		in, out := &in.Refs, &out.Refs
		*out = make([]corev1.ObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFilter.
func (in *ClusterFilter) DeepCopy() *ClusterFilter {
	if in == nil {
		return nil
	}
	out := new(ClusterFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFilterRule) DeepCopyInto(out *ClusterFilterRule) {
	*out = *in
	if in.Exact != nil {
		in, out := &in.Exact, &out.Exact
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFilterRule.
func (in *ClusterFilterRule) DeepCopy() *ClusterFilterRule {
	if in == nil {
		return nil
	}
	out := new(ClusterFilterRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceFilter) DeepCopyInto(out *NamespaceFilter) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Filter != nil {
		in, out := &in.Filter, &out.Filter
		*out = new(NamespaceFilterRule)
		(*in).DeepCopyInto(*out)
	}
	if in.Refs != nil {
		in, out := &in.Refs, &out.Refs
		*out = make([]corev1.ObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceFilter.
func (in *NamespaceFilter) DeepCopy() *NamespaceFilter {
	if in == nil {
		return nil
	}
	out := new(NamespaceFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceFilterRule) DeepCopyInto(out *NamespaceFilterRule) {
	*out = *in
	if in.Exact != nil {
		in, out := &in.Exact, &out.Exact
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceFilterRule.
func (in *NamespaceFilterRule) DeepCopy() *NamespaceFilterRule {
	if in == nil {
		return nil
	}
	out := new(NamespaceFilterRule)
	in.DeepCopyInto(out)
	return out
}
