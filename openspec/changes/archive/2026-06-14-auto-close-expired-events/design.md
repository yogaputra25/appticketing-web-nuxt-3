## Context

Events whose `end_date` has passed remain bookable because no backend handler checks `end_date`. The `finished` status exists in the model but is never auto-assigned and is only set manually by admins. Users can join war queues and reserve tickets for events that have already ended.

## Goals / Non-Goals

**Goals:**
- Prevent war join for events past `end_date` (400 Bad Request)
- Prevent booking reserve for events past `end_date` (400 Bad Request)
- Exclude finished events from public event list
- Return 404 for event detail when `end_date` has passed
- Show "Event Ended" UI on frontend event detail and war pages

**Non-Goals:**
- Automated status change (no cron job or scheduler — real-time date check is sufficient)
- Auto-refund or cancel bookings for finished events
- Admin-facing "finished events" management page

## Decisions

1. **Real-time date check instead of batch job**: Checking `end_date < now()` at every entry point (Join, Reserve, ListPublic, DetailPublic) is simpler, more reliable, and consistent with the existing `start_date` check pattern. No cron, worker, or scheduler needed.

2. **Backend check, not just frontend**: The primary gate is at the API layer. Frontend UI is secondary UX enhancement.

3. **`end_date` comparison uses `time.Now()`**: Consistent with the existing `start_date` check in `Join`. Uses `time.Now().After(event.EndDate)` for the end_date check.

4. **Reserve handler needs event loading**: Currently `Reserve` does not load the event. It will need to call `eventRepo.FindByID` to check `end_date`.

## Risks / Trade-offs

- **Time zone considerations**: `end_date` is stored as `date` (no time) and compared against server's `time.Now()`. If the server is in a different timezone than the event, the cutoff could be off by up to 24 hours. → Mitigation: Ensure server timezone matches event timezone, or store dates with timezone info.
- **Race condition**: A user who joins the war queue just before `end_date` passes could still get a session token and attempt to reserve. → Acceptable: the Reserve handler will also check `end_date` and reject if passed.
