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
package cli_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/NethermindEth/sedge/cli"
	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/dependencies"
	sedge_mocks "github.com/NethermindEth/sedge/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSlashingExport_ValidatorIsRequired(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	tests := []struct {
		name        string
		args        []string
		expectedErr error
	}{
		{
			name:        "without flags",
			args:        []string{},
			expectedErr: cli.ErrInvalidNumberOfArguments,
		},
		{
			name:        "with flags",
			args:        []string{"--network", "sepolia", "--out", "slashing_data.json"},
			expectedErr: cli.ErrInvalidNumberOfArguments,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			depsMgr := sedge_mocks.NewMockDependenciesManager(ctrl)
			slashingExportCmd := cli.SlashingExportCmd(nil, depsMgr)
			slashingExportCmd.SetArgs(tt.args)
			slashingExportCmd.SetOutput(io.Discard)
			err := slashingExportCmd.Execute()
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestSlashingExport_Params(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	customDir := t.TempDir()
	outDir := t.TempDir()

	t.Cleanup(func() {
		if err := os.RemoveAll(customDir); err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(outDir); err != nil {
			t.Fatal(err)
		}
	})

	tests := []struct {
		name          string
		args          []string
		actionOptions actions.SlashingExportOptions
	}{
		{
			name: "validator argument",
			args: []string{"lighthouse"},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "lighthouse",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "network flag",
			args: []string{"lighthouse", "--network", "sepolia"},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "lighthouse",
				Network:         "sepolia",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "network shorthand flag",
			args: []string{"lighthouse", "-n", "sepolia"},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "lighthouse",
				Network:         "sepolia",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "stop-validator flag",
			args: []string{"prysm", "--stop-validator"},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "prysm",
				Network:         "mainnet",
				StopValidator:   true,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "start-validator flag",
			args: []string{"teku", "--start-validator"},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "teku",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  true,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "path flag",
			args: []string{"teku", "--path", customDir},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "teku",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  customDir,
				Out:             filepath.Join(customDir, "slashing_protection.json"),
			},
		},
		{
			name: "path shorthand flag",
			args: []string{"teku", "-p", customDir},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "teku",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  customDir,
				Out:             filepath.Join(customDir, "slashing_protection.json"),
			},
		},
		{
			name: "out flag",
			args: []string{"lodestar", "--out", filepath.Join(outDir, "file.json")},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "lodestar",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(outDir, "file.json"),
			},
		},
		{
			name: "out shorthand flag",
			args: []string{"lodestar", "-o", filepath.Join(outDir, "file.json")},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "lodestar",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(outDir, "file.json"),
			},
		},
		{
			name: "nimbus flag",
			args: []string{"nimbus"},
			actionOptions: actions.SlashingExportOptions{
				ValidatorClient: "nimbus",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				Out:             filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockActions := sedge_mocks.NewMockSedgeActions(ctrl)
			depsMgr := sedge_mocks.NewMockDependenciesManager(ctrl)
			gomock.InOrder(
				depsMgr.EXPECT().Check([]string{dependencies.Docker}).Return([]string{dependencies.Docker, dependencies.DockerCompose}, nil).Times(1),
				depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1),
				depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).Times(1),
				mockActions.EXPECT().ValidateDockerComposeFile(filepath.Join(tt.actionOptions.GenerationPath, "docker-compose.yml")).Return(nil).Times(1),
				mockActions.EXPECT().SetupContainers(actions.SetupContainersOptions{
					GenerationPath: tt.actionOptions.GenerationPath,
					Services:       []string{"validator"},
				}),
				mockActions.EXPECT().ExportSlashingInterchangeData(tt.actionOptions).Times(1),
			)

			slashingExportCmd := cli.SlashingExportCmd(mockActions, depsMgr)
			slashingExportCmd.SetArgs(tt.args)
			slashingExportCmd.SetOutput(io.Discard)
			err := slashingExportCmd.Execute()

			assert.Nil(t, err)
		})
	}
}

func TestSlashingExport_Errors(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	tests := []struct {
		name   string
		args   []string
		run    bool
		err    error
		checks bool
	}{
		{
			name:   "invalid network",
			args:   []string{"lighthouse", "--network", "invalid_network"},
			err:    errors.New("invalid network: invalid_network"),
			checks: false,
		},
		{
			name:   "action error",
			args:   []string{"lighthouse"},
			run:    true,
			err:    errors.New("action error"),
			checks: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockActions := sedge_mocks.NewMockSedgeActions(ctrl)
			depsMgr := sedge_mocks.NewMockDependenciesManager(ctrl)
			if tt.checks {
				depsMgr.EXPECT().Check([]string{dependencies.Docker}).Return([]string{dependencies.Docker, dependencies.DockerCompose}, nil).Times(1)
				depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1)
				depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).Times(1)
				mockActions.EXPECT().ValidateDockerComposeFile(filepath.Join(configs.DefaultAbsSedgeDataPath, "docker-compose.yml")).Return(nil).Times(1)
			}
			if tt.run {
				mockActions.EXPECT().SetupContainers(actions.SetupContainersOptions{
					GenerationPath: configs.DefaultAbsSedgeDataPath,
					Services:       []string{"validator"},
				}).Return(nil).Times(1)
				mockActions.EXPECT().ExportSlashingInterchangeData(gomock.Any()).Return(errors.New("action error")).Times(1)
			}

			slashingExportCmd := cli.SlashingExportCmd(mockActions, depsMgr)
			slashingExportCmd.SetArgs(tt.args)
			slashingExportCmd.SetOutput(io.Discard)
			err := slashingExportCmd.Execute()

			assert.EqualError(t, err, tt.err.Error())
		})
	}
}
