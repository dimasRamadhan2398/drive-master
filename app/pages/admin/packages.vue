<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { ref, computed } from "vue";
import { usePackagesStore } from "~/stores/packages";
import type { Package, Addon } from "~/stores/packages";
import AddPackageAddonModal from "~/components/packages/AddPackageAddonModal.vue";
import EditPackageAddonModal from "~/components/packages/EditPackageAddonModal.vue";
import AddPackageModal from "~/components/packages/AddPackageModal.vue";
import EditPackageModal from "~/components/packages/EditPackageModal.vue";
import type { CreatePackageData } from "~/services/packageService";

const { t } = useI18n()
definePageMeta({ layout: "admin" });

const toast = useToast();
const packagesStore = usePackagesStore();

const showEditModal = ref(false);
const showAddModal = ref(false);
const showAddonModal = ref(false);
const showEditAddonModal = ref(false);
const selectedPackage = ref<Package | null>(null);
const editingAddon = ref<Addon | null>(null);

const packages = computed(() => packagesStore.packages);
const addOns = computed(() => packagesStore.addons);

// Stats computed from store
const totalPackages = computed(() => packagesStore.totalPackages);
const totalSold = computed(() =>
  packagesStore.packages.reduce((sum, p) => sum + p.totalSold, 0),
);
const totalRevenue = computed(() => packagesStore.totalRevenue);

function formatPrice(price: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(price);
}

function editPackage(pkg: Package) {
  selectedPackage.value = {
    ...pkg,
    features: pkg.features,
  };
  showEditModal.value = true;
}

async function togglePackageStatus(pkgId: string) {
  await packagesStore.togglePackageStatus(pkgId);

  var isActive = packagesStore.packages.find((p) => p.id === pkgId)?.isActive;
  if (isActive !== null) {
    toast.add({
      title: isActive ? "Package Activated" : "Package Deactivated",
      icon: isActive ? "i-lucide-check-circle" : "i-lucide-x-circle",
      color: isActive ? "success" : "warning",
    });
  }
}

function duplicatePackage(pkg: Package) {
  packagesStore.duplicatePackage(pkg.id);
  toast.add({
    title: "Package Duplicated",
    description: `Salinan dari "${pkg.name}" telah dibuat.`,
    color: "success",
  });
}

function deletePackage(pkgId: string) {
  if (
    confirm(
      `Anda yakin ingin menghapus paket ini? Aksi ini tidak dapat dibatalkan.`,
    )
  ) {
    packagesStore.deletePackage(pkgId);
    toast.add({
      title: "Package Deleted",
      description: `Paket telah dihapus.`,
      color: "error",
      icon: "i-lucide-trash",
    });
  }
}

function viewPackageSales(pkg: Package) {
  navigateTo(`/admin/sales?packageId=${pkg.id}`);
}

function openEditAddonModal(addon: Addon) {
  editingAddon.value = { ...addon };
  showEditAddonModal.value = true;
}

async function handleAddPackage(pkg: CreatePackageData) {
  try {
    const created = await packagesStore.addPackage(pkg);
    toast.add({
      title: "Paket Ditambahkan",
      description: `"${created.name}" telah dibuat.`,
      color: "success",
    });
    showAddModal.value = false;
  } catch (error) {
    toast.add({
      title: t('common.error'),
      description: "Gagal membuat paket. Silakan coba lagi.",
      color: "error",
    });
  }
}

async function deleteAddon(addonId: string) {
  if (
    confirm(
      `Anda yakin ingin menghapus add-on ini? Aksi ini tidak dapat dibatalkan.`,
    )
  ) {
    const success = await packagesStore.deleteAddon(addonId);
    if (success) {
      toast.add({
        title: "Add-on Dihapus",
        description: `Add-on telah dihapus.`,
        color: "error",
        icon: "i-lucide-trash",
      });
    } else {
      toast.add({
        title: t('common.error'),
        description: "Gagal menghapus add-on. Silakan coba lagi.",
        color: "error",
      });
    }
  }
}

onMounted(() => {
  packagesStore.fetchPackages();
  packagesStore.fetchAddons();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.packages')">
        <template #right>
          <UButton
            icon="i-lucide-plus"
            color="warning"
            :label="t('admin.addNew')"
            @click="showAddModal = true"
          />
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
                <p class="text-md text-muted">Total {{ t('admin.packages') }}</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-green-500/10">
                <UIcon
                  name="i-lucide-shopping-cart"
                  class="size-6 text-green-500"
                />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ totalSold }}</p>
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
                <p class="text-2xl font-bold">
                  {{ formatPrice(totalRevenue) }}
                </p>
                <p class="text-md text-muted">{{ t('admin.totalRevenue') }}</p>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Packages Grid -->
        <div class="grid md:grid-cols-2 gap-6">
          <UCard
            v-for="pkg in packages"
            :key="pkg.id"
            :class="{ 'ring-2 ring-warning': pkg.isPopular }"
          >
            <template #header>
              <div class="flex items-start justify-between">
                <div>
                  <div class="flex items-center gap-2">
                    <h3 class="text-xl font-bold">{{ pkg.name }}</h3>
                    <UBadge
                      v-if="pkg.isPopular"
                      :label="t('packages.popular')"
                      color="warning"
                      size="md"
                    />
                  </div>
                  <p class="text-2xl font-bold text-warning mt-2">
                    {{ formatPrice(pkg.price) }}
                  </p>
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
                  <p class="text-md text-muted">{{ t('billing.sessions') }}</p>
                </div>
                <div class="p-2 rounded-lg bg-muted/50 text-center">
                  <p class="font-bold">{{ pkg.duration }} min</p>
                  <p class="text-md text-muted">Per {{ t('history.session') }}</p>
                </div>
              </div>

              <ul class="space-y-2">
                <li
                  v-for="feature in pkg.features"
                  :key="feature"
                  class="flex items-center gap-2 text-md"
                >
                  <UIcon name="i-lucide-check" class="size-4 text-warning" />
                  {{ feature }}
                </li>
              </ul>

              <USeparator />

              <div class="flex items-center justify-between text-md">
                <span class="text-muted">{{ t('admin.unitsSold') }}</span>
                <span class="font-bold">{{ pkg.totalSold }} {{ t('admin.students').toLowerCase() }}</span>
              </div>
            </div>

            <template #footer>
              <div class="flex gap-2">
                <UButton
                  :label="t('common.edit')"
                  icon="i-lucide-pencil"
                  variant="outline"
                  color="neutral"
                  class="flex-1"
                  @click="editPackage(pkg)"
                />
                <UDropdownMenu
                  :items="[
                    [
                      {
                        label: t('admin.viewSales'),
                        icon: 'i-lucide-chart-bar',
                        onSelect: () => viewPackageSales(pkg),
                      },
                      {
                        label: t('admin.duplicate'),
                        icon: 'i-lucide-copy',
                        onSelect: () => duplicatePackage(pkg),
                      },
                    ],
                    [
                      {
                        label: t('common.delete'),
                        icon: 'i-lucide-trash',
                        color: 'error',
                        onSelect: () => deletePackage(pkg.id),
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
              </div>
            </template>
          </UCard>
        </div>

        <!-- Modals -->
        <!-- Edit Package Modal Component -->
        <EditPackageModal
          v-model:open="showEditModal"
          :package="selectedPackage"
        />

        <!-- Add Package Modal Component -->
        <AddPackageModal
          v-model:open="showAddModal"
          @saved="handleAddPackage"
        />

        <!-- Add Add-on Modal Component -->
        <AddPackageAddonModal v-model:open="showAddonModal" />

        <!-- Edit Add-on Modal Component -->
        <EditPackageAddonModal
          v-model:open="showEditAddonModal"
          :addon="editingAddon"
        />

        <!-- Add-ons Section -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">{{ t('admin.packageAddons') }}</h2>
              <UButton
                :label="t('admin.addon.add')"
                icon="i-lucide-plus"
                size="md"
                color="warning"
                variant="outline"
                @click="showAddonModal = true"
              />
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
                  <UBadge
                    :label="formatPrice(addon.price)"
                    variant="subtle"
                    color="warning"
                  />
                  <UBadge
                    v-if="addon.sessions && addon.sessions > 1"
                    :label="`${addon.sessions} ${t('billing.sessions')}`"
                    variant="subtle"
                    color="info"
                  />
                </div>
                <p class="text-md text-muted mt-1">{{ addon.description }}</p>
              </div>
              <div class="flex items-center gap-2">
                <div class="text-right mr-4">
                  <p class="font-bold">{{ addon.sold }}</p>
                  <p class="text-md text-muted">{{ t('admin.unitsSold').replace('Total ', '') }}</p>
                </div>
                <UButton
                  icon="i-lucide-pencil"
                  color="neutral"
                  variant="ghost"
                  size="md"
                  @click="openEditAddonModal(addon)"
                />
                <UButton
                  icon="i-lucide-trash"
                  color="error"
                  variant="ghost"
                  size="md"
                  @click="deleteAddon(addon.id)"
                />
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
