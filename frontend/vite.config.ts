import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
  plugins: [
    react(),
    tailwindcss()
  ],
  resolve: {
    alias: {
      '@': '/src'
    }
  },
  server: {
    proxy: {
      // Shopping Service (products, cart)
      '/api/products': {
        target: env.VITE_SHOPPING_SERVICE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      '/api/cart': {
        target: env.VITE_SHOPPING_SERVICE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      // Auth Service (login, register, logout)
      '/api/login': {
        target: env.VITE_AUTH_SERVICE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      '/api/register': {
        target: env.VITE_AUTH_SERVICE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      '/api/logout': {
        target: env.VITE_AUTH_SERVICE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      // Checkout Service
      '/api/checkout': {
        target: env.VITE_CHECKOUT_SERVICE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      // Payment Service
      '/api/payments': {
        target: env.VITE_PAYMENT_SERVICE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
  }
})
