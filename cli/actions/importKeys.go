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
package actions

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/images/validator-import/lighthouse"
	"github.com/NethermindEth/sedge/internal/images/validator-import/prysm"
	"github.com/NethermindEth/sedge/internal/images/validator-import/teku"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/internal/pkg/services"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

var ErrInterrupted = errors.New("interrupt")

type ImportValidatorKeysOptions struct {
	ValidatorClient string
	Network         string
	StopValidator   bool
	StartValidator  bool
	From            string
	GenerationPath  string
	ContainerTag    string
	CustomConfig    ImportValidatorKeysCustomOptions
	Distributed     bool
}
type ImportValidatorKeysCustomOptions struct {
	NetworkConfigPath string
	GenesisPath       string
	DeployBlockPath   string
}

func (s *sedgeActions) ImportValidatorKeys(options ImportValidatorKeysOptions) error {
	// Ensure paths are absolute
	if absFrom, err := filepath.Abs(options.From); err != nil {
		return err
	} else {
		options.From = absFrom
	}
	if absGenerationPath, err := filepath.Abs(options.GenerationPath); err != nil {
		return err
	} else {
		options.GenerationPath = absGenerationPath
	}
	validatorCtName := services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
	// Check validator container exists
	_, err := s.dockerServiceManager.ContainerID(validatorCtName)
	if err != nil {
		return err
	}
	previouslyRunning, err := s.dockerServiceManager.IsRunning(validatorCtName)
	if err != nil {
		return err
	}
	// Stop validator
	if previouslyRunning {
		log.Info("Stopping validator client")
		if err := s.dockerServiceManager.Stop(validatorCtName); err != nil {
			return err
		}
	}

	absFrom, err := filepath.Abs(options.From)
	if err != nil {
		return err
	}
	options.From = absFrom
	absGenerationPath, err := filepath.Abs(options.GenerationPath)
	if err != nil {
		return err
	}
	options.GenerationPath = absGenerationPath

	if !isDefaultKeysPath(options.GenerationPath, options.From) {
		if !options.Distributed {
			defaultKeystorePath := filepath.Join(options.GenerationPath, "keystore")
			log.Warnf("The keys path is not the default one, copying the keys to the default path %s", defaultKeystorePath)
			copy.Copy(options.From, defaultKeystorePath)
		}
	}

	if options.Distributed {
		cwd, _ := os.Getwd()
		charonPath := filepath.Join(cwd, ".charon")

		if !isDefaultKeysPath(options.GenerationPath, options.From) {
			charonPath = options.From
			log.Infof("Copying the keys from %s", charonPath)
			options.From = filepath.Join(options.GenerationPath, "keystore")
		}
		defaultCharonPath := filepath.Join(configs.DefaultAbsSedgeDataPath, ".charon")
		// Copy the folder from charonPath to defaultCharonPath
		log.Infof("Copying Charon contents to the default path %s", defaultCharonPath)
		if err := os.MkdirAll(defaultCharonPath, 0o755); err != nil {
			return err
		}
		if err := copy.Copy(charonPath, defaultCharonPath); err != nil {
			return err
		}
		charonValidatorKeysPath := filepath.Join(charonPath, "validator_keys")
		defaultKeystorePath := filepath.Join(configs.DefaultAbsSedgeDataPath, "keystore")
		log.Infof("Copying the keys to the default path %s", defaultKeystorePath)
		if err := os.MkdirAll(defaultKeystorePath, 0o755); err != nil {
			return err
		}

		validatorKeysPath := filepath.Join(defaultKeystorePath, "validator_keys")
		if err := os.MkdirAll(validatorKeysPath, 0o755); err != nil {
			return err
		}

		depositDataPath := filepath.Join(charonPath, "deposit-data.json")
		depositDataPathDest := filepath.Join(defaultKeystorePath, "deposit-data.json")
		if err := copy.Copy(depositDataPath, depositDataPathDest); err != nil {
			return err
		}

		files, err := os.ReadDir(charonValidatorKeysPath)
		if err != nil {
			log.Fatal(err)
		}
		len := len(files)
		for i := 0; i < len/2; i++ {
			keystorePath := filepath.Join(charonValidatorKeysPath, fmt.Sprintf("keystore-%d.json", i))
			validatorPath := filepath.Join(validatorKeysPath, fmt.Sprintf("keystore-%d.json", i))
			if err := copy.Copy(keystorePath, validatorPath); err != nil {
				return err
			}

			keystoreTxtPath := filepath.Join(charonValidatorKeysPath, fmt.Sprintf("keystore-%d.txt", i))
			keystorePasswordPath := filepath.Join(defaultKeystorePath, fmt.Sprintf("keystore-%d.txt", i))
			if err := copy.Copy(keystoreTxtPath, keystorePasswordPath); err != nil {
				return err
			}
		}
		if options.ValidatorClient == "prysm" {
			keystorePasswordPath := filepath.Join(defaultKeystorePath, "keystore_password.txt")
			f, err := os.Create(keystorePasswordPath)
			if err != nil {
				return err
			}
			f.WriteString("prysm-validator-secret")
			defer f.Close()
		}
	}

	var ctID string
	switch options.ValidatorClient {
	case "prysm":
		prysmCtID := ""
		if options.Distributed {
			prysmCtID, err = setupPrysmValidatorImportContainerDV(s.dockerClient, s.commandRunner, s.dockerServiceManager, options)
			if err != nil {
				return err
			}
			ctID = prysmCtID
		} else {
			prysmCtID, err := setupPrysmValidatorImportContainer(s.dockerClient, s.dockerServiceManager, options)
			if err != nil {
				return err
			}
			ctID = prysmCtID
		}
	case "nimbus":
		nimbusCtID, err := setupNimbusValidatorImport(s.dockerClient, s.dockerServiceManager, options)
		if err != nil {
			return err
		}
		ctID = nimbusCtID
	case "lodestar":
		lodestarCtID, err := setupLodestarValidatorImport(s.dockerClient, s.dockerServiceManager, options)
		if err != nil {
			return err
		}
		ctID = lodestarCtID
	case "lighthouse":
		lighthouseCtID, err := setupLighthouseValidatorImport(s.dockerClient, s.commandRunner, options)
		if err != nil {
			return err
		}
		ctID = lighthouseCtID
	case "teku":
		tekuCtID, err := setupTekuValidatorImport(s.dockerClient, s.commandRunner, options)
		if err != nil {
			return err
		}
		ctID = tekuCtID
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedValidatorClient, options.ValidatorClient)
	}
	log.Info("Importing validator keys")
	var runErr error
	if options.ValidatorClient == "nimbus" {
		if !options.Distributed {
			runErr = runAndWaitImportKeysNimbus(s.dockerClient, s.dockerServiceManager, ctID)
		} else {
			runErr = runAndWaitImportKeys(s.dockerClient, s.dockerServiceManager, ctID)
		}
	} else {
		runErr = runAndWaitImportKeys(s.dockerClient, s.dockerServiceManager, ctID)
	}
	// Run validator again
	if (previouslyRunning && !options.StopValidator) || options.StartValidator {
		log.Info("The validator container is being restarted")
		if err := s.dockerServiceManager.Start(validatorCtName); err != nil {
			return err
		}
	}
	if runErr == nil {
		log.Info("Validator keys imported successfully")
	}
	return runErr
}

func isDefaultKeysPath(generationPath, from string) bool {
	return from == filepath.Join(generationPath, "keystore")
}

func setupPrysmValidatorImportContainer(dockerClient client.APIClient, dockerServiceManager DockerServiceManager, options ImportValidatorKeysOptions) (string, error) {
	var (
		validatorCtName       = services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
		validatorImportCtName = services.ContainerNameWithTag(services.ServiceCtValidatorImport, options.ContainerTag)
	)
	validatorImage, err := dockerServiceManager.Image(validatorCtName)
	if err != nil {
		return "", err
	}
	// Mounts
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: options.From,
			Target: "/keystore",
		},
	}
	// CMD
	cmd := []string{
		"--",
		"accounts", "import",
		"--accept-terms-of-use",
		"--keys-dir=/keystore/validator_keys",
		"--wallet-dir=/data/wallet",
		"--wallet-password-file=/keystore/keystore_password.txt",
		"--account-password-file=/keystore/keystore_password.txt",
	}
	// Custom options
	if options.CustomConfig.NetworkConfigPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: options.CustomConfig.NetworkConfigPath,
			Target: "/network_config/config.yml",
		})
		cmd = append(cmd, "--chain-config-file=/network_config/config.yml")
	} else {
		cmd = append(cmd, "--"+options.Network)
	}
	log.Debugf("Creating %s container", validatorImportCtName)
	ct, err := dockerClient.ContainerCreate(context.Background(),
		&container.Config{
			Image: validatorImage,
			Cmd:   cmd,
		},
		&container.HostConfig{
			Mounts:      mounts,
			VolumesFrom: []string{validatorCtName},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		validatorImportCtName,
	)
	if err != nil {
		return "", err
	}
	return ct.ID, nil
}

func setupNimbusValidatorImport(dockerClient client.APIClient, dockerServiceManager DockerServiceManager, options ImportValidatorKeysOptions) (string, error) {
	var (
		// In the case of Nimbus, it's the consensus client the one that import the keys.
		consensusCtName       = services.ContainerNameWithTag(services.DefaultSedgeConsensusClient, options.ContainerTag)
		validatorCtName       = services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
		validatorImportCtName = services.ContainerNameWithTag(services.ServiceCtValidatorImport, options.ContainerTag)
	)
	validatorImage, err := dockerServiceManager.Image(consensusCtName)
	if err != nil {
		return "", err
	}
	// Mounts
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: options.From,
			Target: "/keystore",
		},
	}
	// CMD
	cmd := []string{
		"deposits",
		"import",
		"--data-dir=/data",
		"--method=single-salt",
		"/keystore",
	}
	// Custom options
	if options.CustomConfig.NetworkConfigPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: options.CustomConfig.NetworkConfigPath,
			Target: "/network_config/config.yml",
		})
		cmd = append(cmd, "--config-file=/network_config/config.yml")
	} else {
		cmd = append(cmd, "--network="+options.Network)
	}
	containerConfig := &container.Config{
		Image:        validatorImage,
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		OpenStdin:    true,
		Tty:          true,
	}
	if options.Distributed {
		containerConfig = &container.Config{
			Image: validatorImage,
			Entrypoint: []string{
				"sh", "-c", `
				#!/usr/bin/env bash
				set -e
				tmpkeys="/keystore/validator_keys/tmpkeys"
				mkdir -p ${tmpkeys}
				for f in /keystore/validator_keys/keystore-*.json; do
					echo "Importing key ${f}"
					pwdfile="/keystore/$(basename "$f" .json).txt"
					password=$(cat ${pwdfile})
					echo "Using password file ${pwdfile}"
					echo "Using password ${password}"
					cp "${f}" "${tmpkeys}"
					# Import keystore with password.
						echo "$password" | \
						/home/user/nimbus_beacon_node deposits import \
						--data-dir=/data \
						${tmpkeys}
					filename="$(basename ${f})"
  					rm "${tmpkeys}/${filename}"
				done
			`,
			},
		}
	}

	log.Debugf("Creating %s container", validatorImportCtName)
	ct, err := dockerClient.ContainerCreate(context.Background(),
		containerConfig,
		&container.HostConfig{
			Mounts:      mounts,
			VolumesFrom: []string{consensusCtName, validatorCtName},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		validatorImportCtName,
	)
	if err != nil {
		return "", err
	}
	return ct.ID, nil
}

func setupLodestarValidatorImport(dockerClient client.APIClient, dockerServiceManager DockerServiceManager, options ImportValidatorKeysOptions) (string, error) {
	var (
		validatorCtName       = services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
		validatorImportCtName = services.ContainerNameWithTag(services.ServiceCtValidatorImport, options.ContainerTag)
	)
	validatorImage, err := dockerServiceManager.Image(validatorCtName)
	if err != nil {
		return "", err
	}
	// Mounts
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: options.From,
			Target: "/keystore",
		},
	}
	// CMD
	var preset string
	switch options.Network {
	case "mainnet", "sepolia", "holesky":
		preset = "mainnet"
	case "gnosis", "chiado":
		preset = "gnosis"
	default:
		return "", newUnknownLodestarPresetError(options.Network)
	}
	cmd := []string{
		"validator", "import",
		"--dataDir", "/data",
		"--importKeystores=/keystore/validator_keys",
		"--importKeystoresPassword=/keystore/keystore_password.txt",
	}
	if options.CustomConfig.NetworkConfigPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: options.CustomConfig.NetworkConfigPath,
			Target: "/network_config/config.yaml",
		})
		cmd = append(cmd, "--paramsFile=/network_config/config.yaml")
	} else {
		cmd = append(cmd, "--network", options.Network)
		cmd = append(cmd, "--preset", preset)
	}
	log.Debugf("Creating %s container", validatorImportCtName)
	containerConfig := &container.Config{
		Image: validatorImage,
		Cmd:   cmd,
	}
	if options.Distributed {
		containerConfig = &container.Config{
			Image: validatorImage,
			Entrypoint: []string{
				"sh", "-c", `
				#!/bin/sh
				set -e
				for f in /keystore/validator_keys/keystore-*.json; do
					echo "Importing key ${f}"
					pwdfile="/keystore/$(basename "$f" .json).txt"
					echo "Using password file ${pwdfile}"
					# Import keystore with password.
					node /usr/app/packages/cli/bin/lodestar validator import \
						--dataDir="/data" \
						--importKeystores="$f" \
						--importKeystoresPassword="${pwdfile}"
				done
			`,
			},
		}
	}
	ct, err := dockerClient.ContainerCreate(context.Background(),
		containerConfig,
		&container.HostConfig{
			Mounts:      mounts,
			VolumesFrom: []string{validatorCtName},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		validatorImportCtName,
	)
	if err != nil {
		return "", err
	}
	return ct.ID, nil
}

func setupLighthouseValidatorImport(dockerClient client.APIClient, commandRunner commands.CommandRunner, options ImportValidatorKeysOptions) (string, error) {
	var (
		validatorCtName       = services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
		validatorImportCtName = services.ContainerNameWithTag(services.ServiceCtValidatorImport, options.ContainerTag)
	)
	if options.Network == "chiado" {
		options.Network = "custom"
	}
	// Init build context
	contextDir, err := lighthouse.InitContext()
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrCreatingContextDir, err.Error())
	}
	// Build image
	buildCmd := commandRunner.BuildDockerBuildCMD(commands.DockerBuildOptions{
		Path: contextDir,
		Tag:  "sedge/validator-import-lighthouse",
		Args: map[string]string{
			"NETWORK":    options.Network,
			"LH_VERSION": configs.ClientImages.Validator.Lighthouse.String(),
		},
	})
	log.Infof(configs.RunningCommand, buildCmd.Cmd)
	if _, _, err = commandRunner.RunCMD(buildCmd); err != nil {
		return "", err
	}
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: options.From,
			Target: "/keystore",
		},
	}
	if options.CustomConfig.NetworkConfigPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: options.CustomConfig.NetworkConfigPath,
			Target: "/network_config/config.yaml",
		})
	}
	if options.CustomConfig.GenesisPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: options.CustomConfig.GenesisPath,
			Target: "/network_config/genesis.ssz",
		})
	}
	if options.CustomConfig.DeployBlockPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: options.CustomConfig.DeployBlockPath,
			Target: "/network_config/deploy_block.txt",
		})
	}
	log.Debugf("Creating %s container", validatorImportCtName)
	ct, err := dockerClient.ContainerCreate(context.Background(),
		&container.Config{
			Image: "sedge/validator-import-lighthouse",
		},
		&container.HostConfig{
			Mounts:      mounts,
			VolumesFrom: []string{validatorCtName},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		validatorImportCtName,
	)
	if err != nil {
		return "", err
	}
	return ct.ID, nil
}

func setupTekuValidatorImport(dockerClient client.APIClient, commandRunner commands.CommandRunner, options ImportValidatorKeysOptions) (string, error) {
	var (
		validatorCtName       = services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
		validatorImportCtName = services.ContainerNameWithTag(services.ServiceCtValidatorImport, options.ContainerTag)
	)
	// Init build context
	contextDir, err := teku.InitContext()
	if err != nil {
		return "", err
	}
	// Build image
	buildCmd := commandRunner.BuildDockerBuildCMD(commands.DockerBuildOptions{
		Path: contextDir,
		Tag:  "sedge/validator-import-teku",
	})
	log.Infof(configs.RunningCommand, buildCmd.Cmd)
	if _, _, err := commandRunner.RunCMD(buildCmd); err != nil {
		return "", err
	}
	// Mounts
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: options.From,
			Target: "/keystore",
		},
	}
	if options.CustomConfig.NetworkConfigPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: options.CustomConfig.NetworkConfigPath,
			Target: "/network_config/config.yml",
		})
	}
	log.Debugf("Creating %s container", validatorImportCtName)
	ct, err := dockerClient.ContainerCreate(context.Background(),
		&container.Config{
			Image: "sedge/validator-import-teku",
		},
		&container.HostConfig{
			Mounts:      mounts,
			VolumesFrom: []string{validatorCtName},
			LogConfig: container.LogConfig{
				Type: "json-file",
			},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		validatorImportCtName,
	)
	if err != nil {
		return "", err
	}
	return ct.ID, nil
}

func runAndWaitImportKeys(dockerClient client.APIClient, dockerServiceManager DockerServiceManager, ctID string) error {
	log.Debugf("import keys container id: %s", ctID)
	ctExit, errChan := dockerServiceManager.Wait(ctID, container.WaitConditionNextExit)
	log.Info("The keys import container is starting")
	if err := dockerClient.ContainerStart(context.Background(), ctID, container.StartOptions{}); err != nil {
		return err
	}
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)
	for {
		select {
		case exitResult := <-ctExit:
			logs, err := dockerServiceManager.ContainerLogs(ctID, "Import keys")
			if err != nil {
				return err
			}
			if err = deleteContainer(dockerClient, ctID); err != nil {
				return err
			}
			if exitResult.StatusCode != 0 {
				return newValidatorImportCtBadExitCodeError(ctID, exitResult.StatusCode, logs)
			}
			return nil
		case <-osSignals:
			if err := stopContainer(dockerClient, ctID); err != nil {
				return err
			}
			if err := deleteContainer(dockerClient, ctID); err != nil {
				return err
			}
			return ErrInterrupted
		case exitErr := <-errChan:
			return exitErr
		}
	}
}

func setupPrysmValidatorImportContainerDV(dockerClient client.APIClient, commandRunner commands.CommandRunner, serviceManager DockerServiceManager, options ImportValidatorKeysOptions) (string, error) {
	var (
		validatorCtName       = services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
		validatorImportCtName = services.ContainerNameWithTag(services.ServiceCtValidatorImport, options.ContainerTag)
	)
	// Init build context
	contextDir, err := prysm.InitContext()
	if err != nil {
		return "", err
	}
	// Build image
	buildCmd := commandRunner.BuildDockerBuildCMD(commands.DockerBuildOptions{
		Path: contextDir,
		Tag:  "sedge/prysm-import-teku",
		Args: map[string]string{
			"NETWORK":       options.Network,
			"PRYSM_VERSION": configs.ClientImages.Validator.Prysm.String(),
		},
	})
	log.Infof(configs.RunningCommand, buildCmd.Cmd)
	if _, _, err := commandRunner.RunCMD(buildCmd); err != nil {
		return "", err
	}
	// Mounts
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: options.From,
			Target: "/keystore",
		},
	}
	log.Debugf("Creating %s container", validatorImportCtName)
	containerConfig := &container.Config{
		Image: "sedge/prysm-import-teku",
	}
	ct, err := dockerClient.ContainerCreate(context.Background(),
		containerConfig,
		&container.HostConfig{
			Mounts:      mounts,
			VolumesFrom: []string{validatorCtName},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		validatorImportCtName,
	)
	if err != nil {
		return "", err
	}
	return ct.ID, nil
}

// runAndWaitImportKeysNimbus starts the container in interactive mode and waits for it to finish.
func runAndWaitImportKeysNimbus(dockerClient client.APIClient, dockerServiceManager DockerServiceManager, ctID string) error {
	log.Debugf("Starting interactive container with id: %s", ctID)

	// Attach to the container's input/output for direct interaction
	resp, err := dockerClient.ContainerAttach(context.Background(), ctID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false, // Don't attach previous logs
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	// Put the terminal in raw mode for proper TTY handling
	oldState, err := terminal.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	defer func() {
		// Restore the terminal state immediately when the container finishes
		terminal.Restore(int(syscall.Stdin), oldState)
		// Clear the line again before printing the final success logs
		fmt.Print("\033[2K\r") // Clear the current line in the terminal
	}()

	// Start the container
	if err := dockerClient.ContainerStart(context.Background(), ctID, container.StartOptions{}); err != nil {
		return err
	}

	// Use goroutines to pipe stdin, stdout, and stderr directly to the user's terminal
	go func() {
		// Pipe container stdout and stderr to the terminal
		_, err := io.Copy(os.Stdout, resp.Reader)
		if err != nil {
			log.Errorf("Error piping container output: %v", err)
		}
	}()
	go func() {
		// Pipe terminal input to the container stdin
		_, err := io.Copy(resp.Conn, os.Stdin)
		if err != nil {
			log.Errorf("Error piping user input: %v", err)
		}
	}()

	// Wait for the container to finish execution
	ctExit, errChan := dockerServiceManager.Wait(ctID, container.WaitConditionNextExit)

	// Handle OS interrupts (e.g., Ctrl+C) to gracefully stop the container
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case exitResult := <-ctExit:
		if err = deleteContainer(dockerClient, ctID); err != nil {
			return err
		}
		// Wait for the container to exit normally
		if exitResult.StatusCode != 0 {
			log.Errorf("Container exited with non-zero status: %d", exitResult.StatusCode)
			return newValidatorImportCtBadExitCodeError(ctID, exitResult.StatusCode, "Container logs...")
		}

	case <-osSignals:
		// If the user interrupts (e.g., Ctrl+C), stop the container
		log.Infof("Received interrupt signal, stopping container %s", ctID)
		if err := stopContainer(dockerClient, ctID); err != nil {
			log.Errorf("Error stopping container: %v", err)
		}
		if err = deleteContainer(dockerClient, ctID); err != nil {
			return err
		}
		return ErrInterrupted

	case err := <-errChan:
		return err
	}

	return nil
}
