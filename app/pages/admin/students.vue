<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { computed, ref } from 'vue'

definePageMeta({ layout: 'admin' })

const toast = useToast()
const searchQuery = ref('')
const statusFilter = ref('all')
const showAddModal = ref(false)
const showDetailModal = ref(false)

type Student = {
  id: number
  name: string
  email: string
  phone: string
  package: string
  progress: number
  completedSessions: number
  totalSessions: number
  joinDate: string
  status: 'active' | 'pending' | 'completed'
}

const selectedStudent = ref<Student | null>(null)

const students = ref<Student[]>([
  { id: 1, name: 'John Doe', email: 'john@example.com', phone: '081234567890', package: 'Standard', progress: 40, completedSessions: 4, totalSessions: 10, joinDate: 'Mar 10, 2026', status: 'active' },
  { id: 2, name: 'Sarah Putri', email: 'sarah@example.com', phone: '081234567891', package: 'Pro', progress: 75, completedSessions: 11, totalSessions: 15, joinDate: 'Feb 20, 2026', status: 'active' },
  { id: 3, name: 'Budi Santoso', email: 'budi@example.com', phone: '081234567892', package: 'Starter', progress: 100, completedSessions: 5, totalSessions: 5, joinDate: 'Jan 15, 2026', status: 'completed' },
  { id: 4, name: 'Amanda Chen', email: 'amanda@example.com', phone: '081234567893', package: 'Standard', progress: 20, completedSessions: 2, totalSessions: 10, joinDate: 'Mar 25, 2026', status: 'active' },
  { id: 5, name: 'Ricky Wijaya', email: 'ricky@example.com', phone: '081234567894', package: 'Standard', progress: 0, completedSessions: 0, totalSessions: 10, joinDate: 'Apr 1, 2026', status: 'pending' }
])

const filteredStudents = computed(() => {
  return students.value.filter(student => {
    const matchesSearch = student.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         student.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = statusFilter.value === 'all' || student.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

function getInitials(name: string) {
  return name.split(' ').map(n => n[0]).join('')
}

function bookSessionPage(student: Student) {
  navigateTo('/admin/schedules')
}

function issueCertificatePage(student: Student) {
  navigateTo('/admin/certificates')
}

function viewStudent(student: Student) {
  selectedStudent.value = student
  showDetailModal.value = true
}

function deleteStudent(id: number) {
  students.value = students.value.filter(s => s.id !== id)
  toast.add({ title: 'Student Removed', description: 'The student has been removed from the system.', icon: 'i-lucide-trash', color: 'error' })
}

function getStatusColor(status: string) {
  if (status === 'active') return 'info'
  if (status === 'completed') return 'primary'
  return 'warning'
}

function getStatusLabel(status: string) {
  if (status === 'active') return 'Active'
  if (status === 'completed') return 'Completed'
  return 'Pending'
}

function getPackageColor(pkg: string) {
  return pkg === 'Pro' ? 'neutral' : 'neutral'
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Student Management">
        <template #right>
          <UButton icon="i-lucide-user-plus" color="warning" label="Add Student" @click="showAddModal = true" />
          <!-- Add Student Modal -->
          <UModal v-model:open="showAddModal" title="Add New Student">
            <template #body>
              <div class="space-y-4">
                <UFormField label="Full Name" required>
                  <UInput placeholder="Enter student name" color="warning" class="w-full" icon="i-lucide-user" />
                </UFormField>
                <UFormField label="Email" required>
                  <UInput type="email" placeholder="student@example.com" color="warning" class="w-full" icon="i-lucide-mail" />
                </UFormField>
                <UFormField label="Phone Number" required>
                  <UInput placeholder="081234567890" color="warning" class="w-full" icon="i-lucide-phone" />
                </UFormField>
                <UFormField label="Package" required>
                  <USelect 
                    :items="[
                      { label: 'Free Trial', value: 'free' }, 
                      { label: '6x Training Session', value: '6x' }, 
                      { label: '8x Training Session', value: '8x' },
                      { label: '10x Training Session', value: '10x' },
                      { label: '12x Training Session', value: '12x' }
                    ]" 
                    placeholder="Select package" 
                    class="w-full"
                    color="warning"
                  />
                </UFormField>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton label="Cancel" variant="ghost" color="neutral" @click="showAddModal = false" />
                <UButton label="Add Student" color="warning" icon="i-lucide-user-plus" @click="showAddModal = false" />
              </div>
            </template>
          </UModal>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <UInput v-model="searchQuery" placeholder="Search students..." color="warning" icon="i-lucide-search" class="w-64" />
        </template>
        <template #right>
          <USelect 
            v-model="statusFilter"
            :items="[
              { label: 'All Status', value: 'all' },
              { label: 'Active', value: 'active' },
              { label: 'Pending', value: 'pending' },
              { label: 'Completed', value: 'completed' }
            ]"
            class="w-40"
            color="warning"
          />
          <UButton icon="i-lucide-download" label="Export" color="neutral" variant="outline" />
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
                  <th class="text-left py-3 px-4 font-medium text-muted text-md">Student</th>
                  <th class="text-left py-3 px-4 font-medium text-muted text-md">Package</th>
                  <th class="text-left py-3 px-4 font-medium text-muted text-md">Progress</th>
                  <th class="text-left py-3 px-4 font-medium text-muted text-md">Join Date</th>
                  <th class="text-left py-3 px-4 font-medium text-muted text-md">Status</th>
                  <th class="text-right py-3 px-4 font-medium text-muted text-md">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="student in filteredStudents" :key="student.id" class="border-b border-default hover:bg-muted/30 transition-colors">
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
                    <UBadge :label="student.package" :color="getPackageColor(student.package)" variant="subtle" />
                  </td>
                  <td class="py-3 px-4">
                    <div class="w-32">
                      <div class="flex justify-between text-md mb-1">
                        <span>{{ student.completedSessions }}/{{ student.totalSessions }}</span>
                        <span>{{ student.progress }}%</span>
                      </div>
                      <UProgress :value="student.progress" size="md" />
                    </div>
                  </td>
                  <td class="py-3 px-4 text-md">{{ student.joinDate }}</td>
                  <td class="py-3 px-4">
                    <UBadge :label="getStatusLabel(student.status)" :color="getStatusColor(student.status)" variant="subtle" />
                  </td>
                  <td class="py-3 px-4 text-right">
                    <UDropdownMenu
                      :items="[
                        [
                          { label: 'View Details', icon: 'i-lucide-eye', onSelect: () => viewStudent(student) },
                          { label: 'Edit', icon: 'i-lucide-pencil' },
                          { label: 'Book Session', icon: 'i-lucide-calendar-plus', onSelect: () => bookSessionPage(student) }
                        ],
                        [{ label: 'Issue Certificate', icon: 'i-lucide-award', onSelect: () => issueCertificatePage(student) }],
                        [{ label: 'Delete', icon: 'i-lucide-trash', color: 'error', onSelect: () => deleteStudent(student.id) }]
                      ]"
                    >
                      <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" />
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
                  <UAvatar :text="getInitials(selectedStudent.name)" size="xl" />
                  <div>
                    <h3 class="text-xl font-bold">{{ selectedStudent.name }}</h3>
                    <p class="text-muted">{{ selectedStudent.email }}</p>
                    <UBadge :label="selectedStudent.package + ' Package'" color="warning" class="mt-1" />
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
                    <p class="font-medium">{{ selectedStudent.completedSessions }}/{{ selectedStudent.totalSessions }}</p>
                  </div>
                  <div>
                    <p class="text-md text-muted">Status</p>
                    <UBadge :label="getStatusLabel(selectedStudent.status)" :color="getStatusColor(selectedStudent.status)" />
                  </div>
                </div>
                <div>
                  <p class="text-md text-muted mb-2">Progress</p>
                  <UProgress :value="selectedStudent.progress" />
                  <p class="text-md text-right mt-1">{{ selectedStudent.progress }}% complete</p>
                </div>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton label="Close" variant="ghost" color="neutral" @click="showDetailModal = false" />
                <UButton label="Edit Student" color="warning" icon="i-lucide-pencil" />
              </div>
            </template>
          </UModal>
          <template #footer>
            <div class="flex items-center justify-between">
              <p class="text-md text-muted">Showing {{ filteredStudents.length }} of {{ students.length }} students</p>
              <UPagination :total="students.length" active-color="warning" :items-per-page="10" />
            </div>
          </template>
        </UCard>
      </div>
    </template>

    
    

    
    
  
  </UDashboardPanel>
</template>
