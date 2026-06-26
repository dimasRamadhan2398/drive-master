<script setup lang="ts">
import { z } from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";
import { reactive, ref, onMounted } from "vue";
import { useAuthStore } from "~/stores/auth";

const { t } = useI18n()

definePageMeta({
  layout: "blank",
});

const schema = computed(() => z.object({
  email: z.string().email(t('validation.email')),
  password: z.string().min(8, t('validation.password')),
  remember: z.boolean().optional(),
}));

type Schema = z.output<ReturnType<typeof schema>>;

const state = reactive({
  email: "",
  password: "",
  remember: false,
});

const loading = ref(false);
const error = ref<string | null>(null);

// Direct cookie access to ensure rehydration works
const authToken = useCookie("auth_token");
const userData = useCookie("user_data");
const refreshToken = useCookie("refresh_token");

const authStore = useAuthStore();

// Rehydrate from cookies on mount if valid cookies exist
onMounted(() => {
  if (authToken.value && userData.value) {
    try {
      const user = JSON.parse(userData.value);
      authStore.setAuth(user, authToken.value, refreshToken.value || undefined);
    } catch {
      // Invalid cookie data, clear them
      authToken.value = null;
      userData.value = null;
    }
  }
});

async function onSubmit(_event: FormSubmitEvent<Schema>) {
  loading.value = true;
  error.value = null;

  try {
    await authStore.login({
      email: state.email,
      password: state.password,
    });

    // Navigate based on user role after successful login
    if (authStore.userRole?.toLowerCase().includes("admin")) {
      navigateTo("/admin");
    } else {
      // Check if student has entitlements (has purchased a package)
      const profile = await authStore.fetchMemberProfile();
      const hasEntitlements =
        profile?.entitlements && profile.entitlements.length > 0;

      if (hasEntitlements) {
        navigateTo("/dashboard");
      } else {
        // Redirect to onboarding if no entitlements
        navigateTo("/auth/onboarding");
      }
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : t('auth.loginError');
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
          <h1 class="text-2xl font-bold">{{ t('auth.welcomeBack') }}</h1>
          <p class="text-muted mt-2">{{ t('auth.signInDashboard') }}</p>
        </div>
      </template>

      <UAlert v-if="error" color="error" variant="soft" class="mb-4" :title="error" />

      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
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

        <UFormField name="password" :label="t('auth.password')">
          <UInput
            v-model="state.password"
            type="password"
            :placeholder="t('auth.password')"
            icon="i-lucide-lock"
            size="lg"
            class="w-full"
          />
        </UFormField>

        <div class="flex items-center justify-between">
          <UCheckbox v-model="state.remember" :label="t('auth.rememberMe')" color="warning" />
          <NuxtLink to="/auth/forgot-password" class="text-sm text-warning hover:underline">
            {{ t('auth.forgotPassword') }}
          </NuxtLink>
        </div>

        <UButton
          type="submit"
          :label="t('auth.signIn')"
          color="warning"
          :loading="loading"
          block
          size="lg"
        />
      </UForm>

      <template #footer>
        <div class="text-center space-y-4">
          <p class="text-sm text-muted">
            {{ t('auth.noAccount') }}
            <NuxtLink
              to="/auth/register"
              class="text-warning font-medium hover:underline"
            >
              {{ t('auth.registerHere') }}
            </NuxtLink>
          </p>

          <USeparator :label="t('auth.orContinueWith').split(' ').pop()" />

          <NuxtLink
            to="https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20perlu%20bantuan%20di%20website%20ketika...."
            target="_blank"
            class="block"
          >
            <UButton
              :label="t('auth.contactSupport')"
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
