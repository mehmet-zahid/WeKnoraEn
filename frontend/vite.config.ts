import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'

export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    port: 5173,
    host: true,
    // Proxy configuration for development environment
    // In Docker, use service name 'app', otherwise use localhost
    proxy: {
      '/api': {
        target: process.env.VITE_IS_DOCKER === 'true' ? 'http://app:8080' : 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
