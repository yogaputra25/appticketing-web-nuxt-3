<template>
  <div>
    <NuxtLink to="/admin/events" class="text-sm text-gray-500 hover:text-primary-600 mb-4 inline-block">
      &larr; Kembali ke Events
    </NuxtLink>

    <h2 class="text-xl font-semibold text-gray-900 mb-6">Tambah Event Baru</h2>

    <div class="card p-5 md:p-6 max-w-2xl">
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

          <div v-if="saveError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3">
            {{ saveError }}
          </div>

          <button type="submit" class="btn-primary touch-target" :disabled="saving">
            {{ saving ? 'Menyimpan...' : 'Simpan Event' }}
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

const saving = ref(false)
const saveError = ref('')
const router = useRouter()

const form = reactive({
  title: '',
  venue: '',
  description: '',
  start_date: '',
  end_date: '',
})

async function handleSave() {
  saving.value = true
  saveError.value = ''
  try {
    const api = useApi()
    await api.post('/admin/events', form)
    router.push('/admin/events')
  } catch (err: any) {
    saveError.value = err?.message || 'Gagal menyimpan event'
  } finally {
    saving.value = false
  }
}
</script>
