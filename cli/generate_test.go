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
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/lido/contracts"
	"github.com/NethermindEth/sedge/internal/lido/contracts/mevboostrelaylist"
	"github.com/NethermindEth/sedge/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type subCmd struct {
	name string
	args []string
}

func (cmd subCmd) argsList() []string {
	s := make([]string, 0)
	s = append(s, cmd.name)
	s = append(s, cmd.args...)
	return s
}

type globalFlags struct {
	install        bool
	generationPath string
	network        string
	logging        string
	lidoNode       bool
}

type generateCmdTestCase struct {
	name       string
	subCommand subCmd
	args       GenCmdFlags
	globalArgs globalFlags
	err        error
}

func (flags *globalFlags) argsList() []string {
	s := make([]string, 0)
	if flags.install {
		s = append(s, "--install")
	}
	if flags.generationPath != "" {
		s = append(s, "--path", flags.generationPath)
	}
	if flags.network != "" {
		s = append(s, "--network", flags.network)
	}
	if flags.logging != "" {
		s = append(s, "--logging", flags.logging)
	}
	if flags.lidoNode {
		s = append(s, "--lido")
	}
	return s
}

func (flags *GenCmdFlags) argsList() []string {
	s := make([]string, 0)
	if flags.executionName != "" {
		s = append(s, "-e", flags.executionName)
	}
	if flags.consensusName != "" {
		s = append(s, "-c", flags.consensusName)
	}
	if flags.validatorName != "" {
		s = append(s, "-v", flags.validatorName)
	}
	if flags.checkpointSyncUrl != "" {
		s = append(s, "--checkpoint-sync-url", flags.checkpointSyncUrl)
	}
	if flags.executionAuthUrl != "" {
		s = append(s, "--execution-auth-url", flags.executionAuthUrl)
	}
	if flags.executionApiUrl != "" {
		s = append(s, "--execution-api-url", flags.executionApiUrl)
	}

	if flags.noMev {
		s = append(s, "--no-mev-boost")
	}
	if flags.noValidator {
		s = append(s, "--no-validator")
	}
	if flags.mevImage != "" {
		s = append(s, "--mev-boost-image", flags.mevImage)
	}
	if len(flags.relayURLs) != 0 {
		s = append(s, "--relay-urls", strings.Join(flags.relayURLs, ","))
	}
	if flags.mevBoostUrl != "" {
		s = append(s, "--mev-boost-url", flags.mevBoostUrl)
	}
	if flags.jwtPath != "" {
		s = append(s, "--jwt-secret-path", flags.jwtPath)
	}
	if flags.graffiti != "" {
		s = append(s, "--graffiti", flags.graffiti)
	}
	if flags.consensusApiUrl != "" {
		s = append(s, "--consensus-url", flags.consensusApiUrl)
	}
	if flags.feeRecipient != "" {
		s = append(s, "--fee-recipient", flags.feeRecipient)
	}
	if flags.mapAllPorts {
		s = append(s, "--map-all")
	}
	if flags.customChainSpec != "" {
		s = append(s, "--custom-chainSpec", flags.customChainSpec)
	}
	if flags.customNetworkConfig != "" {
		s = append(s, "--custom-config", flags.customNetworkConfig)
	}
	if flags.customGenesis != "" {
		s = append(s, "--custom-genesis", flags.customGenesis)
	}
	if flags.customDeployBlock != "" {
		s = append(s, "--custom-deploy-block", flags.customDeployBlock)
	}
	if len(flags.customEnodes) > 0 {
		s = append(s, "--execution-bootnodes", strings.Join(flags.customEnodes, ","))
	}
	if len(flags.customEnrs) > 0 {
		s = append(s, "--consensus-bootnodes", strings.Join(flags.customEnrs, ","))
	}
	if len(flags.fallbackEL) > 0 {
		s = append(s, "--fallback-execution-urls", strings.Join(flags.fallbackEL, ","))
	}
	if flags.latestVersion {
		s = append(s, "--latest")
	}
	if flags.distributed {
		s = append(s, "--distributed")
	}
	if flags.distributedValidatorName != "" {
		s = append(s, "--distributedValidator", flags.distributedValidatorName)
	}
	return s
}

func (flags *GenCmdFlags) toString() string {
	return strings.Join(flags.argsList(), " ")
}

func TestGenerateCmd(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)
	tcs := []generateCmdTestCase{
		{
			"Execution, bad number of arguments",
			subCmd{
				name: "execution",
				args: []string{"nethermind", "besu"},
			},
			GenCmdFlags{},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("requires one argument"),
		},
		{
			"Consensus, bad number of arguments",
			subCmd{
				name: "consensus",
				args: []string{"teku", "lodestar"},
			},
			GenCmdFlags{
				executionAuthUrl: "http://localhost:8545",
				executionApiUrl:  "http://localhost:8545",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("requires one argument"),
		},
		{
			"Validator, bad number of arguments",
			subCmd{
				name: "validator",
				args: []string{"teku", "lodestar"},
			},
			GenCmdFlags{
				consensusApiUrl: "http://localhost:4000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("requires one argument"),
		},
		{
			"full-node Fixed clients",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				executionName: "nethermind",
				consensusName: "lighthouse",
				validatorName: "lighthouse",
				feeRecipient:  "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"full-node Random clients, no feeRecipient",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"full-node Random clients, relay-urls, single relay",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				relayURLs:    []string{"https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"full-node Random clients, relay-urls",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				relayURLs:    []string{"https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net", "https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"full-node Random clients, relay-urls",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				relayURLs:    []string{"https://@boost-relay.flashbots.net", "https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"full-node Random clients, invalid relay-urls",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				relayURLs:    []string{"https:/boost-relay.flashbots.net", "https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money{"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "relay", "https:/boost-relay.flashbots.net"),
		},
		{
			"full-node Random clients, invalid relay-url",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				relayURLs:    []string{"boost-relay.flashbots.net"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "relay", "boost-relay.flashbots.net"),
		},
		{
			"full-node Random clients, custom Ckpt sync endpoint",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				checkpointSyncUrl: "https://localhost:8545/api/v1/eth1",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, custom Ckpt sync endpoint",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl:  "https://localhost:8545",
				executionApiUrl:   "http://localhost",
				checkpointSyncUrl: "http://localhost/api/v1/eth1",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, custom Ckpt sync endpoint",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl:  "https://192.168.0.1:8545",
				executionApiUrl:   "http://127.0.0.1",
				checkpointSyncUrl: "http://localhost:7777/api/v1/eth1",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, custom Ckpt sync endpoint",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl:  "https://192.168.0.1:8545/v1/api",
				executionApiUrl:   "http://127.0.0.1/v1/api",
				checkpointSyncUrl: "https://chkp-sync:7777",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, bad custom Ckpt sync endpoint",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl:  "http://localhost:8545",
				executionApiUrl:   "http://localhost:8545",
				checkpointSyncUrl: "./8080",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "checkpoint sync", "./8080"),
		},
		{
			"consensus Random client, bad custom Ckpt sync endpoint",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl:  "http://localhost:8545/..,;",
				executionApiUrl:   "http://localhost:8545/{}",
				checkpointSyncUrl: "44.33.55.66:8080",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "checkpoint sync", "44.33.55.66:8080"),
		},
		{
			"consensus Random client, invalid api url",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://localhost",
				executionApiUrl:  "localhost/8545",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "execution api", "localhost/8545"),
		},
		{
			"consensus Random client, invalid api url",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://localhost",
				executionApiUrl:  "localhost:8545",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "execution api", "localhost:8545"),
		},
		{
			"consensus Random client, invalid auth url",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "htp://localhost:4000",
				executionApiUrl:  "https://localhost:8545",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "execution auth", "htp://localhost:4000"),
		},
		{
			"consensus Random client, valid mev-boost url",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				mevBoostUrl:      "http://mev-boost:3000",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, valid mev-boost url",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				mevBoostUrl:      "http://mev-boost/api/monkey/[spliat]",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, invalid mev-boost url",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				mevBoostUrl:      "mev-boost:3000",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "mev-boost endpoint", "mev-boost:3000"),
		},
		{
			"consensus Random client, invalid mev-boost url",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				mevBoostUrl:      "htp://mev-boost:3000",
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "mev-boost endpoint", "htp://mev-boost:3000"),
		},
		{
			"consensus Random client, custom enrs",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				customEnrs:       []string{"enr:-Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk", "enr:-Ly4QFoZTWR8ulxGVsWydTNGdwEESueIdj-wB6UmmjUcm-AOPxnQi7wprzwcdo7-1jBW_JxELlUKJdJES8TDsbl1EdNlh2F0dG5ldHOI__78_v2bsV-EZXRoMpA2-lATkAAAcf__________gmlkgnY0gmlwhBLYJjGJc2VjcDI1NmsxoQI0gujXac9rMAb48NtMqtSTyHIeNYlpjkbYpWJw46PmYYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA", "enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, custom enr",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				customEnrs:       []string{"enr:-Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			nil,
		},
		{
			"consensus Random client, invalid custom enrs",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				customEnrs:       []string{"enr:Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk", "enr:-Ly4QFoZTWR8ulxGVsWydTNGdwEESueIdj-wB6UmmjUcm-AOPxnQi7wprzwcdo7-1jBW_JxELlUKJdJES8TDsbl1EdNlh2F0dG5ldHOI__78_v2bsV-EZXRoMpA2-lATkAAAcf__________gmlkgnY0gmlwhBLYJjGJc2VjcDI1NmsxoQI0gujXac9rMAb48NtMqtSTyHIeNYlpjkbYpWJw46PmYYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA", "enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidEnrError, "enr:Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk"),
		},
		{
			"consensus Random client, invalid custom enrs",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				customEnrs:       []string{"enr-Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk", "enr:-Ly4QFoZTWR8ulxGVsWydTNGdwEESueIdj-wB6UmmjUcm-AOPxnQi7wprzwcdo7-1jBW_JxELlUKJdJES8TDsbl1EdNlh2F0dG5ldHOI__78_v2bsV-EZXRoMpA2-lATkAAAcf__________gmlkgnY0gmlwhBLYJjGJc2VjcDI1NmsxoQI0gujXac9rMAb48NtMqtSTyHIeNYlpjkbYpWJw46PmYYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA", "enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf(configs.InvalidEnrError, "enr-Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk"),
		},
		{
			"consensus Random client, duplicated custom enrs",
			subCmd{
				name: "consensus",
				args: []string{},
			},
			GenCmdFlags{
				executionAuthUrl: "http://execution:8551",
				executionApiUrl:  "https://execution:8545",
				customEnrs:       []string{"enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA", "enr:-Ly4QFoZTWR8ulxGVsWydTNGdwEESueIdj-wB6UmmjUcm-AOPxnQi7wprzwcdo7-1jBW_JxELlUKJdJES8TDsbl1EdNlh2F0dG5ldHOI__78_v2bsV-EZXRoMpA2-lATkAAAcf__________gmlkgnY0gmlwhBLYJjGJc2VjcDI1NmsxoQI0gujXac9rMAb48NtMqtSTyHIeNYlpjkbYpWJw46PmYYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA", "enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA"},
			},
			globalFlags{
				install: false,
				logging: "",
			},
			fmt.Errorf("%s: %s", configs.ErrDuplicatedBootNode, "enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA"),
		},
		{
			"Wrong sub command",
			subCmd{
				name: "full",
				args: []string{},
			},
			GenCmdFlags{},
			globalFlags{
				install: false,
				network: "",
				logging: "",
			},
			nil,
		},
		{
			"Missing validator client",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				executionName: "nethermind",
				consensusName: "lighthouse",
				feeRecipient:  "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "mainnet",
				logging:        "",
			},
			nil,
		},
		{
			"Good network input",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "mainnet",
				logging:        "",
			},
			nil,
		},
		{
			"Bad network input",
			subCmd{
				name: "consensus",
				args: []string{"lighthouse"},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "wrong",
				logging:        "",
			},
			errors.New("unknown network \"wrong\". Please provide correct network name. Use 'networks' command to see the list of supported networks"),
		},
		{
			"Consensus fixed client",
			subCmd{
				name: "consensus",
				args: []string{"teku"},
			},
			GenCmdFlags{
				executionAuthUrl: "http://localhost:8545/eth",
				executionApiUrl:  "https://execution/eth",
				feeRecipient:     "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			nil,
		},
		{
			"Consensus Missing execution-auth-url",
			subCmd{
				name: "consensus",
				args: []string{"teku"},
			},
			GenCmdFlags{
				executionApiUrl: "http://localhost:8545",
				feeRecipient:    "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("required flag(s) \"execution-auth-url\" not set"),
		},
		{
			"Consensus Missing execution-api-url",
			subCmd{
				name: "consensus",
				args: []string{"teku"},
			},
			GenCmdFlags{
				executionAuthUrl: "http://localhost:8545",
				feeRecipient:     "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("required flag(s) \"execution-api-url\" not set"),
		},
		{
			"Consensus wrong client",
			subCmd{
				name: "consensus",
				args: []string{"wrong"},
			},
			GenCmdFlags{
				executionAuthUrl: "http://localhost:8545",
				executionApiUrl:  "http://localhost:8545",
				feeRecipient:     "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("invalid consensus client"),
		},
		{
			"Validator missing consensus-url",
			subCmd{
				name: "validator",
				args: []string{},
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("required flag(s) \"consensus-url\" not set"),
		},
		{
			"Validator good client",
			subCmd{
				name: "validator",
				args: []string{"teku"},
			},
			GenCmdFlags{
				consensusApiUrl: "http://localhost:4000",
				feeRecipient:    "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			nil,
		},
		{
			"MevBoost",
			subCmd{
				name: "mev-boost",
				args: []string{},
			},
			GenCmdFlags{},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			nil,
		},
		{
			"MevBoost custom relay url, single one",
			subCmd{
				name: "mev-boost",
				args: []string{},
			},
			GenCmdFlags{
				relayURLs: []string{`"https://boost-relay.flashbots.net,"`},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			nil,
		},
		{
			"MevBoost custom relay urls",
			subCmd{
				name: "mev-boost",
				args: []string{},
			},
			GenCmdFlags{
				relayURLs: []string{"https://boost-relay.flashbots.net", "http://@boost-relay.flashbots.net"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			nil,
		},
		{
			"MevBoost custom relay urls",
			subCmd{
				name: "mev-boost",
				args: []string{},
			},
			GenCmdFlags{
				relayURLs: []string{"https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net", "https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			nil,
		},
		{
			"MevBoost invalid custom relay url",
			subCmd{
				name: "mev-boost",
				args: []string{},
			},
			GenCmdFlags{
				relayURLs: []string{"https:/boost-relay.flashbots.net", "https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "relay", "https:/boost-relay.flashbots.net"),
		},
		{
			"MevBoost invalid custom relay url",
			subCmd{
				name: "mev-boost",
				args: []string{},
			},
			GenCmdFlags{
				relayURLs: []string{"boost-relay.flashbots.net", "https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "relay", "boost-relay.flashbots.net"),
		},
		{
			"MevBoost wrong argument",
			subCmd{
				name: "mev-boost",
				args: []string{"wrong"},
			},
			GenCmdFlags{},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "",
				logging:        "",
			},
			errors.New("unknown command \"wrong\" for \"sedge generate mev-boost\""),
		},
		{
			"Execution ",
			subCmd{
				name: "execution",
				args: []string{"nethermind"},
			},
			GenCmdFlags{
				mapAllPorts: true,
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "sepolia",
				logging:        "",
			},
			nil,
		},
		{
			"Execution custom enodes",
			subCmd{
				name: "execution",
				args: []string{"nethermind"},
			},
			GenCmdFlags{
				customEnodes: []string{"enode://ea6d67eb3277d8ae9292fc700fa757ef6d2127c4db9712bcd5eb1341b1d937ac71cc2b15efe3a8496f4fc9fc12156d7ac73d82eb3c0f68928442116030b76f48@3.135.122.4:30303", "enode://c5e1e38709a2eb402557e82e071ccec1c6e2adedb01f7d6afdc80d25f7e9287f954fa9b742f01b1b74a5c532de9476afeb6efdcf5a683672a663204eadb15e45@3.17.46.220:30303"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "mainnet",
				logging:        "",
			},
			nil,
		},
		{
			"Execution custom enode",
			subCmd{
				name: "execution",
				args: []string{"nethermind"},
			},
			GenCmdFlags{
				customEnodes: []string{"enode://ea6d67eb3277d8ae9292fc700fa757ef6d2127c4db9712bcd5eb1341b1d937ac71cc2b15efe3a8496f4fc9fc12156d7ac73d82eb3c0f68928442116030b76f48@3.135.122.4:30303"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "mainnet",
				logging:        "",
			},
			nil,
		},
		{
			"Execution invalid custom enodes",
			subCmd{
				name: "execution",
				args: []string{"nethermind"},
			},
			GenCmdFlags{
				customEnodes: []string{"enode:3.135.122.4:30303", "enode://c5e1e38709a2eb402557e82e071ccec1c6e2adedb01f7d6afdc80d25f7e9287f954fa9b742f01b1b74a5c532de9476afeb6efdcf5a683672a663204eadb15e45@3.17.46.220:30303"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "mainnet",
				logging:        "",
			},
			fmt.Errorf(configs.InvalidEnodeError, "enode:3.135.122.4:30303"),
		},
		{
			"Execution invalid custom enodes",
			subCmd{
				name: "execution",
				args: []string{"nethermind"},
			},
			GenCmdFlags{
				customEnodes: []string{"enode://@3.135.122.4:30303", "enode://c5e1e38709a2eb402557e82e071ccec1c6e2adedb01f7d6afdc80d25f7e9287f954fa9b742f01b1b74a5c532de9476afeb6efdcf5a683672a663204eadb15e45@3.17.46.220:30303"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "mainnet",
				logging:        "",
			},
			fmt.Errorf(configs.InvalidEnodeError, "enode://@3.135.122.4:30303"),
		},
		{
			"Execution duplicated custom enodes",
			subCmd{
				name: "execution",
				args: []string{"nethermind"},
			},
			GenCmdFlags{
				customEnodes: []string{"enode://c5e1e38709a2eb402557e82e071ccec1c6e2adedb01f7d6afdc80d25f7e9287f954fa9b742f01b1b74a5c532de9476afeb6efdcf5a683672a663204eadb15e45@3.17.46.220:30303", "enode://c5e1e38709a2eb402557e82e071ccec1c6e2adedb01f7d6afdc80d25f7e9287f954fa9b742f01b1b74a5c532de9476afeb6efdcf5a683672a663204eadb15e45@3.17.46.220:30303"},
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "mainnet",
				logging:        "",
			},
			fmt.Errorf("%s: %s", configs.ErrDuplicatedBootNode, "enode://c5e1e38709a2eb402557e82e071ccec1c6e2adedb01f7d6afdc80d25f7e9287f954fa9b742f01b1b74a5c532de9476afeb6efdcf5a683672a663204eadb15e45@3.17.46.220:30303"),
		},
		{
			"Execution wrong client on gnosis",
			subCmd{
				name: "execution",
				args: []string{"geth"},
			},
			GenCmdFlags{},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "gnosis",
				logging:        "",
			},
			errors.New("invalid execution client"),
		},
		{
			"Execution ",
			subCmd{
				name: "execution",
				args: []string{"nethermind"},
			},
			GenCmdFlags{
				mapAllPorts: true,
			},
			globalFlags{
				install:        false,
				generationPath: "",
				network:        "sepolia",
				logging:        "",
			},
			nil,
		},
		{
			"Full-node, waitEpoch set",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				waitEpoch:    5,
			},
			globalFlags{
				network: "chiado",
			},
			nil,
		},
		{
			"Full-node, valid Graffiti",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				graffiti:     "sedge-graffiti",
			},
			globalFlags{
				network: "gnosis",
			},
			nil,
		},
		{
			"Full-node, Graffiti too long",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				feeRecipient: "0x0000000000000000000000000000000000000000",
				graffiti:     "sedge-graffiti-sedge",
			},
			globalFlags{
				network: "gnosis",
			},
			fmt.Errorf(configs.ErrGraffitiLength, "sedge-graffiti-sedge", 20),
		},
		{
			"Validator, waitEpoch set",
			subCmd{
				name: "validator",
				args: []string{"lodestar"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				waitEpoch:       50,
				consensusApiUrl: "http://localhost:4000",
			},
			globalFlags{
				network: "holesky",
			},
			nil,
		},
		{
			"Validator, invalid consensus api url",
			subCmd{
				name: "validator",
				args: []string{"lodestar"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "localhost/4000",
			},
			globalFlags{
				network: "gnosis",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "consensus api", "localhost/4000"),
		},
		{
			"Validator, invalid consensus api url",
			subCmd{
				name: "validator",
				args: []string{"teku"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "htp://localhost:4000",
			},
			globalFlags{
				network: "sepolia",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "consensus api", "htp://localhost:4000"),
		},
		{
			"Validator, invalid consensus api url",
			subCmd{
				name: "validator",
				args: []string{"lodestar"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "localhost:4000",
			},
			globalFlags{
				network: "gnosis",
			},
			fmt.Errorf(configs.InvalidUrlFlagError, "consensus api", "localhost:4000"),
		},
		{
			"Validator, valid consensus api url",
			subCmd{
				name: "validator",
				args: []string{"prysm"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "https://localhost:80/dasd,.,",
			},
			globalFlags{
				network: "sepolia",
			},
			nil,
		},
		{
			"Validator, valid consensus api url",
			subCmd{
				name: "validator",
				args: []string{"prysm"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "https://localhost",
			},
			globalFlags{
				network: "sepolia",
			},
			nil,
		},
		{
			"Validator, valid consensus api url",
			subCmd{
				name: "validator",
				args: []string{"lighthouse"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "https://localhost/api/endpoint",
			},
			globalFlags{
				network: "mainnet",
			},
			nil,
		},
		{
			"Validator, valid consensus api url",
			subCmd{
				name: "validator",
				args: []string{"lodestar"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "https://localhost:8000/api/endpoint",
			},
			globalFlags{
				network: "mainnet",
			},
			nil,
		},
		{
			"Validator, graffiti too long",
			subCmd{
				name: "validator",
				args: []string{"teku"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "https://localhost:8000/api/endpoint{}",
				graffiti:        "sedge-graffiti-sedge",
			},
			globalFlags{
				network: "holesky",
			},
			fmt.Errorf(configs.ErrGraffitiLength, "sedge-graffiti-sedge", 20),
		},
		{
			"Validator, valid graffiti",
			subCmd{
				name: "validator",
				args: []string{"lodestar"},
			},
			GenCmdFlags{
				feeRecipient:    "0x0000000000000000000000000000000000000000",
				consensusApiUrl: "https://localhost:8000/api/endpoint",
				graffiti:        "sedge-graffiti",
			},
			globalFlags{
				network: "mainnet",
			},
			nil,
		},
		{
			"Validator blocker not generated with --no-validator flag",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				noValidator:   true,
				executionName: "nethermind",
				consensusName: "teku",
			},
			globalFlags{
				network: "mainnet",
			},
			nil,
		},
		{
			"Full node - Latest version of clients",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				noValidator:   true,
				executionName: "nethermind",
				consensusName: "teku",
				latestVersion: true,
			},
			globalFlags{
				network: "mainnet",
			},
			nil,
		},
		{
			"Execution - Latest version of clients",
			subCmd{
				name: "execution",
			},
			GenCmdFlags{
				latestVersion: true,
			},
			globalFlags{},
			nil,
		},
		{
			"Consensus - Latest version of clients",
			subCmd{
				name: "consensus",
			},
			GenCmdFlags{
				latestVersion:    true,
				executionApiUrl:  "https://localhost:8545",
				executionAuthUrl: "https://localhost:8545",
			},
			globalFlags{},
			nil,
		},
		{
			"Validator - Latest version of clients",
			subCmd{
				name: "validator",
			},
			GenCmdFlags{
				latestVersion:   true,
				consensusApiUrl: "https://localhost:8000/api/endpoint",
			},
			globalFlags{},
			nil,
		},
		{
			"full-node random client Distributed",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				distributed:   true,
				validatorName: "lighthouse",
			},
			globalFlags{
				network: "holesky",
			},
			nil,
		},
		{
			"full-node Fixed clients with DV",
			subCmd{
				name: "full-node",
				args: []string{},
			},
			GenCmdFlags{
				distributed:              true,
				executionName:            "nethermind",
				consensusName:            "lighthouse",
				validatorName:            "lighthouse",
				distributedValidatorName: "charon",
				feeRecipient:             "0x0000000000000000000000000000000000000000",
			},
			globalFlags{
				install: false,
				logging: "",
				network: "holesky",
			},
			nil,
		},
		{
			"Optimism full node",
			subCmd{
				name: "op-full-node",
			},
			GenCmdFlags{},
			globalFlags{},
			nil,
		},
		{
			"Optimism full node with api url",
			subCmd{
				name: "op-full-node",
			},
			GenCmdFlags{
				executionApiUrl: "https://localhost:8545",
				consensusApiUrl: "https://localhost:8000/api/endpoint",
			},
			globalFlags{},
			nil,
		},
		{
			"Optimism full node with only consensus api url",
			subCmd{
				name: "op-full-node",
			},
			GenCmdFlags{
				consensusApiUrl: "https://localhost:8000/api/endpoint",
			},
			globalFlags{},
			nil,
		},
		{
			"Lido Full-node - Holesky without MEV",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				noMev: true,
			},
			globalFlags{
				network:  NetworkHolesky,
				lidoNode: true,
			},
			nil,
		},
		{
			"Lido Full-node - Holesky without MEV",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				noMev: true,
			},
			globalFlags{
				network:  NetworkHolesky,
				lidoNode: true,
			},
			nil,
		},
		{
			"Lido Full-node - Sepolia",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{},
			globalFlags{
				network:  NetworkSepolia,
				lidoNode: true,
			},
			fmt.Errorf(configs.InvalidNetworkForLidoMevBoost, mevboostrelaylist.LidoSupportedNetworksMevBoost()),
		},
		{
			"Lido Full-node - Holesky no validator",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				noValidator: true,
			},
			globalFlags{
				network:  NetworkHolesky,
				lidoNode: true,
			},
			nil,
		},
		{
			"Lido Full-node - unsupported Gnosis",
			subCmd{
				name: "full-node",
			},
			GenCmdFlags{
				noMev: true,
			},
			globalFlags{
				network:  NetworkGnosis,
				lidoNode: true,
			},
			fmt.Errorf(configs.InvalidNetworkForLido, contracts.LidoSupportedNetworks()),
		},
	}

	// TODO: Add test cases for Execution fallback urls
	// TODO: Add test cases for EL and CL bootnodes in full-node

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			descr := fmt.Sprintf("sedge generate %s %s %s", strings.Join(tc.subCommand.argsList(), " "), tc.args.toString(), strings.Join(tc.globalArgs.argsList(), " "))
			sedgeActions := actions.NewSedgeActions(actions.SedgeActionsOptions{})

			rootCmd := RootCmd()
			rootCmd.AddCommand(GenerateCmd(sedgeActions))
			argsL := append([]string{"generate"}, tc.subCommand.argsList()...)
			argsL = append(argsL, tc.args.argsList()...)
			argsL = append(argsL, tc.globalArgs.argsList()...)
			argsL = append(argsL, "-p", t.TempDir())
			rootCmd.SetArgs(argsL)
			rootCmd.SetOutput(io.Discard)

			err := rootCmd.Execute()

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error(), descr)
			} else {
				assert.NoError(t, err, descr)
			}
		})
	}
}

func TestGeneratePathCases(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	// Custom Generation path
	path := t.TempDir()
	descr := fmt.Sprintf("Generation path error, sedge generate execution --path %s", path)
	sedgeActions := actions.NewSedgeActions(actions.SedgeActionsOptions{})

	rootCmd := RootCmd()
	rootCmd.AddCommand(GenerateCmd(sedgeActions))
	argsL := []string{"generate", "execution", "--path", path}
	rootCmd.SetArgs(argsL)
	rootCmd.SetOutput(io.Discard)

	err := rootCmd.Execute()

	assert.NoError(t, err, descr)

	// Init generation path
	path = t.TempDir()
	if err = os.Remove(path); err != nil {
		t.Fatal(err)
	}
	descr = fmt.Sprintf("Init generation path, sedge generate execution --path %s", path)
	sedgeActions = actions.NewSedgeActions(actions.SedgeActionsOptions{})

	rootCmd = RootCmd()
	rootCmd.AddCommand(GenerateCmd(sedgeActions))
	argsL = []string{"generate", "execution", "--path", path}
	rootCmd.SetArgs(argsL)
	rootCmd.SetOutput(io.Discard)

	err = rootCmd.Execute()

	assert.NoError(t, err, descr)

	// Custom jwt secret path, good
	path = t.TempDir()
	descr = fmt.Sprintf("Custom jwt secret path, good, sedge generate execution --jwt-secret-path %s", path)
	err = test.PrepareTestCaseDir(filepath.Join("testdata", "cli_tests", "jwtsecret"), path)
	if err != nil {
		t.Fatalf("Can't build test case: %v", err)
	}
	jwtPath := filepath.Join(path, "jwtsecret")

	sedgeActions = actions.NewSedgeActions(actions.SedgeActionsOptions{})

	rootCmd = RootCmd()
	rootCmd.AddCommand(GenerateCmd(sedgeActions))
	argsL = []string{"generate", "execution", "--path", path, "--jwt-secret-path", jwtPath}
	rootCmd.SetArgs(argsL)
	rootCmd.SetOutput(io.Discard)

	err = rootCmd.Execute()

	assert.NoError(t, err, descr)

	// Custom jwt secret path, error
	path = t.TempDir()
	descr = fmt.Sprintf("Custom jwt secret path, error, sedge generate execution --jwt-secret-path %s", path)
	sedgeActions = actions.NewSedgeActions(actions.SedgeActionsOptions{})

	rootCmd = RootCmd()
	rootCmd.AddCommand(GenerateCmd(sedgeActions))
	argsL = []string{"generate", "execution", "--jwt-secret-path", path}
	rootCmd.SetArgs(argsL)
	rootCmd.SetOutput(io.Discard)

	err = rootCmd.Execute()

	assert.Error(t, err, descr)
}
