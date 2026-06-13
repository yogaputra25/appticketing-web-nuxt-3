<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-gray-900">Events</h2>
      <NuxtLink to="/admin/events/new" class="btn-primary text-sm !py-1.5 !px-3 md:!px-4 touch-target inline-flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span class="hidden sm:inline">Tambah Event</span>
      </NuxtLink>
    </div>

    <div class="card p-4 mb-6 flex flex-col sm:flex-row gap-4 items-start sm:items-center">
      <div class="relative flex-1 w-full sm:max-w-xs">
        <input v-model="search" type="text" class="input !h-[44px] pl-9" placeholder="Cari event..." @input="debouncedSearch" />
        <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      <select v-model="statusFilter" class="input !h-[44px] w-full sm:w-40" @change="fetchEvents(1)">
        <option value="">Semua Status</option>
        <option value="draft">Draft</option>
        <option value="published">Published</option>
        <option value="cancelled">Cancelled</option>
        <option value="finished">Finished</option>
      </select>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg p-4">{{ error }}</div>

    <div v-else-if="events.length === 0" class="text-center py-12">
      <svg class="w-12 h-12 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <p class="text-gray-500 text-sm">Belum ada event</p>
      <p class="text-gray-400 text-xs mt-1">Coba ubah filter atau buat event baru</p>
    </div>

    <template v-else>
      <div class="overflow-x-auto -mx-4 md:-mx-0">
        <div class="inline-block min-w-full align-middle px-4 md:px-0">
          <table class="w-full min-w-[600px]">
            <thead>
              <tr class="border-b border-gray-200 text-left text-sm text-gray-500">
                <th class="pb-3 font-medium">Nama</th>
                <th class="pb-3 font-medium">Status</th>
                <th class="pb-3 font-medium">Tanggal</th>
                <th class="pb-3 font-medium">Tiket Terjual</th>
                <th class="pb-3 font-medium" />
              </tr>
            </thead>
            <tbody>
              <tr v-for="e in events" :key="e.id" class="border-b border-gray-100 hover:bg-gray-50">
                <td class="py-3 text-sm">
                  <NuxtLink :to="`/admin/events/${e.id}`" class="font-medium text-gray-900 hover:text-primary-600">
                    {{ e.title }}
                  </NuxtLink>
                </td>
                <td class="py-3">
                  <span class="badge text-xs" :class="statusBadge(e.status)">{{ e.status }}</span>
                </td>
                <td class="py-3 text-sm text-gray-600">{{ formatDate(e.start_date) }}</td>
                <td class="py-3 text-sm text-gray-600">{{ e.tickets_sold ?? 0 }}</td>
                <td class="py-3 text-right">
                  <div class="flex items-center gap-2 justify-end">
                    <NuxtLink :to="`/admin/events/${e.id}`" class="text-primary-600 hover:text-primary-800 text-sm">Edit</NuxtLink>
                    <button class="text-red-600 hover:text-red-800 text-sm" @click="confirmDelete(e)">Hapus</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-6">
        <button class="btn-outline text-sm !py-1.5" :disabled="page <= 1" @click="fetchEvents(page - 1)">
          &larr; Prev
        </button>
        <span class="text-sm text-gray-500">Halaman {{ page }} dari {{ totalPages }}</span>
        <button class="btn-outline text-sm !py-1.5" :disabled="page >= totalPages" @click="fetchEvents(page + 1)">
          Next &rarr;
        </button>
      </div>
    </template>

    <div v-if="showDeleteModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="showDeleteModal = false">
      <div class="bg-white rounded-xl p-6 max-w-sm w-full">
        <h3 class="font-semibold text-gray-900 mb-2">Hapus Event</h3>
        <p class="text-sm text-gray-600 mb-6">Yakin ingin menghapus <strong>{{ deleteTarget?.title }}</strong>? Tindakan ini tidak bisa dibatalkan.</p>
        <div v-if="deleteError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3 mb-4">{{ deleteError }}</div>
        <div class="flex gap-3 justify-end">
          <button class="btn-outline text-sm" @click="showDeleteModal = false">Batal</button>
          <button class="btn-primary text-sm bg-red-600 hover:bg-red-700" :disabled="deleting" @click="handleDelete">
            {{ deleting ? 'Menghapus...' : 'Hapus' }}
          </button>
        </div>
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
const events = ref<any[]>([])
const page = ref(1)
const totalPages = ref(1)
const search = ref('')
const statusFilter = ref('')
let searchTimeout: ReturnType<typeof setTimeout>

const showDeleteModal = ref(false)
const deleteTarget = ref<any>(null)
const deleting = ref(false)
const deleteError = ref('')

function statusBadge(status: string) {
  const map: Record<string, string> = {
    published: 'bg-green-100 text-green-800',
    draft: 'bg-yellow-100 text-yellow-800',
    cancelled: 'bg-red-100 text-red-800',
    finished: 'bg-gray-100 text-gray-800',
  }
  return map[status] || 'bg-gray-100 text-gray-800'
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => fetchEvents(1), 300)
}

async function fetchEvents(p: number) {
  loading.value = true
  error.value = ''
  page.value = p
  try {
    const api = useApi()
    const params: Record<string, any> = { page: p, per_page: 10 }
    if (search.value) params.search = search.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get<{ data: any[]; total: number; page: number }>('/api/admin/events', params)
    events.value = res.data || []
    totalPages.value = Math.ceil((res.total || 0) / 10) || 1
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat data events'
    events.value = []
  } finally {
    loading.value = false
  }
}

function confirmDelete(e: any) {
  deleteTarget.value = e
  deleteError.value = ''
  showDeleteModal.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  deleteError.value = ''
  try {
    const api = useApi()
    await api.delete(`/api/admin/events/${deleteTarget.value.id}`)
    showDeleteModal.value = false
    deleteTarget.value = null
    await fetchEvents(page.value)
  } catch (err: any) {
    deleteError.value = err?.message || 'Gagal menghapus event'
  } finally {
    deleting.value = false
  }
}

onMounted(() => fetchEvents(1))
</script>
