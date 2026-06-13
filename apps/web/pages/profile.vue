<template>
  <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-24 md:pb-8">
    <div class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-gray-900">Profile</h1>
      <p class="text-gray-500 mt-1 text-sm md:text-base">Kelola informasi akun kamu</p>
    </div>

    <div class="card p-5 md:p-6">
      <form @submit.prevent="handleSave">
        <div class="space-y-4 md:space-y-5">
          <div>
            <label class="label" for="full_name">Nama Lengkap</label>
            <input
              id="full_name"
              v-model="form.full_name"
              type="text"
              class="input !h-[44px]"
              placeholder="Nama lengkap"
            />
          </div>

          <div>
            <label class="label" for="email">Email</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              class="input !h-[44px]"
              placeholder="email@example.com"
              disabled
            />
          </div>

          <div>
            <label class="label" for="phone">No. Telepon</label>
            <input
              id="phone"
              v-model="form.phone"
              type="tel"
              class="input !h-[44px]"
              placeholder="08xxxxxxxxxx"
            />
          </div>

          <div v-if="saveError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3">
            {{ saveError }}
          </div>
          <div v-if="saveSuccess" class="text-sm text-green-600 bg-green-50 rounded-lg p-3">
            Profile berhasil diperbarui
          </div>
        </div>
      </form>
    </div>

    <div class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-4 py-3">
      <button class="btn-primary w-full touch-target" :disabled="saving" @click="handleSave">
        {{ saving ? 'Menyimpan...' : 'Simpan Profile' }}
      </button>
    </div>

    <div class="hidden md:block mt-6">
      <button class="btn-primary px-8 touch-target" :disabled="saving" @click="handleSave">
        {{ saving ? 'Menyimpan...' : 'Simpan Profile' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const auth = useAuthStore()
const saving = ref(false)
const saveError = ref('')
const saveSuccess = ref(false)

const form = reactive({
  full_name: auth.user?.full_name || '',
  email: auth.user?.email || '',
  phone: auth.user?.phone || '',
})

async function handleSave() {
  saving.value = true
  saveError.value = ''
  saveSuccess.value = false
  try {
    const api = useApi()
    await api.put('/api/auth/me', {
      full_name: form.full_name,
      phone: form.phone || undefined,
    })
    await auth.fetchMe()
    saveSuccess.value = true
    setTimeout(() => { saveSuccess.value = false }, 3000)
  } catch (err: any) {
    saveError.value = err?.message || 'Gagal menyimpan profile'
  } finally {
    saving.value = false
  }
}
</script>
