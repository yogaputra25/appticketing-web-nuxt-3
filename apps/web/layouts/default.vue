<template>
  <div class="min-h-screen flex flex-col">
    <nav class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
          <NuxtLink to="/" class="flex items-center gap-2 font-bold text-xl text-primary-600 shrink-0">
            <span class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center text-white text-sm">W</span>
            <span class="hidden sm:inline">War Tiket</span>
          </NuxtLink>

          <div class="hidden md:flex items-center gap-6">
            <NuxtLink to="/events" class="text-gray-600 hover:text-primary-600 font-medium transition-colors">
              Events
            </NuxtLink>
          </div>

          <div class="flex items-center gap-2 md:gap-3">
            <template v-if="auth.isAuthenticated">
              <NuxtLink v-if="auth.isAdmin" to="/admin" class="btn-outline text-sm !py-1.5 hidden sm:inline-flex">
                Dashboard
              </NuxtLink>
              <NuxtLink to="/my/bookings" class="btn-outline text-sm !py-1.5 hidden sm:inline-flex">
                Tiket Saya
              </NuxtLink>
              <div class="relative" @click.outside="showDropdown = false">
                <button
                  @click="showDropdown = !showDropdown"
                  class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-gray-100 transition-colors touch-target"
                >
                  <div class="w-8 h-8 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center text-sm font-bold shrink-0">
                    {{ initials }}
                  </div>
                  <span class="text-sm font-medium text-gray-700 hidden sm:block">{{ auth.user?.full_name }}</span>
                  <svg class="w-4 h-4 text-gray-400 hidden sm:block" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
              <NuxtLink to="/login" class="btn-outline text-sm !py-1.5 !px-3">
                Masuk
              </NuxtLink>
              <NuxtLink to="/register" class="btn-primary text-sm !py-1.5 !px-3">
                Daftar
              </NuxtLink>
            </template>

            <button
              class="md:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors touch-target flex items-center justify-center"
              @click="mobileMenuOpen = !mobileMenuOpen"
              aria-label="Toggle navigation menu"
            >
              <svg v-if="!mobileMenuOpen" class="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
              <svg v-else class="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </nav>

    <Transition name="slide">
      <div
        v-if="mobileMenuOpen"
        class="fixed inset-0 z-40 md:hidden"
      >
        <div class="absolute inset-0 bg-black/50" @click="mobileMenuOpen = false" />
        <div class="absolute right-0 top-0 bottom-0 w-72 max-w-[80vw] bg-white shadow-xl overflow-y-auto">
          <div class="p-4 border-b border-gray-100">
            <div class="flex items-center gap-2 font-bold text-lg text-primary-600">
              <span class="w-7 h-7 bg-primary-600 rounded-lg flex items-center justify-center text-white text-xs">W</span>
              War Tiket
            </div>
          </div>

          <div class="p-4 space-y-1">
            <NuxtLink
              to="/events"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
              @click="mobileMenuOpen = false"
            >
              <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              Events
            </NuxtLink>

            <template v-if="auth.isAuthenticated">
              <NuxtLink
                to="/my/bookings"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
                @click="mobileMenuOpen = false"
              >
                <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
                Tiket Saya
              </NuxtLink>

              <NuxtLink
                to="/profile"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
                @click="mobileMenuOpen = false"
              >
                <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                Profile
              </NuxtLink>

              <hr class="my-2 border-gray-100">

              <button
                @click="handleLogout"
                class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-red-600 hover:bg-red-50 font-medium"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                Logout
              </button>
            </template>

            <template v-else>
              <hr class="my-2 border-gray-100">
              <NuxtLink
                to="/login"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
                @click="mobileMenuOpen = false"
              >
                <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
                </svg>
                Masuk
              </NuxtLink>
              <NuxtLink
                to="/register"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
                @click="mobileMenuOpen = false"
              >
                <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                </svg>
                Daftar
              </NuxtLink>
            </template>
          </div>
        </div>
      </div>
    </Transition>

    <main class="flex-1">
      <slot />
    </main>

    <footer class="bg-white border-t border-gray-200 py-6 md:py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col md:flex-row justify-between items-center gap-3 md:gap-4">
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
const mobileMenuOpen = ref(false)

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
  mobileMenuOpen.value = false
  auth.logout()
  router.push('/')
}
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: opacity 0.2s ease;
}
.slide-enter-from,
.slide-leave-to {
  opacity: 0;
}
</style>
