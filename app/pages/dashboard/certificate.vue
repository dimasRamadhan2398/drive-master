<script setup lang="ts">
import { certificateService } from "~/services/certificateService";

const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

const authStore = useAuthStore()
const isLoading = ref(false)
const issuedCertificates = ref<any[]>([])

const activeEntitlements = computed(() => {
  return authStore.memberEntitlements.filter(e => e.remaining > 0);
});

const formatDate = (dateStr: string) => {
  if (!dateStr) return "";
  const date = new Date(dateStr);
  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
};

const loadCertificates = async () => {
  if (!authStore.userId) return;
  try {
    isLoading.value = true;
    const certs = await certificateService.getMemberCertificates(authStore.userId);
    issuedCertificates.value = certs.map(c => ({
      id: c.id,
      title: c.packageName || "Driving Certification",
      recipientName: c.memberName,
      issuedDate: formatDate(c.issuedDate || c.completedAt),
      status: c.status,
      downloadUrl: "#",
      memberId: c.memberId
    }));
  } catch (err) {
    console.error("Failed to load certificates:", err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
  if (!authStore.hasMemberProfile) {
    await authStore.fetchMemberProfile();
  }
  await loadCertificates();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('certificate.title')">
        <template #right>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Loader -->
        <div v-if="isLoading" class="space-y-4">
          <UCard v-for="n in 2" :key="n">
            <div class="flex items-center gap-4">
              <USkeleton class="h-12 w-12 rounded-xl" />
              <div class="flex-1 space-y-2">
                <USkeleton class="h-5 w-[40%]" />
                <USkeleton class="h-4 w-[70%]" />
              </div>
            </div>
          </UCard>
        </div>

        <template v-else>
          <!-- Active Courses Progress -->
          <div v-if="activeEntitlements.length > 0" class="space-y-4">
            <h2 class="text-lg font-semibold">{{ t('certificate.progress') }}</h2>
            <div class="grid md:grid-cols-2 gap-6">
              <UCard v-for="entitlement in activeEntitlements" :key="entitlement.id">
                <template #header>
                  <div class="flex items-center justify-between">
                    <h3 class="font-bold text-md text-primary">{{ entitlement.packageName || 'Driving Course' }}</h3>
                    <UBadge label="In Progress" color="warning" variant="subtle" />
                  </div>
                </template>
                
                <div class="space-y-4 py-2">
                  <div class="flex justify-between text-sm">
                    <span class="text-muted">{{ t('common.progress') }}</span>
                    <span class="font-medium">
                      {{ entitlement.totalSessions - entitlement.remaining }}/{{ entitlement.totalSessions }} {{ t('billing.sessions') }}
                    </span>
                  </div>
                  <UProgress :model-value="((entitlement.totalSessions - entitlement.remaining) / entitlement.totalSessions) * 100" color="primary" />
                  <p class="text-sm text-muted">
                    {{ t('certificate.remainingSessions', { count: entitlement.remaining }) }}
                  </p>
                </div>

                <template #footer>
                  <div class="flex justify-end">
                    <NuxtLink to="/dashboard/schedule">
                      <UButton :label="t('dashboard.bookNextSession')" icon="i-lucide-calendar-plus" size="sm" />
                    </NuxtLink>
                  </div>
                </template>
              </UCard>
            </div>
          </div>

          <!-- Earned Certificates -->
          <div v-if="issuedCertificates.length > 0" class="space-y-4">
            <h2 class="text-lg font-semibold">{{ t('certificate.title') }}</h2>
            <div class="grid md:grid-cols-2 gap-6">
              <template v-for="certificate in issuedCertificates" :key="certificate.id">
                <CertificateDisplay :certificate="certificate" />
              </template>
            </div>
          </div>

          <!-- Fallback Empty State -->
          <div v-if="activeEntitlements.length === 0 && issuedCertificates.length === 0" class="text-center py-12">
            <div class="mx-auto w-20 h-20 rounded-full bg-muted/50 flex items-center justify-center mb-4">
              <UIcon name="i-lucide-award" class="size-10 text-muted" />
            </div>
            <h3 class="text-xl font-bold mb-2">{{ t('certificate.notAvailable') }}</h3>
            <p class="text-muted max-w-md mx-auto">
              No course enrollments or certificates found. Book your first driving lesson package to start your learning journey!
            </p>
            <div class="mt-6">
              <NuxtLink to="/dashboard/schedule">
                <UButton label="Book Package" icon="i-lucide-calendar-plus" />
              </NuxtLink>
            </div>
          </div>
        </template>

        <!-- Certificate Info -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">{{ t('certificate.aboutTitle') }}</h2>
          </template>

          <div class="grid md:grid-cols-3 gap-6">
            <div class="text-center">
              <div class="mx-auto w-14 h-14 rounded-full bg-warning/10 flex items-center justify-center mb-3">
                <UIcon name="i-lucide-file-badge" class="size-7 text-warning" />
              </div>
              <h3 class="font-semibold mb-1">{{ t('certificate.officialRecognition') }}</h3>
              <p class="text-md text-muted">{{ t('certificate.officialRecognitionDesc') }}</p>
            </div>
            <div class="text-center">
              <div class="mx-auto w-14 h-14 rounded-full bg-warning/10 flex items-center justify-center mb-3">
                <UIcon name="i-lucide-qr-code" class="size-7 text-warning" />
              </div>
              <h3 class="font-semibold mb-1">{{ t('certificate.digitalVerification') }}</h3>
              <p class="text-md text-muted">{{ t('certificate.digitalVerificationDesc') }}</p>
            </div>
            <div class="text-center">
              <div class="mx-auto w-14 h-14 rounded-full bg-warning/10 flex items-center justify-center mb-3">
                <UIcon name="i-lucide-infinity" class="size-7 text-warning" />
              </div>
              <h3 class="font-semibold mb-1">{{ t('certificate.lifetimeValidity') }}</h3>
              <p class="text-md text-muted">{{ t('certificate.lifetimeValidityDesc') }}</p>
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
