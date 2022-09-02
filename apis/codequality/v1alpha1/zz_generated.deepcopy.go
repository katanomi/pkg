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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnalisysMetrics) DeepCopyInto(out *AnalisysMetrics) {
	*out = *in
	if in.Branch != nil {
		in, out := &in.Branch, &out.Branch
		*out = new(CodeChangeMetrics)
		**out = **in
	}
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(CodeChangeMetrics)
		**out = **in
	}
	if in.Ratings != nil {
		in, out := &in.Ratings, &out.Ratings
		*out = make(map[string]AnalysisRating, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Languages != nil {
		in, out := &in.Languages, &out.Languages
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.CodeSize != nil {
		in, out := &in.CodeSize, &out.CodeSize
		*out = new(CodeSize)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnalisysMetrics.
func (in *AnalisysMetrics) DeepCopy() *AnalisysMetrics {
	if in == nil {
		return nil
	}
	out := new(AnalisysMetrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnalysisRating) DeepCopyInto(out *AnalysisRating) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnalysisRating.
func (in *AnalysisRating) DeepCopy() *AnalysisRating {
	if in == nil {
		return nil
	}
	out := new(AnalysisRating)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnalysisResult) DeepCopyInto(out *AnalysisResult) {
	*out = *in
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = new(AnalisysMetrics)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnalysisResult.
func (in *AnalysisResult) DeepCopy() *AnalysisResult {
	if in == nil {
		return nil
	}
	out := new(AnalysisResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeChangeMetrics) DeepCopyInto(out *CodeChangeMetrics) {
	*out = *in
	out.CoverageRate = in.CoverageRate
	out.DuplicationRate = in.DuplicationRate
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeChangeMetrics.
func (in *CodeChangeMetrics) DeepCopy() *CodeChangeMetrics {
	if in == nil {
		return nil
	}
	out := new(CodeChangeMetrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeChangeRates) DeepCopyInto(out *CodeChangeRates) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeChangeRates.
func (in *CodeChangeRates) DeepCopy() *CodeChangeRates {
	if in == nil {
		return nil
	}
	out := new(CodeChangeRates)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeLintIssues) DeepCopyInto(out *CodeLintIssues) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeLintIssues.
func (in *CodeLintIssues) DeepCopy() *CodeLintIssues {
	if in == nil {
		return nil
	}
	out := new(CodeLintIssues)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeLintResult) DeepCopyInto(out *CodeLintResult) {
	*out = *in
	if in.Issues != nil {
		in, out := &in.Issues, &out.Issues
		*out = new(CodeLintIssues)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeLintResult.
func (in *CodeLintResult) DeepCopy() *CodeLintResult {
	if in == nil {
		return nil
	}
	out := new(CodeLintResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeSize) DeepCopyInto(out *CodeSize) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeSize.
func (in *CodeSize) DeepCopy() *CodeSize {
	if in == nil {
		return nil
	}
	out := new(CodeSize)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestCoverage) DeepCopyInto(out *TestCoverage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestCoverage.
func (in *TestCoverage) DeepCopy() *TestCoverage {
	if in == nil {
		return nil
	}
	out := new(TestCoverage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestResult) DeepCopyInto(out *TestResult) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestResult.
func (in *TestResult) DeepCopy() *TestResult {
	if in == nil {
		return nil
	}
	out := new(TestResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UnitTestsResult) DeepCopyInto(out *UnitTestsResult) {
	*out = *in
	if in.Coverage != nil {
		in, out := &in.Coverage, &out.Coverage
		*out = new(TestCoverage)
		**out = **in
	}
	if in.TestResult != nil {
		in, out := &in.TestResult, &out.TestResult
		*out = new(TestResult)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UnitTestsResult.
func (in *UnitTestsResult) DeepCopy() *UnitTestsResult {
	if in == nil {
		return nil
	}
	out := new(UnitTestsResult)
	in.DeepCopyInto(out)
	return out
}