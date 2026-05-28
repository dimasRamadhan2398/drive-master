<template>
    <div class="p-6 space-y-6 w-full">
    <!-- Google Analytics (GA4) Real-Time Insights -->
        <div class="flex items-center justify-between">
        <div>
            <h2 class="text-lg font-bold">Google Analytics (GA4) Real-Time Insights</h2>
            <p class="text-sm text-muted font-normal">Pelacakan traffic situs web dan funnel konversi secara real-time</p>
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
                <span v-else>{{ totalGaUsers }}</span>
                </p>
                <p class="text-sm text-muted">Total Visitors (30 Days)</p>
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
                <span v-else>{{ totalGaPageViews }}</span>
                </p>
                <p class="text-sm text-muted">Page Views (30 Days)</p>
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
                <span v-else>{{ overallConversionRate }}%</span>
                </p>
                <p class="text-sm text-muted">Checkout Conversion Rate</p>
            </div>
            </div>
        </UCard>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
        <!-- Visitor Graph -->
        <UCard>
            <template #header>
            <div class="flex items-center justify-between">
                <h3 class="font-semibold">Visitor Trend</h3>
                <div class="flex items-center gap-4 text-[10px] text-muted">
                <span class="flex items-center gap-1"><span class="size-2 rounded-full bg-orange-500"></span> Page Views</span>
                <span class="flex items-center gap-1"><span class="size-2 rounded-full bg-blue-500"></span> Active Users</span>
                </div>
            </div>
            </template>
            
            <div v-if="gaLoading" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">Loading chart data...</p>
            </div>
            <div v-else-if="gaOverview.length === 0" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">No visitor data available.</p>
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
                <path :d="chartAreaPath" fill="url(#chartGrad)" />
                <path :d="chartPath" fill="none" stroke="rgb(249, 115, 22)" stroke-width="2.5" />
                
                <!-- Users Path -->
                <path :d="chartUsersPath" fill="none" stroke="rgb(59, 130, 246)" stroke-width="2" />
            </svg>
            <div class="flex justify-between text-[9px] text-muted mt-2 px-1">
                <span>{{ gaOverview[0]?.date }}</span>
                <span>{{ gaOverview[Math.floor(gaOverview.length / 2)]?.date }}</span>
                <span>{{ gaOverview[gaOverview.length - 1]?.date }}</span>
            </div>
            </div>
        </UCard>

        <!-- Conversion Funnel -->
        <UCard>
            <template #header>
            <h3 class="font-semibold">Booking Conversion Funnel</h3>
            </template>
            
            <div v-if="gaLoading" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">Loading funnel data...</p>
            </div>
            <div v-else-if="gaFunnel.length === 0" class="h-[150px] flex items-center justify-center">
            <p class="text-sm text-muted">No funnel data available.</p>
            </div>
            <div v-else class="space-y-3.5 py-1">
            <div v-for="(step, idx) in gaFunnel" :key="idx" class="space-y-1">
                <div class="flex items-center justify-between text-[11px] font-medium">
                <span class="capitalize text-muted">{{ getFunnelLabel(step.event_name) }}</span>
                <span>{{ step.count }} events <span class="text-muted text-[9px] font-normal" v-if="idx > 0 && gaFunnel[0]?.count">({{ ((step.count / gaFunnel[0].count) * 100).toFixed(1) }}%)</span></span>
                </div>
                <div class="w-full bg-default/10 rounded-full h-2 overflow-hidden">
                <div 
                    class="h-full rounded-full transition-all duration-500"
                    :class="getFunnelColor(idx)"
                    :style="{ width: `${(step.count / (gaFunnel[0]?.count || 1)) * 100}%` }"
                ></div>
                </div>
            </div>
            </div>
        </UCard>
        </div>
    </div>
</template>

<script setup lang="ts">

definePageMeta({
    layout: "admin",
});

const api = useApiClients()
const gaLoading = ref(false)
const gaOverview = ref<{ date: string; users: number; pageviews: number }[]>([])
const gaFunnel = ref<{ event_name: string; count: number }[]>([])

// Extracted functions outside of computed
function getFunnelLabel(eventName: string) {
  switch (eventName) {
    case 'page_view': return '1. Landing Page Visit'
    case 'view_item': return '2. View Packages'
    case 'begin_checkout': return '3. Start Booking Process'
    case 'purchase': return '4. Successful Payment'
    default: return eventName
  }
}

function getFunnelColor(index: number) {
  switch (index) {
    case 0: return 'bg-primary-500'
    case 1: return 'bg-info-500'
    case 2: return 'bg-warning-500'
    case 3: return 'bg-success-500'
    default: return 'bg-neutral-500'
  }
}

const overallConversionRate = computed(() => {
  const start = gaFunnel.value[0]?.count
  const end = gaFunnel.value[3]?.count
  
  if (!start || !end) return '0.0'
  return ((end / start) * 100).toFixed(1)
})

const totalGaUsers = computed(() => {
  return gaOverview.value.reduce((sum, item) => sum + item.users, 0)
})

const totalGaPageViews = computed(() => {
  return gaOverview.value.reduce((sum, item) => sum + item.pageviews, 0)
})

const chartPath = computed(() => {
  if (gaOverview.value.length === 0) return ''
  const width = 500
  const height = 150
  const maxVal = Math.max(...gaOverview.value.map(d => d.pageviews), 100)
  
  return gaOverview.value.map((d, idx) => {
    const x = (idx / (gaOverview.value.length - 1)) * width
    const y = height - (d.pageviews / maxVal) * height
    return `${idx === 0 ? 'M' : 'L'} ${x.toFixed(1)} ${y.toFixed(1)}`
  }).join(' ')
})

const chartAreaPath = computed(() => {
  if (gaOverview.value.length === 0) return ''
  const width = 500
  const height = 150
  const base = chartPath.value
  return `${base} L ${width} ${height} L 0 ${height} Z`
})

const chartUsersPath = computed(() => {
  if (gaOverview.value.length === 0) return ''
  const width = 500
  const height = 150
  const maxVal = Math.max(...gaOverview.value.map(d => d.pageviews), 100)
  
  return gaOverview.value.map((d, idx) => {
    const x = (idx / (gaOverview.value.length - 1)) * width
    const y = height - (d.users / maxVal) * height
    return `${idx === 0 ? 'M' : 'L'} ${x.toFixed(1)} ${y.toFixed(1)}`
  }).join(' ')
})

async function fetchGAData() {
  gaLoading.value = true
  try {
    const overviewRes = await api.core<{ data: { date: string, users: number, pageviews: number }[] }>('/admin/analytics/overview')
    const funnelRes = await api.core<{ data: { event_name: string, count: number }[] }>('/admin/analytics/funnel')
    
    if (overviewRes && overviewRes.data) {
      gaOverview.value = overviewRes.data
    }
    if (funnelRes && funnelRes.data) {
      gaFunnel.value = funnelRes.data
    }
  } catch (err) {
    console.error("Gagal memuat data Google Analytics 4:", err)
  } finally {
    gaLoading.value = false
  }
}

onMounted(() => {
  fetchGAData()
})
</script>