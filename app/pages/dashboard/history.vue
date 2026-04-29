<script setup lang="ts">
import { computed, ref } from 'vue'

definePageMeta({ layout: 'dashboard' })

const searchQuery = ref('')
const statusFilter = ref('all')

// Mock training history
const trainingHistory = ref([
  {
    id: 1,
    sessionNumber: 4,
    date: 'Mar 25, 2026',
    time: '09:30 AM - 10:30 AM',
    duration: '60 min',
    car: 'BYD Atto 1',
    instructor: 'Pak Ahmad',
    topic: 'Highway Driving Basics',
    status: 'completed',
    notes: 'Good progress on lane changes. Need more practice with highway merging.',
    rating: 4
  },
  {
    id: 2,
    sessionNumber: 3,
    date: 'Mar 22, 2026',
    time: '11:00 AM - 12:00 PM',
    duration: '60 min',
    car: 'BYD Atto 1',
    instructor: 'Bu Sari',
    topic: 'Parking & Maneuvering',
    status: 'completed',
    notes: 'Excellent parallel parking skills. Reverse parking needs improvement.',
    rating: 5
  },
  {
    id: 3,
    sessionNumber: 2,
    date: 'Mar 18, 2026',
    time: '14:30 PM - 15:30 PM',
    duration: '60 min',
    car: 'BYD Atto 1',
    instructor: 'Pak Budi',
    topic: 'City Driving',
    status: 'completed',
    notes: 'Handled traffic lights and intersections well. Good awareness of pedestrians.',
    rating: 4
  },
  {
    id: 4,
    sessionNumber: 1,
    date: 'Mar 15, 2026',
    time: '08:00 AM - 09:00 AM',
    duration: '60 min',
    car: 'BYD Atto 1',
    instructor: 'Pak Ahmad',
    topic: 'Introduction & Basic Controls',
    status: 'completed',
    notes: 'Great first session! Quickly adapted to EV controls and one-pedal driving.',
    rating: 5
  }
])

const filteredHistory = computed(() => {
  return trainingHistory.value.filter(session => {
    const matchesSearch = session.topic.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         session.instructor.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = statusFilter.value === 'all' || session.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

const totalHours = computed(() => {
  return trainingHistory.value.filter(s => s.status === 'completed').length
})
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Training History">
        <template #right>
          <UButton icon="i-lucide-download" label="Export" color="neutral" variant="outline" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <UInput 
            v-model="searchQuery"
            placeholder="Search sessions..."
            icon="i-lucide-search"
            class="w-64"
          />
        </template>
        <template #right>
          <USelect 
            v-model="statusFilter"
            :items="[
              { label: 'All Sessions', value: 'all' },
              { label: 'Completed', value: 'completed' },
              { label: 'Cancelled', value: 'cancelled' }
            ]"
            class="w-40"
          />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Summary Stats -->
        <div class="grid md:grid-cols-3 gap-4">
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-primary/10">
                <UIcon name="i-lucide-check-circle" class="size-6 text-primary" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ trainingHistory.filter(s => s.status === 'completed').length }}</p>
                <p class="text-md text-muted">Sessions Completed</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-blue-500/10">
                <UIcon name="i-lucide-clock" class="size-6 text-blue-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ totalHours }} hrs</p>
                <p class="text-md text-muted">Total Training Time</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-amber-500/10">
                <UIcon name="i-lucide-star" class="size-6 text-amber-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">4.5</p>
                <p class="text-md text-muted">Average Rating</p>
              </div>
            </div>
          </UCard>
        </div>

        <!-- History List -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">Session History</h2>
          </template>

          <div class="space-y-4">
            <div 
              v-for="session in filteredHistory" 
              :key="session.id"
              class="p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
            >
              <div class="flex flex-col lg:flex-row lg:items-start justify-between gap-4">
                <div class="flex items-start gap-4">
                  <div class="p-3 rounded-xl bg-primary/10">
                    <UIcon name="i-lucide-car" class="size-6 text-primary" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2 flex-wrap">
                      <h3 class="font-semibold">Session #{{ session.sessionNumber }}: {{ session.topic }}</h3>
                      <UBadge 
                        :label="session.status === 'completed' ? 'Completed' : 'Cancelled'"
                        :color="session.status === 'completed' ? 'success' : 'error'"
                        variant="subtle"
                        size="md"
                      />
                    </div>
                    <p class="text-md text-muted mt-1">{{ session.date }} | {{ session.time }}</p>
                    <div class="flex items-center gap-4 mt-2 text-md text-muted">
                      <span class="flex items-center gap-1">
                        <UIcon name="i-lucide-car" class="size-4" />
                        {{ session.car }}
                      </span>
                      <span class="flex items-center gap-1">
                        <UIcon name="i-lucide-user" class="size-4" />
                        {{ session.instructor }}
                      </span>
                      <span class="flex items-center gap-1">
                        <UIcon name="i-lucide-clock" class="size-4" />
                        {{ session.duration }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="flex items-center gap-2 lg:flex-col lg:items-end">
                  <div class="flex gap-0.5">
                    <UIcon 
                      v-for="i in 5" 
                      :key="i"
                      name="i-lucide-star" 
                      :class="i <= session.rating ? 'text-amber-500 fill-amber-500' : 'text-muted'" 
                      class="size-4"
                    />
                  </div>
                  <UButton label="View Details" variant="ghost" size="md" icon="i-lucide-eye" />
                </div>
              </div>

              <div v-if="session.notes" class="mt-4 p-3 rounded-lg bg-muted/50">
                <p class="text-md text-muted mb-1">Instructor Notes:</p>
                <p class="text-md">{{ session.notes }}</p>
              </div>
            </div>

            <UEmpty 
              v-if="filteredHistory.length === 0"
              icon="i-lucide-search-x" 
              title="No Sessions Found" 
              description="Try adjusting your search or filter criteria."
            />
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
