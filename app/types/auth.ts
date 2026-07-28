export interface User {
  id: string;
  email: string;
  name: string;
  phone?: string;
  role: "admin" | "student" | "instructor";
  avatar?: string;
  createdAt: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
  remember?: boolean;
}

export interface RegisterData {
  email: string;
  password: string;
  confirmPassword?: string;
  firstName: string;
  lastName: string;
  username?: string;
  phoneNumber?: string;
  dateOfBirth?: string;
  roleId?: number;
}

export interface AuthResponse {
  user: User;
  tokens: AuthTokens;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

// Backend API response types (from user-service)
// Note: All responses are wrapped in { success, message, data }

// Wrapper type for API responses
export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface BackendLoginResponse {
  success: boolean;
  message: string;
  data: {
    // ← add this wrapper
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    user: BackendUserResponse;
  };
}
// Backend user response
export interface BackendUserResponse {
  userId: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  phoneNumber?: string;
  roleId: number;
  role: {
    id: number;
    name: string;
  };
  dateOfBirth: string;
  address?: string;
  image?: string;
  isEmailVerified?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface RegisterResponse {
  user: {
    userId: string;
    username: string;
    email: string;
    firstName: string;
    lastName: string;
    phoneNumber?: string;
    dateOfBirth: string;
    roleId: number;
  };
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}
