import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  root: 'src',
  build: {
    outDir: resolve(__dirname, 'dist'), // 👈 明确指定为 frontend/dist,
    emptyOutDir: true,
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'src/index.html')
      }
    }
  },
  server: {
    port: 3000,
    strictPort: true
  }
})