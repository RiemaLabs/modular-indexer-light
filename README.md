# Modular Indexer (Light) [![Join Nubit Discord Community](https://img.shields.io/discord/916984413944967180?logo=discord&style=flat)](https://discord.gg/5sVBzYa4Sg) [![Follow Nubit On X](https://img.shields.io/twitter/follow/nubit_org)](https://twitter.com/Nubit_org)

<img src="assets/logo.svg" width="400px" alt="Nubit Logo" />

## Background
The Modular Indexer, which both includes the [Committee Indexer](https://github.com/RiemaLabs/modular-indexer-committee) and the `Light Indexer`(This Repo), establishes a fully user-verified execution layer for Bitcoin's meta-protocols. By harnessing Bitcoin's immutable and decentralized nature, it offers a Turing-complete execution layer that transcends the constraints of Bitcoin's script language.

For a detailed understanding, refer to our paper: ["Modular Indexer: Fully User-Verified Execution Layer for Meta-protocols on Bitcoin"](https://eprint.iacr.org/2024/408). Stay updated on the latest progress in our [L1F Discourse Group](https://l1f.discourse.group/t/modular-indexer-fully-user-verified-execution-layer-for-meta-protocols-on-bitcoin/598).

## What is `Light Indexer`?
`Light Indexer` plays a crucial role in this ecosystem. It retrieves the state of Bitcoin's meta-protocol from `Committee Indexer` according to the user's demand. While ensuring obtained states are trustworthy, it is efficient enough to be executed on browsers, mobiles, and other light devices.

## Getting Started

### 1. Requirements
Before stepping into the installation, ensure your machine is equipped with the minimum requirements: (Too easy to be met!)

| Metric       | Minimum Requirements     |
|--------------|------------------------- |
| **CPU**      | Single Core              |
| **Memory**   | 512 MB                   |
| **Disk**     | 30 GB                    |
| **Bandwidth**| Upload/Download 100 KB/s |

### 2. Install Dependence
Light Indexer is built with Golang. You can run your own one by following the procedure below.
`Go` version 1.22.0 is required for running the repository. Please visit the [Golang download Page](https://go.dev/doc/install) to get the latest Golang installed.

Golang is easy to install all dependence. Fetch all required packages by simply running.
```Bash
go mod tidy
```

### 3. Configuration Instructions
To set up Light Indexer, begin by copying the example configuration file. 
```Bash
cp config.example.json config.json
# Edit config.json according to your needs
```
Then, customize it to match your specific requirements as follows.

### Detailed Configuration Instructions:
After copying the `config.example.json` and creating your `config.json`, more detailed information is required. Here's a brief outline of the necessary variables to be configured:

#### Setting Up `report`:
Set up this field to allow your light indexer to upload checkpoints to the Nubit DA Layer and participate in the Pre-Alpha Testnet! To get gasCoupon, please follow the guideline of [Nubit website](https://points.nubit.org).
- `name`: A unique name for your light indexer instance.
- `network`: Specify the network (current: 'Pre-Alpha Testnet').
- `namespaceID`: Your designated namespace identifier. Leave it to empty to create a namespace following the instruction.
- `gasCoupon`: Customized code for managing transaction fees.
- `timeout`: The timeout to upload a checkpoint to the Nubit DA Layer.

#### Setting Up `committeeIndexers`:
As of now, the Light Indexer cannot automatically detect active Committee Indexers. Therefore, the default Committee Indexers that are recognized are those operated officially by Nubit and they are provided by `config.example.json`.

Still, you could add information provided by committee indexer runners:
- **da**:
  - `network`: Specification of the network for DA Layer (current: 'Pre-Alpha Testnet').
  - `namespaceID`: The namespace ID used by the committee indexer.
  - `name`: The name of the committee indexer.
- **s3**:
  - `region`: The AWS S3 region where the committee indexer's S3 bucket is located.
  - `bucket`: The AWS S3 bucket used by the committee indexer.
  - `name`: The name of the committee indexer.

#### Setting Up `verification`:
Set up this field to change the verification process.
- `bitcoinRPC`: The URL of your Bitcoin (mainnet) RPC server for direct blockchain interactions. You have the option to use a public RPC server such as https://bitcoin-mainnet-archive.allthatnode.com, or you can acquire your own through QuickNode.
- `metaProtocol`: Definition of the meta-protocol used (current: 'brc-20').
- `minimalCheckpoint`: The minimum number of checkpoints to be obtained from committee indexers (the validity threshold).

### 4. Running the Program
Run the commands below, and the light indexer will initiate API services and upload checkpoints to DA:
```Bash
go build
./modular-indexer-light
```

## Basic Usage
Light Indexer is optimized for cost-efficiency. This design provides a user-friendly approach for those needing to interact with Bitcoin's meta-protocols (such as brc-20) without expensive data processing.

After successfully launching Light Indexer, you have several functionalities at your disposal. These capabilities can be accessed through [direct API calls](https://docs.nubit.org/modular-indexer/nubit-light-indexer-apis). The brc-20 balances provided by the light indexer are fully verified and trustworthy.

## Useful Links
- :spider_web: <https://www.nubit.org>
- :beetle: <https://github.com/RiemaLabs/modular-indexer-light/issues>
- :book: <https://docs.nubit.org/developer-guides/introduction>

## FAQ
- **How is the set of committee indexers determined?**
    - Committee indexers must publish checkpoints to the DA Layer for access by other participants. Users can maintain their list of committee indexers. Since a light indexer (standard user) can verify the correctness of checkpoints, attackers can be removed from the committee indexer set upon detection of malicious behavior; the judgment of malicious behavior is not based on a 51% vote but on a challenge-proof mechanism. Even if the vast majority of committee indexers are malicious, if there is one honest committee indexer, the correct checkpoint can be calculated/verified, allowing the service to continue.
- **Why do users need to verify data through checkpoints instead of looking at the simple majority of the indexer network?**
    - This would lead to Sybil attacks: joining the indexer network is permissionless, without a staking model or proof of work, so the economic cost of setting up an indexer attacker cluster is very low, requiring only the cost of server resources. This allows attackers to achieve a simple majority at a low economic cost; even by introducing historical reputation proof, without a slashing mechanism, attackers can still achieve a 51% attack at a very low cost.
- **Why are there no attacks like double-spending in the Modular Indexer architecture?**
    - Bitcoin itself provides transaction ordering and finality for meta-protocols (such as BRC-20). It is only necessary to ensure the correctness of the indexer's state transition rules and execution to avoid double-spending attacks (there might be block reorganizations, but indexers can correctly handle them).
- **Why upload checkpoints to the DA Layer instead of a centralized server or Bitcoin?**
    - For a centralized server, if checkpoints are stored on a centralized network, the service loses availability in the event of downtime, and there is also the situation where the centralized server withholds checkpoints submitted by honest indexers, invalidating the 1-of-N trust assumption.
    - For indexers, checkpoints are frequently updated, time-sensitive data:
        - The state of the Indexer updates with block height and block hash, leading to frequent updates of checkpoints (~10 minutes).
        - The cost of publishing data on Bitcoin in terms of transaction fees is very high.
        - The data throughput demand for hundreds or even thousands of meta-protocol indexers storing checkpoints is huge, and the throughput of Bitcoin cannot support it.
- **What are the mainstream meta-protocols on Bitcoin currently?**
    - The mainstream meta-protocols are all based on the Ordinals protocol, which allows users to store raw data on Bitcoin. BRC-20, Bitmap, SatsNames, etc., are mainstream meta-protocols. More meta-protocols and information can be found [here](https://l1f.discourse.group/latest)
- **How does the interaction between Light Indexer and Committee Indexer work?**
    - Light Indexer relies on checkpoints published by Committee Indexer to verify the integrity and correctness of data. Light Indexers fetch these checkpoints from Nubit DA Layer or S3 while Nubit DA layer gurantees the data availability of checkpoints. This interaction ensures reliable data without the need for heavy computation or indexing of the whole history of Bitcoin blocks.
