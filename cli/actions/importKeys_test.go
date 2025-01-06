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
	"embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/internal/pkg/services"
	sedge_mocks "github.com/NethermindEth/sedge/mocks"
	"github.com/NethermindEth/sedge/test"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestImportKeys_ValidatorNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := validatorNotFoundHelper(t, ctrl)

	err := s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{})
	assert.ErrorIs(t, err, services.ErrContainerNotFound)
}

func TestImportKeys_CheckValidatorFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantError := errors.New("error")
	s := checkValidatorFailureHelper(t, ctrl, wantError)

	err := s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{})
	assert.ErrorIs(t, err, wantError)
}

func TestImportKeys_ValidatorStopFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := validatorStopFailureHelper(t, ctrl)

	err := s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{})
	assert.ErrorIs(t, err, services.ErrStoppingContainer)
}

func TestImportKeys_ValidatorRunning(t *testing.T) {
	clients := []string{"prysm", "lodestar"}
	networks := []string{"sepolia", "gnosis"}
	for _, validatorClient := range clients {
		for _, network := range networks {
			t.Run(validatorClient, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				dockerClient := importKeysGoldenPath(t, ctrl, false, false)
				dockerServiceManager := services.NewDockerServiceManager(dockerClient)
				cmdRunner := test.SimpleCMDRunner{}
				s := actions.NewSedgeActions(actions.SedgeActionsOptions{
					DockerClient:         dockerClient,
					DockerServiceManager: dockerServiceManager,
					CommandRunner:        &cmdRunner,
				})

				from, err := setupKeystoreDir(t)
				if err != nil {
					t.Fatal(err)
				}

				generationPath := t.TempDir()

				err = s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{
					ValidatorClient: validatorClient,
					Network:         network,
					From:            from,
					GenerationPath:  generationPath,
				})
				assert.NoError(t, err)
			})
		}
	}
}

func TestImportKeysCustom_ValidatorRunning(t *testing.T) {
	clients := []string{"lighthouse", "teku"}
	for _, validatorClient := range clients {
		t.Run(validatorClient, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dockerClient := importKeysGoldenPath(t, ctrl, true, false)
			dockerServiceManager := services.NewDockerServiceManager(dockerClient)
			cmdRunner := test.SimpleCMDRunner{}
			s := actions.NewSedgeActions(actions.SedgeActionsOptions{
				DockerClient:         dockerClient,
				DockerServiceManager: dockerServiceManager,
				CommandRunner:        &cmdRunner,
			})

			from, err := setupKeystoreDir(t)
			if err != nil {
				t.Fatal(err)
			}

			generationPath := t.TempDir()

			err = s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{
				ValidatorClient: validatorClient,
				Network:         "sepolia",
				From:            from,
				GenerationPath:  generationPath,
			})
			assert.NoError(t, err)
		})
	}
}

func TestImportKeys_UnsupportedClient(t *testing.T) {
	clients := []string{"", "unsupported", "kfjkdshjkr24"}
	for _, validatorClient := range clients {
		t.Run(validatorClient, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := unsupportedClientsHelper(t, ctrl)

			from, err := setupKeystoreDir(t)
			if err != nil {
				t.Fatal(err)
			}

			generationPath := t.TempDir()

			err = s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{
				ValidatorClient: validatorClient,
				Network:         "sepolia",
				From:            from,
				GenerationPath:  generationPath,
			})
			assert.ErrorIs(t, err, actions.ErrUnsupportedValidatorClient)
		})
	}
}

func TestImportKeys_CustomOptions(t *testing.T) {
	tests := []struct {
		client      string
		network     string
		customImage bool
	}{}
	for _, validatorClient := range []string{"prysm", "lodestar"} {
		for _, network := range []string{"sepolia", "mainnet"} {
			tests = append(tests, struct {
				client      string
				network     string
				customImage bool
			}{validatorClient, network, false})
		}
	}
	for _, validatorClient := range []string{"lighthouse", "teku"} {
		for _, network := range []string{"sepolia", "mainnet"} {
			tests = append(tests, struct {
				client      string
				network     string
				customImage bool
			}{validatorClient, network, true})
		}
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.client, tt.network), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dockerClient := importKeysGoldenPath(t, ctrl, tt.customImage, false)
			dockerServiceManager := services.NewDockerServiceManager(dockerClient)
			cmdRunner := test.SimpleCMDRunner{}
			s := actions.NewSedgeActions(actions.SedgeActionsOptions{
				DockerClient:         dockerClient,
				DockerServiceManager: dockerServiceManager,
				CommandRunner:        &cmdRunner,
			})

			from, err := setupKeystoreDir(t)
			if err != nil {
				t.Fatal(err)
			}

			generationPath := t.TempDir()
			customConfigPath := t.TempDir()

			customNetworkFile, err := ioutil.TempFile(customConfigPath, "config.yaml")
			if err != nil {
				t.Fatal(err)
			}
			customNetworkConfigPath := filepath.Join(customConfigPath, "config.yaml")
			defer customNetworkFile.Close()

			customGenesisFile, err := ioutil.TempFile(customConfigPath, "genesis.yaml")
			if err != nil {
				t.Fatal(err)
			}
			customGenesisConfigPath := filepath.Join(customConfigPath, "genesis.yaml")
			defer customGenesisFile.Close()

			customDeployBlockFile, err := ioutil.TempFile(customConfigPath, "deployblock")
			if err != nil {
				t.Fatal(err)
			}
			customDeployBlockPath := filepath.Join(customConfigPath, "deployblock")
			defer customDeployBlockFile.Close()

			err = s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{
				CustomConfig: actions.ImportValidatorKeysCustomOptions{
					NetworkConfigPath: customNetworkConfigPath,
					GenesisPath:       customGenesisConfigPath,
					DeployBlockPath:   customDeployBlockPath,
				},
				ValidatorClient: tt.client,
				Network:         tt.network,
				From:            from,
				GenerationPath:  generationPath,
			})
			assert.NoError(t, err)
		})
	}
}

func TestImportKeys_UnexpectedExitCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerClient := importKeysExitError(t, ctrl)
	dockerServiceManager := services.NewDockerServiceManager(dockerClient)
	cmdRunner := test.SimpleCMDRunner{}
	s := actions.NewSedgeActions(actions.SedgeActionsOptions{
		DockerClient:         dockerClient,
		DockerServiceManager: dockerServiceManager,
		CommandRunner:        &cmdRunner,
	})

	from, err := setupKeystoreDir(t)
	if err != nil {
		t.Fatal(err)
	}

	generationPath := t.TempDir()

	err = s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{
		ValidatorClient: "prysm",
		Network:         "sepolia",
		From:            from,
		GenerationPath:  generationPath,
	})
	assert.ErrorIs(t, err, actions.ErrValidatorImportCtBadExitCode)
}

func TestImportKeys_DistributedMode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerClient := importKeysGoldenPath(t, ctrl, false, true)
	serviceManager := services.NewDockerServiceManager(dockerClient)
	cmdRunner := test.SimpleCMDRunner{}
	s := actions.NewSedgeActions(actions.SedgeActionsOptions{
		DockerClient:         dockerClient,
		DockerServiceManager: serviceManager,
		CommandRunner:        &cmdRunner,
	})

	from, err := setupCharonKeystoreDir(t)
	if err != nil {
		t.Logf("Error setting up keystore dir: %v", err)
		t.Fatal(err)
	}

	generationPath := t.TempDir()

	err = s.ImportValidatorKeys(actions.ImportValidatorKeysOptions{
		ValidatorClient: "lodestar",
		Network:         "holesky",
		GenerationPath:  generationPath,
		Distributed:     true,
		From:            from,
	})
	assert.NoError(t, err)
}

//go:embed testdata/keystore
var keystoreTestData embed.FS

func setupKeystoreDir(t *testing.T) (string, error) {
	t.Helper()
	tempKeystore := t.TempDir()

	baseTestDir := "testdata/keystore"
	dirs := []string{""}
	for len(dirs) > 0 {
		currentDir := dirs[0]
		dirEntries, err := keystoreTestData.ReadDir(path.Join(baseTestDir, currentDir))
		if err != nil {
			return "", err
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				dirs = append(dirs, filepath.Join(currentDir, entry.Name()))
			} else {
				entryData, err := keystoreTestData.ReadFile(path.Join(baseTestDir, currentDir, entry.Name()))
				if err != nil {
					return "", err
				}
				if err := os.MkdirAll(filepath.Join(tempKeystore, currentDir), 0o755); err != nil {
					return "", err
				}
				if err := ioutil.WriteFile(filepath.Join(tempKeystore, currentDir, entry.Name()), entryData, 0o755); err != nil {
					return "", err
				}
			}
		}
		dirs = dirs[1:]
	}
	return tempKeystore, nil
}

//go:embed testdata/charon
var charonKeystoreTestData embed.FS

func setupCharonKeystoreDir(t *testing.T) (string, error) {
	t.Helper()
	tempKeystore := t.TempDir()

	baseTestDir := "testdata/charon"
	dirs := []string{""}
	for len(dirs) > 0 {
		currentDir := dirs[0]
		dirEntries, err := charonKeystoreTestData.ReadDir(path.Join(baseTestDir, currentDir))
		if err != nil {
			return "", err
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				dirs = append(dirs, filepath.Join(currentDir, entry.Name()))
			} else {
				entryData, err := charonKeystoreTestData.ReadFile(path.Join(baseTestDir, currentDir, entry.Name()))
				if err != nil {
					return "", err
				}
				if err := os.MkdirAll(filepath.Join(tempKeystore, currentDir), 0o755); err != nil {
					return "", err
				}
				if err := ioutil.WriteFile(filepath.Join(tempKeystore, currentDir, entry.Name()), entryData, 0o755); err != nil {
					return "", err
				}
			}
		}
		dirs = dirs[1:]
	}
	return tempKeystore, nil
}

// importKeysGoldenPath returns a mocked docker client interface with all the
// required responses for a correct validator import keys container execution.
func importKeysGoldenPath(t *testing.T, ctrl *gomock.Controller, withCustomImage bool, withDistributedOption bool) client.APIClient {
	t.Helper()
	dockerClient := sedge_mocks.NewMockAPIClient(ctrl)

	validatorCtId := "validatorctid"
	validatorImportCtId := "validator-import-ct-id"

	// Mock ContainerList
	dockerClient.EXPECT().
		ContainerList(gomock.Any(), container.ListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", services.DefaultSedgeValidatorClient)),
		}).
		Return([]types.Container{
			{
				ID:    validatorCtId,
				Names: []string{"name-0", "/" + services.DefaultSedgeValidatorClient, "name-2"},
			},
		}, nil)
	// Mock ContainerInspect
	inspectCall := dockerClient.EXPECT().
		ContainerInspect(gomock.Any(), services.DefaultSedgeValidatorClient).
		Return(types.ContainerJSON{
			ContainerJSONBase: &types.ContainerJSONBase{
				ID: validatorCtId,
				State: &types.ContainerState{
					Running: true,
				},
			},
		}, nil)
	if withCustomImage {
		inspectCall.Times(2)
	} else if withDistributedOption {
		inspectCall.MinTimes(2)
	} else {
		inspectCall.Times(3)
	}
	// Mock ContainerStop
	dockerClient.EXPECT().
		ContainerStop(gomock.Any(), validatorCtId, gomock.Any()).
		Return(nil)
	// Mock ContainerCreate
	dockerClient.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), services.ServiceCtValidatorImport).
		Return(container.CreateResponse{ID: validatorImportCtId}, nil).
		Times(1)
	// Mock ContainerStart
	dockerClient.EXPECT().
		ContainerStart(gomock.Any(), validatorImportCtId, gomock.Any()).
		Return(nil).
		Times(1)
	dockerClient.EXPECT().
		ContainerStart(gomock.Any(), services.DefaultSedgeValidatorClient, gomock.Any()).
		Return(nil).
		Times(1)
	// Mock ContainerWait
	exitCh := make(chan container.WaitResponse, 1)
	exitCh <- container.WaitResponse{
		StatusCode: 0,
	}
	dockerClient.EXPECT().
		ContainerWait(gomock.Any(), validatorImportCtId, container.WaitConditionNextExit).
		Return(exitCh, make(chan error)).
		Times(1)
	// Mock container logs
	dockerClient.EXPECT().
		ContainerLogs(gomock.Any(), validatorImportCtId, gomock.Any()).
		Return(ioutil.NopCloser(strings.NewReader("logs")), nil).
		Times(1)
	// Mock ContainerRemove
	dockerClient.EXPECT().
		ContainerRemove(gomock.Any(), validatorImportCtId, container.RemoveOptions{}).
		Return(nil).
		Times(1)

	return dockerClient
}

func importKeysExitError(t *testing.T, ctrl *gomock.Controller) client.APIClient {
	t.Helper()
	dockerClient := sedge_mocks.NewMockAPIClient(ctrl)

	validatorCtId := "validatorctid"
	validatorImportCtId := "validator-import-ct-id"

	// Mock ContainerList
	dockerClient.EXPECT().
		ContainerList(gomock.Any(), container.ListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", services.DefaultSedgeValidatorClient)),
		}).
		Return([]types.Container{
			{
				ID:    validatorCtId,
				Names: []string{"/" + services.DefaultSedgeValidatorClient},
			},
		}, nil)
	// Mock ContainerInspect
	dockerClient.EXPECT().
		ContainerInspect(gomock.Any(), services.DefaultSedgeValidatorClient).
		Return(types.ContainerJSON{
			ContainerJSONBase: &types.ContainerJSONBase{
				ID: validatorCtId,
				State: &types.ContainerState{
					Running: true,
				},
			},
		}, nil).
		Times(3)
	// Mock ContainerStop
	dockerClient.EXPECT().
		ContainerStop(gomock.Any(), validatorCtId, gomock.Any()).
		Return(nil)
	// Mock ContainerCreate
	dockerClient.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), services.ServiceCtValidatorImport).
		Return(container.CreateResponse{ID: validatorImportCtId}, nil).
		Times(1)
	// Mock ContainerStart
	dockerClient.EXPECT().
		ContainerStart(gomock.Any(), validatorImportCtId, gomock.Any()).
		Return(nil).
		Times(1)
	dockerClient.EXPECT().
		ContainerStart(gomock.Any(), services.DefaultSedgeValidatorClient, gomock.Any()).
		Return(nil).
		Times(1)
	// Mock ContainerWait
	exitCh := make(chan container.WaitResponse, 1)
	exitCh <- container.WaitResponse{
		StatusCode: 1,
	}
	dockerClient.EXPECT().
		ContainerWait(gomock.Any(), validatorImportCtId, container.WaitConditionNextExit).
		Return(exitCh, make(chan error)).
		Times(1)
	// Mock container logs
	dockerClient.EXPECT().
		ContainerLogs(gomock.Any(), validatorImportCtId, gomock.Any()).
		Return(ioutil.NopCloser(strings.NewReader("logs")), nil).
		Times(1)
	// Mock ContainerRemove
	dockerClient.EXPECT().
		ContainerRemove(gomock.Any(), validatorImportCtId, container.RemoveOptions{}).
		Return(nil).
		Times(1)

	return dockerClient
}
