<script setup lang="ts">
import { z } from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";
import { reactive, ref } from "vue";
import { authService } from "~/services/authService";

definePageMeta({
  layout: "blank",
});

const schema = z.object({
  email: z.string().email("Please enter a valid email"),
});

type Schema = z.output<typeof schema>;

const state = reactive({
  email: "",
});

const loading = ref(false);
const success = ref(false);
const error = ref<string | null>(null);

async function onSubmit(_event: FormSubmitEvent<Schema>) {
  loading.value = true;
  error.value = null;

  try {
    await authService.forgotPassword(state.email);
    success.value = true;
  } catch (err) {
    error.value = authService.parseError(err);
  } finally {
    loading.value = false;
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
          <h1 class="text-2xl font-bold">Forgot Password</h1>
          <p class="text-muted mt-2">
            Enter your email address and we'll send you a link to reset your password
          </p>
        </div>
      </template>

      <UAlert v-if="error" color="error" variant="soft" class="mb-4" :title="error" />

      <UAlert
        v-if="success"
        color="success"
        variant="soft"
        class="mb-4"
        title="Check your email"
      >
        <template #description>
          If an account exists with that email, we've sent password reset instructions.
        </template>
      </UAlert>

      <UForm
        v-if="!success"
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
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

        <UButton
          type="submit"
          label="Send Reset Link"
          color="warning"
          :loading="loading"
          block
          size="lg"
        />
      </UForm>

      <template #footer>
        <div class="text-center">
          <p class="text-sm text-muted">
            Remember your password?
            <NuxtLink to="/auth/login" class="text-warning font-medium hover:underline">
              Sign in
            </NuxtLink>
          </p>
        </div>
      </template>
    </UCard>
  </div>
</template>
