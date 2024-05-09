import "./wasm_exec.js";
import init from "./modular-indexer-light.wasm?init";
import {expose} from "comlink";

declare const Go: any;

export interface Config {
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

export type Status = "verifying" | "verified" | "unverified";

export interface BalanceOfPkScript {
    error?: string,
    result?: {
        availableBalance: string,
        overallBalance: string
    },
    proof?: string
}

export interface BalanceOfWallet {
    error?: string,
    result?: {
        availableBalance: string,
        overallBalance: string,
        pkscript: string
    },
    proof?: string
}

export interface Checkpoint {
    commitment: string,
    hash: string,
    height: string,
    metaProtocol: string,
    name: string,
    url: string,
    version: string
}

/**
 * Light indexer SDK.
 */
export class SDK {
    /**
     * Load the WebAssembly module and run the light indexer.
     *
     * @param c configuration
     */
    async run(c: Config) {
        const go = new Go();
        go.run(await init(go.importObject));

        lightSetConfig(c);
        lightInitialize();
        lightWarmup();
    }

    /**
     * Get the status of the light indexer.
     */
    async getStatus() {
        return await lightStatus();
    }

    /**
     * Get current BTC block height.
     */
    async getBlockHeight() {
        return await lightGetBlockHeight();
    }

    /**
     * Get balance via PkScript.
     *
     * @param tick token
     * @param pkscript public key script
     */
    async getBalanceOfPkScript(tick: string, pkscript: string) {
        return await lightGetBalanceOfPkScript(tick, pkscript);
    }

    /**
     * Get balance via wallet.
     *
     * @param tick token
     * @param wallet wallet address
     */
    async getBalanceOfWallet(tick: string, wallet: string) {
        return await lightGetBalanceOfWallet(tick, wallet);
    }

    /**
     * Get current checkpoints from all the committee indexers.
     */
    async getCurrentCheckpoints() {
        return await lightGetCurrentCheckpoints();
    }

    /**
     * Get the previous consensus-reached checkpoint.
     */
    async getLastCheckpoint() {
        return await lightGetLastCheckpoint();
    }
}

declare function lightSetConfig(c: Config): void;

declare function lightInitialize(): void;

declare function lightWarmup(): void;

declare function lightStatus(): Promise<Status>;

declare function lightGetBlockHeight(): Promise<number>;

declare function lightGetBalanceOfPkScript(tick: string, pkscript: string): Promise<BalanceOfPkScript>;

declare function lightGetBalanceOfWallet(tick: string, wallet: string): Promise<BalanceOfWallet>;

declare function lightGetCurrentCheckpoints(): Promise<Checkpoint[]>;

declare function lightGetLastCheckpoint(): Promise<Checkpoint>;

expose(SDK);
