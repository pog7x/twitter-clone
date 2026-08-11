<script setup lang="ts">
import type { UpdateUserPayload } from '~/types'
import { getMediaUrl } from '~/utils/format'

// Kept in sync with the limits enforced by the API.
const MAX_NAME = 50
const MAX_DESCRIPTION = 160
const MAX_WEBSITE = 100

type ImageKind = 'pic' | 'pic_cover'

const api = useApi()
const authStore = useAuthStore()
const uiStore = useUiStore()
const config = useRuntimeConfig()
const toast = useToast()

const isSubmitting = ref(false)
const isUploading = ref(false)

const formData = reactive({
  name: '',
  description: '',
  website: '',
})

// Uploads happen as soon as a file is picked; only the resulting media id is
// sent on save, so the API resolves the stored path itself.
const picMediaId = ref<number | null>(null)
const picCoverMediaId = ref<number | null>(null)
const picPreview = ref('')
const picCoverPreview = ref('')

const avatarInputRef = ref<HTMLInputElement | null>(null)
const coverInputRef = ref<HTMLInputElement | null>(null)

const nameError = computed(() => {
  if (formData.name.trim().length === 0) return 'Name is required'
  if (formData.name.length > MAX_NAME) return `At most ${MAX_NAME} characters`
  return undefined
})

const descriptionError = computed(() =>
  formData.description.length > MAX_DESCRIPTION ? `At most ${MAX_DESCRIPTION} characters` : undefined
)

const websiteError = computed(() => {
  if (!formData.website) return undefined
  if (formData.website.length > MAX_WEBSITE) return `At most ${MAX_WEBSITE} characters`

  try {
    const url = new URL(formData.website)
    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      return 'Must start with http:// or https://'
    }
  } catch {
    return 'Please enter a valid URL'
  }

  return undefined
})

const canSubmit = computed(
  () =>
    !nameError.value &&
    !descriptionError.value &&
    !websiteError.value &&
    !isSubmitting.value &&
    !isUploading.value
)

watch(
  () => uiStore.isEditProfileOpen,
  (isOpen) => {
    if (!isOpen || !authStore.user) return

    formData.name = authStore.user.name ?? ''
    formData.description = authStore.user.description ?? ''
    formData.website = authStore.user.website ?? ''

    picMediaId.value = null
    picCoverMediaId.value = null
    picPreview.value = getMediaUrl(authStore.user.pic, config.public.apiBaseUrl as string)
    picCoverPreview.value = getMediaUrl(
      authStore.user.pic_cover,
      config.public.apiBaseUrl as string
    )
  }
)

const handleFileSelect = async (event: Event, kind: ImageKind) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''

  if (!file) return

  isUploading.value = true
  try {
    const body = new FormData()
    body.append('file', file)

    const { result } = await api.uploadMedia(body)
    const link = getMediaUrl(result.link, config.public.apiBaseUrl as string)

    if (kind === 'pic') {
      picMediaId.value = result.media_id
      picPreview.value = link
    } else {
      picCoverMediaId.value = result.media_id
      picCoverPreview.value = link
    }
  } catch {
    toast.add({
      title: 'Error',
      description: 'Failed to upload image. Use a jpeg, png, gif or webp up to 8 MB.',
      color: 'error',
      icon: 'i-lucide-alert-circle',
    })
  } finally {
    isUploading.value = false
  }
}

const handleSubmit = async () => {
  if (!canSubmit.value) return

  const payload: UpdateUserPayload = {
    name: formData.name,
    description: formData.description,
    website: formData.website,
  }

  if (picMediaId.value !== null) payload.pic_media_id = picMediaId.value
  if (picCoverMediaId.value !== null) payload.pic_cover_media_id = picCoverMediaId.value

  isSubmitting.value = true
  try {
    const { result } = await api.updateUser(payload)
    authStore.setUser(result)

    toast.add({ title: 'Profile updated', color: 'success', icon: 'i-lucide-check-circle' })
    uiStore.isEditProfileOpen = false
  } catch {
    toast.add({ title: 'Error', description: 'Failed to update profile', color: 'error' })
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <UModal v-model:open="uiStore.isEditProfileOpen">
    <template #content>
      <div class="p-5">
        <!-- Header -->
        <div class="mb-5 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <UButton
              icon="i-lucide-x"
              variant="ghost"
              color="neutral"
              size="sm"
              aria-label="Close"
              @click="uiStore.isEditProfileOpen = false"
            />
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">Edit profile</h3>
          </div>
          <UButton
            label="Save"
            size="sm"
            :disabled="!canSubmit"
            :loading="isSubmitting"
            @click="handleSubmit"
          />
        </div>

        <!-- Cover with avatar overlay -->
        <div class="mb-12">
          <div class="relative h-36 overflow-hidden rounded-lg bg-gray-200 dark:bg-gray-800">
            <img
              v-if="picCoverPreview"
              :src="picCoverPreview"
              alt="Cover preview"
              class="h-full w-full object-cover"
            />
            <UButton
              icon="i-lucide-camera"
              color="neutral"
              variant="solid"
              class="absolute inset-0 m-auto size-10 justify-center rounded-full opacity-80"
              aria-label="Change cover image"
              :disabled="isUploading"
              @click="coverInputRef?.click()"
            />
            <input
              ref="coverInputRef"
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              hidden
              @change="handleFileSelect($event, 'pic_cover')"
            />
          </div>

          <div class="relative -mt-10 ml-4 w-fit">
            <div class="rounded-full border-4 border-white dark:border-gray-900">
              <UAvatar :src="picPreview" :alt="formData.name" size="3xl" />
            </div>
            <UButton
              icon="i-lucide-camera"
              color="neutral"
              variant="solid"
              class="absolute inset-0 m-auto size-9 justify-center rounded-full opacity-80"
              aria-label="Change avatar"
              :disabled="isUploading"
              @click="avatarInputRef?.click()"
            />
            <input
              ref="avatarInputRef"
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              hidden
              @change="handleFileSelect($event, 'pic')"
            />
          </div>
        </div>

        <!-- Form -->
        <div class="space-y-4">
          <UFormField label="Name" :error="nameError">
            <UInput
              v-model="formData.name"
              placeholder="Your name"
              size="lg"
              :maxlength="MAX_NAME"
              @keydown.enter="handleSubmit"
            />
          </UFormField>

          <UFormField label="Description" :error="descriptionError">
            <UInput
              v-model="formData.description"
              placeholder="Describe yourself"
              size="lg"
              :maxlength="MAX_DESCRIPTION"
              @keydown.enter="handleSubmit"
            />
          </UFormField>

          <UFormField label="Website" :error="websiteError">
            <UInput
              v-model="formData.website"
              placeholder="https://example.com"
              type="url"
              size="lg"
              :maxlength="MAX_WEBSITE"
              @keydown.enter="handleSubmit"
            />
          </UFormField>
        </div>
      </div>
    </template>
  </UModal>
</template>
