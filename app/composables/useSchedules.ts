export interface TimeSlot {
  id: string
  time: string
  duration: string
  car: string
  instructor: string
  student: string | null
  status: 'available' | 'booked' | 'in-progress' | 'completed' | 'blocked'
}

export const useSchedules = () => {
  const slots = useState<TimeSlot[]>('global-schedules', () => [
    { id: '1', time: '07:00', duration: '60 min', car: 'Tesla Model 3', instructor: 'Pak Ahmad', student: 'John Doe', status: 'completed' },
    { id: '2', time: '08:00', duration: '30 min', car: 'BYD Atto 3', instructor: 'Bu Sari', student: null, status: 'blocked' },
    { id: '3', time: '08:30', duration: '60 min', car: 'BYD Atto 3', instructor: 'Bu Sari', student: 'Budi Santoso', status: 'completed' },
    { id: '4', time: '09:30', duration: '60 min', car: 'Tesla Model 3', instructor: 'Pak Ahmad', student: 'Sarah Putri', status: 'in-progress' },
    { id: '5', time: '11:00', duration: '60 min', car: 'Tesla Model 3', instructor: 'Pak Budi', student: null, status: 'available' },
    { id: '6', time: '13:00', duration: '60 min', car: 'BYD Atto 3', instructor: 'Pak Ahmad', student: 'Amanda Chen', status: 'booked' },
    { id: '7', time: '14:30', duration: '60 min', car: 'Tesla Model 3', instructor: 'Bu Sari', student: null, status: 'blocked' },
    { id: '8', time: '16:00', duration: '60 min', car: 'BYD Atto 3', instructor: 'Pak Budi', student: 'Ricky Wijaya', status: 'booked' }
  ])

  const toggleSlotStatus = (id: string) => {
    const slot = slots.value.find(s => s.id === id)
    if (slot && slot.status === 'available') {
      slot.status = 'blocked'
    } else if (slot && slot.status === 'blocked') {
      slot.status = 'available'
    }
  }

  const deleteSlot = (id: string) => {
    slots.value = slots.value.filter(s => s.id !== id)
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
    deleteSlot,
    bookSlot
  }
}
