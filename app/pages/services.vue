<script setup lang="ts">
const services = [
  {
    title: 'Basic Driving Course',
    description: 'Foundation course for complete beginners. Learn the fundamentals of driving in our premium electric vehicles.',
    icon: 'i-lucide-car',
    features: [
      'Vehicle controls and safety features',
      'Basic maneuvering and parking',
      'Traffic rules and regulations',
      'Road signs recognition',
      'EV-specific controls introduction'
    ]
  },
  {
    title: 'Advanced Driving Course',
    description: 'Take your skills to the next level with advanced techniques and real-world scenarios.',
    icon: 'i-lucide-gauge',
    features: [
      'Highway and expressway driving',
      'Defensive driving techniques',
      'Night driving mastery',
      'Weather condition driving',
      'Emergency handling procedures'
    ]
  },
  {
    title: 'EV Specialist Course',
    description: 'Comprehensive training focused on electric vehicle specific features and best practices.',
    icon: 'i-lucide-battery-charging',
    features: [
      'Regenerative braking mastery',
      'One-pedal driving technique',
      'Battery management basics',
      'Charging station protocols',
      'EV maintenance awareness'
    ]
  },
  {
    title: 'Refresher Course',
    description: 'Perfect for licensed drivers who want to brush up their skills or transition to EV driving.',
    icon: 'i-lucide-refresh-cw',
    features: [
      'Skills assessment',
      'Personalized training plan',
      'Updated traffic regulations',
      'EV adaptation for manual drivers',
      'Confidence building sessions'
    ]
  }
]

const methodologies = [
  {
    title: 'Theory Learning',
    description: 'Comprehensive digital materials covering traffic regulations, road safety, and EV-specific knowledge.',
    icon: 'i-lucide-book-open'
  },
  {
    title: 'Simulator Training',
    description: 'Practice in our advanced driving simulators before hitting the real road.',
    icon: 'i-lucide-monitor'
  },
  {
    title: 'Practical Sessions',
    description: 'Hands-on training in our premium electric vehicles with experienced instructors.',
    icon: 'i-lucide-car'
  },
  {
    title: 'Assessment & Certification',
    description: 'Comprehensive evaluation and digital certificate upon successful completion.',
    icon: 'i-lucide-award'
  }
]

const serviceAreas = [
  'Alam Sutera & surrounding areas',
  'Serpong & BSD City',
  'Tangerang City Center',
  'Gading Serpong',
  'Lippo Karawaci',
  'Bintaro Jaya (limited)'
]
</script>

<template>
  <div>
    <!-- Hero -->
    <UPageHero
      title="Our Services"
      description="Comprehensive driving courses designed for the electric future. From beginners to advanced drivers, we have the perfect program for you."
      :links="[
        { label: 'View Packages', to: '/packages', color: 'warning', icon: 'i-lucide-package' },
        { label: 'Book Consultation', to: 'https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi', color: 'primary', variant: 'outline', icon: 'i-simple-icons-whatsapp', external: true }
      ]"
      class="py-12 lg:py-20"
    />

    <!-- Service Types -->
    <UPageSection
      headline="Training Programs"
      title="Choose Your Learning Path"
      description="We offer a variety of courses tailored to different skill levels and learning objectives."
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid md:grid-cols-2 gap-6">
        <UCard v-for="service in services" :key="service.title">
          <template #header>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-xl bg-warning/10">
                <UIcon :name="service.icon" class="size-8 text-warning" />
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ service.title }}</h3>
                <p class="text-muted mt-1">{{ service.description }}</p>
              </div>
            </div>
          </template>

          <ul class="space-y-2">
            <li v-for="feature in service.features" :key="feature" class="flex items-center gap-2">
              <UIcon name="i-lucide-check-circle" class="size-4 text-warning shrink-0" />
              <span class="text-sm">{{ feature }}</span>
            </li>
          </ul>

          <template #footer>
            <NuxtLink to="/packages">
              <UButton label="View Packages" variant="outline" color="neutral" trailingIcon="i-lucide-arrow-right" />
            </NuxtLink>
          </template>
        </UCard>
      </div>
    </UPageSection>

    <!-- Training Methodology -->
    <UPageSection
      headline="How We Teach"
      title="Our Training Methodology"
      description="A structured approach that combines theory, simulation, and practical experience."
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="(method, index) in methodologies" :key="method.title" class="text-center">
          <div class="relative mb-4">
            <div class="mx-auto w-16 h-16 rounded-full bg-warning/10 flex items-center justify-center">
              <UIcon :name="method.icon" class="size-8 text-warning" />
            </div>
            <div class="absolute -top-2 -right-2 w-8 h-8 rounded-full bg-warning text-white flex items-center justify-center text-sm font-bold">
              {{ index + 1 }}
            </div>
          </div>
          <h3 class="font-semibold mb-2">{{ method.title }}</h3>
          <p class="text-sm text-muted">{{ method.description }}</p>
        </div>
      </div>
    </UPageSection>

    <!-- Facilities -->
    <UPageSection
      headline="Premium Fleet"
      title="Learn in the Best Electric Vehicles"
      description="Our fleet consists of premium electric vehicles equipped with dual-control systems for safe learning."
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid md:grid-cols-2 gap-8">
        <UCard class="overflow-hidden">
          <img 
            src="https://images.unsplash.com/photo-1560958089-b8a1929cea89?w=600&auto=format&fit=crop&q=80" 
            alt="Tesla Model 3"
            class="w-full h-48 object-cover -m-4 mb-4"
          />
          <h3 class="text-xl font-bold mb-2">Tesla Model 3</h3>
          <p class="text-muted mb-4">Experience the future of driving with Tesla&apos;s revolutionary autopilot features and instant torque delivery.</p>
          <div class="flex flex-wrap gap-2">
            <UBadge label="Autopilot" variant="subtle" color="warning" />
            <UBadge label="Dual Motor" variant="subtle" color="warning" />
            <UBadge label="Dual Control" variant="subtle" color="warning" />
          </div>
        </UCard>

        <UCard class="overflow-hidden">
          <img 
            src="https://images.unsplash.com/photo-1619767886558-efdc259cde1a?w=600&auto=format&fit=crop&q=80" 
            alt="BYD Atto 3"
            class="w-full h-48 object-cover -m-4 mb-4"
          />
          <h3 class="text-xl font-bold mb-2">BYD Atto 3</h3>
          <p class="text-muted mb-4">A compact SUV perfect for urban driving lessons with excellent visibility and responsive handling.</p>
          <div class="flex flex-wrap gap-2">
            <UBadge label="SUV" variant="subtle" color="warning" />
            <UBadge label="Urban Friendly" variant="subtle" color="warning" />
            <UBadge label="Dual Control" variant="subtle" color="warning" />
          </div>
        </UCard>
      </div>
    </UPageSection>

    <!-- Service Areas -->
    <UPageSection
      headline="Coverage"
      title="Service Areas"
      description="Our training routes cover major areas around Alam Sutera and Tangerang."
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="max-w-3xl mx-auto">
        <UCard>
          <div class="grid sm:grid-cols-2 gap-4">
            <div v-for="area in serviceAreas" :key="area" class="flex items-center gap-3 p-3 rounded-lg hover:bg-muted/50 transition-colors">
              <UIcon name="i-lucide-map-pin" class="size-5 text-warning" />
              <span>{{ area }}</span>
            </div>
          </div>
        </UCard>
        
        <p class="text-center text-muted mt-6">
          Training sessions depart from our Alam Sutera center. Contact us for specific route availability.
        </p>
      </div>
    </UPageSection>

    <!-- CTA -->
    <UPageCTA
      title="Ready to Start?"
      description="Book a free consultation or view our package options to begin your EV driving journey."
      :links="[
        { label: 'View Packages', to: '/packages', color: 'warning', icon: 'i-lucide-package' },
        { label: 'WhatsApp Consultation', to: 'https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi', color: 'primary', variant: 'outline', icon: 'i-simple-icons-whatsapp', external: true }
      ]"
    />
  </div>
</template>
