<script setup lang="ts">
const uiStore = useUiStore()

const handleTweetCreated = () => {
  uiStore.toggleTweetModal()
  uiStore.notifyTweetCreated()
}
</script>

<template>
  <UModal v-model:open="uiStore.isTweetModalOpen">
    <template #content>
      <div class="p-4">
        <div class="mb-3 flex items-center justify-between">
          <UButton
            icon="i-lucide-x"
            variant="ghost"
            color="neutral"
            size="sm"
            aria-label="Close"
            @click="uiStore.toggleTweetModal()"
          />
        </div>
        <TweetComposer @tweet-created="handleTweetCreated" />
      </div>
    </template>
  </UModal>
</template>
