<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { computed, ref } from "vue";
import { useStudentsStore } from "~/stores/students";
import type { Student } from "~/stores/students";

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
  package: "8x",
  status: "pending" as "active" | "pending",
});

// Use store's state instead of local refs for server-side pagination
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
      title: "Not Eligible for Certificate",
      description: `${student.name} has not completed the package yet.`,
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
    title: "Student Removed",
    description: "The student has been removed from the system.",
    icon: "i-lucide-trash",
    color: "error",
  });
}

function addStudent() {
  if (!newStudent.value.email || !newStudent.value.firstName) {
    toast.add({
      title: "Error",
      description: "Nama dan email wajib diisi.",
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
    title: "Student Added",
    description: `${newStudent.value.firstName} telah ditambahkan.`,
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
    package: "8x",
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
    title: "Student Updated",
    description: `${editingStudent.value.name}'s data has been updated.`,
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
  if (status === "active") return "Active";
  if (status === "completed") return "Completed";
  return "Pending";
}

function getPackageColor(pkg: string) {
  return pkg === "8x" ? "warning" : "neutral";
}

onMounted(() => {
  studentsStore.fetchStudents();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Student Management">
        <template #right>
          <UButton
            icon="i-lucide-user-plus"
            color="warning"
            label="Add Student"
            @click="showAddModal = true"
          />
          <!-- Add Student Modal -->
          <UModal v-model:open="showAddModal" title="Add New Student">
            <template #body>
              <div class="space-y-4">
                <div class="grid grid-cols-2 gap-4">
                  <UFormField label="First Name" required>
                    <UInput
                      v-model="newStudent.firstName"
                      placeholder="First name"
                      color="warning"
                      class="w-full"
                      icon="i-lucide-user"
                    />
                  </UFormField>
                  <UFormField label="Last Name">
                    <UInput
                      v-model="newStudent.lastName"
                      placeholder="Last name"
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
                <UFormField label="Phone Number" required>
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
                  label="Cancel"
                  variant="ghost"
                  color="neutral"
                  @click="showAddModal = false"
                />
                <UButton
                  label="Create Student"
                  color="warning"
                  @click="addStudent"
                />
              </div>
            </template>
          </UModal>
          <!-- Edit Student Modal -->
          <UModal v-model:open="showEditModal" title="Edit Student">
            <template #body>
              <div v-if="editingStudent" class="space-y-4">
                <UFormField label="Full Name" required>
                  <UInput
                    v-model="editingStudent.name"
                    placeholder="Enter student name"
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
                <UFormField label="Phone Number" required>
                  <UInput
                    v-model="editingStudent.phone"
                    placeholder="081234567890"
                    color="warning"
                    class="w-full"
                    icon="i-lucide-phone"
                  />
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                  <UFormField label="Package" required>
                    <USelect
                      v-model="editingStudent.package"
                      :items="[
                        { label: '6x', value: '6x' },
                        { label: '8x', value: '8x' },
                        { label: '10x', value: '10x' },
                        { label: '12x', value: '12x' },
                      ]"
                      placeholder="Select package"
                      class="w-full"
                      color="warning"
                    />
                  </UFormField>
                  <UFormField label="Status" required>
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
                  label="Cancel"
                  variant="ghost"
                  color="neutral"
                  @click="showEditModal = false"
                />
                <UButton
                  label="Save Changes"
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
            placeholder="Search students..."
            color="warning"
            icon="i-lucide-search"
            class="w-64"
          />
        </template>
        <template #right>
          <USelect
            v-model="statusFilter"
            :items="[
              { label: 'All Status', value: 'all' },
              { label: 'Active', value: 'active' },
              { label: 'Pending', value: 'pending' },
              { label: 'Completed', value: 'completed' },
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
                    Student
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    Package
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    Progress
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    Join Date
                  </th>
                  <th
                    class="text-left py-3 px-4 font-medium text-muted text-md"
                  >
                    Status
                  </th>
                  <th
                    class="text-right py-3 px-4 font-medium text-muted text-md"
                  >
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="student in studentsList"
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
                    <UBadge
                      :label="student.package"
                      :color="getPackageColor(student.package)"
                      variant="subtle"
                    />
                  </td>
                  <td class="py-3 px-4">
                    <div class="w-32">
                      <div class="flex justify-between text-md mb-1">
                        <span
                          >{{ student.completedSessions }}/{{
                            student.totalSessions
                          }}</span
                        >
                        <span>{{ student.progress }}%</span>
                      </div>
                      <UProgress :value="student.progress" size="md" />
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
                            label: 'View Details',
                            icon: 'i-lucide-eye',
                            onSelect: () => viewStudent(student),
                          },
                          {
                            label: 'Edit',
                            icon: 'i-lucide-pencil',
                            onSelect: () => openEditModal(student),
                          },
                          {
                            label: 'Book Session',
                            icon: 'i-lucide-calendar-plus',
                            onSelect: () => bookSessionPage(student),
                          },
                        ],
                        [
                          {
                            label: 'Issue Certificate',
                            icon: 'i-lucide-award',
                            onSelect: () => issueCertificatePage(student),
                          },
                        ],
                        [
                          {
                            label: 'Delete',
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
          <UModal v-model:open="showDetailModal" title="Student Details">
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
                    <UBadge
                      :label="selectedStudent.package + ' Package'"
                      color="warning"
                      class="mt-1"
                    />
                  </div>
                </div>
                <USeparator />
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <p class="text-md text-muted">Phone</p>
                    <p class="font-medium">{{ selectedStudent.phone }}</p>
                  </div>
                  <div>
                    <p class="text-md text-muted">Join Date</p>
                    <p class="font-medium">{{ selectedStudent.joinDate }}</p>
                  </div>
                  <div>
                    <p class="text-md text-muted">Sessions</p>
                    <p class="font-medium">
                      {{ selectedStudent.completedSessions }}/{{
                        selectedStudent.totalSessions
                      }}
                    </p>
                  </div>
                  <div>
                    <p class="text-md text-muted">Status</p>
                    <UBadge
                      :label="getStatusLabel(selectedStudent.status)"
                      :color="getStatusColor(selectedStudent.status)"
                    />
                  </div>
                </div>
                <div>
                  <p class="text-md text-muted mb-2">Progress</p>
                  <UProgress :value="selectedStudent.progress" />
                  <p class="text-md text-right mt-1">
                    {{ selectedStudent.progress }}% complete
                  </p>
                </div>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton
                  label="Close"
                  variant="ghost"
                  color="neutral"
                  @click="showDetailModal = false"
                />
                <UButton
                  label="Edit Student"
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
                Showing {{ students.length }} of {{ pagination.total }} students
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
