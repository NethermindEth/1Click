<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://user-images.githubusercontent.com/35319980/245197671-464fc898-a6e5-4571-bf34-957c319501ec.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://user-images.githubusercontent.com/35319980/245197581-2cddd292-75c0-4f4a-b26a-a3dcd071ad63.svg">
    <img alt="Nethermind" src="https://user-images.githubusercontent.com/35319980/245197581-2cddd292-75c0-4f4a-b26a-a3dcd071ad63.svg" height="100">
  </picture>
</p>

# Sedge

[![Go Report Card](https://goreportcard.com/badge/github.com/NethermindEth/sedge)](https://goreportcard.com/report/github.com/NethermindEth/sedge)
[![Discord](https://user-images.githubusercontent.com/7288322/34471967-1df7808a-efbb-11e7-9088-ed0b04151291.png)](https://discord.com/invite/PaCMRFdvWT)
[![codecov](https://codecov.io/gh/NethermindEth/sedge/branch/main/graph/badge.svg?token=8FERO4PO1V)](https://codecov.io/gh/NethermindEth/sedge)


Sedge is a one-click node setup tool for PoS network/chain validators and nodes written entirely in the Go programming language. Sedge takes care of the entire on-premise full node setup based on the chosen client, using generated docker-compose scripts based on the desired configuration.

- [Sedge](#sedge)
  - [⚙️ Installation](#️-installation)
    - [Dependencies](#dependencies)
    - [Installation methods](#installation-methods)
    - [**Disclaimer**](#disclaimer)
  - [📜 Documentation](#-documentation)
  - [⚡️ Quick start](#️-quick-start)
  - [💥 Why did we start Sedge?](#-why-did-we-start-sedge)
  - [🔥 What can you do with Sedge today?](#-what-can-you-do-with-sedge-today)
    - [**Disclaimer**](#disclaimer-1)
  - [Supported networks and clients](#supported-networks-and-clients)
    - [Mainnet](#mainnet)
    - [Sepolia](#sepolia)
    - [Holesky](#holesky)
    - [Gnosis](#gnosis)
    - [Chiado (Gnosis testnet)](#chiado-gnosis-testnet)
    - [CL clients with Mev-Boost](#cl-clients-with-mev-boost)
  - [Supported Linux flavours for dependency installation](#supported-linux-flavours-for-dependency-installation)
  - [✅ Roadmap](#-roadmap)
    - [Version 0.1](#version-01)
    - [Version 0.2](#version-02)
    - [Version 0.3](#version-03)
    - [Version 0.4](#version-04)
    - [Version 0.5](#version-05)
    - [Version 0.6](#version-06)
    - [Version 1.0 (Current)](#version-10-current)
    - [Version 1.X](#version-1x)
  - [💪 Want to contribute?](#-want-to-contribute)
  - [⚠️ License](#️-license)

## ⚙️ Installation

### Dependencies

Sedge dependencies are `docker` with the `docker compose` plugin, but if you don't have those installed, Sedge will show instructions to install them or install them for you. Check the [docs](https://docs.sedge.nethermind.io/docs/quickstart/dependencies) for more details.

### Installation methods

Check our [installation guide](https://docs.sedge.nethermind.io/docs/quickstart/install-guide) for detailed instructions on the supported methods:

- Download the binary from the release page
- Using the Homebrew package manager
- Using the Go programming language
- Build from source

### **Disclaimer**

Downloading any binary from the internet risks downloading files that malicious, third-party actors have injected with malware. All users should check that they download the correct, clean binary from a reputable source.

## 📜 Documentation

You can check the [documentation](https://docs.sedge.nethermind.io) for further details.

## ⚡️ Quick start

With `sedge cli` you can go through the entire workflow setup:

1. Generate a `docker-compose.yml` script with randomized clients selection and `.env` file
2. Generate validator keystore, or import it.
3. Check dependencies
4. Execute the `docker-compose.yml` script (only execution and consensus nodes will be executed by default)
5. Validator client will be executed automatically after the consensus node is synced.
  
You can also generate the validator(s) keystore folder using `sedge keys`.

The entire process is interactive. However, Sedge also has a very customizable, non-interactive setup without prompts.

If you want to run the [non-interactive mode](https://docs.sedge.nethermind.io/docs/quickstart/complete-guide#22-non-interactive-setup), you will need to run only four commands (two if you are not running a validator and you don't need the validator keystore), and provide the set of arguments needed for each command.
1. `sedge generate`
2. `sedge keys`
3. `sedge import-key`
4. `sedge run`

Check all the options and flags with `sedge --help`.

## 💥 Why did we start Sedge?

As people who actively deployed validators way before The Merge, we know how hard it is to set up an Ethereum validator:

- You need to procure at least three (compatible) nodes: an execution node (geth, nethermind, erigon, etc), a consensus node, and a validator node (lighthouse, prysm, etc)
- You then need to execute them, connect them, monitor them, and secure the validator keys (which includes staking 32 ETH).
- There may be several valid combinations of clients to choose for your setup, so you need to go through each of the client's docs, evaluate it, get instructions for it and test it. You also need to feel comfortable executing commands in the cli, know docker, and understand basics of networking. On top of this, there are many different settings you must read up on and consider for your client node.
- In the case of working with the Ethereum Mainnet, you are working with real money that can potentially be lost in the event of having downtime or being slashed. To avoid losing real value, you must be aware of and follow best practices on the validator setup, and correctly monitor your nodes.
- Have you heard of MEV-Boost? You most likely want always to be running the latest version, but you also don’t have the time to understand the MEV-Boost architecture ins and outs or how to implement it into your environment successfully.
  
> Enter Sedge

We want Sedge to take care of all of the above for you. With just a few clicks or steps, Sedge can create an entire Ethereum staking architecture that supports client diversity and Ethereum's latest features, while being completely free and open source. We want Sedge to save you from making costly mistakes in this complex setup; along with hours or days of research, reading and testing. We want you to be able to stake easily with or without blockchain knowledge by giving you the tools to help this amazing community (and earn some good money of course 😉).

We want to share our knowledge in this topic and create something that allows everyone to easily and safely set up lots of diverse validators.

We don't want to stop at Ethereum. We also want to help stakers of other PoS networks/chains, so if your favourite chain is not here, you are more than welcome to contribute!

## 🔥 What can you do with Sedge today?

- Select an execution, consensus and validator client node (manually or automatically) and generate a `docker-compose.yml` script with production-tested configurations to run the setup you want.
- Set up only an execution, consensus, validator, or mev-boost instance.
- Don't remember `docker compose` commands or flags for your setup? Check docker logs of the running services with `sedge logs`, and shut them down with `sedge down`
- Generate the keystore folder with `sedge keys`  using our heavily tested own code.

> **Disclaimer:** Users acknowledge that generating the keystore for any network is an unaudited feature of Sedge. Nethermind provides this feature on an ‘as is’ basis and makes no warranties regarding its proper functioning. The use of this feature is at the user’s own risk - Nethermind excludes all liability for any malfunction or loss of money that may occur as the result of an unexpected behavior during the keystore generation.

If you are familiar with `docker`, `docker compose`, and the validator setup, then you can use `sedge cli` or `sedge generate` to create a base docker-compose script with the recommended settings and then edit the script as much as you want. It is a lot more easier than doing everything from scratch!

> Although Sedge supports several clients, some settings may not work if you use clients versions different from the default ones. Please let us know about any issues you encounter!

### **Disclaimer**

While Sedge assists in installing the validator, it is not designed to register or maintain it. Users are solely responsible for ensuring that they monitor and maintain the validator as required, so that they do not incur penalties and/or financial losses. This includes promptly updating the tool to ensure the latest stable releases of clients are used.

Users acknowledge that no warranty is being made of a successful installation. Sedge is a tool and it ultimately depends on you to use it correctly and follow all the best practice guidance, as found in this README and documentation.

## Lido integration

Sedge is integrated with [Lido CSM](https://operatorportal.lido.fi/modules/community-staking-module), with 
permissionless entry, allowing any node operator — and especially community stakers, from solo stakers, to groups of 
friends, to amateur operators — to operate validators by providing an ETH-based bond (security collateral).

Sedge supports the Lido CSM, allowing users to generate validator keys and set up their full nodes with ease. You can 
read more about it in [our documentation](https://docs.sedge.nethermind.io/docs/quickstart/staking-with-lido)! 

## Charon DV integration
Charon is used by stakers to distribute the responsibility of running Ethereum Validators across a number of different instances and client implementations. Setting up and running a full ethereum node with charon, needs some learning curve and compatibility knowledge, in order for the setup to be fully compliant with the charon configuration requirements for different BN-VC combinations. We want to provide a better and guided user experience for setting up a DV with Charon.

Integrating Charon with Sedge would make it easy for stakers to setup and run a DV with Charon without having to go through each individual client setup docs and their compatibility with DVT.

## Supported networks and clients

### Mainnet

| Execution  | Consensus  | Validator  |
| ---------- |------------|------------|
| Geth       | Lighthouse | Lighthouse |
| Nethermind | Lodestar   | Lodestar   |
| Erigon     | Prysm      | Prysm      |
| Besu       | Teku       | Teku       |
|            | Grandine   |            |
|            | Nimbus     | Nimbus     |

### Sepolia

| Execution  | Consensus  | Validator  |
| ---------- | ---------- | ---------- |
| Geth       | Lighthouse | Lighthouse |
| Nethermind | Lodestar   | Lodestar   |
| Erigon     | Prysm      | Prysm      |
| Besu       | Teku       | Teku       |
|            | Grandine   |            |
|            | Nimbus     | Nimbus     |

### Holesky

| Execution  | Consensus  | Validator  |
| ---------- |------------|------------|
| Geth       | Lighthouse | Lighthouse |
| Nethermind | Lodestar   | Lodestar   |
| Erigon     | Teku       | Teku       |
| Besu       | Prysm      | Prysm      |
|            | Grandine   | 
|            | Nimbus     | Nimbus     |

### Gnosis

| Execution   | Consensus  | Validator  |
|-------------| ---------- | ---------- |
| Nethermind  | Lighthouse | Lighthouse |
| Erigon      | Lodestar   | Lodestar   |
|             | Teku       | Teku       |
|            | Nimbus     | Nimbus     |

### Chiado (Gnosis testnet)

| Execution     | Consensus  | Validator  |
|---------------| ---------- | ---------- |
| Nethermind    | Lighthouse | Lighthouse |
| Erigon (soon) | Lodestar   | Lodestar   |
|               | Teku       | Teku       |
|               | Nimbus     | Nimbus     |

### CL clients with Mev-Boost

| Client     | Mev-Boost  | Networks                  |
|------------|------------|---------------------------|
| Lighthouse | yes        | Mainnet, Sepolia, Holesky |
| Lodestar   | yes        | Mainnet, Sepolia, Holesky |
| Prysm      | yes        | Mainnet, Sepolia, Holesky |
| Teku       | yes        | Mainnet, Sepolia, Holesky |
| Nimbus     | yes        | Mainnet, Sepolia, Holesky |

## Supported Linux flavours for dependency installation

| OS             | Versions                |
| -------------- | ----------------------- |
| Ubuntu         | 22.04, 20.04 |
| Debian         | 11,10,9,8               |
| Fedora         | 35,34                   |
| CentOS         | 8                       |
| Arch           | -                       |
| Amazon Linux 2 | -                       |
| Alpine         | 3.15,3.14,3.14.3        |

## ✅ Roadmap

The following roadmap covers the main features and ideas we want to implement but only covers some of what we are planning for this tool. Stay in touch if you are interested.

### Version 0.1

- [x] Generate `docker-compose` scripts and `.env` files for selected clients with a cli tool
- [x] Generate keystore folder with the cli
- [x] Test coverage (unit tests)
- [x] Integrate Kiln network
- [x] Integrate MEV-Boost as an option
- [x] Integrate Ropsten network

### Version 0.2

- [x] Integrate Goerli/Prater network
- [x] Integrate Sepolia network
- [x] Documentation with examples

### Version 0.3

- [x] Integrate Gnosis network
- [x] Prepare for the Merge

### Version 0.4

- [x] Create and handle keystores on our own instead of using staking-deposit-cli
- [x] Improve validator testing
- [x] Bug fixes
- [x] Deprecate Kiln, Ropsten, Denver networks
- [x] Improve support for chiado network (Gnosis testnet)

### Version 0.5

- [x] Support for Gnosis Merge
- [x] Bug fixes

### Version 0.6 

- [x] Besu and Erigon support
- [x] Windows support
- [x] Bug fixes

### Version 1.0 

- [x] Full Ethereum PoS support with MEV-Boost
- [x] Set up and run only one node (execution/consensus/validator)
- [x] Keystore generation
- [x] Slashing protection
  
### Version 1.X (Current)

- [x] Support Erigon on Gnosis
- [x] Support for Lido CSM
- [x] Support for Nimbus client as Consensus and Validator
- [x] Include monitoring tool for alerting, tracking validator balance, and tracking sync progress and status of nodes
- [x] Support for Optimism and Base.
- [ ] More tests!!!


## 💪 Want to contribute?

Please check our [Contributing Guidelines](https://docs.sedge.nethermind.io/docs/guidelines), Code of Conduct and our issues. In case you want to report or suggest something (any help is welcome), please file an issue first so that the main team is aware and can discuss it.

If you know of any good tricks for validator setup that other people could also use well, please consider adding it to Sedge. Your efforts will be greatly appreciated by the community.

## 🛡️ Security 

If you believe you have found a security vulnerability in our code, please report it to us as described in our [security policy](SECURITY.md).

## ⚠️ License

Sedge is a Nethermind free and open-source software licensed under the [Apache 2.0 License](https://github.com/NethermindEth/sedge/blob/main/LICENSE).