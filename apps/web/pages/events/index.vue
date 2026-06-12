<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
    <div class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-gray-900">Events</h1>
      <p class="text-gray-500 mt-1 text-sm md:text-base">Temukan konser dan event favoritmu</p>
    </div>

    <div class="flex flex-col sm:flex-row gap-3 md:gap-4 mb-6 md:mb-8">
      <div class="flex-1 relative">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="search"
          type="text"
          class="input pl-10 !h-[44px]"
          placeholder="Cari event..."
          @input="onSearch"
        />
      </div>
      <select v-model="statusFilter" class="input sm:w-48 !h-[44px]" @change="fetchEvents(1)">
        <option value="">Semua Status</option>
        <option value="published">Published</option>
        <option value="draft">Draft</option>
        <option value="finished">Finished</option>
      </select>
    </div>

    <div v-if="eventStore.loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <template v-else-if="eventStore.events.length === 0">
      <div class="text-center py-16">
        <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <p class="text-gray-500 text-lg">Tidak ada event ditemukan</p>
        <p class="text-gray-400 text-sm mt-1">Coba ubah filter atau kata kunci pencarian</p>
      </div>
    </template>

    <template v-else>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6">
        <EventCard v-for="event in eventStore.events" :key="event.id" :event="event" />
      </div>

      <div v-if="eventStore.pagination.total_pages > 1" class="flex justify-center items-center gap-2 mt-8 md:mt-10">
        <button
          :disabled="currentPage <= 1"
          class="btn-outline !py-2 !px-3 touch-target text-sm"
          @click="fetchEvents(currentPage - 1)"
        >
          &larr; Sebelumnya
        </button>
        <span class="text-sm text-gray-500">
          Halaman {{ currentPage }} dari {{ eventStore.pagination.total_pages }}
        </span>
        <button
          :disabled="currentPage >= eventStore.pagination.total_pages"
          class="btn-outline !py-2 !px-3 touch-target text-sm"
          @click="fetchEvents(currentPage + 1)"
        >
          Selanjutnya &rarr;
        </button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
const eventStore = useEventStore()
const search = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
let searchTimeout: ReturnType<typeof setTimeout> | null = null

function fetchEvents(page: number = 1) {
  currentPage.value = page
  const params: Record<string, any> = { page }
  if (statusFilter.value) params.status = statusFilter.value
  if (search.value.trim()) params.search = search.value.trim()
  eventStore.fetchPublicEvents(params)
}

function onSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => fetchEvents(1), 400)
}

onMounted(() => {
  fetchEvents(1)
})
</script>
