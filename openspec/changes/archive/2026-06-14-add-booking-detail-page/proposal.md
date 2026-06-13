## Why

The "Tiket Saya" booking detail page at `/my/bookings/{id}` cannot be opened properly. The page either crashes on load or shows empty/incomplete data, preventing users from viewing their e-ticket codes, booking status, and event details after a successful booking.

## What Changes

- Add `GET /api/bookings/{id}` API endpoint with full booking details (event, items, e-ticket codes)
- Wire the frontend detail page at `/my/bookings/[id].vue` to fetch and display complete booking data
- Add e-ticket codes display section with proper type-safe rendering
- Add cancel booking action on the detail page
- Add loading, empty, and error states to the detail page

## Capabilities

### New Capabilities
- `booking-detail`: Display complete booking detail for a user including event info, ticket items, status, and e-ticket codes with cancel action.

### Modified Capabilities
(None)

## Impact

- **Backend**: New `GET /api/bookings/{id}` endpoint in `apps/api/internal/handler/booking.go`
- **Frontend**: Dynamic route page at `apps/web/pages/my/bookings/[id].vue`
- **Store**: `bookingStore.fetchBookingDetail` in `apps/web/stores/booking.ts`
