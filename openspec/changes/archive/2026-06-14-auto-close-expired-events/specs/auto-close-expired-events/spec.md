## ADDED Requirements

### Requirement: Events automatically close after end_date
The system SHALL block access to events whose `end_date` has passed, preventing further ticket sales.

#### Scenario: War join blocked for finished event
- **WHEN** user attempts to join war queue for an event where `end_date < now()`
- **THEN** system returns 400 Bad Request with message "event has already ended"

#### Scenario: Booking reserve blocked for finished event
- **WHEN** user attempts to reserve tickets for an event where `end_date < now()`
- **THEN** system returns 400 Bad Request with message "event has already ended"

#### Scenario: Finished event excluded from public list
- **WHEN** client requests public event list
- **THEN** system excludes events where `end_date < now()` from results

#### Scenario: Finished event detail returns 404
- **WHEN** client requests detail for an event where `end_date < now()`
- **THEN** system returns 404 Not Found

#### Scenario: User sees "Event Ended" indicator
- **WHEN** user navigates to event detail or war page for a finished event
- **THEN** frontend shows "Event Ended" message and hides the booking/war action
