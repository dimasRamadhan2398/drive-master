<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive, ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'blank'
})

const schema = z.object({
  email: z.string().email('Please enter a valid email'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  remember: z.boolean().optional()
})

type Schema = z.output<typeof schema>

const state = reactive({
  email: '',
  password: '',
  remember: false
})

const loading = ref(false)
const error = ref<string | null>(null)

// Direct cookie access to ensure rehydration works
const authToken = useCookie('auth_token')
const userData = useCookie('user_data')
const refreshToken = useCookie('refresh_token')

const authStore = useAuthStore()

// Rehydrate from cookies on mount if valid cookies exist
onMounted(() => {
  if (authToken.value && userData.value) {
    try {
      const user = JSON.parse(userData.value)
      authStore.setAuth(user, authToken.value, refreshToken.value || undefined)
    } catch {
      // Invalid cookie data, clear them
      authToken.value = null
      userData.value = null
    }
  }
})

async function onSubmit(_event: FormSubmitEvent<Schema>) {
  loading.value = true
  error.value = null

  try {
    await authStore.login({
      email: state.email,
      password: state.password
    })

    // Navigate based on user role after successful login
    if (authStore.userRole?.toLowerCase().includes('admin')) {
      navigateTo('/admin')
    } else {
      navigateTo('/dashboard')
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] flex items-center justify-center py-12 px-4">
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <div class="flex items-center justify-center gap-2 mb-4">
            <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
          </div>
          <h1 class="text-2xl font-bold">Welcome Back</h1>
          <p class="text-muted mt-2">Sign in to access your member dashboard</p>
        </div>
      </template>

      <UAlert v-if="error" color="error" variant="soft" class="mb-4" :title="error" />

      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormField name="email" label="Email Address">
          <UInput 
            v-model="state.email" 
            type="email"
            placeholder="you@example.com"
            icon="i-lucide-mail"
            size="lg"
            class="w-full"
          />
        </UFormField>

        <UFormField name="password" label="Password">
          <UInput 
            v-model="state.password" 
            type="password"
            placeholder="Enter your password"
            icon="i-lucide-lock"
            size="lg"
            class="w-full"
          />
        </UFormField>

        <div class="flex items-center justify-between">
          <UCheckbox v-model="state.remember" label="Remember me" color="warning" />
          <NuxtLink to="/forgot-password" class="text-sm text-warning hover:underline">
            Forgot password?
          </NuxtLink>
        </div>

        <UButton type="submit" label="Sign In" color="warning" :loading="loading" block size="lg" />
      </UForm>

      <template #footer>
        <div class="text-center space-y-4">
          <p class="text-sm text-muted">
            Don&apos;t have an account?
            <NuxtLink to="/auth/register" class="text-warning font-medium hover:underline">
              Register here
            </NuxtLink>
          </p>
          
          <USeparator label="or" />

          <NuxtLink to="https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi" target="_blank" class="block">
            <UButton 
              label="Contact Support" 
              icon="i-simple-icons-whatsapp" 
              color="primary" 
              variant="outline"
              block
            />
          </NuxtLink>
        </div>
      </template>
    </UCard>
  </div>
</template>
