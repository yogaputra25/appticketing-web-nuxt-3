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

const STORAGE_KEY = 'wartiket_auth'

interface PersistedAuth {
  token: string
  user: User
}

function readPersisted(): PersistedAuth | null {
  if (typeof window === 'undefined') return null
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw)
    if (parsed && typeof parsed.token === 'string' && parsed.user) {
      return parsed as PersistedAuth
    }
  } catch {
    // ignore parse error
  }
  return null
}

function writePersisted(token: string, user: User) {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify({ token, user }))
  } catch {
    // ignore quota / private mode errors
  }
}

function clearPersisted() {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.removeItem(STORAGE_KEY)
  } catch {
    // ignore
  }
}

export const useAuthStore = defineStore('auth', () => {
  // Hydrate from localStorage (only on client; SSR will use null then re-hydrate)
  const initial = readPersisted()
  const token = ref<string | null>(initial?.token ?? null)
  const user = ref<User | null>(initial?.user ?? null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setAuth(t: string, u: User) {
    token.value = t
    user.value = u
    writePersisted(t, u)
  }

  function clearAuth() {
    token.value = null
    user.value = null
    clearPersisted()
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
      // Re-persist in case user data changed
      if (token.value) writePersisted(token.value, u)
    } catch {
      clearAuth()
    }
  }

  function logout() {
    clearAuth()
  }

  return { token, user, loading, isAuthenticated, isAdmin, setAuth, clearAuth, login, register, fetchMe, logout }
})
