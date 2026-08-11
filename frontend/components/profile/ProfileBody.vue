<script setup lang="ts">
interface Props {
  profileId: string
}

const props = defineProps<Props>()

const api = useApi()

const activeTab = ref('tweets')

const tabs = [
  { key: 'tweets', label: 'Posts' },
  { key: 'replies', label: 'Replies' },
  { key: 'media', label: 'Media' },
  { key: 'likes', label: 'Likes' },
]

const { tweets, isLoading, hasMore, refresh, loadMore } = usePaginatedTweets((query) =>
  api.getUserTweets(props.profileId, query)
)

watch(() => props.profileId, refresh)
onMounted(refresh)
</script>

<template>
  <div>
    <!-- Tabs -->
    <div class="flex border-b border-gray-200 dark:border-gray-800">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="flex-1 px-4 py-4 text-center text-sm font-bold transition-colors hover:bg-gray-100 dark:hover:bg-gray-900"
        :class="
          activeTab === tab.key
            ? 'text-gray-900 dark:text-white'
            : 'text-gray-500 dark:text-gray-400'
        "
        @click="activeTab = tab.key"
      >
        <span
          class="relative inline-block pb-3"
          :class="{ 'border-b-2 border-sky-500': activeTab === tab.key }"
        >
          {{ tab.label }}
        </span>
      </button>
    </div>

    <!-- Tweets list -->
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
      class="flex flex-col items-center justify-center py-16 text-gray-500 dark:text-gray-400"
    >
      <p class="text-lg">No posts yet</p>
    </div>
  </div>
</template>
