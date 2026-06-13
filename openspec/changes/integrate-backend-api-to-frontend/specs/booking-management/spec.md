## ADDED Requirements

### Requirement: Admin Booking Detail
The system SHALL allow admins to view the full detail of any booking via `GET /api/admin/bookings/{id}`, used by the admin bookings modal.

#### Scenario: Admin opens booking detail modal
- **WHEN** admin clicks "Detail" on a booking in `/admin/bookings`
- **THEN** the frontend calls `GET /api/admin/bookings/{id}` and the modal shows booking code, event, user, status, total, items, and timestamps

#### Scenario: Admin requests non-existent booking
- **WHEN** admin requests `GET /api/admin/bookings/9999` for a booking that does not exist
- **THEN** backend returns 404 with error message "booking not found"

### Requirement: Frontend My Bookings Uses Real Endpoint
The frontend `/my/bookings` page SHALL fetch from `GET /api/bookings/me` and render the response shape from `BookingRepository.GetByUser`.

#### Scenario: My bookings list renders correctly
- **WHEN** authenticated user opens `/my/bookings`
- **THEN** page shows booking_code, total_amount, status, created_at, and ticket count derived from `e_ticket_codes.length` for each booking

#### Scenario: My booking detail page shows real data
- **WHEN** authenticated user opens `/my/bookings/{id}`
- **THEN** page calls `GET /api/bookings/{id}` and renders the booking returned by the backend including `items[]`, `total_amount`, `status`, `expires_at`, and (when paid) `e_ticket_codes`

### Requirement: Cancel Booking From My Bookings
The frontend SHALL allow users to cancel their own `pending_payment` bookings from `/my/bookings/{id}` by calling `POST /api/bookings/{id}/cancel`.

#### Scenario: User cancels pending booking
- **WHEN** user clicks "Batalkan" on a `pending_payment` booking
- **THEN** frontend calls `POST /api/bookings/{id}/cancel`, shows success message, and the booking status changes to `cancelled`

#### Scenario: User cannot cancel paid booking
- **WHEN** user attempts to cancel a `paid` booking
- **THEN** backend returns 400 with message "Paid bookings cannot be self-cancelled" and frontend surfaces this error
