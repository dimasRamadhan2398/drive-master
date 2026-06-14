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
      title: "Add-on Updated",
      description: `"${editingAddon.value.name}" has been saved.`,
      color: "success",
    });
  } else {
    toast.add({
      title: "Error",
      description: "Add-on not found.",
      color: "error",
    });
  }

  handleClose();
}
</script>

<template>
  <UModal :open="open" title="Edit Add-on" @update:open="(val) => emit('update:open', val)">
    <template #body>
      <div v-if="editingAddon" class="space-y-4">
        <UFormField label="Add-on Name" required>
          <UInput
            v-model="editingAddon.name"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField label="Price (IDR)" required>
          <UInput
            v-model="editingAddon.price"
            type="number"
            :step="50000"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField label="Description">
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
          label="Cancel"
          variant="ghost"
          color="neutral"
          @click="handleClose"
        />
        <UButton
          label="Save Changes"
          icon="i-lucide-save"
          color="warning"
          @click="saveEditedAddon"
        />
      </div>
    </template>
  </UModal>
</template>