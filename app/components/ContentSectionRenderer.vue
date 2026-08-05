<script setup lang="ts">
defineProps<{
  section: {
    type: string
    data: any
  }
}>()
</script>

<template>
  <div>
    <!-- Hero Section -->
    <UPageHero
      v-if="section.type === 'hero'"
      :title="section.data.heading"
      :description="section.data.subheading"
      orientation="horizontal"
      :links="[
        ...(section.data.ctaText ? [{ label: section.data.ctaText, to: section.data.ctaLink || '#', color: 'warning' as const, size: 'lg' as const, icon: 'i-lucide-calendar-check' }] : []),
        ...(section.data.secondaryCtaText ? [{ label: section.data.secondaryCtaText, to: section.data.secondaryCtaLink || '#', color: 'neutral' as const, variant: 'outline' as const, trailingIcon: 'i-lucide-arrow-right', size: 'lg' as const }] : [])
      ]"
    >
      <div v-if="section.data.bgImage" class="relative w-full aspect-video lg:aspect-[4/3] rounded-2xl overflow-hidden bg-elevated shadow-2xl ring ring-default">
        <img 
          :src="section.data.bgImage" 
          alt="Hero Image"
          class="w-full h-full object-cover"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent" />
        <!-- Hero Image Overlay Features -->
        <div
          v-if="section.data.features && section.data.features.length > 0"
          class="absolute bottom-4 left-4 right-4 flex items-center justify-between gap-4"
        >
          <div
            v-for="feat in section.data.features"
            :key="feat.title"
            class="flex items-center gap-2 bg-black/60 backdrop-blur-md rounded-full px-4 py-2 text-white text-md"
          >
            <UIcon
              :name="feat.icon || 'i-lucide-check-circle'"
              class="size-5 text-warning"
            />
            <span>{{ feat.title }}</span>
          </div>
        </div>
      </div>
    </UPageHero>

    <!-- Course Material Section -->
    <UPageSection
      v-else-if="section.type === 'course_material'"
      id="material"
      :headline="section.data.headline"
      :title="section.data.title"
      :description="section.data.description"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <UPageGrid>
        <UPageCard
          v-for="material in section.data.materials"
          :key="material.title"
          :icon="material.icon || 'i-lucide-book-open'"
          :title="material.title"
          :ui="{ leadingIcon: 'text-warning text-3xl' }"
        >
          <template #description>
            <ul class="space-y-2 mt-2">
              <li
                v-for="item in material.description"
                :key="item"
                class="flex items-start gap-2 text-muted"
              >
                <UIcon
                  name="i-lucide-check-circle"
                  class="size-4 text-warning shrink-0 mt-1"
                />
                <span class="text-sm">{{ item }}</span>
              </li>
            </ul>
          </template>
        </UPageCard>
      </UPageGrid>
    </UPageSection>

    <!-- Text Block Section -->
    <UPageSection v-else-if="section.type === 'text'" class="bg-background">
      <div class="max-w-4xl mx-auto prose dark:prose-invert" v-html="section.data.content"></div>
    </UPageSection>

    <!-- Image + Text Section -->
    <UPageSection v-else-if="section.type === 'image_text'" class="bg-muted/30">
      <div class="grid lg:grid-cols-2 gap-8 items-center">
        <div class="rounded-2xl overflow-hidden shadow-lg border border-default h-full min-h-[300px]">
          <img :src="section.data.image" class="w-full h-full object-cover" />
        </div>
        <div class="prose dark:prose-invert text-lg" v-html="section.data.content"></div>
      </div>
    </UPageSection>

    <!-- Specifications Section -->
    <UPageSection
      v-else-if="section.type === 'specifications'"
      :headline="section.data.headline"
      :title="section.data.title"
      :description="section.data.description"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="grid md:grid-cols-1 gap-6">
        <UCard v-for="spec in section.data.items" :key="spec.title">
          <template #header>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-xl bg-warning/10">
                <UIcon :name="spec.icon || 'i-lucide-star'" class="size-8 text-warning" />
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ spec.title }}</h3>
                <p v-if="spec.subtitle" class="text-muted mt-1">{{ spec.subtitle }}</p>
              </div>
            </div>
          </template>

          <ul v-if="spec.description && spec.description.length > 0" class="space-y-2">
            <li v-for="bullet in spec.description" :key="bullet" class="flex items-center gap-2 font-bold">
              <span class="text-sm">{{ bullet }}</span>
            </li>
          </ul>
        </UCard>
      </div>
    </UPageSection>

    <!-- Service Areas Section -->
    <UPageSection
      v-else-if="section.type === 'service_areas'"
      :headline="section.data.headline"
      :title="section.data.title"
      :description="section.data.description"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="max-w-3xl mx-auto">
        <UCard>
          <div class="grid sm:grid-cols-2 gap-4">
            <div v-for="area in section.data.areas" :key="area" class="flex items-center gap-3 p-3 rounded-lg hover:bg-muted/50 transition-colors">
              <UIcon name="i-lucide-map-pin" class="size-5 text-warning" />
              <span>{{ area }}</span>
            </div>
          </div>
        </UCard>
        
        <p v-if="section.data.footer" class="text-center text-muted mt-6">
          {{ section.data.footer }}
        </p>
      </div>
    </UPageSection>

    <!-- Quote Section -->
    <UPageSection
      v-else-if="section.type === 'quote'"
      class="bg-warning/5 border-y border-warning/10"
    >
      <div class="max-w-3xl mx-auto text-center space-y-8">
        <UIcon name="i-lucide-quote" class="size-12 text-warning mx-auto opacity-50" />
        <div class="space-y-6">
          <p class="text-2xl font-medium leading-relaxed italic text-default">
            {{ section.data.quote }}
          </p>
          <p v-if="section.data.description" class="text-lg text-muted leading-relaxed">
            {{ section.data.description }}
          </p>
        </div>
        <div v-if="section.data.ctaText || section.data.secondaryCtaText" class="pt-8 flex flex-wrap justify-center gap-3">
          <UButton
            v-if="section.data.ctaText"
            :label="section.data.ctaText"
            :to="section.data.ctaLink || '#'"
            color="warning"
            size="xl"
            icon="i-lucide-rocket"
          />
          <UButton
            v-if="section.data.secondaryCtaText"
            :label="section.data.secondaryCtaText"
            :to="section.data.secondaryCtaLink || '#'"
            :color="(section.data.secondaryCtaLink?.includes('wa.me') || section.data.secondaryCtaText?.toLowerCase().includes('whatsapp')) ? 'success' : 'neutral'"
            size="xl"
            icon="i-simple-icons-whatsapp"
            variant="outline"
          />
        </div>
      </div>
    </UPageSection>

    <!-- CTA Section -->
    <UPageCTA
      v-else-if="section.type === 'cta'"
      :title="section.data.heading"
      :description="section.data.description"
      :links="[
        ...(section.data.buttonText ? [{
          label: section.data.buttonText,
          to: section.data.buttonLink || '#',
          color: (section.data.buttonIcon?.includes('whatsapp') || section.data.buttonLink?.includes('wa.me')) ? ('success' as const) : ('warning' as const),
          variant: (section.data.buttonIcon?.includes('whatsapp') || section.data.buttonLink?.includes('wa.me')) ? ('outline' as const) : undefined,
          size: 'lg' as const,
          icon: section.data.buttonIcon || undefined
        }] : []),
        ...(section.data.secondaryButtonText ? [{
          label: section.data.secondaryButtonText,
          to: section.data.secondaryButtonLink || '#',
          color: (section.data.secondaryButtonIcon?.includes('whatsapp') || section.data.secondaryButtonLink?.includes('wa.me')) ? ('success' as const) : ('neutral' as const),
          variant: 'outline' as const,
          size: 'lg' as const,
          icon: section.data.secondaryButtonIcon || undefined
        }] : [])
      ]"
      class="bg-warning/5"
    />
  </div>
</template>
