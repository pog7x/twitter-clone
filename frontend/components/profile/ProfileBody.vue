<script setup lang="ts">
import type { Tweet } from '~/types'

interface Props {
  profileId: string
}

const props = defineProps<Props>()

const api = useApi()
const uiStore = useUiStore()
const toast = useToast()

const tweets = ref<Tweet[]>([])
const activeTab = ref('tweets')

const tabs = [
  { key: 'tweets', label: 'Posts' },
  { key: 'replies', label: 'Replies' },
  { key: 'media', label: 'Media' },
  { key: 'likes', label: 'Likes' },
]

const fetchUserTweets = async () => {
  try {
    const response = await api.getUserTweets(props.profileId)
    tweets.value = response.result ?? []
    uiStore.profileTweetCount = tweets.value.length
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to load posts',
      color: 'red',
      icon: 'i-lucide-alert-circle',
    })
  }
}

const handleTweetDeleted = () => {
  fetchUserTweets()
}

const handleTweetUpdated = () => {
  fetchUserTweets()
}

watch(() => props.profileId, fetchUserTweets)
onMounted(fetchUserTweets)
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
        @deleted="handleTweetDeleted"
        @updated="handleTweetUpdated"
      />
    </div>

    <div
      v-else
      class="flex flex-col items-center justify-center py-16 text-gray-500 dark:text-gray-400"
    >
      <p class="text-lg">No posts yet</p>
    </div>
  </div>
</template>
