<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { scheduleService } from '~/services/scheduleService'
import type { SessionResponse } from '~/services/scheduleService'

const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

const authStore = useAuthStore()
const schedulesStore = useSchedulesStore()

const userSessionsList = ref<SessionResponse[]>([])
const isLoadingSessions = ref(false)

const fetchDashboardData = async (userId: string) => {
  await schedulesStore.fetchSchedules({ studentId: userId, limit: 100 })
  try {
    isLoadingSessions.value = true
    const response = await scheduleService.fetchSessions({ studentId: userId, limit: 10 })
    userSessionsList.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch user sessions:', err)
  } finally {
    isLoadingSessions.value = false
  }
}

// ── Fetch / refresh on mount ─────────────────────────────────────────────────
onMounted(async () => {
  await authStore.fetchMemberProfile()
  if (authStore.userId) {
    await fetchDashboardData(authStore.userId)
  }
})

// Watch for userId change in case it loads asynchronously
watch(() => authStore.userId, async (newUserId) => {
  if (newUserId && userSessionsList.value.length === 0) {
    await fetchDashboardData(newUserId)
  }
})

// ── Derived profile data ──────────────────────────────────────────────────────
const profile = computed(() => authStore.memberProfile)
const isLoading = computed(() => authStore.isLoadingProfile)

const firstName = computed(() =>
  authStore.currentUser?.firstName || authStore.memberIdentityFullname?.split(' ')[0] || 'Student'
)

// Aggregate across all active entitlements
const entitlements = computed(() => authStore.memberEntitlements ?? [])

const totalSessions = computed(() =>
  entitlements.value.reduce((sum, e) => sum + (e.totalSessions ?? 0), 0)
)

const completedSessions = computed(() =>
  entitlements.value.reduce((sum, e) => sum + (e.usedSessions ?? 0), 0)
)

const remainingSessions = computed(() =>
  entitlements.value.reduce((sum, e) => sum + (e.remaining ?? 0), 0)
)

const isCourseCompleted = computed(() => {
  return totalSessions.value > 0 && remainingSessions.value === 0
})

const progressPercent = computed(() => {
  const total = totalSessions.value
  if (!total) return 0
  return Math.round((completedSessions.value / total) * 100)
})

// Certificate: eligible when all sessions in at least one entitlement are used
const certificateStatus = computed(() => {
  if (!entitlements.value.length) return t('common.pending')
  if (isCourseCompleted.value) return t('common.completed')
  const hasCompleted = entitlements.value.some(e => e.status === 'completed')
  if (hasCompleted) return t('common.ready')
  if (progressPercent.value >= 100) return t('common.processing')
  return t('common.pending')
})

// Active package name (most recent active entitlement)
const packageName = computed(() => {
  const active = entitlements.value.find(e => e.status === 'active')
  return active?.packageName || t('dashboard.noPackage')
})

// Filter booked and in-progress slots for student
const studentBookedSessions = computed(() => {
  return schedulesStore.slots.filter(
    slot => (slot.status === 'booked' || slot.status === 'in-progress')
  )
})

// Sort by date and time in ascending order to get the closest upcoming session
const sortedUpcomingSessions = computed(() => {
  return [...studentBookedSessions.value].sort((a, b) => {
    const dateTimeA = new Date(`${a.date}T${a.time}`)
    const dateTimeB = new Date(`${b.date}T${b.time}`)
    return dateTimeA.getTime() - dateTimeB.getTime()
  })
})

const nextSession = computed(() => {
  if (sortedUpcomingSessions.value.length === 0) return null
  return sortedUpcomingSessions.value[0]
})

const nextInstructorDetails = computed(() => {
  if (!nextSession.value) return null
  return allInstructors.find(inst => inst.name === nextSession.value.instructor)
})

const showInstructorModal = ref(false)
const showDetails = ref(false)

const sessionDetails = computed(() => ({
  pickup: nextSession.value?.notes || 'Main Lobby of Green Bay Apartments, Pluit'
}))

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const [year, month, day] = dateStr.split('-').map(Number)
  if (!year || !month || !day) return dateStr
  const date = new Date(year, month - 1, day)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const recentActivity = computed(() => {
  // If we have live sessions fetched, map them
  if (userSessionsList.value.length > 0) {
    const completedSessions = userSessionsList.value.filter(s => s.status === 'completed')
    return userSessionsList.value
      .slice(0, 3) // show top 3
      .map((session) => {
        const isCompleted = session.status === 'completed'
        const isInProgress = session.status === 'in-progress'
        
        let title = t('schedule.bookingSuccess')
        if (isCompleted) {
          const completedIndex = completedSessions.findIndex(s => s.id === session.id)
          const sessionNumber = completedIndex !== -1 ? completedSessions.length - completedIndex : 1
          title = `${t('history.session')} #${sessionNumber}`
        }

        let description = `Training Session scheduled`
        if (isInProgress) {
          description = 'Driving session is in progress'
        } else if (isCompleted) {
          description = session.notes && session.notes.trim() ? session.notes : 'Driving session completed'
        }

        return {
          id: session.id,
          type: 'session',
          title,
          description,
          date: formatDate(session.date),
          status: isCompleted ? 'completed' : isInProgress ? 'warning' : 'info'
        }
      })
  }

  const userSlots = schedulesStore.slots.filter(
    slot => slot.status === 'completed' || slot.status === 'booked' || slot.status === 'in-progress'
  )
  
  if (userSlots.length === 0) {
    return []
  }

  const completedSlots = userSlots.filter(s => s.status === 'completed')

  return [...userSlots]
    .sort((a, b) => {
      const dateTimeA = new Date(`${a.date}T${a.time}`)
      const dateTimeB = new Date(`${b.date}T${b.time}`)
      return dateTimeB.getTime() - dateTimeA.getTime()
    })
    .slice(0, 3)
    .map((slot) => {
      const isCompleted = slot.status === 'completed'
      const isInProgress = slot.status === 'in-progress'
      
      let title = t('schedule.bookingSuccess')
      if (isCompleted) {
        const completedIndex = completedSlots.findIndex(s => s.id === slot.id)
        const sessionNumber = completedIndex !== -1 ? completedSlots.length - completedIndex : 1
        title = `${t('history.session')} #${sessionNumber}`
      }

      return {
        id: slot.id,
        type: 'session',
        title,
        description: isInProgress 
          ? 'Driving session is in progress' 
          : isCompleted 
            ? 'Driving session completed' 
            : `Training Session with ${slot.instructor} scheduled`,
        date: formatDate(slot.date),
        status: isCompleted ? 'completed' : isInProgress ? 'warning' : 'info'
      }
    })
})

const quickActions = computed(() => {
  const firstAction = isCourseCompleted.value
    ? { label: t('dashboard.viewCertificate'), icon: 'i-lucide-award', to: '/dashboard/certificate', color: 'warning' as const }
    : { label: t('dashboard.bookSession'), icon: 'i-lucide-calendar-plus', to: '/dashboard/schedule', color: 'warning' as const }

  return [
    firstAction,
    { label: t('dashboard.viewHistory'), icon: 'i-lucide-history', to: '/dashboard/history', color: 'neutral' as const },
    { label: t('dashboard.getSupport'), icon: 'i-simple-icons-whatsapp', to: 'https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi', external: true, color: 'primary' as const }
  ]
})
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
              <h1 class="text-2xl font-bold">{{ t('dashboard.welcome', { name: firstName }) }}</h1>
              <p v-if="isCourseCompleted" class="text-muted mt-1">{{ t('dashboard.courseCompleted') }}</p>
              <p v-else class="text-muted mt-1">{{ t('dashboard.sessionsRemaining', { count: remainingSessions, package: packageName }) }}</p>
            </div>
            <NuxtLink :to="isCourseCompleted ? '/dashboard/certificate' : '/dashboard/schedule'">
              <UButton :label="isCourseCompleted ? t('dashboard.viewCertificate') : t('dashboard.bookNextSession')" :icon="isCourseCompleted ? 'i-lucide-award' : 'i-lucide-calendar-plus'" color="warning" />
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
                <p class="text-2xl font-bold">
                  <USkeleton v-if="isLoading" class="h-8 w-10" />
                  <span v-else>{{ completedSessions }}</span>
                </p>
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
                <p class="text-2xl font-bold">
                  <USkeleton v-if="isLoading" class="h-8 w-10" />
                  <span v-else>{{ remainingSessions }}</span>
                </p>
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
                <p class="text-2xl font-bold">
                  <USkeleton v-if="isLoading" class="h-8 w-12" />
                  <span v-else>{{ progressPercent }}%</span>
                </p>
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
                <p class="text-2xl font-bold">
                  <USkeleton v-if="isLoading" class="h-8 w-20" />
                  <span v-else>{{ certificateStatus }}</span>
                </p>
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
                <UBadge v-if="nextSession" :label="nextSession.status === 'in-progress' ? t('dashboard.inProgress') : t('dashboard.confirmed')" :color="nextSession.status === 'in-progress' ? 'warning' : 'success'" variant="subtle" />
              </div>
            </template>

            <div v-if="nextSession" class="flex flex-col sm:flex-row gap-6">
              <div class="flex-1 space-y-4">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-calendar" class="size-5 text-warning" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">{{ t('dashboard.date') }}</p>
                    <p class="font-medium">{{ nextSession.date }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-clock" class="size-5 text-warning" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">{{ t('dashboard.time') }}</p>
                    <p class="font-medium">{{ nextSession.time }}</p>
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
                    <p class="font-medium">{{ nextSession.car || 'BYD Atto 3' }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-3 group">
                  <div class="p-2 rounded-lg bg-muted">
                    <UIcon name="i-lucide-user" class="size-5 text-warning" />
                  </div>
                  <div class="flex-1">
                    <p class="text-sm text-muted">{{ t('dashboard.instructor') }}</p>
                    <p class="font-medium">{{ nextSession.instructor || 'Instructor' }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="flex flex-col items-center justify-center py-8 text-center space-y-3">
              <div class="p-4 rounded-full bg-muted">
                <UIcon :name="isCourseCompleted ? 'i-lucide-award' : 'i-lucide-calendar-x'" :class="isCourseCompleted ? 'text-warning' : 'text-muted-foreground'" class="size-8" />
              </div>
              <div v-if="isCourseCompleted">
                <p class="font-semibold text-base">{{ t('dashboard.courseCompletedTitle') }}</p>
                <p class="text-sm text-muted mt-1">{{ t('dashboard.courseCompletedDesc') }}</p>
              </div>
              <div v-else>
                <p class="font-semibold text-base">{{ t('schedule.noUpcoming') || 'Tidak ada sesi kelas terdekat' }}</p>
                <p class="text-sm text-muted mt-1">{{ t('schedule.noUpcomingDescription') || 'Silakan pilih tanggal dan pesan kelas mengemudi Anda berikutnya.' }}</p>
              </div>
              <NuxtLink :to="isCourseCompleted ? '/dashboard/certificate' : '/dashboard/schedule'">
                <UButton :label="isCourseCompleted ? t('dashboard.viewCertificate') : t('dashboard.bookNextSession')" :icon="isCourseCompleted ? 'i-lucide-award' : 'i-lucide-calendar-plus'" color="warning" size="sm" class="mt-2" />
              </NuxtLink>
            </div>

            <template #footer v-if="nextSession">
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
                      :stroke-dashoffset="352 - (352 * progressPercent) / 100"
                      stroke-linecap="round"
                    />
                  </svg>
                  <span class="absolute text-2xl font-bold">{{ progressPercent }}%</span>
                </div>
              </div>

              <div class="space-y-2">
                <div class="flex justify-between text-sm">
                  <span class="text-muted">{{ t('common.completed') }}</span>
                  <span class="font-medium">{{ completedSessions }}/{{ totalSessions }}</span>
                </div>
                <UProgress :model-value="progressPercent" />
              </div>
            </div>

            <template #footer>
              <p class="text-sm text-muted text-center">
                <span v-if="isCourseCompleted">{{ t('dashboard.courseFinishedFooter') }}</span>
                <span v-else>{{ t('dashboard.moreSessions', { count: remainingSessions }) }}</span>
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
              <div v-if="isLoadingSessions" class="space-y-3 p-3">
                <div class="flex items-center gap-4">
                  <USkeleton class="h-10 w-10 rounded-lg" />
                  <div class="flex-1 space-y-2">
                    <USkeleton class="h-4 w-[40%]" />
                    <USkeleton class="h-3 w-[70%]" />
                  </div>
                </div>
                <div class="flex items-center gap-4">
                  <USkeleton class="h-10 w-10 rounded-lg" />
                  <div class="flex-1 space-y-2">
                    <USkeleton class="h-4 w-[30%]" />
                    <USkeleton class="h-3 w-[60%]" />
                  </div>
                </div>
              </div>
              <div 
                v-else-if="recentActivity.length > 0"
                v-for="activity in recentActivity" 
                :key="activity.id"
                class="flex items-start gap-4 p-3 rounded-lg hover:bg-muted/50 transition-colors"
              >
                <div 
                  class="p-2 rounded-lg"
                  :class="activity.status === 'completed' ? 'bg-green-500/10' : activity.status === 'warning' ? 'bg-yellow-500/10' : 'bg-blue-500/10'"
                >
                  <UIcon 
                    :name="activity.type === 'session' ? 'i-lucide-car' : 'i-lucide-calendar-check'"
                    :class="activity.status === 'completed' ? 'text-green-500' : activity.status === 'warning' ? 'text-yellow-500' : 'text-blue-500'"
                    class="size-5"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-medium">{{ activity.title }}</p>
                  <p class="text-sm text-muted truncate">{{ activity.description }}</p>
                </div>
                <span class="text-xs text-muted whitespace-nowrap">{{ activity.date }}</span>
              </div>
              <div v-else class="flex flex-col items-center justify-center text-center py-8 px-4 rounded-xl border border-dashed border-gray-200 dark:border-gray-800 bg-gray-50/50 dark:bg-gray-900/30 transition-all duration-300">
                <div class="relative flex items-center justify-center w-12 h-12 rounded-full bg-warning-500/10 text-warning-500 mb-3">
                  <UIcon name="i-lucide-history" class="w-6 h-6" />
                  <div class="absolute inset-0 rounded-full bg-warning-500/5 animate-pulse"></div>
                </div>
                <h3 class="font-semibold text-sm text-gray-900 dark:text-white mb-1">
                  {{ t('history.noHistory') || 'No training history yet' }}
                </h3>
                <p class="text-xs text-muted max-w-[280px] mb-4 animate-fade-in">
                  {{ t('dashboard.noActivityDesc') }}
                </p>
                <NuxtLink to="/dashboard/schedule">
                  <UButton
                    size="xs"
                    color="warning"
                    icon="i-lucide-calendar-plus"
                    class="font-medium shadow-sm hover:scale-[1.02] transition-transform duration-200"
                  >
                    {{ t('dashboard.bookNextSession') }}
                  </UButton>
                </NuxtLink>
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
        <template #body v-if="nextSession">
          <div class="space-y-4 py-2">
            <div class="grid grid-cols-2 gap-4">
              <div class="col-span-2">
                <p class="text-xs font-bold text-muted uppercase mb-1">{{ t('dashboard.courseMaterial') }}</p>
                <p class="font-bold">{{ t('home.material.highway') }}</p>
              </div>
              <div>
                <p class="text-xs font-bold text-muted uppercase mb-1">{{ t('dashboard.time') }}</p>
                <p class="text-sm font-medium text-primary">{{ nextSession.time }}</p>
              </div>
              <div>
                <p class="text-xs font-bold text-muted uppercase mb-1">{{ t('dashboard.instructor') }}</p>
                <p class="text-sm font-medium">{{ nextSession.instructor }}</p>
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
