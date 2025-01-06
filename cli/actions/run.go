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
)

type RunContainersOptions struct {
	GenerationPath string
	Services       []string
	SkipDockerPs   bool
}

func (s *sedgeActions) RunContainers(options RunContainersOptions) error {
	err := s.composeManager.Up(commands.DockerComposeUpOptions{
		Path:     filepath.Join(options.GenerationPath, configs.DefaultDockerComposeScriptName),
		Services: options.Services,
	})
	if err != nil {
		return fmt.Errorf(configs.RunContainersErr, err)
	}
	if !options.SkipDockerPs {
		// Run docker compose ps --filter status=running to show script running containers
		_, err := s.composeManager.PS(commands.DockerComposePsOptions{
			Path:          filepath.Join(options.GenerationPath, configs.DefaultDockerComposeScriptName),
			FilterRunning: true,
		})
		if err != nil {
			return fmt.Errorf(configs.RunContainersErr, err)
		}
	}
	return nil
}
