<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { computed, ref } from "vue";
import { useStudentsStore } from "~/stores/students";
import type { Student } from "~/stores/students";

const { t } = useI18n()
definePageMeta({ layout: "admin" });

const toast = useToast();
const studentsStore = useStudentsStore();

const showAddModal = ref(false);
const showDetailModal = ref(false);
const showEditModal = ref(false);

const selectedStudent = ref<Student | null>(null);
const editingStudent = ref<Student | null>(null);

const newStudent = ref({
  firstName: "",
  lastName: "",
  email: "",
  phoneNumber: "",
  status: "pending" as "active" | "pending",
});

// Use store's state for server-side pagination
const searchQuery = computed({
  get: () => studentsStore.searchQuery,
  set: (value) => studentsStore.setSearchQuery(value),
});

const statusFilter = computed({
  get: () => studentsStore.statusFilter,
  set: (value) => studentsStore.setStatusFilter(value),
});

// Get students and pagination from store
const students = computed(() => studentsStore.students);
const pagination = computed(() => studentsStore.pagination);

// For display purposes - students list is now directly from store
const studentsList = computed(() => studentsStore.students);

const studentsWithProgress = computed(() => {
  return studentsList.value.map((student) => ({
    ...student,
    _progress: getActiveEntitlementProgress(student),
    _sessions: getActiveEntitlementSessions(student),
  }));
});

function getActiveEntitlementProgress(student: Student) {
  if (!student || !student.entitlements) return 0;
  const active = student.entitlements.find((e) => e.status === "active");
  if (!active || !active.totalSessions) return 0;

  const completed = Number(active.usedSessions ?? 0);
  const total = Number(active.totalSessions);

  if (total <= 0 || isNaN(completed) || isNaN(total)) return 0;

  const progress = Math.round((completed / total) * 100);
  return isNaN(progress) || !isFinite(progress) ? 0 : Math.min(100, Math.max(0, progress));
}

function getActiveEntitlementSessions(student: Student) {
  if (!student || !student.entitlements) return { completed: 0, total: 0 };
  const active = student.entitlements.find((e) => e.status === "active");
  if (!active) return { completed: 0, total: 0 };
  return { 
    completed: active.usedSessions ?? 0, 
    total: active.totalSessions ?? 0 
  };
}

function getInitials(name: string) {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("");
}

function bookSessionPage(student: Student) {
  navigateTo(
    `/admin/schedules?studentName=${encodeURIComponent(student.name)}`,
  );
}

function issueCertificatePage(student: Student) {
  if (student.status !== "completed") {
    toast.add({
      title: t('admin.notEligibleCert'),
      description: t('admin.notEligibleCertDesc', { name: student.name }),
      color: "warning",
      icon: "i-lucide-info",
    });
    return;
  }
  navigateTo(`/admin/certificates?issueFor=${student.id}`);
}

function viewStudent(student: Student) {
  selectedStudent.value = student;
  showDetailModal.value = true;
}

function openEditModal(student: Student) {
  editingStudent.value = JSON.parse(JSON.stringify(student));
  showDetailModal.value = false; // Close detail modal if it was open
  showEditModal.value = true;
}

async function deleteStudent(id: string) {
  await studentsStore.deleteStudent(id);
  toast.add({
    title: t('admin.studentRemoved'),
    description: t('admin.studentRemovedDesc'),
    icon: "i-lucide-trash",
    color: "error",
  });
}

function addStudent() {
  if (!newStudent.value.email || !newStudent.value.firstName) {
    toast.add({
      title: t('common.error'),
      description: t('admin.nameEmailRequired'),
      color: "error",
    });
    return;
  }

  studentsStore.addStudent({
    email: newStudent.value.email,
    firstName: newStudent.value.firstName,
    lastName: newStudent.value.lastName,
    phoneNumber: newStudent.value.phoneNumber,
    password: "TempPassword123", // Default password for new students
  });

  toast.add({
    title: t('admin.studentAdded'),
    description: t('admin.studentAddedDesc', { name: newStudent.value.firstName }),
    icon: "i-lucide-check-circle",
    color: "success",
  });

  showAddModal.value = false;
  // Reset form
  newStudent.value = {
    firstName: "",
    lastName: "",
    email: "",
    phoneNumber: "",
    status: "pending",
  };
}

function saveEditedStudent() {
  if (!editingStudent.value) return;

  // Convert Student to UpdateStudentData format
  const nameParts = editingStudent.value.name.split(" ");
  const updateData = {
    firstName: nameParts[0] || "",
    lastName: nameParts.slice(1).join(" ") || "",
    phoneNumber: editingStudent.value.phone,
  };

  studentsStore.updateStudent(editingStudent.value.id, updateData);

  toast.add({
    title: t('admin.studentUpdated'),
    description: t('admin.studentUpdatedDesc', { name: editingStudent.value.name }),
    icon: "i-lucide-check-circle",
    color: "success",
  });

  showEditModal.value = false;
  editingStudent.value = null;
}

function getStatusColor(status: string) {
  if (status === "active") return "info";
  if (status === "completed") return "primary";
  return "warning";
}

function getStatusLabel(status: string) {
  if (status === "active") return t('billing.active');
  if (status === "completed") return t('common.completed');
  return t('common.pending');
}

function getPackageBadgeClass(pkg: string) {
  const pkgLower = pkg.toLowerCase();
  if (pkgLower.includes("gold") || pkgLower.includes("10x"))
    return "bg-gradient-to-r from-yellow-400 to-amber-400 text-amber-950 shadow-[0_0_12px_rgba(250,205,78,0.6)] border border-yellow-300/50 dark:from-amber-400 dark:to-yellow-400 dark:text-amber-950 dark:shadow-[0_0_12px_rgba(250,205,78,0.5)]";
  if (pkgLower.includes("platinum") || pkgLower.includes("12x"))
    return "bg-gradient-to-r from-slate-300 to-slate-400 text-slate-800 shadow-[0_0_12px_rgba(148,163,184,0.4)] border border-slate-200/50 dark:from-slate-500 dark:to-slate-400 dark:text-slate-100 dark:shadow-[0_0_12px_rgba(148,163,184,0.3)]";
  if (pkgLower.includes("silver") || pkgLower.includes("8x"))
    return "bg-gradient-to-r from-gray-200 to-gray-300 text-gray-700 shadow-[0_0_8px_rgba(209,213,219,0.5)] border border-gray-300/50 dark:from-gray-600 dark:to-gray-500 dark:text-gray-100 dark:shadow-[0_0_8px_rgba(156,163,175,0.4)]";
  if (pkgLower.includes("bronze") || pkgLower.includes("6x"))
    return "bg-gradient-to-r from-orange-300 to-orange-400 text-orange-900 shadow-[0_0_12px_rgba(251,146,60,0.5)] border border-orange-200/50 dark:from-orange-600 dark:to-orange-500 dark:text-orange-100 dark:shadow-[0_0_12px_rgba(249,115,22,0.4)]";
  if (pkgLower.includes("basic"))
    return "bg-gradient-to-r from-green-300 to-green-400 text-green-900 shadow-[0_0_10px_rgba(134,239,172,0.5)] border border-green-200/50 dark:from-green-600 dark:to-green-500 dark:text-green-100 dark:shadow-[0_0_10px_rgba(74,222,128,0.4)]";
  return "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200 border border-gray-200 dark:border-gray-600";
}

onMounted(() => {
  studentsStore.fetchStudents();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.manageStudents')">
        <template #right>
          <UButton
            icon="i-lucide-user-plus"
            color="warning"
            :label="t('admin.addNew')"
            @click="() => {
              showAddModal = true
            }"
          />
          <!-- Add Student Modal -->
          <UModal v-model:open="showAddModal" :title="t('admin.addNew')">
            <template #body>
              <div class="space-y-4">
                <div class="grid grid-cols-2 gap-4">
                  <UFormField :label="t('register.form.firstName')" required>
                    <UInput
                      v-model="newStudent.firstName"
                      :placeholder="t('register.form.firstName')"
                      color="warning"
                      class="w-full"
                      icon="i-lucide-user"
                    />
                  </UFormField>
                  <UFormField :label="t('register.form.lastName')">
                    <UInput
                      v-model="newStudent.lastName"
                      :placeholder="t('register.form.lastName')"
                      color="warning"
                      class="w-full"
                    />
                  </UFormField>
                </div>
                <UFormField label="Email" required>
                  <UInput
                    v-model="newStudent.email"
                    type="email"
                    placeholder="student@example.com"
                    color="warning"
                    class="w-full"
                    icon="i-lucide-mail"
                  />
                </UFormField>
                <UFormField :label="t('profile.phone')" required>
                  <UInput
                    v-model="newStudent.phoneNumber"
                    placeholder="+6281234567890"
                    color="warning"
                    class="w-full"
                    icon="i-lucide-phone"
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
                  @click="() => {
                    showAddModal = false
                  }"
                />
                <UButton
                  :label="t('admin.createStudent')"
                  color="warning"
                  @click="addStudent"
                />
              </div>
            </template>
          </UModal>
          <!-- Edit Student Modal -->
          <UModal v-model:open="showEditModal" :title="t('common.edit')">
            <template #body>
              <div v-if="editingStudent" class="space-y-4">
                <UFormField :label="t('profile.fullName')" required>
                  <UInput
                    v-model="editingStudent.name"
                    :placeholder="t('admin.enterStudentName')"
                    color="warning"
                    class="w-full"
                    icon="i-lucide-user"
                  />
                </UFormField>
                <UFormField label="Email" required>
                  <UInput
                    v-model="editingStudent.email"
                    type="email"
                    placeholder="student@example.com"
                    color="warning"
                    class="w-full"
                    icon="i-lucide-mail"
                  />
                </UFormField>
                <UFormField :label="t('profile.phone')" required>
                  <UInput
                    v-model="editingStudent.phone"
                    placeholder="081234567890"
                    color="warning"
                    class="w-full"
                    icon="i-lucide-phone"
                  />
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                  <UFormField :label="t('billing.package')" required>
                    <span
                      :class="getPackageBadgeClass(editingStudent.package)"
                      class="inline-flex items-center px-2.5 py-1.5 rounded-md text-xs font-medium"
                    >
                      {{ editingStudent.package }}
                    </span>
                  </UFormField>
                  <UFormField :label="t('billing.status')" required>
                    <USelect
                      v-model="editingStudent.status"
                      :items="
                        ['active', 'pending', 'completed'].map((s) => ({
                          label: getStatusLabel(s),
                          value: s,
                        }))
                      "
                      class="w-full"
                      color="warning"
                    />
                  </UFormField>
                </div>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton
                  :label="t('common.cancel')"
                  variant="ghost"
                  color="neutral"
                  @click="showEditModal = false"
                />
                <UButton
                  :label="t('admin.saveChanges')"
                  color="warning"
                  @click="saveEditedStudent"
                />
              </div>
            </template>
          </UModal>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <UInput
            v-model="searchQuery"
            :placeholder="t('common.search') + '...'"
            color="warning"
            icon="i-lucide-search"
            class="w-64"
          />
        </template>
        <template #right>
          <USelect
            v-model="statusFilter"
            :items="[
              { label: t('admin.allStatus'), value: 'all' },
              { label: t('billing.active'), value: 'active' },
              { label: t('common.pending'), value: 'pending' },
              { label: t('common.completed'), value: 'completed' },
            ]"
            class="w-40"
            color="warning"
          />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <div class="p-6">
        <UCard>
          <!-- Custom Table -->
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-default">
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    {{ t('admin.students') }}
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    {{ t('billing.package') }}
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    {{ t('common.progress') }}
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    {{ t('admin.joinDate') }}
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    {{ t('billing.status') }}
                  </th>
                  <th
                    class="text-right py-3 px-4 font-medium text-muted text-md"
                  >
                    {{ t('admin.actions') }}
                  </th>
                </tr>
              </thead>
              <tbody>
                 <tr
                  v-for="student in studentsWithProgress"
                  :key="student.id"
                  class="border-b border-default hover:bg-muted/30 transition-colors"
                >
                  <td class="py-3 px-4">
                    <div class="flex items-center gap-3">
                      <UAvatar :text="getInitials(student.name)" size="md" />
                      <div>
                        <p class="font-medium">{{ student.name }}</p>
                        <p class="text-md text-muted">{{ student.email }}</p>
                      </div>
                    </div>
                  </td>
                  <td class="py-3 px-4">
                    <span
                      :class="getPackageBadgeClass(student.package)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ student.package }}
                    </span>
                  </td>
                  <td class="py-3 px-4">
                    <div class="w-32">
                      <div class="flex justify-between text-md mb-1">
                        <span
                          >{{ student._sessions.completed }}/{{
                            student._sessions.total
                          }}</span
                        >
                        <span>{{ student._progress }}%</span>
                      </div>
                      <div class="w-full bg-gray-100 dark:bg-gray-800 rounded-full h-2 overflow-hidden border border-gray-200 dark:border-gray-700">
                        <div 
                          class="bg-primary h-full rounded-full transition-all duration-300"
                          :style="{ width: `${student._progress}%` }"
                        ></div>
                      </div>
                    </div>
                  </td>
                  <td class="py-3 px-4 text-md">{{ student.joinDate }}</td>
                  <td class="py-3 px-4">
                    <UBadge
                      :label="getStatusLabel(student.status)"
                      :color="getStatusColor(student.status)"
                      variant="subtle"
                    />
                  </td>
                  <td class="py-3 px-4 text-right">
                    <UDropdownMenu
                      :items="[
                        [
                          {
                            label: t('dashboard.viewDetails'),
                            icon: 'i-lucide-eye',
                            onSelect: () => viewStudent(student),
                          },
                          {
                            label: t('common.edit'),
                            icon: 'i-lucide-pencil',
                            onSelect: () => openEditModal(student),
                          },
                          {
                            label: t('dashboard.bookSession'),
                            icon: 'i-lucide-calendar-plus',
                            onSelect: () => bookSessionPage(student),
                          },
                        ],
                        [
                          {
                            label: t('admin.certificates'),
                            icon: 'i-lucide-award',
                            onSelect: () => issueCertificatePage(student),
                          },
                        ],
                        [
                          {
                            label: t('common.delete'),
                            icon: 'i-lucide-trash',
                            color: 'error',
                            onSelect: () => deleteStudent(student.id),
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
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <!-- Student Detail Modal -->
          <UModal v-model:open="showDetailModal" :title="t('dashboard.sessionDetail')">
            <template #body>
              <div v-if="selectedStudent" class="space-y-4">
                <div class="flex items-center gap-4">
                  <UAvatar
                    :text="getInitials(selectedStudent.name)"
                    size="xl"
                  />
                  <div>
                    <h3 class="text-xl font-bold">
                      {{ selectedStudent.name }}
                    </h3>
                    <p class="text-muted">{{ selectedStudent.email }}</p>
                    <span
                      :class="getPackageBadgeClass(selectedStudent.package)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium mt-1"
                    >
                      {{ selectedStudent.package }} {{ t('billing.package') }}
                    </span>
                  </div>
                </div>
                <USeparator />
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <p class="text-md text-muted">{{ t('profile.phone') }}</p>
                    <p class="font-medium">{{ selectedStudent.phone }}</p>
                  </div>
                  <div>
                    <p class="text-md text-muted">{{ t('admin.joinDate') }}</p>
                    <p class="font-medium">{{ selectedStudent.joinDate }}</p>
                  </div>
                  <div>
                    <p class="text-md text-muted">{{ t('billing.sessions') }}</p>
                    <p class="font-medium">
                      {{ selectedStudent.completedSessions }}/{{
                        selectedStudent.totalSessions
                      }}
                    </p>
                  </div>
                  <div>
                    <p class="text-md text-muted">{{ t('billing.status') }}</p>
                    <UBadge
                      :label="getStatusLabel(selectedStudent.status)"
                      :color="getStatusColor(selectedStudent.status)"
                    />
                  </div>
                </div>
                <div>
                  <p class="text-md text-muted mb-2">{{ t('common.progress') }}</p>
                  <div class="w-full bg-gray-100 dark:bg-gray-800 rounded-full h-2 overflow-hidden border border-gray-200 dark:border-gray-700">
                    <div 
                      class="bg-primary h-full rounded-full transition-all duration-300"
                      :style="{ width: `${selectedStudent.progress}%` }"
                    ></div>
                  </div>
                  <p class="text-md text-right mt-1">
                    {{ selectedStudent.progress }}% {{ t('common.completed').toLowerCase() }}
                  </p>
                </div>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton
                  :label="t('dashboard.close')"
                  variant="ghost"
                  color="neutral"
                  @click="showDetailModal = false"
                />
                <UButton
                  :label="t('common.edit') + ' ' + t('admin.students')"
                  color="warning"
                  icon="i-lucide-pencil"
                  @click="openEditModal(selectedStudent!)"
                />
              </div>
            </template>
          </UModal>
          <template #footer>
            <div class="flex items-center justify-between">
              <p class="text-md text-muted">
                {{ t('admin.showing', { count: students.length, total: pagination.total }) }}
              </p>
              <UPagination
                v-model="studentsStore.pagination.page"
                :total="pagination.total"
                :items-per-page="pagination.limit"
                active-color="warning"
                @update:page="(page) => studentsStore.setPage(page)"
              />
            </div>
          </template>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
