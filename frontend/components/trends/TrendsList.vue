<script setup lang="ts">
import type { Trend } from '~/types'

const api = useApi()
const toast = useToast()

const trends = ref<Trend[]>([])

const sortedTrends = computed(() =>
  [...trends.value].sort((a, b) => b.tweets_count - a.tweets_count)
)

onMounted(async () => {
  try {
    const response = await api.getTrends()
    trends.value = response.result ?? []
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to load trends',
      color: 'red',
    })
  }
})
</script>

<template>
  <div
    v-if="sortedTrends.length > 0"
    class="overflow-hidden rounded-2xl bg-gray-50 dark:bg-gray-900/50"
  >
    <h3 class="px-4 py-3 text-xl font-bold text-gray-900 dark:text-white">Trends</h3>

    <TrendsItem v-for="trend in sortedTrends" :key="trend.id" :trend="trend" />

    <button
      class="w-full px-4 py-3 text-left text-sm text-sky-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
    >
      Show more
    </button>
  </div>
</template>
