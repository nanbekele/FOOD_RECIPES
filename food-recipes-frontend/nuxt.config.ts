// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  modules: ['@nuxtjs/tailwindcss'],
  css: ['@/assets/css/main.css'],
  runtimeConfig: {
    public: {
      NUXT_PUBLIC_HASURA_ENDPOINT: process.env.NUXT_PUBLIC_HASURA_ENDPOINT || 'http://localhost:8082/v1/graphql',
      NUXT_PUBLIC_BACKEND_ENDPOINT: process.env.NUXT_PUBLIC_BACKEND_ENDPOINT || 'http://localhost:8081',
    }
  }
})
