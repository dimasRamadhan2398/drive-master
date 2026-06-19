<script setup lang="ts">
import { ref, watch } from "vue";
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";

const props = defineProps<{
  open: boolean;
  date: string;
  instructors: Array<{ label: string; value: string }>;
  vehicles: string[];
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
      car: string;
      instructor: string;
    }
  ): void;
}>();

const toast = useToast();

const form = ref({
  time: "08:00",
  duration: "60",
  car: "",
  instructor: "",
});

watch(
  () => props.open,
  (open) => {
    if (!open) {
      form.value = { time: "08:00", duration: "60", car: "", instructor: "" };
    }
  }
);

function handleClose() {
  emit("update:open", false);
}

function handleSave() {
  if (!form.value.time || !form.value.car || !form.value.instructor) {
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
    car: form.value.car,
    instructor: form.value.instructor,
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
        <UFormField label="Date" required>
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
          <UFormField label="Duration">
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
            :items="vehicles"
            v-model="form.car"
            placeholder="Select vehicle"
            class="w-full"
          />
        </UFormField>
        <UFormField label="Instructor" required>
          <USelect
            :items="instructors"
            v-model="form.instructor"
            placeholder="Select instructor"
            class="w-full"
          />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton label="Cancel" variant="ghost" color="neutral" @click="handleClose" />
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
