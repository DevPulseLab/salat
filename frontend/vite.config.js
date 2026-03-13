import { defineConfig, loadEnv } from 'vite'
import path from 'path'
import vue from '@vitejs/plugin-vue'
import mkcert from 'vite-plugin-mkcert'

export default defineConfig(({ mode }) => {

  const env = loadEnv(mode, process.cwd());

  const useHttps = env.VITE_HTTPS !== 'false';

  return {
    base: '/',
    build: {
      outDir: '../public'
    },
    server: {
      https: useHttps
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
  }
})
