import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { fileURLToPath, URL } from 'url'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  optimizeDeps: {
    exclude: ['lucide-react'],
  },
  resolve: {
    alias: {
      // Ini adalah cara ESM yang benar untuk menunjuk ke folder 'src'
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
});
