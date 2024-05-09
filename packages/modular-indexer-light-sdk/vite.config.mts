import {resolve} from "path"
import {defineConfig} from "vite"
import dts from "vite-plugin-dts"

export default defineConfig({
    build: {
        assetsInlineLimit: 1024 * 1024 * 1024,
        lib: {
            entry: resolve(__dirname, "index.ts"),
            formats: ["es"]
        },
    },
    plugins: [
        dts({
            exclude: "vite.config.mts"
        })
    ],
    worker: {
        format: "es"
    }
})
