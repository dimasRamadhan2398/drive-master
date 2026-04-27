<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { ref } from 'vue'

definePageMeta({ layout: 'admin' })

const toast = useToast()
const showAddSlotModal = ref(false)
const selectedDate = ref(new Date('2026-04-10T00:00:00'))

<<<<<<< HEAD
// Mock time slots
const timeSlots = ref([
  { id: '1', time: '08:00', duration: '60 min', car: 'BYD Atto 1', instructor: 'Pak Ahmad', student: 'John Doe', status: 'booked' },
  { id: '2', time: '09:30', duration: '60 min', car: 'BYD Atto 1', instructor: 'Bu Sari', student: 'Sarah Putri', status: 'in-progress' },
  { id: '3', time: '11:00', duration: '60 min', car: 'BYD Atto 1', instructor: 'Pak Budi', student: null, status: 'available' },
  { id: '4', time: '13:00', duration: '60 min', car: 'BYD Atto 1', instructor: 'Pak Ahmad', student: 'Amanda Chen', status: 'booked' },
  { id: '5', time: '14:30', duration: '60 min', car: 'BYD Atto 1', instructor: 'Bu Sari', student: null, status: 'available' },
  { id: '6', time: '16:00', duration: '60 min', car: 'BYD Atto 1', instructor: 'Pak Budi', student: null, status: 'blocked' }
])
=======
const filterInstructor = ref('All Instructors')
const filterVehicle = ref('All Vehicles')

const { slots: timeSlots, toggleSlotStatus: _toggleSlotStatus, deleteSlot: _deleteSlot } = useSchedules()
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815

const instructors = ['Pak Ahmad', 'Bu Sari', 'Pak Budi']
const vehicles = ['BYD Atto 1']

const filteredSlots = computed(() => {
  return timeSlots.value.filter(slot => {
    const matchInst = filterInstructor.value === 'All Instructors' || slot.instructor === filterInstructor.value
    const matchVeh = filterVehicle.value === 'All Vehicles' || slot.car === filterVehicle.value
    return matchInst && matchVeh
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
  } else if (slot?.status === 'available') {
    toast.add({ title: 'Slot Unblocked', color: 'success', icon: 'i-lucide-unlock' })
  }
}

function deleteSlot(slotId: string) {
  _deleteSlot(slotId)
  toast.add({ title: 'Slot Deleted', color: 'error', icon: 'i-lucide-trash' })
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Schedule Management">
        <template #right>
          <UButton icon="i-lucide-plus" color="warning" label="Add Time Slot" @click="showAddSlotModal = true" />
          <!-- Add Slot Modal -->
          <UModal v-model:open="showAddSlotModal" title="Add Time Slot">
            <template #body>
              <div class="space-y-4">
                <UFormField label="Date" required>
                  <UInput type="date" color="warning" class="w-full" :model-value="selectedDate.toISOString().split('T')[0]" />
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                  <UFormField label="Start Time" required>
                    <UInput type="time" color="warning" class="w-full" placeholder="08:00" />
                  </UFormField>
                  <UFormField label="Duration" required>
                    <UInput disabled placeholder="60 minutes" 
                      :items="[
                        { label: '60 minutes', value: '60' }
                      ]"
                      class="w-full"
                      color="warning"
                    />
                  </UFormField>
                </div>
                <UFormField label="Vehicle" required>
                  <USelect :items="vehicles" placeholder="Select vehicle" color="warning" class="w-full" />
                </UFormField>
                <UFormField label="Instructor" required>
                  <USelect :items="instructors" placeholder="Select instructor" color="warning" class="w-full" />
                </UFormField>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton label="Cancel" variant="ghost" color="neutral" @click="showAddSlotModal = false" />
                <UButton label="Add Slot" color="warning" icon="i-lucide-plus" @click="showAddSlotModal = false" />
              </div>
            </template>
          </UModal>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <div class="flex items-center gap-2">
            <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" @click="changeDay(-1)" />
            <span class="font-medium min-w-[160px] text-center">{{ selectedDate.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' }) }}</span>
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
<<<<<<< HEAD
        <div class="grid md:grid-cols-4 gap-4">
=======
        <div class="grid sm:grid-cols-2 lg:grid-cols-5 gap-4">
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
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
              <div class="p-2 rounded-lg bg-blue-500/10">
                <UIcon name="i-lucide-check-square" class="size-5 text-blue-500" />
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
<<<<<<< HEAD
                <div class="flex items-center gap-2">
                  <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" size="md" />
                  <span class="text-md font-medium">April 2026</span>
                  <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" size="md" />
=======
                <div class="flex items-center gap-1">
                  <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" size="xs" @click="changeDay(-1)" />
                  <span class="text-sm font-medium">{{ selectedDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }) }}</span>
                  <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" size="xs" @click="changeDay(1)" />
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
                </div>
              </div>
            </template>
            <!-- Simple custom calendar grid -->
            <div class="grid grid-cols-7 gap-1 mb-2">
<<<<<<< HEAD
              <div v-for="day in ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']" :key="day" class="text-center text-md font-medium text-muted py-2">{{ day }}</div>
            </div>
            <div class="grid grid-cols-7 gap-1">
              <div></div><div></div>
              <button v-for="d in 30" :key="d" :class="['w-full aspect-square rounded-lg text-md font-medium transition-all', selectedDate.getDate() === d ? 'bg-primary text-white' : 'hover:bg-primary/10 cursor-pointer']" @click="selectedDate = new Date(2026, 3, d)">{{ d }}</button>
=======
              <div v-for="day in ['Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su']" :key="day" class="text-center text-xs font-bold text-muted py-2">{{ day }}</div>
            </div>
            <div class="grid grid-cols-7 gap-1">
              <div></div><div></div>
              <button v-for="d in 30" :key="d" :class="['w-full aspect-square rounded-full text-sm font-medium transition-all flex items-center justify-center', selectedDate.getDate() === d ? 'bg-primary text-white shadow-md' : 'hover:bg-muted/50 cursor-pointer']" @click="selectedDate = new Date(2026, selectedDate.getMonth(), d)">{{ d }}</button>
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
            </div>
          </UCard>

          <!-- Time Slots List -->
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">Time Slots</h2>
                <div class="flex gap-2">
<<<<<<< HEAD
                  <UBadge label="Available" color="primary" variant="subtle" />
                  <UBadge label="Booked" color="info" variant="subtle" />
                  <UBadge label="Blocked" color="error" variant="subtle" />
=======
                  <UBadge label="Available" color="success" variant="subtle" />
                  <UBadge label="Booked" color="primary" variant="subtle" />
                  <UBadge label="Completed" color="blue" variant="subtle" />
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
                </div>
              </div>
            </template>

            <div class="space-y-4">
              <div
                v-for="slot in filteredSlots"
                :key="slot.id"
                class="p-4 rounded-lg border border-default hover:shadow-md transition-shadow"
                :class="{
<<<<<<< HEAD
                  'border-l-4 border-l-primary': slot.status === 'available',
                  'border-l-4 border-l-info': slot.status === 'booked',
=======
                  'border-l-4 border-l-green-500': slot.status === 'available',
                  'border-l-4 border-l-primary bg-primary/5': slot.status === 'booked',
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
                  'border-l-4 border-l-amber-500': slot.status === 'in-progress',
                  'border-l-4 border-l-blue-500 bg-blue-500/5': slot.status === 'completed',
                  'border-l-4 border-l-red-500 opacity-60': slot.status === 'blocked'
                }"
              >
                <div class="flex items-center justify-between">
<<<<<<< HEAD
                  <div class="flex items-center gap-4">
                    <div class="text-center min-w-[60px]">
                      <p class="text-lg font-bold">{{ slot.time }}</p>
                      <p class="text-md text-muted">{{ slot.duration }}</p>
=======
                  <div class="flex items-center gap-5">
                    <div class="text-center min-w-[70px]">
                      <p class="text-xl font-black">{{ slot.time }}</p>
                      <p class="text-xs font-semibold text-muted">{{ slot.duration }}</p>
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
                    </div>
                    
                    <USeparator orientation="vertical" class="h-12 hidden sm:block" />
                    
                    <div class="space-y-1">
                      <div class="flex items-center gap-2 font-medium">
                        <UIcon name="i-lucide-car" class="size-4 text-muted" />
                        <span class="text-md">{{ slot.car }}</span>
                      </div>
                      <div class="flex items-center gap-2 mt-1">
<<<<<<< HEAD
                        <UIcon name="i-lucide-user" class="size-4 text-muted" />
                        <span class="text-md">{{ slot.instructor }}</span>
                      </div>
                    </div>
                    <div v-if="slot.student" class="ml-4 pl-4 border-l border-default">
                      <p class="text-md text-muted">Booked by</p>
                      <p class="font-medium">{{ slot.student }}</p>
=======
                        <UIcon name="i-lucide-contact" class="size-4 text-muted" />
                        <span class="text-sm">{{ slot.instructor }}</span>
                      </div>
                    </div>

                    <div v-if="slot.student" class="hidden sm:block ml-4 pl-4 border-l border-default">
                      <p class="text-xs text-muted uppercase font-bold tracking-wider mb-1">Booked by</p>
                      <p class="font-bold flex items-center gap-2">
                        <UIcon name="i-lucide-user" class="size-4 text-primary" />
                        {{ slot.student }}
                      </p>
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
                    </div>
                  </div>

                  <div class="flex items-center gap-4">
                    <UBadge 
<<<<<<< HEAD
                      :label="slot.status === 'available' ? 'Available' : slot.status === 'booked' ? 'Booked' : slot.status === 'in-progress' ? 'In Progress' : 'Blocked'"
                      :color="slot.status === 'available' ? 'primary' : slot.status === 'booked' ? 'info' : slot.status === 'in-progress' ? 'warning' : 'error'"
=======
                      class="hidden sm:flex"
                      :label="slot.status === 'available' ? 'Available' : slot.status === 'booked' ? 'Booked' : slot.status === 'in-progress' ? 'In Progress' : slot.status === 'completed' ? 'Completed' : 'Blocked'"
                      :color="slot.status === 'available' ? 'success' : slot.status === 'booked' ? 'primary' : slot.status === 'in-progress' ? 'warning' : slot.status === 'completed' ? 'blue' : 'error'"
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
                      variant="subtle"
                    />
                    <UDropdownMenu
                      :items="[
                        [
                          { label: slot.status === 'blocked' ? 'Unblock Slot' : 'Block Slot', icon: slot.status === 'blocked' ? 'i-lucide-unlock' : 'i-lucide-lock', onSelect: () => toggleSlotStatus(slot.id) },
                          { label: 'Edit Slot', icon: 'i-lucide-pencil' }
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
    </template>

<<<<<<< HEAD
    
    
=======
    <!-- Add Slot Modal -->
    <UModal v-model:open="showAddSlotModal" title="Add New Time Slot">
      <template #body>
        <div class="space-y-5">
          <UFormField label="Date" required>
            <UInput type="date" :model-value="selectedDate.toISOString().split('T')[0]" />
          </UFormField>
          <div class="grid grid-cols-2 gap-4">
            <UFormField label="Start Time" required>
              <UInput type="time" placeholder="08:00" />
            </UFormField>
            <UFormField label="Duration" required>
              <USelect 
                :items="[
                  { label: '45 minutes', value: '45' },
                  { label: '60 minutes', value: '60' },
                  { label: '90 minutes', value: '90' }
                ]"
                model-value="60"
              />
            </UFormField>
          </div>
          <UFormField label="Vehicle" required>
            <USelect :items="vehicles" placeholder="Select vehicle" />
          </UFormField>
          <UFormField label="Instructor" required>
            <USelect :items="instructors" placeholder="Select instructor" />
          </UFormField>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton label="Cancel" variant="ghost" color="neutral" @click="showAddSlotModal = false" />
          <UButton label="Create Slot" icon="i-lucide-plus" @click="showAddSlotModal = false" />
        </div>
      </template>
    </UModal>
>>>>>>> 9e0209d0057376fb1faa4e5d070cc514f07a8815
  </UDashboardPanel>
</template>
