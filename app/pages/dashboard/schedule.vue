<script setup lang="ts">
definePageMeta({ layout: 'dashboard' })

const toast = useToast()
const selectedDate = ref(15)
const currentMonth = 'April 2026'
const selectedSlot = ref<string | null>(null)
const showBookingModal = ref(false)

const weekDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
const calendarDays = Array.from({ length: 30 }, (_, i) => ({
  day: i + 1,
  available: ![1, 2, 3, 6, 10, 13, 17, 20, 24, 27].includes(i + 1)
}))

const { slots: globalSlots, bookSlot } = useSchedules()

const availableSlots = computed(() => {
  return globalSlots.value.map(slot => ({
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
    car: 'Tesla Model 3',
    instructor: 'Pak Ahmad',
    topic: 'Highway Driving - Advanced'
  },
  { 
    id: 2, 
    sessionNumber: 6,
    date: 'Apr 2, 2026', 
    time: '11:00 AM',
    car: 'BYD Atto 3',
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

function confirmBooking() {
  if (selectedSlot.value) {
    bookSlot(selectedSlot.value, 'John Doe')
  }
  showBookingModal.value = false
  toast.add({
    title: 'Session Booked!',
    description: `Your session on Apr ${selectedDate.value}, 2026 at ${selectedSlotDetails.value?.time} has been confirmed.`,
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
                  <div class="p-3 rounded-xl bg-primary/10">
                    <UIcon name="i-lucide-car" class="size-6 text-primary" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <h3 class="font-semibold">Session #{{ session.sessionNumber }}</h3>
                      <UBadge label="Confirmed" color="success" variant="subtle" size="xs" />
                    </div>
                    <p class="text-sm text-muted">{{ session.topic }}</p>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-default">
                <div>
                  <p class="text-xs text-muted">Date & Time</p>
                  <p class="text-sm font-medium">{{ session.date }}</p>
                  <p class="text-sm">{{ session.time }}</p>
                </div>
                <div>
                  <p class="text-xs text-muted">Vehicle</p>
                  <p class="text-sm font-medium">{{ session.car }}</p>
                  <p class="text-sm text-muted">{{ session.instructor }}</p>
                </div>
              </div>

              <template #footer>
                <div class="flex gap-2">
                  <UButton label="Reschedule" variant="outline" color="neutral" size="sm" icon="i-lucide-calendar-days" />
                  <UButton label="Cancel" variant="ghost" color="error" size="sm" icon="i-lucide-x" />
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
                <h3 class="text-sm font-medium">Select Date</h3>
                <div class="flex items-center gap-2">
                  <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" size="xs" />
                  <span class="text-sm font-medium">{{ currentMonth }}</span>
                  <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" size="xs" />
                </div>
              </div>
              
              <!-- Custom Calendar Grid -->
              <div class="border border-default rounded-lg p-4">
                <div class="grid grid-cols-7 gap-1 mb-2">
                  <div v-for="day in weekDays" :key="day" class="text-center text-xs font-medium text-muted py-2">
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
                      'w-full aspect-square rounded-lg text-sm font-medium transition-all',
                      selectedDate === item.day 
                        ? 'bg-primary text-white' 
                        : item.available && item.day >= 7
                          ? 'hover:bg-primary/10 cursor-pointer'
                          : 'text-muted/50 cursor-not-allowed'
                    ]"
                    @click="item.available && item.day >= 7 && (selectedDate = item.day)"
                  >
                    {{ item.day }}
                  </button>
                </div>
                
                <div class="mt-4 flex items-center gap-4 text-xs">
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
                <h3 class="text-sm font-medium">Available Time Slots</h3>
                <UBadge :label="`Apr ${selectedDate}`" variant="subtle" />
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
                        :name="selectedSlot === slot.id ? 'i-lucide-check-circle' : 'i-lucide-clock'" 
                        :class="selectedSlot === slot.id ? 'text-primary' : 'text-muted'"
                        class="size-5"
                      />
                      <div>
                        <span class="font-semibold">{{ slot.time }}</span>
                        <p class="text-sm text-muted">{{ slot.car }} - {{ slot.instructor }}</p>
                      </div>
                    </div>
                    <UBadge 
                      :label="slot.available ? 'Available' : 'Booked'" 
                      :color="slot.available ? 'success' : 'error'"
                      variant="subtle"
                      size="xs"
                    />
                  </div>
                </button>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="flex items-center justify-between">
              <p v-if="selectedSlot" class="text-sm text-muted">
                Selected: Apr {{ selectedDate }}, 2026 at {{ selectedSlotDetails?.time }} with {{ selectedSlotDetails?.instructor }}
              </p>
              <p v-else class="text-sm text-muted">
                Select a date and time slot to continue
              </p>
              <UButton 
                label="Book Session" 
                :disabled="!selectedSlot"
                icon="i-lucide-check"
                @click="showBookingModal = true"
              />
            </div>
          </template>
        </UCard>
      </div>
    </template>

    <!-- Booking Confirmation Modal -->
    <UModal v-model:open="showBookingModal" title="Confirm Booking">
      <template #body>
        <div class="space-y-4">
          <UAlert icon="i-lucide-info" title="Session Details">
            <template #description>
              <ul class="mt-2 space-y-1 text-sm">
                <li><strong>Date:</strong> April {{ selectedDate }}, 2026</li>
                <li><strong>Time:</strong> {{ selectedSlotDetails?.time }}</li>
                <li><strong>Vehicle:</strong> {{ selectedSlotDetails?.car }}</li>
                <li><strong>Instructor:</strong> {{ selectedSlotDetails?.instructor }}</li>
              </ul>
            </template>
          </UAlert>

          <p class="text-sm text-muted">
            By confirming, you agree to attend this session. Cancellations must be made at least 24 hours in advance.
          </p>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton label="Cancel" variant="ghost" color="neutral" @click="showBookingModal = false" />
          <UButton label="Confirm Booking" icon="i-lucide-check" @click="confirmBooking" />
        </div>
      </template>
    </UModal>
  </UDashboardPanel>
</template>
