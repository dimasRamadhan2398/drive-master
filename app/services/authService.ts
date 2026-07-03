import type {
  LoginCredentials,
  RegisterData,
  BackendLoginResponse,
  BackendUserResponse,
  RegisterResponse,
} from "~/types/auth";
import { useApiClients, type ApiResponse } from "~/composables/useApiClients";

export const authService = {
  async login(
    credentials: LoginCredentials,
  ): Promise<BackendLoginResponse["data"]> {
    const { user, extractData } = useApiClients();
    const response = await user<ApiResponse<BackendLoginResponse["data"]>>(
      "/auth/login",
      {
        method: "POST",
        body: {
          email: credentials.email,
          password: credentials.password,
        },
      },
    );
    return extractData(response);
  },

  async register(data: RegisterData): Promise<RegisterResponse> {
    console.log("[AUTH SERVICE] register() called with:", {
      firstName: data.firstName,
      lastName: data.lastName,
      username: data.username,
      email: data.email,
      phoneNumber: data.phoneNumber,
      dateOfBirth: data.dateOfBirth,
      password: "***",
      confirmPassword: data.confirmPassword ? "***" : undefined,
      roleId: data.roleId,
    });

    const { user, extractData } = useApiClients();

    // Build body with confirmPassword if provided
    const body: any = { ...data };
    if (data.confirmPassword) {
      body.confirmPassword = data.confirmPassword;
    }

    console.log("[AUTH SERVICE] Sending request to /auth/register...");
    const response = await user<ApiResponse<RegisterResponse>>(
      "/auth/register",
      {
        method: "POST",
        body,
      },
    );
    console.log("[AUTH SERVICE] Response received");
    return extractData(response);
  },

  async logout(accessToken?: string | null): Promise<void> {
    if (!accessToken) return;

    const { user } = useApiClients();
    await user("/auth/logout", { method: "POST" });
  },

  async refreshToken(
    refreshToken: string,
  ): Promise<BackendLoginResponse["data"] | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<BackendLoginResponse["data"]>>(
        "/auth/refresh",
        {
          method: "POST",
          body: { refreshToken },
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  async fetchUserById(
    userId: string,
  ): Promise<BackendUserResponse | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<BackendUserResponse>>(
        `/users/${userId}`,
        {
          method: "GET",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  async verifyEmail(token: string): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user("/auth/verify", {
        method: "POST",
        body: { token },
      });
      return true;
    } catch {
      return false;
    }
  },

  async forgotPassword(email: string): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user("/auth/forgot-password", {
        method: "POST",
        body: { email },
      });
      return true;
    } catch {
      return false;
    }
  },

  async resetPassword(token: string, newPassword: string): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user("/auth/reset-password", {
        method: "POST",
        body: { token, password: newPassword },
      });
      return true;
    } catch {
      return false;
    }
  },

  async changePassword(
    accessToken: string,
    oldPassword: string,
    newPassword: string,
  ): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user("/auth/change-password", {
        method: "POST",
        body: { oldPassword, newPassword },
      });
      return true;
    } catch {
      return false;
    }
  },

  // Helper to parse error from API response
  parseError(error: any): string {
    return (
      error?.data?.error?.message ||
      error?.data?.message ||
      error?.message ||
      "An error occurred"
    );
  },
};
