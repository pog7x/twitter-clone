export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore()
  const api = useApi()

  if (to.path === '/login') {
    if (authStore.isAuthenticated) {
      return navigateTo('/')
    }
    return
  }

  if (!authStore.token) {
    return navigateTo('/login')
  }

  if (!authStore.user) {
    try {
      const response = await api.getMe()
      authStore.setUser(response.result)
    } catch {
      authStore.logout()
      return navigateTo('/login')
    }
  }
})
