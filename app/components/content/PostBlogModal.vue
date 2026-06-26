<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/index.js";
import type {
  CreateBlogArticleData,
} from "../../services/contentService";
import RichTextEditor from "./RichTextEditor.vue";
import { computed, ref, watch } from "vue";
import { useAuthStore } from "../../stores/auth";

export interface PostFormData {
  id: number;
  title: string;
  slug: string;
  author: string;
  leadParagraph: string;
  content: string;
  status: "draft" | "published" | "archived";
  featuredImage: string;
  publishing: {
    status: "draft" | "published" | "archived";
    publishedAt?: string;
    scheduledAt?: string;
  };
  attractiveness: {
    isFeatured: boolean;
    isSpotlight: boolean;
    priority: number;
    highlight: boolean;
  };
}

const props = defineProps<{
  open: boolean;
  post?: PostFormData | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", data: CreateBlogArticleData): void;
  (e: "error", message: string): void;
}>();

const toast = useToast();
const authStore = useAuthStore();
const isSaving = ref(false);
const isEditing = computed(() => !!props.post?.id);
const validationError = ref<string | null>(null);

const ACCEPTED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/webp"];
const MAX_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB

const mediaInputRef = ref<HTMLInputElement | null>(null);
const featuredImagePreview = ref<string | null>(null);

const postForm = ref<PostFormData>({
  id: 0,
  title: "",
  slug: "",
  author: "Admin",
  leadParagraph: "",
  content: "",
  status: "draft",
  featuredImage: "",
  publishing: {
    status: "draft",
    publishedAt: undefined,
    scheduledAt: undefined,
  },
  attractiveness: {
    isFeatured: false,
    isSpotlight: false,
    priority: 0,
    highlight: false,
  },
});

function resetForm() {
  postForm.value = {
    id: 0,
    title: "",
    slug: "",
    author: "Admin",
    leadParagraph: "",
    content: "",
    status: "draft",
    featuredImage: "",
    publishing: {
      status: "draft",
      publishedAt: undefined,
      scheduledAt: undefined,
    },
    attractiveness: {
      isFeatured: false,
      isSpotlight: false,
      priority: 0,
      highlight: false,
    },
  };
  featuredImagePreview.value = null;
}

function initForm() {
  if (props.post) {
    // Editing mode: featuredImage is a URL from backend
    postForm.value = {
      id: props.post.id,
      title: props.post.title,
      slug: props.post.slug || "",
      author: props.post.author,
      leadParagraph: props.post.leadParagraph || "",
      content: props.post.content || "",
      status: props.post.status,
      featuredImage: props.post.featuredImage || "",
      publishing: {
        status: props.post.publishing.status,
        publishedAt: props.post.publishing.publishedAt,
        scheduledAt: props.post.publishing.scheduledAt,
      },
      attractiveness: {
        isFeatured: props.post.attractiveness.isFeatured,
        isSpotlight: props.post.attractiveness.isSpotlight,
        priority: props.post.attractiveness.priority,
        highlight: props.post.attractiveness.highlight,
      },
    };
    // For editing, use the URL directly as preview
    featuredImagePreview.value = props.post.featuredImage || null;
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

function triggerMediaUpload() {
  mediaInputRef.value?.click();
}

function handleMediaSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (!input.files || !input.files[0]) return;

  const file = input.files[0];

  if (!ACCEPTED_IMAGE_TYPES.includes(file.type)) {
    toast.add({
      title: "Invalid File",
      description: `"${file.name}" is not a supported format. Use JPG, PNG, or WebP.`,
      color: "error",
    });
    input.value = "";
    return;
  }

  if (file.size > MAX_IMAGE_SIZE) {
    toast.add({
      title: "File Too Large",
      description: `"${file.name}" exceeds the 5MB limit for images.`,
      color: "error",
    });
    input.value = "";
    return;
  }

  const reader = new FileReader();
  reader.onload = (e) => {
    // For create mode: store as base64 data URL (bytes)
    postForm.value.featuredImage = e.target?.result as string;
    featuredImagePreview.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);

  input.value = "";
}

function removeFeaturedImage() {
  postForm.value.featuredImage = "";
  featuredImagePreview.value = null;
}

// Generate slug from title (auto-generate if empty and not editing)
function generateSlug(title: string): string {
  return title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
}

// Parse API error and return user-friendly message
function parseApiError(error: any): string | null {
  if (!error) return null;

  // Handle various error formats
  const errorMessage = error?.message || error?.data?.message || "";

  // Check for duplicate title error
  if (
    errorMessage.toLowerCase().includes("duplicate") ||
    errorMessage.toLowerCase().includes("already exists") ||
    errorMessage.toLowerCase().includes("unique constraint") ||
    errorMessage.toLowerCase().includes("title")
  ) {
    return "A post with this title already exists. Please use a different title.";
  }

  // Generic error
  if (errorMessage) {
    return errorMessage;
  }

  return null;
}

watch(
  () => postForm.value.title,
  (newTitle) => {
    // Clear validation error when user changes title
    if (validationError.value) {
      validationError.value = null;
    }
    if (!isEditing.value && newTitle) {
      postForm.value.slug = generateSlug(newTitle);
    }
  }
);

async function savePost() {
  // Clear previous validation error
  validationError.value = null;

  if (!postForm.value.title.trim()) {
    validationError.value = "Post title is required.";
    toast.add({
      title: "Validation Error",
      description: "Post title is required.",
      color: "error",
    });
    return;
  }

  // Get authorId from authStore
  const authorId = authStore.currentUser?.userId || authStore.user?.userId;
  if (!authorId) {
    validationError.value = "Unable to identify author. Please log in again.";
    toast.add({
      title: "Error",
      description: "Unable to identify author. Please log in again.",
      color: "error",
    });
    return;
  }

  isSaving.value = true;

  try {
    const data: CreateBlogArticleData = {
      title: postForm.value.title,
      slug: postForm.value.slug || undefined,
      author: postForm.value.author,
      authorId: authorId,
      leadParagraph: postForm.value.leadParagraph || undefined,
      content: postForm.value.content || undefined,
      // For create mode: featuredImage is base64/dataURL (bytes)
      // For edit mode: featuredImage is a URL string (already stored)
      featuredImage: postForm.value.featuredImage || undefined,
      publishing: {
        status: postForm.value.publishing.status,
      },
      attractiveness: {
        isFeatured: postForm.value.attractiveness.isFeatured,
        isSpotlight: postForm.value.attractiveness.isSpotlight,
        priority: postForm.value.attractiveness.priority,
        highlight: postForm.value.attractiveness.highlight,
      },
      status: postForm.value.status,
    };

    emit("saved", data);
    emit("update:open", false);
    resetForm();
  } catch (error: any) {
    const errorMessage = parseApiError(error);
    if (errorMessage) {
      validationError.value = errorMessage;
      emit("error", errorMessage);
      toast.add({
        title: "Validation Error",
        description: errorMessage,
        color: "error",
      });
    }
  } finally {
    isSaving.value = false;
  }
}
</script>

<template>
  <UModal
    :open="open"
    title="Post Blog"
    @update:open="(val:any) => emit('update:open', val)"
  >
    <template #content>
      <div class="bg-default rounded-2xl w-full">
        <div class="px-6 py-4 border-b border-default flex items-center justify-between">
          <h3 class="text-base font-semibold">
            {{ isEditing ? "Edit Post" : "New Post" }}
          </h3>
          <UButton
            icon="i-lucide-x"
            color="neutral"
            variant="ghost"
            @click="emit('update:open', false)"
          />
        </div>
        <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
          <UFormField label="Post Title" required :error="validationError!">
            <UInput
              v-model="postForm.title"
              placeholder="e.g. How to Charge Your EV at Home"
              class="w-full"
            />
          </UFormField>
          <UFormField label="Slug">
            <template #hint>
              <span>URL-friendly version (auto-generated from title)</span>
            </template>
            <UInput
              v-model="postForm.slug"
              placeholder="e.g. how-to-charge-your-ev-at-home"
              class="w-full"
            />
          </UFormField>
          <UFormField label="Author">
            <UInput v-model="postForm.author" placeholder="Admin" class="w-full" />
          </UFormField>
          <UFormField label="Lead Paragraph">
            <template #hint>
              <span>Brief summary shown in blog listings</span>
            </template>
            <UTextarea
              v-model="postForm.leadParagraph"
              placeholder="Write a brief summary of your post..."
              class="w-full"
              :rows="3"
            />
          </UFormField>
          <UFormField label="Content">
            <RichTextEditor
              v-model="postForm.content"
              placeholder="Write your blog post content here..."
            />
          </UFormField>
          <UFormField label="Featured Image">
            <template #hint>
              <span>Single featured image for the blog post</span>
            </template>
            <input
              ref="mediaInputRef"
              type="file"
              accept="image/jpeg,image/png,image/webp"
              class="hidden"
              @change="handleMediaSelect"
            />

            <!-- Featured Image Preview -->
            <div v-if="featuredImagePreview" class="relative group rounded-lg border border-default overflow-hidden bg-muted/30">
              <img
                :src="featuredImagePreview"
                alt="Featured Image"
                class="w-full h-48 object-cover"
              />
              <div
                class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
              >
                <UButton
                  icon="i-lucide-trash-2"
                  color="error"
                  variant="ghost"
                  size="sm"
                  @click="removeFeaturedImage"
                  class="text-white"
                />
              </div>
              <div
                class="absolute bottom-0 left-0 right-0 bg-black/60 px-2 py-0.5 flex items-center justify-between"
              >
                <span class="text-[10px] text-white truncate">Featured Image</span>
                <UButton
                  icon="i-lucide-pencil"
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  @click="triggerMediaUpload"
                  class="text-white hover:text-white ml-auto"
                />
              </div>
            </div>

            <!-- Upload Button (shown when no image) -->
            <button
              v-else
              @click="triggerMediaUpload"
              class="w-full border-2 border-dashed border-default rounded-lg p-4 hover:border-primary hover:bg-primary/5 transition-colors cursor-pointer flex flex-col items-center gap-2"
            >
              <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
              <span class="text-sm font-medium text-muted">Click to add featured image</span>
              <span class="text-xs text-muted">JPG, PNG, WebP (max 5MB)</span>
            </button>
          </UFormField>

          <!-- Publishing Section -->
          <div class="border-t border-default pt-4">
            <h4 class="text-sm font-medium mb-3 flex items-center gap-2">
              <UIcon name="i-lucide-send" class="size-4" />
              Publishing
            </h4>
            <UFormField label="Status">
              <USelect
                v-model="postForm.publishing!.status"
                :items="[
                  { label: 'Draft', value: 'draft' },
                  { label: 'Published', value: 'published' },
                  { label: 'Archived', value: 'archived' },
                ]"
                class="w-full"
              />
            </UFormField>
          </div>

          <!-- Attractiveness Section -->
          <div class="border-t border-default pt-4">
            <h4 class="text-sm font-medium mb-3 flex items-center gap-2">
              <UIcon name="i-lucide-star" class="size-4" />
              Attractiveness
            </h4>
            <div class="space-y-3">
              <USwitch
                v-model="postForm.attractiveness!.isFeatured"
                label="Featured"
                color="warning"
              />
              <USwitch
                v-model="postForm.attractiveness!.isSpotlight"
                label="Spotlight"
                color="warning"
              />
              <USwitch
                v-model="postForm.attractiveness!.highlight"
                label="Highlighted"
                color="warning"
              />
              <UFormField label="Priority">
                <UInput
                  v-model="postForm.attractiveness!.priority"
                  type="number"
                  :min="0"
                  :max="100"
                  class="w-full"
                />
              </UFormField>
            </div>
          </div>
        </div>
        <div class="px-6 py-4 border-t border-default flex justify-end gap-3">
          <UButton
            :label="isEditing ? 'Save Changes' : 'Create Post'"
            icon="i-lucide-check"
            :loading="isSaving"
            :disabled="isSaving"
            @click="savePost"
          />
        </div>
      </div>
    </template>
  </UModal>
</template>
