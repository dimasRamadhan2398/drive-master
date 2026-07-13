<script setup lang="ts">
import { certificateService } from "~/services/certificateService";
import { useToast } from "@nuxt/ui/runtime/composables/index.js";

const { t } = useI18n()

const props = defineProps<{
  certificate: {
    id: string;
    title: string;
    recipientName: string;
    issuedDate: string;
    status: string;
    downloadUrl: string;
    memberId?: string;
  };
}>();

const toast = useToast();
const isDownloading = ref(false);

async function downloadCertificate(certId: string) {
  const memberId = props.certificate.memberId;

  if (!memberId) {
    toast.add({
      title: t('common.error'),
      description: "Member ID not found for this certificate.",
      color: "error",
    });
    return;
  }

  isDownloading.value = true;

  try {
    await certificateService.downloadCertificate(memberId, certId);

    toast.add({
      title: t('certificate.downloadStarted'),
      description: t('certificate.downloadStartedDesc'),
      icon: "i-lucide-download",
      color: "success",
    });
  } catch (error) {
    console.error("Failed to download certificate:", error);
    toast.add({
      title: "Download Failed",
      description: "Could not download certificate from server.",
      color: "error",
    });
  } finally {
    isDownloading.value = false;
  }
}

function copyVerificationLink(certId: string) {
  if (import.meta.client) {
    const verificationUrl = `${window.location.origin}/verify-certificate?id=${certId}`;
    navigator.clipboard.writeText(verificationUrl);
    toast.add({
      title: t('certificate.linkCopied'),
      description: t('certificate.linkCopiedDesc'),
      icon: "i-lucide-copy",
      color: "success",
    });
  }
}

async function shareCertificate(certId: string) {
  if (import.meta.client) {
    const verificationUrl = `${window.location.origin}/verify-certificate?id=${certId}`;
    const shareData = {
      title: t('certificate.myDrivingCert'),
      text: t('certificate.shareText'),
      url: verificationUrl,
    };

    if (navigator.share) {
      try {
        await navigator.share(shareData);
      } catch (err: any) {
        if (err.name !== "AbortError") {
          console.error("Share failed:", err);
          toast.add({
            title: t('certificate.shareFailed'),
            description: t('certificate.shareFailedDesc'),
            color: "error",
          });
        }
      }
    } else {
      toast.add({
        title: t('certificate.shareNotSupported'),
        description: t('certificate.shareNotSupportedDesc'),
        color: "info",
        icon: "i-lucide-info",
      });
      copyVerificationLink(certId);
    }
  }
}

defineExpose({
  downloadCertificate,
});
</script>

<template>
  <div>
    <!-- Certificate Preview Card -->
    <UCard>
      <template #header>
        <h2 class="font-semibold">{{ t('certificate.preview') }}</h2>
      </template>

      <div
        class="aspect-[1.414/1] bg-gradient-to-br from-warning/5 to-warning/10 rounded-lg border-2 border-dashed border-warning/30 flex flex-col items-center justify-center p-8 text-center"
      >
        <div class="flex items-center gap-2 mb-4">
          <img
            src="/drive-master-logo-light.png"
            alt="Drive Master Logo"
            class="h-16 dark:hidden"
          />
          <img
            src="/drive-master-logo-dark.jpg"
            alt="Drive Master Logo"
            class="h-16 hidden dark:block"
          />
        </div>

        <p class="text-md text-muted mb-2">{{ t('certificate.certifyThat') }}</p>
        <p class="text-2xl font-bold mb-2">{{ certificate.recipientName }}</p>
        <p class="text-md text-muted mb-4">{{ t('certificate.successfullyCompleted') }}</p>
        <p class="text-lg font-semibold text-warning mb-4">
          {{ certificate.title }}
        </p>

        <div class="mt-4 pt-4 border-t border-dashed border-muted w-full">
          <div class="flex justify-between text-md text-muted">
            <span>{{ t('admin.certId') }}: {{ certificate.id }}</span>
            <span>{{ t('admin.issueDate') }}: {{ certificate.issuedDate }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <UButton
          :label="t('certificate.downloadPdf')"
          icon="i-lucide-download"
          block
          :loading="isDownloading"
          @click="downloadCertificate(certificate.id)"
        />
      </template>
    </UCard>

    <!-- Certificate Details Card -->
    <UCard>
      <template #header>
        <h2 class="font-semibold">{{ t('certificate.details') }}</h2>
      </template>

      <div class="space-y-4">
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">{{ t('admin.certId') }}</span>
          <span class="font-medium font-mono">{{ certificate.id }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">{{ t('admin.content.title') }}</span>
          <span class="font-medium">{{ certificate.title }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">{{ t('certificate.recipient') }}</span>
          <span class="font-medium">{{ certificate.recipientName }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">{{ t('admin.issueDate') }}</span>
          <span class="font-medium">{{ certificate.issuedDate }}</span>
        </div>
        <div class="flex justify-between py-2">
          <span class="text-muted">{{ t('billing.status') }}</span>
          <UBadge :label="t('certificate.valid')" color="success" />
        </div>
      </div>
    </UCard>

    <!-- Verification Card -->
    <UCard>
      <template #header>
        <div class="flex items-center gap-2">
          <UIcon name="i-lucide-shield-check" class="size-5 text-warning" />
          <h2 class="font-semibold">{{ t('certificate.verification') }}</h2>
        </div>
      </template>

      <p class="text-md text-muted mb-4">
        {{ t('certificate.verificationDesc') }}
      </p>

      <div class="flex gap-2">
        <UButton
          :label="t('certificate.copyLink')"
          icon="i-lucide-copy"
          variant="outline"
          color="neutral"
          size="md"
          @click="copyVerificationLink(certificate.id)"
        />
        <UButton
          :label="t('certificate.share')"
          icon="i-lucide-share-2"
          variant="outline"
          color="neutral"
          size="md"
          @click="shareCertificate(certificate.id)"
        />
      </div>
    </UCard>
  </div>
</template>
