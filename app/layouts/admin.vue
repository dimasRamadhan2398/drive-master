<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";

const { t } = useI18n();
const authStore = useAuthStore();
const navItems = computed(() => [
  {
    label: t("admin.overview"),
    icon: "i-lucide-layout-dashboard",
    to: "/admin",
  },
  {
    label: t("admin.students"),
    icon: "i-lucide-users",
    to: "/admin/students",
  },
  {
    label: t("admin.schedules"),
    icon: "i-lucide-calendar",
    to: "/admin/schedules",
  },
  {
    label: t("admin.packages"),
    icon: "i-lucide-package",
    to: "/admin/packages",
  },
  {
    label: t("admin.sales"),
    icon: "i-lucide-package",
    to: "/admin/sales",
  },
  {
    label: t("admin.certificates"),
    icon: "i-lucide-award",
    to: "/admin/certificates",
  },
  {
    label: t("admin.contents"),
    icon: "i-lucide-file-text",
    to: "/admin/content",
  },
  {
    label: t("admin.testimonials"),
    icon: "i-lucide-message-square",
    to: "/admin/testimonials",
  },
  {
    label: t("admin.inquiries"),
    icon: "i-lucide-inbox",
    to: "/admin/inquiries",
  },
  {
    label: t("admin.analytics"),
    icon: "i-lucide-bar-chart-3",
    to: "/admin/analytics",
  },
  {
    label: t("admin.settings"),
    icon: "i-lucide-settings",
    to: "/admin/settings",
  },
]);

const adminMenuItems = computed(() => [
  [
    {
      label: t("admin.settings"),
      icon: "i-lucide-settings",
      to: "/admin/settings",
    },
    {
      label: t("admin.viewWebsite"),
      icon: "i-lucide-external-link",
      to: "/",
      external: true,
    },
  ],
  [
    {
      label: t("admin.signOut"),
      icon: "i-lucide-log-out",
      to: "/admin/login",
      onClick: () => authStore.logout(),
    },
  ],
]);

const admin = {
  name: "Admin User",
  email: "admin@evdriveacademy.id",
  role: "Administrator",
};

// FITUR BARU: Mengaktifkan monitor Smart Alert untuk seluruh halaman Admin
const { startMonitor } = useSmartAlerts();
onMounted(() => {
  startMonitor();
});
</script>

<template>
  <UDashboardGroup>
    <UDashboardSidebar collapsible resizable>
      <template #header="{ collapsed }">
        <NuxtLink to="/admin" class="flex items-center gap-4 py-8">
          <img src="/drive-master-logo-light.png" alt="Drive Master Logo" class="h-10 dark:hidden" />
          <img src="/drive-master-logo-dark.jpg" alt="Drive Master Logo" class="h-10 hidden dark:block" />
          <span v-if="!collapsed" class="font-bold py-16">{{ t("admin.title") }}</span>
        </NuxtLink>
      </template>

      <template #default="{ collapsed }">
        <UNavigationMenu
          :items="navItems"
          orientation="vertical"
          color="warning"
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
            <UAvatar text="AD" size="sm" class="bg-warning text-warning-foreground" />
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
