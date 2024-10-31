import { defineNuxtConfig } from 'nuxt/config'
import wasm from 'vite-plugin-wasm';
export default defineNuxtConfig({
  build: {
    transpile: ['vuetify'],
  },
  css: [
    'vuetify/styles',
  ],
  vite: {
    define: {
      'process.env.DEBUG': false,
    },
    ssr: {
      noExternal: ['vuetify'],
    },
    plugins: [wasm()],
  },
  runtimeConfig: {
    public: {
      baseURL: process.env.BASE_URL || 'http://localhost:8080',
      firebaseApiKey: process.env.FIREBASE_API_KEY,
      firebaseAuthDamain: process.env.FIREBASE_AUTH_DOMAIN,
      firebaseProjectId: process.env.FIREBASE_PROJECT_ID,
      firebaseStorageBucket: process.env.FIREBASE_STORAGE_BUCKET,
      firebaseMessagingSenderId: process.env.FIREBASE_MESSAGING_SENDER_ID,
      firebaseAppId: process.env.FIREBASE_APP_ID,
    },
  },
  typescript: {
    shim: false
  }
})
