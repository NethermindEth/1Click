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
	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RunCmd(cmdRunner commands.CommandRunner, sedgeActions actions.SedgeActions) *cobra.Command {
	// Flags
	var (
		generationPath string
		services       *[]string
	)

	cmd := &cobra.Command{
		Use:   "run [flags]",
		Short: "Run services",
		Long:  "Run all the generated services",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sedgeActions.InstallDependencies(actions.InstallDependenciesOptions{
				Dependencies: []string{"docker"},
				Install:      true,
			}); err != nil {
				log.Fatalf("dependency error: %s", err.Error())
			}
			err := sedgeActions.SetupContainers(actions.SetupContainersOptions{
				GenerationPath: generationPath,
				Services:       *services,
			})
			if err != nil {
				log.Fatalf("error setting up service containers: %s", err.Error())
			}
			err = sedgeActions.RunContainers(actions.RunContainersOptions{
				GenerationPath: generationPath,
				Services:       *services,
			})
			if err != nil {
				log.Fatalf("error starting service containers: %s", err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&generationPath, "path", "p", configs.DefaultDockerComposeScriptsPath, "docker-compose scripts generation path")
	services = cmd.Flags().StringSlice("services", []string{}, "List of services to run. If this flag is not provided, all services will run.")
	return cmd
}
