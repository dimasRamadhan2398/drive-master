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
const formData = ref<{ sections: any[]; [key: string]: any }>(JSON.parse(JSON.stringify(props.page)))
if (!formData.value.sections) {
  formData.value.sections = []
}

// Generate unique ID for new sections
const generateId = () => Math.random().toString(36).substr(2, 9)

const { waLink } = useSettings()

const sectionTypes = [
  { label: t('admin.heroSection'), value: 'hero', icon: 'i-lucide-layout-template' },
  { label: t('admin.specificationsGrid'), value: 'specifications', icon: 'i-lucide-layout-grid' },
  { label: 'Course Material Grid', value: 'course_material', icon: 'i-lucide-book-open' },
  { label: t('admin.serviceAreas'), value: 'service_areas', icon: 'i-lucide-map-pin' },
  { label: t('admin.quoteSection'), value: 'quote', icon: 'i-lucide-quote' },
  { label: t('admin.textBlock'), value: 'text', icon: 'i-lucide-align-left' },
  { label: t('admin.imageText'), value: 'image_text', icon: 'i-lucide-image' },
  { label: t('admin.ctaSection'), value: 'cta', icon: 'i-lucide-megaphone' },
  { label: '⚡ Quick WhatsApp CTA', value: 'cta_whatsapp', icon: 'i-simple-icons-whatsapp' }
]

function fillWhatsAppCTA(section: any, target: 'primary' | 'secondary' = 'secondary') {
  const defaultWaUrl = waLink.value || 'https://wa.me/628119124848'
  if (target === 'primary') {
    section.data.buttonText = 'Chat WhatsApp'
    section.data.buttonLink = defaultWaUrl
    section.data.buttonIcon = 'i-simple-icons-whatsapp'
  } else {
    section.data.secondaryButtonText = 'Chat WhatsApp'
    section.data.secondaryButtonLink = defaultWaUrl
    section.data.secondaryButtonIcon = 'i-simple-icons-whatsapp'
  }
  toast.add({
    title: 'Shortcut WhatsApp Berhasil',
    description: 'Tombol WhatsApp telah terisi secara otomatis.',
    color: 'success'
  })
}

function addSection(type: string) {
  let defaultData = {}
  let actualType = type

  if (type === 'cta_whatsapp') {
    actualType = 'cta'
    defaultData = {
      heading: 'Butuh Informasi Lebih Lanjut?',
      description: 'Hubungi customer service kami secara langsung via WhatsApp untuk konsultasi paket dan pendaftaran.',
      buttonText: 'Lihat Paket',
      buttonLink: '/packages',
      buttonIcon: 'i-lucide-package',
      secondaryButtonText: 'Chat WhatsApp Sekarang',
      secondaryButtonLink: waLink.value || 'https://wa.me/628119124848',
      secondaryButtonIcon: 'i-simple-icons-whatsapp'
    }
  }
  else if (type === 'hero') {
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
  else if (type === 'cta') defaultData = { heading: '', description: '', buttonText: '', buttonLink: '', buttonIcon: '', secondaryButtonText: '', secondaryButtonLink: '', secondaryButtonIcon: '' }
  else if (type === 'course_material') {
    defaultData = {
      headline: '',
      title: '',
      description: '',
      materials: []
    }
  }
  else if (type === 'specifications') {
    defaultData = {
      headline: '',
      title: '',
      description: '',
      items: []
    }
  }
  else if (type === 'service_areas') {
    defaultData = {
      headline: '',
      title: '',
      description: '',
      footer: '',
      areas: []
    }
  }
  else if (type === 'quote') {
    defaultData = {
      quote: '',
      description: '',
      ctaText: '',
      ctaLink: '',
      secondaryCtaText: '',
      secondaryCtaLink: ''
    }
  }

  formData.value.sections.push({
    id: generateId(),
    type: actualType,
    data: defaultData
  })
}

function loadAboutUsTemplate() {
  formData.value.sections = [
    {
      id: generateId(),
      type: "hero",
      data: {
        heading: "Langkah Maju, Belajar dari Masa Depan",
        subheading: "Drive Master bukan hanya tentang mengajar cara mengemudi, ini tentang mendefinisikan ulang standar pendidikan mengemudi di Indonesia. Sebagai pelopor sekolah mengemudi Kendaraan Listrik, kami percaya bahwa pengemudi masa depan harus lahir dari teknologi masa depan—modern, cerdas, dan ramah lingkungan.",
        ctaText: "Mulai Perjalanan Anda",
        ctaLink: "/auth/register",
        secondaryCtaText: "Hubungi Kami",
        secondaryCtaLink: "/contact",
        features: [
          { title: "Pelopor EV Bebas Emisi", icon: "i-lucide-leaf" },
          { title: "Instruktur Bersertifikat", icon: "i-lucide-award" }
        ]
      }
    },
    {
      id: generateId(),
      type: "specifications",
      data: {
        headline: "Keselamatan Utama",
        title: "Prioritas Kami adalah Keselamatan Anda",
        description: "Dalam bisnis kursus mengemudi, keselamatan bukan hanya sebuah fitur; itu adalah fondasi inti kami.",
        items: [
          {
            title: "Instruktur Bersertifikat",
            subtitle: "Instruktur kami adalah profesional berlisensi yang bersertifikat khusus untuk mengoperasikan kendaraan listrik premium.",
            icon: "i-lucide-award",
            description: []
          },
          {
            title: "Teknologi Keselamatan Aktif",
            subtitle: "Memanfaatkan fitur keselamatan bawaan EV seperti Collision Avoidance dan Blind Spot Monitoring untuk meminimalkan risiko.",
            icon: "i-lucide-radar",
            description: []
          }
        ]
      }
    },
    {
      id: generateId(),
      type: "quote",
      data: {
        quote: "Visi kami bukan hanya untuk menghasilkan pengemudi yang bisa memutar kemudi, tetapi untuk membina pengemudi yang cerdas dan aman yang siap merangkul era elektrifikasi.",
        description: "Di Drive Master Indonesia, kami percaya bahwa cara kita belajar mengemudi harus berevolusi seiring dengan evolusi teknologi otomotif. Kami berkomitmen untuk menjadi standar baru dalam pendidikan mengemudi yang ramah lingkungan, memastikan bahwa setiap lulusan memiliki keterampilan mengemudi tingkat tinggi serta kesadaran akan masa depan mobilitas yang berkelanjutan.",
        ctaText: "Mulai Perjalanan Anda",
        ctaLink: "/auth/register",
        secondaryCtaText: "Chat WhatsApp",
        secondaryCtaLink: waLink.value || "https://wa.me/628119124848"
      }
    }
  ]
  toast.add({
    title: "Template Loaded",
    description: "Template About Us berhasil dimuat ke editor.",
    color: "success"
  })
}

function loadContactUsTemplate() {
  formData.value.sections = [
    {
      id: generateId(),
      type: "hero",
      data: {
        heading: "Kami di Sini untuk Membantu",
        subheading: "Punya pertanyaan tentang paket mengemudi EV kami atau penjadwalan? Hubungi tim kami melalui metode di bawah ini.",
        ctaText: "Chat WhatsApp",
        ctaLink: waLink.value || "https://wa.me/628119124848",
        secondaryCtaText: "Lihat FAQ",
        secondaryCtaLink: "/#faq"
      }
    },
    {
      id: generateId(),
      type: "cta",
      data: {
        heading: "Hubungi Customer Service",
        description: "Tim kami siap memberikan informasi detail mengenai jadwal, paket, dan lokasi pelatihan.",
        buttonText: "Lihat Paket",
        buttonLink: "/packages",
        buttonIcon: "i-lucide-package",
        secondaryButtonText: "Chat WhatsApp Sekarang",
        secondaryButtonLink: waLink.value || "https://wa.me/628119124848",
        secondaryButtonIcon: "i-simple-icons-whatsapp"
      }
    }
  ]
  toast.add({
    title: "Template Loaded",
    description: "Template Contact Us berhasil dimuat ke editor.",
    color: "success"
  })
}

function removeSection(index: number | string) {
  formData.value.sections.splice(Number(index), 1)
}

function getSectionTitle(type: string) {
  return sectionTypes.find(t => t.value === type)?.label || 'Unknown Section'
}

function getSectionIcon(type: string) {
  return sectionTypes.find(t => t.value === type)?.icon || 'i-lucide-box'
}

// ==================== DRAG & DROP LOGIC ====================
const dragIndex = ref<number | null>(null)

function onDragStart(index: number | string, event: DragEvent) {
  dragIndex.value = Number(index)
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

function onDrop(targetIndex: number | string) {
  const target = Number(targetIndex)
  if (dragIndex.value === null || dragIndex.value === target) return
  const item = formData.value.sections.splice(dragIndex.value, 1)[0]
  formData.value.sections.splice(target, 0, item)
  dragIndex.value = null
}

function onDragEnd() {
  dragIndex.value = null
}

function addFeature(section: any) {
  if (!section.data.features) section.data.features = []
  section.data.features.push({ title: '', icon: 'i-lucide-check-circle' })
}

function removeFeature(section: any, index: number | string) {
  section.data.features.splice(Number(index), 1)
}

function addMaterial(section: any) {
  if (!section.data.materials) section.data.materials = []
  section.data.materials.push({ title: '', icon: 'i-lucide-book-open', description: [] })
}

function removeMaterial(section: any, index: number | string) {
  section.data.materials.splice(Number(index), 1)
}

function getBulletsText(description: string[] | undefined) {
  if (!description) return ''
  return description.join('\n')
}

function setBulletsText(material: any, text: string) {
  material.description = text.split('\n').map(line => line.trim()).filter(line => line.length > 0)
}

function addSpecification(section: any) {
  if (!section.data.items) section.data.items = []
  section.data.items.push({ title: '', subtitle: '', icon: 'i-lucide-star', description: [] })
}

function removeSpecification(section: any, index: number | string) {
  section.data.items.splice(Number(index), 1)
}

function getAreasText(areas: string[] | undefined) {
  if (!areas) return ''
  return areas.join('\n')
}

function setAreasText(section: any, text: string) {
  section.data.areas = text.split('\n').map(line => line.trim()).filter(line => line.length > 0)
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
        <UButton
          v-if="formData.slug === '/about'"
          label="Load About Us Template"
          icon="i-lucide-sparkles"
          color="warning"
          variant="soft"
          @click="loadAboutUsTemplate"
        />
        <UButton
          v-if="formData.slug === '/contact'"
          label="Load Contact Us Template"
          icon="i-lucide-sparkles"
          color="warning"
          variant="soft"
          @click="loadContactUsTemplate"
        />
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
        <div v-if="formData.slug === '/about'" class="mt-4">
          <UButton
            label="Muat Template About Us"
            icon="i-lucide-sparkles"
            color="warning"
            variant="soft"
            @click="loadAboutUsTemplate"
          />
        </div>
        <div v-if="formData.slug === '/contact'" class="mt-4">
          <UButton
            label="Muat Template Contact Us"
            icon="i-lucide-sparkles"
            color="warning"
            variant="soft"
            @click="loadContactUsTemplate"
          />
        </div>
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
              <!-- WhatsApp Quick Fill Banner -->
              <div class="col-span-2 flex flex-wrap items-center justify-between bg-green-500/10 border border-green-500/30 p-3 rounded-lg gap-2">
                <div class="flex items-center gap-2">
                  <UIcon name="i-simple-icons-whatsapp" class="size-5 text-green-600 dark:text-green-400" />
                  <span class="text-xs font-semibold text-green-700 dark:text-green-300">Shortcut Tombol WhatsApp:</span>
                </div>
                <div class="flex gap-2">
                  <UButton
                    label="+ Set Tombol Utama WA"
                    icon="i-simple-icons-whatsapp"
                    size="xs"
                    color="success"
                    variant="soft"
                    @click="fillWhatsAppCTA(section, 'primary')"
                  />
                  <UButton
                    label="+ Set Tombol Kedua WA"
                    icon="i-simple-icons-whatsapp"
                    size="xs"
                    color="success"
                    variant="solid"
                    @click="fillWhatsAppCTA(section, 'secondary')"
                  />
                </div>
              </div>

              <div class="col-span-2">
                <label class="block text-xs font-medium text-primary mb-1.5">{{ t('admin.ctaSection') }} {{ t('admin.heading') }}</label>
                <UInput v-model="section.data.heading" placeholder="e.g. Ready to start driving?" class="w-full" />
              </div>
              <div class="col-span-2">
                <label class="block text-xs font-medium text-primary mb-1.5">Description / Subtitle</label>
                <UTextarea v-model="section.data.description" placeholder="e.g. Join hundreds of students learning with our BNSP certified instructors." :rows="2" class="w-full" />
              </div>

              <!-- Primary Button -->
              <div class="col-span-2 border-t border-primary/10 pt-3">
                <span class="text-xs font-bold text-primary uppercase tracking-wider block mb-2">Primary Button</span>
                <div class="grid grid-cols-3 gap-3">
                  <div>
                    <label class="block text-[10px] font-bold text-muted uppercase mb-1">{{ t('admin.buttonText') }}</label>
                    <UInput v-model="section.data.buttonText" placeholder="e.g. Lihat Paket" class="w-full" size="sm" />
                  </div>
                  <div>
                    <label class="block text-[10px] font-bold text-muted uppercase mb-1">{{ t('admin.buttonLink') }}</label>
                    <UInput v-model="section.data.buttonLink" placeholder="e.g. /packages" class="w-full" size="sm" />
                  </div>
                  <div>
                    <label class="block text-[10px] font-bold text-muted uppercase mb-1">Button Icon</label>
                    <UInput v-model="section.data.buttonIcon" placeholder="e.g. i-lucide-package" class="w-full" size="sm" />
                  </div>
                </div>
              </div>

              <!-- Secondary Button -->
              <div class="col-span-2 border-t border-primary/10 pt-3">
                <span class="text-xs font-bold text-muted uppercase tracking-wider block mb-2">Secondary Button (Optional)</span>
                <div class="grid grid-cols-3 gap-3">
                  <div>
                    <label class="block text-[10px] font-bold text-muted uppercase mb-1">Button Text</label>
                    <UInput v-model="section.data.secondaryButtonText" placeholder="e.g. Hubungi Kami" class="w-full" size="sm" />
                  </div>
                  <div>
                    <label class="block text-[10px] font-bold text-muted uppercase mb-1">Button Link</label>
                    <UInput v-model="section.data.secondaryButtonLink" placeholder="e.g. /contact" class="w-full" size="sm" />
                  </div>
                  <div>
                    <label class="block text-[10px] font-bold text-muted uppercase mb-1">Button Icon</label>
                    <UInput v-model="section.data.secondaryButtonIcon" placeholder="e.g. i-simple-icons-whatsapp" class="w-full" size="sm" />
                  </div>
                </div>
              </div>
            </div>

            <!-- SPECIFICATIONS GRID FORM -->
            <div v-if="section.type === 'specifications'" class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-xs font-medium text-muted mb-1.5">Headline</label>
                  <UInput v-model="section.data.headline" placeholder="e.g. Layanan Utama" class="w-full" />
                </div>
                <div>
                  <label class="block text-xs font-medium text-muted mb-1.5">Title</label>
                  <UInput v-model="section.data.title" placeholder="e.g. Spesifikasi & Layanan Khusus" class="w-full" />
                </div>
                <div class="col-span-2">
                  <label class="block text-xs font-medium text-muted mb-1.5">Description</label>
                  <UInput v-model="section.data.description" placeholder="e.g. Fasilitas dan pilihan sesi terbaik untuk Anda" class="w-full" />
                </div>
              </div>

              <!-- Specifications Items Manager -->
              <div class="border-t border-default pt-4 mt-2">
                <div class="flex items-center justify-between mb-3">
                  <span class="text-xs font-bold text-muted uppercase tracking-wider">Specification Cards</span>
                  <UButton label="Add Card" icon="i-lucide-plus" size="xs" color="neutral" variant="soft" @click="addSpecification(section)" />
                </div>
                <div class="space-y-4">
                  <div 
                    v-for="(spec, sIdx) in section.data.items" 
                    :key="sIdx" 
                    class="p-4 rounded-lg border border-default bg-muted/20 relative"
                  >
                    <UButton 
                      icon="i-lucide-trash" 
                      color="error" 
                      variant="ghost" 
                      size="xs" 
                      class="absolute top-2 right-2" 
                      @click="removeSpecification(section, sIdx)" 
                    />
                    <div class="grid grid-cols-3 gap-4 mb-3">
                      <div>
                        <label class="block text-[10px] font-bold uppercase text-muted mb-1">Card Title</label>
                        <UInput v-model="spec.title" placeholder="e.g. Layanan Khusus SIM" size="sm" class="w-full" />
                      </div>
                      <div>
                        <label class="block text-[10px] font-bold uppercase text-muted mb-1">Subtitle</label>
                        <UInput v-model="spec.subtitle" placeholder="e.g. Termasuk Fasilitas Lengkap" size="sm" class="w-full" />
                      </div>
                      <div>
                        <label class="block text-[10px] font-bold uppercase text-muted mb-1">Card Icon</label>
                        <UInput v-model="spec.icon" placeholder="e.g. i-lucide-star" size="sm" class="w-full" />
                      </div>
                    </div>
                    <div>
                      <label class="block text-[10px] font-bold uppercase text-muted mb-1">Bullet Points (One per line)</label>
                      <UTextarea 
                        :model-value="getBulletsText(spec.description)" 
                        @update:model-value="setBulletsText(spec, $event)"
                        placeholder="Bullet 1&#10;Bullet 2&#10;Bullet 3" 
                        size="sm" 
                        :rows="3" 
                        class="w-full" 
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- SERVICE AREAS FORM -->
            <div v-if="section.type === 'service_areas'" class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-xs font-medium text-muted mb-1.5">Headline</label>
                  <UInput v-model="section.data.headline" placeholder="e.g. Jangkauan Layanan" class="w-full" />
                </div>
                <div>
                  <label class="block text-xs font-medium text-muted mb-1.5">Title</label>
                  <UInput v-model="section.data.title" placeholder="e.g. Area Jangkauan Kursus Kami" class="w-full" />
                </div>
                <div class="col-span-2">
                  <label class="block text-xs font-medium text-muted mb-1.5">Description</label>
                  <UInput v-model="section.data.description" placeholder="e.g. Kami melayani lokasi di wilayah berikut" class="w-full" />
                </div>
                <div class="col-span-2">
                  <label class="block text-xs font-medium text-muted mb-1.5">Service Area Locations (One per line)</label>
                  <UTextarea 
                    :model-value="getAreasText(section.data.areas)" 
                    @update:model-value="setAreasText(section, $event)"
                    placeholder="Alam Sutera & sekitarnya&#10;Serpong & BSD City&#10;Tangerang Kota" 
                    :rows="4" 
                    class="w-full" 
                  />
                </div>
                <div class="col-span-2">
                  <label class="block text-xs font-medium text-muted mb-1.5">Footer Note / Description</label>
                  <UInput v-model="section.data.footer" placeholder="e.g. Lokasi Anda belum tertera? Hubungi tim kami." class="w-full" />
                </div>
              </div>
            </div>

            <!-- QUOTE / VISION FORM -->
            <div v-if="section.type === 'quote'" class="grid grid-cols-2 gap-4 bg-warning/5 p-4 rounded-lg border border-warning/20">
              <div class="col-span-2">
                <label class="block text-xs font-medium text-warning mb-1.5">Quote / Highlight Text</label>
                <UTextarea v-model="section.data.quote" placeholder="e.g. Driving for a better, greener future." :rows="3" class="w-full" />
              </div>
              <div class="col-span-2">
                <label class="block text-xs font-medium text-muted mb-1.5">Subdescription</label>
                <UTextarea v-model="section.data.description" placeholder="Additional explanation..." :rows="2" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Primary CTA Text</label>
                <UInput v-model="section.data.ctaText" placeholder="e.g. Mulai Sekarang" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Primary CTA Link</label>
                <UInput v-model="section.data.ctaLink" placeholder="e.g. /auth/register" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Secondary CTA Text</label>
                <UInput v-model="section.data.secondaryCtaText" placeholder="e.g. Hubungi Kami" class="w-full" />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted mb-1.5">Secondary CTA Link</label>
                <UInput v-model="section.data.secondaryCtaLink" placeholder="e.g. /contact" class="w-full" />
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
