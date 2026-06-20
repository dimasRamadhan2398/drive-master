<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { usePackagesStore } from "~/stores/packages";
import type { Addon } from "~/stores/packages";

const props = defineProps<{
  open: boolean;
  addon: Addon | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
}>();

const { t } = useI18n()
const toast = useToast();
const packagesStore = usePackagesStore();

const editingAddon = ref<Addon | null>(null);

watch(
  () => props.addon,
  (newAddon) => {
    if (newAddon) {
      editingAddon.value = { ...newAddon };
    }
  },
  { immediate: true },
);

function handleClose() {
  emit("update:open", false);
  editingAddon.value = null;
}

function saveEditedAddon() {
  if (!editingAddon.value) return;

  const updated = packagesStore.updateAddon(
    editingAddon.value.id,
    editingAddon.value,
  );
  if (updated) {
    toast.add({
      title: t('admin.addon.updated'),
      description: `"${editingAddon.value.name}" ${t('admin.studentUpdatedDesc').replace("{name}'s data has been updated.", "telah disimpan.")}`,
      color: "success",
    });
  } else {
    toast.add({
      title: t('common.error'),
      description: t('admin.addon.notFound'),
      color: "error",
    });
  }

  handleClose();
}
</script>

<template>
  <UModal :open="open" :title="t('admin.addon.edit')" @update:open="(val) => emit('update:open', val)">
    <template #body>
      <div v-if="editingAddon" class="space-y-4">
        <UFormField :label="t('admin.addon.name')" required>
          <UInput
            v-model="editingAddon.name"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.price')" required>
          <UInput
            v-model="editingAddon.price"
            type="number"
            :step="50000"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.description')">
          <UTextarea
            v-model="editingAddon.description"
            placeholder="Briefly describe what this add-on includes."
            class="w-full"
            color="warning"
          />
        </UFormField>
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
          @click="saveEditedAddon"
        />
      </div>
    </template>
  </UModal>
</template>