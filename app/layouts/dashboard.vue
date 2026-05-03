<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import { computed } from 'vue'

const navItems = computed<NavigationMenuItem[]>(() => [
  {
    label: 'Dashboard',
    icon: 'i-lucide-layout-dashboard',
    to: '/dashboard'
  },
  {
    label: 'Schedule',
    icon: 'i-lucide-calendar',
    to: '/dashboard/schedule'
  },
  {
    label: 'Training History',
    icon: 'i-lucide-history',
    to: '/dashboard/history'
  },
  {
    label: 'Certificate',
    icon: 'i-lucide-award',
    to: '/dashboard/certificate'
  },
  {
    label: 'Billing',
    icon: 'i-lucide-credit-card',
    to: '/dashboard/billing'
  },
  {
    label: 'Profile',
    icon: 'i-lucide-user',
    to: '/dashboard/profile'
  }
])

const userMenuItems = [
  [
    { label: 'Profile Settings', icon: 'i-lucide-settings', to: '/dashboard/profile' },
    { label: 'Help & Support', icon: 'i-lucide-help-circle', to: 'https://wa.me/6281234567890', external: true }
  ],
  [
    { label: 'Sign Out', icon: 'i-lucide-log-out', to: '/auth/login' }
  ]
]

// Mock user data
const user = {
  name: 'John Doe',
  email: 'john.doe@example.com',
  avatar: 'JD'
}
</script>

<template>
  <UDashboardGroup>
    <UDashboardSidebar 
      collapsible 
      resizable
      :min-size="20"
      :default-size="25"
      :max-size="30"
      :collapsed-size="0"
    >
      <template #header="{ collapsed }">
        <NuxtLink v-if="!collapsed" to="/dashboard" class="flex items-center gap-2 px-2">
          <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
        </NuxtLink>
        <UIcon v-else name="i-simple-icons-nuxtdotjs" class="size-5 text-primary mx-auto" />
        <UDashboardSidebarCollapse variant="subtle" />
      </template>

      <template #default="{ collapsed }">
        <UNavigationMenu
          :items="navItems"
          orientation="vertical"
          color="warning"
          :ui="{ link: collapsed ? 'justify-center' : undefined, }"
        />
      </template>

      <template #footer="{ collapsed }">
        <UDropdownMenu :items="userMenuItems">
          <UButton 
            color="neutral" 
            variant="ghost" 
            class="w-full"
            :class="collapsed ? 'justify-center px-0' : ''"
          >
            <UAvatar :text="user.avatar" size="sm" />
            <template v-if="!collapsed">
              <div class="flex-1 text-left ml-2">
                <p class="text-sm font-medium truncate">{{ user.name }}</p>
                <p class="text-xs text-muted truncate">{{ user.email }}</p>
              </div>
              <UIcon name="i-lucide-chevrons-up-down" class="size-4 text-muted" />
            </template>
          </UButton>
        </UDropdownMenu>
      </template>
    </UDashboardSidebar>

    <UDashboardPanel>
      <slot />
    </UDashboardPanel>
  </UDashboardGroup>
</template>
