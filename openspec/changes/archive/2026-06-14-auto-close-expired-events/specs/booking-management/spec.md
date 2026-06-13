## ADDED Requirements

### Requirement: Reserve blocked for finished events
The system SHALL reject booking reserve requests for events whose `end_date` has passed.

#### Scenario: Reserve for finished event
- **WHEN** user with valid session token attempts to reserve tickets for an event where `end_date < now()`
- **THEN** system returns 400 Bad Request with message "event has already ended"

#### Scenario: Reserve for ongoing event
- **WHEN** user with valid session token attempts to reserve tickets for an event where `end_date >= now()`
- **THEN** system proceeds with normal reservation flow
