<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const navItems = computed<NavigationMenuItem[]>(() => [
  {
    label: 'Overview',
    icon: 'i-lucide-layout-dashboard',
    to: '/admin'
  },
  {
    label: 'Students',
    icon: 'i-lucide-users',
    to: '/admin/students'
  },
  {
    label: 'Schedules',
    icon: 'i-lucide-calendar',
    to: '/admin/schedules'
  },
  {
    label: 'Packages',
    icon: 'i-lucide-package',
    to: '/admin/packages'
  },
  {
    label: 'Certificates',
    icon: 'i-lucide-award',
    to: '/admin/certificates'
  },
  {
    label: 'Content',
    icon: 'i-lucide-file-text',
    to: '/admin/content'
  },
  {
    label: 'Settings',
    icon: 'i-lucide-settings',
    to: '/admin/settings'
  }
])

const adminMenuItems = [
  [
    { label: 'Admin Settings', icon: 'i-lucide-settings', to: '/admin/settings' },
    { label: 'View Website', icon: 'i-lucide-external-link', to: '/', external: true }
  ],
  [
    { label: 'Sign Out', icon: 'i-lucide-log-out', to: '/login' }
  ]
]

const admin = {
  name: 'Admin User',
  email: 'admin@evdriveacademy.id',
  role: 'Administrator'
}
</script>

<template>
  <UDashboardGroup>
    <UDashboardSidebar collapsible resizable>
      <template #header="{ collapsed }">
        <NuxtLink to="/admin" class="flex items-center gap-2 px-2">
          <div class="p-1.5 rounded-lg bg-primary">
            <UIcon name="i-lucide-zap" class="size-5 text-primary-foreground" />
          </div>
          <span v-if="!collapsed" class="font-bold">Admin Panel</span>
        </NuxtLink>
      </template>

      <template #default="{ collapsed }">
        <UNavigationMenu
          :items="navItems"
          orientation="vertical"
          :ui="{ link: collapsed ? 'justify-center' : undefined }"
        />
      </template>

      <template #footer="{ collapsed }">
        <UDropdownMenu :items="adminMenuItems">
          <UButton 
            color="neutral" 
            variant="ghost" 
            class="w-full"
            :class="collapsed ? 'justify-center px-0' : ''"
          >
            <UAvatar text="AD" size="sm" class="bg-primary text-primary-foreground" />
            <template v-if="!collapsed">
              <div class="flex-1 text-left ml-2">
                <p class="text-sm font-medium truncate">{{ admin.name }}</p>
                <p class="text-xs text-muted truncate">{{ admin.role }}</p>
              </div>
              <UIcon name="i-lucide-chevrons-up-down" class="size-4 text-muted" />
            </template>
          </UButton>
        </UDropdownMenu>
      </template>
    </UDashboardSidebar>

    <slot />
  </UDashboardGroup>
</template>
