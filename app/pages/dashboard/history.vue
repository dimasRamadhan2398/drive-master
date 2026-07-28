<script setup lang="ts">
const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

const authStore = useAuthStore()
const schedulesStore = useSchedulesStore()
const instructorsStore = useInstructorsStore()
const vehiclesStore = useVehiclesStore()

const searchQuery = ref('')
const statusFilter = ref('all')
const trainingHistory = ref<any[]>([])
const isLoading = ref(false)

const formatDate = (dateStr: string) => {
  if (!dateStr) return "";
  const [year, month, day] = dateStr.split("-").map(Number);
  if (!year || !month || !day) return dateStr;
  const date = new Date(year, month - 1, day);
  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
};

const getCarName = (s: any) => {
  if (s.carName) return s.carName;
  if (s.car) return s.car;
  if (s.carId) {
    const vehicle = vehiclesStore.getVehicleById(s.carId);
    if (vehicle) return `${vehicle.brand} ${vehicle.model}`;
  }
  if (s.scheduleId) {
    const slot = schedulesStore.slots.find((slot) => slot.id === String(s.scheduleId));
    if (slot && slot.car) return slot.car;
  }
  return "BYD Atto 1";
};

const fetchHistory = async () => {
  if (!authStore.userId) return;
  try {
    isLoading.value = true;
    const sessions = await schedulesStore.fetchUserSessions(authStore.userId);
    
    // Sort the sessions: most recent at the top (descending order)
    const sortedSessions = sessions.slice().sort((a, b) => {
      const dateA = new Date(`${a.date}T${a.time || '00:00'}:00`);
      const dateB = new Date(`${b.date}T${b.time || '00:00'}:00`);
      return dateB.getTime() - dateA.getTime();
    });

    // Map to template format
    trainingHistory.value = sortedSessions.map((s, index, arr) => {
      // Find instructor name
      const instructor = instructorsStore.getInstructorByUserId(s.instructorId);
      const instructorName = s.instructorName || (instructor ? instructor.name : "Instructor");
      
      return {
        id: s.id,
        sessionNumber: arr.length - index,
        date: formatDate(s.date),
        time: s.time,
        duration: `${s.duration} min`,
        car: getCarName(s),
        instructor: instructorName,
        topic: s.notes || "Training Session",
        status: s.status,
        notes: s.notes || "",
        rating: 5
      };
    });
  } catch (err) {
    console.error("Failed to fetch history:", err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
  if (instructorsStore.instructors.length === 0) {
    await instructorsStore.fetchInstructors();
  }
  if (vehiclesStore.vehicles.length === 0) {
    await vehiclesStore.fetchVehicles();
  }
  await fetchHistory();
});

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

const isModalOpen = ref(false)
const selectedSession = ref<any>(null)

const openDetails = (session: any) => {
  selectedSession.value = session
  isModalOpen.value = true
}
</script>


<template>
 <UDashboardPanel>
   <template #header>
     <UDashboardNavbar :title="t('history.title')">
       <template #right>
         <UColorModeButton />
       </template>
     </UDashboardNavbar>


     <UDashboardToolbar>
       <template #left>
         <UInput
           v-model="searchQuery"
           :placeholder="t('history.searchPlaceholder')"
           icon="i-lucide-search"
           class="w-64"
        />
       </template>
       <template #right>
         <USelect
           v-model="statusFilter"
           :items="[
             { label: t('history.allSessions'), value: 'all' },
             { label: t('history.completed'), value: 'completed' },
             { label: t('history.cancelled'), value: 'cancelled' }
           ]"
           class="w-40"
        />
       </template>
     </UDashboardToolbar>
   </template>


   <template #body>
     <div class="p-6 space-y-6">
       <!-- Summary Stats -->
       <div class="grid sm:grid-cols-3 gap-4">
         <UCard>
           <div class="flex items-center gap-4">
             <div class="p-3 rounded-xl bg-primary/10">
               <UIcon name="i-lucide-check-circle" class="size-6 text-primary" />
             </div>
             <div>
               <p class="text-2xl font-bold">{{ trainingHistory.filter(s => s.status === 'completed').length }}</p>
               <p class="text-sm text-muted">{{ t('dashboard.sessionsCompleted') }}</p>
             </div>
           </div>
         </UCard>


         <UCard>
           <div class="flex items-center gap-4">
             <div class="p-3 rounded-xl bg-blue-500/10">
               <UIcon name="i-lucide-clock" class="size-6 text-blue-500" />
             </div>
             <div>
               <p class="text-2xl font-bold">{{ totalHours }} {{ t('common.hours') }}</p>
               <p class="text-sm text-muted">{{ t('history.totalTrainingTime') }}</p>
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
               <p class="text-sm text-muted">{{ t('history.averageRating') }}</p>
             </div>
           </div>
         </UCard>
       </div>


       <!-- History List -->
       <UCard>
         <template #header>
           <h2 class="font-semibold">{{ t('history.sessionHistory') }}</h2>
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
                     <h3 class="font-semibold">{{ t('history.session') }} #{{ session.sessionNumber }}: {{ session.topic }}</h3>
                      <UBadge
                        :label="t('history.' + session.status) || session.status"
                        :color="session.status === 'completed' ? 'success' : (session.status === 'cancelled' ? 'error' : (session.status === 'in-progress' || session.status === 'in_progress' ? 'warning' : 'primary'))"
                        variant="subtle"
                        size="xs"
                      />
                   </div>
                   <p class="text-sm text-muted mt-1">{{ session.date }} | {{ session.time }}</p>
                   <div class="flex items-center gap-4 mt-2 text-sm text-muted">
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
                 <UButton
                   :label="t('history.viewDetails')"
                   variant="ghost"
                   size="xs"
                   icon="i-lucide-eye"
                   @click="openDetails(session)"
                />
               </div>
             </div>


             <div v-if="session.notes" class="mt-4 p-3 rounded-lg bg-muted/50">
               <p class="text-xs text-muted mb-1">{{ t('history.instructorNotes') }}</p>
               <p class="text-sm">{{ session.notes }}</p>
             </div>
           </div>


           <UEmpty
             v-if="filteredHistory.length === 0"
             icon="i-lucide-search-x"
             :title="t('history.noSessionsFound')"
             :description="t('history.adjustSearch')"
          />
         </div>
       </UCard>
     </div>
   </template>
 </UDashboardPanel>


 <ClientOnly>
   <UModal v-model:open="isModalOpen">
     <template #content>
       <UCard v-if="selectedSession" :ui="{ body: { padding: 'p-0' }, ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
         <template #header>
           <div class="flex items-center justify-between">
             <h3 class="text-base font-semibold leading-6">
               {{ t('history.session') }} Details #{{ selectedSession.sessionNumber }}
             </h3>
             <UButton color="neutral" variant="ghost" icon="i-lucide-x" @click="isModalOpen = false" />
           </div>
         </template>


         <div class="p-6 space-y-6">
           <div class="flex items-center gap-4">
             <div class="p-4 rounded-2xl bg-primary/10">
               <UIcon name="i-lucide-car" class="size-8 text-primary" />
             </div>
             <div>
               <h4 class="text-xl font-bold">{{ selectedSession.topic }}</h4>
               <div class="flex items-center gap-2 mt-1">
                  <UBadge
                    :label="t('history.' + selectedSession.status) || selectedSession.status"
                    :color="selectedSession.status === 'completed' ? 'success' : (selectedSession.status === 'cancelled' ? 'error' : (selectedSession.status === 'in-progress' || selectedSession.status === 'in_progress' ? 'warning' : 'primary'))"
                    variant="subtle"
                  />
                 <span class="text-sm text-muted">{{ selectedSession.date }}</span>
               </div>
             </div>
           </div>


           <div class="grid grid-cols-2 gap-6">
             <div class="space-y-1">
               <p class="text-xs text-muted uppercase font-bold tracking-wider">{{ t('dashboard.instructor') }}</p>
               <p class="font-medium flex items-center gap-2">
                 <UIcon name="i-lucide-user" class="size-4 text-primary" />
                 {{ selectedSession.instructor }}
               </p>
             </div>
             <div class="space-y-1">
               <p class="text-xs text-muted uppercase font-bold tracking-wider">{{ t('dashboard.vehicle') }}</p>
               <p class="font-medium flex items-center gap-2">
                 <UIcon name="i-lucide-car" class="size-4 text-primary" />
                 {{ selectedSession.car }}
               </p>
             </div>
             <div class="space-y-1">
               <p class="text-xs text-muted uppercase font-bold tracking-wider">{{ t('dashboard.time') }} & {{ t('packages.duration') }}</p>
               <p class="font-medium flex items-center gap-2">
                 <UIcon name="i-lucide-clock" class="size-4 text-primary" />
                 {{ selectedSession.time }} ({{ selectedSession.duration }})
               </p>
             </div>
             <div class="space-y-1">
               <p class="text-xs text-muted uppercase font-bold tracking-wider">{{ t('instructors.rating') }}</p>
               <div class="flex gap-1 mt-1">
                 <UIcon
                   v-for="i in 5"
                   :key="i"
                   name="i-lucide-star"
                   :class="i <= selectedSession.rating ? 'text-amber-500 fill-amber-500' : 'text-muted'"
                   class="size-4"
                />
               </div>
             </div>
           </div>


           <div v-if="selectedSession.notes" class="p-4 rounded-xl bg-muted/30 border border-default">
             <p class="text-xs text-muted uppercase font-bold tracking-wider mb-2">{{ t('history.instructorNotes') }}</p>
             <p class="text-sm leading-relaxed">{{ selectedSession.notes }}</p>
           </div>
         </div>

       </UCard>
     </template>
   </UModal>
 </ClientOnly>
</template>
