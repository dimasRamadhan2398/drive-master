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
        ...(section.data.ctaText ? [{ label: section.data.ctaText, to: section.data.ctaLink || '#', color: 'warning', size: 'lg', icon: 'i-lucide-calendar-check' }] : []),
        ...(section.data.secondaryCtaText ? [{ label: section.data.secondaryCtaText, to: section.data.secondaryCtaLink || '#', color: 'neutral', variant: 'outline', trailingIcon: 'i-lucide-arrow-right', size: 'lg' }] : [])
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

    <!-- CTA Section -->
    <UPageCTA
      v-else-if="section.type === 'cta'"
      :title="section.data.heading"
      :links="[
        { label: section.data.buttonText, to: section.data.buttonLink || '#', color: 'warning', size: 'lg' }
      ]"
      class="bg-warning/5"
    />
  </div>
</template>
