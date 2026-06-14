<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-24 md:pb-8">
    <NuxtLink to="/my/tickets" class="text-sm text-gray-500 hover:text-primary-600 mb-4 inline-block">
      &larr; Kembali
    </NuxtLink>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="card p-5 text-center">
      <p class="text-red-600 font-medium">Gagal memuat detail tiket</p>
      <p class="text-sm text-gray-500 mt-1">{{ error }}</p>
      <NuxtLink to="/my/tickets" class="text-primary-600 hover:text-primary-700 mt-3 inline-block text-sm">
        &larr; Kembali
      </NuxtLink>
    </div>

    <div v-else-if="!ticket" class="text-center py-16">
      <p class="text-gray-500 text-lg">Tiket tidak ditemukan</p>
    </div>

    <template v-else>
      <div class="card p-5 md:p-6">
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 mb-4">E-Ticket</h1>

        <div class="space-y-3 mb-6">
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Event</span>
            <span class="text-sm font-medium text-gray-900 text-right">{{ ticket.event_title }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Kategori</span>
            <span class="text-sm font-medium text-gray-900">{{ ticket.category_name }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Status</span>
            <span class="badge" :class="statusClass">{{ statusLabel }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Kode Tiket</span>
            <span class="text-sm font-mono font-medium text-gray-900">{{ ticket.ticket_code }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Kode Booking</span>
            <span class="text-sm font-mono font-medium text-gray-900">{{ ticket.booking_code }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100">
            <span class="text-sm text-gray-500">Tanggal Event</span>
            <span class="text-sm font-medium text-gray-900">{{ formatDate(ticket.start_date) }}</span>
          </div>
        </div>

        <div v-if="ticket.status === 'active'" class="border-t border-dashed border-gray-200 pt-6 mt-2">
          <div class="flex flex-col items-center">
            <h3 class="font-semibold text-gray-900 mb-3">QR Code Tiket</h3>
            <div v-if="qrDataUrl" class="bg-white rounded-xl p-3 shadow-sm border border-gray-200">
              <img :src="qrDataUrl" alt="QR Code Tiket" class="w-48 h-48 md:w-56 md:h-56" />
            </div>
            <div v-else class="w-48 h-48 md:w-56 md:h-56 bg-gray-100 rounded-xl flex items-center justify-center">
              <div class="w-6 h-6 border-2 border-primary-600 border-t-transparent rounded-full animate-spin" />
            </div>
            <p class="text-xs text-gray-400 mt-3 text-center max-w-xs">
              Tunjukkan QR code ini di pintu masuk untuk verifikasi
            </p>
          </div>
        </div>

        <div v-else-if="ticket.status === 'used'" class="border-t border-dashed border-gray-200 pt-6 mt-2">
          <div class="text-center py-4">
            <svg class="w-12 h-12 text-gray-300 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <p class="text-sm text-gray-500">Tiket sudah digunakan</p>
            <p v-if="ticket.scanned_at" class="text-xs text-gray-400 mt-1">{{ formatDate(ticket.scanned_at) }}</p>
          </div>
        </div>
      </div>
    </template>

    <div class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-4 py-3">
      <NuxtLink to="/my/tickets" class="btn-outline w-full touch-target flex items-center justify-center text-sm">
        Kembali
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const ticketStore = useTicketStore()
const { generateDataUrl } = useQrCode()

const ticket = computed(() => ticketStore.currentTicket)
const loading = computed(() => ticketStore.loading)
const error = ref('')
const qrDataUrl = ref('')

const statusClass = computed(() => {
  const classes: Record<string, string> = {
    active: 'bg-green-100 text-green-800',
    used: 'bg-gray-100 text-gray-800',
    refunded: 'bg-red-100 text-red-800',
  }
  return classes[ticket.value?.status || ''] || 'bg-gray-100 text-gray-800'
})

const statusLabel = computed(() => {
  const labels: Record<string, string> = {
    active: 'Aktif',
    used: 'Terpakai',
    refunded: 'Refund',
  }
  return labels[ticket.value?.status || ''] || ticket.value?.status || ''
})

function formatDate(dateStr: string | null | undefined) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(async () => {
  try {
    await ticketStore.fetchTicketDetail(Number(route.params.id))
    if (ticket.value?.status === 'active' && ticket.value?.ticket_code) {
      qrDataUrl.value = await generateDataUrl(ticket.value.ticket_code)
    }
  } catch (err: any) {
    error.value = err?.message || 'Terjadi kesalahan saat memuat data tiket'
  }
})
</script>
