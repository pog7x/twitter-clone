<script setup lang="ts">
import type { Tweet } from '~/types'

const api = useApi()
const toast = useToast()
const uiStore = useUiStore()

const tweets = ref<Tweet[]>([])

const fetchTweets = async () => {
  try {
    const response = await api.getTweets()
    tweets.value = response.result ?? []
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to load tweets',
      color: 'red',
      icon: 'i-lucide-alert-circle',
    })
  }
}

const handleTweetDeleted = () => {
  fetchTweets()
}

const handleTweetsUpdated = () => {
  fetchTweets()
}

onMounted(fetchTweets)

watch(() => uiStore.tweetCreatedSignal, fetchTweets)
</script>

<template>
  <div>
    <!-- Page header -->
    <div
      class="sticky top-0 z-10 border-b border-gray-200 bg-white/80 px-4 py-3 backdrop-blur-md dark:border-gray-800 dark:bg-black/80"
    >
      <h2 class="text-xl font-bold text-gray-900 dark:text-white">Home</h2>
    </div>

    <!-- Composer -->
    <TweetComposer @tweet-created="fetchTweets" />

    <!-- Divider -->
    <div class="h-2 bg-gray-100 dark:bg-gray-900/50" />

    <!-- Tweet feed -->
    <div v-if="tweets.length > 0">
      <TweetCard
        v-for="tweet in tweets"
        :key="tweet.id"
        :tweet="tweet"
        @deleted="handleTweetDeleted"
        @updated="handleTweetsUpdated"
      />
    </div>

    <div
      v-else
      class="flex flex-col items-center justify-center px-8 py-16 text-gray-500 dark:text-gray-400"
    >
      <UIcon name="i-lucide-message-circle" class="mb-4 size-12" />
      <p class="text-lg font-medium">No tweets yet</p>
      <p class="text-sm">Be the first to share something!</p>
    </div>
  </div>
</template>
