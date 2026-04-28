<script setup lang="ts">
import { ref } from 'vue'
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'

definePageMeta({
  layout: 'blank'
})

const loading = ref(false)

const schema = z.object({
  fullName: z.string().min(3, 'Name must be at least 3 characters'),
  ktpNumber: z.string().min(16, 'KTP number must be 16 digits').max(16, 'KTP number must be 16 digits'),
  agreedToTerms: z.boolean().refine(val => val === true, 'You must agree to terms')
})

const formData = reactive({
  fullName: '',
  ktpNumber: '',
  agreedToTerms: false
})

async function onSubmit(event: FormSubmitEvent<any>) {
  loading.value = true

  try {
    // Simulate API call to save profile
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    console.log('[v0] Onboarding completed:', formData)
    
    // Redirect to plan selection
    navigateTo('/auth/select-plan')
  } finally {
    loading.value = false
  }
}

const features = [
  {
    icon: 'i-lucide-check-circle',
    title: 'Verified Identity',
    description: 'Your KTP information will be used to issue your official certificate'
  },
  {
    icon: 'i-lucide-gift',
    title: 'Free Trial Access',
    description: 'Unlock your 15-minute free trial session after onboarding'
  },
  {
    icon: 'i-lucide-calendar',
    title: 'Schedule Sessions',
    description: 'Book your training sessions at your preferred times'
  },
  {
    icon: 'i-lucide-award',
    title: 'Earn Certificate',
    description: 'Receive your official Drive Master certificate upon completion'
  }
]
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4 bg-muted/20">
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-12">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon name="i-lucide-user-check" class="size-8 text-warning" />
          <span class="text-xl font-bold">Complete Your Profile</span>
        </div>
        <h1 class="text-3xl font-bold">Onboarding</h1>
        <p class="text-muted mt-3">
          Last step! Verify your identity to unlock your free trial and access all features.
        </p>
      </div>

      <!-- Features Grid -->
      <div class="grid sm:grid-cols-2 gap-4 mb-8">
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
          <UFormField name="fullName" label="Full Name (KTP Registered Name)" required>
            <UInput 
              v-model="formData.fullName"
              placeholder="Enter your full name as shown in your KTP"
              icon="i-lucide-user"
              size="lg"
              class="w-full"
            />
            <template #hint>
              This name will be used on your official certificate
            </template>
          </UFormField>

          <!-- KTP Number -->
          <UFormField name="ktpNumber" label="KTP Number (16 Digits)" required>
            <UInput 
              v-model="formData.ktpNumber"
              placeholder="e.g., 3520123456789012"
              icon="i-lucide-credit-card"
              size="lg"
              class="w-full"
              @input="formData.ktpNumber = formData.ktpNumber.replace(/\D/g, '').slice(0, 16)"
            />
            <template #hint>
              Your KTP data will be verified and used for certificate issuance
            </template>
          </UFormField>

          <!-- Info Alert -->
          <UAlert icon="i-lucide-shield-check" color="primary">
            <template #title>Your Data is Secure</template>
            <template #description>
              Your personal information is encrypted and only used for official certificate issuance. We follow Indonesia's data protection standards.
            </template>
          </UAlert>

          <!-- Terms -->
          <UFormField name="agreedToTerms">
            <UCheckbox v-model="formData.agreedToTerms" color="warning">
              <template #label>
                <span class="text-sm">
                  I agree that my KTP information will be used to issue my official training certificate and follow 
                  <NuxtLink to="/terms" class="text-warning hover:underline">Terms of Service</NuxtLink>
                </span>
              </template>
            </UCheckbox>
          </UFormField>

          <!-- Actions -->
          <div class="flex gap-3 pt-4 border-t">
            <NuxtLink to="/auth/register" class="flex-1">
              <UButton 
                label="Back" 
                color="neutral"
                variant="outline"
                block
              />
            </NuxtLink>
            <UButton 
              type="submit"
              label="Continue to Plans"
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
            <p>Email verification completed ✓</p>
          </div>
        </template>
      </UCard>

      <!-- Help Section
      <div class="mt-8 p-6 bg-muted/30 rounded-lg text-center">
        <p class="text-sm text-muted mb-3">
          Don't have your KTP number handy?
        </p>
        <NuxtLink to="https://wa.me/6281234567890" target="_blank">
          <UButton 
            label="Chat Support on WhatsApp"
            icon="i-simple-icons-whatsapp"
            color="neutral"
            variant="outline"
            size="sm"
          />
        </NuxtLink>
      </div>
    -->
    </div>
  </div>
</template>
