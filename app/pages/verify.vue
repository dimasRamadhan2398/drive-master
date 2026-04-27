<script setup lang="ts">
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'blank'
})

const router = useRouter()

const otp = ref('')
const email = ref('user@example.com') // Would come from registration
const loading = ref(false)
const showResend = ref(false)
const resendCountdown = ref(0)
const errorMessage = ref('')

// OTP input handling
const otpArray = computed({
  get: () => otp.value.split(''),
  set: (value) => {
    otp.value = value.join('').slice(0, 6)
  }
})

// Focus next input on digit entry
const handleOtpInput = (index: number, event: Event) => {
  const input = event.target as HTMLInputElement
  const value = input.value

  if (value.length > 1) {
    // Remove extra characters
    input.value = value.slice(-1)
  }

  if (value && index < 5) {
    // Move to next input
    const nextInput = (document.querySelector(`#otp-${index + 1}`) as HTMLInputElement)
    nextInput?.focus()
  }

  // Update otp ref
  const otpInputs = Array.from(document.querySelectorAll('input[data-otp-input]')) as HTMLInputElement[]
  otp.value = otpInputs.map(input => input.value).join('')
}

// Handle paste
const handlePaste = (event: ClipboardEvent) => {
  event.preventDefault()
  const pastedData = event.clipboardData?.getData('text') || ''
  const digits = pastedData.replace(/\D/g, '').slice(0, 6)

  otp.value = digits
  
  // Focus last input or next empty one
  const otpInputs = Array.from(document.querySelectorAll('input[data-otp-input]')) as HTMLInputElement[]
  otpInputs.forEach((input, index) => {
    input.value = digits[index] || ''
  })
  
  if (digits.length === 6) {
    otpInputs[5]?.blur()
  }
}

// Verify OTP
async function verifyOTP() {
  if (otp.value.length !== 6) {
    errorMessage.value = 'Please enter a complete 6-digit code'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    // Simulate API call for verification
    // In real implementation, this would verify with backend
    await new Promise(resolve => setTimeout(resolve, 1500))

    // Mock verification - in real app, check against backend
    // For demo, accept "123456" as valid
    if (otp.value === '123456' || true) { // Allow any code for demo
      console.log('[v0] OTP verified successfully:', otp.value)
      navigateTo('/onboarding')
    } else {
      errorMessage.value = 'Invalid OTP. Please try again.'
    }
  } catch (error) {
    errorMessage.value = 'Verification failed. Please try again.'
    console.error('[v0] Verification error:', error)
  } finally {
    loading.value = false
  }
}

// Resend OTP
async function resendOTP() {
  loading.value = true
  showResend.value = false
  resendCountdown.value = 30

  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    console.log('[v0] OTP resent to:', email.value)
    
    // Start countdown
    const interval = setInterval(() => {
      resendCountdown.value--
      if (resendCountdown.value <= 0) {
        clearInterval(interval)
        showResend.value = true
      }
    }, 1000)
  } finally {
    loading.value = false
  }
}

// Handle Enter key
const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && otp.value.length === 6) {
    verifyOTP()
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4 bg-muted/20 flex items-center">
    <div class="max-w-md mx-auto w-full">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="flex items-center justify-center gap-2 mb-4">
          <UIcon name="i-lucide-shield-check" class="size-8 text-warning" />
          <span class="text-xl font-bold">Email Verification</span>
        </div>
        <h1 class="text-2xl font-bold">Verify Your Email</h1>
        <p class="text-muted mt-2">We sent a 6-digit code to</p>
        <p class="font-medium">{{ email }}</p>
      </div>

      <UCard>
        <div class="space-y-6">
          <!-- OTP Input Fields -->
          <div>
            <label class="block text-sm font-medium mb-4">Enter Verification Code</label>
            <div class="flex gap-2 justify-between">
              <input
                v-for="i in 6"
                :key="i"
                :id="`otp-${i - 1}`"
                data-otp-input
                type="text"
                inputmode="numeric"
                maxlength="1"
                class="w-12 h-12 text-center text-2xl font-bold border-2 border-default rounded-lg focus:border-warning focus:outline-none transition-colors"
                @input="handleOtpInput(i - 1, $event)"
                @paste="handlePaste"
                @keydown="handleKeyDown"
              />
            </div>
          </div>

          <!-- Error Message -->
          <Transition name="fade">
            <UAlert 
              v-if="errorMessage"
              icon="i-lucide-alert-circle"
              color="error"
              :title="errorMessage"
              class="mb-0"
            />
          </Transition>

          <!-- Info Message -->
          <UAlert 
            icon="i-lucide-info"
            color="warning"
            variant="subtle"
          >
            <template #description>
              <div class="text-sm space-y-1">
                <p>Demo code: <span class="font-mono font-medium">123456</span></p>
                <p>Or paste any 6-digit number</p>
              </div>
            </template>
          </UAlert>

          <!-- Verify Button -->
          <UButton 
            label="Verify Email"
            icon="i-lucide-check"
            :loading="loading"
            :disabled="otp.length !== 6 || loading"
            block
            color="warning"
            size="lg"
            @click="verifyOTP"
          />

          <!-- Resend Code -->
          <div class="text-center">
            <p class="text-sm text-muted mb-3">Didn't receive the code?</p>
            <div v-if="resendCountdown > 0" class="text-sm">
              <span class="text-muted">Resend code in </span>
              <span class="font-medium text-warning">{{ resendCountdown }}s</span>
            </div>
            <UButton
              v-else
              label="Resend Code"
              variant="ghost"
              size="sm"
              color="warning"
              :loading="loading"
              @click="resendOTP"
            />
          </div>
        </div>

        <!-- Footer -->
        <template #footer>
          <div class="text-center text-sm text-muted">
            <p class="mb-2">
              Wrong email? 
              <NuxtLink to="/register" class="text-warning hover:underline">
                Go back to registration
              </NuxtLink>
            </p>
            <p>
              <NuxtLink to="/" class="text-warning hover:underline">
                Back to Home
              </NuxtLink>
            </p>
          </div>
        </template>
      </UCard>

      <!-- Support Info -->
      <div class="mt-6 p-4 bg-muted/50 rounded-lg text-center">
        <p class="text-xs text-muted">
          Check your spam folder if you don't see the email
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
