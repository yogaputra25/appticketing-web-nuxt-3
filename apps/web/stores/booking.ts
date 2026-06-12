import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface BookingItem {
  id: number
  booking_id: number
  ticket_category_id: number
  quantity: number
  unit_price: number
  subtotal: number
}

export interface Booking {
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
  const error = ref('')

  async function reserve(eventId: number, items: { ticket_category_id: number; quantity: number }[]) {
    loading.value = true
    error.value = ''
    try {
      const api = useApi()
      const data = await api.post<Booking>('/api/bookings/reserve', { event_id: eventId, items })
      return data
    } catch (err: any) {
      error.value = err?.message || 'Gagal melakukan reservasi'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchMyBookings() {
    loading.value = true
    try {
      const api = useApi()
      bookings.value = await api.get<Booking[]>('/api/bookings/me')
    } finally {
      loading.value = false
    }
  }

  async function fetchBookingDetail(id: number) {
    loading.value = true
    try {
      const api = useApi()
      currentBooking.value = await api.get<Booking>(`/api/bookings/${id}`)
      return currentBooking.value
    } finally {
      loading.value = false
    }
  }

  return { bookings, currentBooking, loading, error, reserve, fetchMyBookings, fetchBookingDetail }
})
