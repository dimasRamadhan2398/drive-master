<script setup lang="ts">
import { ref, computed } from 'vue'

const { blogPosts } = useContent()
const searchQuery = ref('')

// Only show published posts
const publishedPosts = computed(() =>
  blogPosts.value.filter(p => p.status === 'published')
)

// Search filter
const filteredPosts = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return publishedPosts.value
  return publishedPosts.value.filter(
    p =>
      p.title.toLowerCase().includes(q) ||
      p.author.toLowerCase().includes(q) ||
      (p.content && p.content.toLowerCase().includes(q))
  )
})

// Featured post = latest published post
const featuredPost = computed(() => filteredPosts.value[0] || null)

// Remaining posts (exclude featured)
const remainingPosts = computed(() => filteredPosts.value.slice(1))

// Helper: generate a slug from post title
function postSlug(post: any) {
  return post.title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

// Helper: get first image from media
function getPostThumbnail(post: any): string | null {
  if (!post.media || post.media.length === 0) return null
  const image = post.media.find((m: any) => m.fileType === 'image')
  return image ? image.url : null
}

// Helper: truncate content for excerpt
function getExcerpt(content: string, maxLength = 150): string {
  if (!content) return 'No content available...'
  return content.length > maxLength
    ? content.substring(0, maxLength).trim() + '...'
    : content
}

// Helper: reading time estimate
function getReadingTime(content: string): string {
  if (!content) return '1 min read'
  const words = content.split(/\s+/).length
  const minutes = Math.max(1, Math.ceil(words / 200))
  return `${minutes} min read`
}

// Helper: get author initials for avatar
function getAuthorInitials(author: string): string {
  return author
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .substring(0, 2)
}

// Placeholder gradient colors for posts without images
const gradientColors = [
  'from-amber-500/20 to-orange-600/20',
  'from-emerald-500/20 to-teal-600/20',
  'from-violet-500/20 to-purple-600/20',
  'from-sky-500/20 to-blue-600/20',
  'from-rose-500/20 to-pink-600/20',
  'from-yellow-500/20 to-amber-600/20',
]

function getGradient(index: number) {
  return gradientColors[index % gradientColors.length]
}

useSeoMeta({
  title: 'Articles | Drive Master Academy',
  description: 'Read the latest articles, tips, and updates about EV driving, electric vehicles, and driving education from Drive Master Academy.',
})
</script>

<template>
  <div>
    <!-- Hero Section -->
    <section class="relative overflow-hidden bg-gradient-to-br from-warning/5 via-background to-warning/10">
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-24 -right-24 size-96 bg-warning/5 rounded-full blur-3xl" />
        <div class="absolute -bottom-32 -left-32 size-[500px] bg-warning/3 rounded-full blur-3xl" />
      </div>

      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-24">
        <div class="text-center max-w-3xl mx-auto">
          <UBadge label="Blog & Articles" color="warning" variant="subtle" class="mb-4" />
          <h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight mb-6">
            Insights & 
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-warning to-amber-500">Stories</span>
          </h1>
          <p class="text-lg text-muted max-w-2xl mx-auto mb-10">
            Stay updated with the latest tips on driving education, electric vehicle trends, and stories from our community of learners.
          </p>

          <!-- Search Bar -->
          <div class="max-w-lg mx-auto">
            <UInput
              v-model="searchQuery"
              placeholder="Search articles..."
              icon="i-lucide-search"
              size="xl"
              class="w-full"
              :ui="{ base: 'rounded-full' }"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- Content -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 lg:py-16">

      <!-- Empty State -->
      <div v-if="filteredPosts.length === 0" class="py-24 text-center">
        <div class="inline-flex items-center justify-center size-20 rounded-full bg-muted/30 mb-6">
          <UIcon name="i-lucide-newspaper" class="size-10 text-muted" />
        </div>
        <h2 class="text-2xl font-bold mb-3">
          {{ searchQuery ? 'No articles found' : 'No articles yet' }}
        </h2>
        <p class="text-muted max-w-md mx-auto mb-6">
          {{ searchQuery 
            ? `We couldn't find any articles matching "${searchQuery}". Try a different search term.` 
            : 'Stay tuned! Our team is working on exciting content for you.' 
          }}
        </p>
        <UButton
          v-if="searchQuery"
          label="Clear Search"
          color="warning"
          variant="outline"
          icon="i-lucide-x"
          @click="searchQuery = ''"
        />
      </div>

      <template v-else>
        <!-- Featured Post -->
        <NuxtLink
          v-if="featuredPost"
          :to="`/blog/${featuredPost.id}-${postSlug(featuredPost)}`"
          class="block group mb-12"
        >
          <div class="grid lg:grid-cols-2 gap-6 lg:gap-10 items-center p-2 rounded-3xl border border-default bg-elevated/50 hover:bg-elevated transition-all duration-300 hover:shadow-xl hover:shadow-warning/5">
            <!-- Image / Thumbnail -->
            <div class="relative aspect-[16/10] rounded-2xl overflow-hidden">
              <img
                v-if="getPostThumbnail(featuredPost)"
                :src="getPostThumbnail(featuredPost)!"
                :alt="featuredPost.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              />
              <div
                v-else
                :class="[
                  'w-full h-full bg-gradient-to-br flex items-center justify-center',
                  'from-warning/10 to-amber-600/20'
                ]"
              >
                <UIcon name="i-lucide-newspaper" class="size-20 text-warning/30" />
              </div>
              <div class="absolute top-4 left-4">
                <UBadge label="Featured" color="warning" class="shadow-lg" />
              </div>
            </div>

            <!-- Content -->
            <div class="p-4 lg:p-6 lg:pr-10">
              <div class="flex items-center gap-3 mb-4 text-sm text-muted">
                <div class="flex items-center gap-2">
                  <UAvatar :text="getAuthorInitials(featuredPost.author)" size="xs" />
                  <span>{{ featuredPost.author }}</span>
                </div>
                <span class="size-1 rounded-full bg-muted" />
                <span>{{ featuredPost.date }}</span>
                <span class="size-1 rounded-full bg-muted" />
                <span class="flex items-center gap-1">
                  <UIcon name="i-lucide-clock" class="size-3.5" />
                  {{ getReadingTime(featuredPost.content) }}
                </span>
              </div>

              <h2 class="text-2xl lg:text-3xl font-bold mb-4 group-hover:text-warning transition-colors line-clamp-2">
                {{ featuredPost.title }}
              </h2>

              <p class="text-muted leading-relaxed mb-6 line-clamp-3">
                {{ getExcerpt(featuredPost.content, 250) }}
              </p>

              <div class="flex items-center gap-6">
                <span class="inline-flex items-center gap-2 text-warning font-medium group-hover:gap-3 transition-all">
                  Read Article
                  <UIcon name="i-lucide-arrow-right" class="size-4" />
                </span>
                <span class="flex items-center gap-1.5 text-sm text-muted">
                  <UIcon name="i-lucide-eye" class="size-4" />
                  {{ featuredPost.views }} views
                </span>
              </div>
            </div>
          </div>
        </NuxtLink>

        <!-- Section Divider -->
        <div v-if="remainingPosts.length > 0" class="flex items-center gap-4 mb-10">
          <h2 class="text-xl font-bold whitespace-nowrap">Latest Articles</h2>
          <div class="h-px flex-1 bg-default" />
          <span class="text-sm text-muted whitespace-nowrap">
            {{ publishedPosts.length }} {{ publishedPosts.length === 1 ? 'article' : 'articles' }}
          </span>
        </div>

        <!-- Blog Grid -->
        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          <NuxtLink
            v-for="(post, index) in remainingPosts"
            :key="post.id"
            :to="`/blog/${post.id}-${postSlug(post)}`"
            class="group"
          >
            <div class="flex flex-col h-full rounded-2xl border border-default bg-elevated/30 overflow-hidden hover:bg-elevated hover:shadow-lg hover:shadow-warning/5 transition-all duration-300 hover:-translate-y-1">
              <!-- Image -->
              <div class="relative aspect-[16/10] overflow-hidden">
                <img
                  v-if="getPostThumbnail(post)"
                  :src="getPostThumbnail(post)!"
                  :alt="post.title"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <div
                  v-else
                  :class="[
                    'w-full h-full bg-gradient-to-br flex items-center justify-center',
                    getGradient(index + 1)
                  ]"
                >
                  <UIcon name="i-lucide-file-text" class="size-12 text-muted/30" />
                </div>
              </div>

              <!-- Body -->
              <div class="flex flex-col flex-1 p-5">
                <!-- Meta -->
                <div class="flex items-center gap-2 mb-3 text-xs text-muted">
                  <span class="flex items-center gap-1.5">
                    <UAvatar :text="getAuthorInitials(post.author)" size="3xs" />
                    {{ post.author }}
                  </span>
                  <span class="size-1 rounded-full bg-muted" />
                  <span>{{ post.date }}</span>
                </div>

                <!-- Title -->
                <h3 class="text-lg font-semibold mb-2 group-hover:text-warning transition-colors line-clamp-2">
                  {{ post.title }}
                </h3>

                <!-- Excerpt -->
                <p class="text-sm text-muted leading-relaxed flex-1 line-clamp-3 mb-4">
                  {{ getExcerpt(post.content) }}
                </p>

                <!-- Footer -->
                <div class="flex items-center justify-between pt-4 border-t border-default">
                  <span class="flex items-center gap-1.5 text-xs text-muted">
                    <UIcon name="i-lucide-clock" class="size-3.5" />
                    {{ getReadingTime(post.content) }}
                  </span>
                  <span class="flex items-center gap-1.5 text-xs text-muted">
                    <UIcon name="i-lucide-eye" class="size-3.5" />
                    {{ post.views }}
                  </span>
                  <span class="inline-flex items-center gap-1 text-xs font-medium text-warning opacity-0 group-hover:opacity-100 transition-opacity">
                    Read
                    <UIcon name="i-lucide-arrow-right" class="size-3.5" />
                  </span>
                </div>
              </div>
            </div>
          </NuxtLink>
        </div>
      </template>
    </section>

    <!-- Newsletter CTA -->
    <section class="bg-gradient-to-br from-warning/5 via-background to-warning/10">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-20 text-center">
        <div class="inline-flex items-center justify-center size-16 rounded-2xl bg-warning/10 mb-6">
          <UIcon name="i-lucide-bell-ring" class="size-8 text-warning" />
        </div>
        <h2 class="text-3xl font-bold mb-4">Stay in the Loop</h2>
        <p class="text-muted max-w-xl mx-auto mb-8">
          Get the latest driving tips, EV insights, and academy updates delivered straight to your inbox.
        </p>
        <div class="flex flex-col sm:flex-row gap-3 justify-center max-w-md mx-auto">
          <UInput placeholder="Enter your email" size="lg" class="flex-1" icon="i-lucide-mail" />
          <UButton label="Subscribe" color="warning" size="lg" icon="i-lucide-send" />
        </div>
        <p class="text-xs text-muted mt-4">No spam, unsubscribe anytime.</p>
      </div>
    </section>
  </div>
</template>
