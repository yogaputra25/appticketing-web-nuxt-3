## 1. Backend: Booking detail API endpoint

- [x] 1.1 Add `Detail` handler in `apps/api/internal/handler/booking.go` that parses booking ID from URL, calls `GetByIDWithAssociations`, performs ownership check, returns full booking JSON
- [x] 1.2 Register `GET /bookings/{id}` route in `apps/api/internal/router/router.go` under the authenticated `/bookings` group
- [x] 1.3 Implement `GetByIDWithAssociations` in `apps/api/internal/repository/booking.go` with Preload for Items, User, and Event
- [x] 1.4 Add resilient `e_ticket_codes` JSONB scan — catch unmarshal error and fall back to empty array in `JSONStringList.Scan`

## 2. Frontend: Booking detail store

- [x] 2.1 Add `fetchBookingDetail(id)` method to `apps/web/stores/booking.ts` that calls `GET /api/bookings/{id}` and sets `currentBooking`
- [x] 2.2 Add `cancelBooking(id)` method that calls `POST /api/bookings/{id}/cancel`

## 3. Frontend: Booking detail page

- [x] 3.1 Create dynamic route page at `apps/web/pages/my/bookings/[id].vue` with `onMounted` calling `fetchBookingDetail`
- [x] 3.2 Display booking info: booking code, event title, status badge, ticket count, total amount, expiration date
- [x] 3.3 Display e-ticket codes section with monospace font, hidden when empty
- [x] 3.4 Add cancel booking button with confirmation for pending bookings
- [x] 3.5 Add loading spinner, error display with "Gagal memuat detail booking", and "Booking tidak ditemukan" states
- [x] 3.6 Wrap `onMounted` fetch in try/catch and set local `error` ref on failure

## 4. Frontend: Navigation from list to detail

- [x] 4.1 Ensure `apps/web/pages/my/bookings/index.vue` links to `/my/bookings/{booking.id}` with NuxtLink
- [x] 4.2 Fix Nuxt routing conflict: renamed `bookings.vue` to `bookings/index.vue` so `bookings/[id].vue` route registers properly
