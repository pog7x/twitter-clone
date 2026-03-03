export default defineNuxtConfig({
  ssr: false,

  devtools: { enabled: true },

  modules: ['@nuxt/ui', '@pinia/nuxt'],

  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://127.0.0.1:8080',
    },
  },

  colorMode: {
    preference: 'dark',
    fallback: 'dark',
  },

  css: ['~/assets/css/main.css'],

  components: [
    {
      path: '~/components',
      pathPrefix: false,
    },
  ],

  compatibilityDate: '2025-01-01',
})
