<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { reactive, ref } from 'vue'

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
  mondayEnd: '17:00',
  weekendStart: '08:00',
  weekendEnd: '17:00',
  nightStart: '18:00',
  nightEnd: '20:00',
  // saturdayStart: '08:00',
  // saturdayEnd: '16:00',
  // sundayClosed: true
})

const notificationSettings = reactive({
  emailNotifications: true,
  whatsappNotifications: true,
  reminderHours: 24,
  adminAlerts: true
})

const vehicles = ref([
  { id: 1, name: 'BYD Atto 1', plate: 'B 5678 EV', status: 'active' }
])

const instructors = ref([
  { id: 1, name: 'Mr. Ahmad', phone: '081234567001', status: 'active' },
  { id: 2, name: 'Ms. Sari', phone: '081234567002', status: 'active' },
  { id: 3, name: 'Mr. Budi', phone: '081234567003', status: 'active' }
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
              <span class="w-32 text-md font-medium">Saturday</span>
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
            <!-- <div class="flex items-center gap-4">
              <span class="w-32 text-md font-medium">Sunday</span>
              <USwitch v-model="operatingHours.sundayClosed" label="Closed" />
            </div> -->
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
              <UButton label="Add Vehicle" icon="i-lucide-plus" size="md" variant="outline" color="warning" />
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="vehicle in vehicles" 
              :key="vehicle.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <div class="p-2 rounded-lg bg-warning/10">
                  <UIcon name="i-lucide-car" class="size-5 text-warning" />
                </div>
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
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="md" />
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
              <UButton label="Add Instructor" icon="i-lucide-user-plus" size="md" variant="outline" color="warning" />
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="instructor in instructors" 
              :key="instructor.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <UAvatar :text="instructor.name.split(' ').map((n: string) => n[0]).join('')" size="md" />
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
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="md" />
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
            <USwitch v-model="notificationSettings.adminAlerts" label="Send admin alerts for new registrations" />
            
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
            <div class="flex items-center justify-between p-4 rounded-lg border border-red-200 dark:border-red-900">
              <div>
                <p class="font-medium">Clear Training History</p>
                <p class="text-md text-muted">Remove all completed training session records.</p>
              </div>
              <UButton label="Clear History" icon="i-lucide-trash" color="error" variant="outline" />
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
