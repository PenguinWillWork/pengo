import { defineConfig } from "vite";

// wails serves the app from its own port and proxies this dev server, so the
// hmr socket has to be pinned — otherwise the webview never connects and
// nothing reloads
export default defineConfig({
  server: {
    port: 5173,
    strictPort: true,
    hmr: {
      protocol: "ws",
      host: "localhost",
      port: 5173,
    },
  },
});
