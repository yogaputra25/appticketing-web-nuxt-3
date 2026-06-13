<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-32 md:pb-8">
    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="text-center py-16">
      <p class="text-gray-500 text-lg">{{ error }}</p>
      <NuxtLink to="/my/bookings" class="text-primary-600 hover:text-primary-700 mt-2 inline-block">
        &larr; Kembali
      </NuxtLink>
    </div>

    <template v-else-if="payment">
      <NuxtLink to="/my/bookings" class="text-sm text-gray-500 hover:text-primary-600 mb-4 inline-block">
        &larr; Kembali
      </NuxtLink>

      <div class="card p-5 md:p-6 mb-6">
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 mb-4">Pembayaran</h1>

        <div class="flex justify-between items-center py-2 border-b border-gray-100">
          <span class="text-sm text-gray-500">Kode Pembayaran</span>
          <span class="text-sm font-mono font-medium text-gray-900">{{ payment.payment_code }}</span>
        </div>
        <div class="flex justify-between items-center py-2 border-b border-gray-100">
          <span class="text-sm text-gray-500">Total Pembayaran</span>
          <span class="text-lg font-bold text-gray-900">{{ formatPrice(payment.amount) }}</span>
        </div>
        <div class="flex justify-between items-center py-2 border-b border-gray-100">
          <span class="text-sm text-gray-500">Status</span>
          <span class="badge" :class="statusBadge(payment.status)">{{ payment.status }}</span>
        </div>

        <div class="bg-yellow-50 border border-yellow-200 rounded-xl p-4 mt-4">
          <div class="flex items-center gap-3">
            <svg class="w-6 h-6 text-yellow-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <p class="text-sm font-medium text-yellow-800">Selesaikan pembayaran dalam</p>
              <p class="text-lg font-bold text-yellow-900 tabular-nums">
                <CountdownTimer :target-date="payment.expired_at" @expired="isExpired = true" />
              </p>
            </div>
          </div>
        </div>

        <div v-if="isExpired" class="text-sm text-red-600 bg-red-50 rounded-lg p-3 mt-4">
          Waktu pembayaran telah habis. Pesanan akan dibatalkan.
        </div>

        <div class="space-y-3 mt-6">
          <h3 class="font-semibold text-gray-900">Instruksi Pembayaran</h3>
          <ol class="space-y-3 text-sm text-gray-600">
            <li class="flex gap-2">
              <span class="font-bold text-primary-600 shrink-0">1.</span>
              <span>Transfer ke rekening <strong class="text-gray-900">Bank Mandiri: 123-00-1234567-8</strong> a.n. PT War Tiket</span>
            </li>
            <li class="flex gap-2">
              <span class="font-bold text-primary-600 shrink-0">2.</span>
              <span>Masukkan jumlah <strong class="text-gray-900">{{ formatPrice(payment.amount) }}</strong></span>
            </li>
            <li class="flex gap-2">
              <span class="font-bold text-primary-600 shrink-0">3.</span>
              <span>Konfirmasi pembayaran dengan mengklik tombol "Bayar Sekarang" setelah transfer</span>
            </li>
          </ol>
        </div>

        <div v-if="simulateError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3 mt-4">
          {{ simulateError }}
        </div>
        <div v-if="simulateSuccess" class="text-sm text-green-600 bg-green-50 rounded-lg p-3 mt-4">
          Pembayaran berhasil! Mengarahkan ke detail tiket...
        </div>
      </div>
    </template>

    <div class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-4 py-3 flex gap-2">
      <button
        class="btn-outline flex-1 touch-target text-sm"
        :disabled="isExpired || simulating"
        @click="handleSimulateFail"
      >
        Bayar Gagal
      </button>
      <button
        class="btn-accent flex-1 touch-target"
        :disabled="isExpired || simulating"
        @click="handleSimulateSuccess"
      >
        {{ simulating ? 'Memproses...' : 'Bayar Sekarang' }}
      </button>
    </div>

    <div v-if="!isExpired && !simulateSuccess" class="hidden md:flex gap-3">
      <button class="btn-outline touch-target" :disabled="simulating" @click="handleSimulateFail">
        Bayar Gagal
      </button>
      <button class="btn-accent px-8 touch-target" :disabled="simulating" @click="handleSimulateSuccess">
        {{ simulating ? 'Memproses...' : 'Bayar Sekarang' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const error = ref('')
const isExpired = ref(false)
const payment = ref<any>(null)
const simulating = ref(false)
const simulateError = ref('')
const simulateSuccess = ref(false)

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price)
}

function statusBadge(status: string) {
  const map: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    success: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800',
    expired: 'bg-gray-100 text-gray-800',
    refunded: 'bg-purple-100 text-purple-800',
  }
  return map[status] || 'bg-gray-100 text-gray-800'
}

async function createPayment() {
  const api = useApi()
  const res = await api.post<any>('/api/payments/create', {
    booking_id: Number(route.params.id),
  })
  payment.value = res
  if (res.status !== 'pending') {
    error.value = 'Pembayaran sudah diproses'
  }
}

async function handleSimulateSuccess() {
  if (!payment.value || isExpired.value) return
  simulating.value = true
  simulateError.value = ''
  try {
    const api = useApi()
    const res = await api.post<any>(`/api/payments/${payment.value.id}/simulate`, { action: 'success' })
    simulateSuccess.value = true
    setTimeout(() => {
      router.push(`/my/bookings/${route.params.id}`)
    }, 1500)
  } catch (err: any) {
    simulateError.value = err?.message || 'Gagal memproses pembayaran'
  } finally {
    simulating.value = false
  }
}

async function handleSimulateFail() {
  if (!payment.value || isExpired.value) return
  simulating.value = true
  simulateError.value = ''
  try {
    const api = useApi()
    await api.post<any>(`/api/payments/${payment.value.id}/simulate`, { action: 'fail' })
    router.push(`/my/bookings/${route.params.id}`)
  } catch (err: any) {
    simulateError.value = err?.message || 'Gagal memproses pembayaran'
  } finally {
    simulating.value = false
  }
}

onMounted(async () => {
  try {
    await createPayment()
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat data pembayaran'
  } finally {
    loading.value = false
  }
})
</script>
