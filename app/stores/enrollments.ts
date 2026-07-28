import { defineStore } from "pinia";
import { enrollmentService } from "~/services/enrollmentService";
import type {
  Enrollment,
  EnrollmentBrief,
  EnrollmentStatus,
  CreateEnrollmentData,
  UpdateEnrollmentData,
} from "~/services/enrollmentService";

interface EnrollmentsState {
  enrollments: EnrollmentBrief[];
  currentEnrollment: Enrollment | null;
  isLoading: boolean;
  isCreateLoading: boolean;
  isUpdateLoading: boolean;
  error: string | null;
}

export const useEnrollmentsStore = defineStore("enrollments", {
  state: (): EnrollmentsState => ({
    enrollments: [],
    currentEnrollment: null,
    isLoading: false,
    isCreateLoading: false,
    isUpdateLoading: false,
    error: null,
  }),

  getters: {
    activeEnrollments: (state) =>
      state.enrollments.filter((e) => e.status === "active"),
    pendingEnrollments: (state) =>
      state.enrollments.filter((e) => e.status === "pending"),
    completedEnrollments: (state) =>
      state.enrollments.filter((e) => e.status === "completed"),
    hasActiveEnrollment: (state) =>
      state.enrollments.some((e) => e.status === "active"),
    getEnrollmentById: (state) => (id: string) =>
      state.enrollments.find((e) => e.id === id),
    totalSessions: (state) =>
      state.enrollments.reduce((sum, e) => sum + e.totalSessions, 0),
    completedSessions: (state) =>
      state.enrollments.reduce((sum, e) => sum + e.completedSessions, 0),
    remainingSessions: (state) =>
      state.enrollments.reduce((sum, e) => sum + e.remainingSessions, 0),
  },

  actions: {
    async fetchEnrollments() {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await enrollmentService.fetchMyEnrollments();
        this.enrollments = result;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch enrollments";
        console.error("Error fetching enrollments:", err);
      } finally {
        this.isLoading = false;
      }
    },

    async fetchActiveEnrollment() {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await enrollmentService.fetchActiveEnrollment();
        this.currentEnrollment = result;
        return result;
      } catch (err) {
        this.error =
          err instanceof Error
            ? err.message
            : "Failed to fetch active enrollment";
        console.error("Error fetching active enrollment:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async createEnrollment(
      data: CreateEnrollmentData,
    ): Promise<Enrollment | null> {
      this.isCreateLoading = true;
      this.error = null;

      try {
        const response = await enrollmentService.create(data);
        console.log("[ENROLLMENT STORE] createEnrollment response:", response);

        const newEnrollment = response && typeof response === "object" && "enrollment" in response
          ? (response as any).enrollment
          : response;

        if (newEnrollment && newEnrollment.id) {
          // Normalize the enrollment data to match our interface
          const normalizedEnrollment: Enrollment = {
            id: newEnrollment.id,
            userId: newEnrollment.userId || data.userId,
            userName: newEnrollment.userName || "",
            userEmail: newEnrollment.userEmail || "",
            packageId: newEnrollment.packageId || data.packageId,
            packageName: newEnrollment.packageName || "",
            totalSessions: newEnrollment.totalSessions || 0,
            completedSessions: newEnrollment.completedSessions || 0,
            remainingSessions: newEnrollment.remainingSessions || 0,
            status: newEnrollment.status || "pending",
            enrollmentDate: newEnrollment.enrollmentDate || new Date().toISOString(),
            expiryDate: newEnrollment.expiryDate || newEnrollment.expiresAt || "",
            price: newEnrollment.price || newEnrollment.totalPrice || data.price || 0,
            discountPrice: newEnrollment.discountPrice || data.discountPrice || 0,
            paymentStatus: newEnrollment.paymentStatus || "pending",
            createdAt: newEnrollment.createdAt,
            updatedAt: newEnrollment.updatedAt,
          };

          // Add to enrollments list
          this.enrollments.unshift({
            id: normalizedEnrollment.id,
            packageName: normalizedEnrollment.packageName,
            totalSessions: normalizedEnrollment.totalSessions,
            completedSessions: normalizedEnrollment.completedSessions,
            remainingSessions: normalizedEnrollment.remainingSessions,
            status: normalizedEnrollment.status,
            enrollmentDate: normalizedEnrollment.enrollmentDate,
            expiryDate: normalizedEnrollment.expiryDate,
            paymentStatus: normalizedEnrollment.paymentStatus,
          });
          this.currentEnrollment = normalizedEnrollment;
          return normalizedEnrollment;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to create enrollment";
        console.error("Error creating enrollment:", err);
        return null;
      } finally {
        this.isCreateLoading = false;
      }
    },

    async updateEnrollment(id: string, data: UpdateEnrollmentData) {
      this.isUpdateLoading = true;
      this.error = null;

      try {
        const updated = await enrollmentService.update(id, data);
        if (updated) {
          // Update in list
          const index = this.enrollments.findIndex((e) => e.id === id);
          if (index !== -1) {
            this.enrollments[index] = {
              id: updated.id,
              packageName: updated.packageName,
              totalSessions: updated.totalSessions,
              completedSessions: updated.completedSessions,
              remainingSessions: updated.remainingSessions,
              status: updated.status,
              enrollmentDate: updated.enrollmentDate,
              expiryDate: updated.expiryDate,
              paymentStatus: updated.paymentStatus,
            };
          }
          // Update current if same
          if (this.currentEnrollment?.id === id) {
            this.currentEnrollment = updated;
          }
        }
        return updated;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to update enrollment";
        console.error("Error updating enrollment:", err);
        return null;
      } finally {
        this.isUpdateLoading = false;
      }
    },

    async completeSession(id: string) {
      try {
        const updated = await enrollmentService.completeSession(id);
        if (updated) {
          // Update in list
          const index = this.enrollments.findIndex((e) => e.id === id);
          if (index !== -1) {
            this.enrollments[index]!.completedSessions =
              updated.completedSessions;
            this.enrollments[index]!.remainingSessions =
              updated.remainingSessions;
          }
          // Update current if same
          if (this.currentEnrollment?.id === id) {
            this.currentEnrollment.completedSessions =
              updated.completedSessions;
            this.currentEnrollment.remainingSessions =
              updated.remainingSessions;
          }
        }
        return updated;
      } catch (err) {
        console.error("Error completing session:", err);
        return null;
      }
    },

    async cancelEnrollment(id: string) {
      try {
        const updated = await enrollmentService.cancel(id);
        if (updated) {
          // Update status in list
          const index = this.enrollments.findIndex((e) => e.id === id);
          if (index !== -1) {
            this.enrollments[index]!.status = updated.status;
          }
          // Update current if same
          if (this.currentEnrollment?.id === id) {
            this.currentEnrollment.status = updated.status;
          }
        }
        return updated;
      } catch (err) {
        console.error("Error cancelling enrollment:", err);
        return null;
      }
    },

    // Reset enrollment state after logout
    reset() {
      this.enrollments = [];
      this.currentEnrollment = null;
      this.isLoading = false;
      this.isCreateLoading = false;
      this.isUpdateLoading = false;
      this.error = null;
    },
  },
});
