import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// During development, /api is proxied to the Go backend so the SPA and API
// share an origin. In production the SPA is served by Nginx which proxies /api.
export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    port: 5173,
    proxy: {
      '/api': {
        target: process.env.VITE_API_TARGET || 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
