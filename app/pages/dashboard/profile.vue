<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { reactive, ref, computed, onMounted } from 'vue'

const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

const toast = useToast()
const loading = ref(false)
const authStore = useAuthStore()

const profileSchema = computed(() => z.object({
  fullName: z.string().min(3, t('validation.minLength', { min: 3 })),
  email: z.string().email(t('validation.email')),
  phone: z.string().min(10, t('validation.phone')),
  address: z.string().optional()
}))

const passwordSchema = computed(() => z.object({
  currentPassword: z.string().min(1, t('validation.required')),
  newPassword: z.string().min(8, t('validation.password')),
  confirmPassword: z.string()
}).refine(data => data.newPassword === data.confirmPassword, {
  message: t('validation.passwordMatch'),
  path: ['confirmPassword']
}))

// Live profile data bound to inputs
const profileData = reactive({
  fullName: '',
  email: '',
  phone: '',
  address: ''
})

const passwordData = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const avatarText = computed(() => {
  const names = profileData.fullName.trim().split(' ');
  if (names.length >= 2) {
    return (names[0][0] + names[1][0]).toUpperCase();
  }
  return names[0] ? names[0].substring(0, 2).toUpperCase() : 'JD';
});

const memberInfo = computed(() => {
  const entitlement = authStore.memberProfile?.entitlements?.[0];
  
  const formatDate = (dateStr?: string | null) => {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return {
    memberId: authStore.memberProfile?.id ? `MEM-${authStore.memberProfile.id.slice(0, 8).toUpperCase()}` : '-',
    package: entitlement?.packageName || 'No active package',
    joinDate: formatDate(entitlement?.startDate || authStore.user?.createdAt),
    expiryDate: entitlement?.endDate ? formatDate(entitlement.endDate) : 'Lifetime / No Expiry'
  };
})

onMounted(async () => {
  loading.value = true
  await Promise.all([
    authStore.fetchCurrentUser(),
    authStore.fetchMemberProfile()
  ])
  
  if (authStore.user) {
    profileData.fullName = authStore.memberProfile?.identityFullname || 
                           ((authStore.user.firstName || '') + ' ' + (authStore.user.lastName || '')).trim() || 
                           authStore.user.username;
    profileData.email = authStore.user.email;
    profileData.phone = authStore.user.phoneNumber || '';
    profileData.address = authStore.user.address || '';
  }
  loading.value = false
})

async function updateProfile(event: FormSubmitEvent<any>) {
  loading.value = true
  
  const nameParts = profileData.fullName.trim().split(/\s+/);
  const firstName = nameParts[0] || '';
  const lastName = nameParts.slice(1).join(' ') || firstName;

  try {
    await Promise.all([
      authStore.updateUser({
        firstName,
        lastName,
        phoneNumber: profileData.phone,
        address: profileData.address,
      }),
      authStore.updateMemberProfile({
        identityFullname: profileData.fullName
      })
    ]);

    toast.add({
      title: t('profile.profileUpdated'),
      description: t('profile.profileUpdatedDesc'),
      icon: 'i-lucide-check-circle',
      color: 'success'
    })
  } catch (error: any) {
    toast.add({
      title: 'Failed to update profile',
      description: error?.message || 'Unknown error occurred',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}

async function updatePassword(event: FormSubmitEvent<any>) {
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 1000))
  loading.value = false
  
  passwordData.currentPassword = ''
  passwordData.newPassword = ''
  passwordData.confirmPassword = ''
  
  toast.add({
    title: t('profile.passwordChanged'),
    description: t('profile.passwordChangedDesc'),
    icon: 'i-lucide-check-circle',
    color: 'success'
  })
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('profile.title')">
        <template #right>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6 max-w-4xl">
        <!-- Member Card -->
        <UCard class="bg-gradient-to-r from-warning/10 to-warning/5">
          <div class="flex flex-col md:flex-row items-center gap-6">
            <UAvatar :src="authStore.user?.image" :text="avatarText" size="3xl" class="ring-4 ring-warning/20" />
            <div class="text-center md:text-left">
              <h2 class="text-2xl font-bold">{{ profileData.fullName }}</h2>
              <p class="text-muted">{{ profileData.email }}</p>
              <div class="flex flex-wrap gap-2 mt-3 justify-center md:justify-start">
                <UBadge :label="memberInfo.package" color="warning" />
                <UBadge :label="`ID: ${memberInfo.memberId}`" variant="subtle" />
              </div>
            </div>
          </div>
        </UCard>

        <!-- Membership Info -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">{{ t('profile.membershipInfo') }}</h2>
          </template>

          <div class="grid md:grid-cols-2 gap-4">
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-md text-muted">{{ t('profile.memberId') }}</p>
              <p class="font-medium font-mono">{{ memberInfo.memberId }}</p>
            </div>
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-md text-muted">{{ t('billing.package') }}</p>
              <p class="font-medium">{{ memberInfo.package }}</p>
            </div>
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-md text-muted">{{ t('profile.joinDate') }}</p>
              <p class="font-medium">{{ memberInfo.joinDate }}</p>
            </div>
            <div class="p-4 rounded-lg bg-muted/50">
              <p class="text-md text-muted">{{ t('profile.validUntil') }}</p>
              <p class="font-medium">{{ memberInfo.expiryDate }}</p>
            </div>
          </div>
        </UCard>

        <!-- Personal Information -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">{{ t('profile.personalInfo') }}</h2>
          </template>

          <UForm :schema="profileSchema" :state="profileData" class="space-y-4" @submit="updateProfile">
            <div class="grid md:grid-cols-2 gap-4">
              <UFormField name="fullName" :label="t('profile.fullName')">
                <UInput v-model="profileData.fullName" icon="i-lucide-user" class="w-full"/>
              </UFormField>

              <UFormField name="email" :label="t('profile.email')">
                <UInput v-model="profileData.email" type="email" icon="i-lucide-mail" class="w-full"/>
              </UFormField>

              <UFormField name="phone" :label="t('profile.phone')">
                <UInput v-model="profileData.phone" icon="i-lucide-phone" class="w-full"/>
              </UFormField>

              <UFormField name="address" :label="t('profile.address')">
                <UInput v-model="profileData.address" icon="i-lucide-map-pin" class="w-full"/>
              </UFormField>
            </div>

            <div class="flex justify-end">
              <UButton type="submit" :label="t('profile.saveChanges')" :loading="loading" icon="i-lucide-save" />
            </div>
          </UForm>
        </UCard>

        <!-- Change Password -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">{{ t('profile.changePassword') }}</h2>
          </template>

          <UForm :schema="passwordSchema" :state="passwordData" class="space-y-4" @submit="updatePassword">
            <UFormField name="currentPassword" :label="t('profile.currentPassword')">
              <UInput v-model="passwordData.currentPassword" type="password" icon="i-lucide-lock" class="w-full"/>
            </UFormField>

            <div class="grid md:grid-cols-2 gap-4">
              <UFormField name="newPassword" :label="t('profile.newPassword')">
                <UInput v-model="passwordData.newPassword" type="password" icon="i-lucide-key" class="w-full"/>
              </UFormField>

              <UFormField name="confirmPassword" :label="t('profile.confirmNewPassword')">
                <UInput v-model="passwordData.confirmPassword" type="password" icon="i-lucide-key" class="w-full"/>
              </UFormField>
            </div>

            <div class="flex justify-end">
              <UButton type="submit" :label="t('profile.changePassword')" :loading="loading" variant="outline" color="neutral" icon="i-lucide-shield" />
            </div>
          </UForm>
        </UCard>

        <!-- Notification Settings -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">{{ t('profile.notificationPrefs') }}</h2>
          </template>

          <div class="space-y-4">
            <USwitch :label="t('profile.notifEmail')" :default-checked="true" />
            <USwitch :label="t('profile.notifWa')" :default-checked="true" />
            <USwitch :label="t('profile.notifPromo')" :default-checked="false" />
            <USwitch :label="t('profile.notifNewsletter')" :default-checked="false" />
          </div>
        </UCard>

        <!-- Danger Zone -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2 text-red-500">
              <UIcon name="i-lucide-alert-triangle" class="size-5" />
              <h2 class="font-semibold">{{ t('profile.dangerZone') }}</h2>
            </div>
          </template>

          <div class="flex items-center justify-between">
            <div>
              <p class="font-medium">{{ t('profile.deleteAccount') }}</p>
              <p class="text-md text-muted">{{ t('profile.deleteAccountDesc') }}</p>
            </div>
            <UButton :label="t('profile.deleteAccount')" color="error" variant="outline" icon="i-lucide-trash-2" />
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
