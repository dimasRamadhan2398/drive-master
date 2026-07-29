<script setup lang="ts">
import { ref, watch, computed } from "vue";
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { getAdjustedSlotTime } from "~/composables/useSchedules";

const props = defineProps<{
  open: boolean;
  initialSlot: any | null;
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
  (e: "saved", data: { id: string; time: string; duration: string; carId: string; instructorId: string }): void;
}>();

const { t } = useI18n()
const toast = useToast();

const form = ref({ id: "", time: "08:00", duration: "60", carId: "", instructorId: "" });

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

  const targetDate = props.initialSlot?.date || "";

  let filteredOpts = allOpts;

  if (targetDate === todayStr) {
    filteredOpts = allOpts.filter((opt) => {
      const [h, m] = opt.value.split(":").map(Number);
      const optMins = (h ?? 0) * 60 + (m ?? 0);
      return optMins > currentTotalMins;
    });
  }

  if (form.value.time && !filteredOpts.some((o) => o.value === form.value.time)) {
    const [fh, fm] = form.value.time.split(":").map(Number);
    const formMins = (fh ?? 0) * 60 + (fm ?? 0);
    if (targetDate !== todayStr || formMins > currentTotalMins) {
      filteredOpts.push({ label: form.value.time, value: form.value.time });
      filteredOpts.sort((a, b) => a.value.localeCompare(b.value));
    }
  }

  return filteredOpts;
});

function checkAndAdjustTime(notify: boolean = true) {
  if (!form.value.time || !props.initialSlot?.date) return;
  const { time: adjustedTime, rounded, originalTime } = getAdjustedSlotTime(props.initialSlot.date, form.value.time);
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
  () => [props.open, props.initialSlot],
  ([open, s]) => {
    if (open && s) {
      const rawTime = s.time ?? "08:00";
      const slotDate = s.date ?? "";
      const { time: adjustedTime } = getAdjustedSlotTime(slotDate, rawTime);
      let selectedTime = adjustedTime;
      if (timeOptions.value.length > 0 && !timeOptions.value.some((o) => o.value === selectedTime)) {
        selectedTime = timeOptions.value[0]?.value || adjustedTime;
      }
      form.value.id = s.id ?? "";
      form.value.time = selectedTime;
      form.value.duration = "60";
      form.value.carId = s.carId ?? s.car ?? "";
      form.value.instructorId = s.instructorId ?? s.instructor ?? "";
    } else if (!open) {
      form.value = { id: "", time: "08:00", duration: "60", carId: "", instructorId: "" };
    }
  },
  { immediate: true, deep: true }
);

function handleClose() {
  emit("update:open", false);
}

function handleSave() {
  if (!form.value.time || !form.value.carId || !form.value.instructorId) {
    toast.add({ title: "Error", description: "Please fill all fields", color: "error" });
    return;
  }

  checkAndAdjustTime(true);

  const {
    start,
    end,
    nightStart,
    nightEnd,
    nightEnabled,
    isClosed,
  } = props.operatingHours || { start: "00:00", end: "23:59", nightStart: "00:00", nightEnd: "23:59", nightEnabled: false, isClosed: false };

  if (isClosed) {
    toast.add({ title: "Error", description: "Business is closed on this day", color: "error" });
    return;
  }

  const isDayShift = form.value.time >= start && form.value.time <= end;
  const isNightShift = nightEnabled && form.value.time >= nightStart && form.value.time <= nightEnd;
  if (!isDayShift && !isNightShift) {
    let msg = `Time must be between ${start} - ${end}`;
    if (nightEnabled) msg += ` or ${nightStart} - ${nightEnd}`;
    toast.add({ title: "Error", description: msg, color: "error" });
    return;
  }

  emit("saved", {
    id: form.value.id,
    time: form.value.time,
    duration: form.value.duration,
    carId: form.value.carId,
    instructorId: form.value.instructorId,
  });
}
</script>

<template>
  <UModal :open="open" title="Edit Time Slot" @update:open="(v) => emit('update:open', v)">
    <template #body>
      <div class="space-y-5">
        <div class="grid grid-cols-2 gap-4">
          <UFormField :label="operatingHours.isClosed ? 'Closed Today' : `Start Time`" required>
            <template #hint>
              <div class="text-[10px] text-muted-foreground flex flex-col items-end">
                <span>Day: {{ operatingHours.start }}-{{ operatingHours.end }}</span>
                <span v-if="operatingHours.nightEnabled">Night: {{ operatingHours.nightStart }}-{{ operatingHours.nightEnd }}</span>
              </div>
            </template>
            <USelect v-model="form.time" :items="timeOptions" @change="checkAndAdjustTime(true)" :disabled="operatingHours.isClosed" class="w-full" />
          </UFormField>
          <UFormField :label="t('admin.package.duration')">
            <USelect :items="[{ label: '60 minutes', value: '60' }]" v-model="form.duration" disabled class="w-full" />
          </UFormField>
        </div>
        <UFormField label="Vehicle" required>
          <USelect v-model="form.carId" :items="vehicles" placeholder="Select vehicle" class="w-full" />
        </UFormField>
        <UFormField :label="t('dashboard.instructor')" required>
          <USelect :items="instructors" v-model="form.instructorId" placeholder="Select instructor" class="w-full" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton :label="t('common.cancel')" variant="ghost" color="neutral" @click="handleClose" />
        <UButton :label="t('admin.saveChanges')" icon="i-lucide-save" @click="handleSave" color="warning" />
      </div>
    </template>
  </UModal>
</template>
