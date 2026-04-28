export interface TimeSlot {
  id: string
  date: string // FITUR BARU: Menambahkan field tanggal (YYYY-MM-DD)
  time: string
  duration: string
  car: string
  instructor: string
  student: string | null
  status: 'available' | 'booked' | 'in-progress' | 'completed' | 'blocked'
}

// PERUBAHAN: Fungsi generator sekarang menerima operatingHours agar datanya sinkron
const generateMockData = (hours: any): TimeSlot[] => {
  const result: TimeSlot[] = []
  const year = 2026
  const month = 4 // April
  const daysInMonth = 30 // April memiliki 30 hari
  
  const baseTimes = ['07:00', '08:30', '10:00', '11:30', '13:00', '14:30', '16:00', '17:30', '19:00']
  const instructors = ['Pak Ahmad', 'Bu Sari', 'Pak Budi']
  const cars = ['Tesla Model 3', 'BYD Atto 3']
  const statuses: ('available' | 'available' | 'booked' | 'blocked')[] = ['available', 'available', 'booked', 'blocked']
  
  let idCounter = 1

  for (let day = 1; day <= daysInMonth; day++) {
    // Cek apakah hari ini weekend atau weekday
    const dateObj = new Date(year, month - 1, day)
    const dayOfWeek = dateObj.getDay() // 0 = Sunday, 6 = Saturday
    const isWeekend = dayOfWeek === 0 || dayOfWeek === 6
    
    // Jika Minggu dan disetel tutup, lewati
    if (dayOfWeek === 0 && hours.sundayClosed) continue

    const start = isWeekend ? hours.weekendStart : hours.mondayStart
    const end = isWeekend ? hours.weekendEnd : hours.mondayEnd
    const nightStart = hours.nightStart
    const nightEnd = hours.nightEnd

    // Filter baseTimes yang masuk dalam rentang jam operasional (Day atau Night)
    const allowedTimes = baseTimes.filter(t => {
      const isDay = t >= start && t <= end
      const isNight = hours.nightEnabled && t >= nightStart && t <= nightEnd
      return isDay || isNight
    })
    if (allowedTimes.length === 0) continue

    // Bikin variasi jumlah jadwal per hari
    const slotsCount = Math.min(allowedTimes.length, Math.floor(Math.random() * 3) + 3)
    const shuffledTimes = [...allowedTimes].sort(() => 0.5 - Math.random())
    
    for (let i = 0; i < slotsCount; i++) {
      const dateStr = `${year}-04-${day.toString().padStart(2, '0')}`
      const status = statuses[Math.floor(Math.random() * statuses.length)]
      
      result.push({
        id: idCounter.toString(),
        date: dateStr,
        time: shuffledTimes[i],
        duration: '60 min',
        car: cars[Math.floor(Math.random() * cars.length)],
        instructor: instructors[Math.floor(Math.random() * instructors.length)],
        student: status === 'booked' ? `Siswa ${idCounter}` : null,
        status: status
      })
      idCounter++
    }
  }
  
  return result.sort((a, b) => a.date.localeCompare(b.date) || a.time.localeCompare(b.time))
}

export const useSchedules = () => {
  const { operatingHours } = useSettings()
  const slots = useState<TimeSlot[]>('global-schedules', () => generateMockData(operatingHours.value))

  const toggleSlotStatus = (id: string) => {
    // FITUR BARU: Memungkinkan memblokir status apa pun, dan unblock kembali ke available
    const slot = slots.value.find(s => s.id === id)
    if (slot) {
      if (slot.status === 'blocked') {
        slot.status = 'available'
      } else {
        slot.status = 'blocked'
      }
    }
  }

  // FITUR BARU: Fungsi untuk menambahkan slot baru
  const addSlot = (newSlot: TimeSlot) => {
    slots.value.push(newSlot)
  }

  // FITUR BARU: Fungsi untuk mengubah data slot yang ada
  const editSlot = (id: string, updatedData: Partial<TimeSlot>) => {
    const slot = slots.value.find(s => s.id === id)
    if (slot) {
      Object.assign(slot, updatedData)
    }
  }

  const deleteSlot = (id: string) => {
    slots.value = slots.value.filter(s => s.id !== id)
  }

  // FITUR BARU: Fungsi untuk mengubah status secara manual (Start/Complete session)
  const updateSlotStatus = (id: string, status: TimeSlot['status']) => {
    const slot = slots.value.find(s => s.id === id)
    if (slot) {
      slot.status = status
    }
  }

  const bookSlot = (id: string, studentName: string = 'John Doe') => {
    const slot = slots.value.find(s => s.id === id)
    if (slot && slot.status === 'available') {
      slot.status = 'booked'
      slot.student = studentName
    }
  }

  return {
    slots,
    toggleSlotStatus,
    addSlot, // <-- Ditambahkan ke return
    editSlot, // <-- Ditambahkan ke return
    deleteSlot,
    updateSlotStatus, // <-- Ditambahkan ke return
    bookSlot
  }
}

// FITUR BARU: Singleton monitor untuk Smart Alerts
let monitorInterval: any = null

export const useSmartAlerts = () => {
  const { slots, updateSlotStatus } = useSchedules()
  const notifiedSlots = useState<string[]>('notified-slots', () => [])
  const activeAlerts = useState<TimeSlot[]>('active-alerts', () => [])
  const toast = useToast()

  const checkUpcomingSlots = () => {
    const now = new Date()
    const todayStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`

    slots.value.forEach(slot => {
      // Hanya cek yang sudah di-booked, belum dinotifikasi, dan untuk hari ini
      if (slot.status === 'booked' && slot.date === todayStr) {
        const [h, m] = slot.time.split(':').map(Number)
        const slotTime = new Date()
        slotTime.setHours(h, m, 0, 0)
        
        const diffMinutes = (slotTime.getTime() - now.getTime()) / (1000 * 60)
        
        // Jika kursus akan dimulai dalam 15 menit
        if (diffMinutes > 0 && diffMinutes <= 15) {
          // Tambah ke activeAlerts jika belum ada
          if (!activeAlerts.value.some(a => a.id === slot.id)) {
            activeAlerts.value.push(slot)
          }

          // Kirim Toast hanya sekali
          if (!notifiedSlots.value.includes(slot.id)) {
            notifiedSlots.value.push(slot.id)
            toast.add({
              title: 'Upcoming Session!',
              description: `Kursus ${slot.student} akan dimulai pukul ${slot.time}. Apakah sudah siap?`,
              icon: 'i-lucide-bell-ring',
              color: 'warning',
              timeout: 0, 
              actions: [
                { 
                  label: 'Start Now', 
                  color: 'warning',
                  click: () => {
                    updateSlotStatus(slot.id, 'in-progress')
                    activeAlerts.value = activeAlerts.value.filter(a => a.id !== slot.id)
                    toast.add({ title: 'Session Started', color: 'success', icon: 'i-lucide-play' })
                  }
                }
              ]
            })
          }
        }
      }
    })
    
    // Bersihkan activeAlerts jika status sudah bukan booked (misal sudah started via manual button)
    activeAlerts.value = activeAlerts.value.filter(alert => {
      const currentSlot = slots.value.find(s => s.id === alert.id)
      return currentSlot?.status === 'booked'
    })
  }

  const startMonitor = () => {
    if (import.meta.client && !monitorInterval) {
      console.log('Smart Alert Monitor Started...')
      checkUpcomingSlots() // Cek langsung saat pertama kali jalan
      monitorInterval = setInterval(checkUpcomingSlots, 60000) // Cek setiap 1 menit
    }
  }

  return { startMonitor, activeAlerts }
}
