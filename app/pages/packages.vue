<script setup lang="ts">
const packages = [
  {
    title: 'Starter',
    price: 'Rp 1.500.000',
    priceNote: 'One-time payment',
    description: 'Perfect for beginners looking to get started with driving basics.',
    features: [
      { text: '5 training sessions', included: true },
      { text: '45 minutes per session', included: true },
      { text: 'Basic driving fundamentals', included: true },
      { text: 'Theory materials included', included: true },
      { text: 'Certificate of completion', included: true },
      { text: 'Highway driving', included: false },
      { text: 'EV-specific training', included: false },
      { text: 'Night driving', included: false }
    ],
    highlight: false,
    color: 'neutral' as const
  },
  {
    title: 'Standard',
    price: 'Rp 2.800.000',
    priceNote: 'One-time payment',
    description: 'Our most popular package for comprehensive learning.',
    features: [
      { text: '10 training sessions', included: true },
      { text: '60 minutes per session', included: true },
      { text: 'Basic & advanced techniques', included: true },
      { text: 'Theory materials included', included: true },
      { text: 'Highway & city driving', included: true },
      { text: 'EV-specific training', included: true },
      { text: 'Official EV Certificate', included: true },
      { text: 'Night driving', included: false }
    ],
    highlight: true,
    color: 'primary' as const
  },
  {
    title: 'Pro',
    price: 'Rp 4.500.000',
    priceNote: 'One-time payment',
    description: 'Complete mastery with unlimited support and premium perks.',
    features: [
      { text: '15 training sessions', included: true },
      { text: '90 minutes per session', included: true },
      { text: 'Full driving mastery program', included: true },
      { text: 'Theory materials included', included: true },
      { text: 'Highway & city driving', included: true },
      { text: 'EV-specific training', included: true },
      { text: 'Night & weather driving', included: true },
      { text: 'Premium EV Certificate', included: true }
    ],
    highlight: false,
    color: 'neutral' as const
  }
]

const comparisonFeatures = [
  { name: 'Training Sessions', starter: '5', standard: '10', pro: '15' },
  { name: 'Session Duration', starter: '45 min', standard: '60 min', pro: '90 min' },
  { name: 'Total Training Hours', starter: '3.75 hrs', standard: '10 hrs', pro: '22.5 hrs' },
  { name: 'Theory Materials', starter: true, standard: true, pro: true },
  { name: 'Basic Driving', starter: true, standard: true, pro: true },
  { name: 'Advanced Techniques', starter: false, standard: true, pro: true },
  { name: 'Highway Driving', starter: false, standard: true, pro: true },
  { name: 'Night Driving', starter: false, standard: false, pro: true },
  { name: 'Weather Driving', starter: false, standard: false, pro: true },
  { name: 'EV-Specific Training', starter: false, standard: true, pro: true },
  { name: 'EV Maintenance Basics', starter: false, standard: false, pro: true },
  { name: 'Certificate Type', starter: 'Basic', standard: 'EV Cert', pro: 'Premium' },
  { name: 'Refresher Access', starter: '30 days', standard: '6 months', pro: 'Lifetime' }
]

const addOns = [
  {
    title: 'Extra Session',
    price: 'Rp 350.000',
    description: 'Add more practice time to any package',
    icon: 'i-lucide-plus-circle'
  },
  {
    title: 'Simulator Session',
    price: 'Rp 150.000',
    description: 'Practice in our advanced driving simulator',
    icon: 'i-lucide-monitor'
  },
  {
    title: 'SIM Exam Prep',
    price: 'Rp 500.000',
    description: 'Dedicated preparation for official license test',
    icon: 'i-lucide-graduation-cap'
  }
]
</script>

<template>
  <div>
    <!-- Hero -->
    <UPageHero
      title="Course Packages"
      description="Transparent pricing with no hidden fees. Choose the package that best fits your learning goals and schedule."
      :links="[
        { label: 'Register Now', to: '/register', icon: 'i-lucide-user-plus' },
        { label: 'Compare Packages', to: '#comparison', color: 'neutral', variant: 'outline', icon: 'i-lucide-scale' }
      ]"
      class="py-12 lg:py-20"
    />

    <!-- Package Cards -->
    <UPageSection>
      <div class="grid md:grid-cols-3 gap-6 lg:gap-8">
        <UCard
          v-for="pkg in packages"
          :key="pkg.title"
          :class="[
            'flex flex-col relative',
            pkg.highlight ? 'ring-2 ring-primary shadow-xl' : ''
          ]"
        >
          <UBadge 
            v-if="pkg.highlight" 
            label="Most Popular" 
            color="primary" 
            class="absolute -top-3 left-1/2 -translate-x-1/2"
          />

          <template #header>
            <div class="text-center pt-2">
              <h3 class="text-2xl font-bold">{{ pkg.title }}</h3>
              <div class="mt-4">
                <span class="text-4xl font-bold text-primary">{{ pkg.price }}</span>
              </div>
              <p class="text-sm text-muted mt-1">{{ pkg.priceNote }}</p>
              <p class="text-muted mt-4">{{ pkg.description }}</p>
            </div>
          </template>

          <ul class="space-y-3 flex-1">
            <li 
              v-for="feature in pkg.features" 
              :key="feature.text" 
              class="flex items-center gap-3"
              :class="{ 'opacity-50': !feature.included }"
            >
              <UIcon 
                :name="feature.included ? 'i-lucide-check' : 'i-lucide-x'" 
                :class="feature.included ? 'text-primary' : 'text-muted'" 
                class="size-5 shrink-0" 
              />
              <span class="text-sm">{{ feature.text }}</span>
            </li>
          </ul>

          <template #footer>
            <NuxtLink to="/register" class="w-full">
              <UButton 
                :label="pkg.highlight ? 'Get Started' : 'Choose Plan'"
                :color="pkg.color"
                :variant="pkg.highlight ? undefined : 'outline'"
                block
              />
            </NuxtLink>
          </template>
        </UCard>
      </div>

      <p class="text-center text-muted mt-8">
        All packages include access to our member dashboard and digital certificate upon completion.
      </p>
    </UPageSection>

    <!-- Comparison Table -->
    <UPageSection
      id="comparison"
      headline="Detailed Comparison"
      title="Compare All Features"
      description="A side-by-side comparison of all package features to help you decide."
      class="bg-muted/30"
    >
      <UCard class="overflow-x-auto">
        <table class="w-full min-w-[600px]">
          <thead>
            <tr class="border-b border-default">
              <th class="text-left py-4 px-4 font-semibold">Feature</th>
              <th class="text-center py-4 px-4 font-semibold">Starter</th>
              <th class="text-center py-4 px-4 font-semibold">
                <div class="flex items-center justify-center gap-2">
                  Standard
                  <UBadge label="Popular" size="xs" color="primary" />
                </div>
              </th>
              <th class="text-center py-4 px-4 font-semibold">Pro</th>
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
                <template v-if="typeof feature.starter === 'boolean'">
                  <UIcon 
                    :name="feature.starter ? 'i-lucide-check' : 'i-lucide-minus'" 
                    :class="feature.starter ? 'text-primary' : 'text-muted'" 
                    class="size-5"
                  />
                </template>
                <template v-else>
                  <span class="text-sm">{{ feature.starter }}</span>
                </template>
              </td>
              <td class="py-3 px-4 text-center bg-primary/5">
                <template v-if="typeof feature.standard === 'boolean'">
                  <UIcon 
                    :name="feature.standard ? 'i-lucide-check' : 'i-lucide-minus'" 
                    :class="feature.standard ? 'text-primary' : 'text-muted'" 
                    class="size-5"
                  />
                </template>
                <template v-else>
                  <span class="text-sm font-medium">{{ feature.standard }}</span>
                </template>
              </td>
              <td class="py-3 px-4 text-center">
                <template v-if="typeof feature.pro === 'boolean'">
                  <UIcon 
                    :name="feature.pro ? 'i-lucide-check' : 'i-lucide-minus'" 
                    :class="feature.pro ? 'text-primary' : 'text-muted'" 
                    class="size-5"
                  />
                </template>
                <template v-else>
                  <span class="text-sm">{{ feature.pro }}</span>
                </template>
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-muted/50">
              <td class="py-4 px-4 font-semibold">Price</td>
              <td class="py-4 px-4 text-center font-semibold">Rp 1.500.000</td>
              <td class="py-4 px-4 text-center font-semibold text-primary bg-primary/10">Rp 2.800.000</td>
              <td class="py-4 px-4 text-center font-semibold">Rp 4.500.000</td>
            </tr>
          </tfoot>
        </table>
      </UCard>
    </UPageSection>

    <!-- Add-ons -->
    <UPageSection
      headline="Extras"
      title="Optional Add-ons"
      description="Enhance your learning experience with these additional services."
    >
      <div class="grid md:grid-cols-3 gap-6 max-w-4xl mx-auto">
        <UCard v-for="addon in addOns" :key="addon.title">
          <div class="flex items-start gap-4">
            <div class="p-3 rounded-lg bg-primary/10">
              <UIcon :name="addon.icon" class="size-6 text-primary" />
            </div>
            <div>
              <h3 class="font-semibold">{{ addon.title }}</h3>
              <p class="text-primary font-bold">{{ addon.price }}</p>
              <p class="text-sm text-muted mt-1">{{ addon.description }}</p>
            </div>
          </div>
        </UCard>
      </div>
    </UPageSection>

    <!-- Payment Info -->
    <UPageSection class="bg-muted/30">
      <UCard class="max-w-3xl mx-auto">
        <template #header>
          <div class="flex items-center gap-3">
            <UIcon name="i-lucide-credit-card" class="size-6 text-primary" />
            <h3 class="text-xl font-semibold">Payment Options</h3>
          </div>
        </template>

        <div class="grid sm:grid-cols-2 gap-6">
          <div>
            <h4 class="font-medium mb-3">Accepted Payment Methods</h4>
            <ul class="space-y-2 text-sm text-muted">
              <li class="flex items-center gap-2">
                <UIcon name="i-lucide-building" class="size-4" />
                Bank Transfer (BCA, Mandiri, BNI, BRI)
              </li>
              <li class="flex items-center gap-2">
                <UIcon name="i-lucide-smartphone" class="size-4" />
                Virtual Account (VA)
              </li>
              <li class="flex items-center gap-2">
                <UIcon name="i-lucide-wallet" class="size-4" />
                E-Wallet (GoPay, OVO, DANA)
              </li>
              <li class="flex items-center gap-2">
                <UIcon name="i-lucide-qr-code" class="size-4" />
                QRIS
              </li>
            </ul>
          </div>
          <div>
            <h4 class="font-medium mb-3">Payment Terms</h4>
            <ul class="space-y-2 text-sm text-muted">
              <li class="flex items-center gap-2">
                <UIcon name="i-lucide-check" class="size-4 text-primary" />
                Full payment before first session
              </li>
              <li class="flex items-center gap-2">
                <UIcon name="i-lucide-check" class="size-4 text-primary" />
                Installment available (contact us)
              </li>
              <li class="flex items-center gap-2">
                <UIcon name="i-lucide-check" class="size-4 text-primary" />
                Invoice provided for all transactions
              </li>
            </ul>
          </div>
        </div>

        <template #footer>
          <div class="flex flex-wrap gap-3">
            <NuxtLink to="/register">
              <UButton label="Register & Pay" icon="i-lucide-arrow-right" />
            </NuxtLink>
            <NuxtLink to="https://wa.me/6281234567890" target="_blank">
              <UButton label="Ask About Installments" icon="i-simple-icons-whatsapp" variant="outline" color="neutral" />
            </NuxtLink>
          </div>
        </template>
      </UCard>
    </UPageSection>

    <!-- CTA -->
    <UPageCTA
      title="Questions About Packages?"
      description="Our team is ready to help you choose the right package for your needs."
      :links="[
        { label: 'Register Now', to: '/register', icon: 'i-lucide-user-plus' },
        { label: 'Contact Us', to: 'https://wa.me/6281234567890', color: 'neutral', variant: 'ghost', icon: 'i-simple-icons-whatsapp', external: true }
      ]"
    />
  </div>
</template>
