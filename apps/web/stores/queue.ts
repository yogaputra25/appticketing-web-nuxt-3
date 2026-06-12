import { defineStore } from 'pinia'
import { ref } from 'vue'

interface QueueStatus {
  position: number
  total_before: number
  is_ready: boolean
  token: string
  expires_at: string
}

export const useQueueStore = defineStore('queue', () => {
  const status = ref<QueueStatus | null>(null)
  const loading = ref(false)
  const polling = ref(false)

  async function joinWar(eventId: number) {
    loading.value = true
    try {
      const data = await $fetch<QueueStatus>('/api/war/join', {
        baseURL: useRuntimeConfig().public.apiBase,
        method: 'POST',
        body: { event_id: eventId },
      })
      status.value = data
      return data
    } finally {
      loading.value = false
    }
  }

  async function fetchStatus(eventId: number) {
    const data = await $fetch<QueueStatus>(`/api/war/status?event_id=${eventId}`, {
      baseURL: useRuntimeConfig().public.apiBase,
    })
    status.value = data
    return data
  }

  function reset() {
    status.value = null
    polling.value = false
  }

  return { status, loading, polling, joinWar, fetchStatus, reset }
})
