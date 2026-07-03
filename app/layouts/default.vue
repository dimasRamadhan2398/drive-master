<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import { computed, onMounted } from 'vue';
import { useAuth } from '~/composables/useAuth';
import { useLocale } from '~/composables/useLocale';

const { t, locale, locales, setLocale } = useI18n()
const { pages } = useContent()
const { user, isLoggedIn, logout } = useAuth()
const { switchLocale } = useLocale()
const { waLink, fetchGeneralSettings } = useSettings()

onMounted(() => {
  fetchGeneralSettings()
})

const navItems = computed<NavigationMenuItem[]>(() => {
  const baseItems = [
    { label: t('nav.home'), to: '/' },
    { label: t('nav.services'), to: '/services' },
    { label: t('nav.packages'), to: '/packages' },
    { label: t('nav.instructors'), to: '/instructors' },
    { label: t('nav.article'), to: '/blog' },
    { label: t('nav.about'), to: '/about' },
    { label: t('nav.contact'), to: '/contact' },
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
      label: t('dashboard.welcome', { name: user.value?.name || '' }),
      icon: 'i-lucide-layout-dashboard',
      to: '/dashboard'
    },
    {
      label: t('common.settings'),
      icon: 'i-lucide-settings',
      to: '/dashboard/profile'
    }
  ],
  [
    {
      label: t('auth.logout'),
      icon: 'i-lucide-log-out',
      click: () => logout()
    }
  ]
])

const languageOptions = computed(() => {
  return (locales.value as any[]).map(l => ({
    label: l.name || l.code.toUpperCase(),
    code: l.code,
    icon: l.code === 'id' ? '🇮🇩' : '🇬🇧'
  }))
})

const currentLanguage = computed(() => {
  const current = (locales.value as any[]).find(l => l.code === locale.value)
  return current?.name || locale.value.toUpperCase()
})
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

        <!-- Language Switcher -->
        <UDropdownMenu :items="[languageOptions.map(l => ({
          label: `${l.icon} ${l.label}`,
          onClick: () => setLocale(l.code)
        }))]">
          <UButton variant="ghost" color="neutral" size="sm">
            <span class="text-lg">{{ locale === 'id' ? '🇮🇩' : '🇬🇧' }}</span>
            <UIcon name="i-lucide-chevron-down" class="size-4 ml-1" />
          </UButton>
        </UDropdownMenu>

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
            <UButton :label="t('auth.login')" color="warning" variant="ghost" class="hidden sm:flex" />
          </NuxtLink>
          <NuxtLink to="/auth/register">
            <UButton :label="t('auth.register')" color="warning" />
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
              <UButton :label="t('common.dashboard')" icon="i-lucide-layout-dashboard" color="warning" variant="ghost" block />
            </NuxtLink>
            <UButton :label="t('auth.logout')" icon="i-lucide-log-out" color="error" variant="ghost" block @click="logout" />
          </div>
        </div>
        <div v-else class="flex flex-col gap-2 mt-4 pt-4 border-t border-default">
          <NuxtLink to="/auth/login">
            <UButton :label="t('auth.login')" color="warning" variant="ghost" block />
          </NuxtLink>
          <NuxtLink to="/auth/register">
            <UButton :label="t('auth.register')" color="warning" block />
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
          <img src="/drive-master-logo-light.png" alt="Drive Master Logo" class="h-12 dark:hidden" />
          <img src="/drive-master-logo-dark.jpg" alt="Drive Master Logo" class="h-12 hidden dark:block" />
        </div>
        <p class="text-muted text-sm pl-2">
          {{ t('footer.tagline') }}
        </p>
      </template>
      <template #right>
        <div class="flex flex-col lg:items-end gap-2 items-center">
          <div class="flex items-center gap-2">
            <UButton icon="i-simple-icons-instagram" color="warning" variant="ghost" to="#" target="_blank" aria-label="Instagram" />
            <UButton icon="i-simple-icons-whatsapp" color="warning" variant="ghost" :to="waLink" target="_blank" aria-label="WhatsApp" />
            <UButton icon="i-simple-icons-youtube" color="warning" variant="ghost" to="#" target="_blank" aria-label="YouTube" />
          </div>
          <p class="text-muted text-sm">
            {{ t('footer.copyright', { year: new Date().getFullYear() }) }}
          </p>
        </div>
      </template>
    </UFooter>

    <!-- WhatsApp Floating Button -->
    <NuxtLink 
      :to="waLink" 
      target="_blank"
      class="fixed bottom-6 right-6 z-50"
    >
      <UButton 
        icon="i-simple-icons-whatsapp" 
        size="xl"
        class="rounded-full size-14 !bg-[#25D366] hover:!bg-[#128C7E] text-white shadow-lg flex items-center justify-center"
        aria-label="Chat on WhatsApp"
      />
    </NuxtLink>
  </div>
</template>
