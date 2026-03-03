<script setup lang="ts">
const authStore = useAuthStore()
const uiStore = useUiStore()

const NAV_ITEMS = [
  { label: 'Home', icon: 'i-lucide-home', to: '/', active: true },
  { label: 'Explore', icon: 'i-lucide-hash' },
  { label: 'Notifications', icon: 'i-lucide-bell' },
  { label: 'Messages', icon: 'i-lucide-mail' },
  { label: 'Bookmarks', icon: 'i-lucide-bookmark' },
  { label: 'Lists', icon: 'i-lucide-list' },
]
</script>

<template>
  <div class="flex h-full w-full flex-col justify-between px-3 py-2">
    <nav class="flex flex-col gap-1">
      <!-- Logo -->
      <NuxtLink
        to="/"
        class="mb-2 flex h-12 w-12 items-center justify-center rounded-full transition-colors hover:bg-gray-100 dark:hover:bg-gray-900"
        aria-label="Home"
      >
        <UIcon name="i-simple-icons-x" class="size-7 text-gray-900 dark:text-white" />
      </NuxtLink>

      <!-- Nav items -->
      <SidebarItem
        v-for="item in NAV_ITEMS"
        :key="item.label"
        :label="item.label"
        :icon="item.icon"
        :to="item.to"
        :active="item.active"
      />

      <!-- Profile link -->
      <SidebarItem
        v-if="authStore.user"
        label="Profile"
        icon="i-lucide-user"
        :to="`/profile/${authStore.user.id}`"
        active
      />

      <!-- More menu -->
      <SidebarItem label="More" icon="i-lucide-more-horizontal" />

      <!-- Tweet button -->
      <UButton
        block
        size="lg"
        class="mt-4"
        label="Post"
        @click="uiStore.isMobileMenuOpen = false; uiStore.toggleTweetModal()"
      />
    </nav>

    <!-- Profile popup -->
    <SidebarProfilePopup />
  </div>
</template>
