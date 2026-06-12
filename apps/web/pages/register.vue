<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4 py-6 md:py-8">
    <div class="w-full max-w-md">
      <div class="text-center mb-6 md:mb-8">
        <h1 class="text-2xl md:text-3xl font-bold text-gray-900">Daftar</h1>
        <p class="text-gray-500 mt-2 text-sm md:text-base">Buat akun War Tiket baru</p>
      </div>

      <div class="card p-5 md:p-6">
        <form @submit.prevent="handleSubmit">
          <div class="space-y-4 md:space-y-5">
            <div>
              <label class="label" for="full_name">Nama Lengkap</label>
              <input
                id="full_name"
                v-model="fullName"
                type="text"
                class="input !h-[44px]"
                placeholder="Nama lengkap"
                autocomplete="name"
              />
              <p v-if="errors.full_name" class="text-sm text-red-600 mt-1">{{ errors.full_name }}</p>
            </div>

            <div>
              <label class="label" for="email">Email</label>
              <input
                id="email"
                v-model="email"
                type="email"
                class="input !h-[44px]"
                placeholder="email@example.com"
                autocomplete="email"
              />
              <p v-if="errors.email" class="text-sm text-red-600 mt-1">{{ errors.email }}</p>
            </div>

            <div>
              <label class="label" for="phone">No. Telepon <span class="text-gray-400">(opsional)</span></label>
              <input
                id="phone"
                v-model="phone"
                type="tel"
                class="input !h-[44px]"
                placeholder="08xxxxxxxxxx"
                autocomplete="tel"
              />
            </div>

            <div>
              <label class="label" for="password">Password</label>
              <input
                id="password"
                v-model="password"
                type="password"
                class="input !h-[44px]"
                placeholder="Minimal 8 karakter"
                autocomplete="new-password"
              />
              <p v-if="errors.password" class="text-sm text-red-600 mt-1">{{ errors.password }}</p>
            </div>

            <div>
              <label class="label" for="confirm_password">Konfirmasi Password</label>
              <input
                id="confirm_password"
                v-model="confirmPassword"
                type="password"
                class="input !h-[44px]"
                placeholder="Ulangi password"
                autocomplete="new-password"
              />
              <p v-if="errors.confirm_password" class="text-sm text-red-600 mt-1">{{ errors.confirm_password }}</p>
            </div>

            <p v-if="apiError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3">{{ apiError }}</p>

            <button type="submit" class="btn-primary w-full !h-[44px]" :disabled="auth.loading">
              <svg v-if="auth.loading" class="animate-spin -ml-1 mr-2 h-4 w-4 inline" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
              </svg>
              {{ auth.loading ? 'Memproses...' : 'Daftar' }}
            </button>
          </div>
        </form>

        <p class="text-center text-sm text-gray-500 mt-6">
          Sudah punya akun?
          <NuxtLink to="/login" class="text-primary-600 hover:text-primary-700 font-medium">
            Masuk
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

const fullName = ref('')
const email = ref('')
const phone = ref('')
const password = ref('')
const confirmPassword = ref('')
const errors = ref<Record<string, string>>({})
const apiError = ref('')

const registerSchema = z
  .object({
    full_name: z.string().min(2, 'Nama minimal 2 karakter').max(255, 'Nama terlalu panjang'),
    email: z.string().min(1, 'Email wajib diisi').email('Format email tidak valid'),
    phone: z.string().optional(),
    password: z.string().min(8, 'Password minimal 8 karakter'),
    confirm_password: z.string().min(1, 'Konfirmasi password wajib diisi'),
  })
  .refine((data) => data.password === data.confirm_password, {
    message: 'Password tidak cocok',
    path: ['confirm_password'],
  })

async function handleSubmit() {
  errors.value = {}
  apiError.value = ''

  const result = registerSchema.safeParse({
    full_name: fullName.value,
    email: email.value,
    phone: phone.value || undefined,
    password: password.value,
    confirm_password: confirmPassword.value,
  })

  if (!result.success) {
    const fieldErrors: Record<string, string> = {}
    result.error.errors.forEach((e) => {
      const field = e.path[0] as string
      if (!fieldErrors[field]) {
        fieldErrors[field] = e.message
      }
    })
    errors.value = fieldErrors
    return
  }

  try {
    await auth.register({
      email: email.value,
      password: password.value,
      full_name: fullName.value,
      phone: phone.value || undefined,
    })
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    if (err?.fields) {
      errors.value = err.fields
    }
    apiError.value = err?.message || 'Pendaftaran gagal, silakan coba lagi'
  }
}
</script>
