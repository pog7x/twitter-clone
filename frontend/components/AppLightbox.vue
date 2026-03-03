<script setup lang="ts">
const uiStore = useUiStore()

const currentIndex = ref(0)

const hasMultipleImages = computed(() => uiStore.lightbox.images.length > 1)
const currentImage = computed(() => uiStore.lightbox.images[currentIndex.value])

watch(
  () => uiStore.lightbox.isOpen,
  (isOpen) => {
    if (isOpen) {
      currentIndex.value = uiStore.lightbox.index
    }
  }
)

const handleNext = () => {
  const total = uiStore.lightbox.images.length
  currentIndex.value = (currentIndex.value + 1) % total
}

const handlePrev = () => {
  const total = uiStore.lightbox.images.length
  currentIndex.value = (currentIndex.value - 1 + total) % total
}

const handleClose = () => {
  uiStore.closeLightbox()
}

const handleBackdropClick = (e: MouseEvent) => {
  if (e.target === e.currentTarget) {
    handleClose()
  }
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') handleClose()
  if (e.key === 'ArrowLeft') handlePrev()
  if (e.key === 'ArrowRight') handleNext()
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="uiStore.lightbox.isOpen"
        class="fixed inset-0 z-[9999] flex items-center justify-center bg-black/80"
        role="dialog"
        aria-modal="true"
        aria-label="Image viewer"
        @click="handleBackdropClick"
      >
        <!-- Close button -->
        <UButton
          icon="i-lucide-x"
          variant="ghost"
          color="neutral"
          size="lg"
          class="absolute right-4 top-4 rounded-full"
          aria-label="Close lightbox"
          @click="handleClose"
        />

        <!-- Image -->
        <div class="flex max-h-[85vh] max-w-[75vw] items-center justify-center">
          <img
            v-if="currentImage"
            :src="currentImage"
            alt="Full size image"
            class="max-h-[85vh] max-w-[75vw] select-none object-contain"
          />
        </div>

        <!-- Navigation -->
        <template v-if="hasMultipleImages">
          <UButton
            icon="i-lucide-chevron-left"
            variant="solid"
            color="neutral"
            size="lg"
            class="absolute left-4 top-1/2 -translate-y-1/2 rounded-full"
            aria-label="Previous image"
            @click="handlePrev"
          />
          <UButton
            icon="i-lucide-chevron-right"
            variant="solid"
            color="neutral"
            size="lg"
            class="absolute right-4 top-1/2 -translate-y-1/2 rounded-full"
            aria-label="Next image"
            @click="handleNext"
          />
          <div class="absolute bottom-6 left-1/2 -translate-x-1/2 select-none font-bold text-white">
            {{ currentIndex + 1 }} / {{ uiStore.lightbox.images.length }}
          </div>
        </template>
      </div>
    </Transition>
  </Teleport>
</template>
