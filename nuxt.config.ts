import { defineNuxtConfig } from "nuxt/config";

export default defineNuxtConfig({
  devtools: { enabled: true },
  ssr: true,
  dev: true,
  modules: ["@nuxt/ui", "@pinia/nuxt"],
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
          ? "http://localhost:8080/api/v1"
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
        process.env.NUXT_PUBLIC_GA_MEASUREMENT_ID || "G-MOCK123456",
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
