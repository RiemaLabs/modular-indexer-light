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
To run `Modular Indexer (Light)`, ensure your system meets the following requirements:
- Golang environment, version 1.22.0 or higher
- Basic knowledge of terminal/command line usage
- Internet connection for downloading dependencies and interfacing with the Bitcoin blockchain

Modular Indexer is built with Golang. You can run your own modular Indexer by following the procedure below. `Go` version `1.22.0` is required for running repository. Please visit [Golang download Page](https://go.dev/doc/install) to get latest Golang installed.

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

5. **Basic Usage**
After successfully running Modular Indexer (Light), you can start interacting with the Bitcoin blockchain. Test basic functionalities by executing relevant commands or accessing provided APIs.

