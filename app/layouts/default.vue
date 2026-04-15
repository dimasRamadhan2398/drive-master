<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import { computed } from 'vue';

const navItems = computed<NavigationMenuItem[]>(() => [
  { label: 'Home', to: '/' },
  { label: 'Services', to: '/services' },
  { label: 'Packages', to: '/packages' },
  { label: 'About Us', to: '/about' },
  { label: 'Contact Us', to: '/contact' }
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
        <NuxtLink to="/login">
          <UButton label="Login" color="warning" variant="ghost" class="hidden sm:flex" />
        </NuxtLink>
        <NuxtLink to="/register">
          <UButton label="Register" color="warning" />
        </NuxtLink>
      </template>

      <template #body>
        <UNavigationMenu :items="navItems" orientation="vertical" class="-mx-2.5" />
        <div class="flex flex-col gap-2 mt-4 pt-4 border-t border-default">
          <NuxtLink to="/login">
            <UButton label="Login" color="warning" variant="ghost" block />
          </NuxtLink>
          <NuxtLink to="/register">
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
