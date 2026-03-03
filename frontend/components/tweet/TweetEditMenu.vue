<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

const emit = defineEmits<{
  (e: 'toggle'): void
  (e: 'edit'): void
  (e: 'delete'): void
}>()

const items: DropdownMenuItem[][] = [
  [
    {
      label: 'Edit',
      icon: 'i-lucide-pencil',
      onSelect: () => emit('edit'),
    },
    {
      label: 'Delete',
      icon: 'i-lucide-trash-2',
      color: 'error' as const,
      onSelect: () => emit('delete'),
    },
  ],
]
</script>

<template>
  <UDropdownMenu :items="items">
    <UButton
      icon="i-lucide-more-horizontal"
      variant="ghost"
      color="neutral"
      size="xs"
      aria-label="More options"
    />
  </UDropdownMenu>
</template>
