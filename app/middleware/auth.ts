import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  console.log("[Auth Middleware] isAuthenticated:", authStore.isAuthenticated);
  console.log("[Auth Middleware] user:", authStore.user);
  console.log("[Auth Middleware] accessToken:", !!authStore.accessToken);

  // Check if user is authenticated
  if (!authStore.isAuthenticated) {
    console.log("[Auth Middleware] Not authenticated, redirecting to login");
    return navigateTo('/auth/login')
  }
})