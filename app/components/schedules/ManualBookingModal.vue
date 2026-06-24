<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { Entitlement } from "~/services/scheduleService";
import { useToast, useStudentsStore, useSchedulesStore, useI18n } from "#imports";

const props = defineProps<{
  open: boolean;
  slotId: string;
}>();

const emit = defineEmits<{
  (e: "update:open", val: boolean): void;
  (e: "booked"): void;
}>();

const { t } = useI18n();
const toast = useToast();
const studentsStore = useStudentsStore();
const schedulesStore = useSchedulesStore();

const selectedStudentId = ref("");
const selectedEntitlementId = ref("");
const notes = ref("");
const entitlements = ref<Entitlement[]>([]);
const isLoadingEntitlements = ref(false);
const isSubmitting = ref(false);

const studentOptions = computed(() => {
  return studentsStore.students.map((s) => ({
    label: s.name,
    value: s.id,
  }));
});

const entitlementOptions = computed(() => {
  return entitlements.value.map((e) => ({
    label: `${e.packageName} (${e.remaining} sessions left)`,
    value: e.id,
  }));
});

watch(selectedStudentId, async (newId) => {
  if (newId) {
    isLoadingEntitlements.value = true;
    try {
      entitlements.value = await schedulesStore.fetchUserEntitlements(newId);
      if (entitlements.value.length > 0) {
        selectedEntitlementId.value = entitlements.value[0].id;
      } else {
        selectedEntitlementId.value = "";
      }
    } finally {
      isLoadingEntitlements.value = false;
    }
  } else {
    entitlements.value = [];
    selectedEntitlementId.value = "";
  }
});

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    selectedStudentId.value = "";
    selectedEntitlementId.value = "";
    notes.value = "";
    studentsStore.fetchStudentsNoPagination();
  }
});

async function handleBook() {
  if (!selectedStudentId.value || !selectedEntitlementId.value) {
    toast.add({
      title: "Validation Error",
      description: "Please select a student and an entitlement",
      color: "error"
    });
    return;
  }

  isSubmitting.value = true;
  try {
    const success = await schedulesStore.bookSlot(props.slotId, {
      userId: selectedStudentId.value,
      entitlementId: selectedEntitlementId.value,
      notes: notes.value
    });

    if (success) {
      toast.add({
        title: "Success",
        description: "Slot booked successfully",
        color: "success"
      });
      emit("booked");
      emit("update:open", false);
    } else {
      toast.add({
        title: "Booking Failed",
        description: "Could not book the slot. Please try again.",
        color: "error"
      });
    }
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <UModal :open="open" title="Manual Booking" @update:open="(v) => emit('update:open', v)">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Select Student" required>
          <USelect
            v-model="selectedStudentId"
            :items="studentOptions"
            placeholder="Search student..."
            class="w-full"
          />
        </UFormField>

        <UFormField label="Active Entitlement" required>
          <USelect
            v-model="selectedEntitlementId"
            :items="entitlementOptions"
            :loading="isLoadingEntitlements"
            placeholder="Select entitlement"
            :disabled="!selectedStudentId || entitlements.length === 0"
            class="w-full"
          />
          <template #hint v-if="selectedStudentId && entitlements.length === 0 && !isLoadingEntitlements">
            <span class="text-error text-xs">No active entitlements found for this student</span>
          </template>
        </UFormField>

        <UFormField label="Notes">
          <UTextarea v-model="notes" placeholder="Add booking notes..." class="w-full" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton :label="t('common.cancel')" variant="ghost" color="neutral" @click="emit('update:open', false)" />
        <UButton
          label="Book Slot"
          icon="i-lucide-calendar-check"
          color="warning"
          :loading="isSubmitting"
          @click="handleBook"
        />
      </div>
    </template>
  </UModal>
</template>
