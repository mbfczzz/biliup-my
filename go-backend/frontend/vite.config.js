import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/v1': {
        target: 'http://localhost:19159',
        changeOrigin: true,
      },
      '/bili': {
        target: 'http://localhost:19159',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
  },
})
