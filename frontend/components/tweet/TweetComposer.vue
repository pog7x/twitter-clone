<script setup lang="ts">
import { getMediaUrl } from '~/utils/format'

// Kept in sync with the limits enforced by the API.
const MAX_IMAGES = 4
const MAX_TWEET_LENGTH = 280

const emit = defineEmits<{
  (e: 'tweetCreated'): void
}>()

const api = useApi()
const authStore = useAuthStore()
const config = useRuntimeConfig()
const toast = useToast()

const tweetText = ref('')
const imageFiles = ref<{ url: string; file: File }[]>([])
const isSubmitting = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)

const remainingChars = computed(() => MAX_TWEET_LENGTH - tweetText.value.trim().length)

// An empty text is allowed as long as the post carries an image.
const canSubmit = computed(
  () =>
    (tweetText.value.trim().length > 0 || imageFiles.value.length > 0) &&
    remainingChars.value >= 0 &&
    !isSubmitting.value
)

const canAddMoreImages = computed(() => imageFiles.value.length < MAX_IMAGES)
const remainingSlots = computed(() => MAX_IMAGES - imageFiles.value.length)

const avatarUrl = computed(() =>
  getMediaUrl(authStore.user?.pic, config.public.apiBaseUrl as string)
)

const clearImages = () => {
  imageFiles.value.forEach((image) => URL.revokeObjectURL(image.url))
  imageFiles.value = []
}

const handleSubmit = async () => {
  if (!canSubmit.value) return

  isSubmitting.value = true
  try {
    const mediaIds: number[] = []

    if (imageFiles.value.length > 0) {
      const uploads = imageFiles.value.map((image) => {
        const formData = new FormData()
        formData.append('file', image.file)
        return api.uploadMedia(formData)
      })
      const results = await Promise.all(uploads)
      results.forEach(({ result }) => mediaIds.push(result.media_id))
    }

    await api.createTweet({
      tweet_data: tweetText.value,
      tweet_media_ids: mediaIds,
    })

    toast.add({
      title: 'Post sent!',
      color: 'success',
      icon: 'i-lucide-check-circle',
    })

    tweetText.value = ''
    clearImages()
    emit('tweetCreated')
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to send post',
      color: 'error',
      icon: 'i-lucide-alert-circle',
    })
  } finally {
    isSubmitting.value = false
  }
}

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return

  const filesToAdd = Array.from(input.files).slice(0, remainingSlots.value)

  if (input.files.length > remainingSlots.value) {
    toast.add({
      title: `Maximum ${MAX_IMAGES} images`,
      description: `Only ${remainingSlots.value} more image(s) can be added`,
      color: 'warning',
      icon: 'i-lucide-alert-triangle',
    })
  }

  for (const file of filesToAdd) {
    const url = URL.createObjectURL(file)
    imageFiles.value.push({ url, file })
  }

  input.value = ''
}

const handleRemoveImage = (index: number) => {
  URL.revokeObjectURL(imageFiles.value[index].url)
  imageFiles.value.splice(index, 1)
}

onBeforeUnmount(clearImages)
</script>

<template>
  <div class="flex gap-3 border-b border-gray-200 px-4 pb-2 pt-3 dark:border-gray-800">
    <UAvatar :src="avatarUrl" :alt="authStore.user?.name" size="lg" class="shrink-0" />

    <div class="flex-1">
      <textarea
        v-model="tweetText"
        placeholder="What is happening?!"
        class="min-h-[80px] w-full resize-none bg-transparent text-xl text-gray-900 placeholder-gray-500 outline-none dark:text-white dark:placeholder-gray-500"
        rows="2"
      />

      <!-- Image previews -->
      <div
        v-if="imageFiles.length > 0"
        class="mb-3 grid gap-0.5 overflow-hidden rounded-2xl border border-gray-200 dark:border-gray-700"
        :class="{
          'grid-cols-1': imageFiles.length === 1,
          'grid-cols-2': imageFiles.length >= 2,
        }"
      >
        <div
          v-for="(image, i) in imageFiles"
          :key="i"
          class="relative overflow-hidden"
          :class="{
            'col-span-1 row-span-2': imageFiles.length === 3 && i === 0,
          }"
        >
          <img
            :src="image.url"
            alt="Preview"
            class="w-full object-cover"
            :class="{
              'h-64': imageFiles.length === 1,
              'h-40': imageFiles.length === 2,
              'h-full min-h-[160px]': imageFiles.length === 3 && i === 0,
              'h-[calc(50%-1px)]': imageFiles.length === 3 && i > 0,
              'h-32': imageFiles.length >= 4,
            }"
          />
          <UButton
            icon="i-lucide-x"
            size="xs"
            color="neutral"
            variant="solid"
            class="absolute right-1 top-1"
            aria-label="Remove image"
            @click="handleRemoveImage(i)"
          />
        </div>
      </div>

      <!-- Actions bar -->
      <div class="flex items-center justify-between border-t border-gray-200 pt-3 dark:border-gray-800">
        <div class="flex gap-1">
          <UButton
            icon="i-lucide-image"
            variant="ghost"
            color="primary"
            size="sm"
            :disabled="!canAddMoreImages"
            aria-label="Add images"
            @click="fileInputRef?.click()"
          />
          <input
            ref="fileInputRef"
            type="file"
            accept="image/jpeg,image/png,image/gif,image/webp"
            multiple
            hidden
            @change="handleFileSelect"
          />
          <UButton
            icon="i-lucide-smile"
            variant="ghost"
            color="primary"
            size="sm"
            aria-label="Add emoji"
          />
          <UButton
            icon="i-lucide-calendar"
            variant="ghost"
            color="primary"
            size="sm"
            aria-label="Schedule"
          />
        </div>

        <div class="flex items-center gap-3">
          <span
            v-if="tweetText.length > 0"
            class="text-sm"
            :class="remainingChars < 0 ? 'text-red-500' : 'text-gray-500 dark:text-gray-400'"
          >
            {{ remainingChars }}
          </span>
          <UButton
            label="Post"
            size="sm"
            :disabled="!canSubmit"
            :loading="isSubmitting"
            @click="handleSubmit"
          />
        </div>
      </div>
    </div>
  </div>
</template>
