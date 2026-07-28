import type { ApiResponse, PaginatedResponse } from "~/composables/useApiClients";

// ========== Entitlement ==========

export interface Entitlement {
  id: string;
  memberId: string;
  bookingId: string;
  packageId: string;
  packageName: string;
  isNightSession: boolean;
  isWeekendSession: boolean;
  totalSessions: number;
  remaining: number;
  usedSessions: number;
  startDate: string;
  endDate: string | null;
  status: "active" | "completed" | "expired" | "cancelled";
  createdAt: string;
  updatedAt: string;
}

// ========== Member Profile ==========

export interface MemberProfile {
  userId: string;
  sessionsCompleted: number;
  trainingTime: number; // in minutes
  averageRating: number;
  totalAvailableSessions: number;
  entitlements: Entitlement[];
  identityFullname: string;
}

export interface UpdateMemberProfileData {
  identityFullname?: string;
}

// ========== Member List ==========

export interface MemberListItem {
  userId: string;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  phoneNumber: string;
  image: string;
  dateOfBirth: string;
  address: string;
  roleId: number;
  role: {
    id: number;
    name: string;
  };
  memberProfile?: Partial<MemberProfile>;
  createdAt: string;
}

export interface MemberListResponse {
  members: MemberListItem[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// ========== Member Service ==========

export const memberService = {
  /**
   * Get member profile by user ID
   * GET /members/{userId}/profile
   */
  async getMemberProfile(userId: string): Promise<MemberProfile | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<MemberProfile>>(
        `/members/${userId}/profile`,
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  /**
   * Get current user's member profile
   * GET /members/{userId}/profile
   */
  async getMyProfile(accessToken: string): Promise<MemberProfile | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<MemberProfile>>(
        "/members/me/profile",
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  /**
   * Update member profile
   * PUT /members/{userId}/profile
   */
  async updateMemberProfile(
    userId: string,
    data: UpdateMemberProfileData,
  ): Promise<MemberProfile | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<MemberProfile>>(
        `/members/${userId}/profile`,
        {
          method: "PUT",
          body: data,
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  /**
   * Update user details
   * PUT /users/{id}
   */
  async updateUser(
    userId: string,
    data: {
      firstName: string;
      lastName: string;
      phoneNumber?: string;
      address?: string;
      dateOfBirth?: string;
    },
  ): Promise<any> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<any>>(
        `/users/${userId}`,
        {
          method: "PUT",
          body: data,
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  /**
   * Get all members with pagination (admin)
   * GET /members
   */
  async getAllMembers(params: {
    page?: number;
    limit?: number;
    search?: string;
  } = {}): Promise<MemberListResponse> {
    const { user, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.search) queryParams.set("search", params.search);

    const queryString = queryParams.toString();
    const url = `/members${queryString ? `?${queryString}` : ""}`;

    try {
      const response = await user<PaginatedResponse<MemberListItem>>(url, {
        method: "GET",
      });
      const { data, pagination } = extractPaginatedData(response);
      return { members: data, pagination };
    } catch {
      return { members: [], pagination: { page: 1, limit: 10, total: 0, totalPages: 0 } };
    }
  },

  /**
   * Search members with pagination
   * GET /members/search
   */
  async searchMembers(params: {
    page?: number;
    limit?: number;
    search: string;
  }): Promise<MemberListResponse> {
    const { user, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    queryParams.set("search", params.search);

    const url = `/members/search?${queryParams.toString()}`;

    try {
      const response = await user<PaginatedResponse<MemberListItem>>(url, {
        method: "GET",
      });
      const { data, pagination } = extractPaginatedData(response);
      return { members: data, pagination };
    } catch {
      return { members: [], pagination: { page: 1, limit: 10, total: 0, totalPages: 0 } };
    }
  },

  /**
   * Get recent registrations (admin)
   * GET /members/recent
   */
  async getRecentRegistrations(params: {
    limit?: number;
    fromDate?: string;
    toDate?: string;
  } = {}): Promise<MemberListItem[]> {
    const { user, extractData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.fromDate) queryParams.set("fromDate", params.fromDate);
    if (params.toDate) queryParams.set("toDate", params.toDate);

    const queryString = queryParams.toString();
    const url = `/members/recent${queryString ? `?${queryString}` : ""}`;

    try {
      const response = await user<ApiResponse<MemberListItem[]>>(url, {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return [];
    }
  },

  /**
   * Get member entitlements (sessions/packages)
   * GET /members/{userId}/entitlements
   */
  async getMemberEntitlements(userId: string): Promise<Entitlement[]> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<Entitlement[]>>(
        `/members/${userId}/entitlements`,
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return [];
    }
  },

  /**
   * Get current user's entitlements
   */
  async getMyEntitlements(): Promise<Entitlement[]> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<Entitlement[]>>(
        "/members/me/entitlements",
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return [];
    }
  },
};
