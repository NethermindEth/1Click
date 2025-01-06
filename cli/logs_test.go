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
package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NethermindEth/sedge/configs"
	sedge_mocks "github.com/NethermindEth/sedge/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/NethermindEth/sedge/internal/compose"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/internal/pkg/dependencies"
	"github.com/NethermindEth/sedge/test"
	log "github.com/sirupsen/logrus"
)

func TestLogs(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)
	tests := []struct {
		name  string
		args  []string
		setup func(*sedge_mocks.MockDependenciesManager, *sedge_mocks.MockSedgeActions)
		cmd   commands.CommandRunner
		err   string
	}{
		{
			name: "docker not installed",
			setup: func(depsMgr *sedge_mocks.MockDependenciesManager, sedgeActions *sedge_mocks.MockSedgeActions) {
				depsMgr.EXPECT().Check([]string{"docker"}).Return(nil, []string{"docker"}).Times(1)
			},
			err: "missing dependencies: docker",
		},
		{
			name: "docker compose not installed",
			setup: func(depsMgr *sedge_mocks.MockDependenciesManager, sedgeActions *sedge_mocks.MockSedgeActions) {
				gomock.InOrder(
					depsMgr.EXPECT().Check([]string{"docker"}).Return([]string{"docker"}, nil).Times(1),
					depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1),
					depsMgr.EXPECT().DockerComposeIsInstalled().Return(fmt.Errorf("%w: %s", dependencies.ErrDependencyNotInstalled, "docker-compose")).Times(1),
				)
			},
			err: "dependency not installed: docker-compose",
		},
		{
			name: "docker engine not running",
			setup: func(depsMgr *sedge_mocks.MockDependenciesManager, sedgeActions *sedge_mocks.MockSedgeActions) {
				gomock.InOrder(
					depsMgr.EXPECT().Check([]string{"docker"}).Return([]string{"docker"}, nil).Times(1),
					depsMgr.EXPECT().DockerEngineIsOn().Return(dependencies.ErrDockerEngineIsNotRunning).Times(1),
				)
			},
			err: "docker engine is not running",
		},
		{
			name: "logs stopped by user",
			setup: func(depsMgr *sedge_mocks.MockDependenciesManager, sedgeActions *sedge_mocks.MockSedgeActions) {
				gomock.InOrder(
					depsMgr.EXPECT().Check([]string{"docker"}).Return([]string{"docker"}, nil).Times(1),
					depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1),
					depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).Times(1),
					sedgeActions.EXPECT().ValidateDockerComposeFile(filepath.Join(configs.DefaultAbsSedgeDataPath, "docker-compose.yml")).Return(nil).Times(1),
				)
			},
			cmd: &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, int, error) {
					if strings.Contains(c.Cmd, "docker compose -f") && strings.Contains(c.Cmd, "logs") {
						return "", 130, nil
					}
					return "", 0, nil
				},
				SRunBash: func(bs commands.ScriptFile) (string, error) {
					return "", nil
				},
			},
		},
		{
			name: "error",
			setup: func(depsMgr *sedge_mocks.MockDependenciesManager, sedgeActions *sedge_mocks.MockSedgeActions) {
				gomock.InOrder(
					depsMgr.EXPECT().Check([]string{"docker"}).Return([]string{"docker"}, nil).Times(1),
					depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1),
					depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).Times(1),
					sedgeActions.EXPECT().ValidateDockerComposeFile(filepath.Join(configs.DefaultAbsSedgeDataPath, "docker-compose.yml")).Return(nil).Times(1),
				)
			},
			cmd: &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, int, error) {
					if strings.Contains(c.Cmd, "docker compose -f") && strings.Contains(c.Cmd, "logs") {
						return "", 1, errors.New("error")
					}
					return "", 0, nil
				},
				SRunBash: func(bs commands.ScriptFile) (string, error) {
					return "", nil
				},
			},
			err: "failed to get logs for services . Error: Docker Compose Manager running 'docker compose logs': error. Output: ",
		},
		{
			name: "services arg",
			args: []string{"execution", "consensus"},
			setup: func(depsMgr *sedge_mocks.MockDependenciesManager, sedgeActions *sedge_mocks.MockSedgeActions) {
				gomock.InOrder(
					depsMgr.EXPECT().Check([]string{"docker"}).Return([]string{"docker"}, nil).Times(1),
					depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1),
					depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).Times(1),
					sedgeActions.EXPECT().ValidateDockerComposeFile(filepath.Join(configs.DefaultAbsSedgeDataPath, "docker-compose.yml")).Return(nil).Times(1),
				)
			},
			cmd: &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, int, error) {
					if strings.Contains(c.Cmd, "docker compose -f") && strings.Contains(c.Cmd, "logs") {
						return "", 0, nil
					}
					return "", 0, nil
				},
				SRunBash: func(bs commands.ScriptFile) (string, error) {
					return "", nil
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			composeMgr := compose.NewComposeManager(tc.cmd)

			depsMgr := sedge_mocks.NewMockDependenciesManager(ctrl)
			sedgeActions := sedge_mocks.NewMockSedgeActions(ctrl)

			tc.setup(depsMgr, sedgeActions)

			cmd := LogsCmd(*composeMgr, tc.cmd, sedgeActions, depsMgr)
			cmd.SetOutput(io.Discard)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			if tc.err != "" {
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
