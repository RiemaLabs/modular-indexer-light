# modular-indexer-light-sdk

Modular Indexer (Light) JavaScript SDK.

## Development

In the toplevel project directory, run this `make` command to build the WebAssembly module:

```bash
make GOOS=js GOARCH=wasm packages/modular-indexer-light-sdk/modular-indexer-light.wasm
 ```

And then go back here to start the development server:

```bash
cd packages/modular-indexer-light-sdk/
npm run start
```

Before you publish a new npm package, bump the version first, and then do the dry-run:

```bash
npm pack
```
