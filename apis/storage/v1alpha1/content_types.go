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

package v1alpha1

import "github.com/emicklei/go-restful/v3"

const (
	// FileContentTypeJunitXML for junit xml file
	FileContentTypeJunitXML = "application/vnd.katanomi.junit+xml"
	// FileContentTypeAllureSite for allure site files
	FileContentTypeAllureSite = "application/vnd.katanomi.allure+site"
	// FileContentTypeGoJson for go json file
	FileContentTypeGoJson = "application/vnd.katanomi.gojson+json"
	// FileContentTypeJacocoSite for jacoco site files
	FileContentTypeJacocoSite = "application/vnd.katanomi.jacoco+site"
	// FileContentTypeGolangCoverageTxt for golang coverage text file
	FileContentTypeGolangCoverageTxt = "application/vnd.katanomi.golang-coverage+txt"
)

// SupportedContentTypeList contains content types restful apis should consume
var SupportedContentTypeList = []string{
	// unknown content types will use octet
	restful.MIME_OCTET,
	FileContentTypeJunitXML,
	FileContentTypeAllureSite,
	FileContentTypeGoJson,
	FileContentTypeJacocoSite,
	FileContentTypeGolangCoverageTxt,
}
