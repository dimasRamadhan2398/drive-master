// Extend Nuxt's definePageMeta types to allow string middleware names
declare module '#app/nuxt' {
  interface PageMeta {
    middleware?: Array<NavigationGuard | string> | NavigationGuard | string
  }
}