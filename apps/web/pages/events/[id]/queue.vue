<template>
  <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-20 md:pb-8">
    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="!event" class="text-center py-16">
      <p class="text-gray-500 text-lg">Event tidak ditemukan</p>
      <NuxtLink to="/events" class="text-primary-600 hover:text-primary-700 mt-2 inline-block">
        &larr; Kembali
      </NuxtLink>
    </div>

    <template v-else>
      <div class="text-center mb-6">
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 mb-1">Antrian War Tiket</h1>
        <p class="text-gray-500 text-sm md:text-base">{{ event.title }} — {{ event.venue }}</p>
      </div>

      <QueueStatus
        :position="queueStore.status?.position ?? 0"
        :total-in-queue="queueStore.status?.total_in_queue ?? 0"
        :error="queueError"
      />

      <div v-if="isExpired" class="card p-5 mt-4 text-center">
        <p class="text-red-600 font-medium mb-3">Sesi antrianmu telah berakhir</p>
        <NuxtLink :to="`/events/${event.id}/war`" class="btn-primary touch-target inline-flex items-center justify-center">
          Gabung Antrian Lagi
        </NuxtLink>
      </div>
    </template>

    <div class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-4 py-3">
      <button
        class="btn-outline w-full touch-target flex items-center justify-center gap-2 text-sm"
        @click="handleLeave"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
        Keluar
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const eventStore = useEventStore()
const queueStore = useQueueStore()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const queueError = ref('')
const isExpired = ref(false)
let pollTimer: ReturnType<typeof setInterval> | null = null
const retryCount = ref(0)
const maxRetries = 3

const event = computed(() => eventStore.currentEvent)

async function pollQueue() {
  try {
    const status = await queueStore.fetchStatus(Number(route.params.id))
    retryCount.value = 0

    if (status.is_ready) {
      stopPolling()
      const token = status.session_token || status.token
      router.push(`/events/${route.params.id}/booking?token=${token}`)
    }
  } catch (err: any) {
    retryCount.value++
    if (err?.error === 'unauthorized' || err?.message?.includes('expired')) {
      isExpired.value = true
      stopPolling()
    } else if (retryCount.value > maxRetries) {
      queueError.value = 'Gagal memperbarui posisi. Periksa koneksi internetmu.'
      stopPolling()
    }
  }
}

function startPolling() {
  pollTimer = setInterval(pollQueue, 2000)
  queueStore.polling = true
}

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  queueStore.polling = false
}

async function loadEvent() {
  const id = Number(route.params.id)
  if (isNaN(id)) return
  loading.value = true
  try {
    await eventStore.fetchEventDetail(id)
    startPolling()
  } finally {
    loading.value = false
  }
}

function handleLeave() {
  stopPolling()
  router.push(`/events/${route.params.id}`)
}

onMounted(loadEvent)

onUnmounted(() => {
  stopPolling()
})
</script>
