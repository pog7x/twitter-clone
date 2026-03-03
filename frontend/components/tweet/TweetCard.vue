<script setup lang="ts">
import type { Tweet } from '~/types'
import { formatRelativeTime, getMediaUrl } from '~/utils/format'

interface Props {
  tweet: Tweet
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'deleted'): void
  (e: 'updated'): void
}>()

const api = useApi()
const authStore = useAuthStore()
const uiStore = useUiStore()
const config = useRuntimeConfig()
const toast = useToast()

const isEditing = ref(false)
const editedContent = ref(props.tweet.content)

const avatarUrl = computed(() =>
  getMediaUrl(props.tweet.author?.pic, config.public.apiBaseUrl as string)
)

const tweetImages = computed(() =>
  props.tweet.attachments?.map((path) => getMediaUrl(path, config.public.apiBaseUrl as string)) ?? []
)

const isLikedByUser = computed(() =>
  props.tweet.likes?.some((like) => like.user_id === authStore.user?.id) ?? false
)

const likeCount = computed(() => props.tweet.likes?.length ?? 0)

const handleLikeClick = async () => {
  try {
    if (isLikedByUser.value) {
      await api.unlikeTweet(props.tweet.id)
    } else {
      await api.likeTweet(props.tweet.id)
    }
    emit('updated')
  } catch {
    toast.add({ title: 'Error', description: 'Failed to update like', color: 'error' })
  }
}

const handleDelete = async () => {
  try {
    await api.deleteTweet(props.tweet.id)
    toast.add({ title: 'Post deleted', color: 'success', icon: 'i-lucide-check-circle' })
    emit('deleted')
  } catch {
    toast.add({ title: 'Error', description: 'Failed to delete post', color: 'error' })
  }
}

const handleStartEdit = () => {
  editedContent.value = props.tweet.content
  isEditing.value = true
}

const handleSaveEdit = async () => {
  try {
    await api.updateTweet({ id: props.tweet.id, tweet_data: editedContent.value })
    toast.add({ title: 'Post updated', color: 'success', icon: 'i-lucide-check-circle' })
    isEditing.value = false
    emit('updated')
  } catch {
    toast.add({ title: 'Error', description: 'Failed to update post', color: 'error' })
  }
}

const handleCancelEdit = () => {
  isEditing.value = false
  editedContent.value = props.tweet.content
}

const handleImageClick = (index: number) => {
  uiStore.openLightbox(tweetImages.value, index)
}
</script>

<template>
  <article
    class="flex gap-3 border-b border-gray-200 px-4 py-3 transition-colors hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-gray-900/50"
  >
    <!-- Avatar -->
    <NuxtLink :to="`/profile/${tweet.author?.id}`" class="shrink-0">
      <UAvatar :src="avatarUrl" :alt="tweet.author?.name" size="lg" />
    </NuxtLink>

    <!-- Content -->
    <div class="min-w-0 flex-1">
      <!-- Header -->
      <div class="flex items-start justify-between">
        <div class="flex min-w-0 items-baseline gap-1">
          <NuxtLink
            :to="`/profile/${tweet.author?.id}`"
            class="truncate font-bold text-gray-900 hover:underline dark:text-white"
          >
            {{ tweet.author?.name }}
          </NuxtLink>
          <span class="truncate text-sm text-gray-500 dark:text-gray-400">
            {{ tweet.author?.nickname }}
          </span>
          <span class="text-sm text-gray-500 dark:text-gray-400">&middot;</span>
          <span class="whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
            {{ formatRelativeTime(tweet.created_at) }}
          </span>
        </div>

        <!-- Edit menu -->
        <TweetEditMenu
          @edit="handleStartEdit"
          @delete="handleDelete"
        />
      </div>

      <!-- Body -->
      <div v-if="!isEditing" class="mt-1">
        <p class="whitespace-pre-wrap text-gray-900 dark:text-white">{{ tweet.content }}</p>
      </div>

      <!-- Edit mode -->
      <div v-else class="mt-2">
        <textarea
          v-model="editedContent"
          class="w-full resize-none rounded-lg border border-gray-300 bg-transparent p-3 text-gray-900 outline-none focus:border-sky-500 dark:border-gray-700 dark:text-white"
          rows="3"
        />
        <div class="mt-2 flex justify-end gap-2">
          <UButton variant="outline" color="neutral" size="sm" label="Cancel" @click="handleCancelEdit" />
          <UButton size="sm" label="Save" @click="handleSaveEdit" />
        </div>
      </div>

      <!-- Images -->
      <div
        v-if="tweetImages.length > 0 && !isEditing"
        class="mt-3 grid gap-0.5 overflow-hidden rounded-2xl border border-gray-200 dark:border-gray-700"
        :class="{
          'grid-cols-1': tweetImages.length === 1,
          'grid-cols-2': tweetImages.length >= 2,
        }"
      >
        <div
          v-for="(image, i) in tweetImages"
          :key="i"
          class="cursor-zoom-in overflow-hidden"
          :class="{
            'col-span-1 row-span-2': tweetImages.length === 3 && i === 0,
          }"
          tabindex="0"
          role="button"
          :aria-label="`View image ${i + 1} of ${tweetImages.length}`"
          @click="handleImageClick(i)"
          @keydown.enter="handleImageClick(i)"
        >
          <img
            :src="image"
            alt="Post attachment"
            class="w-full object-cover"
            :class="{
              'max-h-[512px]': tweetImages.length === 1,
              'h-[286px]': tweetImages.length === 2,
              'h-full min-h-[286px]': tweetImages.length === 3 && i === 0,
              'h-[143px]': (tweetImages.length === 3 && i > 0) || tweetImages.length >= 4,
            }"
            loading="lazy"
          />
        </div>
      </div>

      <!-- Actions -->
      <div v-if="!isEditing" class="mt-3 flex max-w-[400px] items-center justify-between">
        <button
          class="group flex items-center gap-1 text-gray-500 transition-colors hover:text-sky-500"
          aria-label="Reply"
        >
          <span class="rounded-full p-2 transition-colors group-hover:bg-sky-500/10">
            <UIcon name="i-lucide-message-circle" class="size-[18px]" />
          </span>
        </button>

        <button
          class="group flex items-center gap-1 text-gray-500 transition-colors hover:text-green-500"
          aria-label="Repost"
        >
          <span class="rounded-full p-2 transition-colors group-hover:bg-green-500/10">
            <UIcon name="i-lucide-repeat-2" class="size-[18px]" />
          </span>
        </button>

        <button
          class="group flex items-center gap-1 transition-colors"
          :class="isLikedByUser ? 'text-pink-500' : 'text-gray-500 hover:text-pink-500'"
          aria-label="Like"
          @click="handleLikeClick"
        >
          <span class="rounded-full p-2 transition-colors group-hover:bg-pink-500/10">
            <UIcon
              name="i-lucide-heart"
              class="size-[18px]"
              :class="{ 'fill-current': isLikedByUser }"
            />
          </span>
          <span v-if="likeCount > 0" class="text-sm">{{ likeCount }}</span>
        </button>

        <button
          class="group flex items-center gap-1 text-gray-500 transition-colors hover:text-sky-500"
          aria-label="Share"
        >
          <span class="rounded-full p-2 transition-colors group-hover:bg-sky-500/10">
            <UIcon name="i-lucide-share" class="size-[18px]" />
          </span>
        </button>
      </div>
    </div>
  </article>
</template>
