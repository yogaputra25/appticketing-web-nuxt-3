# Booking Management

## Purpose

Manage ticket reservations, booking lifecycle (pending → paid/cancelled/expired), and e-ticket generation.

## Requirements

### Requirement: Create Booking (Reservation)
The system SHALL allow authenticated users with a valid booking session to reserve tickets atomically.

#### Scenario: Successful ticket reservation
- **WHEN** user with valid session reserves quantity ≤ available stock for a category
- **THEN** system decrements stock, creates booking with status `pending_payment`, sets expiry 10 minutes, and returns booking details

#### Scenario: Insufficient stock
- **WHEN** user attempts to reserve more than available stock
- **THEN** system returns 409 Conflict with message "Not enough stock" and current available count

#### Scenario: Reservation without valid session
- **WHEN** user without valid booking session token attempts reservation
- **THEN** system returns 401 Unauthorized

### Requirement: View Own Bookings
The system SHALL allow users to view their booking history with status and event details.

#### Scenario: User views own bookings
- **WHEN** authenticated user requests `/api/bookings/me`
- **THEN** system returns paginated list of bookings with event title, ticket category, quantity, status, and total amount

#### Scenario: User views single booking
- **WHEN** user requests `/api/bookings/{id}` for a booking they own
- **THEN** system returns full booking details including e-ticket code if paid

#### Scenario: User attempts to view another user's booking
- **WHEN** user requests a booking ID that belongs to another user
- **THEN** system returns 404 Not Found (does not reveal existence)

### Requirement: Cancel Booking
The system SHALL allow users to cancel a `pending_payment` booking and release the held stock.

#### Scenario: User cancels pending booking
- **WHEN** user cancels a booking in `pending_payment` status
- **THEN** system releases stock back to inventory, marks booking as `cancelled`, and returns success

#### Scenario: Cancel paid booking
- **WHEN** user attempts to cancel a booking in `paid` status
- **THEN** system returns 400 with message "Paid bookings cannot be self-cancelled, contact support"

### Requirement: E-Ticket Generation
The system SHALL generate a unique e-ticket code for each paid booking item.

#### Scenario: E-ticket code issued on successful payment
- **WHEN** payment for a booking is confirmed
- **THEN** system generates a unique alphanumeric code per ticket item and attaches it to the booking

#### Scenario: User views e-ticket
- **WHEN** user accesses a paid booking
- **THEN** system displays the e-ticket codes, event details, and venue information
