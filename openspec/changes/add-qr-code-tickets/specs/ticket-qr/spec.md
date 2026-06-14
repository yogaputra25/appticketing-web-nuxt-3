## ADDED Requirements

### Requirement: Ticket QR Generation
The system SHALL create a unique `Ticket` row for each ticket unit purchased, with a unique ticket code and QR code displayed on the frontend.

#### Scenario: Tickets created on successful payment
- **WHEN** payment for a booking is confirmed as successful
- **THEN** system creates one `Ticket` row per ticket unit (expanding `BookingItem.Quantity`) with a unique alphanumeric `TicketCode` and status `active`

#### Scenario: User views their tickets
- **WHEN** authenticated user requests `GET /api/tickets`
- **THEN** system returns paginated list of their tickets with event name, category name, status, and ticket code

#### Scenario: User views single ticket
- **WHEN** authenticated user requests `GET /api/tickets/{id}` for a ticket they own
- **THEN** system returns full ticket details including event info, ticket code, QR code URL, and scan status

#### Scenario: User sees QR code on ticket detail page
- **WHEN** user opens a ticket detail page
- **THEN** frontend displays a scannable QR code image encoding the verify URL `/tickets/v/{ticketCode}`

### Requirement: Ticket Verification via Scan
The system SHALL allow scanning a ticket's QR code to verify its authenticity and mark it as used.

#### Scenario: Scan valid ticket
- **WHEN** a `POST /api/tickets/verify/{ticketCode}` request is made for a ticket with status `active`
- **THEN** system returns 200 OK with ticket data, changes status to `used`, and sets `scanned_at`

#### Scenario: Scan already-used ticket
- **WHEN** a `POST /api/tickets/verify/{ticketCode}` request is made for a ticket with status `used`
- **THEN** system returns 400 Bad Request with message "ticket has already been used" and includes the original scan timestamp

#### Scenario: Scan invalid ticket code
- **WHEN** a `POST /api/tickets/verify/{ticketCode}` request is made with a non-existent ticket code
- **THEN** system returns 404 Not Found

### Requirement: Scanner Page
The system SHALL provide a scanner page that reads QR codes via device camera.

#### Scenario: User scans QR code at venue
- **WHEN** event staff opens the scanner page, points camera at a ticket QR code
- **THEN** system reads the ticket code from the QR, calls the verify endpoint, and displays valid/invalid result with ticket details
