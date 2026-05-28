<script setup lang="ts">
import { ref } from 'vue'
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'

definePageMeta({
  layout: 'blank'
})

const loading = ref(false)
const showTermsModal = ref(false)

const schema = z.object({
  fullName: z.string().min(3, 'Name must be at least 3 characters'),
  agreedToTerms: z.boolean().refine(val => val === true, 'You must agree to terms')
})

const formData = reactive({
  fullName: '',
  agreedToTerms: false
})

async function onSubmit(event: FormSubmitEvent<any>) {
  loading.value = true

  try {
    // Simulate API call to save profile
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    console.log('Onboarding completed:', formData)
    
    let plan = null
    if (import.meta.client) {
      plan = sessionStorage.getItem('dm_selected_plan')
    }
    
    if (plan) {
      navigateTo(`/auth/payment-method?plan=${plan}`)
    } else {
      navigateTo('/auth/select-plan')
    }
  } finally {
    loading.value = false
  }
}

const features = [
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
                  I agree that my information will be used to issue my official training certificate and follow 
                  <UButton label="Terms of Service" color="warning" variant="ghost" class="underline" @click="showTermsModal = true" />
                  <UModal v-model:open="showTermsModal" title="Terms of Service">
                    <template #body>
                      <div class="prose dark:prose-invert max-w-none space-y-6">
                        <p>
                          Welcome to Drive Master Indonesia. These Terms of Service ("Terms") govern your access to and use of the website <NuxtLink to="/" class="text-warning hover:underline">www.drivemaster.id</NuxtLink> and our driving school services. By accessing or using our services, you agree to be bound by these Terms.
                        </p>

                        <h2 class="text-2xl font-bold">1. Services Provided</h2>
                        <p>Drive Master Indonesia provides professional driving instruction services, including theoretical training and practical on-road sessions. We reserve the right to modify, suspend, or discontinue any part of the services at any time without prior notice.</p>

                        <h2 class="text-2xl font-bold">2. User Accounts</h2>
                        <p>To access certain features of our platform, such as booking sessions, you must register for an account. You agree to:</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>Provide accurate, current, and complete information during the registration process.</li>
                          <li>Maintain the security of your password and accept all risks of unauthorized access to your account.</li>
                          <li>Notify us immediately if you discover or suspect any security breaches related to our services.</li>
                        </ul>

                        <h2 class="text-2xl font-bold">3. Fees and Payments</h2>
                        <p>All prices for our driving packages are listed in Indonesian Rupiah (IDR). Payment obligations include:</p>
                        <ul class="list-disc list-inside ml-4">
                          <li><strong>Payment Processing:</strong> Payments are processed securely via third-party gateways (e.g., Midtrans). We do not store your full financial credentials.</li>
                          <li><strong>Refund Policy:</strong> Requests for refunds are subject to our internal review and are typically only granted if requested at least 48 hours before the start of the first session.</li>
                        </ul>

                        <h2 class="text-2xl font-bold">4. Scheduling and Cancellations</h2>
                        <p>Efficient scheduling is key to our service quality. Our policy includes:</p>
                        <ul class="list-disc list-inside ml-4">
                          <li><strong>Booking:</strong> Sessions must be booked at least 24 hours in advance through the student dashboard.</li>
                          <li><strong>Rescheduling:</strong> You may reschedule a session through our platform at no extra cost if done at least 24 hours before the scheduled time.</li>
                          <li><strong>No-Show:</strong> Failure to attend a scheduled session without prior notice will result in the session being marked as completed and non-refundable.</li>
                        </ul>

                        <h2 class="text-2xl font-bold">5. Student Obligations</h2>
                        <p>As a student of Drive Master Indonesia, you agree to:</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>Possess a valid temporary or permanent driver's permit as required by local law.</li>
                          <li>Follow the instructions of the assigned instructor at all times during practical sessions.</li>
                          <li>Maintain a zero-tolerance policy regarding alcohol or drug use before or during training sessions.</li>
                        </ul>

                        <h2 class="text-2xl font-bold">6. Limitation of Liability</h2>
                        <p>
                          To the maximum extent permitted by law, Drive Master Indonesia shall not be liable for any indirect, incidental, or consequential damages resulting from your use of our services or any interaction with our instructors. While we strive for absolute safety, practical driving involves inherent risks.
                        </p>

                        <h2 class="text-2xl font-bold">7. Changes to Terms</h2>
                        <p>
                          We reserve the right to change or modify these Terms at any time. If we make changes, we will notify you by revising the date at the top of the policy or by posting a notice on our homepage. Your continued use of the services confirms your acceptance of the revised Terms.
                        </p>

                        <h2 class="text-2xl font-bold">8. Contact Us</h2>
                        <p>If you have any questions or concerns regarding these Terms, please reach out to us:</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>By email: <a href="mailto:info@drivemaster.id" class="text-warning hover:underline">info@drivemaster.id</a></li>
                          <li>By phone: +62 812-3456-7890</li>
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
    </div>
  </div>
</template>
