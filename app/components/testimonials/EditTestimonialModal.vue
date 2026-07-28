<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import type { Testimonial } from "~/stores/testimonials";

const { t } = useI18n()
const props = defineProps<{
  open: boolean;
  testimonial: Testimonial | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", testimonial: Testimonial): void;
}>();

const open = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value),
});

const toast = useToast();

// Local copy for editing
const form = ref({
  userName: "",
  userRole: "",
  userImage: "",
  content: "",
  rating: 5,
  tags: "",
  status: "pending" as "draft" | "pending" | "published" | "archived",
  isFeatured: false,
});

// Sync form with prop when testimonial changes
watch(
  () => props.testimonial,
  (testimonial) => {
    if (testimonial) {
      form.value = {
        userName: testimonial.userName,
        userRole: testimonial.userRole,
        userImage: testimonial.userImage,
        content: testimonial.content,
        rating: testimonial.rating,
        tags: testimonial.tags,
        status: testimonial.status as "draft" | "pending" | "published" | "archived",
        isFeatured: testimonial.isFeatured,
      };
    }
  },
  { immediate: true }
);

function handleSave() {
  if (!props.testimonial) return;

  if (!form.value.userName || !form.value.content) {
    toast.add({
      title: t('common.error'),
      description: t('admin.testimonial.nameContentRequired'),
      color: "error",
    });
    return;
  }

  emit("saved", {
    ...props.testimonial,
    ...form.value,
  });
  emit("update:open", false);
}

function handleClose() {
  emit("update:open", false);
}
</script>

<template>
  <UModal v-model:open="open" :title="t('admin.testimonial.edit')">
    <template #body>
      <div v-if="testimonial" class="space-y-4">
        <UFormField :label="t('admin.testimonial.name')" required>
          <UInput
            v-model="form.userName"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.role')">
          <UInput
            v-model="form.userRole"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.photoUrl')">
          <UInput
            v-model="form.userImage"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.content')" required>
          <UTextarea
            v-model="form.content"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.rating')">
          <USelect
            v-model="form.rating"
            :items="
              [1, 2, 3, 4, 5].map((n) => ({
                label: `${n} Star${n > 1 ? 's' : ''}`,
                value: n,
              }))
            "
            class="w-40"
            color="warning"
            :popper="{ placement: 'bottom-start' }"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.tags')">
          <UInput
            v-model="form.tags"
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.status')">
          <USelect
            v-model="form.status"
            :items="[
              { label: 'Draft', value: 'draft' },
              { label: 'Pending', value: 'pending' },
              { label: 'Published', value: 'published' },
              { label: 'Archived', value: 'archived' },
            ]"
            class="w-full"
            color="warning"
            :popper="{ placement: 'bottom-start' }"
          />
        </UFormField>
        <USwitch
          v-model="form.isFeatured"
          :label="t('admin.testimonial.isFeatured')"
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
          @click="handleSave"
        />
      </div>
    </template>
  </UModal>
</template>