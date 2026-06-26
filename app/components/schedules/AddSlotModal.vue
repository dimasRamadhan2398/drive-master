<script setup lang="ts">
import { ref, watch } from "vue";
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { useI18n } from "#imports";

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

watch(
  () => props.open,
  (open) => {
    if (!open) {
      form.value = { time: "08:00", duration: "60", carId: "", instructorId: "" };
    }
  }
);

function handleClose() {
  emit("update:open", false);
}

function handleSave() {
  if (!form.value.time || !form.value.carId || !form.value.instructorId) {
    toast.add({ title: "Error", description: "Please fill all fields", color: "error" });
    return;
  }

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
            <UInput
              type="time"
              v-model="form.time"
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
