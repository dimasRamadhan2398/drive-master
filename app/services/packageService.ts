import type { Package } from "~/stores/packages";
import type { ApiResponse } from "~/composables/useApiClients";

export interface PaginatedPackagesResult {
  packages: Package[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface PackageApiResponse {
  id: string;
  name: string;
  description: string;
  packageType: "bronze" | "silver" | "gold" | "platinum";
  price: number;
  discountPrice: number;
  durationMinutes: number;
  totalSessions: number;
  features: string[];
  highlight: boolean;
  status: "active" | "inactive";
  imageUrl: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreatePackageData {
  name: string;
  description: string;
  packageType: "bronze" | "silver" | "gold" | "platinum";
  price: number;
  discountPrice?: number;
  durationMinutes: number;
  totalSessions: number;
  benefits: string[];
  highlight?: boolean;
  imageUrl?: string;
  isDiscounted: boolean;
}

export interface UpdatePackageData {
  name?: string;
  description?: string;
  packageType?: "bronze" | "silver" | "gold" | "platinum";
  price?: number;
  discountPrice?: number;
  durationMinutes?: number;
  totalSessions?: number;
  features?: string[];
  highlight?: boolean;
  status?: "active" | "inactive";
  imageUrl?: string;
}

export const mapApiToPackage = (item: PackageApiResponse): Package => {
  return {
    id: item.id,
    name: item.name,
    price: item.price,
    discountPrice: item.discountPrice,
    sessions: item.totalSessions,
    duration: item.durationMinutes,
    description: item.description,
    features: item.features || [],
    isActive: item.status === "active",
    isPopular: item.highlight || false,
    packageType: item.packageType,
    imageUrl: item.imageUrl || "",
    totalSold: 0, // Not available in API response
  };
};

export const packageService = {
  // Fetch all packages
  async fetchAll(): Promise<PaginatedPackagesResult> {
    const { core, extractData } = useApiClients();

    const response = await core<{
      success: boolean;
      message: string;
      data: PackageApiResponse[];
    }>("/packages/all?limit=100", {
      method: "GET",
    });

    const data = extractData(response);
    return {
      packages: Array.isArray(data) ? data.map(mapApiToPackage) : [],
      pagination: {
        page: 1,
        limit: data.length,
        total: data.length,
        totalPages: 1,
      },
    };
  },

  // Fetch all without pagination - returns flat array
  async fetchAllFlat(): Promise<Package[]> {
    const result = await this.fetchAll();
    return result.packages;
  },

  async fetchById(id: string): Promise<Package | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<PackageApiResponse>>(
        `/packages/${id}`,
        {
          method: "GET",
        },
      );
      return mapApiToPackage(extractData(response));
    } catch {
      return null;
    }
  },

  async create(data: CreatePackageData): Promise<Package> {
    const { core, extractData } = useApiClients();
    const response = await core<ApiResponse<PackageApiResponse>>(
      "/packages/create",
      {
        method: "POST",
        body: data,
      },
    );
    return mapApiToPackage(extractData(response));
  },

  async update(id: string, data: UpdatePackageData): Promise<Package | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<PackageApiResponse>>(
        `/packages/${id}`,
        {
          method: "PUT",
          body: data,
        },
      );
      return mapApiToPackage(extractData(response));
    } catch {
      return null;
    }
  },

  async delete(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/packages/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  async toggleStatus(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/packages/${id}/toggle-status`, { method: "PATCH" });
      return true;
    } catch {
      return false;
    }
  },
};
