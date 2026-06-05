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
  password: z.string().min(8, 'Password must be at least 8 characters')
})

type Schema = z.output<typeof schema>

const state = reactive({
  email: '',
  password: ''
})

const loading = ref(false)
const error = ref<string | null>(null)
const authStore = useAuthStore()

// Direct cookie access to ensure rehydration works
const authToken = useCookie('auth_token')
const userData = useCookie('user_data')
const refreshToken = useCookie('refresh_token')

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
    // Navigate to admin dashboard after successful login
    navigateTo('/admin')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 bg-gray-50 dark:bg-gray-900">
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <div class="flex items-center justify-center gap-2 mb-4">
            <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
          </div>
          <h1 class="text-2xl font-bold">Admin Portal</h1>
          <p class="text-muted mt-2">Sign in to access the management dashboard</p>
        </div>
      </template>

      <UAlert v-if="error" color="error" variant="soft" class="mb-4" :title="error" />

      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormField name="email" label="Admin Email">
          <UInput 
            v-model="state.email" 
            type="email"
            placeholder="admin@example.com"
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

        <UButton type="submit" label="Sign In to Admin" color="warning" :loading="loading" block size="lg" />
      </UForm>

      <template #footer>
        <div class="text-center">
          <NuxtLink to="/" class="text-sm text-warning hover:underline flex items-center justify-center gap-1">
            <UIcon name="i-lucide-arrow-left" class="size-4" />
            Back to Website
          </NuxtLink>
        </div>
      </template>
    </UCard>
  </div>
</template>
