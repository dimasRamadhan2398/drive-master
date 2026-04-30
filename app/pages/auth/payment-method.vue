<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { computed, reactive, ref, onMounted } from 'vue'

definePageMeta({
  layout: 'blank'
})

const route = useRoute()
const router = useRouter()
const showPrivacyModal = ref(false)
const showTermsModal = ref(false)

const selectedPlan = computed(() => route.query.plan as string || 'eight_package')

const paymentMethods = [
  {
    id: 'va',
    name: 'Virtual Account (VA)',
    description: 'Transfer via BCA, Mandiri, BNI, or BRI virtual account',
    icon: 'i-lucide-building',
    color: 'blue' as const
  },
  {
    id: 'qris',
    name: 'QRIS',
    description: 'Scan QR code with any e-wallet',
    icon: 'i-lucide-qr-code',
    color: 'green' as const
  },
  {
    id: 'bank_transfer',
    name: 'Bank Transfer',
    description: 'Direct transfer from your bank account',
    icon: 'i-lucide-landmark',
    color: 'purple' as const
  },
  {
    id: 'ewallet',
    name: 'E-Wallet',
    description: 'GoPay, OVO, DANA, or LinkAja',
    icon: 'i-lucide-wallet',
    color: 'orange' as const
  }
]

const schema = z.object({
  paymentMethod: z.string().min(1, 'Please select a payment method'),
  email: z.string().email('Please enter a valid email'),
  phone: z.string().min(10, 'Please enter a valid phone number'),
  privacy: z.boolean().refine((val) => val === true, 'You must agree to the privacy policy')
})

const formData = reactive({
  paymentMethod: 'va',
  email: '',
  phone: '',
  privacy: false
})

// Isi otomatis dari data registrasi yang tersimpan di sessionStorage
onMounted(() => {
  if (import.meta.client) {
    const savedEmail = sessionStorage.getItem('dm_reg_email')
    const savedPhone = sessionStorage.getItem('dm_reg_phone')
    if (savedEmail) formData.email = savedEmail
    if (savedPhone) formData.phone = savedPhone
  }
})

const loading = ref(false)

async function onSubmit() {
  loading.value = true
  
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1500))
  
  loading.value = false
  
  // Redirect to payment page
  navigateTo(`/auth/payment?method=${formData.paymentMethod}&plan=${selectedPlan.value}&email=${formData.email}&phone=${formData.phone}`)
}

const packagePrices = {
  six_package: 'Rp 1.750.000',
  six_package_night: 'Rp 1.850.000',
  six_package_weekend: 'Rp 1.850.000',
  six_package_weekend_night: 'Rp 1.950.000',
  eight_package: 'Rp 1.950.000',
  eight_package_night: 'Rp 2.100.000',
  eight_package_weekend: 'Rp 2.100.000',
  eight_package_weekend_night: 'Rp 2.250.000',
  ten_package: 'Rp 2.250.000',
  ten_package_night: 'Rp 2.450.000',
  ten_package_weekend: 'Rp 2.450.000',
  ten_package_weekend_night: 'Rp 2.650.000',
  twelve_package: 'Rp 2.650.000',
  twelve_package_night: 'Rp 2.900.000',
  twelve_package_weekend: 'Rp 2.900.000',
  twelve_package_weekend_night: 'Rp 3.150.000'
}

const packageNames = {
  six_package: '6x Sessions',
  six_package_night: '6x Sessions + Night Session',
  six_package_weekend: '6x Sessions + Weekend Session',
  six_package_weekend_night: '6x Sessions + Weekend & Night Session',
  eight_package: '8x Sessions',
  eight_package_night: '8x Sessions + Night Session',
  eight_package_weekend: '8x Sessions + Weekend Session',
  eight_package_weekend_night: '8x Sessions + Weekend & Night Session',
  ten_package: '10x Sessions',
  ten_package_night: '10x Sessions + Night Session',
  ten_package_weekend: '10x Sessions + Weekend Session',
  ten_package_weekend_night: '10x Sessions + Weekend & Night Session',
  twelve_package: '12x Sessions',
  twelve_package_night: '12x Sessions + Night Session',
  twelve_package_weekend: '12x Sessions + Weekend Session',
  twelve_package_weekend_night: '12x Sessions + Weekend & Night Session'
}
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4 bg-muted/20">
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon name="i-lucide-credit-card" class="size-8 text-warning" />
          <span class="text-xl font-bold">Payment Method</span>
        </div>
        <h1 class="text-2xl font-bold">Complete Your Purchase</h1>
        <p class="text-muted mt-2">
          {{ packageNames[selectedPlan as keyof typeof packageNames] }} - {{ packagePrices[selectedPlan as keyof typeof packagePrices] }}
        </p>
      </div>

      <UCard>
        <UForm 
          :schema="schema" 
          :state="formData" 
          class="space-y-6"
          @submit="onSubmit"
        >
          <!-- Payment Method Selection -->
          <div>
            <label class="block text-sm font-medium mb-4">Choose Payment Method</label>
            <div class="grid sm:grid-cols-2 gap-3">
              <div 
                v-for="method in paymentMethods" 
                :key="method.id"
                class="relative"
              >
                <input 
                  :id="`method-${method.id}`"
                  v-model="formData.paymentMethod"
                  type="radio" 
                  :value="method.id"
                  class="sr-only"
                />
                <label 
                  :for="`method-${method.id}`"
                  class="flex items-start gap-3 p-4 border-2 rounded-lg cursor-pointer transition-all"
                  :class="formData.paymentMethod === method.id 
                    ? 'border-warning bg-warning/5' 
                    : 'border-default hover:border-muted'"
                >
                  <div 
                    class="p-2 rounded-lg shrink-0 mt-1"
                    :class="`bg-${method.color}-500/10`"
                  >
                    <UIcon 
                      :name="method.icon" 
                      class="size-5"
                      :class="`text-${method.color}-500`"
                    />
                  </div>
                  <div class="flex-1">
                    <p class="font-medium">{{ method.name }}</p>
                    <p class="text-xs text-muted mt-1">{{ method.description }}</p>
                  </div>
                  <div 
                    v-if="formData.paymentMethod === method.id"
                    class="p-1 rounded-full bg-warning"
                  >
                    <UIcon name="i-lucide-check" class="size-4 text-white" />
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- Contact Information -->
          <div class="pt-4 border-t">
            <div class="flex items-start justify-between mb-4">
              <h3 class="font-medium">Contact Information</h3>
              <span class="text-xs text-muted flex items-center gap-1">
                <UIcon name="i-lucide-info" class="size-3" />
                Diambil dari data registrasi
              </span>
            </div>
            
            <UFormField name="email" label="Email Address" required>
              <UInput 
                v-model="formData.email"
                type="email"
                placeholder="your@email.com"
                icon="i-lucide-mail"
                size="lg"
                class="w-full"
              />
            </UFormField>

            <UFormField name="phone" label="WhatsApp Phone Number" required class="mt-4">
              <UInput 
                v-model="formData.phone"
                placeholder="08123456789"
                icon="i-lucide-phone"
                size="lg"
                class="w-full"
              />
            </UFormField>

            <UAlert icon="i-lucide-info" color="warning" class="mt-4">
              <template #description>
                We'll send payment instructions to this email and phone number for confirmation. You can edit the fields above if needed.
              </template>
            </UAlert>
          </div>

          <!-- Terms -->
          <UFormField name="privacy">
            <UCheckbox v-model="formData.privacy" color="warning" required>
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

          <!-- Actions -->
          <div class="flex gap-3 pt-4 border-t">
            <NuxtLink to="/auth/select-plan" class="flex-1">
              <UButton 
                label="Back to Packages" 
                color="neutral"
                variant="outline"
                block
              />
            </NuxtLink>
            <UButton 
              type="submit"
              label="Continue to Payment"
              trailingIcon="i-lucide-arrow-right"
              :loading="loading"
              block
              class="flex-1"
              color="warning"
            />
          </div>
        </UForm>

      </UCard>
    </div>
  </div>
</template>
