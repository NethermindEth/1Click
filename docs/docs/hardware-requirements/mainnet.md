---
sidebar_position: 1
id: mainnet-requirements
---

# Mainnet

Some of the most essential requirements for Mainnet can be found in the [Ethereum Launchpad Checklist](https://launchpad.ethereum.org/en/checklist):

### Storage

As you need to keep both execution and consensus layers (consensus and validator client) databases in a single machine, it is recommended to have a solid state storage disk with ~1TB for Mainnet. 

> In our experience 300-500 GB for the Testnets can be a good minimum range. The older the Testnet, the more space you will need.

### CPU and RAM

- Check with client documentation to ensure the hardware you want to use is sufficient and supported.
- Resource usage can vary significantly between clients. Research the different clients if you're working with resource constraints.
- For Grandine consensus client, recommended specifications are:
  - CPU: 4+ cores (8+ cores recommended for optimal performance)
  - RAM: 16GB+ RAM (32GB recommended for optimal performance)
  - Storage: SSD with at least 1TB free space
  - Storage optimization: Grandine supports state and payload pruning to reduce disk usage

### Internet

- Ideally your internet connection should be reliable and as close to 24/7 as possible without interruption.
- Ensure your bandwidth can't be throttled and isn't capped so your node stays in sync and will be ready to validate when called.
- You need enough upload bandwidth too. As of May 2022 this is ~1.2-1.3 GB download and ~0.9-1 GB upload per hour, and is likely to increase.

### Client-Specific Monitoring

Each client exposes metrics differently. Here are the default ports and endpoints:

- Grandine: Metrics available on port 9090 with detailed metrics enabled by default
  - Basic metrics: /metrics
  - Detailed metrics: enabled with --metrics-detailed flag
