## MODIFIED Requirements

### Requirement: E-Ticket Generation
**Previous behavior:** System generated unique alphanumeric codes per booking item and stored them as a flat list in `Booking.ETicketCodes`.
**New behavior:** System creates individual `Ticket` rows per ticket unit with unique codes, replacing the flat e-ticket code list.

#### Scenario: Tickets created on successful payment
- **WHEN** payment for a booking is confirmed
- **THEN** system creates one `Ticket` row per ticket unit (expanding `BookingItem.Quantity`) with a unique `TicketCode` and status `active`

#### Scenario: User views tickets for a booking
- **WHEN** authenticated user views a paid booking detail
- **THEN** system returns booking info along with associated `Ticket` list instead of legacy `e_ticket_codes`

## ADDED Requirements

### Requirement: View My Tickets Dashboard
The system SHALL provide a dedicated page listing all tickets owned by the user across all bookings.

#### Scenario: User opens My Tickets page
- **WHEN** authenticated user navigates to `/my/tickets`
- **THEN** system displays paginated list of tickets with event name, date, category, status (active/used), and a link to view QR code

#### Scenario: Empty ticket list
- **WHEN** authenticated user has no tickets
- **THEN** system displays empty state message "Belum ada tiket"
