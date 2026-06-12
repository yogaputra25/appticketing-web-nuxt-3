<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-28 md:pb-8">
    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="text-center py-16">
      <p class="text-gray-500 text-lg">{{ error }}</p>
      <NuxtLink to="/my/bookings" class="text-primary-600 hover:text-primary-700 mt-2 inline-block">
        &larr; Kembali
      </NuxtLink>
    </div>

    <template v-else>
      <NuxtLink to="/my/bookings" class="text-sm text-gray-500 hover:text-primary-600 mb-4 inline-block">
        &larr; Kembali
      </NuxtLink>

      <div class="card p-5 md:p-6 mb-6">
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 mb-4">Pembayaran</h1>

        <div class="bg-yellow-50 border border-yellow-200 rounded-xl p-4 mb-6">
          <div class="flex items-center gap-3">
            <svg class="w-6 h-6 text-yellow-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <p class="text-sm font-medium text-yellow-800">Selesaikan pembayaran dalam</p>
              <p class="text-lg font-bold text-yellow-900 tabular-nums">
                <CountdownTimer :target-date="expiresAt" @expired="isExpired = true" />
              </p>
            </div>
          </div>
        </div>

        <div v-if="isExpired" class="text-sm text-red-600 bg-red-50 rounded-lg p-3 mb-4">
          Waktu pembayaran telah habis. Pesanan akan dibatalkan.
        </div>

        <div class="space-y-3">
          <h3 class="font-semibold text-gray-900">Instruksi Pembayaran</h3>
          <ol class="space-y-3 text-sm text-gray-600">
            <li class="flex gap-2">
              <span class="font-bold text-primary-600 shrink-0">1.</span>
              <span>Transfer ke rekening <strong class="text-gray-900">Bank Mandiri: 123-00-1234567-8</strong> a.n. PT War Tiket</span>
            </li>
            <li class="flex gap-2">
              <span class="font-bold text-primary-600 shrink-0">2.</span>
              <span>Masukkan jumlah sesuai <strong class="text-gray-900">total pembayaran</strong></span>
            </li>
            <li class="flex gap-2">
              <span class="font-bold text-primary-600 shrink-0">3.</span>
              <span>Konfirmasi pembayaran dengan mengklik tombol "Bayar Sekarang" setelah transfer</span>
            </li>
          </ol>
        </div>
      </div>
    </template>

    <div class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-4 py-3">
      <button
        class="btn-accent w-full touch-target"
        :disabled="isExpired"
        @click="handlePay"
      >
        <CountdownTimer v-if="!isExpired" :target-date="expiresAt" class="mr-2" />
        Bayar Sekarang
      </button>
    </div>

    <div v-if="!isExpired" class="hidden md:block">
      <button class="btn-accent px-8 touch-target" @click="handlePay">
        Bayar Sekarang
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const loading = ref(true)
const error = ref('')
const isExpired = ref(false)
const expiresAt = ref(new Date(Date.now() + 5 * 60000).toISOString())

function handlePay() {
  // placeholder
}

onMounted(async () => {
  try {
    const api = useApi()
    const data = await api.get<any>(`/bookings/${route.params.id}`)
    if (data.expires_at) {
      expiresAt.value = data.expires_at
    }
  } catch {
    error.value = 'Gagal memuat data pembayaran'
  } finally {
    loading.value = false
  }
})
</script>
