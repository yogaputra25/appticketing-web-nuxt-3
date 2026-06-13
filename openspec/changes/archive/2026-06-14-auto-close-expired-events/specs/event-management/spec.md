## ADDED Requirements

### Requirement: Finished events excluded from public list
The system SHALL exclude events whose `end_date` has passed from the public event listing.

#### Scenario: List excludes finished events
- **WHEN** client requests `/api/events?status=published`
- **THEN** system only returns events where `end_date >= now()`

### Requirement: Finished event detail returns 404
The system SHALL return 404 for event detail requests when `end_date` has passed.

#### Scenario: Detail for finished event
- **WHEN** client requests `/api/events/{id}` for an event where `end_date < now()`
- **THEN** system returns 404 Not Found
