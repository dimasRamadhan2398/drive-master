import type { Student } from "~/stores/students";
import type {
  ApiResponse,
  PaginatedResponse,
} from "~/composables/useApiClients";

export interface PaginationParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
}

export interface PaginatedStudentsResult {
  students: Student[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface MemberProfile {
  userId: string;
  sessionsCompleted: number;
  trainingTime: number;
  averageRating: number;
  totalAvailableSessions: number;
  entitlements: Entitlement[];
}

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
  status: "active" | "completed" | "cancelled";
  createdAt: string;
  updatedAt: string;
}

export interface StudentApiResponse {
  userId: string;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  phoneNumber?: string;
  image?: string;
  dateOfBirth?: string;
  address?: string;
  roleId: number;
  role: {
    id: number;
    name: string;
  };
  memberProfile?: MemberProfile;
}

export interface StudentWithStudent extends StudentApiResponse {
  student: StudentApiResponse;
}

export interface CreateStudentData {
  email: string;
  firstName: string;
  lastName: string;
  phoneNumber: string;
  password: string;
  dateOfBirth?: string;
  address?: string;
}

export interface UpdateStudentData {
  firstName?: string;
  lastName?: string;
  phoneNumber?: string;
  address?: string;
  dateOfBirth?: string;
}

const packageSessionMap: { [key: string]: number } = {
  "6x": 6,
  "8x": 8,
  "10x": 10,
  "12x": 12,
};

export const mapApiToStudent = (item: StudentApiResponse): Student => {
  // Get entitlements from memberProfile
  const entitlements = item.memberProfile?.entitlements || [];

  // Find the active entitlement (most recent active package)
  const activeEntitlement = entitlements.find((e) => e.status === "active");

  // Calculate total sessions across all entitlements
  const totalSessions = entitlements.reduce((sum, e) => sum + e.totalSessions, 0);

  // Calculate total completed sessions across all entitlements
  const completedSessions = entitlements.reduce((sum, e) => sum + e.usedSessions, 0);

  // Calculate progress based on completed sessions vs total sessions
  const progress = totalSessions > 0
    ? Math.round((completedSessions / totalSessions) * 100)
    : 0;

  // Determine overall status based on entitlements
  let status: "active" | "pending" | "completed" = "pending";
  const hasActiveEntitlement = entitlements.some((e) => e.status === "active");
  const hasCompletedEntitlement = entitlements.some((e) => e.status === "completed");

  if (hasActiveEntitlement) {
    status = "active";
  } else if (hasCompletedEntitlement || completedSessions > 0) {
    status = "completed";
  }

  // Get the primary package name from the most recent active entitlement
  const packageName = activeEntitlement?.packageName || "No Package";

  // Get join date from the first entitlement or member creation
  const joinDate = activeEntitlement?.startDate
    ? new Date(activeEntitlement.startDate).toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric",
      })
    : new Date().toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric",
      });

  return {
    id: item.userId,
    name: `${item.firstName} ${item.lastName}`.trim(),
    email: item.email,
    phone: item.phoneNumber || "",
    package: packageName,
    progress,
    completedSessions,
    totalSessions,
    joinDate,
    status,
    entitlements,
  };
};

// Generate numeric ID from UUID
function hashCode(str: string): number {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash;
  }
  return Math.abs(hash);
}

export const studentService = {
  // Fetch with pagination - returns students list with pagination info
  async fetchAll(
    params: PaginationParams = {},
  ): Promise<PaginatedStudentsResult> {
    const { user, extractPaginatedData } = useApiClients();

    // Build query params
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.search) queryParams.set("search", params.search);
    if (params.status && params.status !== "all")
      queryParams.set("status", params.status);

    const queryString = queryParams.toString();
    const url = `/members/all${queryString ? `?${queryString}` : ""}`;

    const response = await user<PaginatedResponse<StudentWithStudent>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);

    // Handle both nested and flat response formats
    const students = Array.isArray(data)
      ? data.map((item) => {
          // Handle nested "student" wrapper
          const studentData = "student" in item ? item.student : item;
          return mapApiToStudent(studentData);
        })
      : [];

    return {
      students,
      pagination,
    };
  },

  // Fetch all without pagination - returns flat array (for backward compatibility)
  async fetchAllFlat(): Promise<Student[]> {
    const result = await this.fetchAll();
    return result.students;
  },

  async fetchById(userId: string): Promise<Student | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<StudentApiResponse>>(
        `/members/${userId}`,
        {
          method: "GET",
        },
      );
      return mapApiToStudent(extractData(response));
    } catch {
      return null;
    }
  },

  async create(data: CreateStudentData): Promise<Student> {
    const { user, extractData } = useApiClients();
    const response = await user<ApiResponse<StudentApiResponse>>("/members", {
      method: "POST",
      body: data,
    });
    return mapApiToStudent(extractData(response));
  },

  async update(
    userId: string,
    data: UpdateStudentData,
  ): Promise<Student | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<StudentApiResponse>>(
        `/members/${userId}`,
        {
          method: "PUT",
          body: data,
        },
      );
      return mapApiToStudent(extractData(response));
    } catch {
      return null;
    }
  },

  async delete(userId: string): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user(`/members/${userId}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  async updateProgress(
    userId: string,
    completedSessions: number,
  ): Promise<Student | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<StudentApiResponse>>(
        `/members/${userId}/progress`,
        {
          method: "PATCH",
          body: { sessionsCompleted: completedSessions },
        },
      );
      return mapApiToStudent(extractData(response));
    } catch {
      return null;
    }
  },

  // Utility functions
  getTotalSessions(packageType: string): number {
    return packageSessionMap[packageType] || 0;
  },

  calculateProgress(completedSessions: number, totalSessions: number): number {
    if (totalSessions === 0) return 0;
    return Math.round((completedSessions / totalSessions) * 100);
  },

  formatJoinDate(date?: Date): string {
    return (date || new Date()).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  },

  async searchStudents(query: string): Promise<PaginatedStudentsResult> {
    const { user, extractPaginatedData } = useApiClients();
    const lowerQuery = query.toLowerCase();
    try {
      const response = await user<PaginatedResponse<StudentWithStudent>>(
        `/members/search?page=1&limit=10&search=${encodeURIComponent(query)}`,
        {
          method: "GET",
        },
      );

      const { data, pagination } = extractPaginatedData(response);

      // Handle both nested and flat response formats
      const students = Array.isArray(data)
        ? data.map((item) => {
            const studentData = "student" in item ? item.student : item;
            return mapApiToStudent(studentData);
          })
        : [];

      return {
        students,
        pagination,
      };
    } catch {
      return {
        students: [],
        pagination: {
          page: 1,
          limit: 10,
          total: 0,
          totalPages: 0,
        },
      };
    }
  },
};
