<script setup lang="ts">
const { t } = useI18n()

useSeoMeta({
  title: t('packages.heroTitle') + " | Drive Master Academy",
  description: t('packages.subtitle'),
});

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

const packagesStore = usePackagesStore();

const { promoEndDate } = useSettings();

const packagePlans = computed(() => packagesStore.packages);

const isPromoActive = computed(() => {
  const promoEnd = new Date(promoEndDate.value || "2026-05-31T23:59:59");
  return now.value < promoEnd;
});

const currentPlan = computed(() =>
  packagePlans.value.find((p) => p.id === selectedPlan.value),
);
const discount = computed(() => {
  if (!currentPlan.value || !isPromoActive.value) return 0;
  return currentPlan.value.discountPrice;
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

const comparisonFeatures = computed(() => [
  { name: t('packages.comparison.freeTrial'), pkg6x: true, pkg8x: true, pkg10x: true, pkg12x: true },
  { name: t('packages.comparison.sessions'), pkg6x: "6", pkg8x: "8", pkg10x: "10", pkg12x: "12" },
  {
    name: t('packages.comparison.totalHours'),
    pkg6x: `6 ${t('packages.comparison.hrs')}`,
    pkg8x: `8 ${t('packages.comparison.hrs')}`,
    pkg10x: `10 ${t('packages.comparison.hrs')}`,
    pkg12x: `12 ${t('packages.comparison.hrs')}`,
  },
]);

const addOns = computed(() => [
  {
    title: t('packages.extras.extraSession'),
    price: "Rp 350.000",
    description: t('packages.extras.extraSessionDesc'),
    icon: "i-lucide-plus-circle",
  },
]);

onMounted(async () => {
  packagesStore.fetchPackages();
});
</script>

<template>
  <div>
    <!-- Hero -->
    <UPageHero
      :title="t('packages.heroTitle')"
      :description="t('packages.subtitle')"
      :links="[
        {
          label: t('packages.registerNow'),
          to: '/auth/register',
          color: 'warning',
          icon: 'i-lucide-user-plus',
        },
        {
          label: t('packages.comparePackages'),
          to: '#comparison',
          color: 'neutral',
          variant: 'outline',
        },
      ]"
    />

    <!-- Package Cards -->
    <UPageSection :ui="{ headline: 'text-warning' }">
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
                      >{{ t('home.saveUpTo', { discount: Math.round(((currentPlan?.price || 1) - (currentPlan?.discountPrice || 0)) / (currentPlan?.price || 1) * 100) }) }}</span
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
                      t(`common.${unit}`).substring(0, 3)
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
        <div v-for="plan in packagePlans" :key="plan.id" class="relative">
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
                <UBadge :label="t('packages.popular')" color="warning" />
              </div>

              <template #header>
                <div class="text-center">
                  <h3 class="text-2xl font-bold">{{ plan.name }}</h3>
                  <p class="text-muted text-sm mt-2">
                    {{ t('packages.duration') }} : {{ plan.duration }}
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
                      Rp {{ plan.discountPrice.toLocaleString("id-ID") }}
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
                    :label="`${plan.sessions} ${t('packages.sessions')}`"
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
                    <UButton :label="t('packages.choosePlan')" color="warning" block />
                  </NuxtLink>
                </div>
              </div>
            </UCard>
          </label>
        </div>
      </div>
    </UPageSection>

    <!-- Comparison Table -->
    <UPageSection
      id="comparison"
      :headline="t('packages.comparison.subtitle')"
      :title="t('packages.comparison.title')"
      :description="t('packages.comparison.description')"
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <UCard class="overflow-x-auto">
        <table class="w-full min-w-[600px]">
          <thead>
            <tr class="border-b border-default">
              <th class="text-left py-4 px-4 font-semibold">{{ t('packages.comparison.feature') }}</th>
              <th class="text-center py-4 px-4 font-semibold">6x</th>
              <th class="text-center py-4 px-4 font-semibold">
                <div class="flex items-center justify-center gap-2">
                  8x
                  <UBadge :label="t('packages.popular')" size="xs" color="warning" />
                </div>
              </th>
              <th class="text-center py-4 px-4 font-semibold">10x</th>
              <th class="text-center py-4 px-4 font-semibold">12x</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="feature in comparisonFeatures"
              :key="feature.name"
              class="border-b border-default last:border-0"
            >
              <td class="py-3 px-4 text-sm">{{ feature.name }}</td>
              <td class="py-3 px-4 text-center">
                <template v-if="typeof feature.pkg6x === 'boolean'">
                  <UIcon
                    :name="feature.pkg6x ? 'i-lucide-check' : 'i-lucide-minus'"
                    :class="feature.pkg6x ? 'text-warning' : 'text-muted'"
                    class="size-5"
                  />
                </template>
                <template v-else>
                  <span class="text-sm">{{ feature.pkg6x }}</span>
                </template>
              </td>
              <td class="py-3 px-4 text-center bg-warning/5">
                <template v-if="typeof feature.pkg8x === 'boolean'">
                  <UIcon
                    :name="feature.pkg8x ? 'i-lucide-check' : 'i-lucide-minus'"
                    :class="feature.pkg8x ? 'text-warning' : 'text-muted'"
                    class="size-5"
                  />
                </template>
                <template v-else>
                  <span class="text-sm font-medium">{{ feature.pkg8x }}</span>
                </template>
              </td>
              <td class="py-3 px-4 text-center">
                <template v-if="typeof feature.pkg10x === 'boolean'">
                  <UIcon
                    :name="feature.pkg10x ? 'i-lucide-check' : 'i-lucide-minus'"
                    :class="feature.pkg10x ? 'text-warning' : 'text-muted'"
                    class="size-5"
                  />
                </template>
                <template v-else>
                  <span class="text-sm">{{ feature.pkg10x }}</span>
                </template>
              </td>
              <td class="py-3 px-4 text-center">
                <template v-if="typeof feature.pkg12x === 'boolean'">
                  <UIcon
                    :name="feature.pkg12x ? 'i-lucide-check' : 'i-lucide-minus'"
                    :class="feature.pkg12x ? 'text-warning' : 'text-muted'"
                    class="size-5"
                  />
                </template>
                <template v-else>
                  <span class="text-sm">{{ feature.pkg12x }}</span>
                </template>
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-muted/50">
              <td class="py-4 px-4 font-semibold">{{ t('packages.comparison.price') }}</td>
              <td
                v-for="plan in packagePlans"
                :key="plan.id"
                class="py-4 px-4 text-center font-semibold"
                :class="{ 'text-warning bg-warning/10': plan.isPopular }"
              >
                Rp {{ (plan.isDiscounted ? plan.discountPrice : plan.price).toLocaleString('id-ID') }}
              </td>
            </tr>
          </tfoot>
        </table>
      </UCard>
    </UPageSection>

    <!-- Add-ons -->
    <UPageSection
      :headline="t('packages.extras.subtitle')"
      :title="t('packages.extras.title')"
      :description="t('packages.extras.description')"
      :ui="{ headline: 'text-warning' }"
    >
      <div class="grid md:grid-cols-1 gap-6 max-w-4xl mx-auto">
        <UCard v-for="addon in addOns" :key="addon.title">
          <div class="flex items-start gap-4">
            <div class="p-3 rounded-lg bg-warning/10">
              <UIcon :name="addon.icon" class="size-6 text-warning" />
            </div>
            <div>
              <h3 class="font-semibold">{{ addon.title }}</h3>
              <p class="text-warning font-bold">{{ addon.price }}</p>
              <p class="text-sm text-muted mt-1">{{ addon.description }}</p>
            </div>
          </div>
        </UCard>
      </div>
    </UPageSection>

    <!-- Payment Info Section (Hidden) -->

    <!-- CTA -->
    <UPageCTA
      :title="t('packages.cta.title')"
      :description="t('packages.cta.description')"
      :links="[
        {
          label: t('packages.registerNow'),
          to: '/auth/register',
          color: 'warning',
          icon: 'i-lucide-user-plus',
        },
        {
          label: t('packages.cta.contactUs'),
          to: 'https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi',
          color: 'primary',
          variant: 'outline',
          icon: 'i-simple-icons-whatsapp',
          external: true,
        },
      ]"
    />
  </div>
</template>
