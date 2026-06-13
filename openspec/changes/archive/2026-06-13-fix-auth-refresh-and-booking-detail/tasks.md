## 1. Fix Auth SSR Middleware

- [x] 1.1 Edit `apps/web/middleware/auth.ts` — add `if (import.meta.server) return` at the top
- [x] 1.2 Edit `apps/web/middleware/admin.ts` — add `if (import.meta.server) return` at the top
- [x] 1.3 Edit `apps/web/pages/login.vue` — add auto-redirect on mount if already authenticated (redirect to `route.query.redirect || '/'`)

## 2. Fix Booking Detail Handler

- [x] 2.1 Edit `apps/api/internal/handler/booking.go` — change `GetByID` to `GetByIDWithAssociations` in `Detail` method

## 3. Verify

- [x] 3.1 Build Go code: `go build ./...` in `apps/api/`
- [x] 3.2 Build frontend: `npm run build` (or lint) in `apps/web/`
- [x] 3.3 Manual smoke: login as admin → refresh page → stays logged in → navigate to booking detail → event title visible