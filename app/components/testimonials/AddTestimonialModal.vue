<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { useStudentsStore } from "~/stores/students";
import { useTestimonialsStore } from "~/stores/testimonials";

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", data: any): void;
}>();

const open = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value),
});

const toast = useToast();
const studentsStore = useStudentsStore();
const { t } = useI18n()
const testimonialStore = useTestimonialsStore();
const authStore = useAuthStore();
const isSubmitting = ref(false);

// Student search results formatted for UInputMenu
const studentSuggestions = computed(() => {
  if (!studentsStore.searchResults?.length) return [];
  return studentsStore.searchResults.map((s) => ({
    label: s.name + (s.email ? ` (${s.email})` : ""),
    value: s.name,
  }));
});

// Selected student from search
const selectedStudentOption = ref<
  { label: string; value: string } | undefined
>();

function ratingOptions(): { label: string; value: number }[] {
  var list = [1, 2, 3, 4, 5].map((n) => ({
    label: `${n} Star${n > 1 ? "s" : ""}`,
    value: n,
  }));
  console.log("ada brp list", list.length);
  return list;
}

const form = ref({
  userId: "",
  userName: "",
  userRole: "Student",
  userImage: "",
  content: "",
  rating: 5,
  tags: "",
  status: "pending" as "draft" | "pending" | "published",
  isFeatured: false,
});

// Reset form when modal closes
watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      form.value = {
        userId: "",
        userName: "",
        userRole: "Student",
        userImage: "",
        content: "",
        rating: 5,
        tags: "",
        status: "pending",
        isFeatured: false,
      };
      selectedStudentOption.value = undefined;
      studentsStore.clearSearchResults();
    }
  },
);

function onSearchInput(term: string) {
  // Clear selected option when user starts typing
  selectedStudentOption.value = undefined;

  if (term && term.length >= 2) {
    studentsStore.searchStudents(term);
  } else if (!term) {
    studentsStore.clearSearchResults();
  }
}

function onStudentSelected(opt: { label: string; value: string } | undefined) {
  if (opt) {
    form.value.userName = opt.value;
    const student = studentsStore.searchResults.find(
      (s) => s.name === opt.value,
    );
    if (student) {
      form.value.userId = student.id;
    }
  }
}

async function handleSave() {
  if (!form.value.userName || !form.value.content) {
    toast.add({
      title: t('common.error'),
      description: t('admin.testimonial.nameContentRequired'),
      color: "error",
    });
    return;
  }

  if (!form.value.userId) {
    toast.add({
      title: t('common.error'),
      description: t('admin.testimonial.selectStudent'),
      color: "error",
    });
    return;
  }

  isSubmitting.value = true;

  try {
    const testimonial = await testimonialStore.addTestimonial({
      userId: form.value.userId,
      userName: form.value.userName,
      userImage: form.value.userImage,
      userRole: form.value.userRole,
      content: form.value.content,
      rating: form.value.rating,
      tags: form.value.tags,
      status: form.value.status,
      isFeatured: form.value.isFeatured,
      addedBy: authStore.user?.userId || "unknown",
      addedAt: new Date().toISOString(),
      sortOrder: 0,
    });

    if (testimonial) {
      toast.add({
        title: t('common.success'),
        description: t('admin.testimonial.created'),
        color: "success",
      });
      emit("saved", testimonial);
      emit("update:open", false);
    } else {
      toast.add({
        title: t('common.error'),
        description: "Failed to create testimonial. Please try again.",
        color: "error",
      });
    }
  } catch (error) {
    console.error("Error creating testimonial:", error);
    toast.add({
      title: t('common.error'),
      description: "An error occurred while creating the testimonial.",
      color: "error",
    });
  } finally {
    isSubmitting.value = false;
  }
}

const modalUi = {
  wrapper: "overflow-visible",
  container: "overflow-visible",
  body: "overflow-visible min-h-160",
};

// Custom UI for USelect to prevent dropdown clipping
const selectUi = {
  content:
    "max-h-[min(15rem,var(--reka-select-content-available-height,15rem))] w-full bg-default shadow-lg rounded-md ring ring-default overflow-visible z-50",
};

function handleClose() {
  emit("update:open", false);
}
</script>

<template>
  <UModal
    v-model:open="open"
    :title="t('admin.testimonial.add')"
    :dismissible="true"
    :scrollable="true"
    :ui="modalUi"
  >
    <template #body>
      <div class="space-y-4">
        <UFormField :label="t('admin.testimonial.name')" required>
          <UInputMenu
            v-model="selectedStudentOption"
            placeholder="Customer name"
            :items="studentSuggestions"
            class="w-full"
            color="warning"
            :popper="{
              placement: 'top-start',
              strategy: 'fixed',
              sameWidth: true,
            }"
            @update:search-term="onSearchInput"
            @update:model-value="onStudentSelected"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.role')">
          <UInput
            v-model="form.userRole"
            placeholder="e.g. Student"
            class="w-full"
            color="warning"
            disabled
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.photoUrl')">
          <UInput
            v-model="form.userImage"
            placeholder="https://..."
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.content')" required>
          <UTextarea
            v-model="form.content"
            placeholder="Testimonial content..."
            class="w-full"
            color="warning"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.rating')">
          <USelect
            v-model="form.rating"
            :items="ratingOptions()"
            color="warning"
            :ui="selectUi"
          />
        </UFormField>
        <UFormField :label="t('admin.testimonial.tags')">
          <UInput
            v-model="form.tags"
            placeholder="e.g. SIM A,Professional (comma separated)"
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
            ]"
            class="w-full"
            color="warning"
            :popper="{ placement: 'bottom-start' }"
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
          :label="t('admin.testimonial.add')"
          icon="i-lucide-plus"
          color="warning"
          :loading="isSubmitting"
          @click="handleSave"
        />
      </div>
    </template>
  </UModal>
</template>
