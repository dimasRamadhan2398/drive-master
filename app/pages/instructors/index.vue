<script setup lang="ts">
useSeoMeta({
  title: 'Instructors | Drive Master Academy',
  description: 'Meet our professional instructors, certified and experienced in providing the best driving education.',
})

const instructors = [
  {
    name: 'Mr. Ahmad',
    phone: '081234567001',
    bnsp: 'BNSP-101-2023',
    sim: 'SIM A',
    photoUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&auto=format&fit=crop&q=80',
    experience: 10,
    bio: 'Expert in smooth vehicle control and traffic safety management with extensive knowledge in electric vehicle systems.',
    status: 'active',
    role: 'Senior Instructor',
    specialization: 'Defensive Driving & EV Specialist'
  },
  {
    name: 'Ms. Sari',
    phone: '081234567002',
    bnsp: 'BNSP-102-2022',
    sim: 'SIM A',
    photoUrl: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=400&auto=format&fit=crop&q=80',
    experience: 7,
    bio: 'Specializes in helping beginners master city driving and complex parking maneuvers with patience and precision.',
    status: 'active',
    role: 'Professional Instructor',
    specialization: 'Urban Driving & Parking Expert'
  },
  {
    name: 'Mr. Budi',
    phone: '081234567003',
    bnsp: 'BNSP-105-2024',
    sim: 'SIM A',
    photoUrl: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400&auto=format&fit=crop&q=80',
    experience: 8,
    bio: 'Focuses on high-speed safety, highway ethics, and advanced hazard perception for confident long-distance driving.',
    status: 'active',
    role: 'Safety Specialist',
    specialization: 'Advanced Safety & Highway Driving'
  },
  {
    name: 'Ms. Linda',
    phone: '081234567004',
    bnsp: 'BNSP-108-2023',
    sim: 'SIM A',
    photoUrl: 'https://images.unsplash.com/photo-1580489944761-15a19d654956?w=400&auto=format&fit=crop&q=80',
    experience: 6,
    bio: 'Known for her calm approach in building confidence for first-time drivers and mastering basic controls.',
    status: 'active',
    role: 'Senior Instructor',
    specialization: 'Beginner Basics & Confidence Building'
  }
]
</script>

<template>
  <UContainer class="py-12">
    <UPageHeader
      title="Meet Our Professional Instructors"
      description="Learn from the best. Our certified instructors are dedicated to providing a premium learning experience."
      class="mb-12"
    />

    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
      <UCard 
        v-for="instructor in instructors" 
        :key="instructor.name"
        class="overflow-hidden group"
      >
        <template #header>
          <div class="relative h-64 -m-6 mb-0 overflow-hidden">
            <img 
              :src="instructor.photoUrl" 
              :alt="instructor.name"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
            <div class="absolute bottom-4 left-4">
              <UBadge :label="instructor.role" color="warning" variant="subtle" />
            </div>
          </div>
        </template>

        <div class="space-y-4">
          <div>
            <h3 class="text-xl font-bold">{{ instructor.name }}</h3>
            <p class="text-warning font-medium text-sm">{{ instructor.specialization }}</p>
          </div>

          <p class="text-muted text-sm leading-relaxed">
            {{ instructor.bio }}
          </p>

          <div class="flex items-center gap-4 pt-2 border-t border-default">
            <div class="flex items-center gap-1">
              <UIcon name="i-lucide-clock" class="size-4 text-warning" />
              <span class="text-xs font-semibold">{{ instructor.experience }} Years</span>
            </div>
            <div class="flex items-center gap-1">
              <UIcon name="i-lucide-shield-check" class="size-4 text-warning" />
              <span class="text-xs font-semibold">Certified</span>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="space-y-2">
            <div class="flex items-center gap-2 text-sm">
              <UIcon name="i-lucide-phone" class="size-4 text-muted" />
              <span class="text-muted">{{ instructor.phone }}</span>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <UIcon name="i-lucide-award" class="size-4 text-muted" />
              <span class="text-muted">BNSP: {{ instructor.bnsp }}</span>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <UIcon name="i-lucide-credit-card" class="size-4 text-muted" />
              <span class="text-muted">License: {{ instructor.sim }}</span>
            </div>
          </div>
        </template>
      </UCard>
    </div>

    <UPageCTA
      title="Ready to start your journey?"
      description="Choose your favorite instructor and book your first session today."
      :links="[{ label: 'Register Now', to: '/auth/register', color: 'warning', size: 'lg' }]"
      class="mt-20 bg-warning/5"
    />
  </UContainer>
</template>
