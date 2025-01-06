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
package configs

import (
	"fmt"
	"os"
	"path/filepath"
)

// All the strings that are needed for debugging and info logging, and constant strings.
const (
	DefaultMevBoostEndpoint         = "http://mev-boost"
	DefaultEnvFileName              = ".env"
	CheckingDependencies            = "Checking dependencies: %s"
	DependenciesPending             = "pending dependencies: %s"
	DependenciesOK                  = "All dependencies are installed on host machine"
	GeneratingDockerComposeScript   = "Generating docker-compose script for current selection of clients"
	GeneratingEnvFile               = "Generating environment file for current selection of clients"
	GeneratedDockerComposeScript    = "Generated docker-compose script for current selection of clients"
	GeneratedEnvFile                = "Generated environment file for current selection of clients"
	CleaningGeneratedFiles          = "Cleaning generated docker-compose and environment files"
	CleanedGeneratedFiles           = "Cleaned generated files"
	GenerationEnd                   = "Generation of files successfully, happy staking! You can use now 'sedge run' to start the setup."
	Exiting                         = "Exiting..."
	InstructionsFor                 = "Instructions for %s"
	OSNotSupported                  = "installation not supported for %s"
	ProvideClients                  = "Please provide both execution client and consensus client"
	CreatedFile                     = "Created file %s"
	DefaultSedgeDataFolderName      = "sedge-data"
	ClientNotSupported              = "client %s is not supported. Please use 'clients' command to see the list of supported clients"
	PrintingFile                    = "File %s:"
	SupportedClients                = "Supported clients of type %s: %v"
	ConfigClientsMsg                = "Provided clients of type %s in configuration file: %v"
	RunningDockerCompose            = "Running docker-compose script"
	Component                       = "component"
	RunningCommand                  = "Running command: %s"
	UnableToProceed                 = "Unable to proceed. Please check the logs for more details"
	CheckingDockerEngine            = "Checking if docker engine is on"
	DepositCLIDockerImageUrl        = "nethermindeth/staking-deposit-cli" //"github.com/ethereum/staking-deposit-cli"
	DepositCLIDockerImageName       = "nethermindeth/staking-deposit-cli" //"deposit-cli:local"
	GeneratingKeystores             = "Generating keystores..."
	GeneratingKeystoresLegacy       = "Generating keystore folder"
	KeystoresGenerated              = "Keystores generated successfully"
	GeneratingDepositData           = "Generating deposit data..."
	DepositDataGenerated            = "Deposit data generated successfully"
	KeystorePath                    = "Keystore path: %s"
	KeysFoundAt                     = "If everything went well, your keys can be found at: %s"
	ImageNotFoundBuilding           = "Image %s not found, building it"
	ImageNotFoundPulling            = "Image %s not found, pulling it"
	WithdrawalAddressDefinedWarning = "You have defined a withdrawal address for your validators. Make sure this is intended. Deposits made with this withdrawal address cannot be changed. If you want to change the withdrawal address, you will have to exit your validators and create new validators with the new withdrawal address."

	ReviewKeystorePath = "In case you used custom paths for the 'cli' or the 'keys' commands, please review if the keystore path in the generated .env file points to the generated keystore folder (the .env key should be KEYSTORE_DIR). If not, change the path in the .env file to the correct one."
	NodesSynced        = "Execution and Consensus clients are synced, proceeding to start validator node"
	RemoteNodeNeeded   = `
If you want to run a validator, make sure you have access to a remote or external consensus node.
	
You can use one of your own, a friend's node or providers such as Infura. Edit the .env file accordingly please, check that variable CC_NODE <-> (consensus endpoint) have the correct value. The validator node requires a high available consensus node, and consensus in turn needs a high available execution node.`
	HappySedgingNoRun = `
Your setup is ready. You can run it anytime using the 'sedge run --path %s' command. Feel free to explore the files and make changes, although Sedge is not accountable for any misbehavior or issue caused by any modification done to the setup. Stay tuned for more updates and features!

Happy Sedging!
	`
	HappySedgingRun = `
Your setup is up and running. Thank you for joining and helping the community! You can check the logs of your nodes using the command 'sedge logs --path %s <node_type>'. Stay tuned for more updates and features!

Happy Sedging!
	`
	HappyStakingRun = `
Your full-node is up and running. If you set up new validator keys, you will have to register them. Follow https://launchpad.ethereum.org/ instructions to register the validator addresses of the validators you want to set up and manage using the validator node. If something goes wrong and your validator node is down, don't panic! Short downtimes are pretty decent. Check the logs and try to fix the errors. Ensure errors are not related to connection issues with the consensus nodes or to the keystore folder path or validator data directory. You can start the validator again using the instructions displayed by the tool.
	
Happy Staking!
	`
	ExecutionClientNotSpecifiedWarn         = "Execution client not provided. A random client will be selected. Random client: %s"
	ConsensusClientNotSpecifiedWarn         = "Consensus client not provided. Selecting same pair of clients for consensus and validator clients"
	ValidatorClientNotSpecifiedWarn         = "Validator client not provided. Selecting same pair of clients for consensus and validator clients"
	CLNotSpecifiedWarn                      = "Consensus and validator clients not provided. Selecting same pair of clients for consensus and validator clients using a random client. Random client: %s"
	GeneratingJWTSecret                     = "Generating JWT secret for client authentication"
	JWTSecretGenerated                      = "JWT secret generated"
	CreatingKeystorePassword                = "Creating keystore_password.txt on keystore folder"
	KeystorePasswordCreated                 = "keystore_password.txt on keystore folder created with provided password"
	MnemonicTips                            = "The following mnemonic is going to be used to create the validator keystore. Please save it carefully. It can be used to generate the keystore folder again. If you lose the password and mnemonic, access to your keystore will be lost forever!"
	GeneratingMnemonic                      = "Existing mnemonic not provided. Generating mnemonic for validator keystore:"
	StoreMnemonic                           = "Make sure to store your mnemonic somewhere safe. Losing it could end in the lost of your validators. Press enter to continue" // TODO: improve warning message
	PreparingTekuDatadir                    = "Preparing teku datadirs (must have full read/write/execute permissions to work)"
	GettingContainersIP                     = "Proceeding to get execution and consensus containers IP address for the monitoring tool"
	WaitingForNodesToStart                  = "Waiting a minute for nodes to start"
	CustomExecutionImagesWarning            = "You are using custom images for the execution client!!! Make sure this is intended. Also check these images are correct and available from this device otherwise the setup will fail or have an unexpected behavior."
	CustomConsensusImagesWarning            = "You are using custom images for the consensus client!!! Make sure this is intended. Also check these images are correct and available from this device otherwise the setup will fail or have an unexpected behavior."
	CustomValidatorImagesWarning            = "You are using custom images for the validator client!!! Make sure this is intended. Also check these images are correct and available from this device otherwise the setup will fail or have an unexpected behavior."
	CustomDistributedValidatorImagesWarning = "You are using custom images for the distributed validator client!!! Make sure this is intended. Also check these images are correct and available from this device otherwise the setup will fail or have an unexpected behavior."
	MapAllPortsWarning                      = "You are mapping all ports for the clients!!! Make sure this is intended. This could make the clients vulnerable to attacks. Be sure to setup a firewall."
	CheckpointUrlUsedWarning                = "A Checkpoint Sync Url will be used for the consensus node. Using %s ."
	NoBootnodesFound                        = "No bootnodes found in %s env file "
	UnableToCheckVersion                    = "Unable to check for new Version. Please check manually at " +
		"https://github.com/NethermindEth/sedge/releases, with error:"
	NeedVersionUpdate = "A new Version of sedge is available. Please update to the latest Version. See " +
		"https://github.com/NethermindEth/sedge/releases for more information. Latest detected tag:"
	VersionUpdated             = "You are running the latest version of sedge. Version: "
	Downloading                = "Downloading %s..."
	Copying                    = "Copying %s..."
	GettingCustomChainSpec     = "Getting custom chain spec..."
	GettingCustomGenesis       = "Getting custom genesis..."
	GettingCustomNetworkConfig = "Getting custom network config..."
	WritingCustomDeployBlock   = "Writing custom deploy block..."
)

var DefaultAbsSedgeDataPath string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		// notest
		fmt.Println(err)
	}
	DefaultAbsSedgeDataPath = filepath.Join(cwd, DefaultSedgeDataFolderName)
}
