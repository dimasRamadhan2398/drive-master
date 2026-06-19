import { defineStore } from "pinia";
import { studentService } from "~/services/studentService";
import type {
  CreateStudentData,
  UpdateStudentData,
  Entitlement,
} from "~/services/studentService";

export interface Student {
  id: string;
  name: string;
  email: string;
  phone: string;
  package: string;
  progress: number;
  completedSessions: number;
  totalSessions: number;
  joinDate: string;
  status: "active" | "pending" | "completed";
  entitlements: Entitlement[];
}

interface StudentsState {
  students: Student[];
  searchResults: Student[];
  isLoading: boolean;
  error: string | null;
  // Pagination state
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
  // Client-side search/filter state
  searchQuery: string;
  statusFilter: string;
  // Server-side pagination mode
  useServerPagination: boolean;
}

const packageSessionMap: { [key: string]: number } = {
  "6x": 6,
  "8x": 8,
  "10x": 10,
  "12x": 12,
};

const initialStudents: Student[] = [
  {
    id: "ee61c870-98ab-4c28-a03c-402103f94e56",
    name: "John Doe",
    email: "john@example.com",
    phone: "081234567890",
    package: "Gold Package",
    progress: 40,
    completedSessions: 4,
    totalSessions: 10,
    joinDate: "Mar 10, 2026",
    status: "active",
    entitlements: [
      {
        id: "ent-1",
        memberId: "ee61c870-98ab-4c28-a03c-402103f94e56",
        bookingId: "book-1",
        packageId: "pkg-1",
        packageName: "Gold Package",
        isNightSession: false,
        isWeekendSession: false,
        totalSessions: 10,
        remaining: 6,
        usedSessions: 4,
        startDate: "2026-03-10T10:00:00Z",
        endDate: null,
        status: "active",
        createdAt: "2026-03-10T10:00:00Z",
        updatedAt: "2026-03-10T10:00:00Z",
      },
    ],
  },
  {
    id: "8cb66c0e-a7fd-4626-8cd5-b210198bf74d",
    name: "Sarah Putri",
    email: "sarah@example.com",
    phone: "081234567891",
    package: "Platinum Package",
    progress: 75,
    completedSessions: 11,
    totalSessions: 15,
    joinDate: "Feb 20, 2026",
    status: "active",
    entitlements: [
      {
        id: "ent-2",
        memberId: "8cb66c0e-a7fd-4626-8cd5-b210198bf74d",
        bookingId: "book-2",
        packageId: "pkg-2",
        packageName: "Platinum Package",
        isNightSession: false,
        isWeekendSession: false,
        totalSessions: 15,
        remaining: 4,
        usedSessions: 11,
        startDate: "2026-02-20T10:00:00Z",
        endDate: null,
        status: "active",
        createdAt: "2026-02-20T10:00:00Z",
        updatedAt: "2026-02-20T10:00:00Z",
      },
    ],
  },
  {
    id: "7523a0e1-5be7-4cb5-83af-c8f8c355a53c",
    name: "Budi Santoso",
    email: "budi@example.com",
    phone: "081234567892",
    package: "Silver Package",
    progress: 100,
    completedSessions: 5,
    totalSessions: 5,
    joinDate: "Jan 15, 2026",
    status: "completed",
    entitlements: [
      {
        id: "ent-3",
        memberId: "7523a0e1-5be7-4cb5-83af-c8f8c355a53c",
        bookingId: "book-3",
        packageId: "pkg-3",
        packageName: "Silver Package",
        isNightSession: false,
        isWeekendSession: false,
        totalSessions: 5,
        remaining: 0,
        usedSessions: 5,
        startDate: "2026-01-15T10:00:00Z",
        endDate: "2026-02-15T10:00:00Z",
        status: "completed",
        createdAt: "2026-01-15T10:00:00Z",
        updatedAt: "2026-02-15T10:00:00Z",
      },
    ],
  },
  {
    id: "a1b2c3d4-e5f6-7890-a1b2-c3d4e5f67890",
    name: "Amanda Chen",
    email: "amanda@example.com",
    phone: "081234567893",
    package: "Gold Package",
    progress: 20,
    completedSessions: 2,
    totalSessions: 10,
    joinDate: "Mar 25, 2026",
    status: "active",
    entitlements: [
      {
        id: "ent-4",
        memberId: "a1b2c3d4-e5f6-7890-a1b2-c3d4e5f67890",
        bookingId: "book-4",
        packageId: "pkg-1",
        packageName: "Gold Package",
        isNightSession: false,
        isWeekendSession: false,
        totalSessions: 10,
        remaining: 8,
        usedSessions: 2,
        startDate: "2026-03-25T10:00:00Z",
        endDate: null,
        status: "active",
        createdAt: "2026-03-25T10:00:00Z",
        updatedAt: "2026-03-25T10:00:00Z",
      },
    ],
  },
  {
    id: "f2g3h4i5-j6k7-8901-f2g3-h4i5j6k78901",
    name: "Ricky Wijaya",
    email: "ricky@example.com",
    phone: "081234567894",
    package: "No Package",
    progress: 0,
    completedSessions: 0,
    totalSessions: 0,
    joinDate: "Apr 1, 2026",
    status: "pending",
    entitlements: [],
  },
];

export const useStudentsStore = defineStore("students", {
  state: (): StudentsState => ({
    students: [],
    searchResults: [],
    isLoading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0,
    },
    searchQuery: "",
    statusFilter: "all",
    useServerPagination: true, // Default to server-side pagination
  }),

  getters: {
    activeStudents: (state) =>
      state.students.filter((s) => s.status === "active"),
    pendingStudents: (state) =>
      state.students.filter((s) => s.status === "pending"),
    completedStudents: (state) =>
      state.students.filter((s) => s.status === "completed"),
    totalStudents: (state) => state.students.length,
  },

  actions: {
    // Client-side filtering helper
    applyClientSideFiltering() {
      let filtered = [...initialStudents];

      // Apply search filter
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase();
        filtered = filtered.filter(
          (s) =>
            s.name.toLowerCase().includes(query) ||
            s.email.toLowerCase().includes(query),
        );
      }

      // Apply status filter
      if (this.statusFilter !== "all") {
        filtered = filtered.filter((s) => s.status === this.statusFilter);
      }

      this.students = filtered;
      this.pagination = {
        page: 1,
        limit: this.pagination.limit,
        total: filtered.length,
        totalPages: 1,
      };
    },

    async fetchStudents(page = 1, resetPage = true) {
      this.isLoading = true;
      this.error = null;

      try {
        // If resetting page, set to 1
        if (resetPage) {
          this.pagination.page = page;
        }

        const params: {
          page: number;
          limit: number;
          search?: string;
          status?: string;
        } = {
          page: this.pagination.page,
          limit: this.pagination.limit,
        };

        // Add search/filter params
        if (this.searchQuery) {
          params.search = this.searchQuery;
        }
        if (this.statusFilter !== "all") {
          params.status = this.statusFilter;
        }

        const result = await studentService.fetchAll(params);

        this.students = result.students;
        this.pagination = result.pagination;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch students";
        console.error("Error fetching students:", err);
        // API failed, use dummy data with client-side filtering
        this.applyClientSideFiltering();
      } finally {
        this.isLoading = false;
      }
    },

    async fetchStudentsNoPagination() {
      this.isLoading = true;
      this.error = null;
      try {
        const students = await studentService.fetchAllFlat();
        if (students.length > 0) {
          this.students = students;
        }
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch students";
        console.error("Error fetching students:", err);
        // API failed, use dummy data with client-side filtering
        this.applyClientSideFiltering();
      } finally {
        this.isLoading = false;
      }
    },

    // Set search query and optionally trigger fetch
    setSearchQuery(query: string) {
      this.searchQuery = query;
      if (this.useServerPagination) {
        this.fetchStudents(1, true); // Reset to page 1 on new search
      }
    },

    // Set status filter and optionally trigger fetch
    setStatusFilter(status: string) {
      this.statusFilter = status;
      if (this.useServerPagination) {
        this.fetchStudents(1, true); // Reset to page 1 on new filter
      }
    },

    // Change page for server-side pagination
    setPage(page: number) {
      this.pagination.page = page;
      this.fetchStudents(page, false); // Don't reset page
    },

    async addStudent(data: CreateStudentData) {
      try {
        const newStudent = await studentService.create(data);
        this.students.unshift(newStudent);
        return newStudent;
      } catch {
        // Fallback to local creation if API fails
        return this.addStudentLocal({
          name: `${data.firstName} ${data.lastName}`.trim(),
          email: data.email,
          phone: data.phoneNumber,
          package: "-",
          status: "pending",
          entitlements: [],
        });
      }
    },

    // Local creation fallback (when API is not available)
    addStudentLocal(
      data: Omit<
        Student,
        "id" | "progress" | "completedSessions" | "joinDate" | "totalSessions"
      >,
    ): Student {
      const newId = crypto.randomUUID();
      const totalSessions = studentService.getTotalSessions(data.package);

      const newStudent: Student = {
        id: newId,
        name: data.name,
        email: data.email,
        phone: data.phone,
        package: data.package,
        totalSessions,
        progress: 0,
        completedSessions: 0,
        joinDate: studentService.formatJoinDate(),
        status: data.status || "pending",
        entitlements: [],
      };

      this.students.unshift(newStudent);
      return newStudent;
    },

    async updateStudent(id: string, data: UpdateStudentData) {
      try {
        const updatedStudent = await studentService.update(id, data);
        if (updatedStudent) {
          const index = this.students.findIndex((s) => s.id === id);
          if (index !== -1) {
            this.students[index] = updatedStudent;
          }
          return updatedStudent;
        }
        return null;
      } catch {
        // Fallback to local update if API fails
        // Convert UpdateStudentData to Partial<Student>
        const localData: Partial<Student> = {
          phone: data.phoneNumber,
        };
        return this.updateStudentLocal(id, localData);
      }
    },

    // Local update fallback (when API is not available)
    updateStudentLocal(id: string, data: Partial<Student>): Student | null {
      const index = this.students.findIndex((s) => s.id === id);
      if (index !== -1) {
        const existing = this.students[index];
        if (!existing) return null;

        const updatedData: Student = {
          id: existing.id,
          name: data.name ?? existing.name,
          email: data.email ?? existing.email,
          phone: data.phone ?? existing.phone,
          package: data.package ?? existing.package,
          progress: data.progress ?? existing.progress,
          completedSessions:
            data.completedSessions ?? existing.completedSessions,
          totalSessions: data.totalSessions ?? existing.totalSessions,
          joinDate: data.joinDate ?? existing.joinDate,
          status: data.status ?? existing.status,
          entitlements: data.entitlements ?? existing.entitlements,
        };

        // Update totalSessions if package changed
        if (data.package) {
          updatedData.totalSessions = studentService.getTotalSessions(
            data.package,
          );
          // Recalculate progress
          updatedData.progress = studentService.calculateProgress(
            updatedData.completedSessions,
            updatedData.totalSessions,
          );
        }

        this.students[index] = updatedData;
        return updatedData;
      }
      return null;
    },

    async deleteStudent(id: string) {
      try {
        const success = await studentService.delete(id);
        if (success) {
          this.students = this.students.filter((s) => s.id !== id);
        }
        return success;
      } catch {
        // Fallback to local delete if API fails
        this.students = this.students.filter((s) => s.id !== id);
        return true;
      }
    },

    getStudentById(id: string) {
      return this.students.find((s) => s.id === id);
    },

    filterStudents(searchQuery: string, status: string) {
      return this.students.filter((student) => {
        const matchesSearch =
          student.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          student.email.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus = status === "all" || student.status === status;
        return matchesSearch && matchesStatus;
      });
    },

    async searchStudents(query: string): Promise<Student[]> {
      try {
        const result = await studentService.searchStudents(query);
        // Store results separately from main students list
        this.searchResults = result.students || [];
        return this.searchResults;
      } catch (e) {
        console.error("Error searching students:", e);
        this.searchResults = [];
        return [];
      }
    },

    async updateSessionProgress(id: string, completedSessions: number) {
      try {
        const updatedStudent = await studentService.updateProgress(
          id,
          completedSessions,
        );
        if (updatedStudent) {
          const index = this.students.findIndex((s) => s.id === id);
          if (index !== -1) {
            this.students[index] = updatedStudent;
          }
          return;
        }
      } catch {
        // Fallback to local update if API fails
      }

      // Local progress update fallback
      const student = this.students.find((s) => s.id === id);
      if (student) {
        student.completedSessions = completedSessions;
        student.totalSessions =
          packageSessionMap[student.package] || student.totalSessions;
        student.progress = studentService.calculateProgress(
          completedSessions,
          student.totalSessions,
        );

        // Auto-complete if all sessions done
        if (student.progress >= 100 && student.status !== "completed") {
          student.status = "completed";
        }
      }
    },

    clearSearchResults() {
      this.searchResults = [];
    },

    reset() {
      this.students = [...initialStudents];
      this.isLoading = false;
      this.error = null;
    },
  },
});
