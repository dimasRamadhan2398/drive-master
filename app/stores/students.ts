import { defineStore } from "pinia";
import { studentService } from "~/services/studentService";
import type {
  CreateStudentData,
  UpdateStudentData,
} from "~/services/studentService";

export interface Student {
  id: number;
  name: string;
  email: string;
  phone: string;
  package: string;
  progress: number;
  completedSessions: number;
  totalSessions: number;
  joinDate: string;
  status: "active" | "pending" | "completed";
}

interface StudentsState {
  students: Student[];
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
    id: 1,
    name: "John Doe",
    email: "john@example.com",
    phone: "081234567890",
    package: "8x",
    progress: 40,
    completedSessions: 4,
    totalSessions: 10,
    joinDate: "Mar 10, 2026",
    status: "active",
  },
  {
    id: 2,
    name: "Sarah Putri",
    email: "sarah@example.com",
    phone: "081234567891",
    package: "12x",
    progress: 75,
    completedSessions: 11,
    totalSessions: 15,
    joinDate: "Feb 20, 2026",
    status: "active",
  },
  {
    id: 3,
    name: "Budi Santoso",
    email: "budi@example.com",
    phone: "081234567892",
    package: "6x",
    progress: 100,
    completedSessions: 5,
    totalSessions: 5,
    joinDate: "Jan 15, 2026",
    status: "completed",
  },
  {
    id: 4,
    name: "Amanda Chen",
    email: "amanda@example.com",
    phone: "081234567893",
    package: "8x",
    progress: 20,
    completedSessions: 2,
    totalSessions: 10,
    joinDate: "Mar 25, 2026",
    status: "active",
  },
  {
    id: 5,
    name: "Ricky Wijaya",
    email: "ricky@example.com",
    phone: "081234567894",
    package: "8x",
    progress: 0,
    completedSessions: 0,
    totalSessions: 10,
    joinDate: "Apr 1, 2026",
    status: "pending",
  },
];

export const useStudentsStore = defineStore("students", {
  state: (): StudentsState => ({
    students: [],
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
    async fetchStudents(page = 1, resetPage = true) {
      this.isLoading = true;
      this.error = null;

      try {
        // If resetting page, set to 1
        if (resetPage) {
          this.pagination.page = page;
        }

        const params: { page: number; limit: number; search?: string; status?: string } = {
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
        // API failed, use dummy data
        this.students = [...initialStudents];
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
        // API failed, use dummy data
        this.students = [...initialStudents];
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
          package: "8x",
          status: "pending",
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
      const newId = Math.max(...this.students.map((s) => s.id), 0) + 1;
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
      };

      this.students.unshift(newStudent);
      return newStudent;
    },

    async updateStudent(id: number, data: UpdateStudentData) {
      // Convert numeric ID to userId string (for API)
      const userId = String(id);

      try {
        const updatedStudent = await studentService.update(userId, data);
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
    updateStudentLocal(id: number, data: Partial<Student>): Student | null {
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

    async deleteStudent(id: number) {
      const userId = String(id);
      try {
        const success = await studentService.delete(userId);
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

    getStudentById(id: number) {
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

    async updateSessionProgress(id: number, completedSessions: number) {
      const userId = String(id);
      try {
        const updatedStudent = await studentService.updateProgress(
          userId,
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

    reset() {
      this.students = [...initialStudents];
      this.isLoading = false;
      this.error = null;
    },
  },
});
