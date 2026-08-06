<script setup lang="ts">
import { ref, reactive, computed } from 'vue'

defineProps<{
  section: {
    type: string
    data: any
  }
}>()

const { t } = useI18n()
const toast = useToast()
const api = useApi()

const isSubmitting = ref(false)
const form = reactive({
  name: '',
  email: '',
  subject: '',
  message: ''
})

const subjects = computed(() => [
  { label: t('contact.form.subjects.general') || 'Pertanyaan Umum', value: 'General Inquiry' },
  { label: t('contact.form.subjects.package') || 'Informasi Paket', value: 'Package Information' },
  { label: t('contact.form.subjects.schedule') || 'Jadwal & Sesi', value: 'Scheduling Issue' },
  { label: t('contact.form.subjects.technical') || 'Dukungan Teknis', value: 'Technical Support' }
])

async function handleSubmit() {
  isSubmitting.value = true
  try {
    await api('/contact', {
      method: 'POST',
      body: {
        name: form.name,
        email: form.email,
        subject: form.subject,
        message: form.message
      }
    })
    toast.add({
      title: t('contact.form.success') || 'Pesan Terkirim',
      description: t('contact.form.successDesc') || 'Terima kasih telah menghubungi kami. Tim kami akan segera merespon.',
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    form.name = ''
    form.email = ''
    form.subject = ''
    form.message = ''
  } catch (error) {
    toast.add({
      title: t('common.error') || 'Gagal Mengirim',
      description: 'Gagal mengirim pesan. Silakan coba lagi.',
      color: 'error',
      icon: 'i-lucide-alert-circle'
    })
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div>
    <!-- Hero Section -->
    <UPageHero
      v-if="section.type === 'hero'"
      :title="section.data.heading"
      :description="section.data.subheading"
      :orientation="section.data.bgImage ? 'horizontal' : 'vertical'"
      :links="[
        ...(section.data.ctaText ? [{ label: section.data.ctaText, to: section.data.ctaLink || '#', color: 'warning' as const, size: 'lg' as const, icon: 'i-lucide-calendar-check' }] : []),
        ...(section.data.secondaryCtaText ? [{ label: section.data.secondaryCtaText, to: section.data.secondaryCtaLink || '#', color: 'neutral' as const, variant: 'outline' as const, trailingIcon: 'i-lucide-arrow-right', size: 'lg' as const }] : [])
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
      <div v-else-if="section.data.features && section.data.features.length > 0" class="mt-6 flex flex-wrap justify-center items-center gap-4">
        <div
          v-for="feat in section.data.features"
          :key="feat.title"
          class="inline-flex items-center gap-2 text-warning font-semibold text-sm"
        >
          <UIcon :name="feat.icon || 'i-lucide-leaf'" class="size-4 text-warning" />
          <span>{{ feat.title }}</span>
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

    <!-- Specifications Section -->
    <UPageSection
      v-else-if="section.type === 'specifications'"
      :headline="section.data.headline"
      :title="section.data.title"
      :description="section.data.description"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="grid md:grid-cols-1 gap-6">
        <UCard v-for="spec in section.data.items" :key="spec.title">
          <template #header>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-xl bg-warning/10">
                <UIcon :name="spec.icon || 'i-lucide-star'" class="size-8 text-warning" />
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ spec.title }}</h3>
                <p v-if="spec.subtitle" class="text-muted mt-1">{{ spec.subtitle }}</p>
              </div>
            </div>
          </template>

          <ul v-if="spec.description && spec.description.length > 0" class="space-y-2">
            <li v-for="bullet in spec.description" :key="bullet" class="flex items-center gap-2 font-bold">
              <span class="text-sm">{{ bullet }}</span>
            </li>
          </ul>
        </UCard>
      </div>
    </UPageSection>

    <!-- Service Areas Section -->
    <UPageSection
      v-else-if="section.type === 'service_areas'"
      :headline="section.data.headline"
      :title="section.data.title"
      :description="section.data.description"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="max-w-3xl mx-auto">
        <UCard>
          <div class="grid sm:grid-cols-2 gap-4">
            <div v-for="area in section.data.areas" :key="area" class="flex items-center gap-3 p-3 rounded-lg hover:bg-muted/50 transition-colors">
              <UIcon name="i-lucide-map-pin" class="size-5 text-warning" />
              <span>{{ area }}</span>
            </div>
          </div>
        </UCard>
        
        <p v-if="section.data.footer" class="text-center text-muted mt-6">
          {{ section.data.footer }}
        </p>
      </div>
    </UPageSection>

    <!-- Quote Section -->
    <UPageSection
      v-else-if="section.type === 'quote'"
      class="bg-warning/5 border-y border-warning/10"
    >
      <div class="max-w-3xl mx-auto text-center space-y-8">
        <UIcon name="i-lucide-quote" class="size-12 text-warning mx-auto opacity-50" />
        <div class="space-y-6">
          <p class="text-2xl font-medium leading-relaxed italic text-default">
            {{ section.data.quote }}
          </p>
          <p v-if="section.data.description" class="text-lg text-muted leading-relaxed">
            {{ section.data.description }}
          </p>
        </div>
        <div v-if="section.data.ctaText || section.data.secondaryCtaText" class="pt-8 flex flex-wrap justify-center gap-3">
          <UButton
            v-if="section.data.ctaText"
            :label="section.data.ctaText"
            :to="section.data.ctaLink || '#'"
            color="warning"
            size="xl"
            icon="i-lucide-rocket"
          />
          <UButton
            v-if="section.data.secondaryCtaText"
            :label="section.data.secondaryCtaText"
            :to="section.data.secondaryCtaLink || '#'"
            :color="(section.data.secondaryCtaLink?.includes('wa.me') || section.data.secondaryCtaText?.toLowerCase().includes('whatsapp')) ? 'success' : 'neutral'"
            size="xl"
            icon="i-simple-icons-whatsapp"
            variant="outline"
          />
        </div>
      </div>
    </UPageSection>

    <!-- Contact Info Cards Grid Section -->
    <UPageSection v-else-if="section.type === 'contact_methods'" class="bg-muted/30">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <UPageCard
          v-for="method in section.data.methods"
          :key="method.title"
          :icon="method.icon || 'i-lucide-info'"
          :title="method.title"
          :description="method.description"
          :ui="{ leadingIcon: 'text-warning text-3xl' }"
        >
          <template #footer>
            <UButton
              v-if="method.actionText"
              :label="method.actionText"
              :to="method.actionLink || '#'"
              :target="method.target || '_self'"
              color="warning"
              variant="link"
              class="p-0"
              trailing-icon="i-lucide-arrow-right"
            />
          </template>
        </UPageCard>
      </div>
    </UPageSection>

    <!-- Form & Map Section -->
    <UPageSection
      v-else-if="section.type === 'contact_form_map'"
      :headline="section.data.headline"
      :title="section.data.title"
      :description="section.data.description"
    >
      <div class="grid lg:grid-cols-2 gap-12 items-start">
        <!-- Contact Form -->
        <UCard>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div class="grid sm:grid-cols-2 gap-4">
              <UFormField :label="t('contact.form.name') || 'Nama Lengkap'" required>
                <UInput v-model="form.name" placeholder="John Doe" class="w-full" />
              </UFormField>
              <UFormField :label="t('contact.form.email') || 'Alamat Email'" required>
                <UInput v-model="form.email" type="email" placeholder="john@example.com" class="w-full" />
              </UFormField>
            </div>
            <UFormField :label="t('contact.form.subject') || 'Subjek'" required>
              <USelect
                v-model="form.subject"
                :items="subjects"
                :placeholder="t('contact.form.subjectPlaceholder') || 'Pilih subjek'"
                class="w-full"
              />
            </UFormField>
            <UFormField :label="t('contact.form.message') || 'Pesan'" required>
              <UTextarea v-model="form.message" :placeholder="t('contact.form.messagePlaceholder') || 'Bagaimana kami bisa membantu Anda?'" :rows="5" class="w-full" />
            </UFormField>
            <div class="pt-2">
              <UButton 
                type="submit" 
                :label="t('contact.form.send') || 'Kirim Pesan'"
                icon="i-lucide-send" 
                color="warning" 
                :loading="isSubmitting"
                block 
              />
            </div>
          </form>
        </UCard>

        <!-- Large Map View -->
        <div class="space-y-6">
          <div class="aspect-square rounded-2xl overflow-hidden border border-default shadow-lg bg-elevated">
            <iframe 
              :src="section.data.mapEmbedUrl || 'https://maps.google.com/maps?q=-6.22369663061115,106.66409468196608&z=17&output=embed'" 
              width="100%" 
              height="100%" 
              style="border:0;" 
              allowfullscreen 
              loading="lazy" 
              referrerpolicy="no-referrer-when-downgrade"
            ></iframe>
          </div>
        </div>
      </div>
    </UPageSection>

    <!-- Social Media Section -->
    <UPageSection
      v-else-if="section.type === 'social_media'"
      :headline="section.data.headline"
      :title="section.data.title"
      :description="section.data.description"
      class="bg-muted/30"
    >
      <div class="flex flex-wrap justify-center gap-6">
        <UButton
          v-for="link in section.data.links"
          :key="link.label"
          :icon="link.icon || 'i-lucide-link'"
          :label="link.label"
          color="neutral"
          variant="outline"
          size="lg"
          :to="link.to || '#'"
          target="_blank"
        />
      </div>
    </UPageSection>

    <!-- CTA Section -->
    <UPageCTA
      v-else-if="section.type === 'cta'"
      :title="section.data.heading"
      :description="section.data.description"
      :links="[
        ...(section.data.buttonText ? [{
          label: section.data.buttonText,
          to: section.data.buttonLink || '#',
          color: (section.data.buttonIcon?.includes('whatsapp') || section.data.buttonLink?.includes('wa.me')) ? ('success' as const) : ('warning' as const),
          variant: (section.data.buttonIcon?.includes('whatsapp') || section.data.buttonLink?.includes('wa.me')) ? ('outline' as const) : undefined,
          size: 'lg' as const,
          icon: section.data.buttonIcon || undefined
        }] : []),
        ...(section.data.secondaryButtonText ? [{
          label: section.data.secondaryButtonText,
          to: section.data.secondaryButtonLink || '#',
          color: (section.data.secondaryButtonIcon?.includes('whatsapp') || section.data.secondaryButtonLink?.includes('wa.me')) ? ('success' as const) : ('neutral' as const),
          variant: 'outline' as const,
          size: 'lg' as const,
          icon: section.data.secondaryButtonIcon || undefined
        }] : [])
      ]"
      class="bg-warning/5"
    />
  </div>
</template>
