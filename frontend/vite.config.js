import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// During development, /api is proxied to the Go backend so the SPA and API
// share an origin. In production the built SPA is served by the Go backend
// itself (single-port monolith), so the same relative /api paths just work.
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
