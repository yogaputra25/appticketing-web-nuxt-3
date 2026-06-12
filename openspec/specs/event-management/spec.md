# Event Management

## Purpose

Allow admins to create, update, publish, and delete events. Display published events to the public with pagination.

## Requirements

### Requirement: Create Event
The system SHALL allow admins to create new events with title, description, venue, start datetime, end datetime, and banner image.

#### Scenario: Admin creates event successfully
- **WHEN** admin submits valid event data
- **THEN** system creates event with status `draft` and returns event ID

#### Scenario: Non-admin attempts to create event
- **WHEN** regular user attempts to create event
- **THEN** system returns 403 Forbidden

### Requirement: List Events (Public)
The system SHALL expose a public endpoint to list published/upcoming events with pagination.

#### Scenario: List upcoming published events
- **WHEN** client requests `/api/events?status=published&page=1&limit=20`
- **THEN** system returns paginated list of published events with start_date in the future

#### Scenario: Event detail
- **WHEN** client requests `/api/events/{id}`
- **THEN** system returns event details including ticket categories and prices

### Requirement: Update Event
The system SHALL allow admins to update event details.

#### Scenario: Admin updates event
- **WHEN** admin submits valid update data
- **THEN** system updates event and returns updated record

### Requirement: Publish Event
The system SHALL allow admins to change event status from `draft` to `published`.

#### Scenario: Publish draft event
- **WHEN** admin publishes an event with at least one ticket category and stock
- **THEN** system changes status to `published` and event becomes visible publicly

#### Scenario: Publish event without ticket categories
- **WHEN** admin attempts to publish event with no ticket categories
- **THEN** system returns 400 with message "Add at least one ticket category before publishing"

### Requirement: Delete Event
The system SHALL allow admins to soft-delete events (only events without bookings).

#### Scenario: Delete event with no bookings
- **WHEN** admin deletes an event with zero bookings
- **THEN** system marks event as deleted and hides from public list

#### Scenario: Delete event with existing bookings
- **WHEN** admin attempts to delete event with active bookings
- **THEN** system returns 400 with message "Cannot delete event with existing bookings"
