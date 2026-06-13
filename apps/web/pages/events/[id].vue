<template>
  <NuxtPage v-if="isChildRoute" />
  <div v-else class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="!event" class="text-center py-16">
      <p class="text-gray-500 text-lg">Event tidak ditemukan</p>
      <NuxtLink to="/events" class="text-primary-600 hover:text-primary-700 mt-2 inline-block">
        &larr; Kembali ke daftar event
      </NuxtLink>
    </div>

    <template v-else>
      <NuxtLink to="/events" class="text-sm text-gray-500 hover:text-primary-600 mb-3 md:mb-4 inline-block">
        &larr; Kembali ke daftar event
      </NuxtLink>

      <div class="h-48 md:h-72 bg-gradient-to-br from-primary-500 to-primary-700 rounded-xl flex items-center justify-center mb-6 md:mb-8">
        <span class="text-white text-5xl md:text-8xl font-bold">{{ event.title.charAt(0) }}</span>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8">
        <div class="lg:col-span-2">
          <div class="flex items-center gap-3 mb-2">
            <span
              class="badge"
              :class="event.status === 'published' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'"
            >
              {{ event.status === 'published' ? 'Published' : 'Draft' }}
            </span>
            <span v-if="isEventEnded" class="badge bg-red-100 text-red-800">Event Ended</span>
          </div>

          <h1 class="text-2xl md:text-4xl font-bold text-gray-900 mb-3 md:mb-4">{{ event.title }}</h1>

          <div class="space-y-3 text-gray-600 mb-6">
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-gray-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              {{ event.venue }}
            </div>
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-gray-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              {{ formatDate(event.start_date) }} — {{ formatDate(event.end_date) }}
            </div>
          </div>

          <div v-if="event.description" class="prose prose-gray max-w-none">
            <p>{{ event.description }}</p>
          </div>
        </div>

        <div>
          <div class="card p-5 sticky top-20 md:top-24">
            <h3 class="font-semibold text-gray-900 mb-4">Pilih Tiket</h3>

            <div v-if="!event.categories?.length" class="text-sm text-gray-500 text-center py-4">
              Belum ada kategori tiket tersedia.
            </div>

            <div v-else class="space-y-3">
              <TicketCategoryCard
                v-for="category in event.categories"
                :key="category.id"
                :category="category"
              />

              <NuxtLink
                v-if="mounted && auth.isAuthenticated && hasAvailableStock && !isEventEnded"
                :to="`/events/${event.id}/war`"
                class="btn-accent w-full mt-4 touch-target flex items-center justify-center"
              >
                War Tiket
              </NuxtLink>

              <NuxtLink
                v-else-if="!isEventEnded && (!mounted || !auth.isAuthenticated)"
                :to="`/login?redirect=/events/${event.id}`"
                class="btn-primary w-full mt-4 touch-target flex items-center justify-center"
              >
                {{ mounted ? 'Masuk untuk Membeli' : 'Memuat...' }}
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const eventStore = useEventStore()
const route = useRoute()

const mounted = ref(false)

const isChildRoute = computed(() => {
  return route.name?.toString().startsWith('events-id-') ?? false
})

const event = computed(() => eventStore.currentEvent)
const loading = ref(false)

const hasAvailableStock = computed(() => {
  return event.value?.categories?.some(c => c.available_stock > 0) ?? false
})

const isEventEnded = computed(() => {
  if (!event.value?.end_date) return false
  return new Date(event.value.end_date) < new Date()
})

async function loadEvent() {
  const id = Number(route.params.id)
  if (isNaN(id)) return
  loading.value = true
  try {
    await eventStore.fetchEventDetail(id)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  mounted.value = true
  loadEvent()
})
watch(() => route.params.id, loadEvent)

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>
