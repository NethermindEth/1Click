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

// All the strings that are needed for error logging.
const (
	ReadingInstructionError            = "failed to read instructions from file %s"
	IncorrectClientError               = "incorrect %s client name \"%s\". Please provide correct client name. Use 'clients' command to see the list of supported clients"
	ClosingFileError                   = "failed to close file %s"
	ScriptIsNotRunningError            = "services of docker-compose script provided are not running. Error: %v"
	GettingLogsError                   = "failed to get logs for services %s. Error: %v"
	DockerComposePsReturnedEmptyError  = "'docker compose ps --services' returned empty string"
	InvalidVolumePathError             = "invalid path provided: %s. If you intended to pass a host directory, use absolute path"
	ZipError                           = "all lists must have the same size"
	CommandError                       = "command '%s' throws error: %v"
	DistroInfoError                    = "failed to get linux distribution info. Error: %v"
	EmptyClientMapError                = "is not possible to select a random element from an empty collection"
	NoSupportedClientsError            = "collection of clients given for random choice doesn't have any supported client. Check the target network (flag --network). Use 'clients' command to see the list of supported clients for every supported network"
	NetworkValidationFailedError       = "'network' flag validation failed. Error: %v"
	UnknownNetworkError                = "unknown network \"%s\". Please provide correct network name. Use 'networks' command to see the list of supported networks"
	GenerateJWTSecretError             = "JWT secret generation failed. Error: %v"
	GetPWDError                        = "something failed trying to get current working directory. Error: %v"
	EmptyFeeRecipientError             = "you should provide an Ethereum address for the Fee Recipient"
	KeystorePasswordError              = "keystore password must have more than 8 characters"
	PortOccupationError                = "port occupation check failed. Error: %v"
	DefaultPortInvalidError            = "default %s can not be zero"
	PrintFileError                     = "error printing file content: %v"
	CleaningEnvFileError               = "error cleaning env file: %v"
	CleaningDCFileError                = "error cleaning docker compose file: %v"
	PassphraseReadFileError            = "error reading passphrase file: %v"
	MnemonicReadFileError              = "error reading passphrase file: %v"
	MnemonicGenerationError            = "error creating mnemonic: %v"
	KeyEntryGenerationError            = "error generating keystore: could not read sufficient secure random bytes"
	AESParamsCreationError             = "failed to create AES128CTR params: %w"
	SecretEncryptionError              = "failed to encrypt secret: %w"
	KeystoreOutputExistingError        = "output folder for keystores already exists"
	KeystoreGenerationError            = "error generating keystores: %v"
	KeystoreDerivationError            = "keystore %s cannot be derived, continuing to next keystore"
	KeystoreExistingInWalletError      = "keystore with name \"%s\" already exists"
	KeystoreImportingError             = "failed to import keystore with pubkey %s into output wallet: %v"
	InvalidMnemonicError               = "mnemonic is not valid"
	BadMnemonicError                   = "bad mnemonic: %v"
	ForkVersionDecodeError             = "cannot decode fork version: %v"
	DepositFileWriteError              = "cannot write deposit file: %v"
	KeystoreSecretKeyCreationError     = "failed to create validator private key for path %q: %v"
	WithdrawalSecretKeyCreationError   = "failed to create withdrawal private key for path %q: %v"
	KeystoreSecretKeyConvertionError   = "cannot convert validator priv key: %v"
	DepositDataEncodingError           = "could not encode deposit data to json: %v"
	InvalidLoggingFlag                 = "bad logging flag: %v"
	CannotGenerateSecret               = "cannot generate 32 bytes long secret"
	ShowMnemonicError                  = "error displaying mnemonic: %v"
	InvalidFilePathOrUrl               = "invalid filepath or url: %s"
	CannotGetUrlContent                = "cannot get url %s content: %v"
	CannotReadFileContent              = "cannot read file %s content: %v"
	ErrorCheckingFile                  = "error checking file %s: %v"
	DestFileAlreadyExist               = "destiny file %s already exist"
	ErrorCreatingFile                  = "error creating file %s: %v"
	ErrorDownloadingFile               = "error downloading file from %s: %v"
	ErrorCopyingFile                   = "error copying file from %s: %v"
	ErrorWritingDeployBlockFile        = "error writing custom deploy block file %s: %v"
	InvalidUrlFlagError                = "invalid %s url %s. URL must be in the format http(s)://<host>:<port>/<api>/<endpoint>/... or http://<host>/<api>/<endpoint>/"
	InvalidEnodeError                  = "invalid enode %s. Bootnode must be in the format enode://<node id>@<host>:<port>"
	InvalidEnrError                    = "invalid enr %s. ENR must be in the format enr:<base64 encoded string>"
	InvalidService                     = "provided service %s is invalid"
	SetupContainersErr                 = "error setting up service containers: %w"
	StartingContainersErr              = "error starting service containers: %w"
	ReadingYmlErr                      = "error reading yml file: %w"
	ParsingYmlErr                      = "error parsing yml file, it seems is not a valid docker-compose script: %w"
	ServicesNotFoundErr                = "services not found in the docker-compose script"
	InvalidComposeErr                  = "provided docker-compose script is invalid: %w"
	ErrDuplicatedBootNode              = "duplicated boot node"
	ErrGraffitiLength                  = "graffiti must have 16 characters at most. Provided graffiti %s has %d characters"
	ErrCMDArgsNotSupported             = "command %s does not support arguments. Please use flags instead"
	ErrWithdrawalEth1SecretKeyCreation = "failed to create withdrawal private key for address 0x%s: %v"
	ErrWithdrawalBLSSecretKeyCreation  = "failed to create withdrawal private key for path %q: %v"
	ErrInvalidWithdrawalAddr           = "provided withdrawal address is not a valid Ethereum address"
	IncompatibleLidoAndEth1Withdrawal  = "incompatible flags --lido, and --eth-withdrawal-address can't be used together"
	InvalidNetworkForLido              = "invalid network: Choose valid network for Lido: %v"
	InvalidNetworkForLidoMevBoost      = "invalid network: Choose valid network for Lido with MEV-Boost: %v"
	InvalidNetworkForLidoKeys          = "invalid network: Choose valid network for Lido Withdrawal Address: %v"
	SetUpContainersErr                 = "error setting up containers %w"
	RunContainersErr                   = "error running containers %w"
)
