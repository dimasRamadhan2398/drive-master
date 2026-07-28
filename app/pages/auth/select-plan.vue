<script setup lang="ts">
import { ref, computed, onMounted } from "vue";

const { t } = useI18n();
definePageMeta({
  layout: "blank",
});

const route = useRoute();
const packagesStore = usePackagesStore();
const enrollmentsStore = useEnrollmentsStore();
const authStore = useAuthStore();

// Initialize selected plan type
const selectedPlanId = ref<string | null>(null);
const selectedAddons = ref<string[]>([]);

// Fetch packages and addons on mount
onMounted(async () => {
  if (packagesStore.packages.length === 0) {
    await packagesStore.fetchPackages({ status: "active" });
  }
  if (packagesStore.addons.length === 0) {
    await packagesStore.fetchAddons();
  }
  // Set default selected plan
  if (!selectedPlanId.value && packagesStore.activePackages.length > 0) {
    // Default to the popular package or first active package
    const popular = packagesStore.popularPackages[0];
    selectedPlanId.value = popular?.id || packagesStore.activePackages[0].id;
  }
});

const isExtraSession = (addon: any) => {
  return addon.name.toLowerCase().includes("extra session") || addon.id === "22222222-2222-2222-2222-222222222201";
};

const getAddonCount = (id: string) => {
  return selectedAddons.value.filter((x) => x === id).length;
};

const incrementAddon = (id: string) => {
  selectedAddons.value.push(id);
};

const decrementAddon = (id: string) => {
  const index = selectedAddons.value.indexOf(id);
  if (index !== -1) {
    selectedAddons.value.splice(index, 1);
  }
};

function toggleAddon(addon: any) {
  const id = addon.id;
  if (isExtraSession(addon)) {
    const count = getAddonCount(id);
    if (count > 0) {
      selectedAddons.value = selectedAddons.value.filter((x) => x !== id);
    } else {
      selectedAddons.value.push(id);
    }
  } else {
    const index = selectedAddons.value.indexOf(id);
    if (index === -1) {
      selectedAddons.value.push(id);
    } else {
      selectedAddons.value.splice(index, 1);
    }
  }
}

// Map store packages to plan format with promo pricing
const plans = computed(() => {
  return packagesStore.activePackages.map((pkg: typeof packagesStore.packages[0]) => ({
    id: pkg.id,
    name: pkg.name,
    originalPrice: pkg.price,
    promoPrice: pkg.discountPrice,
    duration: `${pkg.duration} min`,
    sessions: pkg.sessions,
    features: pkg.features || [],
    highlight: pkg.isPopular,
  }));
});

const loading = ref(false);

// Timer Logic
const now = ref(new Date());
if (process.client) {
  setInterval(() => {
    now.value = new Date();
  }, 1000);
}

const { promoEndDate } = useSettings();

const isPromoActive = computed(() => {
  const promoEnd = new Date(promoEndDate.value || "2026-05-31T23:59:59");
  return now.value < promoEnd;
});

const timeLeft = computed(() => {
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

const currentActivePlan = computed(() => route.query.current_plan as string | null);

type PlanType = {
  id: string;
  name: string;
  originalPrice: number;
  promoPrice: number;
  duration: string;
  sessions: number;
  features: string[];
  highlight: boolean;
};

const currentPlan = computed(() => plans.value.find((p: PlanType) => p.id === selectedPlanId.value));

const discount = computed(() => {
  if (!currentPlan.value || !isPromoActive.value) return 0;
  return Math.ceil(
    ((currentPlan.value.originalPrice - currentPlan.value.promoPrice) /
      currentPlan.value.originalPrice) *
      100
  );
});

const groupedSelectedAddons = computed(() => {
  const counts: Record<string, number> = {};
  for (const id of selectedAddons.value) {
    counts[id] = (counts[id] || 0) + 1;
  }
  return Object.entries(counts).map(([id, count]) => {
    const addon = packagesStore.addons.find((a) => a.id === id);
    return {
      id,
      count,
      name: addon ? addon.name : "",
      price: addon ? addon.price * count : 0,
    };
  }).filter((item) => item.name);
});

const totalPrice = computed(() => {
  if (!currentPlan.value) return 0;
  const basePrice = isPromoActive.value ? currentPlan.value.promoPrice : currentPlan.value.originalPrice;
  const addonsPrice = selectedAddons.value.reduce((sum, id) => {
    const addon = packagesStore.addons.find((a) => a.id === id);
    return sum + (addon ? addon.price : 0);
  }, 0);
  return basePrice + addonsPrice;
});

async function selectPlan() {
  if (!currentPlan.value) return;

  loading.value = true;

  try {
    // Get the selected package details
    const selectedPackage = packagesStore.getPackageById(currentPlan.value.id);
    if (!selectedPackage) {
      console.error("[SELECT PLAN] Package not found:", currentPlan.value.id);
      return;
    }

    // Ensure user is authenticated
    if (!authStore.userId) {
      console.error("[SELECT PLAN] User not authenticated");
      return;
    }

    // Create enrollment with userId, packageId, and selected addons
    const enrollment = await enrollmentsStore.createEnrollment({
      userId: authStore.userId,
      packageId: currentPlan.value.id,
      price: selectedPackage.price,
      discountPrice: selectedPackage.discountPrice,
      addOns: selectedAddons.value,
    });

    if (enrollment) {
      console.log("[SELECT PLAN] Enrollment created:", enrollment);
      console.log("[SELECT PLAN] Enrollment ID:", enrollment.id);
      console.log("[SELECT PLAN] Enrollment object keys:", Object.keys(enrollment));

      // Store full enrollment object in session storage for payment page
      if (import.meta.client) {
        const enrollmentData = {
          id: enrollment.id,
          userId: enrollment.userId,
          packageId: enrollment.packageId,
          packageName: enrollment.packageName || "Selected Package",
          price: enrollment.price || enrollment.totalPrice || 0,
          discountPrice: enrollment.discountPrice || 0,
          status: enrollment.status,
          createdAt: enrollment.createdAt,
        };
        console.log("[SELECT PLAN] Storing enrollment to session:", enrollmentData);
        sessionStorage.setItem('dm_enrollment', JSON.stringify(enrollmentData));
        sessionStorage.setItem('dm_enrollment_id', enrollment.id);
        sessionStorage.setItem('dm_selected_plan', currentPlan.value.id);
      }

      // Redirect to payment method with enrollment ID
      navigateTo(`/auth/payment-method?enrollment=${enrollment.id}`);
    } else {
      console.error("[SELECT PLAN] Failed to create enrollment");
      const toast = useToast();
      toast.add({
        title: t("register.errors.enrollmentError") || "Enrollment Failed",
        description: enrollmentsStore.error || t("register.errors.enrollmentFailed") || "Failed to create enrollment. Please try again.",
        color: "error",
      });
    }
  } finally {
    loading.value = false;
  }
}

const freeTrialInfo = computed(() => [
  {
    icon: "i-lucide-gift",
    title: t("freeTrial.heroTitle"),
    description: t("freeTrial.oncePerAccount"),
  },
  {
    icon: "i-lucide-clock",
    title: t("freeTrial.fifteenMinutes"),
    description: t("freeTrial.trialSessionDesc"),
  },
  {
    icon: "i-lucide-zap",
    title: t("freeTrial.fullExperience"),
    description: t("freeTrial.experienceDesc"),
  },
]);
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4 bg-muted/20">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-12">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon name="i-lucide-package" class="size-8 text-warning" />
          <span class="text-xl font-bold">{{ t("packages.heroTitle") }}</span>
        </div>
        <h1 class="text-3xl md:text-4xl font-bold">{{ t("packages.subtitle") }}</h1>
        <p class="text-muted mt-3 max-w-2xl mx-auto">
          {{ t("packages.subtitle") }}
        </p>
      </div>

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
                    {{ t("home.hot") }}
                  </div>
                </div>

                <div>
                  <div
                    class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-warning-500/10 text-warning-600 text-xs font-bold uppercase tracking-widest mb-3"
                  >
                    {{ t("home.flashSale") }}
                  </div>
                  <h2 class="text-2xl md:text-3xl font-black tracking-tight mb-2">
                    {{ t("home.specialPrice") }}
                    <span class="text-warning-500">{{
                      t("home.saveUpTo", { discount: discount })
                    }}</span>
                  </h2>
                  <p class="text-muted text-sm md:text-base max-w-lg">
                    {{ t("home.exclusiveDiscount") }}
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
                  {{ t("home.promoEndsIn") }}
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
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 md:gap-6 mb-12">
        <!-- Loading State -->
        <div v-if="packagesStore.isLoading" class="col-span-full flex justify-center py-12">
          <UIcon name="i-lucide-loader-2" class="size-8 text-warning animate-spin" />
        </div>

        <div v-else v-for="plan in plans" :key="plan.id" class="relative">
          <input
            :id="`plan-${plan.id}`"
            v-model="selectedPlanId"
            type="radio"
            :value="plan.id"
            class="sr-only"
            :disabled="plan.id === currentActivePlan"
          />
          <label
            :for="`plan-${plan.id}`"
            :class="[
              'block h-full',
              plan.id === currentActivePlan ? 'cursor-not-allowed' : 'cursor-pointer',
            ]"
          >
            <UCard
              :class="[
                'h-full flex flex-col transition-all',
                selectedPlanId === plan.id
                  ? 'ring-2 ring-warning shadow-xl'
                  : 'hover:shadow-lg',
                plan.id === currentActivePlan &&
                  'opacity-60 bg-gray-50 dark:bg-gray-800/50',
              ]"
            >
              <div
                v-if="plan.highlight"
                class="absolute -top-3 left-1/2 -translate-x-1/2 z-10"
              >
                <UBadge :label="t('packages.popular')" color="warning" />
              </div>

              <template #header>
                <div class="text-center">
                  <h3 class="text-lg sm:text-2xl font-bold">{{ plan.name }}</h3>
                  <p class="text-muted text-xs sm:text-sm mt-1 sm:mt-2">
                    {{ t("packages.duration") }} : {{ plan.duration }}
                  </p>
                  <div v-if="plan.id === currentActivePlan" class="mt-2">
                    <UBadge :label="t('packages.currentPlan')" color="neutral" variant="subtle" />
                  </div>
                </div>
              </template>

              <div class="flex-1 space-y-4 sm:space-y-6">
                <!-- Pricing -->
                <div class="text-center">
                  <div
                    v-if="isPromoActive && plan.promoPrice < plan.originalPrice"
                    class="space-y-1.5 sm:space-y-2"
                  >
                    <p class="text-xs sm:text-sm text-muted line-through">
                      Rp {{ plan.originalPrice.toLocaleString("id-ID") }}
                    </p>
                    <p class="text-xl sm:text-3xl md:text-4xl font-bold text-warning">
                      Rp {{ plan.promoPrice.toLocaleString("id-ID") }}
                    </p>
                    <p class="text-[10px] sm:text-xs text-green-600 font-medium">
                      {{
                        t("home.saveAmount", {
                          amount: (plan.originalPrice - plan.promoPrice).toLocaleString(
                            "id-ID"
                          ),
                        })
                      }}
                    </p>
                  </div>
                  <div v-else>
                    <p class="text-xl sm:text-3xl md:text-4xl font-bold">
                      Rp {{ plan.originalPrice.toLocaleString("id-ID") }}
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
                <ul class="space-y-2 sm:space-y-3">
                  <li
                    v-for="feature in plan.features"
                    :key="feature"
                    class="flex items-start gap-1.5 sm:gap-3"
                  >
                    <UIcon
                      name="i-lucide-check"
                      class="size-4 sm:size-5 text-warning shrink-0 mt-0.5"
                    />
                    <span class="text-xs sm:text-sm">{{ feature }}</span>
                  </li>
                </ul>
              </div>

              <!-- Selection Indicator -->
              <div v-if="selectedPlanId === plan.id" class="pt-4">
                <div class="flex items-center justify-between">
                  <span class="text-xs sm:text-sm font-medium">{{ t("home.selected") }}</span>
                  <UIcon name="i-lucide-check" class="size-4 text-warning" />
                </div>
              </div>
            </UCard>
          </label>
        </div>
      </div>

      <!-- Add-ons Section -->
      <div class="mb-12">
        <div class="text-center mb-8">
          <h2 class="text-2xl font-bold">{{ t("packages.extras.title") }}</h2>
          <p class="text-muted mt-2">
            {{ t("packages.addonsDesc") }}
          </p>
        </div>

        <div class="flex flex-wrap justify-center gap-6 max-w-4xl mx-auto">
          <!-- Loading State for Add-ons -->
          <div v-if="packagesStore.isLoading" class="w-full flex justify-center items-center py-12">
            <UIcon name="i-lucide-loader-2" class="size-8 text-warning animate-spin" />
          </div>
          
          <UCard
            v-else
            v-for="addon in packagesStore.addons"
            :key="addon.id"
            :class="[
              'cursor-pointer border-2 transition-all w-full sm:w-[280px]',
              selectedAddons.includes(addon.id)
                ? 'border-warning bg-warning/5 shadow-md'
                : 'border-transparent hover:border-warning/30 hover:shadow'
            ]"
            @click="toggleAddon(addon)"
          >
            <div class="flex items-start gap-4 h-full">
              <UCheckbox
                v-if="!isExtraSession(addon)"
                :model-value="selectedAddons.includes(addon.id)"
                color="warning"
                class="pointer-events-none mt-1"
              />
              <UCheckbox
                v-else
                :model-value="getAddonCount(addon.id) > 0"
                color="warning"
                class="pointer-events-none mt-1"
              />
              <div class="flex-1 space-y-1">
                <h3 class="font-semibold text-sm sm:text-base text-foreground">
                  {{ addon.name }}
                </h3>
                <p class="text-warning font-bold text-xs sm:text-sm">
                  Rp {{ addon.price.toLocaleString("id-ID") }}
                </p>
                <p class="text-xs text-muted leading-relaxed">
                  {{ addon.description }}
                </p>

                <div class="pt-2 flex items-center justify-between gap-2 flex-wrap">
                  <div v-if="addon.sessions && !isExtraSession(addon)" class="pt-1">
                    <UBadge
                      :label="t('packages.extraSessionsCount', { count: addon.sessions })"
                      color="warning"
                      variant="subtle"
                      size="xs"
                    />
                  </div>

                  <!-- If it's Extra Session, show the Badge AND the Increment/Decrement controls -->
                  <template v-if="isExtraSession(addon)">
                    <UBadge
                      v-if="getAddonCount(addon.id) > 0"
                      :label="t('packages.extraSessionsCount', { count: getAddonCount(addon.id) })"
                      color="warning"
                      variant="subtle"
                      size="xs"
                    />
                    
                    <!-- Counter Controls -->
                    <div @click.stop class="flex items-center bg-gray-100 dark:bg-gray-800 rounded-lg p-0.5 border border-gray-200 dark:border-gray-700">
                      <UButton
                        icon="i-lucide-minus"
                        size="xs"
                        color="neutral"
                        variant="ghost"
                        :disabled="getAddonCount(addon.id) === 0"
                        @click="decrementAddon(addon.id)"
                      />
                      <span class="px-2.5 text-xs font-semibold text-foreground min-w-[20px] text-center">
                        {{ getAddonCount(addon.id) }}
                      </span>
                      <UButton
                        icon="i-lucide-plus"
                        size="xs"
                        color="warning"
                        variant="ghost"
                        @click="incrementAddon(addon.id)"
                      />
                    </div>
                  </template>
                </div>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- Free Trial Info Section -->
      <div class="mb-12">
        <div class="text-center mb-8">
          <h2 class="text-2xl font-bold">{{ t("freeTrial.heroTitle") }}</h2>
          <p class="text-muted mt-2">
            {{ t("freeTrial.subtitle") }}
          </p>
        </div>

        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <UCard v-for="item in freeTrialInfo" :key="item.title">
            <div class="text-center space-y-3">
              <div class="flex justify-center">
                <div class="p-3 rounded-lg bg-warning/10">
                  <UIcon :name="item.icon" class="size-6 text-warning" />
                </div>
              </div>
              <div>
                <p class="font-semibold text-sm">{{ item.title }}</p>
                <p class="text-xs text-muted mt-1">{{ item.description }}</p>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- CTA Section -->
      <div class="max-w-2xl mx-auto">
        <UCard class="bg-warning/5 border-warning/20">
          <div class="text-center space-y-4">
            <h3 class="text-xl font-bold">{{ t("packages.readyToGetStarted") }}</h3>
            <div class="space-y-1 text-muted">
              <p>
                {{ t("packages.selectedPackage") }}
                <span class="font-bold text-foreground">{{ currentPlan?.name }}</span>
              </p>
              <p v-if="selectedAddons.length > 0" class="text-xs">
                {{ t("packages.selectedAddons") }}
                <span class="font-semibold text-foreground">
                  {{
                    groupedSelectedAddons
                      .map((item) => item.count > 1 ? `${item.name} (${item.count}x)` : item.name)
                      .join(", ")
                  }}
                </span>
              </p>
              <p class="text-lg font-bold text-warning pt-1">
                {{ t("billing.totalPrice") }}: Rp {{ totalPrice.toLocaleString("id-ID") }}
              </p>
            </div>

            <div class="flex gap-3 pt-2">
              <UButton
                :label="t('billing.proceedToPayment')"
                icon="i-lucide-arrow-right"
                color="warning"
                :loading="loading"
                @click="selectPlan"
                block
              />
            </div>

            <p class="text-xs text-muted">
              {{ t("packages.paymentSecuredNote") }}
            </p>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>
