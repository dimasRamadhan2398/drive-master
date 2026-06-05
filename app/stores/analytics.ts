import { defineStore } from "pinia";

export interface GAOverviewItem {
  date: string;
  users: number;
  pageviews: number;
}

export interface GAFunnelStep {
  event_name: string;
  count: number;
}

export interface FunnelLabels {
  [key: string]: string;
}

interface AnalyticsState {
  gaOverview: GAOverviewItem[];
  gaFunnel: GAFunnelStep[];
  isLoading: boolean;
  error: string | null;
  lastFetched: string | null;
}

const funnelLabels: FunnelLabels = {
  page_view: "1. Landing Page Visit",
  view_item: "2. View Packages",
  begin_checkout: "3. Start Booking Process",
  purchase: "4. Successful Payment",
};

const funnelColors = [
  "bg-primary-500",
  "bg-info-500",
  "bg-warning-500",
  "bg-success-500",
];

export const useAnalyticsStore = defineStore("analytics", {
  state: (): AnalyticsState => ({
    gaOverview: [],
    gaFunnel: [],
    isLoading: false,
    error: null,
    lastFetched: null,
  }),

  getters: {
    totalGaUsers: (state) =>
      state.gaOverview.reduce((sum, item) => sum + item.users, 0),

    totalGaPageViews: (state) =>
      state.gaOverview.reduce((sum, item) => sum + item.pageviews, 0),

    overallConversionRate(): string {
      const start = this.gaFunnelData[0]?.count;
      const end = this.gaFunnelData[3]?.count;

      if (!start || !end) return "0.0";
      return ((end / start) * 100).toFixed(1);
    },

    gaFunnelData: (state) => state.gaFunnel,

    chartPath(): string {
      if (this.gaOverviewData.length === 0) return "";
      const width = 500;
      const height = 150;
      const maxVal = Math.max(
        ...this.gaOverviewData.map((d) => d.pageviews),
        100
      );

      return this.gaOverviewData
        .map((d, idx) => {
          const x = (idx / (this.gaOverviewData.length - 1)) * width;
          const y = height - (d.pageviews / maxVal) * height;
          return `${idx === 0 ? "M" : "L"} ${x.toFixed(1)} ${y.toFixed(1)}`;
        })
        .join(" ");
    },

    chartAreaPath(): string {
      if (this.gaOverviewData.length === 0) return "";
      const width = 500;
      const height = 150;
      const base = this.chartPathData;
      return `${base} L ${width} ${height} L 0 ${height} Z`;
    },

    chartUsersPath(): string {
      if (this.gaOverviewData.length === 0) return "";
      const width = 500;
      const height = 150;
      const maxVal = Math.max(
        ...this.gaOverviewData.map((d) => d.pageviews),
        100
      );

      return this.gaOverviewData
        .map((d, idx) => {
          const x = (idx / (this.gaOverviewData.length - 1)) * width;
          const y = height - (d.users / maxVal) * height;
          return `${idx === 0 ? "M" : "L"} ${x.toFixed(1)} ${y.toFixed(1)}`;
        })
        .join(" ");
    },

    gaOverviewData: (state) => state.gaOverview,
    chartPathData(): string {
      return this.chartPath;
    },

    hasData: (state) =>
      state.gaOverview.length > 0 || state.gaFunnel.length > 0,

    funnelStepConversion(): { step: string; count: number; percentage: number }[] {
      if (this.gaFunnelData.length === 0) return [];
      const base = this.gaFunnelData[0]?.count || 1;
      return this.gaFunnelData.map((step, idx) => ({
        step: funnelLabels[step.event_name] || step.event_name,
        count: step.count,
        percentage:
          idx === 0 ? 100 : ((step.count / base) * 100).toFixed(1) + "%",
      }));
    },
  },

  actions: {
    async fetchAnalyticsData() {
      this.isLoading = true;
      this.error = null;

      try {
        const { core } = useApiClients();

        const [overviewRes, funnelRes] = await Promise.all([
          core<{ data: GAOverviewItem[] }>("/admin/analytics/overview"),
          core<{ data: GAFunnelStep[] }>("/admin/analytics/funnel"),
        ]);

        if (overviewRes && overviewRes.data) {
          this.gaOverview = overviewRes.data;
        }
        if (funnelRes && funnelRes.data) {
          this.gaFunnel = funnelRes.data;
        }

        this.lastFetched = new Date().toISOString();
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch analytics";
        console.error("Failed to fetch analytics data:", err);
      } finally {
        this.isLoading = false;
      }
    },

    setMockData(overview: GAOverviewItem[], funnel: GAFunnelStep[]) {
      this.gaOverview = overview;
      this.gaFunnel = funnel;
      this.lastFetched = new Date().toISOString();
    },

    getFunnelLabel(eventName: string): string {
      return funnelLabels[eventName] || eventName;
    },

    getFunnelColor(index: number): string {
      return funnelColors[index] || "bg-neutral-500";
    },

    clearData() {
      this.gaOverview = [];
      this.gaFunnel = [];
      this.error = null;
      this.lastFetched = null;
    },

    reset() {
      this.gaOverview = [];
      this.gaFunnel = [];
      this.isLoading = false;
      this.error = null;
      this.lastFetched = null;
    },
  },
});