## Why

Events whose `end_date` has passed remain in `published` status and are still accessible for war join and booking. Users can attempt to join and reserve tickets for events that have already ended, leading to confusion and invalid reservations.

## What Changes

- Backend: Reject war join if `event.end_date < time.Now()` with 400 "event has already ended"
- Backend: Reject booking reserve if `event.end_date < time.Now()` with 400 "event has already ended"
- Backend: Filter out finished events from public event list (by `end_date`, not just status)
- Backend: Return 404 for event detail if `end_date` has passed
- Backend: Auto-mark events as `finished` via a lightweight scheduler or runtime check
- Frontend: Show "Event Ended" badge/message on event detail and war pages when `end_date` has passed

## Capabilities

### New Capabilities
- `auto-close-expired-events`: Automatically detect and block access to events whose `end_date` has passed, preventing further ticket sales.

### Modified Capabilities
- `ticket-war`: War join MUST reject events past `end_date` with "event has already ended"
- `event-management`: Public event list MUST exclude finished events; event detail MUST return 404 for finished events
- `booking-management`: Booking reserve MUST reject events past `end_date` with 400

## Impact

- **Backend handlers**: `handler/war.go` (Join), `handler/booking.go` (Reserve), `handler/event.go` (ListPublic, DetailPublic), `router/router.go`
- **Backend repository**: `repository/event.go` (new auto-finish check or filter)
- **Frontend pages**: `pages/events/[id].vue` (detail), `pages/events/[id]/war.vue` (war), `components/EventCard.vue` (card badge)
