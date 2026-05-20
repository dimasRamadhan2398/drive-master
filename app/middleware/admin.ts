import { useAuthStore } from "~/stores/auth";

export default defineNuxtRouteMiddleware(() => {
  const authStore = useAuthStore();

  console.log("[Admin Middleware] Checking admin access");
  console.log("[Admin Middleware] userRole:", authStore.userRole);

  // First check if user is authenticated (reuse auth middleware logic)
  if (!authStore.isAuthenticated) {
    console.log("[Admin Middleware] Not authenticated, redirecting to login");
    return navigateTo("/auth/login");
  }

  // Check if user has admin role
  const role = authStore.userRole?.toLowerCase();
  if (role !== "admin") {
    console.log("[Admin Middleware] Not admin, redirecting to home");
    return navigateTo("/");
  }

  console.log("[Admin Middleware] Admin access granted");
});