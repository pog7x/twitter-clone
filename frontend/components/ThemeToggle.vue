<script setup lang="ts">
const colorMode = useColorMode()

const isDark = computed({
  get: () => colorMode.value === 'dark',
  set: (val: boolean) => {
    colorMode.preference = val ? 'dark' : 'light'
  },
})

const handleToggle = () => {
  isDark.value = !isDark.value
}
</script>

<template>
  <button
    class="flex w-full items-center gap-3 px-4 py-3 text-sm text-gray-900 transition-colors hover:bg-gray-100 dark:text-white dark:hover:bg-gray-900"
    aria-label="Toggle dark mode"
    tabindex="0"
    @click="handleToggle"
  >
    <UIcon
      :name="isDark ? 'i-lucide-moon' : 'i-lucide-sun'"
      class="size-5"
    />
    <span>{{ isDark ? 'Dark mode' : 'Light mode' }}</span>
  </button>
</template>
