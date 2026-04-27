<script setup lang="ts">
import { onMounted, ref } from 'vue'

definePageMeta({ layout: 'dashboard' })

// Onboarding modal state - would check if user has completed profile
const showOnboarding = ref(false)
const userHasCompletedOnboarding = ref(false)

onMounted(() => {
  // Check if user needs to complete onboarding
  // In real app, this would check from backend/session
  if (!userHasCompletedOnboarding.value) {
    setTimeout(() => {
      showOnboarding.value = true
    }, 500)
  }
})

function handleOnboardingComplete(data: any) {
  console.log('[v0] User completed onboarding:', data)
  userHasCompletedOnboarding.value = true
  showOnboarding.value = false
}

// Mock data
const userData = {
  name: 'John Doe',
  package: 'Standard Package',
  totalSessions: 10,
  completedSessions: 4,
  remainingSessions: 6,
  progress: 40,
  nextSession: {
    date: 'Tomorrow',
    time: '09:30 AM',
    car: 'Tesla Model 3',
    instructor: 'Pak Ahmad'
  }
}

const recentActivity = [
  { 
    id: 1, 
    type: 'session', 
    title: 'Training Session #4', 
    description: 'Highway driving basics completed',
    date: 'Mar 25, 2026',
    status: 'completed'
  },
  { 
    id: 2, 
    type: 'session', 
    title: 'Training Session #3', 
    description: 'Parking and maneuvering',
    date: 'Mar 22, 2026',
    status: 'completed'
  },
  { 
    id: 3, 
    type: 'booking', 
    title: 'Session Booked', 
    description: 'Training Session #5 scheduled',
    date: 'Mar 20, 2026',
    status: 'info'
  }
]

const quickActions = [
  { label: 'Book Session', icon: 'i-lucide-calendar-plus', to: '/dashboard/schedule', color: 'primary' as const },
  { label: 'View History', icon: 'i-lucide-history', to: '/dashboard/history', color: 'neutral' as const },
  { label: 'Get Support', icon: 'i-simple-icons-whatsapp', to: 'https://wa.me/6281234567890', external: true, color: 'neutral' as const }
]
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Dashboard">
        <template #right>
          <UButton icon="i-lucide-bell" color="neutral" variant="ghost" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <!-- Onboarding Modal -->
      <OnboardingModal 
        v-model="showOnboarding"
        @complete="handleOnboardingComplete"
      />

      <div class="p-6 space-y-6">
        <!-- Welcome Banner -->
        <UCard class="bg-primary/5 border-primary/20">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <h1 class="text-2xl font-bold">Welcome back, {{ userData.name.split(' ')[0] }}!</h1>
              <p class="text-muted mt-1">You have {{ userData.remainingSessions }} sessions remaining in your {{ userData.package }}.</p>
            </div>
            <NuxtLink to="/dashboard/schedule">
              <UButton label="Book Next Session" icon="i-lucide-calendar-plus" />
            </NuxtLink>
          </div>
        </UCard>

        <!-- Stats Cards -->
        <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-primary/10">
                <UIcon name="i-lucide-book-check" class="size-6 text-primary" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ userData.completedSessions }}</p>
                <p class="text-sm text-muted">Sessions Completed</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-amber-500/10">
                <UIcon name="i-lucide-clock" class="size-6 text-amber-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ userData.remainingSessions }}</p>
                <p class="text-sm text-muted">Sessions Remaining</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-blue-500/10">
                <UIcon name="i-lucide-target" class="size-6 text-blue-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ userData.progress }}%</p>
                <p class="text-sm text-muted">Course Progress</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-green-500/10">
                <UIcon name="i-lucide-award" class="size-6 text-green-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">Pending</p>
                <p class="text-sm text-muted">Certificate Status</p>
              </div>
            </div>
          </UCard>
        </div>

        <div class="grid lg:grid-cols-3 gap-6">
          <!-- Next Session Card -->
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">Next Session</h2>
                <UBadge label="Confirmed" color="success" variant="subtle" />
              </div>
            </template>

            <div class="flex flex-col sm:flex-row gap-6">
              <div class="flex-1 space-y-4">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-calendar" class="size-5 text-primary" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">Date</p>
                    <p class="font-medium">{{ userData.nextSession.date }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-clock" class="size-5 text-primary" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">Time</p>
                    <p class="font-medium">{{ userData.nextSession.time }}</p>
                  </div>
                </div>
              </div>

              <div class="flex-1 space-y-4">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-car" class="size-5 text-primary" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">Vehicle</p>
                    <p class="font-medium">{{ userData.nextSession.car }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-user" class="size-5 text-primary" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">Instructor</p>
                    <p class="font-medium">{{ userData.nextSession.instructor }}</p>
                  </div>
                </div>
              </div>
            </div>

            <template #footer>
              <div class="flex gap-3">
                <UButton label="View Details" variant="outline" color="neutral" />
                <UButton label="Reschedule" variant="ghost" color="neutral" icon="i-lucide-calendar-days" />
              </div>
            </template>
          </UCard>

          <!-- Progress Card -->
          <UCard>
            <template #header>
              <h2 class="font-semibold">Course Progress</h2>
            </template>

            <div class="space-y-4">
              <div class="text-center">
                <div class="relative inline-flex items-center justify-center w-32 h-32">
                  <svg class="w-full h-full transform -rotate-90">
                    <circle
                      cx="64"
                      cy="64"
                      r="56"
                      stroke="currentColor"
                      stroke-width="12"
                      fill="none"
                      class="text-muted/20"
                    />
                    <circle
                      cx="64"
                      cy="64"
                      r="56"
                      stroke="currentColor"
                      stroke-width="12"
                      fill="none"
                      class="text-primary"
                      :stroke-dasharray="352"
                      :stroke-dashoffset="352 - (352 * userData.progress) / 100"
                      stroke-linecap="round"
                    />
                  </svg>
                  <span class="absolute text-2xl font-bold">{{ userData.progress }}%</span>
                </div>
              </div>

              <div class="space-y-2">
                <div class="flex justify-between text-sm">
                  <span class="text-muted">Completed</span>
                  <span class="font-medium">{{ userData.completedSessions }}/{{ userData.totalSessions }}</span>
                </div>
                <UProgress :value="userData.progress" />
              </div>
            </div>

            <template #footer>
              <p class="text-sm text-muted text-center">
                {{ userData.remainingSessions }} more sessions to complete your course
              </p>
            </template>
          </UCard>
        </div>

        <!-- Recent Activity & Quick Actions -->
        <div class="grid lg:grid-cols-3 gap-6">
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">Recent Activity</h2>
                <NuxtLink to="/dashboard/history">
                  <UButton label="View All" variant="ghost" size="sm" trailingIcon="i-lucide-arrow-right" />
                </NuxtLink>
              </div>
            </template>

            <div class="space-y-4">
              <div 
                v-for="activity in recentActivity" 
                :key="activity.id"
                class="flex items-start gap-4 p-3 rounded-lg hover:bg-muted/50 transition-colors"
              >
                <div 
                  class="p-2 rounded-lg"
                  :class="activity.status === 'completed' ? 'bg-green-500/10' : 'bg-blue-500/10'"
                >
                  <UIcon 
                    :name="activity.type === 'session' ? 'i-lucide-car' : 'i-lucide-calendar-check'"
                    :class="activity.status === 'completed' ? 'text-green-500' : 'text-blue-500'"
                    class="size-5"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-medium">{{ activity.title }}</p>
                  <p class="text-sm text-muted truncate">{{ activity.description }}</p>
                </div>
                <span class="text-xs text-muted whitespace-nowrap">{{ activity.date }}</span>
              </div>
            </div>
          </UCard>

          <UCard>
            <template #header>
              <h2 class="font-semibold">Quick Actions</h2>
            </template>

            <div class="space-y-3">
              <NuxtLink 
                v-for="action in quickActions" 
                :key="action.label"
                :to="action.to"
                :target="action.external ? '_blank' : undefined"
              >
                <UButton 
                  :label="action.label" 
                  :icon="action.icon" 
                  :color="action.color"
                  :variant="action.color === 'neutral' ? 'outline' : undefined"
                  block
                />
              </NuxtLink>
            </div>

            <template #footer>
              <UAlert icon="i-lucide-info" variant="subtle">
                <template #description>
                  Need help? Contact our support team via WhatsApp.
                </template>
              </UAlert>
            </template>
          </UCard>
        </div>
      </div>
    </template>
  </UDashboardPanel>
</template>
