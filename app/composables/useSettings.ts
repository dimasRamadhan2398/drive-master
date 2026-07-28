export const useSettings = () => {
  // FITUR BARU: State global untuk jam operasional
  const operatingHours = useState('operating-hours', () => ({
    mondayStart: '08:00',
    mondayEnd: '18:00',
    weekendStart: '08:00',
    weekendEnd: '16:00',
    nightStart: '18:00',
    nightEnd: '20:00',
    nightEnabled: true,
    sundayClosed: true
  }))

  const promoEndDate = useCookie('promo-end-date', { default: () => '2026-05-31T23:59:59' })

  const settingsStore = useSettingsStore()

  // Computed WA link built from general settings whatsApp field
  const waLink = computed(() => {
    const number = settingsStore.generalSettings?.whatsApp?.replace(/\D/g, '') ?? '628119124848'
    const normalized = number.startsWith('0') ? `62${number.slice(1)}` : number
    return `https://wa.me/${normalized}?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi`
  })

  // Computed address from general settings
  const address = computed(() => settingsStore.generalSettings?.address ?? null)

  return {
    operatingHours,
    promoEndDate,
    generalSettings: computed(() => settingsStore.generalSettings),
    waLink,
    address,
    fetchGeneralSettings: () => settingsStore.fetchGeneralSettings(),
  }
}

