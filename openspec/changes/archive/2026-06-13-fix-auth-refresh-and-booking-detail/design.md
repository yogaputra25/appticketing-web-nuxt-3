## Context

The app uses Nuxt 3 with SSR. Auth state is persisted to `localStorage` via Pinia store and restored on boot by `plugins/auth.client.ts`. However, route middleware (`auth.ts`, `admin.ts`) runs on the **server** during SSR, where `localStorage` is unavailable, so `auth.isAuthenticated` is always `false` → 302 redirect to `/login` happens before client hydration can restore the token.

Additionally, the booking detail handler uses a repository method (`GetByID`) that does not preload `Event` and `User` associations, leaving those fields empty in the response.

## Goals / Non-Goals

**Goals:**
- Auth state survives page refresh (no redirect to login)
- Booking detail page shows event title and user info
- Login page redirects authenticated users away

**Non-Goals:**
- Replacing localStorage with cookies for SSR auth
- Adding new auth middleware or features
- Refactoring Nuxt SSR mode

## Decisions

### 1. Skip SSR in middleware (not `ssr: false`)

Adding `if (import.meta.server) return` at the top of `auth.ts` and `admin.ts`.

**Why:**
- Minimal change — 1 line per middleware file
- Route guard still works on client-side navigation
- Auth plugin (`auth.client.ts`) already restores token on boot
- No need to disable SSR entirely (would lose SEO benefits for public pages)

**Alternatives considered:**
- Disable SSR globally (`ssr: false` in nuxt.config) → loses SEO for public event pages
- Use `useCookie` instead of localStorage → requires more changes and cookie-based auth flows
- Pinia persistence plugin → adds dependency, overkill for this fix

### 2. Login page auto-redirect

Check `auth.isAuthenticated` on mount and redirect to `route.query.redirect || '/'`.

**Why:**
- Without this, after SSR fix, an authenticated user who navigates to `/login` stays on the form
- Completes the UX flow: refresh protected page → skip SSR → hydrate → still on `/login` → redirect automatically

**Location:** `onMounted` in `login.vue`.

### 3. Booking detail: change repository method

Replace `h.bookings.GetByID(ctx, id)` with `h.bookings.GetByIDWithAssociations(ctx, id)` in `BookingHandler.Detail`.

**Why:**
- `GetByID` only preloads `Items`
- `GetByIDWithAssociations` preloads `Items`, `User`, and `Event`
- No new queries or performance concern — existing method already exists in the repository
- Frontend templates already access `booking.event?.title` — this fix makes it work

## Risks / Trade-offs

- **Risk:** Skipping SSR middleware means protected pages are briefly rendered in "logged out" state during SSR, then flash to "logged in" on hydration
  → **Mitigation:** The auth plugin restores state synchronously from localStorage before mounting, so the flash is imperceptible
- **Risk:** Login page auto-redirect could cause redirect loops if the auth token is invalid
  → **Mitigation:** `fetchMe()` in the auth plugin validates the token on boot; if invalid, `clearAuth()` runs and `isAuthenticated` becomes `false`, so no redirect
- **Trade-off:** SSR protection for protected pages is sacrificed — unauthenticated users could briefly see a shell of protected pages during SSR
  → Acceptable: admin pages have no SEO value, and the content is fetched client-side after hydration anyway

## Migration Plan

1. Edit `apps/web/middleware/auth.ts` — add `if (import.meta.server) return`
2. Edit `apps/web/middleware/admin.ts` — add `if (import.meta.server) return`
3. Edit `apps/web/pages/login.vue` — add auto-redirect in `onMounted`
4. Edit `apps/api/internal/handler/booking.go` — change `GetByID` → `GetByIDWithAssociations`
5. Build & deploy

**Rollback:** Revert commits. Zero data migration needed.
