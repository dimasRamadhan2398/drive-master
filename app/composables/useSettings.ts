export const useSettings = () => {
  const settingsStore = useSettingsStore();

  onMounted(() => {
    if (!settingsStore.isSettingsLoaded) {
      settingsStore.fetchGeneralSettings();
    }
  });

  const parseHours = (str: string | undefined, defaultStart: string, defaultEnd: string) => {
    if (!str) return { start: defaultStart, end: defaultEnd };
    const parts = str.split("-").map((p) => p.trim());
    if (parts.length === 2 && parts[0] && parts[1]) {
      return { start: parts[0], end: parts[1] };
    }
    return { start: defaultStart, end: defaultEnd };
  };

  // Operating hours dynamically synced with General Settings
  const operatingHours = computed(() => {
    const gs = settingsStore.generalSettings;
    const monFri = parseHours(gs?.hoursMonFri, "08:00", "18:00");
    const satSun = parseHours(gs?.hoursSatSun, "08:00", "16:00");
    const night = parseHours(gs?.hoursNightShift, "18:00", "20:00");

    return {
      mondayStart: monFri.start,
      mondayEnd: monFri.end,
      weekendStart: satSun.start,
      weekendEnd: satSun.end,
      nightStart: night.start,
      nightEnd: night.end,
      nightEnabled: true,
      sundayClosed: true,
    };
  });

  const promoEndDate = useCookie("promo-end-date", { default: () => "2026-05-31T23:59:59" });

  // Computed WA link built from general settings whatsApp field
  const waLink = computed(() => {
    const number = settingsStore.generalSettings?.whatsApp?.replace(/\D/g, "") ?? "628119124848";
    const normalized = number.startsWith("0") ? `62${number.slice(1)}` : number;
    return `https://wa.me/${normalized}?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi`;
  });

  // Computed address from general settings
  const address = computed(() => settingsStore.generalSettings?.address ?? null);

  return {
    operatingHours,
    promoEndDate,
    generalSettings: computed(() => settingsStore.generalSettings),
    waLink,
    address,
    fetchGeneralSettings: () => settingsStore.fetchGeneralSettings(),
  };
};

