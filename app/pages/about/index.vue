<script setup lang="ts">
// Page metadata
definePageMeta({
  layout: 'default'
})

useSeoMeta({
  title: 'About Us | Drive Master Indonesia - The Future of Driving',
  description: 'Drive Master is the first premium driving academy in Alam Sutera using 100% Electric Vehicles.'
})

const safetyFeatures = [
  {
    title: 'Certified Instructors',
    description: 'Our instructors are licensed professionals specially certified to operate premium electric vehicles.',
    icon: 'i-lucide-award',
    color: 'warning'
  },
  {
    title: 'Active Safety Tech',
    description: 'Leveraging the inherent safety features of EVs such as Collision Avoidance and Blind Spot Monitoring to minimize risks.',
    icon: 'i-lucide-radar',
    color: 'warning'
  }
]
</script>

<template>
  <div>
    <!-- 1. Hero Statement -->
    <UPageHero
      title="Step Forward, Learn from the Future"
      description="Drive Master is not just about teaching how to drive, it's about redefining the standards of driving education in Indonesia. As a pioneer of Electric Vehicle driving schools, we believe that the drivers of the future must be born from the future of technology—modern, intelligent, and environmentally friendly."
      align="center"
      class="py-16 md:py-24"
    >
      <template #links>
        <div class="flex items-center gap-2 text-warning font-semibold">
          <UIcon name="i-lucide-leaf" class="size-5" />
          <span>Pioneer of Emission-Free Driving Education</span>
        </div>
      </template>
    </UPageHero>

    <!-- 2. Keamanan & Instruktur -->
    <UPageSection
      headline="Safety First"
      title="Our Priority is Your Safety"
      description="In the driving course business, safety is not just a feature; it is our core foundation."
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
            "Our vision is not just to produce drivers who can operate a steering wheel, but to nurture intelligent, safe drivers ready to embrace the electrification era."
          </p>
          <p class="text-lg text-muted leading-relaxed">
            At Drive Master Indonesia, we believe that how we learn to drive must evolve with the evolution of automotive technology. We are committed to becoming the new standard in environmentally friendly driving education, ensuring that every graduate has high-level driving skills as well as an awareness of the future of sustainable mobility.
          </p>
        </div>
        <div class="pt-8">
          <UButton label="Start Your Journey" to="/auth/register" color="warning" size="xl" icon="i-lucide-rocket" />
        </div>
      </div>
    </UPageSection>
  </div>
</template>
