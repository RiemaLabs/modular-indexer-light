# Modular Indexer (Light) [![Join Nubit Discord Community](https://img.shields.io/discord/916984413944967180?logo=discord&style=flat)](https://discord.gg/5sVBzYa4Sg) [![Follow Nubit On X](https://img.shields.io/twitter/follow/nubit_org)](https://twitter.com/Nubit_org)

<img src="assets/logo.svg" width="400px" alt="Nubit Logo" />

## Background
The Modular Indexer, which both includes the [Committee Indexer](https://github.com/RiemaLabs/modular-indexer-committee) and the [Light Indexer](#What-is-Light-Indexer?), establishes a fully user-verified execution layer for Bitcoin's meta-protocols. By harnessing Bitcoin's immutable and decentralized nature, it offers a Turing-complete execution layer that transcends the constraints of Bitcoin's script language.

For a detailed understanding, refer to our paper: ["Modular Indexer: Fully User-Verified Execution Layer for Meta-protocols on Bitcoin"](https://eprint.iacr.org/2024/408). Stay updated on the latest progress in our [L1F Discourse Group](https://l1f.discourse.group/t/modular-indexer-fully-user-verified-execution-layer-for-meta-protocols-on-bitcoin/598).

## What is Light Indexer?
`Light Indexer` plays a crucial role in this ecosystem. It retrieves the state of Bitcoin's meta-protocol from `Committee Indexer` according to the user's demand. While ensuring obtained states are trustworthy, it is efficient enough to be executed on browsers, mobiles, and other light devices.

## Getting Started

### 1. Requirements
Before we stepped into the installation, ensure your machine is equipped with the minimum requirements: (Such low configuration requirements are absolutely insane!)

| Metric       | Minimum Requirements     |
|--------------|------------------------- |
| **CPU**      | Single Core              |
| **Memory**   | 512 MB                   |
| **Disk**     | 30 GB                    |
| **Bandwidth**| Upload/Download 100 KB/s |

### 2. Install Dependence
`Light Indexer` is built with Golang. You can run your own `Light Indexer` by following the procedure below.
`Go` version `1.22.0` is required for running repository. Please visit [Golang download Page](https://go.dev/doc/install) to get latest Golang installed. Golang is easy to install all dependence. Fetch all required package by simply running.
```Bash
go mod tidy
```

### 3. Configuration Instructions
To set up `Light Indexer`, begin by copying the example configuration file. Customize it to match your specific requirements.
```Bash
cp ./config/config.example.json config.json
# Add your own bitcoinRPC!
```
[config.example.json](./config/config.example.json) is the initial configuration for the addresses of the official Nubit committee indexers.

#### Detailed Configuration Instructions:
You may add committee indexers settings and Bitcoin RPC settings. Here's a breakdown of what each section in the configuration file means:

##### a. `committeeIndexers` Configuration:
- **Committee Indexers accessible via DA**:
  - `network`: Specify the network for Nubit DA Layer (current, 'Pre-Alpha Testnet').
  - `namespaceID`: The namespace ID used in the DA Layer.
  - `name`: A unique name for your indexer instance.

- **Committee Indexers accessible by S3**:
  - `region`: The AWS region where your S3 bucket is located.
  - `bucket`: The name of the S3 bucket used by the indexer.
  - `name`: A unique name for your indexer instance.

##### b. `bitcoinRPC` Configuration:
- `bitcoinRPC`: Enter the URL of your Bitcoin RPC server for direct blockchain interactions.

##### c. Additional Configurations:
- `metaProtocol`: Define the meta-protocol used (default, 'brc-20').
- `minimalCheckpoint`: Specify the minimum number of checkpoints required for validation.
- `version`: Indicate the version of your `Light Indexer` setup.

### 4. Running the Program
Build and start the Application:
```Bash
go build
./modular-indexer-light
```
_Please note: When initiating the Light Indexer, the system will automatically generate a private key and save it in the 'private' file. Ensure that you securely store this private key._


## Basic Functionality
After successfully launching `Light Indexer`, you have several functionalities at your disposal for interacting with the Bitcoin blockchain. These capabilities can be accessed through the `Indexer Dashboard` or [direct API calls](https://app.gitbook.com/o/CpG1oV8XXLDnYbUdhhqM/s/RvfNFdIQAghhQdWUByGF/developer-guides/introduction):

#### Querying Wallet Balances
`Light Indexer` is capable of fetching the balance of specified Bitcoin wallets, a crucial feature for tracking transactions and managing wallet assets.

#### Accessing and Verifying Checkpoint Data
Working in conjunction with `Modular Indexer (Committee)`, this tool can access and verify checkpoint data for various meta-protocols, ensuring data accuracy and integrity.

`Light Indexer` is optimized for efficiency, requiring minimal resources. This design provides a user-friendly approach for those needing to interact with Bitcoin's meta-protocols without the complexities of in-depth data processing.

<!-- ## Service API -->

## Useful Links
:spider_web: <https://www.nubit.org>
:beetle: <https://github.com/RiemaLabs/modular-indexer-light/issues>
:book: <https://docs.nubit.org/developer-guides/introduction>

## FAQ
- **How does the interaction between `Light Indexer` and `Modular Indexer (Committee)` work?**
    - `Light Indexer` relies on checkpoints published by Modular Indexer (Committee) to verify the integrity and correctness of data. Light indexers fetch these checkpoints from the DA Layer or S3, enabling accurate and up-to-date blockchain data interaction with minimal computation.

- **How does the `Light Indexer` benefit users in the Modular Indexer ecosystem?**
    - `Light Indexer` stands out for its ability to efficiently process and validate blockchain data. It’s tailored for users requiring quick and reliable access to Bitcoin's meta-protocols, especially those using devices with limited computational resources.

- **What role does the Indexer Dashboard play in the context of `Light Indexer`?**
    - In the context of `Light Indexer`, the Indexer Dashboard acts as a centralized platform for accessing and managing data verified by the `Light Indexer`. It simplifies the process of monitoring and analyzing blockchain data, making it accessible even for users with minimal technical expertise.

- **How does `Light Indexer` ensure the accuracy and trustworthiness of blockchain data?**
    - `Light Indexer` ensures data accuracy and trustworthiness by fetching and verifying checkpoints from the `Modular Indexer (Committee)`. These checkpoints are critical for maintaining the integrity of the data it processes and presents.

- **What makes the `Light Indexer` a suitable choice for mobile and browser-based applications?**
    - The `Light Indexer` is optimized for minimal resource usage and efficient operation, making it an ideal choice for mobile and browser-based applications where computational resources are limited.
