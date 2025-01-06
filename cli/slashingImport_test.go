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

func TestSlashingImport_ValidatorIsRequired(t *testing.T) {
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
			args:        []string{"--network", "sepolia", "--from", "slashing0import.json"},
			expectedErr: cli.ErrInvalidNumberOfArguments,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			depsMgr := sedge_mocks.NewMockDependenciesManager(ctrl)
			slashingImportCmd := cli.SlashingImportCmd(nil, depsMgr)
			slashingImportCmd.SetArgs(tt.args)
			slashingImportCmd.SetOutput(io.Discard)
			err := slashingImportCmd.Execute()
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestSlashingImport_Params(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	tests := []struct {
		name          string
		args          []string
		actionOptions actions.SlashingImportOptions
		customDir     bool
		from          bool
	}{
		{
			name: "validator argument",
			args: []string{"lighthouse"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "lighthouse",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				From:            filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "network flag",
			args: []string{"lighthouse", "--network", "sepolia"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "lighthouse",
				Network:         "sepolia",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				From:            filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "network shorthand flag",
			args: []string{"lighthouse", "-n", "sepolia"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "lighthouse",
				Network:         "sepolia",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				From:            filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "stop-validator flag",
			args: []string{"prysm", "--stop-validator"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "prysm",
				Network:         "mainnet",
				StopValidator:   true,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				From:            filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "start-validator flag",
			args: []string{"teku", "--start-validator"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "teku",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  true,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
				From:            filepath.Join(configs.DefaultAbsSedgeDataPath, "slashing_protection.json"),
			},
		},
		{
			name: "path flag",
			args: []string{"teku"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "teku",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
			},
			customDir: true,
		},
		{
			name: "path shorthand flag",
			args: []string{"teku"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "teku",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
			},
			customDir: true,
		},
		{
			name: "from flag",
			args: []string{"lodestar"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "lodestar",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
			},
			from: true,
		},
		{
			name: "from shorthand flag",
			args: []string{"lodestar"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "lodestar",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
			},
			from: true,
		},
		{
			name: "nimbus",
			args: []string{"nimbus"},
			actionOptions: actions.SlashingImportOptions{
				ValidatorClient: "nimbus",
				Network:         "mainnet",
				StopValidator:   false,
				StartValidator:  false,
				GenerationPath:  configs.DefaultAbsSedgeDataPath,
			},
			from: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.customDir {
				customDir := t.TempDir()
				tt.actionOptions.GenerationPath = customDir
				tt.actionOptions.From = filepath.Join(customDir, "slashing_protection.json")
				tt.args = append(tt.args, "--path", customDir)
			}
			if tt.from {
				from := t.TempDir()
				var file *os.File
				var err error

				if file, err = os.Create(filepath.Join(from, "slashing_protection.json")); err != nil {
					t.Fatal(err)
				}
				defer file.Close()
				tt.actionOptions.From = filepath.Join(from, "slashing_protection.json")
				tt.args = append(tt.args, "--from", filepath.Join(from, "slashing_protection.json"))
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			services := []string{"validator"}

			if tt.actionOptions.ValidatorClient == "nimbus" {
				services = append(services, "consensus")
			}

			mockActions := sedge_mocks.NewMockSedgeActions(ctrl)
			depsMgr := sedge_mocks.NewMockDependenciesManager(ctrl)
			gomock.InOrder(
				depsMgr.EXPECT().Check([]string{dependencies.Docker}).Return([]string{dependencies.Docker, dependencies.DockerCompose}, nil).Times(1),
				depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1),
				depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).Times(1),
				mockActions.EXPECT().ValidateDockerComposeFile(filepath.Join(tt.actionOptions.GenerationPath, "docker-compose.yml")).Return(nil).Times(1),
				mockActions.EXPECT().SetupContainers(actions.SetupContainersOptions{
					GenerationPath: tt.actionOptions.GenerationPath,
					Services:       services,
				}).Return(nil).Times(1),
				mockActions.EXPECT().ImportSlashingInterchangeData(tt.actionOptions).Times(1),
			)

			slashingImportCmd := cli.SlashingImportCmd(mockActions, depsMgr)
			slashingImportCmd.SetArgs(tt.args)
			slashingImportCmd.SetOutput(io.Discard)
			err := slashingImportCmd.Execute()

			assert.Nil(t, err)
		})
	}
}

func TestSlashingImport_Errors(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	tests := []struct {
		name string
		args []string
		run  bool
		err  error
	}{
		{
			name: "invalid network",
			args: []string{"lighthouse", "--network", "invalid_network"},
			err:  errors.New("invalid network: invalid_network"),
		},
		{
			name: "action error",
			args: []string{"lighthouse"},
			run:  true,
			err:  errors.New("action error"),
		},
		{
			name: "nimbus invalid network",
			args: []string{"nimbus", "--network", "invalid_network"},
			err:  errors.New("invalid network: invalid_network"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockActions := sedge_mocks.NewMockSedgeActions(ctrl)
			depsMgr := sedge_mocks.NewMockDependenciesManager(ctrl)

			if tt.run {
				gomock.InOrder(
					depsMgr.EXPECT().Check([]string{dependencies.Docker}).Return([]string{dependencies.Docker, dependencies.DockerCompose}, nil).Times(1),
					depsMgr.EXPECT().DockerEngineIsOn().Return(nil).Times(1),
					depsMgr.EXPECT().DockerComposeIsInstalled().Return(nil).Times(1),
					mockActions.EXPECT().ValidateDockerComposeFile(filepath.Join(configs.DefaultAbsSedgeDataPath, configs.DefaultDockerComposeScriptName)).Return(nil).Times(1),
					mockActions.EXPECT().SetupContainers(actions.SetupContainersOptions{
						GenerationPath: configs.DefaultAbsSedgeDataPath,
						Services:       []string{"validator"},
					}).Return(nil).Times(1),
					mockActions.EXPECT().ImportSlashingInterchangeData(gomock.Any()).Return(errors.New("action error")).Times(1),
				)
			}

			slashingImportCmd := cli.SlashingImportCmd(mockActions, depsMgr)
			slashingImportCmd.SetArgs(tt.args)
			slashingImportCmd.SetOutput(io.Discard)
			err := slashingImportCmd.Execute()

			assert.EqualError(t, err, tt.err.Error())
		})
	}
}
