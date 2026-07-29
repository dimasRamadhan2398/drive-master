<script setup lang="ts">
import { ref, watch, computed } from "vue";
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { useI18n } from "#imports";
import { getAdjustedSlotTime } from "~/composables/useSchedules";

const props = defineProps<{
  open: boolean;
  date: string;
  instructors: Array<{ label: string; value: string }>;
  vehicles: Array<{ label: string; value: string }>;
  operatingHours: {
    start: string;
    end: string;
    nightStart: string;
    nightEnd: string;
    nightEnabled: boolean;
    isClosed: boolean;
  };
}>();

const emit = defineEmits<{
  (e: "update:open", val: boolean): void;
  (
    e: "saved",
    data: {
      date: string;
      time: string;
      duration: string;
      carId: string;
      instructorId: string;
    }
  ): void;
}>();

const { t } = useI18n();
const toast = useToast();

const form = ref({
  time: "08:00",
  duration: "60",
  carId: "",
  instructorId: "",
});

const timeOptions = computed(() => {
  const allOpts: Array<{ label: string; value: string }> = [];
  for (let h = 6; h <= 22; h++) {
    const hh = String(h).padStart(2, "0");
    allOpts.push({ label: `${hh}:00`, value: `${hh}:00` });
    if (h < 22) {
      allOpts.push({ label: `${hh}:30`, value: `${hh}:30` });
    }
  }

  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  const todayStr = `${year}-${month}-${day}`;

  const currentH = now.getHours();
  const currentM = now.getMinutes();
  const currentTotalMins = currentH * 60 + currentM;

  let filteredOpts = allOpts;

  if (props.date === todayStr) {
    filteredOpts = allOpts.filter((opt) => {
      const [h, m] = opt.value.split(":").map(Number);
      const optMins = (h ?? 0) * 60 + (m ?? 0);
      return optMins > currentTotalMins;
    });
  }

  if (form.value.time && !filteredOpts.some((o) => o.value === form.value.time)) {
    const [fh, fm] = form.value.time.split(":").map(Number);
    const formMins = (fh ?? 0) * 60 + (fm ?? 0);
    if (props.date !== todayStr || formMins > currentTotalMins) {
      filteredOpts.push({ label: form.value.time, value: form.value.time });
      filteredOpts.sort((a, b) => a.value.localeCompare(b.value));
    }
  }

  return filteredOpts;
});

function checkAndAdjustTime(notify: boolean = true) {
  if (!form.value.time) return;
  const { time: adjustedTime, rounded, originalTime } = getAdjustedSlotTime(props.date, form.value.time);
  if (rounded && adjustedTime !== originalTime) {
    form.value.time = adjustedTime;
    if (notify) {
      toast.add({
        title: "Time Slot Rounded",
        description: `Selected time (${originalTime}) has already passed today. Automatically rounded to ${adjustedTime}.`,
        color: "info",
        icon: "i-lucide-clock",
      });
    }
  }
}

watch(
  () => [props.open, props.date],
  ([open]) => {
    if (open) {
      const { time: adjustedTime } = getAdjustedSlotTime(props.date, "08:00");
      let selectedTime = adjustedTime;
      if (timeOptions.value.length > 0 && !timeOptions.value.some((o) => o.value === selectedTime)) {
        selectedTime = timeOptions.value[0]?.value || adjustedTime;
      }
      form.value = { time: selectedTime, duration: "60", carId: "", instructorId: "" };
    } else {
      form.value = { time: "08:00", duration: "60", carId: "", instructorId: "" };
    }
  },
  { immediate: true }
);

function handleClose() {
  emit("update:open", false);
}

function handleSave() {
  if (!form.value.time || !form.value.carId || !form.value.instructorId) {
    toast.add({ title: "Error", description: "Please fill all fields", color: "error" });
    return;
  }

  // Ensure time hasn't passed before saving
  checkAndAdjustTime(true);

  const {
    start,
    end,
    nightStart,
    nightEnd,
    nightEnabled,
    isClosed,
  } = props.operatingHours || {
    start: "00:00",
    end: "23:59",
    nightStart: "00:00",
    nightEnd: "23:59",
    nightEnabled: false,
    isClosed: false,
  };

  if (isClosed) {
    toast.add({
      title: "Error",
      description: "Business is closed on this day",
      color: "error",
    });
    return;
  }

  const isDayShift = form.value.time >= start && form.value.time <= end;
  const isNightShift =
    nightEnabled && form.value.time >= nightStart && form.value.time <= nightEnd;
  if (!isDayShift && !isNightShift) {
    let msg = `Time must be between ${start} - ${end}`;
    if (nightEnabled) msg += ` or ${nightStart} - ${nightEnd}`;
    toast.add({ title: "Error", description: msg, color: "error" });
    return;
  }

  emit("saved", {
    date: props.date,
    time: form.value.time,
    duration: form.value.duration,
    carId: form.value.carId,
    instructorId: form.value.instructorId,
  });
}
</script>

<template>
  <UModal
    :open="open"
    title="Add New Time Slot"
    @update:open="(v) => emit('update:open', v)"
  >
    <template #body>
      <div class="space-y-5">
        <UFormField :label="t('dashboard.date')" required>
          <UInput type="date" :model-value="date" disabled class="w-full" />
        </UFormField>
        <div class="grid grid-cols-2 gap-4">
          <UFormField
            :label="operatingHours.isClosed ? 'Closed Today' : `Start Time`"
            required
          >
            <template #hint>
              <div class="text-[10px] text-muted-foreground flex flex-col items-end">
                <span>Day: {{ operatingHours.start }}-{{ operatingHours.end }}</span>
                <span v-if="operatingHours.nightEnabled"
                  >Night: {{ operatingHours.nightStart }}-{{
                    operatingHours.nightEnd
                  }}</span
                >
              </div>
            </template>
            <USelect
              v-model="form.time"
              :items="timeOptions"
              @change="checkAndAdjustTime(true)"
              :disabled="operatingHours.isClosed"
              class="w-full"
            />
          </UFormField>
          <UFormField :label="t('admin.package.duration')">
            <template #hint>
              <div class="text-[10px] flex flex-col items-end text-transparent">
                <span>Day: {{ operatingHours.start }}-{{ operatingHours.end }}</span>
                <span v-if="operatingHours.nightEnabled"
                  >Night: {{ operatingHours.nightStart }}-{{
                    operatingHours.nightEnd
                  }}</span
                >
              </div>
            </template>
            <USelect
              :items="[{ label: '60 minutes', value: '60' }]"
              v-model="form.duration"
              disabled
              class="w-full"
            />
          </UFormField>
        </div>
        <UFormField label="Vehicle" required>
          <USelect
            v-model="form.carId"
            :items="vehicles"
            placeholder="Select vehicle"
            class="w-full"
          />
        </UFormField>
        <UFormField :label="t('dashboard.instructor')" required>
          <USelect
            :items="instructors"
            v-model="form.instructorId"
            placeholder="Select instructor"
            class="w-full"
          />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton
          :label="t('common.cancel')"
          variant="ghost"
          color="neutral"
          @click="handleClose"
        />
        <UButton
          label="Create Slot"
          icon="i-lucide-plus"
          @click="handleSave"
          color="warning"
        />
      </div>
    </template>
  </UModal>
</template>
