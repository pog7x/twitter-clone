<script setup lang="ts">
import { NuxtLink } from '#components'

interface Props {
  label: string
  icon: string
  to?: string
  active?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  to: '',
  active: false,
})

const isNavigable = computed(() => props.active && props.to)
const tag = computed(() => (isNavigable.value ? NuxtLink : 'div'))
</script>

<template>
  <component
    :is="tag"
    :to="isNavigable ? props.to : undefined"
    class="group flex cursor-pointer items-center gap-4 rounded-full px-4 py-3 transition-colors hover:bg-gray-100 dark:hover:bg-gray-900"
    :tabindex="0"
    :aria-label="props.label"
    role="link"
  >
    <UIcon
      :name="props.icon"
      class="size-6 text-gray-900 transition-colors group-hover:text-sky-500 dark:text-white"
    />
    <span
      class="text-xl text-gray-900 transition-colors group-hover:text-sky-500 dark:text-white"
    >
      {{ props.label }}
    </span>
  </component>
</template>
