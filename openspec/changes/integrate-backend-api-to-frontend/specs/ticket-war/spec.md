## MODIFIED Requirements

### Requirement: Join War Queue
The system SHALL place authenticated users into a virtual queue when they attempt to purchase tickets from a sold-out or high-demand event.

#### Scenario: User joins queue for sold-out event
- **WHEN** user clicks "War Tiket" on an event whose categories are all sold out
- **THEN** backend adds user to Redis sorted set queue, assigns a queue token, and returns `{ queued: true, position, token }` (HTTP 200)

#### Scenario: User joins queue for in-stock event
- **WHEN** user clicks "War Tiket" on an event with available stock
- **THEN** backend returns `{ redirect_to_booking: true, message: "Tickets are available. Proceed to booking." }` and frontend routes user directly to `/events/{id}/booking?token=<session_token>` (token issued by queue process)

#### Scenario: Duplicate queue entry
- **WHEN** user already has an active queue token
- **THEN** system returns existing token and position (no duplicate entry)

### Requirement: Queue Position Polling
The system SHALL allow queued users to poll their position and estimated wait time via `GET /api/war/status?event_id={id}`.

#### Scenario: User in queue position 50
- **WHEN** user polls status with valid queue token
- **THEN** backend returns `{ is_ready: false, position, total_in_queue, estimated_wait }` (HTTP 200)

#### Scenario: User's turn arrives
- **WHEN** user reaches position 0 (or threshold)
- **THEN** backend returns `{ is_ready: true, session_token, message: "Your turn! Proceed to booking." }` — the frontend stores `session_token` for the reservation request

#### Scenario: Queue token expired
- **WHEN** user polls with expired or invalid token
- **THEN** backend returns 401 with error and frontend shows "Sesi antrian berakhir" with a button to rejoin

#### Scenario: User not in queue
- **WHEN** user polls status without ever joining the queue
- **THEN** backend returns `{ is_ready: false, queued: false, message: "You are not in the queue for this event." }`

### Requirement: Booking Session Token
The system SHALL issue a short-lived (5 minute) booking session token to users when they are ready to book, limiting the time they can hold stock.

#### Scenario: User has active booking session
- **WHEN** user has a valid session token from `/api/war/status` (or from a `redirect_to_booking` flow)
- **THEN** they can POST to `/api/bookings/reserve` with body `{ event_id, session_token, items: [{ category_id, quantity }] }`

#### Scenario: Booking session expired
- **WHEN** user attempts reservation with expired session token
- **THEN** backend returns 401 and frontend redirects them to rejoin the queue

## ADDED Requirements

### Requirement: Frontend Polling Loop Uses Real Status Endpoint
The `pages/events/[id]/queue.vue` page SHALL poll `GET /api/war/status?event_id={id}` every 2 seconds and react to `is_ready` and `position` correctly.

#### Scenario: Polling continues while waiting
- **WHEN** the queue page is mounted
- **THEN** it starts a 2-second `setInterval` calling `/api/war/status` and updates the position display

#### Scenario: Polling stops on ready
- **WHEN** the status response has `is_ready: true`
- **THEN** the polling stops and the page redirects to `/events/{id}/booking` carrying the `session_token` as a query parameter

#### Scenario: Polling stops on expired
- **WHEN** the status request returns 401 or error indicating expired token
- **THEN** the polling stops and the page shows "Sesi antrian berakhir" with a "Gabung Antrian Lagi" button

### Requirement: Frontend Reservation Sends session_token
The `pages/events/[id]/booking.vue` page SHALL send the `session_token` (from queue status) as part of the `POST /api/bookings/reserve` body so the backend can validate the booking session.

#### Scenario: Booking submitted with valid session
- **WHEN** user has a `session_token` query parameter (from queue or direct redirect) and clicks "Lanjut Bayar"
- **THEN** frontend calls `POST /api/bookings/reserve` with `{ event_id, session_token, items: [{ category_id, quantity }] }`

#### Scenario: Session token is required
- **WHEN** user opens the booking page without a `token` query parameter
- **THEN** the page shows "Sesi booking tidak valid" and a button to start the queue again
