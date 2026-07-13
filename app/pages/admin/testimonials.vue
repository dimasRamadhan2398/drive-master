<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { ref, computed } from "vue";
import { useTestimonialsStore, type Testimonial } from "~/stores/testimonials";
import AddTestimonialModal from "~/components/testimonials/AddTestimonialModal.vue";
import EditTestimonialModal from "~/components/testimonials/EditTestimonialModal.vue";

const { t } = useI18n()
definePageMeta({ layout: "admin" });

const toast = useToast();
const testimonialsStore = useTestimonialsStore();

const showAddModal = ref(false);
const showEditModal = ref(false);
const selectedTestimonial = ref<Testimonial | null>(null);
const filterStatus = ref<"all" | "published" | "pending" | "archived">("all");

const testimonials = computed(() => testimonialsStore.testimonials);

const filteredTestimonials = computed(() => {
  if (filterStatus.value === "all") return testimonials.value;
  return testimonials.value.filter((t) => t.status === filterStatus.value);
});

// Stats
const totalTestimonials = computed(() => testimonialsStore.totalTestimonials);
const totalPublished = computed(() => testimonialsStore.totalPublished);
const totalPending = computed(() => testimonialsStore.totalPending);
const averageRating = computed(() => testimonialsStore.averageRating);

function openEditModal(testimonial: Testimonial) {
  selectedTestimonial.value = { ...testimonial };
  showEditModal.value = true;
}

function toggleFeatured(testimonial: Testimonial) {
  testimonialsStore.toggleFeatured(testimonial.id!);
  toast.add({
    title: testimonial.isFeatured
      ? "Removed from Featured"
      : "Added to Featured",
    icon: "i-lucide-star",
    color: testimonial.isFeatured ? "warning" : "success",
  });
}

function changeStatus(
  testimonialId: string,
  status: "draft" | "pending" | "published" | "archived",
) {
  testimonialsStore.changeStatus(testimonialId, status);
  toast.add({
    title: "Status Updated",
    description: `Testimonial status changed to ${status}`,
    color: "success",
  });
}

function handleAddTestimonial(testimonial: Testimonial) {
  // The testimonial has already been created and added to the store by AddTestimonialModal.vue
}

function handleEditTestimonial(testimonial: Testimonial) {
  testimonialsStore.updateTestimonial(testimonial.id!, testimonial);

  toast.add({
    title: "Testimonial Updated",
    description: `"${testimonial.userName}" has been saved.`,
    color: "success",
  });

  showEditModal.value = false;
  selectedTestimonial.value = null;
}

function deleteTestimonial(testimonialId: string) {
  if (
    confirm(
      "Are you sure you want to delete this testimonial? This action cannot be undone.",
    )
  ) {
    testimonialsStore.deleteTestimonial(testimonialId);
    toast.add({
      title: "Testimonial Deleted",
      description: `Testimonial has been removed.`,
      color: "error",
      icon: "i-lucide-trash",
    });
  }
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

function getStatusColor(
  status: string,
): "success" | "warning" | "neutral" | "info" {
  switch (status) {
    case "published":
      return "success";
    case "pending":
      return "warning";
    case "archived":
      return "neutral";
    default:
      return "info";
  }
}

function getRatingStars(rating: number): string[] {
  return Array(5)
    .fill("")
    .map((_, i) => (i < rating ? "full" : "empty"));
}

function onShowModal() {
  console.log("Opening Add Testimonial Modal");
  showAddModal.value = true;
}

onMounted(() => {
  testimonialsStore.fetchTestimonials();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.testimonials')">
        <template #right>
          <UButton
            icon="i-lucide-plus"
            color="warning"
            :label="t('admin.addNew')"
            @click="onShowModal"
          />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Stats -->
        <div class="grid md:grid-cols-4 gap-4">
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-info/10">
                <UIcon
                  name="i-lucide-message-square"
                  class="size-6 text-info"
                />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ totalTestimonials }}</p>
                <p class="text-md text-muted">Total {{ t('admin.testimonials') }}</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-green-500/10">
                <UIcon
                  name="i-lucide-check-circle"
                  class="size-6 text-green-500"
                />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ totalPublished }}</p>
                <p class="text-md text-muted">{{ t('admin.published') }}</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-amber-500/10">
                <UIcon name="i-lucide-clock" class="size-6 text-amber-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ totalPending }}</p>
                <p class="text-md text-muted">{{ t('admin.pendingReview') }}</p>
              </div>
            </div>
          </UCard>
          <UCard>
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-yellow-500/10">
                <UIcon name="i-lucide-star" class="size-6 text-yellow-500" />
              </div>
              <div>
                <p class="text-2xl font-bold">{{ averageRating }}</p>
                <p class="text-md text-muted">{{ t('admin.averageRating') }}</p>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Filter -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">All {{ t('admin.testimonials') }}</h2>
              <div class="flex items-center gap-2">
                <UButton
                  v-for="status in [
                    'all',
                    'published',
                    'pending',
                    'archived',
                  ] as const"
                  :key="status"
                  :label="
                    status === 'all'
                      ? 'All'
                      : status.charAt(0).toUpperCase() + status.slice(1)
                  "
                  :color="filterStatus === status ? 'warning' : 'neutral'"
                  :variant="filterStatus === status ? 'solid' : 'outline'"
                  size="sm"
                  @click="filterStatus = status"
                />
              </div>
            </div>
          </template>

          <div class="space-y-4">
            <div
              v-for="testimonial in filteredTestimonials"
              :key="testimonial.id"
              class="p-4 rounded-lg border border-default"
              :class="{ 'ring-2 ring-warning': testimonial.isFeatured }"
            >
              <div class="flex items-start gap-4">
                <!-- Avatar -->
                <div class="shrink-0">
                  <img
                    v-if="testimonial.userImage"
                    :src="testimonial.userImage"
                    :alt="testimonial.userName"
                    class="w-12 h-12 rounded-full object-cover"
                  />
                  <div
                    v-else
                    class="w-12 h-12 rounded-full bg-warning/20 flex items-center justify-center"
                  >
                    <span class="text-lg font-bold text-warning">{{
                      testimonial.userName.charAt(0)
                    }}</span>
                  </div>
                </div>

                <!-- Content -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between gap-4">
                    <div>
                      <div class="flex items-center gap-2">
                        <h3 class="font-semibold">
                          {{ testimonial.userName }}
                        </h3>
                        <UBadge
                          :label="testimonial.userRole"
                          variant="subtle"
                          color="neutral"
                          size="sm"
                        />
                        <UBadge
                          :label="testimonial.status"
                          :color="getStatusColor(testimonial.status)"
                          size="sm"
                        />
                        <UIcon
                          v-if="testimonial.isFeatured"
                          name="i-lucide-star"
                          class="size-4 text-warning"
                        />
                      </div>
                      <p class="text-md text-muted mt-1">
                        {{ testimonial.content }}
                      </p>

                      <!-- Rating -->
                      <div class="flex items-center gap-1 mt-2">
                        <template
                          v-for="(_, i) in getRatingStars(testimonial.rating)"
                          :key="i"
                        >
                          <UIcon
                            :name="
                              i < testimonial.rating
                                ? 'i-lucide-star'
                                : 'i-lucide-star'
                            "
                            :class="
                              i < testimonial.rating
                                ? 'text-yellow-500'
                                : 'text-muted'
                            "
                            class="size-4"
                          />
                        </template>
                        <span class="text-md text-muted ml-1"
                          >({{ testimonial.rating }}/5)</span
                        >
                      </div>

                      <!-- Tags -->
                      <div
                        v-if="testimonial.tags"
                        class="flex items-center gap-2 mt-2"
                      >
                        <UBadge
                          v-for="tag in testimonial.tags.split(',')"
                          :key="tag"
                          :label="tag.trim()"
                          variant="outline"
                          color="neutral"
                          size="sm"
                        />
                      </div>

                      <p class="text-xs text-muted mt-2">
                        Added {{ formatDate(testimonial.addedAt) }}
                      </p>
                    </div>

                    <!-- Actions -->
                    <div class="flex items-center gap-2 shrink-0">
                      <USwitch
                        :model-value="testimonial.isFeatured"
                        @update:model-value="toggleFeatured(testimonial)"
                      />
                      <UDropdownMenu
                        :items="[
                          [
                            {
                              label: t('common.edit'),
                              icon: 'i-lucide-pencil',
                              onSelect: () => openEditModal(testimonial),
                            },
                            {
                              label: testimonial.isFeatured
                                ? 'Remove from Featured'
                                : 'Mark as Featured',
                              icon: 'i-lucide-star',
                              onSelect: () => toggleFeatured(testimonial),
                            },
                          ],
                          [
                            {
                              label: testimonial.status === 'published' ? 'Pending' : 'Publish',
                              icon: 'i-lucide-check-circle',
                              onSelect: () =>
                                changeStatus(testimonial.id!, testimonial.status === 'published' ? 'pending' : 'published'),
                            },
                            {
                              label: 'Archive',
                              icon: 'i-lucide-archive',
                              onSelect: () =>
                                changeStatus(testimonial.id!, 'archived'),
                            },
                          ],
                          [
                            {
                              label: t('common.delete'),
                              icon: 'i-lucide-trash',
                              color: 'error',
                              onSelect: () =>
                                deleteTestimonial(testimonial.id!),
                            },
                          ],
                        ]"
                      >
                        <UButton
                          icon="i-lucide-ellipsis-vertical"
                          color="neutral"
                          variant="ghost"
                        />
                      </UDropdownMenu>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div
              v-if="filteredTestimonials.length === 0"
              class="text-center py-8"
            >
              <UIcon
                name="i-lucide-message-square"
                class="size-12 text-muted mx-auto mb-2"
              />
              <p class="text-muted">{{ t('admin.noTestimonials') }}</p>
            </div>
          </div>
        </UCard>

        <!-- Modals -->
        <AddTestimonialModal
          v-model:open="showAddModal"
          @saved="handleAddTestimonial"
        />

        <EditTestimonialModal
          v-model:open="showEditModal"
          :testimonial="selectedTestimonial"
          @saved="handleEditTestimonial"
        />
      </div>
    </template>
  </UDashboardPanel>
</template>
