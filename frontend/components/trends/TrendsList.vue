<script setup lang="ts">
import type { Trend } from '~/types'

// Trends are short, so a page is smaller than the shared feed page size.
const TRENDS_PAGE_SIZE = 5

const api = useApi()
const toast = useToast()

const trends = ref<Trend[]>([])
const isLoading = ref(false)
const hasMore = ref(false)

const sortedTrends = computed(() =>
  [...trends.value].sort((a, b) => b.tweets_count - a.tweets_count)
)

const load = async (offset: number) => {
  if (isLoading.value) return

  isLoading.value = true
  try {
    const { result } = await api.getTrends({ limit: TRENDS_PAGE_SIZE, offset })
    const batch = result ?? []

    trends.value = offset === 0 ? batch : [...trends.value, ...batch]
    hasMore.value = batch.length === TRENDS_PAGE_SIZE
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to load trends',
      color: 'error',
    })
  } finally {
    isLoading.value = false
  }
}

onMounted(() => load(0))
</script>

<template>
  <div
    v-if="sortedTrends.length > 0"
    class="overflow-hidden rounded-2xl bg-gray-50 dark:bg-gray-900/50"
  >
    <h3 class="px-4 py-3 text-xl font-bold text-gray-900 dark:text-white">Trends</h3>

    <TrendsItem v-for="trend in sortedTrends" :key="trend.id" :trend="trend" />

    <button
      v-if="hasMore"
      class="w-full px-4 py-3 text-left text-sm text-sky-500 transition-colors hover:bg-gray-100 disabled:opacity-60 dark:hover:bg-gray-800"
      :disabled="isLoading"
      @click="load(trends.length)"
    >
      Show more
    </button>
  </div>
</template>
