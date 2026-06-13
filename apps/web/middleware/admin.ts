export default defineNuxtRouteMiddleware(() => {
  // Skip during SSR — auth is restored from localStorage on client hydration
  if (import.meta.server) return

  const auth = useAuthStore()

  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }

  if (!auth.isAdmin) {
    return navigateTo('/')
  }
})
