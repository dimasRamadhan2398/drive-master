<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const navItems = computed<NavigationMenuItem[]>(() => [
  { label: 'Services', to: '/services' },
  { label: 'Packages', to: '/packages' },
  { label: 'About', to: '/#about' },
  { label: 'FAQ', to: '/#faq' },
  { label: 'Contact', to: '/#contact' }
])
</script>

<template>
  <div>
    <UHeader>
      <template #title>
        <NuxtLink to="/" class="flex items-center gap-2">
          <UIcon name="i-lucide-zap" class="size-7 text-primary" />
          <span class="font-bold text-xl">EV Drive Academy</span>
        </NuxtLink>
      </template>

      <UNavigationMenu :items="navItems" class="hidden lg:flex" />

      <template #right>
        <UColorModeButton />
        <NuxtLink to="/login">
          <UButton label="Login" color="neutral" variant="ghost" class="hidden sm:flex" />
        </NuxtLink>
        <NuxtLink to="/register">
          <UButton label="Register" />
        </NuxtLink>
      </template>

      <template #body>
        <UNavigationMenu :items="navItems" orientation="vertical" class="-mx-2.5" />
        <div class="flex flex-col gap-2 mt-4 pt-4 border-t border-default">
          <NuxtLink to="/login">
            <UButton label="Login" color="neutral" variant="ghost" block />
          </NuxtLink>
          <NuxtLink to="/register">
            <UButton label="Register" block />
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
          <UIcon name="i-lucide-zap" class="size-5 text-primary" />
          <span class="font-semibold">EV Drive Academy</span>
        </div>
        <p class="text-muted text-sm mt-2">
          The first premium driving academy in Alam Sutera using 100% Electric Vehicles.
        </p>
      </template>
      <template #right>
        <div class="flex flex-col items-end gap-2">
          <div class="flex gap-2">
            <UButton icon="i-simple-icons-instagram" color="neutral" variant="ghost" to="https://instagram.com" target="_blank" aria-label="Instagram" />
            <UButton icon="i-simple-icons-whatsapp" color="neutral" variant="ghost" to="https://wa.me/6281234567890" target="_blank" aria-label="WhatsApp" />
            <UButton icon="i-simple-icons-youtube" color="neutral" variant="ghost" to="https://youtube.com" target="_blank" aria-label="YouTube" />
          </div>
          <p class="text-muted text-sm">
            © {{ new Date().getFullYear() }} EV Drive Academy. All rights reserved.
          </p>
        </div>
      </template>
    </UFooter>

    <!-- WhatsApp Floating Button -->
    <NuxtLink 
      to="https://wa.me/6281234567890?text=Halo%20EV%20Drive%20Academy%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi" 
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
