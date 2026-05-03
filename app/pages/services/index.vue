<script setup lang="ts">
useSeoMeta({
  title: 'Services | Drive Master Academy',
  description: 'Comprehensive driving courses designed for the electric future. From beginners to advanced drivers, we have the perfect program for you.',
})

const services = [
  {
    title: '6x Sessions',
    description: 'Foundation course for complete beginners. Learn the fundamentals of driving in our premium electric vehicles.',
    icon: 'i-lucide-car',
    features: [
      'Free Trial',
      '6x Sessions',
      '60 minutes per session',
      'Night driving session',
      'Weekend driving session',
      'Matic transmission',
      'Free shuttle (Alam Sutra, BSD, Gading Serpong)',
      'Driving license administration',
      'Theory materials included',
      'Certificate of completion',
    ]
  },
  {
    title: '8x Sessions',
    description: 'Take your skills to the next level with advanced techniques and real-world scenarios.',
    icon: 'i-lucide-car',
    features: [
      'Free Trial',
      '8x Sessions',
      '60 minutes per session',
      'Night driving session',
      'Weekend driving session',
      'Matic transmission',
      'Free shuttle (Alam Sutra, BSD, Gading Serpong)',
      'Driving license administration',
      'Theory materials included',
      'Certificate of completion',
    ]
  },
  {
    title: '10x Sessions',
    description: 'Comprehensive Session focused on electric vehicle specific features and best practices.',
    icon: 'i-lucide-car',
    features: [
      'Free Trial',
      '10x Sessions',
      '60 minutes per session',
      'Night driving session',
      'Weekend driving session',
      'Matic transmission',
      'Free shuttle (Alam Sutra, BSD, Gading Serpong)',
      'Driving license administration',
      'Theory materials included',
      'Certificate of completion',
    ]
  },
  {
    title: '12x Sessions',
    description: 'Perfect for licensed drivers who want to brush up their skills or transition to EV driving.',
    icon: 'i-lucide-car',
    features: [
      'Free Trial',
      '12x Sessions',
      '60 minutes per session',
      'Night driving session',
      'Weekend driving session',
      'Matic transmission',
      'Free shuttle (Alam Sutra, BSD, Gading Serpong)',
      'Driving license administration',
      'Theory materials included',
      'Certificate of completion',
    ]
  }
]

const specifications = [
  {
    title: 'Special Services',
    subtitle: 'The Advantage of choosing Us',
    description: [
      'SIM',
      'Free shuttle (Alam Sutra, BSD, Gading Serpong)',
      'Theory materials included',
      'Certificate of completion',
    ],
    icon: 'i-lucide-star'
  },
  {
    title: 'Session Time',
    subtitle: 'Our original session time',
    description: [
      'Session in weekdays(days) open from 08:00 to 17:00',
      '60 minutes per session',
      'Night driving session',
      'Weekend driving session',
    ],
    icon: 'i-lucide-car'
  },
  {
    title: 'Car Transmission',
    subtitle: 'Our car using matic transmission',
    description: [
      'Matic Transmission',
    ],
    icon: 'i-lucide-car'
  },
  {
    title: 'Night Session',
    subtitle: 'Additional price for night session',
    description: [
      'Session at night open from 18:00 to 20:00',
      'Base price 6x sessions + Rp.100,000 for night session',
      'Base price 8x sessions + Rp.150,000 for night session',
      'Base price 10x sessions + Rp.200,000 for night session',
      'Base price 12x sessions + Rp.250,000 for night session',
    ],
    icon: 'i-lucide-moon'
  },
  {
    title: 'Weekend Session',
    subtitle: 'Additional price for weekend session',
    description: [
      'Session in weekend(days) open from 08:00 to 17:00',
      'Base price 6x sessions + Rp.100,000 for session in weekend',
      'Base price 8x sessions + Rp.150,000 for session in weekend',
      'Base price 10x sessions + Rp.200,000 for session in weekend',
      'Base price 12x sessions + Rp.250,000 for session in weekend',
    ],
    icon: 'i-lucide-clock'
  }
]

// Course Material
const courseMaterial = [
  {
    title: 'Material Theory',
    description: [
      'Vehicle Introduction and Basic Controls',
      'Cockpit training (ergonomic seating position, center and side mirror adjustment, and seat belt use)',
      'Introduction to instruments (gas pedal, brake, transmission lever, handbrake, indicator lights in the dashboard)',
      'Safety check (check the condition of tires, oil and radiator water before driving)'
    ],
    icon:'i-lucide-book-open'
  },
  {
    title: 'Initial Control',
    description: [
      'Starting & stopping the engine (standard procedure for starting the engine safely)',
      'Accelerator pedal technique safely (balanced and smooth)',
      'Braking and stopping techniques (smooth braking and how to stop at a certain point precisely)',
    ],
    icon: 'i-lucide-shield-check'
  },
  {
    title: 'Basic Maneuvering Techniques',
    description: [
      'Steering control (steering wheel turning technique when turning quickly)',
      'Reverse (controlling the car to reverse using only the rearview mirror)',
      'Turning at an intersection (technique for taking the correct turning angle to the left or right)',
    ],
    icon: 'i-lucide-radar'
  },
  {
    title: 'Driving Techniques on Uphill & Downhill Roads',
    description: [
      'Start-stop technique on an incline',
      'Start-stop technique on a downhill',
    ],
    icon: 'i-lucide-car'
  },
  {
    title: 'Parking Technic',
    description: [
      'Reverse parking at an angle or straight (enter the parking slot with the car in reverse)',
      'Parallel parking (a technique of inserting a car between two other vehicles parallel to each other)',
    ],
    icon: 'i-lucide-car'
  },
  {
    title: 'Driving on the Highway',
    description: [
      'Road signs and markings (obey traffic signs, no parking signs and road markings)',
      'Driving ethics (using turn signals, maintaining a safe distance, and how to overtake other vehicles correctly)',
      'Blind spot (a technique for checking areas that are not visible in the rearview mirror before changing lanes)',
    ],
    icon: 'i-lucide-car'
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
    />

    <!-- Our Services -->
    <UPageSection
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="grid md:grid-cols-1 gap-6">
        <UCard v-for="specification in specifications" :key="specification.title">
          <template #header>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-xl bg-warning/10">
                <UIcon :name="specification.icon" class="size-8 text-warning" />
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ specification.title }}</h3>
                <p class="text-muted mt-1">{{ specification.subtitle }}</p>
              </div>
            </div>
          </template>

          <ul class="space-y-2">
            <li v-for="spec in specification.description" :key="spec" class="flex items-center gap-2 font-bold">
              <span class="text-sm">{{ spec }}</span>
            </li>
          </ul>

        </UCard>
      </div>
    </UPageSection>

    <!-- Course Material -->
    <UPageSection
      id="material"
      headline="Course material you will study"
      title="Course Material"
      description="Course material that you learn will make you more confident in driving."
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <UPageGrid>
        <UPageCard
          v-for="material in courseMaterial"
          :key="material.title"
          :icon="material.icon"
          :title="material.title"
          :ui="{ leadingIcon: 'text-warning text-3xl' }"
        >
          <template #description>
            <ul class="space-y-2 mt-2">
              <li v-for="item in material.description" :key="item" class="flex items-start gap-2 text-muted">
                <UIcon name="i-lucide-check-circle" class="size-4 text-warning shrink-0 mt-1" />
                <span class="text-sm">{{ item }}</span>
              </li>
            </ul>
          </template>
        </UPageCard>
      </UPageGrid>
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
          Sessions depart from our Alam Sutera center. Contact us for specific route availability.
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
