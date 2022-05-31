// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"

	scapiv1alpha3 "github.com/operator-framework/api/pkg/apis/scorecard/v1alpha3"
	apimanifests "github.com/operator-framework/api/pkg/manifests"
)

// This is the custom scorecard test example binary
// As with the Redhat scorecard test image, the bundle that is under
// test is expected to be mounted so that tests can inspect the
// bundle contents as part of their test implementations.
// The actual test is to be run is named and that name is passed
// as an argument to this binary.  This argument mechanism allows
// this binary to run various tests all from within a single
// test image.

const PodBundleRoot = "/bundle"

var TestMap = map[string]func(*apimanifests.Bundle) scapiv1alpha3.TestStatus{
	"simple-test": SimpleTest,
}

func main() {
	entrypoint := os.Args[1:]
	if len(entrypoint) == 0 {
		log.Fatal("Test name argument is required")
	}

	// Read the pod's untar'd bundle from a well-known path.
	cfg, err := apimanifests.GetBundleFromDir(PodBundleRoot)
	if err != nil {
		log.Fatal(err.Error())
	}

	var result scapiv1alpha3.TestStatus

	// Names of the custom tests which would be passed in the
	// `operator-sdk` command.
	if testFunc, exist := TestMap[entrypoint[0]]; exist {
		result = testFunc(cfg)
	} else {
		result = PrintValidTests()
	}

	// Convert scapiv1alpha3.TestResult to json.
	prettyJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))

}

func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// PrintValidTests will print out full list of test names to give a hint to the end user on what the valid tests are.
func PrintValidTests() scapiv1alpha3.TestStatus {
	result := scapiv1alpha3.TestResult{}
	result.State = scapiv1alpha3.FailState
	result.Errors = make([]string, 0)
	result.Suggestions = make([]string, 0)

	keys := make([]string, len(TestMap))
	idx := 0
	for _, key := range TestMap {
		keys[idx] = reflect.ValueOf(key).String()
		idx++
	}

	str := fmt.Sprintf("Valid tests for this image include: %s", strings.Join(keys, " "))
	result.Errors = append(result.Errors, str)
	return scapiv1alpha3.TestStatus{
		Results: []scapiv1alpha3.TestResult{result},
	}
}

// CustomTest1 and CustomTest2 are example test functions. Relevant operator specific
// test logic is to be implemented in similarly.
func SimpleTest(bundle *apimanifests.Bundle) scapiv1alpha3.TestStatus {
	r := scapiv1alpha3.TestResult{}
	r.Name = GetFuncName()
	r.State = scapiv1alpha3.PassState
	r.Errors = make([]string, 0)
	r.Suggestions = make([]string, 0)
	almExamples := bundle.CSV.GetAnnotations()["alm-examples"]
	if almExamples == "" {
		fmt.Println("no alm-examples in the bundle CSV")
	}

	return wrapResult(r)
}

func wrapResult(r scapiv1alpha3.TestResult) scapiv1alpha3.TestStatus {
	return scapiv1alpha3.TestStatus{
		Results: []scapiv1alpha3.TestResult{r},
	}
}
