<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
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
      <NuxtLink :to="`/events/${event.id}`" class="text-sm text-gray-500 hover:text-primary-600 mb-3 md:mb-4 inline-block">
        &larr; Kembali ke detail event
      </NuxtLink>

      <div class="card overflow-hidden">
        <div v-if="isEventEnded" class="h-48 bg-gradient-to-br from-gray-400 to-gray-600 flex items-center justify-center">
          <div class="text-center text-white px-4">
            <p class="text-2xl sm:text-3xl md:text-4xl font-bold">Event Ended</p>
            <p class="text-sm md:text-base font-medium opacity-80 mt-2">Tiket sudah tidak tersedia</p>
          </div>
        </div>
        <div v-else class="h-48 bg-gradient-to-br from-accent to-primary-700 flex items-center justify-center">
          <div class="text-center text-white px-4">
            <p class="text-sm md:text-base font-medium opacity-80 mb-1">War Tiket Dimulai dalam</p>
            <p class="text-4xl sm:text-5xl md:text-6xl font-bold tabular-nums leading-tight">
              <CountdownTimer
                :target-date="event.start_date"
                @expired="isStarted = true"
              />
            </p>
          </div>
        </div>

        <div class="p-5 md:p-6 text-center">
          <h1 class="text-xl md:text-2xl font-bold text-gray-900 mb-2">{{ event.title }}</h1>
          <p class="text-gray-500 mb-1 text-sm md:text-base">{{ event.venue }}</p>
          <p class="text-sm text-gray-400 mb-6">
            {{ formatDate(event.start_date) }} — {{ formatDate(event.end_date) }}
          </p>

          <div v-if="joiningError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3 mb-4">
            {{ joiningError }}
          </div>

          <button
            v-if="!isEventEnded"
            class="btn-accent w-full sm:w-auto text-base md:text-lg px-8 md:px-12 !min-h-[56px]"
            :disabled="queueStore.loading || joiningInProgress || !isStarted"
            @click="handleJoinWar"
          >
            <svg v-if="queueStore.loading" class="animate-spin -ml-1 mr-2 h-5 w-5 inline" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            {{ queueStore.loading ? 'Memproses...' : joiningInProgress ? 'Mengarahkan...' : 'Mulai War' }}
          </button>

          <p v-if="!isEventEnded" class="text-sm text-gray-400 mt-4">
            Pastikan koneksi internetmu stabil. Sistem akan menempatkanmu dalam antrian secara adil.
          </p>
        </div>
      </div>

      <div class="card p-5 mt-6">
        <h3 class="font-semibold text-gray-900 mb-3">Cara Kerja War Tiket</h3>
        <ol class="space-y-2 text-sm text-gray-600 list-decimal list-inside leading-relaxed">
          <li>Klik tombol <strong>"Mulai War"</strong> saat countdown selesai</li>
          <li>Kamu akan mendapat posisi antrian</li>
          <li>Halaman akan otomatis memperbarui posisimu setiap 2 detik</li>
          <li>Saat giliran tiba, kamu akan diarahkan ke halaman booking</li>
          <li>Selesaikan pemesanan dalam waktu yang tersedia</li>
        </ol>
      </div>
    </template>
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

const isStarted = ref(false)
const joiningInProgress = ref(false)
const joiningError = ref('')

const event = computed(() => eventStore.currentEvent)
const loading = ref(false)

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

async function handleJoinWar() {
  if (!isStarted.value) return
  joiningError.value = ''
  joiningInProgress.value = true
  try {
    const result = await queueStore.joinWar(Number(route.params.id))
    queueStore.reset()

    if ((result as any).redirect_to_booking) {
      router.push(`/events/${route.params.id}/booking?token=${(result as any).session_token}`)
    } else if (result.is_ready || result.position === 0) {
      router.push(`/events/${route.params.id}/booking?token=${result.token}`)
    } else {
      router.push(`/events/${route.params.id}/queue?token=${result.token}`)
    }
  } catch (err: any) {
    if (err?.message) {
      joiningError.value = err.message
    } else {
      joiningError.value = 'Gagal bergabung ke antrian. Silakan coba lagi.'
    }
    joiningInProgress.value = false
  }
}

onMounted(loadEvent)

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>
