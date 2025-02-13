import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    allowedHosts: [ "hoover.tawny-hen.ts.net" ],
    proxy: {
      // string shorthand: http://localhost:5173/foo -> http://localhost:4567/foo
      // '/api': 'http://localhost:9000',
      '/api': 'http://accounts.tawny-hen.ts.net',
    },
  },
})
