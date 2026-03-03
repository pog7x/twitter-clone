<script setup lang="ts">
import type { User } from '~/types'
import { formatJoinDate, getMediaUrl } from '~/utils/format'

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

const isMe = computed(() => authStore.user?.id === Number(props.profileId))

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

const fetchProfile = async () => {
  if (isMe.value && authStore.user) {
    profile.value = authStore.user
    return
  }

  try {
    const response = await api.getUser(props.profileId)
    profile.value = response.result
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to load profile',
      color: 'red',
      icon: 'i-lucide-alert-circle',
    })
  }
}

watch(() => props.profileId, fetchProfile)
onMounted(fetchProfile)

watch(
  () => authStore.user,
  (newUser) => {
    if (isMe.value && newUser) {
      profile.value = newUser
    }
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
    <div class="px-4 pb-4">
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
          <strong>{{ profile.followings?.length ?? 0 }}</strong>
          <span class="text-gray-500 dark:text-gray-400"> Following</span>
        </span>
        <span class="text-gray-900 dark:text-white">
          <strong>{{ profile.followers?.length ?? 0 }}</strong>
          <span class="text-gray-500 dark:text-gray-400"> Followers</span>
        </span>
      </div>
    </div>
  </div>
</template>
