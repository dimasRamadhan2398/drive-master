<script setup lang="ts">
const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

// Mock certificate data
const certificateStatus = ref({
  isEligible: false,
  progress: 40,
  completedSessions: 4,
  totalSessions: 10,
  remainingSessions: 6
})

// Mock issued certificates
const issuedCertificates = ref([
  {
    id: 'EVDA-2026-001234',
    title: 'Drive Master Sessions Course',
    recipientName: 'John Doe',
    issuedDate: 'Apr 15, 2026',
    status: 'issued',
    downloadUrl: '#'
  }
])

// For demo purposes, set certificate as available
const hasCertificate = ref(true)
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
        <!-- Certificate Status -->
        <UCard v-if="!hasCertificate">
          <template #header>
            <h2 class="font-semibold">{{ t('certificate.progress') }}</h2>
          </template>

          <div class="text-center py-8">
            <div class="mx-auto w-20 h-20 rounded-full bg-muted/50 flex items-center justify-center mb-4">
              <UIcon name="i-lucide-award" class="size-10 text-muted" />
            </div>
            <h3 class="text-xl font-bold mb-2">{{ t('certificate.notAvailable') }}</h3>
            <p class="text-muted max-w-md mx-auto">
              {{ t('certificate.notAvailableDesc') }}
            </p>

            <div class="max-w-md mx-auto mt-8 space-y-4">
              <div class="flex justify-between text-md">
                <span class="text-muted">{{ t('common.progress') }}</span>
                <span class="font-medium">{{ certificateStatus.completedSessions }}/{{ certificateStatus.totalSessions }} {{ t('billing.sessions') }}</span>
              </div>
              <UProgress :value="certificateStatus.progress" />
              <p class="text-md text-muted">
                {{ t('certificate.remainingSessions', { count: certificateStatus.remainingSessions }) }}
              </p>
            </div>
          </div>

          <template #footer>
            <div class="flex justify-center">
              <NuxtLink to="/dashboard/schedule">
                <UButton :label="t('dashboard.bookNextSession')" icon="i-lucide-calendar-plus" />
              </NuxtLink>
            </div>
          </template>
        </UCard>

        <!-- Available Certificates -->
        <template v-if="hasCertificate">
          <UAlert icon="i-lucide-award" color="success" :title="t('certificate.congratulations')">
            <template #description>
              {{ t('certificate.readyDownload') }}
            </template>
          </UAlert>

          <div class="grid md:grid-cols-2 gap-6">
            <template v-for="certificate in issuedCertificates" :key="certificate.id">
              <CertificateDisplay :certificate="certificate" />
            </template>
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
