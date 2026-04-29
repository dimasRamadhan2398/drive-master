<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js'
import { ref } from 'vue'

definePageMeta({ layout: 'admin' })

const toast = useToast()
const activeTab = ref('pages')

// ==================== MODAL STATES ====================
const showPageModal = ref(false)
const showPostModal = ref(false)
const showFaqModal = ref(false)
const isEditing = ref(false)

// ==================== FORM DATA ====================
const pageForm = ref({ id: 0, title: '', slug: '', status: 'draft' })
const postForm = ref({ id: 0, title: '', author: 'Admin', content: '', status: 'draft', media: [] as { name: string, type: string, size: string, url: string, fileType: 'image' | 'video' }[] })
const faqForm = ref({ id: 0, question: '', answer: '', order: 0 })

const { pages, blogPosts, faqs } = useContent()
const tabs = [
  { label: 'Pages', value: 'pages', icon: 'i-lucide-file-text' },
  { label: 'Blog Posts', value: 'blog', icon: 'i-lucide-newspaper' },
  { label: 'FAQs', value: 'faq', icon: 'i-lucide-help-circle' }
]

// ==================== HELPER ====================
function todayFormatted() {
  return new Date().toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function generateSlug(title: string) {
  return '/' + title.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')
}

// ==================== PAGES CRUD ====================
const currentEditingPage = ref<any>(null)

function openPageBuilder(page: any) {
  currentEditingPage.value = page
}

function savePageBuilder(updatedPage: any) {
  const idx = pages.value.findIndex(p => p.id === updatedPage.id)
  if (idx !== -1) {
    pages.value[idx] = { ...updatedPage, lastUpdated: todayFormatted() }
  }
  currentEditingPage.value = null
}

function openNewPage() {
  isEditing.value = false
  pageForm.value = { id: 0, title: '', slug: '', status: 'draft' }
  showPageModal.value = true
}

function openEditPage(page: any) {
  isEditing.value = true
  pageForm.value = { id: page.id, title: page.title, slug: page.slug, status: page.status }
  showPageModal.value = true
}

async function savePage() {
  if (!pageForm.value.title.trim()) {
    toast.add({ title: 'Error', description: 'Page title is required.', color: 'error' })
    return
  }

  const isEditingMode = isEditing.value
  // Konsep: endpoint API akan berbeda untuk membuat (POST) dan mengedit (PUT)
  const endpoint = isEditingMode ? `/api/pages/${pageForm.value.id}` : '/api/pages'
  const method = isEditingMode ? 'PUT' : 'POST'

  try {
    // Konsep: data ini akan dikirim ke backend
    const pageData = {
      title: pageForm.value.title,
      slug: pageForm.value.slug || generateSlug(pageForm.value.title),
      status: pageForm.value.status,
    }

    // Konsep: ganti bagian ini dengan panggilan API sesungguhnya, contohnya:
    // const savedPage = await $fetch(endpoint, { method, body: pageData })
    
    // Untuk demonstrasi, kita akan simulasikan respons dari backend
    const savedPage = { ...pageData, id: isEditingMode ? pageForm.value.id : Math.max(...pages.value.map(p => p.id), 0) + 1 }

    // Setelah berhasil, perbarui state lokal
    if (isEditingMode) {
      const idx = pages.value.findIndex(p => p.id === savedPage.id)
      if (idx !== -1) pages.value[idx] = { ...pages.value[idx], ...savedPage, lastUpdated: todayFormatted() }
    } else {
      // Menambahkan properti `sections` agar sesuai dengan struktur data yang ada
      pages.value.push({ ...savedPage, lastUpdated: todayFormatted(), sections: [] })
    }

    toast.add({ title: isEditingMode ? 'Page Updated' : 'Page Created', description: `"${savedPage.title}" has been saved.`, color: 'success' })
    showPageModal.value = false
  } catch (error) {
    console.error('Failed to save page:', error)
    toast.add({ title: 'Save Failed', description: 'Could not save the page. Please try again.', color: 'error' })
  }
}

function previewPage(page: any) {
  toast.add({ title: 'Preview', description: `Opening preview for "${page.title}" (${page.slug})...`, color: 'info' })
  if (!page.slug || page.slug === '/') {
    return window.open('/', 'noopener,noreferrer');
  }

  const targetUrl = page.slug.startsWith('/')
    ? page.slug
    : `/${page.slug}`;

  window.open(targetUrl, '_blank', 'noopener,noreferrer');
}

function deletePage(page: any) {
  pages.value = pages.value.filter(p => p.id !== page.id)
  toast.add({ title: 'Page Deleted', description: `"${page.title}" has been removed.`, color: 'error' })
}

function togglePageStatus(page: any) {
  const p = pages.value.find(pg => pg.id === page.id)
  if (p) {
    p.status = p.status === 'published' ? 'draft' : 'published'
    p.lastUpdated = todayFormatted()
    toast.add({ title: 'Status Updated', description: `"${p.title}" is now ${p.status}.`, color: 'success' })
  }
}

function getPageActions(page: any) {
  return [
    [
      { label: 'Edit Content', icon: 'i-lucide-layout-template', onSelect: () => openPageBuilder(page) },
      { label: 'Edit Settings', icon: 'i-lucide-settings', onSelect: () => openEditPage(page) },
      { label: 'Preview', icon: 'i-lucide-eye', onSelect: () => previewPage(page) },
      { label: page.status === 'published' ? 'Unpublish' : 'Publish', icon: page.status === 'published' ? 'i-lucide-eye-off' : 'i-lucide-globe', onSelect: () => togglePageStatus(page) }
    ],
    [
      { label: 'Delete', icon: 'i-lucide-trash', color: 'error' as const, onSelect: () => deletePage(page) }
    ]
  ]
}

// ==================== BLOG POSTS CRUD ====================
function openNewPost() {
  isEditing.value = false
  postForm.value = { id: 0, title: '', author: 'Admin', content: '', status: 'draft', media: [] }
  showPostModal.value = true
}

function openEditPost(post: any) {
  isEditing.value = true
  postForm.value = { id: post.id, title: post.title, author: post.author, content: post.content || '', status: post.status, media: post.media ? [...post.media] : [] }
  showPostModal.value = true
}

function savePost() {
  if (!postForm.value.title.trim()) {
    toast.add({ title: 'Error', description: 'Post title is required.', color: 'error' })
    return
  }
  if (isEditing.value) {
    const idx = blogPosts.value.findIndex(p => p.id === postForm.value.id)
    if (idx !== -1) {
      blogPosts.value[idx] = { ...blogPosts.value[idx], title: postForm.value.title, author: postForm.value.author, content: postForm.value.content, status: postForm.value.status, media: [...postForm.value.media], date: todayFormatted() }
      toast.add({ title: 'Post Updated', description: `"${postForm.value.title}" has been updated.`, color: 'success' })
    }
  } else {
    const newId = Math.max(...blogPosts.value.map(p => p.id), 0) + 1
    blogPosts.value.push({
      id: newId,
      title: postForm.value.title,
      author: postForm.value.author,
      content: postForm.value.content,
      date: todayFormatted(),
      status: postForm.value.status,
      views: 0,
      media: [...postForm.value.media]
    })
    toast.add({ title: 'Post Created', description: `"${postForm.value.title}" has been created.`, color: 'success' })
  }
  showPostModal.value = false
}

// ==================== MEDIA UPLOAD ====================
const ACCEPTED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp']
const ACCEPTED_VIDEO_TYPES = ['video/mp4', 'video/webm']
const MAX_IMAGE_SIZE = 5 * 1024 * 1024 // 5MB
const MAX_VIDEO_SIZE = 50 * 1024 * 1024 // 50MB
const mediaInputRef = ref<HTMLInputElement | null>(null)

function triggerMediaUpload() {
  mediaInputRef.value?.click()
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function handleMediaSelect(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files) return

  for (const file of Array.from(input.files)) {
    const isImage = ACCEPTED_IMAGE_TYPES.includes(file.type)
    const isVideo = ACCEPTED_VIDEO_TYPES.includes(file.type)

    if (!isImage && !isVideo) {
      toast.add({ title: 'Invalid File', description: `"${file.name}" is not a supported format. Use JPG, PNG, WebP for images or MP4, WebM for videos.`, color: 'error' })
      continue
    }

    if (isImage && file.size > MAX_IMAGE_SIZE) {
      toast.add({ title: 'File Too Large', description: `"${file.name}" exceeds the 5MB limit for images.`, color: 'error' })
      continue
    }

    if (isVideo && file.size > MAX_VIDEO_SIZE) {
      toast.add({ title: 'File Too Large', description: `"${file.name}" exceeds the 50MB limit for videos.`, color: 'error' })
      continue
    }

    const reader = new FileReader()
    reader.onload = (e) => {
      postForm.value.media.push({
        name: file.name,
        type: file.type,
        size: formatFileSize(file.size),
        url: e.target?.result as string,
        fileType: isImage ? 'image' : 'video'
      })
    }
    reader.readAsDataURL(file)
  }

  // Reset input so same file can be re-selected
  input.value = ''
}

function removeMedia(index: number) {
  postForm.value.media.splice(index, 1)
}

function previewPost(post: any) {
  toast.add({ title: 'Preview', description: `Opening preview for "${post.title}"...`, color: 'info' })
  if (!post.slug || post.slug === '/') {
    return window.open('/', 'noopener,noreferrer');
  }

  const targetUrl = post.slug.startsWith('/')
    ? post.slug
    : `/${post.slug}`;

  window.open(targetUrl, '_blank', 'noopener,noreferrer');
}

function deletePost(post: any) {
  blogPosts.value = blogPosts.value.filter(p => p.id !== post.id)
  toast.add({ title: 'Post Deleted', description: `"${post.title}" has been removed.`, color: 'error' })
}

function togglePostStatus(post: any) {
  const p = blogPosts.value.find(bp => bp.id === post.id)
  if (p) {
    p.status = p.status === 'published' ? 'draft' : 'published'
    p.date = todayFormatted()
    toast.add({ title: 'Status Updated', description: `"${p.title}" is now ${p.status}.`, color: 'success' })
  }
}

function getPostActions(post: any) {
  return [
    [
      { label: 'Edit', icon: 'i-lucide-pencil', onSelect: () => openEditPost(post) },
      { label: 'Preview', icon: 'i-lucide-eye', onSelect: () => previewPost(post) },
      { label: post.status === 'published' ? 'Unpublish' : 'Publish', icon: post.status === 'published' ? 'i-lucide-eye-off' : 'i-lucide-globe', onSelect: () => togglePostStatus(post) }
    ],
    [
      { label: 'Delete', icon: 'i-lucide-trash', color: 'error' as const, onSelect: () => deletePost(post) }
    ]
  ]
}

// ==================== FAQ CRUD ====================
function openNewFaq() {
  isEditing.value = false
  faqForm.value = { id: 0, question: '', answer: '', order: 0 }
  showFaqModal.value = true
}

function openEditFaq(faq: any) {
  isEditing.value = true
  faqForm.value = { id: faq.id, question: faq.question, answer: faq.answer, order: faq.order }
  showFaqModal.value = true
}

function saveFaq() {
  if (!faqForm.value.question.trim() || !faqForm.value.answer.trim()) {
    toast.add({ title: 'Error', description: 'Both question and answer are required.', color: 'error' })
    return
  }
  if (isEditing.value) {
    const idx = faqs.value.findIndex(f => f.id === faqForm.value.id)
    if (idx !== -1) {
      faqs.value[idx] = { ...faqForm.value }
      toast.add({ title: 'FAQ Updated', description: 'FAQ has been updated.', color: 'success' })
    }
  } else {
    const newId = Math.max(...faqs.value.map(f => f.id), 0) + 1
    faqs.value.push({
      id: newId,
      question: faqForm.value.question,
      answer: faqForm.value.answer,
      order: faqs.value.length + 1
    })
    toast.add({ title: 'FAQ Created', description: 'New FAQ has been added.', color: 'success' })
  }
  showFaqModal.value = false
}

function deleteFaq(faq: any) {
  faqs.value = faqs.value.filter(f => f.id !== faq.id)
  faqs.value.forEach((f, i) => f.order = i + 1)
  toast.add({ title: 'FAQ Deleted', description: 'FAQ has been removed.', color: 'error' })
}

// ==================== FAQ DRAG & DROP ====================
const dragIndex = ref<number | null>(null)

function onDragStart(index: number, event: DragEvent) {
  dragIndex.value = index
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
  }
}

function onDragOver(event: DragEvent) {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move'
  }
}

function onDrop(targetIndex: number) {
  if (dragIndex.value === null || dragIndex.value === targetIndex) return
  const item = faqs.value.splice(dragIndex.value, 1)[0]
  faqs.value.splice(targetIndex, 0, item)
  faqs.value.forEach((f, i) => f.order = i + 1)
  dragIndex.value = null
  toast.add({ title: 'Reordered', description: 'FAQ order has been updated.', color: 'success' })
}

function onDragEnd() {
  dragIndex.value = null
}

// ==================== HEADER BUTTON HANDLER ====================
function handleHeaderAction() {
  if (activeTab.value === 'pages') openNewPage()
  else if (activeTab.value === 'blog') openNewPost()
  else openNewFaq()
}
</script>

<template>
  <!-- PAGE BUILDER VIEW -->
  <UDashboardPanel v-if="currentEditingPage">
    <div class="h-full w-full overflow-y-auto p-6 bg-background">
      <AdminPageEditor 
        :page="currentEditingPage" 
        @close="currentEditingPage = null" 
        @save="savePageBuilder" 
      />
    </div>
  </UDashboardPanel>

  <!-- MAIN DASHBOARD -->
  <UDashboardPanel v-else>
      <template #header>
      <UDashboardNavbar title="Content Management">
        <template #right>
          <UButton 
            :icon="activeTab === 'pages' ? 'i-lucide-file-plus' : activeTab === 'blog' ? 'i-lucide-pen-square' : 'i-lucide-plus-circle'"
            :label="activeTab === 'pages' ? 'New Page' : activeTab === 'blog' ? 'New Post' : 'Add FAQ'"
            color="warning"
            @click="handleHeaderAction"
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
        <!-- ==================== PAGES TAB ==================== -->
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
                <UDropdownMenu :items="getPageActions(page)">
                  <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" size="sm" />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </UCard>

        <!-- ==================== BLOG POSTS TAB ==================== -->
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
                <UDropdownMenu :items="getPostActions(post)">
                  <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" size="sm" />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </UCard>

        <!-- ==================== FAQS TAB ==================== -->
        <UCard v-if="activeTab === 'faq'">
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">Frequently Asked Questions</h2>
              <p class="text-md text-muted">Drag to reorder</p>
            </div>
          </template>

          <div class="space-y-3">
            <div 
              v-for="(faq, index) in faqs" 
              :key="faq.id"
              draggable="true"
              class="flex items-start gap-4 p-4 rounded-lg border border-default hover:bg-muted/30 transition-colors"
              :class="{ 'opacity-50 border-dashed border-primary': dragIndex === index }"
              @dragstart="onDragStart(index, $event)"
              @dragover="onDragOver"
              @drop="onDrop(index)"
              @dragend="onDragEnd"
            >
              <UIcon name="i-lucide-grip-vertical" class="size-5 text-muted cursor-grab mt-0.5 shrink-0" />
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <UBadge :label="`#${faq.order}`" variant="subtle" size="md" />
                  <h3 class="font-medium">{{ faq.question }}</h3>
                </div>
                <p class="text-md text-muted line-clamp-2">{{ faq.answer }}</p>
              </div>
              <div class="flex gap-1 shrink-0">
                <UButton icon="i-lucide-pencil" color="neutral" variant="ghost" size="sm" @click="openEditFaq(faq)" />
                <UButton icon="i-lucide-trash" color="error" variant="ghost" size="sm" @click="deleteFaq(faq)" />
              </div>
            </div>
          </div>

        </UCard>
      </div>
    </template>
  </UDashboardPanel>

  <!-- ==================== PAGE MODAL ==================== -->
  <ClientOnly>
    <UModal v-model:open="showPageModal">
      <template #content>
        <div class="bg-default rounded-2xl w-full">
          <div class="px-6 py-4 border-b border-default flex items-center justify-between">
            <h3 class="text-base font-semibold">{{ isEditing ? 'Edit Page' : 'New Page' }}</h3>
            <UButton icon="i-lucide-x" color="neutral" variant="ghost" @click="showPageModal = false" />
          </div>
          <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
            <div>
              <label class="block text-sm font-medium mb-1.5">Page Title *</label>
              <UInput v-model="pageForm.title" placeholder="e.g. Contact Us" class="w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Slug</label>
              <UInput v-model="pageForm.slug" :placeholder="pageForm.title ? generateSlug(pageForm.title) : '/contact-us'" class="w-full" />
              <p class="text-xs text-muted mt-1">Leave empty to auto-generate from title</p>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Status</label>
              <USelect v-model="pageForm.status" :items="['draft', 'published']" class="w-full" />
            </div>
          </div>
          <div class="px-6 py-4 border-t border-default flex justify-end gap-3">
            <UButton :label="isEditing ? 'Save Changes' : 'Create Page'" icon="i-lucide-check" @click="savePage" />
          </div>
        </div>
      </template>
    </UModal>
  </ClientOnly>

  <!-- ==================== POST MODAL ==================== -->
  <ClientOnly>
    <UModal v-model:open="showPostModal">
      <template #content>
        <div class="bg-default rounded-2xl w-full">
          <div class="px-6 py-4 border-b border-default flex items-center justify-between">
            <h3 class="text-base font-semibold">{{ isEditing ? 'Edit Post' : 'New Post' }}</h3>
            <UButton icon="i-lucide-x" color="neutral" variant="ghost" @click="showPostModal = false" />
          </div>
          <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
            <div>
              <label class="block text-sm font-medium mb-1.5">Post Title *</label>
              <UInput v-model="postForm.title" placeholder="e.g. How to Charge Your EV at Home" class="w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Author</label>
              <UInput v-model="postForm.author" placeholder="Admin" class="w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Content</label>
              <UTextarea v-model="postForm.content" placeholder="Write your blog post content here..." :rows="6" class="w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Media</label>
              <input 
                ref="mediaInputRef"
                type="file" 
                accept="image/jpeg,image/png,image/webp,video/mp4,video/webm"
                multiple
                class="hidden"
                @change="handleMediaSelect"
              />
              
              <!-- Media Preview Grid -->
              <div v-if="postForm.media.length > 0" class="grid grid-cols-3 gap-3 mb-3">
                <div 
                  v-for="(item, idx) in postForm.media" 
                  :key="idx"
                  class="relative group rounded-lg border border-default overflow-hidden bg-muted/30"
                >
                  <img 
                    v-if="item.fileType === 'image'" 
                    :src="item.url" 
                    :alt="item.name"
                    class="w-full h-24 object-cover"
                  />
                  <div 
                    v-else 
                    class="w-full h-24 flex flex-col items-center justify-center gap-1"
                  >
                    <UIcon name="i-lucide-film" class="size-8 text-muted" />
                    <span class="text-[10px] text-muted truncate max-w-full px-2">{{ item.name }}</span>
                  </div>
                  <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <UButton 
                      icon="i-lucide-trash-2" 
                      color="error" 
                      variant="ghost" 
                      size="xs" 
                      @click="removeMedia(idx)"
                      class="text-white"
                    />
                  </div>
                  <div class="absolute bottom-0 left-0 right-0 bg-black/60 px-2 py-0.5 flex items-center justify-between">
                    <span class="text-[10px] text-white truncate">{{ item.name }}</span>
                    <span class="text-[10px] text-white/70 shrink-0 ml-1">{{ item.size }}</span>
                  </div>
                </div>
              </div>

              <!-- Upload Button -->
              <button 
                @click="triggerMediaUpload"
                class="w-full border-2 border-dashed border-default rounded-lg p-4 hover:border-primary hover:bg-primary/5 transition-colors cursor-pointer flex flex-col items-center gap-2"
              >
                <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
                <span class="text-sm font-medium text-muted">Click to add images or videos</span>
                <span class="text-xs text-muted/60">JPG, PNG, WebP (max 5MB) · MP4, WebM (max 50MB)</span>
              </button>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Status</label>
              <USelect v-model="postForm.status" :items="['draft', 'published']" class="w-full" />
            </div>
          </div>
          <div class="px-6 py-4 border-t border-default flex justify-end gap-3">
            <UButton :label="isEditing ? 'Save Changes' : 'Create Post'" icon="i-lucide-check" @click="savePost" />
          </div>
        </div>
      </template>
    </UModal>
  </ClientOnly>

  <!-- ==================== FAQ MODAL ==================== -->
  <ClientOnly>
    <UModal v-model:open="showFaqModal">
      <template #content>
        <div class="bg-default rounded-2xl w-full">
          <div class="px-6 py-4 border-b border-default flex items-center justify-between">
            <h3 class="text-base font-semibold">{{ isEditing ? 'Edit FAQ' : 'Add New FAQ' }}</h3>
            <UButton icon="i-lucide-x" color="neutral" variant="ghost" @click="showFaqModal = false" />
          </div>
          <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
            <div>
              <label class="block text-sm font-medium mb-1.5">Question *</label>
              <UInput v-model="faqForm.question" placeholder="e.g. What payment methods do you accept?" class="w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Answer *</label>
              <UTextarea v-model="faqForm.answer" placeholder="Type the answer here..." :rows="4" class="w-full" />
            </div>
          </div>
          <div class="px-6 py-4 border-t border-default flex justify-end gap-3">
            <UButton :label="isEditing ? 'Save Changes' : 'Add FAQ'" icon="i-lucide-check" @click="saveFaq" />
          </div>
        </div>
      </template>
    </UModal>
  </ClientOnly>
</template>
