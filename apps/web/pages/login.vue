<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Masuk</h1>
        <p class="text-gray-500 mt-2">Masuk ke akun War Tiket kamu</p>
      </div>

      <div class="card p-6">
        <form @submit.prevent="handleSubmit">
          <div class="space-y-4">
            <div>
              <label class="label" for="email">Email</label>
              <input
                id="email"
                v-model="email"
                type="email"
                class="input"
                placeholder="email@example.com"
                autocomplete="email"
              />
              <p v-if="errors.email" class="text-sm text-red-600 mt-1">{{ errors.email }}</p>
            </div>

            <div>
              <label class="label" for="password">Password</label>
              <input
                id="password"
                v-model="password"
                type="password"
                class="input"
                placeholder="Masukkan password"
                autocomplete="current-password"
              />
              <p v-if="errors.password" class="text-sm text-red-600 mt-1">{{ errors.password }}</p>
            </div>

            <p v-if="apiError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3">{{ apiError }}</p>

            <button type="submit" class="btn-primary w-full py-2.5" :disabled="auth.loading">
              <svg v-if="auth.loading" class="animate-spin -ml-1 mr-2 h-4 w-4 inline" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
              </svg>
              {{ auth.loading ? 'Memproses...' : 'Masuk' }}
            </button>
          </div>
        </form>

        <p class="text-center text-sm text-gray-500 mt-6">
          Belum punya akun?
          <NuxtLink to="/register" class="text-primary-600 hover:text-primary-700 font-medium">
            Daftar
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { z } from 'zod'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const errors = ref<Record<string, string>>({})
const apiError = ref('')

const loginSchema = z.object({
  email: z.string().min(1, 'Email wajib diisi').email('Format email tidak valid'),
  password: z.string().min(1, 'Password wajib diisi'),
})

async function handleSubmit() {
  errors.value = {}
  apiError.value = ''

  const result = loginSchema.safeParse({ email: email.value, password: password.value })
  if (!result.success) {
    const fieldErrors: Record<string, string> = {}
    result.error.errors.forEach((e) => {
      const field = e.path[0] as string
      fieldErrors[field] = e.message
    })
    errors.value = fieldErrors
    return
  }

  try {
    await auth.login(email.value, password.value)
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    apiError.value = err?.message || 'Email atau password salah'
  }
}
</script>
