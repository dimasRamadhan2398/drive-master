import { defineStore } from "pinia";

export interface Package {
  id: number;
  name: string;
  price: number;
  sessions: number;
  duration: number;
  description: string;
  features: string[];
  isActive: boolean;
  isPopular: boolean;
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
}

const initialPackages: Package[] = [
  {
    id: 1,
    name: "6x",
    price: 1750000,
    sessions: 6,
    duration: 60,
    description: "Our most popular package for comprehensive learning",
    features: ["Free Trial", "6x Sessions", "SIM A"],
    isActive: true,
    totalSold: 89,
  },
  {
    id: 2,
    name: "6x + Night Session",
    price: 1850000,
    sessions: 6,
    duration: 60,
    description: "Our most popular package for comprehensive learning",
    features: ["Free Trial", "6x Sessions", "SIM A"],
    isActive: true,
    totalSold: 89,
  },
  {
    id: 3,
    name: "6x + Weekend Session",
    price: 1850000,
    sessions: 6,
    duration: 60,
    description: "Our most popular package for comprehensive learning",
    features: ["Free Trial", "6x Sessions", "SIM A"],
    isActive: true,
    totalSold: 89,
  },
  {
    id: 4,
    name: "6x + Weekend & Night Session",
    price: 1950000,
    sessions: 6,
    duration: 60,
    description: "Our most popular package for comprehensive learning",
    features: ["Free Trial", "6x Sessions", "SIM A"],
    isActive: true,
    totalSold: 89,
  },
  {
    id: 5,
    name: "8x",
    price: 1950000,
    sessions: 8,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "8x Sessions", "SIM A"],
    isActive: true,
    isPopular: true,
    totalSold: 22,
  },
  {
    id: 6,
    name: "8x + Night Session",
    price: 2100000,
    sessions: 8,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "8x Sessions", "SIM A"],
    isActive: true,
    isPopular: true,
    totalSold: 22,
  },
  {
    id: 7,
    name: "8x + Weekend Session",
    price: 2100000,
    sessions: 8,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "8x Sessions", "SIM A"],
    isActive: true,
    isPopular: true,
    totalSold: 22,
  },
  {
    id: 8,
    name: "8x + Weekend & Night Session",
    price: 2250000,
    sessions: 8,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "8x Sessions", "SIM A"],
    isActive: true,
    isPopular: true,
    totalSold: 22,
  },
  {
    id: 9,
    name: "10x",
    price: 2250000,
    sessions: 10,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "10x Sessions", "SIM A"],
    isActive: true,
    totalSold: 22,
  },
  {
    id: 10,
    name: "10x + Night Session",
    price: 2450000,
    sessions: 10,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "10x Sessions", "SIM A"],
    isActive: true,
    totalSold: 22,
  },
  {
    id: 11,
    name: "10x + Weekend Session",
    price: 2450000,
    sessions: 10,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "10x Sessions", "SIM A"],
    isActive: true,
    totalSold: 22,
  },
  {
    id: 12,
    name: "10x + Weekend & Night Session",
    price: 2650000,
    sessions: 10,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "10x Sessions", "SIM A"],
    isActive: true,
    totalSold: 22,
  },
  {
    id: 13,
    name: "12x",
    price: 2650000,
    sessions: 12,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "12x Sessions"],
    isActive: true,
    totalSold: 22,
  },
  {
    id: 14,
    name: "12x + Night Session",
    price: 2900000,
    sessions: 12,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "12x Sessions", "SIM A"],
    isActive: true,
    totalSold: 22,
  },
  {
    id: 15,
    name: "12x + Weekend Session",
    price: 2900000,
    sessions: 12,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "12x Sessions", "SIM A"],
    isActive: true,
    totalSold: 22,
  },
  {
    id: 16,
    name: "12x + Weekend & Night Session",
    price: 3150000,
    sessions: 12,
    duration: 60,
    description: "Complete mastery with unlimited support",
    features: ["Free Trial", "12x Sessions", "SIM A"],
    isActive: true,
    totalSold: 22,
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
    packages: initialPackages,
    addons: initialAddons,
    isLoading: false,
    error: null,
  }),

  getters: {
    activePackages: (state) => state.packages.filter((p) => p.isActive),
    popularPackages: (state) => state.packages.filter((p) => p.isPopular),
    totalPackages: (state) => state.packages.length,
    totalSold: (state) =>
      state.packages.reduce((sum, p) => sum + p.totalSold, 0),
    totalRevenue: (state) =>
      state.packages.reduce((sum, p) => sum + p.price * p.totalSold, 0),
    getPackageById: (state) => (id: number) =>
      state.packages.find((p) => p.id === id),
  },

  actions: {
    addPackage(data: Omit<Package, "id" | "totalSold">) {
      const newId = Math.max(...this.packages.map((p) => p.id), 0) + 1;
      const newPackage: Package = {
        ...data,
        id: newId,
        totalSold: 0,
      };
      this.packages.push(newPackage);
      return newPackage;
    },

    updatePackage(id: number, data: Partial<Package>) {
      const index = this.packages.findIndex((p) => p.id === id);
      if (index !== -1) {
        this.packages[index] = { ...this.packages[index], ...data };
        return this.packages[index];
      }
      return null;
    },

    deletePackage(id: number) {
      this.packages = this.packages.filter((p) => p.id !== id);
    },

    duplicatePackage(id: number) {
      const pkg = this.packages.find((p) => p.id === id);
      if (pkg) {
        const newId = Math.max(...this.packages.map((p) => p.id), 0) + 1;
        const duplicated: Package = {
          ...pkg,
          id: newId,
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

    togglePackageStatus(id: number) {
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
        this.addons[index] = { ...this.addons[index], ...data };
        return this.addons[index];
      }
      return null;
    },

    deleteAddon(id: number) {
      this.addons = this.addons.filter((a) => a.id !== id);
    },

    recordSale(packageId: number) {
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