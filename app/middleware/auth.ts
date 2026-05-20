import { useAuthStore } from "~/stores/auth";

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore();

  console.log("[Auth Middleware] isAuthenticated:", authStore.isAuthenticated);
  console.log("[Auth Middleware] user:", authStore.user);
  console.log("[Auth Middleware] accessToken:", !!authStore.accessToken);

  // Check if user is authenticated
  if (!authStore.isAuthenticated) {
    console.log("[Auth Middleware] Not authenticated, redirecting to login");
    return navigateTo("/auth/login");
  }

  if (to.path === "/auth/login") {
    console.log(
      "[Auth Middleware] Already authenticated, redirecting to dashboard",
    );
    return navigateTo("/dashboard");
  }

  if (authStore.userRole?.includes("admin")) {
    console.log(
      "[Auth Middleware] User role is admin, redirect to admin panel",
    );
    return navigateTo("/admin");
  }
});
