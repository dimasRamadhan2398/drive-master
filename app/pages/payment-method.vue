<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { computed, reactive, ref, onMounted } from 'vue'

definePageMeta({
  layout: 'blank'
})

const route = useRoute()
const router = useRouter()

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
  phone: z.string().min(10, 'Please enter a valid phone number')
})

const formData = reactive({
  paymentMethod: 'va',
  email: '',
  phone: ''
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
  navigateTo(`/payment?method=${formData.paymentMethod}&plan=${selectedPlan.value}&email=${formData.email}&phone=${formData.phone}`)
}

const packagePrices = {
  six_package: 'Rp 1.750.000',
  eight_package: 'Rp 1.950.000',
  ten_package: 'Rp 2.250.000',
  twelve_package: 'Rp 2.650.000'
}

const packageNames = {
  six_package: '6x Sessions',
  eight_package: '8x Sessions',
  ten_package: '10x Sessions',
  twelve_package: '12x Sessions'
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
          <div class="pt-4 border-t">
            <UCheckbox>
              <template #label>
                <span class="text-sm">
                  I agree to the 
                  <NuxtLink to="/terms" class="text-warning hover:underline">Terms of Service</NuxtLink>
                  and
                  <NuxtLink to="/privacy" class="text-warning hover:underline">Privacy Policy</NuxtLink>
                </span>
              </template>
            </UCheckbox>
          </div>

          <!-- Actions -->
          <div class="flex gap-3 pt-4 border-t">
            <NuxtLink to="/select-plan" class="flex-1">
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

        <!-- Footer removed -->
      </UCard>
    </div>
  </div>
</template>
