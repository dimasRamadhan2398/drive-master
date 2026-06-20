<script setup lang="ts">
import { computed } from 'vue'

const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

// Mock data
const userData = {
  name: 'John Doe',
  package: '8x Sessions Package',
  totalSessions: 10,
  completedSessions: 4,
  remainingSessions: 6,
  progress: 40,
  nextSession: {
    date: 'Tomorrow',
    time: '09:30 AM',
    car: 'BYD Atto 1',
    instructor: 'Mr. Ahmad'
  }
}

// FITUR BARU: Data detail instruktor dan state untuk modal
const allInstructors = [
  {
    name: 'Mr. Ahmad',
    phone: '081234567001',
    bnsp: 'BNSP-101-2023',
    sim: 'SIM A',
    photoUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&auto=format&fit=crop&q=80',
    experience: 10,
    bio: 'Expert in smooth vehicle control and traffic safety management with extensive knowledge in electric vehicle systems.',
    status: 'active',
    role: 'Senior Instructor',
    specialization: 'Defensive Driving & EV Specialist'
  },
  {
    name: 'Ms. Sari',
    phone: '081234567002',
    bnsp: 'BNSP-102-2022',
    sim: 'SIM A',
    photoUrl: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=400&auto=format&fit=crop&q=80',
    experience: 7,
    bio: 'Specializes in helping beginners master city driving and complex parking maneuvers with patience and precision.',
    status: 'active',
    role: 'Professional Instructor',
    specialization: 'Urban Driving & Parking Expert'
  }
]

const nextInstructorDetails = computed(() => {
  return allInstructors.find(inst => inst.name === userData.nextSession.instructor)
})

const showInstructorModal = ref(false)

const showDetails = ref(false)
const sessionDetails = {
  pickup: 'Main Lobby of Green Bay Apartments, Pluit'
}

const recentActivity = computed(() => [
  { 
    id: 1, 
    type: 'session', 
    title: `${t('history.session')} #4`,
    description: 'Highway driving basics completed',
    date: 'Mar 25, 2026',
    status: 'completed'
  },
  { 
    id: 2, 
    type: 'session', 
    title: `${t('history.session')} #3`,
    description: 'Parking and maneuvering',
    date: 'Mar 22, 2026',
    status: 'completed'
  },
  { 
    id: 3, 
    type: 'booking', 
    title: t('schedule.bookingSuccess'),
    description: 'Training Session #5 scheduled',
    date: 'Mar 20, 2026',
    status: 'info'
  }
])

const quickActions = computed(() => [
  { label: t('dashboard.bookSession'), icon: 'i-lucide-calendar-plus', to: '/dashboard/schedule', color: 'warning' as const },
  { label: t('dashboard.viewHistory'), icon: 'i-lucide-history', to: '/dashboard/history', color: 'neutral' as const },
  { label: t('dashboard.getSupport'), icon: 'i-simple-icons-whatsapp', to: 'https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi', external: true, color: 'primary' as const }
])
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('common.dashboard')">
        <template #right>
          <UButton icon="i-lucide-bell" color="neutral" variant="ghost" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Welcome Banner -->
        <UCard class="bg-warning/5 border-warning/20">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <h1 class="text-2xl font-bold">{{ t('dashboard.welcome', { name: userData.name.split(' ')[0] }) }}</h1>
              <p class="text-muted mt-1">{{ t('dashboard.sessionsRemaining', { count: userData.remainingSessions, package: userData.package }) }}</p>
            </div>
            <NuxtLink to="/dashboard/schedule">
              <UButton :label="t('dashboard.bookNextSession')" icon="i-lucide-calendar-plus" color="warning" />
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
                <p class="text-sm text-muted">{{ t('dashboard.sessionsCompleted') }}</p>
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
                <p class="text-sm text-muted">{{ t('dashboard.sessionsRemainingCount') }}</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-info/10">
                <UIcon name="i-lucide-target" class="size-6 text-info" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ userData.progress }}%</p>
                <p class="text-sm text-muted">{{ t('dashboard.courseProgress') }}</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-neutral/10">
                <UIcon name="i-lucide-award" class="size-6 text-neutral" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ t('common.pending') }}</p>
                <p class="text-sm text-muted">{{ t('dashboard.certificateStatus') }}</p>
              </div>
            </div>
          </UCard>
        </div>

        <div class="grid lg:grid-cols-3 gap-6">
          <!-- Next Session Card -->
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">{{ t('dashboard.nextSession') }}</h2>
                <UBadge :label="t('dashboard.confirmed')" color="success" variant="subtle" />
              </div>
            </template>

            <div class="flex flex-col sm:flex-row gap-6">
              <div class="flex-1 space-y-4">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-calendar" class="size-5 text-warning" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">{{ t('dashboard.date') }}</p>
                    <p class="font-medium">{{ userData.nextSession.date }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-clock" class="size-5 text-warning" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">{{ t('dashboard.time') }}</p>
                    <p class="font-medium">{{ userData.nextSession.time }}</p>
                  </div>
                </div>
              </div>

              <div class="flex-1 space-y-4">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-car" class="size-5 text-warning" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">{{ t('dashboard.vehicle') }}</p>
                    <p class="font-medium">{{ userData.nextSession.car }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-3 group">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-user" class="size-5 text-warning" />
                  </div>
                  <div class="flex-1">
                    <p class="text-sm text-muted">{{ t('dashboard.instructor') }}</p>
                    <p class="font-medium">{{ userData.nextSession.instructor }}</p>
                  </div>
                  
                  </div>
                </div>
            </div>
            <template #footer>
              <div class="flex gap-3">
                <UButton :label="t('dashboard.viewDetails')" variant="outline" color="neutral" @click="showDetails = true" />
                <UButton :label="t('dashboard.reschedule')" to="/dashboard/schedule" variant="ghost" color="neutral" icon="i-lucide-calendar-days" />
              </div>
            </template>
          </UCard>

          <!-- Progress Card -->
          <UCard>
            <template #header>
              <h2 class="font-semibold">{{ t('dashboard.courseProgressTitle') }}</h2>
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
                  <span class="text-muted">{{ t('common.completed') }}</span>
                  <span class="font-medium">{{ userData.completedSessions }}/{{ userData.totalSessions }}</span>
                </div>
                <UProgress :value="userData.progress" />
              </div>
            </div>

            <template #footer>
              <p class="text-sm text-muted text-center">
                {{ t('dashboard.moreSessions', { count: userData.remainingSessions }) }}
              </p>
            </template>
          </UCard>
        </div>

        <!-- Recent Activity & Quick Actions -->
        <div class="grid lg:grid-cols-3 gap-6">
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">{{ t('dashboard.recentActivity') }}</h2>
                <NuxtLink to="/dashboard/history">
                  <UButton :label="t('dashboard.viewAll')" variant="ghost" color="warning" size="sm" trailingIcon="i-lucide-arrow-right" />
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
              <h2 class="font-semibold">{{ t('dashboard.quickActions') }}</h2>
            </template>

            <div class="grid grid-cols-1 space-y-3">
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
                  {{ t('dashboard.needHelp') }}
                </template>
              </UAlert>
            </template>
          </UCard>
        </div>
      </div>
      <!-- Centered Modal -->
      <UModal v-model:open="showDetails" :title="t('dashboard.sessionDetail')" :ui="{ content: 'm-auto sm:max-w-md' }">
        <template #body>
          <div class="space-y-4 py-2">
            <div class="grid grid-cols-2 gap-4">
              <div class="col-span-2">
                <p class="text-xs font-bold text-muted uppercase mb-1">{{ t('dashboard.courseMaterial') }}</p>
                <p class="font-bold">{{ t('home.material.highway') }}</p>
              </div>
              <div>
                <p class="text-xs font-bold text-muted uppercase mb-1">{{ t('dashboard.time') }}</p>
                <p class="text-sm font-medium text-primary">{{ userData.nextSession.time }}</p>
              </div>
              <div>
                <p class="text-xs font-bold text-muted uppercase mb-1">{{ t('dashboard.instructor') }}</p>
                <p class="text-sm font-medium">{{ userData.nextSession.instructor }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs font-bold text-muted uppercase mb-1">{{ t('dashboard.pickupLocation') }}</p>
                <p class="text-sm font-medium">{{ sessionDetails.pickup }}</p>
              </div>
            </div>
          </div>
        </template>
        <template #footer>
          <UButton :label="t('dashboard.close')" block color="neutral" variant="soft" @click="showDetails = false" />
        </template>
      </UModal>

      
    </template>
  </UDashboardPanel>

  
</template>
