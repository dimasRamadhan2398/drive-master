<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import { computed } from 'vue';
import { useAuth } from '~/composables/useAuth';

const { pages } = useContent()

const { user, isLoggedIn, logout } = useAuth()

const navItems = computed<NavigationMenuItem[]>(() => {
  const baseItems = [
    { label: 'Home', to: '/' },
    { label: 'Services', to: '/services' },
    { label: 'Packages', to: '/packages' },
    { label: 'Instructors', to: '/instructors' },
    { label: 'Article', to: '/blog' },
    { label: 'About Us', to: '/about' },
    { label: 'Contact Us', to: '/contact' },
  ]

  const dynamicItems = pages.value
    .filter(p => p.status === 'published' && 
    p.slug !== '/' && 
    p.slug !== '/services' && 
    p.slug !== '/packages' && 
    p.slug !== '/instructors' && 
    p.slug !== '/blog' && 
    p.slug !== '/about' && 
    p.slug !== '/contact')
    .map(p => ({ label: p.title, to: p.slug }))

  return [...baseItems, ...dynamicItems]
})

const userMenuItems = computed(() => [
  [
    {
      label: user.value?.name || 'Member',
      slot: 'account',
      disabled: true
    }
  ],
  [
    {
      label: 'Dashboard',
      icon: 'i-lucide-layout-dashboard',
      to: '/dashboard'
    },
    {
      label: 'Settings',
      icon: 'i-lucide-settings',
      to: '/dashboard/profile'
    }
  ],
  [
    {
      label: 'Logout',
      icon: 'i-lucide-log-out',
      click: () => logout()
    }
  ]
])
</script>

<template>
  <div>
    <UHeader>
      <template #title>
        <NuxtLink to="/" class="flex items-center gap-2">
          <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
        </NuxtLink>
      </template>

      <UNavigationMenu :items="navItems" color="warning" class="hidden lg:flex" />

      <template #right>
        <UColorModeButton />
        
        <template v-if="isLoggedIn">
          <UDropdownMenu :items="userMenuItems" :ui="{ content: 'w-48' }">
            <UButton variant="ghost" color="neutral" class="p-0.5 rounded-full">
              <UAvatar 
                :src="user?.avatar" 
                :alt="user?.name" 
                size="sm"
                class="ring-2 ring-warning/20"
              />
            </UButton>
            
            <template #account="{ item }">
              <div class="text-left">
                <p class="font-medium text-gray-900 dark:text-white truncate">
                  {{ user?.name }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400 truncate">
                  {{ user?.email }}
                </p>
              </div>
            </template>
          </UDropdownMenu>
        </template>
        <template v-else>
          <NuxtLink to="/auth/login">
            <UButton label="Login" color="warning" variant="ghost" class="hidden sm:flex" />
          </NuxtLink>
          <NuxtLink to="/auth/register">
            <UButton label="Register" color="warning" />
          </NuxtLink>
        </template>
      </template>

      <template #body>
        <UNavigationMenu :items="navItems" orientation="vertical" class="-mx-2.5" />
        
        <div v-if="isLoggedIn" class="mt-4 pt-4 border-t border-default space-y-4">
          <div class="flex items-center gap-3 px-2">
            <UAvatar :src="user?.avatar" :alt="user?.name" size="md" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium truncate">{{ user?.name }}</p>
              <p class="text-xs text-muted truncate">{{ user?.email }}</p>
            </div>
          </div>
          
          <div class="flex flex-col gap-2">
            <NuxtLink to="/dashboard">
              <UButton label="Dashboard" icon="i-lucide-layout-dashboard" color="warning" variant="ghost" block />
            </NuxtLink>
            <UButton label="Logout" icon="i-lucide-log-out" color="error" variant="ghost" block @click="logout" />
          </div>
        </div>
        <div v-else class="flex flex-col gap-2 mt-4 pt-4 border-t border-default">
          <NuxtLink to="/auth/login">
            <UButton label="Login" color="warning" variant="ghost" block />
          </NuxtLink>
          <NuxtLink to="/auth/register">
            <UButton label="Register" color="warning" block />
          </NuxtLink>
        </div>
      </template>
    </UHeader>

    <UMain>
      <slot />
    </UMain>

    <UFooter class="border-t border-default">
      <template #left>
        <div class="flex items-center gap-2">
          <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-12" />
        </div>
        <p class="text-muted text-sm p-10">
          The first premium driving academy in Alam Sutera using 100% Electric Vehicles.
        </p>
      </template>
      <template #right>
        <div class="flex flex-col items-end gap-2">
          <div class="flex gap-2">
            <UButton icon="i-simple-icons-instagram" color="warning" variant="ghost" to="#" target="_blank" aria-label="Instagram" />
            <UButton icon="i-simple-icons-whatsapp" color="warning" variant="ghost" to="https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi" target="_blank" aria-label="WhatsApp" />
            <UButton icon="i-simple-icons-youtube" color="warning" variant="ghost" to="#" target="_blank" aria-label="YouTube" />
          </div>
          <p class="text-muted text-sm">
            © {{ new Date().getFullYear() }} PT Drive Master Indonesia. All rights reserved.
          </p>
        </div>
      </template>
    </UFooter>

    <!-- WhatsApp Floating Button -->
    <NuxtLink 
      to="https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi" 
      target="_blank"
      class="fixed bottom-6 right-6 z-50"
    >
      <UButton 
        icon="i-simple-icons-whatsapp" 
        size="xl"
        class="rounded-full size-14 !bg-[#25D366] hover:!bg-[#128C7E] text-white shadow-lg"
        aria-label="Chat on WhatsApp"
      />
    </NuxtLink>
  </div>
</template>
