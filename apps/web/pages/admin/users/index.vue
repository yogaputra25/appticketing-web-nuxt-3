<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-gray-900">Users</h2>
      <button class="btn-primary text-sm !py-1.5 !px-3 md:!px-4 touch-target inline-flex items-center gap-1" @click="showCreateModal = true">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span class="hidden sm:inline">Tambah Admin</span>
      </button>
    </div>

    <div class="card p-4 mb-6 flex flex-col sm:flex-row gap-4 items-start sm:items-center">
      <div class="relative flex-1 w-full sm:max-w-xs">
        <input v-model="search" type="text" class="input !h-[44px] pl-9" placeholder="Cari user..." @input="debouncedSearch" />
        <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      <select v-model="roleFilter" class="input !h-[44px] w-full sm:w-40" @change="fetchUsers(1)">
        <option value="">Semua Role</option>
        <option value="user">User</option>
        <option value="admin">Admin</option>
      </select>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg p-4">{{ error }}</div>

    <div v-else-if="users.length === 0" class="text-center py-12">
      <svg class="w-12 h-12 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
      </svg>
      <p class="text-gray-500 text-sm">Belum ada user</p>
      <p class="text-gray-400 text-xs mt-1">Coba ubah filter atau kata kunci pencarian</p>
    </div>

    <template v-else>
      <div class="overflow-x-auto -mx-4 md:-mx-0">
        <div class="inline-block min-w-full align-middle px-4 md:px-0">
          <table class="w-full min-w-[600px]">
            <thead>
              <tr class="border-b border-gray-200 text-left text-sm text-gray-500">
                <th class="pb-3 font-medium">Nama</th>
                <th class="pb-3 font-medium">Email</th>
                <th class="pb-3 font-medium">Role</th>
                <th class="pb-3 font-medium">Tanggal Daftar</th>
                <th class="pb-3 font-medium" />
              </tr>
            </thead>
            <tbody>
              <tr v-for="u in users" :key="u.id" class="border-b border-gray-100 hover:bg-gray-50">
                <td class="py-3 text-sm font-medium text-gray-900">{{ u.full_name || '-' }}</td>
                <td class="py-3 text-sm text-gray-600">{{ u.email }}</td>
                <td class="py-3">
                  <span class="badge text-xs" :class="u.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-gray-100 text-gray-800'">
                    {{ u.role }}
                  </span>
                </td>
                <td class="py-3 text-sm text-gray-600">{{ formatDate(u.created_at) }}</td>
                <td class="py-3 text-right">
                  <button v-if="u.role !== 'admin'" class="text-primary-600 hover:text-primary-800 text-sm" @click="promoteToAdmin(u)">Jadikan Admin</button>
                  <span v-else class="text-gray-400 text-sm">Admin</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-6">
        <button class="btn-outline text-sm !py-1.5" :disabled="page <= 1" @click="fetchUsers(page - 1)">&larr; Prev</button>
        <span class="text-sm text-gray-500">Halaman {{ page }} dari {{ totalPages }}</span>
        <button class="btn-outline text-sm !py-1.5" :disabled="page >= totalPages" @click="fetchUsers(page + 1)">Next &rarr;</button>
      </div>
    </template>

    <div v-if="showCreateModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="showCreateModal = false">
      <div class="bg-white rounded-xl p-6 max-w-md w-full">
        <h3 class="font-semibold text-gray-900 mb-4">Tambah Admin Baru</h3>
        <form @submit.prevent="handleCreateAdmin">
          <div class="space-y-4">
            <div>
              <label class="label" for="new-name">Nama</label>
              <input id="new-name" v-model="newAdmin.full_name" type="text" class="input !h-[44px]" required />
            </div>
            <div>
              <label class="label" for="new-email">Email</label>
              <input id="new-email" v-model="newAdmin.email" type="email" class="input !h-[44px]" required />
            </div>
            <div>
              <label class="label" for="new-password">Password</label>
              <input id="new-password" v-model="newAdmin.password" type="password" class="input !h-[44px]" required minlength="6" />
            </div>
            <div v-if="createError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3">{{ createError }}</div>
            <div class="flex gap-3 justify-end">
              <button type="button" class="btn-outline text-sm" @click="showCreateModal = false">Batal</button>
              <button type="submit" class="btn-primary text-sm" :disabled="creating">
                {{ creating ? 'Menyimpan...' : 'Tambah Admin' }}
              </button>
            </div>
          </div>
        </form>
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
const users = ref<any[]>([])
const page = ref(1)
const totalPages = ref(1)
const search = ref('')
const roleFilter = ref('')
let searchTimeout: ReturnType<typeof setTimeout>

const showCreateModal = ref(false)
const creating = ref(false)
const createError = ref('')
const newAdmin = reactive({ full_name: '', email: '', password: '' })

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
  searchTimeout = setTimeout(() => fetchUsers(1), 300)
}

async function fetchUsers(p: number) {
  loading.value = true
  error.value = ''
  page.value = p
  try {
    const api = useApi()
    const params: Record<string, any> = { page: p, per_page: 10 }
    if (search.value) params.search = search.value
    if (roleFilter.value) params.role = roleFilter.value
    const res = await api.get<{ data: any[]; total: number; page: number }>('/api/admin/users', params)
    users.value = res.data || []
    totalPages.value = Math.ceil((res.total || 0) / 10) || 1
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat data users'
    users.value = []
  } finally {
    loading.value = false
  }
}

async function promoteToAdmin(u: any) {
  try {
    const api = useApi()
    await api.put(`/api/admin/users/${u.id}`, { role: 'admin' })
    u.role = 'admin'
  } catch (err: any) {
    error.value = err?.message || 'Gagal mengubah role user'
  }
}

async function handleCreateAdmin() {
  creating.value = true
  createError.value = ''
  try {
    const api = useApi()
    await api.post('/api/admin/users', { ...newAdmin, role: 'admin' })
    showCreateModal.value = false
    newAdmin.full_name = ''
    newAdmin.email = ''
    newAdmin.password = ''
    await fetchUsers(1)
  } catch (err: any) {
    createError.value = err?.message || 'Gagal membuat admin baru'
  } finally {
    creating.value = false
  }
}

onMounted(() => fetchUsers(1))
</script>
