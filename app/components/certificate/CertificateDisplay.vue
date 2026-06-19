<script setup lang="ts">
import { certificateService } from "~/services/certificateService";
import { useToast } from "@nuxt/ui/runtime/composables/index.js";

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
      title: "Error",
      description: "Member ID not found for this certificate.",
      color: "error",
    });
    return;
  }

  isDownloading.value = true;

  try {
    await certificateService.downloadCertificate(memberId, certId);

    toast.add({
      title: "Download Started",
      description: "Your certificate is being downloaded.",
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
      title: "Link Copied",
      description: "Certificate verification link copied to clipboard.",
      icon: "i-lucide-copy",
      color: "success",
    });
  }
}

async function shareCertificate(certId: string) {
  if (import.meta.client) {
    const verificationUrl = `${window.location.origin}/verify-certificate?id=${certId}`;
    const shareData = {
      title: "My Driving Certificate",
      text: "I completed my driving course at Drive Master! Check out my certificate.",
      url: verificationUrl,
    };

    if (navigator.share) {
      try {
        await navigator.share(shareData);
      } catch (err: any) {
        if (err.name !== "AbortError") {
          console.error("Share failed:", err);
          toast.add({
            title: "Share Failed",
            description: "An error occurred while trying to share.",
            color: "error",
          });
        }
      }
    } else {
      toast.add({
        title: "Share Not Supported",
        description: "Verification link copied to clipboard.",
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
        <h2 class="font-semibold">Certificate Preview</h2>
      </template>

      <div
        class="aspect-[1.414/1] bg-gradient-to-br from-warning/5 to-warning/10 rounded-lg border-2 border-dashed border-warning/30 flex flex-col items-center justify-center p-8 text-center"
      >
        <div class="flex items-center gap-2 mb-4">
          <img
            src="/drive-master-logo2.png"
            alt="Drive Master Logo"
            class="h-16"
          />
        </div>

        <p class="text-md text-muted mb-2">This is to certify that</p>
        <p class="text-2xl font-bold mb-2">{{ certificate.recipientName }}</p>
        <p class="text-md text-muted mb-4">has successfully completed the</p>
        <p class="text-lg font-semibold text-warning mb-4">
          {{ certificate.title }}
        </p>

        <div class="mt-4 pt-4 border-t border-dashed border-muted w-full">
          <div class="flex justify-between text-md text-muted">
            <span>Certificate ID: {{ certificate.id }}</span>
            <span>Issued: {{ certificate.issuedDate }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <UButton
          label="Download Certificate (PDF)"
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
        <h2 class="font-semibold">Certificate Details</h2>
      </template>

      <div class="space-y-4">
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">Certificate ID</span>
          <span class="font-medium font-mono">{{ certificate.id }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">Title</span>
          <span class="font-medium">{{ certificate.title }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">Recipient</span>
          <span class="font-medium">{{ certificate.recipientName }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-default">
          <span class="text-muted">Issue Date</span>
          <span class="font-medium">{{ certificate.issuedDate }}</span>
        </div>
        <div class="flex justify-between py-2">
          <span class="text-muted">Status</span>
          <UBadge label="Valid" color="success" />
        </div>
      </div>
    </UCard>

    <!-- Verification Card -->
    <UCard>
      <template #header>
        <div class="flex items-center gap-2">
          <UIcon name="i-lucide-shield-check" class="size-5 text-warning" />
          <h2 class="font-semibold">Verification</h2>
        </div>
      </template>

      <p class="text-md text-muted mb-4">
        This certificate can be verified using the certificate ID. Share your
        achievement on social media or with potential employers.
      </p>

      <div class="flex gap-2">
        <UButton
          label="Copy Link"
          icon="i-lucide-copy"
          variant="outline"
          color="neutral"
          size="md"
          @click="copyVerificationLink(certificate.id)"
        />
        <UButton
          label="Share"
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
