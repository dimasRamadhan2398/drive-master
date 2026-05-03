<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive, ref } from 'vue'
import { useAuth } from '~/composables/useAuth'

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

const { login } = useAuth()

async function onSubmit(event: FormSubmitEvent<Schema>) {
  loading.value = true
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1000))
  
  // Set mock user data
  login({
    name: 'John Doe',
    email: state.email,
    avatar: 'https://i.pravatar.cc/150?u=johndoe'
  })
  
  loading.value = false
  
  // Navigate to dashboard
  navigateTo('/dashboard')
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
