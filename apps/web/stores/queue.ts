import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface QueueStatus {
  position: number
  total_before: number
  is_ready: boolean
  token: string
  expires_at: string
  total_in_queue?: number
}

export const useQueueStore = defineStore('queue', () => {
  const status = ref<QueueStatus | null>(null)
  const loading = ref(false)
  const error = ref('')
  const polling = ref(false)

  async function joinWar(eventId: number) {
    loading.value = true
    error.value = ''
    try {
      const api = useApi()
      const data = await api.post<QueueStatus>('/api/war/join', { event_id: eventId })
      status.value = data
      return data
    } catch (err: any) {
      error.value = err?.message || 'Gagal bergabung ke antrian'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchStatus(eventId: number) {
    try {
      const api = useApi()
      status.value = await api.get<QueueStatus>(`/api/war/status?event_id=${eventId}`)
      return status.value
    } catch (err: any) {
      if (err?.error === 'unauthorized') {
        reset()
      }
      throw err
    }
  }

  function reset() {
    status.value = null
    loading.value = false
    error.value = ''
    polling.value = false
  }

  return { status, loading, error, polling, joinWar, fetchStatus, reset }
})
