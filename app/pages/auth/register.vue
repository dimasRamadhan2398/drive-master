<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive, ref } from 'vue'
import { navigateTo } from 'nuxt/app'

definePageMeta({
  layout: 'blank'
})

const currentStep = ref(0)
const totalSteps = 3

const packageOptions = [
  { label: '6x Session - Rp 1.750.000', value: '6x' },
  { label: '8x Session - Rp 1.950.000 (Popular)', value: '8x' },
  { label: '10x Session - Rp 2.250.000', value: '10x' },
  { label: '12x Session - Rp 2.650.000', value: '12x' }
]

// Step 0: Personal Info
const step0Schema = z.object({
  fullName: z.string().min(3, 'Name must be at least 3 characters'),
  email: z.string().email('Please enter a valid email'),
  phone: z.string().min(10, 'Please enter a valid phone number'),
  birthDate: z.string().min(1, 'Please enter your birth date')
})

// Step 1: Package Selection
const step1Schema = z.object({
  package: z.string().min(1, 'Please select a package'),
  startDate: z.string().optional()
})

// Step 2: Account
const step2Schema = z.object({
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
  terms: z.boolean().refine(val => val === true, 'You must accept the terms')
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword']
})

const formData = reactive({
  // Step 0
  fullName: '',
  email: '',
  phone: '',
  birthDate: '',
  // Step 1
  package: '8x',
  startDate: '',
  // Step 2
  password: '',
  confirmPassword: '',
  terms: false
})

const loading = ref(false)

const stepSchemas = [step0Schema, step1Schema, step2Schema]

async function nextStep() {
  currentStep.value++
}

async function prevStep() {
  currentStep.value--
}

async function onSubmit(event: FormSubmitEvent<any>) {
  if (currentStep.value < totalSteps - 1) {
    nextStep()
    return
  }
  
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 1500))
  loading.value = false

  // Simpan data kontak ke sessionStorage agar bisa digunakan di halaman payment
  if (import.meta.client) {
    sessionStorage.setItem('dm_reg_email', formData.email)
    sessionStorage.setItem('dm_reg_phone', formData.phone)
  }

  console.log('[v0] Registration submitted, redirecting to verify:', formData.email)
  
  navigateTo(`/auth/verify?email=${encodeURIComponent(formData.email)}`)
}

const stepItems = [
  { label: 'Personal Info', icon: 'i-lucide-user' },
  { label: 'Create Account', icon: 'i-lucide-shield-check' },
  { label: 'Select Package', icon: 'i-lucide-package' }
]
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4">
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="flex items-center justify-center gap-2 mb-4">
          <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
        </div>
        <h1 class="text-2xl font-bold">Create Your Account</h1>
        <p class="text-muted mt-2">Join our Drive Master course</p>
      </div>

      <!-- Stepper -->
      <UStepper 
        :items="stepItems" 
        :model-value="currentStep" 
        color="warning"
        class="mb-8"
      />

      <UCard>
        <!-- Step 0: Personal Info -->
        <UForm 
          v-if="currentStep === 0"
          :schema="step0Schema" 
          :state="formData" 
          class="space-y-4"
          @submit="onSubmit"
        >
          <UFormField name="fullName" label="Full Name" required>
            <UInput 
              v-model="formData.fullName"
              placeholder="Enter your full name"
              icon="i-lucide-user"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="email" label="Email Address" required>
            <UInput 
              v-model="formData.email"
              type="email"
              placeholder="you@example.com"
              icon="i-lucide-mail"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="phone" label="Phone Number (WhatsApp)" required>
            <UInput 
              v-model="formData.phone"
              placeholder="08123456789"
              icon="i-lucide-phone"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="birthDate" label="Date of Birth" required>
            <UInput 
              v-model="formData.birthDate"
              type="date"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <div class="flex justify-end pt-4">
            <UButton type="submit" label="Continue" color="warning" trailingIcon="i-lucide-arrow-right" size="lg" />
          </div>
        </UForm>

        <!-- Step 1: Account Creation -->
        <UForm 
          v-if="currentStep === 1"
          :schema="step1Schema" 
          :state="formData" 
          class="space-y-4"
          @submit="onSubmit"
        >
          <UAlert icon="i-lucide-user-check" color="warning">
            <template #title>Almost there, {{ formData.fullName.split(' ')[0] }}!</template>
            <template #description>
              Create your password to secure your account.
            </template>
          </UAlert>

          <UFormField name="password" label="Password" required>
            <UInput 
              v-model="formData.password"
              type="password"
              placeholder="Create a strong password"
              icon="i-lucide-lock"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="confirmPassword" label="Confirm Password" required>
            <UInput 
              v-model="formData.confirmPassword"
              type="password"
              placeholder="Confirm your password"
              icon="i-lucide-lock"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="terms">
            <UCheckbox v-model="formData.terms" color="warning">
              <template #label>
                <span class="text-sm">
                  I agree to the 
                  <NuxtLink to="/terms" class="text-warning hover:underline">Terms of Service</NuxtLink>
                  and
                  <NuxtLink to="/privacy" class="text-warning hover:underline">Privacy Policy</NuxtLink>
                </span>
              </template>
            </UCheckbox>
          </UFormField>

          <div class="flex justify-between pt-4">
            <UButton label="Back" variant="ghost" color="neutral" icon="i-lucide-arrow-left" @click="prevStep" />
            <UButton 
              type="submit" 
              label="Create Account" 
              color="warning"
              :loading="loading"
              icon="i-lucide-check"
              size="lg"
            />
          </div>
        </UForm>

        <!-- Step 2: Package Selection -->
        <UForm 
          v-if="currentStep === 2"
          :schema="step2Schema" 
          :state="formData" 
          class="space-y-4"
          @submit="onSubmit"
        >
          <UFormField name="package" label="Select Package" required>
            <URadioGroup 
              v-model="formData.package"
              :items="packageOptions"
              orientation="vertical"
              color="warning"
            />
          </UFormField>

          <!-- Package Summary -->
          <UAlert 
            v-if="formData.package"
            :icon="formData.package === '8x' ? 'i-lucide-star' : 'i-lucide-info'"
            :color="formData.package === '8x' ? 'warning' : 'neutral'"
          >
            <template #title>
              {{ formData.package === '6x' ? '6x Session' : formData.package === '8x' ? '8x Session (Recommended)' : formData.package === '10x' ? '10x Session' : '12x Session' }}
            </template>
            <template #description>
              <ul class="mt-2 space-y-1 text-sm">
                <template v-if="formData.package === '6x'">
                  <li>Free Trial</li>
                  <li>6x Sessions</li>
                </template>
                <template v-else-if="formData.package === '8x'">
                  <li>Free Trial</li>
                  <li>8x Sessions</li>
                </template>
                <template v-else-if="formData.package === '10x'">
                  <li>Free Trial</li>
                  <li>10x Sessions</li>
                </template>
                <template v-else-if="formData.package === '12x'">
                  <li>Free Trial</li>
                  <li>12x Sessions</li>
                </template>
              </ul>
            </template>
          </UAlert>

          <div class="text-center">
            <p class="text-sm text-muted">
              Let me decide later
              <NuxtLink to="/auth/onboarding" class="text-warning font-medium hover:underline">
                Go to Onboarding page
              </NuxtLink>
            </p>
          </div>

          <!-- <UFormField name="startDate" label="Preferred Start Date (Optional)">
            <UInput 
              v-model="formData.startDate"
              type="date"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField> -->

          <div class="flex justify-between pt-4">
            <UButton label="Back" variant="ghost" color="neutral" icon="i-lucide-arrow-left" @click="prevStep" />
            <UButton type="submit" label="Proceed to Payment" color="warning" trailingIcon="i-lucide-arrow-right" size="lg" />
          </div>
        </UForm>

        <template #footer>
          <div class="text-center">
            <p class="text-sm text-muted">
              Already have an account?
              <NuxtLink to="/auth/login" class="text-warning font-medium hover:underline">
                Sign in
              </NuxtLink>
            </p>
          </div>
        </template>
      </UCard>
    </div>
  </div>
</template>
