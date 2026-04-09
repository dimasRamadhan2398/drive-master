<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'

definePageMeta({ layout: 'dashboard' })

const toast = useToast()
const loading = ref(false)

const profileSchema = z.object({
  fullName: z.string().min(3, 'Name must be at least 3 characters'),
  email: z.string().email('Please enter a valid email'),
  phone: z.string().min(10, 'Please enter a valid phone number'),
  address: z.string().optional()
})

const passwordSchema = z.object({
  currentPassword: z.string().min(1, 'Current password is required'),
  newPassword: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string()
}).refine(data => data.newPassword === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword']
})

// Mock user data
const profileData = reactive({
  fullName: 'John Doe',
  email: 'john.doe@example.com',
  phone: '081234567890',
  address: 'Jl. Alam Sutera No. 123, Tangerang'
})

const passwordData = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const memberInfo = {
  memberId: 'EVDA-MEM-2026-0042',
  package: 'Standard Package',
  joinDate: 'March 10, 2026',
  expiryDate: 'September 10, 2026'
}

async function updateProfile(event: FormSubmitEvent<z.output<typeof profileSchema>>) {
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 1000))
  loading.value = false
  
  toast.add({
    title: 'Profile Updated',
    description: 'Your profile has been successfully updated.',
    icon: 'i-lucide-check-circle',
    color: 'success'
  })
}

async function updatePassword(event: FormSubmitEvent<z.output<typeof passwordSchema>>) {
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 1000))
  loading.value = false
  
  passwordData.currentPassword = ''
  passwordData.newPassword = ''
  passwordData.confirmPassword = ''
  
  toast.add({
    title: 'Password Changed',
    description: 'Your password has been successfully updated.',
    icon: 'i-lucide-check-circle',
    color: 'success'
  })
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Profile Settings">
        <template #right>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6 max-w-4xl">
        <!-- Member Card -->
        <UCard class="bg-gradient-to-r from-primary/10 to-primary/5">
          <div class="flex flex-col sm:flex-row items-center gap-6">
            <UAvatar text="JD" size="3xl" class="ring-4 ring-primary/20" />
            <div class="text-center sm:text-left">
              <h2 class="text-2xl font-bold">{{ profileData.fullName }}</h2>
              <p class="text-muted">{{ profileData.email }}</p>
              <div class="flex flex-wrap gap-2 mt-3 justify-center sm:justify-start">
                <UBadge :label="memberInfo.package" color="primary" />
                <UBadge :label="`ID: ${memberInfo.memberId}`" variant="subtle" />
              </div>
            </div>
          </div>
        </UCard>

        <!-- Membership Info -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">Membership Information</h2>
          </template>

          <div class="grid sm:grid-cols-2 gap-4">
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-sm text-muted">Member ID</p>
              <p class="font-medium font-mono">{{ memberInfo.memberId }}</p>
            </div>
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-sm text-muted">Package</p>
              <p class="font-medium">{{ memberInfo.package }}</p>
            </div>
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-sm text-muted">Join Date</p>
              <p class="font-medium">{{ memberInfo.joinDate }}</p>
            </div>
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-sm text-muted">Package Valid Until</p>
              <p class="font-medium">{{ memberInfo.expiryDate }}</p>
            </div>
          </div>
        </UCard>

        <!-- Personal Information -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">Personal Information</h2>
          </template>

          <UForm :schema="profileSchema" :state="profileData" class="space-y-4" @submit="updateProfile">
            <div class="grid sm:grid-cols-2 gap-4">
              <UFormField name="fullName" label="Full Name">
                <UInput v-model="profileData.fullName" icon="i-lucide-user" />
              </UFormField>

              <UFormField name="email" label="Email Address">
                <UInput v-model="profileData.email" type="email" icon="i-lucide-mail" />
              </UFormField>

              <UFormField name="phone" label="Phone Number">
                <UInput v-model="profileData.phone" icon="i-lucide-phone" />
              </UFormField>

              <UFormField name="address" label="Address">
                <UInput v-model="profileData.address" icon="i-lucide-map-pin" />
              </UFormField>
            </div>

            <div class="flex justify-end">
              <UButton type="submit" label="Save Changes" :loading="loading" icon="i-lucide-save" />
            </div>
          </UForm>
        </UCard>

        <!-- Change Password -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">Change Password</h2>
          </template>

          <UForm :schema="passwordSchema" :state="passwordData" class="space-y-4" @submit="updatePassword">
            <UFormField name="currentPassword" label="Current Password">
              <UInput v-model="passwordData.currentPassword" type="password" icon="i-lucide-lock" />
            </UFormField>

            <div class="grid sm:grid-cols-2 gap-4">
              <UFormField name="newPassword" label="New Password">
                <UInput v-model="passwordData.newPassword" type="password" icon="i-lucide-key" />
              </UFormField>

              <UFormField name="confirmPassword" label="Confirm New Password">
                <UInput v-model="passwordData.confirmPassword" type="password" icon="i-lucide-key" />
              </UFormField>
            </div>

            <div class="flex justify-end">
              <UButton type="submit" label="Update Password" :loading="loading" variant="outline" color="neutral" icon="i-lucide-shield" />
            </div>
          </UForm>
        </UCard>

        <!-- Notification Settings -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">Notification Preferences</h2>
          </template>

          <div class="space-y-4">
            <USwitch label="Email notifications for upcoming sessions" :default-checked="true" />
            <USwitch label="WhatsApp reminders (24 hours before session)" :default-checked="true" />
            <USwitch label="Promotional updates and offers" :default-checked="false" />
            <USwitch label="Newsletter subscription" :default-checked="false" />
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

          <div class="flex items-center justify-between">
            <div>
              <p class="font-medium">Delete Account</p>
              <p class="text-sm text-muted">Permanently delete your account and all associated data.</p>
            </div>
            <UButton label="Delete Account" color="error" variant="outline" icon="i-lucide-trash-2" />
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
