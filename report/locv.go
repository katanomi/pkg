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

package report

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
)

const (
	// TypeLcov is the type of lcov
	TypeLcov ReportType = "lcov"
)

// LcovParser lcov report parser
type LcovParser struct {
	// LineFound found lines
	LineFound int
	// LineHit converage hit lines
	LineHit int
	// BranchFound found branch
	BranchFound int
	// BranchHit converage hit branch
	BranchHit int
}

const (
	// detail: https://ltp.sourceforge.net/coverage/lcov/geninfo.1.php
	// number of instrumented lines
	LineFound = "LF"
	// number of lines with a non-zero execution count
	LineHit = "LH"
	// number of branches found
	BranchFound = "BRF"
	// number of branches hit
	BranchHit = "BRH"
)

// Parse parse lcov report.
func (p *LcovParser) Parse(path string) (result interface{}, err error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		if err = p.parseLine(string(line)); err != nil {
			return nil, fmt.Errorf("invalid lcov text:%s. error: %s", string(line), err.Error())
		}
	}
	return p, nil
}

// parseLine parse lcov report line data.
func (p *LcovParser) parseLine(line string) (err error) {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, ":")
	// line is "TN:" or "end_of_record", ignore parsing to prevent out-of-array problems.
	if len(parts) < 2 {
		return nil
	}

	num := 0
	switch parts[0] {
	case LineFound:
		num, err = strconv.Atoi(parts[1])
		p.LineFound += num
	case LineHit:
		num, err = strconv.Atoi(parts[1])
		p.LineHit += num
	case BranchFound:
		num, err = strconv.Atoi(parts[1])
		p.BranchFound += num
	case BranchHit:
		num, err = strconv.Atoi(parts[1])
		p.BranchHit += num
	default:
		// no action
	}
	return
}

// ConvertToTestCoverage convert to TestCoverage
func (p *LcovParser) ConvertToTestCoverage() v1alpha1.TestCoverage {
	testCoverage := v1alpha1.TestCoverage{}

	testCoverage.Branches = fmt.Sprintf("%.2f", float64(p.BranchHit)/float64(p.BranchFound)*100)
	testCoverage.Lines = fmt.Sprintf("%.2f", float64(p.LineHit)/float64(p.LineFound)*100)

	return testCoverage
}
