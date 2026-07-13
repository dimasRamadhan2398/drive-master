<script setup lang="ts">
const { t } = useI18n()
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
  { label: t('admin.heroSection'), value: 'hero', icon: 'i-lucide-layout-template' },
  { label: t('admin.textBlock'), value: 'text', icon: 'i-lucide-align-left' },
  { label: t('admin.imageText'), value: 'image_text', icon: 'i-lucide-image' },
  { label: t('admin.ctaSection'), value: 'cta', icon: 'i-lucide-megaphone' },
  { label: 'Course Material Grid', value: 'course_material', icon: 'i-lucide-book-open' }
]

function addSection(type: string) {
  let defaultData = {}
  
  if (type === 'hero') {
    defaultData = {
      heading: '',
      subheading: '',
      ctaText: '',
      ctaLink: '',
      secondaryCtaText: '',
      secondaryCtaLink: '',
      bgImage: '',
      features: []
    }
  }
  else if (type === 'text') defaultData = { content: '' }
  else if (type === 'image_text') defaultData = { image: '', content: '' }
  else if (type === 'cta') defaultData = { heading: '', buttonText: '', buttonLink: '' }
  else if (type === 'course_material') {
    defaultData = {
      headline: '',
      title: '',
      description: '',
      materials: []
    }
  }

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

function addFeature(section: any) {
  if (!section.data.features) section.data.features = []
  section.data.features.push({ title: '', icon: 'i-lucide-check-circle' })
}

function removeFeature(section: any, index: number) {
  section.data.features.splice(index, 1)
}

function addMaterial(section: any) {
  if (!section.data.materials) section.data.materials = []
  section.data.materials.push({ title: '', icon: 'i-lucide-book-open', description: [] })
}

function removeMaterial(section: any, index: number) {
  section.data.materials.splice(index, 1)
}

function getBulletsText(description: string[] | undefined) {
  if (!description) return ''
  return description.join('\n')
}

function setBulletsText(material: any, text: string) {
  material.description = text.split('\n').map(line => line.trim()).filter(line => line.length > 0)
}

// ==================== IMAGE UPLOADS ====================
function triggerFileUpload(sectionId: string) {
  const input = document.getElementById(`file-input-${sectionId}`) as HTMLInputElement
  input?.click()
}

function handleFileChange(event: Event, section: any) {
  const input = event.target as HTMLInputElement
  if (!input.files || !input.files[0]) return

  const file = input.files[0]
  if (file.size > 5 * 1024 * 1024) {
    toast.add({
      title: "File Too Large",
      description: `"${file.name}" exceeds the 5 MB limit for images.`,
      color: "error",
    })
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    section.data.bgImage = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

function clearBgImage(section: any) {
  section.data.bgImage = ''
}

function triggerImageUpload(sectionId: string) {
  const input = document.getElementById(`image-input-${sectionId}`) as HTMLInputElement
  input?.click()
}

function handleImageFileChange(event: Event, section: any) {
  const input = event.target as HTMLInputElement
  if (!input.files || !input.files[0]) return

  const file = input.files[0]
  if (file.size > 5 * 1024 * 1024) {
    toast.add({
      title: "File Too Large",
      description: `"${file.name}" exceeds the 5 MB limit for images.`,
      color: "error",
    })
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    section.data.image = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

function clearImage(section: any) {
  section.data.image = ''
}

// ==================== ACTIONS ====================
function handleSave() {
  emit('save', formData.value)
  toast.add({ title: t('admin.pageSaved'), description: t('admin.pageSavedDesc', { title: formData.value.title }), color: 'success' })
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
          <h2 class="text-xl font-semibold">{{ t('admin.editPage') }}: {{ formData.title }}</h2>
          <UBadge :label="formData.status" :color="formData.status === 'published' ? 'success' : 'warning'" variant="subtle" />
        </div>
        <p class="text-sm text-muted ml-11">Path: <code>{{ formData.slug }}</code></p>
      </div>
      <div class="flex items-center gap-3">
        <UButton :label="t('admin.discardChanges')" color="neutral" variant="ghost" @click="handleClose" />
        <UButton :label="t('admin.savePage')" icon="i-lucide-save" @click="handleSave" />
      </div>
    </div>

    <!-- Content Area -->
    <div class="max-w-4xl mx-auto w-full space-y-6 pb-20">
      
      <!-- Empty State -->
      <div v-if="formData.sections.length === 0" class="text-center py-16 border-2 border-dashed border-default rounded-xl">
        <UIcon name="i-lucide-layout-dashboard" class="size-12 text-muted mb-3 mx-auto" />
        <h3 class="text-lg font-medium mb-1">{{ t('admin.noSections') }}</h3>
        <p class="text-muted text-sm mb-4">{{ t('admin.addSectionDesc') }}</p>
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
                <label class="block text-xs font-medium text-muted mb-1.5">{{ t('admin.heading') }}</label>
                <UInput v-model="section.data.heading" placeholder="Main big title" class="w-full" />
              </div>
              <div class="col-span-2">
                <label class="block text-xs font-medium text-muted mb-1.5">{{ t('admin.subheading') }}</label>
                <UTextarea v-model="section.data.subheading" placeholder="Description under the title" :rows="2" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Primary CTA Button Text</label>
                <UInput v-model="section.data.ctaText" placeholder="e.g. Pesan Sesi Pertama" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Primary CTA Link</label>
                <UInput v-model="section.data.ctaLink" placeholder="e.g. /auth/register" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Secondary CTA Button Text</label>
                <UInput v-model="section.data.secondaryCtaText" placeholder="e.g. Lihat Paket" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Secondary CTA Link</label>
                <UInput v-model="section.data.secondaryCtaLink" placeholder="e.g. /packages" class="w-full" />
              </div>
              <div class="col-span-2">
                <label class="block text-xs font-medium text-muted mb-1.5">{{ t('admin.bgImageUrl') }}</label>
                <div class="flex flex-col gap-2">
                  <div v-if="section.data.bgImage" class="relative h-40 rounded-lg overflow-hidden border border-default group">
                    <img :src="section.data.bgImage" class="w-full h-full object-cover" />
                    <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                      <UButton
                        label="Change Image"
                        icon="i-lucide-upload"
                        size="sm"
                        color="neutral"
                        @click="triggerFileUpload(section.id)"
                      />
                      <UButton
                        label="Remove"
                        icon="i-lucide-trash"
                        size="sm"
                        color="error"
                        @click="clearBgImage(section)"
                      />
                    </div>
                  </div>
                  
                  <div class="flex gap-2">
                    <UInput
                      v-model="section.data.bgImage"
                      icon="i-lucide-image"
                      placeholder="Paste image URL or upload local image..."
                      class="flex-1"
                    />
                    <UButton
                      icon="i-lucide-upload"
                      color="neutral"
                      variant="soft"
                      @click="triggerFileUpload(section.id)"
                    />
                    <input
                      :id="`file-input-${section.id}`"
                      type="file"
                      class="hidden"
                      accept="image/*"
                      @change="handleFileChange($event, section)"
                    />
                  </div>
                </div>
              </div>
              <!-- Features manager -->
              <div class="col-span-2 border-t border-default pt-4 mt-2">
                <div class="flex items-center justify-between mb-3">
                  <span class="text-xs font-bold text-muted uppercase tracking-wider">Features</span>
                  <UButton label="Add Feature" icon="i-lucide-plus" size="xs" color="neutral" variant="soft" @click="addFeature(section)" />
                </div>
                <div class="space-y-2">
                  <div v-for="(feat, fIdx) in section.data.features" :key="fIdx" class="flex gap-2 items-center">
                    <UInput v-model="feat.title" placeholder="Feature title" class="flex-1" size="sm" />
                    <UInput v-model="feat.icon" placeholder="Icon class (e.g. i-lucide-car)" class="w-48" size="sm" />
                    <UButton icon="i-lucide-trash" color="error" variant="ghost" size="xs" @click="removeFeature(section, fIdx)" />
                  </div>
                </div>
              </div>
            </div>

            <!-- COURSE MATERIAL GRID FORM -->
            <div v-if="section.type === 'course_material'" class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div class="col-span-2">
                  <label class="block text-xs font-medium text-muted mb-1.5">Headline</label>
                  <UInput v-model="section.data.headline" placeholder="e.g. Materi kursus yang akan Anda pelajari" class="w-full" />
                </div>
                <div>
                  <label class="block text-xs font-medium text-muted mb-1.5">Title</label>
                  <UInput v-model="section.data.title" placeholder="e.g. Materi Kursus" class="w-full" />
                </div>
                <div>
                  <label class="block text-xs font-medium text-muted mb-1.5">Description</label>
                  <UInput v-model="section.data.description" placeholder="Short introduction" class="w-full" />
                </div>
              </div>

              <!-- Materials Manager -->
              <div class="border-t border-default pt-4 mt-2">
                <div class="flex items-center justify-between mb-3">
                  <span class="text-xs font-bold text-muted uppercase tracking-wider">Course Materials / Cards</span>
                  <UButton label="Add Card" icon="i-lucide-plus" size="xs" color="neutral" variant="soft" @click="addMaterial(section)" />
                </div>
                <div class="space-y-4">
                  <div 
                    v-for="(mat, mIdx) in section.data.materials" 
                    :key="mIdx" 
                    class="p-4 rounded-lg border border-default bg-muted/20 relative"
                  >
                    <UButton 
                      icon="i-lucide-trash" 
                      color="error" 
                      variant="ghost" 
                      size="xs" 
                      class="absolute top-2 right-2" 
                      @click="removeMaterial(section, mIdx)" 
                    />
                    <div class="grid grid-cols-2 gap-4 mb-3">
                      <div>
                        <label class="block text-[10px] font-bold uppercase text-muted mb-1">Card Title</label>
                        <UInput v-model="mat.title" placeholder="e.g. Teori Materi" size="sm" class="w-full" />
                      </div>
                      <div>
                        <label class="block text-[10px] font-bold uppercase text-muted mb-1">Card Icon</label>
                        <UInput v-model="mat.icon" placeholder="e.g. i-lucide-book-open" size="sm" class="w-full" />
                      </div>
                    </div>
                    <div>
                      <label class="block text-[10px] font-bold uppercase text-muted mb-1">Bullet Points (One per line)</label>
                      <UTextarea 
                        :model-value="getBulletsText(mat.description)" 
                        @update:model-value="setBulletsText(mat, $event)"
                        placeholder="Bullet 1&#10;Bullet 2&#10;Bullet 3" 
                        size="sm" 
                        :rows="4" 
                        class="w-full" 
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- TEXT BLOCK FORM -->
            <div v-if="section.type === 'text'">
              <label class="block text-xs font-medium text-muted mb-1.5">{{ t('admin.content.label') }}</label>
              <UTextarea v-model="section.data.content" placeholder="Write your paragraph here..." :rows="4" class="w-full" />
            </div>

            <!-- IMAGE + TEXT FORM -->
            <div v-if="section.type === 'image_text'" class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">{{ t('admin.imageUrl') }}</label>
                <div class="flex flex-col gap-2">
                  <div v-if="section.data.image" class="relative h-40 rounded-lg overflow-hidden border border-default group">
                    <img :src="section.data.image" class="w-full h-full object-cover" />
                    <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                      <UButton
                        label="Change Image"
                        icon="i-lucide-upload"
                        size="sm"
                        color="neutral"
                        @click="triggerImageUpload(section.id)"
                      />
                      <UButton
                        label="Remove"
                        icon="i-lucide-trash"
                        size="sm"
                        color="error"
                        @click="clearImage(section)"
                      />
                    </div>
                  </div>

                  <div class="flex gap-2">
                    <UInput
                      v-model="section.data.image"
                      icon="i-lucide-image"
                      placeholder="Paste image URL or upload local image..."
                      class="flex-1"
                    />
                    <UButton
                      icon="i-lucide-upload"
                      color="neutral"
                      variant="soft"
                      @click="triggerImageUpload(section.id)"
                    />
                    <input
                      :id="`image-input-${section.id}`"
                      type="file"
                      class="hidden"
                      accept="image/*"
                      @change="handleImageFileChange($event, section)"
                    />
                  </div>
                </div>
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">{{ t('admin.textContent') }}</label>
                <UTextarea v-model="section.data.content" placeholder="Description text..." class="w-full h-40" />
              </div>
            </div>

            <!-- CTA FORM -->
            <div v-if="section.type === 'cta'" class="grid grid-cols-2 gap-4 bg-primary/5 p-4 rounded-lg border border-primary/20">
              <div class="col-span-2">
                <label class="block text-xs font-medium text-primary mb-1.5">{{ t('admin.ctaSection') }} {{ t('admin.heading') }}</label>
                <UInput v-model="section.data.heading" placeholder="e.g. Ready to start driving?" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-primary mb-1.5">{{ t('admin.buttonText') }}</label>
                <UInput v-model="section.data.buttonText" placeholder="e.g. Contact Us" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-primary mb-1.5">{{ t('admin.buttonLink') }}</label>
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
            sectionTypes.map(st => ({
              label: st.label,
              icon: st.icon,
              onSelect: () => addSection(st.value)
            }))
          ]"
        >
          <UButton :label="t('admin.addSection')" icon="i-lucide-plus" color="primary" variant="soft" size="lg" class="shadow-sm" />
        </UDropdownMenu>
      </div>
      
    </div>
  </div>
</template>
