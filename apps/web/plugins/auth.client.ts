/**
 * Restore auth state on app boot (client only).
 * - The Pinia store is hydrated from localStorage in its initial state
 *   (see `stores/auth.ts`), so token/user are already available synchronously
 *   on first render.
 * - This plugin additionally calls `fetchMe()` to validate the token and
 *   refresh user data from the backend.
 */
export default defineNuxtPlugin(async () => {
  const auth = useAuthStore()
  if (auth.isAuthenticated) {
    // Best-effort: refresh user info; ignore failures (token might be expired
    // and the user will be redirected by useApi on next protected request).
    try {
      await auth.fetchMe()
    } catch {
      // fetchMe already clears auth on failure
    }
  }
})
