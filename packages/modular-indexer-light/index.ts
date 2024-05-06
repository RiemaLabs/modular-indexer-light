declare const Go: any;

interface Config {
    verification: {
        bitcoinRPC: string;
        metaProtocol: string;
        minimalCheckpoint: number;
    };
    committeeIndexers: {
        s3: {
            region: string;
            bucket: string;
            name: string;
        }[];
        da: {
            network: string;
            namespaceID: string;
            name: string;
        }[];
    };
}

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

export declare function lightGetBlockHeight(): Promise<number>;
