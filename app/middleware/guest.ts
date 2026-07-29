import { useAuthStore } from '~/stores/auth';

// Middleware to redirect already authenticated users away from login pages
export default defineNuxtRouteMiddleware(async (to) => {
  // Only apply to non-admin auth pages
  const authPages = ['/auth/login', '/auth/register'];
  if (!authPages.includes(to.path)) {
    return;
  }

  // Rehydrate auth state from cookies if needed (handles SSR/client hydration)
  const authToken = useCookie('auth_token');
  const userData = useCookie('user_data');
  const refreshToken = useCookie('refresh_token');

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

  // Only redirect if user is authenticated
  if (!authStore.isAuthenticated) {
    return;
  }

  // Redirect authenticated users based on role
  if (authStore.userRole?.toLowerCase().includes('admin')) {
    return navigateTo('/admin');
  }

  // Check if student has entitlements before deciding where to redirect
  if (!authStore.memberProfile) {
    await authStore.fetchMemberProfile();
  }

  const hasEntitlements =
    (authStore.memberProfile?.entitlements &&
      authStore.memberProfile.entitlements.length > 0) ||
    (authStore.memberProfile?.sessionsCompleted &&
      authStore.memberProfile.sessionsCompleted > 0);

  if (hasEntitlements) {
    return navigateTo('/dashboard');
  }

  // No entitlements or completed sessions — send to onboarding
  return navigateTo('/auth/onboarding');
});