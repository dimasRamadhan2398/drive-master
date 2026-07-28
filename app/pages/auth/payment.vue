<script setup lang="ts">
import { ref } from 'vue'

const { t } = useI18n()
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
  six_package_night: { name: '6x Session + Night Session', price: 1850000, sessions: 6 },
  six_package_weekend: { name: '6x Session + Weekend Session', price: 1850000, sessions: 6 },
  six_package_weekend_night: { name: '6x Session + Weekend & Night Session', price: 1950000, sessions: 6 },
  eight_package: { name: '8x Session', price: 1950000, sessions: 8 },
  eight_package_night: { name: '8x Session + Night Session', price: 2100000, sessions: 8 },
  eight_package_weekend: { name: '8x Session + Weekend Session', price: 2100000, sessions: 8 },
  eight_package_weekend_night: { name: '8x Session + Weekend & Night Session', price: 2250000, sessions: 8 },
  ten_package: { name: '10x Session', price: 2250000, sessions: 10 },
  ten_package_night: { name: '10x Session + Night Session', price: 2450000, sessions: 10 },
  ten_package_weekend: { name: '10x Session + Weekend Session', price: 2450000, sessions: 10 },
  ten_package_weekend_night: { name: '10x Session + Weekend & Night Session', price: 2650000, sessions: 10 },
  twelve_package: { name: '12x Session', price: 2650000, sessions: 12 },
  twelve_package_night: { name: '12x Session + Night Session', price: 2900000, sessions: 12 },
  twelve_package_weekend: { name: '12x Session + Weekend Session', price: 2900000, sessions: 12 },
  twelve_package_weekend_night: { name: '12x Session + Weekend & Night Session', price: 3150000, sessions: 12 }
}

const paymentMethodDetails = computed(() => ({
  va: {
    title: t('billing.vaNumber'),
    icon: 'i-lucide-building',
    instructions: [
      'You will receive a Virtual Account number via email',
      'Transfer the exact amount to that Virtual Account',
      'Use your name as the transfer description',
      'Payment will be confirmed automatically within minutes',
      'You will receive confirmation via WhatsApp'
    ],
    note: t('billing.vaInstruction')
  },
  qris: {
    title: t('billing.scanQris'),
    icon: 'i-lucide-qr-code',
    instructions: [
      'A QR code will be displayed on the next screen',
      'Scan it with any e-wallet (GoPay, OVO, DANA, LinkAja)',
      'Complete the payment in your e-wallet app',
      'Payment confirmation is instant',
      'You will receive confirmation via WhatsApp'
    ],
    note: t('billing.qrisInstruction')
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
}))

const currentMethod = computed(() => 
  paymentMethodDetails.value[paymentMethod as keyof typeof paymentMethodDetails.value] || paymentMethodDetails.value.va
)

const pkg = computed(() => 
  packageInfo[selectedPlan as keyof typeof packageInfo] || packageInfo.eight_package
)

const paymentGateway = ref('doku') // Dynamic gateway selection: 'midtrans' or 'doku'
const paymentUrl = ref('')

// Initialize payment checkout from backend API
async function initiatePayment() {
  loading.value = true
  
  try {
    // In production, this would call your backend endpoint (e.g. POST /api/v1/payments/transactions)
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Set dynamic fields based on active gateway
    if (paymentGateway.value === 'doku') {
      paymentUrl.value = 'https://sandbox.doku.com/checkout/link/dummy_checkout_url'
      console.log('Doku payment initialized with URL:', paymentUrl.value)
    } else {
      snapToken.value = `mock_token_${Date.now()}`
      console.log('Midtrans payment initiated with token:', snapToken.value)
    }
    
    paymentInitiated.value = true
  } catch (error) {
    console.error('Payment initiation error:', error)
  } finally {
    loading.value = false
  }
}

async function completePayment() {
  loading.value = true
  
  // If Doku is active, trigger Jokul Checkout modal/redirect
  if (paymentGateway.value === 'doku' && paymentUrl.value) {
    loading.value = false
    if (typeof (window as any).loadJokulCheckout === 'function') {
      ;(window as any).loadJokulCheckout(paymentUrl.value)
    } else {
      window.location.href = paymentUrl.value
    }
    return
  }

  // Simulate payment success after 2 seconds for Midtrans/Mock
  await new Promise(resolve => setTimeout(resolve, 2000))
  loading.value = false
  
  // Redirect to payment status with success
  navigateTo(`/auth/payment-status?status=success&plan=${selectedPlan}&email=${email}`)
}

onMounted(() => {
  // Load gateway client scripts dynamically
  if (paymentGateway.value === 'doku') {
    const script = document.createElement('script')
    script.src = 'https://sandbox.doku.com/jokul-checkout-js/v1/jokul-checkout-1.0.0.js'
    script.async = true
    document.head.appendChild(script)
  } else {
    // Load Midtrans Snap script
    const script = document.createElement('script')
    script.src = 'https://app.midtrans.com/snap/snap.js'
    script.setAttribute('data-client-key', 'Mid-client-VRtikOzo00hMYblh')
    script.async = true
    document.head.appendChild(script)
  }
  
  initiatePayment()
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
        <h1 class="text-2xl font-bold">{{ t('billing.makePayment') }}</h1>
        <p class="text-muted mt-2">{{ pkg.name }} - Rp {{ (pkg.price).toLocaleString('id-ID') }}</p>
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <!-- Main Content -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Payment Instructions -->
          <UCard>
            <template #header>
              <h2 class="font-semibold">{{ t('billing.paymentInstruction') }}</h2>
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
                <UIcon name="i-lucide-alert-circle" class="text-warning font-bold shrink-0 mt-0.5" />
                <p class="text-sm text-warning font-bold">{{ currentMethod.note }}</p>
              </div>
            </div>
          </UCard>

          <!-- Payment Method Specific Info -->
          <UCard v-if="paymentMethod === 'va'">
            <template #header>
              <h2 class="font-semibold">{{ t('billing.vaNumber') }}</h2>
            </template>

            <div class="space-y-4">
              <div class="p-4 bg-muted rounded-lg">
                <p class="text-xs text-muted uppercase tracking-wide mb-2">{{ t('billing.vaNumber') }}</p>
                <p class="text-2xl font-mono font-bold break-all">8800123456789</p>
              </div>

              <div class="grid sm:grid-cols-2 gap-4">
                <div>
                  <p class="text-xs text-muted uppercase tracking-wide mb-1">Bank</p>
                  <p class="font-medium">BCA</p>
                </div>
                <div>
                  <p class="text-xs text-muted uppercase tracking-wide mb-1">{{ t('billing.amountDue').split(' ').shift() }}</p>
                  <p class="font-medium">Rp {{ ((pkg.price * 1.1).toLocaleString('id-ID')) }}</p>
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
              <h2 class="font-semibold">{{ t('billing.scanQris') }}</h2>
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
                {{ t('billing.qrisInstruction') }}
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
                <p class="text-2xl font-bold text-warning">Rp {{ ((pkg.price * 1.1).toLocaleString('id-ID')) }}</p>
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
              {{ t('billing.paymentInstructionDesc', { method: 'E-Wallet' }) }}
            </p>

            <UButton 
              label="Open Payment Link"
              icon="i-lucide-external-link"
              color="warning"
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
                <p class="text-sm text-muted">{{ t('billing.package') }}</p>
                <p class="font-medium">{{ pkg.name }}</p>
              </div>

              <div>
                <p class="text-sm text-muted">{{ t('billing.sessions') }}</p>
                <p class="font-medium">{{ pkg.sessions }}</p>
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
                  :label="t('common.confirm')"
                  icon="i-lucide-check"
                  color="warning"
                  :loading="loading"
                  @click="completePayment"
                  block
                />
                <NuxtLink to="/auth/payment-method">
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
