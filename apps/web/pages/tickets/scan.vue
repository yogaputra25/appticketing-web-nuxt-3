<template>
  <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
    <h1 class="text-2xl md:text-3xl font-bold text-gray-900 mb-2">Scan QR Tiket</h1>
    <p class="text-gray-500 text-sm mb-6">Arahkan kamera ke QR code tiket untuk verifikasi</p>

    <div class="card overflow-hidden relative">
      <div v-if="!cameraReady && !cameraError" class="absolute inset-0 z-10 bg-white flex items-center justify-center">
        <div class="text-center p-8">
          <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p class="text-gray-500">Mengakses kamera...</p>
        </div>
      </div>

      <ClientOnly>
        <QrcodeStream
          :paused="scanned"
          :formats="['qr_code']"
          :constraints="{ facingMode: 'environment', width: { min: 360, ideal: 640 }, height: { min: 360, ideal: 640 } }"
          @detect="onDetect"
          @camera-on="cameraReady = true"
          @error="onError"
          class="w-full"
        />
        <template #fallback>
          <div class="p-8 text-center text-gray-500">
            <p>Camera requires browser with JavaScript enabled</p>
          </div>
        </template>
      </ClientOnly>

      <div v-if="scanned" class="absolute inset-0 bg-black/60 flex items-center justify-center z-20">
        <div class="text-center text-white p-4">
          <p class="text-lg font-semibold">QR Code Terbaca</p>
        </div>
      </div>
    </div>

    <div v-if="cameraError" class="card p-5 mt-4 text-center">
      <svg class="w-12 h-12 text-red-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4.5c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
      </svg>
      <p class="text-red-600 font-medium">Gagal mengakses kamera</p>
      <p class="text-sm text-gray-500 mt-1">{{ cameraError }}</p>
      <button class="btn-outline mt-4 text-sm" @click="resetCamera">Coba Lagi</button>
    </div>

    <Transition name="fade">
      <div v-if="result" class="card p-5 mt-4" :class="result.valid ? 'border-green-200 bg-green-50' : 'border-red-200 bg-red-50'">
        <div class="flex items-start gap-4">
          <div v-if="result.valid" class="shrink-0">
            <svg class="w-10 h-10 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div v-else class="shrink-0">
            <svg class="w-10 h-10 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="min-w-0">
            <p class="font-semibold text-lg" :class="result.valid ? 'text-green-800' : 'text-red-800'">
              {{ result.valid ? 'Tiket Valid' : 'Tiket Tidak Valid' }}
            </p>
            <p class="text-sm mt-1" :class="result.valid ? 'text-green-700' : 'text-red-700'">
              {{ result.message }}
            </p>
            <div v-if="result.ticket" class="mt-3 space-y-1 text-sm" :class="result.valid ? 'text-green-800' : 'text-red-800'">
              <p><span class="font-medium">Event:</span> {{ result.ticket.event_title }}</p>
              <p><span class="font-medium">Kategori:</span> {{ result.ticket.category_name }}</p>
              <p><span class="font-medium">Kode:</span> {{ result.ticket.ticket_code }}</p>
              <p v-if="!result.valid && result.ticket.scanned_at">
                <span class="font-medium">Digunakan pada:</span> {{ formatDate(result.ticket.scanned_at) }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <div v-if="scanned" class="mt-4 text-center">
      <button class="btn-primary text-sm" @click="resetScan">Scan Lagi</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { QrcodeStream } from 'vue-qrcode-reader'

definePageMeta({
  middleware: 'auth',
})

const { getVerifyUrl } = useQrCode()
const api = useApi()

const cameraReady = ref(false)
const cameraError = ref('')
const scanned = ref(false)
const result = ref<{ valid: boolean; message: string; ticket?: any } | null>(null)
const verifying = ref(false)

function extractTicketCodeFromUrl(url: string): string | null {
  try {
    const parsed = new URL(url)
    const segments = parsed.pathname.split('/').filter(Boolean)
    const last = segments[segments.length - 1]
    return last || null
  } catch {
    if (url.startsWith('TCK-') || url.length > 10) return url
    return null
  }
}

async function onDetect(detectedCodes: any[]) {
  if (scanned.value || verifying.value) return

  const code = detectedCodes?.[0]
  if (!code?.rawValue) return

  scanned.value = true
  verifying.value = true

  const ticketCode = extractTicketCodeFromUrl(code.rawValue)
  if (!ticketCode) {
    result.value = {
      valid: false,
      message: 'QR code tidak mengandung data tiket yang valid',
    }
    verifying.value = false
    return
  }

  try {
    const data = await api.post<any>(`/api/tickets/verify/${ticketCode}`)
    result.value = {
      valid: true,
      message: 'Tiket berhasil diverifikasi',
      ticket: data.ticket,
    }
  } catch (err: any) {
    result.value = {
      valid: false,
      message: err?.message || 'Gagal memverifikasi tiket',
      ticket: err?.ticket,
    }
  } finally {
    verifying.value = false
  }
}

function onError(err: any) {
  cameraError.value = err?.message || 'Gagal mengakses kamera'
}

function resetScan() {
  scanned.value = false
  result.value = null
}

function resetCamera() {
  cameraReady.value = false
  cameraError.value = ''
  scanned.value = false
  result.value = null
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
