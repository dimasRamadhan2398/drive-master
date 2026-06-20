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
});

function handleClose() {
  emit("update:open", false);
  newAddon.value = { name: "", price: 0, description: "" };
}

function saveNewAddon() {
  if (!newAddon.value.name || newAddon.value.price <= 0) {
    toast.add({
      title: t('common.error'),
      description: t('admin.addon.namePriceRequired'),
      color: "error",
    });
    return;
  }

  packagesStore.addAddon({
    name: newAddon.value.name,
    price: newAddon.value.price,
    description: newAddon.value.description,
  });

  toast.add({
    title: t('admin.addon.added'),
    description: `"${newAddon.value.name}" ${t('admin.studentAddedDesc').replace('{name} telah ditambahkan.', 'telah dibuat.')}`,
    color: "success",
  });

  handleClose();
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
          :label="t('admin.addon.add')"
          icon="i-lucide-plus"
          color="warning"
          @click="saveNewAddon"
        />
      </div>
    </template>
  </UModal>
</template>