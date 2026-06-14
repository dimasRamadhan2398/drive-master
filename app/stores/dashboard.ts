import { dashboardService, type RecentRegistration, type DashboardStats as DashboardStatsType } from "~/services/dashboardService";
import { defineStore } from "pinia";

export interface DashboardState {
  recentRegistrations: RecentRegistration[];
  stats: DashboardStatsType;
  isLoading: boolean;
  error?: string | null;
  lastFetched?: Date | null;
}

export const useDashboardStore = defineStore("dashboard", {
  state: (): DashboardState => ({
    recentRegistrations: [],
    stats: {
      totalUsers: 0,
      totalMembers: 0,
      totalInstructors: 0,
      recentRegistrations: 0,
      activeSessions: 0,
      totalSessions: 0,
      revenueMTD: 0,
      revenueCurrency: "IDR",
      certificatesIssued: 0,
      totalCertifications: 0,
    },
    isLoading: false,
    error: null,
    lastFetched: null,
  }),
  getters: {
    totalRecentRegistrations: (state) => state.recentRegistrations.length,
  },
  actions: {
    async fetchDashboardStats() {
      this.isLoading = true;
      this.error = null;
      try {
        const stats = await dashboardService.fetchDashboardStats();
        this.stats = stats;
        this.lastFetched = new Date();
      } catch (err) {
        this.error = "Failed to load dashboard stats.";
      } finally {
        this.isLoading = false;
      }
    },
    async fetchRecentRegistrations(params: {
      limit?: number;
      fromDate?: string;
      toDate?: string;
    } = {}) {
      this.isLoading = true;
      this.error = null;
      try {
        const response = await dashboardService.fetchRecentRegistrations(params);
        this.recentRegistrations = response;
        this.lastFetched = new Date();
      } catch (err) {
        this.error = "Failed to load recent registrations.";
      } finally {
        this.isLoading = false;
      }
    },
    async fetchDashboardData() {
      this.isLoading = true;
      this.error = null;
      try {
        const [stats, registrations] = await Promise.all([
          dashboardService.fetchDashboardStats(),
          dashboardService.fetchRecentRegistrations({ limit: 10 }),
        ]);
        this.stats = stats;
        this.recentRegistrations = registrations;
        this.lastFetched = new Date();
        return { stats, registrations };
      } catch (err) {
        this.error = "Failed to load dashboard data.";
        return { stats: this.stats, registrations: [] };
      } finally {
        this.isLoading = false;
      }
    },
    clearError() {
      this.error = null;
    },
  },
});