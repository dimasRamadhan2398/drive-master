<script setup lang="ts">
import type { ButtonProps } from '#ui/types'

const { t } = useI18n()
const { waLink } = useSettings()

const heroLinks = computed<ButtonProps[]>(() => [
  { label: t('services.viewPackages'), to: '/packages', color: 'warning', icon: 'i-lucide-package' },
  { label: t('services.bookConsultation'), to: waLink.value, color: 'success', variant: 'outline', icon: 'i-simple-icons-whatsapp', external: true }
])

const ctaLinks = computed<ButtonProps[]>(() => [
  { label: t('services.viewPackages'), to: '/packages', color: 'warning', icon: 'i-lucide-package' },
  { label: t('services.cta.whatsapp'), to: waLink.value, color: 'success', variant: 'outline', icon: 'i-simple-icons-whatsapp', external: true }
])

useSeoMeta({
  title: 'Services | Drive Master Academy',
  description: 'Comprehensive driving courses designed for the electric future. From beginners to advanced drivers, we have the perfect program for you.',
})

const specifications = computed(() => [
  {
    title: t('services.specifications.specialServices'),
    subtitle: t('services.specifications.specialServicesSub'),
    description: [
      t('services.specifications.sim'),
      t('services.specifications.shuttle'),
      t('services.specifications.theory'),
      t('services.specifications.certificate'),
    ],
    icon: 'i-lucide-star'
  },
  {
    title: t('services.specifications.sessionTime'),
    subtitle: t('services.specifications.sessionTimeSub'),
    description: [
      t('services.specifications.weekdayHours'),
      t('services.specifications.sessionDuration'),
      t('home.hoursNight'),
      t('home.hoursWeekend'),
    ],
    icon: 'i-lucide-car'
  },
  {
    title: t('services.specifications.carTransmission'),
    subtitle: t('services.specifications.carTransmissionSub'),
    description: [
      t('services.specifications.matic'),
    ],
    icon: 'i-lucide-car'
  },
  {
    title: t('services.specifications.nightSession'),
    subtitle: t('services.specifications.nightSessionSub'),
    description: [
      t('services.specifications.nightHours'),
      t('services.specifications.nightPrice', { sessions: '6x' }),
      t('services.specifications.nightPrice', { sessions: '8x' }),
      t('services.specifications.nightPrice', { sessions: '10x' }),
      t('services.specifications.nightPrice', { sessions: '12x' }),
    ],
    icon: 'i-lucide-moon'
  },
  {
    title: t('services.specifications.weekendSession'),
    subtitle: t('services.specifications.weekendSessionSub'),
    description: [
      t('services.specifications.weekendHours'),
      t('services.specifications.weekendPrice', { sessions: '6x' }),
      t('services.specifications.weekendPrice', { sessions: '8x' }),
      t('services.specifications.weekendPrice', { sessions: '10x' }),
      t('services.specifications.weekendPrice', { sessions: '12x' }),
    ],
    icon: 'i-lucide-clock'
  }
])

// Course Material
const courseMaterial = computed(() => [
  {
    title: t('home.material.materialTheory'),
    description: [
      t('home.material.materialTheoryDesc'),
      t('home.material.materialTheoryDesc2'),
      t('home.material.materialTheoryDesc3'),
      t('home.material.materialTheoryDesc4'),
    ],
    icon:'i-lucide-book-open'
  },
  {
    title: t('home.material.initialControl'),
    description: [
      t('home.material.initialControlDesc'),
      t('home.material.initialControlDesc2'),
      t('home.material.initialControlDesc3'),
    ],
    icon: 'i-lucide-shield-check'
  },
  {
    title: t('home.material.basicManeuvering'),
    description: [
      t('home.material.basicManeuveringDesc'),
      t('home.material.basicManeuveringDesc2'),
      t('home.material.basicManeuveringDesc3'),
    ],
    icon: 'i-lucide-radar'
  },
  {
    title: t('home.material.uphillDownhill'),
    description: [
      t('home.material.uphillDownhillDesc'),
      t('home.material.uphillDownhillDesc2'),
    ],
    icon: 'i-lucide-car'
  },
  {
    title: t('home.material.parking'),
    description: [
      t('home.material.parkingDesc'),
      t('home.material.parkingDesc2'),
    ],
    icon: 'i-lucide-car'
  },
  {
    title: t('home.material.highway'),
    description: [
      t('home.material.highwayDesc'),
      t('home.material.highwayDesc2'),
      t('home.material.highwayDesc3'),
    ],
    icon: 'i-lucide-car'
  }
])

const serviceAreas = [
  'Alam Sutera & surrounding areas',
  'Serpong & BSD City',
  'Tangerang City Center',
  'Gading Serpong',
  'Lippo Karawaci',
  'Bintaro Jaya (limited)'
]

const contentStore = useContentStore()
const { pages } = useContent()

await useAsyncData('services-page-sections', () => contentStore.fetchPages())

const servicesPage = computed(() =>
  pages.value.find((p) => p.slug === '/services' && p.status === 'published')
)
</script>

<template>
  <div>
    <!-- Dynamic Admin Sections for Services Page -->
    <template v-if="servicesPage && servicesPage.sections.length > 0">
      <ContentSectionRenderer
        v-for="section in servicesPage.sections"
        :key="section.id"
        :section="section"
      />
    </template>
    <template v-else>
      <!-- Hero -->
      <UPageHero
        :title="t('services.heroTitle')"
        :description="t('services.heroDesc')"
        :links="heroLinks"
      />

    <!-- Our Services -->
    <UPageSection
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="grid md:grid-cols-1 gap-6">
        <UCard v-for="specification in specifications" :key="specification.title">
          <template #header>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-xl bg-warning/10">
                <UIcon :name="specification.icon" class="size-8 text-warning" />
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ specification.title }}</h3>
                <p class="text-muted mt-1">{{ specification.subtitle }}</p>
              </div>
            </div>
          </template>

          <ul class="space-y-2">
            <li v-for="spec in specification.description" :key="spec" class="flex items-center gap-2 font-bold">
              <span class="text-sm">{{ spec }}</span>
            </li>
          </ul>

        </UCard>
      </div>
    </UPageSection>

    <!-- Course Material -->
    <UPageSection
      id="material"
      :headline="t('home.courseMaterialHeadline')"
      :title="t('home.courseMaterial')"
      :description="t('home.courseMaterialDesc')"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <UPageGrid>
        <UPageCard
          v-for="material in courseMaterial"
          :key="material.title"
          :icon="material.icon"
          :title="material.title"
          :ui="{ leadingIcon: 'text-warning text-3xl' }"
        >
          <template #description>
            <ul class="space-y-2 mt-2">
              <li v-for="item in material.description" :key="item" class="flex items-start gap-2 text-muted">
                <UIcon name="i-lucide-check-circle" class="size-4 text-warning shrink-0 mt-1" />
                <span class="text-sm">{{ item }}</span>
              </li>
            </ul>
          </template>
        </UPageCard>
      </UPageGrid>
    </UPageSection>

    <!-- Service Areas -->
    <UPageSection
      :headline="t('services.coverage.title')"
      :title="t('services.coverage.headline')"
      :description="t('services.coverage.description')"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="max-w-3xl mx-auto">
        <UCard>
          <div class="grid sm:grid-cols-2 gap-4">
            <div v-for="area in serviceAreas" :key="area" class="flex items-center gap-3 p-3 rounded-lg hover:bg-muted/50 transition-colors">
              <UIcon name="i-lucide-map-pin" class="size-5 text-warning" />
              <span>{{ area }}</span>
            </div>
          </div>
        </UCard>
        
        <p class="text-center text-muted mt-6">
          {{ t('services.coverage.footer') }}
        </p>
      </div>
    </UPageSection>

    <!-- CTA -->
    <UPageCTA
      :title="t('services.cta.title')"
      :description="t('services.cta.description')"
      :links="ctaLinks"
    />
    </template>
  </div>
</template>
