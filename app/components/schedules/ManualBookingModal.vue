<script setup lang="ts">
import { ref, computed } from "vue";
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
const { operatingHours } = useSettings();

// selectedStudent holds the entire selected item { value, label, student }
const selectedStudent = ref<any>(null);
const notes = ref("");
const isSubmitting = ref(false);
const isLoadingStudents = ref(false);

const currentSlot = computed(() => {
  if (!props.slotId) return null;
  return schedulesStore.slots.find((s) => String(s.id) === String(props.slotId)) || null;
});

const isNightSlot = computed(() => {
  if (!currentSlot.value?.time) return false;
  const slotTime = currentSlot.value.time;
  const nightStart = operatingHours.value?.nightStart || "18:00";
  return slotTime >= nightStart;
});

function studentHasNightAccess(student: any): boolean {
  if (!student?.entitlements) return false;
  return student.entitlements.some((e: any) => {
    if (e.status !== "active" || e.remaining <= 0) return false;
    if (e.isNightSession) return true;
    if (e.packageName && /night|malam/i.test(e.packageName)) return true;
    return false;
  });
}

// Check if student has available sessions (active entitlements with remaining > 0)
function hasAvailableSessions(student: any) {
  if (!student?.entitlements) return false;
  return student.entitlements.some(
    (e: any) => e.status === "active" && e.remaining > 0
  );
}

// Format students for USelectMenu - filter out students with no sessions & filter for night session access if night slot
const studentItems = computed(() => {
  return studentsStore.allStudents
    .filter((s) => {
      if (!hasAvailableSessions(s)) return false;
      if (isNightSlot.value && !studentHasNightAccess(s)) return false;
      return true;
    })
    .map((s) => ({
      value: s.id,
      label: s.name,
      student: s,
    }));
});

// Extract the actual student data from selected item
const selectedStudentData = computed(() => {
  return selectedStudent.value?.student || null;
});

// Get the student ID from selected item
const selectedStudentId = computed(() => {
  return selectedStudent.value?.value || null;
});

// Helper function to get active entitlements for any student
function getActiveEntitlements(student: any) {
  if (!student?.entitlements) return [];
  return student.entitlements.filter(
    (e: any) => e.status === "active" && e.remaining > 0
  );
}

// Get the primary entitlement (newest active entitlement of the student)
function getPrimaryEntitlement(student: any) {
  const active = getActiveEntitlements(student);
  if (active.length === 0) return null;

  // Sort active entitlements by newest first (createdAt / startDate descending)
  const sorted = active.slice().sort((a: any, b: any) => {
    const timeA = a.createdAt ? new Date(a.createdAt).getTime() : (a.startDate ? new Date(a.startDate).getTime() : 0);
    const timeB = b.createdAt ? new Date(b.createdAt).getTime() : (b.startDate ? new Date(b.startDate).getTime() : 0);
    return timeB - timeA;
  });

  // If booking a night slot, prefer the newest active entitlement that supports night sessions
  if (isNightSlot.value) {
    const nightEnt = sorted.find((e: any) => e.isNightSession || (e.packageName && /night|malam/i.test(e.packageName)));
    if (nightEnt) return nightEnt;
  }

  return sorted[0];
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

// Get active entitlements for the selected student
const activeEntitlements = computed(() => {
  return getActiveEntitlements(selectedStudentData.value);
});

// Watch for modal open to fetch students
watch(
  () => props.open,
  async (isOpen) => {
    if (isOpen) {
      selectedStudent.value = null;
      notes.value = "";
      if (studentsStore.allStudents.length === 0) {
        isLoadingStudents.value = true;
        await studentsStore.fetchStudentsNoPagination();
        isLoadingStudents.value = false;
      }
    }
  },
  { immediate: true }
);

async function handleBook() {
  const studentData = selectedStudentData.value;
  const studentId = selectedStudentId.value;

  if (!studentData || !studentId) {
    toast.add({
      title: "Validation Error",
      description: "Please select a student",
      color: "error",
    });
    return;
  }

  if (isNightSlot.value && !studentHasNightAccess(studentData)) {
    toast.add({
      title: "Booking Restricted",
      description: "This time slot is in the night shift (18:00+). Selected student does not have the Night Session add-on.",
      color: "error",
    });
    return;
  }

  const entitlements = activeEntitlements.value;
  if (entitlements.length === 0) {
    toast.add({
      title: "Validation Error",
      description: "Student has no active entitlements",
      color: "error",
    });
    return;
  }

  const entitlement = getPrimaryEntitlement(studentData);
  if (!entitlement) {
    toast.add({
      title: "Booking Failed",
      description: "No valid entitlement found",
      color: "error",
    });
    return;
  }

  isSubmitting.value = true;
  try {
    const success = await schedulesStore.bookSlot(props.slotId, {
      userId: studentId,
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
            v-model="selectedStudent"
            :items="studentItems"
            :loading="isLoadingStudents"
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
                      :variant="
                        getPrimaryEntitlement(item.student)?.id === ent.id
                          ? 'solid'
                          : 'subtle'
                      "
                      size="sm"
                      class="font-mono"
                    />
                    <UIcon
                      v-if="
                        hasMultipleEntitlements(item.student) &&
                        getPrimaryEntitlement(item.student)?.id === ent.id
                      "
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
          v-if="selectedStudentData"
          class="border border-default rounded-lg p-4 bg-muted/30"
        >
          <div class="flex items-start justify-between gap-4">
            <!-- Left Side: Student Info -->
            <div class="flex-1 space-y-2">
              <!-- Row 1: Full Name -->
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-user" class="text-muted size-4" />
                <span class="font-semibold text-md">{{ selectedStudentData.name }}</span>
              </div>
              <!-- Row 2: Email -->
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-mail" class="text-muted size-4" />
                <span class="text-sm text-muted">{{ selectedStudentData.email || "N/A" }}</span>
              </div>
              <!-- Row 3: Phone -->
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-phone" class="text-muted size-4" />
                <span class="text-sm text-muted">{{
                  selectedStudentData.phone || "N/A"
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
                    :variant="
                      getPrimaryEntitlement(selectedStudentData)?.id === ent.id
                        ? 'solid'
                        : 'subtle'
                    "
                    size="lg"
                    class="font-mono"
                  />
                  <span class="text-xs text-muted">{{ ent.packageName }}</span>
                  <UIcon
                    v-if="
                      hasMultipleEntitlements(selectedStudentData) &&
                      getPrimaryEntitlement(selectedStudentData)?.id === ent.id
                    "
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

        <!-- Night Session Restriction Alert -->
        <UAlert
          v-if="isNightSlot && selectedStudentData && !studentHasNightAccess(selectedStudentData)"
          color="warning"
          variant="soft"
          icon="i-lucide-moon"
          title="Night Session Restricted"
          description="This time slot is in the night shift. Selected student does not have the Night Session add-on."
        />

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
          :disabled="!selectedStudentData || (isNightSlot && !studentHasNightAccess(selectedStudentData))"
          @click="handleBook"
        />
      </div>
    </template>
  </UModal>
</template>
