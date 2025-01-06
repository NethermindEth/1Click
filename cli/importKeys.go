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
	"path/filepath"

	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/dependencies"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ImportKeysCmd(sedgeActions actions.SedgeActions, depsMgr dependencies.DependenciesManager) *cobra.Command {
	// Flags
	var (
		validatorClient   string
		network           string
		stopValidator     bool
		startValidator    bool
		from              string
		generationPath    string
		customConfigPath  string
		customGenesisPath string
		customDeployBlock string
		distributed       bool
	)

	cmd := &cobra.Command{
		Use:   "import-key [flags] [validator]",
		Short: "Import validator keys",
		Long: `
Import validator client keys, use the 'from' flag to specify the keys location,
and make sure that follows the EIP-2335: BLS12-381 Keystore standard. This command
assumes that the validator client container exists, stopped or not.

This command stops the validator client during the importing process due to the
validator database being locked while it's running but leaves the validator client
in the same state in which it was found. That means if the validator is running/stopped
before the import, then the validator will be running/stopped after the command
is executed, regardless of whether the export fails or not. To force a different
behavior use --start-validator and --stop-validator flags.

The [validator] is a required argument used to specify which validator client, from
all supported by Sedge (lighthouse, lodestar, prysm or teku), is used to import the
validator keys. This is necessary because each client has its own way to achieve
the importation.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrInvalidNumberOfArguments
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			validatorClient = args[0]
			if err := checkDependencies(depsMgr, true, dependencies.Docker); err != nil {
				log.Error("Failed to check dependencies. Run 'sedge deps check' to check dependencies")
				return err
			}
			return sedgeActions.ValidateDockerComposeFile(filepath.Join(generationPath, configs.DefaultDockerComposeScriptName))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options := actions.ImportValidatorKeysOptions{
				ValidatorClient: validatorClient,
				Network:         network,
				StopValidator:   stopValidator,
				StartValidator:  startValidator,
				From:            from,
				GenerationPath:  generationPath,
				ContainerTag:    containerTag,
				CustomConfig: actions.ImportValidatorKeysCustomOptions{
					NetworkConfigPath: customConfigPath,
					GenesisPath:       customGenesisPath,
					DeployBlockPath:   customDeployBlock,
				},
				Distributed: distributed,
			}
			var err error
			if validatorClient == "nimbus" {
				err = sedgeActions.SetupContainers(actions.SetupContainersOptions{
					GenerationPath: generationPath,
					Services:       []string{validator, consensus},
				})
			} else {
				err = sedgeActions.SetupContainers(actions.SetupContainersOptions{
					GenerationPath: generationPath,
					Services:       []string{validator},
				})
			}
			if err != nil {
				return err
			}
			return sedgeActions.ImportValidatorKeys(options)
		},
	}

	cmd.Flags().StringVarP(&network, "network", "n", "mainnet", "network")
	cmd.Flags().StringVarP(&generationPath, "path", "p", configs.DefaultAbsSedgeDataPath, "path to the generation directory")
	cmd.Flags().BoolVar(&startValidator, "start-validator", false, "starts the validator client after import, regardless of the state the validator was in before")
	cmd.Flags().BoolVar(&stopValidator, "stop-validator", false, "stops the validator client after import, regardless of the state the validator was in before")
	cmd.Flags().StringVar(&from, "from", filepath.Join(configs.DefaultAbsSedgeDataPath, "keystore"), "path to the validator keys, must follow the EIP-2335: BLS12-381 Keystore standard")
	cmd.PersistentFlags().StringVar(&containerTag, "container-tag", "", "Container tag to use. If defined, sedge will add to each container and the network, a suffix with the tag. e.g. sedge-validator-client -> sedge-validator-client-<tag>.")
	cmd.Flags().StringVar(&customConfigPath, "custom-config", "", "file path or url to use as custom network config.")
	cmd.Flags().StringVar(&customGenesisPath, "custom-genesis", "", "file path or url to use as custom network genesis.")
	cmd.Flags().StringVar(&customDeployBlock, "custom-deploy-block", "", "custom network deploy block.")
	cmd.Flags().BoolVar(&distributed, "distributed", false, "Import keys generated by Distributed Key Generation (DKG) process")
	return cmd
}
