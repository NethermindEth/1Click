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
package monitoring

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/NethermindEth/sedge/internal/common"
	"github.com/NethermindEth/sedge/internal/monitoring/data"
	mock_locker "github.com/NethermindEth/sedge/internal/monitoring/locker/mocks"
	mocks "github.com/NethermindEth/sedge/internal/monitoring/mocks"
	"github.com/NethermindEth/sedge/internal/monitoring/services/templates"
	"github.com/NethermindEth/sedge/internal/monitoring/services/types"
	"github.com/NethermindEth/sedge/internal/monitoring/utils"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	userDataHome := os.Getenv("XDG_DATA_HOME")
	if userDataHome == "" {
		userHome, err := os.UserHomeDir()
		require.NoError(t, err)
		userDataHome = filepath.Join(userHome, ".local", "share")
	}

	okLocker := func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
		// Create a mock locker
		locker := mock_locker.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		gomock.InOrder(
			locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		return locker
	}

	dotEnvFileWriter := func(t *testing.T, fs afero.Fs, dotenv map[string]string) {
		// Create the dotenv file
		dotenvFile, err := fs.Create(filepath.Join(userDataHome, ".sedge", "monitoring", ".env"))
		require.NoError(t, err)

		// Write the dotenv file
		for key, value := range dotenv {
			_, err := dotenvFile.WriteString(key + "=" + value + "\n")
			require.NoError(t, err)
		}
	}

	tests := []struct {
		name        string
		mocker      func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack, dotenv map[string]string) ([]ServiceAPI, *mocks.MockDockerServiceManager)
		setupDotEnv func(t *testing.T, fs afero.Fs, dotenv map[string]string)
		dotenv      map[string]string
		wantErr     bool
	}{
		{
			name: "ok, 1 service",
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack, dotenv map[string]string) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().ContainerName().Return("node"),
					servicer.EXPECT().SetContainerIP(net.ParseIP("127.0.0.1")).Return(),
				)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node").Return("127.0.0.1", nil)

				return []ServiceAPI{
					servicer,
				}, dockerServiceManager
			},
			setupDotEnv: dotEnvFileWriter,
			dotenv: map[string]string{
				"NODE_PORT": "9000",
			},
		},
		{
			name: "ok, 2 services",
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack, dotenv map[string]string) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service1 := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					service1.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					service1.EXPECT().ContainerName().Return("node1"),
				)

				service2 := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					service2.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					service2.EXPECT().ContainerName().Return("node2"),
				)

				service1.EXPECT().SetContainerIP(net.ParseIP("127.0.0.1")).Return()
				service2.EXPECT().SetContainerIP(net.ParseIP("127.0.0.2")).Return()

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node1").Return("127.0.0.1", nil)
				dockerServiceManager.EXPECT().ContainerIP("node2").Return("127.0.0.2", nil)

				return []ServiceAPI{
					service1, service2,
				}, dockerServiceManager
			},
			setupDotEnv: dotEnvFileWriter,
			dotenv: map[string]string{
				"NODE1_PORT": "9000",
				"NODE2_PORT": "9001",
			},
		},
		{
			name: "error, 1 service, init service error",
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack, dotenv map[string]string) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(errors.New("error")),
				)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, dockerServiceManager
			},
			setupDotEnv: dotEnvFileWriter,
			dotenv: map[string]string{
				"NODE_PORT": "9000",
			},
			wantErr: true,
		},
		{
			name: "error, 1 service, ContainerIP error",
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack, dotenv map[string]string) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().ContainerName().Return("node"),
				)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node").Return("", errors.New("error"))

				return []ServiceAPI{
					servicer,
				}, dockerServiceManager
			},
			setupDotEnv: dotEnvFileWriter,
			dotenv: map[string]string{
				"NODE_PORT": "9000",
			},
			wantErr: true,
		},
		{
			name: "error, 1 service, ContainerIP gives an invalid IP",
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack, dotenv map[string]string) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().ContainerName().Return("node"),
				)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node").Return("nethermind-loves-sedge", nil)

				return []ServiceAPI{
					servicer,
				}, dockerServiceManager
			},
			setupDotEnv: dotEnvFileWriter,
			dotenv: map[string]string{
				"NODE_PORT": "9000",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			// Create a mock controller
			ctrl := gomock.NewController(t)
			locker := okLocker(t, ctrl)

			afs := afero.NewMemMapFs()
			tt.setupDotEnv(t, afs, tt.dotenv)

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{},
				mocks.NewMockComposeManager(ctrl),
				mocks.NewMockDockerServiceManager(ctrl),
				afs,
				locker,
			)

			services, dockerServiceManager := tt.mocker(t, ctrl, manager.stack, tt.dotenv)
			manager.services = services
			manager.dockerServiceManager = dockerServiceManager

			// Init the stack
			err := manager.Init()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestInstallStack(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	userDataHome := os.Getenv("XDG_DATA_HOME")
	if userDataHome == "" {
		userHome, err := os.UserHomeDir()
		require.NoError(t, err)
		userDataHome = filepath.Join(userHome, ".local", "share")
	}

	okLocker := func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
		// Create a mock locker
		locker := mock_locker.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		gomock.InOrder(
			locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		// stack.Installed() lock
		gomock.InOrder(
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		return locker
	}
	onlyNewLocker := func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
		// Create a mock locker
		locker := mock_locker.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker)
		return locker
	}

	// Setup mock http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Listen for the POST request to /-/reload
		if r.URL.Path == "/-/reload" && r.Method == http.MethodPost {
			// All good
			w.WriteHeader(http.StatusOK)
		} else if r.Method != http.MethodGet {
			// Unexpected method
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			// Unexpected path
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	serverPort := strings.Split(server.URL, ":")[2]

	tests := []struct {
		name         string
		mockerLocker func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker
		mocker       func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager)
		wantErr      bool
	}{
		{
			name:         "ok, 1 service, port not occupied",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "9000",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(dotenv),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().Setup(dotenv).Return(nil),
					servicer.EXPECT().ContainerName().Return("node"),
					servicer.EXPECT().SetContainerIP(net.ParseIP("127.0.0.1")).Return(),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Create(commands.DockerComposeCreateOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node").Return("127.0.0.1", nil)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
		},
		{
			name:         "ok, 2 services",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE1_PORT": "9000",
					"NODE2_PORT": "9003",
				}
				service1, service2 := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)

				// Expect the service to be triggered
				gomock.InOrder(
					service1.EXPECT().DotEnv().Return(map[string]string{
						"NODE1_PORT": "9000",
					}),
					service1.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					service1.EXPECT().Setup(dotenv).Return(nil),
					service1.EXPECT().ContainerName().Return("node1"),
				)
				gomock.InOrder(
					service2.EXPECT().DotEnv().Return(map[string]string{
						"NODE2_PORT": "9003",
					}),
					service2.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					service2.EXPECT().Setup(dotenv).Return(nil),
					service2.EXPECT().ContainerName().Return("node2"),
				)
				service1.EXPECT().SetContainerIP(net.ParseIP("168.0.2.1")).Return()
				service2.EXPECT().SetContainerIP(net.ParseIP("168.0.3.1")).Return()

				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Create(commands.DockerComposeCreateOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node1").Return("168.0.2.1", nil)
				dockerServiceManager.EXPECT().ContainerIP("node2").Return("168.0.3.1", nil)

				return []ServiceAPI{
					service1,
					service2,
				}, composeManager, dockerServiceManager
			},
		},
		{
			name:         "ok, 1 service, port occupied",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				// Convert serverPort to int
				p, err := strconv.Atoi(serverPort)
				require.NoError(t, err)
				sp := strconv.Itoa(p + 1)

				dotenv := map[string]string{
					"NODE_PORT": sp,
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(map[string]string{"NODE_PORT": serverPort}),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().Setup(dotenv).Return(nil),
					servicer.EXPECT().ContainerName().Return("node"),
					servicer.EXPECT().SetContainerIP(net.ParseIP("127.1.1.6")).Return(),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Create(commands.DockerComposeCreateOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node").Return("127.1.1.6", nil)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
		},
		{
			name:         "error, 1 service, port not int",
			mockerLocker: onlyNewLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "3RR0R",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				servicer.EXPECT().DotEnv().Return(dotenv)

				composeManager := mocks.NewMockComposeManager(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
		{
			name:         "error, 1 service, invalid port",
			mockerLocker: onlyNewLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "0",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				servicer.EXPECT().DotEnv().Return(dotenv)

				composeManager := mocks.NewMockComposeManager(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
		{
			name:         "error, 1 service, init service error",
			mockerLocker: onlyNewLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "9000",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(dotenv),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(errors.New("error")),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
		{
			name: "error, 1 service, stack setup error",
			mockerLocker: func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(errors.New("error")),
				)
				return locker
			},
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "9000",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(dotenv),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
		{
			name: "error, 1 service, service setup error",
			mockerLocker: func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)
				return locker
			},
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "9000",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(dotenv),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().Setup(dotenv).Return(errors.New("error")),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
		{
			name: "error, 1 service, create error",
			mockerLocker: func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)
				return locker
			},
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "9000",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(dotenv),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().Setup(dotenv).Return(nil),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Create(commands.DockerComposeCreateOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(errors.New("error"))

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
		{
			name: "error, 1 service, run error",
			mockerLocker: func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)
				return locker
			},
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "9000",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(dotenv),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().Setup(dotenv).Return(nil),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Create(commands.DockerComposeCreateOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(errors.New("error"))

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
		{
			name: "error, 1 service, ContainerIP error",
			mockerLocker: func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)
				return locker
			},
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				dotenv := map[string]string{
					"NODE_PORT": "9000",
				}
				servicer := mocks.NewMockServiceAPI(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().DotEnv().Return(dotenv),
					servicer.EXPECT().Init(types.ServiceOptions{
						Stack:  stack,
						Dotenv: dotenv,
					}).Return(nil),
					servicer.EXPECT().Setup(dotenv).Return(nil),
					servicer.EXPECT().ContainerName().Return("node"),
				)

				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Create(commands.DockerComposeCreateOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node").Return("", errors.New("error"))

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			// Create a mock controller
			ctrl := gomock.NewController(t)
			locker := tt.mockerLocker(t, ctrl)

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{},
				mocks.NewMockComposeManager(ctrl),
				mocks.NewMockDockerServiceManager(ctrl),
				afero.NewMemMapFs(),
				locker,
			)

			services, composeManager, dockerServiceManager := tt.mocker(t, ctrl, manager.stack)
			manager.services = services
			manager.composeManager = composeManager
			manager.dockerServiceManager = dockerServiceManager

			// Init the stack
			err := manager.InstallStack()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Check the stack is installed
				installed, err := manager.stack.Installed()
				assert.NoError(t, err)
				assert.True(t, installed)
			}
		})
	}
}

func TestAddAndRemoveTarget(t *testing.T) {
	okLocker := func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker {
		// Create a mock locker
		locker := mock_locker.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		userDataHome := os.Getenv("XDG_DATA_HOME")
		if userDataHome == "" {
			userHome, err := os.UserHomeDir()
			require.NoError(t, err)
			userDataHome = filepath.Join(userHome, ".local", "share")
		}
		locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker)

		return locker
	}

	tests := []struct {
		name          string
		mockerLocker  func(t *testing.T, ctrl *gomock.Controller) *mock_locker.MockLocker
		mocker        func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager)
		target        types.MonitoringTarget
		labels        map[string]string
		dockerNetwork string
		add           bool
		wantErr       bool
	}{
		{
			name:         "add, ok, prometheus service",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				prometheusService := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				gomock.InOrder(
					prometheusService.EXPECT().ContainerName().Return(PrometheusContainerName),
					dockerServiceManager.EXPECT().ContainerNetworks(PrometheusContainerName).Return([]string{"sedge_default"}, nil),
					dockerServiceManager.EXPECT().NetworkConnect(PrometheusContainerName, dockerNetwork).Return(nil),
					prometheusService.EXPECT().AddTarget(target, labels, fmt.Sprintf("%s--%s++%s", labels[InstanceIDLabel], PrometheusContainerName, dockerNetwork)).Return(nil),
				)

				return []ServiceAPI{
					prometheusService,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
				Port: 9000,
			},
			labels: map[string]string{
				InstanceIDLabel: "becks",
			},
			add: true,
		},
		{
			name:         "add, ok, 1 service",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				gomock.InOrder(
					service.EXPECT().ContainerName().Return("service1"),
					service.EXPECT().AddTarget(target, labels, fmt.Sprintf("%s--service1++%s", labels[InstanceIDLabel], dockerNetwork)).Return(nil),
				)

				return []ServiceAPI{
					service,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
				Port: 9000,
				Path: "/custom/path",
			},
			labels: map[string]string{
				InstanceIDLabel: "heineken",
			},
			add: true,
		},
		{
			name:         "add, ok, 1 services, prometheus was already added to network",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				prometheusService := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				gomock.InOrder(
					prometheusService.EXPECT().ContainerName().Return(PrometheusContainerName),
					dockerServiceManager.EXPECT().ContainerNetworks(PrometheusContainerName).Return([]string{"sedge_default", dockerNetwork}, nil),
					prometheusService.EXPECT().AddTarget(target, labels, fmt.Sprintf("%s--%s++%s", labels[InstanceIDLabel], PrometheusContainerName, dockerNetwork)).Return(nil),
				)

				return []ServiceAPI{
					prometheusService,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
				Port: 9000,
			},
			labels: map[string]string{
				InstanceIDLabel: "grolsch",
			},
			add: true,
		},
		{
			name:         "add, ok, 2 services, 1 AddTarget error",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				prometheusService, service2 := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				gomock.InOrder(
					prometheusService.EXPECT().ContainerName().Return(PrometheusContainerName),
					dockerServiceManager.EXPECT().ContainerNetworks(PrometheusContainerName).Return([]string{"sedge_default"}, nil),
					dockerServiceManager.EXPECT().NetworkConnect(PrometheusContainerName, dockerNetwork).Return(nil),
					prometheusService.EXPECT().AddTarget(target, labels, fmt.Sprintf("%s--%s++%s", labels[InstanceIDLabel], PrometheusContainerName, dockerNetwork)).Return(nil),
					service2.EXPECT().ContainerName().Return("node2"),
					service2.EXPECT().AddTarget(target, labels, fmt.Sprintf("%s--node2++%s", labels[InstanceIDLabel], dockerNetwork)).Return(errors.New("error")),
				)

				return []ServiceAPI{
					prometheusService,
					service2,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
				Port: 9000,
			},
			labels: map[string]string{
				InstanceIDLabel: "carlsberg",
			},
			wantErr: true,
			add:     true,
		},
		{
			name:         "add, ok, 2 services, prometheus ContainerNetworks error",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				prometheusService, service2 := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				gomock.InOrder(
					service2.EXPECT().ContainerName().Return("node2"),
					service2.EXPECT().AddTarget(target, labels, fmt.Sprintf("%s--node2++%s", labels[InstanceIDLabel], dockerNetwork)).Return(nil),
					prometheusService.EXPECT().ContainerName().Return(PrometheusContainerName),
					dockerServiceManager.EXPECT().ContainerNetworks(PrometheusContainerName).Return(nil, errors.New("error")),
				)

				return []ServiceAPI{
					service2,
					prometheusService,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
				Port: 9000,
			},
			labels: map[string]string{
				InstanceIDLabel: "amstel",
			},
			wantErr: true,
			add:     true,
		},
		{
			name:         "add, ok, 2 services, prometheus NetworkConnect error",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service1, prometheusService := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				gomock.InOrder(
					service1.EXPECT().ContainerName().Return("node1"),
					service1.EXPECT().AddTarget(target, labels, fmt.Sprintf("%s--node1++%s", labels[InstanceIDLabel], dockerNetwork)).Return(nil),
					prometheusService.EXPECT().ContainerName().Return(PrometheusContainerName),
					dockerServiceManager.EXPECT().ContainerNetworks(PrometheusContainerName).Return([]string{"sedge_default"}, nil),
					dockerServiceManager.EXPECT().NetworkConnect(PrometheusContainerName, dockerNetwork).Return(errors.New("error")),
				)

				return []ServiceAPI{
					service1,
					prometheusService,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
				Port: 9000,
			},
			labels: map[string]string{
				InstanceIDLabel: "budweiser",
			},
			wantErr: true,
			add:     true,
		},
		{
			name:         "remove, ok, 1 service",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				servicer := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					servicer.EXPECT().RemoveTarget(labels[InstanceIDLabel]).Return(dockerNetwork, nil),
					servicer.EXPECT().ContainerName().Return("node"),
					dockerServiceManager.EXPECT().NetworkDisconnect("node", dockerNetwork).Return(nil),
				)

				return []ServiceAPI{
					servicer,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
			},
			labels: map[string]string{
				InstanceIDLabel: "stella",
			},
		},
		{
			name:         "remove, ok, 2 services, one of them was already removed from network",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service1, service2 := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					service1.EXPECT().RemoveTarget(labels[InstanceIDLabel]).Return(dockerNetwork, nil),
					service1.EXPECT().ContainerName().Return("node1"),
					dockerServiceManager.EXPECT().NetworkDisconnect("node1", dockerNetwork).Return(nil),
					service2.EXPECT().RemoveTarget(labels[InstanceIDLabel]).Return(dockerNetwork, nil),
					service2.EXPECT().ContainerName().Return("node2"),
					dockerServiceManager.EXPECT().NetworkDisconnect("node2", dockerNetwork).Return(assert.AnError),
				)

				return []ServiceAPI{
					service1,
					service2,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "localhost",
			},
			labels: map[string]string{
				InstanceIDLabel: "corona",
			},
		},
		{
			name:         "remove, ok, 2 services, 1 RemoveTarget error",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service1, service2 := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					service1.EXPECT().RemoveTarget(labels[InstanceIDLabel]).Return(dockerNetwork, nil),
					service1.EXPECT().ContainerName().Return("node1"),
					dockerServiceManager.EXPECT().NetworkDisconnect("node1", dockerNetwork).Return(nil),
					service2.EXPECT().RemoveTarget(labels[InstanceIDLabel]).Return("", assert.AnError),
				)

				return []ServiceAPI{
					service1,
					service2,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "remove",
			},
			labels: map[string]string{
				InstanceIDLabel: "perla",
			},
			wantErr: true,
		},
		{
			name:         "remove, ok, 2 services, 1 NetworkDisconnect error",
			mockerLocker: okLocker,
			mocker: func(t *testing.T, ctrl *gomock.Controller, labels map[string]string, dockerNetwork string, target types.MonitoringTarget) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service1, service2 := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				// Expect the service to be triggered
				gomock.InOrder(
					service1.EXPECT().RemoveTarget(labels[InstanceIDLabel]).Return(dockerNetwork, nil),
					service1.EXPECT().ContainerName().Return("node1"),
					dockerServiceManager.EXPECT().NetworkDisconnect("node1", dockerNetwork).Return(nil),
					service2.EXPECT().RemoveTarget(labels[InstanceIDLabel]).Return(dockerNetwork, nil),
					service2.EXPECT().ContainerName().Return("node2"),
					dockerServiceManager.EXPECT().NetworkDisconnect("node2", dockerNetwork).Return(assert.AnError),
				)

				return []ServiceAPI{
					service1,
					service2,
				}, dockerServiceManager
			},
			target: types.MonitoringTarget{
				Host: "remove",
			},
			labels: map[string]string{
				InstanceIDLabel: "yeungling",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			locker := tt.mockerLocker(t, ctrl)

			services, dockerServiceManager := tt.mocker(t, ctrl, tt.labels, tt.dockerNetwork, tt.target)

			// Create a monitoring manager
			manager := NewMonitoringManager(
				services,
				mocks.NewMockComposeManager(ctrl),
				dockerServiceManager,
				afero.NewMemMapFs(),
				locker,
			)

			var err error
			if tt.add {
				// Add the target
				err = manager.AddTarget(tt.target, tt.labels, tt.dockerNetwork)
			} else {
				// Remove the target
				err = manager.RemoveTarget(tt.labels[InstanceIDLabel])
			}
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRun(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	userDataHome := os.Getenv("XDG_DATA_HOME")
	if userDataHome == "" {
		userHome, err := os.UserHomeDir()
		require.NoError(t, err)
		userDataHome = filepath.Join(userHome, ".local", "share")
	}
	composePath := filepath.Join(userDataHome, ".sedge", "monitoring", "docker-compose.yml")

	tests := []struct {
		name    string
		mocker  func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager)
		wantErr bool
	}{
		{
			name: "ok",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				service1, service2 := mocks.NewMockServiceAPI(ctrl), mocks.NewMockServiceAPI(ctrl)
				service1.EXPECT().ContainerName().Return("node1")
				service2.EXPECT().ContainerName().Return("node2")

				service1.EXPECT().SetContainerIP(net.ParseIP("168.0.2.1")).Return()
				service2.EXPECT().SetContainerIP(net.ParseIP("168.0.3.1")).Return()

				composeManager := mocks.NewMockComposeManager(ctrl)
				// Expect the compose manager to be triggered
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: composePath}).Return(nil)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node1").Return("168.0.2.1", nil)
				dockerServiceManager.EXPECT().ContainerIP("node2").Return("168.0.3.1", nil)

				return []ServiceAPI{
					service1, service2,
				}, composeManager, dockerServiceManager
			},
		},
		{
			name: "up error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				composeManager := mocks.NewMockComposeManager(ctrl)
				// Expect the compose manager to be triggered
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: composePath}).Return(errors.New("error"))
				return []ServiceAPI{mocks.NewMockServiceAPI(ctrl)}, composeManager, mocks.NewMockDockerServiceManager(ctrl)
			},
			wantErr: true,
		},
		{
			name: "error, ContainerIP error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockComposeManager, *mocks.MockDockerServiceManager) {
				servicer := mocks.NewMockServiceAPI(ctrl)
				servicer.EXPECT().ContainerName().Return("node")

				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: composePath}).Return(nil)

				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP("node").Return("", errors.New("error"))

				return []ServiceAPI{
					servicer,
				}, composeManager, dockerServiceManager
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory filesystem
			afs := afero.NewMemMapFs()

			// Create a mock controller
			ctrl := gomock.NewController(t)

			// Create a mock locker
			locker := mock_locker.NewMockLocker(ctrl)
			// Expect the lock to be acquired
			locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker)

			services, composeManager, dockerServiceManager := tt.mocker(t, ctrl)

			// Create a monitoring manager
			manager := NewMonitoringManager(
				services,
				composeManager,
				dockerServiceManager,
				afs,
				locker,
			)

			// Run the stack
			err := manager.Run()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStop(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	userDataHome := os.Getenv("XDG_DATA_HOME")
	if userDataHome == "" {
		userHome, err := os.UserHomeDir()
		require.NoError(t, err)
		userDataHome = filepath.Join(userHome, ".local", "share")
	}
	composePath := filepath.Join(userDataHome, ".sedge", "monitoring", "docker-compose.yml")

	tests := []struct {
		name    string
		mocker  func(t *testing.T, ctrl *gomock.Controller) *mocks.MockComposeManager
		wantErr bool
	}{
		{
			name: "ok",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *mocks.MockComposeManager {
				composeManager := mocks.NewMockComposeManager(ctrl)
				// Expect the compose manager to be triggered
				gomock.InOrder(
					composeManager.EXPECT().Down(commands.DockerComposeDownOptions{Path: composePath}).Return(nil),
				)
				return composeManager
			},
		},
		{
			name: "down error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *mocks.MockComposeManager {
				composeManager := mocks.NewMockComposeManager(ctrl)
				// Expect the compose manager to be triggered
				composeManager.EXPECT().Down(commands.DockerComposeDownOptions{Path: composePath}).Return(errors.New("error"))
				return composeManager
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)

			// Create a mock locker
			locker := mock_locker.NewMockLocker(ctrl)
			// Expect the lock to be acquired
			locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker)

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{mocks.NewMockServiceAPI(ctrl)},
				tt.mocker(t, ctrl),
				mocks.NewMockDockerServiceManager(ctrl),
				afero.NewMemMapFs(),
				locker,
			)

			// Stop the stack
			err := manager.Stop()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	userDataHome := os.Getenv("XDG_DATA_HOME")
	if userDataHome == "" {
		userHome, err := os.UserHomeDir()
		require.NoError(t, err)
		userDataHome = filepath.Join(userHome, ".local", "share")
	}

	tests := []struct {
		name    string
		mocker  func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockDockerServiceManager)
		want    common.Status
		wantErr bool
	}{
		{
			name: "ok",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				services := []ServiceAPI{
					mocks.NewMockServiceAPI(ctrl),
					mocks.NewMockServiceAPI(ctrl),
					mocks.NewMockServiceAPI(ctrl),
				}
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				for i, service := range services {
					mockService := service.(*mocks.MockServiceAPI)
					containerName := fmt.Sprintf("service%d", i+1)
					mockService.EXPECT().ContainerName().Return(containerName)
					dockerServiceManager.EXPECT().ContainerStatus(containerName).Return(common.Running, nil)
				}

				return services, dockerServiceManager
			},
			want: common.Running,
		},
		{
			name: "error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				service.EXPECT().ContainerName().Return("service1")
				dockerServiceManager.EXPECT().ContainerStatus("service1").Return(common.Unknown, errors.New("error"))

				return []ServiceAPI{service}, dockerServiceManager
			},
			want:    common.Unknown,
			wantErr: true,
		},
		{
			name: "Restarting",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				services := []ServiceAPI{
					mocks.NewMockServiceAPI(ctrl),
					mocks.NewMockServiceAPI(ctrl),
				}
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				for i, service := range services {
					mockService := service.(*mocks.MockServiceAPI)
					containerName := fmt.Sprintf("service%d", i+1)
					mockService.EXPECT().ContainerName().Return(containerName)
					dockerServiceManager.EXPECT().ContainerStatus(containerName).Return(common.Restarting, nil)
				}

				return services, dockerServiceManager
			},
			want: common.Restarting,
		},
		{
			name: "Paused",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				service.EXPECT().ContainerName().Return("service1")
				dockerServiceManager.EXPECT().ContainerStatus("service1").Return(common.Paused, nil)

				return []ServiceAPI{service}, dockerServiceManager
			},
			want:    common.Broken,
			wantErr: true,
		},
		{
			name: "Exited",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				service.EXPECT().ContainerName().Return("service1")
				dockerServiceManager.EXPECT().ContainerStatus("service1").Return(common.Exited, nil)

				return []ServiceAPI{service}, dockerServiceManager
			},
			want:    common.Broken,
			wantErr: true,
		},
		{
			name: "Dead",
			mocker: func(t *testing.T, ctrl *gomock.Controller) ([]ServiceAPI, *mocks.MockDockerServiceManager) {
				service := mocks.NewMockServiceAPI(ctrl)
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)

				service.EXPECT().ContainerName().Return("service1")
				dockerServiceManager.EXPECT().ContainerStatus("service1").Return(common.Dead, nil)

				return []ServiceAPI{service}, dockerServiceManager
			},
			want:    common.Broken,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			locker := mock_locker.NewMockLocker(ctrl)
			locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker)

			services, dockerServiceManager := tt.mocker(t, ctrl)
			manager := NewMonitoringManager(
				services,
				mocks.NewMockComposeManager(ctrl),
				dockerServiceManager,
				afero.NewMemMapFs(),
				locker,
			)

			status, err := manager.Status()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, status)
		})
	}
}

func TestInstallationStatus(t *testing.T) {
	userDataHome := os.Getenv("XDG_DATA_HOME")
	if userDataHome == "" {
		userHome, err := os.UserHomeDir()
		require.NoError(t, err)
		userDataHome = filepath.Join(userHome, ".local", "share")
	}

	tests := []struct {
		name    string
		mocker  func(t *testing.T, ctrl *gomock.Controller) (afero.Fs, *mock_locker.MockLocker)
		want    common.Status
		wantErr bool
	}{
		{
			name: "installed",
			mocker: func(t *testing.T, ctrl *gomock.Controller) (afero.Fs, *mock_locker.MockLocker) {
				fs := afero.NewMemMapFs()
				// Recreate installed monitoring
				err := fs.MkdirAll(filepath.Join(userDataHome, ".sedge", "monitoring"), 0o755)
				require.NoError(t, err)
				_, err = fs.Create(filepath.Join(userDataHome, ".sedge", "monitoring", "docker-compose.yml"))
				require.NoError(t, err)
				_, err = fs.Create(filepath.Join(userDataHome, ".sedge", "monitoring", ".env"))
				require.NoError(t, err)

				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)
				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)

				return fs, locker
			},
			want: common.Installed,
		},
		{
			name: "not installed",
			mocker: func(t *testing.T, ctrl *gomock.Controller) (afero.Fs, *mock_locker.MockLocker) {
				fs := afero.NewMemMapFs()

				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)
				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)

				return fs, locker
			},
			want: common.NotInstalled,
		},
		{
			name: "lock error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) (afero.Fs, *mock_locker.MockLocker) {
				fs := afero.NewMemMapFs()

				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)
				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(errors.New("error")),
				)

				return fs, locker
			},
			want:    common.Unknown,
			wantErr: true,
		},
		{
			name: "unlock error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) (afero.Fs, *mock_locker.MockLocker) {
				fs := afero.NewMemMapFs()

				// Create a mock locker
				locker := mock_locker.NewMockLocker(ctrl)
				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(errors.New("error")),
				)

				return fs, locker
			},
			want:    common.Unknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)

			fs, locker := tt.mocker(t, ctrl)

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{mocks.NewMockServiceAPI(ctrl)},
				mocks.NewMockComposeManager(ctrl),
				mocks.NewMockDockerServiceManager(ctrl),
				fs,
				locker,
			)

			status, err := manager.InstallationStatus()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, status)
		})
	}
}

func TestCleanup(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	userDataHome := os.Getenv("XDG_DATA_HOME")
	if userDataHome == "" {
		userHome, err := os.UserHomeDir()
		require.NoError(t, err)
		userDataHome = filepath.Join(userHome, ".local", "share")
	}
	composePath := filepath.Join(userDataHome, ".sedge", "monitoring", "docker-compose.yml")

	tests := []struct {
		name      string
		mocker    func(t *testing.T, ctrl *gomock.Controller) (*mocks.MockComposeManager, *mock_locker.MockLocker)
		noInstall bool
		wantErr   bool
	}{
		{
			name: "ok",
			mocker: func(t *testing.T, ctrl *gomock.Controller) (*mocks.MockComposeManager, *mock_locker.MockLocker) {
				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Down(commands.DockerComposeDownOptions{Path: composePath}).Return(nil)

				locker := mock_locker.NewMockLocker(ctrl)
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)

				return composeManager, locker
			},
		},
		{
			name: "down error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) (*mocks.MockComposeManager, *mock_locker.MockLocker) {
				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Down(commands.DockerComposeDownOptions{Path: composePath}).Return(errors.New("error"))

				locker := mock_locker.NewMockLocker(ctrl)
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)

				return composeManager, locker
			},
			wantErr: true,
		},
		{
			name: "ok, no install",
			mocker: func(t *testing.T, ctrl *gomock.Controller) (*mocks.MockComposeManager, *mock_locker.MockLocker) {
				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Down(commands.DockerComposeDownOptions{Path: composePath}).Return(nil)

				locker := mock_locker.NewMockLocker(ctrl)
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker),
				)

				return composeManager, locker
			},
			noInstall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory filesystem
			afs := afero.NewMemMapFs()

			// Create a mock controller
			ctrl := gomock.NewController(t)

			composeMgr, locker := tt.mocker(t, ctrl)

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{mocks.NewMockServiceAPI(ctrl)},
				composeMgr,
				mocks.NewMockDockerServiceManager(ctrl),
				afs,
				locker,
			)

			if !tt.noInstall {
				err := manager.stack.Setup(map[string]string{"NODE_NAME": "test"}, templates.Services)
				require.NoError(t, err)
			}

			err := manager.Cleanup()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Check that monitoring stack has been removed
				exists, err := afero.DirExists(afs, filepath.Join(userDataHome, ".sedge", "monitoring"))
				assert.NoError(t, err)
				assert.False(t, exists)
			}
		})
	}
}

func TestServiceEndpoints(t *testing.T) {
	want := map[string]string{
		GrafanaContainerName:      "http://grafana:3005",
		PrometheusContainerName:   "http://prometheus:9095",
		NodeExporterContainerName: "http://node-exporter:9105",
	}

	// Setup mocks
	ctrl := gomock.NewController(t)
	grafanaMock := mocks.NewMockServiceAPI(ctrl)
	promMock := mocks.NewMockServiceAPI(ctrl)
	nodeExporterMock := mocks.NewMockServiceAPI(ctrl)

	// Expect the service to be triggered
	gomock.InOrder(
		grafanaMock.EXPECT().ContainerName().Return(GrafanaContainerName),
		grafanaMock.EXPECT().Endpoint().Return("http://grafana:3005"),
		promMock.EXPECT().ContainerName().Return(PrometheusContainerName),
		promMock.EXPECT().Endpoint().Return("http://prometheus:9095"),
		nodeExporterMock.EXPECT().ContainerName().Return(NodeExporterContainerName),
		nodeExporterMock.EXPECT().Endpoint().Return("http://node-exporter:9105"),
	)

	// Init monitoring manager and services
	manager := MonitoringManager{
		services: []ServiceAPI{grafanaMock, promMock, nodeExporterMock},
	}

	// Check endpoints
	endpoints := manager.ServiceEndpoints()
	assert.Equal(t, want, endpoints)
}

func TestAddService(t *testing.T) {
	tests := []struct {
		name          string
		mocker        func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) (ServiceAPI, ComposeManager, DockerServiceManager)
		expectedError bool
	}{
		{
			name: "Add Lido exporter service (valid case)",
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) (ServiceAPI, ComposeManager, DockerServiceManager) {
				// service expectations
				service := mocks.NewMockServiceAPI(ctrl)
				service.EXPECT().ContainerName().Return("sedge_lido_exporter").AnyTimes()
				service.EXPECT().DotEnv().Return(map[string]string{"LIDO_EXPORTER_PORT": "8080"}).AnyTimes()
				service.EXPECT().Init(gomock.Any()).Return(nil).AnyTimes()
				service.EXPECT().Setup(gomock.Any()).Return(nil).AnyTimes()
				service.EXPECT().Name().Return("lido_exporter").AnyTimes()
				service.EXPECT().Endpoint().Return("http://sedge_lido_exporter:8080").AnyTimes()
				service.EXPECT().SetContainerIP(gomock.Any()).AnyTimes()
				service.EXPECT().AddTarget(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

				// dockerServiceManager
				dockerServiceManager := mocks.NewMockDockerServiceManager(ctrl)
				dockerServiceManager.EXPECT().ContainerIP(GrafanaContainerName).Return("168.0.2.1", nil).AnyTimes()
				dockerServiceManager.EXPECT().ContainerIP(PrometheusContainerName).Return("168.0.2.2", nil).AnyTimes()
				dockerServiceManager.EXPECT().ContainerIP(NodeExporterContainerName).Return("168.0.2.3", nil).AnyTimes()
				dockerServiceManager.EXPECT().ContainerIP("sedge_lido_exporter").Return("168.0.2.4", nil).AnyTimes()
				dockerServiceManager.EXPECT().ContainerNetworks(PrometheusContainerName).Return([]string{"sedge_default", SedgeNetworkName}, nil)

				// composeManager
				composeManager := mocks.NewMockComposeManager(ctrl)
				composeManager.EXPECT().Create(commands.DockerComposeCreateOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)
				composeManager.EXPECT().Up(commands.DockerComposeUpOptions{Path: filepath.Join(stack.Path(), "docker-compose.yml")}).Return(nil)

				return service, composeManager, dockerServiceManager
			},
			expectedError: false,
		},
		{
			name: "Add already existing service",
			mocker: func(t *testing.T, ctrl *gomock.Controller, stack *data.MonitoringStack) (ServiceAPI, ComposeManager, DockerServiceManager) {
				service := mocks.NewMockServiceAPI(ctrl)
				service.EXPECT().ContainerName().Return("sedge_prometheus").AnyTimes()
				service.EXPECT().Name().Return(PrometheusContainerName).AnyTimes()
				composeManager := mocks.NewMockComposeManager(ctrl)

				return service, composeManager, mocks.NewMockDockerServiceManager(ctrl)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ctrl := gomock.NewController(t)

			userDataHome := os.Getenv("XDG_DATA_HOME")
			if userDataHome == "" {
				userHome, err := os.UserHomeDir()
				require.NoError(t, err)
				userDataHome = filepath.Join(userHome, ".local", "share")
			}
			// Create a mock locker
			locker := mock_locker.NewMockLocker(ctrl)
			// Expect the lock to be acquired
			locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join(userDataHome, ".sedge", "monitoring", ".lock")}).Return(locker)
			locker.EXPECT().Lock().Return(nil).AnyTimes()
			locker.EXPECT().Locked().Return(true).AnyTimes()
			locker.EXPECT().Unlock().Return(nil).AnyTimes()

			// Expect the services to be triggered
			grafanaMock := mocks.NewMockServiceAPI(ctrl)
			promMock := mocks.NewMockServiceAPI(ctrl)
			nodeExporterMock := mocks.NewMockServiceAPI(ctrl)
			grafanaMock.EXPECT().ContainerName().Return(GrafanaContainerName).AnyTimes()
			grafanaMock.EXPECT().SetContainerIP(gomock.Any()).AnyTimes()
			grafanaMock.EXPECT().AddTarget(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			promMock.EXPECT().ContainerName().Return(PrometheusContainerName).AnyTimes()
			promMock.EXPECT().SetContainerIP(gomock.Any()).AnyTimes()
			promMock.EXPECT().AddTarget(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			nodeExporterMock.EXPECT().ContainerName().Return(NodeExporterContainerName).AnyTimes()
			nodeExporterMock.EXPECT().SetContainerIP(gomock.Any()).AnyTimes()
			nodeExporterMock.EXPECT().AddTarget(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			// Create a filesystem
			fs := afero.NewMemMapFs()
			_, err := fs.Create(filepath.Join(userDataHome, ".sedge", "monitoring", "docker-compose.yml"))
			require.NoError(t, err)
			// Create the dotenv file
			dotenvFile, err := fs.Create(filepath.Join(userDataHome, ".sedge", "monitoring", ".env"))
			require.NoError(t, err)
			dotenv := map[string]string{
				"NODE_PORT": "9000",
			}
			// Write the dotenv file
			for key, value := range dotenv {
				_, err := dotenvFile.WriteString(key + "=" + value + "\n")
				require.NoError(t, err)
			}

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{grafanaMock, promMock, nodeExporterMock},
				mocks.NewMockComposeManager(ctrl),
				mocks.NewMockDockerServiceManager(ctrl),
				fs,
				locker,
			)

			service, composeManager, dockerServiceManager := tt.mocker(t, ctrl, manager.stack)
			manager.composeManager = composeManager
			manager.dockerServiceManager = dockerServiceManager
			err = manager.AddService(service)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, manager.services, service)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, manager.services, service)
			}
		})
	}
}

func TestUpdateDockerComposeFile(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	baseTemplate := `{{/* docker-compose_base.tmpl */}}
{{ define "docker-compose" }}
services:
  base_service:
    image: base:latest
{{ template "additional_services" . }}
networks:
  sedge:
    name: sedge-network
    external: true
{{ end }}
`

	lidoExporterTemplate := `{{/* lido_exporter.tmpl */}}
{{ define "additional_services" }}
  {{ if .LidoExporter }}
  lido_exporter:
    image: lido_exporter:latest
    networks:
      - sedge
  {{ end }}
{{ end }}
`

	tests := []struct {
		name            string
		serviceName     string
		baseTemplate    string
		serviceTemplate string
		expectedCompose string
		wantErr         bool
	}{
		{
			name:            "add lido exporter service",
			serviceName:     "lido_exporter",
			baseTemplate:    baseTemplate,
			serviceTemplate: lidoExporterTemplate,
			expectedCompose: `services:
  base_service:
    image: base:latest
  lido_exporter:
    image: lido_exporter:latest
    networks:
      - sedge
networks:
  sedge:
    name: sedge-network
    external: true
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			ctrl := gomock.NewController(t)
			mockLocker := mock_locker.NewMockLocker(ctrl)
			mockService := mocks.NewMockServiceAPI(ctrl)

			// Mock the locker behavior
			mockLocker.EXPECT().New(gomock.Any()).Return(mockLocker).AnyTimes()
			mockLocker.EXPECT().Lock().Return(nil).AnyTimes()
			mockLocker.EXPECT().Locked().Return(true).AnyTimes()
			mockLocker.EXPECT().Unlock().Return(nil).AnyTimes()

			// Mock the service behavior
			mockService.EXPECT().Name().Return(tt.serviceName).AnyTimes()

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{},
				mocks.NewMockComposeManager(ctrl),
				mocks.NewMockDockerServiceManager(ctrl),
				fs,
				mockLocker,
			)

			// Create a mock monitoring filesystem
			mockMonitoringFs := afero.NewMemMapFs()
			err := afero.WriteFile(mockMonitoringFs, "services/docker-compose_base.tmpl", []byte(tt.baseTemplate), 0o644)
			require.NoError(t, err)
			err = afero.WriteFile(mockMonitoringFs, "services/"+tt.serviceName+".tmpl", []byte(tt.serviceTemplate), 0o644)
			require.NoError(t, err)

			// Call the function being tested
			err = manager.updateDockerComposeFile(mockService, afero.NewIOFS(mockMonitoringFs))

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				content, err := manager.stack.ReadFile("docker-compose.yml")
				assert.NoError(t, err)
				// Normalize whitespace and remove empty lines
				normalizeWhitespace := func(s string) string {
					lines := strings.Split(s, "\n")
					var nonEmptyLines []string
					for _, line := range lines {
						trimmed := strings.TrimSpace(line)
						if trimmed != "" {
							nonEmptyLines = append(nonEmptyLines, trimmed)
						}
					}
					return strings.Join(nonEmptyLines, " ")
				}

				expectedNormalized := normalizeWhitespace(tt.expectedCompose)
				actualNormalized := normalizeWhitespace(string(content))

				assert.Equal(t, expectedNormalized, actualNormalized)
			}
		})
	}
}

func TestUpdateEnvFile(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	tests := []struct {
		name        string
		initialEnv  string
		newEnv      map[string]string
		expectedEnv string
		wantErr     bool
	}{
		{
			name:        "add new variables",
			initialEnv:  "EXISTING_VAR=value\n",
			newEnv:      map[string]string{"NEW_VAR": "new_value"},
			expectedEnv: "EXISTING_VAR=value\nNEW_VAR=new_value\n",
			wantErr:     false,
		},
		{
			name:        "update existing variable",
			initialEnv:  "EXISTING_VAR=old_value\n",
			newEnv:      map[string]string{"EXISTING_VAR": "new_value"},
			expectedEnv: "EXISTING_VAR=new_value\n",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			ctrl := gomock.NewController(t)
			mockLocker := mock_locker.NewMockLocker(ctrl)

			// Mock the locker behavior
			mockLocker.EXPECT().New(gomock.Any()).Return(mockLocker).AnyTimes()
			mockLocker.EXPECT().Lock().Return(nil).AnyTimes()
			mockLocker.EXPECT().Locked().Return(true).AnyTimes()
			mockLocker.EXPECT().Unlock().Return(nil).AnyTimes()

			// Create a monitoring manager
			manager := NewMonitoringManager(
				[]ServiceAPI{},
				mocks.NewMockComposeManager(ctrl),
				mocks.NewMockDockerServiceManager(ctrl),
				fs,
				mockLocker,
			)

			err := afero.WriteFile(fs, filepath.Join(manager.stack.Path(), ".env"), []byte(tt.initialEnv), 0o644)
			require.NoError(t, err)

			err = manager.updateEnvFile(tt.newEnv)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				content, err := afero.ReadFile(fs, filepath.Join(manager.stack.Path(), ".env"))
				assert.NoError(t, err)

				// Normalize and sort the lines for comparison
				expectedLines := strings.Split(strings.TrimSpace(tt.expectedEnv), "\n")
				actualLines := strings.Split(strings.TrimSpace(string(content)), "\n")

				sort.Strings(expectedLines)
				sort.Strings(actualLines)

				assert.Equal(t, expectedLines, actualLines)
			}
		})
	}
}

func TestValidateNewService(t *testing.T) {
	tests := []struct {
		name          string
		service       func(t *testing.T, ctrl *gomock.Controller) ServiceAPI
		expectedError error
	}{
		{
			name: "new service",
			service: func(t *testing.T, ctrl *gomock.Controller) ServiceAPI {
				serviceMock := mocks.NewMockServiceAPI(ctrl)

				serviceMock.EXPECT().ContainerName().Return("SedgeService").AnyTimes()

				return serviceMock
			},
			expectedError: nil,
		},
		{
			name: "existing service",
			service: func(t *testing.T, ctrl *gomock.Controller) ServiceAPI {
				promMock := mocks.NewMockServiceAPI(ctrl)

				promMock.EXPECT().ContainerName().Return(PrometheusContainerName).AnyTimes()

				return promMock
			},
			expectedError: errors.New("service sedge_prometheus already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ctrl := gomock.NewController(t)
			grafanaMock := mocks.NewMockServiceAPI(ctrl)
			promMock := mocks.NewMockServiceAPI(ctrl)
			nodeExporterMock := mocks.NewMockServiceAPI(ctrl)

			// Expect the service to be triggered

			grafanaMock.EXPECT().ContainerName().Return(GrafanaContainerName).AnyTimes()
			promMock.EXPECT().ContainerName().Return(PrometheusContainerName).AnyTimes()
			nodeExporterMock.EXPECT().ContainerName().Return(NodeExporterContainerName).AnyTimes()

			// Init monitoring manager and services
			manager := MonitoringManager{
				services: []ServiceAPI{grafanaMock, promMock, nodeExporterMock},
			}
			err := manager.validateNewService(tt.service(t, ctrl))
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMakeTarget(t *testing.T) {
	manager := &MonitoringManager{}

	tests := []struct {
		name           string
		service        func(t *testing.T, ctrl *gomock.Controller) ServiceAPI
		expectedTarget types.MonitoringTarget
		expectedLabels map[string]string
		expectError    bool
	}{
		{
			name: "Valid HTTP endpoint",
			service: func(t *testing.T, ctrl *gomock.Controller) ServiceAPI {
				serviceMock := mocks.NewMockServiceAPI(ctrl)
				serviceMock.EXPECT().ContainerName().Return("SedgeService").AnyTimes()
				serviceMock.EXPECT().Endpoint().Return("http://sedge_service:8080").AnyTimes()
				return serviceMock
			},
			expectedTarget: types.MonitoringTarget{Host: "SedgeService", Port: 8080, Path: "/metrics"},
			expectedLabels: map[string]string{InstanceIDLabel: "SedgeService"},
			expectError:    false,
		},
		{
			name: "Valid HTTPS endpoint",
			service: func(t *testing.T, ctrl *gomock.Controller) ServiceAPI {
				serviceMock := mocks.NewMockServiceAPI(ctrl)
				serviceMock.EXPECT().ContainerName().Return("Service").AnyTimes()
				serviceMock.EXPECT().Endpoint().Return("https://service:8084").AnyTimes()
				return serviceMock
			},
			expectedTarget: types.MonitoringTarget{Host: "Service", Port: 8084, Path: "/metrics"},
			expectedLabels: map[string]string{InstanceIDLabel: "Service"},
			expectError:    false,
		},
		{
			name: "Invalid endpoint format",
			service: func(t *testing.T, ctrl *gomock.Controller) ServiceAPI {
				serviceMock := mocks.NewMockServiceAPI(ctrl)
				serviceMock.EXPECT().ContainerName().Return("Service").AnyTimes()
				serviceMock.EXPECT().Endpoint().Return("invalid endpoint").AnyTimes()
				return serviceMock
			},
			expectedTarget: types.MonitoringTarget{},
			expectedLabels: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			monitoringTarget, labels, err := manager.makeTarget(tt.service(t, ctrl))

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTarget, monitoringTarget)
				assert.Equal(t, tt.expectedLabels, labels)
			}
		})
	}
}
