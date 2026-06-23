import { defineStore } from "pinia";
import { settingsService } from "~/services/settingsService";
import type { GeneralSettings, UpdateGeneralSettingsData } from "~/services/settingsService";

interface SettingsState {
  generalSettings: GeneralSettings | null;
  isLoading: boolean;
  error: string | null;
}

export const useSettingsStore = defineStore("settings", {
  state: (): SettingsState => ({
    generalSettings: null,
    isLoading: false,
    error: null,
  }),

  getters: {
    getGeneralSettings: (state) => state.generalSettings,
    isSettingsLoaded: (state) => state.generalSettings !== null,
  },

  actions: {
    async fetchGeneralSettings() {
      this.isLoading = true;
      this.error = null;
      try {
        const settings = await settingsService.fetchGeneralSettings();
        if (settings) {
          this.generalSettings = settings;
        }
        return settings;
      } catch (err) {
        this.error = err instanceof Error ? err.message : "Failed to fetch general settings";
        console.error("Error fetching general settings:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async updateGeneralSettings(data: UpdateGeneralSettingsData) {
      this.isLoading = true;
      this.error = null;
      try {
        const updated = await settingsService.updateGeneralSettings(data);
        if (updated) {
          this.generalSettings = updated;
          return updated;
        }
        return null;
      } catch (err) {
        this.error = err instanceof Error ? err.message : "Failed to update general settings";
        console.error("Error updating general settings:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    reset() {
      this.generalSettings = null;
      this.isLoading = false;
      this.error = null;
    },
  },
});
