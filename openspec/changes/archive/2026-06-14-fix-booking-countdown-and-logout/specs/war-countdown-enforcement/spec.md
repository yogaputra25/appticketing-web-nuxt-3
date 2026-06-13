## ADDED Requirements

### Requirement: War join disabled before event start

The "Mulai War" button SHALL be disabled until the countdown timer reaches zero (event start time). The backend SHALL also reject war join requests made before `event.start_date`.

#### Scenario: Button disabled during countdown

- **WHEN** the war page loads and countdown timer is still running
- **THEN** the "Mulai War" button SHALL be disabled (not clickable)
- **AND** the button SHALL show "Mulai War" text (unchanged)

#### Scenario: Button enabled after countdown

- **WHEN** the countdown timer reaches zero
- **THEN** the "Mulai War" button SHALL become enabled
- **AND** clicking it SHALL call the war join API

#### Scenario: Backend rejects early war join

- **WHEN** a client sends `POST /api/war/join` before `event.start_date`
- **THEN** the API SHALL respond with HTTP 400 and message "event has not started yet"

#### Scenario: Backend allows war join after event start

- **WHEN** a client sends `POST /api/war/join` after or at `event.start_date`
- **THEN** the API SHALL process the request normally (no time-related rejection)

### Requirement: Booking page unavailable before event start

The booking page SHALL show a "belum dimulai" message instead of the booking form when the event's `start_date` has not yet been reached.

#### Scenario: Booking page before event start

- **WHEN** a user navigates to the booking page and `event.start_date` is in the future
- **THEN** the booking form SHALL NOT be displayed
- **AND** a message "Event belum dimulai" SHALL be shown
- **AND** an error SHALL NOT be thrown (page loads gracefully)

#### Scenario: Booking page after event start

- **WHEN** a user navigates to the booking page and `event.start_date` has passed
- **THEN** the booking form SHALL be displayed as normal
