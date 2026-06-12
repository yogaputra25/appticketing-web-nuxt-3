import { useAuthStore } from '~/stores/auth'

interface ApiError {
  error: string
  message?: string
  fields?: Record<string, string>
}

export function useApi() {
  const config = useRuntimeConfig()
  const auth = useAuthStore()

  async function request<T>(url: string, opts: Record<string, any> = {}): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(opts.headers || {}),
    }

    if (auth.token) {
      headers['Authorization'] = `Bearer ${auth.token}`
    }

    try {
      const data = await $fetch<T>(url, {
        baseURL: config.public.apiBase,
        ...opts,
        headers,
      })
      return data
    } catch (err: any) {
      if (err?.response?.status === 401) {
        auth.clearAuth()
        navigateTo('/login')
      }
      const apiError: ApiError = err?.data || { error: 'unknown', message: err?.message }
      throw apiError
    }
  }

  return {
    get: <T>(url: string, params?: Record<string, any>) =>
      request<T>(url, { method: 'GET', params }),

    post: <T>(url: string, body?: any) =>
      request<T>(url, { method: 'POST', body }),

    put: <T>(url: string, body?: any) =>
      request<T>(url, { method: 'PUT', body }),

    delete: <T>(url: string) =>
      request<T>(url, { method: 'DELETE' }),
  }
}
