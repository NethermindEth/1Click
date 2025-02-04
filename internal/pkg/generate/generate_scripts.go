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
package generate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/NethermindEth/sedge/internal/pkg/env"

	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/clients"
	"github.com/NethermindEth/sedge/internal/utils"
	"github.com/NethermindEth/sedge/templates"
	log "github.com/sirupsen/logrus"
)

const (
	execution            = "execution"
	consensus            = "consensus"
	validator            = "validator"
	optimism             = "optimism"
	opExecution          = "opexecution"
	validatorImport      = "validator-import"
	mevBoost             = "mev-boost"
	configConsensus      = "config_consensus"
	empty                = "empty"
	distributedValidator = "distributedValidator"
	charon               = "charon"
)

// validateClients validates each client in GenData
func validateClients(gd *GenData) error {
	c := clients.ClientInfo{Network: gd.Network}
	if err := validateConsensus(gd, &c); err != nil {
		return err
	}
	if err := validateExecution(gd, &c); err != nil {
		return err
	}
	if err := validateValidator(gd, &c); err != nil {
		return err
	}
	if err := validateDistributedValidator(gd, &c); err != nil {
		return err
	}
	return nil
}

// validateValidator validates the validator client in GenData
func validateValidator(gd *GenData, c *clients.ClientInfo) error {
	if gd.ValidatorClient == nil {
		return nil
	}
	validatorClients, err := c.SupportedClients(validator)
	if err != nil {
		return ErrUnableToGetClientsInfo
	}
	if !utils.Contains(validatorClients, gd.ValidatorClient.Name) {
		// Only check for Grandine specifically if it's not in the supported list
		// This provides a more helpful error message
		if gd.ValidatorClient.Name == "grandine" {
			return fmt.Errorf("grandine is a consensus-only client and currently not supported as a validator")
		}
		return ErrValidatorClientNotValid
	}
	return nil
}

// validateDistributedValidator validates the validator client in GenData
func validateDistributedValidator(gd *GenData, c *clients.ClientInfo) error {
	if gd.DistributedValidatorClient == nil {
		return nil
	}
	distributedValidatorClients, err := c.SupportedClients(distributedValidator)
	if err != nil {
		return ErrUnableToGetClientsInfo
	}
	if !utils.Contains(distributedValidatorClients, gd.DistributedValidatorClient.Name) {
		return ErrDistributedValidatorClientNotValid
	}
	return nil
}

// validateExecution validates the execution client in GenData
func validateExecution(gd *GenData, c *clients.ClientInfo) error {
	if gd.ExecutionClient == nil {
		return nil
	}
	executionClients, err := c.SupportedClients(execution)
	if err != nil {
		return ErrUnableToGetClientsInfo
	}
	if !utils.Contains(executionClients, gd.ExecutionClient.Name) {
		return ErrExecutionClientNotValid
	}
	return nil
}

// validateConsensus validates the consensus client in GenData
func validateConsensus(gd *GenData, c *clients.ClientInfo) error {
	if gd.ConsensusClient == nil {
		return nil
	}

	consensusClients, err := c.SupportedClients(consensus)
	if err != nil {
		return ErrUnableToGetClientsInfo
	}
	if !utils.Contains(consensusClients, gd.ConsensusClient.Name) {
		return ErrConsensusClientNotValid
	}
	return nil
}

// mapClients convert genData clients to clients.Clients
func mapClients(gd *GenData) map[string]*clients.Client {
	cls := map[string]*clients.Client{
		execution:            gd.ExecutionClient,
		consensus:            gd.ConsensusClient,
		validator:            gd.ValidatorClient,
		optimism:             gd.OptimismClient,
		opExecution:          gd.ExecutionOPClient,
		distributedValidator: gd.DistributedValidatorClient,
	}

	return cls
}

// ComposeFile generates a docker-compose file with the provided GenData
func ComposeFile(gd *GenData, at io.Writer) error {
	// Check empty data
	if gd == nil {
		return ErrEmptyData
	}
	err := validateClients(gd)
	if err != nil {
		return err
	}
	// Check for port occupation
	defaultsPorts := map[string]uint16{
		"ELDiscovery":       configs.DefaultDiscoveryPortEL,
		"ELMetrics":         configs.DefaultMetricsPortEL,
		"ELApi":             configs.DefaultApiPortEL,
		"ELAuth":            configs.DefaultAuthPortEL,
		"ELWS":              configs.DefaultWSPortEL,
		"CLDiscovery":       configs.DefaultDiscoveryPortCL,
		"CLMetrics":         configs.DefaultMetricsPortCL,
		"CLApi":             configs.DefaultApiPortCL,
		"CLAdditionalApi":   configs.DefaultAdditionalApiPortCL,
		"VLMetrics":         configs.DefaultMetricsPortVL,
		"MevPort":           configs.DefaultMevPort,
		"DVDiscovery":       configs.DefaultDiscoveryPortDV,
		"DVMetrics":         configs.DefaultMetricsPortDV,
		"DVApi":             configs.DefaultApiPortDV,
		"ApiPortELOP":       configs.DefaultApiPortELOP,
		"AuthPortELOP":      configs.DefaultAuthPortELOP,
		"DiscoveryPortELOP": configs.DefaultDiscoveryPortELOP,
		"MetricsPortELOP":   configs.DefaultMetricsPortELOP,
	}
	ports, err := utils.AssignPorts("localhost", defaultsPorts)
	if err != nil {
		// notest
		return fmt.Errorf(configs.PortOccupationError, err)
	}
	gd.Ports = ports

	rawBaseTmp, err := templates.Services.ReadFile(strings.Join([]string{"services", "docker-compose_base.tmpl"}, "/"))
	if err != nil {
		return err
	}

	baseTmp, err := template.New("docker-compose").Parse(string(rawBaseTmp))
	if err != nil {
		return err
	}

	cls := mapClients(gd)
	networkConfig := configs.NetworksConfigs()[gd.Network]
	for tmpKind, client := range cls {
		var name string
		if client == nil {
			name = empty
		} else {
			name = client.Name
		}
		tmp, err := templates.Services.ReadFile(strings.Join([]string{
			"services",
			networkConfig.NetworkService,
			tmpKind,
			name + ".tmpl",
		}, "/"))
		if err != nil {
			return err
		}
		_, err = baseTmp.Parse(string(tmp))
		if err != nil {
			return err
		}
	}
	validatorBlockerTemplate := "validator-blocker"

	// Parse validator-blocker template
	tmp, err := templates.Services.ReadFile(strings.Join([]string{"services", validatorBlockerTemplate + ".tmpl"}, "/"))
	if err != nil {
		return err
	}
	if _, err = baseTmp.Parse(string(tmp)); err != nil {
		return err
	}

	// Check for splitted network flags
	splittedNetwork, err := env.CheckVariableBase(env.ReSPLITTED, gd.Network)
	if err != nil {
		return err
	}

	// Check vars related to Consensus service
	var xeeVersion, clCheckpointSyncUrl bool

	if cls[execution] != nil {
		gd.ExecutionClient.Endpoint = configs.OnPremiseExecutionURL
	}

	if cls[consensus] != nil {
		gd.ConsensusClient.Endpoint = configs.OnPremiseConsensusURL
		// Check for XEE_VERSION in teku
		xeeVersion, err = env.CheckVariable(env.ReXEEV, gd.Network, "consensus", gd.ConsensusClient.Name)
		if err != nil {
			return err
		}

		clCheckpointSyncUrl, err = env.CheckVariable(env.ReCHECKPOINT, gd.Network, "consensus", gd.ConsensusClient.Name)
		if err != nil {
			return err
		}
	}

	// Check for CL Bootnode nodes
	if len(gd.CCBootnodes) == 0 {
		gd.CCBootnodes = configs.NetworksConfigs()[gd.Network].DefaultCCBootnodes
	}

	// Check for EL Bootnode nodes
	if len(gd.ECBootnodes) == 0 {
		gd.ECBootnodes = configs.NetworksConfigs()[gd.Network].DefaultECBootnodes
	}
	var mevSupported bool
	if cls[validator] != nil {
		mevSupported, err = env.CheckVariable(env.ReMEV, gd.Network, "validator", gd.ValidatorClient.Name)
		if err != nil {
			return err
		}
	}
	// If consensus is running with the validator, and the MevBoostEndpoint is not set, set it to the default value
	if cls[consensus] != nil && cls[validator] != nil && gd.MevBoostEndpoint == "" && gd.Mev {
		mevSupported, err = env.CheckVariable(env.ReMEV, gd.Network, "validator", gd.ConsensusClient.Name)
		if err != nil {
			return err
		}
		if mevSupported {
			gd.MevBoostEndpoint = fmt.Sprintf("%s:%v", configs.DefaultMevBoostEndpoint, gd.Ports["MevPort"])
		}
	}
	gd.MevBoostService = slices.Contains(gd.Services, "mev-boost")

	if gd.Distributed {
		// Check for distributed validator
		if cls[distributedValidator] != nil {
			gd.DistributedValidatorClient.Endpoint = configs.OnPremiseDistributedValidatorURL
		}
	}
	consensusApiUrl := gd.ConsensusApiUrl
	if cls[consensus] != nil && consensusApiUrl == "" {
		consensusApiUrl = fmt.Sprintf("%s:%v", endpointOrEmpty(cls[consensus]), gd.Ports["CLApi"])
	}

	// Set network prefix
	networkPrefix := "op"
	if gd.IsBase {
		networkPrefix = "base"
	}

	data := DockerComposeData{
		Services:            gd.Services,
		Network:             gd.Network,
		Distributed:         gd.Distributed,
		XeeVersion:          xeeVersion,
		Mev:                 networkConfig.SupportsMEVBoost && (gd.MevBoostService || (mevSupported && gd.Mev)),
		MevBoostOnValidator: gd.MevBoostService || (mevSupported && gd.Mev) || gd.MevBoostOnValidator,
		MevPort:             gd.Ports["MevPort"],
		MevBoostEndpoint:    gd.MevBoostEndpoint,
		MevImage:            gd.MevImage,
		CheckpointSyncUrl:   gd.CheckpointSyncUrl,
		FeeRecipient:        gd.FeeRecipient,
		ElDiscoveryPort:     gd.Ports["ELDiscovery"],
		ElMetricsPort:       gd.Ports["ELMetrics"],
		ElApiPort:           gd.Ports["ELApi"],
		ElAuthPort:          gd.Ports["ELAuth"],
		ElWsPort:            gd.Ports["ELWS"],
		ClDiscoveryPort:     gd.Ports["CLDiscovery"],
		ClMetricsPort:       gd.Ports["CLMetrics"],
		ClApiPort:           gd.Ports["CLApi"],
		ClAdditionalApiPort: gd.Ports["CLAdditionalApi"],
		VlMetricsPort:       gd.Ports["VLMetrics"],
		ElOPApiPort:         gd.Ports["ApiPortELOP"],
		ElOPAuthPort:        gd.Ports["AuthPortELOP"],
		ElOPDiscoveryPort:   gd.Ports["DiscoveryPortELOP"],
		ElOPMetricsPort:     gd.Ports["MetricsPortELOP"],
		FallbackELUrls:      gd.FallbackELUrls,
		ElExtraFlags:        gd.ElExtraFlags,
		ElOPExtraFlags:      gd.ElOpExtraFlags,
		OPExtraFlags:        gd.OpExtraFlags,
		NetworkPrefix:       networkPrefix,
		ClExtraFlags:        gd.ClExtraFlags,
		VlExtraFlags:        gd.VlExtraFlags,
		DvExtraFlags:        gd.DvExtraFlags,
		ECBootnodes:         strings.Join(gd.ECBootnodes, ","),
		CCBootnodes:         strings.Join(gd.CCBootnodes, ","),
		CCBootnodesList:     gd.CCBootnodes,
		MapAllPorts:         gd.MapAllPorts,
		SplittedNetwork:     splittedNetwork,
		ClCheckpointSyncUrl: clCheckpointSyncUrl,
		LoggingDriver:       gd.LoggingDriver,
		VLStartGracePeriod:  gd.VLStartGracePeriod,
		CustomNetwork:       gd.Network == configs.NetworkCustom, // Used custom templates
		CustomConsensusConfigs: gd.CustomNetworkConfigPath != "" ||
			gd.CustomGenesisPath != "" ||
			gd.CustomDeployBlockPath != "", // Have custom configs paths
		CustomChainSpecPath:     gd.CustomChainSpecPath,     // Path to chainspec.json
		CustomNetworkConfigPath: gd.CustomNetworkConfigPath, // Path to config.yaml
		CustomGenesisPath:       gd.CustomGenesisPath,       // Path to genesis.ssz
		CustomDeployBlockPath:   gd.CustomDeployBlockPath,   // Path to deploy_block.txt
		UID:                     os.Geteuid(),
		GID:                     os.Getegid(),
		ContainerTag:            gd.ContainerTag,
		ConsensusApiURL:         consensusApiUrl,
		DVDiscoveryPort:         gd.Ports["DVDiscovery"],
		DVMetricsPort:           gd.Ports["DVMetrics"],
		DVApiPort:               gd.Ports["DVApi"],
	}

	// Save to writer
	err = baseTmp.Execute(at, data)
	if err != nil {
		return err
	}

	return nil
}

// EnvFile generates a .env file with the provided GenData
func EnvFile(gd *GenData, at io.Writer) error {
	rawBaseTmp, err := templates.Envs.ReadFile(strings.Join([]string{"envs", gd.Network, "env_base.tmpl"}, "/"))
	if err != nil {
		return ErrTemplateNotFound
	}

	baseTmp, err := template.New("env").Parse(string(rawBaseTmp))
	if err != nil {
		return err
	}

	cls := mapClients(gd)
	networkConfig := configs.NetworksConfigs()[gd.Network]
	for tmpKind, client := range cls {
		var tmp []byte
		if client == nil {
			tmp, err = templates.Services.ReadFile(strings.Join([]string{
				"services",
				networkConfig.NetworkService,
				tmpKind,
				"empty.tmpl",
			}, "/"))
			if err != nil {
				return err
			}
		} else {
			tmp, err = templates.Envs.ReadFile(strings.Join([]string{"envs", gd.Network, tmpKind, client.Name + ".tmpl"}, "/"))
			if err != nil {
				return err
			}
		}
		_, err = baseTmp.Parse(string(tmp))
		if err != nil {
			return err
		}
	}
	executionApiUrl := gd.ExecutionApiUrl
	executionAuthUrl := gd.ExecutionAuthUrl
	if cls[execution] != nil {
		if executionApiUrl == "" {
			executionApiUrl = fmt.Sprintf("%s:%v", cls[execution].Endpoint, gd.Ports["ELApi"])
		}
		if executionAuthUrl == "" {
			executionAuthUrl = fmt.Sprintf("%s:%v", cls[execution].Endpoint, gd.Ports["ELAuth"])
		}
	}
	consensusApiUrl := gd.ConsensusApiUrl
	consensusAdditionalApiUrl := consensusApiUrl
	if cls[consensus] != nil && consensusApiUrl == "" {
		consensusAdditionalApiUrl = fmt.Sprintf("%s:%v", endpointOrEmpty(cls[consensus]), gd.Ports["CLAdditionalApi"])
		consensusApiUrl = fmt.Sprintf("%s:%v", endpointOrEmpty(cls[consensus]), gd.Ports["CLApi"])

		// Prysm urls must be without http:// or https://
		if cls[validator] != nil && cls[validator].Name == "prysm" {
			consensusAdditionalApiUrl = fmt.Sprintf("%s:%v", "consensus", gd.Ports["CLAdditionalApi"])
		}
	} else {
		if cls[consensus] != nil && cls[consensus].Name == "prysm" {
			consensusAdditionalApiUrl = fmt.Sprintf("%s:%v", "consensus", gd.Ports["CLAdditionalApi"])
		} else if cls[validator] != nil && cls[validator].Name == "prysm" {
			// Strip the http:// or https:// from the url
			consensusAdditionalApiUrl = strings.TrimPrefix(consensusAdditionalApiUrl, "http://")
			consensusAdditionalApiUrl = strings.TrimPrefix(consensusAdditionalApiUrl, "https://")
		} else {
			consensusAdditionalApiUrl = consensusApiUrl
		}
	}

	if cls[optimism] != nil {
		gd.ExecutionOPClient.Endpoint = configs.OnPremiseOpExecutionURL
	}

	executionOPApiUrl := ""
	if cls[optimism] != nil {
		executionOPApiUrl = fmt.Sprintf("%s:%v", endpointOrEmpty(gd.ExecutionOPClient), gd.Ports["ApiPortELOP"])
	}
	var mevSupported bool
	if cls[validator] != nil {
		mevSupported, err = env.CheckVariable(env.ReMEV, gd.Network, "validator", gd.ValidatorClient.Name)
		if err != nil {
			return err
		}
	}
	// If consensus is running with the validator, and the MevBoostEndpoint is not set, set it to the default value
	if cls[consensus] != nil && cls[validator] != nil && gd.MevBoostEndpoint == "" && gd.Mev {
		mevSupported, err = env.CheckVariable(env.ReMEV, gd.Network, "validator", gd.ConsensusClient.Name)
		if err != nil {
			return err
		}
		if mevSupported {
			gd.MevBoostEndpoint = fmt.Sprintf("%s:%v", configs.DefaultMevBoostEndpoint, gd.Ports["MevPort"])
		}
	}

	graffiti := gd.Graffiti
	if graffiti == "" {
		graffiti = generateGraffiti(gd.ExecutionClient, gd.ConsensusClient, gd.ValidatorClient)
	}

	if len(gd.RelayURLs) == 0 {
		gd.RelayURLs = configs.NetworksConfigs()[gd.Network].RelayURLs
	}

	if gd.CheckpointSyncUrl == "" {
		gd.CheckpointSyncUrl = configs.NetworksConfigs()[gd.Network].CheckpointSyncURL
	}
	elOpImage := ""
	if gd.ExecutionOPClient != nil {
		elOpImage = imageOrEmpty(gd.ExecutionOPClient, gd.LatestVersion)
	}
	opImageVersion := ""
	opSequencerHttp := ""
	rethNetwork := ""
	if gd.OptimismClient != nil {
		opImageVersion = imageOrEmpty(cls[optimism], gd.LatestVersion)
		if gd.IsBase {
			opSequencerHttp = "https://" + gd.Network + "-sequencer.base.org"
		} else {
			opSequencerHttp = "https://" + gd.Network + "-sequencer.optimism.io"
		}
		if gd.Network == configs.NetworkMainnet {
			if gd.IsBase {
				rethNetwork = "base"
			} else {
				rethNetwork = "optimism"
			}
		}
		if gd.Network == configs.NetworkSepolia {
			if gd.IsBase {
				rethNetwork = "base-sepolia"
			} else {
				rethNetwork = "optimism-sepolia"
			}
		}

	}

	distributedValidatorApiUrl := ""
	if gd.Distributed {
		// Check for distributed validator
		if cls[distributedValidator] != nil {
			distributedValidatorApiUrl = fmt.Sprintf("%s:%v", cls[distributedValidator].Endpoint, gd.Ports["DVApi"])
		}
	}

	data := EnvData{
		Services:                   gd.Services,
		Mev:                        networkConfig.SupportsMEVBoost && (gd.MevBoostService || (mevSupported && gd.Mev) || gd.MevBoostOnValidator),
		ElImage:                    imageOrEmpty(cls[execution], gd.LatestVersion),
		ElDataDir:                  "./" + configs.ExecutionDir,
		CcImage:                    imageOrEmpty(cls[consensus], gd.LatestVersion),
		CcDataDir:                  "./" + configs.ConsensusDir,
		VlImage:                    imageOrEmpty(cls[validator], gd.LatestVersion),
		VlDataDir:                  "./" + configs.ValidatorDir,
		ExecutionApiURL:            executionApiUrl,
		ExecutionAuthURL:           executionAuthUrl,
		ConsensusApiURL:            consensusApiUrl,
		ConsensusAdditionalApiURL:  consensusAdditionalApiUrl,
		FeeRecipient:               gd.FeeRecipient,
		JWTSecretPath:              gd.JWTSecretPath,
		ExecutionEngineName:        nameOrEmpty(cls[execution]),
		ConsensusClientName:        nameOrEmpty(cls[consensus]),
		KeystoreDir:                "./" + configs.KeystoreDir,
		Graffiti:                   graffiti,
		RelayURLs:                  strings.Join(gd.RelayURLs, ","),
		CheckpointSyncUrl:          gd.CheckpointSyncUrl,
		Distributed:                gd.Distributed,
		DistributedValidatorApiUrl: distributedValidatorApiUrl,
		DvDataDir:                  "./" + configs.DistributedValidatorDir,
		DvImage:                    imageOrEmpty(cls[distributedValidator], gd.LatestVersion),
		ExecutionOPApiURL:          executionOPApiUrl,
		JWTOPSecretPath:            gd.JWTSecretOP,
		OPImageVersion:             opImageVersion,
		ElOpImage:                  elOpImage,
		ElOPAuthPort:               gd.Ports["AuthPortELOP"],
		OpSequencerHttp:            opSequencerHttp,
		RethNetwork:                rethNetwork,
	}

	// Save to writer
	err = baseTmp.Execute(at, data)
	if err != nil {
		return err
	}

	return nil
}

type CustomConfigsSources struct {
	ChainSpecSrc     string
	NetworkConfigSrc string
	GenesisSrc       string
	DeployBlockSrc   string
}

type CustomNetworkConfigsData struct {
	ChainSpecPath     string
	NetworkConfigPath string
	GenesisPath       string
	DeployBlockPath   string
}

func CustomNetworkConfigs(generationPath, network string, sources CustomConfigsSources) (CustomNetworkConfigsData, error) {
	var customNetworkConfigsData CustomNetworkConfigsData
	networkData, ok := configs.NetworksConfigs()[network]
	if !ok {
		return customNetworkConfigsData, fmt.Errorf(configs.UnknownNetworkError, network)
	}
	valueOrDefault := func(value, def string) string {
		if value != "" {
			return value
		}
		return def
	}
	chainSpecSrc := valueOrDefault(sources.ChainSpecSrc, networkData.DefaultCustomChainSpecSrc)
	networkConfigSrc := valueOrDefault(sources.NetworkConfigSrc, networkData.DefaultCustomConfigSrc)
	genesisSrc := valueOrDefault(sources.GenesisSrc, networkData.DefaultCustomGenesisSrc)
	deployBlock := valueOrDefault(sources.DeployBlockSrc, networkData.DefaultCustomDeployBlock)

	// Check if any custom config is needed
	if chainSpecSrc == "" && networkConfigSrc == "" && genesisSrc == "" && deployBlock == "" {
		return customNetworkConfigsData, nil
	}

	// Setup destination folder
	destFolder := filepath.Join(generationPath, configs.CustomNetworkConfigsFolder)
	if _, err := os.Stat(destFolder); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(destFolder, os.ModePerm)
			if err != nil {
				return customNetworkConfigsData, err
			}
		} else {
			return customNetworkConfigsData, err
		}
	}

	if chainSpecSrc != "" {
		customNetworkConfigsData.ChainSpecPath = filepath.Join(destFolder, configs.ExecutionNetworkConfigFileName)
		log.Info(configs.GettingCustomChainSpec)
		err := utils.DownloadOrCopy(chainSpecSrc, customNetworkConfigsData.ChainSpecPath, true)
		if err != nil {
			return customNetworkConfigsData, err
		}
		customNetworkConfigsData.ChainSpecPath, err = filepath.Abs(customNetworkConfigsData.ChainSpecPath)
		if err != nil {
			return customNetworkConfigsData, err
		}
	}

	if networkConfigSrc != "" {
		customNetworkConfigsData.NetworkConfigPath = filepath.Join(destFolder, configs.ConsensusNetworkConfigFileName)
		log.Info(configs.GettingCustomNetworkConfig)
		err := utils.DownloadOrCopy(networkConfigSrc, customNetworkConfigsData.NetworkConfigPath, true)
		if err != nil {
			return customNetworkConfigsData, err
		}
		customNetworkConfigsData.NetworkConfigPath, err = filepath.Abs(customNetworkConfigsData.NetworkConfigPath)
		if err != nil {
			return customNetworkConfigsData, err
		}
	}

	if genesisSrc != "" {
		customNetworkConfigsData.GenesisPath = filepath.Join(destFolder, configs.ConsensusNetworkGenesisFileName)
		log.Info(configs.GettingCustomGenesis)
		err := utils.DownloadOrCopy(genesisSrc, customNetworkConfigsData.GenesisPath, true)
		if err != nil {
			return customNetworkConfigsData, err
		}
		customNetworkConfigsData.GenesisPath, err = filepath.Abs(customNetworkConfigsData.GenesisPath)
		if err != nil {
			return customNetworkConfigsData, err
		}
	}

	if deployBlock != "" {
		customNetworkConfigsData.DeployBlockPath = filepath.Join(destFolder, configs.ConsensusNetworkDeployBlockFileName)
		log.Info(configs.WritingCustomDeployBlock)
		err := os.WriteFile(customNetworkConfigsData.DeployBlockPath, []byte(deployBlock), os.ModePerm)
		if err != nil {
			return customNetworkConfigsData, fmt.Errorf(configs.ErrorWritingDeployBlockFile, customNetworkConfigsData.DeployBlockPath, err)
		}
		customNetworkConfigsData.DeployBlockPath, err = filepath.Abs(customNetworkConfigsData.DeployBlockPath)
		if err != nil {
			return customNetworkConfigsData, err
		}
	}

	return customNetworkConfigsData, nil
}

// endpointOrEmpty returns the endpoint of the client if it is not nil, otherwise returns an empty string
func endpointOrEmpty(cls *clients.Client) string {
	if cls != nil {
		return cls.Endpoint
	}
	return ""
}

// nameOrEmpty returns the name of the client if it is not nil, otherwise returns an empty string
func nameOrEmpty(cls *clients.Client) string {
	if cls != nil {
		return cls.Name
	}
	return ""
}

// generateGraffiti generates a graffiti string based on the execution, consensus and validator clients
func generateGraffiti(execution, consensus, validator *clients.Client) string {
	if consensus != nil && execution != nil {
		if validator != nil {
			if consensus.Name == validator.Name {
				return strings.Join([]string{nameOrEmpty(execution), nameOrEmpty(consensus)}, "-")
			}
		}
	}
	return joinIfNotEmpty(nameOrEmpty(execution), nameOrEmpty(consensus), nameOrEmpty(validator))
}

// joinIfNotEmpty joins the strings if they are not empty
func joinIfNotEmpty(strs ...string) string {
	var result []string
	for _, str := range strs {
		if str != "" {
			result = append(result, str)
		}
	}
	return strings.Join(result, "-")
}

// imageOrEmpty returns the image of the client if it is not nil, otherwise returns an empty string
func imageOrEmpty(cls *clients.Client, latest bool) string {
	if cls != nil {
		if latest && !cls.Modified {
			if cls.Name == "nimbus" {
				splits := strings.Split(cls.Image, ":")
				splits[len(splits)-1] = "multiarch-latest"
				return strings.Join(splits, ":")

			} else {
				splits := strings.Split(cls.Image, ":")
				splits[len(splits)-1] = "latest"
				return strings.Join(splits, ":")
			}
		}
		return cls.Image
	}
	return ""
}
