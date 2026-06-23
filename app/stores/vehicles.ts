import { defineStore } from "pinia";
import { vehicleService } from "~/services/vehicleService";
import type {
  Car,
  CreateCarData,
  UpdateCarData,
  VehicleFilterParams,
  VehicleStatus,
} from "~/services/vehicleService";

interface VehiclesState {
  vehicles: Car[];
  isLoading: boolean;
  error: string | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export const useVehiclesStore = defineStore("vehicles", {
  state: (): VehiclesState => ({
    vehicles: [],
    isLoading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 100,
      total: 0,
      totalPages: 0,
    },
  }),

  getters: {
    // Get all vehicles
    allVehicles: (state) => state.vehicles,

    // Get available vehicles
    availableVehicles: (state) =>
      state.vehicles.filter((v) => v.status === "available"),

    // Get vehicles by status
    vehiclesByStatus: (state) => (status: VehicleStatus) =>
      state.vehicles.filter((v) => v.status === status),

    // Get vehicle by ID
    getVehicleById: (state) => (id: string) =>
      state.vehicles.find((v) => v.id === id),

    // Get vehicle options for dropdown (name: "Brand Model")
    vehicleOptions: (state) =>
      state.vehicles
        .filter((v) => v.status === "available")
        .map((v) => ({
          label: `${v.brand} ${v.model}`,
          value: v.id,
        })),

    // Check if vehicles are loaded
    isLoaded: (state) => state.vehicles.length > 0,
  },

  actions: {
    // ==================== DATA FETCHING ====================

    async fetchVehicles(params?: VehicleFilterParams) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await vehicleService.fetchCars({
          page: this.pagination.page,
          limit: this.pagination.limit,
          ...params,
        });

        this.vehicles = result.cars;
        this.pagination = {
          page: result.page,
          limit: result.limit,
          total: result.total,
          totalPages: result.totalPages,
        };
        return result.cars;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch vehicles";
        console.error("Error fetching vehicles:", err);
        return [];
      } finally {
        this.isLoading = false;
      }
    },

    async fetchVehicleById(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const vehicle = await vehicleService.fetchCarById(id);
        return vehicle;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch vehicle";
        console.error("Error fetching vehicle:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    // ==================== CRUD OPERATIONS ====================

    async createVehicle(data: CreateCarData) {
      this.isLoading = true;
      this.error = null;

      try {
        const vehicle = await vehicleService.createCar(data);
        if (vehicle) {
          this.vehicles.unshift(vehicle);
          return vehicle;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to create vehicle";
        console.error("Error creating vehicle:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async updateVehicle(id: string, data: UpdateCarData) {
      this.isLoading = true;
      this.error = null;

      try {
        const vehicle = await vehicleService.updateCar(id, data);
        if (vehicle) {
          const index = this.vehicles.findIndex((v) => v.id === id);
          if (index !== -1) {
            this.vehicles[index] = vehicle;
          }
          return vehicle;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to update vehicle";
        console.error("Error updating vehicle:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async deleteVehicle(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const success = await vehicleService.deleteCar(id);
        if (success) {
          this.vehicles = this.vehicles.filter((v) => v.id !== id);
        }
        return success;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to delete vehicle";
        console.error("Error deleting vehicle:", err);
        return false;
      } finally {
        this.isLoading = false;
      }
    },

    // ==================== UTILITY OPERATIONS ====================

    // Set pagination
    setPage(page: number) {
      this.pagination.page = page;
    },

    // Set filters and refetch
    async setFilters(params: VehicleFilterParams) {
      await this.fetchVehicles(params);
    },

    // Reset state
    reset() {
      this.vehicles = [];
      this.isLoading = false;
      this.error = null;
      this.pagination = {
        page: 1,
        limit: 100,
        total: 0,
        totalPages: 0,
      };
    },
  },
});
