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
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NethermindEth/sedge/cli/actions"
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

type downCmdTestCase struct {
	generationPath string
	composeMgr     compose.ComposeManager
	runner         commands.CommandRunner
	depsMgr        dependencies.DependenciesManager
	sedgeActions   actions.SedgeActions
	fdOut          *bytes.Buffer
	isErr          bool
	ctrl           *gomock.Controller
}

func buildDownTestCase(t *testing.T, caseName string, isErr bool, path string) *downCmdTestCase {
	tc := downCmdTestCase{}

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		t.Fatalf("Can't build test case: %v", err)
	}

	tc.ctrl = gomock.NewController(t)
	depsMgr := sedge_mocks.NewMockDependenciesManager(tc.ctrl)
	sedgeActions := sedge_mocks.NewMockSedgeActions(tc.ctrl)
	depsMgr.EXPECT().Check([]string{dependencies.Docker}).Return([]string{dependencies.Docker}, nil).AnyTimes()
	depsMgr.EXPECT().DockerEngineIsOn().Return(nil).AnyTimes()
	depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).AnyTimes()
	sedgeActions.EXPECT().ValidateDockerComposeFile(filepath.Join(path, "docker-compose.yml")).Return(nil).AnyTimes()
	tc.depsMgr = depsMgr
	tc.sedgeActions = sedgeActions

	err := test.PrepareTestCaseDir(filepath.Join("testdata", "down_tests", caseName, configs.DefaultSedgeDataFolderName), path)
	if err != nil {
		t.Fatalf("Can't build test case: %v", err)
	}

	// TODO: allow runner edition
	tc.runner = &test.SimpleCMDRunner{
		SRunCMD: func(c commands.Command) (string, int, error) {
			return "", 0, nil
		},
		SRunBash: func(bs commands.ScriptFile) (string, error) {
			return "", nil
		},
	}
	composeMgr := compose.NewComposeManager(tc.runner)
	tc.composeMgr = *composeMgr

	tc.generationPath = path
	tc.fdOut = new(bytes.Buffer)
	tc.isErr = isErr
	return &tc
}

func TestDownCmd(t *testing.T) {
	configPath := t.TempDir()
	tcs := []downCmdTestCase{
		*buildDownTestCase(t, "case_1", false, filepath.Join(configPath, configs.DefaultSedgeDataFolderName)),
	}

	for _, tc := range tcs {
		rootCmd := RootCmd()
		rootCmd.AddCommand(DownCmd(tc.composeMgr, tc.runner, tc.sedgeActions, tc.depsMgr))
		rootCmd.SetArgs([]string{"down", "--path", tc.generationPath})
		rootCmd.SetOut(tc.fdOut)
		log.SetOutput(tc.fdOut)

		descr := "sedge down"
		err := rootCmd.Execute()
		if tc.isErr && err == nil {
			t.Errorf("%s expected to fail", descr)
		} else if !tc.isErr && err != nil {
			t.Errorf("%s failed: %v", descr, err)
		}
		tc.ctrl.Finish()
	}
}

func TestDown_Error(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)
	tests := []struct {
		name       string
		runner     commands.CommandRunner
		err        func(*downCmdTestCase, string) string
		customPath string
	}{
		{
			name: "docker compose ps error, PreCheck error",
			runner: &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, int, error) {
					if strings.Contains(c.Cmd, "docker compose") && strings.Contains(c.Cmd, "ps") {
						return "", 1, errors.New("runner error")
					}
					return "", 0, nil
				},
				SRunBash: func(bs commands.ScriptFile) (string, error) {
					return "", nil
				},
			},
			err: func(tt *downCmdTestCase, path string) string {
				return "services of docker-compose script provided are not running. Error: runner error"
			},
		},
		{
			name: "docker compose ps --status running error, Check Containers error",
			runner: &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, int, error) {
					if strings.Contains(c.Cmd, "docker compose") && strings.Contains(c.Cmd, "ps") && strings.Contains(c.Cmd, "--filter status=running") {
						return "", 1, errors.New("runner error")
					}
					return "", 0, nil
				},
				SRunBash: func(bs commands.ScriptFile) (string, error) {
					return "", nil
				},
			},
			err: func(tt *downCmdTestCase, path string) string {
				return "services of docker-compose script provided are not running. Error: runner error"
			},
		},
		{
			name: "docker compose down error",
			runner: &test.SimpleCMDRunner{
				SRunCMD: func(c commands.Command) (string, int, error) {
					if strings.Contains(c.Cmd, "docker compose") && strings.Contains(c.Cmd, "ps") && strings.Contains(c.Cmd, "--filter status=running") {
						return "", 0, nil
					}
					if strings.Contains(c.Cmd, "docker compose") && strings.Contains(c.Cmd, "down") {
						return "", 1, errors.New("runner error")
					}
					return "", 0, nil
				},
				SRunBash: func(bs commands.ScriptFile) (string, error) {
					return "", nil
				},
			},
			err: func(tt *downCmdTestCase, path string) string {
				return fmt.Sprintf("command 'docker compose -f %s down' throws error: runner error", filepath.Join(tt.generationPath, "docker-compose.yml"))
			},
		},
		// {
		// 	name: "Generation path error",
		// 	runner: &test.SimpleCMDRunner{
		// 		SRunCMD: func(c commands.Command) (string, int, error) {
		// 			return "", 0, nil
		// 		},
		// 		SRunBash: func(bs commands.ScriptFile) (string, error) {
		// 			return "", nil
		// 		},
		// 	},
		// 	err: func(tt *downCmdTestCase, path string) string {
		// 		return fmt.Sprintf(configs.DockerComposeScriptNotFoundError, path, configs.DefaultAbsSedgeDataPath)
		// 	},
		// 	customPath: t.TempDir(),
		// },
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configPath := t.TempDir()
			pathFlag := filepath.Join(configPath, configs.DefaultSedgeDataFolderName)
			if tc.customPath != "" {
				pathFlag = tc.customPath
			}
			tt := buildDownTestCase(t, "case_1", true, pathFlag)
			downCmd := DownCmd(tt.composeMgr, tc.runner, tt.sedgeActions, tt.depsMgr)
			downCmd.SetArgs([]string{"--path", pathFlag})
			downCmd.SetOut(io.Discard)
			err := downCmd.Execute()
			if err != nil {
				assert.EqualError(t, err, tc.err(tt, pathFlag))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
