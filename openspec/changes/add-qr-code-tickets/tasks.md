## 1. Backend: Ticket Model & Repository

- [x] 1.1 Add `Ticket` model to `apps/api/internal/model/model.go` with fields: ID, BookingID, TicketCode (unique), CategoryName, Status (active/used/refunded), ScannedAt, CreatedAt
- [x] 1.2 Create `TicketRepository` in `apps/api/internal/repository/ticket.go` with methods: Create, ListByUser, FindByCode, FindByID, MarkAsUsed

## 2. Backend: Update Payment Flow

- [x] 2.1 In `apps/api/internal/handler/payment.go` `ConfirmPayment`, replace `generateETicketCodes` with logic to create Ticket rows (expand BookingItem.Quantity → N Ticket rows) via TicketRepository
- [x] 2.2 Add `TicketRepository` dependency to `PaymentHandler` struct and constructor
- [x] 2.3 Update router.go to wire `TicketRepository` into `PaymentHandler`

## 3. Backend: Ticket API Endpoints

- [x] 3.1 Create `TicketHandler` in `apps/api/internal/handler/ticket.go` with `ListMy` (GET /api/tickets), `Detail` (GET /api/tickets/{id}), `Verify` (POST /api/tickets/verify/{code})
- [x] 3.2 Register TicketHandler routes in router.go

## 4. Frontend: QR & Scanner Dependencies

- [x] 4.1 Install `qrcode` and `vue-qrcode-reader` npm packages in `apps/web`
- [x] 4.2 Create QR code utility composable (generate QR data URL from ticket code)

## 5. Frontend: My Tickets Page

- [x] 5.1 Create `apps/web/pages/my/tickets/index.vue` — list user's tickets with event name, date, category, status
- [ ] 5.2 Create `apps/web/pages/my/tickets/[id].vue` — ticket detail with QR code display

## 6. Frontend: Scanner Page

- [x] 6.1 Create `apps/web/pages/tickets/scan.vue` — scan QR code via camera and verify

## 7. Frontend: Navigation

- [x] 7.1 Add "Tiket Saya" link to user dashboard navigation
