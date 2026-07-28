<script setup lang="ts">
import { computed, h, ref, resolveComponent, onMounted } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { useStudentsStore } from '~/stores/students'
import { certificateService } from '~/services/certificateService'
import { useApiClients } from '~/composables/useApiClients'

const { t } = useI18n()
definePageMeta({ layout: 'admin' })

const toast = useToast()
const studentsStore = useStudentsStore()
const searchQuery = ref('')
const showIssueModal = ref(false)
const selectedStudent = ref<any>(null)
const isLoading = ref(false)

type EligibleStudent = {
  id: string
  name: string
  email: string
  package: string
  packageId: string
  entitlementId: string
  completedDate: string
}

type Certificate = {
  id: string
  certNumber: string
  studentName: string
  email: string
  package: string
  issueDate: string
  status: 'issued' | 'revoked'
}

const eligibleStudents = ref<EligibleStudent[]>([])
const issuedCertificates = ref<Certificate[]>([])

async function loadData() {
  isLoading.value = true
  try {
    // 1. Fetch all issued certificates
    const certsResponse = await certificateService.getAllCertificates()
    issuedCertificates.value = certsResponse.map(c => ({
      id: c.id,
      certNumber: c.certNumber,
      studentName: c.memberName,
      email: c.memberEmail,
      package: c.packageName,
      issueDate: new Date(c.issuedDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }),
      status: c.status === 'revoked' ? 'revoked' : 'issued'
    }))

    // 2. Fetch all students to see who completed a package
    await studentsStore.fetchStudentsNoPagination()
    const completedStudentsList = studentsStore.allStudents.filter(s => s.status === 'completed' || s.progress >= 100)

    // 3. Filter out students who already have an active certificate for their package
    const tempEligible: EligibleStudent[] = []
    for (const student of completedStudentsList) {
      for (const entitlement of student.entitlements) {
        if (entitlement.remaining === 0 || entitlement.status === 'completed') {
          const alreadyIssued = certsResponse.some(c => 
            c.memberId === student.id && 
            c.packageName === entitlement.packageName && 
            c.status != 'revoked'
          )
          if (!alreadyIssued) {
            tempEligible.push({
              id: student.id,
              name: student.name,
              email: student.email,
              package: entitlement.packageName,
              packageId: entitlement.packageId,
              entitlementId: entitlement.id,
              completedDate: new Date(entitlement.updatedAt || entitlement.createdAt).toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
                year: 'numeric'
              })
            })
          }
        }
      }
    }
    eligibleStudents.value = tempEligible
  } catch (error) {
    console.error('Failed to load certificates data:', error)
    toast.add({ title: 'Error', description: 'Failed to load certificate records.', color: 'error' })
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadData()
})

const filteredCertificates = computed(() => {
  return issuedCertificates.value.filter(cert => {
    return cert.studentName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
           cert.certNumber.toLowerCase().includes(searchQuery.value.toLowerCase())
  })
})

async function issueCertificate(studentId: string | null) {
  const targetStudent = studentId 
    ? eligibleStudents.value.find(s => s.id === studentId) 
    : selectedStudent.value?.student

  if (targetStudent) {
    try {
      await certificateService.issueCertificate({
        memberId: targetStudent.id,
        packageId: targetStudent.packageId,
        packageName: targetStudent.package,
        entitlementId: targetStudent.entitlementId
      })
      toast.add({ title: 'Certificate Issued', description: `Certificate has been issued to ${targetStudent.name}.`, icon: 'i-lucide-award', color: 'success' })
      showIssueModal.value = false
      selectedStudent.value = null
      await loadData()
    } catch (error) {
      console.error('Failed to issue certificate:', error)
      toast.add({ title: 'Error', description: 'Failed to issue certificate.', color: 'error' })
    }
  } else {
    toast.add({ title: 'Error', description: 'Please select a student to issue certificate.', color: 'error' })
  }
}

async function revokeCertificate(certId: string) {
  try {
    await certificateService.revokeCertificate(certId)
    toast.add({ title: 'Certificate Revoked', description: `Certificate has been successfully revoked.`, icon: 'i-lucide-x-circle', color: 'success' })
    await loadData()
  } catch (error) {
    console.error('Failed to revoke certificate:', error)
    toast.add({ title: 'Error', description: 'Failed to revoke certificate.', color: 'error' })
  }
}

async function viewCertificate(certId: string) {
  try {
    const url = await certificateService.getCertificatePDFBlobUrl(certId)
    window.open(url, '_blank')
  } catch (error) {
    console.error('Failed to view certificate:', error)
    toast.add({ title: 'Error', description: 'Failed to load certificate preview.', color: 'error' })
  }
}

async function downloadCertificatePDF(certId: string, certNumber: string) {
  try {
    await certificateService.downloadCertificatePDF(certId, certNumber)
    toast.add({ title: 'Success', description: 'Certificate PDF downloaded successfully.', color: 'success' })
  } catch (error) {
    console.error('Failed to download certificate PDF:', error)
    toast.add({ title: 'Error', description: 'Failed to download certificate PDF.', color: 'error' })
  }
}

const columns: TableColumn<Certificate>[] = [
  {
    accessorKey: 'certNumber',
    header: t('admin.certId'),
    cell: ({ row }) => h('span', { class: 'font-mono text-md' }, row.getValue('certNumber'))
  },
  {
    accessorKey: 'studentName',
    header: t('admin.students'),
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
    header: t('billing.package'),
    cell: ({ row }) => {
      const Badge = resolveComponent('UBadge')
      return h(Badge, { label: row.getValue('package') as string, color: 'neutral', variant: 'subtle' })
    }
  },
  { accessorKey: 'issueDate', header: t('admin.issueDate') },
  {
    accessorKey: 'status',
    header: t('billing.status'),
    cell: ({ row }) => {
      const Badge = resolveComponent('UBadge')
      const status = row.getValue('status') as string
      return h(Badge, { label: status === 'issued' ? t('admin.active') : t('admin.revoked'), color: status === 'issued' ? 'success' : 'error', variant: 'subtle' })
    }
  },
  {
    id: 'actions',
    header: "",
    cell: ({ row }) => {
      const DropdownMenu = resolveComponent('UDropdownMenu')
      const Button = resolveComponent('UButton')
      const certId = row.original.id
      const certNumber = row.original.certNumber
      const items = [
        [
          { label: 'View Certificate', icon: 'i-lucide-eye', onSelect: () => viewCertificate(certId) },
          { label: 'Download PDF', icon: 'i-lucide-download', onSelect: () => downloadCertificatePDF(certId, certNumber) }
        ],
        [
          { label: 'Revoke Certificate', icon: 'i-lucide-x-circle', color: 'error', onSelect: () => revokeCertificate(certId) }
        ]
      ]
      return h(DropdownMenu, { items }, () => h(Button, { icon: 'i-lucide-ellipsis-vertical', color: 'neutral', variant: 'ghost' }))
    }
  }
]
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.certificates')">
        <template #right>
          <UButton icon="i-lucide-file-badge" color="warning" :label="t('admin.issueCert')" @click="showIssueModal = true" />
          <!-- Issue Certificate Modal -->
          <UModal v-model:open="showIssueModal" :title="t('admin.issueNewCert')">
            <template #body>
              <div class="space-y-4">
                <UFormField :label="t('admin.selectStudent')" required>
                  <USelectMenu v-model="selectedStudent" :items="eligibleStudents.map(s => ({ label: s.name, value: s.id, student: s }))" placeholder="Search and select student..." searchable color="warning" class="w-full"/>
                </UFormField>
                <UFormField :label="t('admin.certType')" required>
                  <USelect :items="[{ label: 'Basic Completion Certificate', value: 'basic' }, { label: 'Premium Certificate', value: 'premium' }]" color="warning" class="w-full" />
                </UFormField>
                <UAlert icon="i-lucide-info" color="warning" variant="subtle">
                  <template #description>{{ t('admin.certAutoNote') }}</template>
                </UAlert>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton :label="t('common.cancel')" variant="ghost" color="neutral" @click="showIssueModal = false; selectedStudent = null" />
                <UButton :label="t('admin.issueCert')" color="warning" icon="i-lucide-award" @click="issueCertificate(null)" />
              </div>
            </template>
          </UModal>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
      
      <UDashboardToolbar>
        <template #left>
          <UInput v-model="searchQuery" :placeholder="t('common.search') + '...'" icon="i-lucide-search" class="w-64" color="warning"/>
        </template>
        <template #right>
          <UButton icon="i-lucide-download" :label="t('admin.exportAll')" color="neutral" variant="outline" />
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
              <div><p class="text-2xl font-bold">{{ issuedCertificates.length }}</p><p class="text-md text-muted">{{ t('admin.totalIssued') }}</p></div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-amber-500/10"><UIcon name="i-lucide-clock" class="size-6 text-amber-500" /></div>
              <div><p class="text-2xl font-bold">{{ eligibleStudents.length }}</p><p class="text-md text-muted">{{ t('admin.pendingIssuance') }}</p></div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-green-500/10"><UIcon name="i-lucide-check-circle" class="size-6 text-green-500" /></div>
              <div><p class="text-2xl font-bold">{{ issuedCertificates.filter(c => c.status === 'issued').length }}</p><p class="text-md text-muted">{{ t('admin.activeCerts') }}</p></div>
            </div>
          </UCard>
        </div>

        <!-- Pending Issuance -->
        <UCard v-if="eligibleStudents.length > 0">
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-bell" class="size-5 text-amber-500" />
              <h2 class="font-semibold">{{ t('admin.pendingCertIssuance') }}</h2>
              <UBadge :label="eligibleStudents.length.toString()" color="warning" />
            </div>
          </template>
          <div class="space-y-3">
            <div v-for="student in eligibleStudents" :key="student.id" class="flex items-center justify-between p-4 rounded-lg border border-default bg-amber-500/5">
              <div class="flex items-center gap-4">
                <UAvatar :text="student.name.split(' ').map((n: string) => n[0]).join('')" />
                <div><p class="font-medium">{{ student.name }}</p><p class="text-md text-muted">{{ student.email }}</p></div>
                <UBadge :label="student.package + ' ' + t('billing.package')" color="warning" variant="subtle" />
                <span class="text-md text-muted">{{ t('common.completed') }}: {{ student.completedDate }}</span>
              </div>
              <UButton :label="t('admin.issueCert')" color="warning" icon="i-lucide-award" @click="issueCertificate(student.id)" />
              
            </div>
            <!-- Issue Certificate Modal Component -->
            
          </div>
        </UCard>

        <!-- Issued Certificates -->
        <UCard>
          <template #header><h2 class="font-semibold">{{ t('admin.certificates') }}</h2></template>
          <UTable :data="filteredCertificates" color="warning" :columns="columns" />
          <template #footer>
            <div class="flex items-center justify-between">
              <p class="text-md text-muted">{{ t('admin.showingCerts', { count: filteredCertificates.length }) }}</p>
              <UPagination :total="issuedCertificates.length" active-color="warning" :items-per-page="10" />
            </div>
          </template>
        </UCard>
      </div>
    </template>

    

  </UDashboardPanel>
</template>
