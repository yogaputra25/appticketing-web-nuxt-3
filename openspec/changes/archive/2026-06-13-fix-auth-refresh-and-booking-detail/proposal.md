## Why

Users are unable to use the app reliably because:
1. **Login state lost on page refresh** — SSR middleware redirects to `/login` before localStorage-based auth can be restored on the client.
2. **Booking detail page shows empty event title** — the backend handler uses a repository method that does not preload the `Event` association.
3. **Add event flow broken** — frontend sends dates in `YYYY-MM-DD` format but the backend DTO uses `time.Time` which only accepts RFC3339.

Issue #3 was already fixed in a prior session; this change captures the remaining fixes.

## What Changes

- **Login page auto-redirect** — if user is already authenticated, redirect away from `/login` to home
- **Auth middleware skip SSR** — prevent 302 redirect during server-side rendering when localStorage is unavailable
- **Booking detail handler** — use `GetByIDWithAssociations` instead of `GetByID` so `Event` and `User` are preloaded

## Capabilities

### New Capabilities
- `ssr-auth-hydration`: Fix auth state persistence across SSR hydration boundary

### Modified Capabilities
*(none — no spec-level requirement changes, only implementation fixes)*

## Impact

- `apps/web/middleware/auth.ts` — skip SSR check
- `apps/web/middleware/admin.ts` — skip SSR check
- `apps/web/pages/login.vue` — add auto-redirect when already authenticated
- `apps/api/internal/handler/booking.go` — change `GetByID` to `GetByIDWithAssociations` in Detail handler
