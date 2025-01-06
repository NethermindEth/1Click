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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/dependencies"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")

func SlashingImportCmd(sedgeActions actions.SedgeActions, depsMgr dependencies.DependenciesManager) *cobra.Command {
	// Flags
	var (
		validatorClient string
		network         string
		stopValidator   bool
		startValidator  bool
		generationPath  string
		from            string
		containerTag    string
	)

	cmd := &cobra.Command{
		Use:   "slashing-import [flags] [validator]",
		Short: "Import slashing protection data",
		Long: `
Import Slashing Protection Interchange Format (EIP-3076) data. This command assumes
that the validator client container exists, stopped or not and that its database
is already initialized. The validator database is initialized if the validator is
running or has already run but is stopped, and also after importing the validator keys.
This command stops the validator client during the importing process due to the
validator database being locked while it's running but leaves the validator client
in the same state in which it was found. That means if the validator is running/stopped
before the import, then the validator will be running/stopped after the command
is executed, regardless of whether the export fails or not. To force a different
behavior use --start-validator and --stop-validator flags.
The [validator] is a required argument used to specify which validator client, from
all supported by Sedge (lighthouse, lodestar, prysm or teku), is used to import the
Slashing Protection data. This is necessary because each client has its own way to
achieve the importation.`,
		Example: `
sedge slashing-import --from slashing-data.json prysm
sedge slashing-import --from slashing-data.json --stop-validator lodestar 
sedge slashing-import --from slashing-data.json --start-validator lighthouse`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return ErrInvalidNumberOfArguments
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if f, err := os.Stat(generationPath); os.IsNotExist(err) || !f.IsDir() {
				return fmt.Errorf("generation path %s does not exist or is not a directory", generationPath)
			}
			if from == "" {
				from = filepath.Join(generationPath, "slashing_protection.json")
			} else {
				if f, err := os.Stat(from); os.IsNotExist(err) || f.IsDir() {
					return fmt.Errorf("slashing protection data file %s does not exist or is a directory", from)
				}
			}
			if err := configs.NetworkCheck(network); err != nil {
				return err
			}
			validatorClient = args[0]
			if err := checkDependencies(depsMgr, true, dependencies.Docker); err != nil {
				log.Error("Failed to check dependencies. Run 'sedge deps check' to check dependencies")
				return err
			}
			return sedgeActions.ValidateDockerComposeFile(filepath.Join(generationPath, configs.DefaultDockerComposeScriptName))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			services := []string{validator}
			// If the validator client is nimbus, we need to import the consensus service as well
			// The consensus service is not imported by default because it's not necessary for the other clients
			// On Nimbus, the consensus service is necessary to import the slashing protection data
			if validatorClient == "nimbus" {
				services = append(services, consensus)
			}
			err := sedgeActions.SetupContainers(actions.SetupContainersOptions{
				GenerationPath: generationPath,
				Services:       services,
			})
			if err != nil {
				return err
			}
			err = sedgeActions.ImportSlashingInterchangeData(actions.SlashingImportOptions{
				ValidatorClient: validatorClient,
				Network:         network,
				StopValidator:   stopValidator,
				StartValidator:  startValidator,
				GenerationPath:  generationPath,
				From:            from,
				ContainerTag:    containerTag,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&network, "network", "n", "mainnet", "network")
	cmd.Flags().BoolVar(&startValidator, "start-validator", false, "starts the validator client after import, regardless of the state the validator was in before")
	cmd.Flags().BoolVar(&stopValidator, "stop-validator", false, "stops the validator client after import, regardless of the state the validator was in before")
	cmd.Flags().StringVarP(&generationPath, "path", "p", configs.DefaultAbsSedgeDataPath, "path to the generation directory")
	cmd.Flags().StringVarP(&from, "from", "f", "", "path to the JSON file in the EIP-3076 format with the slashing protection data to import (default: <generation-dir>/slashing_protection.json)")
	cmd.PersistentFlags().StringVar(&containerTag, "container-tag", "", "Container tag to use. If defined, sedge will add to each container and the network, a suffix with the tag. e.g. sedge-validator-client -> sedge-validator-client-<tag>.")
	return cmd
}
