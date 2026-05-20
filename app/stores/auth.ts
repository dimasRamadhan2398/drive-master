import { defineStore } from "pinia";
import type {
  LoginCredentials,
  RegisterData,
  BackendLoginResponse,
  BackendUserResponse,
  RegisterResponse,
} from "~/types/auth";

interface AuthStoreState {
  user: BackendUserResponse | null;
  accessToken: string | null;
  refreshTokenValue: string | null;
  isLoading: boolean;
  error: string | null;
}

export const useAuthStore = defineStore("auth", {
  state: (): AuthStoreState => {
    // rehydrate from cookie on init
    const tokenCookie = useCookie("auth_token");
    const refreshCookie = useCookie("refresh_token");
    const userCookie = useCookie("user_data");

    let user: BackendUserResponse | null = null;
    try {
      user = userCookie.value ? JSON.parse(userCookie.value) : null;
    } catch {
      user = null;
    }

    return {
      user,
      accessToken: tokenCookie.value || null,
      refreshTokenValue: refreshCookie.value || null,
      isLoading: false,
      error: null,
    };
  },
  getters: {
    isAuthenticated: (state) => !!state.user && !!state.accessToken,
    currentUser: (state) => state.user,
    userRole: (state) => state.user?.role?.name ?? null,
  },
  actions: {
    // Dedicated action to set authentication state
    setAuth(
      user: BackendUserResponse,
      accessToken: string,
      refreshToken?: string,
    ) {
      this.user = user;
      this.accessToken = accessToken;
      this.refreshTokenValue = refreshToken || null;

      // Store tokens in cookies
      this.setTokenCookie(accessToken);
      if (refreshToken) {
        this.setRefreshTokenCookie(refreshToken);
      }
      // Persist user data for rehydration
      this.setUserCookie(user);

      console.log("[Auth] setAuth completed:", {
        user: user,
        hasToken: !!accessToken,
        isAuthenticated: this.isAuthenticated,
      });
    },

    // Transform backend user to frontend user format
    transformUser(backendUser: BackendUserResponse) {
      return {
        userId: backendUser.userId,
        email: backendUser.email,
        firstName: backendUser.firstName,
        lastName: backendUser.lastName,
        phone: backendUser.phoneNumber,
        role: backendUser.role.name as "admin" | "student" | "instructor",
        createdAt: backendUser.createdAt,
      };
    },

    async login(credentials: LoginCredentials) {
      this.isLoading = true;
      this.error = null;

      try {
        const { user } = useApiClients();
        const response = await user<BackendLoginResponse>("/auth/login", {
          method: "POST",
          body: {
            email: credentials.email,
            password: credentials.password,
          },
        });

        console.log("[Auth] Login response:", response);

        // Use dedicated setAuth action for proper state management
        this.setAuth(
          response.data.user,
          response.data.accessToken,
          response.data.refreshToken,
        );

        return response;
      } catch (error: any) {
        // Extract error message from backend response
        const errorMessage =
          error?.data?.error?.message ||
          error?.data?.message ||
          error?.message ||
          "Login failed";
        this.error = errorMessage;

        // Re-throw with the extracted message so the component can use it
        const enhancedError = new Error(errorMessage);
        throw enhancedError;
      } finally {
        this.isLoading = false;
      }
    },

    async register(data: RegisterData) {
      this.isLoading = true;
      this.error = null;

      try {
        const { user } = useApiClients();
        const response = await user<RegisterResponse>("/auth/register", {
          method: "POST",
          body: data,
        });

        this.user = {
          userId: response.userId,
          username: response.username,
          email: response.email,
          firstName: response.firstName,
          lastName: response.lastName,
          phoneNumber: response.phoneNumber,
          roleId: response.role.id,
          role: {
            id: response.role.id,
            name: response.role.name,
          },
          createdAt: response.createdAt,
          updatedAt: response.createdAt,
          dateOfBirth: response.dateOfBirth,
          isEmailVerified: response.isEmailVerified,
        };

        // Use setAuth for consistency
        this.setAuth(this.user, response.accessToken);

        return response;
      } catch (error: any) {
        this.error = error?.data?.error?.message || "Registration failed";
        throw error;
      } finally {
        this.isLoading = false;
      }
    },

    async logout() {
      try {
        const { user } = useApiClients();
        if (this.accessToken) {
          await user("/auth/logout", { method: "POST" });
        }
      } catch {
        // Continue with logout even if API call fails
      } finally {
        this.clearAuth();
      }
    },

    async refreshToken() {
      try {
        const { user } = useApiClients();
        const refreshTokenValue = this.getRefreshTokenCookie();

        if (!refreshTokenValue) {
          this.clearAuth();
          return false;
        }

        const response = await user<BackendLoginResponse>("/auth/refresh", {
          method: "POST",
          body: {
            refreshToken: refreshTokenValue,
          },
        });

        // Use setAuth for consistent state management
        this.setAuth(
          response.data.user,
          response.data.accessToken,
          response.data.refreshToken,
        );

        return true;
      } catch {
        this.clearAuth();
        return false;
      }
    },

    async fetchCurrentUser() {
      if (!this.accessToken) return null;

      try {
        const { user } = useApiClients();
        const userData = await user<BackendUserResponse>("/auth/me", {
          method: "GET",
        });
        this.user = userData;
        return userData;
      } catch {
        return null;
      }
    },

    setTokenCookie(token: string) {
      const cookie = useCookie("auth_token", {
        maxAge: 60 * 60 * 24 * 7,
        secure: true,
        sameSite: "lax",
      });
      cookie.value = token; // works on both server and client
    },

    getRefreshTokenCookie() {
      if (import.meta.client) {
        const cookie = useCookie("refresh_token");
        return cookie.value;
      }
      return null;
    },

    setRefreshTokenCookie(token: string) {
      const cookie = useCookie("refresh_token", {
        maxAge: 60 * 60 * 24 * 30,
        secure: true,
        sameSite: "lax",
      });
      cookie.value = token;
    },

    setUserCookie(user: BackendUserResponse) {
      const cookie = useCookie("user_data", {
        maxAge: 60 * 60 * 24 * 7, // 7 days, same as access token
        secure: true,
        sameSite: "lax",
      });
      cookie.value = JSON.stringify(user);
    },

    clearAuth() {
      this.user = null;
      this.accessToken = null;
      this.refreshTokenValue = null;
      this.error = null;

      if (import.meta.client) {
        const tokenCookie = useCookie("auth_token");
        const refreshCookie = useCookie("refresh_token");
        const userCookie = useCookie("user_data");
        tokenCookie.value = null;
        refreshCookie.value = null;
        userCookie.value = null;
      }
    },

    clearError() {
      this.error = null;
    },
  },
});
