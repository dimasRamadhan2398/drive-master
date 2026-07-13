<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'

const { t } = useI18n()
definePageMeta({
  layout: 'blank'
})

const auth = useAuthStore()
const loading = ref(false)
const showTermsModal = ref(false)

// Fetch member profile on mount
onMounted(async () => {
  if (auth.userId && !auth.memberProfile) {
    await auth.fetchMemberProfile()
  }
  // Pre-fill form if identity fullname already exists
  if (auth.memberIdentityFullname) {
    formData.fullName = auth.memberIdentityFullname
  }
})

const schema = computed(() => z.object({
  fullName: z.string().min(3, t('validation.minLength', { min: 3 })),
  agreedToTerms: z.boolean().refine(val => val === true, t('register.validation.termsRequired'))
}))

const formData = reactive({
  fullName: '',
  agreedToTerms: false
})

async function onSubmit(_event: FormSubmitEvent<any>) {
  loading.value = true

  try {
    // Update member profile with identity fullname
    const result = await auth.updateMemberProfile({
      identityFullname: formData.fullName
    })

    if (!result) {
      console.error('[ONBOARDING] Failed to update member profile')
      // Continue anyway since API might not be available
    } else {
      console.log('[ONBOARDING] Member profile updated successfully:', result)
    }

    let enrollmentId = null
    if (import.meta.client) {
      enrollmentId = sessionStorage.getItem('dm_enrollment_id')
    }

    if (enrollmentId) {
      navigateTo(`/auth/payment-method?enrollment=${enrollmentId}`)
    } else {
      navigateTo('/auth/select-plan')
    }
  } finally {
    loading.value = false;
  }
}

const features = computed(() => [
  {
    icon: 'i-lucide-gift',
    title: t('freeTrial.readyToSchedule'),
    description: t('freeTrial.heroDesc').split('.').shift()
  },
  {
    icon: 'i-lucide-calendar',
    title: t('schedule.title'),
    description: t('schedule.subtitle')
  },
  {
    icon: 'i-lucide-award',
    title: t('certificate.title'),
    description: t('certificate.notAvailableDesc').split('.').shift()
  }
])
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4 bg-muted/20">
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-12">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon name="i-lucide-user-check" class="size-8 text-warning" />
          <span class="text-xl font-bold">{{ t('auth.completeProfile') }}</span>
        </div>
        <h1 class="text-3xl font-bold">{{ t('auth.onboarding') }}</h1>
        <p class="text-muted mt-3">
          {{ t('auth.onboardingDesc') }}
        </p>
      </div>

      <!-- Features Grid -->
      <div class="grid sm:grid-cols-3 gap-4 mb-8">
        <div 
          v-for="feature in features" 
          :key="feature.title"
          class="flex gap-3 p-4 rounded-lg bg-muted/30"
        >
          <UIcon :name="feature.icon" class="size-5 text-warning shrink-0 mt-0.5" />
          <div>
            <p class="font-medium text-sm">{{ feature.title }}</p>
            <p class="text-xs text-muted mt-1">{{ feature.description }}</p>
          </div>
        </div>
      </div>

      <!-- Form -->
      <UCard>
        <UForm 
          :schema="schema" 
          :state="formData" 
          class="space-y-6"
          @submit="onSubmit"
        >
          <!-- Full Name -->
          <UFormField name="fullName" :label="t('auth.fullNameKtp')" required>
            <UInput 
              v-model="formData.fullName"
              :placeholder="t('auth.fullNamePlaceholder')"
              icon="i-lucide-user"
              size="lg"
              class="w-full"
            />
            <template #hint>
              {{ t('auth.nameHint') }}
            </template>
          </UFormField>

          <!-- Info Alert -->
          <UAlert icon="i-lucide-shield-check" color="primary">
            <template #title>{{ t('auth.dataSecure') }}</template>
            <template #description>
              {{ t('auth.dataSecureDesc') }}
            </template>
          </UAlert>

          <!-- Terms -->
          <UFormField name="agreedToTerms">
            <UCheckbox v-model:model-value="formData.agreedToTerms" color="warning">
              <template #label>
                <span class="text-sm">
                  {{ t('auth.agreeOnboarding') }}
                  <UButton :label="t('register.terms.termsOfService')" color="warning" variant="ghost" class="underline mx-0 px-0 h-auto py-0 text-sm" @click="showTermsModal = true" />
                  <UModal v-model:open="showTermsModal" :title="t('register.terms.termsOfService')">
                    <template #body>
                      <div class="prose dark:prose-invert max-w-none space-y-6">
                        <p>
                          {{ t('register.tosContent.welcome', { url: 'www.drivemaster.id' }) }}
                        </p>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.servicesTitle') }}</h2>
                        <p>
                          {{ t('register.tosContent.servicesDesc') }}
                        </p>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.accountsTitle') }}</h2>
                        <p>
                          {{ t('register.tosContent.accountsDesc') }}
                        </p>
                        <ul class="list-disc list-inside ml-4">
                          <li>{{ t('register.tosContent.accountsList1') }}</li>
                          <li>{{ t('register.tosContent.accountsList2') }}</li>
                          <li>{{ t('register.tosContent.accountsList3') }}</li>
                        </ul>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.feesTitle') }}</h2>
                        <p>{{ t('register.tosContent.feesDesc') }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li><strong>{{ t('register.tosContent.feesList1').split(':')[0] }}:</strong> {{ t('register.tosContent.feesList1').split(':')[1] }}</li>
                          <li><strong>{{ t('register.tosContent.feesList2').split(':')[0] }}:</strong> {{ t('register.tosContent.feesList2').split(':')[1] }}</li>
                        </ul>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.schedulingTitle') }}</h2>
                        <p>{{ t('register.tosContent.schedulingDesc') }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li><strong>{{ t('register.tosContent.schedulingList1').split(':')[0] }}:</strong> {{ t('register.tosContent.schedulingList1').split(':')[1] }}</li>
                          <li><strong>{{ t('register.tosContent.schedulingList2').split(':')[0] }}:</strong> {{ t('register.tosContent.schedulingList2').split(':')[1] }}</li>
                          <li><strong>{{ t('register.tosContent.schedulingList3').split(':')[0] }}:</strong> {{ t('register.tosContent.schedulingList3').split(':')[1] }}</li>
                        </ul>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.obligationsTitle') }}</h2>
                        <p>{{ t('register.tosContent.obligationsDesc') }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>{{ t('register.tosContent.obligationsList1') }}</li>
                          <li>{{ t('register.tosContent.obligationsList2') }}</li>
                          <li>{{ t('register.tosContent.obligationsList3') }}</li>
                        </ul>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.liabilityTitle') }}</h2>
                        <p>{{ t('register.tosContent.liabilityDesc') }}</p>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.changesTitle') }}</h2>
                        <p>{{ t('register.tosContent.changesDesc') }}</p>

                        <h2 class="text-2xl font-bold">{{ t('register.tosContent.contactTitle') }}</h2>
                        <p>{{ t('register.tosContent.contactDesc') }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>{{ t('register.tosContent.contactEmail', { email: 'info@drivemaster.id' }) }}</li>
                          <li>{{ t('register.tosContent.contactPhone') }}</li>
                        </ul>
                      </div>
                    </template>
                  </UModal>
                </span>
              </template>
            </UCheckbox>
          </UFormField>

          <!-- Actions -->
          <div class="flex gap-3 pt-4 border-t">
            <NuxtLink to="/auth/register" class="flex-1">
              <UButton 
                :label="t('auth.back')"
                color="neutral"
                variant="outline"
                block
              />
            </NuxtLink>
            <UButton 
              type="submit"
              :label="t('auth.continueToPlans')"
              trailingIcon="i-lucide-arrow-right"
              color="warning"
              :loading="loading"
              block
              class="flex-1"
            />
          </div>
        </UForm>

        <template #footer>
          <div class="text-center text-sm text-muted">
            <p>{{ t('auth.emailVerified') }}</p>
          </div>
        </template>
      </UCard>
    </div>
  </div>
</template>
