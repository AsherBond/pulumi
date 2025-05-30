// Copyright 2016-2025, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ciutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectVars(t *testing.T) {
	buildNumber := "123"
	buildID := "87638724"
	systemAndEnvVars := map[SystemName]map[string]string{
		// Since the `pulumi/pulumi` repo runs on Travis,
		// we set the TRAVIS env var to an empty string for all test cases
		// except for the Travis one itself.
		// This way when the unit test runs on Travis, we don't pick-up Travis env vars.
		AzurePipelines: {
			"TRAVIS":         "",
			"TF_BUILD":       "true",
			"BUILD_BUILDID":  buildNumber,
			"GITHUB_ACTIONS": "",
		},
		Buildkite: {
			"GITHUB_ACTIONS":         "",
			"TRAVIS":                 "",
			"BUILDKITE":              "true",
			"BUILDKITE_BUILD_NUMBER": buildNumber,
			"BUILDKITE_BUILD_ID":     buildID,
		},
		CircleCI: {
			"TRAVIS":           "",
			"CIRCLECI":         "true",
			"CIRCLE_BUILD_NUM": buildNumber,
			"GITHUB_ACTIONS":   "",
		},
		Codefresh: {
			"TRAVIS":         "",
			"CF_BUILD_URL":   "https://g.codefresh.io/build/99f5d825577e23c56f8c6b2a",
			"CF_BUILD_ID":    buildNumber,
			"GITHUB_ACTIONS": "",
		},
		GenericCI: {
			"TRAVIS":             "",
			"PULUMI_CI_SYSTEM":   "generic-ci-system",
			"PULUMI_CI_BUILD_ID": buildNumber,
			"GITHUB_ACTIONS":     "",
		},
		GitLab: {
			"TRAVIS":          "",
			"GITLAB_CI":       "true",
			"CI_PIPELINE_ID":  buildID,
			"CI_PIPELINE_IID": buildNumber,
			"GITHUB_ACTIONS":  "",
		},
		Travis: {
			"TRAVIS":            "true",
			"TRAVIS_JOB_ID":     buildID,
			"TRAVIS_JOB_NUMBER": buildNumber,
			"GITHUB_ACTIONS":    "",
		},
	}

	for system := range systemAndEnvVars {
		t.Run(fmt.Sprintf("Test_%v_Detection", system), func(t *testing.T) {
			envVars := systemAndEnvVars[system]
			for envVar := range envVars {
				t.Setenv(envVar, envVars[envVar])
			}
			vars := DetectVars()
			// For CI systems where the build number and the ID are the same,
			// only the BuildID is set and is considered to be the build "number".
			if vars.BuildNumber == "" {
				assert.Equal(t,
					buildNumber, vars.BuildID,
					"%v did not set the expected build ID %v in the Vars struct.", system, buildNumber)
			} else {
				assert.Equal(t,
					buildID, vars.BuildID,
					"%v did not set the expected build ID %v in the Vars struct.", system, buildNumber)
				assert.Equal(t,
					buildNumber, vars.BuildNumber,
					"%v did not set the expected build number %v in the Vars struct.", system, buildNumber)
			}
		})
	}
}

func TestDetectVarsBaseCI(t *testing.T) {
	systemAndEnvVars := map[SystemName]map[string]string{
		// Since the `pulumi/pulumi` repo runs on Travis,
		// we set the TRAVIS env var to an empty string for all test cases
		// except for the Travis one itself.
		// This way when the unit test runs on Travis, we don't pick-up Travis env vars.
		AppVeyor: {
			"TRAVIS":         "",
			"APPVEYOR":       "true",
			"GITHUB_ACTIONS": "",
		},
		Codeship: {
			"TRAVIS":         "",
			"CI_NAME":        "codeship",
			"GITHUB_ACTIONS": "",
		},
	}

	for system := range systemAndEnvVars {
		t.Run(fmt.Sprintf("Test_%v_Detection", system), func(t *testing.T) {
			envVars := systemAndEnvVars[system]
			for envVar := range envVars {
				t.Setenv(envVar, envVars[envVar])
			}
			vars := DetectVars()
			assert.Equal(t,
				string(system), string(vars.Name),
				"%v did not set the expected CI system name in the Vars struct.", system)
		})
	}
}

func TestDetectVarsDisableCIDetection(t *testing.T) {
	t.Setenv("PULUMI_DISABLE_CI_DETECTION", "nonEmptyString")
	t.Setenv("TRAVIS", "true")
	t.Setenv("TRAVIS_JOB_ID", "1234")

	v := DetectVars()
	assert.Equal(t, "", v.BuildID)
}

func TestDetectVarsPrefersGenericCI(t *testing.T) {
	name := "generic-ci-system"

	t.Setenv("PULUMI_CI_SYSTEM", name)
	t.Setenv("TRAVIS", "true")

	v := DetectVars()

	assert.Equal(t, name, string(v.Name))
}
