# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- New cli flag --distributed for running cluster with Charon distributed validator

### Changed
- Override `--latest` flag to not use the latest version of the image in the clients if image is specified

## [v1.7.2] - 2024-11-12

### Fixed
- Fix Erigon image source
- Fix installation script path 

## [v1.7.1] - 2024-11-1

### Added
- Updates on support for `op-geth` and `op-reth` clients on Optimism and Base networks.

### Changed
- Updated Sedge's Docker commands internal functionality.
- Changed the `--op-execution-image` flag to `--op-execution` on the `generate op-full-node` command.

## [v1.7.0] - 2024-10-24

### Added
- Support Mainnet for Lido CSM setup.

### Fixed
- Fix missing equals sign when setting builder on Lodestar.

## [v1.6.0] - 2024-10-18


### Added
- New command `lido-status` to display data of Lido Node Operator.
- New command `monitoring` to run monitoring stack setup with Grafana, Prometheus, Node Exporter and Lido Exporter.
- Security policy.
- Support for Nimbus as Consensus and Validator client. 
- Documentation for Lido Exporter and Monitoring Stack.

### Changed
- Update Go version from 1.21 to 1.22.
- Update documentation versions and dependencies.
- Update Lido Mainnet Vetted Relays List.
- Remove the need of users cloning the repository to run `make generate` and `make compile`, rolling back to only `make compile`
- Update client images.

### Fixed
-  Teku and Lighthouse import keys container error on Windows.
-  Security issues on dependencies.
-  Typos on documentation.
-  Fix `sync-mode` command on besu.

## [v1.5.0] - 2024-09-06

### Added
- Add support for Optimism and Base, using Nethermind Client on Mainnet and Sepolia.

### Changed
- Update client images to latest versions.

#### Fixed
- Remove Peer upper limit of peers on CL

##  [v1.4.0] - 2024-07-10

### Fixed
- Remove arguments for Lighthouse `--eth1` flag.

### Added

- Add support for MEV-boost on Holesky.
- New flag `--lido` to `generate` command for Lido CSM setup.
- New Sedge setup flow with `sedge cli` command for Lido CSM setup.
- Support for `sedge keys` to generate 0x01 withdrawal credentials.
- Support for installing Docker in Ubuntu 24.04 LTS.

### Changed

- Update client images to Dencun-ready versions.
- Renamed `--eth1-withdrawal-address` flag from `sedge keys` to `--eth-withdrawal-address`.
- Update client images to latest versions.

### Removed

- Removed support for Goerli.


## [v1.3.2] - 2024-03-08

### Changed

- Update client images to Dencun-ready versions. 

## [v1.3.1] - 2024-02-14

### Fixed

- Remove the Lodestar `--eth1.providerUrls` flag
- Replace the Lighthouse `--eth1-endpoints` flag with `--eth1`.

### Changed

- Update client images.

## [v1.3.0] - 2023-12-11

### Added

- New `--latest` flag to `sedge generate` subcommands to use the latest version of the image in the clients.
- Integrate Holesky network.

### Fixed

- Remove TTD of the params, test cases, command flags and geth patch for custom testnets.

### Changed

- Remove default bootnodes for networks managed by clients.
- Remove custom configurations for Chiado and instead use the clients configuration.
- Update client images.

## [v1.2.3] - 2023-08-16

### Added

- Support for Erigon in Gnosis.

### Fixed

- Issue in Chiado templates with checkpoint sync.
- Fix show error when using containers tags.

### Changed

- Update client images.
- Moved xdai to gnosis on Nethermind config.

## [v1.2.2] - 2023-07-24

### Changed

- Update client images to Shapella-ready version for Gnosis.

### Fixed

- Issue in Chiado templates with checkpoint sync.
- Fix show error when using containers tags.
- Fix chiado checkpoint sync url.

## [v1.2.1] - 2023-06-22

### Changed

- Update client images.
- Update Nethermind's configuration to show new logging format and colors.

### Fixed

- Missing symbol in validator-blocker container main command.

## [v1.2.0] - 2023-06-06

### Added

- New command to show Sedge container info: `sedge show`.
- Add default [checkpoint sync url]( https://checkpoint.chiadochain.net ) for Chiado.
- New `--skip-pull` flag to skip pulling docker images when running `sedge run`.

### Changed

- Update client images.

### Fixed

- Change validator blocker container image to [busybox](https://hub.docker.com/_/busybox).
- Erigon command line flags.

## [v1.1.0] - 2023-04-07

### Added

- Unit tests.

### Changed

- Update client images to Shapella-ready version.
- Update Nethermind client settings. Use default JSON-RPC modules, Prunning Cache size, and Snap Sync mode.
- Update Goerli checkpoint sync to use [this]( https://goerli.checkpoint-sync.ethpandaops.io).

### Fixed

- Fix Erigon command line flags.
- Checkpoint sync URL prompt is not longer mandatory.

## [v1.0.0] - 2023-3-23

### Added

- Add `generate` command to generate only an execution, consensus, or validator setup code.
- Support import slashing protection interchange data [EIP-3076](https://eips.ethereum.org/EIPS/eip-3076)
  while running the setup with the `cli` command or with the new command `slashing-import`.
- Support export slashing protection interchange data [EIP-3076](https://eips.ethereum.org/EIPS/eip-3076)
  with the new command `slashing-export`.
- Support the new command `run`, used to run all the services generated.
- Support for PPA packaging (apt install).
- Add flag `--container-tag` to add a suffix to sedge containers name.
- Support the new command `import-key` to import validator keys in an existing configuration.
- New command to check dependencies: `sedge deps check`
- New command to install dependencies: `sedge deps install`
- All the commands that needs dependencies will check if they are installed
  as a pre-requisite. If not, the command will fail and suggest to run `sedge deps check`.
- Documentation for all supported networks and clients.
- Add guides and descriptions for new commands in the Documentation.
- When generating new mnemonic, show it without a trace in the Terminal.

### Changed

- Improved test coverage
- Updated client versions
- Refactor `sedge cli` command to use the new sedge actions in a more interactive way.
- Rename sedge data directory name to `sedge-data` instead of `docker-compose-scripts`.
- Update mev-boost relay URLs.
- Fixed validator restart failure. Validator never restarts, but has a better and safer start-up method:
  - Validator waits a grace period (2 epochs) before starting.
  - Validator waits for the consensus sync to finish: verifying that the `/eth/v1/node/health` consensus endpoint
      returns with the HTTP code 200 in the newly added `validator-blocker` docker-compose service. This replace the
      older track sync method.
- Updated installation script for docker and docker compose in Linux flavours.
- Updated Checkpoint Sync Url for Goerli.
- Use new created action for JWT secret handling.
- Removed `v` from the tag on release scripts. Now we should use `Tag=v1.0.0` instead of `Tag=1.0.0`.

### Removed

- Removed config file dependency.
- Removed `prysm` from consensus templates on Gnosis as supported client.

### Fixed

- Add missing params at teku validator template.
- Remove double params at prysm validator template.

## [v0.6.0] - 2022-12-23

### Added

- Support for Erigon as Execution Client.
- Support for Besu as Execution Client.
- Support for Windows.
- Add search functionality to documentation using Algolia.

### Changed

- Completely replace keystore generation mechanism for Ethereum mainnet.

### Fixed

- Revert Erigon Image from `v2.32.0` to `v2.29.0`.
- Bad `{{end}}` tag on docker-compose_base template.
- Inconsistent behavior of `AssignPorts` function test.
- Dependencies install script bug for Ubuntu 22.04.

## [v0.5.1] - 2022-12-2

### Fixed

- Fix missing preset flag for Lodestar validator service.
- Support custom configuration for Lighthouse validator-import service.

## [v0.5.0] - 2022-11-29

### Added

- Add checkpoint sync url for Chiado.

### Changed

- Update Gnosis and Chiado networks default clients images to merge ready versions.
- Update client versions.
- Update checkpoint sync url for gnosis.

### Fixed

- Add missing `depends_on` tag to lodestar validator.
- Fix wrong fork version in Gnosis network config.

## [v0.4.0] - 2022-10-25

### Added

- Check for new Sedge releases on GitHub. Sedge will now report if there is a new version released.
- Add `--mev-boost-image` flag to specify which Mev Boost docker image sedge uses.
- Set mev-boost as default on mainnet.
- Documentation for How to Run a Validator Node on Sedge by yourself, once generated.
- Documentation for How to Run a Validator or Full Node on Chiado, Gnosis testnet.
- validator-import service for Lodestar.

### Changed

- Drop support for deprecated networks. (Kiln, Ropsten, Denver).
- Update Chiado testnet configs.
- Sedge uses its own mechanism for keystore generation with `sedge keys` instead of the staking-deposit-cli tool. This
  is experimental and staking-deposit-cli is still being used for mainnet.
- Updated docker image tags for all clients except Teku.

### Fixed

- Print generated files as string instead of []byte.
- keystore_password.txt permissions issue.

## [v0.3.0] - 2022-09-15

### Added

- Add support for gnosis merge testnets named chiado and denver.
- Add support for custom testnets by allowing custom remote config and genesis files.
- Add `--no-validator` flag to exclude the validator node from the full node setup.
- Gnosis Network support.
- Add `--graffiti` to allow overriding the default graffiti used by validators nodes.
- Allow the extra flags to overwrite fixed template images commands flags. Latest flag apparition will be used.
- Add Homebrew installation method.
- Add Checkpoint Sync for Ethereum Mainnet.

### Changed

- Remove Nethermind metrics configuration.
- Prysm and Lodestar can use now an undefined number of Bootnodes.
- Allow providing a number for `--tail` flag of `logs` command, e.g: `sedge logs --tail 50`
- Remove `mainnet` service templates and use `merge` service templates for Mainnet network.
- `checkpoint-sync-url` can now be set separately for every consensus client of a network using `CHECKPOINT_SYNC_URL` as
  env variabe in the consensus client templates.
- Use fixed docker images with latest tested versions.
- Change Prater to Goerli.

### Fixed

- Check for compose plugin installed. The check was running compose with a wrong path flag.
- Errors in validators templates.

## [v0.2.0] - 2022-08-4

### Added

- Support for Goerli/Prater network.
- Docusaurus documentation engine with initial version of docs.
- CI workflows to test and deploy documentation to Github Pages.
- Goerli/Prater support.
- Sepolia support.
- Mev-boost settings for all the currently supported CL clients.
- Relays URLs for all the testnets.
- Mev-boost set on by default in Prater.

### Changed

- Update Readme. Some instructions were moved to the docs. Roadmap updated.
- Updated Teku mev-boost configuration.

## [v0.1.2] - 2022-07-21

### Added

- Table with mev-boost support information on the Readme.
- Table with OS support for dependency installation.
- More instructions and descriptions about the tool on the Readme.

### Changed

- Update Roadmap.

### Fixed

- Fix bad checks for compose availability.
- Stop running Lighthouse with mev-boost settings on Ropsten. Official and stable docker image doesn't support mev-boost
  yet.
- Fix compose installation. Now compose is installed on root user, as sedge runs everything using sudo currently.
- Fix track sync of nodes stoping early.

## [v0.1.1] - 2022-07-20

### Changed

- Update Roadmap.

### Fixed

- Fix error when trying to generate jwtsecret on unexisting folder.

## [v0.1.0] - 2022-07-15

### Added

- Create cli tool able to set up a Ethereum based validator in an on-premise way.
- Generate `docker-compose` scripts and `.env` files for selected clients with the `cli` command.
- Generate keystore folder with the cli using `keys` command. The inner tool used for this is
  the [staking-deposit-cli](https://github.com/ethereum/staking-deposit-cli) tool.
- Generate `jwtsecret` for post-merge networks.
- Test coverage (unit tests).
- Integrate Kiln network.
- Integrate Ropsten network.
- Integrate MEV-Boost as an option.
