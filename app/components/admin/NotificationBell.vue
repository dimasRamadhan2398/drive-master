<script setup lang="ts">
const { activeAlerts } = useSmartAlerts()
const { updateSlotStatus } = useSchedules()
const toast = useToast()

function startSession(id: string, student: string) {
  updateSlotStatus(id, 'in-progress')
  toast.add({ 
    title: 'Session Started', 
    description: `Kursus untuk ${student} resmi dimulai.`, 
    color: 'success', 
    icon: 'i-lucide-play' 
  })
}
</script>

<template>
  <UDropdownMenu 
    v-if="activeAlerts.length > 0"
    :items="[]" 
    :ui="{ content: 'w-80 p-0' }"
  >
    <!-- Bell Icon with Ping Badge -->
    <UButton
      icon="i-lucide-bell"
      color="neutral"
      variant="ghost"
      class="relative"
    >
      <span class="absolute top-1.5 right-1.5 flex size-2">
        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
        <span class="relative inline-flex rounded-full size-2 bg-red-500"></span>
      </span>
    </UButton>

    <template #content>
      <div class="p-4 border-b border-default flex items-center justify-between">
        <h3 class="font-bold text-sm">Upcoming Sessions</h3>
        <UBadge :label="activeAlerts.length.toString()" color="error" variant="subtle" size="sm" />
      </div>
      
      <div class="max-h-60 overflow-y-auto">
        <div 
          v-for="alert in activeAlerts" 
          :key="alert.id"
          class="p-3 border-b border-default last:border-0 hover:bg-muted/50 transition-colors"
        >
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-bold text-primary">{{ alert.time }}</span>
            <span class="text-[10px] text-muted font-medium uppercase tracking-wider">Starts Soon</span>
          </div>
          <p class="text-sm font-medium mb-3">{{ alert.student }}</p>
          <UButton 
            label="Start Session" 
            icon="i-lucide-play" 
            size="xs" 
            block 
            color="warning"
            @click="startSession(alert.id, alert.student || 'Siswa')"
          />
        </div>
      </div>
    </template>
  </UDropdownMenu>

  <!-- Empty state bell if no alerts -->
  <UButton
    v-else
    icon="i-lucide-bell"
    color="neutral"
    variant="ghost"
    class="opacity-50"
  />
</template>
