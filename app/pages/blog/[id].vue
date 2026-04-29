<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed, ref } from 'vue'

const route = useRoute()
const { blogPosts } = useContent()

// Extract post ID from route param (format: "1-slug-title")
const postId = computed(() => {
  const param = route.params.id as string
  return parseInt(param.split('-')[0])
})

const post = computed(() =>
  blogPosts.value.find(p => p.id === postId.value && p.status === 'published')
)

// Increment views on load
if (post.value) {
  post.value.views = (post.value.views || 0) + 1
}

// Related posts (exclude current)
const relatedPosts = computed(() =>
  blogPosts.value
    .filter(p => p.id !== postId.value && p.status === 'published')
    .slice(0, 3)
)

// Media lightbox
const lightboxOpen = ref(false)
const lightboxIndex = ref(0)

function openLightbox(index: number) {
  lightboxIndex.value = index
  lightboxOpen.value = true
}

// Helpers
function getAuthorInitials(author: string): string {
  return author
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .substring(0, 2)
}

function getReadingTime(content: string): string {
  if (!content) return '1 min read'
  const words = content.split(/\s+/).length
  const minutes = Math.max(1, Math.ceil(words / 200))
  return `${minutes} min read`
}

function postSlug(p: any) {
  return p.title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

function getExcerpt(content: string, maxLength = 120): string {
  if (!content) return 'No content available...'
  return content.length > maxLength
    ? content.substring(0, maxLength).trim() + '...'
    : content
}

function getPostThumbnail(p: any): string | null {
  if (!p.media || p.media.length === 0) return null
  const image = p.media.find((m: any) => m.fileType === 'image')
  return image ? image.url : null
}

const gradientColors = [
  'from-amber-500/20 to-orange-600/20',
  'from-emerald-500/20 to-teal-600/20',
  'from-violet-500/20 to-purple-600/20',
]

function getGradient(index: number) {
  return gradientColors[index % gradientColors.length]
}

// Format content paragraphs
const contentParagraphs = computed(() => {
  if (!post.value?.content) return []
  return post.value.content.split('\n').filter((p: string) => p.trim())
})

// Images from media
const postImages = computed(() => {
  if (!post.value?.media) return []
  return post.value.media.filter(m => m.fileType === 'image')
})

// Videos from media
const postVideos = computed(() => {
  if (!post.value?.media) return []
  return post.value.media.filter(m => m.fileType === 'video')
})

// SEO
if (post.value) {
  useSeoMeta({
    title: `${post.value.title} - Drive Master Blog`,
    description: getExcerpt(post.value.content, 160),
  })
}

if (!post.value && import.meta.server) {
  throw createError({ statusCode: 404, statusMessage: 'Post Not Found' })
}
</script>

<template>
  <div>
    <!-- 404 State -->
    <div v-if="!post" class="py-32 text-center flex flex-col items-center justify-center">
      <UIcon name="i-lucide-file-question" class="size-16 text-muted mb-4" />
      <h1 class="text-4xl font-bold mb-2">Post Not Found</h1>
      <p class="text-muted text-lg mb-8">
        This article doesn't exist or has been unpublished.
      </p>
      <NuxtLink to="/blog">
        <UButton label="Back to Blog" color="warning" size="lg" icon="i-lucide-arrow-left" />
      </NuxtLink>
    </div>

    <template v-else>
      <!-- Hero / Header -->
      <section class="relative overflow-hidden bg-gradient-to-br from-warning/5 via-background to-warning/10">
        <div class="absolute inset-0 pointer-events-none">
          <div class="absolute -top-24 -right-24 size-96 bg-warning/5 rounded-full blur-3xl" />
        </div>

        <div class="relative max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 lg:py-20">
          <!-- Breadcrumb -->
          <nav class="flex items-center gap-2 text-sm text-muted mb-8">
            <NuxtLink to="/" class="hover:text-foreground transition-colors">Home</NuxtLink>
            <UIcon name="i-lucide-chevron-right" class="size-4" />
            <NuxtLink to="/blog" class="hover:text-foreground transition-colors">Blog</NuxtLink>
            <UIcon name="i-lucide-chevron-right" class="size-4" />
            <span class="text-foreground truncate max-w-[200px]">{{ post.title }}</span>
          </nav>

          <!-- Title -->
          <h1 class="text-3xl sm:text-4xl lg:text-5xl font-bold tracking-tight mb-6 leading-tight">
            {{ post.title }}
          </h1>

          <!-- Meta -->
          <div class="flex flex-wrap items-center gap-4 text-sm text-muted">
            <div class="flex items-center gap-3">
              <UAvatar :text="getAuthorInitials(post.author)" size="md" />
              <div>
                <p class="font-medium text-foreground">{{ post.author }}</p>
                <p class="text-xs">Author</p>
              </div>
            </div>

            <span class="h-8 w-px bg-default" />

            <div class="flex items-center gap-1.5">
              <UIcon name="i-lucide-calendar" class="size-4" />
              <span>{{ post.date }}</span>
            </div>

            <span class="h-8 w-px bg-default hidden sm:block" />

            <div class="flex items-center gap-1.5">
              <UIcon name="i-lucide-clock" class="size-4" />
              <span>{{ getReadingTime(post.content) }}</span>
            </div>

            <span class="h-8 w-px bg-default hidden sm:block" />

            <div class="flex items-center gap-1.5">
              <UIcon name="i-lucide-eye" class="size-4" />
              <span>{{ post.views }} views</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Article Content -->
      <article class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 lg:py-16">
        <!-- Featured Image -->
        <div v-if="postImages.length > 0" class="mb-10 rounded-2xl overflow-hidden border border-default shadow-lg">
          <img
            :src="postImages[0].url"
            :alt="post.title"
            class="w-full max-h-[500px] object-cover cursor-pointer hover:scale-[1.02] transition-transform duration-500"
            @click="openLightbox(0)"
          />
        </div>

        <!-- Text Content -->
        <div class="prose prose-lg dark:prose-invert max-w-none">
          <p
            v-for="(paragraph, idx) in contentParagraphs"
            :key="idx"
            class="text-base sm:text-lg leading-relaxed text-muted mb-6"
          >
            {{ paragraph }}
          </p>
        </div>

        <!-- Media Gallery -->
        <div v-if="postImages.length > 1 || postVideos.length > 0" class="mt-12">
          <h3 class="text-xl font-bold mb-6 flex items-center gap-2">
            <UIcon name="i-lucide-images" class="size-5 text-warning" />
            Media Gallery
          </h3>

          <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <!-- Images (skip first if already shown as featured) -->
            <div
              v-for="(img, idx) in postImages.slice(1)"
              :key="'img-' + idx"
              class="relative group aspect-[4/3] rounded-xl overflow-hidden border border-default cursor-pointer"
              @click="openLightbox(idx + 1)"
            >
              <img
                :src="img.url"
                :alt="img.name"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              />
              <div class="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors flex items-center justify-center">
                <UIcon name="i-lucide-zoom-in" class="size-8 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
              </div>
            </div>

            <!-- Videos -->
            <div
              v-for="(vid, idx) in postVideos"
              :key="'vid-' + idx"
              class="relative aspect-[4/3] rounded-xl overflow-hidden border border-default bg-black"
            >
              <video
                :src="vid.url"
                controls
                class="w-full h-full object-contain"
              />
            </div>
          </div>
        </div>

        <!-- Share / Actions -->
        <div class="mt-12 pt-8 border-t border-default">
          <div class="flex flex-wrap items-center justify-between gap-4">
            <div class="flex items-center gap-3">
              <span class="text-sm text-muted font-medium">Share this article:</span>
              <UButton icon="i-lucide-link" color="neutral" variant="ghost" size="sm" />
              <UButton icon="i-simple-icons-whatsapp" color="neutral" variant="ghost" size="sm" />
              <UButton icon="i-simple-icons-x" color="neutral" variant="ghost" size="sm" />
            </div>
            <NuxtLink to="/blog">
              <UButton label="Back to Blog" color="warning" variant="outline" icon="i-lucide-arrow-left" />
            </NuxtLink>
          </div>
        </div>
      </article>

      <!-- Related Posts -->
      <section v-if="relatedPosts.length > 0" class="bg-muted/20">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-20">
          <div class="text-center mb-10">
            <h2 class="text-2xl font-bold mb-2">More Articles</h2>
            <p class="text-muted">Continue reading from our blog</p>
          </div>

          <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
            <NuxtLink
              v-for="(related, index) in relatedPosts"
              :key="related.id"
              :to="`/blog/${related.id}-${postSlug(related)}`"
              class="group"
            >
              <div class="flex flex-col h-full rounded-2xl border border-default bg-elevated/30 overflow-hidden hover:bg-elevated hover:shadow-lg hover:shadow-warning/5 transition-all duration-300 hover:-translate-y-1">
                <!-- Image -->
                <div class="relative aspect-[16/10] overflow-hidden">
                  <img
                    v-if="getPostThumbnail(related)"
                    :src="getPostThumbnail(related)!"
                    :alt="related.title"
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                  />
                  <div
                    v-else
                    :class="[
                      'w-full h-full bg-gradient-to-br flex items-center justify-center',
                      getGradient(index)
                    ]"
                  >
                    <UIcon name="i-lucide-file-text" class="size-12 text-muted/30" />
                  </div>
                </div>

                <!-- Body -->
                <div class="flex flex-col flex-1 p-5">
                  <div class="flex items-center gap-2 mb-3 text-xs text-muted">
                    <span>{{ related.author }}</span>
                    <span class="size-1 rounded-full bg-muted" />
                    <span>{{ related.date }}</span>
                  </div>
                  <h3 class="text-lg font-semibold mb-2 group-hover:text-warning transition-colors line-clamp-2">
                    {{ related.title }}
                  </h3>
                  <p class="text-sm text-muted flex-1 line-clamp-2">
                    {{ getExcerpt(related.content) }}
                  </p>
                </div>
              </div>
            </NuxtLink>
          </div>
        </div>
      </section>

      <!-- Lightbox Modal -->
      <ClientOnly>
        <UModal v-model:open="lightboxOpen">
          <template #content>
            <div class="relative bg-black rounded-2xl overflow-hidden">
              <div class="flex items-center justify-between p-4">
                <span class="text-white text-sm">
                  {{ lightboxIndex + 1 }} / {{ postImages.length }}
                </span>
                <UButton
                  icon="i-lucide-x"
                  color="neutral"
                  variant="ghost"
                  class="text-white"
                  @click="lightboxOpen = false"
                />
              </div>
              <div class="flex items-center justify-center p-4 pt-0">
                <img
                  v-if="postImages[lightboxIndex]"
                  :src="postImages[lightboxIndex].url"
                  :alt="postImages[lightboxIndex].name"
                  class="max-h-[70vh] max-w-full object-contain rounded-lg"
                />
              </div>
              <div class="flex justify-center gap-4 p-4">
                <UButton
                  icon="i-lucide-chevron-left"
                  color="neutral"
                  variant="ghost"
                  class="text-white"
                  :disabled="lightboxIndex === 0"
                  @click="lightboxIndex--"
                />
                <UButton
                  icon="i-lucide-chevron-right"
                  color="neutral"
                  variant="ghost"
                  class="text-white"
                  :disabled="lightboxIndex >= postImages.length - 1"
                  @click="lightboxIndex++"
                />
              </div>
            </div>
          </template>
        </UModal>
      </ClientOnly>
    </template>
  </div>
</template>
