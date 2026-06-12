<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <section class="text-center py-16 md:py-24">
      <h1 class="text-4xl md:text-6xl font-bold text-gray-900 mb-4">
        War Tiket
      </h1>
      <p class="text-lg md:text-xl text-gray-600 max-w-2xl mx-auto mb-8">
        Platform pemesanan tiket konser dengan sistem antrian yang adil. Dapatkan tiket untuk konser favoritmu tanpa khawatir ketinggalan.
      </p>
      <div class="flex gap-4 justify-center">
        <NuxtLink to="/events" class="btn-primary text-lg px-8 py-3">
          Lihat Events
        </NuxtLink>
        <NuxtLink to="/register" class="btn-outline text-lg px-8 py-3">
          Daftar Sekarang
        </NuxtLink>
      </div>
    </section>

    <section class="py-12">
      <div class="flex items-center justify-between mb-8">
        <h2 class="text-2xl font-bold text-gray-900">Event Mendatang</h2>
        <NuxtLink to="/events" class="text-primary-600 hover:text-primary-700 font-medium">
          Lihat Semua &rarr;
        </NuxtLink>
      </div>
      <div v-if="eventStore.loading" class="flex justify-center py-12">
        <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
      </div>
      <p v-else-if="eventStore.events.length === 0" class="text-center text-gray-500 py-12">
        Belum ada event tersedia.
      </p>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <NuxtLink
          v-for="event in eventStore.events"
          :key="event.id"
          :to="`/events/${event.id}`"
          class="card p-0 overflow-hidden hover:shadow-md transition-shadow"
        >
          <div class="h-40 bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center">
            <span class="text-white text-3xl font-bold">{{ event.title.charAt(0) }}</span>
          </div>
          <div class="p-4">
            <div class="flex items-center gap-2 mb-2">
              <span class="badge"
                :class="event.status === 'published' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'"
              >
                {{ event.status === 'published' ? 'Published' : 'Draft' }}
              </span>
            </div>
            <h3 class="font-semibold text-gray-900 mb-1">{{ event.title }}</h3>
            <p class="text-sm text-gray-500">{{ event.venue }}</p>
            <p class="text-sm text-gray-400 mt-1">{{ formatDate(event.start_date) }}</p>
          </div>
        </NuxtLink>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const eventStore = useEventStore()

onMounted(() => {
  eventStore.fetchPublicEvents({ status: 'published', per_page: 6 })
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>
