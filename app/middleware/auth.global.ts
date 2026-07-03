import { useAuthStore } from "~/stores/auth";
import { useTokenValidator } from "~/composables/useTokenValidator";

export default defineNuxtRouteMiddleware(async (to: any) => {
  // Rehydrate auth state from cookies if needed (handles SSR/client hydration)
  const authToken = useCookie("auth_token");
  const userData = useCookie("user_data");
  const refreshToken = useCookie("refresh_token");

  const authStore = useAuthStore();
  const tokenValidator = useTokenValidator();

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

  // Public routes that don't require authentication
  const publicRoutes = [
    "/",
    "/auth/login",
    "/auth/register",
    "/auth/verify",
    "/auth/forgot-password",
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

  // Check if token is expired or invalid - redirect to appropriate login page
  if (authStore.accessToken) {
    const loginRedirect = tokenValidator.getLoginRedirectPath(to.path);

    if (tokenValidator.isTokenExpired(authStore.accessToken)) {
      // Token is expired - clear auth and redirect only if on a protected route
      tokenValidator.handleInvalidToken(loginRedirect);
      if (!isPublicRoute) {
        return navigateTo(loginRedirect);
      }
    }
  }

  // Handle admin routes
  if (to.path.startsWith("/admin")) {
    // If accessing admin login
    if (to.path === "/admin/login") {
      // Redirect authenticated admin users away from login
      if (authStore.isAuthenticated && authStore.userRole?.toLowerCase().includes("admin")) {
        return navigateTo("/admin");
      }
      // Allow unauthenticated users to access login
      return;
    }

    // Check if authenticated
    if (!authStore.isAuthenticated) {
      return navigateTo("/admin/login");
    }

    // Check if user has admin role
    const role = authStore.userRole?.toLowerCase() || "";
    if (!role.includes("admin")) {
      return navigateTo("/");
    }

    return;
  }

  // If user is authenticated as admin and trying to access auth pages, redirect to admin
  if (
    authStore.isAuthenticated &&
    authStore.userRole?.toLowerCase().includes("admin")
  ) {
    if (to.path.startsWith("/auth/") || to.path === "/") {
      return navigateTo("/admin");
    }
  }

  // Allow public routes without authentication
  if (isPublicRoute) {
    return;
  }

  // Check if user is authenticated
  if (!authStore.isAuthenticated) {
    return navigateTo("/auth/login");
  }

  // Check if accessing dashboard routes - verify member has entitlements
  if (to.path.startsWith("/dashboard")) {
    // Skip entitlement check if user just completed payment — entitlements
    // may not yet be provisioned due to async Kafka processing.
    const justPaid = to.query.just_paid === "true";
    if (justPaid) {
      return;
    }

    // Fetch member profile if not already loaded
    if (!authStore.memberProfile) {
      await authStore.fetchMemberProfile();
    }

    // Check if member has entitlements (purchased a package)
    const hasEntitlements =
      authStore.memberProfile?.entitlements &&
      authStore.memberProfile.entitlements.length > 0;

    // Redirect to onboarding if no entitlements
    if (!hasEntitlements) {
      return navigateTo("/auth/onboarding");
    }
  }
});
