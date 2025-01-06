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
	"io"
	"testing"

	"github.com/NethermindEth/sedge/internal/common"
	"github.com/NethermindEth/sedge/internal/monitoring"
	sedge_mocks "github.com/NethermindEth/sedge/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitoringCmd(t *testing.T) {
	tcs := []struct {
		name   string
		flags  []string
		mocker func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager
		isErr  bool
	}{
		{
			name:  "valid monitoring init: default",
			flags: []string{"init", "default"},
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				mockManager := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					mockManager.EXPECT().InstallationStatus().Return(common.NotInstalled, nil).AnyTimes(),
					mockManager.EXPECT().InstallStack().Return(nil).AnyTimes(),
					mockManager.EXPECT().Status().Return(common.Created, nil).AnyTimes(),
					mockManager.EXPECT().Run().Return(nil).AnyTimes(),
					mockManager.EXPECT().Init().Return(nil).AnyTimes(),
				)
				return mockManager
			},
			isErr: false,
		},
		{
			name:  "valid monitoring clean",
			flags: []string{"clean"},
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				mockManager := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					mockManager.EXPECT().InstallationStatus().Return(common.Installed, nil).AnyTimes(),
					mockManager.EXPECT().Cleanup().Return(nil).AnyTimes(),
				)
				return mockManager
			},
			isErr: false,
		},
		{
			name:  "valid monitoring init: lido",
			flags: []string{"init", "lido", "--node-operator-id", "1"},
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				mockManager := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					mockManager.EXPECT().InstallationStatus().Return(common.NotInstalled, nil).AnyTimes(),
					mockManager.EXPECT().InstallStack().Return(nil).AnyTimes(),
					mockManager.EXPECT().Status().Return(common.Created, nil).AnyTimes(),
					mockManager.EXPECT().Run().Return(nil).AnyTimes(),
					mockManager.EXPECT().Init().Return(nil).AnyTimes(),
					mockManager.EXPECT().AddService(gomock.Any()).Return(nil).AnyTimes(),
				)
				return mockManager
			},
			isErr: false,
		},
		{
			name:  "invalid monitoring init: lido, no nodeID or reward address",
			flags: []string{"init", "lido"},
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				return sedge_mocks.NewMockMonitoringManager(ctrl)
			},
			isErr: true,
		},
		{
			name:  "valid monitoring init: lido, Mainnet",
			flags: []string{"init", "lido", "--node-operator-id", "1", "--network", "mainnet"},
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				mockManager := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					mockManager.EXPECT().InstallationStatus().Return(common.NotInstalled, nil).AnyTimes(),
					mockManager.EXPECT().InstallStack().Return(nil).AnyTimes(),
					mockManager.EXPECT().Status().Return(common.Created, nil).AnyTimes(),
					mockManager.EXPECT().Run().Return(nil).AnyTimes(),
					mockManager.EXPECT().Init().Return(nil).AnyTimes(),
					mockManager.EXPECT().AddService(gomock.Any()).Return(nil).AnyTimes(),
				)
				return mockManager
			},
			isErr: false,
		},
		{
			name:  "invalid monitoring init: lido, no nodeID or reward address, Mainnet",
			flags: []string{"init", "lido", "--network", "mainnet"},
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				return sedge_mocks.NewMockMonitoringManager(ctrl)
			},
			isErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			log.SetOutput(io.Discard)
			logsOut := new(bytes.Buffer)
			tableOut := new(bytes.Buffer)
			rootCmd := RootCmd()
			rootCmd.SetOut(tableOut)
			rootCmd.AddCommand(MonitoringCmd(tc.mocker(t, ctrl)))
			argsL := append([]string{"monitoring"}, tc.flags...)
			rootCmd.SetArgs(argsL)
			initLogging()
			log.SetOutput(logsOut)
			rootCmd.SetOutput(io.Discard)
			err := rootCmd.Execute()
			if tc.isErr && err == nil {
				t.Error("sedge monitoring expected to fail")
			} else if !tc.isErr && err != nil {
				t.Errorf("sedge monitoring failed: %v", err)
			}
		})
	}
}

func TestInitMonitoring(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	tests := []struct {
		name    string
		mocker  func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager
		wantErr bool
	}{
		{
			name: "monitoring -> prev: not installed, after: installation status error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.NotInstalled, errors.New("installation status error")),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
		{
			name: "monitoring -> prev: not installed, after: installed and started",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.NotInstalled, nil),
					monitoringMgr.EXPECT().InstallStack().Return(nil),
					monitoringMgr.EXPECT().Status().Return(common.Running, nil),
					monitoringMgr.EXPECT().Init().Return(nil),
				)
				return monitoringMgr
			},
		},
		{
			name: "monitoring -> prev: not installed, after: installation failed",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.NotInstalled, nil),
					monitoringMgr.EXPECT().InstallStack().Return(monitoring.ErrInstallingMonitoringMngr),
					monitoringMgr.EXPECT().Cleanup().Return(nil),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
		{
			name: "monitoring -> prev: not installed, after: installation failed, cleanup error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.NotInstalled, nil),
					monitoringMgr.EXPECT().InstallStack().Return(monitoring.ErrInstallingMonitoringMngr),
					monitoringMgr.EXPECT().Cleanup().Return(errors.New("cleanup error")),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
		{
			name: "monitoring -> prev: not installed, after: installation failed but no cleanup needed",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.NotInstalled, nil),
					monitoringMgr.EXPECT().InstallStack().Return(errors.New("init error")),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
		{
			name: "monitoring -> prev: installed and running, after: installed and running",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Running, nil),
					monitoringMgr.EXPECT().Init().Return(nil),
				)
				return monitoringMgr
			},
		},
		{
			name: "monitoring -> prev: installed and created, after: installed and running",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Created, nil),
					monitoringMgr.EXPECT().Run().Return(nil),
					monitoringMgr.EXPECT().Init().Return(nil),
				)
				return monitoringMgr
			},
		},
		{
			name: "monitoring -> prev: installed and created, after: installed and run-error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Created, nil),
					monitoringMgr.EXPECT().Run().Return(errors.New("run error")),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
		{
			name: "monitoring -> prev: installed, after: installed and status error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Unknown, errors.New("status error")),
					monitoringMgr.EXPECT().Run().Return(nil),
					monitoringMgr.EXPECT().Init().Return(nil),
				)
				return monitoringMgr
			},
			wantErr: false,
		},
		{
			name: "monitoring -> prev: installed and restarting, after: installed and restarting",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Restarting, nil),
					monitoringMgr.EXPECT().Init().Return(nil),
				)
				return monitoringMgr
			},
		},
		{
			name: "monitoring -> prev: installed and broken, after: installed and re-run",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Broken, nil),
					monitoringMgr.EXPECT().Run().Return(nil),
					monitoringMgr.EXPECT().Init().Return(nil),
				)
				return monitoringMgr
			},
		},
		{
			name: "monitoring -> prev: installed and broken, after: installed and run error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Broken, nil),
					monitoringMgr.EXPECT().Run().Return(errors.New("run error")),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
		{
			name: "monitoring -> prev: installed and created, after: monitoring stack initialization error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Status().Return(common.Created, nil),
					monitoringMgr.EXPECT().Run().Return(nil),
					monitoringMgr.EXPECT().Init().Return(monitoring.ErrInitializingMonitoringMngr),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock monitoring manager.
			ctrl := gomock.NewController(t)

			// Get monitoring manager mock
			monitoringMgr := tt.mocker(t, ctrl)

			err := InitMonitoring(true, true, monitoringMgr, nil)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCleanMonitoring(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	tests := []struct {
		name    string
		mocker  func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager
		wantErr bool
	}{
		{
			name: "monitoring -> prev: not installed, after: nothing to do",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				monitoringMgr.EXPECT().InstallationStatus().Return(common.NotInstalled, nil)
				return monitoringMgr
			},
		},
		{
			name: "monitoring -> prev: installed, after: uninstalled",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Cleanup().Return(nil),
				)
				return monitoringMgr
			},
		},
		{
			name: "monitoring -> prev: installed, after: uninstalled failed",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				gomock.InOrder(
					monitoringMgr.EXPECT().InstallationStatus().Return(common.Installed, nil),
					monitoringMgr.EXPECT().Cleanup().Return(assert.AnError),
				)
				return monitoringMgr
			},
			wantErr: true,
		},
		{
			name: "monitoring -> installation status error",
			mocker: func(t *testing.T, ctrl *gomock.Controller) *sedge_mocks.MockMonitoringManager {
				monitoringMgr := sedge_mocks.NewMockMonitoringManager(ctrl)
				monitoringMgr.EXPECT().InstallationStatus().Return(common.Unknown, assert.AnError)
				return monitoringMgr
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock monitoring manager.
			ctrl := gomock.NewController(t)

			// Get monitoring manager mock
			monitoringMgr := tt.mocker(t, ctrl)

			err := CleanMonitoring(monitoringMgr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
