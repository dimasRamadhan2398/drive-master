<script setup lang="ts">
import { ref, computed } from 'vue'

useSeoMeta({
  title: 'Home | Drive Master Academy',
  description: 'Drive Master Academy offers comprehensive EV and manual driving courses in Alam Sutera with expert instructors and certified programs.',
})

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

const selectedPlan = ref<
'six_package' | 
'six_package_night' | 
'six_package_weekend' | 
'six_package_weekend_night' | 
'eight_package' | 
'eight_package_night' | 
'eight_package_weekend' | 
'eight_package_weekend_night' | 
'ten_package' | 
'ten_package_night' | 
'ten_package_weekend' | 
'ten_package_weekend_night' | 
'twelve_package' | 
'twelve_package_night' | 
'twelve_package_weekend' | 
'twelve_package_weekend_night'>('eight_package')

const now = ref(new Date())
if (process.client) {
  setInterval(() => {
    now.value = new Date()
  }, 1000)
}

const isPromoActive = computed(() => {
  const promoEnd = new Date('2026-05-20T23:59:59')
  return now.value < promoEnd
})

const currentPlan = computed(() => plans.find(p => p.id === selectedPlan.value))
const discount = computed(() => {
  if (!currentPlan.value || !isPromoActive.value) return 0
  return Math.ceil(((currentPlan.value.originalPrice - currentPlan.value.promoPrice) / currentPlan.value.originalPrice) * 100)
})

const timeLeft = computed(() => {
  // Target: Akhir dari jam saat ini (Top of the next hour)
  const nextHour = new Date(now.value)
  nextHour.setHours(nextHour.getHours() + 1, 0, 0, 0)
  
  const diff = nextHour.getTime() - now.value.getTime()

  if (diff <= 0) {
    return { hours: 0, minutes: 0, seconds: 0 }
  }

  const totalSeconds = Math.floor(diff / 1000)
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60

  return { hours, minutes, seconds }
})

const plans = [
  {
    id: 'six_package',
    name: '6x',
    originalPrice: 1950000,
    promoPrice: 1750000,
    duration: '3 Months',
    sessions: 6,
    features: [
      'Free Trial',
      '6 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'six_package_night',
    name: '6x + Night Session',
    originalPrice: 2050000,
    promoPrice: 1850000,
    duration: '3 Months',
    sessions: 6,
    features: [
      'Free Trial',
      '6 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'six_package_weekend',
    name: '6x + Weekend Session',
    originalPrice: 2050000,
    promoPrice: 1850000,
    duration: '3 Months',
    sessions: 6,
    features: [
      'Free Trial',
      '6 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'six_package_weekend_night',
    name: '6x + Weekend & Night Session',
    originalPrice: 2150000,
    promoPrice: 1950000,
    duration: '3 Months',
    sessions: 6,
    features: [
      'Free Trial',
      '6 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'eight_package',
    name: '8x',
    originalPrice: 2150000,
    promoPrice: 1950000,
    duration: '3 Months',
    sessions: 8,
    features: [
      'Free Trial',
      '8 training sessions',
      'SIM A'
    ],
    highlight: true
  },
  {
    id: 'eight_package_night',
    name: '8x + Night Session',
    originalPrice: 2300000,
    promoPrice: 2100000,
    duration: '3 Months',
    sessions: 8,
    features: [
      'Free Trial',
      '8 training sessions',
      'SIM A'
    ],
    highlight: true
  },
  {
    id: 'eight_package_weekend',
    name: '8x + Weekend Session',
    originalPrice: 2300000,
    promoPrice: 2100000,
    duration: '3 Months',
    sessions: 8,
    features: [
      'Free Trial',
      '8 training sessions',
      'SIM A'
    ],
    highlight: true
  },
  {
    id: 'eight_package_weekend_night',
    name: '8x + Weekend & Night Session',
    originalPrice: 2450000,
    promoPrice: 2250000,
    duration: '3 Months',
    sessions: 8,
    features: [
      'Free Trial',
      '8 training sessions',
      'SIM A'
    ],
    highlight: true
  },
  {
    id: 'ten_package',
    name: '10x',
    originalPrice: 2500000,
    promoPrice: 2250000,
    duration: '3 Months',
    sessions: 10,
    features: [
      'Free Trial',
      '10 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'ten_package_night',
    name: '10x + Night Session',
    originalPrice: 2700000,
    promoPrice: 2450000,
    duration: '3 Months',
    sessions: 10,
    features: [
      'Free Trial',
      '10 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'ten_package_weekend',
    name: '10x + Weekend Session',
    originalPrice: 2700000,
    promoPrice: 2450000,
    duration: '3 Months',
    sessions: 10,
    features: [
      'Free Trial',
      '10 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'ten_package_weekend_night',
    name: '10x + Weekend & Night Session',
    originalPrice: 2900000,
    promoPrice: 2650000,
    duration: '3 Months',
    sessions: 10,
    features: [
      'Free Trial',
      '10 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'twelve_package',
    name: '12x',
    originalPrice: 2950000,
    promoPrice: 2650000,
    duration: '3 Months',
    sessions: 12,
    features: [
      'Free Trial',
      '12 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'twelve_package_night',
    name: '12x + Night Session',
    originalPrice: 3250000,
    promoPrice: 2900000,
    duration: '3 Months',
    sessions: 12,
    features: [
      'Free Trial',
      '12 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'twelve_package_weekend',
    name: '12x + Weekend Session',
    originalPrice: 3250000,
    promoPrice: 2900000,
    duration: '3 Months',
    sessions: 12,
    features: [
      'Free Trial',
      '12 training sessions',
      'SIM A'
    ],
    highlight: false
  },
  {
    id: 'twelve_package_weekend_night',
    name: '12x + Weekend & Night Session',
    originalPrice: 3500000,
    promoPrice: 3150000,
    duration: '3 Months',
    sessions: 12,
    features: [
      'Free Trial',
      '12 training sessions',
      'SIM A'
    ],
    highlight: false
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
    content: 'The instructors are incredibly patient and the BYD Atto 1 cars are amazing to learn in. Got my license on the first try!'
  },
  {
    name: 'Amanda Chen',
    role: 'Business Owner',
    avatar: 'AC',
    content: 'Premium service from start to finish. The booking system is seamless and the Alam Sutera location is very convenient.'
  }
]

// Dynamic content from admin
const { faqs, pages } = useContent()

const homePage = computed(() => pages.value.find(p => p.slug === '/' && p.status === 'published'))

// FAQ items
const faqItems = computed(() => {
  return faqs.value.map((faq, index) => ({
    label: faq.question,
    content: faq.answer,
    defaultOpen: index === 0
  }))
})

const { slots: globalSlots } = useSchedules()

const selectedDate = ref(15)
const selectedSlot = ref<string | null>(null)
const currentDate = ref(new Date('2026-04-10T00:00:00'))

const currentMonth = computed(() => {
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
      instructor: slot.instructor,
      available: slot.status === 'available'
    }))
})

const instructors = [
  {
    name: 'Mr. Ahmad',
    role: 'Senior Instructor',
    specialization: 'Defensive Driving',
    image: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&auto=format&fit=crop&q=80'
  },
  {
    name: 'Ms. Sari',
    role: 'Professional Instructor',
    specialization: 'Urban Driving',
    image: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=400&auto=format&fit=crop&q=80'
  },
  {
    name: 'Mr. Budi',
    role: 'Safety Specialist',
    specialization: 'Highway Safety',
    image: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400&auto=format&fit=crop&q=80'
  }
]
</script>

<template>
  <div>
    <!-- Dynamic Admin Sections for Home Page -->
    <template v-if="homePage && homePage.sections.length > 0">
      <ContentSectionRenderer 
        v-for="section in homePage.sections" 
        :key="section.id" 
        :section="{ type: section.type, data: section }" 
      />
    </template>

    <!-- Hero Section -->
    <UPageHero
      title="Master the Road, Drive the Future"
      description="The first premium driving academy in Alam Sutera using 100% Electric Vehicles. Experience smooth, silent, and sustainable learning."
      orientation="horizontal"
      :links="[
        { label: 'Book Your First Session', to: '/auth/register', color: 'warning', icon: 'i-lucide-calendar-check', size: 'lg' },
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
      <!-- Premium Promo Banner with Real-time Countdown -->
      <Transition
        enter-active-class="transition duration-500 ease-out"
        enter-from-class="opacity-0 translate-y-4"
        enter-to-class="opacity-100 translate-y-0"
      >
        <div 
          v-if="isPromoActive"
          class="mb-12 relative group"
        >
          <!-- Background Glow Effect -->
          <div class="absolute -inset-1 bg-gradient-to-r from-warning-500 to-orange-600 rounded-3xl blur opacity-25 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"></div>
          
          <div class="relative bg-white dark:bg-gray-900 border border-warning-500/20 rounded-2xl overflow-hidden shadow-2xl">
            <div class="flex flex-col lg:flex-row">
              <!-- Left Side: Promo Info -->
              <div class="flex-1 p-6 md:p-8 flex items-center gap-6">
                <div class="relative">
                  <div class="size-20 rounded-2xl bg-warning-500 flex items-center justify-center shadow-lg shadow-warning-500/30 animate-pulse">
                    <UIcon name="i-lucide-zap" class="size-10 text-white" />
                  </div>
                  <div class="absolute -top-2 -right-2 bg-red-600 text-white text-[10px] font-black px-2 py-1 rounded-full uppercase tracking-tighter">Hot</div>
                </div>
                
                <div>
                  <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-warning-500/10 text-warning-600 text-xs font-bold uppercase tracking-widest mb-3">
                    Flash Sale Active
                  </div>
                  <h2 class="text-2xl md:text-3xl font-black tracking-tight mb-2">
                    Special Price: <span class="text-warning-500">Save up to {{ discount }}%</span>
                  </h2>
                  <p class="text-muted text-sm md:text-base max-w-lg">
                    Exclusive discount for new members. Start your journey today with our professional instructors.
                  </p>
                </div>
              </div>

              <!-- Right Side: Countdown Timer -->
              <div class="lg:w-[320px] bg-warning-50 dark:bg-warning-500/5 p-6 md:p-8 flex flex-col items-center justify-center border-t lg:border-t-0 lg:border-l border-warning-500/10">
                <p class="text-xs font-bold text-warning-600 uppercase tracking-[0.2em] mb-4">Promo Ends In:</p>
                
                <div class="flex gap-3">
                  <div v-for="(val, unit) in { hours: timeLeft.hours, mins: timeLeft.minutes, secs: timeLeft.seconds }" :key="unit" class="text-center">
                    <div class="size-14 md:size-16 bg-white dark:bg-gray-800 border border-warning-500/20 rounded-xl shadow-inner flex items-center justify-center mb-1">
                      <span class="text-2xl font-black text-foreground">{{ String(val).padStart(2, '0') }}</span>
                    </div>
                    <span class="text-[10px] font-bold text-muted uppercase">{{ unit }}</span>
                  </div>
                </div>
              </div>
              
            </div>
          </div>
        </div>
      </Transition>

      <!-- Plan Cards -->
      <div class="grid md:grid-cols-4 gap-6 mb-12">
        <div 
          v-for="plan in plans" 
          :key="plan.id"
          class="relative"
        >
          <input 
            :id="`plan-${plan.id}`"
            v-model="selectedPlan"
            type="radio" 
            :value="plan.id"
            class="sr-only"
          />
          <label 
            :for="`plan-${plan.id}`"
            class="block h-full cursor-pointer"
          >
            <UCard
              :class="[
                'h-full flex flex-col transition-all',
                selectedPlan === plan.id 
                  ? 'ring-2 ring-warning shadow-xl' 
                  : 'hover:shadow-lg'
              ]"
            >
              <div v-if="plan.highlight" class="absolute -top-3 left-1/2 -translate-x-1/2 z-10">
                <UBadge label="Most Popular" color="warning" />
              </div>

              <template #header>
                <div class="text-center">
                  <h3 class="text-2xl font-bold">{{ plan.name }}</h3>
                  <p class="text-muted text-sm mt-2">Package Duration : {{ plan.duration }}</p>
                </div>
              </template>

              <div class="flex-1 space-y-6">
                <!-- Pricing -->
                <div class="text-center">
                  <div v-if="isPromoActive && plan.promoPrice < plan.originalPrice" class="space-y-2">
                    <p class="text-sm text-muted line-through">
                      Rp {{ plan.originalPrice.toLocaleString('id-ID') }}
                    </p>
                    <p class="text-4xl font-bold text-warning">
                      Rp {{ plan.promoPrice.toLocaleString('id-ID') }}
                    </p>
                    <p class="text-xs text-green-600 font-medium">
                      Save Rp {{ (plan.originalPrice - plan.promoPrice).toLocaleString('id-ID') }}
                    </p>
                  </div>
                  <div v-else>
                    <p class="text-4xl font-bold">
                      Rp {{ plan.originalPrice.toLocaleString('id-ID') }}
                    </p>
                  </div>
                </div>

                <!-- Sessions Badge -->
                <div class="text-center">
                  <UBadge :label="`${plan.sessions} Sessions`" color="warning" variant="subtle" />
                </div>

                <!-- Features -->
                <ul class="space-y-3">
                  <li 
                    v-for="feature in plan.features" 
                    :key="feature"
                    class="flex items-start gap-3"
                  >
                    <UIcon name="i-lucide-check" class="size-5 text-warning shrink-0 mt-0.5" />
                    <span class="text-sm">{{ feature }}</span>
                  </li>
                </ul>
              </div>

              <!-- Selection Indicator -->
              <div v-if="selectedPlan === plan.id" class="pt-4">
                <div class="flex flex-col gap-3">
                  <NuxtLink :to="{ path: '/auth/register', query: { plan: plan.id } }">
                    <UButton label="Choose this Plan" color="warning" block />
                  </NuxtLink>
                </div>
              </div>
            </UCard>
          </label>
        </div>
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
            <!-- PERUBAHAN: Kalender dinamis untuk Home Page -->
            <div v-for="(item, index) in calendarDays" :key="index">
              <button
                v-if="item.day !== null"
                :disabled="!item.available"
                :class="[
                  'w-full aspect-square rounded-lg text-md font-medium transition-all',
                  selectedDate === item.day 
                    ? 'bg-primary text-white' 
                    : item.available && item.day >= 7
                      ? 'hover:bg-primary/10 cursor-pointer'
                      : 'text-muted/50 cursor-not-allowed'
                ]"
                @click="item.available && (selectedDate = item.day)"
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
            <NuxtLink :to="{ path: '/auth/register', query: { plan: selectedPlan } }" class="w-full">
              <UButton 
                label="Choose this Schedule" 
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

    <!-- Instructors Section -->
    <UPageSection
      headline="Our Instructor"
      title="Meet Our Professional Instructors"
      description="Learn from certified experts who are passionate about teaching and safety."
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid md:grid-cols-3 gap-8">
        <UCard 
          v-for="instructor in instructors" 
          :key="instructor.name"
          class="overflow-hidden group text-center"
          :ui="{ body: 'p-0' }"
        >
          <div class="relative h-64 overflow-hidden">
            <img 
              :src="instructor.image" 
              :alt="instructor.name"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          </div>
          
          <div class="p-6">
            <h3 class="text-xl font-bold">{{ instructor.name }}</h3>
            <p class="text-warning font-medium text-sm mb-4">{{ instructor.role }}</p>
            <p class="text-muted text-sm">{{ instructor.specialization }}</p>
          </div>
        </UCard>
      </div>

      <div class="mt-12 text-center">
        <NuxtLink to="/instructors">
          <UButton 
            label="View All Instructors" 
            color="warning" 
            variant="outline" 
            size="lg" 
            trailing-icon="i-lucide-arrow-right" 
          />
        </NuxtLink>
      </div>
    </UPageSection>

    <!-- Location Section -->
    <UPageSection
      id="contact"
      headline="Location"
      title="Conveniently Located in Alam Sutera"
      description="Our training center is strategically located in Alam Sutera, easily accessible from Tangerang and Jakarta."
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
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
        <div class="h-auto lg:h-full min-h-auto rounded-2xl overflow-hidden bg-elevated border border-default">
          <iframe src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3966.1806061048765!2d106.65588507475077!3d-6.2399118937483635!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69fbc070b4d71d%3A0x8b1a633faf5dbd46!2sALAM%20SUTERA!5e0!3m2!1sen!2sid!4v1776223155011!5m2!1sen!2sid" width="600" height="550" style="border:0;" allowfullscreen loading="lazy" referrerpolicy="no-referrer-when-downgrade"></iframe>
        </div>
      </div>
    </UPageSection>

    <!-- Testimonials -->
    <UPageSection
      headline="Testimonials"
      title="What Our Students Say"
      description="Join hundreds of satisfied students who have learned to drive with us."
      :ui="{ headline: 'text-warning' }"
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
      class="bg-muted/30"
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
        { label: 'Start Your Journey', to: '/auth/register', color: 'warning', icon: 'i-lucide-rocket', size: 'lg' },
        { label: 'View All Packages', to: '/packages', color: 'neutral', variant: 'ghost', trailingIcon: 'i-lucide-arrow-right' }
      ]"
      class="bg-warning/5"
    />
  </div>
</template>
