<script setup lang="ts">
import { ref, onMounted } from 'vue'

definePageMeta({
  layout: 'blank'
})

const route = useRoute()
const router = useRouter()

const status = ref(route.query.status as string || 'pending')
const email = ref(route.query.email as string || '')
const plan = ref(route.query.plan as string || 'standard')

// Simulate waiting for confirmation
onMounted(() => {
  // After 3 seconds, auto-redirect to dashboard
  if (status.value === 'success') {
    setTimeout(() => {
      navigateTo('/dashboard')
    }, 5000)
  }
})

const isSuccess = computed(() => status.value === 'success')
const isFailed = computed(() => status.value === 'failed')

const packageNames = {
  starter: 'Starter Package (5 sessions)',
  standard: 'Standard Package (10 sessions)',
  pro: 'Pro Package (15 sessions)'
}

const nextSteps = {
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
}
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
          <h1 class="text-3xl font-bold">Payment Successful!</h1>
          <p class="text-muted mt-2">
            Your registration for {{ packageNames[plan as keyof typeof packageNames] || 'Standard Package' }} is confirmed.
          </p>
        </div>

        <!-- Confirmation Details -->
        <UCard class="text-left">
          <template #header>
            <h2 class="font-semibold">Confirmation Details</h2>
          </template>

          <div class="space-y-4">
            <div>
              <p class="text-xs text-muted uppercase tracking-wide mb-1">Confirmation Email</p>
              <p class="font-medium break-all">{{ email }}</p>
            </div>

            <div class="pt-4 border-t">
              <p class="text-xs text-muted uppercase tracking-wide mb-2">What Happens Next:</p>
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
              <p>Check your email for detailed instructions</p>
            </div>
          </template>
        </UCard>

        <!-- Free Trial Banner -->
        <UCard class="bg-amber-500/10 border-amber-500/20">
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-gift" class="size-5 text-warning" />
              <h3 class="font-semibold text-warning">Free Trial Session</h3>
            </div>
          </template>

          <p class="text-sm text-warning">
            You have access to a <span class="font-bold">free 15-minute trial session</span>. This is a one-time complimentary session available for registered users to experience our services at no cost. Use it wisely!
          </p>

          <template #footer>
            <NuxtLink to="/dashboard/free-trial">
              <UButton 
                label="View Free Trial Details"
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
          <NuxtLink to="/dashboard">
            <UButton 
              label="Go to Dashboard"
              icon="i-lucide-arrow-right"
              block
            />
          </NuxtLink>
          <p class="text-xs text-muted text-center">
            Redirecting in 5 seconds...
          </p>
        </div>

        <!-- Support -->
        <UAlert icon="i-lucide-info" color="primary" variant="subtle">
          <template #description>
            Need help? Contact us on WhatsApp at <span class="font-medium">+62 812-3456-7890</span>
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
          <h1 class="text-3xl font-bold">Payment Failed</h1>
          <p class="text-muted mt-2">
            Unfortunately, your payment could not be processed.
          </p>
        </div>

        <!-- Error Details -->
        <UCard class="text-left bg-red-500/5 border-red-500/20">
          <template #header>
            <h2 class="font-semibold text-red-900">What Went Wrong?</h2>
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
          <NuxtLink :to="`/auth/payment-method?plan=${plan}`">
            <UButton 
              label="Try Again"
              icon="i-lucide-rotate-cw"
              block
            />
          </NuxtLink>
          <NuxtLink to="/packages">
            <UButton 
              label="Back to Packages"
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
            Still having issues? Contact our support team via WhatsApp
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
          <h1 class="text-3xl font-bold">Processing Payment</h1>
          <p class="text-muted mt-2">
            Please wait while we confirm your payment. This may take a few moments.
          </p>
        </div>

        <UCard>
          <p class="text-sm text-muted">
            Do not close this window. You will be automatically redirected once payment is confirmed.
          </p>
        </UCard>
      </div>
    </div>
  </div>
</template>
