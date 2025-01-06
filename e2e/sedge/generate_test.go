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
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NethermindEth/sedge/internal/utils"

	"github.com/NethermindEth/sedge/internal/lido/contracts"
	"github.com/NethermindEth/sedge/internal/lido/contracts/mevboostrelaylist"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	base "github.com/NethermindEth/sedge/e2e"
)

func TestE2E_Generate_FullNode_GoerliNotSupported(t *testing.T) {
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
		func(t *testing.T, binaryPath, dataDirPath string) {
			runErr = base.RunSedge(t, binaryPath, "generate", "full-node", "--network", "goerli")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.Error(t, runErr, "generate command should fail without arguments")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_Generate_FullNode_Lido_GnosisNotSupported(t *testing.T) {
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
			runErr = base.RunSedge(t, binaryPath, "generate", "--lido", "full-node", "--network", "gnosis")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.Error(t, runErr, "generate command should fail")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_Generate_FullNode_Lido_SepoliaNotSupported(t *testing.T) {
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
			runErr = base.RunSedge(t, binaryPath, "generate", "--lido", "full-node", "--network", "sepolia")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			// Sepolia not supported Lido with MEV-Boost
			assert.Error(t, runErr, "generate command should fail")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_Generate_FullNode_Lido_Sepolia_NoMEV(t *testing.T) {
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
			runErr = base.RunSedge(t, binaryPath, "generate", "--lido", "full-node", "--network", "sepolia", "--no-mev-boost")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr, "generate command should succeed with the given arguments")
			generateDataFilePath := filepath.Join(dataDirPath, ".env")
			assert.FileExists(t, generateDataFilePath, ".env file should be created")

			// Read .env file
			envMap, err := godotenv.Read(generateDataFilePath)
			assert.NoError(t, err, "should be able to read .env file")

			// Read FEE_RECIPIENT value
			feeRecipient, exists := envMap["FEE_RECIPIENT"]
			assert.True(t, exists, "FEE_RECIPIENT should exist in .env file")
			expectedFeeRecipient, ok := contracts.FeeRecipient("sepolia")
			assert.True(t, ok, "FeeRecipient should be found")
			assert.Equal(t, expectedFeeRecipient, feeRecipient, "FEE_RECIPIENT value should match expected value")

			// Read RELAY_URLS value
			_, exists = envMap["RELAY_URLS"]
			assert.False(t, exists, "RELAY_URLS shouldn't exist in .env file")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_Generate_FullNode_Lido_Mainnet(t *testing.T) {
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
			runErr = base.RunSedge(t, binaryPath, "generate", "--lido", "full-node", "--network", "mainnet")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr, "generate command should succeed")
			generateDataFilePath := filepath.Join(dataDirPath, ".env")
			assert.FileExists(t, generateDataFilePath, ".env file should be created")

			// Read .env file
			envMap, err := godotenv.Read(generateDataFilePath)
			assert.NoError(t, err, "should be able to read .env file")

			// Read FEE_RECIPIENT value
			feeRecipient, exists := envMap["FEE_RECIPIENT"]
			assert.True(t, exists, "FEE_RECIPIENT should exist in .env file")
			expectedFeeRecipient, ok := contracts.FeeRecipient("mainnet")
			assert.True(t, ok, "FeeRecipient should be found")
			assert.Equal(t, expectedFeeRecipient, feeRecipient, "FEE_RECIPIENT value should match expected value")

			// Read RELAY_URLS value
			relayURLs, exists := envMap["RELAY_URLS"]
			assert.True(t, exists, "RELAY_URLS should exist in .env file")
			relayURLsList := strings.Split(relayURLs, ",")
			expectedRelayURLs, _ := mevboostrelaylist.RelaysURI("mainnet")
			assert.Equal(t, expectedRelayURLs, relayURLsList, "RELAY_URLS value should match expected value")
		},
	)
	// Run test case
	e2eTest.run()
}

func TestE2E_Generate_FullNode_Nimbus_Mainnet(t *testing.T) {
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
			runErr = base.RunSedge(t, binaryPath, "generate", "full-node", "--network", "mainnet", "-v", "nimbus")
		},
		// Assert
		func(t *testing.T, dataDirPath string) {
			assert.NoError(t, runErr, "generate command should succeed")
			generateDataFilePath := filepath.Join(dataDirPath, ".env")
			assert.FileExists(t, generateDataFilePath, ".env file should be created")

			// Read .env file
			envMap, err := godotenv.Read(generateDataFilePath)
			assert.NoError(t, err, "should be able to read .env file")

			// Read CC_IMAGE_VERSION value
			ccImageVersion, exists := envMap["CC_IMAGE_VERSION"]
			assert.True(t, exists, "FEE_RECIPIENT should exist in .env file")
			assert.Contains(t, ccImageVersion, "nimbus", "CC_IMAGE_VERSION value should contain 'nimbus'")

			// Read VL_IMAGE_VERSION value
			vlImageVersion, exists := envMap["VL_IMAGE_VERSION"]
			assert.True(t, exists, "VL_IMAGE_VERSION should exist in .env file")
			assert.Contains(t, vlImageVersion, "nimbus", "VL_IMAGE_VERSION value should contain 'nimbus'")

			// Check that the docker-compose contains the consensus and validator services
			dockerComposeFilePath := filepath.Join(dataDirPath, "docker-compose.yml")
			assert.FileExists(t, dockerComposeFilePath, "docker-compose.yml file should be created")
			err = utils.ValidateCompose(dockerComposeFilePath)
			assert.NoError(t, err, "docker-compose file should be valid")
			composeServices, err := utils.LoadDockerComposeServices(dockerComposeFilePath)
			for _, service := range []string{"execution", "consensus", "validator"} {
				assert.Contains(t, composeServices, service, fmt.Sprintf("docker-compose file should contain service %s", service))
			}
		},
	)
	// Run test case
	e2eTest.run()
}
