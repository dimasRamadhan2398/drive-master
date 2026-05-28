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

  return {
    operatingHours,
    promoEndDate
  }
}
