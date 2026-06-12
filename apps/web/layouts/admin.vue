<template>
  <div class="min-h-screen flex">
    <Transition name="slide-sidebar">
      <aside
        v-if="sidebarOpen || !isMobile"
        class="bg-gray-900 text-white flex flex-col shrink-0 fixed md:static inset-y-0 left-0 z-50 w-64 transition-transform md:translate-x-0"
        :class="isMobile && !sidebarOpen ? '-translate-x-full' : ''"
      >
        <div class="p-4 border-b border-gray-800 flex items-center justify-between">
          <NuxtLink to="/admin" class="flex items-center gap-2 font-bold text-lg" @click="sidebarOpen = false">
            <span class="w-8 h-8 bg-accent rounded-lg flex items-center justify-center text-white text-sm">W</span>
            <span class="hidden sm:inline">Admin Panel</span>
          </NuxtLink>
          <button
            v-if="isMobile"
            class="p-1 rounded-lg hover:bg-gray-800 transition-colors"
            @click="sidebarOpen = false"
            aria-label="Close sidebar"
          >
            <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <nav class="flex-1 p-4 space-y-1 overflow-y-auto">
          <NuxtLink
            v-for="item in menuItems"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors"
            :class="isActive(item.to) ? 'bg-primary-600 text-white' : 'text-gray-300 hover:bg-gray-800 hover:text-white'"
            @click="isMobile && (sidebarOpen = false)"
          >
            <component :is="item.icon" class="w-5 h-5 shrink-0" />
            {{ item.label }}
          </NuxtLink>
        </nav>

        <div class="p-4 border-t border-gray-800">
          <NuxtLink to="/" class="flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors" @click="isMobile && (sidebarOpen = false)">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali ke Website
          </NuxtLink>
        </div>
      </aside>
    </Transition>

    <div v-if="isMobile && sidebarOpen" class="fixed inset-0 bg-black/50 z-40" @click="sidebarOpen = false" />

    <div class="flex-1 flex flex-col min-w-0">
      <header class="bg-white shadow-sm border-b border-gray-200 h-16 flex items-center px-4 md:px-6 gap-4">
        <button
          class="md:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors touch-target flex items-center justify-center"
          @click="sidebarOpen = true"
          aria-label="Open sidebar"
        >
          <svg class="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>

        <h1 class="text-lg font-semibold text-gray-900 truncate">{{ pageTitle }}</h1>

        <div class="ml-auto flex items-center gap-3 shrink-0">
          <div class="w-8 h-8 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center text-sm font-bold">
            {{ initials }}
          </div>
          <span class="text-sm font-medium text-gray-700 hidden sm:block truncate max-w-[120px]">{{ auth.user?.full_name }}</span>
        </div>
      </header>

      <main class="flex-1 p-4 md:p-6 bg-gray-50 overflow-y-auto">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const route = useRoute()
const { isMobile } = useViewport()
const sidebarOpen = ref(false)

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/admin': 'Dashboard',
    '/admin/events': 'Events',
    '/admin/bookings': 'Bookings',
    '/admin/payments': 'Payments',
    '/admin/users': 'Users',
  }
  return titles[route.path] || 'Admin Panel'
})

const initials = computed(() => {
  if (!auth.user?.full_name) return '?'
  return auth.user.full_name
    .split(' ')
    .map(w => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
})

interface MenuItem {
  label: string
  to: string
  icon: any
}

const menuItems: MenuItem[] = [
  {
    label: 'Dashboard',
    to: '/admin',
    icon: defineComponent({
      template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg>',
    }),
  },
  {
    label: 'Events',
    to: '/admin/events',
    icon: defineComponent({
      template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>',
    }),
  },
  {
    label: 'Bookings',
    to: '/admin/bookings',
    icon: defineComponent({
      template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>',
    }),
  },
  {
    label: 'Payments',
    to: '/admin/payments',
    icon: defineComponent({
      template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>',
    }),
  },
  {
    label: 'Users',
    to: '/admin/users',
    icon: defineComponent({
      template: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" /></svg>',
    }),
  },
]

function isActive(path: string) {
  if (path === '/admin') return route.path === '/admin'
  return route.path.startsWith(path)
}
</script>

<style scoped>
.slide-sidebar-enter-active,
.slide-sidebar-leave-active {
  transition: transform 0.25s ease;
}
.slide-sidebar-enter-from,
.slide-sidebar-leave-to {
  transform: translateX(-100%);
}
</style>
