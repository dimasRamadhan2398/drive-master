import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  console.log("[Guest Middleware] isAuthenticated:", authStore.isAuthenticated);
  console.log("[Guest Middleware] user:", authStore.user);
  console.log("[Guest Middleware] accessToken:", !!authStore.accessToken);

  // Redirect to dashboard if already authenticated
  if (authStore.isAuthenticated) {
    return navigateTo('/dashboard')
  }
})