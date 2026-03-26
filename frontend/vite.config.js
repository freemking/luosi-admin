import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => ({
  plugins: [vue()],
  server: {
    port: 5173
  },
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true
      }
    }
  },
  optimizeDeps: {
    disabled: false,
    force: true
  },
  // Make YAML config files available as raw imports
  assetsInclude: ['**/*.yaml']
}))
