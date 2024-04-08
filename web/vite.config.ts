import path from 'path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const pathSrc = path.resolve(__dirname, 'src')
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
    resolve: {
        alias: {
            '@': pathSrc,
        },
    },
  server: {
    proxy: {
      "/api/v1": {
        target: "http://localhost:8542",
        changeOrigin: true,
      }
    }
  }
})
