// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-09-01',
  devtools: { enabled: true },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt',
  ],

  css: ['~/assets/css/main.css'],

  app: {
    head: {
      title: 'War Tiket — Pemesanan Tiket Konser',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Platform pemesanan tiket konser dengan sistem war tiket yang adil.' },
      ],
      htmlAttrs: { lang: 'id' },
    },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
      siteName: process.env.NUXT_PUBLIC_SITE_NAME || 'War Tiket',
    },
  },

  typescript: {
    strict: true,
  },
})
