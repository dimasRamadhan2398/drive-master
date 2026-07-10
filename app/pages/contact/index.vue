<script setup lang="ts">
const { t } = useI18n()

// Page metadata
definePageMeta({
  layout: 'default'
})

useSeoMeta({
  title: t('contact.title') + ' | Drive Master - Get in Touch',
  description: 'Contact Drive Master in Alam Sutera. Reach out via WhatsApp, phone, email, or visit our training center.'
})

const settingsStore = useSettingsStore()

onMounted(() => {
  settingsStore.fetchGeneralSettings()
})

const waLink = computed(() => {
  const number = settingsStore.generalSettings?.whatsApp?.replace(/\D/g, '') ?? '628119124848'
  const normalized = number.startsWith('0') ? `62${number.slice(1)}` : number
  return `https://wa.me/${normalized}?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi`
})

const operatingHoursStr = computed(() => {
  const s = settingsStore.generalSettings
  if (!s) return `${t('home.hoursWeekday')} | ${t('home.hoursWeekend')} | ${t('home.hoursNight')}`
  
  const weekday = s.hoursMonFri ? `Mon-Fri: ${s.hoursMonFri}` : t('home.hoursWeekday')
  const weekend = s.hoursSatSun ? `Sat-Sun: ${s.hoursSatSun}` : t('home.hoursWeekend')
  const night = s.hoursNightShift ? `Night: ${s.hoursNightShift}` : t('home.hoursNight')
  
  return `${weekday} | ${weekend} | ${night}`
})

const contactMethods = computed(() => [
  {
    title: t('contact.trainingCenter'),
    description: settingsStore.generalSettings?.address ?? t('home.address'),
    icon: 'i-lucide-map-pin',
    action: { label: t('contact.getDirections'), to: 'https://maps.app.goo.gl/RpSdkpjs4RZg2ZY77', target: '_blank' }
  },
  {
    title: t('contact.whatsappSupport'),
    description: `${settingsStore.generalSettings?.whatsApp ?? '+62 811-9124-848'} (Available 08:00 - 18:00)`,
    icon: 'i-simple-icons-whatsapp',
    action: { label: t('contact.chatNow'), to: waLink.value, target: '_blank' }
  },
  {
    title: t('contact.emailAddress'),
    description: settingsStore.generalSettings?.email ?? 'info@evdriveacademy.id',
    icon: 'i-lucide-mail',
    action: { label: t('contact.sendEmail'), to: `mailto:${settingsStore.generalSettings?.email ?? 'info@evdriveacademy.id'}` }
  },
  {
    title: t('contact.operatingHours'),
    description: operatingHoursStr.value,
    icon: 'i-lucide-clock',
    action: { label: t('contact.viewFaq'), to: '/#faq' }
  }
])

const form = reactive({
  name: '',
  email: '',
  subject: '',
  message: ''
})

const isSubmitting = ref(false)
const toast = useToast()

async function handleSubmit() {
  isSubmitting.value = true
  // Mock API call
  await new Promise(resolve => setTimeout(resolve, 1500))
  
  toast.add({
    title: t('contact.form.success'),
    description: t('contact.form.successDesc'),
    color: 'success',
    icon: 'i-lucide-check-circle'
  })
  
  form.name = ''
  form.email = ''
  form.subject = ''
  form.message = ''
  isSubmitting.value = false
}

const subjects = computed(() => [
  { label: t('contact.form.subjects.general'), value: 'General Inquiry' },
  { label: t('contact.form.subjects.package'), value: 'Package Information' },
  { label: t('contact.form.subjects.schedule'), value: 'Scheduling Issue' },
  { label: t('contact.form.subjects.technical'), value: 'Technical Support' }
])
</script>

<template>
  <div>
    <!-- Hero Section -->
    <UPageHero
      :title="t('contact.heroTitle')"
      :description="t('contact.heroDesc')"
      align="center"
      class="py-16 md:py-24"
    />

    <!-- Contact Info Grid -->
    <UPageSection class="bg-muted/30">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <UPageCard
          v-for="method in contactMethods"
          :key="method.title"
          :icon="method.icon"
          :title="method.title"
          :description="method.description"
          :ui="{ leadingIcon: 'text-warning text-3xl' }"
        >
          <template #footer>
            <UButton v-bind="method.action" color="warning" variant="link" class="p-0" trailing-icon="i-lucide-arrow-right" />
          </template>
        </UPageCard>
      </div>
    </UPageSection>

    <!-- Form & Map Section -->
    <UPageSection
      :headline="t('contact.form.headline')"
      :title="t('contact.form.title')"
      :description="t('contact.form.description')"
    >
      <div class="grid lg:grid-cols-2 gap-12 items-start">
        <!-- Contact Form -->
        <UCard>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div class="grid sm:grid-cols-2 gap-4">
              <UFormField :label="t('contact.form.name')" required>
                <UInput v-model="form.name" placeholder="John Doe" class="w-full" />
              </UFormField>
              <UFormField :label="t('contact.form.email')" required>
                <UInput v-model="form.email" type="email" placeholder="john@example.com" class="w-full" />
              </UFormField>
            </div>
            <UFormField :label="t('contact.form.subject')" required>
              <USelect
                v-model="form.subject"
                :items="subjects"
                :placeholder="t('contact.form.subjectPlaceholder')"
                class="w-full"
              />
            </UFormField>
            <UFormField :label="t('contact.form.message')" required>
              <UTextarea v-model="form.message" :placeholder="t('contact.form.messagePlaceholder')" :rows="5" class="w-full" />
            </UFormField>
            <div class="pt-2">
              <UButton 
                type="submit" 
                :label="t('contact.form.send')"
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
              src="https://maps.google.com/maps?q=-6.22369663061115,106.66409468196608&z=17&output=embed" 
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
      :headline="t('contact.social.headline')"
      :title="t('contact.social.title')"
      :description="t('contact.social.description')"
      class="bg-muted/30"
    >
      <div class="flex flex-wrap justify-center gap-6">
        <UButton icon="i-simple-icons-instagram" label="Instagram" color="neutral" variant="outline" size="lg" />
        <UButton icon="i-simple-icons-youtube" label="YouTube" color="neutral" variant="outline" size="lg" />
        <UButton icon="i-simple-icons-tiktok" label="TikTok" color="neutral" variant="outline" size="lg" />
        <UButton icon="i-simple-icons-facebook" label="Facebook" color="neutral" variant="outline" size="lg" />
      </div>
    </UPageSection>
  </div>
</template>
