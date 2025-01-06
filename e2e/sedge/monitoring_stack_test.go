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
package e2e

import (
	"runtime"
	"testing"

	base "github.com/NethermindEth/sedge/e2e"
	"github.com/stretchr/testify/assert"
)

func skipIfNotAMD64(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("Skipping test on non-AMD64 architecture")
	}
}

var grafanaOnCallContainers = []string{"engine", "celery", "redis", "oncall_setup"}

// TestMonitoringStack_Init tests that the monitoring stack is not initialized if the user does not run the init-monitoring command
func TestE2E_MonitoringStack_NotInitialized(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		nil,
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "--help")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)

			checkMonitoringStackNotInstalled(t)
			checkMonitoringStackContainersNotRunning(t, grafanaOnCallContainers...)
		},
	)
	// Run test case
	e2eTest.run()
}

// TestMonitoringStack_Init tests the monitoring stack initialization
func TestE2E_MonitoringStack_Init(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		nil,
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "default")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)
			checkMonitoringStackDir(t)
			checkPrometheusDir(t)
			checkMonitoringStackContainers(t, grafanaOnCallContainers...)
			checkPrometheusTargetsUp(t, "sedge_node_exporter:9100")
			checkGrafanaHealth(t)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_NotReinstalled(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		grafanaContainerID      string
		prometheusContainerID   string
		nodeExporterContainerID string
		runErr                  error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		func(t *testing.T, sedgePath string) error {
			err := base.RunCommand(t, sedgePath, "sedge", "monitoring", "init", "default")
			if err != nil {
				return err
			}
			grafanaContainerID, err = getContainerIDByName("sedge_grafana")
			if err != nil {
				return err
			}
			prometheusContainerID, err = getContainerIDByName("sedge_prometheus")
			if err != nil {
				return err
			}
			nodeExporterContainerID, err = getContainerIDByName("sedge_node_exporter")
			return err
		},
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)

			checkMonitoringStackDir(t)
			checkPrometheusDir(t)
			checkMonitoringStackContainers(t)
			checkGrafanaHealth(t)
			newGrafanaContainerID, err := getContainerIDByName("sedge_grafana")
			assert.NoError(t, err)
			assert.Equal(t, grafanaContainerID, newGrafanaContainerID, "grafana container ID has changed")
			newPrometheusContainerID, err := getContainerIDByName("sedge_prometheus")
			assert.NoError(t, err)
			assert.Equal(t, prometheusContainerID, newPrometheusContainerID, "prometheus container ID has changed")
			newNodeExporterContainerID, err := getContainerIDByName("sedge_node_exporter")
			assert.NoError(t, err)
			assert.Equal(t, nodeExporterContainerID, newNodeExporterContainerID, "node-exporter container ID has changed")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_Clean(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		func(t *testing.T, sedgePath string) error {
			return base.RunCommand(t, sedgePath, "sedge", "monitoring", "init", "default")
		},
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "clean")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)

			// Check that monitoring stack directory is removed
			assert.NoDirExists(t, dataDirPath)

			// Check that monitoring stack containers are removed
			checkMonitoringStackContainersNotRunning(t, grafanaOnCallContainers...)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_CleanNonExistent(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		nil,
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "clean")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)

			// Check that monitoring stack directory doesn't exist
			assert.NoDirExists(t, dataDirPath)

			// Check that monitoring stack containers don't exist
			checkMonitoringStackContainersNotRunning(t, grafanaOnCallContainers...)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido_ValidID(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		nil,
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido", "--node-operator-id", "1")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)
			checkMonitoringStackDir(t)
			checkPrometheusDir(t)
			checkMonitoringStackContainers(t, "sedge_lido_exporter")
			checkPrometheusTargetsUp(t, "sedge_lido_exporter:8080", "sedge_node_exporter:9100")
			checkGrafanaHealth(t)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_CleanLido(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		func(t *testing.T, sedgePath string) error {
			return base.RunCommand(t, sedgePath, "sedge", "monitoring", "init", "lido", "--node-operator-id", "10")
		},
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "clean")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)

			// Check that monitoring stack directory is removed
			assert.NoDirExists(t, dataDirPath)

			// Check that monitoring stack containers are removed
			checkMonitoringStackContainersNotRunning(t, "sedge_lido_exporter")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido_InvalidAddress(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		func(t *testing.T, sedgePath string) error {
			return base.RunCommand(t, sedgePath, "sedge", "monitoring", "clean")
		},
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido", "--reward-address", "lol_what_a_reward_address")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.Error(t, runErr)

			checkMonitoringStackNotInstalled(t)
			checkMonitoringStackContainersNotRunning(t, grafanaOnCallContainers...)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido_OccupiedPort(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		nil,
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido", "--node-operator-id", "10", "--port", "9090")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.Error(t, runErr)
			checkContainerNotRunning(t, "sedge_lido_exporter")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		nil,
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido",
				"--rpc-endpoints", "https://endpoints.omniatech.io/v1/eth/holesky/public,https://ethereum-holesky-rpc.publicnode.com",
				"--ws-endpoints", "https://ethereum-holesky-rpc.publicnode.com,wss://ethereum-holesky-rpc.publicnode.com",
				"--port", "9989",
				"--scrape-time", "30s",
				"--network", "holesky",
				"--node-operator-id", "250",
				"--reward-address", "0x22bA5CaFB5E26E6Fe51f330294209034013A5A4c",
			)
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)
			checkMonitoringStackDir(t)
			checkPrometheusDir(t)
			checkMonitoringStackContainers(t, "sedge_lido_exporter")
			checkPrometheusTargetsUp(t, "sedge_lido_exporter:9989", "sedge_node_exporter:9100")
			checkGrafanaHealth(t)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido_InvalidNodeID(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		func(t *testing.T, sedgePath string) error {
			return base.RunCommand(t, sedgePath, "sedge", "monitoring", "clean")
		},
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido", "--node-operator-id", "-1")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.Error(t, runErr)

			checkMonitoringStackNotInstalled(t)
			checkMonitoringStackContainersNotRunning(t, grafanaOnCallContainers...)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido_InvalidNodeID_Mainnet(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		func(t *testing.T, sedgePath string) error {
			return base.RunCommand(t, sedgePath, "sedge", "monitoring", "clean")
		},
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido", "--node-operator-id", "-20", "--network", "mainnet")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.Error(t, runErr)

			checkMonitoringStackNotInstalled(t)
			checkMonitoringStackContainersNotRunning(t, grafanaOnCallContainers...)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido_InvalidAddress_Mainnet(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		func(t *testing.T, sedgePath string) error {
			return base.RunCommand(t, sedgePath, "sedge", "monitoring", "clean")
		},
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido", "--reward-address", "xrewardx", "--network", "mainnet")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.Error(t, runErr)

			checkMonitoringStackNotInstalled(t)
			checkMonitoringStackContainersNotRunning(t, grafanaOnCallContainers...)
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_MonitoringStack_InitLido_ValidID_Mainnet(t *testing.T) {
	skipIfNotAMD64(t)
	// Test context
	var (
		runErr error
	)
	// Build test case
	e2eTest := newE2ESedgeTestCase(
		t,
		// Arrange
		nil,
		// Act
		func(t *testing.T, binaryPath string, dataDirPath string) {
			runErr = base.RunCommand(t, binaryPath, "sedge", "monitoring", "init", "lido", "--node-operator-id", "1", "--network", "mainnet")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr)
			checkMonitoringStackDir(t)
			checkPrometheusDir(t)
			checkMonitoringStackContainers(t, "sedge_lido_exporter")
			checkPrometheusTargetsUp(t, "sedge_lido_exporter:8080", "sedge_node_exporter:9100")
			checkGrafanaHealth(t)
		},
	)
	// Run test case
	e2eTest.run()
}
