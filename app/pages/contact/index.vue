<script setup lang="ts">
// Page metadata
definePageMeta({
  layout: 'default'
})

useSeoMeta({
  title: 'Contact Us | Drive Master - Get in Touch',
  description: 'Contact Drive Master in Alam Sutera. Reach out via WhatsApp, phone, email, or visit our training center.'
})

const contactMethods = [
  {
    title: 'Training Center',
    description: 'Jl. Alam Sutera Boulevard No. 123, Alam Sutera, Tangerang 15143',
    icon: 'i-lucide-map-pin',
    action: { label: 'Get Directions', to: 'https://maps.app.goo.gl/RpSdkpjs4RZg2ZY77', target: '_blank' }
  },
  {
    title: 'WhatsApp Support',
    description: '+62 811-9124-848 (Available 08:00 - 18:00)',
    icon: 'i-simple-icons-whatsapp',
    action: { label: 'Chat Now', to: 'https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi', target: '_blank' }
  },
  {
    title: 'Email Address',
    description: 'info@evdriveacademy.id',
    icon: 'i-lucide-mail',
    action: { label: 'Send Email', to: 'mailto:info@evdriveacademy.id' }
  },
  {
    title: 'Operating Hours',
    description: 'Weekdays: 08:00 - 17:00 | Weekend: 08:00 - 17:00 | Night Shift: 18:00-20:00',
    icon: 'i-lucide-clock',
    action: { label: 'View FAQ', to: '/#faq' }
  }
]

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
    title: 'Message Sent!',
    description: "Thank you for reaching out. We'll get back to you within 24 hours.",
    color: 'success',
    icon: 'i-lucide-check-circle'
  })
  
  form.name = ''
  form.email = ''
  form.subject = ''
  form.message = ''
  isSubmitting.value = false
}
</script>

<template>
  <div>
    <!-- Hero Section -->
    <UPageHero
      title="We&apos;re Here to Help"
      description="Have questions about our EV driving packages or scheduling? Reach out to our team via any of the methods below or fill out the form."
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
      headline="Get in Touch"
      title="Send us a Message"
      description="Fill out the form below and our customer success team will contact you shortly."
    >
      <div class="grid lg:grid-cols-2 gap-12 items-start">
        <!-- Contact Form -->
        <UCard>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div class="grid sm:grid-cols-2 gap-4">
              <UFormField label="Full Name" required>
                <UInput v-model="form.name" placeholder="John Doe" class="w-full" />
              </UFormField>
              <UFormField label="Email Address" required>
                <UInput v-model="form.email" type="email" placeholder="john@example.com" class="w-full" />
              </UFormField>
            </div>
            <UFormField label="Subject" required>
              <USelect
                v-model="form.subject"
                :items="['General Inquiry', 'Package Information', 'Scheduling Issue', 'Technical Support']"
                placeholder="Select a subject"
                class="w-full"
              />
            </UFormField>
            <UFormField label="Message" required>
              <UTextarea v-model="form.message" placeholder="How can we help you?" :rows="5" class="w-full" />
            </UFormField>
            <div class="pt-2">
              <UButton 
                type="submit" 
                label="Send Message" 
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
              src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3966.1806061048765!2d106.65588507475077!3d-6.2399118937483635!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69fbc070b4d71d%3A0x8b1a633faf5dbd46!2sALAM%20SUTERA!5e0!3m2!1sen!2sid!4v1776223155011!5m2!1sen!2sid" 
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
      headline="Social Media"
      title="Join Our Community"
      description="Follow us on social media for driving tips, EV news, and student success stories."
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
