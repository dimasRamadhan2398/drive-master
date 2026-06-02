import { useAuthStore } from "~/stores/auth";

export default defineNuxtRouteMiddleware((to) => {
  // Rehydrate auth state from cookies if needed (handles SSR/client hydration)
  const authToken = useCookie("auth_token");
  const userData = useCookie("user_data");
  const refreshToken = useCookie("refresh_token");

  const authStore = useAuthStore();

  // Rehydrate from cookies if not already authenticated
  if (!authStore.isAuthenticated && authToken.value && userData.value) {
    try {
      const user = JSON.parse(userData.value);
      authStore.setAuth(user, authToken.value, refreshToken.value || undefined);
    } catch {
      // Invalid cookie data, clear them
      authToken.value = null;
      userData.value = null;
    }
  }

  // Handle admin routes
  if (to.path.startsWith("/admin")) {
    // Skip admin login page (handled by guest middleware)
    if (to.path === "/admin/login") {
      return;
    }

    // Check if authenticated
    if (!authStore.isAuthenticated) {
      return navigateTo("/admin/login");
    }

    // Check if user has admin role
    const role = authStore.userRole?.toLowerCase();
    if (role !== "admin") {
      return navigateTo("/");
    }

    return;
  }

  // Public routes that don't require authentication
  const publicRoutes = [
    "/",
    "/auth/login",
    "/auth/register",
    "/auth/verify",
    "/auth/onboarding",
    "/auth/select-plan",
    "/auth/payment",
    "/auth/payment-status",
    "/auth/payment-method",
    "/packages",
    "/instructors",
    "/services",
    "/about",
    "/contact",
    "/blog",
  ];

  // Check if current path is public
  const isPublicRoute = publicRoutes.some(
    (route) => to.path === route || to.path.startsWith(route + "/"),
  );

  // Allow public routes without authentication
  if (isPublicRoute) {
    return;
  }

  // Check if user is authenticated
  if (!authStore.isAuthenticated) {
    return navigateTo("/auth/login");
  }

  // Check if user has admin role - redirect to admin panel
  if (authStore.userRole?.toLowerCase().includes("admin")) {
    return navigateTo("/admin");
  }
});
