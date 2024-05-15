# Modular Indexer (Light) JavaScript SDK [![Join Nubit Discord Community](https://img.shields.io/discord/916984413944967180?logo=discord&style=flat)](https://discord.gg/5sVBzYa4Sg) [![Follow Nubit On X](https://img.shields.io/twitter/follow/nubit_org)](https://twitter.com/Nubit_org)

<img src="../../.github/logo.svg" width="400px" alt="Nubit Logo" />

> 👛 User-verifiability shapes into your wallet!

JavaScript SDK for running [Modular Indexer (Light)] inside your wallets, websites, Chrome extensions, and more!

*See our demos for the [UniSat wallet extension]. More wallet integration on the way!*

[Modular Indexer (Light)]: ../../.github/README.md

[UniSat wallet extension]: https://github.com/unisat-wallet/extension/pull/196

## Installation

> [!IMPORTANT]
> Our JavaScript SDK is **only available on the web (browsers)**, since it introduces web workers for better
> performance, and they are currently not supported on JavaScript runtimes like Node.js.

Install via `npm` or other package managers:

```bash
npm i @nubit/modular-indexer-light-sdk
```

## Examples

The following example shows how to get the verified balance with a given tick and wallet address.

```typescript
import {create} from "@nubit/modular-indexer-light-sdk";

// Create the SDK first.
const sdk = await create();

// Run SDK with configurations.
await sdk.run({
    "committeeIndexers": {
        "s3": [
            {
                "region": "us-west-2",
                "bucket": "nubit-modular-indexer-brc-20",
                "name": "nubit-official-00"
            }
        ],
        "da": []
    },
    "verification": {
        "bitcoinRPC": "https://bitcoin-mainnet-archive.allthatnode.com",
        "metaProtocol": "brc-20",
        "minimalCheckpoint": 1
    },
});

// Get the SDK status, e.g. `verifying`.
console.log(await sdk.getStatus());

// Some code here to wait for the SDK status being `verified`...

// If the SDK status becomes `verified`, get the balance.
console.log(await sdk.getBalanceOfWallet("ordi", "123abc456def"));
```

## How it works?

### Modular Indexer for the greater good

Check out [Modular Indexer (Committee)] for the technical details about our user-verifiable execution layer.

[Modular Indexer (Committee)]: https://github.com/RiemaLabs/modular-indexer-committee

### WebAssembly for seamless integration

Modular Indexer (Light) is a web service and library in Go, and Go officially supports compiling code into WebAssembly.
With this, the possibility of integrating Modular Indexer (Light) everywhere becomes reality.

### Web workers for smoother interactivity

Loading and warming up a WebAssembly module inside the main JavaScript thread would be a disaster, where it could freeze
the UI and annoy our users.

We start a new web worker to tame the beast of it and communicate with the worker via posting messages.

## APIs

### `create(): Promise<SDK>`

Create the SDK instance.

### `SDK.run(c: Config): Promise<void>`

Load the WebAssembly module and get everything started.

### `SDK.getStatus(): Promise<Status>`

Get the SDK status, it could be:

* `"verifying"`: Still verifying the checkpoints from committee indexers
* `"verified"`: All checkpoints are verified and consistent to retrieve user data for good
* `"unverified"`: Checkpoints seem inconsistent and a further reconstruction of checkpoints starts to run

### `SDK.getBlockHeight(): Promise<number>`

Get the current Bitcoin block height.

Throws an error if SDK is still verifying.

### `SDK.getBalanceOfPkScript(tick: string, pkscript: string): Promise<BalanceOfPkScript>`

Get verified balances via PkScript.

Throws an error if SDK is still verifying.

### `SDK.getBalanceOfWallet(tick: string, wallet: string): Promise<BalanceOfWallet>`

Get verified balances via a wallet address.

Throws an error if SDK is still verifying.

### `SDK.getCurrentCheckpoints(): Promise<Checkpoint[]>`

Get current checkpoints from all the committee indexers, for introspection. A checkpoint contains useful information
like commitments to track down malicious ones.

Throws an error if SDK is still verifying.

### `SDK.getLastCheckpoint(): Promise<Checkpoint>`

Get the previous checkpoint that is proven to be consistent. This checkpoint could be used internally to reconstruct a
new and trusted checkpoint when some malicious committee indexers exist.

Throws an error if SDK is still verifying.

## Development

In the toplevel project directory, run this `make` command to build the WebAssembly module:

```bash
make GOOS=js GOARCH=wasm packages/modular-indexer-light-sdk/modular-indexer-light.wasm
 ```

And then go back here to start the development server:

```bash
cd packages/modular-indexer-light-sdk/
npm run dev
```

Before publishing a new package, bump the version first, and then do the dry-run:

```bash
npm run build
npm pack # to check the list of bundled files
npm publish --access=public
```
