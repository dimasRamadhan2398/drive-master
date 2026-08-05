import { defineStore } from "pinia";
import { packageService } from "~/services/packageService";
import { addonService } from "~/services/addonService";
import type {
  CreatePackageData,
  UpdatePackageData,
} from "~/services/packageService";

export interface Package {
  id: string; // UUID from API
  name: string;
  price: number;
  discountPrice: number;
  sessions: number;
  duration: number; // in minutes
  description: string;
  features: string[];
  isActive: boolean;
  isPopular: boolean; // mapped from 'highlight' in API
  packageType: "bronze" | "silver" | "gold" | "platinum";
  imageUrl: string;
  totalSold: number;
}

export interface Addon {
  id: string; // UUID from API
  name: string; // mapped from 'title' in API
  price: number;
  description: string;
  sold: number;
  sessions?: number;
  status?: "active" | "inactive";
  imageUrl?: string;
  sortOrder?: number;
}

interface PackagesState {
  packages: Package[];
  addons: Addon[];
  isLoading: boolean;
  isAddLoading: boolean;
  error: string | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
  filteredByDate: Package[];
}

const initialPackages: Package[] = [
  // 6 sessions (bronze)
  {
    id: "11111111-1111-1111-1111-111111111101",
    name: "6x",
    price: 1500000,
    discountPrice: 1350000,
    sessions: 6,
    duration: 90,
    description: "6 training sessions with SIM A certification",
    features: ["Free Trial", "6 training sessions", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "bronze",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111102",
    name: "6x + Night Session",
    price: 1750000,
    discountPrice: 1600000,
    sessions: 6,
    duration: 90,
    description:
      "6 training sessions including night driving with SIM A certification",
    features: ["Free Trial", "6 training sessions", "Night Session", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "bronze",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111103",
    name: "6x + Weekend Session",
    price: 1850000,
    discountPrice: 1700000,
    sessions: 6,
    duration: 90,
    description:
      "6 training sessions including weekend sessions with SIM A certification",
    features: ["Free Trial", "6 training sessions", "Weekend Session", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "bronze",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111104",
    name: "6x + Weekend & Night Session",
    price: 1950000,
    discountPrice: 1750000,
    sessions: 6,
    duration: 90,
    description:
      "6 training sessions including weekend and night sessions with SIM A certification",
    features: [
      "Free Trial",
      "6 training sessions",
      "Weekend Session",
      "Night Session",
      "SIM A",
    ],
    isActive: true,
    isPopular: false,
    packageType: "bronze",
    imageUrl: "",
    totalSold: 0,
  },
  // 8 sessions (silver)
  {
    id: "11111111-1111-1111-1111-111111111201",
    name: "8x",
    price: 1950000,
    discountPrice: 1750000,
    sessions: 8,
    duration: 90,
    description: "8 training sessions with SIM A certification",
    features: ["Free Trial", "8 training sessions", "SIM A"],
    isActive: true,
    isPopular: true,
    packageType: "silver",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111202",
    name: "8x + Night Session",
    price: 2100000,
    discountPrice: 1900000,
    sessions: 8,
    duration: 90,
    description:
      "8 training sessions including night driving with SIM A certification",
    features: ["Free Trial", "8 training sessions", "Night Session", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "silver",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111203",
    name: "8x + Weekend Session",
    price: 2100000,
    discountPrice: 1900000,
    sessions: 8,
    duration: 90,
    description:
      "8 training sessions including weekend sessions with SIM A certification",
    features: ["Free Trial", "8 training sessions", "Weekend Session", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "silver",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111204",
    name: "8x + Weekend & Night Session",
    price: 2250000,
    discountPrice: 2050000,
    sessions: 8,
    duration: 90,
    description:
      "8 training sessions including weekend and night sessions with SIM A certification",
    features: [
      "Free Trial",
      "8 training sessions",
      "Weekend Session",
      "Night Session",
      "SIM A",
    ],
    isActive: true,
    isPopular: false,
    packageType: "silver",
    imageUrl: "",
    totalSold: 0,
  },
  // 10 sessions (gold)
  {
    id: "11111111-1111-1111-1111-111111111301",
    name: "10x",
    price: 2250000,
    discountPrice: 2050000,
    sessions: 10,
    duration: 90,
    description: "10 training sessions with SIM A certification",
    features: ["Free Trial", "10 training sessions", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "gold",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111302",
    name: "10x + Night Session",
    price: 2450000,
    discountPrice: 2250000,
    sessions: 10,
    duration: 90,
    description:
      "10 training sessions including night driving with SIM A certification",
    features: ["Free Trial", "10 training sessions", "Night Session", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "gold",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111303",
    name: "10x + Weekend Session",
    price: 2450000,
    discountPrice: 2250000,
    sessions: 10,
    duration: 90,
    description:
      "10 training sessions including weekend sessions with SIM A certification",
    features: [
      "Free Trial",
      "10 training sessions",
      "Weekend Session",
      "SIM A",
    ],
    isActive: true,
    isPopular: false,
    packageType: "gold",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111304",
    name: "10x + Weekend & Night Session",
    price: 2650000,
    discountPrice: 2450000,
    sessions: 10,
    duration: 90,
    description:
      "10 training sessions including weekend and night sessions with SIM A certification",
    features: [
      "Free Trial",
      "10 training sessions",
      "Weekend Session",
      "Night Session",
      "SIM A",
    ],
    isActive: true,
    isPopular: false,
    packageType: "gold",
    imageUrl: "",
    totalSold: 0,
  },
  // 12 sessions (platinum)
  {
    id: "11111111-1111-1111-1111-111111111401",
    name: "12x",
    price: 2650000,
    discountPrice: 2450000,
    sessions: 12,
    duration: 90,
    description: "12 training sessions with SIM A certification",
    features: ["Free Trial", "12 training sessions", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "platinum",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111402",
    name: "12x + Night Session",
    price: 2900000,
    discountPrice: 2650000,
    sessions: 12,
    duration: 90,
    description:
      "12 training sessions including night driving with SIM A certification",
    features: ["Free Trial", "12 training sessions", "Night Session", "SIM A"],
    isActive: true,
    isPopular: false,
    packageType: "platinum",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111403",
    name: "12x + Weekend Session",
    price: 2900000,
    discountPrice: 2650000,
    sessions: 12,
    duration: 90,
    description:
      "12 training sessions including weekend sessions with SIM A certification",
    features: [
      "Free Trial",
      "12 training sessions",
      "Weekend Session",
      "SIM A",
    ],
    isActive: true,
    isPopular: false,
    packageType: "platinum",
    imageUrl: "",
    totalSold: 0,
  },
  {
    id: "11111111-1111-1111-1111-111111111404",
    name: "12x + Weekend & Night Session",
    price: 3150000,
    discountPrice: 2900000,
    sessions: 12,
    duration: 90,
    description:
      "12 training sessions including weekend and night sessions with SIM A certification",
    features: [
      "Free Trial",
      "12 training sessions",
      "Weekend Session",
      "Night Session",
      "SIM A",
    ],
    isActive: true,
    isPopular: false,
    packageType: "platinum",
    imageUrl: "",
    totalSold: 0,
  },
];

const initialAddons: Addon[] = [
  {
    id: "00000000-0000-0000-0000-000000000001",
    name: "Extra Session",
    price: 350000,
    description: "Additional training session (90 mins)",
    sold: 0,
    sessions: 1,
    status: "active",
  },
  {
    id: "00000000-0000-0000-0000-000000000002",
    name: "SIM A Express Processing",
    price: 750000,
    description: "Priority handling and administrative assistant for SIM A",
    sold: 0,
    sessions: 0,
    status: "active",
  },
  {
    id: "00000000-0000-0000-0000-000000000003",
    name: "Night Driving Practice",
    price: 450000,
    description: "Special night driving coaching with high visibility vehicles",
    sold: 0,
    sessions: 2,
    status: "active",
  },
  {
    id: "00000000-0000-0000-0000-000000000004",
    name: "Highway & Toll Driving Module",
    price: 500000,
    description: "Advanced highway speed & toll lane maneuvering practice",
    sold: 0,
    sessions: 2,
    status: "active",
  },
];

// Helper to get today's date in YYYY-MM-DD format
const getTodayDate = (): string => {
  const today = new Date();
  const year = today.getFullYear();
  const month = (today.getMonth() + 1).toString().padStart(2, "0");
  const day = today.getDate().toString().padStart(2, "0");
  return `${year}-${month}-${day}`;
};

export const usePackagesStore = defineStore("packages", {
  state: (): PackagesState => ({
    packages: [],
    addons: [],
    isLoading: false,
    isAddLoading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0,
    },
    filteredByDate: [],
  }),

  getters: {
    activePackages: (state) => state.packages.filter((p) => p.isActive),
    popularPackages: (state) => state.packages.filter((p) => p.isPopular),
    totalPackages: (state) => state.packages.length,
    totalRevenue: (state) =>
      state.packages.reduce((sum, p) => sum + p.discountPrice * p.totalSold, 0),
    getPackageById: (state) => (id: string) =>
      state.packages.find((p) => p.id === id),
    getPackagesByType: (state) => (type: Package["packageType"]) =>
      state.packages.filter((p) => p.packageType === type),
  },

  actions: {
    // Helper to sort packages by sessions (6 to 12)
    sortPackagesBySessions(packages: Package[]): Package[] {
      return [...packages].sort((a, b) => {
        // First sort by sessions
        if (a.sessions !== b.sessions) {
          return a.sessions - b.sessions;
        }
        // Then sort by price within the same session count
        return a.discountPrice - b.discountPrice;
      });
    },

    async fetchPackages(params?: { status?: string }) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await packageService.fetchAll(params);

        this.packages = this.sortPackagesBySessions(result.packages);
        this.pagination = result.pagination;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch packages";
        console.error("Error fetching packages:", err);
        // API failed, use dummy data (already sorted)
        this.packages = [...initialPackages];
      } finally {
        this.isLoading = false;
      }
    },
    async addPackage(data: CreatePackageData) {
      this.isAddLoading = true;
      try {
        const newPackage = await packageService.create(data);
        this.packages.unshift(newPackage);
        return newPackage;
      } catch {
        // Fallback to local creation if API fails
        return this.addPackageLocal({
          name: data.name,
          price: data.price,
          discountPrice: data.discountPrice ?? data.price,
          sessions: data.totalSessions,
          duration: data.durationMinutes,
          description: data.description,
          features: data.benefits,
          isActive: true,
          isPopular: data.highlight ?? false,
          packageType: data.packageType,
          imageUrl: data.imageUrl ?? "",
        });
      } finally {
        this.isAddLoading = false;
      }
    },

    // Local creation fallback (when API is not available)
    addPackageLocal(data: Omit<Package, "id" | "totalSold">) {
      const newPackage: Package = {
        ...data,
        id: crypto.randomUUID(),
        totalSold: 0,
      };
      this.packages.push(newPackage);
      return newPackage;
    },

    async updatePackage(id: string, data: UpdatePackageData) {
      try {
        const updatedPackage = await packageService.update(id, data);
        if (updatedPackage) {
          const index = this.packages.findIndex((p) => p.id === id);
          if (index !== -1) {
            this.packages[index] = updatedPackage;
          }
          return updatedPackage;
        }
        return null;
      } catch {
        // Fallback to local update if API fails
        return this.updatePackageLocal(id, data);
      }
    },

    // Local update fallback (when API is not available)
    updatePackageLocal(id: string, data: Partial<Package>): Package | null {
      const index = this.packages.findIndex((p) => p.id === id);
      if (index !== -1) {
        this.packages[index] = { ...this.packages[index], ...data } as Package;
        return this.packages[index];
      }
      return null;
    },

    async deletePackage(id: string) {
      try {
        const success = await packageService.delete(id);
        if (success) {
          this.packages = this.packages.filter((p) => p.id !== id);
        }
        return success;
      } catch {
        // Fallback to local delete if API fails
        this.packages = this.packages.filter((p) => p.id !== id);
        return true;
      }
    },

    duplicatePackage(id: string) {
      const pkg = this.packages.find((p) => p.id === id);
      if (pkg) {
        const duplicated: Package = {
          ...pkg,
          id: crypto.randomUUID(),
          name: `${pkg.name} (Copy)`,
          isPopular: false,
          isActive: true,
          totalSold: 0,
        };
        this.packages.push(duplicated);
        return duplicated;
      }
      return null;
    },

    async togglePackageStatus(id: string) {
      try {
        const success = await packageService.toggleStatus(id);
        if (success) {
          const pkg = this.packages.find((p) => p.id === id);
          if (pkg) {
            pkg.isActive = !pkg.isActive;
          }
        }
      } catch {
        // Fallback to local toggle if API fails
        const pkg = this.packages.find((p) => p.id === id);
        if (pkg) {
          pkg.isActive = !pkg.isActive;
        }
      }
    },

    async fetchAddons() {
      this.isLoading = true;
      this.error = null;

      try {
        const data = await addonService.fetchAll();
        this.addons = data;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch addons";
        console.error("Error fetching addons:", err);
        // API failed, use initial data
        this.addons = [...initialAddons];
      } finally {
        this.isLoading = false;
      }
    },

    async addAddon(data: { name: string; price: number; description?: string; sessions?: number }) {
      this.isAddLoading = true;
      try {
        const newAddon = await addonService.create({
          title: data.name,
          price: data.price,
          description: data.description || "",
          sessions: data.sessions || 1,
        });
        this.addons.unshift(newAddon);
        return newAddon;
      } catch {
        // Fallback to local creation if API fails
        return this.addAddonLocal(data);
      } finally {
        this.isAddLoading = false;
      }
    },

    // Local creation fallback (when API is not available)
    addAddonLocal(data: { name: string; price: number; description?: string; sessions?: number }) {
      const newAddon: Addon = {
        id: crypto.randomUUID(),
        name: data.name,
        price: data.price,
        description: data.description || "",
        sold: 0,
        sessions: data.sessions || 1,
        status: "active",
      };
      this.addons.push(newAddon);
      return newAddon;
    },

    async updateAddon(id: string, data: { name?: string; price?: number; description?: string; sessions?: number; status?: "active" | "inactive" }) {
      try {
        const updatedAddon = await addonService.update(id, {
          title: data.name,
          description: data.description,
          price: data.price,
          sessions: data.sessions,
          status: data.status,
        });
        if (updatedAddon) {
          const index = this.addons.findIndex((a) => a.id === id);
          if (index !== -1) {
            this.addons[index] = updatedAddon;
          }
          return updatedAddon;
        }
        return null;
      } catch {
        // Fallback to local update if API fails
        return this.updateAddonLocal(id, data);
      }
    },

    // Local update fallback (when API is not available)
    updateAddonLocal(id: string, data: Partial<Addon>): Addon | null {
      const index = this.addons.findIndex((a) => a.id === id);
      if (index !== -1) {
        this.addons[index] = { ...this.addons[index], ...data } as Addon;
        return this.addons[index];
      }
      return null;
    },

    async deleteAddon(id: string) {
      try {
        const success = await addonService.delete(id);
        if (success) {
          this.addons = this.addons.filter((a) => a.id !== id);
        }
        return success;
      } catch {
        // Fallback to local delete if API fails
        this.addons = this.addons.filter((a) => a.id !== id);
        return true;
      }
    },

    async toggleAddonStatus(id: string) {
      try {
        const success = await addonService.toggleStatus(id);
        if (success) {
          const addon = this.addons.find((a) => a.id === id);
          if (addon) {
            addon.status = addon.status === "active" ? "inactive" : "active";
          }
        }
        return success;
      } catch {
        // Fallback to local toggle if API fails
        const addon = this.addons.find((a) => a.id === id);
        if (addon) {
          addon.status = addon.status === "active" ? "inactive" : "active";
        }
        return true;
      }
    },

    recordSale(packageId: string) {
      const pkg = this.packages.find((p) => p.id === packageId);
      if (pkg) {
        pkg.totalSold++;
      }
    },

    reset() {
      this.packages = [...initialPackages];
      this.addons = [...initialAddons];
      this.isLoading = false;
      this.error = null;
    },
  },
});
