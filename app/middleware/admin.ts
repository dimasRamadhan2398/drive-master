import { useAuthStore } from "~/stores/auth";
import { useAuth } from "~/composables/useAuth";

export default defineNuxtRouteMiddleware(() => {
  const authStore = useAuthStore();
  const { isAdminLoggedIn } = useAuth();

  console.log("[Admin Middleware] Checking admin access");
  console.log("[Admin Middleware] userRole:", authStore.userRole);
  console.log("[Admin Middleware] isAdminLoggedIn:", isAdminLoggedIn.value);

  // Check if admin is logged in via simple auth composable (used by admin login)
  if (isAdminLoggedIn.value) {
    console.log("[Admin Middleware] Admin access granted via useAuth");
    return;
  }

  // Check if user is authenticated via Pinia store
  if (!authStore.isAuthenticated) {
    console.log("[Admin Middleware] Not authenticated, redirecting to admin login");
    return navigateTo("/admin/login");
  }

  // Check if user has admin role
  const role = authStore.userRole?.toLowerCase();
  if (role !== "admin") {
    console.log("[Admin Middleware] Not admin, redirecting to home");
    return navigateTo("/");
  }

  console.log("[Admin Middleware] Admin access granted");
});