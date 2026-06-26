<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/index.js";
import { computed, ref, onMounted, reactive, h, resolveComponent } from "vue";
import type { TableColumn } from "@nuxt/ui";
import { useContentStore } from "~/stores/content";
import type { Faq, BlogPost } from "~/stores/content";
import type {
  CreateBlogArticleData,
  UpdateBlogPostData,
} from "~/services/contentService";
import PostBlogModal from "~/components/content/PostBlogModal.vue";
import type { PostFormData } from "~/types/content-form";

const { t } = useI18n();
definePageMeta({ layout: "admin" });

const contentStore = useContentStore();
const toast = useToast();

// Tabs - Only Blog and FAQ (static array for proper UTabs reactivity)
const tabs = [
  { value: "blog", label: "Blog", icon: "i-lucide-newspaper", slot: "blog" },
  { value: "faqs", label: "FAQs", icon: "i-lucide-help-circle", slot: "faqs" },
];

const activeTab = ref("blog");

// ==================== BLOG SECTION ====================

const blogPosts = computed(() => {
  return contentStore.blogPosts;
});

const blogColumns: TableColumn<BlogPost>[] = [
  { accessorKey: "title", header: t("admin.content.title") },
  { accessorKey: "author", header: "Author" },
  {
    accessorKey: "publishing.status",
    header: t("admin.content.status"),
    cell: ({ row }) => {
      const status = row.getValue("publishing.status") as string;
      return h(resolveComponent("UBadge"), {
        label: status?.toUpperCase(),
        color:
          status === "published" ? "success" : status === "draft" ? "neutral" : "error",
        variant: "subtle",
      });
    },
  },
  { accessorKey: "views", header: "Views" },
  {
    accessorKey: "attractiveness.isFeatured",
    header: "Featured",
    cell: ({ row }) => {
      const isFeatured = row.getValue("attractiveness.isFeatured") as boolean;
      return h(resolveComponent("UBadge"), {
        label: isFeatured ? "FEATURED" : "—",
        color: isFeatured ? "warning" : "neutral",
        variant: "subtle",
      });
    },
  },
  { id: "actions" },
];

// Blog Modal State
const isBlogModalOpen = ref(false);
const editingPost = ref<PostFormData | null>(null);

const isEditing = computed(() => !!editingPost.value);

function openNewBlog() {
  editingPost.value = null;
  isBlogModalOpen.value = true;
}

function openEditBlog(post: BlogPost) {
  editingPost.value = {
    id: post.id,
    title: post.title,
    slug: post.slug || "",
    author: post.author,
    leadParagraph: "",
    content: post.content,
    media: post.media || [],
    publishing: {
      status: post.publishing?.status || post.status,
      publishedAt: post.publishing?.publishedAt,
      scheduledAt: post.publishing?.scheduledAt,
    },
    attractiveness: {
      isFeatured: post.attractiveness?.isFeatured || false,
      isSpotlight: post.attractiveness?.isSpotlight || false,
      priority: post.attractiveness?.priority || 0,
      highlight: post.attractiveness?.highlight || false,
    },
  };
  isBlogModalOpen.value = true;
}

async function handleBlogSaved(data: CreateBlogArticleData) {
  try {
    if (editingPost.value) {
      // Update existing post
      const updateData: UpdateBlogPostData = {
        title: data.title,
        slug: data.slug,
        author: data.author,
        content: data.content,
        featuredImage: data.featuredImage,
        publishing: data.publishing,
        attractiveness: data.attractiveness,
      };
      await contentStore.updatePost(editingPost.value.id, updateData);
      toast.add({ title: "Blog Post Updated", color: "success" });
    } else {
      // Create new post - pass full data including authorId
      const result = await contentStore.addPost({
        title: data.title,
        slug: data.slug,
        author: data.author,
        authorId: data.authorId,
        leadParagraph: data.leadParagraph,
        content: data.content,
        featuredImage: data.featuredImage,
        publishing: data.publishing,
        attractiveness: data.attractiveness,
      });

      // Check if the result indicates an error (store returns null on failure)
      if (!result) {
        // Error was already handled by the store, but we can add UI feedback here
        return;
      }

      toast.add({ title: "Blog Post Created", color: "success" });
    }
    // Refresh the list
    contentStore.fetchBlogPosts();
  } catch (err: any) {
    // Handle API errors
    const errorMessage = err?.message || err?.data?.message || "";

    if (
      errorMessage.toLowerCase().includes("duplicate") ||
      errorMessage.toLowerCase().includes("already exists") ||
      errorMessage.toLowerCase().includes("unique constraint") ||
      errorMessage.toLowerCase().includes("title")
    ) {
      toast.add({
        title: "Validation Error",
        description:
          "A post with this title already exists. Please use a different title.",
        color: "error",
      });
    } else {
      toast.add({
        title: t("common.error"),
        description: "Failed to save blog post.",
        color: "error",
      });
    }
  }
}

function handleBlogError(message: string) {
  toast.add({
    title: "Validation Error",
    description: message,
    color: "error",
  });
}

async function deleteBlog(id: string | number) {
  if (confirm("Are you sure you want to delete this blog post?")) {
    try {
      await contentStore.deletePost(id);
      toast.add({ title: t("toast.deleteSuccess"), color: "error" });
    } catch (err) {
      toast.add({ title: t("common.error"), color: "error" });
    }
  }
}

// ==================== FAQ SECTION ====================

const faqs = computed(() => contentStore.faqs);
const isFaqsLoading = ref(false);

// FAQ Modal State
const isFaqEditing = ref(false);
const editingFaqId = ref<string | null>(null);

const faqColumns: TableColumn<Faq>[] = [
  { accessorKey: "question", header: t("admin.content.question") },
  {
    accessorKey: "isActive",
    header: t("admin.content.status"),
    cell: ({ row }) => {
      const isActive = row.getValue("isActive") as boolean;
      return h(resolveComponent("UBadge"), {
        label: isActive ? "ACTIVE" : "INACTIVE",
        color: isActive ? "success" : "neutral",
        variant: "subtle",
      });
    },
  },
  { id: "actions" },
];

async function fetchFAQs() {
  isFaqsLoading.value = true;
  try {
    await contentStore.fetchFaqs();
  } catch (err) {
    toast.add({
      title: t("common.error"),
      description: "Failed to load FAQs.",
      color: "error",
    });
  } finally {
    isFaqsLoading.value = false;
  }
}

// FAQ Modal
const isFAQModalOpen = ref(false);

const faqForm = reactive({
  question: "",
  answer: "",
  category: "General",
  status: "published" as "draft" | "published",
});

function openNewFAQ() {
  isFaqEditing.value = false;
  editingFaqId.value = null;
  faqForm.question = "";
  faqForm.answer = "";
  faqForm.category = "General";
  faqForm.status = "published";
  isFAQModalOpen.value = true;
}

function openEditFAQ(faq: Faq) {
  isFaqEditing.value = true;
  editingFaqId.value = faq.id;
  faqForm.question = faq.question;
  faqForm.answer = faq.answer;
  faqForm.category = "";
  faqForm.status = faq.isActive ? "published" : "draft";
  isFAQModalOpen.value = true;
}

async function saveFAQ() {
  try {
    if (isFaqEditing.value && editingFaqId.value) {
      await contentStore.updateFaq(editingFaqId.value, {
        question: faqForm.question,
        answer: faqForm.answer,
        isActive: faqForm.status === "published",
      });
      toast.add({ title: "FAQ Updated", color: "success" });
    } else {
      await contentStore.createFaq({
        question: faqForm.question,
        answer: faqForm.answer,
        isActive: faqForm.status === "published",
      });
      toast.add({ title: "FAQ Created", color: "success" });
    }
    isFAQModalOpen.value = false;
    fetchFAQs();
  } catch (err) {
    toast.add({
      title: t("common.error"),
      description: "Failed to save FAQ.",
      color: "error",
    });
  }
}

async function deleteFAQ(id: string) {
  if (confirm(t("admin.content.deleteConfirmFaq"))) {
    try {
      await contentStore.deleteFaq(id);
      toast.add({ title: t("toast.deleteSuccess"), color: "error" });
    } catch (err) {
      toast.add({ title: t("common.error"), color: "error" });
    }
  }
}

onMounted(() => {
  contentStore.fetchBlogPosts();
  contentStore.fetchFaqs();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.content.label')">
        <template #right>
          <UButton
            v-if="activeTab === 'blog'"
            icon="i-lucide-plus"
            color="warning"
            :label="t('admin.content.createBlog')"
            @click="openNewBlog"
          />
          <UButton
            v-else
            icon="i-lucide-plus"
            color="warning"
            :label="t('admin.content.addFaq')"
            @click="openNewFAQ"
          />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6">
        <UTabs v-model="activeTab" :items="tabs" class="w-full">
          <!-- Blog Tab -->
          <template #blog>
            <UCard class="mt-6">
              <UTable
                :columns="blogColumns"
                :data="blogPosts"
                :loading="contentStore.isLoading"
              >
                <template #actions-cell="{ row }">
                  <div class="flex justify-end gap-2">
                    <UButton
                      icon="i-lucide-pencil"
                      variant="ghost"
                      color="neutral"
                      @click="openEditBlog(row.original)"
                    />
                    <UButton
                      icon="i-lucide-trash"
                      variant="ghost"
                      color="error"
                      @click="deleteBlog(row.original.id)"
                    />
                  </div>
                </template>
              </UTable>
            </UCard>
          </template>

          <!-- FAQs Tab -->
          <template #faqs>
            <UCard class="mt-6">
              <UTable :columns="faqColumns" :data="faqs" :loading="isFaqsLoading">
                <template #actions-cell="{ row }">
                  <div class="flex justify-end gap-2">
                    <UButton
                      icon="i-lucide-pencil"
                      variant="ghost"
                      color="neutral"
                      @click="openEditFAQ(row.original)"
                    />
                    <UButton
                      icon="i-lucide-trash"
                      variant="ghost"
                      color="error"
                      @click="deleteFAQ(row.original.id)"
                    />
                  </div>
                </template>
              </UTable>
            </UCard>
          </template>
        </UTabs>
      </div>
      <!-- FAQ Modal -->
      <UModal
        v-model:open="isFAQModalOpen"
        :title="isEditing ? t('admin.content.editFaq') : t('admin.content.addFaq')"
      >
        <template #body>
          <div class="space-y-4">
            <UFormField :label="t('admin.content.question')" required>
              <UInput v-model="faqForm.question" color="warning" class="w-full" />
            </UFormField>
            <UFormField :label="t('admin.content.answer')" required>
              <UTextarea
                v-model="faqForm.answer"
                color="warning"
                class="w-full"
                :rows="4"
              />
            </UFormField>
            <UFormField :label="t('admin.content.category')">
              <UInput v-model="faqForm.category" color="warning" class="w-full" />
            </UFormField>
            <UFormField :label="t('admin.content.status')">
              <USelect
                v-model="faqForm.status"
                :items="[
                  { label: 'Draft', value: 'draft' },
                  { label: 'Published', value: 'published' },
                ]"
                color="warning"
                class="w-full"
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
              @click="isFAQModalOpen = false"
            />
            <UButton
              :label="t('admin.content.saveFaq')"
              color="warning"
              @click="saveFAQ"
            />
          </div>
        </template>
      </UModal>
    </template>
  </UDashboardPanel>

  <!-- Post Blog Modal -->
  <PostBlogModal
    :open="isBlogModalOpen"
    :post="editingPost"
    @update:open="(val) => (isBlogModalOpen = val)"
    @saved="handleBlogSaved"
    @error="handleBlogError"
  />
</template>
