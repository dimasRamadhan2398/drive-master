<script setup lang="ts">
import { ref } from 'vue'

definePageMeta({
  layout: 'blank'
})

const route = useRoute()
const router = useRouter()

const paymentMethod = route.query.method as string || 'va'
const selectedPlan = route.query.plan as string || 'standard'
const email = route.query.email as string || ''
const phone = route.query.phone as string || ''

const loading = ref(false)
const paymentInitiated = ref(false)
const snapToken = ref<string | null>(null)

const packageInfo = {
  six_package: { name: '6x Session', price: 1750000, sessions: 6 },
  eight_package: { name: '8x Session', price: 1950000, sessions: 8 },
  ten_package: { name: '10x Session', price: 2250000, sessions: 10 },
  twelve_package: { name: '12x Session', price: 2650000, sessions: 12 }
}

const paymentMethodDetails = {
  va: {
    title: 'Virtual Account',
    icon: 'i-lucide-building',
    instructions: [
      'You will receive a Virtual Account number via email',
      'Transfer the exact amount to that Virtual Account',
      'Use your name as the transfer description',
      'Payment will be confirmed automatically within minutes',
      'You will receive confirmation via WhatsApp'
    ],
    note: 'Make sure to transfer the exact amount to ensure automatic confirmation.'
  },
  qris: {
    title: 'QRIS',
    icon: 'i-lucide-qr-code',
    instructions: [
      'A QR code will be displayed on the next screen',
      'Scan it with any e-wallet (GoPay, OVO, DANA, LinkAja)',
      'Complete the payment in your e-wallet app',
      'Payment confirmation is instant',
      'You will receive confirmation via WhatsApp'
    ],
    note: 'All major e-wallets are supported. Scan with your preferred app.'
  },
  bank_transfer: {
    title: 'Bank Transfer',
    icon: 'i-lucide-landmark',
    instructions: [
      'You will receive bank details via email',
      'Transfer from your bank account using ATM, mobile banking, or internet banking',
      'Use your name as the transfer description',
      'Confirmation within 1-2 hours during business hours',
      'You will receive confirmation via WhatsApp'
    ],
    note: 'Supported banks: BCA, Mandiri, BNI, BRI, and others.'
  },
  ewallet: {
    title: 'E-Wallet',
    icon: 'i-lucide-wallet',
    instructions: [
      'A payment link will be sent to your email',
      'Click the link or scan the displayed QR code',
      'Select your preferred e-wallet (GoPay, OVO, DANA)',
      'Complete payment in the e-wallet app',
      'Instant confirmation'
    ],
    note: 'E-wallets accepted: GoPay, OVO, DANA, LinkAja.'
  }
}

const currentMethod = computed(() => 
  paymentMethodDetails[paymentMethod as keyof typeof paymentMethodDetails] || paymentMethodDetails.va
)

const pkg = computed(() => 
  packageInfo[selectedPlan as keyof typeof packageInfo] || packageInfo.eight_package
)

// Mock Midtrans initialization
async function initiateMidtransPayment() {
  loading.value = true
  
  try {
    // In real implementation, this would call your backend to get Snap token
    // For now, we'll simulate it
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Mock snap token
    snapToken.value = `mock_token_${Date.now()}`
    paymentInitiated.value = true
    
    console.log('Midtrans payment initiated with token:', snapToken.value)
  } catch (error) {
    console.error('Payment initiation error:', error)
  } finally {
    loading.value = false
  }
}

async function completePayment() {
  // Simulate payment success after 2 seconds
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 2000))
  loading.value = false
  
  // Redirect to payment status with success
  navigateTo(`/payment-status?status=success&plan=${selectedPlan}&email=${email}`)
}

onMounted(() => {
  // Auto-initiate payment in real implementation
  initiateMidtransPayment()
})
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-12 bg-muted/20">
    <div class="max-w-full mx-auto">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon :name="currentMethod.icon" class="size-8 text-warning" />
          <span class="text-xl font-bold">{{ currentMethod.title }}</span>
        </div>
        <h1 class="text-2xl font-bold">Complete Payment</h1>
        <p class="text-muted mt-2">{{ pkg.name }} - Rp {{ (pkg.price).toLocaleString('id-ID') }}</p>
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <!-- Main Content -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Payment Instructions -->
          <UCard>
            <template #header>
              <h2 class="font-semibold">Payment Instructions</h2>
            </template>

            <ol class="space-y-4">
              <li 
                v-for="(instruction, index) in currentMethod.instructions" 
                :key="index"
                class="flex gap-4"
              >
                <div class="flex-shrink-0 flex items-center justify-center w-8 h-8 rounded-full bg-warning/10 text-warning font-medium text-sm">
                  {{ index + 1 }}
                </div>
                <div class="flex-1 pt-1">
                  <p class="text-sm">{{ instruction }}</p>
                </div>
              </li>
            </ol>

            <div class="mt-6 p-4 bg-amber-500/10 rounded-lg border border-amber-500/20">
              <div class="flex gap-2">
                <UIcon name="i-lucide-alert-circle" class="text-amber-500 shrink-0 mt-0.5" />
                <p class="text-sm text-amber-900">{{ currentMethod.note }}</p>
              </div>
            </div>
          </UCard>

          <!-- Payment Method Specific Info -->
          <UCard v-if="paymentMethod === 'va'">
            <template #header>
              <h2 class="font-semibold">Your Virtual Account Details</h2>
            </template>

            <div class="space-y-4">
              <div class="p-4 bg-muted rounded-lg">
                <p class="text-xs text-muted uppercase tracking-wide mb-2">Virtual Account Number</p>
                <p class="text-2xl font-mono font-bold break-all">8800123456789</p>
              </div>

              <div class="grid sm:grid-cols-2 gap-4">
                <div>
                  <p class="text-xs text-muted uppercase tracking-wide mb-1">Bank</p>
                  <p class="font-medium">BCA</p>
                </div>
                <div>
                  <p class="text-xs text-muted uppercase tracking-wide mb-1">Amount</p>
                  <p class="font-medium">Rp {{ (pkg.price).toLocaleString('id-ID') }}</p>
                </div>
              </div>

              <UButton 
                label="Copy Account Number"
                icon="i-lucide-copy"
                color="neutral"
                variant="outline"
                block
              />
            </div>
          </UCard>

          <UCard v-else-if="paymentMethod === 'qris'">
            <template #header>
              <h2 class="font-semibold">Scan QRIS Code</h2>
            </template>

            <div class="space-y-4">
              <div class="p-6 bg-muted rounded-lg flex items-center justify-center">
                <div class="w-48 h-48 bg-white rounded-lg border-4 border-warning flex items-center justify-center">
                  <div class="text-center">
                    <UIcon name="i-lucide-qr-code" class="size-24 text-muted mx-auto" />
                    <p class="text-xs text-muted mt-2">QR Code Placeholder</p>
                  </div>
                </div>
              </div>

              <p class="text-sm text-muted text-center">
                Open your e-wallet and scan this QR code to complete payment
              </p>
            </div>
          </UCard>

          <UCard v-else-if="paymentMethod === 'bank_transfer'">
            <template #header>
              <h2 class="font-semibold">Bank Transfer Details</h2>
            </template>

            <div class="space-y-4">
              <div class="grid sm:grid-cols-2 gap-4">
                <div>
                  <p class="text-xs text-muted uppercase tracking-wide mb-1">Bank Name</p>
                  <p class="font-medium">BCA (Bank Central Asia)</p>
                </div>
                <div>
                  <p class="text-xs text-muted uppercase tracking-wide mb-1">Account Number</p>
                  <p class="font-medium font-mono">1234567890</p>
                </div>
              </div>

              <div>
                <p class="text-xs text-muted uppercase tracking-wide mb-1">Account Holder</p>
                <p class="font-medium">EV Drive Academy</p>
              </div>

              <div>
                <p class="text-xs text-muted uppercase tracking-wide mb-1">Amount to Transfer</p>
                <p class="text-2xl font-bold text-warning">Rp {{ (pkg.price).toLocaleString('id-ID') }}</p>
              </div>

              <UButton 
                label="Copy Bank Details"
                icon="i-lucide-copy"
                color="neutral"
                variant="outline"
                block
              />
            </div>
          </UCard>

          <UCard v-else>
            <template #header>
              <h2 class="font-semibold">E-Wallet Payment</h2>
            </template>

            <p class="text-sm text-muted mb-4">
              A payment link has been sent to {{ email }}. Click the link to proceed with payment using your preferred e-wallet.
            </p>

            <UButton 
              label="Open Payment Link"
              icon="i-lucide-external-link"
              block
            />
          </UCard>

          <!-- Status Message -->
          <UAlert 
            v-if="!paymentInitiated"
            icon="i-lucide-loader-circle"
            color="neutral"
            class="animate-pulse"
          >
            <template #title>Preparing payment...</template>
          </UAlert>
        </div>

        <!-- Order Summary Sidebar -->
        <div>
          <UCard class="sticky top-4">
            <template #header>
              <h2 class="font-semibold">Order Summary</h2>
            </template>

            <div class="space-y-4">
              <div>
                <p class="text-sm text-muted">Package</p>
                <p class="font-medium">{{ pkg.name }}</p>
              </div>

              <div>
                <p class="text-sm text-muted">Sessions</p>
                <p class="font-medium">{{ pkg.sessions }} sessions</p>
              </div>

              <div class="pt-4 border-t space-y-3">
                <div class="flex justify-between">
                  <span class="text-sm">Subtotal</span>
                  <span class="font-medium">Rp {{ (pkg.price).toLocaleString('id-ID') }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-sm">Tax (10%)</span>
                  <span class="font-medium">Rp {{ ((pkg.price * 0.1).toLocaleString('id-ID')) }}</span>
                </div>
              </div>

              <div class="pt-3 border-t flex justify-between">
                <span class="font-semibold">Total</span>
                <span class="font-bold text-lg text-warning">Rp {{ ((pkg.price * 1.1).toLocaleString('id-ID')) }}</span>
              </div>
            </div>

            <template #footer>
              <div class="space-y-3">
                <UButton 
                  label="Confirm Payment Sent"
                  icon="i-lucide-check"
                  color="warning"
                  :loading="loading"
                  @click="completePayment"
                  block
                />
                <NuxtLink to="/payment-method">
                  <UButton 
                    label="Change Method"
                    color="neutral"
                    variant="outline"
                    icon="i-lucide-arrow-left"
                    block
                  />
                </NuxtLink>
              </div>
            </template>
          </UCard>
        </div>
      </div>
    </div>
  </div>
</template>
