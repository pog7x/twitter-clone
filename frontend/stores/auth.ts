import { defineStore } from 'pinia'
import type { User } from '~/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = useCookie('access_token', { maxAge: 60 * 60 * 24 })

  const isAuthenticated = computed(() => !!user.value && !!token.value)

  const setUser = (userData: User) => {
    user.value = userData
  }

  const setToken = (newToken: string) => {
    token.value = newToken
  }

  const logout = () => {
    user.value = null
    token.value = null
  }

  return {
    user,
    token,
    isAuthenticated,
    setUser,
    setToken,
    logout,
  }
})
