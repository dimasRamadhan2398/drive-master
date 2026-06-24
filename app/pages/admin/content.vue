<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { computed, ref, onMounted, reactive, h, resolveComponent } from "vue";
import type { TableColumn } from "@nuxt/ui";
import { useContentStore } from "~/stores/content";
import type { Faq, BlogPost } from "~/stores/content";

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

interface BlogFormData {
  title: string;
  slug: string;
  author: string;
  authorId: string;
  leadParagraph: string;
  content: string;
  isFeatured: boolean;
  isSpotlight: boolean;
  priority: number;
  highlight: boolean;
  status: "draft" | "published" | "archived";
}

const blogPosts = computed(() => contentStore.blogPosts);

const blogColumns: TableColumn<BlogPost>[] = [
  { accessorKey: "title", header: t("admin.content.title") },
  { accessorKey: "author", header: "Author" },
  {
    accessorKey: "status",
    header: t("admin.content.status"),
    cell: ({ row }) => {
      const status = row.getValue("status") as string;
      return h(resolveComponent("UBadge"), {
        label: status.toUpperCase(),
        color:
          status === "published" ? "success" : status === "draft" ? "neutral" : "error",
        variant: "subtle",
      });
    },
  },
  { accessorKey: "views", header: "Views" },
  { id: "actions" },
];

const isBlogModalOpen = ref(false);
const isEditing = ref(false);
const editingId = ref<string | number | null>(null);

const blogForm = reactive<BlogFormData>({
  title: "",
  slug: "",
  author: "",
  authorId: "00000000-0000-0000-0000-000000000000",
  leadParagraph: "",
  content: "",
  isFeatured: false,
  isSpotlight: false,
  priority: 0,
  highlight: false,
  status: "draft",
});

// Generate slug from title
const generateSlug = (title: string): string => {
  return title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
};

function openNewBlog() {
  isEditing.value = false;
  editingId.value = null;
  blogForm.title = "";
  blogForm.slug = "";
  blogForm.author = "";
  blogForm.authorId = "00000000-0000-0000-0000-000000000000";
  blogForm.leadParagraph = "";
  blogForm.content = "";
  blogForm.isFeatured = false;
  blogForm.isSpotlight = false;
  blogForm.priority = 0;
  blogForm.highlight = false;
  blogForm.status = "draft";
  isBlogModalOpen.value = true;
}

function openEditBlog(post: BlogPost) {
  isEditing.value = true;
  editingId.value = post.id;
  blogForm.title = post.title;
  blogForm.slug = post.slug || "";
  blogForm.author = post.author;
  blogForm.authorId = "00000000-0000-0000-0000-000000000000";
  blogForm.leadParagraph = "";
  blogForm.content = post.content;
  blogForm.isFeatured = post.attractiveness?.isFeatured || false;
  blogForm.isSpotlight = post.attractiveness?.isSpotlight || false;
  blogForm.priority = post.attractiveness?.priority || 0;
  blogForm.highlight = post.attractiveness?.highlight || false;
  blogForm.status = post.status;
  isBlogModalOpen.value = true;
}

function updateSlugFromTitle() {
  if (!isEditing.value && blogForm.title) {
    blogForm.slug = generateSlug(blogForm.title);
  }
}

async function saveBlog() {
  try {
    const blogData = {
      title: blogForm.title,
      slug: blogForm.slug || undefined,
      author: blogForm.author,
      content: blogForm.content,
      publishing: {
        status: blogForm.status,
      },
      attractiveness: {
        isFeatured: blogForm.isFeatured,
        isSpotlight: blogForm.isSpotlight,
        priority: blogForm.priority,
        highlight: blogForm.highlight,
      },
    };

    if (isEditing.value && editingId.value) {
      await contentStore.updatePost(editingId.value, blogData);
      toast.add({ title: "Blog Post Updated", color: "success" });
    } else {
      await contentStore.addPost({
        title: blogForm.title,
        slug: blogForm.slug || undefined,
        author: blogForm.author,
        content: blogForm.content,
        status: blogForm.status,
        authorId: blogForm.authorId,
      });
      toast.add({ title: "Blog Post Created", color: "success" });
    }
    isBlogModalOpen.value = false;
    contentStore.fetchBlogPosts();
  } catch (err) {
    toast.add({
      title: t("common.error"),
      description: "Failed to save blog post.",
      color: "error",
    });
  }
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
  isEditing.value = false;
  editingId.value = null;
  faqForm.question = "";
  faqForm.answer = "";
  faqForm.category = "General";
  faqForm.status = "published";
  isFAQModalOpen.value = true;
}

function openEditFAQ(faq: Faq) {
  isEditing.value = true;
  editingId.value = faq.id;
  faqForm.question = faq.question;
  faqForm.answer = faq.answer;
  faqForm.category = "";
  faqForm.status = faq.isActive ? "published" : "draft";
  isFAQModalOpen.value = true;
}

async function saveFAQ() {
  try {
    if (isEditing.value && editingId.value) {
      await contentStore.updateFaq(editingId.value as string, {
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

      <!-- Blog Modal -->
      <UModal
        v-model:open="isBlogModalOpen"
        :title="isEditing ? t('admin.content.editBlog') : t('admin.content.createBlog')"
      >
        <template #body>
          <div class="space-y-4">
            <!-- Title -->
            <UFormField :label="t('admin.content.title')" required>
              <UInput
                v-model="blogForm.title"
                color="warning"
                class="w-full"
                @input="updateSlugFromTitle"
              />
            </UFormField>

            <!-- Slug -->
            <UFormField :label="t('admin.content.slug')">
              <UInput v-model="blogForm.slug" color="warning" class="w-full" />
            </UFormField>

            <!-- Author -->
            <UFormField :label="t('admin.content.author')">
              <UInput v-model="blogForm.author" color="warning" class="w-full" />
            </UFormField>

            <!-- Lead Paragraph -->
            <UFormField :label="t('admin.content.leadParagraph')">
              <UTextarea
                v-model="blogForm.leadParagraph"
                color="warning"
                class="w-full"
                :rows="2"
              />
            </UFormField>

            <!-- Content -->
            <UFormField :label="t('admin.content.content')">
              <UTextarea
                v-model="blogForm.content"
                color="warning"
                class="w-full"
                :rows="8"
              />
            </UFormField>

            <!-- Status -->
            <UFormField :label="t('admin.content.status')">
              <USelect
                v-model="blogForm.status"
                :items="[
                  { label: 'Draft', value: 'draft' },
                  { label: 'Published', value: 'published' },
                  { label: 'Archived', value: 'archived' },
                ]"
                color="warning"
                class="w-full"
              />
            </UFormField>

            <!-- Attractiveness Section -->
            <div class="border-t pt-4 mt-4">
              <p class="text-sm font-medium mb-3">
                {{ t("admin.content.attractiveness") }}
              </p>

              <div class="grid grid-cols-2 gap-4">
                <!-- Priority -->
                <UFormField :label="t('admin.content.priority')">
                  <UInput
                    v-model.number="blogForm.priority"
                    type="number"
                    color="warning"
                    class="w-full"
                    min="0"
                  />
                </UFormField>

                <!-- Toggles -->
                <div class="flex flex-col gap-3">
                  <div class="flex items-center gap-2">
                    <USwitch v-model="blogForm.isFeatured" color="warning" />
                    <span class="text-sm">{{ t("admin.content.isFeatured") }}</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <USwitch v-model="blogForm.highlight" color="warning" />
                    <span class="text-sm">{{ t("admin.content.highlight") }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
        <template #footer>
          <div class="flex justify-end gap-3">
            <UButton
              :label="t('common.cancel')"
              variant="ghost"
              color="neutral"
              @click="isBlogModalOpen = false"
            />
            <UButton :label="t('common.save')" color="warning" @click="saveBlog" />
          </div>
        </template>
      </UModal>

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
</template>
