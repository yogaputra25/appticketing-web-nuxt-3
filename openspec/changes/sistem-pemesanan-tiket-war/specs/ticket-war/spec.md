## ADDED Requirements

### Requirement: Join War Queue
The system SHALL place authenticated users into a virtual queue when they attempt to purchase tickets from a sold-out or high-demand event.

#### Scenario: User joins queue for sold-out event
- **WHEN** user clicks "War Tiket" on an event whose categories are all sold out
- **THEN** system adds user to Redis sorted set queue, assigns a queue token, and returns current position

#### Scenario: User joins queue for in-stock event
- **WHEN** user clicks "War Tiket" on an event with available stock
- **THEN** system does NOT queue and redirects user directly to the booking flow

#### Scenario: Duplicate queue entry
- **WHEN** user already has an active queue token
- **THEN** system returns existing token and position (no duplicate entry)

### Requirement: Queue Position Polling
The system SHALL allow queued users to poll their position and estimated wait time.

#### Scenario: User in queue position 50
- **WHEN** user polls `/api/queue/status` with valid token
- **THEN** system returns current position, total in queue, and estimated wait time

#### Scenario: User's turn arrives
- **WHEN** user reaches position 0 (or threshold)
- **THEN** system returns `is_ready: true` and a 5-minute booking session token

#### Scenario: Queue token expired
- **WHEN** user polls with expired or invalid token
- **THEN** system returns 401 with message "Queue session expired, please rejoin"

### Requirement: Rate Limiting on War Endpoint
The system SHALL rate-limit the war/queue join endpoint to prevent abuse.

#### Scenario: User exceeds join attempts
- **WHEN** user makes more than 5 join requests within 60 seconds
- **THEN** system returns 429 Too Many Requests

#### Scenario: Bot detection via CAPTCHA
- **WHEN** traffic from a single IP exceeds threshold (100 requests/min)
- **THEN** system requires CAPTCHA verification before allowing further queue joins

### Requirement: Queue Position Advancement
The system SHALL process queue users in FIFO order and advance them when previous users complete or time out.

#### Scenario: User completes booking and leaves queue
- **WHEN** user completes reservation
- **THEN** system removes their queue entry and advances the next user

#### Scenario: User abandons queue
- **WHEN** user does not poll or proceed within 2 minutes of being ready
- **THEN** system removes their queue entry and advances the next user

### Requirement: Booking Session Token
The system SHALL issue a short-lived (5 minute) booking session token to users when they are ready to book, limiting the time they can hold stock.

#### Scenario: User has active booking session
- **WHEN** user has a valid session token
- **THEN** they can POST to `/api/bookings/reserve` to attempt ticket reservation

#### Scenario: Booking session expired
- **WHEN** user attempts reservation with expired session token
- **THEN** system returns 401 and instructs them to rejoin the queue
