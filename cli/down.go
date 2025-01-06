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
	"fmt"
	"path/filepath"

	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/compose"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/internal/pkg/dependencies"
	"github.com/NethermindEth/sedge/internal/utils"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func DownCmd(composeManager compose.ComposeManager, cmdRunner commands.CommandRunner, a actions.SedgeActions, depsMgr dependencies.DependenciesManager) *cobra.Command {
	// Flags
	var generationPath string
	// Build command
	cmd := &cobra.Command{
		Use:   "down [flags]",
		Short: "Shutdown sedge running containers",
		Long:  `Shutdown sedge running containers using docker compose CLI. Shortcut for 'docker compose -f <script> down'`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkDependencies(depsMgr, true, dependencies.Docker); err != nil {
				log.Error("Failed to check dependencies. Run 'sedge deps check' to check dependencies")
				return err
			}
			return a.ValidateDockerComposeFile(filepath.Join(generationPath, configs.DefaultDockerComposeScriptName))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := utils.CheckContainers(cmdRunner, generationPath); err != nil {
				return err
			}

			err := composeManager.Down(commands.DockerComposeDownOptions{
				Path: filepath.Join(generationPath, configs.DefaultDockerComposeScriptName),
			})
			if err != nil {
				return fmt.Errorf("error shutting down contaiers %w", err)
			}

			return nil
		},
	}
	// Bind flags
	cmd.Flags().StringVarP(&generationPath, "path", "p", configs.DefaultAbsSedgeDataPath, "generation path for sedge data")
	return cmd
}
