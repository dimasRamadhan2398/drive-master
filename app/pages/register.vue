<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive, ref } from 'vue'
import { navigateTo } from 'nuxt/app'

const currentStep = ref(1)
const totalSteps = 3

const packageOptions = [
  { label: 'Starter - Rp 1.500.000', value: 'starter' },
  { label: 'Standard - Rp 2.800.000 (Popular)', value: 'standard' },
  { label: 'Pro - Rp 4.500.000', value: 'pro' }
]

// Step 1: Personal Info
const step1Schema = z.object({
  fullName: z.string().min(3, 'Name must be at least 3 characters'),
  email: z.string().email('Please enter a valid email'),
  phone: z.string().min(10, 'Please enter a valid phone number'),
  birthDate: z.string().min(1, 'Please enter your birth date')
})

// Step 2: Package Selection
const step2Schema = z.object({
  package: z.string().min(1, 'Please select a package'),
  startDate: z.string().optional()
})

// Step 3: Account
const step3Schema = z.object({
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
  terms: z.boolean().refine(val => val === true, 'You must accept the terms')
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword']
})

const formData = reactive({
  // Step 1
  fullName: '',
  email: '',
  phone: '',
  birthDate: '',
  // Step 2
  package: 'standard',
  startDate: '',
  // Step 3
  password: '',
  confirmPassword: '',
  terms: false
})

const loading = ref(false)

const stepSchemas = [step1Schema, step2Schema, step3Schema]

async function nextStep() {
  currentStep.value++
}

async function prevStep() {
  currentStep.value--
}

async function onSubmit(event: FormSubmitEvent<any>) {
  if (currentStep.value < totalSteps) {
    nextStep()
    return
  }
  
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 1500))
  loading.value = false
  
  navigateTo('/dashboard')
}

const stepItems = [
  { label: 'Personal Info', icon: 'i-lucide-user' },
  { label: 'Select Package', icon: 'i-lucide-package' },
  { label: 'Create Account', icon: 'i-lucide-shield-check' }
]
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4">
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon name="i-lucide-zap" class="size-8 text-warning" />
          <span class="text-xl font-bold">EV Drive Academy</span>
        </div>
        <h1 class="text-2xl font-bold">Create Your Account</h1>
        <p class="text-muted mt-2">Join our premium EV driving academy</p>
      </div>

      <!-- Stepper -->
      <UStepper 
        :items="stepItems" 
        :model-value="currentStep" 
        color="warning"
        class="mb-8"
      />

      <UCard>
        <!-- Step 1: Personal Info -->
        <UForm 
          v-if="currentStep === 1"
          :schema="step1Schema" 
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
            />
          </UFormField>

          <UFormField name="email" label="Email Address" required>
            <UInput 
              v-model="formData.email"
              type="email"
              placeholder="you@example.com"
              icon="i-lucide-mail"
              size="lg"
            />
          </UFormField>

          <UFormField name="phone" label="Phone Number (WhatsApp)" required>
            <UInput 
              v-model="formData.phone"
              placeholder="08123456789"
              icon="i-lucide-phone"
              size="lg"
            />
          </UFormField>

          <UFormField name="birthDate" label="Date of Birth" required>
            <UInput 
              v-model="formData.birthDate"
              type="date"
              size="lg"
            />
          </UFormField>

          <div class="flex justify-end pt-4">
            <UButton type="submit" label="Continue" color="warning" trailingIcon="i-lucide-arrow-right" size="lg" />
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
            :icon="formData.package === 'standard' ? 'i-lucide-star' : 'i-lucide-info'"
            :color="formData.package === 'standard' ? 'warning' : 'neutral'"
          >
            <template #title>
              {{ formData.package === 'starter' ? 'Starter Package' : formData.package === 'standard' ? 'Standard Package (Recommended)' : 'Pro Package' }}
            </template>
            <template #description>
              <ul class="mt-2 space-y-1 text-sm">
                <template v-if="formData.package === 'starter'">
                  <li>5 sessions x 45 minutes</li>
                  <li>Basic driving fundamentals</li>
                  <li>Certificate of completion</li>
                </template>
                <template v-else-if="formData.package === 'standard'">
                  <li>10 sessions x 60 minutes</li>
                  <li>Highway & city driving</li>
                  <li>EV-specific training + Official Certificate</li>
                </template>
                <template v-else>
                  <li>15 sessions x 90 minutes</li>
                  <li>Complete mastery program</li>
                  <li>Premium Certificate + Lifetime refresher</li>
                </template>
              </ul>
            </template>
          </UAlert>

          <UFormField name="startDate" label="Preferred Start Date (Optional)">
            <UInput 
              v-model="formData.startDate"
              type="date"
              size="lg"
            />
          </UFormField>

          <div class="flex justify-between pt-4">
            <UButton label="Back" variant="ghost" color="neutral" icon="i-lucide-arrow-left" @click="prevStep" />
            <UButton type="submit" label="Continue" color="warning" trailingIcon="i-lucide-arrow-right" size="lg" />
          </div>
        </UForm>

        <!-- Step 3: Account Creation -->
        <UForm 
          v-if="currentStep === 3"
          :schema="step3Schema" 
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
            />
          </UFormField>

          <UFormField name="confirmPassword" label="Confirm Password" required>
            <UInput 
              v-model="formData.confirmPassword"
              type="password"
              placeholder="Confirm your password"
              icon="i-lucide-lock"
              size="lg"
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

        <template #footer>
          <div class="text-center">
            <p class="text-sm text-muted">
              Already have an account?
              <NuxtLink to="/login" class="text-warning font-medium hover:underline">
                Sign in
              </NuxtLink>
            </p>
          </div>
        </template>
      </UCard>
    </div>
  </div>
</template>
