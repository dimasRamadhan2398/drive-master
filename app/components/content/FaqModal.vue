<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { computed, ref, watch } from "vue";

export interface FaqFormData {
  id: string;
  question: string;
  answer: string;
  sortOrder: number;
}

export interface CreateFaqData {
  question: string;
  answer: string;
  sortOrder?: number;
}

const props = defineProps<{
  open: boolean;
  faq?: FaqFormData | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", data: CreateFaqData): void;
}>();

const toast = useToast();
const isSaving = ref(false);
const isEditing = computed(() => !!props.faq?.id);

const faqForm = ref<FaqFormData>({
  id: "",
  question: "",
  answer: "",
  sortOrder: 0,
});

function resetForm() {
  faqForm.value = {
    id: "",
    question: "",
    answer: "",
    sortOrder: 0,
  };
}

function initForm() {
  if (props.faq) {
    faqForm.value = {
      id: props.faq.id,
      question: props.faq.question,
      answer: props.faq.answer,
      sortOrder: props.faq.sortOrder,
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

async function saveFaq() {
  if (!faqForm.value.question.trim() || !faqForm.value.answer.trim()) {
    toast.add({
      title: "Error",
      description: "Question and answer are required.",
      color: "error",
    });
    return;
  }

  isSaving.value = true;

  try {
    const data: CreateFaqData = {
      question: faqForm.value.question,
      answer: faqForm.value.answer,
      sortOrder: faqForm.value.sortOrder,
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
    :title="isEditing ? 'Edit FAQ' : 'Add New FAQ'"
    @update:open="(val: any) => emit('update:open', val)"
  >
    <template #content>
      <div class="bg-default rounded-2xl w-full">
        <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
          <UFormField label="Question" required>
            <UInput
              v-model="faqForm.question"
              placeholder="e.g. What payment methods do you accept?"
              class="w-full"
            />
          </UFormField>
          <UFormField label="Answer" required>
            <UTextarea
              v-model="faqForm.answer"
              placeholder="Type the answer here..."
              :rows="4"
              class="w-full"
            />
          </UFormField>
        </div>
        <div class="px-6 py-4 border-t border-default flex justify-end gap-3">
          <UButton
            :label="isEditing ? 'Save Changes' : 'Add FAQ'"
            icon="i-lucide-check"
            :loading="isSaving"
            :disabled="isSaving"
            @click="saveFaq"
          />
        </div>
      </div>
    </template>
  </UModal>
</template>
