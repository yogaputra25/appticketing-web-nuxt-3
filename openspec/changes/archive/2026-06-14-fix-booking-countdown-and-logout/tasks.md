## 1. Backend: War join start_date enforcement

- [x] 1.1 Add `event.start_date > time.Now()` validation in `apps/api/internal/handler/war.go` `Join` handler — reject with 400 if event has not started

## 2. Frontend: Countdown gating on war page

- [x] 2.1 Add `:disabled="... || !isStarted"` to "Mulai War" button in `apps/web/pages/events/[id]/war.vue`
- [x] 2.2 Add `if (!isStarted.value) return` at the top of `handleJoinWar` in `apps/web/pages/events/[id]/war.vue`

## 3. Frontend: Booking page start_date guard

- [x] 3.1 Add `eventNotStarted` computed in `apps/web/pages/events/[id]/booking.vue` that checks if `event.start_date` is in the future
- [x] 3.2 Show "Event belum dimulai" message and hide booking form when `eventNotStarted` is true

## 4. Backend: Resilient e_ticket_codes scan

- [x] 4.1 In `apps/api/internal/repository/booking.go` `GetByIDWithAssociations`, catch `e_ticket_codes` scan error and reset to empty `JSONStringList`

## 5. Frontend: Booking detail error handling

- [x] 5.1 Wrap `onMounted` in `apps/web/pages/my/bookings/[id].vue` with try/catch and set a local `error` ref
- [x] 5.2 Add `v-if="error"` error display block with "Gagal memuat data booking" message
- [x] 5.3 Add type guard in template for `e_ticket_codes` — if string, wrap in array for display

## 6. Frontend: Logout button visibility

- [x] 6.1 Add a visible "Logout" button/link next to the user avatar in `apps/web/layouts/default.vue` desktop navbar for authenticated users
