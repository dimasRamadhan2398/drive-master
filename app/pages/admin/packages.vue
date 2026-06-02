<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { ref, computed } from 'vue'
import { usePackagesStore } from '~/stores/packages'
import type { Package, Addon } from '~/stores/packages'

definePageMeta({ layout: 'admin' })

const toast = useToast()
const packagesStore = usePackagesStore()

const showEditModal = ref(false)
const showAddModal = ref(false)
const showAddonModal = ref(false)
const showEditAddonModal = ref(false)
const selectedPackage = ref<Package | null>(null)

const packages = computed(() => packagesStore.packages)
const addOns = computed(() => packagesStore.addons)

// Stats computed from store
const totalPackages = computed(() => packagesStore.totalPackages)
const totalSold = computed(() => packagesStore.totalSold)
const totalRevenue = computed(() => packagesStore.totalRevenue)

const newPackage = ref({
  name: '',
  price: 0,
  sessions: 0,
  duration: 60,
  description: '',
  features: '', // Akan menjadi string dari textarea, dipisah per baris saat disimpan
  isPopular: false
})

const newAddon = ref({
  name: '',
  price: 0,
  description: ''
})

const editingAddon = ref<Addon | null>(null)

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

function editPackage(pkg: Package) {
  selectedPackage.value = {
    ...pkg,
    features: Array.isArray(pkg.features) ? pkg.features.join('\n') : ''
  }
  showEditModal.value = true
}

function togglePackageStatus(pkgId: number) {
  const isActive = packagesStore.togglePackageStatus(pkgId)
  if (isActive !== null) {
    toast.add({
      title: isActive ? 'Package Activated' : 'Package Deactivated',
      icon: isActive ? 'i-lucide-check-circle' : 'i-lucide-x-circle',
      color: isActive ? 'success' : 'warning'
    })
  }
}

function saveNewPackage() {
  if (!newPackage.value.name || newPackage.value.price <= 0) {
    toast.add({ title: 'Error', description: 'Nama paket dan harga yang valid wajib diisi.', color: 'error' })
    return
  }

  packagesStore.addPackage({
    name: newPackage.value.name,
    price: newPackage.value.price,
    sessions: newPackage.value.sessions,
    duration: newPackage.value.duration,
    description: newPackage.value.description,
    features: newPackage.value.features.split('\n').filter(f => f.trim() !== ''),
    isActive: true,
    isPopular: newPackage.value.isPopular,
  })

  toast.add({ title: 'Paket Ditambahkan', description: `"${newPackage.value.name}" telah dibuat.`, color: 'success' })

  showAddModal.value = false
  newPackage.value = {
    name: '',
    price: 0,
    sessions: 0,
    duration: 60,
    description: '',
    features: '',
    isPopular: false
  }
}

function saveEditedPackage() {
  if (!selectedPackage.value) return

  packagesStore.updatePackage(selectedPackage.value.id, {
    name: selectedPackage.value.name,
    price: selectedPackage.value.price,
    sessions: selectedPackage.value.sessions,
    duration: selectedPackage.value.duration,
    description: selectedPackage.value.description,
    features: selectedPackage.value.features.split('\n').filter((f: string) => f.trim() !== ''),
    isPopular: selectedPackage.value.isPopular,
  })

  toast.add({ title: 'Package Updated', description: `"${selectedPackage.value.name}" has been saved.`, color: 'success' })
  showEditModal.value = false
  selectedPackage.value = null
}

function duplicatePackage(pkg: Package) {
  packagesStore.duplicatePackage(pkg.id)
  toast.add({ title: 'Package Duplicated', description: `Salinan dari "${pkg.name}" telah dibuat.`, color: 'success' });
}

function deletePackage(pkgId: number) {
  if (confirm(`Anda yakin ingin menghapus paket ini? Aksi ini tidak dapat dibatalkan.`)) {
    packagesStore.deletePackage(pkgId)
    toast.add({ title: 'Package Deleted', description: `Paket telah dihapus.`, color: 'error', icon: 'i-lucide-trash' });
  }
}

function viewPackageSales(pkg: Package) {
  navigateTo(`/admin/sales?packageId=${pkg.id}`);
}

function saveNewAddon() {
  if (!newAddon.value.name || newAddon.value.price <= 0) {
    toast.add({ title: 'Error', description: 'Nama dan harga add-on wajib diisi.', color: 'error' })
    return
  }

  packagesStore.addAddon({
    name: newAddon.value.name,
    price: newAddon.value.price,
    description: newAddon.value.description,
  })

  toast.add({ title: 'Add-on Ditambahkan', description: `"${newAddon.value.name}" telah dibuat.`, color: 'success' })

  showAddonModal.value = false
  newAddon.value = { name: '', price: 0, description: '' }
}

function openEditAddonModal(addon: Addon) {
  editingAddon.value = { ...addon }
  showEditAddonModal.value = true
}

function saveEditedAddon() {
  if (!editingAddon.value) return

  const updated = packagesStore.updateAddon(editingAddon.value.id, editingAddon.value)
  if (updated) {
    toast.add({ title: 'Add-on Updated', description: `"${editingAddon.value.name}" has been saved.`, color: 'success' })
  } else {
    toast.add({ title: 'Error', description: 'Add-on not found.', color: 'error' })
  }

  showEditAddonModal.value = false
  editingAddon.value = null
}

</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Package Management">
        <template #right>
          <UButton icon="i-lucide-plus" color="warning" label="Add Package" @click="showAddModal = true" />
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
                <p class="text-2xl font-bold">{{ totalPackages }}</p>
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
                <p class="text-2xl font-bold">{{ totalSold }}</p>
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
                <p class="text-2xl font-bold">{{ formatPrice(totalRevenue) }}</p>
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
                      { label: 'View Sales', icon: 'i-lucide-chart-bar', onSelect: () => viewPackageSales(pkg) },
                      { label: 'Duplicate', icon: 'i-lucide-copy', onSelect: () => duplicatePackage(pkg) }
                    ],
                    [
                      { label: 'Delete', icon: 'i-lucide-trash', color: 'error', onSelect: () => deletePackage(pkg.id) }
                    ]
                  ]"
                >
                  <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" />
                </UDropdownMenu>
                
              </div>

            </template>
          </UCard>
        </div>

        <!-- Modals -->
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
              <UFormField label="Features">
                <template #hint>
                  <span>Masukkan satu fitur per baris.</span>
                </template>
                <UTextarea v-model="selectedPackage.features" placeholder="Free Trial&#10;15x Sessions" class="w-full" color="warning" />
              </UFormField>
              <USwitch v-model="selectedPackage.isPopular" label="Mark as Popular" class="w-full" color="warning" />
            </div>
          </template>
          <template #footer>
            <div class="flex justify-end gap-3">
              <UButton label="Cancel" variant="ghost" color="neutral" @click="showEditModal = false" />
              <UButton label="Save Changes" icon="i-lucide-save" color="warning" @click="saveEditedPackage" />
            </div>
          </template>
        </UModal>

        <!-- Add Package Modal -->
        <UModal v-model:open="showAddModal" title="Add New Package">
          <template #body>
            <div class="space-y-4">
              <UFormField label="Package Name" required>
                <UInput v-model="newPackage.name" placeholder="e.g. 15x Sessions" class="w-full" color="warning"/>
              </UFormField>
              <UFormField label="Price (IDR)" required>
                <UInput v-model="newPackage.price" type="number" :step="100000" class="w-full" color="warning"/>
              </UFormField>
              <div class="grid grid-cols-2 gap-4">
                <UFormField label="Sessions" required>
                  <UInput v-model="newPackage.sessions" type="number" class="w-full" color="warning" />
                </UFormField>
                <UFormField label="Duration (min)" required>
                  <UInput v-model="newPackage.duration" type="number" :step="15" class="w-full" color="warning" />
                </UFormField>
              </div>
              <UFormField label="Description">
                <UTextarea v-model="newPackage.description" class="w-full" color="warning" />
              </UFormField>
              <UFormField label="Features">
                <template #hint>
                  <span>Masukkan satu fitur per baris.</span>
                </template>
                <UTextarea v-model="newPackage.features" placeholder="Free Trial&#10;15x Sessions" class="w-full" color="warning" />
              </UFormField>
              <USwitch v-model="newPackage.isPopular" label="Mark as Popular" class="w-full" color="warning" />
            </div>
          </template>
          <template #footer>
            <div class="flex justify-end gap-3">
              <UButton label="Cancel" variant="ghost" color="neutral" @click="showAddModal = false" />
              <UButton label="Create Package" icon="i-lucide-plus" color="warning" @click="saveNewPackage" />
            </div>
          </template>
        </UModal>

        <!-- Add Add-on Modal -->
        <UModal v-model:open="showAddonModal" title="Add New Add-on">
          <template #body>
            <div class="space-y-4">
              <UFormField label="Add-on Name" required>
                <UInput v-model="newAddon.name" placeholder="e.g. SIM A Assistance" class="w-full" color="warning"/>
              </UFormField>
              <UFormField label="Price (IDR)" required>
                <UInput v-model="newAddon.price" type="number" :step="50000" class="w-full" color="warning"/>
              </UFormField>
              <UFormField label="Description">
                <UTextarea v-model="newAddon.description" placeholder="Briefly describe what this add-on includes." class="w-full" color="warning" />
              </UFormField>
            </div>
          </template>
          <template #footer>
            <div class="flex justify-end gap-3">
              <UButton label="Cancel" variant="ghost" color="neutral" @click="showAddonModal = false" />
              <UButton label="Create Add-on" icon="i-lucide-plus" color="warning" @click="saveNewAddon" />
            </div>
          </template>
        </UModal>

        <!-- Edit Add-on Modal -->
        <UModal v-model:open="showEditAddonModal" title="Edit Add-on">
          <template #body>
            <div v-if="editingAddon" class="space-y-4">
              <UFormField label="Add-on Name" required>
                <UInput v-model="editingAddon.name" class="w-full" color="warning"/>
              </UFormField>
              <UFormField label="Price (IDR)" required>
                <UInput v-model="editingAddon.price" type="number" :step="50000" class="w-full" color="warning"/>
              </UFormField>
              <UFormField label="Description">
                <UTextarea v-model="editingAddon.description" placeholder="Briefly describe what this add-on includes." class="w-full" color="warning" />
              </UFormField>
            </div>
          </template>
          <template #footer>
            <div class="flex justify-end gap-3">
              <UButton label="Cancel" variant="ghost" color="neutral" @click="showEditAddonModal = false" />
              <UButton label="Save Changes" icon="i-lucide-save" color="warning" @click="saveEditedAddon" />
            </div>
          </template>
        </UModal>

        <!-- Add-ons Section -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">Package Add-ons</h2>
              <UButton label="Add Add-on" icon="i-lucide-plus" size="md" color="warning" variant="outline" @click="showAddonModal = true" />
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
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="md" @click="openEditAddonModal(addon)" />
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
