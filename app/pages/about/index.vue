<script setup lang="ts">
const { t } = useI18n()

// Page metadata
definePageMeta({
  layout: 'default'
})

useSeoMeta({
  title: t('about.title') + ' | Drive Master Indonesia - The Future of Driving',
  description: t('footer.tagline')
})

const { waLink, fetchGeneralSettings } = useSettings()
const contentStore = useContentStore()
const { pages } = useContent()

await useAsyncData('about-page-sections', () => contentStore.fetchPages())

const aboutPage = computed(() =>
  pages.value.find((p) => p.slug === '/about' && p.status === 'published')
)

onMounted(() => {
  fetchGeneralSettings()
})

const safetyFeatures = [
  {
    title: t('about.safety.certifiedInstructors'),
    description: t('about.safety.certifiedInstructorsDesc'),
    icon: 'i-lucide-award',
    color: 'warning'
  },
  {
    title: t('about.safety.safetyTech'),
    description: t('about.safety.safetyTechDesc'),
    icon: 'i-lucide-radar',
    color: 'warning'
  }
]
</script>

<template>
  <div>
    <!-- Dynamic Admin Sections for About Page -->
    <template v-if="aboutPage && aboutPage.sections.length > 0">
      <ContentSectionRenderer
        v-for="section in aboutPage.sections"
        :key="section.id"
        :section="section"
      />
    </template>
    <template v-else>
      <!-- 1. Hero Statement -->
      <UPageHero
        :title="t('about.heroTitle')"
        :description="t('about.heroDesc')"
        align="center"
        class="py-16 md:py-24"
      >
        <template #links>
          <div class="flex items-center gap-2 text-warning font-semibold">
            <UIcon name="i-lucide-leaf" class="size-5" />
            <span>{{ t('about.pioneer') }}</span>
          </div>
        </template>
      </UPageHero>

      <!-- 2. Keamanan & Instruktur -->
      <UPageSection
        :headline="t('about.safety.headline')"
        :title="t('about.safety.title')"
        :description="t('about.safety.description')"
        class="bg-muted/30"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <UPageCard
            v-for="feature in safetyFeatures"
            :key="feature.title"
            :icon="feature.icon"
            :title="feature.title"
            :description="feature.description"
            :ui="{ 
              leadingIcon: 'text-warning text-4xl mb-4',
              title: 'text-xl font-bold mb-2'
            }"
            class="hover:shadow-lg transition-all"
          />
          </div>
      </UPageSection>

      <!-- 4. Visi Singkat -->
      <UPageSection
        class="bg-warning/5 border-y border-warning/10"
      >
        <div class="max-w-3xl mx-auto text-center space-y-8">
          <UIcon name="i-lucide-quote" class="size-12 text-warning mx-auto opacity-50" />
          <div class="space-y-6">
            <p class="text-2xl font-medium leading-relaxed italic text-default">
              {{ t('about.vision.quote') }}
            </p>
            <p class="text-lg text-muted leading-relaxed">
              {{ t('about.vision.description') }}
            </p>
          </div>
          <div class="pt-8 flex flex-wrap justify-center gap-3">
            <UButton :label="t('about.startJourney')" to="/auth/register" color="warning" size="xl" icon="i-lucide-rocket" />
            <UButton :label="t('contact.chatNow')" :to="waLink" target="_blank" color="neutral" size="xl" icon="i-simple-icons-whatsapp" variant="outline" external />
          </div>
        </div>
      </UPageSection>
    </template>
  </div>
</template>
