<template>
  <div class="min-h-screen flex flex-col">
    <nav class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
          <NuxtLink to="/" class="flex items-center gap-2 font-bold text-xl text-primary-600">
            <span class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center text-white text-sm">W</span>
            War Tiket
          </NuxtLink>

          <div class="hidden md:flex items-center gap-6">
            <NuxtLink to="/events" class="text-gray-600 hover:text-primary-600 font-medium transition-colors">
              Events
            </NuxtLink>
          </div>

          <div class="flex items-center gap-3">
            <template v-if="auth.isAuthenticated">
              <NuxtLink v-if="auth.isAdmin" to="/admin" class="btn-outline text-sm !py-1.5">
                Dashboard
              </NuxtLink>
              <NuxtLink to="/my/bookings" class="btn-outline text-sm !py-1.5">
                Tiket Saya
              </NuxtLink>
              <div class="relative" @click.outside="showDropdown = false">
                <button
                  @click="showDropdown = !showDropdown"
                  class="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div class="w-8 h-8 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center text-sm font-bold">
                    {{ initials }}
                  </div>
                  <span class="text-sm font-medium text-gray-700 hidden sm:block">{{ auth.user?.full_name }}</span>
                  <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                <div v-if="showDropdown" class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-200 py-1 z-50">
                  <NuxtLink to="/profile" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" @click="showDropdown = false">
                    Profile
                  </NuxtLink>
                  <hr class="my-1 border-gray-100">
                  <button
                    @click="handleLogout"
                    class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50"
                  >
                    Logout
                  </button>
                </div>
              </div>
            </template>
            <template v-else>
              <NuxtLink to="/login" class="btn-outline text-sm !py-1.5">
                Masuk
              </NuxtLink>
              <NuxtLink to="/register" class="btn-primary text-sm !py-1.5">
                Daftar
              </NuxtLink>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <main class="flex-1">
      <slot />
    </main>

    <footer class="bg-white border-t border-gray-200 py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col md:flex-row justify-between items-center gap-4">
          <div class="flex items-center gap-2 font-bold text-lg text-gray-700">
            <span class="w-7 h-7 bg-primary-600 rounded-lg flex items-center justify-center text-white text-xs">W</span>
            War Tiket
          </div>
          <p class="text-sm text-gray-500">
            &copy; {{ new Date().getFullYear() }} War Tiket. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()
const showDropdown = ref(false)

const initials = computed(() => {
  if (!auth.user?.full_name) return '?'
  return auth.user.full_name
    .split(' ')
    .map(w => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
})

function handleLogout() {
  showDropdown.value = false
  auth.logout()
  router.push('/')
}
</script>
