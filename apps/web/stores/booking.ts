import { defineStore } from 'pinia'
import { ref } from 'vue'

interface BookingItem {
  id: number
  booking_id: number
  ticket_category_id: number
  quantity: number
  unit_price: number
  subtotal: number
}

interface Booking {
  id: number
  booking_code: string
  user_id: number
  event_id: number
  total_amount: number
  status: string
  expires_at?: string | null
  e_ticket_codes: string[]
  items?: BookingItem[]
  created_at: string
  updated_at: string
}

export const useBookingStore = defineStore('booking', () => {
  const bookings = ref<Booking[]>([])
  const currentBooking = ref<Booking | null>(null)
  const loading = ref(false)

  async function reserve(eventId: number, items: { ticket_category_id: number; quantity: number }[]) {
    loading.value = true
    try {
      const data = await $fetch<Booking>('/api/bookings/reserve', {
        baseURL: useRuntimeConfig().public.apiBase,
        method: 'POST',
        body: { event_id: eventId, items },
      })
      return data
    } finally {
      loading.value = false
    }
  }

  async function fetchMyBookings() {
    loading.value = true
    try {
      const data = await $fetch<Booking[]>('/api/bookings/me', {
        baseURL: useRuntimeConfig().public.apiBase,
      })
      bookings.value = data
    } finally {
      loading.value = false
    }
  }

  async function fetchBookingDetail(id: number) {
    loading.value = true
    try {
      const data = await $fetch<Booking>(`/api/bookings/${id}`, {
        baseURL: useRuntimeConfig().public.apiBase,
      })
      currentBooking.value = data
      return data
    } finally {
      loading.value = false
    }
  }

  return { bookings, currentBooking, loading, reserve, fetchMyBookings, fetchBookingDetail }
})
