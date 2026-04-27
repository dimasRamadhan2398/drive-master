<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { ref } from 'vue'

definePageMeta({ layout: 'admin' })

const toast = useToast()
const activeTab = ref('pages')
const showEditModal = ref(false)

// Mock pages data
const pages = ref([
  { id: 1, title: 'Home Page', slug: '/', lastUpdated: 'Mar 25, 2026', status: 'published' },
  { id: 2, title: 'Services', slug: '/services', lastUpdated: 'Mar 20, 2026', status: 'published' },
  { id: 3, title: 'Packages', slug: '/packages', lastUpdated: 'Mar 18, 2026', status: 'published' },
  { id: 4, title: 'About Us', slug: '/about', lastUpdated: 'Mar 15, 2026', status: 'draft' }
])

// Mock blog posts
const blogPosts = ref([
  { id: 1, title: '5 Tips for Learning to Drive in an EV', author: 'Admin', date: 'Mar 25, 2026', status: 'published', views: 234 },
  { id: 2, title: 'The Benefits of Electric Vehicle Driving', author: 'Admin', date: 'Mar 20, 2026', status: 'published', views: 189 },
  { id: 3, title: 'Understanding One-Pedal Driving', author: 'Admin', date: 'Mar 15, 2026', status: 'draft', views: 0 }
])

// Mock FAQs
const faqs = ref([
  { id: 1, question: 'Do I need prior driving experience?', answer: 'No, our courses are designed for complete beginners...', order: 1 },
  { id: 2, question: 'What type of electric vehicles do you use?', answer: 'We use premium electric vehicles including Tesla Model 3...', order: 2 },
  { id: 3, question: 'How do I book my training sessions?', answer: 'After registering and purchasing a package...', order: 3 }
])

const tabs = [
  { label: 'Pages', value: 'pages', icon: 'i-lucide-file-text' },
  { label: 'Blog Posts', value: 'blog', icon: 'i-lucide-newspaper' },
  { label: 'FAQs', value: 'faq', icon: 'i-lucide-help-circle' }
]
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="Content Management">
        <template #right>
          <UButton 
            :icon="activeTab === 'pages' ? 'i-lucide-file-plus' : activeTab === 'blog' ? 'i-lucide-pen-square' : 'i-lucide-plus-circle'"
            :label="activeTab === 'pages' ? 'New Page' : activeTab === 'blog' ? 'New Post' : 'Add FAQ'"
            color="warning"
          />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <UTabs 
            v-model="activeTab"
            :items="tabs"
            class="py-4"
            color="warning"
          />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <div class="p-6">
        <!-- Pages Tab -->
        <UCard v-if="activeTab === 'pages'">
          <template #header>
            <h2 class="font-semibold">Website Pages</h2>
          </template>

          <div class="space-y-3">
            <div 
              v-for="page in pages" 
              :key="page.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
            >
              <div class="flex items-center gap-4">
                <UIcon name="i-lucide-file-text" class="size-5 text-muted" />
                <div>
                  <p class="font-medium">{{ page.title }}</p>
                  <code class="text-md bg-muted px-2 py-0.5 rounded">{{ page.slug }}</code>
                </div>
              </div>
              <div class="flex items-center gap-4">
                <span class="text-md text-muted">{{ page.lastUpdated }}</span>
                <UBadge 
                  :label="page.status"
                  :color="page.status === 'published' ? 'success' : 'neutral'"
                  variant="subtle"
                />
                <div class="flex gap-1">
                  <UButton icon="i-lucide-eye" color="neutral" variant="ghost" size="md" />
                  <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="md" />
                </div>
              </div>
            </div>
          </div>
        </UCard>

        <!-- Blog Posts Tab -->
        <UCard v-if="activeTab === 'blog'">
          <template #header>
            <h2 class="font-semibold">Blog Posts</h2>
          </template>

          <div class="space-y-3">
            <div 
              v-for="post in blogPosts" 
              :key="post.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
            >
              <div class="flex-1">
                <p class="font-medium">{{ post.title }}</p>
                <div class="flex items-center gap-4 mt-1 text-md text-muted">
                  <span>By {{ post.author }}</span>
                  <span>{{ post.date }}</span>
                  <span class="flex items-center gap-1">
                    <UIcon name="i-lucide-eye" class="size-4" />
                    {{ post.views }}
                  </span>
                </div>
              </div>
              <div class="flex items-center gap-4">
                <UBadge 
                  :label="post.status"
                  :color="post.status === 'published' ? 'success' : 'neutral'"
                  variant="subtle"
                />
                <UDropdownMenu
                  :items="[
                    [
                      { label: 'Edit', icon: 'i-lucide-pencil' },
                      { label: 'Preview', icon: 'i-lucide-eye' }
                    ],
                    [
                      { label: 'Delete', icon: 'i-lucide-trash', color: 'error' }
                    ]
                  ]"
                >
                  <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" size="md" />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </UCard>

        <!-- FAQs Tab -->
        <UCard v-if="activeTab === 'faq'">
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">Frequently Asked Questions</h2>
              <p class="text-md text-muted">Drag to reorder</p>
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="faq in faqs" 
              :key="faq.id"
              class="flex items-start gap-4 p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
            >
              <UIcon name="i-lucide-grip-vertical" class="size-5 text-muted cursor-move mt-0.5" />
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <UBadge :label="`#${faq.order}`" variant="subtle" size="md" />
                  <h3 class="font-medium">{{ faq.question }}</h3>
                </div>
                <p class="text-md text-muted line-clamp-2">{{ faq.answer }}</p>
              </div>
              <div class="flex gap-1">
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="md" />
                <UButton icon="i-lucide-trash" color="error" variant="ghost" size="md" />
              </div>
            </div>
          </div>

          <template #footer>
            <UButton label="Add FAQ" icon="i-lucide-plus" variant="outline" color="neutral" />
          </template>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
