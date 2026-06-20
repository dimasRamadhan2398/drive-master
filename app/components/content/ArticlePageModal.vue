<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";

export interface PageFormData {
  id: number;
  title: string;
  slug: string;
  status: "draft" | "published";
}

export interface CreatePageData {
  title: string;
  slug: string;
  status: "draft" | "published";
}

const { t } = useI18n()
const props = defineProps<{
  open: boolean;
  page?: PageFormData | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", data: CreatePageData): void;
}>();

const toast = useToast();
const isSaving = ref(false);
const isEditing = computed(() => !!props.page?.id);

const pageForm = ref<PageFormData>({
  id: 0,
  title: "",
  slug: "",
  status: "draft",
});

function generateSlug(title: string) {
  return title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
}

function resetForm() {
  pageForm.value = {
    id: 0,
    title: "",
    slug: "",
    status: "draft",
  };
}

function initForm() {
  if (props.page) {
    pageForm.value = {
      id: props.page.id,
      title: props.page.title,
      slug: props.page.slug,
      status: props.page.status,
    };
  } else {
    resetForm();
  }
}

watch(
  () => props.open,
  (newVal) => {
    if (newVal) {
      initForm();
    }
  }
);

async function savePage() {
  if (!pageForm.value.title.trim()) {
    toast.add({
      title: "Error",
      description: "Page title is required.",
      color: "error",
    });
    return;
  }

  isSaving.value = true;

  try {
    const data: CreatePageData = {
      title: pageForm.value.title,
      slug: pageForm.value.slug || generateSlug(pageForm.value.title),
      status: pageForm.value.status,
    };

    emit("saved", data);
    emit("update:open", false);
    resetForm();
  } finally {
    isSaving.value = false;
  }
}
</script>

<template>
  <UModal
    :open="open"
    :title="isEditing ? 'Edit Page' : 'New Page'"
    @update:open="(val) => emit('update:open', val)"
  >
    <template #content>
      <div class="bg-default rounded-2xl w-full">
        <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
          <UFormField label="Page Title" required>
            <UInput
              v-model="pageForm.title"
              placeholder="e.g. Contact Us"
              class="w-full"
            />
          </UFormField>
          <UFormField label="Slug">
            <UInput
              v-model="pageForm.slug"
              :placeholder="pageForm.title ? generateSlug(pageForm.title) : '/contact-us'"
              class="w-full"
            />
            <template #hint>
              <span>{{ t('admin.content.slugHint') }}</span>
            </template>
          </UFormField>
          <UFormField label="Status">
            <USelect
              v-model="pageForm.status"
              :items="['draft', 'published']"
              class="w-full"
            />
          </UFormField>
        </div>
        <div class="px-6 py-4 border-t border-default flex justify-end gap-3">
          <UButton
            :label="isEditing ? 'Save Changes' : 'Create Page'"
            icon="i-lucide-check"
            :loading="isSaving"
            :disabled="isSaving"
            @click="savePage"
          />
        </div>
      </div>
    </template>
  </UModal>
</template>
