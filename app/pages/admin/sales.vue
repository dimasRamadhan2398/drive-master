<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import { useSalesStore } from '~/stores/sales'
import { usePackagesStore } from '~/stores/packages'

definePageMeta({ layout: 'admin' })

const route = useRoute()
const salesStore = useSalesStore()
const packagesStore = usePackagesStore()

const packageId = computed(() => route.query.packageId ? Number(route.query.packageId) : null)

// Get packages from store
const packages = computed(() => packagesStore.packages)

// Date filter
const startDate = ref(salesStore.startDate)
const endDate = ref(salesStore.endDate)

// Watch for changes and update store
watch([startDate, endDate], ([newStart, newEnd]) => {
  salesStore.setDateRange(newStart, newEnd)
})

// Filtered transactions from store
const filteredTransactions = computed(() => salesStore.filteredTransactions)

// Data for specific package view
const selectedPackage = computed(() => packages.value.find(p => p.id === packageId.value))
const packageTransactions = computed(() => {
    if (!packageId.value) return []
    return salesStore.transactionsByPackage(packageId.value)
})
const packageTotalRevenue = computed(() => packageTransactions.value.reduce((sum, t) => sum + t.amount, 0))
const packageTotalSales = computed(() => packageTransactions.value.length)

// Data for general dashboard view
const overallTotalRevenue = computed(() => salesStore.filteredTotalRevenue)
const overallTotalSales = computed(() => salesStore.filteredTotalSales)

// Sales breakdown by package
const salesByPackage = computed(() => {
  return packages.value.map(pkg => {
    const pkgTransactions = filteredTransactions.value.filter(t => t.packageId === pkg.id)
    const totalSales = pkgTransactions.length
    const revenue = pkgTransactions.reduce((sum, t) => sum + t.amount, 0)
    return {
      ...pkg,
      totalSales,
      revenue
    }
  }).sort((a, b) => b.revenue - a.revenue)
})

// Transactions to display in table
const displayTransactions = computed(() => {
  return packageId.value ? packageTransactions.value : filteredTransactions.value
})

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
}

// Table columns (Nuxt UI v3 style)
const columns = computed<TableColumn<any>[]>(() => {
  const base: TableColumn<any>[] = [
    { accessorKey: 'id', header: 'ID' },
    { accessorKey: 'studentName', header: 'Student Name' },
  ]

  if (!packageId.value) {
    base.push({ accessorKey: 'packageName', header: 'Package' })
  }

  return [
    ...base,
    { accessorKey: 'purchaseDateFormatted', header: 'Date' },
    { accessorKey: 'amountFormatted', header: 'Amount' },
    { accessorKey: 'status', header: 'Status' }
  ]
})

// Data for table, formatted
const tableData = computed(() => displayTransactions.value.map(t => {
  const pkg = packages.value.find(p => p.id === t.packageId)
  return {
    ...t,
    packageName: pkg ? pkg.name : 'N/A',
    purchaseDateFormatted: formatDate(t.purchaseDate),
    amountFormatted: formatPrice(t.amount)
  }
}))

function clearFilter() {
  startDate.value = ''
  endDate.value = ''
  salesStore.clearDateFilter()
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar>
        <template #left>
          <!-- Header Tampilan Paket Spesifik -->
          <div v-if="packageId" class="flex items-center gap-4">
            <UButton to="/admin/packages" icon="i-lucide-arrow-left" color="neutral" variant="ghost" aria-label="Back to packages" />
            <div>
              <p class="text-sm text-muted">Sales Report</p>
              <h1 v-if="selectedPackage" class="text-xl font-bold">{{ selectedPackage.name }} Package</h1>
              <h1 v-else class="text-xl font-bold">Loading...</h1>
            </div>
          </div>
          <!-- Header Tampilan Dasbor Umum -->
          <div v-else>
            <h1 class="text-xl font-bold">Sales Dashboard</h1>
            <p class="text-sm text-muted">Comprehensive overview of all package sales</p>
          </div>
        </template>
        <template #right>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <!-- Filter Tanggal (Visible in both views) -->
      <div class="px-6 pt-6 flex items-end gap-4 pb-2">
        <UFormGroup label="Start Date">
          <UInput type="date" v-model="startDate" />
        </UFormGroup>
        <UFormGroup label="End Date">
          <UInput type="date" v-model="endDate" />
        </UFormGroup>
        <UButton v-if="startDate || endDate" label="Clear Filter" color="neutral" variant="soft" @click="clearFilter" />
      </div>

      <!-- Tampilan Body Paket Spesifik -->
      <div v-if="packageId">
        <div v-if="selectedPackage" class="p-6 space-y-6">
          <!-- Stats -->
          <div class="grid md:grid-cols-2 gap-6">
            <UCard>
              <div class="flex items-center gap-4">
                <div class="p-3 rounded-xl bg-green-500/10">
                  <UIcon name="i-lucide-shopping-cart" class="size-6 text-green-500" />
                </div>
                <div>
                  <p class="text-2xl font-bold">{{ packageTotalSales }}</p>
                  <p class="text-md text-muted">Total Sales</p>
                </div>
              </div>
            </UCard>
            <UCard>
              <div class="flex items-center gap-4">
                <div class="p-3 rounded-xl bg-amber-500/10">
                  <UIcon name="i-lucide-banknote" class="size-6 text-amber-500" />
                </div>
                <div>
                  <p class="text-2xl font-bold">{{ formatPrice(packageTotalRevenue) }}</p>
                  <p class="text-md text-muted">Total Revenue</p>
                </div>
              </div>
            </UCard>
          </div>
          <!-- Transactions Table -->
          <UCard>
            <template #header>
              <h2 class="font-semibold">Transaction History for {{ selectedPackage.name }}</h2>
            </template>
            <UTable :data="tableData" :columns="columns" />
          </UCard>
        </div>
        <div v-else class="p-6 text-center">
          <UAlert icon="i-lucide-alert-circle" color="error" title="Package not found" description="The package you are looking for does not exist or has been removed." />
          <UButton to="/admin/sales" label="Back to Dashboard" color="neutral" variant="ghost" class="mt-4" />
        </div>
      </div>

      <!-- Tampilan Body Dasbor Umum -->
      <div v-else class="p-6 space-y-8">
        <!-- Overview Stats -->
        <div class="grid md:grid-cols-3 gap-6">
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-info/10">
                <UIcon name="i-lucide-package" class="size-6 text-info" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ packages.length }}</p>
                <p class="text-md text-muted">Active Packages</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-green-500/10">
                <UIcon name="i-lucide-shopping-cart" class="size-6 text-green-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ overallTotalSales }}</p>
                <p class="text-md text-muted">Total Units Sold</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-amber-500/10">
                <UIcon name="i-lucide-banknote" class="size-6 text-amber-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ formatPrice(overallTotalRevenue) }}</p>
                <p class="text-md text-muted">Lifetime Revenue</p>
              </div>
            </div>
          </UCard>
        </div>

        

        <!-- Sales by Package Dashboard -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold">Performance by Package</h2>
            <p class="text-sm text-muted">Click a package to view details</p>
          </div>
          <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
            <NuxtLink 
              v-for="pkg in salesByPackage" 
              :key="pkg.id" 
              :to="{ path: '/admin/sales', query: { packageId: pkg.id } }"
              class="block group"
            >
              <UCard class="hover:ring-2 hover:ring-warning transition-all cursor-pointer h-full">
                <div class="space-y-4">
                  <div class="flex items-center justify-between">
                    <UBadge :label="pkg.name" color="neutral" variant="subtle" size="md" />
                    <UIcon name="i-lucide-trending-up" class="size-4 text-green-500 opacity-0 group-hover:opacity-100 transition-opacity" />
                  </div>
                  <div>
                    <p class="text-sm text-muted">Revenue</p>
                    <p class="text-xl font-bold">{{ formatPrice(pkg.revenue) }}</p>
                  </div>
                  <div class="flex items-center justify-between pt-2 border-t border-default">
                    <span class="text-xs text-muted">{{ pkg.totalSales }} Sales</span>
                    <span class="text-xs font-medium text-warning group-hover:underline">View Report</span>
                  </div>
                </div>
              </UCard>
            </NuxtLink>
          </div>
        </div>

        <!-- Transactions Table -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">Recent Transactions</h2>
            </div>
          </template>
          <UTable :data="tableData" :columns="columns" />
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
