<script setup lang="ts">
const uiStore = useUiStore()
const router = useRouter()

router.afterEach(() => {
  uiStore.isMobileMenuOpen = false
})
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-black">
    <div class="mx-auto flex w-full max-w-[1280px]">
      <!-- Sidebar -->
      <aside class="sticky top-0 hidden h-screen w-[275px] shrink-0 md:flex">
        <AppSidebar />
      </aside>

      <!-- Mobile menu overlay -->
      <Teleport to="body">
        <Transition name="fade">
          <div
            v-if="uiStore.isMobileMenuOpen"
            class="fixed inset-0 z-50 bg-black/50 md:hidden"
            @click="uiStore.isMobileMenuOpen = false"
          >
            <Transition name="slide">
              <div
                v-if="uiStore.isMobileMenuOpen"
                class="h-full w-[275px] bg-white dark:bg-black"
                @click.stop
              >
                <AppSidebar />
              </div>
            </Transition>
          </div>
        </Transition>
      </Teleport>

      <!-- Main content -->
      <main
        class="min-h-screen w-full max-w-[600px] border-x border-gray-200 dark:border-gray-800"
      >
        <slot />
      </main>

      <!-- Right sidebar -->
      <aside class="sticky top-0 hidden h-screen w-[350px] shrink-0 pl-6 pt-2 lg:block">
        <div class="flex h-full flex-col gap-4">
          <AppSearchBar />
          <TrendsList />
        </div>
      </aside>
    </div>

    <!-- Mobile bottom bar -->
    <div class="fixed bottom-4 right-4 z-40 md:hidden">
      <UButton
        icon="i-lucide-menu"
        size="xl"
        variant="solid"
        class="shadow-lg"
        aria-label="Open menu"
        @click="uiStore.isMobileMenuOpen = true"
      />
    </div>

    <!-- Tweet modal -->
    <TweetModal />
  </div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 200ms ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
.slide-enter-active,
.slide-leave-active {
  transition: transform 200ms ease;
}
.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}
</style>
