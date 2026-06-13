## 1. Backend: War join end_date check

- [x] 1.1 Add `if time.Now().After(event.EndDate)` check in `apps/api/internal/handler/war.go` `Join` handler after start_date check — return 400 "event has already ended"

## 2. Backend: Booking reserve end_date check

- [x] 2.1 In `apps/api/internal/handler/booking.go` `Reserve` handler, load the event via `eventRepo.FindByID` and reject with 400 if `time.Now().After(event.EndDate)`
- [x] 2.2 Add `EventRepository` dependency to `BookingHandler` struct and constructor

## 3. Backend: Event list/detail end_date filter

- [x] 3.1 In `apps/api/internal/handler/event.go` `ListPublic`, add `end_date >= now()` filter to exclude finished events from public listing
- [x] 3.2 In `apps/api/internal/handler/event.go` `DetailPublic`, return 404 if `time.Now().After(event.EndDate)`

## 4. Frontend: Event ended UI

- [x] 4.1 In `apps/web/pages/events/[id].vue`, add "Event Ended" badge/message when `new Date(event.end_date) < new Date()` and hide the "War Tiket" button
- [x] 4.2 In `apps/web/pages/events/[id]/war.vue`, add "Event Ended" message when `new Date(event.end_date) < new Date()` and hide the join button and countdown
