## Context

The war/booking flow has three distinct issues found during E2E testing:

1. **Countdown bypass**: `war.vue` declares an `isStarted` ref that is set to `true` when the countdown timer emits `expired`, but the ref is never read anywhere — the "Mulai War" button is always enabled and `handleJoinWar` has no guard. The backend `war/join` handler only checks `event.Status == "published"` without comparing `event.start_date` to `time.Now()`.

2. **Booking detail crash**: `UpdateETicketCodes` in `booking.go` previously passed `model.JSONStringList(codes)` to GORM's `Update()`, which serialized the value as a raw string instead of a JSON array (e.g., `TCK-...` instead of `["TCK-..."]`). Existing bookings that went through `simulate` now have corrupted `e_ticket_codes` JSONB values. When `GetByIDWithAssociations` scans them, `JSONStringList.Scan` fails because `json.Unmarshal("\"TCK-...\"", &[]string{})` returns an error. The `onMounted` handler in the frontend `[id].vue` lacks a try/catch, so the error becomes an unhandled promise rejection that blanks the page.

3. **Poor logout affordance**: The desktop navbar hides the logout button inside a dropdown that opens only on avatar click — there is no direct logout button visible, making it hard for users to find how to sign out.

## Goals / Non-Goals

**Goals:**
- Prevent war join and booking before event start time (frontend + backend enforcement)
- Enable booking detail page to load for all bookings regardless of `e_ticket_codes` data state
- Provide a clearly visible logout button in the desktop navbar for authenticated users
- Repair existing corrupted `e_ticket_codes` data so historical bookings display correctly

**Non-Goals:**
- No changes to the countdown timer component itself (it works correctly)
- No redesign of the navbar or authentication system
- No changes to the mobile layout (already has a visible logout button)
- No data migration script — repair is done inline during scan

## Decisions

### 1. Frontend countdown enforcement: `:disabled` + early return
- Add `:disabled="... || !isStarted"` to the "Mulai War" button so it stays disabled until countdown ends
- Add `if (!isStarted.value) return` at the top of `handleJoinWar` as a defensive check
- **Alternative considered**: Removing the button entirely until countdown ends — rejected because keeping it visible (disabled) gives users feedback on when they can act

### 2. Backend `start_date` validation in war/join
- Add a check after loading the event: `if (time.Now().Before(event.StartDate)) { httputil.BadRequest("event has not started yet") }`
- This is defense-in-depth: even if the frontend gate is bypassed or an API client calls directly, the backend rejects early joins
- **Alternative considered**: Rejecting with 403 Forbidden — rejected; 400 BadRequest is more descriptive for client-side validation errors

### 3. Resilient `e_ticket_codes` scan
- In `GetByIDWithAssociations`, catch the JSON scan error on `e_ticket_codes` and fall back to an empty array `[]`
- In the frontend booking detail page, add a type check: `if (typeof booking.e_ticket_codes === 'string')` → wrap in an array for display
- **Alternative considered**: Running a one-off SQL migration to fix corrupted rows — rejected because it requires manual execution; inline resilience is zero-touch

### 4. Booking detail `onMounted` error handling
- Wrap `onMounted` in try/catch and set an `error` ref to display a user-friendly message instead of crashing
- **Alternative considered**: Letting Nuxt error boundary handle it — rejected because it shows a generic error page with no recovery actions

### 5. Logout button placement
- Add a "Logout" text button next to the user avatar in the desktop navbar, styled as a subtle link
- Keep the existing dropdown logout as a secondary option
- **Alternative considered**: Moving logout into a settings/profile page — rejected because it adds friction; logout should be one click from anywhere

## Risks / Trade-offs

- **Corrupted data may still cause issues in other query paths** → The fallback is only in `GetByIDWithAssociations`; other queries (e.g., admin booking list) may still fail on corrupted rows. This is acceptable for now since the booking detail is the primary user-facing path.
- **`isStarted` ref is reactive but `handleJoinWar` is async** → The early return check happens synchronously at the top of the function, so there's no race between the check and the API call.
- **Backend `start_date` check adds a DB query** → The event is already loaded (for category/sold-out checks), so there's no additional query cost.
