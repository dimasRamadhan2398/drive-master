<script setup lang="ts">
import { useToast } from "@nuxt/ui/runtime/composables/useToast.js";
import { reactive, ref } from "vue";

const { t } = useI18n();
definePageMeta({ layout: "admin" });

const toast = useToast();
const settingsStore = useSettingsStore();
const instructorsStore = useInstructorsStore();
const vehiclesStore = useVehiclesStore();

// General settings - synced with store
const generalSettings = reactive({
  businessName: "",
  email: "",
  phone: "",
  fax: "",
  whatsApp: "",
  address: "",
  hoursMonFri: "",
  hoursSatSun: "",
  hoursNightShift: "",
  promoEndDate: null as string | null,
  notifyEmail: false,
  notifySms: false,
  notifyWhatsApp: false,
});

const notificationSettings = reactive({
  emailNotifications: false,
  whatsappNotifications: false,
  reminderHours: 24,
  adminAlerts: false,
  newUserRegistration: false,
  newPackagePurchase: false,
});

// Computed ref for vehicles from store
const vehicles = computed(() => vehiclesStore.vehicles);

// Vehicle CRUD State
const isVehicleModalOpen = ref(false);
const isEditingVehicle = ref(false);
const editingVehicleId = ref<string | null>(null);
const vehicleImageRef = ref<HTMLInputElement | null>(null);
const vehicleImageFile = ref<File | null>(null);
const vehicleForm = reactive({
  brand: "",
  model: "",
  year: new Date().getFullYear(),
  licensePlate: "",
  color: "",
  transmission: "automatic" as "manual" | "automatic",
  status: "available" as "available" | "in_use" | "maintenance" | "retired",
  mileage: 0,
  imageUrl: "",
  notes: "",
});

function openNewVehicle() {
  isEditingVehicle.value = false;
  editingVehicleId.value = null;
  vehicleImageFile.value = null;
  Object.assign(vehicleForm, {
    brand: "",
    model: "",
    year: new Date().getFullYear(),
    licensePlate: "",
    color: "",
    transmission: "automatic",
    status: "available",
    mileage: 0,
    imageUrl: "",
    notes: "",
  });
  isVehicleModalOpen.value = true;
}

function openEditVehicle(vehicle: any) {
  isEditingVehicle.value = true;
  editingVehicleId.value = vehicle.id;
  vehicleImageFile.value = null;
  Object.assign(vehicleForm, {
    brand: vehicle.brand,
    model: vehicle.model,
    year: vehicle.year,
    licensePlate: vehicle.licensePlate,
    color: vehicle.color,
    transmission: vehicle.transmission,
    status: vehicle.status,
    mileage: vehicle.mileage,
    imageUrl: vehicle.imageUrl,
    notes: vehicle.notes,
  });
  isVehicleModalOpen.value = true;
}

async function saveVehicle() {
  if (!vehicleForm.brand || !vehicleForm.model || !vehicleForm.year) {
    toast.add({
      title: t("common.error"),
      description: "Please fill in brand, model, and year",
      color: "error",
    });
    return;
  }

  if (isEditingVehicle.value && editingVehicleId.value) {
    const result = await vehiclesStore.updateVehicle(editingVehicleId.value, {
      brand: vehicleForm.brand,
      model: vehicleForm.model,
      year: vehicleForm.year,
      licensePlate: vehicleForm.licensePlate,
      color: vehicleForm.color,
      transmission: vehicleForm.transmission,
      status: vehicleForm.status,
      mileage: vehicleForm.mileage,
      imageUrl: vehicleForm.imageUrl,
      notes: vehicleForm.notes,
    });

    if (result) {
      toast.add({
        title: "Vehicle Updated",
        description: "Vehicle details saved.",
        color: "success",
      });
      isVehicleModalOpen.value = false;
    } else {
      toast.add({
        title: t("common.error"),
        description: "Failed to update vehicle.",
        color: "error",
      });
    }
  } else {
    const result = await vehiclesStore.createVehicle({
      brand: vehicleForm.brand,
      model: vehicleForm.model,
      year: vehicleForm.year,
      licensePlate: vehicleForm.licensePlate,
      color: vehicleForm.color,
      transmission: vehicleForm.transmission,
      status: vehicleForm.status,
      imageUrl: vehicleForm.imageUrl,
      notes: vehicleForm.notes,
    });

    if (result) {
      toast.add({
        title: "Vehicle Added",
        description: "New vehicle created.",
        color: "success",
      });
      isVehicleModalOpen.value = false;
    } else {
      toast.add({
        title: t("common.error"),
        description: "Failed to create vehicle.",
        color: "error",
      });
    }
  }
}

async function deleteVehicle() {
  if (editingVehicleId.value) {
    const success = await vehiclesStore.deleteVehicle(editingVehicleId.value);
    if (success) {
      toast.add({
        title: "Vehicle Deleted",
        description: "Vehicle has been removed.",
        color: "error",
      });
      isVehicleModalOpen.value = false;
    } else {
      toast.add({
        title: t("common.error"),
        description: "Failed to delete vehicle.",
        color: "error",
      });
    }
  }
}

function triggerVehicleImagePicker() {
  vehicleImageRef.value?.click();
}

function clearVehicleImage() {
  vehicleForm.imageUrl = "";
  vehicleImageFile.value = null;
  if (vehicleImageRef.value) {
    vehicleImageRef.value.value = "";
  }
}

function handleVehicleImageUpload(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];

  if (!file) return;

  if (!["image/jpeg", "image/png", "image/webp"].includes(file.type)) {
    toast.add({
      title: "Invalid File",
      description: "Please upload JPG, PNG, or WebP.",
      color: "error",
    });
    input.value = "";
    return;
  }

  if (file.size > 5 * 1024 * 1024) {
    toast.add({
      title: "File Too Large",
      description: "Image exceeds 5MB limit.",
      color: "error",
    });
    input.value = "";
    return;
  }

  vehicleImageFile.value = file;

  const reader = new FileReader();
  reader.onload = (e) => {
    vehicleForm.imageUrl = e.target?.result as string;
  };
  reader.readAsDataURL(file);
  input.value = "";
}

// Instructor CRUD State
const isInstructorModalOpen = ref(false);
const isEditingInstructor = ref(false);
const instructorUserId = ref("");
const instructorPhotoUrl = ref("");
const instructorFile = ref<File | null>(null);
const instructorForm = ref({
  firstName: "",
  lastName: "",
  phoneNumber: "",
  email: "",
  bnspCertificateNumber: "",
  licenseNumber: "",
  yearsOfExperience: 0,
  description: "",
  status: "active",
});

// Computed ref for template
const instructors = computed(() => instructorsStore.instructors);

function openNewInstructor() {
  isEditingInstructor.value = false;
  instructorUserId.value = "";
  instructorPhotoUrl.value = "";
  instructorFile.value = null;
  instructorForm.value = {
    firstName: "",
    lastName: "",
    phoneNumber: "",
    email: "",
    bnspCertificateNumber: "",
    licenseNumber: "",
    yearsOfExperience: 0,
    description: "",
    status: "active",
  };
  isInstructorModalOpen.value = true;
}

function openEditInstructor(instructor: any) {
  isEditingInstructor.value = true;
  instructorUserId.value = instructor.userId || "";
  instructorPhotoUrl.value = instructor.image || "";
  const nameParts = instructor.name?.split(" ") || [];
  instructorForm.value = {
    firstName: nameParts[0] || "",
    lastName: nameParts.slice(1).join(" ") || "",
    phoneNumber: instructor.phone || "",
    email: instructor.email || "",
    bnspCertificateNumber: instructor.certifications?.[0]?.replace("BNSP: ", "") || "",
    licenseNumber: "",
    yearsOfExperience: instructor.yearsOfExperience || 0,
    description: instructor.description || "",
    status: instructor.status || "active",
  };
  isInstructorModalOpen.value = true;
}

async function saveInstructor() {
  const data: any = {
    firstName: instructorForm.value.firstName,
    lastName: instructorForm.value.lastName,
    phoneNumber: instructorForm.value.phoneNumber,
    description: instructorForm.value.description,
    yearsOfExperience: instructorForm.value.yearsOfExperience,
  };

  if (instructorForm.value.bnspCertificateNumber) {
    data.bnspCertificateNumber = instructorForm.value.bnspCertificateNumber;
  }
  if (instructorForm.value.licenseNumber) {
    data.licenseNumber = instructorForm.value.licenseNumber;
  }

  if (isEditingInstructor.value) {
    const result = await instructorsStore.updateInstructor(instructorUserId.value, data);
    if (result) {
      toast.add({
        title: "Instructor Updated",
        description: "Instructor details saved.",
        color: "success",
      });
    } else {
      toast.add({
        title: "Update Failed",
        description: "Could not update instructor.",
        color: "error",
      });
    }
  } else {
    const createData = {
      ...data,
      username:
        instructorForm.value.firstName.toLowerCase().replace(/\s/g, "") + Date.now(),
      password: "TempPass123!",
      email: `${instructorForm.value.firstName.toLowerCase()}@example.com`,
    };
    const result = await instructorsStore.createInstructor(
      createData,
      instructorFile.value || undefined
    );
    if (result) {
      toast.add({
        title: "Instructor Added",
        description: "New instructor created.",
        color: "success",
      });
    } else {
      toast.add({
        title: "Create Failed",
        description: "Could not create instructor.",
        color: "error",
      });
    }
  }
  isInstructorModalOpen.value = false;
}

async function deleteInstructor() {
  const success = await instructorsStore.deleteInstructor(instructorUserId.value);
  if (success) {
    toast.add({
      title: "Instructor Deleted",
      description: "Instructor has been removed.",
      color: "error",
    });
    isInstructorModalOpen.value = false;
  }
}

// ==================== IMAGE UPLOAD ====================
const MAX_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB
const instructorImageRef = ref<HTMLInputElement | null>(null);

function triggerInstructorImagePicker() {
  instructorImageRef.value?.click();
}

function clearInstructorImage() {
  instructorPhotoUrl.value = "";
  instructorFile.value = null;
  if (instructorImageRef.value) {
    instructorImageRef.value.value = "";
  }
}

function handleInstructorImageUpload(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];

  if (!file) return;

  if (!["image/jpeg", "image/png", "image/webp"].includes(file.type)) {
    toast.add({
      title: "Invalid File",
      description: "Please upload JPG, PNG, or WebP.",
      color: "error",
    });
    input.value = "";
    return;
  }

  if (file.size > MAX_IMAGE_SIZE) {
    toast.add({
      title: "File Too Large",
      description: "Image exceeds 5MB limit.",
      color: "error",
    });
    input.value = "";
    return;
  }

  instructorFile.value = file;

  const reader = new FileReader();
  reader.onload = (e) => {
    instructorPhotoUrl.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
  input.value = "";
}

async function saveSettings() {
  try {
    const result = await settingsStore.updateGeneralSettings({
      businessName: generalSettings.businessName,
      email: generalSettings.email,
      phone: generalSettings.phone,
      fax: generalSettings.fax,
      whatsApp: generalSettings.whatsApp,
      address: generalSettings.address,
      hoursMonFri: generalSettings.hoursMonFri,
      hoursSatSun: generalSettings.hoursSatSun,
      hoursNightShift: generalSettings.hoursNightShift,
      promoEndDate: generalSettings.promoEndDate,
      notifyEmail: notificationSettings.emailNotifications,
      notifySms: notificationSettings.whatsappNotifications,
      notifyWhatsApp: notificationSettings.whatsappNotifications,
    });

    if (result) {
      toast.add({
        title: t("common.success"),
        description: "Settings saved successfully.",
        icon: "i-lucide-check-circle",
        color: "success",
      });
    } else {
      toast.add({
        title: t("common.error"),
        description: "Failed to save settings.",
        color: "error",
      });
    }
  } catch {
    toast.add({
      title: t("common.error"),
      description: "Failed to save settings.",
      color: "error",
    });
  }
}

async function fetchSettings() {
  const settings = await settingsStore.fetchGeneralSettings();
  if (settings) {
    generalSettings.businessName = settings.businessName;
    generalSettings.email = settings.email;
    generalSettings.phone = settings.phone;
    generalSettings.fax = settings.fax;
    generalSettings.whatsApp = settings.whatsApp;
    generalSettings.address = settings.address;
    generalSettings.hoursMonFri = settings.hoursMonFri;
    generalSettings.hoursSatSun = settings.hoursSatSun;
    generalSettings.hoursNightShift = settings.hoursNightShift;
    generalSettings.promoEndDate = settings.promoEndDate;
    notificationSettings.emailNotifications = settings.notifyEmail;
    notificationSettings.whatsappNotifications = settings.notifyWhatsApp;
  }
}

onMounted(() => {
  fetchSettings();
  instructorsStore.fetchInstructors();
  vehiclesStore.fetchVehicles();
});
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar :title="t('admin.settings')">
        <template #right>
          <UButton
            :label="t('admin.saveChanges')"
            color="warning"
            icon="i-lucide-save"
            @click="saveSettings"
          />
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6 max-w-full">
        <!-- General Settings -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-settings" class="size-5 text-warning" />
              <h2 class="font-semibold">{{ t("admin.generalSettings") }}</h2>
            </div>
          </template>

          <div class="grid md:grid-cols-2 gap-4">
            <UFormField :label="t('admin.businessName')">
              <UInput
                v-model="generalSettings.businessName"
                icon="i-lucide-building"
                class="w-full"
                color="warning"
              />
            </UFormField>
            <UFormField label="Email">
              <UInput
                v-model="generalSettings.email"
                type="email"
                icon="i-lucide-mail"
                class="w-full"
                color="warning"
              />
            </UFormField>
            <UFormField :label="t('profile.phone')">
              <UInput
                v-model="generalSettings.phone"
                icon="i-lucide-phone"
                class="w-full"
                color="warning"
              />
            </UFormField>
            <UFormField label="Fax">
              <UInput
                v-model="generalSettings.fax"
                icon="i-lucide-printer"
                class="w-full"
                color="warning"
              />
            </UFormField>
            <UFormField :label="t('admin.whatsappNumber')">
              <UInput
                v-model="generalSettings.whatsApp"
                icon="i-simple-icons-whatsapp"
                class="w-full"
                color="warning"
              />
            </UFormField>
            <UFormField :label="t('admin.promoEndDate')">
              <UInput
                :model-value="generalSettings.promoEndDate || ''"
                type="datetime-local"
                class="w-full"
                color="warning"
                @update:model-value="generalSettings.promoEndDate = $event || null"
              />
            </UFormField>
            <UFormField :label="t('profile.address')" class="md:col-span-2">
              <UTextarea
                v-model="generalSettings.address"
                class="w-full"
                color="warning"
              />
            </UFormField>
          </div>
        </UCard>

        <!-- Operating Hours -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-clock" class="size-5 text-warning" />
              <h2 class="font-semibold">{{ t("home.operatingHours") }}</h2>
            </div>
          </template>

          <div class="space-y-4">
            <div class="flex items-center gap-4">
              <span class="w-32 text-md font-medium">{{ t("admin.mondayFriday") }}</span>
              <UInput
                v-model="generalSettings.hoursMonFri"
                type="text"
                placeholder="08:00 - 17:00"
                class="w-full"
                color="warning"
              />
            </div>
            <div class="flex items-center gap-4">
              <span class="w-32 text-sm font-medium">{{
                t("admin.saturdaySunday")
              }}</span>
              <UInput
                v-model="generalSettings.hoursSatSun"
                type="text"
                placeholder="09:00 - 15:00"
                class="w-full"
                color="warning"
              />
            </div>
            <div class="flex items-center gap-4">
              <span class="w-32 text-md font-medium">{{ t("admin.nightShift") }}</span>
              <UInput
                v-model="generalSettings.hoursNightShift"
                type="text"
                placeholder="19:00 - 22:00"
                class="w-full"
                color="warning"
              />
            </div>
          </div>
        </UCard>

        <!-- Vehicles -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-car" class="size-5 text-warning" />
                <h2 class="font-semibold">{{ t("dashboard.vehicle") }}s</h2>
              </div>
              <UButton
                :label="t('admin.addVehicle')"
                icon="i-lucide-plus"
                size="sm"
                variant="outline"
                color="neutral"
                @click="openNewVehicle"
              />
            </div>
          </template>

          <div class="space-y-3">
            <div
              v-for="vehicle in vehicles"
              :key="vehicle.id"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <UAvatar
                  :src="vehicle.imageUrl || undefined"
                  :text="
                    !vehicle.imageUrl
                      ? (vehicle.brand?.[0] ?? '') + (vehicle.model?.[0] ?? '')
                      : undefined
                  "
                  size="md"
                />
                <div>
                  <p class="font-medium">{{ vehicle.brand }} {{ vehicle.model }}</p>
                  <p class="text-sm text-muted">
                    {{ vehicle.licensePlate || "No plate" }} - {{ vehicle.year }}
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <UBadge
                  :label="
                    vehicle.status === 'available'
                      ? t('billing.active')
                      : vehicle.status === 'in_use'
                      ? 'In Use'
                      : vehicle.status === 'maintenance'
                      ? 'Maintenance'
                      : 'Retired'
                  "
                  :color="
                    vehicle.status === 'available'
                      ? 'success'
                      : vehicle.status === 'in_use'
                      ? 'info'
                      : vehicle.status === 'maintenance'
                      ? 'warning'
                      : 'neutral'
                  "
                  variant="subtle"
                />
                <UButton
                  icon="i-lucide-pencil"
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  @click="openEditVehicle(vehicle)"
                />
              </div>
            </div>
            <div v-if="vehicles.length === 0" class="text-center py-8 text-muted">
              <UIcon name="i-lucide-car" class="size-8 mx-auto mb-2 opacity-50" />
              <p>No vehicles added yet</p>
            </div>
          </div>
        </UCard>

        <!-- Instructors -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-users" class="size-5 text-warning" />
                <h2 class="font-semibold">
                  {{ t("admin.students").replace("Murid", "Instruktur") }}
                </h2>
              </div>
              <UButton
                :label="t('admin.addInstructor')"
                icon="i-lucide-user-plus"
                size="sm"
                variant="outline"
                color="neutral"
                @click="openNewInstructor"
              />
            </div>
          </template>

          <div class="space-y-3">
            <div
              v-for="instructor in instructors"
              :key="instructor.userId"
              class="flex items-center justify-between p-4 rounded-lg border border-default"
            >
              <div class="flex items-center gap-4">
                <UAvatar
                  :src="instructor.image || undefined"
                  :text="
                    !instructor.image
                      ? instructor.name
                          .split(' ')
                          .map((n: string) => n[0])
                          .join('')
                      : undefined
                  "
                  size="sm"
                />
                <div>
                  <p class="font-medium">{{ instructor.name }}</p>
                  <p class="text-md text-muted">{{ instructor.phone }}</p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <UBadge
                  :label="
                    instructor.status === 'active' ? t('billing.active') : 'Inactive'
                  "
                  :color="instructor.status === 'active' ? 'info' : 'neutral'"
                  variant="subtle"
                />
                <UButton
                  icon="i-lucide-pencil"
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  @click="openEditInstructor(instructor)"
                />
              </div>
            </div>
          </div>
        </UCard>

        <!-- Notification Settings -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-bell" class="size-5 text-warning" />
              <h2 class="font-semibold">{{ t("profile.notificationPrefs") }}</h2>
            </div>
          </template>

          <div class="space-y-4">
            <USwitch
              v-model="notificationSettings.emailNotifications"
              :label="t('profile.notifEmail').replace('mendatang', 'kepada murid')"
            />
            <USwitch
              v-model="notificationSettings.whatsappNotifications"
              :label="
                t('profile.notifWa').replace('(24 jam sebelum sesi)', 'kepada murid')
              "
            />
            <USwitch
              v-model="notificationSettings.adminAlerts"
              label="Send admin alerts for general events"
            />
            <USwitch
              v-model="notificationSettings.newUserRegistration"
              label="Notify when new user registers"
            />
            <USwitch
              v-model="notificationSettings.newPackagePurchase"
              label="Notify when user buys a package / becomes a member"
            />

            <UFormField label="Send reminders before session (hours)">
              <UInput
                v-model="notificationSettings.reminderHours"
                type="number"
                min="1"
                max="72"
                class="w-full"
              />
            </UFormField>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>

  <!-- Vehicle Modal -->
  <ClientOnly>
    <UModal v-model:open="isVehicleModalOpen">
      <template #content>
        <div class="bg-default rounded-2xl w-full">
          <div
            class="px-6 py-4 border-b border-default flex items-center justify-between"
          >
            <h3 class="text-base font-semibold">
              {{
                isEditingVehicle
                  ? t("common.edit") + " " + t("dashboard.vehicle")
                  : t("admin.addNew") + " " + t("dashboard.vehicle")
              }}
            </h3>
            <UButton
              icon="i-lucide-x"
              color="neutral"
              variant="ghost"
              @click="isVehicleModalOpen = false"
            />
          </div>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium mb-1.5">{{
                t("admin.vehiclePhoto")
              }}</label>
              <input
                ref="vehicleImageRef"
                type="file"
                accept="image/jpeg,image/png,image/webp"
                class="hidden"
                @change="handleVehicleImageUpload"
              />
              <div
                class="border-2 border-dashed border-default rounded-lg p-4 text-center cursor-pointer hover:border-primary transition-colors"
                @click="triggerVehicleImagePicker"
              >
                <div v-if="vehicleForm.imageUrl" class="relative w-full h-32">
                  <img
                    :src="vehicleForm.imageUrl"
                    class="w-full h-full object-cover rounded-md"
                  />
                  <UButton
                    icon="i-lucide-trash-2"
                    color="error"
                    size="xs"
                    class="absolute top-2 right-2"
                    @click.stop="clearVehicleImage"
                  />
                </div>
                <div v-else class="flex flex-col items-center gap-2 py-4">
                  <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
                  <span class="text-sm text-muted">{{ t("admin.clickToUpload") }}</span>
                </div>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-1.5">Brand</label>
                <UInput
                  v-model="vehicleForm.brand"
                  placeholder="e.g. Toyota"
                  class="w-full"
                />
              </div>
              <div>
                <label class="block text-sm font-medium mb-1.5">Model</label>
                <UInput
                  v-model="vehicleForm.model"
                  placeholder="e.g. Yaris"
                  class="w-full"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-1.5">Year</label>
                <UInput
                  v-model="vehicleForm.year"
                  type="number"
                  placeholder="e.g. 2024"
                  class="w-full"
                />
              </div>
              <div>
                <label class="block text-sm font-medium mb-1.5">{{
                  t("admin.licensePlate")
                }}</label>
                <UInput
                  v-model="vehicleForm.licensePlate"
                  placeholder="e.g. B 1234 EV"
                  class="w-full"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-1.5">Color</label>
                <UInput
                  v-model="vehicleForm.color"
                  placeholder="e.g. Black"
                  class="w-full"
                />
              </div>
              <div>
                <label class="block text-sm font-medium mb-1.5">Transmission</label>
                <USelect
                  v-model="vehicleForm.transmission"
                  :items="[
                    { label: 'Automatic', value: 'automatic' },
                    { label: 'Manual', value: 'manual' },
                  ]"
                  class="w-full"
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">{{
                t("admin.status")
              }}</label>
              <USelect
                v-model="vehicleForm.status"
                :items="[
                  { label: 'Available', value: 'available' },
                  { label: 'In Use', value: 'in_use' },
                  { label: 'Maintenance', value: 'maintenance' },
                  { label: 'Unavailable', value: 'unavailable' },
                ]"
                class="w-full"
              />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Notes</label>
              <UTextarea
                v-model="vehicleForm.notes"
                placeholder="Additional notes about this vehicle..."
                :rows="2"
                class="w-full"
              />
            </div>
          </div>
          <div
            class="px-6 py-4 border-t border-default flex justify-between items-center"
          >
            <div>
              <UButton
                v-if="isEditingVehicle"
                :label="t('common.delete')"
                icon="i-lucide-trash"
                color="error"
                variant="ghost"
                @click="deleteVehicle"
              />
            </div>
            <div class="flex gap-3">
              <UButton
                :label="t('common.cancel')"
                color="neutral"
                variant="outline"
                @click="isVehicleModalOpen = false"
              />
              <UButton
                :label="
                  isEditingVehicle
                    ? t('admin.saveChanges')
                    : t('admin.addNew').replace('Tambah Baru', 'Tambah Kendaraan')
                "
                icon="i-lucide-check"
                @click="saveVehicle"
              />
            </div>
          </div>
        </div>
      </template>
    </UModal>
  </ClientOnly>

  <!-- Instructor Modal -->
  <ClientOnly>
    <UModal v-model:open="isInstructorModalOpen">
      <template #content>
        <div class="bg-default rounded-2xl w-full">
          <div
            class="px-6 py-4 border-b border-default flex items-center justify-between"
          >
            <h3 class="text-base font-semibold">
              {{
                isEditingInstructor
                  ? t("common.edit") + " Instruktur"
                  : t("admin.addNew").replace("Tambah Baru", "Tambah Instruktur")
              }}
            </h3>
            <UButton
              icon="i-lucide-x"
              color="neutral"
              variant="ghost"
              @click="isInstructorModalOpen = false"
            />
          </div>
          <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
            <div class="grid grid-cols-2 gap-4">
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("register.form.firstName")
                }}</label>
                <UInput
                  v-model="instructorForm.firstName"
                  placeholder="e.g. Budi"
                  class="w-full"
                />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("register.form.lastName")
                }}</label>
                <UInput
                  v-model="instructorForm.lastName"
                  placeholder="e.g. Santoso"
                  class="w-full"
                />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("profile.phone")
                }}</label>
                <UInput
                  v-model="instructorForm.phoneNumber"
                  placeholder="e.g. 08123456789"
                  class="w-full"
                />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("instructors.yearsExperience")
                }}</label>
                <UInput
                  v-model="instructorForm.yearsOfExperience"
                  type="number"
                  min="0"
                  class="w-full"
                />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("admin.bnspNo")
                }}</label>
                <UInput
                  v-model="instructorForm.bnspCertificateNumber"
                  placeholder="e.g. BNSP-101-2023"
                  class="w-full"
                />
              </div>
              <div class="col-span-2 sm:col-span-1">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("admin.licenseNo")
                }}</label>
                <UInput
                  v-model="instructorForm.licenseNumber"
                  placeholder="e.g. SIM A / 12345678"
                  class="w-full"
                />
              </div>
              <div class="col-span-2">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("admin.instructorPhoto")
                }}</label>
                <input
                  ref="instructorImageRef"
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  class="hidden"
                  @change="handleInstructorImageUpload"
                />
                <div
                  class="border-2 border-dashed border-default rounded-lg p-4 text-center cursor-pointer hover:border-primary transition-colors"
                  @click="instructorImageRef?.click()"
                >
                  <div v-if="instructorPhotoUrl" class="relative w-full h-32">
                    <img
                      :src="instructorPhotoUrl"
                      class="w-full h-full object-contain rounded-md"
                    />
                    <UButton
                      icon="i-lucide-trash-2"
                      color="error"
                      size="xs"
                      class="absolute top-2 right-2"
                      @click.stop="
                        instructorPhotoUrl = '';
                        instructorFile = null;
                      "
                    />
                  </div>
                  <div v-else class="flex flex-col items-center gap-2 py-4">
                    <UIcon name="i-lucide-image-plus" class="size-6 text-muted" />
                    <span class="text-sm text-muted">{{ t("admin.clickToUpload") }}</span>
                  </div>
                </div>
              </div>
              <div class="col-span-2">
                <label class="block text-sm font-medium mb-1.5">{{
                  t("admin.bioDescription")
                }}</label>
                <UTextarea
                  v-model="instructorForm.description"
                  placeholder="Short description of the instructor..."
                  :rows="3"
                  class="w-full"
                />
              </div>
            </div>
          </div>
          <div
            class="px-6 py-4 border-t border-default flex justify-between items-center"
          >
            <div>
              <UButton
                v-if="isEditingInstructor"
                :label="t('common.delete')"
                icon="i-lucide-trash"
                color="error"
                variant="ghost"
                @click="deleteInstructor"
              />
            </div>
            <div class="flex gap-3">
              <UButton
                :label="t('common.cancel')"
                color="neutral"
                variant="outline"
                @click="isInstructorModalOpen = false"
              />
              <UButton
                :label="
                  isEditingInstructor
                    ? t('admin.saveChanges')
                    : t('admin.addNew').replace('Tambah Baru', 'Tambah Instruktur')
                "
                icon="i-lucide-check"
                @click="saveInstructor"
              />
            </div>
          </div>
        </div>
      </template>
    </UModal>
  </ClientOnly>
</template>
