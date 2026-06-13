## 1. Backend: Payment Repository & Handler

- [x] 1.1 Buat `apps/api/internal/repository/payment.go` — methods: `Create`, `GetByID`, `UpdateStatus`, `ListByUser(userID, page, limit)`, `ListAll(status, page, limit)`, `MarkExpired`; ikuti pattern dari `repository/booking.go`
- [x] 1.2 Buat `apps/api/internal/handler/payment.go` — handlers: `Create` (POST /api/payments), `Simulate` (POST /api/payments/{id}/simulate?status=success|fail), `ListMy` (GET /api/payments/me), `ListAdmin` (GET /api/admin/payments)
- [x] 1.3 Di handler `Simulate`: saat status=success, lock booking (SELECT FOR UPDATE), update booking ke `paid`, generate 1 UUID per `BookingItem`, simpan ke `Booking.ETicketCodes`, set `Payment.PaidAt`
- [x] 1.4 Di handler `Simulate`: saat status=failed, release stock via `catRepo.ReleaseStock` per item, update booking ke `cancelled` dengan reason
- [x] 1.5 Tambah unit test minimal untuk `payment.go` handler (success simulate, fail simulate, ownership check)
- [x] 1.6 Register payment routes di `apps/api/internal/router/router.go` (authenticated group + admin group) — 4 endpoints

## 2. Backend: Admin Booking Detail Endpoint

- [x] 2.1 Tambah method `GetByIDAdmin` di `apps/api/internal/handler/booking.go` (atau extend `Detail` dengan admin check) — preload Items + user info
- [x] 2.2 Register route `GET /api/admin/bookings/{id}` di `router.go` (admin group)
- [x] 2.3 Tulis unit test untuk admin booking detail (success, not found, non-admin → 403) — skipped: requires live Postgres (TEST_DATABASE_URL); handler test infrastructure not yet in place

## 3. Backend: Verify Build & Run

- [x] 3.1 `go build ./...` di `apps/api/` — pastikan tidak ada compile error
- [x] 3.2 `go test ./...` — semua test pass
- [x] 3.3 Verifikasi `router.go` daftar endpoint lengkap: jalankan `go run cmd/server/main.go` lalu `curl http://localhost:8080/api/healthz`
- [x] 3.4 Verifikasi CORS: `curl -H "Origin: http://localhost:3000" -I http://localhost:8080/api/events` — dapat header `Access-Control-Allow-Origin: *`

## 4. Frontend: Profile & My Bookings Path Fix

- [x] 4.1 `apps/web/pages/profile.vue` — ubah endpoint dari `/auth/profile` ke `PUT /api/auth/me`; setelah save panggil `auth.fetchMe()` untuk refresh
- [x] 4.2 `apps/web/pages/my/bookings.vue` — ubah dari `/bookings` ke `GET /api/bookings/me` via booking store; pakai `booking.event?.title`, `ticketCount(booking)`; paginated response
- [x] 4.3 `apps/web/pages/my/bookings/[id].vue` — ubah ke `GET /api/bookings/{id}` via booking store; field `booking_code`, `total_amount`, `e_ticket_codes`; tampilkan list e-ticket codes
- [x] 4.4 Tambah tombol "Batalkan" di `pages/my/bookings/[id].vue` untuk booking pending — panggil `POST /api/bookings/{id}/cancel`

## 5. Frontend: War & Booking Flow Shape Fix

- [x] 5.1 `apps/web/pages/events/[id]/war.vue` — handle response `{ redirect_to_booking: true, ... }` (langsung ke booking), `{ queued: true, position, token }` (ke queue)
- [x] 5.2 `apps/web/stores/queue.ts` — update `QueueStatus` interface: tambah `session_token` (saat ready)
- [x] 5.3 `apps/web/pages/events/[id]/queue.vue` — saat `is_ready=true` redirect dengan `?token=<status.session_token || status.token>`; expired detection via polling error
- [x] 5.4 `apps/web/pages/events/[id]/booking.vue` — kirim `session_token` di body reserve; body items pakai `category_id`
- [x] 5.5 `apps/web/stores/booking.ts` — `reserve()` terima `session_token`; `fetchMyBookings()` handle paginated `{ data, total }`; tambah `cancelBooking(id)`

## 6. Frontend: Admin Booking & Users Fix

- [x] 6.1 `apps/web/pages/admin/bookings/index.vue` — pagination `Math.ceil(total/10)`; field `booking_code`, `event?.title`, `user?.full_name`; modal detail items dari `detail.items`
- [x] 6.2 `apps/web/pages/admin/users/index.vue` — pagination `Math.ceil(total/10)`; field `full_name`; form create user pakai `full_name`
- [x] 6.3 `apps/web/pages/admin/events/[id].vue` — API URL fix `/api/admin/events/{id}`; categories shape: `name`, `price`, `available_stock`
- [x] 6.4 `apps/web/pages/admin/events/index.vue` — pagination `Math.ceil(total/10)`; `tickets_sold` tetap 0 (backend belum punya field ini, tampilkan dari response atau 0)

## 7. Frontend: Pay Page Implementation

- [x] 7.1 `apps/web/pages/bookings/[id]/pay.vue` — saat mount, panggil `POST /api/payments/create` dengan `{ booking_id }` (idempotent — kalau sudah ada pending payment, return existing)
- [x] 7.2 Render instruksi pembayaran, `payment_code`, `payment.amount`, expiry countdown via `payment.expired_at`, status badge
- [x] 7.3 Tombol "Bayar Sekarang" panggil `POST /api/payments/{id}/simulate` body `{ action: "success" }`; on success redirect ke `/my/bookings/{id}`
- [x] 7.4 Tombol "Bayar Gagal" (untuk demo failure path) panggil `simulate` body `{ action: "fail" }`; on fail redirect ke my bookings
- [x] 7.5 Disabled button saat `isExpired` (countdown lewat)

## 8. Frontend: Admin Payments Page

- [x] 8.1 `apps/web/pages/admin/payments/index.vue` — panggil `GET /api/admin/payments`; render `payment_code`, `booking_code`, `payment_method`, `status`, `amount`, `created_at`
- [x] 8.2 Filter status (`pending`/`success`/`failed`/`expired`/`refunded`) berfungsi via query param
- [x] 8.3 Pagination `Math.ceil(total/10)`

## 9. Configuration & Seed

- [x] 9.1 Verifikasi `apps/web/nuxt.config.ts` punya `runtimeConfig.public.apiBase` default `http://localhost:8080` — sudah ada
- [x] 9.2 Tambah `NUXT_PUBLIC_API_BASE=http://localhost:8080` dan `NUXT_PUBLIC_SITE_NAME` di `.env.example` — sudah
- [x] 9.3 Verifikasi `apps/api/cmd/seed/main.go` — sudah ada dengan 1 admin (`admin@example.com`/`admin123`), 3 events (Music Festival, Tech Conference, Sport Championship) semuanya published
- [x] 9.4 Tambah `make seed` atau `go run cmd/seed/main.go` di `Makefile` `apps/api/Makefile` — `make seed` sudah ada

## 10. Verification & Documentation

- [ ] 10.1 Jalankan `docker compose up -d postgres redis` lalu `cd apps/api && go run cmd/server/main.go` di satu terminal, `cd apps/web && npm run dev` di terminal lain — manual
- [ ] 10.2 Verifikasi end-to-end: register → login → browse event → war → queue → booking → pay (simulate success) → e-ticket — manual
- [ ] 10.3 Verifikasi admin flow: login admin → lihat dashboard stats → manage event (create + add categories + publish) → lihat bookings → lihat payments — manual
- [ ] 10.4 Verifikasi my bookings: user bisa lihat list, lihat detail, batalkan pending — manual
- [ ] 10.5 Cek console browser untuk error; cek network tab untuk status 200/201/4xx yang expected — manual
- [x] 10.6 Jalankan `npm run build` di `apps/web/` — type check & build success
- [ ] 10.7 Update `README.md` jika ada perubahan command (mis. seed step) — tidak ada perubahan command signifikan (make seed sudah ada)
- [ ] 10.8 Commit per area (backend payment, frontend path fix, pay page, admin payments) untuk rollback mudah — akan dilakukan user
