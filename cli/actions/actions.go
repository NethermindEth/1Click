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
	"github.com/NethermindEth/sedge/internal/compose"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/internal/pkg/generate"
	"github.com/docker/docker/client"
)

//go:generate mockgen -package=sedge_mocks -destination=../../mocks/sedgeActions.go github.com/NethermindEth/sedge/cli/actions SedgeActions
type SedgeActions interface {
	GetCommandRunner() commands.CommandRunner
	ImportSlashingInterchangeData(SlashingImportOptions) error
	ExportSlashingInterchangeData(SlashingExportOptions) error
	SetupContainers(SetupContainersOptions) error
	RunContainers(RunContainersOptions) error
	Generate(GenerateOptions) (generate.GenData, error)
	CreateJWTSecrets(CreateJWTSecretOptions) (string, error)
	ImportValidatorKeys(ImportValidatorKeysOptions) error
	ValidateDockerComposeFile(path string, services ...string) error
	GetContainersData(GetContainersDataOptions) (ContainersData, error)
}

type sedgeActions struct {
	dockerClient         client.APIClient
	dockerServiceManager DockerServiceManager
	commandRunner        commands.CommandRunner
	composeManager       compose.ComposeManager
}

type SedgeActionsOptions struct {
	DockerClient         client.APIClient
	DockerServiceManager DockerServiceManager
	CommandRunner        commands.CommandRunner
	ComposeManager       compose.ComposeManager
}

func NewSedgeActions(options SedgeActionsOptions) SedgeActions {
	return &sedgeActions{
		dockerClient:         options.DockerClient,
		dockerServiceManager: options.DockerServiceManager,
		commandRunner:        options.CommandRunner,
		composeManager:       options.ComposeManager,
	}
}

func (s *sedgeActions) GetCommandRunner() commands.CommandRunner {
	return s.commandRunner
}
