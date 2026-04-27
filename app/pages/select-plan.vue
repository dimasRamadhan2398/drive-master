<script setup lang="ts">
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'blank'
})

const selectedPlan = ref<'six_package' | 'eight_package' | 'ten_package' | 'twelve_package'>('eight_package')
const loading = ref(false)

// Check if promo is active (dynamic)
const isPromoActive = computed(() => {
  const now = new Date()
  const promoEnd = new Date('2026-05-20T23:59:59') // Can be dynamic from backend
  return now < promoEnd
})

const promoDuration = computed(() => {
  const now = new Date()
  const promoEnd = new Date('2026-05-20T23:59:59')
  const daysLeft = Math.ceil((promoEnd.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return daysLeft > 0 ? daysLeft : 0
})

const plans = [
  {
    id: 'six_package',
    name: '6x',
    originalPrice: 1750000,
    // promoPrice: 1600000,
    sessions: 6,
    features: [
      'Free Trial',
      '6 training sessions'
    ],
    highlight: false
  },
  {
    id: 'eight_package',
    name: '8x',
    originalPrice: 1950000,
    // promoPrice: 1750000,
    sessions: 8,
    features: [
      'Free Trial',
      '8 training sessions'
    ],
    highlight: true
  },
  {
    id: 'ten_package',
    name: '10x',
    originalPrice: 2250000,
    // promoPrice: 2100000,
    sessions: 10,
    features: [
      'Free Trial',
      '10 training sessions'
    ],
    highlight: false
  },
  {
    id: 'twelve_package',
    name: '12x',
    originalPrice: 2650000,
    // promoPrice: 2500000,
    sessions: 12,
    features: [
      'Free Trial',
      '12 training sessions'
    ],
    highlight: false
  }
]

const currentPlan = computed(() => plans.find(p => p.id === selectedPlan.value))
// const discount = computed(() => {
//   if (!currentPlan.value || !isPromoActive.value) return 0
//   return Math.round(((currentPlan.value.originalPrice - currentPlan.value.promoPrice) / currentPlan.value.originalPrice) * 100)
// })

async function selectPlan() {
  if (!currentPlan.value) return

  loading.value = true

  try {
    // Simulate validation
    await new Promise(resolve => setTimeout(resolve, 1000))
    console.log('Plan selected:', selectedPlan.value)

    // Redirect to payment method
    navigateTo(`/payment-method?plan=${selectedPlan.value}`)
  } finally {
    loading.value = false
  }
}

const freeTrialInfo = [
  {
    icon: 'i-lucide-gift',
    title: 'One-Time Offer',
    description: 'Available once per account for new registered users'
  },
  {
    icon: 'i-lucide-clock',
    title: '15 Minutes',
    description: 'Session duration for all user types (free and paid)'
  },
  {
    icon: 'i-lucide-check-circle',
    title: 'No Commitment',
    description: 'Experience our service risk-free before payment'
  },
  {
    icon: 'i-lucide-zap',
    title: 'Full Experience',
    description: 'Same professional instructor and Premium vehicles'
  }
]
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4 bg-muted/20">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-12">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon name="i-lucide-package" class="size-8 text-warning" />
          <span class="text-xl font-bold">Select Your Package</span>
        </div>
        <h1 class="text-3xl md:text-4xl font-bold">Choose Your Training Package</h1>
        <p class="text-muted mt-3 max-w-2xl mx-auto">
          Select the package that best fits your learning goals. All packages include access to our premium vehicles and experienced instructors.
        </p>
      </div>

      <!-- Promo Banner
      <Transition
        enter-active-class="transition duration-300"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
      >
        <UCard 
          v-if="isPromoActive"
          class="mb-8 bg-gradient-to-r from-amber-500/10 to-orange-500/10 border-amber-500/20"
        >
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-amber-500/10 shrink-0">
              <UIcon name="i-lucide-zap" class="size-6 text-amber-600" />
            </div>
            <div class="flex-1">
              <h3 class="font-bold text-amber-900">Limited Time Offer!</h3>
              <p class="text-sm text-amber-800">
                Get up to {{ discount }}% discount on all packages. Offer ends in {{ promoDuration }} days.
              </p>
            </div>
            <div class="text-right shrink-0">
              <p class="text-2xl font-bold text-amber-600">{{ discount }}%</p>
              <p class="text-xs text-amber-700">OFF</p>
            </div>
          </div>
        </UCard>
      </Transition>
    -->

      <!-- Plan Cards -->
      <div class="grid md:grid-cols-4 gap-6 mb-12">
        <div 
          v-for="plan in plans" 
          :key="plan.id"
          class="relative"
        >
          <input 
            :id="`plan-${plan.id}`"
            v-model="selectedPlan"
            type="radio" 
            :value="plan.id"
            class="sr-only"
          />
          <label 
            :for="`plan-${plan.id}`"
            class="block h-full cursor-pointer"
          >
            <UCard
              :class="[
                'h-full flex flex-col transition-all',
                selectedPlan === plan.id 
                  ? 'ring-2 ring-warning shadow-xl' 
                  : 'hover:shadow-lg'
              ]"
            >
              <div v-if="plan.highlight" class="absolute -top-3 left-1/2 -translate-x-1/2 z-10">
                <UBadge label="Most Popular" color="warning" />
              </div>

              <template #header>
                <div class="text-center">
                  <h3 class="text-2xl font-bold">{{ plan.name }}</h3>
                  <!-- <p class="text-muted text-sm mt-2">{{ plan.duration }}</p> -->
                </div>
              </template>

              <div class="flex-1 space-y-6">
                <!-- Pricing -->
                <div class="text-center">
                  <div>
                    <p class="text-4xl font-bold text-warning">
                      Rp {{ plan.originalPrice.toLocaleString('id-ID') }}
                    </p>
                  </div>
                  <!-- <div v-if="isPromoActive && plan.promoPrice < plan.originalPrice" class="space-y-2">
                    <p class="text-sm text-muted line-through">
                      Rp {{ plan.originalPrice.toLocaleString('id-ID') }}
                    </p>
                    <p class="text-4xl font-bold text-warning">
                      Rp {{ plan.promoPrice.toLocaleString('id-ID') }}
                    </p>
                    <p class="text-xs text-green-600 font-medium">
                      Save Rp {{ (plan.originalPrice - plan.promoPrice).toLocaleString('id-ID') }}
                    </p>
                  </div>
                  <div v-else>
                    <p class="text-4xl font-bold">
                      Rp {{ plan.originalPrice.toLocaleString('id-ID') }}
                    </p>
                  </div> -->
                </div>

                <!-- Sessions Badge -->
                <div class="text-center">
                  <UBadge :label="`${plan.sessions} Sessions`" color="warning" variant="subtle" />
                </div>

                <!-- Features -->
                <ul class="space-y-3">
                  <li 
                    v-for="feature in plan.features" 
                    :key="feature"
                    class="flex items-start gap-3"
                  >
                    <UIcon name="i-lucide-check" class="size-5 text-warning shrink-0 mt-0.5" />
                    <span class="text-sm">{{ feature }}</span>
                  </li>
                </ul>
              </div>

              <!-- Selection Indicator -->
              <div v-if="selectedPlan === plan.id" class="pt-4 border-t">
                <div class="flex items-center justify-between">
                  <span class="text-sm font-medium">Selected</span>
                  <div class="p-1 rounded-full bg-warning">
                    <UIcon name="i-lucide-check" class="size-4 text-white" />
                  </div>
                </div>
              </div>
            </UCard>
          </label>
        </div>
      </div>

      <!-- Free Trial Info Section -->
      <div class="mb-12">
        <div class="text-center mb-8">
          <h2 class="text-2xl font-bold">Free Trial Session Included</h2>
          <p class="text-muted mt-2">
            Available for both free and paid registered users - one time offer
          </p>
        </div>

        <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <UCard v-for="item in freeTrialInfo" :key="item.title">
            <div class="text-center space-y-3">
              <div class="flex justify-center">
                <div class="p-3 rounded-lg bg-warning/10">
                  <UIcon :name="item.icon" class="size-6 text-warning" />
                </div>
              </div>
              <div>
                <p class="font-semibold text-sm">{{ item.title }}</p>
                <p class="text-xs text-muted mt-1">{{ item.description }}</p>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- CTA Section -->
      <div class="max-w-2xl mx-auto">
        <UCard class="bg-warning/5 border-warning/20">
          <div class="text-center space-y-4">
            <h3 class="text-xl font-bold">Ready to Get Started?</h3>
            <p class="text-muted">
              Selected package: <span class="font-bold">{{ currentPlan?.name }} - Rp {{ currentPlan?.originalPrice?.toLocaleString('id-ID') }}</span>
            </p>
            <!-- <p class="text-muted">
              Selected package: <span class="font-bold">{{ currentPlan?.name }} - Rp {{ (isPromoActive ? currentPlan?.promoPrice : currentPlan?.originalPrice)?.toLocaleString('id-ID') }}</span>
            </p> -->

            <div class="flex gap-3 pt-2">
              <UButton 
                label="Proceed to Payment"
                icon="i-lucide-arrow-right"
                color="warning"
                :loading="loading"
                @click="selectPlan"
                block
              />
            </div>

            <p class="text-xs text-muted">
              Payment secured by Midtrans. You can cancel anytime before first session.
            </p>
          </div>
        </UCard>
      </div>

      <!-- Additional Info
      <div class="mt-8 text-center">
        <p class="text-sm text-muted">
          Have questions?
          <NuxtLink to="/packages" class="text-warning hover:underline">
            View detailed package comparison
          </NuxtLink>
        </p>
      </div>
    -->
    </div>
  </div>
</template>
