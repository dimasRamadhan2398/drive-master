<script setup lang="ts">
import { computed, h, ref, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'

definePageMeta({ layout: 'admin' })

const toast = useToast()
const searchQuery = ref('')
const showIssueModal = ref(false)

type EligibleStudent = {
  id: number
  name: string
  email: string
  package: string
  completedDate: string
}

type Certificate = {
  id: string
  studentName: string
  email: string
  package: string
  issueDate: string
  status: 'issued' | 'revoked'
}

const eligibleStudents = ref<EligibleStudent[]>([
  { id: 1, name: 'Budi Santoso', email: 'budi@example.com', package: '6x', completedDate: 'Mar 28, 2026' },
  { id: 2, name: 'Maria Garcia', email: 'maria@example.com', package: '8x', completedDate: 'Mar 30, 2026' }
])

const issuedCertificates = ref<Certificate[]>([
  { id: 'EVDA-2026-001230', studentName: 'Ahmad Rizky', email: 'ahmad@example.com', package: '12x', issueDate: 'Mar 20, 2026', status: 'issued' },
  { id: 'EVDA-2026-001229', studentName: 'Linda Wijaya', email: 'linda@example.com', package: '10x', issueDate: 'Mar 18, 2026', status: 'issued' },
  { id: 'EVDA-2026-001228', studentName: 'Kevin Tanaka', email: 'kevin@example.com', package: '6x', issueDate: 'Mar 15, 2026', status: 'issued' }
])

const filteredCertificates = computed(() => {
  return issuedCertificates.value.filter(cert => {
    return cert.studentName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
           cert.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  })
})

function issueCertificate(studentId: number) {
  const student = eligibleStudents.value.find(s => s.id === studentId)
  if (student) {
    const newCertId = `EVDA-2026-00${1231 + issuedCertificates.value.length}`
    issuedCertificates.value.unshift({
      id: newCertId,
      studentName: student.name,
      email: student.email,
      package: student.package,
      issueDate: new Date().toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }),
      status: 'issued'
    })
    eligibleStudents.value = eligibleStudents.value.filter(s => s.id !== studentId)
    toast.add({ title: 'Certificate Issued', description: `Certificate ${newCertId} has been issued to ${student.name}.`, icon: 'i-lucide-award', color: 'success' })
  }
}

function revokeCertificate(certId: string) {
  const cert = issuedCertificates.value.find(c => c.id === certId)
  if (cert) {
    cert.status = 'revoked'
    toast.add({ title: 'Certificate Revoked', description: `Certificate ${certId} has been revoked.`, icon: 'i-lucide-x-circle', color: 'error' })
  }
}

const columns: TableColumn<Certificate>[] = [
  {
    accessorKey: 'id',
    header: 'Certificate ID',
    cell: ({ row }) => h('span', { class: 'font-mono text-md' }, row.getValue('id'))
  },
  {
    accessorKey: 'studentName',
    header: 'Student',
    cell: ({ row }) => {
      const Avatar = resolveComponent('UAvatar')
      const name = row.getValue('studentName') as string
      const email = row.original.email
      const initials = name.split(' ').map((n: string) => n[0]).join('')
      return h('div', { class: 'flex items-center gap-3' }, [
        h(Avatar, { text: initials, size: 'md' }),
        h('div', {}, [
          h('p', { class: 'font-medium' }, name),
          h('p', { class: 'text-md text-muted' }, email)
        ])
      ])
    }
  },
  {
    accessorKey: 'package',
    header: 'Package',
    cell: ({ row }) => {
      const Badge = resolveComponent('UBadge')
      return h(Badge, { label: row.getValue('package') as string, color: 'neutral', variant: 'subtle' })
    }
  },
  { accessorKey: 'issueDate', header: 'Issue Date' },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      const Badge = resolveComponent('UBadge')
      const status = row.getValue('status') as string
      return h(Badge, { label: status === 'issued' ? 'Active' : 'Revoked', color: status === 'issued' ? 'success' : 'error', variant: 'subtle' })
    }
  },
  {
    id: 'actions',
    header: '',
    cell: ({ row }) => {
      const DropdownMenu = resolveComponent('UDropdownMenu')
      const Button = resolveComponent('UButton')
      const certId = row.original.id
      const items = [
        [{ label: 'View Certificate', icon: 'i-lucide-eye' }, { label: 'Download PDF', icon: 'i-lucide-download' }, { label: 'Resend Email', icon: 'i-lucide-mail' }],
        [{ label: 'Revoke Certificate', icon: 'i-lucide-x-circle', color: 'error', onSelect: () => revokeCertificate(certId) }]
      ]
      return h(DropdownMenu, { items }, () => h(Button, { icon: 'i-lucide-ellipsis-vertical', color: 'neutral', variant: 'ghost' }))
    }
  }
]
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Certificate Management">
        <template #right>
          <UButton icon="i-lucide-file-badge" color="warning" label="Issue Certificate" @click="showIssueModal = true" />
          <!-- Issue Certificate Modal -->
          <UModal v-model:open="showIssueModal" title="Issue New Certificate">
            <template #body>
              <div class="space-y-4">
                <UFormField label="Select Student" required>
                  <USelectMenu :items="eligibleStudents.map(s => ({ label: s.name, value: s.id.toString() }))" placeholder="Search and select student..." searchable color="warning" class="w-full"/>
                </UFormField>
                <UFormField label="Certificate Type" required>
                  <USelect :items="[{ label: 'Basic Completion Certificate', value: 'basic' }, { label: 'Premium Certificate', value: 'premium' }]" color="warning" class="w-full" />
                </UFormField>
                <UAlert icon="i-lucide-info" color="warning" variant="subtle">
                  <template #description>The certificate will be generated automatically and sent to the student&apos;s email.</template>
                </UAlert>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton label="Cancel" variant="ghost" color="neutral" @click="showIssueModal = false" />
                <UButton label="Issue Certificate" color="warning" icon="i-lucide-award" @click="showIssueModal = false" />
              </div>
            </template>
          </UModal>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
      
      <UDashboardToolbar>
        <template #left>
          <UInput v-model="searchQuery" placeholder="Search certificates..." icon="i-lucide-search" class="w-64" color="warning"/>
        </template>
        <template #right>
          <UButton icon="i-lucide-download" label="Export All" color="neutral" variant="outline" />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Stats -->
        <div class="grid md:grid-cols-3 gap-4">
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-info/10"><UIcon name="i-lucide-award" class="size-6 text-info" /></div>
              <div><p class="text-2xl font-bold">{{ issuedCertificates.length }}</p><p class="text-md text-muted">Total Issued</p></div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-amber-500/10"><UIcon name="i-lucide-clock" class="size-6 text-amber-500" /></div>
              <div><p class="text-2xl font-bold">{{ eligibleStudents.length }}</p><p class="text-md text-muted">Pending Issuance</p></div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-green-500/10"><UIcon name="i-lucide-check-circle" class="size-6 text-green-500" /></div>
              <div><p class="text-2xl font-bold">{{ issuedCertificates.filter(c => c.status === 'issued').length }}</p><p class="text-md text-muted">Active Certificates</p></div>
            </div>
          </UCard>
        </div>

        <!-- Pending Issuance -->
        <UCard v-if="eligibleStudents.length > 0">
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-bell" class="size-5 text-amber-500" />
              <h2 class="font-semibold">Pending Certificate Issuance</h2>
              <UBadge :label="eligibleStudents.length.toString()" color="warning" />
            </div>
          </template>
          <div class="space-y-3">
            <div v-for="student in eligibleStudents" :key="student.id" class="flex items-center justify-between p-4 rounded-lg border border-default bg-amber-500/5">
              <div class="flex items-center gap-4">
                <UAvatar :text="student.name.split(' ').map((n: string) => n[0]).join('')" />
                <div><p class="font-medium">{{ student.name }}</p><p class="text-md text-muted">{{ student.email }}</p></div>
                <UBadge :label="student.package + ' Package'" color="warning" variant="subtle" />
                <span class="text-md text-muted">Completed: {{ student.completedDate }}</span>
              </div>
              <UButton label="Issue Certificate" color="warning" icon="i-lucide-award" @click="issueCertificate(student.id)" />
              
            </div>
            <!-- Issue Certificate Modal Component -->
            
          </div>
        </UCard>

        <!-- Issued Certificates -->
        <UCard>
          <template #header><h2 class="font-semibold">Issued Certificates</h2></template>
          <UTable :data="filteredCertificates" color="warning" :columns="columns" />
          <template #footer>
            <div class="flex items-center justify-between">
              <p class="text-md text-muted">Showing {{ filteredCertificates.length }} certificates</p>
              <UPagination :total="issuedCertificates.length" active-color="warning" :items-per-page="10" />
            </div>
          </template>
        </UCard>
      </div>
    </template>

    

  </UDashboardPanel>
</template>
