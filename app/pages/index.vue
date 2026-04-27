<script setup lang="ts">
import { ref, computed } from 'vue'

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

// Course packages
const packages = [
  {
    title: '6x',
    price: 'Rp 1.750.000',
    description: 'Perfect for beginners looking to get started',
    features: [
      'Free Trial',
      '6x training sessions',
    ],
    highlight: false,
    button: { label: 'Get Started', color: 'neutral' as const, variant: 'outline' as const }
  },
  {
    title: '8x',
    price: 'Rp 1.950.000',
    description: 'Our most popular package for comprehensive learning',
    features: [
      'Free Trial',
      '8x training sessions',
    ],
    highlight: true,
    button: { label: 'Most Popular', color: 'warning' as const }
  },
  {
    title: '10x',
    price: 'Rp 2.250.000',
    description: 'Complete mastery with unlimited support',
    features: [
      'Free Trial',
      '10x training sessions',
    ],
    highlight: false,
    button: { label: 'Get Started', color: 'neutral' as const, variant: 'outline' as const }
  },
  {
    title: '12x',
    price: 'Rp 2.650.000',
    description: 'Perfect drivers looking to master driving',
    features: [
      'Free Trial',
      '12x training sessions',
    ],
    highlight: false,
    button: { label: 'Get Started', color: 'neutral' as const, variant: 'outline' as const }
  }
]

// Testimonials
const testimonials = [
  {
    name: 'Sarah Putri',
    role: 'University Student',
    avatar: 'SP',
    content: 'Learning to drive in an EV was so much easier than I expected! No stalling meant I could focus on the road, not the clutch. Highly recommend!'
  },
  {
    name: 'Budi Santoso',
    role: 'Working Professional',
    avatar: 'BS',
    content: 'The instructors are incredibly patient and the Tesla cars are amazing to learn in. Got my license on the first try!'
  },
  {
    name: 'Amanda Chen',
    role: 'Business Owner',
    avatar: 'AC',
    content: 'Premium service from start to finish. The booking system is seamless and the Alam Sutera location is very convenient.'
  }
]

// FAQ items
const faqItems = [
  {
    label: 'Do I need prior driving experience?',
    content: 'No, our courses are designed for complete beginners as well as those looking to improve their skills. EVs are actually easier to learn in!',
    defaultOpen: true
  },
  {
    label: 'What type of electric vehicles do you use?',
    content: 'We use premium electric vehicles including Tesla Model 3 and BYD Atto 3, both equipped with dual-control systems for safe learning.'
  },
  {
    label: 'How do I book my training sessions?',
    content: 'After registering and purchasing a package, you can book your sessions through our online calendar system. Choose from available time slots based on instructor and car availability.'
  },
  {
    label: 'What is the EV Certificate?',
    content: 'Our EV Certificate is a digital certification that validates your proficiency in driving electric vehicles. It covers unique EV aspects like regenerative braking and charging protocols.'
  },
  {
    label: 'Can I reschedule my sessions?',
    content: 'Yes, you can reschedule up to 24 hours before your scheduled session through your member dashboard or by contacting us via WhatsApp.'
  },
  {
    label: 'Do you provide pick-up services?',
    content: 'Training sessions start from our Alam Sutera location. We recommend arranging your own transportation to the training center.'
  }
]

const { slots: globalSlots } = useSchedules()

const timeSlots = computed(() => {
  return globalSlots.value.map(slot => ({
    time: slot.time,
    car: slot.car,
    instructor: slot.instructor,
    available: slot.status === 'available'
  }))
})

const selectedDate = ref(15)
const selectedSlot = ref<string | null>(null)
const currentMonth = 'April 2026'
const weekDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
const calendarDays = [
  { day: 1, available: false },
  { day: 2, available: false },
  { day: 3, available: false },
  { day: 4, available: true },
  { day: 5, available: true },
  { day: 6, available: false },
  { day: 7, available: true },
  { day: 8, available: true },
  { day: 9, available: true },
  { day: 10, available: false },
  { day: 11, available: true },
  { day: 12, available: true },
  { day: 13, available: false },
  { day: 14, available: true },
  { day: 15, available: true },
  { day: 16, available: true },
  { day: 17, available: false },
  { day: 18, available: true },
  { day: 19, available: true },
  { day: 20, available: false },
  { day: 21, available: true },
  { day: 22, available: true },
  { day: 23, available: true },
  { day: 24, available: false },
  { day: 25, available: true },
  { day: 26, available: true },
  { day: 27, available: false },
  { day: 28, available: true },
  { day: 29, available: true },
  { day: 30, available: true }
]
</script>

<template>
  <div>
    <!-- Hero Section -->
    <UPageHero
      title="Master the Road, Drive the Future"
      description="The first premium driving academy in Alam Sutera using 100% Electric Vehicles. Experience smooth, silent, and sustainable learning."
      orientation="horizontal"
      :links="[
        { label: 'Book Your First Session', to: '/register', color: 'warning', icon: 'i-lucide-calendar-check', size: 'lg' },
        { label: 'View Packages', to: '/packages', color: 'neutral', variant: 'outline', trailingIcon: 'i-lucide-arrow-right', size: 'lg' }
      ]"
    >
      <div class="relative w-full aspect-video lg:aspect-[4/3] rounded-2xl overflow-hidden bg-elevated shadow-2xl ring ring-default">
        <img 
          src="https://images.unsplash.com/photo-1593941707882-a5bba14938c7?w=800&auto=format&fit=crop&q=80" 
          alt="Electric Vehicle"
          class="w-full h-full object-cover"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent" />
        <div class="absolute bottom-4 left-4 right-4 flex items-center justify-between">
          <div class="flex items-center gap-2 bg-black/60 backdrop-blur-md rounded-full px-4 py-2 text-white text-md">
            <UIcon name="i-lucide-battery-charging" class="size-5 text-warning" />
            <span>100% Electric</span>
          </div>
          <div class="flex items-center gap-2 bg-black/60 backdrop-blur-md rounded-full px-4 py-2 text-white text-md">
            <UIcon name="i-lucide-shield-check" class="size-5 text-warning" />
            <span>Dual Controls</span>
          </div>
        </div>
      </div>
    </UPageHero>

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

    <!-- Pricing Section -->
    <UPageSection
      id="pricing"
      headline="Pricing"
      title="Choose Your Learning Path"
      description="Flexible packages designed to match your goals and schedule."
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid md:grid-cols-4 gap-6 lg:gap-8">
        <UCard
          v-for="pkg in packages"
          :key="pkg.title"
          :class="[
            'flex flex-col',
            pkg.highlight ? 'ring-2 ring-warning shadow-xl scale-[1.02]' : ''
          ]"
        >
          <template #header>
            <div class="text-center">
              <UBadge v-if="pkg.highlight" label="Most Popular" color="warning" class="mb-2" />
              <h3 class="text-xl font-bold">{{ pkg.title }}</h3>
              <div class="mt-2">
                <span class="text-3xl font-bold text-warning">{{ pkg.price }}</span>
              </div>
              <p class="text-muted text-md mt-2">{{ pkg.description }}</p>
            </div>
          </template>
          
          <ul class="space-y-3 flex-1">
            <li v-for="feature in pkg.features" :key="feature" class="flex items-start gap-2">
              <UIcon name="i-lucide-check" class="size-5 text-warning shrink-0 mt-0.5" />
              <span class="text-md">{{ feature }}</span>
            </li>
          </ul>

          <template #footer>
            <NuxtLink to="/register" class="w-full">
              <UButton v-bind="pkg.button" block />
            </NuxtLink>
          </template>
        </UCard>
      </div>
    </UPageSection>

    <!-- Interactive Booking Preview -->
    <UPageSection
      headline="Easy Booking"
      title="Book Your Sessions with Ease"
      description="Our intuitive booking system shows real-time availability for our electric vehicles."
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="grid lg:grid-cols-2 gap-8 items-start">
        <!-- Calendar -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-calendar" class="size-5 text-warning" />
                <h3 class="font-semibold">Select Date</h3>
              </div>
              <div class="flex items-center gap-2">
                <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" size="md" />
                <span class="text-md font-medium">{{ currentMonth }}</span>
                <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" size="md" />
              </div>
            </div>
          </template>
          
          <!-- Custom Calendar Grid -->
          <div class="grid grid-cols-7 gap-1 mb-2">
            <div v-for="day in weekDays" :key="day" class="text-center text-md font-medium text-muted py-2">
              {{ day }}
            </div>
          </div>
          <div class="grid grid-cols-7 gap-1">
            <!-- Empty cells for days before month starts (April 2026 starts on Wednesday) -->
            <div></div>
            <div></div>
            <div v-for="item in calendarDays" :key="item.day">
              <button
                :disabled="!item.available || item.day < 7"
                :class="[
                  'w-full aspect-square rounded-lg text-md font-medium transition-all',
                  selectedDate === item.day 
                    ? 'bg-primary text-white' 
                    : item.available && item.day >= 7
                      ? 'hover:bg-primary/10 cursor-pointer'
                      : 'text-muted/50 cursor-not-allowed'
                ]"
                @click="item.available && item.day >= 7 && (selectedDate = item.day)"
              >
                {{ item.day }}
              </button>
            </div>
          </div>
          
          <div class="mt-4 flex items-center gap-4 text-md">
            <div class="flex items-center gap-2">
              <div class="size-3 rounded bg-primary/10 border border-primary/30"></div>
              <span class="text-muted">Available</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="size-3 rounded bg-primary"></div>
              <span class="text-muted">Selected</span>
            </div>
          </div>
        </UCard>

        <!-- Time Slots -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-clock" class="size-5 text-warning" />
                <h3 class="font-semibold">Available Slots</h3>
              </div>
              <UBadge :label="`Apr ${selectedDate}`" variant="subtle" />
            </div>
          </template>
          
          <div class="space-y-3">
            <button
              v-for="slot in timeSlots"
              :key="slot.time"
              :disabled="!slot.available"
              :class="[
                'w-full p-4 rounded-lg border transition-all text-left',
                slot.available 
                  ? selectedSlot === slot.time 
                    ? 'border-primary bg-primary/10' 
                    : 'border-default hover:border-primary cursor-pointer'
                  : 'border-default bg-muted/50 opacity-50 cursor-not-allowed'
              ]"
              @click="slot.available && (selectedSlot = slot.time)"
            >
              <div class="flex items-center justify-between">
                <div>
                  <span class="font-semibold">{{ slot.time }}</span>
                  <p class="text-md text-muted">{{ slot.instructor }} - {{ slot.car }}</p>
                </div>
                <UBadge 
                  :label="slot.available ? 'Available' : 'Booked'" 
                  :color="slot.available ? 'success' : 'error'"
                  variant="subtle"
                />
              </div>
            </button>
          </div>

          <template #footer>
            <NuxtLink to="/register" class="w-full">
              <UButton 
                label="Continue to Registration" 
                icon="i-lucide-arrow-right"
                color="warning"
                trailing
                :disabled="!selectedSlot"
                block
              />
            </NuxtLink>
          </template>
        </UCard>
      </div>
    </UPageSection>

    <!-- Location Section -->
    <UPageSection
      id="contact"
      headline="Location"
      title="Conveniently Located in Alam Sutera"
      description="Our training center is strategically located in Alam Sutera, easily accessible from Tangerang and Jakarta."
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid lg:grid-cols-2 gap-8">
        <div class="space-y-6">
          <UCard>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-lg bg-warning/10">
                <UIcon name="i-lucide-map-pin" class="size-6 text-warning" />
              </div>
              <div>
                <h3 class="font-semibold mb-1">Training Center</h3>
                <p class="text-muted text-md">
                  Jl. Alam Sutera Boulevard No. 123<br />
                  Alam Sutera, Tangerang 15143<br />
                  Banten, Indonesia
                </p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-lg bg-warning/10">
                <UIcon name="i-lucide-clock" class="size-6 text-warning" />
              </div>
              <div>
                <h3 class="font-semibold mb-1">Operating Hours</h3>
                <p class="text-muted text-md">
                  Monday - Friday: 08:00 - 17:00<br />
                  Saturday - Sunday: 08:00 - 17:00<br />
                  Night Shift: 18:00 - 20:00
                </p>
              </div>
            </div>
          </UCard>

          <UCard>
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-lg bg-warning/10">
                <UIcon name="i-lucide-phone" class="size-6 text-warning" />
              </div>
              <div>
                <h3 class="font-semibold mb-1">Contact Us</h3>
                <p class="text-muted text-md">
                  Phone: +62 812-3456-7890<br />
                  Email: info@evdriveacademy.id
                </p>
              </div>
            </div>
          </UCard>

          <div class="flex gap-3">
            <NuxtLink to="https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi" target="_blank" class="flex-1">
              <UButton 
                icon="i-simple-icons-whatsapp" 
                label="Chat on WhatsApp"
                class="!bg-[#25D366] hover:!bg-[#128C7E] text-white"
                block
              />
            </NuxtLink>
            <NuxtLink to="https://maps.app.goo.gl/RpSdkpjs4RZg2ZY77" target="_blank" class="flex-1">
              <UButton 
                icon="i-lucide-navigation" 
                label="Get Directions"
                color="neutral"
                variant="outline"
                block
              />
            </NuxtLink>
          </div>
        </div>

        <!-- Map Placeholder -->
        <div class="h-[400px] lg:h-full min-h-[400px] rounded-2xl overflow-hidden bg-elevated border border-default">
          <iframe src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3966.1806061048765!2d106.65588507475077!3d-6.2399118937483635!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69fbc070b4d71d%3A0x8b1a633faf5dbd46!2sALAM%20SUTERA!5e0!3m2!1sen!2sid!4v1776223155011!5m2!1sen!2sid" width="600" height="450" style="border:0;" allowfullscreen="" loading="lazy" referrerpolicy="no-referrer-when-downgrade"></iframe>
        </div>
      </div>
    </UPageSection>

    <!-- Testimonials -->
    <UPageSection
      headline="Testimonials"
      title="What Our Students Say"
      description="Join hundreds of satisfied students who have learned to drive with us."
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <div class="grid md:grid-cols-3 gap-6">
        <UCard v-for="testimonial in testimonials" :key="testimonial.name">
          <div class="flex items-center gap-3 mb-4">
            <UAvatar :text="testimonial.avatar" size="lg" />
            <div>
              <p class="font-semibold">{{ testimonial.name }}</p>
              <p class="text-md text-muted">{{ testimonial.role }}</p>
            </div>
          </div>
          <div class="flex gap-0.5 mb-3">
            <UIcon v-for="i in 5" :key="i" name="i-lucide-star" class="size-4 text-amber-500 fill-amber-500" />
          </div>
          <p class="text-muted">"{{ testimonial.content }}"</p>
        </UCard>
      </div>
    </UPageSection>

    <!-- FAQ Section -->
    <UPageSection
      id="faq"
      headline="FAQ"
      title="Frequently Asked Questions"
      description="Find answers to common questions about our EV driving courses."
      :ui="{ headline: 'text-warning' }"
    >
      <div class="max-w-3xl mx-auto">
        <UAccordion :items="faqItems" />
      </div>
    </UPageSection>

    <!-- CTA Section -->
    <UPageCTA
      title="Ready to Drive the Future?"
      description="Book your first session today and experience the joy of learning in a premium electric vehicle."
      :links="[
        { label: 'Start Your Journey', to: '/register', color: 'warning', icon: 'i-lucide-rocket', size: 'lg' },
        { label: 'View All Packages', to: '/packages', color: 'neutral', variant: 'ghost', trailingIcon: 'i-lucide-arrow-right' }
      ]"
      class="bg-warning/5"
    />
  </div>
</template>
