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
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/services"
	"github.com/NethermindEth/sedge/internal/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	log "github.com/sirupsen/logrus"
)

const SlashingImportFile string = "slashing_protection.json"

type SlashingImportOptions struct {
	ValidatorClient string
	Network         string
	StopValidator   bool
	StartValidator  bool
	GenerationPath  string
	From            string
	ContainerTag    string
}

func (s *sedgeActions) ImportSlashingInterchangeData(options SlashingImportOptions) error {
	validatorContainerName := services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
	consensusContainerName := services.ContainerNameWithTag(services.DefaultSedgeConsensusClient, options.ContainerTag)
	slashingContainerName := services.ContainerNameWithTag(services.ServiceCtSlashingData, options.ContainerTag)
	// Check validator container exists
	_, err := s.dockerServiceManager.ContainerID(validatorContainerName)
	if err != nil {
		return err
	}
	previouslyRunning, err := s.dockerServiceManager.IsRunning(validatorContainerName)
	if err != nil {
		return err
	}
	// Stop validator
	if previouslyRunning {
		log.Info("Stopping validator client...")
		if err := s.dockerServiceManager.Stop(validatorContainerName); err != nil {
			return err
		}
		log.Info("Validator client stopped.")
	}

	// Copy slashing data to generation path
	slashingDataPath := path.Join(options.GenerationPath, configs.ValidatorDir, SlashingImportFile)
	log.Debugf("Copying slashing data file from %s to %s", options.From, slashingDataPath)
	if err := utils.CopyFile(options.From, slashingDataPath); err != nil {
		return err
	}

	var cmd []string
	switch options.ValidatorClient {
	case "prysm":
		cmd = []string{
			"--", "slashing-protection-history",
			"import",
			"--accept-terms-of-use",
			"--" + options.Network,
			"--datadir=/data",
			"--slashing-protection-json-file=/data/slashing_protection.json",
		}
	case "lighthouse":
		cmd = []string{
			"lighthouse", "account", "validator", "slashing-protection", "import",
			"--network", options.Network,
			"--datadir", "/data",
			"/data/slashing_protection.json",
		}
	case "lodestar":
		cmd = []string{
			"validator", "slashing-protection", "import",
			"--network", options.Network,
			"--dataDir", "/data",
			"--file", "/data/slashing_protection.json",
		}
	case "teku":
		cmd = []string{
			"slashing-protection",
			"import",
			"--data-path=/data",
			"--from=/data/slashing_protection.json",
		}
	case "nimbus":
		cmd = []string{
			"slashingdb",
			"import",
			"/data/slashing_protection.json",
			"--validators-dir=/data",
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedValidatorClient, options.ValidatorClient)
	}
	log.Infof("Importing slashing data to client %s from %s", options.ValidatorClient, options.From)
	if err := runSlashingContainer(s.dockerClient, s.dockerServiceManager, cmd, validatorContainerName, consensusContainerName, slashingContainerName, options.ValidatorClient == "nimbus"); err != nil {
		return err
	}

	// Run validator again
	if (previouslyRunning && !options.StopValidator) || options.StartValidator {
		log.Info("The validator container is being restarted")
		if err := s.dockerServiceManager.Start(validatorContainerName); err != nil {
			return err
		}
		log.Info("Validator started.")
	}
	return nil
}

type SlashingExportOptions struct {
	ValidatorClient string
	Network         string
	StopValidator   bool
	StartValidator  bool
	GenerationPath  string
	Out             string
	ContainerTag    string
}

func (s *sedgeActions) ExportSlashingInterchangeData(options SlashingExportOptions) error {
	validatorContainerName := services.ContainerNameWithTag(services.DefaultSedgeValidatorClient, options.ContainerTag)
	consensusContainerName := services.ContainerNameWithTag(services.DefaultSedgeConsensusClient, options.ContainerTag)
	slashingContainerName := services.ContainerNameWithTag(services.ServiceCtSlashingData, options.ContainerTag)
	// Check validator container exists
	_, err := s.dockerServiceManager.ContainerID(validatorContainerName)
	if err != nil {
		return err
	}
	previouslyRunning, err := s.dockerServiceManager.IsRunning(validatorContainerName)
	if err != nil {
		return err
	}
	// Stop validator client
	if previouslyRunning {
		log.Info("Stopping validator client")
		if err := s.dockerServiceManager.Stop(validatorContainerName); err != nil {
			return err
		}
		log.Info("Validator client stopped.")
	}

	var cmd []string
	switch options.ValidatorClient {
	case "prysm":
		cmd = []string{
			"--", "slashing-protection-history",
			"export",
			"--accept-terms-of-use",
			"--" + options.Network,
			"--datadir=/data",
			"--slashing-protection-export-dir=/data",
		}
	case "lighthouse":
		cmd = []string{
			"lighthouse", "account", "validator", "slashing-protection", "export",
			"--network", options.Network,
			"--datadir", "/data",
			"/data/slashing_protection.json",
		}
	case "lodestar":
		cmd = []string{
			"validator", "slashing-protection", "export",
			"--network", options.Network,
			"--dataDir", "/data/validator",
			"--file", "/data/validator/slashing_protection.json",
		}
	case "teku":
		cmd = []string{
			"slashing-protection",
			"export",
			"--data-path=/data",
			"--to=/data/slashing_protection.json",
		}
	case "nimbus":
		cmd = []string{
			"slashingdb",
			"export",
			"/data/slashing_protection.json",
			"--validators-dir=/data",
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedValidatorClient, options.ValidatorClient)
	}
	log.Infof("Exporting slashing data from client %s", options.ValidatorClient)
	if err := runSlashingContainer(s.dockerClient, s.dockerServiceManager, cmd, validatorContainerName, consensusContainerName, slashingContainerName, options.ValidatorClient == "nimbus"); err != nil {
		return err
	}
	copyFrom := filepath.Join(options.GenerationPath, configs.ValidatorDir, "slashing_protection.json")
	log.Debugf("Copying slashing data file from %s to %s", copyFrom, options.Out)
	if err := utils.CopyFile(copyFrom, options.Out); err != nil {
		return err
	}

	// Run validator again
	if (previouslyRunning && !options.StopValidator) || options.StartValidator {
		log.Info("The validator container is being restarted...")
		if err := s.dockerServiceManager.Start(validatorContainerName); err != nil {
			return err
		}
		log.Info("Validator started.")
	}
	return nil
}

func runSlashingContainer(dockerClient client.APIClient, dockerServiceManager DockerServiceManager, cmd []string,
	validatorContainerName string, consensusContainerName string, slashingContainerName string, isNimbus bool,
) error {
	slashingImage := ""
	var err error
	if isNimbus {
		slashingImage, err = dockerServiceManager.Image(consensusContainerName)
		if err != nil {
			return err
		}
	} else {
		slashingImage, err = dockerServiceManager.Image(validatorContainerName)
		if err != nil {
			return err
		}
	}
	log.Debugf("Creating %s container", services.ServiceCtSlashingData)
	ct, err := dockerClient.ContainerCreate(context.Background(),
		&container.Config{
			Image: slashingImage,
			Cmd:   cmd,
		},
		&container.HostConfig{
			VolumesFrom: []string{validatorContainerName},
			LogConfig: container.LogConfig{
				Type: "json-file",
			},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		slashingContainerName,
	)
	if err != nil {
		return err
	}
	log.Debugf("Slashing protection container id: %s", ct.ID)
	ctExit, errChan := dockerServiceManager.Wait(slashingContainerName, container.WaitConditionNextExit)
	log.Info("The slashing protection container is starting...")
	if err := dockerClient.ContainerStart(context.Background(), ct.ID, container.StartOptions{}); err != nil {
		return err
	}
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)
	for {
		select {
		case exitResult := <-ctExit:
			logs, err := dockerServiceManager.ContainerLogs(ct.ID, "Slashing protection")
			if err != nil {
				return err
			}
			if err = deleteContainer(dockerClient, ct.ID); err != nil {
				return err
			}
			if exitResult.StatusCode != 0 {
				return fmt.Errorf(`%w: slashing protection container with id %s ends with status code %d. Here are the logs for more details: %s`, ErrValidatorImportCtBadExitCode, ct.ID, exitResult.StatusCode, logs)
			}
			log.Info("The slashing container ends successfully.")
			return nil
		case <-osSignals:
			if err := stopContainer(dockerClient, ct.ID); err != nil {
				return err
			}
			if err := deleteContainer(dockerClient, ct.ID); err != nil {
				return err
			}
			return ErrInterrupted
		case exitErr := <-errChan:
			return exitErr
		}
	}
}

func deleteContainer(dockerClient client.APIClient, containerVar string) error {
	log.Debugf("Removing container %s", containerVar)
	if err := dockerClient.ContainerRemove(context.Background(), containerVar, container.RemoveOptions{}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("error removing container %s: %w", containerVar, err)
	}
	return nil
}

func stopContainer(dockerClient client.APIClient, ctID string) error {
	log.Debugf("Stopping container %s", ctID)
	wait := 30
	if err := dockerClient.ContainerStop(context.Background(), ctID, container.StopOptions{Timeout: &wait}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("error stopping container %s: %w", ctID, err)
	}
	return nil
}
