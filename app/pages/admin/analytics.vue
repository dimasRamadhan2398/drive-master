<template>
    <div class="p-6 space-y-6 w-full">
    <!-- Google Analytics (GA4) Real-Time Insights -->
        <div class="flex items-center justify-between">
        <div>
            <h2 class="text-lg font-bold">{{ t('admin.gaInsights') }}</h2>
            <p class="text-sm text-muted font-normal">{{ t('admin.gaInsightsDesc') }}</p>
        </div>
        <UButton icon="i-lucide-refresh-cw" color="neutral" variant="ghost" size="xs" :loading="gaLoading" @click="fetchGAData" />
        </div>

        <div class="grid md:grid-cols-3 gap-6">
        <UCard class="relative overflow-hidden">
            <div class="flex items-center gap-4">
            <div class="p-3 rounded-xl bg-primary-500/10">
                <UIcon name="i-lucide-users" class="size-6 text-primary-500" />
            </div>
            <div>
                <p class="text-2xl font-bold">
                <span v-if="gaLoading">...</span>
                <span v-else>{{ analyticsStore.totalGaUsers }}</span>
                </p>
                <p class="text-sm text-muted">{{ t('admin.totalVisitors') }}</p>
            </div>
            </div>
        </UCard>

        <UCard class="relative overflow-hidden">
            <div class="flex items-center gap-4">
            <div class="p-3 rounded-xl bg-success-500/10">
                <UIcon name="i-lucide-eye" class="size-6 text-success-500" />
            </div>
            <div>
                <p class="text-2xl font-bold">
                <span v-if="gaLoading">...</span>
                <span v-else>{{ analyticsStore.totalGaPageViews }}</span>
                </p>
                <p class="text-sm text-muted">{{ t('admin.pageViews') }}</p>
            </div>
            </div>
        </UCard>

        <UCard class="relative overflow-hidden">
            <div class="flex items-center gap-4">
            <div class="p-3 rounded-xl bg-amber-500/10">
                <UIcon name="i-lucide-trending-up" class="size-6 text-amber-500" />
            </div>
            <div>
                <p class="text-2xl font-bold">
                <span v-if="gaLoading">...</span>
                <span v-else>{{ analyticsStore.overallConversionRate }}%</span>
                </p>
                <p class="text-sm text-muted">{{ t('admin.conversionRate') }}</p>
            </div>
            </div>
        </UCard>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
        <!-- Visitor Graph -->
        <UCard>
            <template #header>
            <div class="flex items-center justify-between">
                <h3 class="font-semibold">{{ t('admin.visitorTrend') }}</h3>
                <div class="flex items-center gap-4 text-[10px] text-muted">
                <span class="flex items-center gap-1"><span class="size-2 rounded-full bg-orange-500"></span> {{ t('admin.pageViews') }}</span>
                <span class="flex items-center gap-1"><span class="size-2 rounded-full bg-blue-500"></span> {{ t('admin.activeUsers') }}</span>
                </div>
            </div>
            </template>

            <div v-if="gaLoading" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">{{ t('admin.loadingChart') }}</p>
            </div>
            <div v-else-if="analyticsStore.gaOverviewData.length === 0" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">{{ t('admin.noVisitorData') }}</p>
            </div>
            <div v-else class="relative pt-2">
            <svg viewBox="0 0 500 150" class="w-full h-[150px] overflow-visible">
                <defs>
                <linearGradient id="chartGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stop-color="rgb(249, 115, 22)" stop-opacity="0.15" />
                    <stop offset="100%" stop-color="rgb(249, 115, 22)" stop-opacity="0.0" />
                </linearGradient>
                </defs>

                <!-- Grid Lines -->
                <line x1="0" y1="37.5" x2="500" y2="37.5" stroke="currentColor" class="text-default/5" stroke-dasharray="4" />
                <line x1="0" y1="75" x2="500" y2="75" stroke="currentColor" class="text-default/5" stroke-dasharray="4" />
                <line x1="0" y1="112.5" x2="500" y2="112.5" stroke="currentColor" class="text-default/5" stroke-dasharray="4" />
                <line x1="0" y1="150" x2="500" y2="150" stroke="currentColor" class="text-default/5" />

                <!-- Pageviews Area & Path -->
                <path :d="analyticsStore.chartAreaPath" fill="url(#chartGrad)" />
                <path :d="analyticsStore.chartPathData" fill="none" stroke="rgb(249, 115, 22)" stroke-width="2.5" />

                <!-- Users Path -->
                <path :d="analyticsStore.chartUsersPath" fill="none" stroke="rgb(59, 130, 246)" stroke-width="2" />
            </svg>
            <div class="flex justify-between text-[9px] text-muted mt-2 px-1">
                <span>{{ analyticsStore.gaOverviewData[0]?.date }}</span>
                <span>{{ analyticsStore.gaOverviewData[Math.floor(analyticsStore.gaOverviewData.length / 2)]?.date }}</span>
                <span>{{ analyticsStore.gaOverviewData[analyticsStore.gaOverviewData.length - 1]?.date }}</span>
            </div>
            </div>
        </UCard>

        <!-- Conversion Funnel -->
        <UCard>
            <template #header>
            <h3 class="font-semibold">{{ t('admin.funnelTitle') }}</h3>
            </template>

            <div v-if="gaLoading" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">{{ t('admin.loadingFunnel') }}</p>
            </div>
            <div v-else-if="analyticsStore.gaFunnelData.length === 0" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">{{ t('admin.noFunnelData') }}</p>
            </div>
            <div v-else class="space-y-3.5 py-1">
            <div v-for="(step, idx) in analyticsStore.gaFunnelData" :key="idx" class="space-y-1">
                <div class="flex items-center justify-between text-[11px] font-medium">
                <span class="capitalize text-muted">{{ getFunnelLabel(step.event_name) }}</span>
                <span>{{ step.count }} events <span class="text-muted text-[9px] font-normal" v-if="idx > 0 && analyticsStore.gaFunnelData[0]?.count">({{ ((step.count / analyticsStore.gaFunnelData[0].count) * 100).toFixed(1) }}%)</span></span>
                </div>
                <div class="w-full bg-default/10 rounded-full h-2 overflow-hidden">
                <div
                    class="h-full rounded-full transition-all duration-500"
                    :class="getFunnelColor(idx)"
                    :style="{ width: `${(step.count / (analyticsStore.gaFunnelData[0]?.count || 1)) * 100}%` }"
                ></div>
                </div>
            </div>
            </div>
        </UCard>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useAnalyticsStore } from '~/stores/analytics'

const { t } = useI18n()
definePageMeta({
    layout: "admin",
});

const analyticsStore = useAnalyticsStore()
const gaLoading = computed(() => analyticsStore.isLoading)

function getFunnelLabel(eventName: string) {
  return analyticsStore.getFunnelLabel(eventName)
}

function getFunnelColor(index: number) {
  return analyticsStore.getFunnelColor(index)
}

async function fetchGAData() {
  await analyticsStore.fetchAnalyticsData()
}

onMounted(() => {
  fetchGAData()
})
</script>