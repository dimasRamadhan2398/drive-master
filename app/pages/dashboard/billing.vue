<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { paymentService } from '~/services/paymentService'

const { t } = useI18n()
definePageMeta({ layout: 'dashboard' })

const authStore = useAuthStore()

const paymentsList = ref<any[]>([])
const isLoadingPayments = ref(false)
const showBillingModal = ref(false)
const selectedPaymentMethod = ref('va')

const fetchPayments = async (userId: string) => {
  try {
    isLoadingPayments.value = true
    const res = await paymentService.listUserPayments(userId)
    paymentsList.value = res.payments || []
  } catch (err) {
    console.error('Failed to fetch payments:', err)
  } finally {
    isLoadingPayments.value = false
  }
}

onMounted(async () => {
  await authStore.fetchMemberProfile()
  if (authStore.userId) {
    fetchPayments(authStore.userId)
  }
})

const entitlements = computed(() => authStore.memberEntitlements ?? [])

// Find active or completed entitlements to show as package
const activeEntitlement = computed(() => {
  return entitlements.value.find(e => e.status === 'active') || entitlements.value[0] || null
})

// Dynamic user data
const userData = computed(() => {
  const active = activeEntitlement.value
  const total = entitlements.value.reduce((sum, e) => sum + (e.totalSessions ?? 0), 0)
  const completed = entitlements.value.reduce((sum, e) => sum + (e.usedSessions ?? 0), 0)
  const remaining = entitlements.value.reduce((sum, e) => sum + (e.remaining ?? 0), 0)
  
  return {
    name: authStore.currentUser?.firstName + ' ' + (authStore.currentUser?.lastName || ''),
    email: authStore.currentUser?.email || '',
    package: active?.packageName || t('dashboard.noPackage'),
    totalSessions: total,
    completedSessions: completed,
    remainingSessions: remaining,
    progress: total ? Math.round((completed / total) * 100) : 0,
  }
})

// Find any pending payment
const pendingPayment = computed(() => {
  return paymentsList.value.find(p => p.status === 'pending') || null
})

// Dynamic billing data
const billingData = computed(() => {
  const pending = pendingPayment.value
  const active = activeEntitlement.value
  
  if (pending) {
    return {
      status: 'payment-due',
      amount: formatCurrency(pending.amount),
      dueDate: pending.expiresAt ? formatDate(pending.expiresAt) : t('common.pending'),
      package: getPaymentPlanName(pending),
      paymentUrl: pending.paymentUrl || '',
      nextRenewalDate: pending.expiresAt ? formatDate(pending.expiresAt) : t('common.pending')
    }
  }
  
  // Default paid/active status
  return {
    status: active ? 'paid' : 'inactive',
    amount: formatCurrency(2250000), // Default mock fallback amount
    dueDate: active?.endDate ? formatDate(active.endDate) : t('common.pending'),
    package: active?.packageName || t('dashboard.noPackage'),
    paymentUrl: '',
    nextRenewalDate: active?.endDate ? formatDate(active.endDate) : t('common.pending')
  }
})

const paymentMethods = computed(() => [
  {
    id: 'va',
    name: 'Virtual Account (VA)',
    description: 'Transfer via BCA, Mandiri, BNI, or BRI virtual account',
    icon: 'i-lucide-building',
    color: 'blue' as const
  },
  {
    id: 'qris',
    name: 'QRIS',
    description: 'Scan QR code with any e-wallet',
    icon: 'i-lucide-qr-code',
    color: 'green' as const
  },
  {
    id: 'bank_transfer',
    name: 'Bank Transfer',
    description: 'Direct transfer from your bank account',
    icon: 'i-lucide-landmark',
    color: 'purple' as const
  },
  {
    id: 'ewallet',
    name: 'E-Wallet',
    description: 'GoPay, OVO, DANA, or LinkAja',
    icon: 'i-lucide-wallet',
    color: 'orange' as const
  }
])

const pricingTabs = computed(() => [
  { label: '6x Sessions', icon: 'i-lucide-package' },
  { label: '8x Sessions', icon: 'i-lucide-package' },
  { label: '10x Sessions', icon: 'i-lucide-package', default: true },
  { label: '12x Sessions', icon: 'i-lucide-package' }
])


function getBillingStatusColor() {
  return billingData.value.status === 'payment-due' ? 'warning' : 'success'
}


function getBillingStatusLabel() {
  if (billingData.value.status === 'payment-due') return t('billing.paymentDue')
  if (billingData.value.status === 'paid') return t('billing.paid')
  return t('billing.active')
}

const planMap: { [key: string]: string } = {
  '6x Sessions': 'six_package',
  '8x Sessions': 'eight_package',
  '10x Sessions': 'ten_package',
  '12x Sessions': 'twelve_package',
  '6x Session Package': 'six_package',
  '8x Session Package': 'eight_package',
  '10x Session Package': 'ten_package',
  '12x Session Package': 'twelve_package',
  'Bronze Package': 'six_package',
  'Silver Package': 'eight_package',
  'Gold Package': 'ten_package',
  'Platinum Package': 'twelve_package'
}

function makePayment() {
  if (billingData.value.status === 'payment-due' && billingData.value.paymentUrl) {
    // Open payment URL directly
    window.location.href = billingData.value.paymentUrl
    return
  }

  const planKey = planMap[billingData.value.package] || 'ten_package'

  navigateTo({
    path: '/auth/payment',
    query: {
      method: selectedPaymentMethod.value,
      plan: planKey,
      email: userData.value.email
    }
  })
}

// Helper methods
function formatCurrency(value: number) {
  if (value === undefined || value === null) return 'Rp 0'
  return `Rp ${value.toLocaleString('id-ID')}`
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function getPaymentPlanName(payment: any) {
  const entitlement = entitlements.value.find(e => e.bookingId === payment.enrollmentId)
  if (entitlement) return entitlement.packageName
  return t('billing.package') || 'Driving Package'
}

function getPaymentStatusColor(status: string) {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    failed: 'error',
    cancelled: 'neutral',
    expired: 'neutral'
  }
  return map[status] || 'primary'
}
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="t('billing.title')">
                <template #right>
                    <UButton icon="i-lucide-bell" color="neutral" variant="ghost" />
                    <UColorModeButton />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
    <div class="space-y-6">
       <!-- Billing Notification Alert -->
       <div v-if="billingData.status === 'payment-due'" class="p-4 rounded-lg border-l-4 border-warning bg-warning/5">
         <div class="flex items-start gap-4">
           <div class="p-2 rounded-lg bg-warning/10 shrink-0">
             <UIcon name="i-lucide-alert-circle" class="size-6 text-warning" />
           </div>
           <div class="flex-1">
             <h3 class="font-semibold text-foreground">{{ t('billing.paymentDue') }}</h3>
             <p class="text-sm text-muted mt-1">
               {{ t('billing.paymentDueDesc', { amount: billingData.amount, date: billingData.dueDate }) }}
             </p>
             <p class="text-xs text-muted mt-2">{{ billingData.package }}</p>
             <div class="flex flex-wrap gap-3 mt-4">
               <UButton :label="t('billing.proceedToPayment')" size="sm" icon="i-lucide-credit-card" color="warning" @click="showBillingModal = true"/>
             </div>
           </div>
           <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs" />
         </div>
       </div>


        <!-- Billing Overview Card (Alternative: when payment is not due) -->
        <div v-else class="p-4 rounded-lg border border-default bg-elevated hover:bg-muted/30 transition-colors">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="p-3 rounded-lg bg-success/10">
                <UIcon :name="userData.remainingSessions === 0 && userData.totalSessions > 0 ? 'i-lucide-award' : 'i-lucide-check-circle'" class="size-6 text-success" />
              </div>
              <div>
                <h3 class="font-semibold">
                  {{ userData.remainingSessions === 0 && userData.totalSessions > 0 ? t('billing.packageCompleted') : t('billing.subscriptionActive') }}
                </h3>
                <p class="text-sm text-muted">
                  {{ userData.remainingSessions === 0 && userData.totalSessions > 0 ? t('billing.allSessionsUsed') : t('billing.nextRenewal', { date: billingData.nextRenewalDate }) }}
                </p>
              </div>
            </div>
            <NuxtLink v-if="userData.remainingSessions === 0 && userData.totalSessions > 0" to="/auth/select-plan">
              <UButton :label="t('billing.buyNewPackage')" color="warning" size="sm" />
            </NuxtLink>
            <UButton v-else :label="t('billing.manageBilling')" variant="outline" color="neutral" size="sm" @click="showBillingModal = true" />
          </div>
        </div>

       <!-- Welcome Banner -->
       <UCard class="bg-primary/5 border-primary/20">
         <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
           <div>
             <h2 class="text-xl font-semibold">
               {{ userData.remainingSessions === 0 && userData.totalSessions > 0 ? t('billing.packageCompleted') : t('billing.sessionsRemaining', { count: userData.remainingSessions }) }}
             </h2>
             <p class="text-muted mt-1">{{ t('billing.inPackage', { package: userData.package }) }}</p>
           </div>
           <NuxtLink v-if="userData.remainingSessions === 0 && userData.totalSessions > 0" to="/auth/select-plan">
             <UButton :label="t('billing.buyNewPackage')" color="warning" icon="i-lucide-shopping-cart" />
           </NuxtLink>
           <NuxtLink v-else to="/dashboard/schedule">
            <UButton :label="t('dashboard.bookNextSession')" color="warning" icon="i-lucide-calendar-plus" />
           </NuxtLink>
         </div>
       </UCard>

       <!-- Billing Plan Section -->
       <UCard>
         <template #header>
           <h2 class="font-semibold">{{ t('billing.packagePricing') }}</h2>
         </template>
         
         <!-- Package Tabs -->
         <div class="mb-4">
           <UHorizontalNavigation :links="pricingTabs" />
         </div>
         
         <!-- Package Details -->
         <div class="p-4 border border-default rounded-lg bg-muted/30">
           <div class="flex items-center justify-between mb-4">
             <div class="flex items-center gap-3">
               <UIcon :name="userData.remainingSessions === 0 && userData.totalSessions > 0 ? 'i-lucide-award' : 'i-lucide-check-circle'" class="size-6 text-success" />
               <h3 class="font-semibold">
                 {{ billingData.package }} - 
                 {{ userData.remainingSessions === 0 && userData.totalSessions > 0 ? t('billing.packageCompleted') : t('billing.active') }}
               </h3>
             </div>
             <div class="text-right">
               <p class="text-2xl font-bold text-primary">{{ billingData.amount }}</p>
               <p class="text-sm text-muted">{{ t('billing.nextBilling', { date: billingData.nextRenewalDate }) }}</p>
             </div>
           </div>
           
           <div class="grid md:grid-cols-2 gap-4 mb-4">
             <div class="text-sm">
               <div class="text-muted">{{ t('billing.sessions') }}</div>
                <div class="font-medium">{{ userData.totalSessions || 10 }} {{ t('billing.sessions') }}</div>
             </div>
             <div class="text-sm">
               <div class="text-muted">{{ t('billing.totalCost') }}</div>
               <div class="font-medium">{{ billingData.amount }}</div>
             </div>
           </div>
           <NuxtLink :to="`/auth/select-plan?current_plan=${planMap[billingData.package]}`">
            <UButton :label="t('billing.changePackage')" color="warning" icon="i-lucide-credit-card" />
           </NuxtLink>
         </div>
       </UCard>

        <!-- Payment History -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">{{ t('billing.paymentHistory') }}</h2>
          </template>
          
          <div v-if="isLoadingPayments" class="space-y-4">
            <USkeleton class="h-10 w-full" />
            <USkeleton class="h-10 w-full" />
          </div>
          <div v-else-if="paymentsList.length > 0" class="divide-y divide-default">
            <div 
              v-for="payment in paymentsList" 
              :key="payment.id" 
              class="flex items-center justify-between py-3 first:pt-0 last:pb-0"
            >
              <div>
                <p class="font-medium">{{ getPaymentPlanName(payment) }}</p>
                <p class="text-xs text-muted">{{ formatDate(payment.createdAt) }}</p>
                <p class="text-xs text-muted font-mono mt-0.5">{{ payment.orderId }}</p>
              </div>
              <div class="text-right">
                <p class="font-semibold">{{ formatCurrency(payment.amount) }}</p>
                <UBadge 
                  :label="t(`payment.status.${payment.status}`) || payment.status" 
                  :color="getPaymentStatusColor(payment.status)" 
                  variant="subtle" 
                  size="xs" 
                />
              </div>
            </div>
          </div>
          <div v-else class="text-center py-6 text-sm text-muted">
            No payment history found.
          </div>
        </UCard> 
    </div>

   <!-- Billing Modal -->
   <UModal v-model:open="showBillingModal" :title="t('billing.title')">
     <template #body>
       <div class="space-y-6">
         <!-- Current Due -->
         <div class="border border-default rounded-lg p-4 bg-muted/30">
           <div class="flex items-center justify-between mb-4">
             <h3 class="font-semibold flex items-center gap-2">
               <UIcon name="i-lucide-receipt" class="size-5" />
               {{ t('billing.paymentDue') }}
             </h3>
             <UBadge :label="getBillingStatusLabel()" :color="getBillingStatusColor()" variant="subtle" />
           </div>
           <div class="space-y-3">
             <div class="flex items-center justify-between">
               <span class="text-muted">{{ t('billing.amountDue') }}</span>
               <span class="font-semibold text-lg">{{ billingData.amount }}</span>
             </div>
             <div class="flex items-center justify-between">
               <span class="text-muted">{{ t('billing.dueDate') }}</span>
               <span class="font-medium">{{ billingData.dueDate }}</span>
             </div>
             <div class="flex items-center justify-between">
               <span class="text-muted">{{ t('billing.package') }}</span>
               <span class="font-medium">{{ billingData.package }}</span>
             </div>
             <div class="flex items-center justify-between">
               <span class="text-muted">{{ t('billing.nextRenewal', { date: '' }).replace(' on ', '') }}</span>
               <span class="font-medium">{{ billingData.nextRenewalDate }}</span>
             </div>
           </div>
         </div>

         <!-- Payment Methods -->
         <div>
           <h3 class="font-semibold mb-3 flex items-center gap-2">
             <UIcon name="i-lucide-credit-card" class="size-5" />
             {{ t('billing.paymentMethods') }}
           </h3>
           <div class="space-y-2">
             <button 
               v-for="method in paymentMethods" 
               :key="method.id"
               @click="selectedPaymentMethod = method.id"
               class="w-full p-3 rounded-lg border-2 transition-colors text-left"
               :class="selectedPaymentMethod === method.id ? 'border-primary bg-primary/5' : 'border-default hover:border-primary/50 bg-background'"
             >
               <div class="flex items-center gap-3">
                 <div class="p-2 rounded" :class="`bg-${method.color}-500/10`">
                   <UIcon :name="method.icon" class="size-5" :class="`text-${method.color}-500`" />
                 </div>
                 <div>
                   <p class="font-medium">{{ method.name }}</p>
                   <p class="text-xs text-muted">{{ method.description }}</p>
                 </div>
                 <UIcon v-if="selectedPaymentMethod === method.id" name="i-lucide-check" class="size-5 text-primary ml-auto" />
               </div>
             </button>
           </div>

           <!-- Payment Details (Direct Payment) -->
           <div v-if="selectedPaymentMethod" class="mt-4 p-4 rounded-lg bg-muted/30 border border-default">
             <div v-if="selectedPaymentMethod === 'va'" class="space-y-3">
               <p class="text-sm font-medium">{{ t('billing.vaNumber') }}</p>
               <div class="flex items-center justify-between p-3 bg-background rounded border border-default">
                 <span class="font-mono font-bold text-lg tracking-wider">8801 2345 6789</span>
                 <UButton icon="i-lucide-copy" size="sm" color="neutral" variant="ghost" />
               </div>
               <p class="text-xs text-muted">{{ t('billing.vaInstruction') }}</p>
             </div>
             <div v-else-if="selectedPaymentMethod === 'qris'" class="space-y-3 flex flex-col items-center">
               <p class="text-sm font-medium w-full text-left">{{ t('billing.scanQris') }}</p>
               <div class="p-2 bg-white rounded-lg inline-block">
                 <UIcon name="i-lucide-qr-code" class="size-32 text-black" />
               </div>
               <p class="text-xs text-muted text-center">{{ t('billing.qrisInstruction') }}</p>
             </div>
             <div v-else class="space-y-3">
               <p class="text-sm font-medium">{{ t('billing.paymentInstruction') }}</p>
               <p class="text-sm text-muted">{{ t('billing.paymentInstructionDesc', { method: paymentMethods.find(m => m.id === selectedPaymentMethod)?.name }) }}</p>
             </div>
           </div>
         </div>

         <!-- Billing Info -->
         <UAlert icon="i-lucide-info" :title="t('billing.billingInfo')" variant="subtle">
           <template #description>
             {{ t('billing.billingInfoDesc') }}
           </template>
         </UAlert>
        
       </div>
     </template>
     <template #footer>
       <div class="flex justify-end gap-3">
         <UButton :label="t('billing.makePayment')" color="warning" icon="i-lucide-credit-card" @click="makePayment" />
       </div>
     </template>
   </UModal>
   </template>
</UDashboardPanel>
</template>
