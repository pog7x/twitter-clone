<script setup lang="ts">
useHead({
  title: 'Twitter Clone',
  meta: [{ name: 'description', content: 'A Twitter/X clone built with Nuxt' }],
  htmlAttrs: { lang: 'en' },
})
</script>

<template>
  <UApp>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <AppLightbox />
  </UApp>
</template>
