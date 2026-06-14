import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface TicketWithBooking {
  id: number
  booking_id: number
  ticket_code: string
  category_name: string
  status: string
  scanned_at?: string | null
  created_at: string
  booking_code: string
  event_id: number
  event_title: string
  event_venue: string
  start_date?: string | null
  end_date?: string | null
}

export const useTicketStore = defineStore('ticket', () => {
  const tickets = ref<TicketWithBooking[]>([])
  const currentTicket = ref<TicketWithBooking | null>(null)
  const total = ref(0)
  const loading = ref(false)
  const error = ref('')

  async function fetchMyTickets(page = 1, limit = 20) {
    loading.value = true
    try {
      const api = useApi()
      const res = await api.get<{ data: TicketWithBooking[]; total: number }>('/api/tickets', { page, limit })
      tickets.value = res.data || []
      total.value = res.total || 0
    } finally {
      loading.value = false
    }
  }

  async function fetchTicketDetail(id: number) {
    loading.value = true
    try {
      const api = useApi()
      currentTicket.value = await api.get<TicketWithBooking>(`/api/tickets/${id}`)
      return currentTicket.value
    } finally {
      loading.value = false
    }
  }

  return { tickets, currentTicket, total, loading, error, fetchMyTickets, fetchTicketDetail }
})
