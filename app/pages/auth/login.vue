<script setup lang="ts">
<<<<<<< HEAD
import { z } from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";
import { reactive, ref } from "vue";
import { useAuthStore } from "~/stores/auth";
import guest from "~/middleware/guest";

definePageMeta({
  layout: "blank",
  middleware: [guest],
});

const schema = z.object({
  email: z.string().email("Please enter a valid email"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  remember: z.boolean().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive({
  email: "",
  password: "",
  remember: false,
});

const loading = ref(false);
const errorMessage = ref("");

const authStore = useAuthStore();
const toast = useToast();

async function onSubmit(event: FormSubmitEvent<Schema>) {
  loading.value = true;
  errorMessage.value = "";

  try {
    await authStore.login({
      email: state.email,
      password: state.password,
      remember: state.remember,
    });

    toast.add({
      title: "Welcome back!",
      description: "You have successfully signed in.",
      color: "success",
    });

    navigateTo("/dashboard");
  } catch (error: any) {
    const errorMsg =
      error?.message || error?.data?.error?.message || "Login failed";
    errorMessage.value = errorMsg;
    toast.add({
      title: "Login Failed",
      description: errorMsg,
      color: "error",
    });
  } finally {
    loading.value = false;
  }
=======
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
>>>>>>> main
}
</script>

<template>
<<<<<<< HEAD
  <div
    class="min-h-[calc(100vh-200px)] flex items-center justify-center py-12 px-4"
  >
=======
  <div class="min-h-[calc(100vh-200px)] flex items-center justify-center py-12 px-4">
>>>>>>> main
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <div class="flex items-center justify-center gap-2 mb-4">
<<<<<<< HEAD
            <img
              src="/drive-master-logo2.png"
              alt="Drive Master Logo"
              class="h-16"
            />
=======
            <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
>>>>>>> main
          </div>
          <h1 class="text-2xl font-bold">Welcome Back</h1>
          <p class="text-muted mt-2">Sign in to access your member dashboard</p>
        </div>
      </template>

<<<<<<< HEAD
      <UAlert v-if="errorMessage" color="error" variant="subtle" class="mb-4">
        {{ errorMessage }}
      </UAlert>

      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
        <UFormField name="email" label="Email Address">
          <UInput
            v-model="state.email"
=======
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormField name="email" label="Email Address">
          <UInput 
            v-model="state.email" 
>>>>>>> main
            type="email"
            placeholder="you@example.com"
            icon="i-lucide-mail"
            size="lg"
            class="w-full"
          />
        </UFormField>

        <UFormField name="password" label="Password">
<<<<<<< HEAD
          <UInput
            v-model="state.password"
=======
          <UInput 
            v-model="state.password" 
>>>>>>> main
            type="password"
            placeholder="Enter your password"
            icon="i-lucide-lock"
            size="lg"
            class="w-full"
          />
        </UFormField>

        <div class="flex items-center justify-between">
<<<<<<< HEAD
          <UCheckbox
            v-model="state.remember"
            label="Remember me"
            color="warning"
          />
          <NuxtLink
            to="/forgot-password"
            class="text-sm text-warning hover:underline"
          >
=======
          <UCheckbox v-model="state.remember" label="Remember me" color="warning" />
          <NuxtLink to="/forgot-password" class="text-sm text-warning hover:underline">
>>>>>>> main
            Forgot password?
          </NuxtLink>
        </div>

<<<<<<< HEAD
        <UButton
          type="submit"
          label="Sign In"
          color="warning"
          :loading="loading"
          block
          size="lg"
        />
=======
        <UButton type="submit" label="Sign In" color="warning" :loading="loading" block size="lg" />
>>>>>>> main
      </UForm>

      <template #footer>
        <div class="text-center space-y-4">
          <p class="text-sm text-muted">
            Don&apos;t have an account?
<<<<<<< HEAD
            <NuxtLink
              to="/auth/register"
              class="text-warning font-medium hover:underline"
            >
              Register here
            </NuxtLink>
          </p>

          <USeparator label="or" />

          <NuxtLink
            to="https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi"
            target="_blank"
            class="block"
          >
            <UButton
              label="Contact Support"
              icon="i-simple-icons-whatsapp"
              color="primary"
=======
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
>>>>>>> main
              variant="outline"
              block
            />
          </NuxtLink>
        </div>
      </template>
    </UCard>
  </div>
</template>
