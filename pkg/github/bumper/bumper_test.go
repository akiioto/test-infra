/*
Copyright 2019 The Kubernetes Authors.

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

package bumper

import (
	"os"
	"strings"
	"testing"
)

func TestValidateOptions(t *testing.T) {
	emptyStr := ""
	cases := []struct {
		name                string
		githubToken         *string
		githubOrg           *string
		githubRepo          *string
		gerrit              *bool
		gerritAuthor        *string
		gerritPRIdentifier  *string
		gerritHostRepo      *string
		gerritCookieFile    *string
		remoteName          *string
		skipPullRequest     *bool
		signoff             *bool
		err                 bool
		upstreamBaseChanged bool
	}{
		{
			name: "Everything correct",
			err:  false,
		},
		{
			name:        "GitHubToken must not be empty when SkipPullRequest is false",
			githubToken: &emptyStr,
			err:         true,
		},
		{
			name:       "remoteName must not be empty when SkipPullRequest is false",
			remoteName: &emptyStr,
			err:        true,
		},
		{
			name:      "GitHubOrg cannot be empty when SkipPullRequest is false",
			githubOrg: &emptyStr,
			err:       true,
		},
		{
			name:       "GitHubRepo cannot be empty when SkipPullRequest is false",
			githubRepo: &emptyStr,
			err:        true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defaultOption := &Options{
				GitHubOrg:       "whatever-org",
				GitHubRepo:      "whatever-repo",
				GitHubLogin:     "whatever-login",
				GitHubToken:     "whatever-token",
				GitName:         "whatever-name",
				GitEmail:        "whatever-email",
				RemoteName:      "whatever-name",
				SkipPullRequest: false,
				Signoff:         false,
			}

			if tc.skipPullRequest != nil {
				defaultOption.SkipPullRequest = *tc.skipPullRequest
			}
			if tc.signoff != nil {
				defaultOption.Signoff = *tc.signoff
			}
			if tc.githubToken != nil {
				defaultOption.GitHubToken = *tc.githubToken
			}
			if tc.remoteName != nil {
				defaultOption.RemoteName = *tc.remoteName
			}
			if tc.githubOrg != nil {
				defaultOption.GitHubOrg = *tc.githubOrg
			}
			if tc.githubRepo != nil {
				defaultOption.GitHubRepo = *tc.githubRepo
			}

			err := validateOptions(defaultOption)
			t.Logf("err is: %v", err)
			if err == nil && tc.err {
				t.Errorf("Expected to get an error for %#v but got nil", defaultOption)
			}
			if err != nil && !tc.err {
				t.Errorf("Expected to not get an error for %#v but got %v", defaultOption, err)
			}
		})
	}
}

func writeToFile(t *testing.T, path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Errorf("write file %s dir with error '%v'", path, err)
	}
}

func TestGetAssignment(t *testing.T) {
	cases := []struct {
		description          string
		assignTo             string
		oncallURL            string
		oncallGroup          string
		oncallServerResponse string
		expectResKeyword     string
	}{
		{
			description:      "AssignTo takes precedence over oncall settings",
			assignTo:         "some-user",
			expectResKeyword: "/cc @some-user",
		},
		{
			description:      "No assign to",
			assignTo:         "",
			expectResKeyword: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			res := getAssignment(tc.assignTo)
			if !strings.Contains(res, tc.expectResKeyword) {
				t.Errorf("Expect the result %q contains keyword %q but it does not", res, tc.expectResKeyword)
			}
		})
	}
}
# (2025-03-04)