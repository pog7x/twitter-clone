<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()
const config = useRuntimeConfig()

const isOpen = ref(false)
const popupRef = ref<HTMLElement | null>(null)

const handleToggle = () => {
  isOpen.value = !isOpen.value
}

const handleClickOutside = (event: MouseEvent) => {
  if (!popupRef.value) return
  if (!popupRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

watch(isOpen, (open) => {
  if (open) {
    document.addEventListener('click', handleClickOutside, { capture: true })
  } else {
    document.removeEventListener('click', handleClickOutside, { capture: true })
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside, { capture: true })
})

const handleViewProfile = () => {
  isOpen.value = false
  router.push(`/profile/${authStore.user?.id}`)
}

const handleLogout = () => {
  authStore.logout()
  isOpen.value = false
  router.push('/login')
}

const avatarUrl = computed(() => {
  if (!authStore.user?.pic) return ''
  return `${config.public.apiBaseUrl}${authStore.user.pic}`
})
</script>

<template>
  <div v-if="authStore.user" ref="popupRef" class="relative mb-3">
    <!-- Popup menu -->
    <Transition name="fade">
      <div
        v-if="isOpen"
        class="absolute bottom-full left-0 mb-3 w-full rounded-2xl border border-gray-200 bg-white shadow-lg dark:border-gray-800 dark:bg-black"
      >
        <div class="flex items-center gap-3 p-3">
          <UAvatar :src="avatarUrl" :alt="authStore.user.name" size="sm" />
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-bold text-gray-900 dark:text-white">
              {{ authStore.user.name }}
            </p>
            <p class="truncate text-sm text-gray-500 dark:text-gray-400">
              {{ authStore.user.nickname }}
            </p>
          </div>
          <UIcon name="i-lucide-check" class="size-5 text-sky-500" />
        </div>

        <div class="border-t border-gray-200 dark:border-gray-800">
          <button
            class="w-full px-4 py-3 text-left text-sm text-gray-900 transition-colors hover:bg-gray-100 dark:text-white dark:hover:bg-gray-900"
            aria-label="View profile"
            tabindex="0"
            @click="handleViewProfile"
          >
            View profile
          </button>
        </div>

        <div class="border-t border-gray-200 dark:border-gray-800">
          <button
            class="w-full px-4 py-3 text-left text-sm text-gray-900 transition-colors hover:bg-gray-100 dark:text-white dark:hover:bg-gray-900"
            @click="handleLogout"
          >
            Log out <span class="text-sky-500">{{ authStore.user.nickname }}</span>
          </button>
        </div>

        <div class="border-t border-gray-200 dark:border-gray-800">
          <ThemeToggle />
        </div>
      </div>
    </Transition>

    <!-- Profile button -->
    <button
      class="flex w-full items-center gap-3 rounded-full p-3 transition-colors hover:bg-gray-100 dark:hover:bg-gray-900"
      aria-label="Account menu"
      tabindex="0"
      @click="handleToggle"
    >
      <UAvatar :src="avatarUrl" :alt="authStore.user.name" size="sm" />
      <div class="hidden min-w-0 flex-1 text-left xl:block">
        <p class="truncate text-sm font-bold text-gray-900 dark:text-white">
          {{ authStore.user.name }}
        </p>
        <p class="truncate text-sm text-gray-500 dark:text-gray-400">
          {{ authStore.user.nickname }}
        </p>
      </div>
      <UIcon
        name="i-lucide-more-horizontal"
        class="hidden size-5 text-gray-500 xl:block"
      />
    </button>
  </div>
</template>
