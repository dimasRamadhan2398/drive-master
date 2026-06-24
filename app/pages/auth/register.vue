<script setup lang="ts">
import { z } from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";
import { reactive, ref, computed, onMounted, watch } from "vue";
import { navigateTo, useRoute } from "nuxt/app";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "~/stores/auth";
import { usePackagesStore } from "~/stores/packages";
import { useEnrollmentsStore } from "~/stores/enrollments";

definePageMeta({
  layout: "blank",
});

const { t } = useI18n();
const route = useRoute();
const authStore = useAuthStore();
const packagesStore = usePackagesStore();
const enrollmentsStore = useEnrollmentsStore();

const planFromQuery = computed(() => route.query.plan as string | undefined);

const currentStep = ref(0);
const totalSteps = computed(() => (planFromQuery.value ? 2 : 3));

const showPrivacyModal = ref(false);
const showTermsModal = ref(false);

// Fetch packages on mount
onMounted(() => {
  packagesStore.fetchPackages();
});

// Convert packages to radio options
const packageOptions = computed(() => {
  return packagesStore.activePackages.map((pkg) => ({
    label: `${pkg.name} - Rp ${pkg.discountPrice.toLocaleString("id-ID")}${
      pkg.isPopular ? ` (${t("packages.popular")})` : ""
    }`,
    value: pkg.id,
  }));
});

// Get selected package details
const selectedPackage = computed(() => {
  return packagesStore.getPackageById(formData.package);
});

// Step 0: Personal Info
const step0Schema = computed(() =>
  z.object({
    firstName: z.string().min(1, t("register.validation.firstNameRequired")),
    lastName: z.string().min(1, t("register.validation.lastNameRequired")),
    email: z.string().email(t("register.validation.emailRequired")),
    phone: z.string().min(10, t("register.validation.phoneRequired")),
    birthDate: z.string().min(1, t("register.validation.birthDateRequired")),
  })
);

// Step 1: Package Selection
const step1Schema = computed(() =>
  z.object({
    package: z.string().min(1, t("register.validation.packageRequired")),
    startDate: z.string().optional(),
  })
);

// Step 2: Account
const step2Schema = computed(() =>
  z
    .object({
      password: z.string().min(8, t("register.validation.passwordMinLength")),
      confirmPassword: z.string(),
      terms: z
        .boolean()
        .refine((val) => val === true, t("register.validation.termsRequired")),
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: t("register.validation.passwordsNotMatch"),
      path: ["confirmPassword"],
    })
);

const formData = reactive({
  // Step 0
  firstName: "",
  lastName: "",
  email: "",
  phone: "",
  birthDate: "",
  // Step 1
  package: "",
  startDate: "",
  // Step 2
  password: "",
  confirmPassword: "",
  terms: false,
});

const loading = ref(false);

async function nextStep() {
  currentStep.value++;
}

async function prevStep() {
  currentStep.value--;
}

async function onSubmit(_event: FormSubmitEvent<any>) {
  console.log("[REGISTER PAGE] onSubmit() called, currentStep:", currentStep.value);

  if (currentStep.value === 0) {
    nextStep();
    return;
  }

  if (currentStep.value === 1) {
    // This is always the account creation step — register here
    loading.value = true;
    const toast = useToast();
    try {
      const username = formData.email.split("@")[0];

      const registerResponse = await authStore.register({
        firstName: formData.firstName,
        lastName: formData.lastName,
        username,
        email: formData.email,
        phoneNumber: formData.phone,
        dateOfBirth: formData.birthDate,
        password: formData.password,
        confirmPassword: formData.confirmPassword,
        roleId: 2,
      });
      console.log("[REGISTER PAGE] Registration successful:", registerResponse);

      if (import.meta.client) {
        sessionStorage.setItem("dm_reg_email", formData.email);
        sessionStorage.setItem("dm_reg_phone", formData.phone);
      }

      if (!planFromQuery.value) {
        navigateTo(`/auth/verify?email=${encodeURIComponent(formData.email)}`);
      } else {
        // Proceed to package selection step
        nextStep();
      }
    } catch (error: any) {
      console.log("[REGISTER PAGE] Registration failed:", error);
      toast.add({
        title: t("register.errors.registrationFailed"),
        description: error?.message || t("register.errors.tryAgain"),
        color: "error",
      });
    } finally {
      loading.value = false;
    }
    return;
  }

  if (currentStep.value === 2) {
    // Package selection step — create enrollment and navigate to verify
    loading.value = true;
    const toast = useToast();
    try {
      if (import.meta.client) {
        sessionStorage.setItem("dm_selected_plan", formData.package);
      }

      // Create enrollment after successful registration
      if (formData.package && authStore.user?.userId) {
        const selectedPkg = selectedPackage.value;
        if (selectedPkg) {
          await enrollmentsStore.createEnrollment({
            userId: authStore.user.userId,
            packageId: formData.package,
            price: selectedPkg.price,
            discountPrice: selectedPkg.discountPrice,
            startDate: formData.startDate || undefined,
          });
          console.log(
            "[REGISTER PAGE] Enrollment created:",
            enrollmentsStore.currentEnrollment
          );
        }
      }

      navigateTo(`/auth/verify?email=${encodeURIComponent(formData.email)}`);
    } catch (error: any) {
      console.log("[REGISTER PAGE] Error creating enrollment:", error);
      toast.add({
        title: t("register.errors.enrollmentError"),
        description: error?.message || t("register.errors.enrollmentFailed"),
        color: "error",
      });
    } finally {
      loading.value = false;
    }
    return;
  }
}

const stepItems = computed(() => {
  const items = [
    { label: t("register.steps.personalInfo"), icon: "i-lucide-user" },
    { label: t("register.steps.createAccount"), icon: "i-lucide-shield-check" },
  ];
  if (!planFromQuery.value) {
    items.push({ label: t("register.steps.selectPackage"), icon: "i-lucide-package" });
  }
  return items;
});

// If plan is passed in query, auto-select it
watch(
  () => planFromQuery.value,
  (newPlan) => {
    if (newPlan) {
      formData.package = newPlan;
    }
  },
  { immediate: true }
);
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] py-12 px-4">
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="flex items-center justify-center gap-2 mb-4">
          <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
        </div>
        <h1 class="text-2xl font-bold">{{ t("register.title") }}</h1>
        <p class="text-muted mt-2">{{ t("register.subtitle") }}</p>
      </div>

      <!-- Stepper -->
      <UStepper
        :items="stepItems"
        :model-value="currentStep"
        color="warning"
        class="mb-8"
      />

      <UCard>
        <!-- Step 0: Personal Info -->
        <UForm
          v-if="currentStep === 0"
          :schema="step0Schema"
          :state="formData"
          class="space-y-4"
          @submit="onSubmit"
        >
          <div class="grid grid-cols-2 gap-4">
            <UFormField name="firstName" :label="t('register.form.firstName')" required>
              <UInput
                v-model="formData.firstName"
                :placeholder="t('register.placeholder.firstName')"
                icon="i-lucide-user"
                size="lg"
                class="w-full"
                color="warning"
              />
            </UFormField>

            <UFormField name="lastName" :label="t('register.form.lastName')" required>
              <UInput
                v-model="formData.lastName"
                :placeholder="t('register.placeholder.lastName')"
                icon="i-lucide-user"
                size="lg"
                class="w-full"
                color="warning"
              />
            </UFormField>
          </div>

          <UFormField name="email" :label="t('register.form.email')" required>
            <UInput
              v-model="formData.email"
              type="email"
              :placeholder="t('register.placeholder.email')"
              icon="i-lucide-mail"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="phone" :label="t('register.form.phone')" required>
            <UInput
              v-model="formData.phone"
              :placeholder="t('register.placeholder.phone')"
              icon="i-lucide-phone"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="birthDate" :label="t('register.form.birthDate')" required>
            <UInput
              v-model="formData.birthDate"
              type="date"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <div class="flex justify-end pt-4">
            <UButton
              type="submit"
              :label="t('register.buttons.continue')"
              color="warning"
              trailingIcon="i-lucide-arrow-right"
              size="lg"
            />
          </div>
        </UForm>

        <!-- Step 1: Account Creation -->
        <UForm
          v-if="currentStep === 1"
          :schema="step2Schema"
          :state="formData"
          class="space-y-4"
          @submit="onSubmit"
        >
          <UAlert icon="i-lucide-user-check" color="warning">
            <template #title>{{
              t("register.almostThere", { name: formData.firstName })
            }}</template>
            <template #description>
              {{ t("register.createPassword") }}
            </template>
          </UAlert>

          <UFormField name="password" :label="t('register.form.password')" required>
            <UInput
              v-model="formData.password"
              type="password"
              :placeholder="t('register.placeholder.password')"
              icon="i-lucide-lock"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField
            name="confirmPassword"
            :label="t('register.form.confirmPassword')"
            required
          >
            <UInput
              v-model="formData.confirmPassword"
              type="password"
              :placeholder="t('register.placeholder.confirmPassword')"
              icon="i-lucide-lock"
              size="lg"
              class="w-full"
              color="warning"
            />
          </UFormField>

          <UFormField name="terms">
            <UCheckbox v-model:model-value="formData.terms" color="warning">
              <template #label>
                <span class="text-sm flex items-center gap-1">
                  {{ t("register.terms.agree") }}
                  <UButton
                    :label="t('register.terms.termsOfService')"
                    color="warning"
                    variant="ghost"
                    class="underline mx-0 px-0 h-auto py-0 text-sm"
                    @click="showTermsModal = true"
                  />
                  <UModal
                    v-model:open="showTermsModal"
                    :title="t('register.terms.termsTitle')"
                  >
                    <template #body>
                      <div class="prose dark:prose-invert max-w-none space-y-6">
                        <p>
                          {{
                            t("register.tosContent.welcome", {
                              url: "www.drivemaster.id",
                            })
                          }}
                        </p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.servicesTitle") }}
                        </h2>
                        <p>
                          {{ t("register.tosContent.servicesDesc") }}
                        </p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.accountsTitle") }}
                        </h2>
                        <p>
                          {{ t("register.tosContent.accountsDesc") }}
                        </p>
                        <ul class="list-disc list-inside ml-4">
                          <li>{{ t("register.tosContent.accountsList1") }}</li>
                          <li>{{ t("register.tosContent.accountsList2") }}</li>
                          <li>{{ t("register.tosContent.accountsList3") }}</li>
                        </ul>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.feesTitle") }}
                        </h2>
                        <p>{{ t("register.tosContent.feesDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>
                            <strong
                              >{{
                                t("register.tosContent.feesList1").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.tosContent.feesList1").split(":")[1] }}
                          </li>
                          <li>
                            <strong
                              >{{
                                t("register.tosContent.feesList2").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.tosContent.feesList2").split(":")[1] }}
                          </li>
                        </ul>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.schedulingTitle") }}
                        </h2>
                        <p>{{ t("register.tosContent.schedulingDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>
                            <strong
                              >{{
                                t("register.tosContent.schedulingList1").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.tosContent.schedulingList1").split(":")[1] }}
                          </li>
                          <li>
                            <strong
                              >{{
                                t("register.tosContent.schedulingList2").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.tosContent.schedulingList2").split(":")[1] }}
                          </li>
                          <li>
                            <strong
                              >{{
                                t("register.tosContent.schedulingList3").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.tosContent.schedulingList3").split(":")[1] }}
                          </li>
                        </ul>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.obligationsTitle") }}
                        </h2>
                        <p>{{ t("register.tosContent.obligationsDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>{{ t("register.tosContent.obligationsList1") }}</li>
                          <li>{{ t("register.tosContent.obligationsList2") }}</li>
                          <li>{{ t("register.tosContent.obligationsList3") }}</li>
                        </ul>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.liabilityTitle") }}
                        </h2>
                        <p>{{ t("register.tosContent.liabilityDesc") }}</p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.changesTitle") }}
                        </h2>
                        <p>{{ t("register.tosContent.changesDesc") }}</p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.tosContent.contactTitle") }}
                        </h2>
                        <p>{{ t("register.tosContent.contactDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>
                            {{
                              t("register.tosContent.contactEmail", {
                                email: "info@drivemaster.id",
                              })
                            }}
                          </li>
                          <li>{{ t("register.tosContent.contactPhone") }}</li>
                        </ul>
                      </div>
                    </template>
                  </UModal>
                  {{ t("register.terms.and") }}
                  <UButton
                    :label="t('register.terms.privacyPolicy')"
                    color="warning"
                    variant="ghost"
                    class="underline mx-0 px-0 h-auto py-0 text-sm"
                    @click="showPrivacyModal = true"
                  />
                  <UModal
                    v-model:open="showPrivacyModal"
                    :title="t('register.terms.privacyTitle')"
                  >
                    <template #body>
                      <div class="prose dark:prose-invert max-w-none space-y-6">
                        <p>
                          {{
                            t("register.privacyContent.welcome", {
                              url: "www.drivemaster.id",
                            })
                          }}
                        </p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.privacyContent.collectTitle") }}
                        </h2>
                        <p>{{ t("register.privacyContent.collectDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>
                            <strong
                              >{{
                                t("register.privacyContent.collectList1").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.privacyContent.collectList1").split(":")[1] }}
                          </li>
                          <li>
                            <strong
                              >{{
                                t("register.privacyContent.collectList2").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.privacyContent.collectList2").split(":")[1] }}
                          </li>
                          <li>
                            <strong
                              >{{
                                t("register.privacyContent.collectList3").split(":")[0]
                              }}:</strong
                            >
                            {{ t("register.privacyContent.collectList3").split(":")[1] }}
                          </li>
                        </ul>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.privacyContent.useTitle") }}
                        </h2>
                        <p>{{ t("register.privacyContent.useDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>{{ t("register.privacyContent.useList1") }}</li>
                          <li>{{ t("register.privacyContent.useList2") }}</li>
                          <li>{{ t("register.privacyContent.useList3") }}</li>
                          <li>{{ t("register.privacyContent.useList4") }}</li>
                          <li>{{ t("register.privacyContent.useList5") }}</li>
                        </ul>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.privacyContent.disclosureTitle") }}
                        </h2>
                        <p>{{ t("register.privacyContent.disclosureDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>
                            <strong
                              >{{
                                t("register.privacyContent.disclosureList1").split(
                                  ":"
                                )[0]
                              }}:</strong
                            >
                            {{
                              t("register.privacyContent.disclosureList1").split(":")[1]
                            }}
                          </li>
                          <li>
                            <strong
                              >{{
                                t("register.privacyContent.disclosureList2").split(
                                  ":"
                                )[0]
                              }}:</strong
                            >
                            {{
                              t("register.privacyContent.disclosureList2").split(":")[1]
                            }}
                          </li>
                          <li>
                            <strong
                              >{{
                                t("register.privacyContent.disclosureList3").split(
                                  ":"
                                )[0]
                              }}:</strong
                            >
                            {{
                              t("register.privacyContent.disclosureList3").split(":")[1]
                            }}
                          </li>
                        </ul>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.privacyContent.securityTitle") }}
                        </h2>
                        <p>{{ t("register.privacyContent.securityDesc") }}</p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.privacyContent.rightsTitle") }}
                        </h2>
                        <p>{{ t("register.privacyContent.rightsDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>{{ t("register.privacyContent.rightsList1") }}</li>
                          <li>{{ t("register.privacyContent.rightsList2") }}</li>
                          <li>{{ t("register.privacyContent.rightsList3") }}</li>
                          <li>{{ t("register.privacyContent.rightsList4") }}</li>
                          <li>{{ t("register.privacyContent.rightsList5") }}</li>
                        </ul>
                        <p>
                          {{
                            t("register.privacyContent.rightsContact", {
                              email: "info@drivemaster.id",
                            })
                          }}
                        </p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.privacyContent.changesTitle") }}
                        </h2>
                        <p>{{ t("register.privacyContent.changesDesc") }}</p>

                        <h2 class="text-2xl font-bold">
                          {{ t("register.privacyContent.contactTitle") }}
                        </h2>
                        <p>{{ t("register.privacyContent.contactDesc") }}</p>
                        <ul class="list-disc list-inside ml-4">
                          <li>
                            {{
                              t("register.privacyContent.contactEmail", {
                                email: "info@drivemaster.id",
                              })
                            }}
                          </li>
                          <li>{{ t("register.privacyContent.contactPhone") }}</li>
                        </ul>
                      </div>
                    </template>
                  </UModal>
                </span>
              </template>
            </UCheckbox>
          </UFormField>

          <div class="flex justify-between pt-4">
            <UButton
              :label="t('register.buttons.back')"
              variant="ghost"
              color="neutral"
              icon="i-lucide-arrow-left"
              @click="prevStep"
            />
            <UButton
              type="submit"
              :label="t('register.buttons.createAccount')"
              color="warning"
              :loading="loading"
              icon="i-lucide-check"
              size="lg"
            />
          </div>
        </UForm>

        <!-- Step 2: Package Selection -->
        <UForm
          v-if="currentStep === 2"
          :schema="step1Schema"
          :state="formData"
          class="space-y-4"
          @submit="onSubmit"
        >
          <UFormField name="package" :label="t('register.form.package')" required>
            <URadioGroup
              v-model="formData.package"
              :items="packageOptions"
              orientation="vertical"
              color="warning"
            />
          </UFormField>

          <!-- Package Summary -->
          <UAlert
            v-if="selectedPackage"
            :icon="selectedPackage.isPopular ? 'i-lucide-star' : 'i-lucide-info'"
            :color="selectedPackage.isPopular ? 'warning' : 'neutral'"
          >
            <template #title>
              {{ selectedPackage.name }}
            </template>
            <template #description>
              <ul class="mt-2 space-y-1 text-sm">
                <li v-for="(feature, idx) in selectedPackage.features" :key="idx">
                  {{ feature }}
                </li>
                <li class="font-semibold mt-2">
                  {{ selectedPackage.sessions }} {{ t("register.sessions") }}
                </li>
                <li class="font-semibold text-warning">
                  Rp {{ selectedPackage.discountPrice.toLocaleString("id-ID") }}
                </li>
              </ul>
            </template>
          </UAlert>

          <div class="text-center">
            <p class="text-sm text-muted">
              {{ t("register.letMeDecide") }}
              <NuxtLink
                to="/auth/onboarding"
                class="text-warning font-medium hover:underline"
              >
                {{ t("register.goToOnboarding") }}
              </NuxtLink>
            </p>
          </div>

          <div class="flex justify-between pt-4">
            <UButton
              :label="t('register.buttons.back')"
              variant="ghost"
              color="neutral"
              icon="i-lucide-arrow-left"
              @click="prevStep"
            />
            <UButton
              type="submit"
              :label="t('register.buttons.proceedToPayment')"
              color="warning"
              trailingIcon="i-lucide-arrow-right"
              size="lg"
            />
          </div>
        </UForm>

        <template #footer>
          <div class="text-center">
            <p class="text-sm text-muted">
              {{ t("register.alreadyHaveAccount") }}
              <NuxtLink to="/auth/login" class="text-warning font-medium hover:underline">
                {{ t("register.signIn") }}
              </NuxtLink>
            </p>
          </div>
        </template>
      </UCard>
    </div>
  </div>
</template>
