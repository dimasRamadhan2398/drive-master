import { defineNuxtConfig } from "nuxt/config";
import react from "@vitejs/plugin-react";

export default defineNuxtConfig({
  devtools: { enabled: true },
  ssr: true,
  dev: true,
  modules: ["@nuxt/ui", "@pinia/nuxt", "nuxt-gtag", "@nuxtjs/i18n"],
  i18n: {
    locales: [
      { code: "id", iso: "id-ID", name: "Bahasa Indonesia", file: "id.json" },
      { code: "en", iso: "en-US", name: "English", file: "en.json" },
    ],
    defaultLocale: "id",
    langDir: "locales",
    strategy: "no_prefix",
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: "i18n_redirected",
      redirectOn: "root",
    },
  },
  hooks: {
    "pages:extend"(pages) {
      function setMiddleware(pageList: any[]) {
        for (const page of pageList) {
          if (page.path.startsWith("admin")) {
            page.meta ||= {};
            page.meta.middleware = ["admin"];
          }
          if (page.children) {
            setMiddleware(page.children);
          }
        }
      }
      setMiddleware(pages);
    },
  },
  css: ["~/assets/css/main.css"],
  colorMode: {
    preference: "dark",
  },
  runtimeConfig: {
    public: {
      apiBase:
        process.env.NUXT_PUBLIC_MODE == "dev"
          ? process.env.NUXT_PUBLIC_API_BASE_URL
          : "https://api.drivemaster.id/api/v1",
      // userApiBase:
      //   process.env.NUXT_PUBLIC_USER_API_BASE ||
      //   (process.env.NUXT_PUBLIC_API_BASE_URL
      //     ? process.env.NUXT_PUBLIC_API_BASE_URL + "/api/v1/users"
      //     : "http://localhost:8001/api/v1"),
      // coreApiBase:
      //   process.env.NUXT_PUBLIC_CORE_API_BASE ||
      //   (process.env.NUXT_PUBLIC_API_BASE_URL
      //     ? process.env.NUXT_PUBLIC_API_BASE_URL + "/api/v1/core"
      //     : "http://localhost:8002/api/v1"),
      // bookingApiBase:
      //   process.env.NUXT_PUBLIC_BOOKING_API_BASE ||
      //   (process.env.NUXT_PUBLIC_API_BASE_URL
      //     ? process.env.NUXT_PUBLIC_API_BASE_URL + "/api/v1/bookings"
      //     : "http://localhost:8003/api/v1"),
      gaMeasurementId:
        process.env.NUXT_PUBLIC_GA_MEASUREMENT_ID || "G-07PS1N5DZ5",
      gaPropertyId: process.env.NUXT_PUBLIC_GA_PROPERTY_ID || "G-539969879",
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
  gtag: {
    id: process.env.NUXT_PUBLIC_GA_MEASUREMENT_ID || "G-07PS1N5DZ5",
    config: {
      page_title: "Drive Master Indonesia - Premium Driving Academy",
      send_page_view: true,
    },
  },
  vite: {
    plugins: [react()],
    server: {
      watch: {
        usePolling: true,
      },
    },
  },
  compatibilityDate: "2026-04-07",
});
