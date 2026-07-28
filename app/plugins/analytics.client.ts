import { defineNuxtPlugin, useRuntimeConfig } from "#app"

/**
 * analytics.client.ts
 *
 * NOTE: `nuxt-gtag` v4 (configured in nuxt.config.ts) already handles:
 *  - Loading the gtag script
 *  - Sending the initial page_view
 *  - Tracking SPA route changes internally
 *
 * Previously, this plugin was manually injecting a second gtag script,
 * which CONFLICTED with nuxt-gtag and prevented data from appearing in GA.
 *
 * This plugin now only validates the config and logs a warning if the
 * measurement ID is missing. No manual script injection is needed.
 *
 * For custom event tracking, use the `useGtag()` composable from nuxt-gtag
 * directly in your components/pages.
 */
export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const measurementId = config.public.gaMeasurementId as string | undefined

  if (!measurementId) {
    console.warn("[Analytics] NUXT_PUBLIC_GA_MEASUREMENT_ID is not set. Google Analytics will not track data.")
  }
})
