<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()
const { pages } = useContent()

// the route.params.slug is an array for catch-all [...slug]
const slugPath = '/' + (Array.isArray(route.params.slug) ? route.params.slug.join('/') : (route.params.slug || ''))

const page = computed(() => {
  return pages.value.find(p => p.slug === slugPath && p.status === 'published')
})

if (!page.value && process.server) {
  throw createError({ statusCode: 404, statusMessage: 'Page Not Found' })
}
</script>

<template>
  <div v-if="page">
    <ContentSectionRenderer 
      v-for="section in page.sections" 
      :key="section.id" 
      :section="section" 
    />
    <div v-if="page.sections.length === 0" class="py-20 text-center text-muted">
      <h1 class="text-3xl font-bold mb-4">{{ page.title }}</h1>
      <p>This page has no content yet.</p>
    </div>
  </div>
  <div v-else class="py-32 text-center flex flex-col items-center justify-center">
    <UIcon name="i-lucide-file-question" class="size-16 text-muted mb-4" />
    <h1 class="text-4xl font-bold mb-2">404 - Page Not Found</h1>
    <p class="text-muted text-lg mb-8">The page you are looking for does not exist or has been unpublished.</p>
    <NuxtLink to="/">
      <UButton label="Return Home" color="primary" size="lg" icon="i-lucide-home" />
    </NuxtLink>
  </div>
</template>
