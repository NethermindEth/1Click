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
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/clients"
	"github.com/NethermindEth/sedge/internal/pkg/env"
	"github.com/NethermindEth/sedge/internal/utils"
	"github.com/NethermindEth/sedge/templates"
	log "github.com/sirupsen/logrus"
)

/*
GenerateScripts :
This function is responsible for generating docker-compose files for execution, consensus and
validator clients.

params :-
a. gd GenerationData
Data object containing clients whose script are to be generated, path of generated scripts and special options for the clients configuration.

returns :-
a. string
Execution client json-rpc API port
b. string
Consensus client HTTP API port
a. error
Error if any
*/
func GenerateScripts(gd GenerationData) (result GenerationResults, err error) {
	// Create scripts directory if not exists
	if _, err := os.Stat(gd.GenerationPath); os.IsNotExist(err) {
		err = os.MkdirAll(gd.GenerationPath, 0o755)
		if err != nil {
			return GenerationResults{}, err
		}
	}

	// Check for port occupation
	defaultsPorts := map[string]string{
		"ELDiscovery":     configs.DefaultDiscoveryPortEL,
		"ELMetrics":       configs.DefaultMetricsPortEL,
		"ELApi":           configs.DefaultApiPortEL,
		"ELAuth":          configs.DefaultAuthPortEL,
		"ELWS":            configs.DefaultWSPortEL,
		"CLDiscovery":     configs.DefaultDiscoveryPortCL,
		"CLMetrics":       configs.DefaultMetricsPortCL,
		"CLApi":           configs.DefaultApiPortCL,
		"CLAdditionalApi": configs.DefaultAdditionalApiPortCL,
		"VLMetrics":       configs.DefaultMetricsPortVL,
		"MevPort":         configs.DefaultMevPort,
	}
	ports, err := utils.AssignPorts("localhost", defaultsPorts)
	if err != nil {
		return GenerationResults{}, fmt.Errorf(configs.PortOccupationError, err)
	}
	gd.Ports = ports
	// External endpoints will be configured here. Also Ports should be updated with external ports
	gd.ExecutionClient.Endpoint = configs.OnPremiseExecutionURL
	gd.ConsensusClient.Endpoint = configs.OnPremiseConsensusURL

	log.Info(configs.GeneratingEnvFile)
	envFilePath, err := generateEnvFile(gd)
	if err != nil {
		return GenerationResults{}, err
	}

	if err := CleanEnvFile(envFilePath); err != nil {
		return GenerationResults{}, err
	}

	log.Info(configs.GeneratingDockerComposeScript)
	dockerComposePath, err := generateDockerComposeScripts(gd, envFilePath)
	if err != nil {
		return GenerationResults{}, err
	}

	if err := CleanDockerCompose(dockerComposePath); err != nil {
		return GenerationResults{}, err
	}

	return GenerationResults{
		ELPort:            ports["ELApi"],
		CLPort:            ports["CLApi"],
		EnvFilePath:       envFilePath,
		DockerComposePath: dockerComposePath,
	}, nil
}

/*
generateDockerComposeScripts :
This function is responsible for generating docker-compose scripts for execution, consensus and
validator clients.

params :-
a. executionClient string
Execution client whose script is to be generated
b. consensusClient string
Execution client whose script is to be generated
c. validatorClient string
Execution client whose script is to be generated
d. path string
Path of generated scripts

returns :-
a. error
Error if any
*/
func generateDockerComposeScripts(gd GenerationData, envFilePath string) (dockerComposePath string, err error) {
	rawBaseTmp, err := templates.Services.ReadFile(strings.Join([]string{"services", "docker-compose_base.tmpl"}, "/"))
	if err != nil {
		return
	}

	baseTmp, err := template.New("docker-compose").Parse(string(rawBaseTmp))
	if err != nil {
		return
	}

	clients := map[string]clients.Client{
		"execution": gd.ExecutionClient,
		"consensus": gd.ConsensusClient,
		"validator": gd.ValidatorClient,
	}
	for tmpKind, client := range clients {
		name := client.Name
		if client.Omited {
			name = "empty"
		}
		tmp, err := templates.Services.ReadFile(strings.Join([]string{
			"services",
			configs.NetworksConfigs[gd.Network].NetworkService,
			tmpKind,
			name + ".tmpl",
		}, "/"))
		if err != nil {
			return "", err
		}
		_, err = baseTmp.Parse(string(tmp))
		if err != nil {
			return "", err
		}
	}

	// Parse validator-blocker template
	tmp, err := templates.Services.ReadFile(filepath.Join("services", "validator-blocker.tmpl"))
	if err != nil {
		return "", err
	}
	if _, err := baseTmp.Parse(string(tmp)); err != nil {
		return "", err
	}

	// Check for TTD in generated env file
	TTD, err := env.GetTTD(envFilePath)
	if err != nil {
		return "", err
	}

	// Check for splitted network flags
	splittedNetwork, err := env.CheckVariableBase(env.ReSPLITTED, gd.Network)
	if err != nil {
		return "", err
	}

	// Check for custom network config
	ccRemoteCfg, err := env.CheckVariable(env.ReCONFIG, gd.Network, "consensus", gd.ConsensusClient.Name)
	if err != nil {
		return "", err
	}
	if ccRemoteCfg { // Have to download custom configs
		if gd.CustomNetworkConfigPath == "" { // Setup paths to use downloaded configs
			gd.CustomNetworkConfigPath = "./config.yaml"
		} else { // Overriding custom configs, no need to download
			ccRemoteCfg = false
		}
	}
	ccRemoteGen, err := env.CheckVariable(env.ReGENESIS, gd.Network, "consensus", gd.ConsensusClient.Name)
	if err != nil {
		return "", err
	}
	if ccRemoteGen { // Have to download custom configs
		if gd.CustomGenesisPath == "" { // Setup paths to use downloaded configs
			gd.CustomGenesisPath = "./genesis.ssz"
		} else { // Overriding custom configs, no need to download
			ccRemoteGen = false
		}
	}
	ccRemoteDpl, err := env.CheckVariable(env.ReDEPLOY, gd.Network, "consensus", gd.ConsensusClient.Name)
	if err != nil {
		return "", err
	}
	if ccRemoteDpl { // Have to download custom configs
		if gd.CustomDeployBlockPath == "" { // Setup paths to use downloaded configs
			gd.CustomDeployBlockPath = "./deploy_block.txt"
		} else { // Overriding custom configs, no need to download
			ccRemoteDpl = false
		}
	}

	vlRemoteCfg, err := env.CheckVariable(env.ReCONFIG, gd.Network, "validator", gd.ValidatorClient.Name)
	if err != nil {
		return "", err
	}
	if vlRemoteCfg { // Have to download custom configs
		if gd.CustomNetworkConfigPath == "" { // Setup paths to use downloaded configs
			gd.CustomNetworkConfigPath = "config.yaml"
		} else { // Overriding custom configs, no need to download
			vlRemoteCfg = false
		}
	}
	vlRemoteGen, err := env.CheckVariable(env.ReGENESIS, gd.Network, "validator", gd.ValidatorClient.Name)
	if err != nil {
		return "", err
	}
	if vlRemoteGen { // Have to download custom configs
		if gd.CustomGenesisPath == "" { // Setup paths to use downloaded configs
			gd.CustomGenesisPath = "./genesis.ssz"
		} else { // Overriding custom configs, no need to download
			vlRemoteGen = false
		}
	}
	vlRemoteDpl, err := env.CheckVariable(env.ReDEPLOY, gd.Network, "validator", gd.ValidatorClient.Name)
	if err != nil {
		return "", err
	}
	if vlRemoteDpl { // Have to download custom configs
		if gd.CustomDeployBlockPath == "" { // Setup paths to use downloaded configs
			gd.CustomDeployBlockPath = "./deploy_block.txt"
		} else { // Overriding custom configs, no need to download
			vlRemoteDpl = false
		}
	}

	// Check for XEE_VERSION in teku
	xeeVersion, err := env.CheckVariable(env.ReXEEV, gd.Network, "consensus", gd.ConsensusClient.Name)
	if err != nil {
		return "", err
	}

	// Check for Mev
	mev, err := env.CheckVariable(env.ReMEV, gd.Network, "validator", gd.ValidatorClient.Name)
	if err != nil {
		return "", err
	}

	// Check for CC Bootnode nodes
	ccBootnodes, err := env.GetCCBootnodes(envFilePath)
	if err != nil {
		return "", err
	}

	// Check for Bootnode nodes
	ecBootnodes, err := env.GetECBootnodes(envFilePath)
	if err != nil {
		return "", err
	}

	clCheckpointSyncUrl, err := env.CheckVariable(env.ReCHECKPOINT, gd.Network, "consensus", gd.ConsensusClient.Name)
	if err != nil {
		return "", err
	}

	data := DockerComposeData{
		Services:            gd.Services,
		TTD:                 TTD != "",
		XeeVersion:          xeeVersion,
		Mev:                 mev && gd.Mev,
		MevPort:             gd.Ports["MevPort"],
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
		FallbackELUrls:      gd.FallbackELUrls,
		ElExtraFlags:        gd.ElExtraFlags,
		ClExtraFlags:        gd.ClExtraFlags,
		VlExtraFlags:        gd.VlExtraFlags,
		ECBootnodes:         ecBootnodes,
		CCBootnodes:         ccBootnodes,
		MapAllPorts:         gd.MapAllPorts,
		SplittedNetwork:     splittedNetwork,
		ClCheckpointSyncUrl: clCheckpointSyncUrl,
		LoggingDriver:       gd.LoggingDriver,
		VLStartGracePeriod:  gd.VLStartGracePeriod,
		// Custom configs data
		CustomCfgDownload: ccRemoteCfg || ccRemoteGen || ccRemoteDpl ||
			vlRemoteCfg || vlRemoteGen || vlRemoteDpl, // Need to download something for consensus node
		RemoteCfg:     ccRemoteCfg || vlRemoteCfg,               // Download config.yaml
		RemoteGen:     ccRemoteGen || vlRemoteGen,               // Download genesis.ssz
		RemoteDpl:     ccRemoteDpl || vlRemoteDpl,               // Download deploy_block.txt
		CustomNetwork: gd.Network == configs.CustomNetwork.Name, // Used custom templates
		CustomConsensusConfigs: gd.CustomNetworkConfigPath != "" ||
			gd.CustomGenesisPath != "" ||
			gd.CustomDeployBlockPath != "", // Have custom configs paths
		CustomChainSpecPath:     gd.CustomChainSpecPath,     // Path to chainspec.json
		CustomNetworkConfigPath: gd.CustomNetworkConfigPath, // Path to config.yaml
		CustomGenesisPath:       gd.CustomGenesisPath,       // Path to genesis.ssz
		CustomDeployBlockPath:   gd.CustomDeployBlockPath,   // Path to deploy_block.txt
	}

	dockerComposePath = filepath.Join(gd.GenerationPath, configs.DefaultDockerComposeScriptName)

	err = writeTemplateToFile(baseTmp, dockerComposePath, data, false)
	if err != nil {
		return "", fmt.Errorf(configs.GeneratingScriptsError, gd.ExecutionClient.Name, gd.ConsensusClient.Name, gd.ValidatorClient.Name, err)
	}

	return dockerComposePath, nil
}

/*
generateEnvFile :
This function is responsible for generating the environment variable for the
generated docker-compose scripts for execution, consensus and
validator clients.

params :-
a. executionClient string
Execution client whose script was generated
b. consensusClient string
Execution client whose script was generated
c. validatorClient string
Execution client whose script was generated
d. path string
Path of generated scripts

returns :-
a. error
Error if any
*/
func generateEnvFile(gd GenerationData) (envFilePath string, err error) {
	rawBaseTmp, err := templates.Envs.ReadFile(strings.Join([]string{"envs", gd.Network, "env_base.tmpl"}, "/"))
	if err != nil {
		return
	}

	baseTmp, err := template.New("env").Parse(string(rawBaseTmp))
	if err != nil {
		return
	}

	clients := map[string]clients.Client{
		"execution": gd.ExecutionClient,
		"consensus": gd.ConsensusClient,
		"validator": gd.ValidatorClient,
	}
	for tmpKind, client := range clients {
		var tmp []byte
		if client.Omited {
			tmp, err = templates.Services.ReadFile(strings.Join([]string{
				"services",
				configs.NetworksConfigs[gd.Network].NetworkService,
				tmpKind,
				"empty.tmpl",
			}, "/"))
			if err != nil {
				return "", err
			}
		} else {
			tmp, err = templates.Envs.ReadFile(strings.Join([]string{"envs", gd.Network, tmpKind, client.Name + ".tmpl"}, "/"))
			if err != nil {
				return "", err
			}
		}
		_, err = baseTmp.Parse(string(tmp))
		if err != nil {
			return "", err
		}
	}

	// TODO: Use OS wise delimiter for these data structs
	data := EnvData{
		ElImage:                   gd.ExecutionClient.Image,
		ElDataDir:                 "./" + configs.ExecutionDir,
		CcImage:                   gd.ConsensusClient.Image,
		CcDataDir:                 "./" + configs.ConsensusDir,
		VlImage:                   gd.ValidatorClient.Image,
		VlDataDir:                 "./" + configs.ValidatorDir,
		ExecutionApiURL:           gd.ExecutionClient.Endpoint + ":" + gd.Ports["ELApi"],
		ExecutionAuthURL:          gd.ExecutionClient.Endpoint + ":" + gd.Ports["ELAuth"],
		ConsensusApiURL:           gd.ConsensusClient.Endpoint + ":" + gd.Ports["CLApi"],
		ConsensusAdditionalApiURL: gd.ConsensusClient.Endpoint + ":" + gd.Ports["CLAdditionalApi"],
		FeeRecipient:              gd.FeeRecipient,
		JWTSecretPath:             gd.JWTSecretPath,
		ExecutionEngineName:       gd.ExecutionClient.Name,
		ConsensusClientName:       gd.ConsensusClient.Name,
		KeystoreDir:               "./" + configs.KeystoreDir,
		Graffiti:                  gd.Graffiti,
		ECBootnodes:               gd.ECBootnodes,
		CCBootnodes:               gd.CCBootnodes,
		CustomTTD:                 gd.CustomTTD,
		CustomDeployBlock:         gd.CustomDeployBlock,
	}

	// Fix prysm rpc url
	if gd.ValidatorClient.Name == "prysm" {
		data.ConsensusAdditionalApiURL = fmt.Sprintf("%s:%s", "consensus", gd.Ports["CLAdditionalApi"])
	}

	envFilePath = filepath.Join(gd.GenerationPath, ".env")

	err = writeTemplateToFile(baseTmp, envFilePath, data, false)
	if err != nil {
		return "", fmt.Errorf(configs.GeneratingScriptsError, gd.ExecutionClient.Name, gd.ConsensusClient.Name, gd.ValidatorClient.Name, err)
	}

	return envFilePath, nil
}

/*
writeTemplateToFile :
Write template to `file`.

params :-
a. template *template.Template
Template to be written
b. file string
File's complete path
c. data interface{}
Data object to be applied to `template`
d. append bool
True to append the template to `file`. False to create it or overwrite it.

returns :-
a. err error
Error if any
*/
func writeTemplateToFile(template *template.Template, file string, data interface{}, append bool) (err error) {
	var f *os.File

	if append {
		f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0o666)
		if err != nil {
			return fmt.Errorf(configs.CreatingFileError, file, err)
		}
	} else {
		f, err = os.Create(file)
		if err != nil {
			return fmt.Errorf(configs.CreatingFileError, file, err)
		}
	}

	// Just closing a file without checking any closing errors is a bad practice
	defer func() {
		cerr := f.Close()
		if err == nil && cerr != nil {
			log.Errorf(configs.ClosingFileError, file)
			err = cerr
		}
	}()

	err = template.Execute(f, data)
	if err != nil {
		return
	}

	return nil
}
