import type {
  ApiResponse,
  PaginatedResponse,
} from "~/composables/useApiClients";

export interface UserResponse {
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
}

export interface MemberProfileResponse {
  userId: string;
  sessionsCompleted: number;
  trainingTime: number;
  averageRating: number;
  totalAvailableSessions: number;
}

export interface UserWithProfileResponse extends UserResponse {
  memberProfile?: MemberProfileResponse;
  instructorProfile?: any;
}

export interface Testimonial {
  id: string;
  userId: string;
  userName: string;
  userImage: string;
  userRole: string;
  content: string;
  rating: number;
  tags: string;
  status: string;
  isFeatured: boolean;
  sortOrder: number;
  addedBy: string;
  createdAt: string;
  updatedAt: string;
}

export interface RecentRegistration {
  userId: string;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  phoneNumber: string;
  image: string;
  dateOfBirth: string;
  roleId: number;
  createdAt: string;
}

export interface UserDashboardStats {
  totalUsers: number;
  totalMembers: number;
  totalInstructors: number;
  recentRegistrations: number;
  growthTotalUsers: number;
  growthTotalMembers: number;
  growthTotalInstructors: number;
  growthRecentRegistrations: number;
  activeSessions: number;
  totalSessions: number;
  revenueMTD: number;
  revenueCurrency: string;
  certificatesIssued: number;
  totalCertifications: number;
}

export interface UserDashboardStatsAPI {
  totalUsers: number;
  totalMembers: number;
  totalInstructors: number;
  recentRegistrations: number;
  growthTotalUsers: number;
  growthTotalMembers: number;
  growthTotalInstructors: number;
  growthRecentRegistrations: number;
  activeSessions?: number;
  totalSessions?: number;
  revenueMTD?: number;
  revenueCurrency?: string;
  certificatesIssued?: number;
  totalCertifications?: number;
}

export interface DashboardStats {
  totalUsers: number;
  totalMembers: number;
  totalInstructors: number;
  recentRegistrations: number;
  activeSessions: number;
  totalSessions: number;
  growthRecentRegistrations: number;
  growthTotalUsers: number;
  growthTotalMembers: number;
  growthTotalInstructors: number;
  revenueMTD: number;
  revenueCurrency: string;
  certificatesIssued: number;
  totalCertifications: number;
}

// Helper function to map API response to UserDashboardStats
function mapApiToData(data: UserDashboardStatsAPI): UserDashboardStats {
  return {
    totalUsers: data.totalUsers ?? 0,
    totalMembers: data.totalMembers ?? 0,
    totalInstructors: data.totalInstructors ?? 0,
    recentRegistrations: data.recentRegistrations ?? 0,
    // Multiply growth percentages by 100 since they're stored as decimals
    growthTotalUsers: (data.growthTotalUsers ?? 0) * 100,
    growthTotalMembers: (data.growthTotalMembers ?? 0) * 100,
    growthTotalInstructors: (data.growthTotalInstructors ?? 0) * 100,
    growthRecentRegistrations: (data.growthRecentRegistrations ?? 0) * 100,
    activeSessions: data.activeSessions ?? 0,
    totalSessions: data.totalSessions ?? 0,
    revenueMTD: data.revenueMTD ?? 0,
    revenueCurrency: data.revenueCurrency ?? "IDR",
    certificatesIssued: data.certificatesIssued ?? 0,
    totalCertifications: data.totalCertifications ?? 0,
  };
}

export const dashboardService = {
  // GET /admin/dashboard/stats - Fetch unified dashboard stats (admin only)
  // async fetchDashboardStats(): Promise<DashboardStats> {
  //   const { user, extractData } = useApiClients();

  //   try {
  //     const response = await user<ApiResponse<DashboardStats>>(
  //       "/dashboard/stats",
  //       { method: "GET" },
  //     );
  //     return extractData(response);
  //   } catch {
  //     return {
  //       totalUsers: 0,
  //       totalMembers: 0,
  //       totalInstructors: 0,
  //       recentRegistrations: 0,
  //       activeSessions: 0,
  //       totalSessions: 0,
  //       revenueMTD: 0,
  //       revenueCurrency: "IDR",
  //       certificatesIssued: 0,
  //       totalCertifications: 0,
  //     };
  //   }
  // },

  async fetchUserDashboardStats(): Promise<UserDashboardStats> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<UserDashboardStatsAPI>>(
        "/dashboard/stats",
        { method: "GET" },
      );
      const rawData = extractData(response);
      return mapApiToData(rawData);
    } catch (err) {
      console.error("[dashboardService] fetchUserDashboardStats failed:", err);
      return {
        totalUsers: 0,
        totalMembers: 0,
        totalInstructors: 0,
        recentRegistrations: 0,
        growthRecentRegistrations: 0,
        growthTotalUsers: 0,
        growthTotalMembers: 0,
        growthTotalInstructors: 0,
        activeSessions: 0,
        totalSessions: 0,
        revenueMTD: 0,
        revenueCurrency: "IDR",
        certificatesIssued: 0,
        totalCertifications: 0,
      };
    }
  },

  // GET /members/recent - Fetch recent registrations (admin only)
  async fetchRecentRegistrations(
    params: {
      limit?: number;
      fromDate?: string;
      toDate?: string;
    } = {},
  ): Promise<RecentRegistration[]> {
    const { user, extractData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.fromDate) queryParams.set("fromDate", params.fromDate);
    if (params.toDate) queryParams.set("toDate", params.toDate);

    const queryString = queryParams.toString();
    const url = `/members/recent${queryString ? `?${queryString}` : ""}`;

    try {
      const response = await user<ApiResponse<RecentRegistration[]>>(url, {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return [];
    }
  },

  // GET /users/all - Fetch all users (admin only)
  async fetchAllUsers(): Promise<UserWithProfileResponse[]> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<UserWithProfileResponse[]>>(
        "/users/all",
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return [];
    }
  },

  // GET /members/all - Fetch all members with pagination (admin only)
  async fetchAllMembers(
    params: {
      page?: number;
      limit?: number;
    } = {},
  ): Promise<{ members: UserWithProfileResponse[]; pagination: any }> {
    const { user, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));

    const queryString = queryParams.toString();
    const url = `/members/all${queryString ? `?${queryString}` : ""}`;

    try {
      const response = await user<PaginatedResponse<UserWithProfileResponse>>(
        url,
        { method: "GET" },
      );
      const { data, pagination } = extractPaginatedData(response);
      return { members: data, pagination };
    } catch {
      return { members: [], pagination: null };
    }
  },
};
