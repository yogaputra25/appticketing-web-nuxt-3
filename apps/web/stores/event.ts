import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface TicketCategory {
  id: number
  event_id: number
  name: string
  description?: string | null
  price: number
  total_stock: number
  available_stock: number
  max_per_user: number
  sale_start_at?: string | null
  sale_end_at?: string | null
  created_at: string
  updated_at: string
}

export interface Event {
  id: number
  title: string
  description?: string | null
  venue: string
  start_date: string
  end_date: string
  banner_url?: string | null
  status: string
  categories?: TicketCategory[]
  created_at: string
  updated_at: string
}

export interface PaginatedResponse<T> {
  data: T[]
  page: number
  per_page: number
  total: number
  total_pages: number
}

export const useEventStore = defineStore('event', () => {
  const events = ref<Event[]>([])
  const currentEvent = ref<Event | null>(null)
  const loading = ref(false)
  const pagination = ref({ page: 1, per_page: 12, total: 0, total_pages: 0 })

  async function fetchPublicEvents(params?: { page?: number; per_page?: number; search?: string; status?: string }) {
    loading.value = true
    try {
      const data = await $fetch<PaginatedResponse<Event>>('/api/events', {
        baseURL: useRuntimeConfig().public.apiBase,
        params,
      })
      events.value = data.data
      pagination.value = { page: data.page, per_page: data.per_page, total: data.total, total_pages: data.total_pages }
    } finally {
      loading.value = false
    }
  }

  async function fetchEventDetail(id: number) {
    loading.value = true
    try {
      const data = await $fetch<Event>(`/api/events/${id}`, {
        baseURL: useRuntimeConfig().public.apiBase,
      })
      currentEvent.value = data
      return data
    } finally {
      loading.value = false
    }
  }

  return { events, currentEvent, loading, pagination, fetchPublicEvents, fetchEventDetail }
})
