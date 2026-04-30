<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive, ref, computed } from 'vue'
import { navigateTo, useRoute } from 'nuxt/app'

definePageMeta({
  layout: 'blank'
})

const route = useRoute()
const planFromQuery = computed(() => route.query.plan as string | undefined)

const currentStep = ref(0)
const totalSteps = computed(() => planFromQuery.value ? 2 : 3)

const showPrivacyModal = ref(false)
const showTermsModal = ref(false)

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
  if (currentStep.value < totalSteps.value - 1) {
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
    if (planFromQuery.value) {
      sessionStorage.setItem('dm_selected_plan', planFromQuery.value)
    }
  }

  console.log('[v0] Registration submitted, redirecting to verify:', formData.email)
  
  navigateTo(`/auth/verify?email=${encodeURIComponent(formData.email)}`)
}

const stepItems = computed(() => {
  const items = [
    { label: 'Personal Info', icon: 'i-lucide-user' },
    { label: 'Create Account', icon: 'i-lucide-shield-check' }
  ]
  if (!planFromQuery.value) {
    items.push({ label: 'Select Package', icon: 'i-lucide-package' })
  }
  return items
})
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
                  and
                  <UButton label="Privacy Policy" color="warning" variant="ghost" class="underline" @click="showPrivacyModal = true" />
                  <UModal v-model:open="showPrivacyModal" title="Privacy Policy">
                      <template #body>
                        <div class="prose dark:prose-invert max-w-none space-y-6">
                          <p>
                            Welcome to Drive Master Indonesia. We are committed to protecting your privacy and ensuring you have a positive experience on our website and in using our services. This Privacy Policy outlines how we collect, use, disclose, and safeguard your information when you visit our website <NuxtLink to="/" class="text-warning hover:underline">www.drivemaster.id</NuxtLink> and use our driving school services.
                          </p>

                          <h2 class="text-2xl font-bold">1. Information We Collect</h2>
                          <p>We may collect personal information that you voluntarily provide to us when you register for our services, make a purchase, or interact with our website. This includes:</p>
                          <ul class="list-disc list-inside ml-4">
                            <li><strong>Personal Identification Information:</strong> Name, email address, phone number, physical address, date of birth, and driver's license details.</li>
                            <li><strong>Payment Information:</strong> Details required for processing payments, such as credit/debit card numbers (processed securely by third-party payment gateways like Midtrans).</li>
                            <li><strong>Usage Data:</strong> Information about how you access and use our website, including IP address, browser type, pages viewed, and time spent on pages.</li>
                          </ul>

                          <h2 class="text-2xl font-bold">2. How We Use Your Information</h2>
                          <p>The information we collect is used for various purposes, including:</p>
                          <ul class="list-disc list-inside ml-4">
                            <li>To provide and maintain our services, including scheduling driving lessons and managing your account.</li>
                            <li>To process transactions and send you related information, including purchase confirmations and invoices.</li>
                            <li>To improve our website and services based on your feedback and usage patterns.</li>
                            <li>To send you marketing and promotional communications (if you have opted in).</li>
                            <li>To comply with legal obligations and resolve disputes.</li>
                          </ul>

                          <h2 class="text-2xl font-bold">3. Disclosure of Your Information</h2>
                          <p>We may share your information with third parties in the following situations:</p>
                          <ul class="list-disc list-inside ml-4">
                            <li><strong>Service Providers:</strong> With third-party vendors, consultants, and other service providers who perform services for us or on our behalf (e.g., payment processing, email delivery, hosting services).</li>
                            <li><strong>Legal Requirements:</strong> If required to do so by law or in response to valid requests by public authorities (e.g., a court order or government agency).</li>
                            <li><strong>Business Transfers:</strong> In connection with, or during negotiations of, any merger, sale of company assets, financing, or acquisition of all or a portion of our business to another company.</li>
                          </ul>

                          <h2 class="text-2xl font-bold">4. Data Security</h2>
                          <p>
                            We implement reasonable security measures to protect your personal information from unauthorized access, use, alteration, and disclosure. However, no method of transmission over the Internet or electronic storage is 100% secure, and we cannot guarantee absolute security.
                          </p>

                          <h2 class="text-2xl font-bold">5. Your Data Protection Rights</h2>
                          <p>Depending on your location, you may have the following rights regarding your personal data:</p>
                          <ul class="list-disc list-inside ml-4">
                            <li>The right to access, update, or delete the information we have on you.</li>
                            <li>The right to object to our processing of your personal data.</li>
                            <li>The right to request that we restrict the processing of your personal information.</li>
                            <li>The right to data portability.</li>
                            <li>The right to withdraw consent at any time.</li>
                          </ul>
                          <p>To exercise any of these rights, please contact us at <a href="mailto:info@drivemaster.id" class="text-warning hover:underline">info@drivemaster.id</a>.</p>

                          <h2 class="text-2xl font-bold">6. Changes to This Privacy Policy</h2>
                          <p>
                            We may update our Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page and updating the "Last Updated" date. You are advised to review this Privacy Policy periodically for any changes.
                          </p>

                          <h2 class="text-2xl font-bold">7. Contact Us</h2>
                          <p>If you have any questions about this Privacy Policy, please contact us:</p>
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
