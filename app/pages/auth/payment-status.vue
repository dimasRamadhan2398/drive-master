<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { usePaymentsStore } from '~/stores/payments'

const { t } = useI18n()
definePageMeta({
  layout: 'blank'
})

const route = useRoute()
const router = useRouter()
const paymentsStore = usePaymentsStore()

const status = ref(route.query.status as string || 'pending')
const email = ref(route.query.email as string || '')
const plan = ref(route.query.plan as string || 'standard')
const orderId = ref(route.query.orderId as string || '')
const enrollmentId = ref<string | null>(null)

const isSuccess = computed(() => status.value === 'success' || status.value === 'paid')
const isFailed = computed(() => status.value === 'failed')

let pollInterval: any = null

import { paymentService } from '~/services/paymentService'

const isCheckingStatus = ref(false)
const isSimulating = ref(false)

const pollStatus = async (oid: string) => {
  let attempts = 0
  const maxAttempts = 15 // Poll for 30 seconds max
  
  const check = async () => {
    attempts++
    const res = await paymentsStore.checkPaymentStatus(oid)
    if (res === 'success' || res === 'paid') {
      status.value = 'success'
      if (pollInterval) clearInterval(pollInterval)
      setTimeout(() => {
        navigateTo('/dashboard?just_paid=true')
      }, 3000)
    } else if (res === 'failed' || res === 'expired' || res === 'cancelled') {
      status.value = 'failed'
      if (pollInterval) clearInterval(pollInterval)
    } else if (attempts >= maxAttempts) {
      if (pollInterval) clearInterval(pollInterval)
    }
  }

  await check()
  
  if (status.value === 'pending') {
    pollInterval = setInterval(check, 2000)
  }
}

const manualCheckStatus = async () => {
  if (!orderId.value) return
  isCheckingStatus.value = true
  const toast = useToast()
  try {
    const res = await paymentsStore.checkPaymentStatus(orderId.value)
    if (res === 'success' || res === 'paid') {
      status.value = 'success'
      toast.add({
        title: 'Payment Confirmed!',
        description: 'Your payment has been successfully processed.',
        color: 'success'
      })
      setTimeout(() => {
        navigateTo('/dashboard?just_paid=true')
      }, 2000)
    } else if (res === 'failed' || res === 'expired' || res === 'cancelled') {
      status.value = 'failed'
    } else {
      toast.add({
        title: 'Payment Still Pending',
        description: 'Payment is still processing. Please try again in a few moments or click Simulate Payment.',
        color: 'info'
      })
    }
  } finally {
    isCheckingStatus.value = false
  }
}

const simulatePayment = async () => {
  if (!orderId.value) return
  isSimulating.value = true
  const toast = useToast()
  try {
    const ok = await paymentService.simulate(orderId.value)
    if (ok) {
      status.value = 'success'
      toast.add({
        title: 'Payment Simulated Successfully!',
        description: 'Payment marked as completed. Redirecting to dashboard...',
        color: 'success'
      })
      setTimeout(() => {
        navigateTo('/dashboard?just_paid=true')
      }, 2000)
    } else {
      toast.add({
        title: 'Simulation Failed',
        description: 'Could not simulate payment. Please try refreshing status.',
        color: 'error'
      })
    }
  } catch (err: any) {
    toast.add({
      title: 'Simulation Error',
      description: err?.message || 'Error simulating payment.',
      color: 'error'
    })
  } finally {
    isSimulating.value = false
  }
}

onMounted(() => {
  if (import.meta.client) {
    enrollmentId.value = sessionStorage.getItem("dm_enrollment_id") || null
  }

  // Resolve orderId from sessionStorage fallback if not in query
  if (!orderId.value && import.meta.client) {
    orderId.value = sessionStorage.getItem("dm_order_id") || ''
  }

  // Trigger status check & polling if we have orderId
  if (orderId.value) {
    paymentsStore.fetchPaymentByOrderId(orderId.value).then((payment) => {
      if (payment && payment.enrollmentId) {
        enrollmentId.value = payment.enrollmentId
        if (import.meta.client) {
          sessionStorage.setItem("dm_enrollment_id", payment.enrollmentId)
        }
      }
    })
    pollStatus(orderId.value)
  } else if (status.value === 'success') {
    setTimeout(() => {
      navigateTo('/dashboard?just_paid=true')
    }, 5000)
  }
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
  }
})

const packageNames = {
  starter: 'Starter Package (5 sessions)',
  standard: 'Standard Package (10 sessions)',
  pro: 'Pro Package (15 sessions)'
}

const nextSteps = computed(() => ({
  success: [
    'Check your email for payment confirmation and access details',
    'Confirm your full name and ID via WhatsApp for KTP verification',
    'Check the dashboard for your free trial session access',
    'Book your first paid session within 7 days to activate your package',
    'Receive your training schedule and instructor contact'
  ],
  failed: [
    'Payment was not processed',
    'Verify your payment amount and account details',
    'Try again with a different payment method',
    'Contact us via WhatsApp for assistance'
  ]
}))
</script>

<template>
  <div class="min-h-screen py-12 px-4 bg-gradient-to-b from-muted/20 to-background">
    <div class="max-w-md mx-auto">
      <!-- Success State -->
      <div v-if="isSuccess" class="space-y-6 text-center">
        <!-- Icon -->
        <div class="flex justify-center">
          <div class="w-24 h-24 rounded-full bg-green-500/10 flex items-center justify-center animate-bounce">
            <UIcon name="i-lucide-check-circle-2" class="size-12 text-green-500" />
          </div>
        </div>

        <!-- Message -->
        <div>
          <h1 class="text-3xl font-bold">{{ t('auth.paymentSuccessful') }}</h1>
          <p class="text-muted mt-2">
            {{ t('auth.paymentConfirmed', { plan: packageNames[plan as keyof typeof packageNames] || 'Standard Package' }) }}
          </p>
        </div>

        <!-- Confirmation Details -->
        <UCard class="text-left">
          <template #header>
            <h2 class="font-semibold">{{ t('auth.confirmationDetails') }}</h2>
          </template>

          <div class="space-y-4">
            <div>
              <p class="text-xs text-muted uppercase tracking-wide mb-1">{{ t('auth.confirmationEmail') }}</p>
              <p class="font-medium break-all">{{ email }}</p>
            </div>

            <div class="pt-4 border-t">
              <p class="text-xs text-muted uppercase tracking-wide mb-2">{{ t('auth.whatHappensNext') }}</p>
              <ol class="space-y-2">
                <li v-for="(step, index) in nextSteps.success" :key="index" class="flex gap-3 text-sm">
                  <span class="font-semibold text-primary shrink-0">{{ index + 1 }}.</span>
                  <span>{{ step }}</span>
                </li>
              </ol>
            </div>
          </div>

          <template #footer>
            <div class="text-sm text-muted">
              <p>{{ t('auth.checkEmailInstructions') }}</p>
            </div>
          </template>
        </UCard>

        <!-- Free Trial Banner -->
        <UCard class="bg-amber-500/10 border-amber-500/20">
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-gift" class="size-5 text-warning" />
              <h3 class="font-semibold text-warning">{{ t('freeTrial.title') }}</h3>
            </div>
          </template>

          <p class="text-sm text-warning">
            {{ t('freeTrial.heroDesc') }}
          </p>

          <template #footer>
            <NuxtLink to="/dashboard/free-trial">
              <UButton 
                :label="t('auth.viewFreeTrialDetails')"
                color="warning"
                variant="soft"
                block
                icon="i-lucide-arrow-right"
              />
            </NuxtLink>
          </template>
        </UCard>

        <!-- Action Buttons -->
        <div class="space-y-3 pt-4">
          <NuxtLink to="/dashboard?just_paid=true">
            <UButton 
              :label="t('auth.goToDashboard')"
              icon="i-lucide-arrow-right"
              block
              color="warning"
            />
          </NuxtLink>
          <p class="text-xs text-muted text-center py-6">
            {{ t('auth.redirectingIn') }}
          </p>
        </div>

        <!-- Support -->
        <UAlert icon="i-lucide-info" color="primary" variant="subtle">
          <template #description>
            {{ t('dashboard.needHelp') }}
          </template>
        </UAlert>
      </div>

      <!-- Failed State -->
      <div v-else-if="isFailed" class="space-y-6 text-center">
        <!-- Icon -->
        <div class="flex justify-center">
          <div class="w-24 h-24 rounded-full bg-red-500/10 flex items-center justify-center">
            <UIcon name="i-lucide-x-circle" class="size-12 text-red-500" />
          </div>
        </div>

        <!-- Message -->
        <div>
          <h1 class="text-3xl font-bold">{{ t('auth.paymentFailed') }}</h1>
          <p class="text-muted mt-2">
            Unfortunately, your payment could not be processed.
          </p>
        </div>

        <!-- Error Details -->
        <UCard class="text-left bg-red-500/5 border-red-500/20">
          <template #header>
            <h2 class="font-semibold text-red-900">{{ t('auth.whatWentWrong') }}</h2>
          </template>

          <ul class="space-y-2">
            <li v-for="(step, index) in nextSteps.failed" :key="index" class="flex gap-3 text-sm text-red-900">
              <UIcon name="i-lucide-alert-circle" class="size-4 shrink-0 mt-0.5" />
              <span>{{ step }}</span>
            </li>
          </ul>
        </UCard>

        <!-- Action Buttons -->
        <div class="space-y-3 pt-4">
          <NuxtLink :to="enrollmentId ? `/auth/payment-method?enrollment=${enrollmentId}` : '/auth/select-plan'">
            <UButton 
              :label="t('auth.tryAgain')"
              icon="i-lucide-rotate-cw"
              block
            />
          </NuxtLink>
          <NuxtLink to="/packages">
            <UButton 
              :label="t('auth.backToPackages')"
              color="neutral"
              variant="outline"
              icon="i-lucide-arrow-left"
              block
            />
          </NuxtLink>
        </div>

        <!-- Support -->
        <UAlert icon="i-lucide-help-circle" color="primary" variant="subtle">
          <template #description>
            {{ t('dashboard.needHelp') }}
          </template>
        </UAlert>
      </div>

      <!-- Pending State -->
      <div v-else class="space-y-6 text-center">
        <!-- Loading Animation -->
        <div class="flex justify-center">
          <div class="w-24 h-24 rounded-full bg-primary/10 flex items-center justify-center animate-spin">
            <UIcon name="i-lucide-loader-circle" class="size-12 text-primary" />
          </div>
        </div>

        <!-- Message -->
        <div>
          <h1 class="text-3xl font-bold">{{ t('auth.processingPayment') }}</h1>
          <p class="text-muted mt-2">
            Please wait while we confirm your payment. Order ID: <span class="font-mono text-sm">{{ orderId }}</span>
          </p>
        </div>

        <UCard class="space-y-3">
          <p class="text-sm text-muted">
            {{ t('auth.dontCloseWindow') }}
          </p>

          <div class="space-y-2 pt-2">
            <UButton 
              label="Check Status / Refresh"
              icon="i-lucide-rotate-cw"
              block
              color="primary"
              :loading="isCheckingStatus"
              @click="manualCheckStatus"
            />

            <UButton 
              label="Simulate Payment Success"
              icon="i-lucide-zap"
              block
              color="warning"
              variant="soft"
              :loading="isSimulating"
              @click="simulatePayment"
            />

            <NuxtLink to="/dashboard">
              <UButton 
                label="Go to Dashboard"
                icon="i-lucide-arrow-right"
                block
                color="neutral"
                variant="outline"
                class="mt-2"
              />
            </NuxtLink>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>
