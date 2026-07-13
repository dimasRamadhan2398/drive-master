<script setup lang="ts">
import { ref, computed } from 'vue'

const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

const freeTrial = ref({
  used: false,
  expiresAt: null as Date | null,
  sessionDate: null as Date | null,
  duration: 15, // minutes
  status: 'available' // available, scheduled, completed, expired
})

const trialSession = ref({
  id: 'trial-001',
  date: null as Date | null,
  time: '14:00',
  duration: 15,
  instructor: 'Pak Ahmad',
  vehicle: 'BYD Atto 1',
  location: 'Alam Sutera, Jakarta'
})

const timeSlots = [
  '09:00', '09:30', '10:00', '10:30', '11:00', '11:30',
  '14:00', '14:30', '15:00', '15:30', '16:00', '16:30',
  '17:00', '17:30', '18:00'
]

const selectedDate = ref<string>('')
const selectedTime = ref<string>('')
const loading = ref(false)

const isExpired = computed(() => {
  if (!freeTrial.value.expiresAt) return false
  return new Date() > freeTrial.value.expiresAt
})

const daysUntilExpiry = computed(() => {
  if (!freeTrial.value.expiresAt) return null
  const daysLeft = Math.ceil((freeTrial.value.expiresAt.getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24))
  return daysLeft > 0 ? daysLeft : 0
})

async function scheduleFreeTrial() {
  if (!selectedDate.value || !selectedTime.value) return

  loading.value = true

  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500))

    trialSession.value.date = new Date(selectedDate.value)
    trialSession.value.time = selectedTime.value
    freeTrial.value.status = 'scheduled'
    freeTrial.value.sessionDate = new Date(selectedDate.value)

    console.log('[v0] Free trial scheduled:', trialSession.value)
  } finally {
    loading.value = false
  }
}

async function cancelFreeTrial() {
  loading.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    selectedDate.value = ''
    selectedTime.value = ''
    freeTrial.value.status = 'available'
    freeTrial.value.sessionDate = null
  } finally {
    loading.value = false
  }
}

// Set minimum date to tomorrow
const minDate = computed(() => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return tomorrow.toISOString().split('T')[0]
})

// Set maximum date to 30 days from now (expiry period)
const maxDate = computed(() => {
  const maxDay = new Date()
  maxDay.setDate(maxDay.getDate() + 30)
  return maxDay.toISOString().split('T')[0]
})

// Simulate checking if user has free trial (from user state/session)
const userStatus = ref<'free' | 'paid'>('paid') // Would come from actual user data
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('freeTrial.title')">
        <template #right>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 max-w-4xl mx-auto space-y-6">
        <!-- Header Card -->
        <UCard class="bg-gradient-to-r from-amber-500/10 to-orange-500/10 border-amber-500/20">
          <div class="flex items-start gap-4">
            <div class="p-3 rounded-lg bg-amber-500/10 shrink-0">
              <UIcon name="i-lucide-gift" class="size-6 text-amber-600" />
            </div>
            <div class="flex-1">
              <h1 class="text-2xl font-bold">{{ t('freeTrial.heroTitle') }}</h1>
              <p class="text-muted mt-1">
                {{ t('freeTrial.heroDesc') }}
              </p>
            </div>
          </div>
        </UCard>

        <!-- Trial Info Cards -->
        <div class="grid sm:grid-cols-3 gap-4">
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-primary/10">
                <UIcon name="i-lucide-clock" class="size-5 text-primary" />
              </div>
              <div>
                <p class="text-xs text-muted uppercase tracking-wide">{{ t('freeTrial.duration') }}</p>
                <p class="font-bold">15 {{ t('common.minutes') }}</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-green-500/10">
                <UIcon name="i-lucide-check" class="size-5 text-green-600" />
              </div>
              <div>
                <p class="text-xs text-muted uppercase tracking-wide">{{ t('freeTrial.cost') }}</p>
                <p class="font-bold">{{ t('freeTrial.free') }}</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-3">
              <div 
                class="p-2 rounded-lg"
                :class="isExpired ? 'bg-red-500/10' : 'bg-blue-500/10'"
              >
                <UIcon 
                  :name="isExpired ? 'i-lucide-x' : 'i-lucide-calendar'"
                  :class="isExpired ? 'text-red-600' : 'text-blue-600'"
                  class="size-5"
                />
              </div>
              <div>
                <p class="text-xs text-muted uppercase tracking-wide">
                  {{ isExpired ? t('freeTrial.status') : t('freeTrial.expiresIn') }}
                </p>
                <p v-if="isExpired" class="font-bold text-red-600">{{ t('freeTrial.expired') }}</p>
                <p v-else class="font-bold">{{ daysUntilExpiry }} {{ t('freeTrial.days') }}</p>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Status Badge -->
        <div class="flex items-center gap-2">
          <UBadge 
            v-if="freeTrial.status === 'available' && !isExpired"
            :label="t('freeTrial.readyToSchedule')"
            color="success"
            variant="subtle"
          />
          <UBadge 
            v-else-if="freeTrial.status === 'scheduled'"
            :label="t('freeTrial.scheduled')"
            color="primary"
            variant="subtle"
          />
          <UBadge 
            v-else-if="freeTrial.status === 'completed'"
            :label="t('freeTrial.completed')"
            color="success"
            variant="subtle"
          />
          <UBadge 
            v-else
            :label="t('freeTrial.expired')"
            color="error"
            variant="subtle"
          />
        </div>

        <!-- Main Content -->
        <div v-if="!isExpired && freeTrial.status !== 'completed'" class="grid lg:grid-cols-3 gap-6">
          <!-- Booking Form -->
          <div class="lg:col-span-2">
            <UCard v-if="freeTrial.status === 'available'">
              <template #header>
                <h2 class="font-semibold">{{ t('freeTrial.scheduleTitle') }}</h2>
              </template>

              <div class="space-y-6">
                <!-- Date Selection -->
                <div>
                  <label class="block text-sm font-medium mb-3">{{ t('freeTrial.selectDate') }}</label>
                  <UInput 
                    v-model="selectedDate"
                    type="date"
                    :min="minDate"
                    :max="maxDate"
                    size="lg"
                    icon="i-lucide-calendar"
                  />
                  <p class="text-xs text-muted mt-2">{{ t('freeTrial.availableFrom') }}</p>
                </div>

                <!-- Time Selection -->
                <div v-if="selectedDate">
                  <label class="block text-sm font-medium mb-3">{{ t('freeTrial.selectTime') }}</label>
                  <div class="grid grid-cols-3 gap-2">
                    <button
                      v-for="time in timeSlots"
                      :key="time"
                      class="p-2 rounded-lg border-2 transition-all text-sm font-medium"
                      :class="selectedTime === time 
                        ? 'border-primary bg-primary/10 text-primary' 
                        : 'border-default hover:border-primary/30'"
                      @click="selectedTime = time"
                    >
                      {{ time }}
                    </button>
                  </div>
                </div>

                <!-- Session Details Preview -->
                <div v-if="selectedDate && selectedTime" class="p-4 bg-muted/30 rounded-lg space-y-3">
                  <h3 class="font-medium">{{ t('freeTrial.sessionDetails') }}</h3>
                  <div class="grid sm:grid-cols-2 gap-4 text-sm">
                    <div>
                      <p class="text-muted">{{ t('dashboard.date') }} & {{ t('dashboard.time') }}</p>
                      <p class="font-medium">
                        {{ new Date(selectedDate).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
                        at {{ selectedTime }}
                      </p>
                    </div>
                    <div>
                      <p class="text-muted">{{ t('freeTrial.duration') }}</p>
                      <p class="font-medium">15 {{ t('common.minutes') }}</p>
                    </div>
                    <div>
                      <p class="text-muted">{{ t('dashboard.vehicle') }}</p>
                      <p class="font-medium">BYD Atto 1</p>
                    </div>
                    <div>
                      <p class="text-muted">{{ t('dashboard.instructor') }}</p>
                      <p class="font-medium">Mr. Ahmad</p>
                    </div>
                  </div>
                </div>

                <!-- Action Button -->
                <div class="pt-4 border-t">
                  <UButton 
                    :label="t('freeTrial.confirmBooking')"
                    icon="i-lucide-check"
                    :loading="loading"
                    :disabled="!selectedDate || !selectedTime || loading"
                    @click="scheduleFreeTrial"
                    block
                  />
                </div>
              </div>
            </UCard>

            <!-- Scheduled Session View -->
            <UCard v-else-if="freeTrial.status === 'scheduled'">
              <template #header>
                <div class="flex items-center justify-between">
                  <h2 class="font-semibold">{{ t('freeTrial.sessionScheduled') }}</h2>
                  <UBadge :label="t('freeTrial.confirmed')" color="success" />
                </div>
              </template>

              <div class="space-y-6">
                <!-- Scheduled Details -->
                <div class="space-y-4">
                  <div class="flex items-center gap-4 p-4 bg-primary/5 rounded-lg">
                    <UIcon name="i-lucide-check-circle-2" class="size-6 text-primary shrink-0" />
                    <div>
                      <p class="font-medium">{{ t('freeTrial.sessionConfirmed') }}</p>
                      <p class="text-sm text-muted">{{ t('freeTrial.confirmationSent') }}</p>
                    </div>
                  </div>

                  <div class="grid gap-4">
                    <div class="flex items-center gap-3 p-3 rounded-lg bg-muted/30">
                      <UIcon name="i-lucide-calendar" class="size-5 text-primary" />
                      <div class="flex-1">
                        <p class="text-sm text-muted">{{ t('dashboard.date') }} & {{ t('dashboard.time') }}</p>
                        <p class="font-medium">
                          {{ new Date(trialSession.date || '').toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
                          at {{ trialSession.time }}
                        </p>
                      </div>
                    </div>

                    <div class="flex items-center gap-3 p-3 rounded-lg bg-muted/30">
                      <UIcon name="i-lucide-clock" class="size-5 text-primary" />
                      <div class="flex-1">
                        <p class="text-sm text-muted">{{ t('freeTrial.duration') }}</p>
                        <p class="font-medium">{{ trialSession.duration }} {{ t('common.minutes') }}</p>
                      </div>
                    </div>

                    <div class="flex items-center gap-3 p-3 rounded-lg bg-muted/30">
                      <UIcon name="i-lucide-car" class="size-5 text-primary" />
                      <div class="flex-1">
                        <p class="text-sm text-muted">{{ t('dashboard.vehicle') }}</p>
                        <p class="font-medium">{{ trialSession.vehicle }}</p>
                      </div>
                    </div>

                    <div class="flex items-center gap-3 p-3 rounded-lg bg-muted/30">
                      <UIcon name="i-lucide-user" class="size-5 text-primary" />
                      <div class="flex-1">
                        <p class="text-sm text-muted">{{ t('dashboard.instructor') }}</p>
                        <p class="font-medium">{{ trialSession.instructor }}</p>
                      </div>
                    </div>

                    <div class="flex items-center gap-3 p-3 rounded-lg bg-muted/30">
                      <UIcon name="i-lucide-map-pin" class="size-5 text-primary" />
                      <div class="flex-1">
                        <p class="text-sm text-muted">{{ t('dashboard.pickupLocation') }}</p>
                        <p class="font-medium">{{ trialSession.location }}</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- What to Expect -->
                <div class="pt-4 border-t space-y-3">
                  <h3 class="font-medium">{{ t('freeTrial.whatToExpect') }}</h3>
                  <ul class="space-y-2 text-sm">
                    <li class="flex gap-2">
                      <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0 mt-0.5" />
                      <span>{{ t('freeTrial.meetInstructor') }}</span>
                    </li>
                    <li class="flex gap-2">
                      <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0 mt-0.5" />
                      <span>{{ t('freeTrial.safetyBriefing') }}</span>
                    </li>
                    <li class="flex gap-2">
                      <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0 mt-0.5" />
                      <span>{{ t('freeTrial.drivingExperience') }}</span>
                    </li>
                    <li class="flex gap-2">
                      <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0 mt-0.5" />
                      <span>{{ t('freeTrial.discussion') }}</span>
                    </li>
                  </ul>
                </div>

                <!-- Reschedule -->
                <div class="pt-4 border-t flex gap-3">
                  <UButton 
                    :label="t('freeTrial.cancelBooking')"
                    icon="i-lucide-trash"
                    color="error"
                    variant="outline"
                    :loading="loading"
                    @click="cancelFreeTrial"
                    block
                  />
                </div>
              </div>
            </UCard>
          </div>

          <!-- Sidebar Info -->
          <div class="space-y-6">
            <!-- Requirements -->
            <UCard>
              <template #header>
                <h3 class="font-semibold">{{ t('freeTrial.requirements') }}</h3>
              </template>

              <ul class="space-y-3 text-sm">
                <li class="flex gap-2">
                  <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0" />
                  <span>{{ t('freeTrial.validLicense') }}</span>
                </li>
                <li class="flex gap-2">
                  <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0" />
                  <span>{{ t('freeTrial.atLeast18') }}</span>
                </li>
                <li class="flex gap-2">
                  <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0" />
                  <span>{{ t('freeTrial.closedShoes') }}</span>
                </li>
                <li class="flex gap-2">
                  <UIcon name="i-lucide-check" class="size-4 text-primary shrink-0" />
                  <span>{{ t('freeTrial.beEarly') }}</span>
                </li>
              </ul>
            </UCard>

            <!-- After Trial -->
            <UCard>
              <template #header>
                <h3 class="font-semibold">{{ t('freeTrial.afterTrial') }}</h3>
              </template>

              <p class="text-sm text-muted mb-4">
                {{ t('freeTrial.afterTrialDesc') }}
              </p>

              <ul class="space-y-2 text-sm">
                <li class="flex gap-2">
                  <UIcon name="i-lucide-arrow-right" class="size-4 text-primary shrink-0" />
                  <span>{{ t('freeTrial.bookPaid') }}</span>
                </li>
                <li class="flex gap-2">
                  <UIcon name="i-lucide-arrow-right" class="size-4 text-primary shrink-0" />
                  <span>{{ t('freeTrial.choosePkg') }}</span>
                </li>
                <li class="flex gap-2">
                  <UIcon name="i-lucide-arrow-right" class="size-4 text-primary shrink-0" />
                  <span>{{ t('freeTrial.getCert') }}</span>
                </li>
              </ul>

              <UButton 
                :label="t('freeTrial.viewPackages')"
                to="/packages"
                variant="outline"
                block
                class="mt-4"
              />
            </UCard>
          </div>
        </div>

        <!-- Expired State -->
        <UAlert v-if="isExpired" icon="i-lucide-alert-circle" color="error">
          <template #title>{{ t('freeTrial.periodExpired') }}</template>
          <template #description>
            {{ t('freeTrial.periodExpiredDesc') }}
          </template>
        </UAlert>

        <!-- Completed State -->
        <UAlert v-if="freeTrial.status === 'completed'" icon="i-lucide-check-circle-2" color="success">
          <template #title>{{ t('freeTrial.trialCompleted') }}</template>
          <template #description>
            {{ t('freeTrial.trialCompletedDesc') }}
          </template>
        </UAlert>
      </div>
    </template>
  </UDashboardPanel>
</template>
