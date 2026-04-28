<script setup lang="ts">
import { ref } from 'vue'


definePageMeta({ layout: 'dashboard' })

// Billing data
const userData = {
  name: 'John Doe',
  package: '10x Sessions Package',
  totalSessions: 10,
  completedSessions: 4,
  remainingSessions: 6,
  progress: 40,
  nextSession: {
    date: 'Tomorrow',
    time: '09:30 AM',
    car: 'BYD Atto 1',
    instructor: 'Pak Ahmad'
  }
}

const billingData = {
 status: 'payment-due',
 daysRemaining: 3,
 amount: 'Rp 2,250,000',
 dueDate: 'Apr 24, 2026',
 lastPaymentDate: 'Mar 25, 2026',
 package: '10x Session Package',
 nextRenewalDate: 'May 25, 2026'
}

const showBillingModal = ref(false)
const selectedPaymentMethod = ref('va')

const paymentMethods = [
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
]

const pricingTabs = [
  { label: '6x Sessions', icon: 'i-lucide-package' },
  { label: '8x Sessions', icon: 'i-lucide-package' },
  { label: '10x Sessions', icon: 'i-lucide-package', default: true },
  { label: '12x Sessions', icon: 'i-lucide-package' }
]


function getBillingStatusColor() {
 return billingData.status === 'payment-due' ? 'warning' : 'success'
}


function getBillingStatusLabel() {
 if (billingData.status === 'payment-due') return 'Payment Due'
 if (billingData.status === 'paid') return 'Paid'
 return 'Active'
}
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar title="Billing & Payments">
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
             <h3 class="font-semibold text-foreground">Payment Due</h3>
             <p class="text-sm text-muted mt-1">
               Your payment of <span class="font-medium">{{ billingData.amount }}</span> is due on <span class="font-medium">{{ billingData.dueDate }}</span> ({{ billingData.daysRemaining }} days remaining).
             </p>
             <p class="text-xs text-muted mt-2">{{ billingData.package }}</p>
             <div class="flex flex-wrap gap-3 mt-4">
               <UButton label="Proceed to Payment" size="sm" icon="i-lucide-credit-card" color="warning" @click="showBillingModal = true"/>
               <!-- <UButton label="View Invoice" size="sm" variant="outline" color="neutral" /> -->
               <!-- <UButton label="Download Receipt" size="sm" variant="ghost" icon="i-lucide-download" color="neutral" /> -->
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
               <UIcon name="i-lucide-check-circle" class="size-6 text-success" />
             </div>
             <div>
               <h3 class="font-semibold">Your subscription is active</h3>
               <p class="text-sm text-muted">Next renewal on {{ billingData.nextRenewalDate }}</p>
             </div>
           </div>
           <UButton label="Manage Billing" variant="outline" color="neutral" size="sm" @click="showBillingModal = true" />
         </div>
       </div>

       <!-- Welcome Banner -->
       <UCard class="bg-primary/5 border-primary/20">
         <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
           <div>
             <h2 class="text-xl font-semibold">You have {{ userData.remainingSessions }} sessions remaining</h2>
             <p class="text-muted mt-1">in your {{ userData.package }}.</p>
           </div>
           <NuxtLink to="/dashboard/schedule">
            <UButton label="Book Next Session" color="warning" icon="i-lucide-calendar-plus" />
           </NuxtLink>
         </div>
       </UCard>

       <!-- Billing Plan Section -->
       <UCard>
         <template #header>
           <h2 class="font-semibold">Your Package & Pricing</h2>
         </template>
         
         <!-- Package Tabs -->
         <div class="mb-4">
           <UHorizontalNavigation :links="pricingTabs" />
         </div>
         
         <!-- Package Details -->
         <div class="p-4 border border-default rounded-lg bg-muted/30">
           <div class="flex items-center justify-between mb-4">
             <div class="flex items-center gap-3">
               <UIcon name="i-lucide-check-circle" class="size-6 text-success" />
               <h3 class="font-semibold">{{ billingData.package }} - Active</h3>
             </div>
             <div class="text-right">
               <p class="text-2xl font-bold text-primary">{{ billingData.amount }}</p>
               <p class="text-sm text-muted">Next billing: {{ billingData.nextRenewalDate }}</p>
             </div>
           </div>
           
           <div class="grid md:grid-cols-2 gap-4 mb-4">
             <div class="text-sm">
               <div class="text-muted">Sessions</div>
               <div class="font-medium">10 Sessions</div>
             </div>
             <div class="text-sm">
               <div class="text-muted">Total Cost</div>
               <div class="font-medium">{{ billingData.amount }}</div>
             </div>
           </div>
           <NuxtLink to="/select-plan">
            <UButton label="Upgrade Package" color="warning" icon="i-lucide-credit-card" />
           </NuxtLink>
         </div>
       </UCard>

       <!-- Last Payment -->
       <UCard>
         <template #header>
           <h2 class="font-semibold">Payment History</h2>
         </template>
           
           <div class="space-y-3">
             <div class="flex items-center justify-between pb-3">
               <div>
                 <p class="font-medium">{{ billingData.package }}</p>
                 <p class="text-xs text-muted">{{ billingData.lastPaymentDate }}</p>
               </div>
               <div class="text-right">
                 <p class="font-semibold">{{ billingData.amount }}</p>
                 <UBadge label="Paid" color="success" variant="subtle" size="xs" />
               </div>
             </div>
           
         </div>
       </UCard> 
    </div>

   <!-- Billing Modal -->
   <UModal v-model:open="showBillingModal" title="Billing & Payment">
     <template #body>
       <div class="space-y-6">
         <!-- Current Due -->
         <div class="border border-default rounded-lg p-4 bg-muted/30">
           <div class="flex items-center justify-between mb-4">
             <h3 class="font-semibold flex items-center gap-2">
               <UIcon name="i-lucide-receipt" class="size-5" />
               Current Bill
             </h3>
             <UBadge :label="getBillingStatusLabel()" :color="getBillingStatusColor()" variant="subtle" />
           </div>
           <div class="space-y-3">
             <div class="flex items-center justify-between">
               <span class="text-muted">Amount Due</span>
               <span class="font-semibold text-lg">{{ billingData.amount }}</span>
             </div>
             <div class="flex items-center justify-between">
               <span class="text-muted">Due Date</span>
               <span class="font-medium">{{ billingData.dueDate }}</span>
             </div>
             <div class="flex items-center justify-between">
               <span class="text-muted">Package</span>
               <span class="font-medium">{{ billingData.package }}</span>
             </div>
             <div class="flex items-center justify-between">
               <span class="text-muted">Next Renewal</span>
               <span class="font-medium">{{ billingData.nextRenewalDate }}</span>
             </div>
           </div>
         </div>

         <!-- Payment Methods -->
         <div>
           <h3 class="font-semibold mb-3 flex items-center gap-2">
             <UIcon name="i-lucide-credit-card" class="size-5" />
             Payment Methods
           </h3>
           <div class="space-y-2">
             <button 
               v-for="method in paymentMethods" 
               :key="method.id"
               @click="selectedPaymentMethod = method.id"
               class="w-full p-3 rounded-lg border-2 transition-colors text-left"
               :class="selectedPaymentMethod === method.id ? 'border-warning bg-warning/5' : 'border-default hover:border-warning/50 bg-background'"
             >
               <div class="flex items-center gap-3">
                 <div class="p-2 rounded" :class="`bg-${method.color}-500/10`">
                   <UIcon :name="method.icon" class="size-5" :class="`text-${method.color}-500`" />
                 </div>
                 <div>
                   <p class="font-medium">{{ method.name }}</p>
                   <p class="text-xs text-muted">{{ method.description }}</p>
                 </div>
                 <UIcon v-if="selectedPaymentMethod === method.id" name="i-lucide-check" class="size-5 text-warning ml-auto" />
               </div>
             </button>
           </div>

           <!-- Payment Details (Direct Payment) -->
           <div v-if="selectedPaymentMethod" class="mt-4 p-4 rounded-lg bg-muted/30 border border-default">
             <div v-if="selectedPaymentMethod === 'va'" class="space-y-3">
               <p class="text-sm font-medium">Virtual Account Number</p>
               <div class="flex items-center justify-between p-3 bg-background rounded border border-default">
                 <span class="font-mono font-bold text-lg tracking-wider">8801 2345 6789</span>
                 <UButton icon="i-lucide-copy" size="sm" color="neutral" variant="ghost" />
               </div>
               <p class="text-xs text-muted">Please transfer the exact amount before the due date.</p>
             </div>
             <div v-else-if="selectedPaymentMethod === 'qris'" class="space-y-3 flex flex-col items-center">
               <p class="text-sm font-medium w-full text-left">Scan QR Code</p>
               <div class="p-2 bg-white rounded-lg inline-block">
                 <UIcon name="i-lucide-qr-code" class="size-32 text-black" />
               </div>
               <p class="text-xs text-muted text-center">Open your e-wallet or banking app and scan the QR code above.</p>
             </div>
             <div v-else class="space-y-3">
               <p class="text-sm font-medium">Payment Instructions</p>
               <p class="text-sm text-muted">Please follow the instructions sent to your email to complete the payment via {{ paymentMethods.find(m => m.id === selectedPaymentMethod)?.name }}.</p>
             </div>
           </div>
         </div>

         <!-- Billing Info -->
         <UAlert icon="i-lucide-info" title="Billing Information" variant="subtle">
           <template #description>
             Invoices are sent to your registered email. You can download, view, or manage your subscription anytime.
           </template>
         </UAlert>
        
       </div>
     </template>
     <template #footer>
       <div class="flex justify-end gap-3">
         <UButton label="Make Payment" color="warning" icon="i-lucide-credit-card" @click="showBillingModal = false" />
       </div>
     </template>
   </UModal>
   </template>
</UDashboardPanel>
</template>
