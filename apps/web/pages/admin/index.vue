<template>
  <div>
    <h2 class="text-xl font-semibold text-gray-900 mb-6">Overview</h2>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg p-4 mb-6">
      {{ error }}
    </div>

    <template v-else>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6 mb-8">
        <div class="card p-5">
          <p class="text-sm text-gray-500 mb-1">Total Events</p>
          <p class="text-3xl font-bold text-gray-900">{{ stats.total_events }}</p>
        </div>
        <div class="card p-5">
          <p class="text-sm text-gray-500 mb-1">Total Bookings</p>
          <p class="text-3xl font-bold text-gray-900">{{ stats.total_bookings }}</p>
        </div>
        <div class="card p-5">
          <p class="text-sm text-gray-500 mb-1">Revenue</p>
          <p class="text-3xl font-bold text-gray-900">{{ formatPrice(stats.total_revenue) }}</p>
        </div>
        <div class="card p-5">
          <p class="text-sm text-gray-500 mb-1">Total Users</p>
          <p class="text-3xl font-bold text-gray-900">{{ stats.total_users }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="card p-5">
          <h3 class="font-semibold text-gray-900 mb-4">Recent Bookings</h3>
          <div v-if="recentBookings.length === 0" class="text-sm text-gray-500 text-center py-4">
            Belum ada bookings
          </div>
          <div v-else class="space-y-3">
            <div v-for="b in recentBookings" :key="b.id" class="flex items-center justify-between text-sm">
              <div class="min-w-0">
                <p class="font-medium text-gray-900 truncate">{{ b.event_title }}</p>
                <p class="text-gray-500">{{ formatDate(b.created_at) }}</p>
              </div>
              <span class="badge shrink-0" :class="statusBadge(b.status)">{{ b.status }}</span>
            </div>
          </div>
        </div>

        <div class="card p-5">
          <h3 class="font-semibold text-gray-900 mb-4">Upcoming Events</h3>
          <div v-if="upcomingEvents.length === 0" class="text-sm text-gray-500 text-center py-4">
            Tidak ada event mendatang
          </div>
          <div v-else class="space-y-3">
            <div v-for="e in upcomingEvents" :key="e.id" class="flex items-center justify-between text-sm">
              <div class="min-w-0">
                <NuxtLink :to="`/admin/events/${e.id}`" class="font-medium text-gray-900 hover:text-primary-600 truncate block">
                  {{ e.title }}
                </NuxtLink>
                <p class="text-gray-500">{{ formatDate(e.start_date) }}</p>
              </div>
              <span class="badge shrink-0" :class="statusBadge(e.status)">{{ e.status }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'admin',
  layout: 'admin',
})

interface AdminStats {
  total_events: number
  total_bookings: number
  total_revenue: number
  total_users: number
}

const loading = ref(true)
const error = ref('')
const stats = ref<AdminStats>({ total_events: 0, total_bookings: 0, total_revenue: 0, total_users: 0 })
const recentBookings = ref<any[]>([])
const upcomingEvents = ref<any[]>([])

function statusBadge(status: string) {
  const map: Record<string, string> = {
    published: 'bg-green-100 text-green-800',
    draft: 'bg-yellow-100 text-yellow-800',
    cancelled: 'bg-red-100 text-red-800',
    finished: 'bg-gray-100 text-gray-800',
    confirmed: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    paid: 'bg-green-100 text-green-800',
    expired: 'bg-gray-100 text-gray-800',
  }
  return map[status] || 'bg-gray-100 text-gray-800'
}

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price)
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

onMounted(async () => {
  try {
    const api = useApi()
    const [statsData, bookingsData, eventsData] = await Promise.all([
      api.get<AdminStats>('/api/admin/stats').catch(() => ({ total_events: 0, total_bookings: 0, total_revenue: 0, total_users: 0 })),
      api.get<any[]>('/api/admin/bookings?per_page=5').catch(() => []),
      api.get<any>('/api/events?per_page=5&status=published').catch(() => ({ data: [] })),
    ])
    stats.value = statsData
    recentBookings.value = bookingsData
    upcomingEvents.value = eventsData.data || []
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat data dashboard'
  } finally {
    loading.value = false
  }
})
</script>
