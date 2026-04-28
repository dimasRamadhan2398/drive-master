<script setup lang="ts">
import { icons } from '@iconify-json/lucide/index.js'
import { ref } from 'vue'

// EV Advantages
const evAdvantages = [
  {
    title: 'One-Pedal Driving',
    description: 'Master the intuitive one-pedal driving experience. Brake and accelerate with a single pedal for smoother control.',
    icon:'i-lucide-disc'
  },
  {
    title: 'No Engine Stalling',
    description: 'Say goodbye to stalling anxiety. Electric motors provide instant torque with zero risk of stalling.',
    icon: 'i-lucide-shield-check'
  },
  {
    title: 'High-Tech Safety',
    description: 'Learn in vehicles equipped with advanced safety features including collision avoidance and lane assist.',
    icon: 'i-lucide-radar'
  },
  {
    title: 'Eco-Friendly Learning',
    description: 'Start your driving journey with zero emissions. Learn to drive while contributing to a sustainable future.',
    icon: 'i-lucide-leaf'
  }
]

// Course packages
const packages = [
  {
    title: 'Free Trial',
    price: 'Rp 0',
    description: 'Perfect for user who want to try our service',
    features: [
      '1 training sessions',
      '60 minutes per session',
      'Basic driving fundamentals',
      'Theory materials included',
    ],
    highlight: false,
    button: { label: 'Try Now', color: 'neutral' as const, variant: 'outline' as const }
  },
  {
    title: 'Starter',
    price: 'Rp 1.500.000',
    description: 'Perfect for beginners looking to get started',
    features: [
      '5 training sessions',
      '45 minutes per session',
      'Basic driving fundamentals',
      'Theory materials included',
      'Certificate of completion'
    ],
    highlight: false,
    button: { label: 'Get Started', color: 'neutral' as const, variant: 'outline' as const }
  },
  {
    title: 'Standard',
    price: 'Rp 2.800.000',
    description: 'Our most popular package for comprehensive learning',
    features: [
      '10 training sessions',
      '60 minutes per session',
      'Advanced driving techniques',
      'Highway & city driving',
      'EV-specific training',
      'Official EV Certificate'
    ],
    highlight: true,
    button: { label: 'Most Popular', color: 'warning' as const }
  },
  {
    title: 'Pro',
    price: 'Rp 4.500.000',
    description: 'Complete mastery with unlimited support',
    features: [
      '15 training sessions',
      '90 minutes per session',
      'Full driving mastery program',
      'Night & weather driving',
      'EV maintenance basics',
      'Premium EV Certificate',
      'Lifetime refresher access'
    ],
    highlight: false,
    button: { label: 'Go Pro', color: 'neutral' as const, variant: 'outline' as const }
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

// FITUR BARU: Logika Kalender Dinamis
const currentDate = ref(new Date('2026-04-10T00:00:00'))
const selectedDate = ref(15)
const selectedSlot = ref<string | null>(null)
const currentMonthStr = computed(() => {
  return currentDate.value.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
})
const currentMonthShortStr = computed(() => {
  return currentDate.value.toLocaleDateString('en-US', { month: 'short' })
})
const weekDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

// FITUR BARU: Kalender dinamis untuk halaman Home
const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  const firstDay = new Date(year, month, 1).getDay()
  const emptyDays = firstDay === 0 ? 6 : firstDay - 1
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  
  const days = []
  for (let i = 0; i < emptyDays; i++) {
    days.push({ day: null, available: false })
  }
  for (let i = 1; i <= daysInMonth; i++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`
    const isAvailable = globalSlots.value.some(s => s.date === dateStr && s.status === 'available')
    days.push({ day: i, available: isAvailable })
  }
  return days
})

function changeMonth(offset: number) {
  const newDate = new Date(currentDate.value)
  newDate.setMonth(newDate.getMonth() + offset)
  currentDate.value = newDate
}

const timeSlots = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = String(currentDate.value.getMonth() + 1).padStart(2, '0')
  const day = String(selectedDate.value).padStart(2, '0')
  const dateStr = `${year}-${month}-${day}`
  
  return globalSlots.value
    .filter(slot => slot.date === dateStr)
    .map(slot => ({
      time: slot.time,
      car: slot.car,
      available: slot.status === 'available'
    }))
})
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
          <div class="flex items-center gap-2 bg-black/60 backdrop-blur-sm rounded-full px-4 py-2 text-white text-sm">
            <UIcon name="i-lucide-battery-charging" class="size-5 text-warning" />
            <span>100% Electric</span>
          </div>
          <div class="flex items-center gap-2 bg-black/60 backdrop-blur-sm rounded-full px-4 py-2 text-white text-sm">
            <UIcon name="i-lucide-shield-check" class="size-5 text-warning" />
            <span>Dual Controls</span>
          </div>
        </div>
      </div>
    </UPageHero>

    <!-- EV Advantages Section -->
    <UPageSection
      id="about"
      headline="Why Learn in an EV?"
      title="The Electric Advantage"
      description="Discover why learning to drive in an electric vehicle gives you a head start over traditional driving schools."
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <UPageGrid>
        <UPageCard
          v-for="advantage in evAdvantages"
          :key="advantage.title"
          :icon="advantage.icon"
          :title="advantage.title"
          :description="advantage.description"
          :ui="{ leadingIcon: 'text-warning text-3xl' }"
        />
      </UPageGrid>
    </UPageSection>

    <!-- Premium Fleet Section -->
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
              <p class="text-muted text-sm mt-2">{{ pkg.description }}</p>
            </div>
          </template>
          
          <ul class="space-y-3 flex-1">
            <li v-for="feature in pkg.features" :key="feature" class="flex items-start gap-2">
              <UIcon name="i-lucide-check" class="size-5 text-warning shrink-0 mt-0.5" />
              <span class="text-sm">{{ feature }}</span>
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
                <UButton icon="i-lucide-chevron-left" variant="ghost" color="neutral" size="xs" @click="changeMonth(-1)" />
                <span class="text-sm font-medium">{{ currentMonthStr }}</span>
                <UButton icon="i-lucide-chevron-right" variant="ghost" color="neutral" size="xs" @click="changeMonth(1)" />
              </div>
            </div>
          </template>
          
          <!-- Custom Calendar Grid -->
          <div class="grid grid-cols-7 gap-1 mb-2">
            <div v-for="day in weekDays" :key="day" class="text-center text-xs font-medium text-muted py-2">
              {{ day }}
            </div>
          </div>
          <div class="grid grid-cols-7 gap-1">
            <!-- PERUBAHAN: Kalender dinamis untuk Home Page -->
            <div v-for="(item, index) in calendarDays" :key="index">
              <button
                v-if="item.day !== null"
                :disabled="!item.available"
                :class="[
                  'w-full aspect-square rounded-lg text-sm font-medium transition-all',
                  selectedDate === item.day 
                    ? 'bg-warning text-white' 
                    : item.available 
                      ? 'hover:bg-warning/10 cursor-pointer'
                      : 'text-muted/50 cursor-not-allowed'
                ]"
                @click="item.available && (selectedDate = item.day)"
              >
                {{ item.day }}
              </button>
            </div>
          </div>
          
          <div class="mt-4 flex items-center gap-4 text-xs">
            <div class="flex items-center gap-2">
              <div class="size-3 rounded bg-warning/10 border border-warning/30"></div>
              <span class="text-muted">Available</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="size-3 rounded bg-warning"></div>
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
              <!-- PERUBAHAN: Memperbaiki teks bulan statis -->
              <UBadge :label="`${currentMonthShortStr} ${selectedDate}`" variant="subtle" />
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
                    ? 'border-warning bg-warning/10' 
                    : 'border-default hover:border-warning cursor-pointer'
                  : 'border-default bg-muted/50 opacity-50 cursor-not-allowed'
              ]"
              @click="slot.available && (selectedSlot = slot.time)"
            >
              <div class="flex items-center justify-between">
                <div>
                  <span class="font-semibold">{{ slot.time }}</span>
                  <p class="text-sm text-muted">{{ slot.car }}</p>
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
                <p class="text-muted text-sm">
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
                <p class="text-muted text-sm">
                  Monday - Saturday: 08:00 - 18:00<br />
                  Sunday: Closed
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
                <p class="text-muted text-sm">
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
              <p class="text-sm text-muted">{{ testimonial.role }}</p>
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

    <!-- Member Dashboard Preview -->
    <UPageSection
      headline="Member Area"
      title="Track Your Progress"
      description="Get access to your personalized dashboard where you can view sessions, track progress, and download certificates."
      :ui="{ headline: 'text-warning' }"
    >
      <UCard class="max-w-3xl mx-auto overflow-hidden">
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <UAvatar text="JD" />
              <div>
                <p class="font-semibold">Welcome back, John!</p>
                <p class="text-sm text-muted">Standard Package</p>
              </div>
            </div>
            <UBadge label="Active" color="success" />
          </div>
        </template>

        <div class="grid sm:grid-cols-3 gap-4 mb-6">
          <div class="text-center p-4 rounded-lg bg-muted/50">
            <p class="text-3xl font-bold text-primary">4</p>
            <p class="text-sm text-muted">Sessions Completed</p>
          </div>
          <div class="text-center p-4 rounded-lg bg-muted/50">
            <p class="text-3xl font-bold text-warning">6</p>
            <p class="text-sm text-muted">Sessions Remaining</p>
          </div>
          <div class="text-center p-4 rounded-lg bg-muted/50">
            <p class="text-3xl font-bold text-info">40%</p>
            <p class="text-sm text-muted">Course Progress</p>
          </div>
        </div>

        <UProgress :value="40" color="warning" class="mb-6" />

        <div class="flex flex-wrap gap-3">
          <NuxtLink to="/dashboard">
            <UButton label="Go to Dashboard" color="warning" icon="i-lucide-layout-dashboard" />
          </NuxtLink>
          <UButton label="Book Next Session" icon="i-lucide-calendar-plus" variant="outline" color="neutral" />
          <UButton label="Download Certificate" icon="i-lucide-download" variant="ghost" color="neutral" disabled />
        </div>
      </UCard>
    </UPageSection>
  </div>
</template>
