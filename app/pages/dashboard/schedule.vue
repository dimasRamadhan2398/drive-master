<script setup lang="ts">
import { computed, ref, onMounted } from "vue";

const { t } = useI18n();
definePageMeta({ layout: "dashboard" });

const toast = useToast();
const authStore = useAuthStore();
const schedulesStore = useSchedulesStore();
const enrollmentsStore = useEnrollmentsStore();

// FITUR BARU: Logika Kalender Dinamis
const currentDate = ref(new Date("2026-04-10T00:00:00")); // Default start di April 2026
const selectedDate = ref(15);
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

const globalSlots = computed(() => schedulesStore.slots);

onMounted(async () => {
  await Promise.all([
    schedulesStore.initialize(),
    enrollmentsStore.fetchActiveEnrollment()
  ])
})

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

// Real upcoming sessions from store
const upcomingSessions = computed(() => {
  return [...schedulesStore.slots]
    .filter(s => s.status === 'booked' || s.status === 'in-progress')
    .sort((a, b) => new Date(a.date + 'T' + a.time).getTime() - new Date(b.date + 'T' + b.time).getTime())
    .map((s, index) => ({
      ...s,
      sessionNumber: index + 1,
      topic: s.notes || "General Driving Session"
    }))
});

const selectedSlotDetails = computed(() => {
  return globalSlots.value.find((s) => s.id === selectedSlot.value);
});

// FITUR BARU: Reschedule & Cancel State
const showRescheduleModal = ref(false);
const showCancelModal = ref(false);
const sessionToReschedule = ref<any>(null);
const sessionToCancel = ref<any>(null);

const rescheduleDate = ref(15);
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
  rescheduleDate.value = 15;
  rescheduleSlot.value = null;
  showRescheduleModal.value = true;
}

async function confirmReschedule() {
  const newSlot = rescheduleSlotDetails.value;
  const currentEnrollment = enrollmentsStore.currentEnrollment;

  if (newSlot && sessionToReschedule.value && currentEnrollment) {
    const oldSlotId = sessionToReschedule.value.id;

    try {
      // 1. Cancel the old booking
      await schedulesStore.cancelBooking(oldSlotId);

      // 2. Book the new slot
      await schedulesStore.bookSlot(newSlot.id, {
        userId: authStore.user?.userId || "",
        entitlementId: currentEnrollment.id,
        notes: sessionToReschedule.value.topic
      });

      showRescheduleModal.value = false;
      toast.add({
        title: t("schedule.rescheduleSuccess"),
        description: t("schedule.rescheduleSuccessDesc"),
        color: "success",
      });
    } catch (error) {
      toast.add({
        title: "Reschedule Failed",
        description: "An error occurred while rescheduling your session.",
        color: "error"
      });
    }
  }
}

function openCancelModal(session: any) {
  sessionToCancel.value = session;
  showCancelModal.value = true;
}

async function confirmCancel() {
  if (sessionToCancel.value) {
    const canceledSlotId = sessionToCancel.value.id;
    try {
      await schedulesStore.cancelBooking(canceledSlotId);
      showCancelModal.value = false;
      toast.add({
        title: t("schedule.cancelSuccess"),
        description: t("schedule.cancelSuccessDesc"),
        color: "neutral",
      });
    } catch (error) {
      toast.add({
        title: "Cancellation Failed",
        description: "An error occurred while cancelling your session.",
        color: "error"
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
  const currentEnrollment = enrollmentsStore.currentEnrollment;

  if (selectedSlot.value && selectedSlotDetails.value && currentEnrollment) {
    const bookedSlotId = selectedSlot.value;

    try {
      await schedulesStore.bookSlot(bookedSlotId, {
        userId: authStore.user?.userId || "",
        entitlementId: currentEnrollment.id,
        notes: "General Driving Session"
      });

      showBookingModal.value = false;
      toast.add({
        title: t("schedule.bookingSuccess"),
        description: t("schedule.bookingSuccessDesc", {
          date: selectedSlotDetails.value.date,
          time: selectedSlotDetails.value.time,
        }),
        icon: "i-lucide-check-circle",
        color: "success",
      });
      selectedSlot.value = null;
    } catch (error) {
      toast.add({
        title: "Booking Failed",
        description: "An error occurred while booking your session.",
        color: "error"
      });
    }
  } else if (!currentEnrollment) {
    toast.add({
      title: "No Active Enrollment",
      description: "Please purchase a package first to book a session.",
      color: "error"
    });
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

          <div v-if="upcomingSessions.length > 0" class="grid md:grid-cols-2 gap-4">
            <UCard v-for="session in upcomingSessions" :key="session.id">
              <div class="flex items-start justify-between">
                <div class="flex items-center gap-4">
                  <div class="p-3 rounded-xl bg-warning/10">
                    <UIcon name="i-lucide-car" class="size-6 text-warning" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <h3 class="font-semibold">
                        {{ t("history.session") }} #{{ session.sessionNumber }}
                      </h3>
                      <UBadge
                        :label="t('dashboard.confirmed')"
                        color="success"
                        variant="subtle"
                        size="md"
                      />
                    </div>
                    <p class="text-md text-muted">{{ session.topic }}</p>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-default">
                <div>
                  <p class="text-md text-muted">
                    {{ t("dashboard.date") }} & {{ t("dashboard.time") }}
                  </p>
                  <p class="text-md font-medium">{{ session.date }}</p>
                  <p class="text-md">{{ session.time }}</p>
                </div>
                <div>
                  <p class="text-md text-muted">{{ t("dashboard.instructor") }}</p>
                  <p class="text-md font-medium">{{ session.instructor }}</p>
                  <p class="text-md text-muted">{{ session.car }}</p>
                </div>
              </div>

              <template #footer>
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

          <UEmpty
            v-else
            icon="i-lucide-calendar-x"
            :title="t('schedule.noUpcoming')"
            :description="t('schedule.noUpcoming')"
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
              <div class="border border-default rounded-lg p-4">
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
                  <!-- Empty cells for days before month starts (April 2026 starts on Wednesday) -->
                  <div></div>
                  <div></div>
                  <div
                    v-for="(item, idx) in calendarDays"
                    :key="idx"
                    class="w-full aspect-square flex items-center justify-center"
                  >
                    <button
                      v-if="item.day !== null"
                      :disabled="!item.available"
                      :class="[
                        'w-full h-full rounded-lg text-sm font-medium transition-all',
                        selectedDate === item.day
                          ? 'bg-primary text-white'
                          : item.available
                          ? 'hover:bg-primary/10 cursor-pointer'
                          : 'text-muted/50 cursor-not-allowed',
                      ]"
                      @click="item.available && (selectedDate = item.day)"
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
                <UBadge :label="`Apr ${selectedDate}`" color="primary" variant="subtle" />
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
                    date: `Apr ${selectedDate}`,
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
                            <strong>{{ t("dashboard.date") }}:</strong> April
                            {{ selectedDate }}, 2026
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
