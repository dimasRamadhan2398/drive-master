import type {
  ApiResponse,
  PaginatedResponse,
} from "~/composables/useApiClients";

// Schedule Status Types
export type ScheduleStatus =
  | "available"
  | "booked"
  | "in-progress"
  | "completed"
  | "blocked";

// Schedule Types (mirroring backend DTOs)
export interface Schedule {
  id: string;
  date: string;
  startTime: string;
  endTime: string;
  duration: number; // in minutes
  vehicleId: string;
  vehicleName: string;
  instructorId: string;
  instructorName: string;
  studentId?: string;
  studentName?: string;
  status: ScheduleStatus;
  bookingId?: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

// Brief response for lists
export interface ScheduleBrief {
  id: string;
  date: string;
  startTime: string;
  endTime: string;
  duration: number;
  vehicleName: string;
  instructorName: string;
  studentName?: string;
  status: ScheduleStatus;
}

export interface ScheduleListResponse {
  schedules: ScheduleBrief[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// Request DTOs (matching backend)
export interface CreateScheduleData {
  date: string;
  startTime: string;
  duration: number;
  vehicleId: string;
  instructorId: string;
  notes?: string;
}

export interface UpdateScheduleData {
  date?: string;
  startTime?: string;
  duration?: number;
  vehicleId?: string;
  instructorId?: string;
  status?: ScheduleStatus;
  notes?: string;
}

export interface BookSlotData {
  studentId: string;
  notes?: string;
}

// Filter params
export interface ScheduleFilterParams {
  page?: number;
  limit?: number;
  date?: string;
  startDate?: string;
  endDate?: string;
  instructorId?: string;
  vehicleId?: string;
  status?: ScheduleStatus;
  studentId?: string;
  sortBy?: "date" | "startTime" | "createdAt";
  sortOrder?: "asc" | "desc";
}

export interface AvailableScheduleParams {
  date?: string;
  startDate?: string;
  endDate?: string;
  instructorId?: string;
  vehicleId?: string;
  duration?: number;
}

export const scheduleService = {
  // ==================== SCHEDULE CRUD METHODS ====================

  // GET /schedules/all - Get all schedules
  async fetchAll(
    params: ScheduleFilterParams = {},
  ): Promise<ScheduleListResponse> {
    const { user, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.date) queryParams.set("date", params.date);
    if (params.startDate) queryParams.set("startDate", params.startDate);
    if (params.endDate) queryParams.set("endDate", params.endDate);
    if (params.instructorId)
      queryParams.set("instructorId", params.instructorId);
    if (params.vehicleId) queryParams.set("vehicleId", params.vehicleId);
    if (params.status) queryParams.set("status", params.status);
    if (params.studentId) queryParams.set("studentId", params.studentId);
    if (params.sortBy) queryParams.set("sortBy", params.sortBy);
    if (params.sortOrder) queryParams.set("sortOrder", params.sortOrder);

    const queryString = queryParams.toString();
    const url = `/schedules/all${queryString ? `?${queryString}` : ""}`;

    const response = await user<PaginatedResponse<Schedule>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      schedules: Array.isArray(data)
        ? data.map<ScheduleBrief>((s: Schedule) => ({
            id: s.id,
            date: s.date,
            startTime: s.startTime,
            endTime: s.endTime,
            duration: s.duration,
            vehicleName: s.vehicleName,
            instructorName: s.instructorName,
            studentName: s.studentName,
            status: s.status,
          }))
        : [],
      ...pagination,
    };
  },

  // POST /schedules/create - Create a new schedule
  async create(data: CreateScheduleData): Promise<Schedule | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Schedule>>("/schedules/create", {
        method: "POST",
        body: data,
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // GET /schedules/filter - Get filtered schedules
  async fetchFiltered(
    params: ScheduleFilterParams = {},
  ): Promise<ScheduleListResponse> {
    const { user, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.date) queryParams.set("date", params.date);
    if (params.startDate) queryParams.set("startDate", params.startDate);
    if (params.endDate) queryParams.set("endDate", params.endDate);
    if (params.instructorId)
      queryParams.set("instructorId", params.instructorId);
    if (params.vehicleId) queryParams.set("vehicleId", params.vehicleId);
    if (params.status) queryParams.set("status", params.status);
    if (params.studentId) queryParams.set("studentId", params.studentId);
    if (params.sortBy) queryParams.set("sortBy", params.sortBy);
    if (params.sortOrder) queryParams.set("sortOrder", params.sortOrder);

    const queryString = queryParams.toString();
    const url = `/schedules/filter${queryString ? `?${queryString}` : ""}`;

    const response = await user<PaginatedResponse<Schedule>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      schedules: Array.isArray(data)
        ? data.map<ScheduleBrief>((s: Schedule) => ({
            id: s.id,
            date: s.date,
            startTime: s.startTime,
            endTime: s.endTime,
            duration: s.duration,
            vehicleName: s.vehicleName,
            instructorName: s.instructorName,
            studentName: s.studentName,
            status: s.status,
          }))
        : [],
      ...pagination,
    };
  },

  // GET /schedules/available - Get available schedules
  async fetchAvailable(
    params: AvailableScheduleParams = {},
  ): Promise<ScheduleBrief[]> {
    const { user, extractData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.date) queryParams.set("date", params.date);
    if (params.startDate) queryParams.set("startDate", params.startDate);
    if (params.endDate) queryParams.set("endDate", params.endDate);
    if (params.instructorId)
      queryParams.set("instructorId", params.instructorId);
    if (params.vehicleId) queryParams.set("vehicleId", params.vehicleId);
    if (params.duration) queryParams.set("duration", String(params.duration));

    const queryString = queryParams.toString();
    const url = `/schedules/available${queryString ? `?${queryString}` : ""}`;

    try {
      const response = await user<ApiResponse<ScheduleBrief[]>>(url, {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return [];
    }
  },

  // GET /schedules/:id - Get schedule by ID
  async fetchById(id: string): Promise<Schedule | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Schedule>>(`/schedules/${id}`, {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // PUT /schedules/:id - Update schedule
  async update(id: string, data: UpdateScheduleData): Promise<Schedule | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Schedule>>(`/schedules/${id}`, {
        method: "PUT",
        body: data,
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // DELETE /schedules/:id - Delete schedule
  async delete(id: string): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user(`/schedules/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  // ==================== BOOKING METHODS ====================

  // POST /schedules/:id/book - Book a slot
  async bookSlot(id: string, data: BookSlotData): Promise<Schedule | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Schedule>>(
        `/schedules/${id}/book`,
        {
          method: "POST",
          body: data,
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /schedules/:id/cancel - Cancel booking
  async cancelBooking(id: string): Promise<Schedule | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Schedule>>(
        `/schedules/${id}/cancel`,
        {
          method: "POST",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // ==================== UTILITY METHODS ====================

  // Calculate end time from start time and duration
  calculateEndTime(startTime: string, duration: number): string {
    const [hours, minutes] = startTime.split(":").map(Number);
    const totalMinutes = hours! * 60 + minutes! + duration;
    const endHours = Math.floor(totalMinutes / 60) % 24;
    const endMinutes = totalMinutes % 60;
    return `${String(endHours).padStart(2, "0")}:${String(endMinutes).padStart(2, "0")}`;
  },

  // Format duration for display
  formatDuration(duration: number): string {
    if (duration >= 60) {
      const hours = Math.floor(duration / 60);
      const mins = duration % 60;
      return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
    }
    return `${duration} min`;
  },

  // Get status color for UI
  getStatusColor(
    status: ScheduleStatus,
  ): "primary" | "info" | "warning" | "success" | "error" | "neutral" {
    const colorMap: Record<
      ScheduleStatus,
      "primary" | "info" | "warning" | "success" | "error" | "neutral"
    > = {
      available: "primary",
      booked: "info",
      "in-progress": "warning",
      completed: "neutral",
      blocked: "error",
    };
    return colorMap[status] || "primary";
  },

  // Get status label for UI
  getStatusLabel(status: ScheduleStatus): string {
    const labelMap: Record<ScheduleStatus, string> = {
      available: "Available",
      booked: "Booked",
      "in-progress": "In Progress",
      completed: "Completed",
      blocked: "Blocked",
    };
    return labelMap[status] || status;
  },
};
