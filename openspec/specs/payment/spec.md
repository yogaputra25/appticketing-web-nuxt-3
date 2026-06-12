# Payment

## Purpose

Handle payment creation, simulated payment confirmation, expiry, and payment history for users and admins.

## Requirements

### Requirement: Create Payment Order
The system SHALL create a payment order for a pending booking and return payment instructions.

#### Scenario: User initiates payment for pending booking
- **WHEN** user has a pending booking (within TTL) and calls `/api/payments/create`
- **THEN** system creates payment record with status `pending`, returns payment ID, amount, expiry, and a simulated payment URL

#### Scenario: Payment creation for non-pending booking
- **WHEN** user attempts to create payment for a booking that is not in `pending_payment` status
- **THEN** system returns 400 with message "Booking is not in pending payment state"

### Requirement: Simulated Payment Confirmation
The system SHALL expose a simulation endpoint to confirm or fail a payment (replacing real gateway callback in v1).

#### Scenario: Simulate successful payment
- **WHEN** admin/test calls `/api/payments/{id}/simulate?status=success`
- **THEN** system marks payment as `success`, updates booking to `paid`, issues e-ticket code, and returns confirmation

#### Scenario: Simulate failed payment
- **WHEN** payment simulation returns `failed`
- **THEN** system marks payment as `failed`, releases ticket stock back to inventory, and updates booking to `cancelled`

### Requirement: Payment Expiry
The system SHALL automatically mark unpaid payments as expired when they exceed the payment window.

#### Scenario: Payment window expires
- **WHEN** a payment remains `pending` for more than 10 minutes (booking TTL)
- **THEN** background job marks payment as `expired`, booking as `cancelled`, and releases ticket stock

### Requirement: Payment History
The system SHALL allow users to view their payment history.

#### Scenario: User views own payments
- **WHEN** authenticated user requests `/api/payments/me`
- **THEN** system returns list of payments for that user with status, amount, and associated event details

#### Scenario: Admin views all payments
- **WHEN** admin requests `/api/admin/payments`
- **THEN** system returns paginated list of all payments with filters by status and date
