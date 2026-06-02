import { defineStore } from "pinia";

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
    students: initialStudents,
    isLoading: false,
    error: null,
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
    addStudent(
      data: Omit<Student, "id" | "progress" | "completedSessions" | "joinDate" | "totalSessions">
    ) {
      const newId = Math.max(...this.students.map((s) => s.id), 0) + 1;
      const totalSessions = packageSessionMap[data.package] || 0;

      const newStudent: Student = {
        ...data,
        id: newId,
        totalSessions,
        progress: 0,
        completedSessions: 0,
        joinDate: new Date().toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
      };

      this.students.unshift(newStudent);
      return newStudent;
    },

    updateStudent(id: number, data: Partial<Student>) {
      const index = this.students.findIndex((s) => s.id === id);
      if (index !== -1) {
        const updatedData = { ...this.students[index], ...data };

        // Update totalSessions if package changed
        if (data.package) {
          updatedData.totalSessions = packageSessionMap[data.package] || 0;
          // Recalculate progress
          updatedData.progress =
            updatedData.totalSessions > 0
              ? Math.round(
                  (updatedData.completedSessions / updatedData.totalSessions) *
                    100
                )
              : 0;
        }

        this.students[index] = updatedData;
        return updatedData;
      }
      return null;
    },

    deleteStudent(id: number) {
      this.students = this.students.filter((s) => s.id !== id);
    },

    getStudentById(id: number) {
      return this.students.find((s) => s.id === id);
    },

    filterStudents(searchQuery: string, status: string) {
      return this.students.filter((student) => {
        const matchesSearch =
          student.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          student.email.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus =
          status === "all" || student.status === status;
        return matchesSearch && matchesStatus;
      });
    },

    updateSessionProgress(id: number, completedSessions: number) {
      const student = this.students.find((s) => s.id === id);
      if (student) {
        student.completedSessions = completedSessions;
        student.totalSessions =
          packageSessionMap[student.package] || student.totalSessions;
        student.progress =
          student.totalSessions > 0
            ? Math.round((completedSessions / student.totalSessions) * 100)
            : 0;

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