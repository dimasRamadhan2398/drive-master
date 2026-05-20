<script setup lang="ts">
import { ref } from 'vue'
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'

const isOpen = ref(false)
const step = ref(1)
const loading = ref(false)

const schema = z.object({
  fullName: z.string().min(3, 'Name must be at least 3 characters'),
  ktpNumber: z.string().min(16, 'KTP number must be 16 digits').max(16, 'KTP number must be 16 digits')
})

const formData = reactive({
  fullName: '',
  ktpNumber: ''
})

const props = defineProps({
  modelValue: Boolean
})

const emit = defineEmits(['update:modelValue', 'complete']())

const show = computed({
  get: () => props.modelValue,
  set: (value) => {
    emit('update:modelValue', value)
  }
})

async function onSubmit(event: FormSubmitEvent<any>) {
  if (step.value < 2) {
    step.value++
    return
  }

  loading.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1500))
    console.log('[v0] Onboarding modal completed')
    emit('complete', formData)
    show.value = false
    step.value = 1
  } finally {
    loading.value = false
  }
}

function closeModal() {
  show.value = false
  step.value = 1
}
</script>

<template>
  <UModal v-model="show" prevent-close>
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold">Complete Your Profile</h2>
          <UButton 
            color="gray" 
            variant="ghost" 
            icon="i-lucide-x" 
            @click="closeModal"
            :disabled="loading"
          />
        </div>
      </template>

      <!-- Stepper -->
      <div class="mb-6">
        <div class="flex gap-2">
          <div 
            v-for="i in 2"
            :key="i"
            :class="[
              'h-1 flex-1 rounded-full transition-colors',
              step >= i ? 'bg-primary' : 'bg-muted/30'
            ]"
          />
        </div>
      </div>

      <UForm 
        :schema="schema" 
        :state="formData" 
        class="space-y-4"
        @submit="onSubmit"
      >
        <!-- Step 1: Full Name -->
        <div v-if="step === 1" class="space-y-4">
          <div>
            <h3 class="font-semibold mb-2">What's your full name?</h3>
            <p class="text-sm text-muted mb-4">
              This will be used on your official certificate
            </p>
          </div>

          <UFormField name="fullName" label="Full Name" required>
            <UInput 
              v-model="formData.fullName"
              placeholder="Enter your full name"
              icon="i-lucide-user"
              autofocus
              size="lg"
            />
          </UFormField>
        </div>

        <!-- Step 2: KTP Number -->
        <div v-if="step === 2" class="space-y-4">
          <div>
            <h3 class="font-semibold mb-2">Verify with Your KTP</h3>
            <p class="text-sm text-muted mb-4">
              Enter your 16-digit KTP number for certificate verification
            </p>
          </div>

          <UFormField name="ktpNumber" label="KTP Number (16 Digits)" required>
            <UInput 
              v-model="formData.ktpNumber"
              placeholder="e.g., 3520123456789012"
              icon="i-lucide-credit-card"
              autofocus
              size="lg"
              @input="formData.ktpNumber = formData.ktpNumber.replace(/\D/g, '').slice(0, 16)"
            />
          </UFormField>

          <UAlert icon="i-lucide-shield-check" color="primary" size="sm">
            <template #description>
              Your data is encrypted and secure
            </template>
          </UAlert>
        </div>

        <!-- Actions -->
        <div class="flex gap-3 pt-4 border-t">
          <UButton 
            v-if="step > 1"
            type="button"
            label="Back"
            color="neutral"
            variant="outline"
            @click="step--"
            block
          />
          <UButton 
            type="submit"
            :label="step === 2 ? 'Complete' : 'Next'"
            :loading="loading"
            :icon="step === 2 ? 'i-lucide-check' : 'i-lucide-arrow-right'"
            :block="step === 1"
            :class="step === 1 ? 'col-span-2' : ''"
          />
        </div>
      </UForm>
    </UCard>
  </UModal>
</template>
