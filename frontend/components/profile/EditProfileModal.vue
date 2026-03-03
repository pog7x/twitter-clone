<script setup lang="ts">
const api = useApi()
const authStore = useAuthStore()
const uiStore = useUiStore()
const toast = useToast()

const isSubmitting = ref(false)

const formData = reactive({
  name: '',
  description: '',
  website: '',
})

const isFormValid = computed(() => formData.name.length > 1 && formData.description.length > 2)

const isUrlValid = computed(() => {
  if (!formData.website) return true
  try {
    new URL(formData.website)
    return true
  } catch {
    return false
  }
})

const canSubmit = computed(() => isFormValid.value && isUrlValid.value && !isSubmitting.value)

watch(
  () => uiStore.isEditProfileOpen,
  (isOpen) => {
    if (isOpen && authStore.user) {
      formData.name = authStore.user.name ?? ''
      formData.description = authStore.user.description ?? ''
      formData.website = authStore.user.website ?? ''
    }
  }
)

const handleSubmit = async () => {
  if (!canSubmit.value) return

  isSubmitting.value = true
  try {
    await api.updateUser(formData)

    if (authStore.user) {
      authStore.setUser({
        ...authStore.user,
        name: formData.name,
        description: formData.description,
        website: formData.website,
      })
    }

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

        <!-- Form -->
        <div class="space-y-4">
          <UFormField label="Name">
            <UInput
              v-model="formData.name"
              placeholder="Your name"
              size="lg"
              @keydown.enter="handleSubmit"
            />
          </UFormField>

          <UFormField label="Description">
            <UInput
              v-model="formData.description"
              placeholder="Describe yourself"
              size="lg"
              @keydown.enter="handleSubmit"
            />
          </UFormField>

          <UFormField label="Website" :error="!isUrlValid ? 'Please enter a valid URL' : undefined">
            <UInput
              v-model="formData.website"
              placeholder="https://example.com"
              type="url"
              size="lg"
              @keydown.enter="handleSubmit"
            />
          </UFormField>
        </div>
      </div>
    </template>
  </UModal>
</template>
