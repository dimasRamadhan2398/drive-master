<script setup lang="ts">
import { useToast } from '@nuxt/ui/runtime/composables/index.js'
import { ref } from 'vue'
import { jsPDF } from 'jspdf'
import html2canvas from 'html2canvas'

definePageMeta({ layout: 'dashboard' })

const toast = useToast()
const certificatePreview = ref<HTMLElement[] | null>(null)
const isDownloading = ref(false)

// Mock certificate data
const certificateStatus = ref({
  isEligible: false,
  progress: 40,
  completedSessions: 4,
  totalSessions: 10,
  remainingSessions: 6
})

// Mock issued certificates
const issuedCertificates = ref([
  {
    id: 'EVDA-2026-001234',
    title: 'Drive Master Sessions Course',
    recipientName: 'John Doe',
    issuedDate: 'Apr 15, 2026',
    status: 'issued',
    downloadUrl: '#'
  }
])

// For demo purposes, set certificate as available
const hasCertificate = ref(true)

async function downloadCertificate(certId: string) {
  if (!certificatePreview.value || certificatePreview.value.length === 0) {
    toast.add({ title: 'Error', description: 'Elemen pratinjau sertifikat tidak ditemukan.', color: 'error' })
    return
  }

  isDownloading.value = true
  // Menargetkan elemen pratinjau pertama dari v-for
  const elementToCapture = certificatePreview.value?.[0]

  if (!elementToCapture) {
    isDownloading.value = false
    return
  }

  try {
    const canvas = await html2canvas(elementToCapture, {
      scale: 2, // Meningkatkan skala untuk kualitas gambar yang lebih baik
      useCORS: true,
      backgroundColor: null // Menggunakan background dari elemen itu sendiri
    })

    const imgData = canvas.toDataURL('image/png')

    // Membuat PDF dengan orientasi landscape (A4)
    const pdf = new jsPDF({ orientation: 'landscape', unit: 'mm', format: 'a4' })
    const pdfWidth = pdf.internal.pageSize.getWidth()
    const pdfHeight = pdf.internal.pageSize.getHeight()

    // Menambahkan gambar ke PDF
    pdf.addImage(imgData, 'PNG', 0, 0, pdfWidth, pdfHeight)
    pdf.save(`certificate-${certId}.pdf`)

    toast.add({ title: 'Unduhan Dimulai', description: 'Sertifikat Anda sedang diunduh.', icon: 'i-lucide-download', color: 'success' })
  } catch (error) {
    console.error('Gagal membuat PDF:', error)
    toast.add({ title: 'Gagal Mengunduh', description: 'Tidak dapat membuat file PDF.', color: 'error' })
  } finally {
    isDownloading.value = false
  }
}

function copyVerificationLink(certId: string) {
  if (import.meta.client) {
    // Membangun URL verifikasi yang dinamis
    const verificationUrl = `${window.location.origin}/verify-certificate?id=${certId}`
    navigator.clipboard.writeText(verificationUrl)
    toast.add({
      title: 'Tautan Disalin',
      description: 'Tautan verifikasi sertifikat telah disalin ke papan klip.',
      icon: 'i-lucide-copy',
      color: 'success'
    })
  }
}

async function shareCertificate(certId: string) {
  if (import.meta.client) {
    const verificationUrl = `${window.location.origin}/verify-certificate?id=${certId}`
    const shareData = {
      title: 'Sertifikat Mengemudi Saya',
      text: 'Saya telah menyelesaikan kursus mengemudi di Drive Master! Lihat sertifikat saya.',
      url: verificationUrl
    }

    if (navigator.share) {
      try {
        await navigator.share(shareData)
      } catch (err: any) {
        // Tidak melakukan apa-apa jika pengguna membatalkan dialog share (AbortError)
        if (err.name !== 'AbortError') {
          console.error('Share failed:', err)
          toast.add({ title: 'Gagal Membagikan', description: 'Terjadi kesalahan saat mencoba membagikan.', color: 'error' })
        }
      }
    } else {
      // Fallback untuk desktop atau browser yang tidak mendukung
      toast.add({ title: 'Fitur Share Tidak Didukung', description: 'Tautan verifikasi disalin ke clipboard.', color: 'info', icon: 'i-lucide-info' })
      copyVerificationLink(certId)
    }
  }
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="My Certificate">
        <template #right>
          <UColorModeButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="p-6 space-y-6">
        <!-- Certificate Status -->
        <UCard v-if="!hasCertificate">
          <template #header>
            <h2 class="font-semibold">Certificate Progress</h2>
          </template>

          <div class="text-center py-8">
            <div class="mx-auto w-20 h-20 rounded-full bg-muted/50 flex items-center justify-center mb-4">
              <UIcon name="i-lucide-award" class="size-10 text-muted" />
            </div>
            <h3 class="text-xl font-bold mb-2">Certificate Not Yet Available</h3>
            <p class="text-muted max-w-md mx-auto">
              Complete all your training sessions to receive your official EV Driving Certificate.
            </p>

            <div class="max-w-md mx-auto mt-8 space-y-4">
              <div class="flex justify-between text-md">
                <span class="text-muted">Progress</span>
                <span class="font-medium">{{ certificateStatus.completedSessions }}/{{ certificateStatus.totalSessions }} sessions</span>
              </div>
              <UProgress :value="certificateStatus.progress" />
              <p class="text-md text-muted">
                {{ certificateStatus.remainingSessions }} more sessions to complete
              </p>
            </div>
          </div>

          <template #footer>
            <div class="flex justify-center">
              <NuxtLink to="/dashboard/schedule">
                <UButton label="Book Next Session" icon="i-lucide-calendar-plus" />
              </NuxtLink>
            </div>
          </template>
        </UCard>

        <!-- Available Certificates -->
        <template v-if="hasCertificate">
          <UAlert icon="i-lucide-award" color="success" title="Congratulations!">
            <template #description>
              You have successfully completed your training and your certificate is ready for download.
            </template>
          </UAlert>

          <div class="grid md:grid-cols-2 gap-6">
            <!-- Certificate Preview -->
            <UCard>
              <template #header>
                <h2 class="font-semibold">Certificate Preview</h2>
              </template>

              <div ref="certificatePreview" v-for="certificate in issuedCertificates" :key="certificate.id" class="aspect-[1.414/1] bg-gradient-to-br from-warning/5 to-warning/10 rounded-lg border-2 border-dashed border-warning/30 flex flex-col items-center justify-center p-8 text-center">
                <div class="flex items-center gap-2 mb-4">
                  <img src="/drive-master-logo2.png" alt="Drive Master Logo" class="h-16" />
                </div>
                
                <p class="text-md text-muted mb-2">This is to certify that</p>
                <p class="text-2xl font-bold mb-2">{{ certificate.recipientName }}</p>
                <p class="text-md text-muted mb-4">has successfully completed the</p>
                <p class="text-lg font-semibold text-warning mb-4">{{ certificate.title }}</p>
                
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
                  @click="issuedCertificates[0] && downloadCertificate(issuedCertificates[0].id)"
                />
              </template>
            </UCard>

            <!-- Certificate Details -->
            <div class="space-y-6">
              <UCard>
                <template #header>
                  <h2 class="font-semibold">Certificate Details</h2>
                </template>

                <div v-for="certificate in issuedCertificates" :key="certificate.id" class="space-y-4">
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

              <UCard>
                <template #header>
                  <div class="flex items-center gap-2">
                    <UIcon name="i-lucide-shield-check" class="size-5 text-warning" />
                    <h2 class="font-semibold">Verification</h2>
                  </div>
                </template>

                <p class="text-md text-muted mb-4">
                  This certificate can be verified using the certificate ID. Share your achievement on social media or with potential employers.
                </p>

                <div class="flex gap-2">
                  <UButton label="Copy Link" icon="i-lucide-copy" variant="outline" color="neutral" size="md" @click="issuedCertificates[0] && copyVerificationLink(issuedCertificates[0].id)" />
                  <UButton label="Share" icon="i-lucide-share-2" variant="outline" color="neutral" size="md" @click="issuedCertificates[0] && shareCertificate(issuedCertificates[0].id)" />
                </div>
              </UCard>
            </div>
          </div>
        </template>

        <!-- Certificate Info -->
        <UCard>
          <template #header>
            <h2 class="font-semibold">About Drive Master Certificates</h2>
          </template>

          <div class="grid md:grid-cols-3 gap-6">
            <div class="text-center">
              <div class="mx-auto w-14 h-14 rounded-full bg-warning/10 flex items-center justify-center mb-3">
                <UIcon name="i-lucide-file-badge" class="size-7 text-warning" />
              </div>
              <h3 class="font-semibold mb-1">Official Recognition</h3>
              <p class="text-md text-muted">Our certificates are recognized by industry partners and employers.</p>
            </div>
            <div class="text-center">
              <div class="mx-auto w-14 h-14 rounded-full bg-warning/10 flex items-center justify-center mb-3">
                <UIcon name="i-lucide-qr-code" class="size-7 text-warning" />
              </div>
              <h3 class="font-semibold mb-1">Digital Verification</h3>
              <p class="text-md text-muted">Each certificate includes a unique ID for easy online verification.</p>
            </div>
            <div class="text-center">
              <div class="mx-auto w-14 h-14 rounded-full bg-warning/10 flex items-center justify-center mb-3">
                <UIcon name="i-lucide-infinity" class="size-7 text-warning" />
              </div>
              <h3 class="font-semibold mb-1">Lifetime Validity</h3>
              <p class="text-md text-muted">Your certificate never expires and remains valid indefinitely.</p>
            </div>
          </div>
        </UCard>
      </div>
    </template>
  </UDashboardPanel>
</template>
