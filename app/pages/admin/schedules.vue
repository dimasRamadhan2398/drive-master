<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { ref, computed } from 'vue'
import { useRoute, navigateTo } from 'nuxt/app'

definePageMeta({ layout: 'admin' })

const route = useRoute()
const toast = useToast()
const studentNameToBook = computed(() => route.query.studentName as string | undefined)
const showAddSlotModal = ref(false)
const selectedDate = ref(new Date('2026-04-10T00:00:00'))

// FIX: Define localDateStr to format the selected date
const localDateStr = computed(() => {
  const date = selectedDate.value;
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const day = date.getDate().toString().padStart(2, '0');
  return `${year}-${month}-${day}`;
});

const filterInstructor = ref('All Instructors')
const filterVehicle = ref('All Vehicles')

const { slots: timeSlots, toggleSlotStatus: _toggleSlotStatus, deleteSlot: _deleteSlot, updateSlotStatus, addSlot, editSlot, bookSlot } = useSchedules()

const instructors = ['Mr. Ahmad', 'Ms. Sari', 'Mr. Budi']
const vehicles = ['BYD Atto 1']

// FITUR BARU: Sinkronisasi jam operasional dari settings
const { operatingHours } = useSettings()
const currentDayOperatingHours = computed(() => {
  const day = selectedDate.value.getDay() // 0 = Sunday, 6 = Saturday
  const isWeekend = day === 0 || day === 6
  
  if (isWeekend) {
    return {
      start: operatingHours.value.weekendStart,
      end: operatingHours.value.weekendEnd,
      nightStart: operatingHours.value.nightStart,
      nightEnd: operatingHours.value.nightEnd,
      nightEnabled: operatingHours.value.nightEnabled,
      isClosed: day === 0 && operatingHours.value.sundayClosed
    }
  }
  
  return {
    start: operatingHours.value.mondayStart,
    end: operatingHours.value.mondayEnd,
    nightStart: operatingHours.value.nightStart,
    nightEnd: operatingHours.value.nightEnd,
    nightEnabled: operatingHours.value.nightEnabled,
    isClosed: false
  }
})

const filteredSlots = computed(() => {
  const dateStr = localDateStr.value // FIX: Menggunakan localDateStr
  return timeSlots.value.filter(slot => {
    const matchDate = slot.date === dateStr
    const matchInst = filterInstructor.value === 'All Instructors' || slot.instructor === filterInstructor.value
    const matchVeh = filterVehicle.value === 'All Vehicles' || slot.car === filterVehicle.value
    return matchDate && matchInst && matchVeh
  })
})

function changeDay(offset: number) {
  const nextDate = new Date(selectedDate.value)
  nextDate.setDate(selectedDate.value.getDate() + offset)
  selectedDate.value = nextDate
}

function toggleSlotStatus(slotId: string) {
  _toggleSlotStatus(slotId)
  const slot = timeSlots.value.find(s => s.id === slotId)
  if (slot?.status === 'blocked') {
    toast.add({ title: 'Slot Blocked', color: 'warning', icon: 'i-lucide-lock' })
  } else {
    toast.add({ title: `Slot is now ${slot?.status}`, color: 'success', icon: 'i-lucide-check' })
  }
}

function deleteSlot(slotId: string) {
  _deleteSlot(slotId)
  toast.add({ title: 'Slot Deleted', color: 'error', icon: 'i-lucide-trash' })
}

// FITUR BARU: Logika In-Progress Session
const isToday = computed(() => {
  const today = new Date()
  const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
  return localDateStr.value === todayStr
})

function startSession(slotId: string) {
  if (confirm('Apakah kursus sudah mau jalan? Sesi akan berubah menjadi "In Progress".')) {
    updateSlotStatus(slotId, 'in-progress')
    toast.add({ 
      title: 'Session Started', 
      description: 'Driving session is now in progress.', 
      color: 'warning', 
      icon: 'i-lucide-play' 
    })
  }
}

function completeSession(slotId: string) {
  if (confirm('Apakah kursus sudah selesai? Sesi akan ditandai sebagai "Completed".')) {
    updateSlotStatus(slotId, 'completed')
    toast.add({ 
      title: 'Session Completed', 
      description: 'Driving session has been finished.', 
      color: 'neutral', 
      icon: 'i-lucide-check-circle' 
    })
  }
}

function handleManualBooking(slotId: string) {
  const studentName = studentNameToBook.value ? decodeURIComponent(studentNameToBook.value) : null
  if (studentName) {
    bookSlot(slotId, studentName)
    toast.add({ title: 'Slot Booked', description: `Berhasil booking untuk ${studentName}`, color: 'success' })
    // Kembali ke halaman schedule tanpa query untuk membersihkan state
    navigateTo('/admin/schedules')
  } else {
    const name = prompt('Masukkan nama siswa untuk booking manual:', 'Siswa Manual')
    if (name) {
      bookSlot(slotId, name)
      toast.add({ title: 'Slot Booked', description: `Berhasil booking untuk ${name}`, color: 'success' })
    }
  }
}

function handleSlotClick(slot: any) {
  if (studentNameToBook.value && slot.status === 'available') {
    const studentName = decodeURIComponent(studentNameToBook.value as string)
    bookSlot(slot.id, studentName)
    toast.add({ title: 'Slot Booked', description: `Berhasil booking untuk ${studentName}`, color: 'success', icon: 'i-lucide-calendar-check' })
    // Kembali ke halaman schedule tanpa query untuk membersihkan state
    navigateTo('/admin/schedules')
  }
}

function cancelBooking(slotId: string) {
  if (confirm('Batalkan booking dan kembalikan slot menjadi Available?')) {
    updateSlotStatus(slotId, 'available')
    toast.add({ title: 'Booking Cancelled', description: 'Slot is now available again.', color: 'neutral' })
  }
}

// FITUR BARU: Add Slot State & Logic
const newSlotData = ref({
  time: '08:00',
  duration: '60',
  car: '',
  instructor: ''
})

function saveNewSlot() {
  if (!newSlotData.value.time || !newSlotData.value.car || !newSlotData.value.instructor) {
    toast.add({ title: 'Error', description: 'Please fill all fields', color: 'error' })
    return
  }
  
  // FITUR BARU: Validasi jam operasional (Termasuk Night Shift)
  const { start, end, nightStart, nightEnd, nightEnabled, isClosed } = currentDayOperatingHours.value
  if (isClosed) {
    toast.add({ title: 'Error', description: 'Business is closed on this day', color: 'error' })
    return
  }

  const isDayShift = newSlotData.value.time >= start && newSlotData.value.time <= end
  const isNightShift = nightEnabled && newSlotData.value.time >= nightStart && newSlotData.value.time <= nightEnd

  if (!isDayShift && !isNightShift) {
    let msg = `Time must be between ${start} - ${end}`
    if (nightEnabled) msg += ` or ${nightStart} - ${nightEnd}`
    toast.add({ title: 'Error', description: msg, color: 'error' })
    return
  }
  
  const id = Date.now().toString()
  addSlot({
    id,
    date: localDateStr.value, // FIX: Menggunakan localDateStr
    time: newSlotData.value.time,
    duration: newSlotData.value.duration + ' min',
    car: newSlotData.value.car,
    instructor: newSlotData.value.instructor,
    student: null,
    status: 'available'
  })
  
  toast.add({ title: 'New Slot Added', color: 'success', icon: 'i-lucide-check' })
  showAddSlotModal.value = false
  // Reset
  newSlotData.value = { time: '08:00', duration: '60', car: '', instructor: '' }
}

// FITUR BARU: Edit Slot State & Logic
const showEditSlotModal = ref(false)
const editSlotData = ref({
  id: '',
  time: '',
  duration: '60',
  car: '',
  instructor: ''
})

function openEditModal(slot: any) {
  if (slot.status !== 'available') {
    toast.add({ title: 'Error', description: 'Hanya slot dengan status Available yang bisa diedit', color: 'error' })
    return
  }
  editSlotData.value = {
    id: slot.id,
    time: slot.time,
    duration: slot.duration.replace(' min', ''),
    car: slot.car,
    instructor: slot.instructor
  }
  showEditSlotModal.value = true
}

function saveEditSlot() {
  if (!editSlotData.value.time || !editSlotData.value.car || !editSlotData.value.instructor) {
    toast.add({ title: 'Error', description: 'Please fill all fields', color: 'error' })
    return
  }

  // Validasi jam operasional (Termasuk Night Shift)
  const { start, end, nightStart, nightEnd, nightEnabled, isClosed } = currentDayOperatingHours.value
  if (isClosed) {
    toast.add({ title: 'Error', description: 'Business is closed on this day', color: 'error' })
    return
  }

  const isDayShift = editSlotData.value.time >= start && editSlotData.value.time <= end
  const isNightShift = nightEnabled && editSlotData.value.time >= nightStart && editSlotData.value.time <= nightEnd

  if (!isDayShift && !isNightShift) {
    let msg = `Time must be between ${start} - ${end}`
    if (nightEnabled) msg += ` or ${nightStart} - ${nightEnd}`
    toast.add({ title: 'Error', description: msg, color: 'error' })
    return
  }

  editSlot(editSlotData.value.id, {
    time: editSlotData.value.time,
    duration: editSlotData.value.duration + ' min',
    car: editSlotData.value.car,
    instructor: editSlotData.value.instructor
  })
  
  toast.add({ title: 'Slot Updated', color: 'success', icon: 'i-lucide-check' })
  showEditSlotModal.value = false
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UAlert
        v-if="studentNameToBook"
        icon="i-lucide-user-check"
        color="info"
        variant="subtle"
        class="m-4"
        title="Mode Booking"
      >
        <template #description>
          <span>Pilih slot untuk <strong>{{ decodeURIComponent(studentNameToBook) }}</strong>. <UButton variant="link" :padded="false" @click="navigateTo('/admin/schedules')">Batalkan</UButton></span>
        </template>
      </UAlert>
      <UDashboardNavbar title="Schedule Management">
        <template #right>
          <UButton icon="i-lucide-plus" color="warning" label="Add Time Slot" @click="showAddSlotModal = true" />
          <UColorModeButton />
          <!-- Add Slot Modal -->
          <UModal v-model:open="showAddSlotModal" title="Add New Time Slot">
            <template #body>
              <div class="space-y-5">
                <UFormField label="Date" required>
                  <UInput type="date" :model-value="localDateStr" disabled class="w-full"/>
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                  <UFormField :label="currentDayOperatingHours.isClosed ? 'Closed Today' : `Start Time`" required>
                    <template #hint>
                      <div class="text-[10px] text-muted-foreground flex flex-col items-end">
                        <span>Day: {{ currentDayOperatingHours.start }}-{{ currentDayOperatingHours.end }}</span>
                        <span v-if="currentDayOperatingHours.nightEnabled">Night: {{ currentDayOperatingHours.nightStart }}-{{ currentDayOperatingHours.nightEnd }}</span>
                      </div>
                    </template>
                    <UInput 
                      type="time" 
                      v-model="newSlotData.time" 
                      :disabled="currentDayOperatingHours.isClosed"
                      class="w-full"
                    />
                  </UFormField>
                  <UFormField label="Duration" >
                    <USelect 
                      :items="[   
                        { label: '60 minutes', value: '60' },
                      ]"
                      v-model="newSlotData.duration"
                      disabled
                      class="w-full"
                    />
                  </UFormField>
                </div>
                <UFormField label="Vehicle" required>
                  <USelect :items="vehicles" v-model="newSlotData.car" placeholder="Select vehicle" class="w-full" />
                </UFormField>
                <UFormField label="Instructor" required>
                  <USelect :items="instructors" v-model="newSlotData.instructor" placeholder="Select instructor" class="w-full"/>
                </UFormField>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton label="Create Slot" icon="i-lucide-plus" @click="saveNewSlot" color="warning"/>
              </div>
            </template>
          </UModal>
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <div class="flex items-center gap-2">
            <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" @click="changeDay(-1)" />
            <div class="relative flex items-center justify-center min-w-[160px] hover:bg-gray-100 rounded py-1 transition-colors cursor-pointer">
              <span class="font-medium text-center pointer-events-none">{{ selectedDate.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' }) }}</span>
              <input 
                type="date" 
                class="absolute inset-0 opacity-0 cursor-pointer w-full h-full"
                :value="localDateStr"
                @input="e => { if ((e.target as HTMLInputElement).value) selectedDate = new Date((e.target as HTMLInputElement).value + 'T00:00:00') }"
              />
            </div>
            <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" @click="changeDay(1)" />
          </div>
        </template>
        <template #right>
          <UButton icon="i-lucide-calendar" label="Today" color="neutral" variant="outline" @click="selectedDate = new Date()" />
          <USelect 
            v-model="filterInstructor"
            :items="['All Instructors', ...instructors]"
            placeholder="Filter by instructor"
            class="w-48"
            color="warning"
          />
          <USelect 
            v-model="filterVehicle"
            :items="['All Vehicles', ...vehicles]"
            placeholder="Filter by vehicle"
            class="w-44"
            color="warning"
          />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Stats -->
        <div class="grid sm:grid-cols-2 lg:grid-cols-5 gap-4">
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-green-500/10">
                <UIcon name="i-lucide-check-circle" class="size-5 text-green-500" />
              </div>
              <div>
                <p class="text-xl font-bold">{{ timeSlots.filter(s => s.status === 'available').length }}</p>
                <p class="text-md text-muted">Available</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-info/10">
                <UIcon name="i-lucide-calendar-check" class="size-5 text-info" />
              </div>
              <div>
                <p class="text-xl font-bold">{{ timeSlots.filter(s => s.status === 'booked').length }}</p>
                <p class="text-md text-muted">Booked</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-amber-500/10">
                <UIcon name="i-lucide-play-circle" class="size-5 text-amber-500" />
              </div>
              <div>
                <p class="text-xl font-bold">{{ timeSlots.filter(s => s.status === 'in-progress').length }}</p>
                <p class="text-md text-muted">In Progress</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-neutral-500/10">
                <UIcon name="i-lucide-check-square" class="size-5 text-neutral-500" />
              </div>
              <div>
                <p class="text-xl font-bold">{{ timeSlots.filter(s => s.status === 'completed').length }}</p>
                <p class="text-sm text-muted">Completed</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-red-500/10">
                <UIcon name="i-lucide-lock" class="size-5 text-red-500" />
              </div>
              <div>
                <p class="text-xl font-bold">{{ timeSlots.filter(s => s.status === 'blocked').length }}</p>
                <p class="text-md text-muted">Blocked</p>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Schedule Grid -->
        <div class="grid lg:grid-cols-3 gap-6">
          <!-- Calendar -->
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">Select Date</h2>
                <div class="flex items-center gap-1">
                  <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" size="xs" @click="changeDay(-1)" />
                  <span class="text-sm font-medium">{{ selectedDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }) }}</span>
                  <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" size="xs" @click="changeDay(1)" />
                </div>
              </div>
            </template>
            <!-- Simple custom calendar grid -->
            <div class="grid grid-cols-7 gap-1 mb-2">
              <div v-for="day in ['Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su']" :key="day" class="text-center text-xs font-bold text-muted py-2">{{ day }}</div>
            </div>
            <div class="grid grid-cols-7 gap-1">
              <div></div><div></div>
              <button v-for="d in 30" :key="d" :class="['w-full aspect-square rounded-full text-sm font-medium transition-all flex items-center justify-center', selectedDate.getDate() === d ? 'bg-primary text-white shadow-md' : 'hover:bg-muted/50 cursor-pointer']" @click="selectedDate = new Date(2026, selectedDate.getMonth(), d)">{{ d }}</button>

            </div>
          </UCard>

          <!-- Time Slots List -->
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">Time Slots</h2>
                <div class="flex gap-2">
                  <UBadge label="Available" color="success" variant="subtle" />
                  <UBadge label="Booked" color="info" variant="subtle" />
                  <UBadge label="Completed" color="neutral" variant="subtle" />
                  <UBadge label="Blocked" color="error" variant="subtle" />
                </div>
              </div>
            </template>

            <div class="space-y-4">
              <div
                v-for="slot in filteredSlots"
                :key="slot.id"
                class="p-4 rounded-lg border border-default hover:shadow-md transition-shadow"
                :class="{
                  'border-l-4 border-l-primary': slot.status === 'available',
                  'border-l-4 border-l-info bg-info/5': slot.status === 'booked',
                  'border-l-4 border-l-amber-500': slot.status === 'in-progress',
                  'border-l-4 border-l-neutral-500 bg-neutral-500/5': slot.status === 'completed',
                  'border-l-4 border-l-red-500 opacity-60': slot.status === 'blocked',
                  'cursor-pointer hover:bg-primary/5': studentNameToBook && slot.status === 'available'
                }"
                @click="handleSlotClick(slot)"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-5">
                    <div class="text-center min-w-[70px]">
                      <p class="text-xl font-black flex items-center justify-center gap-1">
                        {{ slot.time }}
                        <UIcon v-if="slot.time >= operatingHours.nightStart" name="i-lucide-moon" class="size-3 text-indigo-500" />
                      </p>
                      <p class="text-xs font-semibold text-muted">{{ slot.duration }}</p>
                    </div>
                    
                    <USeparator orientation="vertical" class="h-12 hidden sm:block" />
                    
                    <div class="space-y-1">
                      <div class="flex items-center gap-2 font-medium">
                        <UIcon name="i-lucide-car" class="size-4 text-muted" />
                        <span class="text-md">{{ slot.car }}</span>
                      </div>
                      <div class="flex items-center gap-2 mt-1">
                        <UIcon name="i-lucide-contact" class="size-4 text-muted" />
                        <span class="text-sm">{{ slot.instructor }}</span>
                        <UBadge v-if="slot.time >= operatingHours.nightStart" label="Night Session" variant="subtle" size="xs" color="neutral" class="ml-2" />
                      </div>
                    </div>

                    <div v-if="slot.student" class="hidden sm:block ml-4 pl-4 border-l border-default">
                      <p class="text-xs text-muted uppercase font-bold tracking-wider mb-1">Booked by</p>
                      <p class="font-bold flex items-center gap-2">
                        <UIcon name="i-lucide-user" class="size-4 text-info" />
                        {{ slot.student }}
                      </p>
                    </div>

                    <!-- FITUR BARU: Tombol Start/Complete Session (Hanya muncul di hari H) -->
                    <div v-if="slot.status === 'booked' && isToday" class="flex items-center">
                      <UButton 
                        label="Start Session" 
                        icon="i-lucide-play" 
                        size="sm" 
                        color="warning" 
                        variant="soft"
                        @click="startSession(slot.id)"
                      />
                    </div>
                    <div v-if="slot.status === 'in-progress' && isToday" class="flex items-center">
                      <UButton 
                        label="Complete" 
                        icon="i-lucide-check-circle" 
                        size="sm" 
                        color="primary" 
                        variant="soft"
                        @click="completeSession(slot.id)"
                      />
                    </div>
                  </div>

                  <div class="flex items-center gap-4">
                    <UBadge 
                      class="hidden sm:flex"
                      :label="slot.status === 'available' ? 'Available' : slot.status === 'booked' ? 'Booked' : slot.status === 'in-progress' ? 'In Progress' : slot.status === 'completed' ? 'Completed' : 'Blocked'"
                      :color="slot.status === 'available' ? 'primary' : slot.status === 'booked' ? 'info' : slot.status === 'in-progress' ? 'warning' : slot.status === 'completed' ? 'neutral' : 'error'"
                      variant="subtle"
                    />
                    <UDropdownMenu
                      :items="[
                        [
                          { label: slot.status === 'blocked' ? 'Unblock Slot' : 'Block Slot', icon: slot.status === 'blocked' ? 'i-lucide-unlock' : 'i-lucide-lock', onSelect: () => toggleSlotStatus(slot.id), disabled: slot.status === 'in-progress' || slot.status === 'completed' },
                          { label: 'Book Manually', icon: 'i-lucide-user-plus', onSelect: () => handleManualBooking(slot.id), disabled: slot.status !== 'available' },
                          { label: 'Edit Slot', icon: 'i-lucide-pencil', onSelect: () => openEditModal(slot), disabled: slot.status !== 'available' }
                        ],
                        [
                          { label: 'Start Session', icon: 'i-lucide-play', color: 'warning', onSelect: () => startSession(slot.id), disabled: slot.status !== 'booked' },
                          { label: 'Complete Session', icon: 'i-lucide-check-circle', color: 'primary', onSelect: () => completeSession(slot.id), disabled: slot.status !== 'in-progress' },
                          { label: 'Cancel Booking', icon: 'i-lucide-user-minus', color: 'neutral', onSelect: () => cancelBooking(slot.id), disabled: slot.status !== 'booked' }
                        ],
                        [
                          { label: 'Delete Slot', icon: 'i-lucide-trash', color: 'error', onSelect: () => deleteSlot(slot.id) }
                        ]
                      ]"
                    >
                      <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" />
                      
                    </UDropdownMenu>
                    
                  </div>
                </div>
              </div>

              <div v-if="filteredSlots.length === 0" class="text-center py-10 border-2 border-dashed border-default rounded-xl">
                <UIcon name="i-lucide-calendar-x" class="size-10 text-muted mx-auto mb-3" />
                <p class="text-muted font-medium">No time slots found for these filters.</p>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- FITUR BARU: Edit Slot Modal -->
    <UModal v-model:open="showEditSlotModal" title="Edit Time Slot">
      <template #body>
        <div class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <UFormField :label="currentDayOperatingHours.isClosed ? 'Closed Today' : `Start Time`" required>
              <template #hint>
                <div class="text-[10px] text-muted-foreground flex flex-col items-end">
                  <span>Day: {{ currentDayOperatingHours.start }}-{{ currentDayOperatingHours.end }}</span>
                  <span v-if="currentDayOperatingHours.nightEnabled">Night: {{ currentDayOperatingHours.nightStart }}-{{ currentDayOperatingHours.nightEnd }}</span>
                </div>
              </template>
              <UInput type="time" v-model="editSlotData.time" :disabled="currentDayOperatingHours.isClosed" class="w-full"/>
            </UFormField>
            <UFormField label="Duration">
              <USelect 
                :items="[                
                  { label: '60 minutes', value: '60' },
                ]"
                v-model="editSlotData.duration"
                disabled
                class="w-full"
              />
            </UFormField>
          </div>
          <UFormField label="Vehicle" required>
            <USelect :items="vehicles" v-model="editSlotData.car" placeholder="Select vehicle" class="w-full"/>
          </UFormField>
          <UFormField label="Instructor" required>
            <USelect :items="instructors" v-model="editSlotData.instructor" placeholder="Select instructor" class="w-full"/>
          </UFormField>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton label="Save Changes" icon="i-lucide-save" @click="saveEditSlot" color="warning"/>
        </div>
      </template>
    </UModal>

    </template>

  </UDashboardPanel>
</template>
