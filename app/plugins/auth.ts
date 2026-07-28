import { useAuthStore } from "~/stores/auth";

export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore();

  // Restore session from cookie if token exists
  if (import.meta.client) {
    const tokenCookie = useCookie("auth_token");
    if (tokenCookie.value && !authStore.isAuthenticated) {
      try {
        await authStore.fetchCurrentUser();
      } catch {
        // Token invalid, clear it
        authStore.clearAuth();
      }
    }
  }
});
