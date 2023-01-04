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
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/test"
	"github.com/stretchr/testify/assert"
)

func TestRunContainers(t *testing.T) {
	tests := []struct {
		name            string
		options         actions.RunContainersOptions
		expectedCommand string
	}{
		{
			name: "with services",
			options: actions.RunContainersOptions{
				GenerationPath: filepath.Join("a", "b", "c", "d"),
				Services:       []string{"validator", "consensus", "execution"},
			},
			expectedCommand: fmt.Sprintf("docker compose -f %s up -d validator consensus execution", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
		},
		{
			name: "without services",
			options: actions.RunContainersOptions{
				GenerationPath: filepath.Join("a", "b", "c", "d"),
				Services:       []string{},
			},
			expectedCommand: fmt.Sprintf("docker compose -f %s up -d", filepath.Join("a", "b", "c", "d", "docker-compose.yml")),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			commandRunner := &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, error) {
					assert.Equal(t, tc.expectedCommand, c.Cmd)
					return "", nil
				},
			}
			sedgeActions := actions.NewSedgeActions(nil, nil, commandRunner)
			sedgeActions.RunContainers(tc.options)
		})
	}
}
