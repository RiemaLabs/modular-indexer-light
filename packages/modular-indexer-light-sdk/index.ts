declare const Go: any;

/**
 * Load the WebAssembly module to make all the APIs available.
 */
export async function lightLoadWasm() {
    const go = new Go();
    const response = await fetch('./modular-indexer-light.wasm');
    const buffer = await response.arrayBuffer();
    const result = await WebAssembly.instantiate(buffer, go.importObject);
    go.run(result.instance);
}

interface Config {
    verification: {
        bitcoinRPC: string,
        metaProtocol: string,
        minimalCheckpoint: number
    },
    committeeIndexers: {
        s3: {
            region: string,
            bucket: string,
            name: string
        }[],
        da: {
            network: string,
            namespaceID: string,
            name: string
        }[]
    }
}

/**
 * Set the configuration for the light indexer.
 *
 * @param c configuration
 */
export declare function lightSetConfig(c: Config): void;

/**
 * Initialize and start the syncing and verification loop.
 */
export declare function lightInitialize(): void;

type Status = "verifying" | "verified" | "unverified";

/**
 * Get the status of the light indexer.
 */
export declare function lightStatus(): Status;

/**
 * Get current BTC block height.
 */
export declare function lightGetBlockHeight(): Promise<number>;

interface BalanceOfPkScript {
    error?: string,
    result?: {
        availableBalance: string,
        overallBalance: string
    },
    proof?: string
}

/**
 * Get balance via PkScript.
 *
 * @param tick token
 * @param pkscript public key script
 */
export declare function lightGetBalanceOfPkScript(tick: string, pkscript: string): Promise<BalanceOfPkScript>;

interface BalanceOfWallet {
    error?: string,
    result?: {
        availableBalance: string,
        overallBalance: string,
        pkscript: string
    },
    proof?: string
}

/**
 * Get balance via wallet.
 *
 * @param tick token
 * @param wallet wallet address
 */
export declare function lightGetBalanceOfWallet(tick: string, wallet: string): Promise<BalanceOfWallet>;

interface Checkpoint {
    commitment: string,
    hash: string,
    height: string,
    metaProtocol: string,
    name: string,
    url: string,
    version: string
}

/**
 * Get current checkpoints from all the committee indexers.
 */
export declare function lightGetCurrentCheckpoints(): Promise<Checkpoint[]>;

/**
 * Get the previous consensus-reached checkpoint.
 */
export declare function lightGetLastCheckpoint(): Promise<Checkpoint>;
