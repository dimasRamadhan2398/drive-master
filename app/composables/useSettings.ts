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

  return {
    operatingHours
  }
}
