import type {
  ApiResponse,
  PaginatedResponse,
} from "~/composables/useApiClients";

// Vehicle Status Types
export type VehicleStatus = "available" | "in_use" | "maintenance" | "retired";

// Transmission Types
export type TransmissionType = "manual" | "automatic";

// Vehicle Types (matching backend DTOs)
export interface Car {
  id: string;
  brand: string;
  model: string;
  year: number;
  licensePlate: string;
  color: string;
  transmission: TransmissionType;
  status: VehicleStatus;
  mileage: number;
  imageUrl: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

// Request DTOs (matching backend)
export interface CreateCarData {
  brand: string;
  model: string;
  year: number;
  licensePlate?: string;
  color?: string;
  transmission?: TransmissionType;
  status?: VehicleStatus;
  imageUrl?: string;
  imageBase64?: string;
  notes?: string;
}

export interface UpdateCarData {
  brand?: string;
  model?: string;
  year?: number;
  licensePlate?: string;
  color?: string;
  transmission?: TransmissionType;
  status?: VehicleStatus;
  mileage?: number;
  imageUrl?: string;
  imageBase64?: string;
  notes?: string;
}

// Helper to process image data - if starts with "http", use imageUrl, otherwise use imageBase64
export function processCarImage(imageData?: string): { imageUrl?: string; imageBase64?: string } | undefined {
  if (!imageData) return undefined;

  if (imageData.startsWith("http")) {
    return { imageUrl: imageData };
  }
  return { imageBase64: imageData };
}

export interface VehicleFilterParams {
  page?: number;
  limit?: number;
  status?: VehicleStatus;
  transmission?: TransmissionType;
  brand?: string;
  search?: string;
}

// Helper to get status color for UI
export const getVehicleStatusColor = (
  status: VehicleStatus,
): "success" | "info" | "warning" | "neutral" | "error" => {
  const colorMap: Record<VehicleStatus, "success" | "info" | "warning" | "neutral" | "error"> = {
    available: "success",
    in_use: "info",
    maintenance: "warning",
    retired: "neutral",
  };
  return colorMap[status] || "neutral";
};

// Helper to get status label for UI
export const getVehicleStatusLabel = (status: VehicleStatus): string => {
  const labelMap: Record<VehicleStatus, string> = {
    available: "Available",
    in_use: "In Use",
    maintenance: "Maintenance",
    retired: "Retired",
  };
  return labelMap[status] || status;
};

export const vehicleService = {
  // ==================== VEHICLE CRUD METHODS ====================

  // GET /cars - Get all vehicles
  async fetchCars(
    params: VehicleFilterParams = {},
  ): Promise<{ cars: Car[]; total: number; page: number; limit: number; totalPages: number }> {
    const { core, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.status) queryParams.set("status", params.status);
    if (params.transmission)
      queryParams.set("transmission", params.transmission);
    if (params.brand) queryParams.set("brand", params.brand);
    if (params.search) queryParams.set("search", params.search);

    const queryString = queryParams.toString();
    const url = `/cars${queryString ? `?${queryString}` : ""}`;

    try {
      const response = await core<PaginatedResponse<Car>>(url, {
        method: "GET",
      });
      const { data, pagination } = extractPaginatedData(response);
      return {
        cars: Array.isArray(data) ? data : [],
        ...pagination,
      };
    } catch {
      return { cars: [], total: 0, page: 1, limit: 100, totalPages: 0 };
    }
  },

  // GET /cars/:id - Get vehicle by ID
  async fetchCarById(id: string): Promise<Car | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Car>>(`/cars/${id}`, {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /cars - Create a new vehicle
  async createCar(data: CreateCarData): Promise<Car | null> {
    const { core, extractData } = useApiClients();
    try {
      // Process image data - use imageUrl for http URLs, imageBase64 for base64 strings
      const imageResult = processCarImage(data.imageUrl || data.imageBase64);
      const requestBody = {
        ...data,
        ...imageResult,
      };

      const response = await core<ApiResponse<Car>>("/cars", {
        method: "POST",
        body: requestBody,
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // PUT /cars/:id - Update vehicle
  async updateCar(id: string, data: UpdateCarData): Promise<Car | null> {
    const { core, extractData } = useApiClients();
    try {
      // Process image data - use imageUrl for http URLs, imageBase64 for base64 strings
      const imageResult = processCarImage(data.imageUrl || data.imageBase64);
      const requestBody = {
        ...data,
        ...imageResult,
      };

      const response = await core<ApiResponse<Car>>(`/cars/${id}`, {
        method: "PUT",
        body: requestBody,
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // DELETE /cars/:id - Delete vehicle
  async deleteCar(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/cars/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },
};
