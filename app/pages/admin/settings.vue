<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const toast = useToast()

// Mock settings
const generalSettings = reactive({
  businessName: 'EV Drive Academy',
  email: 'info@evdriveacademy.id',
  phone: '+62 812-3456-7890',
  address: 'Jl. Alam Sutera Boulevard No. 123, Tangerang 15143',
  whatsappNumber: '6281234567890'
})

const operatingHours = reactive({
  mondayStart: '08:00',
  mondayEnd: '18:00',
  saturdayStart: '08:00',
  saturdayEnd: '16:00',
  sundayClosed: true
})

const notificationSettings = reactive({
  emailNotifications: true,
  whatsappNotifications: true,
  reminderHours: 24,
  adminAlerts: true
})

const vehicles = ref([
  { id: 1, name: 'Tesla Model 3', plate: 'B 1234 EV', status: 'active' },
  { id: 2, name: 'BYD Atto 3', plate: 'B 5678 EV', status: 'active' }
])

const instructors = ref([
  { id: 1, name: 'Pak Ahmad', phone: '081234567001', status: 'active' },
  { id: 2, name: 'Bu Sari', phone: '081234567002', status: 'active' },
  { id: 3, name: 'Pak Budi', phone: '081234567003', status: 'active' }
])

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
          <UButton label="Save All Changes" icon="i-lucide-save" @click="saveSettings" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6 max-w-4xl">
        <!-- General Settings -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-settings" class="size-5 text-primary" />
              <h2 class="font-semibold">General Settings</h2>
            </div>
          </template>

          <div class="grid sm:grid-cols-2 gap-4">
            <UFormField label="Business Name">
              <UInput v-model="generalSettings.businessName" icon="i-lucide-building" />
            </UFormField>
            <UFormField label="Email">
              <UInput v-model="generalSettings.email" type="email" icon="i-lucide-mail" />
            </UFormField>
            <UFormField label="Phone Number">
              <UInput v-model="generalSettings.phone" icon="i-lucide-phone" />
            </UFormField>
            <UFormField label="WhatsApp Number">
              <UInput v-model="generalSettings.whatsappNumber" icon="i-simple-icons-whatsapp" />
            </UFormField>
            <UFormField label="Address" class="sm:col-span-2">
              <UTextarea v-model="generalSettings.address" />
            </UFormField>
          </div>
        </UCard>

        <!-- Operating Hours -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-clock" class="size-5 text-primary" />
              <h2 class="font-semibold">Operating Hours</h2>
            </div>
          </template>

          <div class="space-y-4">
            <div class="flex items-center gap-4">
              <span class="w-32 text-sm font-medium">Monday - Friday</span>
              <UInput v-model="operatingHours.mondayStart" type="time" class="w-32" />
              <span class="text-muted">to</span>
              <UInput v-model="operatingHours.mondayEnd" type="time" class="w-32" />
            </div>
            <div class="flex items-center gap-4">
              <span class="w-32 text-sm font-medium">Saturday</span>
              <UInput v-model="operatingHours.saturdayStart" type="time" class="w-32" />
              <span class="text-muted">to</span>
              <UInput v-model="operatingHours.saturdayEnd" type="time" class="w-32" />
            </div>
            <div class="flex items-center gap-4">
              <span class="w-32 text-sm font-medium">Sunday</span>
              <USwitch v-model="operatingHours.sundayClosed" label="Closed" />
            </div>
          </div>
        </UCard>

        <!-- Vehicles -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-car" class="size-5 text-primary" />
                <h2 class="font-semibold">Vehicles</h2>
              </div>
              <UButton label="Add Vehicle" icon="i-lucide-plus" size="sm" variant="outline" color="neutral" />
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="vehicle in vehicles" 
              :key="vehicle.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <div class="p-2 rounded-lg bg-primary/10">
                  <UIcon name="i-lucide-car" class="size-5 text-primary" />
                </div>
                <div>
                  <p class="font-medium">{{ vehicle.name }}</p>
                  <p class="text-sm text-muted">{{ vehicle.plate }}</p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <UBadge 
                  :label="vehicle.status === 'active' ? 'Active' : 'Maintenance'"
                  :color="vehicle.status === 'active' ? 'success' : 'warning'"
                  variant="subtle"
                />
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="xs" />
              </div>
            </div>
          </div>
        </UCard>

        <!-- Instructors -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-users" class="size-5 text-primary" />
                <h2 class="font-semibold">Instructors</h2>
              </div>
              <UButton label="Add Instructor" icon="i-lucide-user-plus" size="sm" variant="outline" color="neutral" />
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="instructor in instructors" 
              :key="instructor.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <UAvatar :text="instructor.name.split(' ').map((n: string) => n[0]).join('')" size="sm" />
                <div>
                  <p class="font-medium">{{ instructor.name }}</p>
                  <p class="text-sm text-muted">{{ instructor.phone }}</p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <UBadge 
                  :label="instructor.status === 'active' ? 'Active' : 'Inactive'"
                  :color="instructor.status === 'active' ? 'success' : 'neutral'"
                  variant="subtle"
                />
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="xs" />
              </div>
            </div>
          </div>
        </UCard>

        <!-- Notification Settings -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-bell" class="size-5 text-primary" />
              <h2 class="font-semibold">Notification Settings</h2>
            </div>
          </template>

          <div class="space-y-4">
            <USwitch v-model="notificationSettings.emailNotifications" label="Send email notifications to students" />
            <USwitch v-model="notificationSettings.whatsappNotifications" label="Send WhatsApp reminders to students" />
            <USwitch v-model="notificationSettings.adminAlerts" label="Send admin alerts for new registrations" />
            
            <UFormField label="Send reminders before session (hours)">
              <UInput v-model="notificationSettings.reminderHours" type="number" min="1" max="72" class="w-32" />
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
                <p class="text-sm text-muted">Download all system data as a backup file.</p>
              </div>
              <UButton label="Export" icon="i-lucide-download" color="neutral" variant="outline" />
            </div>
            <div class="flex items-center justify-between p-4 rounded-lg border border-red-200 dark:border-red-900">
              <div>
                <p class="font-medium">Clear Training History</p>
                <p class="text-sm text-muted">Remove all completed training session records.</p>
              </div>
              <UButton label="Clear History" icon="i-lucide-trash" color="error" variant="outline" />
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
