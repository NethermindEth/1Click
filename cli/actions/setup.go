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
	"fmt"
	"path/filepath"

	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	log "github.com/sirupsen/logrus"
)

type SetupContainersOptions struct {
	GenerationPath string
	Services       []string
	SkipPull       bool
}

func (s *sedgeActions) SetupContainers(options SetupContainersOptions) error {
	log.Info("Setting up containers")
	err := s.composeManager.Build(commands.DockerComposeBuildOptions{
		Path:     filepath.Join(options.GenerationPath, configs.DefaultDockerComposeScriptName),
		Services: options.Services,
	})
	if err != nil {
		return fmt.Errorf(configs.SetUpContainersErr, err)
	}
	if !options.SkipPull {
		err := s.composeManager.Pull(commands.DockerComposePullOptions{
			Path:     filepath.Join(options.GenerationPath, configs.DefaultDockerComposeScriptName),
			Services: options.Services,
		})
		if err != nil {
			return fmt.Errorf(configs.SetUpContainersErr, err)
		}
	} else {
		log.Warn("Skipping 'docker compose pull' step")
	}
	err = s.composeManager.Create(commands.DockerComposeCreateOptions{
		Path:     filepath.Join(options.GenerationPath, configs.DefaultDockerComposeScriptName),
		Services: options.Services,
	})
	if err != nil {
		return fmt.Errorf(configs.SetUpContainersErr, err)
	}
	return nil
}
