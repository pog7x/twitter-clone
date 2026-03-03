<script setup lang="ts">
import type { Trend } from '~/types'
import { formatCount } from '~/utils/format'

interface Props {
  trend: Trend
}

const props = defineProps<Props>()
</script>

<template>
  <div
    class="cursor-pointer px-4 py-3 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
    role="link"
    tabindex="0"
    :aria-label="`Trend: ${trend.name}`"
  >
    <p class="text-xs text-gray-500 dark:text-gray-400">Trending</p>
    <p class="font-bold text-gray-900 dark:text-white">{{ props.trend.name }}</p>
    <p class="text-xs text-gray-500 dark:text-gray-400">
      {{ formatCount(props.trend.tweets_count) }} posts
    </p>
  </div>
</template>
