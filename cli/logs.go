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
	"strings"

	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/compose"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/internal/pkg/dependencies"
	"github.com/NethermindEth/sedge/internal/utils"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func LogsCmd(composeManager compose.ComposeManager, cmdRunner commands.CommandRunner, sedgeActions actions.SedgeActions, depsMgr dependencies.DependenciesManager) *cobra.Command {
	// Flags
	var (
		generationPath string
		tail           int
	)
	// Build command
	cmd := &cobra.Command{
		Use:   "logs [flags] [services]",
		Short: "Get running container logs",
		Long: `Get running container logs using docker-compose CLI. If no services are provided, the logs of all running services will be displayed.
	
	By default will run 'docker compose -f <script> logs --follow <service>'`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkDependencies(depsMgr, true, dependencies.Docker); err != nil {
				log.Error("Failed to check dependencies. Run 'sedge deps check' to check dependencies")
				return err
			}
			return sedgeActions.ValidateDockerComposeFile(filepath.Join(generationPath, configs.DefaultDockerComposeScriptName))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rawServices, err := utils.CheckContainers(cmdRunner, generationPath)
			if err != nil {
				return err
			}

			file := filepath.Join(generationPath, configs.DefaultDockerComposeScriptName)
			// Get logs from docker compose script services
			services := strings.Split(rawServices, "\n")
			// Remove empty string resulting of spliting the last blank line of rawServices
			services = services[:len(services)-1]
			if len(args) > 0 {
				services = args
			}

			err = composeManager.Logs(commands.DockerComposeLogsOptions{
				Path:     file,
				Services: services,
				Follow:   tail == 0,
				Tail:     tail,
			})
			if err != nil {
				return fmt.Errorf(configs.GettingLogsError, strings.Join(services, " "), err)
			}
			return nil
		},
	}
	// Bind flags
	cmd.Flags().StringVarP(&generationPath, "path", "p", configs.DefaultAbsSedgeDataPath, "generation path for sedge data")
	cmd.Flags().IntVarP(&tail, "tail", "t", 0, "Tail the number of desired logs. If not set, or set to 0, logs are followed.")

	return cmd
}
