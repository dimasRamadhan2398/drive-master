<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { ref } from "vue";
import FaqModal from "~/components/content/FaqModal.vue";
import ArticlePageModal from "~/components/content/ArticlePageModal.vue";
import PostBlogModal from "~/components/content/PostBlogModal.vue";
import { useContentStore } from "~/stores/content";
import type { Page, BlogPost, Faq, BlogPostMedia } from "~/stores/content";
import type { CreateBlogPostData } from "~/services/contentService";
import type { CreateFaqData } from "~/components/content/FaqModal.vue";
import type { CreatePageData } from "~/components/content/ArticlePageModal.vue";

definePageMeta({ layout: "admin" });

const toast = useToast();
const contentStore = useContentStore();
const activeTab = ref("pages");

// ==================== MODAL STATES ====================
const showPageModal = ref(false);
const showPostModal = ref(false);
const showFaqModal = ref(false);
const isEditing = ref(false);

// ==================== FORM DATA ====================
const pageForm = ref({ id: 0, title: "", slug: "", status: "draft" as "draft" | "published" | "archived" });
const postForm = ref({
  id: 0,
  title: "",
  author: "Admin",
  content: "",
  status: "draft" as "draft" | "published" | "archived",
  media: [] as BlogPostMedia[],
});
const faqForm = ref({ id: "", question: "", answer: "", sortOrder: 0 });

// Get data from store
const pages = computed(() => contentStore.pages);
const blogPosts = computed(() => contentStore.blogPosts);
const faqs = computed(() => contentStore.sortedFaqs);

const tabs = [
  { label: "Pages", value: "pages", icon: "i-lucide-file-text" },
  { label: "Blog Posts", value: "blog", icon: "i-lucide-newspaper" },
  { label: "FAQs", value: "faq", icon: "i-lucide-help-circle" },
];

// ==================== PAGES CRUD ====================
const currentEditingPage = ref<Page | null>(null);

function openPageBuilder(page: Page) {
  currentEditingPage.value = page;
}

function savePageBuilder(updatedPage: Page) {
  contentStore.updatePage(updatedPage.id, updatedPage);
  currentEditingPage.value = null;
}

function openNewPage() {
  isEditing.value = false;
  pageForm.value = { id: 0, title: "", slug: "", status: "draft" };
  showPageModal.value = true;
}

function openEditPage(page: Page) {
  isEditing.value = true;
  pageForm.value = {
    id: page.id,
    title: page.title,
    slug: page.slug,
    status: page.status,
  };
  showPageModal.value = true;
}

function handlePageSaved(data: CreatePageData) {
  if (isEditing.value) {
    contentStore.updatePage(pageForm.value.id, data);
    toast.add({
      title: "Page Updated",
      description: `"${data.title}" has been updated.`,
      color: "success",
    });
  } else {
    contentStore.addPage(data);
    toast.add({
      title: "Page Created",
      description: `"${data.title}" has been created.`,
      color: "success",
    });
  }
}

function previewPage(page: Page) {
  toast.add({
    title: "Preview",
    description: `Opening preview for "${page.title}" (${page.slug})...`,
    color: "info",
  });
  if (!page.slug || page.slug === "/") {
    return window.open("/", "noopener,noreferrer");
  }

  const targetUrl = page.slug.startsWith("/") ? page.slug : `/${page.slug}`;

  window.open(targetUrl, "_blank", "noopener,noreferrer");
}

function deletePage(page: Page) {
  contentStore.deletePage(page.id);
  toast.add({
    title: "Page Deleted",
    description: `"${page.title}" has been removed.`,
    color: "error",
  });
}

function togglePageStatus(page: Page) {
  const newStatus = contentStore.togglePageStatus(page.id);
  if (newStatus) {
    toast.add({
      title: "Status Updated",
      description: `"${page.title}" is now ${newStatus}.`,
      color: "success",
    });
  }
}

function getPageActions(page: Page) {
  return [
    [
      {
        label: "Edit Content",
        icon: "i-lucide-layout-template",
        onSelect: () => openPageBuilder(page),
      },
      {
        label: "Edit Settings",
        icon: "i-lucide-settings",
        onSelect: () => openEditPage(page),
      },
      {
        label: "Preview",
        icon: "i-lucide-eye",
        onSelect: () => previewPage(page),
      },
      {
        label: page.status === "published" ? "Unpublish" : "Publish",
        icon:
          page.status === "published" ? "i-lucide-eye-off" : "i-lucide-globe",
        onSelect: () => togglePageStatus(page),
      },
    ],
    [
      {
        label: "Delete",
        icon: "i-lucide-trash",
        color: "error" as const,
        onSelect: () => deletePage(page),
      },
    ],
  ];
}

// ==================== BLOG POSTS CRUD ====================
function openNewPost() {
  isEditing.value = false;
  postForm.value = {
    id: 0,
    title: "",
    author: "Admin",
    content: "",
    status: "draft",
    media: [],
  };
  showPostModal.value = true;
}

function openEditPost(post: BlogPost) {
  isEditing.value = true;
  postForm.value = {
    id: post.id,
    title: post.title,
    author: post.author,
    content: post.content || "",
    status: post.status,
    media: post.media ? [...post.media] : [],
  };
  showPostModal.value = true;
}

async function handlePostSaved(data: CreateBlogPostData) {
  const postData = {
    title: data.title,
    author: data.author,
    content: data.content,
    status: data.publishing?.status || "draft",
    media: data.media,
    publishing: data.publishing,
    attractiveness: data.attractiveness,
  };

  if (isEditing.value) {
    await contentStore.updatePost(postForm.value.id, postData);
    toast.add({
      title: "Post Updated",
      description: `"${data.title}" has been updated.`,
      color: "success",
    });
  } else {
    await contentStore.addPost(postData);
    toast.add({
      title: "Post Created",
      description: `"${data.title}" has been created.`,
      color: "success",
    });
  }
}

function previewPost(post: BlogPost) {
  toast.add({
    title: "Preview",
    description: `Opening preview for "${post.title}"...`,
    color: "info",
  });
  if (!post.date || post.date === "/") {
    return window.open("/", "noopener,noreferrer");
  }
  window.open(post.date, "_blank", "noopener,noreferrer");
}

async function deletePost(post: BlogPost) {
  await contentStore.deletePost(post.id);
  toast.add({
    title: "Post Deleted",
    description: `"${post.title}" has been removed.`,
    color: "error",
  });
}

function togglePostStatus(post: BlogPost) {
  const newStatus = contentStore.togglePostStatus(post.id);
  if (newStatus) {
    toast.add({
      title: "Status Updated",
      description: `"${post.title}" is now ${newStatus}.`,
      color: "success",
    });
  }
}

function getPostActions(post: BlogPost) {
  return [
    [
      {
        label: "Edit",
        icon: "i-lucide-pencil",
        onSelect: () => openEditPost(post),
      },
      {
        label: "Preview",
        icon: "i-lucide-eye",
        onSelect: () => previewPost(post),
      },
      {
        label: post.status === "published" ? "Unpublish" : "Publish",
        icon:
          post.status === "published" ? "i-lucide-eye-off" : "i-lucide-globe",
        onSelect: () => togglePostStatus(post),
      },
    ],
    [
      {
        label: "Delete",
        icon: "i-lucide-trash",
        color: "error" as const,
        onSelect: () => deletePost(post),
      },
    ],
  ];
}

// ==================== FAQ CRUD ====================
function openNewFaq() {
  isEditing.value = false;
  faqForm.value = { id: "", question: "", answer: "", sortOrder: 0 };
  showFaqModal.value = true;
}

function openEditFaq(faq: Faq) {
  isEditing.value = true;
  faqForm.value = {
    id: faq.id,
    question: faq.question,
    answer: faq.answer,
    sortOrder: faq.sortOrder,
  };
  showFaqModal.value = true;
}

function handleFaqSaved(data: CreateFaqData) {
  if (isEditing.value) {
    contentStore.updateFaq(faqForm.value.id, {
      question: data.question,
      answer: data.answer,
    });
    toast.add({
      title: "FAQ Updated",
      description: "FAQ has been updated.",
      color: "success",
    });
  } else {
    contentStore.createFaq({
      question: data.question,
      answer: data.answer,
    });
    toast.add({
      title: "FAQ Created",
      description: "New FAQ has been added.",
      color: "success",
    });
  }
}

function deleteFaq(faq: Faq) {
  contentStore.deleteFaq(faq.id);
  toast.add({
    title: "FAQ Deleted",
    description: "FAQ has been removed.",
    color: "error",
  });
}

// ==================== FAQ DRAG & DROP ====================
const dragIndex = ref<number | null>(null);

function onDragStart(index: number, event: DragEvent) {
  dragIndex.value = index;
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = "move";
  }
}

function onDragOver(event: DragEvent) {
  event.preventDefault();
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = "move";
  }
}

async function onDrop(targetIndex: number) {
  if (dragIndex.value === null || dragIndex.value === targetIndex) return;
  await contentStore.reorderFaqs(dragIndex.value, targetIndex);
  dragIndex.value = null;
  toast.add({
    title: "Reordered",
    description: "FAQ order has been updated.",
    color: "success",
  });
}

function onDragEnd() {
  dragIndex.value = null;
}

// ==================== HEADER BUTTON HANDLER ====================
function handleHeaderAction() {
  if (activeTab.value === "pages") openNewPage();
  else if (activeTab.value === "blog") openNewPost();
  else openNewFaq();
}

onMounted(() => {
  // contentStore.fetchPages()
  contentStore.fetchBlogPosts();
  contentStore.fetchFaqs();
});
</script>

<template>
  <UDashboardPanel v-if="currentEditingPage">
    <div class="h-full w-full overflow-y-auto p-6 bg-background">
      <AdminPageEditor
        :page="currentEditingPage"
        @close="currentEditingPage = null"
        @save="savePageBuilder"
      />
    </div>
  </UDashboardPanel>

  <UDashboardPanel v-else>
    <template #header>
      <UDashboardNavbar title="Content Management">
        <template #right>
          <UButton
            :icon="
              activeTab === 'pages'
                ? 'i-lucide-file-plus'
                : activeTab === 'blog'
                  ? 'i-lucide-pen-square'
                  : 'i-lucide-plus-circle'
            "
            :label="
              activeTab === 'pages'
                ? 'New Page'
                : activeTab === 'blog'
                  ? 'New Post'
                  : 'Add FAQ'
            "
            color="warning"
            @click="handleHeaderAction"
          />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <UTabs
            v-model="activeTab"
            :items="tabs"
            class="py-4"
            color="warning"
          />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <div class="p-6">
        <UCard v-if="activeTab === 'pages'">
          <template #header>
            <h2 class="font-semibold">Website Pages</h2>
          </template>

          <div class="space-y-3">
            <div
              v-for="page in pages"
              :key="page.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
            >
              <div class="flex items-center gap-4">
                <UIcon name="i-lucide-file-text" class="size-5 text-muted" />
                <div>
                  <p class="font-medium">{{ page.title }}</p>
                  <code class="text-md bg-muted px-2 py-0.5 rounded">{{
                    page.slug
                  }}</code>
                </div>
              </div>
              <div class="flex items-center gap-4">
                <span class="text-md text-muted">{{ page.lastUpdated }}</span>
                <UBadge
                  :label="page.status"
                  :color="page.status === 'published' ? 'success' : 'neutral'"
                  variant="subtle"
                />
                <UDropdownMenu :items="getPageActions(page)">
                  <UButton
                    icon="i-lucide-ellipsis-vertical"
                    color="neutral"
                    variant="ghost"
                    size="sm"
                  />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </UCard>
        <UCard v-if="activeTab === 'blog'">
          <template #header>
            <h2 class="font-semibold">Blog Posts</h2>
          </template>

          <div class="space-y-3">
            <div
              v-for="post in blogPosts"
              :key="post.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
            >
              <div class="flex-1">
                <p class="font-medium">{{ post.title }}</p>
                <div class="flex items-center gap-4 mt-1 text-md text-muted">
                  <span>By {{ post.author }}</span>
                  <span>{{ post.date }}</span>
                  <span class="flex items-center gap-1">
                    <UIcon name="i-lucide-eye" class="size-4" />
                    {{ post.views }}
                  </span>
                </div>
              </div>
              <div class="flex items-center gap-4">
                <UBadge
                  :label="post.status"
                  :color="post.status === 'published' ? 'success' : 'neutral'"
                  variant="subtle"
                />
                <UDropdownMenu :items="getPostActions(post)">
                  <UButton
                    icon="i-lucide-ellipsis-vertical"
                    color="neutral"
                    variant="ghost"
                    size="sm"
                  />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </UCard>

        <UCard v-if="activeTab === 'faq'">
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">Frequently Asked Questions</h2>
              <p class="text-md text-muted">Drag to reorder</p>
            </div>
          </template>

          <div class="space-y-3">
            <div
              v-for="(faq, index) in faqs"
              :key="faq.id"
              draggable="true"
              class="flex items-start gap-4 p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
              :class="{
                'opacity-50 border-dashed border-primary': dragIndex === index,
              }"
              @dragstart="onDragStart(index, $event)"
              @dragover="onDragOver"
              @drop="onDrop(index)"
              @dragend="onDragEnd"
            >
              <UIcon
                name="i-lucide-grip-vertical"
                class="size-5 text-muted cursor-grab mt-0.5 shrink-0"
              />
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <UBadge
                    :label="`#${faq.sortOrder}`"
                    variant="subtle"
                    size="md"
                  />
                  <h3 class="font-medium">{{ faq.question }}</h3>
                </div>
                <p class="text-md text-muted line-clamp-2">{{ faq.answer }}</p>
              </div>
              <div class="flex gap-1 shrink-0">
                <UButton
                  icon="i-lucide-pencil"
                  color="neutral"
                  variant="ghost"
                  size="sm"
                  @click="openEditFaq(faq)"
                />
                <UButton
                  icon="i-lucide-trash"
                  color="error"
                  variant="ghost"
                  size="sm"
                  @click="deleteFaq(faq)"
                />
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
  <ClientOnly>
    <!-- <ArticlePageModal
      v-model:open="showPageModal"
      :page="isEditing ? pageForm : null"
      @saved="handlePageSaved"
    /> -->
  </ClientOnly>
  <ClientOnly>
    <PostBlogModal
      v-model:open="showPostModal"
      :post="isEditing ? postForm : null"
      @saved="handlePostSaved"
    />
  </ClientOnly>
  <ClientOnly>
    <FaqModal
      v-model:open="showFaqModal"
      :faq="isEditing ? faqForm : null"
      @saved="handleFaqSaved"
    />
  </ClientOnly>
</template>
