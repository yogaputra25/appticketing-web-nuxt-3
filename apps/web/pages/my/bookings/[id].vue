<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-24 md:pb-8">
    <NuxtLink to="/my/bookings" class="text-sm text-gray-500 hover:text-primary-600 mb-4 inline-block">
      &larr; Kembali
    </NuxtLink>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="card p-5 text-center">
      <p class="text-red-600 font-medium">Gagal memuat detail booking</p>
      <p class="text-sm text-gray-500 mt-1">{{ error }}</p>
      <NuxtLink to="/my/bookings" class="text-primary-600 hover:text-primary-700 mt-3 inline-block text-sm">
        &larr; Kembali
      </NuxtLink>
    </div>

    <div v-else-if="!booking" class="text-center py-16">
      <p class="text-gray-500 text-lg">Booking tidak ditemukan</p>
    </div>

    <template v-else>
      <div class="card p-5 md:p-6">
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 mb-4">E-Ticket</h1>

        <div class="space-y-3 mb-6">
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Kode Booking</span>
            <span class="text-sm font-mono font-medium text-gray-900">{{ booking.booking_code }}</span>
          </div>
          <div v-if="booking.event" class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Event</span>
            <span class="text-sm font-medium text-gray-900 text-right">{{ booking.event.title }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Status</span>
            <span class="badge" :class="statusClass">{{ booking.status }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Tiket</span>
            <span class="text-sm font-medium text-gray-900">{{ ticketCount }} tiket</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Total</span>
            <span class="text-sm font-bold text-gray-900">{{ formatPrice(booking.total_amount) }}</span>
          </div>
          <div v-if="booking.expires_at && booking.status === 'pending'" class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Batas Bayar</span>
            <span class="text-sm font-medium text-gray-900">{{ formatDate(booking.expires_at) }}</span>
          </div>
        </div>

        <div v-if="Array.isArray(booking.e_ticket_codes) && booking.e_ticket_codes.length" class="border-t border-dashed border-gray-200 pt-4 mt-4">
          <h3 class="font-semibold text-gray-900 mb-3">E-Ticket Codes</h3>
          <div class="space-y-2">
            <div v-for="(code, i) in booking.e_ticket_codes" :key="i" class="bg-gray-50 rounded-lg px-4 py-3 font-mono text-sm text-gray-800 break-all">
              {{ code }}
            </div>
          </div>
          <p class="text-center text-xs text-gray-400 mt-3">Tunjukkan kode ini di pintu masuk</p>
        </div>
      </div>

      <div v-if="booking.status === 'pending'" class="mt-4">
        <div class="card p-4 flex flex-col sm:flex-row gap-3">
          <NuxtLink :to="`/bookings/${booking.id}/pay`" class="btn-primary flex-1 touch-target flex items-center justify-center text-sm">
            Lanjut Bayar
          </NuxtLink>
          <button class="btn-outline flex-1 touch-target flex items-center justify-center text-sm text-red-600 border-red-200 hover:bg-red-50" :disabled="cancelling" @click="handleCancel">
            {{ cancelling ? 'Membatalkan...' : 'Batalkan Pesanan' }}
          </button>
        </div>
        <div v-if="cancelError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3 mt-3">
          {{ cancelError }}
        </div>
      </div>
    </template>

    <div class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-4 py-3 flex gap-3">
      <NuxtLink to="/my/bookings" class="btn-outline flex-1 touch-target flex items-center justify-center text-sm">
        Kembali
      </NuxtLink>
      <NuxtLink v-if="booking?.status === 'pending'" :to="`/bookings/${booking.id}/pay`" class="btn-primary flex-1 touch-target flex items-center justify-center text-sm">
        Lanjut Bayar
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const bookingStore = useBookingStore()

const booking = computed(() => bookingStore.currentBooking)
const loading = computed(() => bookingStore.loading)
const cancelling = ref(false)
const cancelError = ref('')
const error = ref('')

const statusClass = computed(() => {
  const classes: Record<string, string> = {
    paid: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    cancelled: 'bg-red-100 text-red-800',
    expired: 'bg-gray-100 text-gray-800',
  }
  return classes[booking.value?.status || ''] || 'bg-gray-100 text-gray-800'
})

const ticketCount = computed(() => {
  if (!booking.value) return 0
  if (Array.isArray(booking.value.e_ticket_codes) && booking.value.e_ticket_codes.length) return booking.value.e_ticket_codes.length
  if (booking.value.items) return booking.value.items.reduce((s: number, i: any) => s + i.quantity, 0)
  return 0
})

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price)
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function handleCancel() {
  cancelling.value = true
  cancelError.value = ''
  try {
    await bookingStore.cancelBooking(Number(route.params.id))
  } catch (err: any) {
    cancelError.value = err?.message || 'Gagal membatalkan booking'
  } finally {
    cancelling.value = false
  }
}

onMounted(async () => {
  try {
    await bookingStore.fetchBookingDetail(Number(route.params.id))
  } catch (err: any) {
    error.value = err?.message || 'Terjadi kesalahan saat memuat data booking'
  }
})
</script>
