## ADDED Requirements

### Requirement: Join blocked for finished events
The system SHALL reject war join requests for events whose `end_date` has passed.

#### Scenario: Join finished event
- **WHEN** user clicks "Mulai War" on an event where `end_date < now()`
- **THEN** system returns 400 Bad Request with message "event has already ended"

#### Scenario: Join ongoing event
- **WHEN** user clicks "Mulai War" on an event where `end_date >= now()`
- **THEN** system proceeds with normal join flow
