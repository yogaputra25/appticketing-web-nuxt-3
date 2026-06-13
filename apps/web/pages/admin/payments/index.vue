<template>
  <div>
    <h2 class="text-xl font-semibold text-gray-900 mb-6">Payments</h2>

    <div class="card p-4 mb-6 flex flex-col sm:flex-row gap-4 items-start sm:items-center">
      <div class="relative flex-1 w-full sm:max-w-xs">
        <input v-model="search" type="text" class="input !h-[44px] pl-9" placeholder="Cari payment..." @input="debouncedSearch" />
        <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      <select v-model="statusFilter" class="input !h-[44px] w-full sm:w-40" @change="fetchPayments(1)">
        <option value="">Semua Status</option>
        <option value="pending">Pending</option>
        <option value="success">Success</option>
        <option value="failed">Failed</option>
        <option value="expired">Expired</option>
        <option value="refunded">Refunded</option>
      </select>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg p-4">{{ error }}</div>

    <div v-else-if="payments.length === 0" class="text-center py-12">
      <svg class="w-12 h-12 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
      </svg>
      <p class="text-gray-500 text-sm">Belum ada payment</p>
      <p class="text-gray-400 text-xs mt-1">Coba ubah filter atau kata kunci pencarian</p>
    </div>

    <template v-else>
      <div class="overflow-x-auto -mx-4 md:-mx-0">
        <div class="inline-block min-w-full align-middle px-4 md:px-0">
          <table class="w-full min-w-[600px]">
            <thead>
              <tr class="border-b border-gray-200 text-left text-sm text-gray-500">
                <th class="pb-3 font-medium">ID</th>
                <th class="pb-3 font-medium">Booking</th>
                <th class="pb-3 font-medium">Metode</th>
                <th class="pb-3 font-medium">Status</th>
                <th class="pb-3 font-medium">Jumlah</th>
                <th class="pb-3 font-medium">Tanggal</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in payments" :key="p.id" class="border-b border-gray-100 hover:bg-gray-50">
                <td class="py-3 text-sm text-gray-600 font-mono">{{ p.payment_code || p.id }}</td>
                <td class="py-3 text-sm text-gray-600">{{ p.booking_code || p.booking_id?.slice(0, 8) || '-' }}</td>
                <td class="py-3 text-sm text-gray-600">{{ p.payment_method || '-' }}</td>
                <td class="py-3">
                  <span class="badge text-xs" :class="statusBadge(p.status)">{{ p.status }}</span>
                </td>
                <td class="py-3 text-sm text-gray-600">{{ formatPrice(p.amount || 0) }}</td>
                <td class="py-3 text-sm text-gray-600">{{ formatDate(p.created_at || p.paid_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-6">
        <button class="btn-outline text-sm !py-1.5" :disabled="page <= 1" @click="fetchPayments(page - 1)">&larr; Prev</button>
        <span class="text-sm text-gray-500">Halaman {{ page }} dari {{ totalPages }}</span>
        <button class="btn-outline text-sm !py-1.5" :disabled="page >= totalPages" @click="fetchPayments(page + 1)">Next &rarr;</button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'admin',
  layout: 'admin',
})

const loading = ref(true)
const error = ref('')
const payments = ref<any[]>([])
const page = ref(1)
const totalPages = ref(1)
const search = ref('')
const statusFilter = ref('')
let searchTimeout: ReturnType<typeof setTimeout>

function statusBadge(status: string) {
  const map: Record<string, string> = {
    paid: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    failed: 'bg-red-100 text-red-800',
    expired: 'bg-gray-100 text-gray-800',
    refunded: 'bg-purple-100 text-purple-800',
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
  searchTimeout = setTimeout(() => fetchPayments(1), 300)
}

async function fetchPayments(p: number) {
  loading.value = true
  error.value = ''
  page.value = p
  try {
    const api = useApi()
    const params: Record<string, any> = { page: p, per_page: 10 }
    if (search.value) params.search = search.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get<{ data: any[]; total: number; page: number }>('/api/admin/payments', params)
    payments.value = res.data || []
    totalPages.value = Math.ceil((res.total || 0) / 10) || 1
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat data payments'
    payments.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchPayments(1))
</script>
