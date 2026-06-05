import { defineStore } from "pinia";
import { packageService } from "~/services/packageService";

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
  id: number;
  name: string;
  price: number;
  description: string;
  sold: number;
}

interface PackagesState {
  packages: Package[];
  addons: Addon[];
  isLoading: boolean;
  error: string | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

const initialPackages: Package[] = [
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
    totalSold: 89,
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
    totalSold: 45,
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
    totalSold: 62,
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
    totalSold: 38,
  },
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
    totalSold: 22,
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
    totalSold: 15,
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
    totalSold: 18,
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
    totalSold: 12,
  },
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
    totalSold: 22,
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
    totalSold: 15,
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
    totalSold: 20,
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
    totalSold: 10,
  },
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
    totalSold: 22,
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
    totalSold: 18,
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
    totalSold: 25,
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
    totalSold: 30,
  },
];

const initialAddons: Addon[] = [
  {
    id: 1,
    name: "Extra Session",
    price: 350000,
    description: "Additional training session",
    sold: 34,
  },
];

export const usePackagesStore = defineStore("packages", {
  state: (): PackagesState => ({
    packages: [],
    addons: initialAddons,
    isLoading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0,
    },
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
    async fetchPackages() {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await packageService.fetchAll();

        this.packages = result.packages;
        this.pagination = result.pagination;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch packages";
        console.error("Error fetching packages:", err);
        // API failed, use dummy data
        this.packages = [...initialPackages];
      } finally {
        this.isLoading = false;
      }
    },
    addPackage(data: Omit<Package, "id" | "totalSold">) {
      const newPackage: Package = {
        ...data,
        id: crypto.randomUUID(),
        totalSold: 0,
      };
      this.packages.push(newPackage);
      return newPackage;
    },

    updatePackage(id: string, data: Partial<Package>) {
      const index = this.packages.findIndex((p) => p.id === id);
      if (index !== -1) {
        this.packages[index] = { ...this.packages[index], ...data } as Package;
        return this.packages[index];
      }
      return null;
    },

    deletePackage(id: string) {
      this.packages = this.packages.filter((p) => p.id !== id);
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

    togglePackageStatus(id: string) {
      const pkg = this.packages.find((p) => p.id === id);
      if (pkg) {
        pkg.isActive = !pkg.isActive;
        return pkg.isActive;
      }
      return null;
    },

    addAddon(data: Omit<Addon, "id" | "sold">) {
      const newId = Math.max(...this.addons.map((a) => a.id), 0) + 1;
      const newAddon: Addon = {
        ...data,
        id: newId,
        sold: 0,
      };
      this.addons.push(newAddon);
      return newAddon;
    },

    updateAddon(id: number, data: Partial<Addon>) {
      const index = this.addons.findIndex((a) => a.id === id);
      if (index !== -1) {
        this.addons[index] = { ...this.addons[index], ...data } as Addon;
        return this.addons[index];
      }
      return null;
    },

    deleteAddon(id: number) {
      this.addons = this.addons.filter((a) => a.id !== id);
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
