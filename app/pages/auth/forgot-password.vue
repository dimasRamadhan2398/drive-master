<script setup lang="ts">
import { z } from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";
import { reactive, ref } from "vue";
import { authService } from "~/services/authService";

const { t } = useI18n()

definePageMeta({
  layout: "blank",
});

const schema = computed(() => z.object({
  email: z.string().email(t('validation.email')),
}));

type Schema = z.output<ReturnType<typeof schema>>;

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
          <h1 class="text-2xl font-bold">{{ t('auth.forgotPasswordTitle') }}</h1>
          <p class="text-muted mt-2">
            {{ t('auth.forgotPasswordDesc') }}
          </p>
        </div>
      </template>

      <UAlert v-if="error" color="error" variant="soft" class="mb-4" :title="error" />

      <UAlert
        v-if="success"
        color="success"
        variant="soft"
        class="mb-4"
        :title="t('auth.checkEmail')"
      >
        <template #description>
          {{ t('auth.checkEmailDesc') }}
        </template>
      </UAlert>

      <UForm
        v-if="!success"
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
        <UFormField name="email" :label="t('auth.email')">
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
          :label="t('auth.sendResetLink')"
          color="warning"
          :loading="loading"
          block
          size="lg"
        />
      </UForm>

      <template #footer>
        <div class="text-center">
          <p class="text-sm text-muted">
            {{ t('auth.rememberPassword') }}
            <NuxtLink to="/auth/login" class="text-warning font-medium hover:underline">
              {{ t('auth.signIn') }}
            </NuxtLink>
          </p>
        </div>
      </template>
    </UCard>
  </div>
</template>
