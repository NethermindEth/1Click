/*
Copyright 2022 Nethermind

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
package actions_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/internal/compose"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/test"
	"github.com/stretchr/testify/assert"
)

func TestSetupContainers(t *testing.T) {
	tests := []struct {
		name              string
		options           actions.SetupContainersOptions
		expectedBuildCmd  string
		expectedPullCmd   string
		expectedCreateCmd string
	}{
		{
			name: "with services",
			options: actions.SetupContainersOptions{
				GenerationPath: filepath.Join("a", "b", "c", "d"),
				Services:       []string{"validator", "consensus", "execution"},
			},
			expectedBuildCmd:  fmt.Sprintf("docker compose -f %s build validator consensus execution", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
			expectedPullCmd:   fmt.Sprintf("docker compose -f %s pull validator consensus execution", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
			expectedCreateCmd: fmt.Sprintf("docker compose -f %s create validator consensus execution", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
		},
		{
			name: "without services",
			options: actions.SetupContainersOptions{
				GenerationPath: filepath.Join("a", "b", "c", "d"),
				Services:       []string{},
			},
			expectedBuildCmd:  fmt.Sprintf("docker compose -f %s build", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
			expectedPullCmd:   fmt.Sprintf("docker compose -f %s pull", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
			expectedCreateCmd: fmt.Sprintf("docker compose -f %s create", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
		},
		{
			name: "skip-pull",
			options: actions.SetupContainersOptions{
				GenerationPath: filepath.Join("a", "b", "c", "d"),
				Services:       []string{"execution", "consensus"},
				SkipPull:       true,
			},
			expectedBuildCmd:  fmt.Sprintf("docker compose -f %s build execution consensus", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
			expectedPullCmd:   "",
			expectedCreateCmd: fmt.Sprintf("docker compose -f %s create execution consensus", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			commandRunner := &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, int, error) {
					assert.Contains(t, []string{tc.expectedBuildCmd, tc.expectedPullCmd, tc.expectedCreateCmd}, c.Cmd)
					return "", 0, nil
				},
			}
			composeMgr := compose.NewComposeManager(commandRunner)
			sedgeActions := actions.NewSedgeActions(actions.SedgeActionsOptions{
				CommandRunner:  commandRunner,
				ComposeManager: *composeMgr,
			})
			sedgeActions.SetupContainers(tc.options)
		})
	}
}
