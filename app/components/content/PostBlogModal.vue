<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/index.js";
import { computed, ref, watch } from "vue";
import { useAuthStore } from "../../stores/auth";
import RichTextEditor from "./RichTextEditor.vue";
import type { Publishing, Attractiveness } from "../../services/contentService";

// ============================================================
// PostFormData – what the parent passes in for edit mode
// ============================================================
export interface PostFormData {
  id: string;
  title: string;
  slug: string;
  author: string;
  leadParagraph: string;
  content: string;
  featuredImage: string;
  publishing: Publishing;
  attractiveness: Attractiveness;
}

// ============================================================
// SavedPayload – emitted to the parent upon save
// ============================================================
export interface SavedPayload {
  id?: string; // present in edit mode
  title: string;
  slug?: string;
  author?: string;
  authorId: string;
  leadParagraph?: string;
  content?: string;
  /** data URI string when a new file was selected; existing URL when not changed */
  featuredImage?: string;
  /** Original filename of the newly selected image (for format validation) */
  featuredImageFileName?: string;
  publishing: Publishing;
  attractiveness: Attractiveness;
}

const props = defineProps<{
  open: boolean;
  post?: PostFormData | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", data: SavedPayload): void;
  (e: "error", message: string): void;
}>();

const toast = useToast();
const authStore = useAuthStore();
const isSaving = ref(false);
const isEditing = computed(() => !!props.post?.id);
const validationError = ref<string | null>(null);

const ACCEPTED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/webp"];
const MAX_IMAGE_SIZE = 5 * 1024 * 1024; // 5 MB

const mediaInputRef = ref<HTMLInputElement | null>(null);
const featuredImagePreview = ref<string | null>(null);
/** Filename of the newly selected image file (undefined = no new file selected) */
const newImageFileName = ref<string | undefined>(undefined);

// ============================================================
// Form state
// ============================================================
interface FormState {
  id: string;
  title: string;
  slug: string;
  author: string;
  leadParagraph: string;
  content: string;
  featuredImage: string;
  publishing: Publishing;
  attractiveness: Attractiveness;
}

const postForm = ref<FormState>({
  id: "",
  title: "",
  slug: "",
  author: "Admin",
  leadParagraph: "",
  content: "",
  featuredImage: "",
  publishing: { status: "draft", publishedAt: null, scheduledAt: null },
  attractiveness: {
    isFeatured: false,
    isSpotlight: false,
    priority: 0,
    highlight: false,
  },
});

function resetForm() {
  postForm.value = {
    id: "",
    title: "",
    slug: "",
    author: "Admin",
    leadParagraph: "",
    content: "",
    featuredImage: "",
    publishing: { status: "draft", publishedAt: null, scheduledAt: null },
    attractiveness: {
      isFeatured: false,
      isSpotlight: false,
      priority: 0,
      highlight: false,
    },
  };
  featuredImagePreview.value = null;
  newImageFileName.value = undefined;
  validationError.value = null;
}

function initForm() {
  if (props.post) {
    postForm.value = {
      id: props.post.id,
      title: props.post.title,
      slug: props.post.slug || "",
      author: props.post.author || "Admin",
      leadParagraph: props.post.leadParagraph || "",
      content: props.post.content || "",
      featuredImage: props.post.featuredImage || "",
      publishing: {
        status: props.post.publishing?.status || "draft",
        publishedAt: props.post.publishing?.publishedAt ?? null,
        scheduledAt: props.post.publishing?.scheduledAt ?? null,
      },
      attractiveness: {
        isFeatured: props.post.attractiveness?.isFeatured ?? false,
        isSpotlight: props.post.attractiveness?.isSpotlight ?? false,
        priority: props.post.attractiveness?.priority ?? 0,
        highlight: props.post.attractiveness?.highlight ?? false,
      },
    };
    // Show existing image as preview
    featuredImagePreview.value = props.post.featuredImage || null;
    newImageFileName.value = undefined;
  } else {
    resetForm();
  }
  validationError.value = null;
}

watch(
  () => props.open,
  (newVal) => {
    if (newVal) initForm();
  }
);

// ============================================================
// Image handling
// ============================================================
function triggerMediaUpload() {
  mediaInputRef.value?.click();
}

function handleMediaSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (!input.files || !input.files[0]) return;

  const file = input.files[0]!;

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
      description: `"${file.name}" exceeds the 5 MB limit for images.`,
      color: "error",
    });
    input.value = "";
    return;
  }

  // Store original filename for backend format validation
  newImageFileName.value = file.name;

  const reader = new FileReader();
  reader.onload = (e) => {
    // Store full data URI; the service/store will strip the prefix for upload
    postForm.value.featuredImage = e.target?.result as string;
    featuredImagePreview.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
  input.value = "";
}

function removeFeaturedImage() {
  postForm.value.featuredImage = "";
  featuredImagePreview.value = null;
  newImageFileName.value = undefined;
}

// ============================================================
// Slug generation
// ============================================================
function generateSlug(title: string): string {
  return title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
}

watch(
  () => postForm.value.title,
  (newTitle) => {
    if (validationError.value) validationError.value = null;
    if (!isEditing.value && newTitle) {
      postForm.value.slug = generateSlug(newTitle);
    }
  }
);

// ============================================================
// Save
// ============================================================
async function savePost() {
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
    const payload: SavedPayload = {
      ...(isEditing.value ? { id: postForm.value.id } : {}),
      title: postForm.value.title,
      slug: postForm.value.slug || undefined,
      author: postForm.value.author || undefined,
      authorId,
      leadParagraph: postForm.value.leadParagraph || undefined,
      content: postForm.value.content || undefined,
      featuredImage: postForm.value.featuredImage || undefined,
      featuredImageFileName: newImageFileName.value,
      publishing: { ...postForm.value.publishing },
      attractiveness: { ...postForm.value.attractiveness },
    };

    emit("saved", payload);
    emit("update:open", false);
    resetForm();
  } catch (error: any) {
    const msg = error?.message || "An unexpected error occurred.";
    validationError.value = msg;
    emit("error", msg);
    toast.add({ title: "Error", description: msg, color: "error" });
  } finally {
    isSaving.value = false;
  }
}
</script>

<template>
  <UModal
    :open="open"
    title="Post Blog"
    @update:open="(val: any) => emit('update:open', val)"
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
          <!-- Title -->
          <UFormField label="Post Title" required :error="validationError!">
            <UInput
              v-model="postForm.title"
              placeholder="e.g. How to Charge Your EV at Home"
              class="w-full"
            />
          </UFormField>

          <!-- Slug -->
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

          <!-- Author -->
          <UFormField label="Author">
            <UInput v-model="postForm.author" placeholder="Admin" class="w-full" />
          </UFormField>

          <!-- Lead Paragraph -->
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

          <!-- Content -->
          <UFormField label="Content">
            <RichTextEditor
              v-model="postForm.content"
              placeholder="Write your blog post content here..."
            />
          </UFormField>

          <!-- Featured Image -->
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

            <!-- Preview -->
            <div
              v-if="featuredImagePreview"
              class="relative group rounded-lg border border-default overflow-hidden bg-muted/30"
            >
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
                  class="text-white"
                  @click="removeFeaturedImage"
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
                  class="text-white hover:text-white ml-auto"
                  @click="triggerMediaUpload"
                />
              </div>
            </div>

            <!-- Upload Button -->
            <button
              v-else
              class="w-full border-2 border-dashed border-default rounded-lg p-4 hover:border-primary hover:bg-primary/5 transition-colors cursor-pointer flex flex-col items-center gap-2"
              @click="triggerMediaUpload"
            >
              <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
              <span class="text-sm font-medium text-muted"
                >Click to add featured image</span
              >
              <span class="text-xs text-muted">JPG, PNG, WebP (max 5 MB)</span>
            </button>
          </UFormField>

          <!-- Publishing -->
          <div class="border-t border-default pt-4">
            <h4 class="text-sm font-medium mb-3 flex items-center gap-2">
              <UIcon name="i-lucide-send" class="size-4" />
              Publishing
            </h4>
            <UFormField label="Status">
              <USelect
                v-model="postForm.publishing.status"
                :items="[
                  { label: 'Draft', value: 'draft' },
                  { label: 'Published', value: 'published' },
                  { label: 'Archived', value: 'archived' },
                ]"
                class="w-full"
              />
            </UFormField>
          </div>

          <!-- Attractiveness -->
          <div class="border-t border-default pt-4">
            <h4 class="text-sm font-medium mb-3 flex items-center gap-2">
              <UIcon name="i-lucide-star" class="size-4" />
              Attractiveness
            </h4>
            <div class="space-y-3">
              <USwitch
                v-model="postForm.attractiveness.isFeatured"
                label="Featured"
                color="warning"
              />
              <USwitch
                v-model="postForm.attractiveness.isSpotlight"
                label="Spotlight"
                color="warning"
              />
              <USwitch
                v-model="postForm.attractiveness.highlight"
                label="Highlighted"
                color="warning"
              />
              <UFormField label="Priority">
                <UInput
                  v-model="postForm.attractiveness.priority"
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
