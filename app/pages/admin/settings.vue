<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { reactive, ref } from 'vue'

definePageMeta({ layout: 'admin' })

const toast = useToast()

// Mock settings
const generalSettings = reactive({
  businessName: 'Drive Master Academy',
  email: 'info@drivemasteracademy.id',
  phone: '+62 812-3456-7890',
  address: 'Jl. Alam Sutera Boulevard No. 123, Tangerang 15143',
  whatsappNumber: '6281234567890'
})

const operatingHours = reactive({
  mondayStart: '08:00',
  mondayEnd: '18:00',
  weekendStart: '08:00',
  weekendEnd: '16:00',
  nightStart: '18:00',
  nightEnd: '20:00',
  sundayClosed: true
})

const notificationSettings = reactive({
  emailNotifications: true,
  whatsappNotifications: true,
  reminderHours: 24,
  adminAlerts: true,
  newUserRegistration: true,
  newPackagePurchase: true
})

const vehicles = ref([
  { id: 1, name: 'BYD Atto 1', plate: 'B 1234 EV', status: 'active', photoUrl: '' },
])

// Vehicle CRUD State
const isVehicleModalOpen = ref(false)
const isEditingVehicle = ref(false)
const vehicleForm = ref({ id: 0, name: '', plate: '', status: 'active', photoUrl: '' })

function openNewVehicle() {
  isEditingVehicle.value = false
  vehicleForm.value = { id: 0, name: '', plate: '', status: 'active', photoUrl: '' }
  isVehicleModalOpen.value = true
}

function openEditVehicle(vehicle: any) {
  isEditingVehicle.value = true
  vehicleForm.value = { ...vehicle }
  isVehicleModalOpen.value = true
}

function saveVehicle() {
  if (isEditingVehicle.value) {
    const idx = vehicles.value.findIndex(v => v.id === vehicleForm.value.id)
    if (idx !== -1) vehicles.value[idx] = { ...vehicleForm.value }
    toast.add({ title: 'Vehicle Updated', description: 'Vehicle details saved.', color: 'success' })
  } else {
    vehicles.value.push({
      ...vehicleForm.value,
      id: Math.max(0, ...vehicles.value.map(v => v.id)) + 1
    })
    toast.add({ title: 'Vehicle Added', description: 'New vehicle created.', color: 'success' })
  }
  isVehicleModalOpen.value = false
}

function deleteVehicle() {
  vehicles.value = vehicles.value.filter(v => v.id !== vehicleForm.value.id)
  toast.add({ title: 'Vehicle Deleted', description: 'Vehicle has been removed.', color: 'error' })
  isVehicleModalOpen.value = false
}

const instructors = ref([
  { id: 1, name: 'Mr. Ahmad', phone: '081234567001', bnsp: 'BNSP-101-2023', sim: 'SIM A', photoUrl: '', experience: 5, bio: 'Expert in defensive driving techniques.', status: 'active' },
  { id: 2, name: 'Ms. Sari', phone: '081234567002', bnsp: 'BNSP-102-2022', sim: 'SIM A', photoUrl: '', experience: 8, bio: 'Patient and friendly, great for beginners.', status: 'active' },
  { id: 3, name: 'Mr. Budi', phone: '081234567003', bnsp: 'BNSP-105-2024', sim: 'SIM A', photoUrl: '', experience: 3, bio: 'Specialist in night driving and bad weather conditions.', status: 'active' }
])

// Instructor CRUD State
const isInstructorModalOpen = ref(false)
const isEditingInstructor = ref(false)
const instructorForm = ref({ id: 0, name: '', phone: '', bnsp: '', sim: '', photoUrl: '', experience: 0, bio: '', status: 'active' })

function openNewInstructor() {
  isEditingInstructor.value = false
  instructorForm.value = { id: 0, name: '', phone: '', bnsp: '', sim: '', photoUrl: '', experience: 0, bio: '', status: 'active' }
  isInstructorModalOpen.value = true
}

function openEditInstructor(instructor: any) {
  isEditingInstructor.value = true
  instructorForm.value = { ...instructor }
  isInstructorModalOpen.value = true
}

function saveInstructor() {
  if (isEditingInstructor.value) {
    const idx = instructors.value.findIndex(i => i.id === instructorForm.value.id)
    if (idx !== -1) instructors.value[idx] = { ...instructorForm.value }
    toast.add({ title: 'Instructor Updated', description: 'Instructor details saved.', color: 'success' })
  } else {
    instructors.value.push({
      ...instructorForm.value,
      id: Math.max(0, ...instructors.value.map(i => i.id)) + 1
    })
    toast.add({ title: 'Instructor Added', description: 'New instructor created.', color: 'success' })
  }
  isInstructorModalOpen.value = false
}

function deleteInstructor() {
  instructors.value = instructors.value.filter(i => i.id !== instructorForm.value.id)
  toast.add({ title: 'Instructor Deleted', description: 'Instructor has been removed.', color: 'error' })
  isInstructorModalOpen.value = false
}

// ==================== IMAGE UPLOAD ====================
const MAX_IMAGE_SIZE = 5 * 1024 * 1024 // 5MB

function handleImageUpload(event: Event, targetForm: any) {
  const input = event.target as HTMLInputElement
  if (!input.files || input.files.length === 0) return

  const file = input.files[0]
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
    toast.add({ title: 'Invalid File', description: 'Please upload JPG, PNG, or WebP.', color: 'error' })
    return
  }
  if (file.size > MAX_IMAGE_SIZE) {
    toast.add({ title: 'File Too Large', description: 'Image exceeds 5MB limit.', color: 'error' })
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    targetForm.photoUrl = e.target?.result as string
  }
  reader.readAsDataURL(file)
  input.value = ''
}

const vehicleImageRef = ref<HTMLInputElement | null>(null)
const instructorImageRef = ref<HTMLInputElement | null>(null)

function saveSettings() {
  toast.add({
    title: 'Settings Saved',
    description: 'Your settings have been updated successfully.',
    icon: 'i-lucide-check-circle',
    color: 'success'
  })
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Settings">
        <template #right>
          <UButton label="Save All Changes" color="warning" icon="i-lucide-save" @click="saveSettings" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6 max-w-full">
        <!-- General Settings -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-settings" class="size-5 text-warning" />
              <h2 class="font-semibold">General Settings</h2>
            </div>
          </template>

          <div class="grid md:grid-cols-2 gap-4">
            <UFormField label="Business Name">
              <UInput v-model="generalSettings.businessName" icon="i-lucide-building" class="w-full" color="warning"/>
            </UFormField>
            <UFormField label="Email">
              <UInput v-model="generalSettings.email" type="email" icon="i-lucide-mail" class="w-full" color="warning"/>
            </UFormField>
            <UFormField label="Phone Number">
              <UInput v-model="generalSettings.phone" icon="i-lucide-phone" class="w-full" color="warning"/>
            </UFormField>
            <UFormField label="WhatsApp Number">
              <UInput v-model="generalSettings.whatsappNumber" icon="i-simple-icons-whatsapp" class="w-full" color="warning"/>
            </UFormField>
            <UFormField label="Address" class="md:col-span-2">
              <UTextarea v-model="generalSettings.address" class="w-full" color="warning" />
            </UFormField>
          </div>
        </UCard>

        <!-- Operating Hours -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-clock" class="size-5 text-warning" />
              <h2 class="font-semibold">Operating Hours</h2>
            </div>
          </template>

          <div class="space-y-4">
            <div class="flex items-center gap-4">
              <span class="w-32 text-md font-medium">Monday - Friday</span>
              <UInput v-model="operatingHours.mondayStart" type="time" class="w-full" color="warning"/>
              <span class="text-muted">to</span>
              <UInput v-model="operatingHours.mondayEnd" type="time" class="w-full" color="warning" />
            </div>
            <div class="flex items-center gap-4">
              <span class="w-32 text-md font-medium">Saturday & Sunday</span>
              <UInput v-model="operatingHours.weekendStart" type="time" class="w-full" color="warning" />
              <span class="text-muted">to</span>
              <UInput v-model="operatingHours.weekendEnd" type="time" class="w-full" color="warning" />
            </div>
            <div class="flex items-center gap-4">
              <span class="w-32 text-md font-medium">Night Shift</span>
              <UInput v-model="operatingHours.nightStart" type="time" class="w-full" color="warning" />
              <span class="text-muted">to</span>
              <UInput v-model="operatingHours.nightEnd" type="time" class="w-full" color="warning" />
            </div>
          </div>
        </UCard>

        <!-- Vehicles -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-car" class="size-5 text-warning" />
                <h2 class="font-semibold">Vehicles</h2>
              </div>
              <UButton label="Add Vehicle" icon="i-lucide-plus" size="sm" variant="outline" color="warning" @click="openNewVehicle" />
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="vehicle in vehicles" 
              :key="vehicle.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <UAvatar :src="vehicle.photoUrl || undefined" :text="!vehicle.photoUrl ? vehicle.name.split(' ').map((n: string) => n[0]).join('') : undefined" size="md" />
                <div>
                  <p class="font-medium">{{ vehicle.name }}</p>
                  <p class="text-md text-muted">{{ vehicle.plate }}</p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <UBadge 
                  :label="vehicle.status === 'active' ? 'Active' : 'Maintenance'"
                  :color="vehicle.status === 'active' ? 'info' : 'error'"
                  variant="subtle"
                />
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="sm" @click="openEditVehicle(vehicle)" />
              </div>
            </div>
          </div>
        </UCard>

        <!-- Instructors -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-users" class="size-5 text-warning" />
                <h2 class="font-semibold">Instructors</h2>
              </div>
              <UButton label="Add Instructor" icon="i-lucide-user-plus" size="sm" variant="outline" color="warning" @click="openNewInstructor" />
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="instructor in instructors" 
              :key="instructor.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <UAvatar :src="instructor.photoUrl || undefined" :text="!instructor.photoUrl ? instructor.name.split(' ').map((n: string) => n[0]).join('') : undefined" size="md" />
                <div>
                  <p class="font-medium">{{ instructor.name }}</p>
                  <p class="text-md text-muted">{{ instructor.phone }}</p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <UBadge 
                  :label="instructor.status === 'active' ? 'Active' : 'Inactive'"
                  :color="instructor.status === 'active' ? 'info' : 'neutral'"
                  variant="subtle"
                />
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="sm" @click="openEditInstructor(instructor)" />
              </div>
            </div>
          </div>
        </UCard>

        <!-- Notification Settings -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-bell" class="size-5 text-warning" />
              <h2 class="font-semibold">Notification Settings</h2>
            </div>
          </template>

          <div class="space-y-4">
            <USwitch v-model="notificationSettings.emailNotifications" label="Send email notifications to students" />
            <USwitch v-model="notificationSettings.whatsappNotifications" label="Send WhatsApp reminders to students" />
            <USwitch v-model="notificationSettings.adminAlerts" label="Send admin alerts for general events" />
            <USwitch v-model="notificationSettings.newUserRegistration" label="Notify when new user registers" />
            <USwitch v-model="notificationSettings.newPackagePurchase" label="Notify when user buys a package / becomes a member" />
            
            <UFormField label="Send reminders before session (hours)">
              <UInput v-model="notificationSettings.reminderHours" type="number" min="1" max="72" class="w-full" />
            </UFormField>
          </div>
        </UCard>

        <!-- Danger Zone -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2 text-red-500">
              <UIcon name="i-lucide-alert-triangle" class="size-5" />
              <h2 class="font-semibold">Danger Zone</h2>
            </div>
          </template>

          <div class="space-y-4">
            <div class="flex items-center justify-between p-4 rounded-lg border border-red-200 dark:border-red-900">
              <div>
                <p class="font-medium">Export All Data</p>
                <p class="text-md text-muted">Download all system data as a backup file.</p>
              </div>
              <UButton label="Export" icon="i-lucide-download" color="neutral" variant="outline" />
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>

  <!-- Vehicle Modal -->
  <ClientOnly>
    <UModal v-model:open="isVehicleModalOpen">
      <template #content>
        <div class="bg-default rounded-2xl w-full">
          <div class="px-6 py-4 border-b border-default flex items-center justify-between">
            <h3 class="text-base font-semibold">{{ isEditingVehicle ? 'Edit Vehicle' : 'Add New Vehicle' }}</h3>
            <UButton icon="i-lucide-x" color="neutral" variant="ghost" @click="isVehicleModalOpen = false" />
          </div>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium mb-1.5">Vehicle Photo</label>
              <input ref="vehicleImageRef" type="file" accept="image/jpeg,image/png,image/webp" class="hidden" @change="e => handleImageUpload(e, vehicleForm)" />
              <div class="border-2 border-dashed border-default rounded-lg p-4 text-center cursor-pointer hover:border-primary transition-colors" @click="vehicleImageRef?.click()">
                <div v-if="vehicleForm.photoUrl" class="relative w-full h-32">
                  <img :src="vehicleForm.photoUrl" class="w-full h-full object-cover rounded-md" />
                  <UButton icon="i-lucide-trash-2" color="error" size="xs" class="absolute top-2 right-2" @click.stop="vehicleForm.photoUrl = ''" />
                </div>
                <div v-else class="flex flex-col items-center gap-2 py-4">
                  <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
                  <span class="text-sm text-muted">Click to upload photo</span>
                </div>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Vehicle Name</label>
              <UInput v-model="vehicleForm.name" placeholder="e.g. Tesla Model 3" class="w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">License Plate</label>
              <UInput v-model="vehicleForm.plate" placeholder="e.g. B 1234 EV" class="w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Status</label>
              <USelect v-model="vehicleForm.status" :items="['active', 'maintenance']" class="w-full" />
            </div>
          </div>
          <div class="px-6 py-4 border-t border-default flex justify-between items-center">
            <div>
              <UButton v-if="isEditingVehicle" label="Delete" icon="i-lucide-trash" color="error" variant="ghost" @click="deleteVehicle" />
            </div>
            <div class="flex gap-3">
              <UButton label="Cancel" color="neutral" variant="outline" @click="isVehicleModalOpen = false" />
              <UButton :label="isEditingVehicle ? 'Save Changes' : 'Add Vehicle'" icon="i-lucide-check" @click="saveVehicle" />
            </div>
          </div>
        </div>
      </template>
    </UModal>
  </ClientOnly>

  <!-- Instructor Modal -->
  <ClientOnly>
    <UModal v-model:open="isInstructorModalOpen">
      <template #content>
        <div class="bg-default rounded-2xl w-full">
          <div class="px-6 py-4 border-b border-default flex items-center justify-between">
            <h3 class="text-base font-semibold">{{ isEditingInstructor ? 'Edit Instructor' : 'Add New Instructor' }}</h3>
            <UButton icon="i-lucide-x" color="neutral" variant="ghost" @click="isInstructorModalOpen = false" />
          </div>
          <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
            <div class="grid grid-cols-2 gap-4">
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">Full Name</label>
                <UInput v-model="instructorForm.name" placeholder="e.g. Budi Santoso" class="w-full" />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">Phone Number</label>
                <UInput v-model="instructorForm.phone" placeholder="e.g. 08123456789" class="w-full" />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">BNSP Certificate No.</label>
                <UInput v-model="instructorForm.bnsp" placeholder="e.g. BNSP-101-2023" class="w-full" />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">SIM Type / Number</label>
                <UInput v-model="instructorForm.sim" placeholder="e.g. SIM A / 12345678" class="w-full" />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">Years of Experience</label>
                <UInput v-model="instructorForm.experience" type="number" min="0" class="w-full" />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">Status</label>
                <USelect v-model="instructorForm.status" :items="['active', 'inactive']" class="w-full" />
              </div>
              <div class="col-span-2">
                <label class="block text-sm font-medium mb-1.5">Instructor Photo</label>
                <input ref="instructorImageRef" type="file" accept="image/jpeg,image/png,image/webp" class="hidden" @change="e => handleImageUpload(e, instructorForm)" />
                <div class="border-2 border-dashed border-default rounded-lg p-4 text-center cursor-pointer hover:border-primary transition-colors" @click="instructorImageRef?.click()">
                  <div v-if="instructorForm.photoUrl" class="relative w-full h-32">
                    <img :src="instructorForm.photoUrl" class="w-full h-full object-contain rounded-md" />
                    <UButton icon="i-lucide-trash-2" color="error" size="xs" class="absolute top-2 right-2" @click.stop="instructorForm.photoUrl = ''" />
                  </div>
                  <div v-else class="flex flex-col items-center gap-2 py-4">
                    <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
                    <span class="text-sm text-muted">Click to upload photo</span>
                  </div>
                </div>
              </div>
              <div class="col-span-2">
                <label class="block text-sm font-medium mb-1.5">Bio Description</label>
                <UTextarea v-model="instructorForm.bio" placeholder="Short description of the instructor..." :rows="3" class="w-full" />
              </div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-default flex justify-between items-center">
            <div>
              <UButton v-if="isEditingInstructor" label="Delete" icon="i-lucide-trash" color="error" variant="ghost" @click="deleteInstructor" />
            </div>
            <div class="flex gap-3">
              <UButton label="Cancel" color="neutral" variant="outline" @click="isInstructorModalOpen = false" />
              <UButton :label="isEditingInstructor ? 'Save Changes' : 'Add Instructor'" icon="i-lucide-check" @click="saveInstructor" />
            </div>
          </div>
        </div>
      </template>
    </UModal>
  </ClientOnly>
</template>
