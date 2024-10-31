import { ref } from 'vue'
import { useNuxtApp } from 'nuxt/app'
import {
    signInWithEmailAndPassword,
    signOut,
    getIdToken,
    onAuthStateChanged,
  } from 'firebase/auth'
import type { User } from 'firebase/auth'
export const useAuth = () => {
  const nuxtApp = useNuxtApp()
  const auth = nuxtApp.$auth
  const user = ref<User | null>(null)

  onAuthStateChanged(auth, (currentUser) => {
    user.value = currentUser
  })

  const login = async (email: string, password: string): Promise<void> => {
    await signInWithEmailAndPassword(auth, email, password)
  }

  const logout = async (): Promise<void> => {
    await signOut(auth)
  }

  const getToken = async (): Promise<string | null> => {
    return user.value ? await getIdToken(user.value) : null
  }

  return { user, login, logout, getToken }
}
