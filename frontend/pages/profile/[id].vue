<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const uiStore = useUiStore()

const profileId = computed(() => String(route.params.id))
</script>

<template>
  <div>
    <!-- Page header -->
    <div
      class="sticky top-0 z-10 flex items-center gap-4 border-b border-gray-200 bg-white/80 px-4 py-2 backdrop-blur-md dark:border-gray-800 dark:bg-black/80"
    >
      <UButton
        icon="i-lucide-arrow-left"
        variant="ghost"
        color="neutral"
        size="sm"
        class="rounded-full"
        aria-label="Go back"
        @click="router.push('/')"
      />
      <div>
        <h2 class="text-xl font-bold text-gray-900 dark:text-white">Profile</h2>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          {{ uiStore.profileTweetCount }} tweets
        </p>
      </div>
    </div>

    <ProfileHeader :profile-id="profileId" />
    <ProfileBody :profile-id="profileId" />
    <EditProfileModal />
  </div>
</template>
