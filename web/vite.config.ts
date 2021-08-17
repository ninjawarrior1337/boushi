import {defineConfig} from "vite"
import WindiCSS from "vite-plugin-windicss"
import svelte from "@sveltejs/vite-plugin-svelte"

export default defineConfig({
    plugins: [
        svelte({}),
        WindiCSS()
    ],
    build: {
        outDir: "../data",
        emptyOutDir: true
    }
})