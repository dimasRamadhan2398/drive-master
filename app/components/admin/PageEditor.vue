<script setup lang="ts">
const props = defineProps({
  page: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'save'])
const toast = useToast()

// Working copy of the page data to avoid mutating props directly until save
const formData = ref(JSON.parse(JSON.stringify(props.page)))
if (!formData.value.sections) {
  formData.value.sections = []
}

// Generate unique ID for new sections
const generateId = () => Math.random().toString(36).substr(2, 9)

const sectionTypes = [
  { label: 'Hero Section', value: 'hero', icon: 'i-lucide-layout-template' },
  { label: 'Text Block', value: 'text', icon: 'i-lucide-align-left' },
  { label: 'Image + Text', value: 'image_text', icon: 'i-lucide-image' },
  { label: 'Call to Action (CTA)', value: 'cta', icon: 'i-lucide-megaphone' }
]

function addSection(type: string) {
  let defaultData = {}
  
  if (type === 'hero') defaultData = { heading: '', subheading: '', ctaText: '', bgImage: '' }
  else if (type === 'text') defaultData = { content: '' }
  else if (type === 'image_text') defaultData = { image: '', content: '' }
  else if (type === 'cta') defaultData = { heading: '', buttonText: '', buttonLink: '' }

  formData.value.sections.push({
    id: generateId(),
    type,
    data: defaultData
  })
}

function removeSection(index: number) {
  formData.value.sections.splice(index, 1)
}

function getSectionTitle(type: string) {
  return sectionTypes.find(t => t.value === type)?.label || 'Unknown Section'
}

function getSectionIcon(type: string) {
  return sectionTypes.find(t => t.value === type)?.icon || 'i-lucide-box'
}

// ==================== DRAG & DROP LOGIC ====================
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
  const item = formData.value.sections.splice(dragIndex.value, 1)[0]
  formData.value.sections.splice(targetIndex, 0, item)
  dragIndex.value = null
}

function onDragEnd() {
  dragIndex.value = null
}

// ==================== ACTIONS ====================
function handleSave() {
  emit('save', formData.value)
  toast.add({ title: 'Page Saved', description: `Sections for "${formData.value.title}" have been saved.`, color: 'success' })
}

function handleClose() {
  emit('close')
}
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center justify-between pb-6 mb-6 border-b border-default">
      <div>
        <div class="flex items-center gap-3 mb-1">
          <UButton icon="i-lucide-arrow-left" color="neutral" variant="ghost" @click="handleClose" />
          <h2 class="text-xl font-semibold">Editing: {{ formData.title }}</h2>
          <UBadge :label="formData.status" :color="formData.status === 'published' ? 'success' : 'warning'" variant="subtle" />
        </div>
        <p class="text-sm text-muted ml-11">Path: <code>{{ formData.slug }}</code></p>
      </div>
      <div class="flex items-center gap-3">
        <UButton label="Discard Changes" color="neutral" variant="ghost" @click="handleClose" />
        <UButton label="Save Page" icon="i-lucide-save" @click="handleSave" />
      </div>
    </div>

    <!-- Content Area -->
    <div class="max-w-4xl mx-auto w-full space-y-6 pb-20">
      
      <!-- Empty State -->
      <div v-if="formData.sections.length === 0" class="text-center py-16 border-2 border-dashed border-default rounded-xl">
        <UIcon name="i-lucide-layout-dashboard" class="size-12 text-muted mb-3 mx-auto" />
        <h3 class="text-lg font-medium mb-1">No sections yet</h3>
        <p class="text-muted text-sm mb-4">Start building your page by adding a section below.</p>
      </div>

      <!-- Sections List -->
      <div class="space-y-4">
        <div 
          v-for="(section, index) in formData.sections" 
          :key="section.id"
          draggable="true"
          class="bg-default rounded-xl border border-default shadow-sm overflow-hidden transition-opacity"
          :class="{ 'opacity-40 border-primary border-dashed': dragIndex === index }"
          @dragstart="onDragStart(index, $event)"
          @dragover="onDragOver"
          @drop="onDrop(index)"
          @dragend="onDragEnd"
        >
          <!-- Section Header -->
          <div class="bg-muted/30 px-4 py-3 border-b border-default flex items-center justify-between group">
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-grip-horizontal" class="size-5 text-muted cursor-grab active:cursor-grabbing hover:text-foreground transition-colors" />
              <UIcon :name="getSectionIcon(section.type)" class="size-4 text-primary" />
              <h4 class="font-medium text-sm">{{ getSectionTitle(section.type) }}</h4>
            </div>
            <UButton icon="i-lucide-trash-2" color="error" variant="ghost" size="xs" class="opacity-0 group-hover:opacity-100 transition-opacity" @click="removeSection(index)" />
          </div>

          <!-- Section Body (Dynamic Forms) -->
          <div class="p-5">
            <!-- HERO FORM -->
            <div v-if="section.type === 'hero'" class="grid grid-cols-2 gap-4">
              <div class="col-span-2">
                <label class="block text-xs font-medium text-muted mb-1.5">Heading</label>
                <UInput v-model="section.data.heading" placeholder="Main big title" class="w-full" />
              </div>
              <div class="col-span-2">
                <label class="block text-xs font-medium text-muted mb-1.5">Subheading</label>
                <UTextarea v-model="section.data.subheading" placeholder="Description under the title" :rows="2" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">CTA Button Text</label>
                <UInput v-model="section.data.ctaText" placeholder="e.g. Get Started" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Background Image URL</label>
                <UInput v-model="section.data.bgImage" icon="i-lucide-image" placeholder="https://..." class="w-full" />
              </div>
            </div>

            <!-- TEXT BLOCK FORM -->
            <div v-if="section.type === 'text'">
              <label class="block text-xs font-medium text-muted mb-1.5">Content</label>
              <UTextarea v-model="section.data.content" placeholder="Write your paragraph here..." :rows="4" class="w-full" />
            </div>

            <!-- IMAGE + TEXT FORM -->
            <div v-if="section.type === 'image_text'" class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Image URL</label>
                <div class="border-2 border-dashed border-default rounded-lg p-4 text-center h-32 flex flex-col items-center justify-center relative overflow-hidden">
                  <template v-if="section.data.image">
                    <img :src="section.data.image" class="absolute inset-0 w-full h-full object-cover opacity-30" />
                    <UInput v-model="section.data.image" class="relative z-10 w-[90%]" size="sm" placeholder="URL..." />
                  </template>
                  <template v-else>
                    <UIcon name="i-lucide-image-plus" class="size-6 text-muted mb-2" />
                    <UInput v-model="section.data.image" class="w-full" size="sm" placeholder="Paste image URL here" />
                  </template>
                </div>
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Text Content</label>
                <UTextarea v-model="section.data.content" placeholder="Description text..." class="w-full h-32" />
              </div>
            </div>

            <!-- CTA FORM -->
            <div v-if="section.type === 'cta'" class="grid grid-cols-2 gap-4 bg-primary/5 p-4 rounded-lg border border-primary/20">
              <div class="col-span-2">
                <label class="block text-xs font-medium text-primary mb-1.5">CTA Heading</label>
                <UInput v-model="section.data.heading" placeholder="e.g. Ready to start driving?" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-primary mb-1.5">Button Text</label>
                <UInput v-model="section.data.buttonText" placeholder="e.g. Contact Us" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-primary mb-1.5">Button Link</label>
                <UInput v-model="section.data.buttonLink" placeholder="e.g. /contact" class="w-full" />
              </div>
            </div>

          </div>
        </div>
      </div>

      <!-- Add Section Button -->
      <div class="flex justify-center mt-8">
        <UDropdownMenu
          :items="[
            sectionTypes.map(t => ({ 
              label: t.label, 
              icon: t.icon, 
              onSelect: () => addSection(t.value) 
            }))
          ]"
        >
          <UButton label="Add Section" icon="i-lucide-plus" color="primary" variant="soft" size="lg" class="shadow-sm" />
        </UDropdownMenu>
      </div>
      
    </div>
  </div>
</template>
