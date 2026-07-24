<script setup lang="ts">
import { useToast, useInstructorsStore, useVehiclesStore, useSchedulesStore, useI18n, useLocale, useSettings } from "#imports";
import { ref, computed, onMounted } from "vue";
import type { ScheduleSlot } from "~/stores/schedules";
import { useRoute, navigateTo } from "#imports";
import AddSlotModal from "~/components/schedules/AddSlotModal.vue";
import EditSlotModal from "~/components/schedules/EditSlotModal.vue";
import ManualBookingModal from "~/components/schedules/ManualBookingModal.vue";

const { t } = useI18n();
const { locale } = useLocale();
definePageMeta({ layout: "admin" });

const route = useRoute();
const toast = useToast();
const instructorsStore = useInstructorsStore();
const vehiclesStore = useVehiclesStore();
const schedulesStore = useSchedulesStore();
const studentNameToBook = computed(() => route.query.studentName as string | undefined);
const showAddSlotModal = ref(false);
const showManualBookingModal = ref(false);
const selectedBookingSlotId = ref("");
const selectedDate = ref(new Date());

// FIX: Define localDateStr to format the selected date
const localDateStr = computed(() => {
  const date = selectedDate.value;
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  return `${year}-${month}-${day}`;
});

const filterInstructor = ref("All Instructors");
const filterVehicle = ref("All Vehicles");

// Use store's slots for display
const timeSlots = computed(() => {
  return schedulesStore.slots;
});

const instructors = computed(() => {
  return instructorsStore.instructors.map((value) => {
    return {
      label: value.name,
      value: value.userId,
    };
  });
});

const instructorOptions = computed(() => [
  { label: "All Instructors", value: "All Instructors" },
  ...instructors.value,
]);

const vehicleOptionsList = computed(() => {
  return vehiclesStore.vehicles.map((v) => ({
    label: `${v.brand} ${v.model}`,
    value: v.id,
  }));
});

// FITUR BARU: Sinkronisasi jam operasional dari settings
const { operatingHours } = useSettings();
const currentDayOperatingHours = computed(() => {
  const day = selectedDate.value.getDay(); // 0 = Sunday, 6 = Saturday
  const isWeekend = day === 0 || day === 6;

  if (isWeekend) {
    return {
      start: operatingHours.value.weekendStart,
      end: operatingHours.value.weekendEnd,
      nightStart: operatingHours.value.nightStart,
      nightEnd: operatingHours.value.nightEnd,
      nightEnabled: operatingHours.value.nightEnabled,
      isClosed: day === 0 && operatingHours.value.sundayClosed,
    };
  }

  return {
    start: operatingHours.value.mondayStart,
    end: operatingHours.value.mondayEnd,
    nightStart: operatingHours.value.nightStart,
    nightEnd: operatingHours.value.nightEnd,
    nightEnabled: operatingHours.value.nightEnabled,
    isClosed: false,
  };
});

const adminCalendarDays = computed(() => {
  const year = selectedDate.value.getFullYear();
  const month = selectedDate.value.getMonth();

  const firstDay = new Date(year, month, 1).getDay();
  const emptyDays = firstDay === 0 ? 6 : firstDay - 1;
  const daysInMonth = new Date(year, month + 1, 0).getDate();

  const days = [];
  for (let i = 0; i < emptyDays; i++) {
    days.push({ day: null, hasSlots: false });
  }

  for (let i = 1; i <= daysInMonth; i++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, "0")}-${String(i).padStart(2, "0")}`;
    const hasSlots = timeSlots.value.some((s) => s.date === dateStr);
    days.push({
      day: i,
      hasSlots,
    });
  }

  return days;
});

const filteredSlots = computed(() => {
  const dateStr = localDateStr.value; // FIX: Menggunakan localDateStr
  return timeSlots.value.filter((slot) => {
    const matchDate = slot.date === dateStr;
    const matchInst =
      filterInstructor.value === "All Instructors" ||
      slot.instructorId === filterInstructor.value ||
      slot.instructor === filterInstructor.value;
    const matchVeh =
      filterVehicle.value === "All Vehicles" ||
      slot.carId === filterVehicle.value ||
      slot.car === filterVehicle.value;
    return matchDate && matchInst && matchVeh;
  });
});

function changeDay(offset: number) {
  const nextDate = new Date(selectedDate.value);
  nextDate.setDate(selectedDate.value.getDate() + offset);
  selectedDate.value = nextDate;
}

async function toggleSlotStatus(slotId: string) {
  const currentSlot = timeSlots.value.find((s) => s.id === slotId);
  if (!currentSlot) return;

  const newStatus = currentSlot.status === "blocked" ? "available" : "blocked";
  const success = await schedulesStore.updateSlot(slotId, { status: newStatus });

  if (success) {
    if (newStatus === "blocked") {
      toast.add({ title: "Slot Blocked", color: "warning", icon: "i-lucide-lock" });
    } else {
      toast.add({
        title: `Slot is now ${newStatus}`,
        color: "success",
        icon: "i-lucide-check",
      });
    }
  }
}

async function deleteSlot(slotId: string) {
  if (confirm("Are you sure you want to delete this slot?")) {
    const success = await schedulesStore.deleteSlot(slotId);
    if (success) {
      toast.add({ title: "Slot Deleted", color: "error", icon: "i-lucide-trash" });
    }
  }
}

// FITUR BARU: Logika In-Progress Session
const isToday = computed(() => {
  const today = new Date();
  const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(
    2,
    "0"
  )}-${String(today.getDate()).padStart(2, "0")}`;
  return localDateStr.value === todayStr;
});

async function startSession(slotId: string) {
  if (
    confirm('Apakah kursus sudah mau jalan? Sesi akan berubah menjadi "In Progress".')
  ) {
    const success = await schedulesStore.startSession(slotId);
    if (success) {
      toast.add({
        title: "Session Started",
        description: "Driving session is now in progress.",
        color: "warning",
        icon: "i-lucide-play",
      });
    }
  }
}

async function completeSession(slotId: string) {
  if (confirm('Apakah kursus sudah selesai? Sesi akan ditandai sebagai "Completed".')) {
    const success = await schedulesStore.completeSession(slotId);
    if (success) {
      toast.add({
        title: "Session Completed",
        description: "Driving session has been finished.",
        color: "neutral",
        icon: "i-lucide-check-circle",
      });
    }
  }
}

function handleManualBooking(slotId: string) {
  selectedBookingSlotId.value = slotId;
  showManualBookingModal.value = true;
}

function handleSlotClick(slot: any) {
  if (studentNameToBook.value && slot.status === "available") {
    const studentName = decodeURIComponent(studentNameToBook.value as string);
    schedulesStore.bookSlotLocal(slot.id, studentName);
    toast.add({
      title: "Slot Booked",
      description: `Berhasil booking untuk ${studentName}`,
      color: "success",
      icon: "i-lucide-calendar-check",
    });
    // Kembali ke halaman schedule tanpa query untuk membersihkan state
    navigateTo("/admin/schedules");
  }
}

async function cancelBooking(slotId: string) {
  if (confirm("Batalkan booking dan kembalikan slot menjadi Available?")) {
    const success = await schedulesStore.cancelBooking(slotId);
    if (success) {
      toast.add({
        title: "Booking Cancelled",
        description: "Slot is now available again.",
        color: "neutral",
      });
    }
  }
}

// FITUR BARU: Add Slot State & Logic
async function handleAddSlot(form: {
  date?: string;
  time: string;
  duration: string;
  carId: string;
  instructorId: string;
}) {
  if (!form.time || !form.carId || !form.instructorId) {
    toast.add({
      title: t("common.error"),
      description: "Please fill all fields",
      color: "error",
    });
    return;
  }

  const {
    start,
    end,
    nightStart,
    nightEnd,
    nightEnabled,
    isClosed,
  } = currentDayOperatingHours.value;
  if (isClosed) {
    toast.add({
      title: t("common.error"),
      description: "Business is closed on this day",
      color: "error",
    });
    return;
  }

  const isDayShift = form.time >= start && form.time <= end;
  const isNightShift = nightEnabled && form.time >= nightStart && form.time <= nightEnd;

  if (!isDayShift && !isNightShift) {
    let msg = `Time must be between ${start} - ${end}`;
    if (nightEnabled) msg += ` or ${nightStart} - ${nightEnd}`;
    toast.add({ title: t("common.error"), description: msg, color: "error" });
    return;
  }

  const id = Date.now().toString();
  await schedulesStore.createSlot({
    date: form.date ?? localDateStr.value,
    time: form.time,
    duration: parseInt(form.duration),
    carId: form.carId,
    instructorId: form.instructorId,
  });
  await fetchAdminSchedules();

  toast.add({ title: "New Slot Added", color: "success", icon: "i-lucide-check" });
  showAddSlotModal.value = false;
}

// FITUR BARU: Edit Slot State & Logic
const showEditSlotModal = ref(false);
const selectedEditSlot = ref<any>(null);

function openEditModal(slot: ScheduleSlot) {
  if (slot.status !== "available") {
    toast.add({
      title: t("common.error"),
      description: "Hanya slot dengan status Available yang bisa diedit",
      color: "error",
    });
    return;
  }
  selectedEditSlot.value = {
    ...slot,
    car: slot.carId,
    instructor: slot.instructorId
  };
  showEditSlotModal.value = true;
}

async function handleEditSlot(updated: {
  id: string;
  time: string;
  duration: string;
  carId: string;
  instructorId: string;
}) {
  if (!updated.time || !updated.carId || !updated.instructorId) {
    toast.add({
      title: t("common.error"),
      description: "Please fill all fields",
      color: "error",
    });
    return;
  }

  const {
    start,
    end,
    nightStart,
    nightEnd,
    nightEnabled,
    isClosed,
  } = currentDayOperatingHours.value;
  if (isClosed) {
    toast.add({
      title: t("common.error"),
      description: "Business is closed on this day",
      color: "error",
    });
    return;
  }

  const isDayShift = updated.time >= start && updated.time <= end;
  const isNightShift =
    nightEnabled && updated.time >= nightStart && updated.time <= nightEnd;

  if (!isDayShift && !isNightShift) {
    let msg = `Time must be between ${start} - ${end}`;
    if (nightEnabled) msg += ` or ${nightStart} - ${nightEnd}`;
    toast.add({ title: t("common.error"), description: msg, color: "error" });
    return;
  }

  await schedulesStore.updateSlot(updated.id, {
    time: updated.time,
    duration: parseInt(updated.duration),
    carId: updated.carId,
    instructorId: updated.instructorId,
  });

  toast.add({ title: "Slot Updated", color: "success", icon: "i-lucide-check" });
  showEditSlotModal.value = false;
}

const isLoadingAdminSchedules = ref(false);

const fetchAdminSchedules = async () => {
  const date = selectedDate.value;
  const year = date.getFullYear();
  const month = date.getMonth(); // 0-11

  // Start of month: YYYY-MM-01
  const startDayStr = `${year}-${String(month + 1).padStart(2, "0")}-01`;

  // End of month
  const lastDay = new Date(year, month + 1, 0).getDate();
  const endDayStr = `${year}-${String(month + 1).padStart(2, "0")}-${String(lastDay).padStart(2, "0")}`;

  try {
    isLoadingAdminSchedules.value = true;
    await schedulesStore.fetchSchedules({
      startDate: startDayStr,
      endDate: endDayStr,
      limit: 500, // retrieve all slots for the month
    });
  } catch (err) {
    console.error("Failed to fetch admin schedules:", err);
  } finally {
    isLoadingAdminSchedules.value = false;
  }
};

watch(
  () => [selectedDate.value.getFullYear(), selectedDate.value.getMonth()],
  async () => {
    await fetchAdminSchedules();
  },
  { immediate: true }
);

onMounted(() => {
  instructorsStore.fetchInstructors();
  vehiclesStore.fetchVehicles();
});

// Debug: Watch for slot changes
watch(
  () => schedulesStore.slots,
  (newSlots) => {
    console.log('[AdminSchedules] slots changed:', newSlots.map(s => ({ id: s.id, status: s.status })));
  },
  { deep: true }
);
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UAlert
        v-if="studentNameToBook"
        icon="i-lucide-user-check"
        color="info"
        variant="subtle"
        class="m-4"
        :title="t('dashboard.bookingMode')"
      >
        <template #description>
          <span
            >{{ t("admin.selectStudent").replace("Pilih Murid", "Pilih slot untuk") }}
            <strong>{{ decodeURIComponent(studentNameToBook) }}</strong
            >.
            <UButton
              variant="link"
              :padded="false"
              @click="navigateTo('/admin/schedules')"
              >{{ t("schedule.cancel") }}</UButton
            ></span
          >
        </template>
      </UAlert>
      <UDashboardNavbar :title="t('admin.schedules')">
        <template #right>
          <UButton
            icon="i-lucide-plus"
            color="warning"
            :label="t('admin.addNew').replace('Tambah Baru', 'Tambah Slot')"
            @click="showAddSlotModal = true"
          />
          <UColorModeButton />
          <AddSlotModal
            v-model:open="showAddSlotModal"
            :date="localDateStr"
            :instructors="instructors"
            :vehicles="vehicleOptionsList"
            :operatingHours="currentDayOperatingHours"
            @saved="handleAddSlot"
          />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <div class="flex items-center gap-2">
            <UButton
              icon="i-lucide-chevron-left"
              variant="ghost"
              color="neutral"
              @click="changeDay(-1)"
            />
            <div
              class="relative flex items-center justify-center min-w-[160px] hover:bg-gray-100 rounded py-1 transition-colors cursor-pointer"
            >
              <span class="font-medium text-center pointer-events-none">{{
                selectedDate.toLocaleDateString(locale === "id" ? "id-ID" : "en-US", {
                  month: "long",
                  day: "numeric",
                  year: "numeric",
                })
              }}</span>
              <input
                type="date"
                class="absolute inset-0 opacity-0 cursor-pointer w-full h-full"
                :value="localDateStr"
                @input="e => { if ((e.target as HTMLInputElement).value) selectedDate = new Date((e.target as HTMLInputElement).value + 'T00:00:00') }"
              />
            </div>
            <UButton
              icon="i-lucide-chevron-right"
              variant="ghost"
              color="neutral"
              @click="changeDay(1)"
            />
          </div>
        </template>
        <template #right>
          <UButton
            icon="i-lucide-calendar"
            :label="t('time.today')"
            color="neutral"
            variant="outline"
            @click="selectedDate = new Date()"
          />
          <USelect
            v-model="filterInstructor"
            :items="instructorOptions"
            :placeholder="t('common.filter') + ' ' + t('dashboard.instructor')"
            class="w-48"
            color="warning"
          />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Stats -->
        <!-- <div class="grid sm:grid-cols-2 lg:grid-cols-5 gap-4">
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-green-500/10">
                <UIcon name="i-lucide-check-circle" class="size-5 text-green-500" />
              </div>
              <div>
                <p class="text-xl font-bold">
                  {{ timeSlots.filter((s) => s.status === "available").length }}
                </p>
                <p class="text-md text-muted">Available</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-info/10">
                <UIcon name="i-lucide-calendar-check" class="size-5 text-info" />
              </div>
              <div>
                <p class="text-xl font-bold">
                  {{ timeSlots.filter((s) => s.status === "booked").length }}
                </p>
                <p class="text-md text-muted">Booked</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-amber-500/10">
                <UIcon name="i-lucide-play-circle" class="size-5 text-amber-500" />
              </div>
              <div>
                <p class="text-xl font-bold">
                  {{ timeSlots.filter((s) => s.status === "in-progress").length }}
                </p>
                <p class="text-md text-muted">In Progress</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-neutral-500/10">
                <UIcon name="i-lucide-check-square" class="size-5 text-neutral-500" />
              </div>
              <div>
                <p class="text-xl font-bold">
                  {{ timeSlots.filter((s) => s.status === "completed").length }}
                </p>
                <p class="text-sm text-muted">Completed</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-lg bg-red-500/10">
                <UIcon name="i-lucide-lock" class="size-5 text-red-500" />
              </div>
              <div>
                <p class="text-xl font-bold">
                  {{ timeSlots.filter((s) => s.status === "blocked").length }}
                </p>
                <p class="text-md text-muted">Blocked</p>
              </div>
            </div>
          </UCard>
        </div> -->

        <!-- Schedule Grid -->
        <div class="grid lg:grid-cols-3 gap-6">
          <!-- Calendar -->
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">{{ t("schedule.selectDate") }}</h2>
                <div class="flex items-center gap-1">
                  <UButton
                    icon="i-lucide-chevron-left"
                    variant="ghost"
                    color="neutral"
                    size="xs"
                    @click="changeDay(-1)"
                  />
                  <span class="text-sm font-medium">{{
                    selectedDate.toLocaleDateString(
                      locale === "id" ? "id-ID" : "en-US",
                      {
                        month: "long",
                        year: "numeric",
                      }
                    )
                  }}</span>
                  <UButton
                    icon="i-lucide-chevron-right"
                    variant="ghost"
                    color="neutral"
                    size="xs"
                    @click="changeDay(1)"
                  />
                </div>
              </div>
            </template>
            <!-- Simple custom calendar grid -->
            <div class="grid grid-cols-7 gap-1 mb-2">
              <div
                v-for="day in locale === 'id'
                  ? ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min']
                  : ['Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su']"
                :key="day"
                class="text-center text-xs font-bold text-muted py-2"
              >
                {{ day }}
              </div>
            </div>
            <div class="grid grid-cols-7 gap-1">
              <div
                v-for="(item, idx) in adminCalendarDays"
                :key="idx"
                class="w-full aspect-square flex items-center justify-center"
              >
                <button
                  v-if="item.day !== null"
                  :class="[
                    'w-full h-full rounded-lg text-sm font-medium transition-all cursor-pointer flex items-center justify-center relative',
                    selectedDate.getDate() === item.day
                      ? 'bg-primary text-white shadow-md'
                      : item.hasSlots
                      ? 'hover:bg-primary/10 text-primary-600 dark:text-primary-400 font-semibold'
                      : 'hover:bg-muted/50 text-gray-700 dark:text-gray-300',
                  ]"
                  @click="selectedDate = new Date(selectedDate.getFullYear(), selectedDate.getMonth(), item.day)"
                >
                  {{ item.day }}
                  <span
                    v-if="item.hasSlots && selectedDate.getDate() !== item.day"
                    class="absolute bottom-1 w-1.5 h-1.5 rounded-full bg-primary"
                  ></span>
                </button>
              </div>
            </div>
          </UCard>

          <!-- Time Slots List -->
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">{{ t("schedule.availableSlots") }}</h2>
                <div class="flex gap-2">
                  <UBadge
                    :label="t('common.available')"
                    color="success"
                    variant="subtle"
                  />
                  <UBadge :label="t('home.booked')" color="info" variant="subtle" />
                  <UBadge
                    :label="t('common.completed')"
                    color="neutral"
                    variant="subtle"
                  />
                  <UBadge
                    :label="t('admin.blockSlot').replace('Blokir ', '')"
                    color="error"
                    variant="subtle"
                  />
                </div>
              </div>
            </template>

            <div class="space-y-4">
              <div
                v-for="slot in filteredSlots"
                :key="slot.id"
                class="p-4 rounded-lg border border-default hover:shadow-md transition-shadow"
                :class="{
                  'border-l-4 border-l-primary': slot.status === 'available',
                  'border-l-4 border-l-info bg-info/5': slot.status === 'booked',
                  'border-l-4 border-l-amber-500': slot.status === 'in-progress',
                  'border-l-4 border-l-neutral-500 bg-neutral-500/5':
                    slot.status === 'completed',
                  'border-l-4 border-l-red-500 opacity-60': slot.status === 'blocked',
                  'cursor-pointer hover:bg-primary/5':
                    studentNameToBook && slot.status === 'available',
                }"
                @click="handleSlotClick(slot)"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-5">
                    <div class="text-center min-w-[70px]">
                      <p
                        class="text-xl font-black flex items-center justify-center gap-1"
                      >
                        {{ slot.time }}
                        <UIcon
                          v-if="slot.time >= operatingHours.nightStart"
                          name="i-lucide-moon"
                          class="size-3 text-indigo-500"
                        />
                      </p>
                      <p class="text-xs font-semibold text-muted">{{ slot.duration }}</p>
                    </div>

                    <USeparator orientation="vertical" class="h-12 hidden sm:block" />

                    <div class="space-y-1">
                      <div class="flex items-center gap-2 font-medium">
                        <UIcon name="i-lucide-car" class="size-4 text-muted" />
                        <span class="text-md">{{ slot.car }}</span>
                      </div>
                      <div class="flex items-center gap-2 mt-1">
                        <UIcon name="i-lucide-contact" class="size-4 text-muted" />
                        <span class="text-sm">{{ slot.instructor }}</span>
                        <UBadge
                          v-if="slot.time >= operatingHours.nightStart"
                          label="Night Session"
                          variant="subtle"
                          size="xs"
                          color="neutral"
                          class="ml-2"
                        />
                      </div>
                    </div>

                    <div
                      v-if="slot.student"
                      class="hidden sm:block ml-4 pl-4 border-l border-default"
                    >
                      <p
                        class="text-xs text-muted uppercase font-bold tracking-wider mb-1"
                      >
                        Booked by
                      </p>
                      <p class="font-bold flex items-center gap-2">
                        <UIcon name="i-lucide-user" class="size-4 text-info" />
                        {{ slot.student }}
                      </p>
                    </div>

                    <!-- FITUR BARU: Tombol Start/Complete Session (Hanya muncul di hari H) -->
                    <div
                      v-if="slot.status === 'booked' && isToday"
                      class="flex items-center"
                    >
                      <UButton
                        :label="t('admin.startSession')"
                        icon="i-lucide-play"
                        size="sm"
                        color="warning"
                        variant="soft"
                        @click="startSession(slot.id)"
                      />
                    </div>
                    <div
                      v-if="slot.status === 'in-progress' && isToday"
                      class="flex items-center"
                    >
                      <UButton
                        :label="t('common.completed')"
                        icon="i-lucide-check-circle"
                        size="sm"
                        color="primary"
                        variant="soft"
                        @click="completeSession(slot.id)"
                      />
                    </div>
                  </div>

                  <div class="flex items-center gap-4">
                    <UBadge
                      class="hidden sm:flex"
                      :label="
                        slot.status === 'available'
                          ? t('common.available')
                          : slot.status === 'booked'
                          ? t('home.booked')
                          : slot.status === 'in-progress'
                          ? 'In Progress'
                          : slot.status === 'completed'
                          ? t('common.completed')
                          : 'Blocked'
                      "
                      :color="
                        slot.status === 'available'
                          ? 'primary'
                          : slot.status === 'booked'
                          ? 'info'
                          : slot.status === 'in-progress'
                          ? 'warning'
                          : slot.status === 'completed'
                          ? 'neutral'
                          : 'error'
                      "
                      variant="subtle"
                    />
                    <UDropdownMenu
                      :items="[
                        [
                          {
                            label:
                              slot.status === 'blocked'
                                ? t('admin.unblockSlot')
                                : t('admin.blockSlot'),
                            icon:
                              slot.status === 'blocked'
                                ? 'i-lucide-unlock'
                                : 'i-lucide-lock',
                            onSelect: () => toggleSlotStatus(slot.id),
                            disabled:
                              slot.status === 'in-progress' ||
                              slot.status === 'completed',
                          },
                          {
                            label: t('admin.addNew').replace(
                              'Tambah Baru',
                              'Booking Manual'
                            ),
                            icon: 'i-lucide-user-plus',
                            onSelect: () => handleManualBooking(slot.id),
                            disabled: slot.status !== 'available',
                          },
                          {
                            label: t('common.edit') + ' Slot',
                            icon: 'i-lucide-pencil',
                            onSelect: () => openEditModal(slot),
                            disabled: slot.status !== 'available',
                          },
                        ],
                        [
                          {
                            label: t('admin.startSession'),
                            icon: 'i-lucide-play',
                            color: 'warning',
                            onSelect: () => startSession(slot.id),
                            disabled: slot.status !== 'booked',
                          },
                          {
                            label: t('admin.completeSession'),
                            icon: 'i-lucide-check-circle',
                            color: 'primary',
                            onSelect: () => completeSession(slot.id),
                            disabled: slot.status !== 'in-progress',
                          },
                          {
                            label: t('schedule.cancelSession'),
                            icon: 'i-lucide-user-minus',
                            color: 'neutral',
                            onSelect: () => cancelBooking(slot.id),
                            disabled: slot.status !== 'booked',
                          },
                        ],
                        [
                          {
                            label: t('common.delete') + ' Slot',
                            icon: 'i-lucide-trash',
                            color: 'error',
                            onSelect: () => deleteSlot(slot.id),
                          },
                        ],
                      ]"
                    >
                      <UButton
                        icon="i-lucide-ellipsis-vertical"
                        color="neutral"
                        variant="ghost"
                      />
                    </UDropdownMenu>
                  </div>
                </div>
              </div>

              <div
                v-if="filteredSlots.length === 0"
                class="text-center py-10 border-2 border-dashed border-default rounded-xl"
              >
                <UIcon
                  name="i-lucide-calendar-x"
                  class="size-10 text-muted mx-auto mb-3"
                />
                <p class="text-muted font-medium">
                  {{ t("blog.noArticles").replace("artikel", "slot") }}.
                </p>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <EditSlotModal
        v-model:open="showEditSlotModal"
        :initialSlot="selectedEditSlot"
        :instructors="instructors"
        :vehicles="vehicleOptionsList"
        :operatingHours="currentDayOperatingHours"
        @saved="handleEditSlot"
      />

      <ManualBookingModal
        v-model:open="showManualBookingModal"
        :slotId="selectedBookingSlotId"
        @booked="schedulesStore.fetchByDate(localDateStr)"
      />
    </template>
  </UDashboardPanel>
</template>
