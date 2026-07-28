import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  server: {
    port: 3002,
    host: true,
    allowedHosts: ['.monkeycode-ai.online', '3002-cd407fcd09c4d17d.monkeycode-ai.online'],
  },
  vite: {
    plugins: [tailwindcss()],
    server: {
      allowedHosts: ['.monkeycode-ai.online', '3002-cd407fcd09c4d17d.monkeycode-ai.online'],
    },
  },
  devToolbar: { enabled: false },
});
