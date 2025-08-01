package main

import (
	"reflect"
	"testing"

	"github.com/kyma-project/test-infra/pkg/tags"
	"go.uber.org/zap"
)

func Test_ParseConfig(t *testing.T) {
	tc := []struct {
		name           string
		config         string
		expectedConfig Config
		expectErr      bool
	}{
		{
			name: "parsed full config one repo",
			config: `registry: kyma-project.io/prod-registry
dev-registry: dev.kyma-project.io/dev-registry
default-commit-tag:
  name: default_tag
  value: v{{ .Date }}-{{ .ShortSHA }}
  validation: ^(v[0-9]{8}-[0-9a-f]{8})$
default-pr-tag:
  name: default_tag
  value: pr-{{ .PRNumber }}
  validation: ^(PR-[0-9]+)$`,
			expectedConfig: Config{
				Registry:         []string{"kyma-project.io/prod-registry"},
				DevRegistry:      []string{"dev.kyma-project.io/dev-registry"},
				DefaultCommitTag: tags.Tag{Name: "default_tag", Value: `v{{ .Date }}-{{ .ShortSHA }}`, Validation: `^(v[0-9]{8}-[0-9a-f]{8})$`},
				DefaultPRTag:     tags.Tag{Name: "default_tag", Value: `pr-{{ .PRNumber }}`, Validation: `^(PR-[0-9]+)$`},
			},
		},
		{
			name: "parsed full config with multiple repos",
			config: `registry:
- kyma-project.io/prod-registry
- kyma-project.io/second-registry
dev-registry:
- dev.kyma-project.io/dev-registry
- dev.kyma-project.io/second-registry
default-commit-tag:
  name: default_tag
  value: v{{ .Date }}-{{ .ShortSHA }}
  validation: ^(v[0-9]{8}-[0-9a-f]{8})$
default-pr-tag:
  name: default_tag
  value: pr-{{ .PRNumber }}
  validation: ^(PR-[0-9]+)$`,
			expectedConfig: Config{
				Registry:         []string{"kyma-project.io/prod-registry", "kyma-project.io/second-registry"},
				DevRegistry:      []string{"dev.kyma-project.io/dev-registry", "dev.kyma-project.io/second-registry"},
				DefaultCommitTag: tags.Tag{Name: "default_tag", Value: `v{{ .Date }}-{{ .ShortSHA }}`, Validation: `^(v[0-9]{8}-[0-9a-f]{8})$`},
				DefaultPRTag:     tags.Tag{Name: "default_tag", Value: `pr-{{ .PRNumber }}`, Validation: `^(PR-[0-9]+)$`},
			},
		},
		{
			name:           "malformed yaml file, fail",
			config:         `garbage:malformed:config`,
			expectedConfig: Config{},
			expectErr:      true,
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			o := Config{}
			err := o.ParseConfig([]byte(c.config))
			if err != nil && !c.expectErr {
				t.Errorf("caught error, but didn't want to: %v", err)
			}
			if !reflect.DeepEqual(o, c.expectedConfig) {
				t.Errorf("%v != %v", o, c.expectedConfig)
			}
		})
	}
}

func TestLoadGitStateConfig(t *testing.T) {
	tc := []struct {
		name        string
		options     options
		env         map[string]string
		gitState    GitStateConfig
		expectError bool
	}{
		{
			name: "Load data from event payload for github pull_request_target",
			options: options{
				ciSystem: GithubActions,
			},
			env: map[string]string{
				"GITHUB_EVENT_PATH": "./test_fixture/pull_request_target_reopened.json",
				"GITHUB_EVENT_NAME": "pull_request_target",
			},
			gitState: GitStateConfig{
				RepositoryName:    "test-infra",
				RepositoryOwner:   "kyma-project",
				JobType:           "presubmit",
				PullRequestNumber: 10410,
				BaseCommitSHA:     "4b91c74a2aa9aeeb4a265cf1ffe2dd54812b4124",
				PullHeadCommitSHA: "8d0172d980653a377317a8bff9a6bb6ec2334801",
				isPullRequest:     true,
			},
		},
		{
			name: "Load data from event payload for github push event",
			options: options{
				ciSystem: GithubActions,
			},
			env: map[string]string{
				"GITHUB_EVENT_PATH": "./test_fixture/push_event.json",
				"GITHUB_EVENT_NAME": "push",
			},
			gitState: GitStateConfig{
				RepositoryName:  "test-infra",
				RepositoryOwner: "KacperMalachowski",
				JobType:         "postsubmit",
				BaseCommitSHA:   "d42f5051757b3e0699eb979d7581404e36fc0eee",
				isPullRequest:   false,
			},
		},
		{
			name: "Load data from event payload for github workflow_dispatch event",
			options: options{
				ciSystem: GithubActions,
			},
			env: map[string]string{
				"GITHUB_EVENT_PATH": "./test_fixture/workflow_dispatch.json",
				"GITHUB_EVENT_NAME": "workflow_dispatch",
				"GITHUB_SHA":        "d42f5051757b3e0699eb979d7581404e36fc0eee",
				"GITHUB_REF":        "refs/heads/main",
			},
			gitState: GitStateConfig{
				RepositoryName:  "test-infra",
				RepositoryOwner: "KacperMalachowski",
				JobType:         "workflow_dispatch",
				BaseCommitSHA:   "d42f5051757b3e0699eb979d7581404e36fc0eee",
				BaseCommitRef:   "refs/heads/main",
				isPullRequest:   false,
			},
		},
		{
			name: "Load data from event payload for github schedule event",
			options: options{
				ciSystem: GithubActions,
			},
			env: map[string]string{
				"GITHUB_EVENT_PATH": "./test_fixture/schedule.json",
				"GITHUB_EVENT_NAME": "schedule",
				"GITHUB_SHA":        "d42f5051757b3e0699eb979d7581404e36fc0eee",
				"GITHUB_REF":        "refs/heads/main",
			},
			gitState: GitStateConfig{
				RepositoryName:  "test-infra",
				RepositoryOwner: "KacperMalachowski",
				JobType:         "schedule",
				BaseCommitSHA:   "d42f5051757b3e0699eb979d7581404e36fc0eee",
				BaseCommitRef:   "refs/heads/main",
				isPullRequest:   false,
			},
		},
		{
			name: "Unknown ci system, return err",
			options: options{
				ciSystem: "",
			},
			env: map[string]string{
				"GITHUB_EVENT_PATH": "./test_fixture/pull_request_target_reopened.json",
				"GITHUB_EVENT_NAME": "pull_request_target",
			},
			expectError: true,
			gitState:    GitStateConfig{},
		},
		{
			name: "Unsupported github event, err",
			options: options{
				ciSystem: GithubActions,
			},
			env: map[string]string{
				"GITHUB_EVENT_PATH": "./test_fixture/pull_request_target_reopened.json",
				"GITHUB_EVENT_NAME": "pull_request",
			},
			expectError: true,
			gitState:    GitStateConfig{},
		},
		{
			name: "load data from event payload for github merge_group event",
			options: options{
				ciSystem: GithubActions,
			},
			env: map[string]string{
				"GITHUB_EVENT_PATH": "./test_fixture/merge-group_event.json",
				"GITHUB_EVENT_NAME": "merge_group",
				"GITHUB_SHA":        "659bf74f7b4ecab07d9398eec554217b51bad738",
				"GITHUB_REF":        "refs/heads/main",
			},
			gitState: GitStateConfig{
				RepositoryName:    "test-infra",
				RepositoryOwner:   "kyma-project",
				JobType:           "merge_group",
				BaseCommitSHA:     "659bf74f7b4ecab07d9398eec554217b51bad738",
				BaseCommitRef:     "refs/heads/main",
				isPullRequest:     true,
				PullHeadCommitSHA: "e47034172c36d3e5fb407b5ba57adf0f7868599d",
			},
		},
		{
			name: "load data from push event for jenkins",
			options: options{
				ciSystem: Jenkins,
			},
			env: map[string]string{
				"CHANGE_BRANCH": "refs/heads/main",
				"JENKINS_HOME":  "/some/absolute/path",
				"GIT_URL":       "github.com/kyma-project/test-infra.git",
				"GIT_COMMIT":    "1234",
			},
			gitState: GitStateConfig{
				RepositoryName:  "test-infra",
				RepositoryOwner: "kyma-project",
				JobType:         "postsubmit",
				BaseCommitSHA:   "1234",
			},
		},
		{
			name: "load data from pull request event for jenkins",
			options: options{
				ciSystem: Jenkins,
			},
			env: map[string]string{
				"CHANGE_BRANCH":   "refs/heads/main",
				"JENKINS_HOME":    "/some/absolute/path",
				"CHANGE_ID":       "14",
				"GIT_URL":         "github.com/kyma-project/test-infra.git",
				"GIT_COMMIT":      "1234",
				"CHANGE_BASE_SHA": "4321",
			},
			gitState: GitStateConfig{
				RepositoryName:    "test-infra",
				RepositoryOwner:   "kyma-project",
				JobType:           "presubmit",
				BaseCommitSHA:     "4321",
				BaseCommitRef:     "refs/heads/main",
				PullRequestNumber: 14,
				PullHeadCommitSHA: "1234",
				isPullRequest:     true,
			},
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			// Prepare env vars
			for key, value := range c.env {
				t.Setenv(key, value)
			}

			// Setup logger
			zapLogger, err := zap.NewDevelopment()
			if err != nil {
				t.Errorf("Failed to initialize logger: %s", err)
			}
			logger := zapLogger.Sugar()

			// Load git state
			state, err := LoadGitStateConfig(logger, c.options.ciSystem)
			if err != nil && !c.expectError {
				t.Errorf("unexpected error occured %s", err)
			}
			if err == nil && c.expectError {
				t.Error("expected error, but not occured")
			}

			if !reflect.DeepEqual(state, c.gitState) {
				t.Errorf("LoadGitStateConfigFromEnv(): Got %v, but expected %v", state, c.gitState)
			}
		})
	}
}

type mockEnv map[string]string

func (e mockEnv) mockGetenv(key string) string {
	return e[key]
}

func (e mockEnv) mockLookupEnv(key string) (string, bool) {
	env := e.mockGetenv(key)
	return env, env != ""
}

func Test_determineCISystem(t *testing.T) {
	tc := []struct {
		name      string
		env       mockEnv
		ciSystem  CISystem
		expectErr bool
	}{
		{
			name: "detect running in github actions",
			env: mockEnv{
				"GITHUB_ACTIONS": "true",
			},
			ciSystem: GithubActions,
		},
		{
			name: "detect running in jenkins",
			env: mockEnv{
				"JENKINS_HOME": "/some/absolute/path",
			},
			ciSystem: Jenkins,
		},
		{
			name: "unknown ci system",
			env: mockEnv{
				// Prevent false positivie detection of CI system running test
				"GITHUB_ACTIONS": "false",
				"PROW_JOB_ID":    "",
			},
			ciSystem:  "",
			expectErr: true,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			ciSystem, err := determineUsedCISystem(c.env.mockGetenv, c.env.mockLookupEnv)
			if err != nil && !c.expectErr {
				t.Errorf("got unexpected error: %s", err)
			}
			if err == nil && c.expectErr {
				t.Error("error expected, but no one occured")
			}

			if ciSystem != c.ciSystem {
				t.Errorf("determineCISystem(): Got %s, but expected %s", ciSystem, c.ciSystem)
			}
		})
	}
}
