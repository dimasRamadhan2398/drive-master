<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import type { CreatePackageData } from "~/services/packageService";
import { packageService } from "~/services/packageService";

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", data: any): void;
}>();

const toast = useToast();
const isAddLoading = ref(false);

const packageTypes = [
  { label: "Bronze", value: "bronze" },
  { label: "Silver", value: "silver" },
  { label: "Gold", value: "gold" },
  { label: "Platinum", value: "platinum" },
];

const newPackage = ref({
  name: "",
  price: 0,
  discountPrice: 0,
  isDiscountActive: false,
  packageType: "bronze",
  sessions: 0,
  duration: 60,
  description: "",
  features: "",
  highlight: false,
});

function handleClose() {
  emit("update:open", false);
  newPackage.value = {
    name: "",
    price: 0,
    discountPrice: 0,
    isDiscountActive: false,
    packageType: "bronze",
    sessions: 0,
    duration: 60,
    description: "",
    features: "",
    highlight: false,
  };
}

async function saveNewPackage(e: Event) {
  e.preventDefault();
  
  if (!newPackage.value.name || newPackage.value.price <= 0) {
    toast.add({
      title: "Error",
      description: "Nama paket dan harga yang valid wajib diisi.",
      color: "error",
    });
    return;
  }

  // Convert features string to array (one per line)
  const featuresArray = newPackage.value.features
    ? newPackage.value.features.split("\n").filter((f) => f.trim() !== "")
    : [];

  const pkg: CreatePackageData = {
    name: newPackage.value.name,
    description: newPackage.value.description,
    packageType: newPackage.value.packageType as "bronze" | "silver" | "gold" | "platinum",
    price: newPackage.value.price,
    discountPrice: newPackage.value.isDiscountActive ? newPackage.value.discountPrice : undefined,
    durationMinutes: newPackage.value.duration,
    totalSessions: newPackage.value.sessions,
    benefits: featuresArray,
    highlight: newPackage.value.highlight,
    isDiscounted: newPackage.value.isDiscountActive,
  };

  emit("saved", pkg);
  // isAddLoading.value = true;

  // try {
  //   const created = await packageService.create(pkg);
  //   toast.add({
  //     title: "Paket Ditambahkan",
  //     description: `"${created.name}" telah dibuat.`,
  //     color: "success",
  //   });
  //   emit("saved", created);
  //   handleClose();
  // } catch (error) {
  //   toast.add({
  //     title: "Error",
  //     description: "Gagal membuat paket. Silakan coba lagi.",
  //     color: "error",
  //   });
  // } finally {
  //   isAddLoading.value = false;
  // }
}
</script>

<template>
  <UModal
    :open="open"
    :title="t('admin.package.add')"
    @update:open="(val) => emit('update:open', val)"
  >
    <template #body>
      <div class="space-y-4">
        <UFormField :label="t('admin.package.name')" required>
          <UInput
            v-model="newPackage.name"
            placeholder="e.g. 15x Sessions"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.type')" required>
          <USelect
            v-model="newPackage.packageType"
            :items="packageTypes"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <div class="grid grid-cols-2 gap-4">
          <UFormField :label="t('admin.package.price')" required>
            <UInput
              v-model="newPackage.price"
              type="number"
              :step="100000"
              class="w-full"
              color="warning"
            />
          </UFormField>
          <UFormField :label="t('admin.package.discountPrice')">
            <UInput
              v-model="newPackage.discountPrice"
              type="number"
              :step="100000"
              class="w-full"
              color="warning"
              :disabled="!newPackage.isDiscountActive"
            />
          </UFormField>
        </div>
        <USwitch
          v-model="newPackage.isDiscountActive"
          :label="t('admin.package.enableDiscount')"
          class="w-full"
          color="warning"
        />
        <div class="grid grid-cols-2 gap-4">
          <UFormField :label="t('admin.package.sessions')" required>
            <UInput
              v-model="newPackage.sessions"
              type="number"
              class="w-full"
              color="warning"
            />
          </UFormField>
          <UFormField :label="t('admin.package.duration')" required>
            <UInput
              v-model="newPackage.duration"
              type="number"
              :step="15"
              class="w-full"
              color="warning"
            />
          </UFormField>
        </div>
        <UFormField :label="t('admin.package.description')">
          <UTextarea
            v-model="newPackage.description"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.features')">
          <template #hint>
            <span>{{ t('admin.package.featuresHint') }}</span>
          </template>
          <UTextarea
            v-model="newPackage.features"
            placeholder="Free Trial&#10;15x Sessions"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <USwitch
          v-model="newPackage.highlight"
          :label="t('admin.package.markPopular')"
          class="w-full"
          color="warning"
        />
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton
          :label="t('common.cancel')"
          variant="ghost"
          color="neutral"
          @click="handleClose"
        />
        <UButton
          :label="t('admin.package.add')"
          icon="i-lucide-plus"
          color="warning"
          :loading="isAddLoading"
          :disabled="isAddLoading"
          @click="saveNewPackage"
        />
      </div>
    </template>
  </UModal>
</template>
