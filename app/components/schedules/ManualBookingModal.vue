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
const notes = ref("");
const isSubmitting = ref(false);

// Format students for USelectMenu
const studentItems = computed(() => {
  return studentsStore.students.map((s: any) => ({
    value: s.id,
    label: `${s.firstName} ${s.lastName}`,
    student: s,
  }));
});

// Helper function to get active entitlements for any student
function getActiveEntitlements(student: any) {
  if (!student?.entitlements) return [];
  return student.entitlements.filter((e: any) => e.status === "active" && e.remaining > 0);
}

// Get the primary entitlement (first one or the one with most remaining sessions)
function getPrimaryEntitlement(student: any) {
  const active = getActiveEntitlements(student);
  if (active.length === 0) return null;
  if (active.length === 1) return active[0];
  // Return the entitlement with the most remaining sessions
  return active.reduce((prev: any, curr: any) =>
    curr.remaining > prev.remaining ? curr : prev
  );
}

// Check if student has multiple active entitlements
function hasMultipleEntitlements(student: any) {
  return getActiveEntitlements(student).length > 1;
}

// Helper function to determine badge color based on usage
function getEntitlementColor(ent: any): "success" | "warning" | "error" {
  if (!ent || !ent.totalSessions) return "error";
  const usagePercent = ent.usedSessions / ent.totalSessions;
  return usagePercent >= 0.8 ? "warning" : "success";
}

// Selected student data for the card display
const selectedStudent = computed(() => {
  if (!selectedStudentId.value) return null;
  return studentsStore.students.find((s: any) => s.id === selectedStudentId.value);
});

// Get active entitlements for the selected student
const activeEntitlements = computed(() => {
  return getActiveEntitlements(selectedStudent.value);
});

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      selectedStudentId.value = "";
      notes.value = "";
      studentsStore.fetchStudentsNoPagination();
    }
  }
);

async function handleBook() {
  if (!selectedStudentId.value || activeEntitlements.value.length === 0) {
    toast.add({
      title: "Validation Error",
      description: "Please select a student with active entitlements",
      color: "error",
    });
    return;
  }

  isSubmitting.value = true;
  try {
    // Use the primary entitlement (the one with most remaining sessions)
    const entitlement = getPrimaryEntitlement(selectedStudent.value);
    if (!entitlement) {
      toast.add({
        title: "Booking Failed",
        description: "No valid entitlement found",
        color: "error",
      });
      return;
    }

    const success = await schedulesStore.bookSlot(props.slotId, {
      userId: selectedStudentId.value,
      entitlementId: entitlement.id,
      notes: notes.value,
    });

    if (success) {
      toast.add({
        title: "Success",
        description: "Slot booked successfully",
        color: "success",
      });
      emit("booked");
      emit("update:open", false);
    } else {
      toast.add({
        title: "Booking Failed",
        description: "Could not book the slot. Please try again.",
        color: "error",
      });
    }
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <UModal
    :open="open"
    title="Manual Booking"
    class="max-w-2xl"
    @update:open="(v : any) => emit('update:open', v)"
  >
    <template #body>
      <div class="space-y-4">
        <!-- Student Selection Dropdown -->
        <UFormField label="Select Student" required>
          <USelectMenu
            v-model="selectedStudentId"
            :items="studentItems"
            placeholder="Select student..."
            icon="i-lucide-user"
            color="neutral"
            class="w-full"
          >
            <template #item="{ item }">
              <div class="flex items-start justify-between gap-4 w-full py-1">
                <!-- Left: Student Info -->
                <div class="flex-1 min-w-0 space-y-1">
                  <div class="flex items-center gap-2">
                    <span class="font-medium truncate">{{ item.student.name }}</span>
                    <UBadge
                      v-if="item.student.memberProfile?.nickname"
                      :label="item.student.memberProfile.nickname"
                      variant="subtle"
                      color="neutral"
                      size="sm"
                    />
                  </div>
                  <div class="text-xs text-muted truncate">
                    {{ item.student.email || "No email" }}
                  </div>
                  <div class="text-xs text-muted truncate">
                    {{ item.student.phone || "No phone" }}
                  </div>
                </div>

                <!-- Right: Package Info -->
                <div class="flex flex-col items-end gap-1 shrink-0">
                  <div
                    v-for="ent in getActiveEntitlements(item.student).slice(0, 2)"
                    :key="ent.id"
                    class="flex items-center gap-1"
                  >
                    <UBadge
                      :label="`${ent.remaining}/${ent.totalSessions}`"
                      :color="
                        ent.remaining > 3
                          ? 'success'
                          : ent.remaining > 1
                          ? 'warning'
                          : 'error'
                      "
                      :variant="getPrimaryEntitlement(item.student)?.id === ent.id ? 'solid' : 'subtle'"
                      size="sm"
                      class="font-mono"
                    />
                    <UIcon
                      v-if="hasMultipleEntitlements(item.student) && getPrimaryEntitlement(item.student)?.id === ent.id"
                      name="i-lucide-star"
                      class="size-3 text-warning"
                    />
                  </div>
                  <span
                    v-if="getActiveEntitlements(item.student).length === 0"
                    class="text-xs text-error"
                  >
                    No sessions
                  </span>
                </div>
              </div>
            </template>
          </USelectMenu>
        </UFormField>

        <!-- Student Info Card -->
        <div
          v-if="selectedStudent"
          class="border border-default rounded-lg p-4 bg-muted/30"
        >
          <div class="flex items-start justify-between gap-4">
            <!-- Left Side: Student Info -->
            <div class="flex-1 space-y-2">
              <!-- Row 1: Full Name -->
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-user" class="text-muted size-4" />
                <span class="font-semibold text-md">{{ selectedStudent.name }} </span>
                <!-- <UBadge
                  v-if="selectedStudent.memberProfile?.nickname"
                  :label="selectedStudent.memberProfile.nickname"
                  variant="subtle"
                  color="neutral"
                  size="sm"
                /> -->
              </div>
              <!-- Row 2: Email -->
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-mail" class="text-muted size-4" />
                <span class="text-sm text-muted">{{
                  selectedStudent.email || "N/A"
                }}</span>
              </div>
              <!-- Row 3: Phone -->
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-phone" class="text-muted size-4" />
                <span class="text-sm text-muted">{{
                  selectedStudent.phoneNumber || "N/A"
                }}</span>
              </div>
            </div>

            <!-- Right Side: Package Info -->
            <div class="flex flex-col gap-2">
              <template v-for="ent in activeEntitlements" :key="ent.id">
                <div class="flex items-center gap-2">
                  <UBadge
                    :label="`${ent.usedSessions}/${ent.totalSessions}`"
                    :color="getEntitlementColor(ent)"
                    :variant="getPrimaryEntitlement(selectedStudent)?.id === ent.id ? 'solid' : 'subtle'"
                    size="lg"
                    class="font-mono"
                  />
                  <span class="text-xs text-muted">{{ ent.packageName }}</span>
                  <UIcon
                    v-if="hasMultipleEntitlements(selectedStudent) && getPrimaryEntitlement(selectedStudent)?.id === ent.id"
                    name="i-lucide-star"
                    class="size-3 text-warning"
                  />
                </div>
              </template>
              <div v-if="activeEntitlements.length === 0" class="text-xs text-error">
                No active entitlements
              </div>
            </div>
          </div>
        </div>

        <UFormField label="Notes">
          <UTextarea v-model="notes" placeholder="Add booking notes..." class="w-full" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton
          :label="t('common.cancel')"
          variant="ghost"
          color="neutral"
          @click="emit('update:open', false)"
        />
        <UButton
          label="Book Slot"
          icon="i-lucide-calendar-check"
          color="warning"
          :loading="isSubmitting"
          :disabled="!selectedStudentId || activeEntitlements.length === 0"
          @click="handleBook"
        />
      </div>
    </template>
  </UModal>
</template>
