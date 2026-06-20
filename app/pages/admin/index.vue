<script setup lang="ts">
import { h, ref, onMounted, computed } from "vue";
import type { TableColumn } from "@nuxt/ui";
import { useDashboardStore } from "~/stores/dashboard";

const { t } = useI18n()
definePageMeta({ layout: "admin" });

const dashboardStore = useDashboardStore();
const UBadge = resolveComponent("UBadge");
const UButton = resolveComponent("UButton");
const UAvatar = resolveComponent("UAvatar");
const UDropdownMenu = resolveComponent("UDropdownMenu");

// Fetch dashboard data on mount
onMounted(async () => {
  await dashboardStore.fetchDashboardData();
});

// Computed stats from dashboard store
const stats = computed(() => {
  const dashboardStats = dashboardStore.stats;
  return [
    {
      label: t('admin.totalStudents'),
      value: dashboardStats.totalMembers.toString() || "0",
      change: `${dashboardStats.growthTotalMembers * 100}%`,
      icon: "i-lucide-users",
      color: "primary",
    },
    {
      label: t('admin.activeSessions'),
      value: dashboardStats.activeSessions.toString() || "0",
      change: "+5%",
      icon: "i-lucide-calendar-check",
      color: "blue",
    },
    {
      label: t('admin.revenueMtd'),
      value: formatCurrency(dashboardStats.revenueMTD),
      change: "+18%",
      icon: "i-lucide-banknote",
      color: "amber",
    },
    {
      label: t('admin.certsIssued'),
      value: dashboardStats.certificatesIssued.toString() || "0",
      change: "+8%",
      icon: "i-lucide-award",
      color: "primary",
    },
  ];
});

// Helper to format currency
const formatCurrency = (value: number): string => {
  if (value >= 1000000) {
    return `Rp ${(value / 1000000).toFixed(1)}M`;
  } else if (value >= 1000) {
    return `Rp ${(value / 1000).toFixed(0)}K`;
  }
  return `Rp ${value}`;
};

// Mock today's sessions
type Session = {
  id: number;
  time: string;
  student: string;
  car: string;
  instructor: string;
  status: "completed" | "in-progress" | "upcoming";
};

const todaySessions = ref<Session[]>([
  {
    id: 1,
    time: "08:00",
    student: "John Doe",
    car: "BYD Atto 1",
    instructor: "Mr. Ahmad",
    status: "completed",
  },
  {
    id: 2,
    time: "09:30",
    student: "Jane Smith",
    car: "BYD Atto 1",
    instructor: "Ms. Sari",
    status: "in-progress",
  },
  {
    id: 3,
    time: "11:00",
    student: "Alex Johnson",
    car: "BYD Atto 1",
    instructor: "Mr. Budi",
    status: "upcoming",
  },
  {
    id: 4,
    time: "13:00",
    student: "Maria Garcia",
    car: "BYD Atto 1",
    instructor: "Mr. Ahmad",
    status: "upcoming",
  },
]);

const sessionColumns: TableColumn<Session>[] = [
  { accessorKey: "time", header: t('dashboard.time') },
  { accessorKey: "student", header: t('admin.students') },
  { accessorKey: "car", header: t('dashboard.vehicle') },
  { accessorKey: "instructor", header: t('dashboard.instructor') },
  {
    accessorKey: "status",
    header: t('billing.status'),
    cell: ({ row }) => {
      const status = row.getValue("status") as string;
      const color = status === "completed" ? "success" : status === "in-progress" ? "info" : "neutral";
      const label = status === "completed" ? t('common.completed') : status === "in-progress" ? "In Progress" : t('common.pending');
      return h(UBadge, { label, color, variant: "subtle", size: "md" });
    },
  },
];

// Recent registrations from store
const recentRegistrations = computed(() => {
  return dashboardStore.recentRegistrations.map((reg) => ({
    id: reg.userId,
    name: `${reg.firstName} ${reg.lastName}`,
    email: reg.email,
    package: reg.roleId === 2 ? "Member" : "Standard",
    date: reg.createdAt,
    status: "active" as const,
  }));
});

type Registration = {
  id: string;
  name: string;
  email: string;
  package: string;
  date: string;
  status: "pending" | "active";
};

const registrationColumns: TableColumn<Registration>[] = [
  {
    accessorKey: "name",
    header: t('admin.students'),
    cell: ({ row }) => {
      const name = row.getValue("name") as string;
      const initials = name
        .split(" ")
        .map((n: string) => n[0])
        .join("")
        .toUpperCase();
      return h("div", { class: "flex items-center gap-3" }, [
        h(UAvatar, { text: initials, size: "md" }),
        h("span", { class: "font-medium" }, name),
      ]);
    },
  },
  { accessorKey: "email", header: t('auth.email') },
  {
    accessorKey: "package",
    header: t('billing.package'),
    cell: ({ row }) => {
      const pkg = row.getValue("package") as string;
      return h(UBadge, { label: pkg, color: "neutral", variant: "subtle", size: "md" });
    },
  },
  { accessorKey: "date", header: t('register.form.startDate').replace(' (Opsional)', '').replace(' (Optional)', '') },
  {
    accessorKey: "status",
    header: t('billing.status'),
    cell: ({ row }) => {
      const status = row.getValue("status") as string;
      const color = status === "active" ? "info" : "warning";
      const label = status === "active" ? t('billing.active') : t('common.pending');
      return h(UBadge, { label, color, variant: "subtle", size: "md" });
    },
  },
  {
    id: "actions",
    header: "",
    cell: () => {
      const items = [
        [
          { label: t('dashboard.viewDetails'), icon: "i-lucide-eye" },
          { label: t('common.edit'), icon: "i-lucide-pencil" },
        ],
        [{ label: t('common.delete'), icon: "i-lucide-trash", color: "error" }],
      ];
      return h(UDropdownMenu, { items }, () =>
        h(UButton, {
          icon: "i-lucide-ellipsis-vertical",
          color: "neutral",
          variant: "ghost",
          size: "md",
        }),
      );
    },
  },
];

// Mock quick actions
const quickActions = computed(() => [
  { label: t('admin.addNew'), icon: "i-lucide-user-plus", to: "/admin/students" },
  {
    label: t('admin.schedules'),
    icon: "i-lucide-calendar-plus",
    to: "/admin/schedules",
  },
  {
    label: t('admin.certificates'),
    icon: "i-lucide-file-badge",
    to: "/admin/certificates",
  },
]);
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.title')">
        <template #right>
          <UInput
            :placeholder="t('common.search') + '...'"
            icon="i-lucide-search"
            color="warning"
            class="w-64 hidden md:flex"
          />
          <UButton icon="i-lucide-bell" color="neutral" variant="ghost">
            <UChip color="error" size="md">3</UChip>
          </UButton>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Stats Grid -->
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
          <UCard v-for="stat in stats" :key="stat.label">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-md text-muted">{{ stat.label }}</p>
                <p class="text-2xl font-bold mt-1">{{ stat.value }}</p>
                <p class="text-md text-green-500 mt-1">
                  {{ stat.change }} {{ t('admin.fromLastMonth') }}
                </p>
              </div>
              <div
                class="p-3 rounded-xl"
                :class="{
                  'bg-primary/10': stat.color === 'primary',
                  'bg-green-500/10': stat.color === 'success',
                  'bg-amber-500/10': stat.color === 'amber',
                  'bg-blue-500/10': stat.color === 'blue',
                }"
              >
                <UIcon
                  :name="stat.icon"
                  class="size-6"
                  :class="{
                    'text-primary': stat.color === 'primary',
                    'text-green-500': stat.color === 'success',
                    'text-amber-500': stat.color === 'amber',
                    'text-blue-500': stat.color === 'blue',
                  }"
                />
              </div>
            </div>
          </UCard>
        </div>

        <div class="grid lg:grid-cols-3 gap-6">
          <!-- Today's Sessions -->
          <UCard class="lg:col-span-2">
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">{{ t('admin.todaysSessions') }}</h2>
                <NuxtLink to="/admin/schedules">
                  <UButton
                    :label="t('common.viewAll')"
                    color="neutral"
                    variant="ghost"
                    size="md"
                    trailing-icon="i-lucide-arrow-right"
                  />
                </NuxtLink>
              </div>
            </template>

            <UTable :data="todaySessions" :columns="sessionColumns" />
          </UCard>

          <!-- Quick Actions -->
          <UCard>
            <template #header>
              <h2 class="font-semibold">{{ t('dashboard.quickActions') }}</h2>
            </template>

            <div class="grid grid-cols-1 space-y-3">
              <NuxtLink
                v-for="action in quickActions"
                :key="action.label"
                :to="action.to"
              >
                <UButton
                  :label="action.label"
                  :icon="action.icon"
                  color="warning"
                  variant="outline"
                  block
                />
              </NuxtLink>
            </div>

            <template #footer>
              <div class="space-y-2">
                <h3 class="font-medium text-md">{{ t('dashboard.vehicle') }} Status</h3>
                <div
                  class="flex items-center justify-between p-2 rounded-lg bg-muted/50"
                >
                  <span class="text-md">BYD Atto 1</span>
                  <UBadge
                    :label="t('common.available')"
                    color="primary"
                    variant="subtle"
                    size="md"
                  />
                </div>
              </div>
            </template>
          </UCard>
        </div>

        <!-- Recent Registrations -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold">{{ t('admin.recentReg') }}</h2>
              <NuxtLink to="/admin/students">
                <UButton
                  :label="t('admin.manageStudents')"
                  color="neutral"
                  variant="ghost"
                  size="md"
                  trailing-icon="i-lucide-arrow-right"
                />
              </NuxtLink>
            </div>
          </template>

          <UTable :data="recentRegistrations" :columns="registrationColumns" />
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
