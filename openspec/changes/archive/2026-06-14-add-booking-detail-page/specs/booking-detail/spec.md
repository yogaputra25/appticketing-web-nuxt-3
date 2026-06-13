## ADDED Requirements

### Requirement: User can view booking detail
The system SHALL provide a booking detail page accessible at `/my/bookings/{id}` for authenticated users.

#### Scenario: View booking detail
- **WHEN** user navigates to `/my/bookings/{id}`
- **THEN** system fetches booking data from `GET /api/bookings/{id}`
- **THEN** system displays booking code, event title, status, ticket count, total amount, and expiration date

#### Scenario: Booking not found
- **WHEN** user navigates to `/my/bookings/{id}` with a non-existent or non-owned booking ID
- **THEN** system shows "Booking tidak ditemukan" error message
- **THEN** user can navigate back to the bookings list

#### Scenario: API failure
- **WHEN** the API request fails with a server error
- **THEN** system shows "Gagal memuat detail booking" error message with details

#### Scenario: Loading state
- **WHEN** the API request is in progress
- **THEN** system shows a loading spinner

### Requirement: User can view e-ticket codes
The system SHALL display e-ticket codes for paid bookings.

#### Scenario: Display e-ticket codes
- **WHEN** booking status is "paid"
- **THEN** system shows each e-ticket code in a monospace font, separated per ticket

#### Scenario: No e-ticket codes available
- **WHEN** booking has no e-ticket codes (e.g., pending or cancelled)
- **THEN** system hides the e-ticket codes section

### Requirement: User can cancel a pending booking
The system SHALL allow users to cancel their own pending bookings.

#### Scenario: Cancel pending booking
- **WHEN** user clicks "Batalkan Pesanan" on a pending booking
- **THEN** system calls `POST /api/bookings/{id}/cancel`
- **THEN** system updates the page to show the booking as cancelled

#### Scenario: Cancel non-pending booking
- **WHEN** user tries to cancel a paid or already cancelled booking
- **THEN** system shows "Paid bookings cannot be self-cancelled" error message

### Requirement: Booking detail API returns complete data
The `GET /api/bookings/{id}` endpoint SHALL return a booking object with preloaded event, items, and e-ticket codes.

#### Scenario: Full response structure
- **WHEN** authenticated user requests their own booking
- **THEN** response includes: id, booking_code, event (with title, venue, dates), items (with category, quantity, price), total_amount, status, expires_at, e_ticket_codes (string array), created_at

#### Scenario: Ownership check
- **WHEN** authenticated user requests another user's booking
- **THEN** system returns 404 "booking not found"
