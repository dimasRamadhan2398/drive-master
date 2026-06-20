<script setup lang="ts">
const { t } = useI18n()

useSeoMeta({
  title: t('instructors.title') + " | Drive Master Academy",
  description: t('instructors.subtitle'),
});

const instructorsStore = useInstructorsStore();

const instructors = computed(() => instructorsStore.instructors);

onMounted(() => {
  instructorsStore.fetchInstructors();
});
</script>

<template>
  <UContainer class="py-12">
    <UPageHeader
      :title="t('instructors.heroTitle')"
      :description="t('instructors.heroDesc')"
      class="mb-12"
    />

    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
      <UCard
        v-for="instructor in instructors"
        :key="instructor.userId"
        class="overflow-hidden group"
      >
        <template #header>
          <div class="relative h-64 -m-6 mb-0 overflow-hidden">
            <img
              v-if="instructor.image"
              :src="instructor.image"
              :alt="instructor.name"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
            />
            <div
              v-else
              class="w-full h-full bg-gradient-to-br from-warning/20 to-warning/40 flex items-center justify-center"
            >
              <UIcon name="i-lucide-user" class="size-16 text-warning/60" />
            </div>
            <div
              class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent"
            />
            <div class="absolute bottom-4 left-4">
              <UBadge
                :label="instructor.status"
                :color="instructor.status === 'active' ? 'success' : 'warning'"
                variant="subtle"
              />
            </div>
          </div>
        </template>

        <div class="space-y-4">
          <div>
            <h3 class="text-xl font-bold">{{ instructor.name }}</h3>
            <div class="flex items-center gap-2 text-warning">
              <UIcon name="i-lucide-star" class="size-4" />
              <span class="text-sm font-medium">{{ instructor.rating.toFixed(1) }} {{ t('instructors.rating') }}</span>
            </div>
          </div>

          <p class="text-muted text-sm leading-relaxed line-clamp-3">
            {{ instructor.bio || t('instructors.noBio') }}
          </p>

          <div class="flex items-center gap-4 pt-2 border-t border-default">
            <div class="flex items-center gap-1">
              <UIcon name="i-lucide-clock" class="size-4 text-warning" />
              <span class="text-xs font-semibold"
                >{{ instructor.yearsOfExperience }} {{ t('home.yearsExperience') }}</span
              >
            </div>
            <div class="flex items-center gap-1">
              <UIcon name="i-lucide-users" class="size-4 text-warning" />
              <span class="text-xs font-semibold">{{ instructor.totalStudents }} {{ t('instructors.students') }}</span>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="space-y-2">
            <div class="flex items-center gap-2 text-sm">
              <UIcon name="i-lucide-phone" class="size-4 text-muted" />
              <span class="text-muted">{{ instructor.phone || t('instructors.noPhone') }}</span>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <UIcon name="i-lucide-mail" class="size-4 text-muted" />
              <span class="text-muted truncate">{{ instructor.email }}</span>
            </div>
            <div v-if="instructor.certifications.length" class="flex items-center gap-2 text-sm">
              <UIcon name="i-lucide-shield-check" class="size-4 text-muted" />
              <span class="text-muted">{{ instructor.certifications[0] }}</span>
            </div>
          </div>
        </template>
      </UCard>
    </div>

    <UPageCTA
      :title="t('instructors.cta.title')"
      :description="t('instructors.cta.description')"
      :links="[
        {
          label: t('auth.register'),
          to: '/auth/register',
          color: 'warning',
          size: 'lg',
        },
      ]"
      class="mt-20 bg-warning/5"
    />
  </UContainer>
</template>
