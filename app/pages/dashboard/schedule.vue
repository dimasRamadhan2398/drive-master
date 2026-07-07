<script setup lang="ts">
import { computed, ref, onMounted, watch } from "vue";
import { useSchedulesStore } from "~/stores/schedules";
import { useAuthStore } from "~/stores/auth";
import { scheduleService } from "~/services/scheduleService";
import type { SessionResponse } from "~/services/scheduleService";

const { t } = useI18n();
definePageMeta({ layout: "dashboard" });

const toast = useToast();

// FITUR BARU: Logika Kalender Dinamis
const currentDate = ref(new Date());
const selectedDate = ref(new Date().getDate());
const selectedSlot = ref<string | null>(null);
const showBookingModal = ref(false);

const weekDays = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

const currentMonthStr = computed(() => {
  return currentDate.value.toLocaleDateString("en-US", {
    month: "long",
    year: "numeric",
  });
});

// PERUBAHAN: Memperbaiki teks bulan statis
const currentMonthShortStr = computed(() => {
  return currentDate.value.toLocaleDateString("en-US", { month: "short" });
});

const schedulesStore = useSchedulesStore();
const globalSlots = computed(() => schedulesStore.slots);

const isLoadingCalendar = ref(false);

const fetchSchedulesForCalendar = async () => {
  const year = currentDate.value.getFullYear();
  const month = currentDate.value.getMonth(); // 0-11
  
  // Start of month: YYYY-MM-01
  const startDayStr = `${year}-${String(month + 1).padStart(2, "0")}-01`;
  
  // End of month
  const lastDay = new Date(year, month + 1, 0).getDate();
  const endDayStr = `${year}-${String(month + 1).padStart(2, "0")}-${String(lastDay).padStart(2, "0")}`;
  
  try {
    isLoadingCalendar.value = true;
    await schedulesStore.fetchSchedules({
      startDate: startDayStr,
      endDate: endDayStr,
      limit: 200, // retrieve all slots for the month
    });
  } catch (err) {
    console.error("Failed to fetch schedules for calendar:", err);
  } finally {
    isLoadingCalendar.value = false;
  }
};

watch(currentDate, async () => {
  await fetchSchedulesForCalendar();
}, { immediate: true });

// FITUR BARU: Kalender merender hari secara dinamis berdasarkan bulan yang sedang dipilih
const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear();
  const month = currentDate.value.getMonth(); // 0-11

  // Cari tahu tanggal 1 bulan ini jatuh di hari apa (0 = Sun, 1 = Mon)
  const firstDay = new Date(year, month, 1).getDay();
  // Konversi ke format Senin=0, Minggu=6
  const emptyDays = firstDay === 0 ? 6 : firstDay - 1;

  const daysInMonth = new Date(year, month + 1, 0).getDate();

  const days = [];

  // Selipkan hari kosong sebelum tanggal 1
  for (let i = 0; i < emptyDays; i++) {
    days.push({ day: null, available: false });
  }

  // Generate tanggal 1 sampai akhir bulan
  for (let i = 1; i <= daysInMonth; i++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, "0")}-${String(i).padStart(
      2,
      "0"
    )}`;
    // Periksa apakah ada slot available di tanggal ini dari global state
    const isAvailable = globalSlots.value.some(
      (s) => s.date === dateStr && s.status === "available"
    );

    days.push({
      day: i,
      available: isAvailable,
    });
  }

  return days;
});

// FITUR BARU: Fungsi untuk navigasi bulan di kalender
function changeMonth(offset: number) {
  const newDate = new Date(currentDate.value);
  newDate.setMonth(newDate.getMonth() + offset);
  currentDate.value = newDate;
}

// FITUR BARU: Hanya tampilkan slot yang sesuai dengan tanggal yang dipilih di kalender
const availableSlots = computed(() => {
  const year = currentDate.value.getFullYear();
  const month = String(currentDate.value.getMonth() + 1).padStart(2, "0");
  const day = String(selectedDate.value).padStart(2, "0");
  const dateStr = `${year}-${month}-${day}`;

  return globalSlots.value
    .filter((slot) => slot.date === dateStr)
    .map((slot) => ({
      ...slot,
      available: slot.status === "available",
    }));
});



const authStore = useAuthStore();

const userSessionsList = ref<SessionResponse[]>([]);
const isLoadingSessions = ref(false);
const currentPage = ref(1);
const itemsPerPage = ref(4); // 4 sessions per page
const totalSessions = ref(0);

const activeEntitlement = computed(() => {
  return authStore.memberEntitlements.find(e => e.status === "active") || authStore.memberEntitlements[0];
});

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

const getInstructorName = (session: SessionResponse) => {
  if (session.scheduleId) {
    const slot = globalSlots.value.find((s) => s.id === String(session.scheduleId));
    if (slot) return slot.instructor;
  }
  return "Instructor";
};

const getCarName = (session: SessionResponse) => {
  if (session.scheduleId) {
    const slot = globalSlots.value.find((s) => s.id === String(session.scheduleId));
    if (slot) return slot.car;
  }
  return "BYD Atto 3";
};

const fetchSessions = async () => {
  if (!authStore.userId) return;
  try {
    isLoadingSessions.value = true;
    const response = await scheduleService.fetchSessions({
      studentId: authStore.userId,
      page: currentPage.value,
      limit: itemsPerPage.value,
    });
    const sessions = response.data || [];
    userSessionsList.value = sessions.slice().sort((a, b) => {
      const dateA = new Date(`${a.date}T${a.time || '00:00'}:00`);
      const dateB = new Date(`${b.date}T${b.time || '00:00'}:00`);
      
      const isAActive = a.status === 'booked' || a.status === 'in-progress' || a.status === 'in_progress' || a.status === 'scheduled';
      const isBActive = b.status === 'booked' || b.status === 'in-progress' || b.status === 'in_progress' || b.status === 'scheduled';

      if (isAActive && !isBActive) return -1;
      if (!isAActive && isBActive) return 1;

      if (isAActive && isBActive) {
        // Upcoming active sessions: sort ascending (closest to today first)
        return dateA.getTime() - dateB.getTime();
      }

      // Past inactive sessions: sort descending (most recent first)
      return dateB.getTime() - dateA.getTime();
    });
    totalSessions.value = response.total || 0;
  } catch (err) {
    console.error("Failed to fetch sessions:", err);
  } finally {
    isLoadingSessions.value = false;
  }
};

onMounted(async () => {
  if (!authStore.hasMemberProfile) {
    await authStore.fetchMemberProfile();
  }
  await fetchSessions();
});

watch(currentPage, async () => {
  await fetchSessions();
});

const selectedSlotDetails = computed(() => {
  return globalSlots.value.find((s) => s.id === selectedSlot.value);
});

const showRescheduleModal = ref(false);
const showCancelModal = ref(false);
const sessionToReschedule = ref<any>(null);
const sessionToCancel = ref<any>(null);

const rescheduleDate = ref(new Date().getDate());
const rescheduleSlot = ref<string | null>(null);

// Computed untuk slot di modal reschedule
const rescheduleAvailableSlots = computed(() => {
  const year = currentDate.value.getFullYear();
  const month = String(currentDate.value.getMonth() + 1).padStart(2, "0");
  const day = String(rescheduleDate.value).padStart(2, "0");
  const dateStr = `${year}-${month}-${day}`;

  return globalSlots.value
    .filter((slot) => slot.date === dateStr)
    .map((slot) => ({
      ...slot,
      available: slot.status === "available",
    }));
});

const rescheduleSlotDetails = computed(() => {
  return globalSlots.value.find((s) => s.id === rescheduleSlot.value);
});

function openRescheduleModal(session: any) {
  sessionToReschedule.value = session;
  if (session && session.date) {
    const parts = session.date.split('-').map(Number);
    rescheduleDate.value = parts[2] || new Date().getDate();
  } else {
    rescheduleDate.value = new Date().getDate();
  }
  rescheduleSlot.value = null;
  showRescheduleModal.value = true;
}

async function confirmReschedule() {
  const newSlot = rescheduleSlotDetails.value;
  const session = sessionToReschedule.value;
  if (newSlot && session) {
    try {
      // 1. Cancel the old slot if it has a scheduleId
      if (session.scheduleId) {
        await scheduleService.cancelBooking(String(session.scheduleId));
      }
      
      // 2. Book the new slot
      if (authStore.userId) {
        const entitlement = activeEntitlement.value;
        if (!entitlement) {
          throw new Error("No active entitlement found");
        }
        await scheduleService.bookSlot(newSlot.id, {
          userId: authStore.userId,
          entitlementId: entitlement.id,
          notes: session.notes || "Rescheduled driving session"
        });
      }

      // 3. Make the old slot available and the new slot booked in store
      if (session.scheduleId) {
        schedulesStore.updateSlotStatus(String(session.scheduleId), "available");
      }
      schedulesStore.updateSlotStatus(newSlot.id, "booked");

      // 4. Refresh sessions list
      await fetchSessions();
      
      showRescheduleModal.value = false;
      toast.add({
        title: t("schedule.rescheduleSuccess"),
        description: t("schedule.rescheduleSuccessDesc"),
        color: "success",
      });
    } catch (err) {
      console.error("Failed to reschedule session:", err);
      toast.add({
        title: "Error",
        description: "Failed to reschedule session. Please try again.",
        color: "error",
      });
    }
  }
}

function openCancelModal(session: any) {
  sessionToCancel.value = session;
  showCancelModal.value = true;
}

async function confirmCancel() {
  const session = sessionToCancel.value;
  if (session) {
    try {
      // 1. Cancel the booking via API if scheduleId exists
      if (session.scheduleId) {
        await scheduleService.cancelBooking(String(session.scheduleId));
      }

      // 2. Update status in store (for fallback/sync)
      if (session.scheduleId) {
        schedulesStore.updateSlotStatus(String(session.scheduleId), "available");
      }

      // 3. Refresh sessions list
      await fetchSessions();
      
      showCancelModal.value = false;
      toast.add({
        title: t("schedule.cancelSuccess"),
        description: t("schedule.cancelSuccessDesc"),
        color: "neutral",
      });
    } catch (err) {
      console.error("Failed to cancel session:", err);
      toast.add({
        title: "Error",
        description: "Failed to cancel session. Please try again.",
        color: "error",
      });
    }
  }
}

function selectSlot(slotId: string) {
  const slot = globalSlots.value.find((s) => s.id === slotId);
  if (slot?.status === "available") {
    selectedSlot.value = slotId;
  }
}

async function confirmBooking() {
  if (selectedSlot.value && selectedSlotDetails.value) {
    const bookedSlotId = selectedSlot.value;
    const bookedSlotDetails = selectedSlotDetails.value;

    if (!authStore.userId) {
      toast.add({ title: "Error", description: "User not authenticated", color: "error" });
      return;
    }

    const entitlement = activeEntitlement.value;
    if (!entitlement) {
      toast.add({ 
        title: "Error", 
        description: "No active package found to book this session. Please purchase a package first.", 
        color: "error" 
      });
      return;
    }

    try {
      // 1. Call API to book the slot
      const bookedSlot = await scheduleService.bookSlot(bookedSlotId, {
        userId: authStore.userId,
        entitlementId: entitlement.id,
        notes: "Booked via dashboard schedule page"
      });

      if (!bookedSlot) {
        throw new Error("API booking failed");
      }

      // 2. Update local slot status
      schedulesStore.updateSlotStatus(bookedSlotId, "booked");

      // 3. Refresh sessions list
      await fetchSessions();

      showBookingModal.value = false;
      toast.add({
        title: t("schedule.bookingSuccess"),
        description: t("schedule.bookingSuccessDesc", {
          date: formatDate(bookedSlotDetails.date),
          time: bookedSlotDetails.time,
        }),
        icon: "i-lucide-check-circle",
        color: "success",
      });
      selectedSlot.value = null;
    } catch (err) {
      console.error("Failed to book session:", err);
      toast.add({
        title: "Booking Failed",
        description: "The time slot might have just been taken. Please choose another slot.",
        color: "error",
      });
    }
  }
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('schedule.title')">
        <template #right>
          <UButton icon="i-lucide-bell" color="neutral" variant="ghost" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Upcoming Sessions -->
        <div>
          <h2 class="text-lg font-semibold mb-4">{{ t("schedule.upcomingSessions") }}</h2>

          <!-- Loading State -->
          <div v-if="isLoadingSessions" class="grid md:grid-cols-2 gap-4">
            <UCard v-for="n in 2" :key="n">
              <div class="flex items-start justify-between">
                <div class="flex items-center gap-4 w-full">
                  <USkeleton class="h-12 w-12 rounded-xl" />
                  <div class="flex-1 space-y-2">
                    <USkeleton class="h-5 w-[40%]" />
                    <USkeleton class="h-4 w-[70%]" />
                  </div>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-default">
                <div class="space-y-2">
                  <USkeleton class="h-3 w-[50%]" />
                  <USkeleton class="h-4 w-[80%]" />
                </div>
                <div class="space-y-2">
                  <USkeleton class="h-3 w-[50%]" />
                  <USkeleton class="h-4 w-[80%]" />
                </div>
              </div>
            </UCard>
          </div>

          <!-- Sessions List -->
          <div v-else-if="userSessionsList.length > 0" class="space-y-4">
            <div class="grid md:grid-cols-2 gap-4">
              <UCard v-for="(session, index) in userSessionsList" :key="session.id">
                <div class="flex items-start justify-between">
                  <div class="flex items-center gap-4">
                    <div class="p-3 rounded-xl bg-warning/10">
                      <UIcon name="i-lucide-car" class="size-6 text-warning" />
                    </div>
                    <div>
                      <div class="flex items-center gap-2">
                        <h3 class="font-semibold">
                          {{ t("history.session") }} #{{ (activeEntitlement?.usedSessions ?? 0) + 1 + index }}
                        </h3>
                        <UBadge
                          :label="t('history.' + session.status) || session.status"
                          :color="scheduleService.getStatusColor(session.status as any)"
                          variant="subtle"
                          size="md"
                        />
                      </div>
                      <p class="text-md text-muted truncate max-w-[200px]">
                        {{ session.notes || 'Training Session' }}
                      </p>
                    </div>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-default">
                  <div>
                    <p class="text-md text-muted">
                      {{ t("dashboard.date") }} & {{ t("dashboard.time") }}
                    </p>
                    <p class="text-md font-medium">{{ formatDate(session.date) }}</p>
                    <p class="text-md">{{ session.time }}</p>
                  </div>
                  <div>
                    <p class="text-md text-muted">{{ t("dashboard.instructor") }}</p>
                    <p class="text-md font-medium">{{ getInstructorName(session) }}</p>
                    <p class="text-md text-muted">{{ getCarName(session) }}</p>
                  </div>
                </div>

                <template #footer v-if="session.status === 'booked' || session.status === 'in-progress'">
                  <div class="flex gap-2">
                    <UButton
                      :label="t('dashboard.reschedule')"
                      variant="outline"
                      color="warning"
                      size="md"
                      icon="i-lucide-calendar-days"
                      @click="openRescheduleModal(session)"
                    />
                    <UButton
                      :label="t('common.cancel')"
                      variant="ghost"
                      color="error"
                      size="md"
                      icon="i-lucide-x"
                      @click="openCancelModal(session)"
                    />
                  </div>
                </template>
              </UCard>
            </div>

            <!-- Pagination Control -->
            <div class="flex justify-end pt-2">
              <UPagination
                v-model="currentPage"
                :total="totalSessions"
                :items-per-page="itemsPerPage"
                active-color="warning"
              />
            </div>
          </div>

          <UEmpty
            v-else
            icon="i-lucide-calendar-x"
            :title="t('schedule.noUpcoming')"
          />
        </div>

        <!-- Book New Session -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">{{ t("schedule.bookNewSession") }}</h2>
          </template>

          <div class="grid lg:grid-cols-2 gap-8">
            <!-- Calendar -->
            <div>
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-md font-medium">{{ t("schedule.selectDate") }}</h3>
                <div class="flex items-center gap-2">
                  <UButton
                    icon="i-lucide-chevron-left"
                    variant="ghost"
                    color="neutral"
                    size="md"
                    @click="changeMonth(-1)"
                  />
                  <span class="text-md font-medium">{{ currentMonthStr }}</span>
                  <UButton
                    icon="i-lucide-chevron-right"
                    variant="ghost"
                    color="neutral"
                    size="md"
                    @click="changeMonth(1)"
                  />
                </div>
              </div>

              <!-- Custom Calendar Grid -->
              <div class="border border-default rounded-lg p-4 relative">
                <!-- Loading Overlay -->
                <div v-if="isLoadingCalendar" class="absolute inset-0 bg-white/70 dark:bg-gray-950/70 z-10 flex items-center justify-center rounded-lg">
                  <UIcon name="i-lucide-loader-2" class="w-8 h-8 animate-spin text-warning-500" />
                </div>
                <div class="grid grid-cols-7 gap-1 mb-2">
                  <div
                    v-for="day in weekDays"
                    :key="day"
                    class="text-center text-md font-medium text-muted py-2"
                  >
                    {{ day }}
                  </div>
                </div>
                <div class="grid grid-cols-7 gap-1">
                  <div
                    v-for="(item, idx) in calendarDays"
                    :key="idx"
                    class="w-full aspect-square flex items-center justify-center"
                  >
                    <button
                      v-if="item.day !== null"
                      :class="[
                        'w-full h-full rounded-lg text-sm font-medium transition-all cursor-pointer',
                        selectedDate === item.day
                          ? 'bg-primary text-white shadow-sm'
                          : item.available
                          ? 'hover:bg-primary/10 text-primary-600 dark:text-primary-400 font-semibold'
                          : 'hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-700 dark:text-gray-300',
                      ]"
                      @click="selectedDate = item.day"
                    >
                      {{ item.day }}
                    </button>
                  </div>
                </div>

                <div class="mt-4 flex items-center gap-4 text-md">
                  <div class="flex items-center gap-2">
                    <div
                      class="size-3 rounded bg-primary/10 border border-primary/30"
                    ></div>
                    <span class="text-muted">{{ t("common.available") }}</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <div class="size-3 rounded bg-primary"></div>
                    <span class="text-muted">{{ t("home.selected") }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Time Slots -->
            <div>
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-md font-medium">{{ t("schedule.availableSlots") }}</h3>
                <UBadge :label="`${currentMonthShortStr} ${selectedDate}`" color="primary" variant="subtle" />
              </div>

              <div class="space-y-3">
                <button
                  v-for="slot in availableSlots"
                  :key="slot.id"
                  :disabled="!slot.available"
                  :class="[
                    'w-full p-4 rounded-lg border text-left transition-all',
                    slot.available
                      ? selectedSlot === slot.id
                        ? 'border-primary bg-primary/10'
                        : 'border-default hover:border-primary cursor-pointer'
                      : 'border-default bg-muted/50 opacity-50 cursor-not-allowed',
                  ]"
                  @click="selectSlot(slot.id)"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3">
                      <UIcon
                        :name="
                          selectedSlot === slot.id
                            ? 'i-lucide-check-circle'
                            : 'i-lucide-circle'
                        "
                        :class="selectedSlot === slot.id ? 'text-primary' : 'text-muted'"
                        class="size-5"
                      />
                      <div>
                        <span class="font-semibold">{{ slot.time }}</span>
                        <p class="text-md text-muted">
                          {{ slot.instructor }} - {{ slot.car }}
                        </p>
                      </div>
                    </div>
                    <UBadge
                      :label="slot.available ? t('common.available') : t('home.booked')"
                      :color="slot.available ? 'success' : 'error'"
                      variant="subtle"
                      size="md"
                    />
                  </div>
                </button>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="flex items-center justify-between">
              <p v-if="selectedSlot" class="text-md text-muted">
                {{
                  t("schedule.selectedInfo", {
                    date: `${currentMonthShortStr} ${selectedDate}`,
                    time: selectedSlotDetails?.time,
                    instructor: selectedSlotDetails?.instructor,
                  })
                }}
              </p>
              <p v-else class="text-md text-muted">
                {{ t("schedule.selectDateCont") }}
              </p>
              <UButton
                :label="t('schedule.bookNow')"
                :disabled="!selectedSlot"
                icon="i-lucide-check"
                color="warning"
                @click="showBookingModal = true"
              />
              <!-- Booking Confirmation Modal -->
              <UModal
                v-model:open="showBookingModal"
                :title="t('schedule.confirmBooking')"
              >
                <template #body>
                  <div class="space-y-4">
                    <UAlert
                      icon="i-lucide-info"
                      color="warning"
                      :title="t('schedule.sessionDetails')"
                    >
                      <template #description>
                        <ul class="mt-2 space-y-1 text-md">
                          <li>
                            <strong>{{ t("dashboard.date") }}:</strong>
                            {{ selectedDate }} {{ currentMonthStr }}
                          </li>
                          <li>
                            <strong>{{ t("dashboard.time") }}:</strong>
                            {{ selectedSlotDetails?.time }}
                          </li>
                          <li>
                            <strong>{{ t("dashboard.vehicle") }}:</strong>
                            {{ selectedSlotDetails?.car }}
                          </li>
                          <li>
                            <strong>{{ t("dashboard.instructor") }}:</strong>
                            {{ selectedSlotDetails?.instructor }}
                          </li>
                        </ul>
                      </template>
                    </UAlert>

                    <p class="text-md text-muted">
                      {{ t("register.terms.agree") }}
                      {{ t("register.terms.termsOfService") }}.
                    </p>
                  </div>
                </template>
                <template #footer>
                  <div class="flex justify-end gap-3">
                    <UButton
                      :label="t('common.cancel')"
                      variant="ghost"
                      color="neutral"
                      @click="showBookingModal = false"
                    />
                    <UButton
                      :label="t('schedule.confirmBooking')"
                      color="warning"
                      icon="i-lucide-check"
                      @click="confirmBooking"
                    />
                  </div>
                </template>
              </UModal>
            </div>
          </template>
        </UCard>

        <!-- Reschedule Modal -->
        <UModal
          v-model:open="showRescheduleModal"
          :title="t('schedule.rescheduleTitle')"
          class="max-w-2xl"
        >
          <template #body>
            <div class="grid md:grid-cols-2 gap-6">
              <div>
                <h4 class="text-sm font-medium mb-3">
                  {{ t("schedule.chooseNewDate") }}
                </h4>
                <div class="border border-default rounded-lg p-3">
                  <div class="grid grid-cols-7 gap-1">
                    <div
                      v-for="(item, idx) in calendarDays"
                      :key="idx"
                      class="aspect-square"
                    >
                      <button
                        v-if="item.day !== null"
                        class="w-full h-full text-xs rounded-md transition-all"
                        :class="[
                          rescheduleDate === item.day
                            ? 'bg-primary text-white'
                            : 'hover:bg-primary/10',
                        ]"
                        @click="rescheduleDate = item.day"
                      >
                        {{ item.day }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              <div>
                <h4 class="text-sm font-medium mb-3">
                  {{ t("schedule.availableSlots") }} ({{ t("dashboard.date") }}
                  {{ rescheduleDate }})
                </h4>
                <div class="space-y-2 max-h-[300px] overflow-y-auto pr-1">
                  <button
                    v-for="slot in rescheduleAvailableSlots"
                    :key="slot.id"
                    :disabled="!slot.available"
                    class="w-full p-3 text-left border rounded-lg transition-all text-sm"
                    :class="[
                      rescheduleSlot === slot.id
                        ? 'border-primary bg-primary/5'
                        : 'border-default hover:border-primary',
                      !slot.available ? 'opacity-40 cursor-not-allowed' : '',
                    ]"
                    @click="rescheduleSlot = slot.id"
                  >
                    <div class="font-bold">{{ slot.time }}</div>
                    <div class="text-xs text-muted">{{ slot.instructor }}</div>
                  </button>
                </div>
              </div>
            </div>
          </template>
          <template #footer>
            <div class="flex justify-between items-center w-full">
              <span v-if="rescheduleSlotDetails" class="text-xs text-muted"
                >Your New Time: {{ rescheduleSlotDetails.time }}</span
              >
              <div class="flex gap-2">
                <UButton
                  :label="t('common.cancel')"
                  variant="ghost"
                  color="neutral"
                  @click="showRescheduleModal = false"
                />
                <UButton
                  :label="t('schedule.confirmReschedule')"
                  color="warning"
                  :disabled="!rescheduleSlot"
                  @click="confirmReschedule"
                />
              </div>
            </div>
          </template>
        </UModal>

        <!-- Cancel Confirmation Modal -->
        <UModal v-model:open="showCancelModal" :title="t('schedule.cancelSession')">
          <template #body>
            <div class="text-center py-4">
              <div
                class="size-16 rounded-full bg-error/10 flex items-center justify-center mx-auto mb-4"
              >
                <UIcon name="i-lucide-alert-triangle" class="size-8 text-error" />
              </div>
              <h3 class="text-lg font-bold">{{ t("common.confirm") }}</h3>
              <p class="text-sm text-muted mt-1">
                {{ t("schedule.confirmCancel") }}
              </p>
            </div>
          </template>
          <template #footer>
            <div class="flex justify-center gap-3 w-full">
              <UButton
                :label="t('common.cancel')"
                variant="outline"
                color="neutral"
                class="flex-1"
                @click="showCancelModal = false"
              />
              <UButton
                :label="t('schedule.confirm')"
                color="error"
                class="flex-1"
                @click="confirmCancel"
              />
            </div>
          </template>
        </UModal>
      </div>
    </template>
  </UDashboardPanel>
</template>
