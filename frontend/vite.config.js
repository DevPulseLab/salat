import { defineConfig } from 'vite'
import path from 'path'
import vue from '@vitejs/plugin-vue'
import mkcert from 'vite-plugin-mkcert'

export default defineConfig({
  base: '/public/',
  build: {
    outDir: '../public'
  },
  server: {
    https: true
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  plugins: [
    vue(),
    mkcert()
  ],
})
