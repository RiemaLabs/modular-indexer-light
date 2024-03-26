# Modular Indexer (Light) [![Join Nubit Discord Community](https://img.shields.io/discord/916984413944967180?logo=discord&style=flat)](https://discord.gg/5sVBzYa4Sg) [![Follow Nubit On X](https://img.shields.io/twitter/follow/nubit_org)](https://twitter.com/Nubit_org)

<img src="assets/logo.svg" width="400px" alt="Nubit Logo" />

***Warning!*** *This release is specifically for the Pre-alpha Testnet and may include changes that are not backward compatible in the future.*

## Background
The Modular Indexer, which both includes the [Modular Indexer (Committee)](https://github.com/RiemaLabs/modular-indexer-committee) and the `Modular Indexer (Light)`, establishes a fully user-verified execution layer for Bitcoin's meta-protocols. By harnessing Bitcoin's immutable and decentralized nature, it offers a Turing-complete execution layer that transcends the constraints of Bitcoin's script language.

The `Modular Indexer (Light)` plays a crucial role in this ecosystem. It acts as a nimble counterpart to the `Modular Indexer (Committee)`, focusing on real-time interaction with the Bitcoin blockchain and validation of data integrity using state-of-the-art cryptographic methods. Even in the presence of adversarial conditions, `Modular Indexer (Light)` ensures reliable and secure connections between Bitcoin and intricate applications such as BRC-20, thus driving the Bitcoin ecosystem's advancement.

For a detailed understanding, refer to our paper: ["Modular Indexer: Fully User-Verified Execution Layer for Meta-protocols on Bitcoin"](https://eprint.iacr.org/2024/408). Stay updated on the latest progress in our [L1F Discourse Group](https://l1f.discourse.group/t/modular-indexer-fully-user-verified-execution-layer-for-meta-protocols-on-bitcoin/598).

## What is Modular Indexer (Light)?
`Modular Indexer (Light)` is an integral, lightweight counterpart within the Modular Indexer ecosystem, primarily designed to facilitate real-time verification and processing of Bitcoin's meta-protocol data. This component interfaces directly with the Bitcoin blockchain, ensuring up-to-date synchronization and providing users with reliable, timely information. It stands out for its streamlined nature, enabling ease of use and accessibility, especially for those with limited computational resources. In essence, `Modular Indexer (Light)` acts as a bridge, simplifying complex protocol interactions while upholding the decentralized trust of the Bitcoin ecosystem. This makes it an ideal solution for users seeking efficient, reliable access to advanced Bitcoin functionalities without the complexities of intensive data processing.

## Getting Started
Welcome to the setup guide for `Modular Indexer (Light)`. This section will guide you through the necessary steps to get your Modular Indexer (Light) up and running.

### System Requirements
`Modular Indexer (Light)` is crafted for ease of use and efficient interaction with the Bitcoin blockchain. It's tailored for systems with basic configurations, requiring only a Golang environment (version 1.22.0 or higher) and basic command line knowledge. An internet connection is essential for downloading dependencies and establishing connections to the blockchain. Significantly lighter on resources compared to its counterpart, `Modular Indexer (Committee)`, it seamlessly handles interactions with blockchain data without demanding sophisticated hardware. This synergy between `Modular Indexer (Light)` and `Modular Indexer (Committee)` ensures an accessible yet robust approach to working with Bitcoin's meta-protocols.

For installation, start by ensuring your system has the latest version of Golang, available from the [Golang download page](https://go.dev/doc/install).

### Installation Steps
1. **Clone the Repository**
```Bash
git clone https://github.com/RiemaLabs/modular-indexer-light.git
cd modular-indexer-light
```

2. **Install Dependencies**
Once you have cloned the repository, install all necessary dependencies.
```Bash
go mod tidy
```

3. **Configuration Instructions**
Prepare config.json: Copy the example configuration file and tailor it according to your setup.
```Bash
cp config.example.json config.json
# Edit config.json with your specific settings
```

Edit Configuration: Ensure the bitCoinRpc and committeeIndexerApi sections in config.json are set with the correct URLs and credentials, particularly the bitCoinRpc and committeeIndexerApi sections.
```json
{
  "bitCoinRpc": {
    "host": "YOUR_BITCOIN_RPC_HOST",
    "user": "YOUR_RPC_USER",
    "password": "YOUR_RPC_PASSWORD"
  },
  "committeeIndexerApi": {
    "name": "YOUR_INDEXER_NAME",
    "url": "YOUR_INDEXER_API_URL"
  }
}
```

### `bitCoinRpc` Configuration:
- `host`: The URL of your Bitcoin RPC server.
- `user`: The username for accessing the Bitcoin RPC server.
- `password`: The password associated with the specified user.

### `committeeIndexerApi` Configuration:
- `name`: A unique name for your indexer instance.
- `url`: The endpoint URL of the Modular Indexer (Committee) or any equivalent service.

4. **Running the Program**
Build the Application: Compile the Modular Indexer (Light) source code.
```Bash
go build
```

Run Modular Indexer (Light): Start the application. You can also include additional command flags as needed.
```Bash
./modular-indexer-light
```

### 5. Basic Usage
After successfully launching `Modular Indexer (Light)`, you have several functionalities at your disposal for interacting with the Bitcoin blockchain. These capabilities can be accessed through the `Indexer Dashboard` or [direct API calls](https://app.gitbook.com/o/CpG1oV8XXLDnYbUdhhqM/s/RvfNFdIQAghhQdWUByGF/developer-guides/introduction):

#### Interacting with Bitcoin RPC
The application can interact with the Bitcoin network via RPC calls, allowing operations like fetching the latest block height or retrieving block hashes.

#### Querying Wallet Balances
`Modular Indexer (Light)` is capable of fetching the balance of specified Bitcoin wallets, a crucial feature for tracking transactions and managing wallet assets.

#### Accessing and Verifying Checkpoint Data
Working in conjunction with `Modular Indexer (Committee)`, this tool can access and verify checkpoint data for various meta-protocols, ensuring data accuracy and integrity.

#### Real-time Data Processing
The application efficiently processes and analyzes blockchain data in real-time, enabling up-to-date blockchain interactions and analytics.

As `Modular Indexer (Light)` is designed for efficiency and minimal resource usage, it provides a streamlined and accessible approach for users requiring interaction with Bitcoin's meta-protocols without intensive data processing needs.

<!-- ## Service API -->

## Useful Links
:spider_web: <https://www.nubit.org>  
:octocat: <https://github.com/Wechaty/wechaty>  
:beetle: <https://github.com/RiemaLabs/modular-indexer-light/issues>  
:book: <https://docs.nubit.org/developer-guides/introduction>  

## FAQ
- **Is there a consensus mechanism among modular-indexer-committees?**
    - No, within the modular-indexer-committee, only one honest indexer needs to be available in the network to satisfy the 1-of-N trust assumption. This allows the light indexer to detect checkpoint inconsistencies and proceed with the verification process.
- **How is the set of modular-indexer-committees determined?**
    - Modular-indexer-committees must publish checkpoints to the DA Layer for access by other participants. Users can maintain their list of trusted committees. Since the light indexer verifies checkpoint correctness, malicious actors can be removed from the set upon detection; judgment of malicious behavior relies on a challenge-proof mechanism, not a 51% vote.
- **Why do users need to verify data through checkpoints instead of a simple majority of the indexer network?**
    - Relying on a simple majority could lead to Sybil attacks, as joining the indexer network is permissionless. Without a staking model or proof of work, attackers can gain a majority at a low cost. Therefore, verifying data through checkpoints prevents these low-cost attacks.
- **Why are there no attacks like double-spending in the Modular Indexer architecture?**
    - Bitcoin itself ensures transaction ordering and finality for meta-protocols (current: BRC-20). By ensuring the correctness of state transition rules and execution, indexers avoid double-spending attacks, including handling block reorganizations.
- **Why upload checkpoints to the DA Layer instead of a centralized server or directly to Bitcoin?**
    - Using a centralized server risks downtime and withholding of honest checkpoints, breaking the 1-of-N trust assumption. Uploading checkpoints to Bitcoin is costly due to transaction fees and the high data throughput demand, which Bitcoin's throughput cannot support.
- **What kind of ecosystem support has this proposal received?**
    - This proposal by Nubit aims to support and build the Bitcoin ecosystem. We've exchanged ideas with ecosystem partners and aim to promote the progress of the modular indexer architecture jointly.
- **How does the interaction between Modular Indexer (Light) and Modular Indexer (Committee) work?**
    - Modular Indexer (Light) relies on checkpoints published by Modular Indexer (Committee) to verify the integrity and correctness of data. Light indexers fetch these checkpoints from the DA Layer or S3 to confirm the accuracy of blockchain data. This interaction ensures reliable data without the need for heavy computation.
- **What role does the Indexer Dashboard play in the Modular Indexer ecosystem?**
    - The Indexer Dashboard is an interface that interacts with both Modular Indexer (Light) and (Committee), providing a user-friendly way to access and visualize blockchain data and checkpoints. It enhances user interaction, making it easy to monitor and manage data across the Modular Indexer network.
    
<!-- ## License -->