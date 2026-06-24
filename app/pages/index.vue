<script setup lang="ts">
import { ref, computed } from "vue";
const { t, locale } = useI18n();

useSeoMeta({
  title: t('nav.home') + " | Drive Master Academy",
  description: t('home.hero.description'),
});

// Course Material with i18n
const courseMaterial = computed(() => [
  {
    title: t("home.material.materialTheory"),
    description: [
      t("home.material.materialTheoryDesc", ""),
      t("home.material.materialTheoryDesc2", ""),
      t("home.material.materialTheoryDesc3", ""),
      t("home.material.materialTheoryDesc4", ""),
    ],
    icon: "i-lucide-book-open",
  },
  {
    title: t("home.material.initialControl"),
    description: [
      t("home.material.initialControlDesc", ""),
      t("home.material.initialControlDesc2", ""),
      t("home.material.initialControlDesc3", ""),
    ],
    icon: "i-lucide-shield-check",
  },
  {
    title: t("home.material.basicManeuvering"),
    description: [
      t("home.material.basicManeuveringDesc", ""),
      t("home.material.basicManeuveringDesc2", ""),
      t("home.material.basicManeuveringDesc3", ""),
    ],
    icon: "i-lucide-radar",
  },
  {
    title: t("home.material.uphillDownhill"),
    description: [
      t("home.material.uphillDownhillDesc", ""),
      t("home.material.uphillDownhillDesc2", ""),
    ],
    icon: "i-lucide-car",
  },
  {
    title: t("home.material.parking"),
    description: [t("home.material.parkingDesc", ""), t("home.material.parkingDesc2", "")],
    icon: "i-lucide-car",
  },
  {
    title: t("home.material.highway"),
    description: [
      t("home.material.highwayDesc", ""),
      t("home.material.highwayDesc2", ""),
      t("home.material.highwayDesc3", ""),
    ],
    icon: "i-lucide-car",
  },
]);

const selectedPlan = ref<
  | "six_package"
  | "six_package_night"
  | "six_package_weekend"
  | "six_package_weekend_night"
  | "eight_package"
  | "eight_package_night"
  | "eight_package_weekend"
  | "eight_package_weekend_night"
  | "ten_package"
  | "ten_package_night"
  | "ten_package_weekend"
  | "ten_package_weekend_night"
  | "twelve_package"
  | "twelve_package_night"
  | "twelve_package_weekend"
  | "twelve_package_weekend_night"
>("eight_package");

const now = ref(new Date());
if (process.client) {
  setInterval(() => {
    now.value = new Date();
  }, 1000);
}

const isPromoActive = computed(() => {
  // Promo diubah ke tanggal di masa depan (misalnya 2026-12-31) agar banner muncul
  const promoEnd = new Date("2026-12-31T23:59:59");
  return now.value < promoEnd;
});

const currentPlan = computed(() =>
  plans.find((p) => p.id === selectedPlan.value),
);
const discount = computed(() => {
  if (!currentPlan.value || !isPromoActive.value) return 0;
  return Math.ceil(
    ((currentPlan.value.originalPrice - currentPlan.value.promoPrice) /
      currentPlan.value.originalPrice) *
      100,
  );
});

const translatedSaveAmount = computed(() => {
  return t('home.saveUpTo', { discount: discount.value });
});

const timeLeft = computed(() => {
  // Target: Akhir dari jam saat ini (Top of the next hour)
  const nextHour = new Date(now.value);
  nextHour.setHours(nextHour.getHours() + 1, 0, 0, 0);

  const diff = nextHour.getTime() - now.value.getTime();

  if (diff <= 0) {
    return { hours: 0, minutes: 0, seconds: 0 };
  }

  const totalSeconds = Math.floor(diff / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  return { hours, minutes, seconds };
});

const plans = [
  {
    id: "six_package",
    name: "6x",
    originalPrice: 1950000,
    promoPrice: 1750000,
    duration: "3 Months",
    sessions: 6,
    features: ["Free Trial", "6 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "six_package_night",
    name: "6x + Night Session",
    originalPrice: 2050000,
    promoPrice: 1850000,
    duration: "3 Months",
    sessions: 6,
    features: ["Free Trial", "6 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "six_package_weekend",
    name: "6x + Weekend Session",
    originalPrice: 2050000,
    promoPrice: 1850000,
    duration: "3 Months",
    sessions: 6,
    features: ["Free Trial", "6 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "six_package_weekend_night",
    name: "6x + Weekend & Night Session",
    originalPrice: 2150000,
    promoPrice: 1950000,
    duration: "3 Months",
    sessions: 6,
    features: ["Free Trial", "6 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "eight_package",
    name: "8x",
    originalPrice: 2150000,
    promoPrice: 1950000,
    duration: "3 Months",
    sessions: 8,
    features: ["Free Trial", "8 training sessions", "SIM A"],
    highlight: true,
  },
  {
    id: "eight_package_night",
    name: "8x + Night Session",
    originalPrice: 2300000,
    promoPrice: 2100000,
    duration: "3 Months",
    sessions: 8,
    features: ["Free Trial", "8 training sessions", "SIM A"],
    highlight: true,
  },
  {
    id: "eight_package_weekend",
    name: "8x + Weekend Session",
    originalPrice: 2300000,
    promoPrice: 2100000,
    duration: "3 Months",
    sessions: 8,
    features: ["Free Trial", "8 training sessions", "SIM A"],
    highlight: true,
  },
  {
    id: "eight_package_weekend_night",
    name: "8x + Weekend & Night Session",
    originalPrice: 2450000,
    promoPrice: 2250000,
    duration: "3 Months",
    sessions: 8,
    features: ["Free Trial", "8 training sessions", "SIM A"],
    highlight: true,
  },
  {
    id: "ten_package",
    name: "10x",
    originalPrice: 2500000,
    promoPrice: 2250000,
    duration: "3 Months",
    sessions: 10,
    features: ["Free Trial", "10 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "ten_package_night",
    name: "10x + Night Session",
    originalPrice: 2700000,
    promoPrice: 2450000,
    duration: "3 Months",
    sessions: 10,
    features: ["Free Trial", "10 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "ten_package_weekend",
    name: "10x + Weekend Session",
    originalPrice: 2700000,
    promoPrice: 2450000,
    duration: "3 Months",
    sessions: 10,
    features: ["Free Trial", "10 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "ten_package_weekend_night",
    name: "10x + Weekend & Night Session",
    originalPrice: 2900000,
    promoPrice: 2650000,
    duration: "3 Months",
    sessions: 10,
    features: ["Free Trial", "10 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "twelve_package",
    name: "12x",
    originalPrice: 2950000,
    promoPrice: 2650000,
    duration: "3 Months",
    sessions: 12,
    features: ["Free Trial", "12 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "twelve_package_night",
    name: "12x + Night Session",
    originalPrice: 3250000,
    promoPrice: 2900000,
    duration: "3 Months",
    sessions: 12,
    features: ["Free Trial", "12 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "twelve_package_weekend",
    name: "12x + Weekend Session",
    originalPrice: 3250000,
    promoPrice: 2900000,
    duration: "3 Months",
    sessions: 12,
    features: ["Free Trial", "12 training sessions", "SIM A"],
    highlight: false,
  },
  {
    id: "twelve_package_weekend_night",
    name: "12x + Weekend & Night Session",
    originalPrice: 3500000,
    promoPrice: 3150000,
    duration: "3 Months",
    sessions: 12,
    features: ["Free Trial", "12 training sessions", "SIM A"],
    highlight: false,
  },
];

// Testimonials
const testimonials = [
  {
    name: "Sarah Putri",
    role: "University Student",
    avatar: "SP",
    content:
      "Learning to drive in an EV was so much easier than I expected! No stalling meant I could focus on the road, not the clutch. Highly recommend!",
  },
  {
    name: "Budi Santoso",
    role: "Working Professional",
    avatar: "BS",
    content:
      "The instructors are incredibly patient and the BYD Atto 1 cars are amazing to learn in. Got my license on the first try!",
  },
  {
    name: "Amanda Chen",
    role: "Business Owner",
    avatar: "AC",
    content:
      "Premium service from start to finish. The booking system is seamless and the Alam Sutera location is very convenient.",
  },
];

// Dynamic content from admin
const { faqs, pages } = useContent();

const homePage = computed(() =>
  pages.value.find((p) => p.slug === "/" && p.status === "published"),
);

// FAQ items
const faqItems = computed(() => {
  return faqs.value.map((faq, index) => ({
    label: faq.question,
    content: faq.answer,
    defaultOpen: index === 0,
  }));
});

// Schedules - Use store to fetch from API
const schedulesStore = useSchedulesStore();

// Loading state for schedules
const schedulesLoading = ref(false);

// Helper to format date to YYYY-MM-DD
const formatDateString = (date: Date): string => {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  return `${year}-${month}-${day}`;
};

// Today's date
const todayDate = formatDateString(new Date());

// Calendar state - use today's date
const selectedDate = ref(new Date().getDate());
const selectedSlot = ref<string | null>(null);
const currentDate = ref(new Date());

// Fetch schedules when date changes
const fetchSchedulesForDate = async (dateStr: string) => {
  schedulesLoading.value = true;
  try {
    await schedulesStore.fetchByDate(dateStr);
  } finally {
    schedulesLoading.value = false;
  }
};

// Watch for calendar date changes and fetch schedules
watch(
  () => [currentDate.value.getMonth(), currentDate.value.getFullYear()],
  async () => {
    // Fetch today's schedules when month/year changes
    await fetchSchedulesForDate(todayDate);
  },
);

// Also watch for selectedDate day changes
watch(selectedDate, async (newDay) => {
  const dateStr = `${currentDate.value.getFullYear()}-${String(currentDate.value.getMonth() + 1).padStart(2, "0")}-${String(newDay).padStart(2, "0")}`;
  await fetchSchedulesForDate(dateStr);
});

// Calendar navigation
function changeMonth(offset: number) {
  const newDate = new Date(currentDate.value);
  newDate.setMonth(newDate.getMonth() + offset);
  currentDate.value = newDate;
  // Reset selected date to 1 when changing month to avoid invalid dates
  selectedDate.value = 1;
}

// Map store slots to TimeSlot format for UI compatibility
const globalSlots = computed(() => schedulesStore.slots);

const currentMonth = computed(() => {
  return currentDate.value.toLocaleDateString(locale.value === 'id' ? 'id-ID' : 'en-US', {
    month: "long",
    year: "numeric",
  });
});
const currentMonthShortStr = computed(() => {
  return currentDate.value.toLocaleDateString(locale.value === 'id' ? 'id-ID' : 'en-US', { month: "short" });
});
const weekDays = computed(() => {
  return locale.value === 'id'
    ? ["Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"]
    : ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
});

// FITUR BARU: Kalender dinamis untuk halaman Home
const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear();
  const month = currentDate.value.getMonth();
  const firstDay = new Date(year, month, 1).getDay();
  const emptyDays = firstDay === 0 ? 6 : firstDay - 1;
  const daysInMonth = new Date(year, month + 1, 0).getDate();

  const days = [];
  for (let i = 0; i < emptyDays; i++) {
    days.push({ day: null, available: false });
  }
  for (let i = 1; i <= daysInMonth; i++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, "0")}-${String(i).padStart(2, "0")}`;
    const isAvailable = globalSlots.value.some(
      (s) => s.date === dateStr && s.status === "available",
    );
    days.push({ day: i, available: isAvailable });
  }
  return days;
});

const timeSlots = computed(() => {
  const year = currentDate.value.getFullYear();
  const month = String(currentDate.value.getMonth() + 1).padStart(2, "0");
  const day = String(selectedDate.value).padStart(2, "0");
  const dateStr = `${year}-${month}-${day}`;

  return globalSlots.value
    .filter((slot) => slot.date === dateStr)
    .map((slot) => ({
      time: slot.time,
      car: slot.car,
      instructor: slot.instructor,
      available: slot.status === "available",
    }));
});

const instructorsStore = useInstructorsStore();
const instructorsList = computed(() => instructorsStore.instructors);

const testimonialsStore = useTestimonialsStore();
const testimonialsList = computed(() => testimonialsStore.testimonials);

const packagesStore = usePackagesStore();
const packagesList = computed(() => packagesStore.packages);

onMounted(async () => {
  instructorsStore.fetchInstructors();
  testimonialsStore.fetchTestimonials();
  packagesStore.fetchPackages();
  // Fetch today's schedules from API
  await fetchSchedulesForDate(todayDate);
});
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
      :title="t('home.hero.title')"
      :description="t('home.hero.description')"
      orientation="horizontal"
      :links="[
        {
          label: t('home.bookFirstSession'),
          to: '/auth/register',
          color: 'warning',
          icon: 'i-lucide-calendar-check',
          size: 'lg',
        },
        {
          label: t('home.viewPackages'),
          to: '/packages',
          color: 'neutral',
          variant: 'outline',
          trailingIcon: 'i-lucide-arrow-right',
          size: 'lg',
        },
      ]"
    >
      <div
        class="relative w-full aspect-video lg:aspect-[4/3] rounded-2xl overflow-hidden bg-elevated shadow-2xl ring ring-default"
      >
        <img
          src="https://images.unsplash.com/photo-1593941707882-a5bba14938c7?w=800&auto=format&fit=crop&q=80"
          alt="Electric Vehicle"
          class="w-full h-full object-cover"
        />
        <div
          class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"
        />
        <div
          class="absolute bottom-4 left-4 right-4 flex items-center justify-between"
        >
          <div
            class="flex items-center gap-2 bg-black/60 backdrop-blur-md rounded-full px-4 py-2 text-white text-md"
          >
            <UIcon
              name="i-lucide-battery-charging"
              class="size-5 text-warning"
            />
            <span>{{ t('home.electricVehicle') }}</span>
          </div>
          <div
            class="flex items-center gap-2 bg-black/60 backdrop-blur-md rounded-full px-4 py-2 text-white text-md"
          >
            <UIcon name="i-lucide-shield-check" class="size-5 text-warning" />
            <span>{{ t('home.dualControls') }}</span>
          </div>
        </div>
      </div>
    </UPageHero>

    <!-- Course Material -->
    <UPageSection
      id="material"
      :headline="t('home.courseMaterialHeadline')"
      :title="t('home.courseMaterial')"
      :description="t('home.courseMaterialDesc')"
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

    <!-- Pricing Section -->
    <UPageSection
      id="pricing"
      :headline="t('home.pricing')"
      :title="t('home.choosePath')"
      :description="t('home.pricingDesc')"
      :ui="{ headline: 'text-warning' }"
    >
      <!-- Premium Promo Banner with Real-time Countdown -->
      <Transition
        enter-active-class="transition duration-500 ease-out"
        enter-from-class="opacity-0 translate-y-4"
        enter-to-class="opacity-100 translate-y-0"
      >
        <div v-if="isPromoActive" class="mb-12 relative group">
          <!-- Background Glow Effect -->
          <div
            class="absolute -inset-1 bg-gradient-to-r from-warning-500 to-orange-600 rounded-3xl blur opacity-25 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"
          ></div>

          <div
            class="relative bg-white dark:bg-gray-900 border border-warning-500/20 rounded-2xl overflow-hidden shadow-2xl"
          >
            <div class="flex flex-col lg:flex-row">
              <!-- Left Side: Promo Info -->
              <div class="flex-1 p-6 md:p-8 flex items-center gap-6">
                <div class="relative">
                  <div
                    class="size-20 rounded-2xl bg-warning-500 flex items-center justify-center shadow-lg shadow-warning-500/30 animate-pulse"
                  >
                    <UIcon name="i-lucide-zap" class="size-10 text-white" />
                  </div>
                  <div
                    class="absolute -top-2 -right-2 bg-red-600 text-white text-[10px] font-black px-2 py-1 rounded-full uppercase tracking-tighter"
                  >
                    {{ t('home.hot') }}
                  </div>
                </div>

                <div>
                  <div
                    class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-warning-500/10 text-warning-600 text-xs font-bold uppercase tracking-widest mb-3"
                  >
                    {{ t('home.flashSale') }}
                  </div>
                  <h2
                    class="text-2xl md:text-3xl font-black tracking-tight mb-2"
                  >
                    {{ t('home.specialPrice') }}
                    <span class="text-warning-500"
                      >{{ translatedSaveAmount }}</span
                    >
                  </h2>
                  <p class="text-muted text-sm md:text-base max-w-lg">
                    {{ t('home.exclusiveDiscount') }}
                  </p>
                </div>
              </div>

              <!-- Right Side: Countdown Timer -->
              <div
                class="lg:w-[320px] bg-warning-50 dark:bg-warning-500/5 p-6 md:p-8 flex flex-col items-center justify-center border-t lg:border-t-0 lg:border-l border-warning-500/10"
              >
                <p
                  class="text-xs font-bold text-warning-600 uppercase tracking-[0.2em] mb-4"
                >
                  {{ t('home.promoEndsIn') }}
                </p>

                <div class="flex gap-3">
                  <div
                    v-for="(val, unit) in {
                      hours: timeLeft.hours,
                      minutes: timeLeft.minutes,
                      seconds: timeLeft.seconds,
                    }"
                    :key="unit"
                    class="text-center"
                  >
                    <div
                      class="size-14 md:size-16 bg-white dark:bg-gray-800 border border-warning-500/20 rounded-xl shadow-inner flex items-center justify-center mb-1"
                    >
                      <span class="text-2xl font-black text-foreground">{{
                        String(val).padStart(2, "0")
                      }}</span>
                    </div>
                    <span class="text-[10px] font-bold text-muted uppercase">{{
                      t(`common.${unit}`)
                    }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Plan Cards -->
      <div class="grid md:grid-cols-4 gap-6 mb-12">
        <div v-for="plan in packagesList" :key="plan.id" class="relative">
          <input
            :id="`plan-${plan.id}`"
            v-model="selectedPlan"
            type="radio"
            :value="plan.id"
            class="sr-only"
          />
          <label :for="`plan-${plan.id}`" class="block h-full cursor-pointer">
            <UCard
              :class="[
                'h-full flex flex-col transition-all',
                selectedPlan === plan.id
                  ? 'ring-2 ring-warning shadow-xl'
                  : 'hover:shadow-lg',
              ]"
            >
              <div
                v-if="plan.isPopular"
                class="absolute -top-3 left-1/2 -translate-x-1/2 z-10"
              >
                <UBadge :label="t('home.mostPopular')" color="warning" />
              </div>

              <template #header>
                <div class="text-center">
                  <h3 class="text-2xl font-bold">{{ plan.name }}</h3>
                  <p class="text-muted text-sm mt-2">
                    {{ t('home.packageDuration') }} {{ plan.duration }}
                  </p>
                </div>
              </template>

              <div class="flex-1 space-y-6">
                <!-- Pricing -->
                <div class="text-center">
                  <div
                    v-if="isPromoActive && plan.discountPrice < plan.price"
                    class="space-y-2"
                  >
                    <p class="text-sm text-muted line-through">
                      Rp {{ plan.price.toLocaleString("id-ID") }}
                    </p>
                    <p class="text-4xl font-bold text-warning">
                      Rp {{ plan.price.toLocaleString("id-ID") }}
                    </p>
                    <p class="text-xs text-green-600 font-medium">
                      {{ t('home.saveAmount', { amount: (plan.price - plan.discountPrice).toLocaleString("id-ID") }) }}
                    </p>
                  </div>
                  <div v-else>
                    <p class="text-4xl font-bold">
                      Rp {{ plan.price.toLocaleString("id-ID") }}
                    </p>
                  </div>
                </div>

                <!-- Sessions Badge -->
                <div class="text-center">
                  <UBadge
                    :label="`${plan.sessions} ${t('home.sessions')}`"
                    color="warning"
                    variant="subtle"
                  />
                </div>

                <!-- Features -->
                <ul class="space-y-3">
                  <li
                    v-for="feature in plan.features"
                    :key="feature"
                    class="flex items-start gap-3"
                  >
                    <UIcon
                      name="i-lucide-check"
                      class="size-5 text-warning shrink-0 mt-0.5"
                    />
                    <span class="text-sm">{{ feature }}</span>
                  </li>
                </ul>
              </div>

              <!-- Selection Indicator -->
              <div v-if="selectedPlan === plan.id" class="pt-4">
                <div class="flex flex-col gap-3">
                  <NuxtLink
                    :to="{ path: '/auth/register', query: { plan: plan.id } }"
                  >
                    <UButton :label="t('home.chooseThisPlan')" color="warning" block />
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
      :headline="t('home.easyBooking')"
      :title="t('home.bookSessions')"
      :description="t('home.bookingDesc')"
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
                <h3 class="font-semibold">{{ t('home.selectDate') }}</h3>
              </div>
              <div class="flex items-center gap-2">
                <UButton
                  icon="i-lucide-chevron-left"
                  variant="ghost"
                  color="neutral"
                  size="md"
                  @click="changeMonth(-1)"
                />
                <span class="text-md font-medium min-w-[100px] text-center">{{ currentMonth }}</span>
                <UButton
                  icon="i-lucide-chevron-right"
                  variant="ghost"
                  color="neutral"
                  size="md"
                  @click="changeMonth(1)"
                />
              </div>
            </div>
          </template>

          <div v-if="schedulesLoading" class="py-8">
            <div class="animate-pulse space-y-4">
              <div class="grid grid-cols-7 gap-1">
                <div v-for="i in 7" :key="i" class="h-8 bg-muted rounded"></div>
              </div>
              <div class="grid grid-cols-7 gap-1">
                <div v-for="i in 35" :key="i" class="aspect-square bg-muted rounded-lg"></div>
              </div>
            </div>
          </div>

          <!-- Custom Calendar Grid -->
          <template v-else>
            <div class="grid grid-cols-7 gap-1 mb-2">
              <div
                v-for="day in weekDays"
                :key="day"
                class="text-center text-md font-medium text-muted py-2"
              >
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
                        : 'text-muted/50 cursor-not-allowed',
                  ]"
                  @click="item.available && (selectedDate = item.day)"
                >
                  {{ item.day }}
                </button>
              </div>
            </div>

            <div class="mt-4 flex items-center gap-4 text-md">
              <div class="flex items-center gap-2">
                <div
                  class="size-3 rounded bg-primary/10 border border-primary/30"
                ></div>
                <span class="text-muted">{{ t('home.available') }}</span>
              </div>
              <div class="flex items-center gap-2">
                <div class="size-3 rounded bg-primary"></div>
                <span class="text-muted">{{ t('home.selected') }}</span>
              </div>
            </div>
          </template>
        </UCard>

        <!-- Time Slots -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-clock" class="size-5 text-warning" />
                <h3 class="font-semibold">{{ t('home.availableSlots') }}</h3>
              </div>
              <div class="flex items-center gap-2">
                <span v-if="schedulesLoading" class="text-xs text-muted">
                  <UIcon name="i-lucide-loader-2" class="size-3 animate-spin inline" />
                  Loading...
                </span>
                <UBadge
                  v-else
                  :label="`${currentMonthShortStr} ${selectedDate}`"
                  variant="subtle"
                />
              </div>
            </div>
          </template>

          <div v-if="schedulesLoading" class="space-y-3">
            <div
              v-for="i in 4"
              :key="i"
              class="w-full p-4 rounded-lg border border-default animate-pulse"
            >
              <div class="flex items-center justify-between">
                <div class="space-y-2">
                  <div class="h-4 w-16 bg-muted rounded"></div>
                  <div class="h-3 w-32 bg-muted rounded"></div>
                </div>
                <div class="h-6 w-16 bg-muted rounded"></div>
              </div>
            </div>
          </div>

          <div v-else class="space-y-3">
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
                  : 'border-default bg-muted/50 opacity-50 cursor-not-allowed',
              ]"
              @click="slot.available && (selectedSlot = slot.time)"
            >
              <div class="flex items-center justify-between">
                <div>
                  <span class="font-semibold">{{ slot.time }}</span>
                  <p class="text-md text-muted">
                    {{ slot.instructor }} - {{ slot.car }}
                  </p>
                </div>
                <UBadge
                  :label="slot.available ? t('home.available') : t('home.booked')"
                  :color="slot.available ? 'success' : 'error'"
                  variant="subtle"
                />
              </div>
            </button>

            <div
              v-if="timeSlots.length === 0"
              class="text-center py-8"
            >
              <UIcon name="i-lucide-calendar-x" class="size-8 text-muted mx-auto mb-2" />
              <p class="text-muted text-sm">{{ t('home.noSlotsAvailable') || 'No slots available for this date' }}</p>
            </div>
          </div>

          <template #footer>
            <NuxtLink
              :to="{ path: '/auth/register', query: { plan: selectedPlan } }"
              class="w-full"
            >
              <UButton
                :label="t('home.chooseSchedule')"
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
      :headline="t('home.instructors')"
      :title="t('home.meetInstructors')"
      :description="t('home.instructorsDesc')"
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid md:grid-cols-3 gap-8">
        <UCard
          v-for="instructor in instructorsList"
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
            <div
              class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
            />
          </div>

          <div class="p-6">
            <h3 class="text-xl font-bold">{{ instructor.name }}</h3>
            <p class="text-warning font-medium text-sm mb-4">
              {{ instructor.yearsOfExperience }} {{ t('home.yearsExperience') }}
            </p>
            <p class="text-muted text-sm">{{ instructor.bio }}</p>
          </div>
        </UCard>
      </div>

      <div class="mt-12 text-center">
        <NuxtLink to="/instructors">
          <UButton
            :label="t('home.viewAllInstructors')"
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
      :headline="t('home.location')"
      :title="t('home.convenientlyLocated')"
      :description="t('home.locationDesc')"
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
                <h3 class="font-semibold mb-1">{{ t('home.trainingCenter') }}</h3>
                <p class="text-muted text-md">
                  {{ t('home.address') }}
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
                <h3 class="font-semibold mb-1">{{ t('home.operatingHours') }}</h3>
                <p class="text-muted text-md">
                  {{ t('home.hoursWeekday') }}<br />
                  {{ t('home.hoursWeekend') }}<br />
                  {{ t('home.hoursNight') }}
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
                <h3 class="font-semibold mb-1">{{ t('home.contactUs') }}</h3>
                <p class="text-muted text-md">
                  {{ t('home.phone') }}: +62 812-3456-7890<br />
                  {{ t('home.email') }}: info@evdriveacademy.id
                </p>
              </div>
            </div>
          </UCard>

          <div class="flex gap-3">
            <NuxtLink
              to="https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi"
              target="_blank"
              class="flex-1"
            >
              <UButton
                icon="i-simple-icons-whatsapp"
                :label="t('home.chatWhatsApp')"
                class="!bg-[#25D366] hover:!bg-[#128C7E] text-white"
                block
              />
            </NuxtLink>
            <NuxtLink
              to="https://maps.app.goo.gl/RpSdkpjs4RZg2ZY77"
              target="_blank"
              class="flex-1"
            >
              <UButton
                icon="i-lucide-navigation"
                :label="t('home.getDirections')"
                color="neutral"
                variant="outline"
                block
              />
            </NuxtLink>
          </div>
        </div>

        <!-- Map Placeholder -->
        <div
          class="h-auto lg:h-full min-h-auto rounded-2xl overflow-hidden bg-elevated border border-default"
        >
          <iframe
            src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3966.1806061048765!2d106.65588507475077!3d-6.2399118937483635!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69fbc070b4d71d%3A0x8b1a633faf5dbd46!2sALAM%20SUTERA!5e0!3m2!1sen!2sid!4v1776223155011!5m2!1sen!2sid"
            width="600"
            height="550"
            style="border: 0"
            allowfullscreen
            loading="lazy"
            referrerpolicy="no-referrer-when-downgrade"
          ></iframe>
        </div>
      </div>
    </UPageSection>

    <!-- Testimonials -->
    <UPageSection
      :headline="t('home.testimonials')"
      :title="t('home.whatStudentsSay')"
      :description="t('home.testimonialsDesc')"
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid md:grid-cols-3 gap-6">
        <UCard
          v-for="testimonial in testimonialsList"
          :key="testimonial.userName"
        >
          <div class="flex items-center gap-3 mb-4">
            <UAvatar :text="testimonial.userImage" size="lg" />
            <div>
              <p class="font-semibold">{{ testimonial.userName }}</p>
              <p class="text-md text-muted">{{ testimonial.userRole }}</p>
            </div>
          </div>
          <div class="flex gap-0.5 mb-3">
            <UIcon
              v-for="i in testimonial.rating"
              :key="i"
              name="i-lucide-star"
              class="size-4 text-amber-500 fill-amber-500"
            />
          </div>
          <p class="text-muted">"{{ testimonial.content }}"</p>
        </UCard>
      </div>
    </UPageSection>

    <!-- FAQ Section -->
    <UPageSection
      id="faq"
      :headline="t('home.faq')"
      :title="t('home.frequentlyAsked')"
      :description="t('home.faqDesc')"
      :ui="{ headline: 'text-warning', body: 'w-full min-w-full', container: 'max-w-full px-0'}"
      class="bg-muted/30"
    >
      <div class="max-w-3xl mx-auto">
        <UAccordion :items="faqItems" :ui="{ body: 'w-full min-w-full', label: 'w-full min-w-full block', header: 'max-w-none w-full px-0', }" />
      </div>
    </UPageSection>

    <!-- CTA Section -->
    <UPageCTA
      :title="t('home.readyToDrive')"
      :description="t('home.readyToDriveDesc')"
      :links="[
        {
          label: t('home.startJourney'),
          to: '/auth/register',
          color: 'warning',
          icon: 'i-lucide-rocket',
          size: 'lg',
        },
        {
          label: t('home.viewAllPackages'),
          to: '/packages',
          color: 'neutral',
          variant: 'ghost',
          trailingIcon: 'i-lucide-arrow-right',
        },
      ]"
      class="bg-warning/5"
    />
  </div>
</template>
