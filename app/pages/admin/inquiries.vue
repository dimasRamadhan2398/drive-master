<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useApiClients } from "~/composables/useApiClients";
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";

const { t } = useI18n();
definePageMeta({ layout: "admin" });

const toast = useToast();
const { core, extractData } = useApiClients();

interface Inquiry {
  id: string;
  name: string;
  email: string;
  subject: string;
  message: string;
  createdAt: string;
  updatedAt?: string;
}

const inquiries = ref<Inquiry[]>([]);
const isLoading = ref(false);
const selectedInquiry = ref<Inquiry | null>(null);
const isDetailOpen = ref(false);

const columns = [
  { accessorKey: "name", key: "name", header: "Name", label: "Name" },
  { accessorKey: "email", key: "email", header: "Email", label: "Email" },
  { accessorKey: "subject", key: "subject", header: "Subject", label: "Subject" },
  { accessorKey: "createdAt", key: "createdAt", header: "Date", label: "Date" },
  { accessorKey: "actions", key: "actions", header: "Actions", label: "Actions" }
];

async function fetchInquiries() {
  isLoading.value = true;
  try {
    const response = await core("/contact");
    const data = extractData(response);
    inquiries.value = Array.isArray(data) ? data : [];
    return;
  } catch (error) {
    console.warn("Primary API client failed for /contact, attempting local API gateway fallback...", error);
    try {
      const localResponse = await $fetch("http://localhost:8088/api/v1/contact");
      const data = extractData(localResponse);
      inquiries.value = Array.isArray(data) ? data : [];
      return;
    } catch (fallbackErr) {
      try {
        const directResponse = await $fetch("http://localhost:8012/api/v1/contact");
        const data = extractData(directResponse);
        inquiries.value = Array.isArray(data) ? data : [];
        return;
      } catch (directErr) {
        console.error("All contact fetch attempts failed:", directErr);
      }
    }
    toast.add({
      title: "Error",
      description: "Failed to fetch contact inquiries",
      color: "error"
    });
  } finally {
    isLoading.value = false;
  }
}

function openDetail(inquiry: Inquiry) {
  selectedInquiry.value = inquiry;
  isDetailOpen.value = true;
}

function formatDate(dateStr: string): string {
  if (!dateStr) return "-";
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}

function getRowItem(row: any): Inquiry {
  return row?.original ? row.original : row;
}

onMounted(() => {
  fetchInquiries();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.inquiries')">
        <template #right>
          <UButton
            icon="i-lucide-refresh-cw"
            color="neutral"
            variant="ghost"
            :loading="isLoading"
            @click="fetchInquiries"
          />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6">
        <UCard>
          <UTable
            :data="inquiries"
            :rows="inquiries"
            :columns="columns"
            :loading="isLoading"
          >
            <template #name-cell="{ row }">
              <span class="font-medium">{{ getRowItem(row).name }}</span>
            </template>
            <template #name-data="{ row }">
              <span class="font-medium">{{ getRowItem(row).name }}</span>
            </template>

            <template #email-cell="{ row }">
              <a :href="'mailto:' + getRowItem(row).email" class="text-primary hover:underline">{{ getRowItem(row).email }}</a>
            </template>
            <template #email-data="{ row }">
              <a :href="'mailto:' + getRowItem(row).email" class="text-primary hover:underline">{{ getRowItem(row).email }}</a>
            </template>

            <template #subject-cell="{ row }">
              <div class="max-w-md cursor-pointer" @click="openDetail(getRowItem(row))">
                <span class="font-medium hover:text-primary transition-colors">{{ getRowItem(row).subject }}</span>
                <p class="text-sm text-muted mt-1 truncate">{{ getRowItem(row).message }}</p>
              </div>
            </template>
            <template #subject-data="{ row }">
              <div class="max-w-md cursor-pointer" @click="openDetail(getRowItem(row))">
                <span class="font-medium hover:text-primary transition-colors">{{ getRowItem(row).subject }}</span>
                <p class="text-sm text-muted mt-1 truncate">{{ getRowItem(row).message }}</p>
              </div>
            </template>

            <template #createdAt-cell="{ row }">
              <span class="text-sm text-muted">{{ formatDate(getRowItem(row).createdAt) }}</span>
            </template>
            <template #createdAt-data="{ row }">
              <span class="text-sm text-muted">{{ formatDate(getRowItem(row).createdAt) }}</span>
            </template>

            <template #actions-cell="{ row }">
              <UButton
                icon="i-lucide-eye"
                color="neutral"
                variant="ghost"
                size="sm"
                @click="openDetail(getRowItem(row))"
              />
            </template>
            <template #actions-data="{ row }">
              <UButton
                icon="i-lucide-eye"
                color="neutral"
                variant="ghost"
                size="sm"
                @click="openDetail(getRowItem(row))"
              />
            </template>

            <template #empty-state>
              <div class="text-center py-8">
                <UIcon name="i-lucide-inbox" class="size-12 text-muted mx-auto mb-2" />
                <p class="text-muted">No contact messages found</p>
              </div>
            </template>
          </UTable>
        </UCard>
      </div>

      <!-- Detail Modal -->
      <UModal v-model:open="isDetailOpen" :title="selectedInquiry?.subject || 'Inquiry Details'">
        <template #content>
          <div v-if="selectedInquiry" class="p-6 space-y-4">
            <div class="flex items-center justify-between pb-4 border-b border-default">
              <div>
                <h3 class="font-semibold text-lg">{{ selectedInquiry.name }}</h3>
                <a :href="'mailto:' + selectedInquiry.email" class="text-sm text-primary">{{ selectedInquiry.email }}</a>
              </div>
              <span class="text-xs text-muted">{{ formatDate(selectedInquiry.createdAt) }}</span>
            </div>
            <div>
              <p class="text-xs text-muted uppercase tracking-wider mb-1 font-semibold">Subject</p>
              <p class="font-medium text-default">{{ selectedInquiry.subject }}</p>
            </div>
            <div>
              <p class="text-xs text-muted uppercase tracking-wider mb-1 font-semibold">Message</p>
              <div class="p-4 rounded-lg bg-muted/10 text-default whitespace-pre-wrap leading-relaxed text-sm">
                {{ selectedInquiry.message }}
              </div>
            </div>
            <div class="flex justify-end gap-2 pt-4 border-t border-default">
              <UButton
                :href="'mailto:' + selectedInquiry.email + '?subject=Re: ' + encodeURIComponent(selectedInquiry.subject)"
                icon="i-lucide-mail"
                color="primary"
                label="Reply via Email"
              />
              <UButton
                color="neutral"
                variant="outline"
                label="Close"
                @click="isDetailOpen = false"
              />
            </div>
          </div>
        </template>
      </UModal>
    </template>
  </UDashboardPanel>
</template>

