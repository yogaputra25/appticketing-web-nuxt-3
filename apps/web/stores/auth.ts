import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface User {
  id: number
  email: string
  full_name: string
  phone?: string | null
  role: 'user' | 'admin'
  created_at: string
  updated_at: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setAuth(t: string, u: User) {
    token.value = t
    user.value = u
  }

  function clearAuth() {
    token.value = null
    user.value = null
  }

  async function login(email: string, password: string) {
    loading.value = true
    try {
      const data = await $fetch<{ user: User; token: string }>('/api/auth/login', {
        baseURL: useRuntimeConfig().public.apiBase,
        method: 'POST',
        body: { email, password },
      })
      setAuth(data.token, data.user)
      return data
    } finally {
      loading.value = false
    }
  }

  async function register(payload: { email: string; password: string; full_name: string; phone?: string }) {
    loading.value = true
    try {
      const data = await $fetch<{ user: User; token: string }>('/api/auth/register', {
        baseURL: useRuntimeConfig().public.apiBase,
        method: 'POST',
        body: payload,
      })
      setAuth(data.token, data.user)
      return data
    } finally {
      loading.value = false
    }
  }

  async function fetchMe() {
    if (!token.value) return
    try {
      const u = await $fetch<User>('/api/auth/me', {
        baseURL: useRuntimeConfig().public.apiBase,
        headers: { Authorization: `Bearer ${token.value}` },
      })
      user.value = u
    } catch {
      clearAuth()
    }
  }

  function logout() {
    clearAuth()
  }

  return { token, user, loading, isAuthenticated, isAdmin, setAuth, clearAuth, login, register, fetchMe, logout }
})
