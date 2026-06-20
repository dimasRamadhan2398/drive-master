<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import type {
  BlogPostMedia,
  Publishing,
  Attractiveness,
  CreateBlogPostData,
} from "~/services/contentService";
import RichTextEditor from "./RichTextEditor.vue";

export interface PostFormData {
  id: number;
  title: string;
  author: string;
  content: string;
  status: "draft" | "published" | "archived";
  media: BlogPostMedia[];
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

const { t } = useI18n()
const props = defineProps<{
  open: boolean;
  post?: PostFormData | null;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "saved", data: CreateBlogPostData): void;
}>();

const toast = useToast();
const isSaving = ref(false);
const isEditing = computed(() => !!props.post?.id);

const ACCEPTED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/webp"];
const ACCEPTED_VIDEO_TYPES = ["video/mp4", "video/webm"];
const MAX_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB
const MAX_VIDEO_SIZE = 50 * 1024 * 1024; // 50MB

const mediaInputRef = ref<HTMLInputElement | null>(null);

const postForm = ref<PostFormData>({
  id: 0,
  title: "",
  author: "Admin",
  content: "",
  status: "draft",
  media: [],
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
    author: "Admin",
    content: "",
    status: "draft",
    media: [],
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
}

function initForm() {
  if (props.post) {
    postForm.value = {
      id: props.post.id,
      title: props.post.title,
      author: props.post.author,
      content: props.post.content || "",
      status: props.post.status,
      media: props.post.media ? [...props.post.media] : [],
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

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

function triggerMediaUpload() {
  mediaInputRef.value?.click();
}

function handleMediaSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (!input.files) return;

  for (const file of Array.from(input.files)) {
    const isImage = ACCEPTED_IMAGE_TYPES.includes(file.type);
    const isVideo = ACCEPTED_VIDEO_TYPES.includes(file.type);

    if (!isImage && !isVideo) {
      toast.add({
        title: "Invalid File",
        description: `"${file.name}" is not a supported format. Use JPG, PNG, WebP for images or MP4, WebM for videos.`,
        color: "error",
      });
      continue;
    }

    if (isImage && file.size > MAX_IMAGE_SIZE) {
      toast.add({
        title: "File Too Large",
        description: `"${file.name}" exceeds the 5MB limit for images.`,
        color: "error",
      });
      continue;
    }

    if (isVideo && file.size > MAX_VIDEO_SIZE) {
      toast.add({
        title: "File Too Large",
        description: `"${file.name}" exceeds the 50MB limit for videos.`,
        color: "error",
      });
      continue;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      postForm.value.media.push({
        name: file.name,
        type: file.type,
        size: formatFileSize(file.size),
        url: e.target?.result as string,
        fileType: isImage ? "image" : "video",
      });
    };
    reader.readAsDataURL(file);
  }

  input.value = "";
}

function removeMedia(index: number) {
  postForm.value.media.splice(index, 1);
}

async function savePost() {
  if (!postForm.value.title.trim()) {
    toast.add({
      title: "Error",
      description: "Post title is required.",
      color: "error",
    });
    return;
  }

  isSaving.value = true;

  try {
    const data: CreateBlogPostData = {
      title: postForm.value.title,
      author: postForm.value.author,
      content: postForm.value.content,
      media: [...postForm.value.media],
      publishing: {
        status: postForm.value.publishing.status,
        publishedAt: postForm.value.publishing.publishedAt,
        scheduledAt: postForm.value.publishing.scheduledAt,
      },
      attractiveness: {
        isFeatured: postForm.value.attractiveness.isFeatured,
        isSpotlight: postForm.value.attractiveness.isSpotlight,
        priority: postForm.value.attractiveness.priority,
        highlight: postForm.value.attractiveness.highlight,
      },
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
  <UModal :open="open" title="Post Blog" @update:open="(val) => emit('update:open', val)">
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
          <UFormField label="Post Title" required>
            <UInput
              v-model="postForm.title"
              placeholder="e.g. How to Charge Your EV at Home"
              class="w-full"
            />
          </UFormField>
          <UFormField label="Author">
            <UInput v-model="postForm.author" placeholder="Admin" class="w-full" />
          </UFormField>
          <UFormField label="Content">
            <RichTextEditor
              v-model="postForm.content"
              placeholder="Write your blog post content here..."
            />
          </UFormField>
          <UFormField label="Media">
            <template #hint>
              <span>{{ t('blog.mediaHint') }}</span>
            </template>
            <input
              ref="mediaInputRef"
              type="file"
              accept="image/jpeg,image/png,image/webp,video/mp4,video/webm"
              multiple
              class="hidden"
              @change="handleMediaSelect"
            />

            <!-- Media Preview Grid -->
            <div v-if="postForm.media.length > 0" class="grid grid-cols-3 gap-3 mb-3">
              <div
                v-for="(item, idx) in postForm.media"
                :key="idx"
                class="relative group rounded-lg border border-default overflow-hidden bg-muted/30"
              >
                <img
                  v-if="item.fileType === 'image'"
                  :src="item.url"
                  :alt="item.name"
                  class="w-full h-24 object-cover"
                />
                <div
                  v-else
                  class="w-full h-24 flex flex-col items-center justify-center gap-1"
                >
                  <UIcon name="i-lucide-film" class="size-8 text-muted" />
                  <span class="text-[10px] text-muted truncate max-w-full px-2">{{
                    item.name
                  }}</span>
                </div>
                <div
                  class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                >
                  <UButton
                    icon="i-lucide-trash-2"
                    color="error"
                    variant="ghost"
                    size="xs"
                    @click="removeMedia(idx)"
                    class="text-white"
                  />
                </div>
                <div
                  class="absolute bottom-0 left-0 right-0 bg-black/60 px-2 py-0.5 flex items-center justify-between"
                >
                  <span class="text-[10px] text-white truncate">{{ item.name }}</span>
                  <span class="text-[10px] text-white/70 shrink-0 ml-1">{{
                    item.size
                  }}</span>
                </div>
              </div>
            </div>

            <!-- Upload Button -->
            <button
              @click="triggerMediaUpload"
              class="w-full border-2 border-dashed border-default rounded-lg p-4 hover:border-primary hover:bg-primary/5 transition-colors cursor-pointer flex flex-col items-center gap-2"
            >
              <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
              <span class="text-sm font-medium text-muted"
                >{{ t('blog.clickToAddMedia') }}</span
              >
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
