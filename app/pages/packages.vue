<script setup lang="ts">
const packages = [
  {
    title: '6x',
    price: 'Rp 1.500.000',
    // priceNote: 'One-time payment',
    description: 'Perfect for beginners looking to get started with driving basics.',
    features: [
      { text: 'Free Trial', included: true },
      { text: '6x Sessions', included: true },
    ],
    highlight: false,
    color: 'neutral' as const
  },
  {
    title: '8x',
    price: 'Rp 1.950.000',
    // priceNote: 'One-time payment',
    description: 'Our most popular package for comprehensive learning.',
    features: [
      { text: 'Free Trial', included: true },
      { text: '8x Sessions', included: true },
    ],
    highlight: true,
    color: 'warning' as const
  },
  {
    title: '10x',
    price: 'Rp 2.250.000',
    // priceNote: 'One-time payment',
    description: 'Complete mastery with unlimited support and premium perks.',
    features: [
      { text: 'Free Trial', included: true },
      { text: '10x Sessions', included: true },
    ],
    highlight: false,
    color: 'neutral' as const
  },
  {
    title: '12x',
    price: 'Rp 2.650.000',
    // priceNote: 'One-time payment',
    description: 'Complete mastery with unlimited support and premium perks.',
    features: [
      { text: 'Free Trial', included: true },
      { text: '12x Sessions', included: true },
    ],
    highlight: false,
    color: 'neutral' as const
  }
]

const comparisonFeatures = [
  { name: 'Free Trial', pkg6x: true, pkg8x: true, pkg10x: true, pkg12x: true },
  { name: 'Sessions', pkg6x: '6', pkg8x: '8', pkg10x: '10', pkg12x: '12' },
  { name: 'Total Hours', pkg6x: '6 hrs', pkg8x: '8 hrs', pkg10x: '10 hrs', pkg12x: '12 hrs' },
]

const addOns = [
  {
    title: 'Extra Session',
    price: 'Rp 350.000',
    description: 'Add more practice time to any package',
    icon: 'i-lucide-plus-circle'
  },
]
</script>

<template>
  <div>
    <!-- Hero -->
    <UPageHero
      title="Course Packages"
      description="Transparent pricing with no hidden fees. Choose the package that best fits your learning goals and schedule."
      :links="[
        { label: 'Register Now', to: '/register', color: 'warning', icon: 'i-lucide-user-plus' },
        { label: 'Compare Packages', to: '#comparison', color: 'neutral', variant: 'outline' }
      ]"
    />

    <!-- Package Cards -->
    <UPageSection :ui="{ headline: 'text-warning' }">
      <div class="grid md:grid-cols-4 gap-6 lg:gap-8">
        <UCard
          v-for="pkg in packages"
          :key="pkg.title"
          :class="[
            'flex flex-col relative',
            pkg.highlight ? 'ring-2 ring-warning shadow-xl' : ''
          ]"
        >
          <UBadge 
            v-if="pkg.highlight" 
            label="Most Popular" 
            color="warning" 
            class="absolute -top-0.5 left-1/2 -translate-x-1/2"
          />

          <template #header>
            <div class="text-center pt-2">
              <h3 class="text-2xl font-bold">{{ pkg.title }}</h3>
              <div class="mt-4">
                <span class="text-4xl font-bold text-warning">{{ pkg.price }}</span>
              </div>
              <!-- <p class="text-sm text-muted mt-1">{{ pkg.priceNote }}</p> -->
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
                :class="feature.included ? 'text-warning' : 'text-muted'" 
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
      :ui="{ headline: 'text-warning' }"
      class="bg-muted/30"
    >
      <UCard class="overflow-x-auto">
        <table class="w-full min-w-[600px]">
          <thead>
            <tr class="border-b border-default">
              <th class="text-left py-4 px-4 font-semibold">Feature</th>
              <th class="text-center py-4 px-4 font-semibold">6x</th>
              <th class="text-center py-4 px-4 font-semibold">
                <div class="flex items-center justify-center gap-2">
                  8x
                  <UBadge label="Popular" size="xs" color="warning" />
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
              <td class="py-4 px-4 font-semibold">Price</td>
              <td class="py-4 px-4 text-center font-semibold">Rp 1.750.000</td>
              <td class="py-4 px-4 text-center font-semibold text-warning bg-warning/10">Rp 1.950.000</td>
              <td class="py-4 px-4 text-center font-semibold">Rp 2.250.000</td>
              <td class="py-4 px-4 text-center font-semibold">Rp 2.650.000</td>
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
      title="Questions About Packages?"
      description="Our team is ready to help you choose the right package for your needs."
      :links="[
        { label: 'Register Now', to: '/register', color: 'warning', icon: 'i-lucide-user-plus' },
        { label: 'Contact Us', to: 'https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi', color: 'primary', variant: 'outline', icon: 'i-simple-icons-whatsapp', external: true }
      ]"
    />
  </div>
</template>
