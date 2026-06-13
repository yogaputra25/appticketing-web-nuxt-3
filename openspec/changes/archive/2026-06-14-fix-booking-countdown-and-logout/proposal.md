## Why

Three UX bugs prevent users from completing the booking flow and managing their account: (1) the war countdown timer is displayed but never enforced, so users can join the war before the event starts; (2) the booking detail page fails to load for paid/pending bookings because corrupted e-ticket JSONB data causes a database scan error; (3) users cannot easily log out because the logout button is hidden inside an avatar dropdown with no visual affordance.

## What Changes

- Enforce countdown timer: disable "Mulai War" button and gate `handleJoinWar` until countdown expires (`isStarted`); also validate `event.start_date` in the backend `war/join` handler to reject early joins
- **BREAKING** (DB data): Fix corrupted `e_ticket_codes` JSONB by adding a migration to repair existing records and making the scan resilient to malformed data
- Fix booking detail `onMounted` to catch API errors gracefully instead of crashing
- Add a visible logout button in the desktop navbar (instead of only inside the avatar dropdown)

## Capabilities

### New Capabilities
- `war-countdown-enforcement`: Gate war join on countdown expiry (frontend + backend), prevent booking before event start

### Modified Capabilities
*(No existing spec-level requirement changes)*

## Impact

- `apps/web/pages/events/[id]/war.vue`: wire `isStarted` into button `:disabled` and `handleJoinWar`
- `apps/web/pages/my/bookings/[id].vue`: wrap `onMounted` in try/catch
- `apps/web/pages/my/bookings/[id].vue`: handle `e_ticket_codes` display even when scan returns string
- `apps/web/layouts/default.vue`: add standalone logout button for authenticated users on desktop
- `apps/api/internal/handler/war.go`: add `event.start_date > now` check in `Join` handler
- `apps/api/internal/repository/booking.go`: make `GetByIDWithAssociations` resilient to `e_ticket_codes` scan errors
- `apps/api/internal/handler/booking.go`: add try/catch in `Detail` handler to fall back gracefully on scan errors
