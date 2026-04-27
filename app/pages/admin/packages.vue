<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { ref } from 'vue'

definePageMeta({ layout: 'admin' })

const toast = useToast()
const showEditModal = ref(false)
const selectedPackage = ref<any>(null)

// Mock packages data
const packages = ref([
  {
    id: 1,
    name: '6x',
    price: 1750000,
    sessions: 6,
    duration: 60,
    description: 'Our most popular package for comprehensive learning',
    features: [
      'Free Trial',
      '6x Sessions',
    ],
    isActive: true,
    isPopular: true,
    totalSold: 89
  },
  {
    id: 2,
    name: '8x',
    price: 1950000,
    sessions: 8,
    duration: 60,
    description: 'Complete mastery with unlimited support',
    features: [
      'Free Trial',
      '8x Sessions',
    ],
    isActive: true,
    totalSold: 22
  },
  {
    id: 3,
    name: '10x',
    price: 2250000,
    sessions: 10,
    duration: 60,
    description: 'Complete mastery with unlimited support',
    features: [
      'Free Trial',
      '10x Sessions',
    ],
    isActive: true,
    totalSold: 22
  },
  {
    id: 4,
    name: '12x',
    price: 2650000,
    sessions: 12,
    duration: 60,
    description: 'Complete mastery with unlimited support',
    features: [
      'Free Trial',
      '12x Sessions',
    ],
    isActive: true,
    totalSold: 22
  }
])

// Mock add-ons data
const addOns = ref([
  { id: 1, name: 'Extra Session', price: 350000, description: 'Additional training session', sold: 34 },
  // { id: 2, name: 'Simulator Session', price: 150000, description: 'Practice in driving simulator', sold: 22 },
  // { id: 3, name: 'SIM Exam Prep', price: 500000, description: 'License test preparation', sold: 18 }
])

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

function editPackage(pkg: any) {
  selectedPackage.value = { ...pkg }
  showEditModal.value = true
}

function togglePackageStatus(pkgId: number) {
  const pkg = packages.value.find(p => p.id === pkgId)
  if (pkg) {
    pkg.isActive = !pkg.isActive
    toast.add({
      title: pkg.isActive ? 'Package Activated' : 'Package Deactivated',
      icon: pkg.isActive ? 'i-lucide-check-circle' : 'i-lucide-x-circle',
      color: pkg.isActive ? 'success' : 'warning'
    })
  }
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Package Management">
        <template #right>
          <UButton icon="i-lucide-plus" color="warning" label="Add Package" />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Stats -->
        <div class="grid md:grid-cols-3 gap-4">
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-info/10">
                <UIcon name="i-lucide-package" class="size-6 text-info" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ packages.length }}</p>
                <p class="text-md text-muted">Total Packages</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-green-500/10">
                <UIcon name="i-lucide-shopping-cart" class="size-6 text-green-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ packages.reduce((sum, p) => sum + p.totalSold, 0) }}</p>
                <p class="text-md text-muted">Total Sold</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-amber-500/10">
                <UIcon name="i-lucide-banknote" class="size-6 text-amber-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ formatPrice(packages.reduce((sum, p) => sum + (p.price * p.totalSold), 0)) }}</p>
                <p class="text-md text-muted">Total Revenue</p>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Packages Grid -->
        <div class="grid md:grid-cols-2 gap-6">
          <UCard v-for="pkg in packages" :key="pkg.id" :class="{ 'ring-2 ring-warning': pkg.isPopular }">
            <template #header>
              <div class="flex items-start justify-between">
                <div>
                  <div class="flex items-center gap-2">
                    <h3 class="text-xl font-bold">{{ pkg.name }}</h3>
                    <UBadge v-if="pkg.isPopular" label="Popular" color="warning" size="md" />
                  </div>
                  <p class="text-2xl font-bold text-warning mt-2">{{ formatPrice(pkg.price) }}</p>
                </div>
                <USwitch 
                  :model-value="pkg.isActive" 
                  @update:model-value="togglePackageStatus(pkg.id)"
                />
              </div>
            </template>

            <div class="space-y-4">
              <p class="text-md text-muted">{{ pkg.description }}</p>
              
              <div class="grid grid-cols-2 gap-3 text-md">
                <div class="p-2 rounded-lg bg-muted/50 text-center">
                  <p class="font-bold">{{ pkg.sessions }}</p>
                  <p class="text-md text-muted">Sessions</p>
                </div>
                <div class="p-2 rounded-lg bg-muted/50 text-center">
                  <p class="font-bold">{{ pkg.duration }} min</p>
                  <p class="text-md text-muted">Per Session</p>
                </div>
              </div>

              <ul class="space-y-2">
                <li v-for="feature in pkg.features" :key="feature" class="flex items-center gap-2 text-md">
                  <UIcon name="i-lucide-check" class="size-4 text-warning" />
                  {{ feature }}
                </li>
              </ul>

              <USeparator />

              <div class="flex items-center justify-between text-md">
                <span class="text-muted">Total Sold</span>
                <span class="font-bold">{{ pkg.totalSold }} students</span>
              </div>
            </div>

            <template #footer>
              <div class="flex gap-2">
                <UButton label="Edit" icon="i-lucide-pencil" variant="outline" color="neutral" class="flex-1" @click="editPackage(pkg)" />
                <UDropdownMenu
                  :items="[
                    [
                      { label: 'View Sales', icon: 'i-lucide-chart-bar' },
                      { label: 'Duplicate', icon: 'i-lucide-copy' }
                    ],
                    [
                      { label: 'Delete', icon: 'i-lucide-trash', color: 'error' }
                    ]
                  ]"
                >
                  <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" />
                </UDropdownMenu>
                
              </div>

            </template>
          </UCard>
          <!-- Edit Package Modal -->
          <UModal v-model:open="showEditModal" title="Edit Package">
            <template #body>
              <div v-if="selectedPackage" class="space-y-4">
                <UFormField label="Package Name" required>
                  <UInput v-model="selectedPackage.name" class="w-full" color="warning"/>
                </UFormField>
                <UFormField label="Price (IDR)" required>
                  <UInput v-model="selectedPackage.price" type="number" :step="100000" class="w-full" color="warning"/>
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                  <UFormField label="Sessions" required>
                    <UInput v-model="selectedPackage.sessions" type="number" class="w-full" color="warning" />
                  </UFormField>
                  <UFormField label="Duration (min)" required>
                    <UInput v-model="selectedPackage.duration" type="number" :step="15" class="w-full" color="warning" />
                  </UFormField>
                </div>
                <UFormField label="Description">
                  <UTextarea v-model="selectedPackage.description" class="w-full" color="warning" />
                </UFormField>
                <USwitch v-model="selectedPackage.isPopular" label="Mark as Popular" class="w-full" color="warning" />
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-3">
                <UButton label="Cancel" variant="ghost" color="neutral" @click="showEditModal = false" />
                <UButton label="Save Changes" icon="i-lucide-save" color="warning" @click="showEditModal = false" />
              </div>
            </template>
          </UModal>
        </div>

        <!-- Add-ons Section -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">Package Add-ons</h2>
              <UButton label="Add Add-on" icon="i-lucide-plus" size="md" color="warning" variant="outline" />
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="addon in addOns" 
              :key="addon.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex-1">
                <div class="flex items-center gap-3">
                  <p class="font-medium">{{ addon.name }}</p>
                  <UBadge :label="formatPrice(addon.price)" variant="subtle" color="warning" />
                </div>
                <p class="text-md text-muted mt-1">{{ addon.description }}</p>
              </div>
              <div class="flex items-center gap-4">
                <div class="text-right">
                  <p class="font-bold">{{ addon.sold }}</p>
                  <p class="text-md text-muted">Sold</p>
                </div>
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="md" />
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
