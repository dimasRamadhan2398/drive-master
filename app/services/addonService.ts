import type { Addon } from "~/stores/packages";
import type { ApiResponse } from "~/composables/useApiClients";

export interface AddonApiResponse {
  id: string;
  title: string;
  description: string;
  price: number;
  sessions: number;
  status: "active" | "inactive";
  imageUrl: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAddonData {
  title: string;
  description?: string;
  price: number;
  sessions?: number;
  imageUrl?: string;
  sortOrder?: number;
}

export interface UpdateAddonData {
  title?: string;
  description?: string;
  price?: number;
  sessions?: number;
  status?: "active" | "inactive";
  imageUrl?: string;
  sortOrder?: number;
}

export const mapApiToAddon = (item: AddonApiResponse): Addon => {
  return {
    id: item.id,
    name: item.title,
    price: item.price,
    description: item.description,
    sold: 0, // Not available in API response
    sessions: item.sessions,
    status: item.status,
    imageUrl: item.imageUrl,
    sortOrder: item.sortOrder,
  };
};

export const addonService = {
  // Fetch all addons
  async fetchAll(): Promise<Addon[]> {
    const { core, extractData } = useApiClients();

    const response = await core<{
      success: boolean;
      message: string;
      data: AddonApiResponse[];
    }>("/addons/all", {
      method: "GET",
    });

    const data = extractData(response);
    return Array.isArray(data) ? data.map(mapApiToAddon) : [];
  },

  // Fetch active addons
  async fetchActive(): Promise<Addon[]> {
    const { core, extractData } = useApiClients();

    const response = await core<{
      success: boolean;
      message: string;
      data: AddonApiResponse[];
    }>("/addons/active", {
      method: "GET",
    });

    const data = extractData(response);
    return Array.isArray(data) ? data.map(mapApiToAddon) : [];
  },

  async fetchById(id: string): Promise<Addon | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<AddonApiResponse>>(
        `/addons/${id}`,
        {
          method: "GET",
        },
      );
      return mapApiToAddon(extractData(response));
    } catch {
      return null;
    }
  },

  async create(data: CreateAddonData): Promise<Addon> {
    const { core, extractData } = useApiClients();
    const response = await core<ApiResponse<AddonApiResponse>>(
      "/addons/create",
      {
        method: "POST",
        body: data,
      },
    );
    return mapApiToAddon(extractData(response));
  },

  async update(id: string, data: UpdateAddonData): Promise<Addon | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<AddonApiResponse>>(
        `/addons/${id}`,
        {
          method: "PUT",
          body: data,
        },
      );
      return mapApiToAddon(extractData(response));
    } catch {
      return null;
    }
  },

  async delete(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/addons/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  async toggleStatus(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/addons/toggle-status/${id}`, { method: "PUT" });
      return true;
    } catch {
      return false;
    }
  },
};
