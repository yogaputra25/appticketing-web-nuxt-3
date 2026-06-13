## Context

The booking detail page at `/my/bookings/{id}` is currently non-functional. Backend and frontend scaffolding exist but the end-to-end flow needs to be completed:

- Backend `GET /api/bookings/{id}` needs to return full booking data with preloaded event, items, and e-ticket codes
- Frontend `/my/bookings/[id].vue` needs proper API integration, loading states, error handling, and display of all booking data
- E-ticket codes from the JSONB column must be serialized and deserialized correctly

## Goals / Non-Goals

**Goals:**
- Fully functional booking detail page showing: booking code, event title, status, ticket count, total amount, expiration date, e-ticket codes
- Cancel booking action from the detail page
- Robust error handling (API failure, corrupted data, network errors)
- Loading and error states for UX

**Non-Goals:**
- Payment flow integration from the detail page (handled elsewhere)
- Admin booking detail view
- Download/print e-ticket feature

## Decisions

1. **Preload Event + Items on backend**: `GetByIDWithAssociations` preloads `Items`, `User`, and `Event` via GORM so the frontend gets a complete booking object in one request.
2. **Ownership check returns 404**: If the booking does not belong to the authenticated user, return 404 (not 403) to avoid leaking booking existence.
3. **Frontend uses Pinia store**: Booking data fetched through `bookingStore.fetchBookingDetail` with a reactive `currentBooking` ref consumed by the detail page.
4. **`e_ticket_codes` as JSONB with resilient scan**: Custom `JSONStringList` type handles deserialization gracefully — if the stored value is invalid JSON (e.g., a plain string instead of an array), it falls back to an empty array.

## Risks / Trade-offs

- **Corrupted e_ticket_codes in legacy data**: Fixed via resilient `JSONStringList.Scan` that catches unmarshal errors and returns empty array instead of failing.
- **Nuxt 3.21+ parent-child routing**: Dynamic route at `pages/my/bookings/[id].vue` coexists with `pages/my/bookings.vue`. Verified this works as independent routes in Nuxt 3.
