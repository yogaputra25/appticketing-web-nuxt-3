<template>
  <div>
    <h2 class="text-xl font-semibold text-gray-900 mb-6">Bookings</h2>

    <div class="card p-4 mb-6 flex flex-col sm:flex-row gap-4 items-start sm:items-center">
      <div class="relative flex-1 w-full sm:max-w-xs">
        <input v-model="search" type="text" class="input !h-[44px] pl-9" placeholder="Cari booking..." @input="debouncedSearch" />
        <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      <select v-model="statusFilter" class="input !h-[44px] w-full sm:w-40" @change="fetchBookings(1)">
        <option value="">Semua Status</option>
        <option value="pending">Pending</option>
        <option value="confirmed">Confirmed</option>
        <option value="cancelled">Cancelled</option>
        <option value="expired">Expired</option>
      </select>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg p-4">{{ error }}</div>

    <div v-else-if="bookings.length === 0" class="text-center py-12">
      <svg class="w-12 h-12 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      </svg>
      <p class="text-gray-500 text-sm">Belum ada booking</p>
      <p class="text-gray-400 text-xs mt-1">Coba ubah filter atau kata kunci pencarian</p>
    </div>

    <template v-else>
      <div class="overflow-x-auto -mx-4 md:-mx-0">
        <div class="inline-block min-w-full align-middle px-4 md:px-0">
          <table class="w-full min-w-[600px]">
            <thead>
              <tr class="border-b border-gray-200 text-left text-sm text-gray-500">
                <th class="pb-3 font-medium">ID</th>
                <th class="pb-3 font-medium">Event</th>
                <th class="pb-3 font-medium">User</th>
                <th class="pb-3 font-medium">Status</th>
                <th class="pb-3 font-medium">Total</th>
                <th class="pb-3 font-medium" />
              </tr>
            </thead>
            <tbody>
              <tr v-for="b in bookings" :key="b.id" class="border-b border-gray-100 hover:bg-gray-50">
                <td class="py-3 text-sm text-gray-600 font-mono">{{ b.booking_code || b.id }}</td>
                <td class="py-3 text-sm font-medium text-gray-900">{{ b.event?.title || '-' }}</td>
                <td class="py-3 text-sm text-gray-600">{{ b.user?.full_name || b.user?.email || '-' }}</td>
                <td class="py-3">
                  <span class="badge text-xs" :class="statusBadge(b.status)">{{ b.status }}</span>
                </td>
                <td class="py-3 text-sm text-gray-600">{{ formatPrice(b.total_amount || b.total || 0) }}</td>
                <td class="py-3 text-right">
                  <button class="text-primary-600 hover:text-primary-800 text-sm" @click="viewDetail(b)">Detail</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-6">
        <button class="btn-outline text-sm !py-1.5" :disabled="page <= 1" @click="fetchBookings(page - 1)">&larr; Prev</button>
        <span class="text-sm text-gray-500">Halaman {{ page }} dari {{ totalPages }}</span>
        <button class="btn-outline text-sm !py-1.5" :disabled="page >= totalPages" @click="fetchBookings(page + 1)">Next &rarr;</button>
      </div>
    </template>

    <div v-if="detail" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="detail = null">
      <div class="bg-white rounded-xl p-6 max-w-lg w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-semibold text-gray-900">Detail Booking</h3>
          <button class="text-gray-400 hover:text-gray-600" @click="detail = null">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <dl class="space-y-3 text-sm">
          <div class="flex justify-between"><dt class="text-gray-500">Kode</dt><dd class="font-medium text-gray-900 font-mono">{{ detail.booking_code || detail.id }}</dd></div>
          <div class="flex justify-between"><dt class="text-gray-500">Event</dt><dd class="font-medium text-gray-900">{{ detail.event?.title || '-' }}</dd></div>
          <div class="flex justify-between"><dt class="text-gray-500">User</dt><dd class="font-medium text-gray-900">{{ detail.user?.full_name || detail.user?.email || '-' }}</dd></div>
          <div class="flex justify-between"><dt class="text-gray-500">Status</dt><dd><span class="badge text-xs" :class="statusBadge(detail.status)">{{ detail.status }}</span></dd></div>
          <div class="flex justify-between"><dt class="text-gray-500">Total</dt><dd class="font-medium text-gray-900">{{ formatPrice(detail.total_amount || 0) }}</dd></div>
          <div class="flex justify-between"><dt class="text-gray-500">Dibuat</dt><dd class="text-gray-900">{{ formatDate(detail.created_at) }}</dd></div>
          <div v-if="detail.items?.length" class="pt-3 border-t border-gray-100">
            <dt class="text-gray-500 mb-2">Tiket</dt>
            <dd v-for="t in detail.items" :key="t.id" class="flex justify-between text-sm ml-4">
              <span>{{ t.ticket_category_id ? 'Kategori #' + t.ticket_category_id : 'Tiket' }} x{{ t.quantity }}</span>
              <span class="text-gray-900">{{ formatPrice(t.subtotal || 0) }}</span>
            </dd>
          </div>
        </dl>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'admin',
  layout: 'admin',
})

const loading = ref(true)
const error = ref('')
const bookings = ref<any[]>([])
const page = ref(1)
const totalPages = ref(1)
const search = ref('')
const statusFilter = ref('')
const detail = ref<any>(null)
let searchTimeout: ReturnType<typeof setTimeout>

function statusBadge(status: string) {
  const map: Record<string, string> = {
    confirmed: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    cancelled: 'bg-red-100 text-red-800',
    expired: 'bg-gray-100 text-gray-800',
  }
  return map[status] || 'bg-gray-100 text-gray-800'
}

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

function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => fetchBookings(1), 300)
}

async function fetchBookings(p: number) {
  loading.value = true
  error.value = ''
  page.value = p
  try {
    const api = useApi()
    const params: Record<string, any> = { page: p, per_page: 10 }
    if (search.value) params.search = search.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get<{ data: any[]; total: number; page: number }>('/api/admin/bookings', params)
    bookings.value = res.data || []
    totalPages.value = Math.ceil((res.total || 0) / 10) || 1
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat data bookings'
    bookings.value = []
  } finally {
    loading.value = false
  }
}

async function viewDetail(b: any) {
  try {
    const api = useApi()
    detail.value = await api.get(`/api/admin/bookings/${b.id}`)
  } catch {
    detail.value = b
  }
}

onMounted(() => fetchBookings(1))
</script>
