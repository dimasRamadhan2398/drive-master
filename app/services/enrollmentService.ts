import type { ApiResponse } from "~/composables/useApiClients";

// Enrollment Status Types
export type EnrollmentStatus =
  | "pending"
  | "active"
  | "completed"
  | "cancelled"
  | "expired";

// Enrollment interface
export interface Enrollment {
  id: string;
  userId: string;
  userName: string;
  userEmail: string;
  packageId: string;
  packageName: string;
  totalSessions: number;
  completedSessions: number;
  remainingSessions: number;
  status: EnrollmentStatus;
  enrollmentDate: string;
  expiryDate: string;
  expiresAt?: string; // Backend field name
  startDate?: string;
  endDate?: string;
  price: number;
  discountPrice: number;
  totalPrice?: number; // Backend field name
  paymentStatus: "pending" | "paid" | "failed" | "refunded";
  paymentMethod?: string;
  transactionId?: string;
  notes?: string;
  paidAt?: string | null;
  createdAt: string;
  updatedAt: string;
}

// Brief response for lists
export interface EnrollmentBrief {
  id: string;
  packageName: string;
  totalSessions: number;
  completedSessions: number;
  remainingSessions: number;
  status: EnrollmentStatus;
  enrollmentDate: string;
  expiryDate: string;
  paymentStatus: "pending" | "paid" | "failed" | "refunded";
}

// Enrollment list response
export interface EnrollmentListResponse {
  enrollments: EnrollmentBrief[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// Create enrollment data
export interface CreateEnrollmentData {
  userId: string;
  packageId: string;
  price: number;
  discountPrice: number;
  startDate?: string;
  notes?: string;
  addOns?: string[];
}

// Update enrollment data
export interface UpdateEnrollmentData {
  status?: EnrollmentStatus;
  completedSessions?: number;
  startDate?: string;
  endDate?: string;
  notes?: string;
  paymentStatus?: "pending" | "paid" | "failed" | "refunded";
  paymentMethod?: string;
  transactionId?: string;
}

// Enrollment service
export const enrollmentService = {
  // GET /enrollments - Get all enrollments for current user
  async fetchAll(): Promise<EnrollmentBrief[]> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<EnrollmentBrief[]>>(
        "/enrollments",
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return [];
    }
  },

  // GET /enrollments/my - Get current user's enrollments
  async fetchMyEnrollments(): Promise<EnrollmentBrief[]> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<EnrollmentBrief[]>>(
        "/enrollments/my",
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return [];
    }
  },

  // GET /enrollments/active - Get active enrollment
  async fetchActiveEnrollment(): Promise<Enrollment | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<Enrollment>>(
        "/enrollments/active",
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /enrollments - Create a new enrollment
  async create(data: CreateEnrollmentData): Promise<Enrollment | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<any>("/enrollments", {
        method: "POST",
        body: data,
      });
      const dataExtracted = extractData(response);
      if (dataExtracted && typeof dataExtracted === "object" && "enrollment" in dataExtracted) {
        return (dataExtracted as any).enrollment as Enrollment;
      }
      return dataExtracted as Enrollment;
    } catch (error) {
      console.error("[EnrollmentService] Error creating enrollment:", error);
      return null;
    }
  },

  // GET /enrollments/:id - Get enrollment by ID
  async fetchById(id: string): Promise<Enrollment | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<any>(
        `/enrollments/${id}`,
        { method: "GET" },
      );
      const dataExtracted = extractData(response);
      if (dataExtracted && typeof dataExtracted === "object" && "enrollment" in dataExtracted) {
        return (dataExtracted as any).enrollment as Enrollment;
      }
      return dataExtracted as Enrollment;
    } catch {
      return null;
    }
  },

  // PUT /enrollments/:id - Update enrollment
  async update(
    id: string,
    data: UpdateEnrollmentData,
  ): Promise<Enrollment | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<Enrollment>>(
        `/enrollments/${id}`,
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

  // PATCH /enrollments/:id/complete-session - Mark a session as completed
  async completeSession(id: string): Promise<Enrollment | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<Enrollment>>(
        `/enrollments/${id}/complete-session`,
        { method: "PATCH" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // PATCH /enrollments/:id/cancel - Cancel enrollment
  async cancel(id: string): Promise<Enrollment | null> {
    const { user, extractData } = useApiClients();

    try {
      const response = await user<ApiResponse<Enrollment>>(
        `/enrollments/${id}/cancel`,
        { method: "PATCH" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // DELETE /enrollments/:id - Delete enrollment
  async delete(id: string): Promise<boolean> {
    const { user } = useApiClients();

    try {
      await user(`/enrollments/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  // Helper: Get status color for UI
  getStatusColor(
    status: EnrollmentStatus,
  ): "primary" | "info" | "warning" | "success" | "error" | "neutral" {
    const colorMap: Record<
      EnrollmentStatus,
      "primary" | "info" | "warning" | "success" | "error" | "neutral"
    > = {
      pending: "warning",
      active: "success",
      completed: "primary",
      cancelled: "error",
      expired: "neutral",
    };
    return colorMap[status] || "primary";
  },

  // Helper: Get status label for UI
  getStatusLabel(status: EnrollmentStatus): string {
    const labelMap: Record<EnrollmentStatus, string> = {
      pending: "Pending",
      active: "Active",
      completed: "Completed",
      cancelled: "Cancelled",
      expired: "Expired",
    };
    return labelMap[status] || status;
  },

  // Helper: Get payment status color for UI
  getPaymentStatusColor(
    status: "pending" | "paid" | "failed" | "refunded",
  ): "primary" | "info" | "warning" | "success" | "error" | "neutral" {
    const colorMap: Record<
      "pending" | "paid" | "failed" | "refunded",
      "primary" | "info" | "warning" | "success" | "error" | "neutral"
    > = {
      pending: "warning",
      paid: "success",
      failed: "error",
      refunded: "info",
    };
    return colorMap[status] || "primary";
  },

  // Helper: Get payment status label for UI
  getPaymentStatusLabel(
    status: "pending" | "paid" | "failed" | "refunded",
  ): string {
    const labelMap: Record<"pending" | "paid" | "failed" | "refunded", string> =
      {
        pending: "Pending",
        paid: "Paid",
        failed: "Failed",
        refunded: "Refunded",
      };
    return labelMap[status] || status;
  },
};