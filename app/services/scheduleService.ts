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
  time: string;
  endTime: string;
  duration: number; // in minutes
  carId: string;
  carName: string;
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

// Brief response for lists (API response format)
export interface ScheduleBrief {
  id: string;
  date: string;
  time: string;
  endTime: string;
  duration: number;
  carId: string;
  carName: string;
  instructorId: string;
  instructorName: string;
  studentName?: string;
  status: ScheduleStatus;
}

// Raw API response format (from backend)
export interface ScheduleApiResponse {
  id: string;
  date: string;
  time: string;
  duration: number;
  instructorId?: string;
  instructorName?: string;
  carId?: number;
  carName?: string; // API uses "carName" instead of "vehicleName"
  userId?: string;
  userName?: string; // API uses "userName" instead of "studentName"
  bookingId?: string;
  status: ScheduleStatus;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

// Map API response to ScheduleBrief format
export const mapApiToScheduleBrief = (item: ScheduleApiResponse): ScheduleBrief => {
  return {
    id: item.id,
    date: item.date,
    time: item.time,
    endTime: "", // Not provided in API response
    duration: item.duration,
    carId: item.carId?.toString() || "",
    carName: item.carName || "",
    instructorId: item.instructorId || "",
    instructorName: item.instructorName || "",
    studentName: item.userName || undefined, // Map "userName" to "studentName"
    status: item.status,
  };
};

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
  time: string;
  duration: number;
  carId: string;
  instructorId: string;
  notes?: string;
}

export interface UpdateScheduleData {
  date?: string;
  time?: string;
  duration?: number;
  carId?: string;
  instructorId?: string;
  status?: ScheduleStatus;
  notes?: string;
}

export interface BookSlotData {
  userId: string;
  entitlementId: string;
  notes?: string;
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
  status: string;
  createdAt: string;
  updatedAt: string;
}

// Filter params
export interface ScheduleFilterParams {
  page?: number;
  limit?: number;
  date?: string;
  startDate?: string;
  endDate?: string;
  instructorId?: string;
  carId?: string;
  status?: ScheduleStatus;
  studentId?: string;
  sortBy?: "date" | "time" | "createdAt";
  sortOrder?: "asc" | "desc";
}

export interface AvailableScheduleParams {
  date?: string;
  startDate?: string;
  endDate?: string;
  instructorId?: string;
  carId?: string;
  duration?: number;
}

// Helper to format date to YYYY-MM-DD
export const formatDateString = (date: Date): string => {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  return `${year}-${month}-${day}`;
};

export const scheduleService = {
  // ==================== SCHEDULE CRUD METHODS ====================

  // GET /schedules/all - Get all schedules
  async fetchAll(
    params: ScheduleFilterParams = {},
  ): Promise<ScheduleListResponse> {
    const { core, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.date) queryParams.set("date", params.date);
    if (params.startDate) queryParams.set("startDate", params.startDate);
    if (params.endDate) queryParams.set("endDate", params.endDate);
    if (params.instructorId)
      queryParams.set("instructorId", params.instructorId);
    if (params.carId) queryParams.set("carId", params.carId);
    if (params.status) queryParams.set("status", params.status);
    if (params.studentId) queryParams.set("studentId", params.studentId);
    if (params.sortBy) queryParams.set("sortBy", params.sortBy);
    if (params.sortOrder) queryParams.set("sortOrder", params.sortOrder);

    const queryString = queryParams.toString();
    const url = `/schedules/all${queryString ? `?${queryString}` : ""}`;

    const response = await core<PaginatedResponse<ScheduleApiResponse>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      schedules: Array.isArray(data)
        ? data.map<ScheduleBrief>((s: ScheduleApiResponse) => mapApiToScheduleBrief(s))
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

  // POST /schedules/:id/start - Start session
  async startSession(id: string): Promise<Schedule | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Schedule>>(`/schedules/${id}/start`, {
        method: "POST",
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /schedules/:id/complete - Complete session
  async completeSession(id: string): Promise<Schedule | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Schedule>>(
        `/schedules/${id}/complete`,
        {
          method: "POST",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // GET /entitlements/user/:userId/active - Get active entitlements
  async fetchActiveEntitlements(userId: string): Promise<Entitlement[]> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<Entitlement[]>>(
        `/entitlements/user/${userId}/active`,
        { method: "GET" },
      );
      return extractData(response) || [];
    } catch {
      return [];
    }
  },

  // GET /schedules/filter - Get filtered schedules
  async fetchFiltered(
    params: ScheduleFilterParams = {},
  ): Promise<ScheduleListResponse> {
    const { core, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.date) queryParams.set("date", params.date);
    if (params.startDate) queryParams.set("startDate", params.startDate);
    if (params.endDate) queryParams.set("endDate", params.endDate);
    if (params.instructorId)
      queryParams.set("instructorId", params.instructorId);
    if (params.carId) queryParams.set("carId", params.carId);
    if (params.status) queryParams.set("status", params.status);
    if (params.studentId) queryParams.set("studentId", params.studentId);
    if (params.sortBy) queryParams.set("sortBy", params.sortBy);
    if (params.sortOrder) queryParams.set("sortOrder", params.sortOrder);

    const queryString = queryParams.toString();
    const url = `/schedules/filter${queryString ? `?${queryString}` : ""}`;

    const response = await core<PaginatedResponse<ScheduleApiResponse>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      schedules: Array.isArray(data)
        ? data.map<ScheduleBrief>((s: ScheduleApiResponse) => mapApiToScheduleBrief(s))
        : [],
      ...pagination,
    };
  },

  // GET /schedules/filter?date=YYYY-MM-DD - Fetch schedules by date (defaults to today)
  async fetchByDate(date?: string): Promise<ScheduleBrief[]> {
    const { core, extractData } = useApiClients();

    // Default to today if no date provided
    const targetDate = date || formatDateString(new Date());

    try {
      const response = await core<PaginatedResponse<ScheduleApiResponse[]>>(
        `/schedules/filter?date=${targetDate}`,
        { method: "GET" },
      );
      // Handle paginated response format: { data: [...], pagination: {...} }
      const schedulesData = "data" in response && Array.isArray(response.data)
        ? response.data
        : extractData(response);
      // Map API response to ScheduleBrief format
      return Array.isArray(schedulesData) ? schedulesData.map(mapApiToScheduleBrief) : [];
    } catch (err) {
      console.error("Error fetching schedules by date:", err);
      return [];
    }
  },

  // GET /schedules/filter?date=YYYY-MM-DD - Fetch today's schedules (convenience method)
  async todaySessions(): Promise<ScheduleBrief[]> {
    return this.fetchByDate();
  },

  // GET /schedules/available - Get available schedules
  async fetchAvailable(
    params: AvailableScheduleParams = {},
  ): Promise<ScheduleBrief[]> {
    const { core, extractData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.date) queryParams.set("date", params.date);
    if (params.startDate) queryParams.set("startDate", params.startDate);
    if (params.endDate) queryParams.set("endDate", params.endDate);
    if (params.instructorId)
      queryParams.set("instructorId", params.instructorId);
    if (params.carId) queryParams.set("carId", params.carId);
    if (params.duration) queryParams.set("duration", String(params.duration));

    const queryString = queryParams.toString();
    const url = `/schedules/available${queryString ? `?${queryString}` : ""}`;

    try {
      const response = await core<ApiResponse<ScheduleApiResponse[]>>(url, {
        method: "GET",
      });
      const data = extractData(response);
      // Map API response to ScheduleBrief format
      return Array.isArray(data) ? data.map(mapApiToScheduleBrief) : [];
    } catch (err) {
      console.error("Error fetching available schedules:", err);
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
