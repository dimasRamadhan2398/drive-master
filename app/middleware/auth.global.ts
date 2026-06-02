import { useAuthStore } from "~/stores/auth";

export default defineNuxtRouteMiddleware((to) => {
  // Check if we're on the admin login page - skip this middleware for admin routes
  if (to.path.startsWith("/admin")) {
    return;
  }

  const authStore = useAuthStore();

  console.log("[Auth Middleware] isAuthenticated:", authStore.isAuthenticated);
  console.log("[Auth Middleware] user:", authStore.user);
  console.log("[Auth Middleware] accessToken:", !!authStore.accessToken);

  // Check if user is authenticated
  if (!authStore.isAuthenticated) {
    console.log("[Auth Middleware] Not authenticated, redirecting to login");
    return navigateTo("/auth/login");
  }

  // Check if user has admin role - redirect to admin panel
  if (authStore.userRole?.toLowerCase().includes("admin")) {
    console.log(
      "[Auth Middleware] User role is admin, redirect to admin panel",
    );
    return navigateTo("/admin");
  }

  // If already on login page and authenticated, redirect to dashboard
  if (to.path === "/auth/login") {
    console.log(
      "[Auth Middleware] Already authenticated, redirecting to dashboard",
    );
    return navigateTo("/dashboard");
  }
});
