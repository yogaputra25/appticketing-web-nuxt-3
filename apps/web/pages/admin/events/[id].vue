<template>
  <div>
    <NuxtLink to="/admin/events" class="text-sm text-gray-500 hover:text-primary-600 mb-4 inline-block">
      &larr; Kembali ke Events
    </NuxtLink>

    <h2 class="text-xl font-semibold text-gray-900 mb-6">Edit Event</h2>

    <div v-if="loading" class="flex justify-center py-8">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else class="card p-5 md:p-6 max-w-2xl">
      <form @submit.prevent="handleSave">
        <div class="space-y-4">
          <div>
            <label class="label" for="title">Judul Event</label>
            <input id="title" v-model="form.title" type="text" class="input !h-[44px]" placeholder="Nama event" />
          </div>

          <div>
            <label class="label" for="venue">Lokasi</label>
            <input id="venue" v-model="form.venue" type="text" class="input !h-[44px]" placeholder="Nama venue" />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="label" for="start_date">Tanggal Mulai</label>
              <input id="start_date" v-model="form.start_date" type="date" class="input !h-[44px]" />
            </div>
            <div>
              <label class="label" for="end_date">Tanggal Selesai</label>
              <input id="end_date" v-model="form.end_date" type="date" class="input !h-[44px]" />
            </div>
          </div>

          <div>
            <label class="label" for="description">Deskripsi</label>
            <textarea id="description" v-model="form.description" class="input min-h-[100px]" placeholder="Deskripsi event" />
          </div>

          <div>
            <h3 class="font-semibold text-gray-900 mb-3">Kategori Tiket</h3>
            <div class="space-y-2">
              <div v-for="(cat, i) in form.categories" :key="i" class="card p-3 flex flex-col sm:flex-row gap-3">
                <input v-model="cat.name" type="text" class="input !h-[44px] flex-1" placeholder="Nama kategori" />
                <input v-model="cat.price" type="number" class="input !h-[44px] sm:w-32" placeholder="Harga" />
                <input v-model="cat.stock" type="number" class="input !h-[44px] sm:w-24" placeholder="Stok" />
                <button type="button" class="text-red-500 hover:text-red-700 touch-target shrink-0" @click="form.categories.splice(i, 1)">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
            <button type="button" class="btn-outline text-sm mt-2 touch-target" @click="form.categories.push({ name: '', price: 0, stock: 0 })">
              + Tambah Kategori
            </button>
          </div>

          <div v-if="saveError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3">
            {{ saveError }}
          </div>

          <button type="submit" class="btn-primary touch-target" :disabled="saving">
            {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'admin',
  layout: 'admin',
})

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const saving = ref(false)
const saveError = ref('')

const form = reactive({
  title: '',
  venue: '',
  description: '',
  start_date: '',
  end_date: '',
  categories: [] as { name: string; price: number; stock: number }[],
})

async function loadEvent() {
  try {
    const api = useApi()
    const data = await api.get<any>(`/admin/events/${route.params.id}`)
    form.title = data.title || ''
    form.venue = data.venue || ''
    form.description = data.description || ''
    form.start_date = data.start_date?.slice(0, 10) || ''
    form.end_date = data.end_date?.slice(0, 10) || ''
    form.categories = data.categories?.map((c: any) => ({
      name: c.name,
      price: c.price,
      stock: c.available_stock || c.stock || 0,
    })) || []
  } catch {
    saveError.value = 'Gagal memuat data event'
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  saveError.value = ''
  try {
    const api = useApi()
    await api.put(`/admin/events/${route.params.id}`, form)
    router.push('/admin/events')
  } catch (err: any) {
    saveError.value = err?.message || 'Gagal menyimpan event'
  } finally {
    saving.value = false
  }
}

onMounted(loadEvent)
</script>
