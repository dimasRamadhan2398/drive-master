import { defineStore } from "pinia";
import { scheduleService } from "~/services/scheduleService";
import type {
  Schedule,
  ScheduleBrief,
  ScheduleStatus,
  CreateScheduleData,
  UpdateScheduleData,
  ScheduleFilterParams,
  AvailableScheduleParams,
} from "~/services/scheduleService";

// UI-specific slot interface (for component use)
export interface ScheduleSlot {
  id: string;
  date: string;
  time: string;
  duration: string;
  car: string;
  instructor: string;
  student: string | null;
  status: ScheduleStatus;
}

export interface ScheduleFormFullData {
  date?: string;
  startTime?: string;
  duration?: number;
  vehicleId?: string;
  instructorId?: string;
  vehicleName?: string;
  instructorName?: string;
  status?: ScheduleStatus;
  notes?: string;
}

// Map API Schedule to UI ScheduleSlot
export const mapScheduleToSlot = (
  schedule: Schedule | ScheduleBrief,
): ScheduleSlot => {
  return {
    id: schedule.id,
    date: schedule.date,
    time: schedule.startTime,
    duration:
      schedule.duration >= 60
        ? `${Math.floor(schedule.duration / 60)}h ${schedule.duration % 60 > 0 ? `${schedule.duration % 60}m` : ""}`.trim()
        : `${schedule.duration} min`,
    car: "vehicleName" in schedule ? schedule.vehicleName : "",
    instructor: "instructorName" in schedule ? schedule.instructorName : "",
    student: "studentName" in schedule ? schedule.studentName || null : null,
    status: schedule.status,
  };
};

// Map UI ScheduleSlot to API format
export const mapSlotToCreateData = (slot: ScheduleSlot): CreateScheduleData => {
  const durationNum = parseInt(slot.duration.replace(/[^0-9]/g, "")) || 60;
  return {
    date: slot.date,
    startTime: slot.time,
    duration: durationNum,
    vehicleId: "", // Will be resolved by UI layer
    instructorId: "", // Will be resolved by UI layer
  };
};

interface SchedulesState {
  slots: ScheduleSlot[];
  isLoading: boolean;
  error: string | null;
  isInitialized: boolean;
  // Filter state
  selectedDate: string;
  filterInstructor: string;
  filterVehicle: string;
  // Pagination
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// Helper to format date to YYYY-MM-DD
const formatDateString = (date: Date): string => {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  return `${year}-${month}-${day}`;
};

// Default sample slots for fallback (uses today's date)
const getInitialSlots = (): ScheduleSlot[] => {
  const today = formatDateString(new Date());
  return [
    {
      id: "1",
      date: today,
      time: "08:00",
      duration: "60 min",
      car: "BYD Atto 1",
      instructor: "Mr. Ahmad",
      student: null,
      status: "available",
    },
    {
      id: "2",
      date: today,
      time: "09:00",
      duration: "60 min",
      car: "BYD Atto 1",
      instructor: "Mr. Ahmad",
      student: "John Doe",
      status: "booked",
    },
    {
      id: "3",
      date: today,
      time: "10:00",
      duration: "60 min",
      car: "BYD Atto 1",
      instructor: "Ms. Sari",
      student: "Sarah Putri",
      status: "in-progress",
    },
    {
      id: "4",
      date: today,
      time: "14:00",
      duration: "60 min",
      car: "BYD Atto 1",
      instructor: "Mr. Budi",
      student: null,
      status: "blocked",
    },
  ];
};

// Initialize with empty array (will be populated from API)
const initialSlots: ScheduleSlot[] = [];

// Helper to get initial slots with today's date
const getInitialSlotsForFallback = (): ScheduleSlot[] => getInitialSlots();

export const useSchedulesStore = defineStore("schedules", {
  state: (): SchedulesState => ({
    slots: [],
    isLoading: false,
    error: null,
    isInitialized: false,
    selectedDate: formatDateString(new Date()),
    filterInstructor: "All Instructors",
    filterVehicle: "All Vehicles",
    pagination: {
      page: 1,
      limit: 50,
      total: 0,
      totalPages: 0,
    },
  }),

  getters: {
    // Get filtered slots based on date and filters
    filteredSlots: (state): ScheduleSlot[] => {
      return state.slots.filter((slot) => {
        const matchDate = slot.date === state.selectedDate;
        const matchInst =
          state.filterInstructor === "All Instructors" ||
          slot.instructor === state.filterInstructor;
        const matchVeh =
          state.filterVehicle === "All Vehicles" ||
          slot.car === state.filterVehicle;
        return matchDate && matchInst && matchVeh;
      });
    },

    // Get slots by status
    availableSlots: (state): ScheduleSlot[] =>
      state.slots.filter((s) => s.status === "available"),
    bookedSlots: (state): ScheduleSlot[] =>
      state.slots.filter((s) => s.status === "booked"),
    inProgressSlots: (state): ScheduleSlot[] =>
      state.slots.filter((s) => s.status === "in-progress"),
    completedSlots: (state): ScheduleSlot[] =>
      state.slots.filter((s) => s.status === "completed"),
    blockedSlots: (state): ScheduleSlot[] =>
      state.slots.filter((s) => s.status === "blocked"),

    // Today's sessions - slots for today
    todaySessions: (state): ScheduleSlot[] => {
      const today = formatDateString(new Date());
      return state.slots.filter((slot) => slot.date === today);
    },

    // Stats for current view
    currentStats: (state) => {
      const filtered = state.slots.filter(
        (slot) => slot.date === state.selectedDate,
      );
      return {
        available: filtered.filter((s) => s.status === "available").length,
        booked: filtered.filter((s) => s.status === "booked").length,
        inProgress: filtered.filter((s) => s.status === "in-progress").length,
        completed: filtered.filter((s) => s.status === "completed").length,
        blocked: filtered.filter((s) => s.status === "blocked").length,
      };
    },

    // Today's stats
    todayStats: (state) => {
      const today = formatDateString(new Date());
      const todaySlots = state.slots.filter((slot) => slot.date === today);
      return {
        available: todaySlots.filter((s) => s.status === "available").length,
        booked: todaySlots.filter((s) => s.status === "booked").length,
        inProgress: todaySlots.filter((s) => s.status === "in-progress").length,
        completed: todaySlots.filter((s) => s.status === "completed").length,
        blocked: todaySlots.filter((s) => s.status === "blocked").length,
        total: todaySlots.length,
      };
    },

    // Get slot by ID
    getSlotById:
      (state) =>
      (id: string): ScheduleSlot | undefined => {
        return state.slots.find((s) => s.id === id);
      },
  },

  actions: {
    // ==================== DATA FETCHING ====================

    async fetchSchedules(params?: ScheduleFilterParams) {
      this.isLoading = true;
      this.error = null;

      try {
        const queryParams: ScheduleFilterParams = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          ...params,
        };

        const result = await scheduleService.fetchAll(queryParams);

        this.slots = result.schedules.map(mapScheduleToSlot);
        this.pagination = {
          page: result.page,
          limit: result.limit,
          total: result.total,
          totalPages: result.totalPages,
        };
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch schedules";
        console.error("Error fetching schedules:", err);
        // Fallback to sample data
        this.slots = [...getInitialSlots()];
      } finally {
        this.isLoading = false;
      }
    },

    async fetchSchedulesByDate(date: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await scheduleService.fetchFiltered({
          date,
          limit: 100,
        });

        this.slots = result.schedules.map(mapScheduleToSlot);
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch schedules";
        console.error("Error fetching schedules:", err);
        // Fallback to sample data filtered by date
        this.slots = getInitialSlots().filter((s) => s.date === date);
      } finally {
        this.isLoading = false;
      }
    },

    // Fetch schedules by date using the dedicated endpoint (defaults to today)
    async fetchByDate(date?: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const schedules = await scheduleService.fetchByDate(date);
        this.slots = schedules.map(mapScheduleToSlot);
      } catch (err) {
        this.error =
          err instanceof Error
            ? err.message
            : "Failed to fetch schedules by date";
        console.error("Error fetching schedules by date:", err);
        // Fallback to sample data
        this.slots = [...getInitialSlots()];
      } finally {
        this.isLoading = false;
      }
    },

    // Fetch today's sessions from API
    async fetchTodaySessions() {
      return this.fetchByDate();
    },

    // Initialize store and auto-fetch schedules for today
    async initialize() {
      if (this.isInitialized) return;

      this.isLoading = true;
      this.error = null;

      try {
        const schedules = await scheduleService.fetchByDate();
        this.slots = schedules.map(mapScheduleToSlot);
        this.isInitialized = true;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to initialize schedules";
        console.error("Error initializing schedules:", err);
        // Fallback to sample data
        this.slots = [...getInitialSlots()];
        this.isInitialized = true;
      } finally {
        this.isLoading = false;
      }
    },

    async fetchAvailableSchedules(params?: AvailableScheduleParams) {
      this.isLoading = true;
      this.error = null;

      try {
        const available = await scheduleService.fetchAvailable(params);
        return available.map(mapScheduleToSlot);
      } catch (err) {
        this.error =
          err instanceof Error
            ? err.message
            : "Failed to fetch available schedules";
        console.error("Error fetching available schedules:", err);
        return getInitialSlots().filter((s) => s.status === "available");
      } finally {
        this.isLoading = false;
      }
    },

    // ==================== CRUD OPERATIONS ====================

    async createSlot(data: ScheduleFormFullData): Promise<ScheduleSlot | null> {
      if (
        !data.date ||
        !data.startTime ||
        !data.duration ||
        !data.vehicleId ||
        !data.instructorId
      ) {
        return null;
      }
      const createData: CreateScheduleData = {
        date: data.date,
        startTime: data.startTime,
        duration: data.duration,
        vehicleId: data.vehicleId,
        instructorId: data.instructorId,
        notes: data.notes,
      };
      try {
        const schedule = await scheduleService.create(createData);
        if (schedule) {
          const slot = mapScheduleToSlot(schedule);
          this.slots.push(slot);
          return slot;
        }
        return null;
      } catch {
        // Fallback to local creation
        return this.createSlotLocal(createData);
      }
    },

    // Local creation fallback
    createSlotLocal(data: CreateScheduleData): ScheduleSlot {
      const slot: ScheduleSlot = {
        id: crypto.randomUUID(),
        date: data.date,
        time: data.startTime,
        duration: `${data.duration} min`,
        car: "Vehicle", // Will be resolved by UI
        instructor: "Instructor", // Will be resolved by UI
        student: null,
        status: "available",
      };
      this.slots.push(slot);
      return slot;
    },

    async updateSlot(
      id: string,
      data: ScheduleFormFullData,
    ): Promise<ScheduleSlot | null> {
      try {
        const updateData: UpdateScheduleData = {
          date: data.date,
          startTime: data.startTime,
          duration: data.duration,
          vehicleId: data.vehicleId,
          instructorId: data.instructorId,
          status: data.status,
          notes: data.notes,
        };
        const schedule = await scheduleService.update(id, updateData);
        if (schedule) {
          const slot = mapScheduleToSlot(schedule);
          const index = this.slots.findIndex((s) => s.id === id);
          if (index !== -1) {
            this.slots[index] = slot;
          }
          return slot;
        }
        return null;
      } catch {
        return null;
      }
    },

    // Local update fallback
    updateSlotLocal(
      id: string,
      data: Partial<ScheduleSlot>,
    ): ScheduleSlot | null {
      const index = this.slots.findIndex((s) => s.id === id);
      if (index !== -1) {
        const existing = this.slots[index];
        if (!existing) return null;
        const updated: ScheduleSlot = {
          ...existing,
          id: data.id ?? existing.id,
          date: data.date ?? existing.date,
          time: data.time ?? existing.time,
          duration: data.duration ?? existing.duration,
          car: data.car ?? existing.car,
          instructor: data.instructor ?? existing.instructor,
          student: data.student ?? existing.student,
          status: data.status ?? existing.status,
        };
        this.slots[index] = updated;
        return updated;
      }
      return null;
    },

    async deleteSlot(id: string): Promise<boolean> {
      try {
        const success = await scheduleService.delete(id);
        if (success) {
          this.slots = this.slots.filter((s) => s.id !== id);
        }
        return success;
      } catch {
        // Fallback to local delete
        this.slots = this.slots.filter((s) => s.id !== id);
        return true;
      }
    },

    // ==================== SLOT OPERATIONS ====================

    addSlot(slot: ScheduleSlot) {
      // Check for duplicate
      const exists = this.slots.some(
        (s) =>
          s.date === slot.date && s.time === slot.time && s.car === slot.car,
      );
      if (exists) {
        console.warn(
          "Slot already exists for this date/time/vehicle combination",
        );
        return null;
      }
      this.slots.push(slot);
      return slot;
    },

    editSlot(
      id: string,
      data: Partial<Omit<ScheduleSlot, "id" | "status" | "student">>,
    ): ScheduleSlot | null {
      const index = this.slots.findIndex((s) => s.id === id);
      if (index !== -1) {
        const slot = this.slots[index];
        if (!slot) return null;
        const updated = { ...slot, ...data };
        this.slots[index] = updated;
        return updated;
      }
      return null;
    },

    updateSlotStatus(id: string, status: ScheduleStatus): ScheduleSlot | null {
      const index = this.slots.findIndex((s) => s.id === id);
      if (index !== -1) {
        const slot = this.slots[index];
        if (!slot) return null;
        slot.status = status;
        // Clear student if status becomes available
        if (status === "available") {
          slot.student = null;
        }
        return slot;
      }
      return null;
    },

    toggleSlotStatus(id: string) {
      const slot = this.slots.find((s) => s.id === id);
      if (!slot) return null;

      const newStatus: ScheduleStatus =
        slot.status === "available" ? "blocked" : "available";
      return this.updateSlotStatus(id, newStatus);
    },

    // ==================== BOOKING OPERATIONS ====================

    async bookSlot(
      id: string,
      studentName: string,
    ): Promise<ScheduleSlot | null> {
      try {
        const schedule = await scheduleService.bookSlot(id, {
          studentId: "", // Will be resolved by backend
        });
        if (schedule) {
          const slot = mapScheduleToSlot(schedule);
          const index = this.slots.findIndex((s) => s.id === id);
          if (index !== -1) {
            this.slots[index] = slot;
          }
          return slot;
        }
        return null;
      } catch {
        // Fallback to local booking
        return this.bookSlotLocal(id, studentName);
      }
    },

    // Local booking fallback
    bookSlotLocal(id: string, studentName: string): ScheduleSlot | null {
      const index = this.slots.findIndex((s) => s.id === id);
      if (index !== -1) {
        const slot = this.slots[index];
        if (!slot) return null;
        if (slot.status !== "available") {
          console.warn("Cannot book a slot that is not available");
          return null;
        }
        slot.student = studentName;
        slot.status = "booked";
        return slot;
      }
      return null;
    },

    async cancelBooking(id: string): Promise<ScheduleSlot | null> {
      try {
        const schedule = await scheduleService.cancelBooking(id);
        if (schedule) {
          const slot = mapScheduleToSlot(schedule);
          const index = this.slots.findIndex((s) => s.id === id);
          if (index !== -1) {
            this.slots[index] = slot;
          }
          return slot;
        }
        return null;
      } catch {
        // Fallback to local cancel
        return this.cancelBookingLocal(id);
      }
    },

    // Local cancel booking fallback
    cancelBookingLocal(id: string): ScheduleSlot | null {
      const index = this.slots.findIndex((s) => s.id === id);
      if (index !== -1) {
        const slot = this.slots[index];
        if (!slot) return null;
        slot.student = null;
        slot.status = "available";
        return slot;
      }
      return null;
    },

    // ==================== SESSION OPERATIONS ====================

    startSession(id: string): ScheduleSlot | null {
      const index = this.slots.findIndex((s) => s.id === id);
      if (index !== -1) {
        const slot = this.slots[index];
        if (!slot) return null;
        if (slot.status !== "booked") {
          console.warn("Can only start a session for a booked slot");
          return null;
        }
        slot.status = "in-progress";
        return slot;
      }
      return null;
    },

    completeSession(id: string): ScheduleSlot | null {
      const index = this.slots.findIndex((s) => s.id === id);
      if (index !== -1) {
        const slot = this.slots[index];
        if (!slot) return null;
        if (slot.status !== "in-progress") {
          console.warn("Can only complete a session that is in progress");
          return null;
        }
        slot.status = "completed";
        return slot;
      }
      return null;
    },

    // ==================== FILTER OPERATIONS ====================

    setSelectedDate(date: string) {
      this.selectedDate = date;
    },

    setFilterInstructor(instructor: string) {
      this.filterInstructor = instructor;
    },

    setFilterVehicle(vehicle: string) {
      this.filterVehicle = vehicle;
    },

    clearFilters() {
      this.filterInstructor = "All Instructors";
      this.filterVehicle = "All Vehicles";
    },

    // ==================== UTILITY OPERATIONS ====================

    formatDateString,

    getSlotStats(date?: string): {
      available: number;
      booked: number;
      inProgress: number;
      completed: number;
      blocked: number;
    } {
      const targetDate = date || this.selectedDate;
      const filtered = this.slots.filter((s) => s.date === targetDate);
      return {
        available: filtered.filter((s) => s.status === "available").length,
        booked: filtered.filter((s) => s.status === "booked").length,
        inProgress: filtered.filter((s) => s.status === "in-progress").length,
        completed: filtered.filter((s) => s.status === "completed").length,
        blocked: filtered.filter((s) => s.status === "blocked").length,
      };
    },

    reset() {
      this.slots = [];
      this.isLoading = false;
      this.error = null;
      this.isInitialized = false;
      this.selectedDate = formatDateString(new Date());
      this.filterInstructor = "All Instructors";
      this.filterVehicle = "All Vehicles";
    },
  },
});
