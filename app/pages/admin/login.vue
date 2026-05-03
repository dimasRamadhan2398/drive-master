<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive, ref } from 'vue'

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
const { adminLogin } = useAuth()

async function onSubmit(event: FormSubmitEvent<Schema>) {
  loading.value = true
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1000))
  
  // Set mock admin data
  adminLogin({
    name: 'Admin User',
    email: state.email,
    role: 'Administrator'
  })
  
  loading.value = false
  
  // Navigate to admin dashboard
  navigateTo('/admin')
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
