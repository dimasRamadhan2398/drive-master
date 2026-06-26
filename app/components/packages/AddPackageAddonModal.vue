<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { usePackagesStore } from "~/stores/packages";

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
}>();

const { t } = useI18n()
const toast = useToast();
const packagesStore = usePackagesStore();

const newAddon = ref({
  name: "",
  price: 0,
  description: "",
  sessions: 1,
});

const isSubmitting = ref(false);

function handleClose() {
  emit("update:open", false);
  newAddon.value = { name: "", price: 0, description: "", sessions: 1 };
}

async function saveNewAddon() {
  if (!newAddon.value.name || newAddon.value.price <= 0) {
    toast.add({
      title: t('common.error'),
      description: t('admin.addon.namePriceRequired'),
      color: "error",
    });
    return;
  }

  isSubmitting.value = true;
  try {
    const result = await packagesStore.addAddon({
      name: newAddon.value.name,
      price: newAddon.value.price,
      description: newAddon.value.description,
      sessions: newAddon.value.sessions,
    });

    if (result) {
      toast.add({
        title: t('admin.addon.added'),
        description: `"${newAddon.value.name}" telah dibuat.`,
        color: "success",
      });
      handleClose();
    } else {
      toast.add({
        title: t('common.error'),
        description: "Gagal membuat add-on. Silakan coba lagi.",
        color: "error",
      });
    }
  } catch (error) {
    console.error("Error creating addon:", error);
    toast.add({
      title: t('common.error'),
      description: "Gagal membuat add-on. Silakan coba lagi.",
      color: "error",
    });
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <UModal :open="open" :title="t('admin.addon.add')" @update:open="(val) => emit('update:open', val)">
    <template #body>
      <div class="space-y-4">
        <UFormField :label="t('admin.addon.name')" required>
          <UInput
            v-model="newAddon.name"
            placeholder="e.g. SIM A Assistance"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.price')" required>
          <UInput
            v-model="newAddon.price"
            type="number"
            :step="50000"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.description')">
          <UTextarea
            v-model="newAddon.description"
            placeholder="Briefly describe what this add-on includes."
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.package.sessions')">
          <UInput
            v-model="newAddon.sessions"
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
          :label="t('admin.addon.add')"
          icon="i-lucide-plus"
          color="warning"
          :loading="isSubmitting"
          @click="saveNewAddon"
        />
      </div>
    </template>
  </UModal>
</template>