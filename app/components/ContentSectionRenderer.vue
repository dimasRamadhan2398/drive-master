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
      :links="section.data.ctaText ? [
        { label: section.data.ctaText, to: section.data.ctaLink || '#', color: 'warning', size: 'lg' }
      ] : undefined"
    >
      <div v-if="section.data.bgImage" class="relative w-full aspect-video lg:aspect-[4/3] rounded-2xl overflow-hidden bg-elevated shadow-2xl ring ring-default">
        <img 
          :src="section.data.bgImage" 
          alt="Hero Image"
          class="w-full h-full object-cover"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent" />
      </div>
    </UPageHero>

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
