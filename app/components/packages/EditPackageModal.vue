<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { usePackagesStore } from "~/stores/packages";
import type { Package } from "~/stores/packages";

const { t } = useI18n()
const props = defineProps<{
  open: boolean;
  package: Package | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
}>();

const toast = useToast();
const packagesStore = usePackagesStore();

const selectedPackage = ref<Package | null>(null);
const formattedPackages = ref("");

watch(
  () => props.package,
  (newPackage) => {
    if (newPackage) {
      selectedPackage.value = { ...newPackage, features: newPackage.features };
      formattedPackages.value = newPackage.features.join("\n");
    }
  },
  { immediate: true },
);

watch(
  () => formattedPackages.value,
  (newValue) => {
    if (selectedPackage.value) {
      selectedPackage.value.features = newValue.split("\n").filter((f) => f.trim() !== "");
    }
  },
);

function handleClose() {
  emit("update:open", false);
  selectedPackage.value = null;
  formattedPackages.value = "";
}

function saveEditedPackage() {
  if (!selectedPackage.value) return;

  packagesStore.updatePackage(selectedPackage.value.id, {
    name: selectedPackage.value.name,
    price: selectedPackage.value.price,
    discountPrice: selectedPackage.value.price,
    description: selectedPackage.value.description,
    features: selectedPackage.value.features,
  });

  toast.add({
    title: "Package Updated",
    description: `"${selectedPackage.value.name}" has been saved.`,
    color: "success",
  });

  handleClose();
}
</script>

<template>
  <UModal :open="open" :title="t('admin.package.edit')" @update:open="(val) => emit('update:open', val)">
    <template #body>
      <div v-if="selectedPackage" class="space-y-4">
        <UFormField :label="t('admin.package.name')" required>
          <UInput
            v-model="selectedPackage.name"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.price')" required>
          <UInput
            v-model="selectedPackage.price"
            type="number"
            :step="100000"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <div class="grid grid-cols-2 gap-4">
          <UFormField :label="t('admin.package.sessions')" required>
            <UInput
              v-model="selectedPackage.sessions"
              type="number"
              class="w-full"
              color="warning"
            />
          </UFormField>
          <UFormField :label="t('admin.package.duration')" required>
            <UInput
              v-model="selectedPackage.duration"
              type="number"
              :step="15"
              class="w-full"
              color="warning"
            />
          </UFormField>
        </div>
        <UFormField :label="t('admin.package.description')">
          <UTextarea
            v-model="selectedPackage.description"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.features')">
          <template #hint>
            <span>{{ t('admin.package.featuresHint') }}</span>
          </template>
          <UTextarea
            v-model="formattedPackages"
            placeholder="Free Trial&#10;15x Sessions"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <USwitch
          v-model="selectedPackage.isPopular"
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
          :label="t('admin.saveChanges')"
          icon="i-lucide-save"
          color="warning"
          @click="saveEditedPackage"
        />
      </div>
    </template>
  </UModal>
</template>