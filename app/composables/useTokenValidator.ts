import { useAuthStore } from "~/stores/auth";

interface TokenPayload {
  exp?: number;
  iat?: number;
  sub?: string;
  email?: string;
  role?: string;
}

export const useTokenValidator = () => {
  /**
   * Decode a JWT token and return its payload
   */
  const decodeToken = (token: string): TokenPayload | null => {
    try {
      const parts = token.split(".");
      if (parts.length !== 3) return null;

      const payload = parts[1];
      // URL-safe base64 decode
      const decoded = atob(payload.replace(/-/g, "+").replace(/_/g, "/"));
      return JSON.parse(decoded);
    } catch {
      return null;
    }
  };

  /**
   * Check if a token is expired
   */
  const isTokenExpired = (token: string): boolean => {
    const payload = decodeToken(token);
    if (!payload || !payload.exp) return true;

    const currentTime = Math.floor(Date.now() / 1000);
    return payload.exp < currentTime;
  };

  /**
   * Check if a token will expire soon (within 5 minutes)
   */
  const isTokenExpiringSoon = (token: string, minutesThreshold = 5): boolean => {
    const payload = decodeToken(token);
    if (!payload || !payload.exp) return true;

    const currentTime = Math.floor(Date.now() / 1000);
    const thresholdSeconds = minutesThreshold * 60;
    return payload.exp < currentTime + thresholdSeconds;
  };

  /**
   * Get remaining time until token expiration (in seconds)
   */
  const getTokenRemainingTime = (token: string): number => {
    const payload = decodeToken(token);
    if (!payload || !payload.exp) return 0;

    const currentTime = Math.floor(Date.now() / 1000);
    return Math.max(0, payload.exp - currentTime);
  };

  /**
   * Validate a token and handle auth state if invalid
   * Returns true if valid, false if invalid/expired
   */
  const validateToken = async (token: string): Promise<boolean> => {
    // Check if token exists
    if (!token) {
      return false;
    }

    // Check if token is expired
    if (isTokenExpired(token)) {
      return false;
    }

    // Check if token is structurally valid
    const payload = decodeToken(token);
    if (!payload) {
      return false;
    }

    // Optionally verify with server (for more security)
    // This is commented out as it adds extra API call
    // Uncomment if you need server-side token validation
    /*
    try {
      const authStore = useAuthStore();
      const userData = await authService.fetchCurrentUser(token);
      if (!userData) {
        authStore.clearAuth();
        return false;
      }
    } catch {
      return false;
    }
    */

    return true;
  };

  /**
   * Handle token invalidation with proper redirect
   */
  const handleInvalidToken = (redirectPath: string = "/auth/login") => {
    const authStore = useAuthStore();
    authStore.clearAuth();
  };

  /**
   * Check if current route is admin route
   */
  const isAdminRoute = (path: string): boolean => {
    return path.startsWith("/admin");
  };

  /**
   * Get the appropriate login redirect path based on current route
   */
  const getLoginRedirectPath = (currentPath: string): string => {
    return isAdminRoute(currentPath) ? "/admin/login" : "/auth/login";
  };

  /**
   * Check and handle invalid token, redirecting to appropriate login page
   */
  const checkAndRedirectIfInvalid = async (currentPath: string): Promise<boolean> => {
    const authStore = useAuthStore();
    const token = authStore.accessToken;

    if (!token) {
      handleInvalidToken(getLoginRedirectPath(currentPath));
      return false;
    }

    const isValid = await validateToken(token);
    if (!isValid) {
      handleInvalidToken(getLoginRedirectPath(currentPath));
      return false;
    }

    return true;
  };

  return {
    decodeToken,
    isTokenExpired,
    isTokenExpiringSoon,
    getTokenRemainingTime,
    validateToken,
    handleInvalidToken,
    isAdminRoute,
    getLoginRedirectPath,
    checkAndRedirectIfInvalid,
  };
};