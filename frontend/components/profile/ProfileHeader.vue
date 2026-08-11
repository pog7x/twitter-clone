<script setup lang="ts">
import type { User } from '~/types'
import { formatCount, formatJoinDate, getMediaUrl } from '~/utils/format'

interface Props {
  profileId: string
}

const props = defineProps<Props>()

const api = useApi()
const authStore = useAuthStore()
const uiStore = useUiStore()
const config = useRuntimeConfig()
const toast = useToast()

const profile = ref<User | null>(null)

const isMe = computed(() => profile.value?.is_me ?? false)

const coverUrl = computed(() =>
  getMediaUrl(profile.value?.pic_cover, config.public.apiBaseUrl as string)
)

const avatarUrl = computed(() =>
  getMediaUrl(profile.value?.pic, config.public.apiBaseUrl as string)
)

const joinedDate = computed(() =>
  profile.value?.created_at ? formatJoinDate(profile.value.created_at) : ''
)

const websiteDisplay = computed(() => {
  if (!profile.value?.website) return null
  try {
    const url = new URL(profile.value.website)
    return { display: url.host, full: profile.value.website }
  } catch {
    return null
  }
})

// The profile is always read from the API: counts such as tweets_count are
// computed server side and would go stale if taken from the auth store.
const fetchProfile = async () => {
  try {
    const { result } = await api.getUser(props.profileId)
    profile.value = result
    uiStore.profileTweetCount = result.tweets_count
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to load profile',
      color: 'error',
      icon: 'i-lucide-alert-circle',
    })
  }
}

watch(() => props.profileId, fetchProfile)
onMounted(fetchProfile)

// Pick up edits made through the profile modal.
watch(
  () => authStore.user,
  () => {
    if (isMe.value) fetchProfile()
  },
  { deep: true }
)
</script>

<template>
  <div v-if="profile">
    <!-- Cover -->
    <div class="h-48 bg-gray-200 dark:bg-gray-800">
      <img
        v-if="coverUrl"
        :src="coverUrl"
        alt="Cover"
        class="h-full w-full object-cover"
      />
    </div>

    <!-- Header section -->
    <div class="px-4 pt-3 pb-4">
      <!-- Avatar + edit button -->
      <div class="flex items-end justify-between">
        <div class="-mt-16 rounded-full border-4 border-white dark:border-black">
          <UAvatar
            :src="avatarUrl"
            :alt="profile.name"
            size="3xl"
            :ui="{ rounded: 'rounded-full' }"
          />
        </div>

        <UButton
          v-if="isMe"
          variant="outline"
          color="neutral"
          label="Edit profile"
          class="rounded-full font-bold"
          @click="uiStore.isEditProfileOpen = true"
        />
      </div>

      <!-- Name & username -->
      <div class="mt-3">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white">
          {{ profile.name }}
        </h2>
        <p class="text-gray-500 dark:text-gray-400">{{ profile.nickname }}</p>
      </div>

      <!-- Bio -->
      <p v-if="profile.description" class="mt-3 text-gray-900 dark:text-white">
        {{ profile.description }}
      </p>

      <!-- Meta info -->
      <div class="mt-3 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-gray-500 dark:text-gray-400">
        <a
          v-if="websiteDisplay"
          :href="websiteDisplay.full"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center gap-1 text-sky-500 hover:underline"
        >
          <UIcon name="i-lucide-link" class="size-4" />
          {{ websiteDisplay.display }}
        </a>
        <span v-if="joinedDate" class="flex items-center gap-1">
          <UIcon name="i-lucide-calendar" class="size-4" />
          Joined {{ joinedDate }}
        </span>
      </div>

      <!-- Follower counts -->
      <div class="mt-3 flex gap-4 text-sm">
        <span class="text-gray-900 dark:text-white">
          <strong>{{ formatCount(profile.followings_count) }}</strong>
          <span class="text-gray-500 dark:text-gray-400"> Following</span>
        </span>
        <span class="text-gray-900 dark:text-white">
          <strong>{{ formatCount(profile.followers_count) }}</strong>
          <span class="text-gray-500 dark:text-gray-400"> Followers</span>
        </span>
      </div>
    </div>
  </div>
</template>
