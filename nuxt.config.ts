import { defineNuxtConfig } from "nuxt/config";

export default defineNuxtConfig({
  devtools: { enabled: false },
  modules: ["@nuxt/ui", "@pinia/nuxt"],
  css: ["~/assets/css/main.css"],
  colorMode: {
    preference: "dark",
  },
  runtimeConfig: {
    public: {
      userApiBase:
        process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8001/api/v1",
      coreApiBase:
        process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8002/api/v1",
      bookingApiBase:
        process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8003/api/v1",
    },
  },
  app: {
    head: {
      title: "Drive Master Indonesia",
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
        //{ name: 'description', content: 'The first premium driving academy in Alam Sutera using 100% Electric Vehicles. Experience smooth, silent, and sustainable learning.' }
      ],
      link: [
        { rel: "icon", type: "image/svg+xml", href: "/drive-master-icon.svg" },
      ],
    },
  },
  vite: {
    server: {
      watch: {
        usePolling: true,
      },
    },
  },
  compatibilityDate: "2026-04-07",
});
