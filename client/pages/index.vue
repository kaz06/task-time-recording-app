<template>
  <v-app>
    <v-main>
      <v-container>
        <v-form @submit.prevent="handleLogin">
          <v-text-field v-model="email" label="Email" required></v-text-field>
          <v-text-field
            v-model="password"
            label="Password"
            type="password"
            required
          ></v-text-field>
          <v-btn type="submit">Login</v-btn>
        </v-form>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { signInWithEmailAndPassword } from 'firebase/auth'
import { useNuxtApp, useRouter } from 'nuxt/app'

const email = ref('')
const password = ref('')
const nuxtApp = useNuxtApp()
const auth = nuxtApp.$auth
const router = useRouter()

const handleLogin = async () => {
  try {
    await signInWithEmailAndPassword(auth, email.value, password.value)
    console.log('Login successful')
    router.push('/task-register')
  } catch (error) {
    console.error('Login error:', error)
  }
}
</script>