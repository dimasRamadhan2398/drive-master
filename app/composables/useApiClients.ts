import { useAuthStore } from "~/stores/auth";
import { useTokenValidator } from "./useTokenValidator";

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> {
  success: boolean;
  message: string;
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export const useApiClients = () => {
  const config = useRuntimeConfig();
  const authStore = useAuthStore();
  const tokenValidator = useTokenValidator();

  // Handle 401 errors and redirect to appropriate login page
  const handleUnauthorized = () => {
    const currentPath = useRoute().path;
    tokenValidator.handleInvalidToken(
      tokenValidator.getLoginRedirectPath(currentPath),
    );
  };

  // Create a reactive fetcher that includes auth headers dynamically
  const createFetcher = () => {
    return $fetch.create({
      baseURL: config.public.apiBase as string,
      headers: {
        get Authorization() {
          return authStore.accessToken ? `Bearer ${authStore.accessToken}` : "";
        },
      },
      onResponse({ response }) {
        // Auto-extract data from wrapped response
        if (
          response._data &&
          typeof response._data === "object" &&
          "data" in response._data
        ) {
          return response._data;
        }
        return response._data;
      },
      onResponseError({ response }) {
        // Handle 401 Unauthorized
        if (response.status === 401) {
          handleUnauthorized();
        }
      },
      onRequestError({ error }) {
        // Handle network errors that might indicate auth issues
        if (error?.cause === 401) {
          handleUnauthorized();
        }
      },
    });
  };

  // Extract data from API response wrapper
  const extractData = <T>(response: ApiResponse<T>): T => {
    return response.data;
  };

  // Extract data and pagination from paginated response
  const extractPaginatedData = <T>(response: PaginatedResponse<T>) => {
    return {
      data: response.data,
      pagination: response.pagination,
    };
  };

  return {
    user: createFetcher(),
    core: createFetcher(),
    booking: createFetcher(),
    extractData,
    extractPaginatedData,
  };
};

// Helper function to unwrap API responses
export const unwrapResponse = <T>(response: ApiResponse<T> | any): T => {
  if (response && typeof response === "object" && "data" in response) {
    return response.data;
  }
  return response as T;
};
