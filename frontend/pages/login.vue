<script setup lang="ts">
definePageMeta({
  layout: 'auth',
})

const api = useApi()
const authStore = useAuthStore()
const router = useRouter()
const toast = useToast()

const isSubmitting = ref(false)
const credentials = reactive({
  username: '',
  password: '',
})

const isFormValid = computed(() => credentials.username.length > 0 && credentials.password.length > 0)

const handleLogin = async () => {
  if (!isFormValid.value || isSubmitting.value) return

  isSubmitting.value = true
  try {
    const loginResponse = await api.login(credentials)
    if (!loginResponse.success) return

    authStore.setToken(loginResponse.result)

    const meResponse = await api.getMe()
    authStore.setUser(meResponse.result)

    router.push('/')
  } catch {
    toast.add({
      title: 'Authentication failed',
      description: 'Invalid username or password',
      color: 'error',
      icon: 'i-lucide-alert-circle',
    })
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="w-full max-w-[400px] px-6">
    <div class="mb-8">
      <UIcon name="i-simple-icons-x" class="mb-6 size-10 text-gray-900 dark:text-white" />
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Sign in</h1>
    </div>

    <form class="space-y-5" @submit.prevent="handleLogin">
      <UFormField label="Username">
        <UInput
          v-model="credentials.username"
          placeholder="Enter your username"
          size="lg"
          icon="i-lucide-user"
          autocomplete="username"
          autofocus
        />
      </UFormField>

      <UFormField label="Password">
        <UInput
          v-model="credentials.password"
          type="password"
          placeholder="Enter your password"
          size="lg"
          icon="i-lucide-lock"
          autocomplete="current-password"
          @keydown.enter="handleLogin"
        />
      </UFormField>

      <UButton
        type="submit"
        block
        size="lg"
        :loading="isSubmitting"
        :disabled="!isFormValid"
        label="Sign in"
      />
    </form>

    <div class="mt-6 text-center">
      <p class="text-sm text-gray-500 dark:text-gray-400">
        <span class="cursor-pointer text-sky-500 hover:underline">Forgot password?</span>
        <span class="mx-2">&middot;</span>
        <span class="cursor-pointer text-sky-500 hover:underline">Sign up</span>
      </p>
    </div>
  </div>
</template>
