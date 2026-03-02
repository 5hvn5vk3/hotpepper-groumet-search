import react from "@vitejs/plugin-react-swc";
import tailwindcss from "@tailwindcss/vite";
import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";
export default defineConfig({
    plugins: [react(), tailwindcss(), tsconfigPaths()],
    server: {
        port: 3000,
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
            },
        },
    },
    test: {
        globals: true,
        environment: "happy-dom",
        setupFiles: "./vitest.setup.ts",
    },
});
