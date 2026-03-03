<script setup lang="ts">
const searchQuery = ref('')
</script>

<template>
  <div class="sticky top-0 z-10 bg-white pb-3 pt-1 dark:bg-black">
    <UInput
      v-model="searchQuery"
      placeholder="Search"
      icon="i-lucide-search"
      size="lg"
      aria-label="Search"
    >
      <template v-if="searchQuery" #trailing>
        <UButton
          icon="i-lucide-x"
          variant="ghost"
          color="neutral"
          size="2xs"
          aria-label="Clear search"
          @click="searchQuery = ''"
        />
      </template>
    </UInput>
  </div>
</template>
