<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-24 md:pb-8">
    <NuxtLink to="/my/bookings" class="text-sm text-gray-500 hover:text-primary-600 mb-4 inline-block">
      &larr; Kembali
    </NuxtLink>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="!booking" class="text-center py-16">
      <p class="text-gray-500 text-lg">Booking tidak ditemukan</p>
    </div>

    <template v-else>
      <div class="card p-5 md:p-6">
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 mb-4">E-Ticket</h1>

        <div class="space-y-3 mb-6">
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Event</span>
            <span class="text-sm font-medium text-gray-900 text-right">{{ booking.event_title }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Status</span>
            <span class="badge" :class="statusClass">{{ booking.status }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Tiket</span>
            <span class="text-sm font-medium text-gray-900">{{ booking.ticket_count }} tiket</span>
          </div>
          <div v-if="booking.total_price" class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Total</span>
            <span class="text-sm font-bold text-gray-900">{{ formatPrice(booking.total_price) }}</span>
          </div>
        </div>

        <div class="flex justify-center py-6 border-t border-dashed border-gray-200">
          <div class="bg-gray-100 w-48 h-48 flex items-center justify-center rounded-xl">
            <div class="text-center">
              <svg class="w-16 h-16 text-gray-400 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4b2 2 0 012 2v3a2 2 0 01-2 2h-4a2 2 0 01-2-2v-3a2 2 0 012-2z" />
              </svg>
              <p class="text-xs text-gray-500">QR Code</p>
            </div>
          </div>
        </div>
        <p class="text-center text-xs text-gray-400 mt-2">Tunjukkan QR code ini di pintu masuk</p>
      </div>
    </template>

    <div class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-4 py-3 flex gap-3">
      <NuxtLink to="/my/bookings" class="btn-outline flex-1 touch-target flex items-center justify-center text-sm">
        Kembali
      </NuxtLink>
      <button class="btn-primary flex-1 touch-target flex items-center justify-center text-sm" @click="handleDownload">
        Download
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const route = useRoute()

interface Booking {
  id: number
  event_title: string
  ticket_count: number
  total_price: number
  status: string
  created_at: string
}

const booking = ref<Booking | null>(null)
const loading = ref(true)

const statusClass = computed(() => {
  const classes: Record<string, string> = {
    confirmed: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    cancelled: 'bg-red-100 text-red-800',
    expired: 'bg-gray-100 text-gray-800',
  }
  return classes[booking.value?.status || ''] || 'bg-gray-100 text-gray-800'
})

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price)
}

function handleDownload() {
  // placeholder
}

onMounted(async () => {
  try {
    const api = useApi()
    const data = await api.get<Booking>(`/bookings/${route.params.id}`)
    booking.value = data
  } catch {
    // silent
  } finally {
    loading.value = false
  }
})
</script>
