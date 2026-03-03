import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', () => {
  const isLoading = ref(false)
  const isTweetModalOpen = ref(false)
  const isEditProfileOpen = ref(false)
  const isMobileMenuOpen = ref(false)
  const profileTweetCount = ref(0)

  const lightbox = ref({
    isOpen: false,
    images: [] as string[],
    index: 0,
  })

  const tweetCreatedSignal = ref(0)

  const toggleTweetModal = () => {
    isTweetModalOpen.value = !isTweetModalOpen.value
  }

  const notifyTweetCreated = () => {
    tweetCreatedSignal.value++
  }

  const openLightbox = (images: string[], index: number) => {
    lightbox.value = { isOpen: true, images, index }
  }

  const closeLightbox = () => {
    lightbox.value = { isOpen: false, images: [], index: 0 }
  }

  return {
    isLoading,
    isTweetModalOpen,
    isEditProfileOpen,
    isMobileMenuOpen,
    profileTweetCount,
    lightbox,
    tweetCreatedSignal,
    toggleTweetModal,
    notifyTweetCreated,
    openLightbox,
    closeLightbox,
  }
})
