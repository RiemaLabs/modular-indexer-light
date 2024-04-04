# Modular Indexer (Light) [![Join Nubit Discord Community](https://img.shields.io/discord/916984413944967180?logo=discord&style=flat)](https://discord.gg/5sVBzYa4Sg) [![Follow Nubit On X](https://img.shields.io/twitter/follow/nubit_org)](https://twitter.com/Nubit_org)

<img src="assets/logo.svg" width="400px" alt="Nubit Logo" />

## Background
The Modular Indexer, which both includes the [Committee Indexer](https://github.com/RiemaLabs/modular-indexer-committee) and the `Light Indexer`(This Repo), establishes a fully user-verified execution layer for Bitcoin's meta-protocols. By harnessing Bitcoin's immutable and decentralized nature, it offers a Turing-complete execution layer that transcends the constraints of Bitcoin's script language.

For a detailed understanding, refer to our paper: ["Modular Indexer: Fully User-Verified Execution Layer for Meta-protocols on Bitcoin"](https://eprint.iacr.org/2024/408). Stay updated on the latest progress in our [L1F Discourse Group](https://l1f.discourse.group/t/modular-indexer-fully-user-verified-execution-layer-for-meta-protocols-on-bitcoin/598).

## What is `Light Indexer`?
`Light Indexer` plays a crucial role in this ecosystem. It retrieves the state of Bitcoin's meta-protocol from `Committee Indexer` according to the user's demand. While ensuring obtained states are trustworthy, it is efficient enough to be executed on browsers, mobiles, and other light devices.

## Getting Started
Welcome to the setup.

### 1. Requirements
Before we stepped into the installation, ensure your machine is equipped with the minimum requirements: (Such low configuration requirements are absolutely insane!)

| Metric       | Minimum Requirements     |
|--------------|------------------------- |
| **CPU**      | Single Core              |
| **Memory**   | 512 MB                   |
| **Disk**     | 30 GB                    |
| **Bandwidth**| Upload/Download 100 KB/s |

`Light Indexer` is crafted for ease of use and efficient interaction with the Bitcoin blockchain. It's tailored for systems with basic configurations, requiring only a Golang environment (version 1.22.0 or higher) and basic command line knowledge. An internet connection is essential for downloading dependencies and establishing connections to the blockchain. Significantly lighter on resources compared to its counterpart, `Modular Indexer (Committee)`, it seamlessly handles interactions with blockchain data without demanding sophisticated hardware. This synergy between `Light Indexer` and `Modular Indexer (Committee)` ensures an accessible yet robust approach to working with Bitcoin's meta-protocols.
For installation, start by ensuring your system has the latest version of Golang, available from the [Golang download page](https://go.dev/doc/install).

### 2. Install Dependence
`Light Indexer` is built with Golang. You can run your own `Light Indexer` by following the procedure below.
`Go` version `1.22.0` is required for running repository. Please visit [Golang download Page](https://go.dev/doc/install) to get latest Golang installed.
Golang is easy to install all dependence. Fetch all required package by simply running.
```Bash
go mod tidy
```

### 3. Configuration Instructions
Configure `config.json`: To set up `Modular Indexer (Light)`, begin by copying the example configuration file. Customize it to match your specific requirements.
```Bash
cp config.example.json config.json
# Edit config.json according to your needs

```json
{
  "committeeIndexers": {
    "s3": [
      {
        "region": "us-west-2",
        "bucket": "nubit-modular-indexer-brc-20",
        "name": "nubit-official-02"
      }
    ],
    "da": [
      {
        "network": "Pre-Alpha Testnet",
        "namespaceID": "0x00000003",
        "name": "nubit-official-00"
      },
      {
        "network": "Pre-Alpha Testnet",
        "namespaceID": "0x00000005",
        "name": "nubit-official-01"
      }
    ]
  },
  "bitcoinRPC": "",
  "metaProtocol": "brc-20",
  "minimalCheckpoint": 2,
  "version": "v0.1.0-rc.1"
}
```

#### Detailed Configuration Instructions:
After copying the `config.example.json` to create your `config.json`, you'll need to provide detailed information for both Committee Indexer settings and Bitcoin RPC settings. Here's a breakdown of what each section in the configuration file means:

##### a. Setting Up `committeeIndexers`:
- **S3 Indexers**:
  - `region`: The AWS region where your S3 bucket is located.
  - `bucket`: The name of the S3 bucket used by the indexer.
  - `name`: A unique name for your indexer instance.
- **DA Indexers**:
  - `network`: Specify the network for DA Layer (e.g., 'Pre-Alpha Testnet').
  - `namespaceID`: The namespace ID used in the DA Layer.
  - `name`: A unique name for your indexer instance.

##### b. `bitcoinRPC` Configuration:
- `bitcoinRPC`: Enter the URL of your Bitcoin RPC server for direct blockchain interactions.

##### c. Additional Configurations:
- `metaProtocol`: Define the meta-protocol used (e.g., 'brc-20').
- `minimalCheckpoint`: Specify the minimum number of checkpoints required for validation.
- `version`: Indicate the version of your Modular Indexer (Light) setup.

### 4. Running the Program
Build the Application: Compile the Modular Indexer (Light) source code.
```Bash
go build
```

Run Modular Indexer (Light): Start the application. You can also include additional command flags as needed.
```Bash
./modular-indexer-light
```
_Please note: When initiating the Light Indexer, the system will automatically generate a private key and save it in the 'private' file. Ensure that you securely store this private key._


## Basic Usage
After successfully launching `Light Indexer`, you have several functionalities at your disposal for interacting with the Bitcoin blockchain. These capabilities can be accessed through the `Indexer Dashboard` or [direct API calls](https://app.gitbook.com/o/CpG1oV8XXLDnYbUdhhqM/s/RvfNFdIQAghhQdWUByGF/developer-guides/introduction):

#### Querying Wallet Balances
`Light Indexer` is capable of fetching the balance of specified Bitcoin wallets, a crucial feature for tracking transactions and managing wallet assets.

#### Accessing and Verifying Checkpoint Data
Working in conjunction with `Modular Indexer (Committee)`, this tool can access and verify checkpoint data for various meta-protocols, ensuring data accuracy and integrity.

`Light Indexer` is optimized for efficiency, requiring minimal resources. This design provides a user-friendly approach for those needing to interact with Bitcoin's meta-protocols without the complexities of in-depth data processing.

<!-- ## Service API -->

## Useful Links
:spider_web: <https://www.nubit.org>  
:octocat: <https://github.com/Wechaty/wechaty>  
:beetle: <https://github.com/RiemaLabs/modular-indexer-light/issues>  
:book: <https://docs.nubit.org/developer-guides/introduction>  

## FAQ
- **How does the interaction between `Modular Indexer (Light)` and `Modular Indexer (Committee)` work?**
    - Modular Indexer (Light) relies on checkpoints published by Modular Indexer (Committee) to verify the integrity and correctness of data. Light indexers fetch these checkpoints from the DA Layer or S3, enabling accurate and up-to-date blockchain data interaction with minimal computation.

- **How does the `Modular Indexer (Light)` benefit users in the Modular Indexer ecosystem?**
    - `Modular Indexer (Light)` stands out for its ability to efficiently process and validate blockchain data. It’s tailored for users requiring quick and reliable access to Bitcoin's meta-protocols, especially those using devices with limited computational resources.

- **What role does the Indexer Dashboard play in the context of `Modular Indexer (Light)`?**
    - In the context of `Modular Indexer (Light)`, the Indexer Dashboard acts as a centralized platform for accessing and managing data verified by the `Light Indexer`. It simplifies the process of monitoring and analyzing blockchain data, making it accessible even for users with minimal technical expertise.

- **How does `Modular Indexer (Light)` ensure the accuracy and trustworthiness of blockchain data?**
    - `Modular Indexer (Light)` ensures data accuracy and trustworthiness by fetching and verifying checkpoints from the `Modular Indexer (Committee)`. These checkpoints are critical for maintaining the integrity of the data it processes and presents.

- **What makes the `Modular Indexer (Light)` a suitable choice for mobile and browser-based applications?**
    - The `Modular Indexer (Light)` is optimized for minimal resource usage and efficient operation, making it an ideal choice for mobile and browser-based applications where computational resources are limited.

<!-- ## License -->