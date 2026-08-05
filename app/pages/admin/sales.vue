<script setup lang="ts">
import { ref, computed, onMounted, watch, h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import { useSalesStore } from '~/stores/sales'
import { usePackagesStore } from '~/stores/packages'
import { useStudentsStore } from '~/stores/students'

const { t } = useI18n()
definePageMeta({ layout: 'admin' })

const route = useRoute()
const salesStore = useSalesStore()
const packagesStore = usePackagesStore()
const studentsStore = useStudentsStore()

onMounted(async () => {
  await Promise.all([
    packagesStore.fetchPackages(),
    packagesStore.fetchAddons(),
    studentsStore.fetchStudentsNoPagination()
  ])
  await salesStore.fetchTransactions()
})

const packageId = computed(() => route.query.packageId ? String(route.query.packageId) : null)
const addonId = computed(() => route.query.addonId ? String(route.query.addonId) : null)

// Get packages & addons from store
const packages = computed(() => packagesStore.packages)
const addons = computed(() => packagesStore.addons)

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

// Data for specific addon view
const selectedAddon = computed(() => addons.value.find(a => a.id === addonId.value))
const addonTransactions = computed(() => {
  if (!addonId.value) return []
  return salesStore.transactionsByAddon(addonId.value)
})
const addonTotalRevenue = computed(() => {
  if (!addonTransactions.value.length && selectedAddon.value) {
    return (selectedAddon.value.sold || 0) * selectedAddon.value.price
  }
  return addonTransactions.value.reduce((sum, t) => sum + t.amount, 0)
})
const addonTotalSales = computed(() => {
  if (!addonTransactions.value.length && selectedAddon.value) {
    return selectedAddon.value.sold || 0
  }
  return addonTransactions.value.length
})

// Data for general dashboard view
const overallTotalRevenue = computed(() => salesStore.filteredTotalRevenue)
const overallTotalSales = computed(() => salesStore.filteredTotalSales)

// Sales breakdown by package
const salesByPackage = computed(() => {
  return packages.value.map(pkg => {
    const pkgTransactions = filteredTransactions.value.filter(t => t.packageId === pkg.id)
    const salesFromTx = pkgTransactions.length
    const revFromTx = pkgTransactions.reduce((sum, t) => sum + t.amount, 0)
    const totalSales = salesFromTx > 0 ? salesFromTx : (pkg.totalSold || 0)
    const revenue = revFromTx > 0 ? revFromTx : (pkg.totalSold || 0) * (pkg.discountPrice || pkg.price)
    return {
      ...pkg,
      totalSales,
      revenue
    }
  }).sort((a, b) => b.revenue - a.revenue)
})

// Sales breakdown by addon
const salesByAddon = computed(() => {
  return addons.value.map(addon => {
    const addonTx = filteredTransactions.value.filter(t => t.addonId === addon.id)
    const salesFromTx = addonTx.length
    const revFromTx = addonTx.reduce((sum, t) => sum + t.amount, 0)
    const totalSales = salesFromTx > 0 ? salesFromTx : (addon.sold || 0)
    const revenue = revFromTx > 0 ? revFromTx : (addon.sold || 0) * addon.price
    return {
      ...addon,
      totalSales,
      revenue
    }
  }).sort((a, b) => b.revenue - a.revenue)
})

// Transactions to display in table
const displayTransactions = computed(() => {
  if (packageId.value) return packageTransactions.value
  if (addonId.value) return addonTransactions.value
  return filteredTransactions.value
})

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
}

// Table columns
const columns = computed<TableColumn<any>[]>(() => {
  const base: TableColumn<any>[] = [
    { accessorKey: 'id', header: 'ID' },
    { accessorKey: 'studentName', header: t('admin.students').replace('Murid', 'Nama Murid') },
  ]

  if (!packageId.value && !addonId.value) {
    base.push({ accessorKey: 'itemName', header: 'Item' })
    base.push({
      accessorKey: 'itemType',
      header: 'Type',
      cell: ({ row }: any) => {
        const isPackage = row.original.itemType === 'package'
        return h(resolveComponent('UBadge'), {
          label: isPackage ? 'PACKAGE' : 'ADDON',
          color: isPackage ? 'warning' : 'info',
          variant: 'subtle'
        })
      }
    })
  }

  return [
    ...base,
    { accessorKey: 'purchaseDateFormatted', header: t('dashboard.date') },
    { accessorKey: 'amountFormatted', header: t('billing.amount') },
    { accessorKey: 'status', header: t('billing.status') }
  ]
})

// Data for table, formatted
const tableData = computed(() => displayTransactions.value.map(t => {
  let itemName = 'N/A'
  let itemType: 'package' | 'addon' = t.addonId ? 'addon' : 'package'

  if (t.packageId) {
    const pkg = packages.value.find(p => p.id === t.packageId)
    if (pkg) itemName = pkg.name
  } else if (t.addonId) {
    const addon = addons.value.find(a => a.id === t.addonId)
    if (addon) itemName = addon.name
  }

  return {
    ...t,
    itemName,
    itemType,
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
            <UButton to="/admin/sales" icon="i-lucide-arrow-left" color="neutral" variant="ghost" aria-label="Back to sales" />
            <div>
              <p class="text-sm text-muted">{{ t('admin.salesReport') }}</p>
              <h1 v-if="selectedPackage" class="text-xl font-bold">{{ selectedPackage.name }} {{ t('billing.package') }}</h1>
              <h1 v-else class="text-xl font-bold">{{ t('common.loading') }}</h1>
            </div>
          </div>

          <!-- Header Tampilan Addon Spesifik -->
          <div v-else-if="addonId" class="flex items-center gap-4">
            <UButton to="/admin/sales" icon="i-lucide-arrow-left" color="neutral" variant="ghost" aria-label="Back to sales" />
            <div>
              <p class="text-sm text-muted">Addon Sales Report</p>
              <h1 v-if="selectedAddon" class="text-xl font-bold">{{ selectedAddon.name }} Addon</h1>
              <h1 v-else class="text-xl font-bold">{{ t('common.loading') }}</h1>
            </div>
          </div>

          <!-- Header Tampilan Dasbor Umum -->
          <div v-else>
            <h1 class="text-xl font-bold">{{ t('admin.salesDashboard') }}</h1>
            <p class="text-sm text-muted">Performance analytics for packages and addons</p>
          </div>
        </template>
        <template #right>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <!-- Filter Tanggal (Visible in all views) -->
      <div class="px-6 pt-6 flex items-end gap-4 pb-2">
        <UFormField :label="t('admin.startDate')">
          <UInput type="date" v-model="startDate" />
        </UFormField>
        <UFormField :label="t('admin.endDate')">
          <UInput type="date" v-model="endDate" />
        </UFormField>
        <UButton v-if="startDate || endDate" :label="t('admin.clearFilter')" color="neutral" variant="soft" @click="clearFilter" />
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
                  <p class="text-md text-muted">{{ t('admin.unitsSold') }}</p>
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
              <h2 class="font-semibold">{{ t('admin.transactionHistoryFor', { name: selectedPackage.name }) }}</h2>
            </template>
            <UTable :data="tableData" :columns="columns" />
          </UCard>
        </div>
        <div v-else class="p-6 text-center">
          <UAlert icon="i-lucide-alert-circle" color="error" title="Package not found" description="The package you are looking for does not exist or has been removed." />
          <UButton to="/admin/sales" label="Back to Dashboard" color="neutral" variant="ghost" class="mt-4" />
        </div>
      </div>

      <!-- Tampilan Body Addon Spesifik -->
      <div v-else-if="addonId">
        <div v-if="selectedAddon" class="p-6 space-y-6">
          <!-- Stats -->
          <div class="grid md:grid-cols-2 gap-6">
            <UCard>
              <div class="flex items-center gap-4">
                <div class="p-3 rounded-xl bg-info/10">
                  <UIcon name="i-lucide-layers" class="size-6 text-info" />
                </div>
                <div>
                  <p class="text-2xl font-bold">{{ addonTotalSales }}</p>
                  <p class="text-md text-muted">Units Sold</p>
                </div>
              </div>
            </UCard>
            <UCard>
              <div class="flex items-center gap-4">
                <div class="p-3 rounded-xl bg-amber-500/10">
                  <UIcon name="i-lucide-banknote" class="size-6 text-amber-500" />
                </div>
                <div>
                  <p class="text-2xl font-bold">{{ formatPrice(addonTotalRevenue) }}</p>
                  <p class="text-md text-muted">Total Addon Revenue</p>
                </div>
              </div>
            </UCard>
          </div>
          <!-- Transactions Table -->
          <UCard>
            <template #header>
              <h2 class="font-semibold">Transaction History for Addon: {{ selectedAddon.name }}</h2>
            </template>
            <UTable :data="tableData" :columns="columns" />
          </UCard>
        </div>
        <div v-else class="p-6 text-center">
          <UAlert icon="i-lucide-alert-circle" color="error" title="Addon not found" description="The addon you are looking for does not exist or has been removed." />
          <UButton to="/admin/sales" label="Back to Dashboard" color="neutral" variant="ghost" class="mt-4" />
        </div>
      </div>

      <!-- Tampilan Body Dasbor Umum -->
      <div v-else class="p-6 space-y-8">
        <!-- Overview Stats -->
        <div class="grid md:grid-cols-4 gap-6">
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-info/10">
                <UIcon name="i-lucide-package" class="size-6 text-info" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ packages.length }}</p>
                <p class="text-sm text-muted">Active Packages</p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-purple-500/10">
                <UIcon name="i-lucide-layers" class="size-6 text-purple-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ addons.length }}</p>
                <p class="text-sm text-muted">Active Addons</p>
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
                <p class="text-sm text-muted">{{ t('admin.unitsSold') }}</p>
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
                <p class="text-sm text-muted">{{ t('admin.lifetimeRevenue') }}</p>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Sales by Package Section -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-package" class="size-5 text-warning" />
              <h2 class="text-lg font-bold">{{ t('admin.performanceByPackage') }}</h2>
            </div>
            <p class="text-sm text-muted">{{ t('admin.clickPackageDetails') }}</p>
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
                    <p class="text-xs text-muted">Revenue</p>
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

        <!-- Sales by Addon Section -->
        <div class="space-y-4 pt-4 border-t border-default">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-layers" class="size-5 text-info" />
              <h2 class="text-lg font-bold">Performance by Addon</h2>
            </div>
            <p class="text-sm text-muted">Click an addon card to filter transactions</p>
          </div>
          <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
            <NuxtLink 
              v-for="addon in salesByAddon" 
              :key="addon.id" 
              :to="{ path: '/admin/sales', query: { addonId: addon.id } }"
              class="block group"
            >
              <UCard class="hover:ring-2 hover:ring-info transition-all cursor-pointer h-full">
                <div class="space-y-4">
                  <div class="flex items-center justify-between">
                    <UBadge :label="addon.name" color="info" variant="subtle" size="md" />
                    <UIcon name="i-lucide-trending-up" class="size-4 text-green-500 opacity-0 group-hover:opacity-100 transition-opacity" />
                  </div>
                  <div>
                    <p class="text-xs text-muted">Revenue ({{ formatPrice(addon.price) }}/unit)</p>
                    <p class="text-xl font-bold">{{ formatPrice(addon.revenue) }}</p>
                  </div>
                  <div class="flex items-center justify-between pt-2 border-t border-default">
                    <span class="text-xs text-muted">{{ addon.totalSales }} Units Sold</span>
                    <span class="text-xs font-medium text-info group-hover:underline">View Report</span>
                  </div>
                </div>
              </UCard>
            </NuxtLink>
          </div>
        </div>

        <!-- Recent Transactions Table -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">{{ t('admin.recentTransactions') }}</h2>
            </div>
          </template>
          <UTable :data="tableData" :columns="columns" />
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
