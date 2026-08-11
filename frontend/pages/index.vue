<script setup lang="ts">
const api = useApi()
const uiStore = useUiStore()

const { tweets, isLoading, hasMore, refresh, loadMore } = usePaginatedTweets((query) =>
  api.getTweets(query)
)

onMounted(refresh)

watch(() => uiStore.tweetCreatedSignal, refresh)
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
    <TweetComposer @tweet-created="refresh" />

    <!-- Divider -->
    <div class="h-2 bg-gray-100 dark:bg-gray-900/50" />

    <!-- Tweet feed -->
    <div v-if="tweets.length > 0">
      <TweetCard
        v-for="tweet in tweets"
        :key="tweet.id"
        :tweet="tweet"
        @deleted="refresh"
        @updated="refresh"
      />

      <div v-if="hasMore" class="flex justify-center p-4">
        <UButton
          variant="ghost"
          color="primary"
          label="Show more"
          :loading="isLoading"
          @click="loadMore"
        />
      </div>
    </div>

    <div
      v-else-if="!isLoading"
      class="flex flex-col items-center justify-center px-8 py-16 text-gray-500 dark:text-gray-400"
    >
      <UIcon name="i-lucide-message-circle" class="mb-4 size-12" />
      <p class="text-lg font-medium">No tweets yet</p>
      <p class="text-sm">Be the first to share something!</p>
    </div>
  </div>
</template>
