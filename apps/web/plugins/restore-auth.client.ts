import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(() => {
  const auth = useAuthStore()

  const raw = localStorage.getItem('wartiket_auth')
  if (!raw) return

  try {
    const parsed = JSON.parse(raw)
    if (parsed && typeof parsed.token === 'string' && parsed.user && !auth.token) {
      auth.token = parsed.token
      auth.user = parsed.user
    }
  } catch {
    // ignore
  }
})
