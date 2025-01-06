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

	"github.com/NethermindEth/sedge/cli/actions"
	sedgeOpts "github.com/NethermindEth/sedge/internal/pkg/options"
	"github.com/spf13/cobra"
)

var ErrCustomFlagsUsedWithoutCustomNetwork = errors.New("custom flags used without --network custom")

func validateCustomNetwork(flags *CustomFlags, net string) error {
	if net != "custom" {
		if len(flags.customChainSpec) != 0 || len(flags.customNetworkConfig) != 0 ||
			len(flags.customGenesis) != 0 || len(flags.customDeployBlock) != 0 {
			// TODO add error on expected place
			return ErrCustomFlagsUsedWithoutCustomNetwork
		}
	}
	return nil
}

func FullNodeSubCmd(sedgeAction actions.SedgeActions) *cobra.Command {
	var flags GenCmdFlags

	cmd := &cobra.Command{
		Use:   "full-node [flags]",
		Short: "Generate a full node config, with or without a validator",
		Long: `Generate a docker-compose and an environment file with a full node configuration.

It will not generate a validator configuration if the --no-validator flag is set to true.

On mainnet and sepolia mev-boost will be activated by default unless you run it with --no-mev-boost flag.

If you don't provide a execution, consensus or validator client, it will be chosen randomly. If one of the consensus or validator is provided, but the other one is omitted, then the same pair of clients will be used for both consensus and validator.

Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client, for example 'sedge generate full-node --execution nethermind:docker.image'. If you want to use the default docker image, just use the client name`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateCustomNetwork(&flags.CustomFlags, network); err != nil {
				return err
			}
			if err := preValidationGenerateCmd(network, logging, &flags); err != nil {
				return err
			}
			opts := sedgeOpts.CreateSedgeOptions(nodeType())
			settings := sedgeOpts.OptionSettings{
				Network:         network,
				MEVBoostEnabled: !flags.noMev && !flags.noValidator,
			}
			return opts.ValidateSettings(settings)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			services := []string{execution, consensus}
			if !flags.noValidator {
				services = append(services, validator)
			}
			if !flags.noMev && !flags.noValidator {
				services = append(services, mevBoost)
			}
			if flags.consensusName == "" {
				flags.consensusName = flags.validatorName
			} else if flags.validatorName == "" {
				flags.validatorName = flags.consensusName
			}
			if flags.distributed {
				services = append(services, distributedValidator)
			}
			return runGenCmd(cmd.OutOrStdout(), &flags, sedgeAction, services)
		},
	}
	// Bind flags
	cmd.Flags().StringVarP(&flags.consensusName, "consensus", "c", "", "Consensus engine client, e.g. teku, lodestar, prysm, lighthouse, Nimbus. Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client. If you want to use the default docker image, just use the client name")
	cmd.Flags().StringVarP(&flags.executionName, "execution", "e", "", "Execution engine client, e.g. geth, nethermind, besu, erigon. Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client. If you want to use the default docker image, just use the client name")
	cmd.Flags().StringVarP(&flags.validatorName, "validator", "v", "", "Validator engine client, e.g. teku, lodestar, prysm, lighthouse, Nimbus. Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client. If you want to use the default docker image, just use the client name")
	cmd.Flags().BoolVar(&flags.distributed, "distributed", false, "Deploy a node configured to run as part of a Distributed Validator Cluster.")
	cmd.Flags().StringVarP(&flags.distributedValidatorName, "distributedValidator", "d", "", "Distributed Validator client, e.g. charon. Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client. If you want to use the default docker image, just use the client name")
	cmd.Flags().BoolVar(&flags.latestVersion, "latest", false, "Use the latest version of clients. This sets the \"latest\" tag on the client's docker images. Latest version might not work.")
	cmd.Flags().StringVar(&flags.checkpointSyncUrl, "checkpoint-sync-url", "", "Initial state endpoint (trusted synced consensus endpoint) for the consensus client to sync from a finalized checkpoint. Provide faster sync process for the consensus client and protect it from long-range attacks affored by Weak Subjetivity. Each network has a default checkpoint sync url.")
	cmd.Flags().StringVar(&flags.feeRecipient, "fee-recipient", "", "Suggested fee recipient. Is a 20-byte Ethereum address which the execution layer might choose to set as the coinbase and the recipient of other fees or rewards. There is no guarantee that an execution node will use the suggested fee recipient to collect fees, it may use any address it chooses. It is assumed that an honest execution node will use the suggested fee recipient, but users should note this trust assumption.\n"+
		"Note: When setting up a Lido node, fee recipient address will be automatically set by the system.")
	cmd.Flags().BoolVar(&flags.noMev, "no-mev-boost", false, "Not use mev-boost if supported")
	cmd.Flags().StringVarP(&flags.mevImage, "mev-boost-image", "m", "", "Custom docker image to use for Mev Boost. Example: 'sedge generate full-node --mev-boost-image flashbots/mev-boost:latest-portable'")
	cmd.Flags().StringSliceVar(&flags.relayURLs, "relay-urls", []string{}, "List of comma separated relay URLs used to connect to mev relay. Example: 'sedge generate full-node --relay-urls=https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net,https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money'\n"+
		"Note: When setting up a Lido node, the provided relay URLs will be automatically set by the system.")
	cmd.Flags().BoolVar(&flags.noValidator, "no-validator", false, "Exclude the validator from the full node setup. Designed for execution and consensus nodes setup without a validator node. Exclude also the validator from other flags. If set, mev-boost will not be used.")
	cmd.Flags().StringVar(&flags.jwtPath, "jwt-secret-path", "", "Path to the JWT secret file")
	cmd.Flags().StringVar(&flags.graffiti, "graffiti", "", "Graffiti to be used by the validator")
	cmd.Flags().BoolVar(&flags.mapAllPorts, "map-all", false, "Map all clients ports to host. Use with care. Useful to allow remote access to the clients")
	cmd.Flags().StringSliceVar(&flags.fallbackEL, "fallback-execution-urls", []string{}, "Fallback/backup execution endpoints for the consensus client. Not supported by Teku. Example: 'sedge generate full-node -r --fallback-execution=https://mainnet.infura.io/v3/YOUR-PROJECT-ID,https://eth-mainnet.alchemyapi.io/v2/YOUR-PROJECT-ID'")
	cmd.Flags().StringArrayVar(&flags.elExtraFlags, "el-extra-flag", []string{}, "Additional flag to configure the execution client service in the generated docker-compose script. Example: 'sedge generate full-node --el-extra-flag \"<flag1>=value1\" --el-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringArrayVar(&flags.clExtraFlags, "cl-extra-flag", []string{}, "Additional flag to configure the consensus client service in the generated docker-compose script. Example: 'sedge generate full-node --cl-extra-flag \"<flag1>=value1\" --cl-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringArrayVar(&flags.vlExtraFlags, "vl-extra-flag", []string{}, "Additional flag to configure the validator client service in the generated docker-compose script. Example: 'sedge generate full-node --vl-extra-flag \"<flag1>=value1\" --vl-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringArrayVar(&flags.dvExtraFlags, "dv-extra-flag", []string{}, "Additional flag to configure the distributed validator client service in the generated docker-compose script. Example: 'sedge generate full-node --distributed --dv-extra-flag \"<flag1>=value1\" --dv-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringVar(&flags.customChainSpec, "custom-chainSpec", "", "File path or url to use as custom network chainSpec for execution client.")
	cmd.Flags().StringVar(&flags.customNetworkConfig, "custom-config", "", "File path or url to use as custom network config file for consensus client.")
	cmd.Flags().StringVar(&flags.customGenesis, "custom-genesis", "", "File path or url to use as custom network genesis for consensus client.")
	cmd.Flags().StringVar(&flags.customDeployBlock, "custom-deploy-block", "", "Custom network deploy block to use for consensus client.")
	cmd.Flags().IntVar(&flags.waitEpoch, "wait-epoch", 1, "Number of epochs to wait before starting and restarting of the validator client.")
	cmd.Flags().StringSliceVar(&flags.customEnodes, "execution-bootnodes", []string{}, "List of comma separated enodes to use as custom network peers for execution client.")
	cmd.Flags().StringSliceVar(&flags.customEnrs, "consensus-bootnodes", []string{}, "List of comma separated enrs to use as custom network peers for consensus client.")
	cmd.Flags().SortFlags = false
	return cmd
}

func OpFullNodeSubCmd(sedgeAction actions.SedgeActions) *cobra.Command {
	var flags GenCmdFlags
	cmd := &cobra.Command{
		Use:   "op-full-node [flags]",
		Short: "Generate a full node config for Optimism or Base",
		Long: `Generate a docker-compose and an environment file with a full node configuration for Optimism or Base networks.

This command sets up an Optimism or Base full node, which includes an execution client, a consensus client, a Optimism consensus client, and an Optimism node.

If you don't provide images for your clients, they will be chosen randomly. You can specify custom images for the Optimism and other nodes.

Use the --base flag to generate a configuration for a Base node (which is built on Optimism).

The command allows you to use external execution and consensus APIs instead of running your own nodes, by providing the respective URLs.

Additionally, you can use the syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client, for example 'sedge generate op-full-node --execution nethermind:custom.image'. If you want to use the default docker image, just use the client name.

This command does not generate a validator configuration, as Optimism and Base use different validation mechanisms compared to standard Ethereum networks.`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateCustomNetwork(&flags.CustomFlags, network); err != nil {
				return err
			}
			return preValidationGenerateCmd(network, logging, &flags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			services := []string{execution, consensus, optimism, opExecution}
			return runGenCmd(cmd.OutOrStdout(), &flags, sedgeAction, services)
		},
	}
	// Bind flags
	cmd.Flags().StringVarP(&flags.consensusName, "consensus", "c", "", "Consensus engine client, e.g. teku, lodestar, prysm, lighthouse, Nimbus. Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client. If you want to use the default docker image, just use the client name")
	cmd.Flags().StringVar(&flags.optimismName, "op-image", "", "Optimism consensus client image.")
	cmd.Flags().StringVar(&flags.optimismExecutionName, "op-execution", "", "Optimism Execution client to be used, op-nethermind, op-geth, or op-reth. Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client. If you want to use the default docker image, just use the client name.")
	cmd.Flags().StringVarP(&flags.executionName, "execution", "e", "", "Execution engine client, e.g. geth, nethermind, besu, erigon. Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client. If you want to use the default docker image, just use the client name")
	cmd.Flags().StringVarP(&flags.executionApiUrl, "execution-api-url", "", "", "Set execution api url. If Set, will omit the creation of execution and beacon nodes, and only create optimism nodes.")
	cmd.Flags().StringVarP(&flags.consensusApiUrl, "consensus-url", "", "", "Set consensus api url. If Set, will omit the creation of execution and beacon nodes, and only create optimism nodes.")
	cmd.Flags().BoolVar(&flags.latestVersion, "latest", false, "Use the latest version of clients. This sets the \"latest\" tag on the client's docker images. Latest version might not work.")
	cmd.Flags().StringVar(&flags.checkpointSyncUrl, "checkpoint-sync-url", "", "Initial state endpoint (trusted synced consensus endpoint) for the consensus client to sync from a finalized checkpoint. Provide faster sync process for the consensus client and protect it from long-range attacks affored by Weak Subjetivity. Each network has a default checkpoint sync url.")
	cmd.Flags().StringVar(&flags.feeRecipient, "fee-recipient", "", "Suggested fee recipient. Is a 20-byte Ethereum address which the execution layer might choose to set as the coinbase and the recipient of other fees or rewards. There is no guarantee that an execution node will use the suggested fee recipient to collect fees, it may use any address it chooses. It is assumed that an honest execution node will use the suggested fee recipient, but users should note this trust assumption")
	cmd.Flags().StringVar(&flags.jwtPath, "jwt-secret-path", "", "Path to the JWT secret file")
	cmd.Flags().BoolVar(&flags.mapAllPorts, "map-all", false, "Map all clients ports to host. Use with care. Useful to allow remote access to the clients")
	cmd.Flags().BoolVar(&flags.isBase, "base", false, "If set, will generate the docker-compose file for Base L2 config")
	cmd.Flags().StringSliceVar(&flags.fallbackEL, "fallback-execution-urls", []string{}, "Fallback/backup execution endpoints for the consensus client. Not supported by Teku. Example: 'sedge generate full-node -r --fallback-execution=https://mainnet.infura.io/v3/YOUR-PROJECT-ID,https://eth-mainnet.alchemyapi.io/v2/YOUR-PROJECT-ID'")
	cmd.Flags().StringArrayVar(&flags.elExtraFlags, "el-extra-flag", []string{}, "Additional flag to configure the execution client service in the generated docker-compose script. Example: 'sedge generate full-node --el-extra-flag \"<flag1>=value1\" --el-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringArrayVar(&flags.elOpExtraFlags, "el-op-extra-flag", []string{}, "Additional flag to configure the execution client for optimism service in the generated docker-compose script. Example: 'sedge generate full-node --el-extra-flag \"<flag1>=value1\" --el-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringArrayVar(&flags.opExtraFlags, "op-extra-flag", []string{}, "Additional flag to configure the optimism client service in the generated docker-compose script. Example: 'sedge generate full-node --el-extra-flag \"<flag1>=value1\" --el-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringArrayVar(&flags.clExtraFlags, "cl-extra-flag", []string{}, "Additional flag to configure the consensus client service in the generated docker-compose script. Example: 'sedge generate full-node --cl-extra-flag \"<flag1>=value1\" --cl-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringSliceVar(&flags.customEnodes, "execution-bootnodes", []string{}, "List of comma separated enodes to use as custom network peers for execution client.")
	cmd.Flags().StringSliceVar(&flags.customEnrs, "consensus-bootnodes", []string{}, "List of comma separated enrs to use as custom network peers for consensus client.")

	cmd.Flags().SortFlags = false
	return cmd
}

func ExecutionSubCmd(sedgeAction actions.SedgeActions) *cobra.Command {
	var flags GenCmdFlags

	cmd := &cobra.Command{
		Use:   "execution [flags] [args]",
		Short: "Generate a execution node config",
		Long: "Generate a docker-compose and an environment file with a execution node configuration.\n" +
			"Valid args: name of execution clients according to network\n\n" +
			"Should be one of: nethermind, geth, besu, erigon. If you don't provide one, it will chosen randomly.\n" +
			"Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client, for example 'sedge generate execution nethermind:docker.image'. If you want to use the default docker image, just use the client name",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if cobra.ExactArgs(1)(cmd, args) != nil {
					return errors.New("requires one argument")
				}
				flags.executionName = args[0]
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateCustomNetwork(&flags.CustomFlags, network); err != nil {
				return err
			}
			if err := preValidationGenerateCmd(network, logging, &flags); err != nil {
				return err
			}
			opts := sedgeOpts.CreateSedgeOptions(nodeType())
			settings := sedgeOpts.OptionSettings{
				Network:         network,
				MEVBoostEnabled: false, // MEV Boost is not supported for execution nodes only
			}
			return opts.ValidateSettings(settings)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenCmd(cmd.OutOrStdout(), &flags, sedgeAction, []string{execution})
		},
	}

	// Bind flags
	cmd.Flags().BoolVar(&flags.latestVersion, "latest", false, "Use the latest version of clients. This sets the \"latest\" tag on the client's docker images. Latest version might not work.")
	cmd.Flags().StringVar(&flags.jwtPath, "jwt-secret-path", "", "Path to the JWT secret file")
	cmd.Flags().BoolVar(&flags.mapAllPorts, "map-all", false, "Map all clients ports to host. Use with care. Useful to allow remote access to the clients")
	cmd.Flags().StringVar(&flags.customChainSpec, "custom-chainSpec", "", "File path or url to use as custom network chainSpec for execution client.")
	cmd.Flags().StringSliceVar(&flags.customEnodes, "execution-bootnodes", []string{}, "List of comma separated enodes to use as custom network peers for execution client.")
	cmd.Flags().StringArrayVar(&flags.elExtraFlags, "el-extra-flag", []string{}, "Additional flag to configure the execution client service in the generated docker-compose script. Example: 'sedge generate consensus--el-extra-flag \"<flag1>=value1\" --el-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().SortFlags = false
	return cmd
}

func ConsensusSubCmd(sedgeAction actions.SedgeActions) *cobra.Command {
	var flags GenCmdFlags

	cmd := &cobra.Command{
		Use:   "consensus [flags] --execution-api-url <URL> --execution-auth-url <URL> [args]",
		Short: "Generate a consensus node config",
		Long: "Generate a docker-compose and an environment file with a consensus node configuration\n" +
			"Valid args: name of execution clients according to network\n\n" +
			"Should be one of: lighthouse, teku, prysm, lodestar. If you don't provide one, it will chosen randomly.\n" +
			"Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client, for example 'sedge generate consensus prysm:docker.image'. If you want to use the default docker image, just use the client name" +
			"\n\n" +
			"Required flags:\n" +
			"- '--execution-api-url'\n" +
			"- '--execution-auth-url'",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if cobra.ExactArgs(1)(cmd, args) != nil {
					return errors.New("requires one argument")
				}
				flags.consensusName = args[0]
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateCustomNetwork(&flags.CustomFlags, network); err != nil {
				return err
			}
			if err := preValidationGenerateCmd(network, logging, &flags); err != nil {
				return err
			}
			opts := sedgeOpts.CreateSedgeOptions(nodeType())
			settings := sedgeOpts.OptionSettings{
				Network:         network,
				MEVBoostEnabled: false, // MEV Boost is not supported for consensus nodes only
			}
			return opts.ValidateSettings(settings)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenCmd(cmd.OutOrStdout(), &flags, sedgeAction, []string{consensus})
		},
	}
	// Bind flags
	cmd.Flags().BoolVar(&flags.latestVersion, "latest", false, "Use the latest version of clients. This sets the \"latest\" tag on the client's docker images. Latest version might not work.")
	cmd.Flags().StringVar(&flags.executionApiUrl, "execution-api-url", "", "Execution API endpoint for the consensus client. Example: 'sedge generate consensus -r --execution-api-url=https://api.url.endpoint'")
	cmd.Flags().StringVar(&flags.executionAuthUrl, "execution-auth-url", "", "Execution AUTH endpoint for the consensus client. Example: 'sedge generate consensus -r --execution-auth-url=https://auth.url.endpoint'")
	cmd.Flags().StringVar(&flags.checkpointSyncUrl, "checkpoint-sync-url", "", "Initial state endpoint (trusted synced consensus endpoint) for the consensus client to sync from a finalized checkpoint. Provide faster sync process for the consensus client and protect it from long-range attacks affored by Weak Subjetivity. Each network has a default checkpoint sync url.")
	cmd.Flags().StringVar(&flags.feeRecipient, "fee-recipient", "", "Suggested fee recipient. Is a 20-byte Ethereum address which the execution layer might choose to set as the coinbase and the recipient of other fees or rewards. There is no guarantee that an execution node will use the suggested fee recipient to collect fees, it may use any address it chooses. It is assumed that an honest execution node will use the suggested fee recipient, but users should note this trust assumption.\n"+
		"Note: When setting up a Lido node, fee recipient address will be automatically set by the system.")
	cmd.Flags().StringVar(&flags.jwtPath, "jwt-secret-path", "", "Path to the JWT secret file")
	cmd.Flags().StringVar(&flags.mevBoostUrl, "mev-boost-url", "", "If you are running a mev boost node, and you want to connect to it, you need to set mev-boost-url, if not set, node will not load any mev boost related config.")
	cmd.Flags().BoolVar(&flags.mapAllPorts, "map-all", false, "Map all clients ports to host. Use with care. Useful to allow remote access to the clients")
	cmd.Flags().StringSliceVar(&flags.fallbackEL, "fallback-execution-urls", []string{}, "Fallback/backup execution endpoints for the consensus client. Not supported by Teku. Example: 'sedge generate consensus --fallback-execution=https://mainnet.infura.io/v3/YOUR-PROJECT-ID,https://eth-mainnet.alchemyapi.io/v2/YOUR-PROJECT-ID'")
	cmd.Flags().StringArrayVar(&flags.clExtraFlags, "cl-extra-flag", []string{}, "Additional flag to configure the consensus client service in the generated docker-compose script. Example: 'sedge generate consensus --cl-extra-flag \"<flag1>=value1\" --cl-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	cmd.Flags().StringVar(&flags.customNetworkConfig, "custom-config", "", "File path or url to use as custom network config file for consensus client.")
	cmd.Flags().StringVar(&flags.customGenesis, "custom-genesis", "", "File path or url to use as custom network genesis for consensus client.")
	cmd.Flags().StringVar(&flags.customDeployBlock, "custom-deploy-block", "", "Custom network deploy block to use for consensus client.")
	cmd.Flags().StringSliceVar(&flags.customEnrs, "consensus-bootnodes", []string{}, "List of comma separated enrs to use as custom network peers for consensus client.")
	err := cmd.MarkFlagRequired("execution-api-url")
	if err != nil {
		return nil
	}
	err = cmd.MarkFlagRequired("execution-auth-url")
	if err != nil {
		return nil
	}
	cmd.Flags().SortFlags = false
	return cmd
}

func ValidatorSubCmd(sedgeAction actions.SedgeActions) *cobra.Command {
	var flags GenCmdFlags

	cmd := &cobra.Command{
		Use:   "validator [flags] --consensus-url <URL> [args]",
		Short: "Generate a validator node config",
		Long: "Generate a docker-compose and an environment file with a validator node configuration\n" +
			"Valid args: name of execution clients according to network\n\n" +
			"Should be one of: lighthouse, teku, prysm, lodestar. If you don't provide one, it will chosen randomly.\n" +
			"Additionally, you can use this syntax '<CLIENT>:<DOCKER_IMAGE>' to override the docker image used for the client, for example 'sedge generate validator prysm:docker.image'. If you want to use the default docker image, just use the client name" +
			"\n\n" +
			"Required flags:\n" +
			"- `--consensus-url`",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if cobra.ExactArgs(1)(cmd, args) != nil {
					return errors.New("requires one argument")
				}
				flags.validatorName = args[0]
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateCustomNetwork(&flags.CustomFlags, network); err != nil {
				return err
			}
			if err := preValidationGenerateCmd(network, logging, &flags); err != nil {
				return err
			}
			opts := sedgeOpts.CreateSedgeOptions(nodeType())
			settings := sedgeOpts.OptionSettings{
				Network:         network,
				MEVBoostEnabled: !flags.noMev && !flags.noValidator,
			}
			return opts.ValidateSettings(settings)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.noMev = true
			return runGenCmd(cmd.OutOrStdout(), &flags, sedgeAction, []string{validator})
		},
	}
	// Bind flags
	cmd.Flags().BoolVar(&flags.latestVersion, "latest", false, "Use the latest version of clients. This sets the \"latest\" tag on the client's docker images. Latest version might not work.")
	cmd.Flags().StringVar(&flags.consensusApiUrl, "consensus-url", "", "Consensus endpoint for the validator client to connect to. Example: 'sedge generate validator --consensus-url http://localhost:4000'")
	cmd.Flags().StringVar(&flags.feeRecipient, "fee-recipient", "", "Suggested fee recipient. Is a 20-byte Ethereum address which the execution layer might choose to set as the coinbase and the recipient of other fees or rewards. There is no guarantee that an execution node will use the suggested fee recipient to collect fees, it may use any address it chooses. It is assumed that an honest execution node will use the suggested fee recipient, but users should note this trust assumption.\n"+
		"Note: When setting up a Lido node, fee recipient address will be automatically set by the system.")
	cmd.Flags().StringVar(&flags.graffiti, "graffiti", "", "Graffiti to be used by the validator")
	cmd.Flags().BoolVar(&flags.mevBoostOnVal, "mev-boost", false, "Use mev-boost while turning on validator node")
	cmd.Flags().StringVar(&flags.customNetworkConfig, "custom-config", "", "File path or url to use as custom network config file for consensus client.")
	cmd.Flags().StringVar(&flags.customGenesis, "custom-genesis", "", "File path or url to use as custom network genesis for consensus client.")
	cmd.Flags().StringVar(&flags.customDeployBlock, "custom-deploy-block", "", "Custom network deploy block to use for consensus client.")
	cmd.Flags().IntVar(&flags.waitEpoch, "wait-epoch", 1, "Number of epochs to wait before starting and restarting of the validator client.")
	cmd.Flags().StringArrayVar(&flags.vlExtraFlags, "vl-extra-flag", []string{}, "Additional flag to configure the validator client service in the generated docker-compose script. Example: 'sedge generate validator --vl-extra-flag \"<flag1>=value1\" --vl-extra-flag \"<flag2>=\\\"value2\\\"\"'")
	err := cmd.MarkFlagRequired("consensus-url")
	if err != nil {
		return nil
	}
	cmd.Flags().SortFlags = false
	return cmd
}

func MevBoostSubCmd(sedgeAction actions.SedgeActions) *cobra.Command {
	var flags GenCmdFlags

	cmd := &cobra.Command{
		Use:   "mev-boost [flags]",
		Short: "Generate a mev-boost node config",
		Long:  "Generate a docker-compose and an environment file with a mev-boost node configuration",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preValidationGenerateCmd(network, logging, &flags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenCmd(cmd.OutOrStdout(), &flags, sedgeAction, []string{mevBoost})
		},
	}
	// Bind flags
	cmd.Flags().StringSliceVar(&flags.relayURLs, "relay-urls", []string{}, "List of comma separated relay URLs used to connect to mev relay. Example: 'sedge generate full-node --relay-urls=https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net,https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money'\n"+
		"Note: When setting up a Lido node, the provided relay URLs will be automatically set by the system.")
	cmd.Flags().StringVarP(&flags.mevImage, "mev-boost-image", "m", "", "Custom docker image to use for Mev Boost. Example: 'sedge generate mev-boost --mev-boost-image flashbots/mev-boost:latest-portable'")
	cmd.Flags().StringVarP(&network, "network", "n", "mainnet", "Target network. e.g. mainnet, sepolia etc.")
	cmd.Flags().SortFlags = false
	return cmd
}
