<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8 pb-24 md:pb-8">
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
      <div class="flex items-center justify-between mb-4 md:mb-6">
        <div>
          <h1 class="text-xl md:text-2xl font-bold text-gray-900">Booking Tiket</h1>
          <p class="text-gray-500 text-sm md:text-base">{{ event.title }} — {{ event.venue }}</p>
        </div>
        <NuxtLink v-if="hasActiveSession" :to="`/events/${event.id}`" class="text-sm text-primary-600 hover:text-primary-700 shrink-0">
          Detail Event &rarr;
        </NuxtLink>
      </div>

      <div v-if="hasActiveSession" class="card p-3 md:p-4 mb-4 md:mb-6 flex items-center justify-between">
        <div class="flex items-center gap-2 text-sm text-gray-600">
          <svg class="w-5 h-5 text-yellow-500 shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
          </svg>
          <span class="text-xs md:text-sm">Sesi booking akan berakhir dalam <strong class="text-gray-900 tabular-nums">
            <CountdownTimer :target-date="sessionExpiresAt" @expired="sessionExpired = true" />
          </strong></span>
        </div>
      </div>

      <div v-if="sessionExpired" class="card p-5 mb-6 text-center">
        <p class="text-red-600 font-medium mb-3">Sesi booking telah berakhir</p>
        <NuxtLink :to="`/events/${event.id}/war`" class="btn-primary touch-target inline-flex items-center justify-center">
          Mulai War Lagi
        </NuxtLink>
      </div>

      <template v-if="!sessionExpired && hasActiveSession">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div class="lg:col-span-2 space-y-3">
            <h3 class="font-semibold text-gray-900">Pilih Jumlah Tiket</h3>

            <div v-for="cat in event.categories" :key="cat.id" class="card p-4">
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <div>
                  <h4 class="font-semibold text-gray-900">{{ cat.name }}</h4>
                  <p class="text-sm text-gray-500">{{ formatPrice(cat.price) }} / tiket</p>
                  <p v-if="cat.available_stock <= 5 && cat.available_stock > 0" class="text-xs text-red-500 mt-0.5">
                    Sisa {{ cat.available_stock }} tiket!
                  </p>
                </div>
                <div class="flex items-center gap-3">
                  <button
                    class="w-10 h-10 md:w-9 md:h-9 rounded-full border border-gray-300 flex items-center justify-center text-gray-600 hover:bg-gray-50 disabled:opacity-30 touch-target"
                    :disabled="(quantities[cat.id] || 0) === 0"
                    @click="decrement(cat.id)"
                  >
                    &minus;
                  </button>
                  <span class="w-8 md:w-10 text-center font-semibold tabular-nums text-base">{{ quantities[cat.id] || 0 }}</span>
                  <button
                    class="w-10 h-10 md:w-9 md:h-9 rounded-full border border-gray-300 flex items-center justify-center text-gray-600 hover:bg-gray-50 disabled:opacity-30 touch-target"
                    :disabled="(quantities[cat.id] || 0) >= Math.min(cat.max_per_user, cat.available_stock)"
                    @click="increment(cat.id)"
                  >
                    +
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="hidden md:block">
            <div class="card p-5 sticky top-24">
              <h3 class="font-semibold text-gray-900 mb-4">Ringkasan</h3>

              <div v-if="totalItems === 0" class="text-sm text-gray-500 text-center py-4">
                Belum ada tiket dipilih
              </div>

              <template v-else>
                <div class="space-y-2 mb-4">
                  <div
                    v-for="item in selectedItems"
                    :key="item.categoryId"
                    class="flex justify-between text-sm"
                  >
                    <span class="text-gray-600">{{ item.name }} x{{ item.qty }}</span>
                    <span class="font-medium text-gray-900">{{ formatPrice(item.subtotal) }}</span>
                  </div>
                </div>

                <hr class="my-3">

                <div class="flex justify-between font-semibold text-gray-900 mb-4">
                  <span>Total</span>
                  <span>{{ formatPrice(totalPrice) }}</span>
                </div>

                <div v-if="bookingError" class="text-sm text-red-600 bg-red-50 rounded-lg p-3 mb-3">
                  {{ bookingError }}
                </div>

                <button
                  class="btn-accent w-full touch-target"
                  :disabled="bookingStore.loading || totalItems === 0"
                  @click="handleReserve"
                >
                  <svg v-if="bookingStore.loading" class="animate-spin -ml-1 mr-2 h-4 w-4 inline" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  {{ bookingStore.loading ? 'Memproses...' : 'Lanjut Bayar' }}
                </button>
              </template>
            </div>
          </div>
        </div>

        <StickyBottomBar :visible="totalItems > 0">
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-xs text-gray-500">Total</p>
              <p class="text-lg font-bold text-gray-900">{{ formatPrice(totalPrice) }}</p>
            </div>
            <button
              class="btn-accent flex-1 touch-target"
              :disabled="bookingStore.loading || totalItems === 0"
              @click="handleReserve"
            >
              {{ bookingStore.loading ? 'Memproses...' : 'Lanjut Bayar' }}
            </button>
          </div>
        </StickyBottomBar>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const eventStore = useEventStore()
const bookingStore = useBookingStore()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const sessionExpired = ref(false)
const bookingError = ref('')
const quantities = ref<Record<number, number>>({})

const event = computed(() => eventStore.currentEvent)
const sessionToken = computed(() => (route.query.token as string) || '')
const hasActiveSession = computed(() => !!sessionToken.value && !sessionExpired.value)

const sessionExpiresAt = computed(() => {
  if (!event.value?.start_date) return new Date(Date.now() + 5 * 60000).toISOString()
  return new Date(Date.now() + 5 * 60000).toISOString()
})

const selectedItems = computed(() => {
  if (!event.value?.categories) return []
  return event.value.categories
    .filter(cat => (quantities.value[cat.id] || 0) > 0)
    .map(cat => ({
      categoryId: cat.id,
      name: cat.name,
      qty: quantities.value[cat.id] || 0,
      subtotal: (quantities.value[cat.id] || 0) * cat.price,
    }))
})

const totalItems = computed(() => selectedItems.value.reduce((sum, i) => sum + i.qty, 0))

const totalPrice = computed(() => selectedItems.value.reduce((sum, i) => sum + i.subtotal, 0))

function increment(categoryId: number) {
  quantities.value[categoryId] = (quantities.value[categoryId] || 0) + 1
}

function decrement(categoryId: number) {
  const current = quantities.value[categoryId] || 0
  if (current > 0) {
    quantities.value[categoryId] = current - 1
  }
}

async function handleReserve() {
  bookingError.value = ''
  const items = selectedItems.value.map(i => ({
    category_id: i.categoryId,
    quantity: i.qty,
  }))

  try {
    const booking = await bookingStore.reserve(Number(route.params.id), items, sessionToken.value)
    router.push(`/bookings/${booking.id}/pay`)
  } catch (err: any) {
    bookingError.value = err?.message || 'Gagal melakukan reservasi'
  }
}

async function loadEvent() {
  const id = Number(route.params.id)
  if (isNaN(id)) return
  loading.value = true
  try {
    await eventStore.fetchEventDetail(id)
    if (event.value?.categories) {
      event.value.categories.forEach(cat => {
        quantities.value[cat.id] = 0
      })
    }
  } finally {
    loading.value = false
  }
}

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price)
}

onMounted(loadEvent)
</script>
