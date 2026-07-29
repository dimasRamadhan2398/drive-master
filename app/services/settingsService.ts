import type {
  ApiResponse,
} from "~/composables/useApiClients";

// General Settings Types (matching backend DTOs)
export interface GeneralSettings {
  id: string;
  businessName: string;
  email: string;
  phone: string;
  fax: string;
  whatsApp: string;
  instagram: string;
  youtube: string;
  mapDirection: string;
  address: string;
  hoursMonFri: string;
  hoursSatSun: string;
  hoursNightShift: string;
  promoEndDate: string | null;
  notifyEmail: boolean;
  notifySms: boolean;
  notifyWhatsApp: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateGeneralSettingsData {
  businessName?: string;
  email?: string;
  phone?: string;
  fax?: string;
  whatsApp?: string;
  instagram?: string;
  youtube?: string;
  mapDirection?: string;
  address?: string;
  hoursMonFri?: string;
  hoursSatSun?: string;
  hoursNightShift?: string;
  promoEndDate?: string | null;
  notifyEmail?: boolean;
  notifySms?: boolean;
  notifyWhatsApp?: boolean;
}

export const settingsService = {
  // GET /general-settings - Get general settings
  async fetchGeneralSettings(): Promise<GeneralSettings | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<GeneralSettings>>(
        "/general-settings",
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // PUT /general-settings - Update general settings
  async updateGeneralSettings(
    data: UpdateGeneralSettingsData,
  ): Promise<GeneralSettings | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<GeneralSettings>>(
        "/general-settings",
        {
          method: "PUT",
          body: data,
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },
};
