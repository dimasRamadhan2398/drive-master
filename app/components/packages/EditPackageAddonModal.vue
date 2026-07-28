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

const editingAddon = ref<{
  id: string;
  name: string;
  price: number;
  description: string;
  sessions: number;
} | null>(null);

const isSubmitting = ref(false);

watch(
  () => props.addon,
  (newAddon) => {
    if (newAddon) {
      editingAddon.value = {
        id: newAddon.id,
        name: newAddon.name,
        price: newAddon.price,
        description: newAddon.description,
        sessions: newAddon.sessions || 1,
      };
    }
  },
  { immediate: true },
);

function handleClose() {
  emit("update:open", false);
  editingAddon.value = null;
}

async function saveEditedAddon() {
  if (!editingAddon.value) return;

  isSubmitting.value = true;
  try {
    const result = await packagesStore.updateAddon(
      editingAddon.value.id,
      {
        name: editingAddon.value.name,
        price: editingAddon.value.price,
        description: editingAddon.value.description,
        sessions: editingAddon.value.sessions,
      },
    );

    if (result) {
      toast.add({
        title: t('admin.addon.updated'),
        description: `"${editingAddon.value.name}" telah disimpan.`,
        color: "success",
      });
      handleClose();
    } else {
      toast.add({
        title: t('common.error'),
        description: t('admin.addon.notFound'),
        color: "error",
      });
    }
  } catch (error) {
    console.error("Error updating addon:", error);
    toast.add({
      title: t('common.error'),
      description: "Gagal menyimpan add-on. Silakan coba lagi.",
      color: "error",
    });
  } finally {
    isSubmitting.value = false;
  }
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
        <UFormField :label="t('admin.package.sessions')">
          <UInput
            v-model="editingAddon.sessions"
            type="number"
            :min="1"
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
          :disabled="isSubmitting"
          @click="handleClose"
        />
        <UButton
          :label="t('admin.saveChanges')"
          icon="i-lucide-save"
          color="warning"
          :loading="isSubmitting"
          @click="saveEditedAddon"
        />
      </div>
    </template>
  </UModal>
</template>