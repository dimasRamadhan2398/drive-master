<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { computed, ref } from 'vue'

definePageMeta({ layout: 'dashboard' })

const toast = useToast()

// FITUR BARU: Logika Kalender Dinamis
const currentDate = ref(new Date('2026-04-10T00:00:00')) // Default start di April 2026
const selectedDate = ref(15)
const selectedSlot = ref<string | null>(null)
const showBookingModal = ref(false)

const weekDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

const currentMonthStr = computed(() => {
  return currentDate.value.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
})

// PERUBAHAN: Memperbaiki teks bulan statis
const currentMonthShortStr = computed(() => {
  return currentDate.value.toLocaleDateString('en-US', { month: 'short' })
})

const { slots: globalSlots, bookSlot } = useSchedules()

// FITUR BARU: Kalender merender hari secara dinamis berdasarkan bulan yang sedang dipilih
const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth() // 0-11
  
  // Cari tahu tanggal 1 bulan ini jatuh di hari apa (0 = Sun, 1 = Mon)
  const firstDay = new Date(year, month, 1).getDay()
  // Konversi ke format Senin=0, Minggu=6
  const emptyDays = firstDay === 0 ? 6 : firstDay - 1
  
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  
  const days = []
  
  // Selipkan hari kosong sebelum tanggal 1
  for (let i = 0; i < emptyDays; i++) {
    days.push({ day: null, available: false })
  }
  
  // Generate tanggal 1 sampai akhir bulan
  for (let i = 1; i <= daysInMonth; i++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`
    // Periksa apakah ada slot available di tanggal ini dari global state
    const isAvailable = globalSlots.value.some(s => s.date === dateStr && s.status === 'available')
    
    days.push({
      day: i,
      available: isAvailable
    })
  }
  
  return days
})

// FITUR BARU: Fungsi untuk navigasi bulan di kalender
function changeMonth(offset: number) {
  const newDate = new Date(currentDate.value)
  newDate.setMonth(newDate.getMonth() + offset)
  currentDate.value = newDate
}

// FITUR BARU: Hanya tampilkan slot yang sesuai dengan tanggal yang dipilih di kalender
const availableSlots = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = String(currentDate.value.getMonth() + 1).padStart(2, '0')
  const day = String(selectedDate.value).padStart(2, '0')
  const dateStr = `${year}-${month}-${day}`
  
  return globalSlots.value
    .filter(slot => slot.date === dateStr)
    .map(slot => ({
      ...slot,
      available: slot.status === 'available'
    }))
})

// Mock upcoming sessions
const upcomingSessions = ref([
  { 
    id: 1, 
    sessionNumber: 5,
    date: 'Mar 28, 2026', 
    time: '09:30 AM',
    car: 'BYD Atto 1',
    instructor: 'Pak Ahmad',
    topic: 'Highway Driving - Advanced'
  },
  { 
    id: 2, 
    sessionNumber: 6,
    date: 'Apr 2, 2026', 
    time: '11:00 AM',
    car: 'BYD Atto 1',
    instructor: 'Bu Sari',
    topic: 'Night Driving Introduction'
  }
])

const selectedSlotDetails = computed(() => {
  return availableSlots.value.find(s => s.id === selectedSlot.value)
})

function selectSlot(slotId: string) {
  const slot = availableSlots.value.find(s => s.id === slotId)
  if (slot?.available) {
    selectedSlot.value = slotId
  }
}

// PERUBAHAN: Memperbaiki teks bulan statis
function confirmBooking() {
  if (selectedSlot.value) {
    bookSlot(selectedSlot.value, 'John Doe')
  }
  showBookingModal.value = false
  toast.add({
    title: 'Session Booked!',
    description: `Your session on ${currentMonthStr.value.split(' ')[0]} ${selectedDate.value}, ${currentDate.value.getFullYear()} at ${selectedSlotDetails.value?.time} has been confirmed.`,
    icon: 'i-lucide-check-circle',
    color: 'success'
  })
  selectedSlot.value = null
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="My Schedule">
        <template #right>
          <UButton icon="i-lucide-bell" color="neutral" variant="ghost" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Upcoming Sessions -->
        <div>
          <h2 class="text-lg font-semibold mb-4">Upcoming Sessions</h2>
          
          <div v-if="upcomingSessions.length > 0" class="grid md:grid-cols-2 gap-4">
            <UCard v-for="session in upcomingSessions" :key="session.id">
              <div class="flex items-start justify-between">
                <div class="flex items-center gap-4">
                  <div class="p-3 rounded-xl bg-warning/10">
                    <UIcon name="i-lucide-car" class="size-6 text-warning" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <h3 class="font-semibold">Session #{{ session.sessionNumber }}</h3>
                      <UBadge label="Confirmed" color="success" variant="subtle" size="md" />
                    </div>
                    <p class="text-md text-muted">{{ session.topic }}</p>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-default">
                <div>
                  <p class="text-md text-muted">Date & Time</p>
                  <p class="text-md font-medium">{{ session.date }}</p>
                  <p class="text-md">{{ session.time }}</p>
                </div>
                <div>
                  <p class="text-md text-muted">Instructor</p>
                  <p class="text-md font-medium ">{{ session.instructor }}</p>
                  <p class="text-md text-muted">{{ session.car }}</p>
                </div>
              </div>

              <template #footer>
                <div class="flex gap-2">
                  <UButton label="Reschedule" variant="outline" color="warning" size="md" icon="i-lucide-calendar-days" />
                  <UButton label="Cancel" variant="ghost" color="error" size="md" icon="i-lucide-x" />
                </div>
              </template>
            </UCard>
          </div>

          <UEmpty v-else icon="i-lucide-calendar-x" title="No Upcoming Sessions" description="You don't have any scheduled sessions. Book one below!" />
        </div>

        <!-- Book New Session -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">Book a New Session</h2>
          </template>

          <div class="grid lg:grid-cols-2 gap-8">
            <!-- Calendar -->
            <div>
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-md font-medium">Select Date</h3>
                <div class="flex items-center gap-2">
                  <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" size="md" />
                  <span class="text-md font-medium">{{ currentMonth }}</span>
                  <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" size="md" />
                </div>
              </div>
              
              <!-- Custom Calendar Grid -->
              <div class="border border-default rounded-lg p-4">
                <div class="grid grid-cols-7 gap-1 mb-2">
                  <div v-for="day in weekDays" :key="day" class="text-center text-md font-medium text-muted py-2">
                    {{ day }}
                  </div>
                </div>
                <div class="grid grid-cols-7 gap-1">
                  <!-- Empty cells for days before month starts (April 2026 starts on Wednesday) -->
                  <div></div>
                  <div></div>
                  <button 
                    v-for="item in calendarDays" 
                    :key="item.day"
                    :disabled="!item.available || item.day < 7"
                    :class="[
                      'w-full aspect-square rounded-lg text-md font-medium transition-all',
                      selectedDate === item.day 
                        ? 'bg-primary text-white' 
                        : item.available && item.day >= 7
                          ? 'hover:bg-primary/10 cursor-pointer'
                          : 'text-muted/50 cursor-not-allowed'
                    ]"
                    @click="item.available && item.day >= 7 && (selectedDate = item.day)"
                  >
                    <button 
                      v-if="item.day !== null"
                      :disabled="!item.available"
                      :class="[
                        'w-full aspect-square rounded-lg text-sm font-medium transition-all',
                        selectedDate === item.day 
                          ? 'bg-primary text-white' 
                          : item.available 
                            ? 'hover:bg-primary/10 cursor-pointer'
                            : 'text-muted/50 cursor-not-allowed'
                      ]"
                      @click="item.available && (selectedDate = item.day)"
                    >
                      {{ item.day }}
                    </button>
                  </div>
                </div>
                
                <div class="mt-4 flex items-center gap-4 text-md">
                  <div class="flex items-center gap-2">
                    <div class="size-3 rounded bg-primary/10 border border-primary/30"></div>
                    <span class="text-muted">Available</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <div class="size-3 rounded bg-primary"></div>
                    <span class="text-muted">Selected</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Time Slots -->
            <div>
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-md font-medium">Available Time Slots</h3>
                <UBadge :label="`Apr ${selectedDate}`" color="primary" variant="subtle" />
              </div>

              <div class="space-y-3">
                <button
                  v-for="slot in availableSlots"
                  :key="slot.id"
                  :disabled="!slot.available"
                  :class="[
                    'w-full p-4 rounded-lg border text-left transition-all',
                    slot.available 
                      ? selectedSlot === slot.id 
                        ? 'border-primary bg-primary/10' 
                        : 'border-default hover:border-primary cursor-pointer'
                      : 'border-default bg-muted/50 opacity-50 cursor-not-allowed'
                  ]"
                  @click="selectSlot(slot.id)"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3">
                      <UIcon 
                        :name="selectedSlot === slot.id ? 'i-lucide-check-circle' : 'i-lucide-circle'" 
                        :class="selectedSlot === slot.id ? 'text-primary' : 'text-muted'"
                        class="size-5"
                      />
                      <div>
                        <span class="font-semibold">{{ slot.time }}</span>
                        <p class="text-md text-muted">{{ slot.instructor }} - {{ slot.car }}</p>
                      </div>
                    </div>
                    <UBadge 
                      :label="slot.available ? 'Available' : 'Booked'" 
                      :color="slot.available ? 'success' : 'error'"
                      variant="subtle"
                      size="md"
                    />
                  </div>
                </button>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="flex items-center justify-between">
              <p v-if="selectedSlot" class="text-md text-muted">
                Selected: Apr {{ selectedDate }}, 2026 at {{ selectedSlotDetails?.time }} with {{ selectedSlotDetails?.instructor }}
              </p>
              <p v-else class="text-md text-muted">
                Select a date and time slot to continue
              </p>
              <UButton 
                label="Book Session" 
                :disabled="!selectedSlot"
                icon="i-lucide-check"
                color="warning"
                @click="showBookingModal = true"
              />
              <!-- Booking Confirmation Modal -->
              <UModal v-model:open="showBookingModal" title="Confirm Booking">
                <template #body>
                  <div class="space-y-4">
                    <UAlert icon="i-lucide-info" color="warning" title="Session Details">
                      <template #description>
                        <ul class="mt-2 space-y-1 text-md">
                          <li><strong>Date:</strong> April {{ selectedDate }}, 2026</li>
                          <li><strong>Time:</strong> {{ selectedSlotDetails?.time }}</li>
                          <li><strong>Vehicle:</strong> {{ selectedSlotDetails?.car }}</li>
                          <li><strong>Instructor:</strong> {{ selectedSlotDetails?.instructor }}</li>
                        </ul>
                      </template>
                    </UAlert>

                    <p class="text-md text-muted">
                      By confirming, you agree to attend this session. Cancellations must be made at least 24 hours in advance.
                    </p>
                  </div>
                </template>
                <template #footer>
                  <div class="flex justify-end gap-3">
                    <UButton label="Cancel" variant="ghost" color="neutral" @click="showBookingModal = false" />
                    <UButton label="Confirm Booking" color="warning" icon="i-lucide-check" @click="confirmBooking" />
                  </div>
                </template>
              </UModal>
            </div>
          </template>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>