import { defineNuxtConfig } from "nuxt/config";

export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: ['@nuxt/ui'],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api'
    }
  },
  app: {
    head: {
      title: 'Drive Master Indonesia',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        //{ name: 'description', content: 'The first premium driving academy in Alam Sutera using 100% Electric Vehicles. Experience smooth, silent, and sustainable learning.' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/drive-master-icon.svg' }
      ]
    }
  },
  vite: {
    server: {
      watch: {
        usePolling: true
      }
    }
  },
  compatibilityDate: '2026-04-07'
})
