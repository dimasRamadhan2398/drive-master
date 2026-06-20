<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { computed, ref, onMounted } from "vue";
import type { TableColumn } from "@nuxt/ui";
import { contentService } from "~/services/contentService";
import type { FAQ, Page } from "~/services/contentService";

const { t } = useI18n()
definePageMeta({ layout: "admin" });

const toast = useToast();

// Tabs
const tabs = computed(() => [
  { label: t('admin.content').replace('Konten', 'Pages'), icon: "i-lucide-files", slot: "pages" },
  { label: "FAQs", icon: "i-lucide-help-circle", slot: "faqs" },
]);

const activeTab = ref(0);

// --- Pages Data ---
const pages = ref<Page[]>([]);
const isPagesLoading = ref(false);

const pageColumns: TableColumn<Page>[] = [
  { accessorKey: "title", header: t('admin.content.title') },
  { accessorKey: "slug", header: t('admin.content.slug') },
  {
    accessorKey: "status",
    header: t('admin.content.status'),
    cell: ({ row }) => {
      const status = row.getValue("status") as string;
      return h(resolveComponent("UBadge"), {
        label: status.toUpperCase(),
        color: status === "published" ? "success" : "neutral",
        variant: "subtle",
      });
    },
  },
  { accessorKey: "lastModified", header: t('admin.content.lastModified') },
  { id: "actions" },
];

async function fetchPages() {
  isPagesLoading.value = true;
  try {
    const data = await contentService.getPages();
    pages.value = data;
  } catch (err) {
    toast.add({
      title: t('common.error'),
      description: "Failed to load pages.",
      color: "error",
    });
  } finally {
    isPagesLoading.value = false;
  }
}

// --- FAQs Data ---
const faqs = ref<FAQ[]>([]);
const isFaqsLoading = ref(false);

const faqColumns: TableColumn<FAQ>[] = [
  { accessorKey: "question", header: t('admin.content.question') },
  { accessorKey: "category", header: t('admin.content.category') },
  {
    accessorKey: "status",
    header: t('admin.content.status'),
    cell: ({ row }) => {
      const status = row.getValue("status") as string;
      return h(resolveComponent("UBadge"), {
        label: status.toUpperCase(),
        color: status === "published" ? "success" : "neutral",
        variant: "subtle",
      });
    },
  },
  { id: "actions" },
];

async function fetchFAQs() {
  isFaqsLoading.value = true;
  try {
    const data = await contentService.getFAQs();
    faqs.value = data;
  } catch (err) {
    toast.add({
      title: t('common.error'),
      description: "Failed to load FAQs.",
      color: "error",
    });
  } finally {
    isFaqsLoading.value = false;
  }
}

// Modals
const isPageModalOpen = ref(false);
const isFAQModalOpen = ref(false);
const isEditing = ref(false);
const editingId = ref<string | number | null>(null);

const pageForm = reactive({
  title: "",
  slug: "",
  status: "draft" as "draft" | "published",
});

const faqForm = reactive({
  question: "",
  answer: "",
  category: "General",
  status: "published" as "draft" | "published",
});

function openNewPage() {
  isEditing.value = false;
  pageForm.title = "";
  pageForm.slug = "/";
  pageForm.status = "draft";
  isPageModalOpen.value = true;
}

function openEditPage(page: Page) {
  isEditing.value = true;
  editingId.value = page.id;
  pageForm.title = page.title;
  pageForm.slug = page.slug;
  pageForm.status = page.status;
  isPageModalOpen.value = true;
}

function openNewFAQ() {
  isEditing.value = false;
  faqForm.question = "";
  faqForm.answer = "";
  faqForm.category = "General";
  faqForm.status = "published";
  isFAQModalOpen.value = true;
}

function openEditFAQ(faq: FAQ) {
  isEditing.value = true;
  editingId.value = faq.id;
  faqForm.question = faq.question;
  faqForm.answer = faq.answer;
  faqForm.category = faq.category;
  faqForm.status = faq.status;
  isFAQModalOpen.value = true;
}

async function savePage() {
  try {
    if (isEditing.value && editingId.value) {
      await contentService.updatePage(editingId.value as string, {
        title: pageForm.title,
        slug: pageForm.slug,
        status: pageForm.status,
      });
      toast.add({ title: "Page Updated", color: "success" });
    } else {
      await contentService.createPage({
        title: pageForm.title,
        slug: pageForm.slug,
        status: pageForm.status,
      });
      toast.add({ title: "Page Created", color: "success" });
    }
    isPageModalOpen.value = false;
    fetchPages();
  } catch (err) {
    toast.add({
      title: t('common.error'),
      description: "Failed to save page.",
      color: "error",
    });
  }
}

async function saveFAQ() {
  try {
    if (isEditing.value && editingId.value) {
      await contentService.updateFAQ(editingId.value as number, {
        question: faqForm.question,
        answer: faqForm.answer,
        category: faqForm.category,
        status: faqForm.status,
      });
      toast.add({ title: "FAQ Updated", color: "success" });
    } else {
      await contentService.createFAQ({
        question: faqForm.question,
        answer: faqForm.answer,
        category: faqForm.category,
        status: faqForm.status,
      });
      toast.add({ title: "FAQ Created", color: "success" });
    }
    isFAQModalOpen.value = false;
    fetchFAQs();
  } catch (err) {
    toast.add({
      title: t('common.error'),
      description: "Failed to save FAQ.",
      color: "error",
    });
  }
}

async function deletePage(id: string) {
  if (confirm(t('admin.content.deleteConfirmPage'))) {
    try {
      await contentService.deletePage(id);
      toast.add({ title: t('toast.deleteSuccess'), color: "error" });
      fetchPages();
    } catch (err) {
      toast.add({ title: t('common.error'), color: "error" });
    }
  }
}

async function deleteFAQ(id: number) {
  if (confirm(t('admin.content.deleteConfirmFaq'))) {
    try {
      await contentService.deleteFAQ(id);
      toast.add({ title: t('toast.deleteSuccess'), color: "error" });
      fetchFAQs();
    } catch (err) {
      toast.add({ title: t('common.error'), color: "error" });
    }
  }
}

onMounted(() => {
  fetchPages();
  fetchFAQs();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.content')">
        <template #right>
          <UButton
            v-if="activeTab === 0"
            icon="i-lucide-plus"
            color="warning"
            :label="t('admin.content.createPage')"
            @click="openNewPage"
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
          <!-- Pages Tab -->
          <template #pages>
            <UCard class="mt-6">
              <UTable
                :columns="pageColumns"
                :data="pages"
                :loading="isPagesLoading"
              >
                <template #actions-cell="{ row }">
                  <div class="flex justify-end gap-2">
                    <UButton
                      icon="i-lucide-pencil"
                      variant="ghost"
                      color="neutral"
                      @click="openEditPage(row.original)"
                    />
                    <UButton
                      icon="i-lucide-trash"
                      variant="ghost"
                      color="error"
                      @click="deletePage(row.original.id)"
                    />
                  </div>
                </template>
              </UTable>
            </UCard>
          </template>

          <!-- FAQs Tab -->
          <template #faqs>
            <UCard class="mt-6">
              <UTable
                :columns="faqColumns"
                :data="faqs"
                :loading="isFaqsLoading"
              >
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

      <!-- Page Modal -->
      <UModal v-model:open="isPageModalOpen" :title="isEditing ? t('admin.content.editPage') : t('admin.content.createPage')">
        <template #body>
          <div class="space-y-4">
            <UFormField :label="t('admin.content.title')" required>
              <UInput v-model="pageForm.title" color="warning" class="w-full" />
            </UFormField>
            <UFormField :label="t('admin.content.slug')" required>
              <UInput v-model="pageForm.slug" color="warning" class="w-full" />
            </UFormField>
            <UFormField :label="t('admin.content.status')">
              <USelect
                v-model="pageForm.status"
                :items="['draft', 'published']"
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
              @click="isPageModalOpen = false"
            />
            <UButton :label="t('admin.content.savePage')" color="warning" @click="savePage" />
          </div>
        </template>
      </UModal>

      <!-- FAQ Modal -->
      <UModal v-model:open="isFAQModalOpen" :title="isEditing ? t('admin.content.editFaq') : t('admin.content.addFaq')">
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
                :items="['draft', 'published']"
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
            <UButton :label="t('admin.content.saveFaq')" color="warning" @click="saveFAQ" />
          </div>
        </template>
      </UModal>
    </template>
  </UDashboardPanel>
</template>
