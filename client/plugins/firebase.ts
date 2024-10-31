import { defineNuxtPlugin, useRuntimeConfig } from 'nuxt/app'
import { initializeApp, type FirebaseOptions } from 'firebase/app'
import { getAuth } from 'firebase/auth'

export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig()

  const firebaseConfig: FirebaseOptions = {
    apiKey: config.public.firebaseApiKey as string,
    authDomain: config.public.firebaseAuthDomain as string,
    projectId: config.public.firebaseProjectId as string,
    storageBucket: config.public.firebaseStorageBucket as string,
    messagingSenderId: config.public.firebaseMessagingSenderId as string,
    appId: config.public.firebaseAppId as string,
    measurementId: config.public.firebaseMeasurementId as string,
  }

  const app = initializeApp(firebaseConfig)
  const auth = getAuth(app)

  nuxtApp.provide('auth', auth)
})

