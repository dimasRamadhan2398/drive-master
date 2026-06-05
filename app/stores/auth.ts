import { defineStore } from "pinia";
import { authService } from "~/services/authService";
import type {
  LoginCredentials,
  RegisterData,
  BackendUserResponse,
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
        const authData = await authService.login(credentials);
        this.setAuth(
          authData.user,
          authData.accessToken,
          authData.refreshToken,
        );
        return authData;
      } catch (error: any) {
        this.error = authService.parseError(error);
        const enhancedError = new Error(this.error);
        throw enhancedError;
      } finally {
        this.isLoading = false;
      }
    },

    async register(data: RegisterData) {
      console.log(
        "[AUTH STORE] register() called with data:",
        JSON.stringify({
          firstName: data.firstName,
          lastName: data.lastName,
          username: data.username,
          email: data.email,
          phoneNumber: data.phoneNumber,
          dateOfBirth: data.dateOfBirth,
          password: data.password ? "***" : undefined,
          roleId: data.roleId,
        }),
      );

      this.isLoading = true;
      this.error = null;

      try {
        console.log("[AUTH STORE] Calling authService.register()...");
        const response = await authService.register(data);
        console.log(
          "[AUTH STORE] authService.register() succeeded, response:",
          JSON.stringify({
            userId: response.user?.userId,
            username: response.user?.username,
            email: response.user?.email,
          }),
        );

        this.user = {
          userId: response.user.userId,
          username: response.user.username,
          email: response.user.email,
          firstName: response.user.firstName,
          lastName: response.user.lastName,
          phoneNumber: response.user.phoneNumber,
          roleId: response.user.roleId,
          role: {
            id: response.user.roleId,
            name: "member",
          },
          dateOfBirth: response.user.dateOfBirth,
        };

        this.setAuth(this.user, response.accessToken, response.refreshToken);
        return response;
      } catch (error: any) {
        console.log(
          "[AUTH STORE] authService.register() failed with error:",
          error,
        );
        this.error = authService.parseError(error);
        console.log("[AUTH STORE] Parsed error message:", this.error);
        throw error;
      } finally {
        console.log(
          "[AUTH STORE] register() completed, isLoading set to false",
        );
        this.isLoading = false;
      }
    },

    async logout() {
      try {
        await authService.logout(this.accessToken);
      } catch {
        // Continue with logout even if API call fails
      } finally {
        this.clearAuth();
      }
    },

    async refreshToken() {
      if (!this.refreshTokenValue) {
        this.clearAuth();
        return false;
      }

      try {
        const authData = await authService.refreshToken(this.refreshTokenValue);
        if (authData) {
          this.setAuth(
            authData.user,
            authData.accessToken,
            authData.refreshToken,
          );
          return true;
        }
        this.clearAuth();
        return false;
      } catch {
        this.clearAuth();
        return false;
      }
    },

    async fetchCurrentUser() {
      if (!this.accessToken) return null;

      try {
        const userData = await authService.fetchCurrentUser(this.accessToken);
        if (userData) {
          this.user = userData;
        }
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
      cookie.value = token;
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
        maxAge: 60 * 60 * 24 * 7,
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
