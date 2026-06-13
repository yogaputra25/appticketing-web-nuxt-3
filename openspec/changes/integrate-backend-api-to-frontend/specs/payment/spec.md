## ADDED Requirements

### Requirement: Backend Payment Handler Implementation
The system SHALL implement the payment handlers and routes in `apps/api/internal/handler/payment.go` and register them in `apps/api/internal/router/router.go`, fulfilling the `payment` capability requirements that were specified in change `sistem-pemesanan-tiket-war` but not yet implemented in backend code.

#### Scenario: Payment handler files exist
- **WHEN** the backend is built
- **THEN** the file `apps/api/internal/handler/payment.go` exists and is referenced from `router.go` for routes `POST /api/payments`, `POST /api/payments/{id}/simulate`, `GET /api/payments/me`, and `GET /api/admin/payments`

#### Scenario: Payment repository exists
- **WHEN** the backend is built
- **THEN** the file `apps/api/internal/repository/payment.go` exists with at minimum: `Create`, `GetByID`, `UpdateStatus`, `ListByUser`, `ListAll` (for admin), and `MarkExpired` methods following the pattern of `repository/booking.go`

### Requirement: POST /api/payments Creates Payment Order
The system SHALL create a `pending` payment record for a `pending_payment` booking via `POST /api/payments` and return the payment ID, amount, expiry timestamp, and a `redirect_url` (or payment URL) used by the frontend pay page.

#### Scenario: User creates payment for pending booking
- **WHEN** authenticated user calls `POST /api/payments` with body `{ booking_id }` and the booking is in `pending_payment` status and belongs to the user
- **THEN** backend creates a `pending` payment (or returns existing pending payment), decrements nothing, and returns `{ id, booking_id, amount, status: "pending", expires_at, redirect_url }` (HTTP 201 or 200)

#### Scenario: Cannot pay for non-pending booking
- **WHEN** user calls `POST /api/payments` for a booking in `paid` or `cancelled` status
- **THEN** backend returns 400 with message "Booking is not in pending payment state"

#### Scenario: Cannot pay for another user's booking
- **WHEN** user calls `POST /api/payments` with a `booking_id` they do not own
- **THEN** backend returns 404

### Requirement: POST /api/payments/{id}/simulate Confirms Payment
The system SHALL expose `POST /api/payments/{id}/simulate?status=success` (or `failed`) to confirm or fail a payment in development — replacing a real gateway callback for v1.

#### Scenario: Simulate successful payment
- **WHEN** authenticated user (owner of the payment) calls `POST /api/payments/{id}/simulate?status=success`
- **THEN** backend marks payment as `success`, updates booking to `paid`, generates one e-ticket code per booking item, sets `paid_at`, and returns the updated payment

#### Scenario: Simulate failed payment
- **WHEN** user calls the simulate endpoint with `status=failed`
- **THEN** backend marks payment as `failed`, releases ticket stock, and updates booking to `cancelled` with reason "payment failed"

### Requirement: GET /api/payments/me Returns User Payments
The system SHALL return the authenticated user's payment history via `GET /api/payments/me`.

#### Scenario: User views own payment history
- **WHEN** authenticated user calls `GET /api/payments/me`
- **THEN** backend returns a paginated list of payments belonging to the user with `id`, `booking_id`, `booking_code`, `amount`, `status`, `payment_method`, `created_at`, `paid_at`

### Requirement: GET /api/admin/payments Lists All Payments
The system SHALL allow admins to list all payments with filters via `GET /api/admin/payments?status=...&page=...&limit=...`.

#### Scenario: Admin lists all payments
- **WHEN** admin calls `GET /api/admin/payments?page=1&limit=20`
- **THEN** backend returns paginated `{ data, total, page, limit }` list of all payments

#### Scenario: Admin filters payments by status
- **WHEN** admin calls `GET /api/admin/payments?status=paid`
- **THEN** backend returns only payments with `status = "paid"`

### Requirement: Frontend Pay Page Uses Real Payment Endpoint
The `pages/bookings/[id]/pay.vue` page SHALL call the real payment endpoints and stop being a placeholder.

#### Scenario: Pay page loads payment status
- **WHEN** user opens `/bookings/{id}/pay` for a pending booking
- **THEN** frontend creates a payment via `POST /api/payments` (or fetches existing one) and displays amount, expiry countdown, and the payment instructions

#### Scenario: Pay button triggers simulate
- **WHEN** user clicks "Bayar Sekarang"
- **THEN** frontend calls `POST /api/payments/{id}/simulate?status=success` and on success redirects to `/my/bookings/{id}` showing the e-ticket

#### Scenario: Booking expired shows cancel message
- **WHEN** the booking is past `expires_at`
- **THEN** pay page shows "Waktu pembayaran telah habis" and disables the pay button

### Requirement: Frontend Admin Payments Page Uses Real Endpoint
The `pages/admin/payments/index.vue` page SHALL fetch payments from `GET /api/admin/payments` and render the list with the correct field names.

#### Scenario: Admin payments list renders
- **WHEN** admin opens `/admin/payments`
- **THEN** page calls `GET /api/admin/payments`, displays payment_code, booking_code, payment_method, status badge, amount, and created_at for each row
