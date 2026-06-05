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
  const completedSessions = item.memberProfile?.sessionsCompleted || 0;
  const totalSessions =
    item.memberProfile?.totalAvailableSessions || packageSessionMap["8x"] || 8;
  const progress =
    totalSessions > 0
      ? Math.round((completedSessions / totalSessions) * 100)
      : 0;

  // Determine status based on progress
  let status: "active" | "pending" | "completed" = "pending";
  if (completedSessions > 0 && progress < 100) {
    status = "active";
  } else if (progress >= 100) {
    status = "completed";
  }

  return {
    id: hashCode(item.userId),
    name: `${item.firstName} ${item.lastName}`.trim(),
    email: item.email,
    phone: item.phoneNumber || "",
    package: "8x", // Default package
    progress,
    completedSessions,
    totalSessions,
    joinDate: item.dateOfBirth
      ? new Date(item.dateOfBirth).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        })
      : new Date().toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
    status,
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

    const response = await user<PaginatedResponse<StudentApiResponse>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      students: Array.isArray(data) ? data.map(mapApiToStudent) : [],
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
};
