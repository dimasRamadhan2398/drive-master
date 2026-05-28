import { defineNuxtPlugin, useHead, useRouter, useRuntimeConfig } from "#app"

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const measurementId = config.public.gaMeasurementId

  if (!measurementId) {
    console.warn("Google Analytics Measurement ID is not configured. Analytics tracking will not be active.")
    return
  }

  // Inject Google Analytics Scripts into Page Head
  useHead({
    script: [
      {
        src: `https://www.googletagmanager.com/gtag/js?id=${measurementId}`,
        async: true,
      },
      {
        innerHTML: `
          window.dataLayer = window.dataLayer || [];
          function gtag(){dataLayer.push(arguments);}
          window.gtag = gtag;
          gtag('js', new Date());
          gtag('config', '${measurementId}', {
            page_path: window.location.pathname,
          });
        `,
        type: 'text/javascript',
      },
    ],
  })

  // Listen to router route changes to track page_views
  const router = useRouter()
  router.afterEach((to) => {
    const win = window as any
    if (typeof win.gtag === 'function') {
      win.gtag('config', measurementId, {
        page_path: to.fullPath,
      })
    }
  })
})
