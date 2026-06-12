## 1. Project Setup & Infrastructure

- [x] 1.1 Inisialisasi repository monorepo (folder `apps/web` untuk Nuxt, `apps/api` untuk Go, `docker-compose.yml`)
- [x] 1.2 Setup Docker Compose dengan service: postgres, redis, api, web
- [x] 1.3 Setup environment variables (.env.example) untuk database, redis, JWT secret, payment config
- [x] 1.4 Inisialisasi project Nuxt 3 di `apps/web` dengan Tailwind CSS, Pinia, dan konfigurasi awal
- [x] 1.5 Inisialisasi project Go di `apps/api` dengan Chi router, GORM, dan koneksi database
- [x] 1.6 Setup golang-migrate untuk database migrations

## 2. Database Schema & Migrations

- [x] 2.1 Buat migration tabel `users` (id, email, password_hash, full_name, phone, role, timestamps)
- [x] 2.2 Buat migration tabel `events` (id, title, description, venue, start_date, end_date, banner_url, status, timestamps)
- [x] 2.3 Buat migration tabel `ticket_categories` (id, event_id, name, price, total_stock, available_stock, timestamps) + index pada event_id
- [x] 2.4 Buat migration tabel `bookings` (id, user_id, event_id, total_amount, status, expires_at, e_ticket_codes JSONB, timestamps) + index pada user_id dan status
- [x] 2.5 Buat migration tabel `booking_items` (id, booking_id, ticket_category_id, quantity, unit_price)
- [x] 2.6 Buat migration tabel `payments` (id, booking_id, amount, status, payment_method, paid_at, timestamps) + index pada booking_id
- [x] 2.7 Buat migration tabel `queue_tokens` (id, user_id, event_id, token, status, created_at, expires_at) untuk tracking queue

## 3. Backend: Auth Module (user-auth)

- [x] 3.1 Implementasi model User di GORM dengan bcrypt password hashing
- [x] 3.2 Implementasi repository layer untuk User (Create, FindByEmail, FindByID, Update)
- [x] 3.3 Implementasi JWT utility (GenerateToken, ValidateToken, middleware)
- [x] 3.4 Implementasi handler POST `/api/auth/register` dengan validasi
- [x] 3.5 Implementasi handler POST `/api/auth/login` dengan credential check
- [x] 3.6 Implementasi handler GET `/api/auth/me` (protected)
- [x] 3.7 Implementasi handler PUT `/api/auth/me` untuk update profile
- [x] 3.8 Implementasi role-based middleware untuk admin endpoints
- [x] 3.9 Tulis unit tests untuk auth handlers

## 4. Backend: Event Management (event-management)

- [x] 4.1 Implementasi model Event di GORM
- [x] 4.2 Implementasi repository Event (CRUD + filter by status, pagination)
- [x] 4.3 Implementasi handler POST `/api/admin/events` (admin only) untuk create event
- [x] 4.4 Implementasi handler GET `/api/events` (public) dengan pagination & filter
- [x] 4.5 Implementasi handler GET `/api/events/{id}` dengan ticket categories embedded
- [x] 4.6 Implementasi handler PUT `/api/admin/events/{id}` untuk update
- [x] 4.7 Implementasi handler POST `/api/admin/events/{id}/publish` dengan validasi
- [x] 4.8 Implementasi handler DELETE `/api/admin/events/{id}` dengan soft-delete & booking check

## 5. Backend: Ticket Inventory (ticket-inventory)

- [x] 5.1 Implementasi model TicketCategory di GORM dengan hooks (auto-set available_stock = total_stock on create)
- [x] 5.2 Implementasi repository TicketCategory (Create, ListByEvent, Update, GetForUpdate dengan SELECT FOR UPDATE)
- [x] 5.3 Implementasi handler POST `/api/admin/events/{eventId}/categories` untuk create category
- [x] 5.4 Implementasi handler GET `/api/events/{eventId}/categories` (public) dengan sold-out flag
- [x] 5.4 Implementasi handler PUT `/api/admin/categories/{id}` untuk update (dengan validasi stock ≥ sold)
- [x] 5.5 Implementasi service `ReserveStock(ctx, categoryID, qty)` dengan row-level lock + atomic decrement
- [x] 5.6 Implementasi service `ReleaseStock(categoryID, qty)` untuk restore stock
- [x] 5.7 Tulis concurrency tests untuk reserve stock (simulate 1000 concurrent requests)

## 6. Backend: Ticket War Queue (ticket-war)

- [ ] 6.1 Setup Redis client di Go dengan connection pool
- [ ] 6.2 Implementasi service `JoinQueue(userID, eventID)` dengan Redis ZADD ke sorted set
- [ ] 6.3 Implementasi service `GetQueuePosition(userID, eventID)` dengan ZRANK
- [ ] 6.4 Implementasi handler POST `/api/war/join` dengan rate limiting (5 req/min per user)
- [ ] 6.5 Implementasi handler GET `/api/war/status` (polling) yang return posisi & is_ready
- [ ] 6.6 Implementasi service `ProcessQueue` background worker yang advance user ke booking session
- [ ] 6.7 Implementasi service `IssueBookingSessionToken(userID, eventID)` dengan TTL 5 menit (Redis)
- [ ] 6.8 Implementasi cleanup job untuk queue token yang abandoned (>2 menit tidak aktif)
- [ ] 6.9 Tulis integration tests untuk queue flow end-to-end

## 7. Backend: Booking Management (booking-management)

- [ ] 7.1 Implementasi model Booking & BookingItem di GORM
- [ ] 7.2 Implementasi repository Booking (Create dengan transaction, GetByUser, GetByID, UpdateStatus)
- [ ] 7.3 Implementasi handler POST `/api/bookings/reserve` dengan validasi session token + atomic stock decrement
- [ ] 7.4 Implementasi handler GET `/api/bookings/me` (paginated, user only)
- [ ] 7.5 Implementasi handler GET `/api/bookings/{id}` dengan ownership check
- [ ] 7.6 Implementasi handler POST `/api/bookings/{id}/cancel` untuk pending bookings
- [ ] 7.7 Implementasi background job `ExpirePendingBookings` (jalan tiap 1 menit) yang release stock
- [ ] 7.8 Tulis tests untuk booking flow: reserve → cancel → stock restored

## 8. Backend: Payment Module (payment)

- [ ] 8.1 Implementasi model Payment di GORM
- [ ] 8.2 Implementasi repository Payment (Create, GetByID, UpdateStatus, ListByUser)
- [ ] 8.3 Implementasi handler POST `/api/payments/create` untuk initiate payment
- [ ] 8.4 Implementasi handler POST `/api/payments/{id}/simulate` (test mode) untuk success/fail
- [ ] 8.5 Implementasi service `ConfirmPayment` yang update booking ke paid dan generate e-ticket codes (UUID)
- [ ] 8.6 Implementasi handler GET `/api/payments/me` (user) dan GET `/api/admin/payments` (admin)
- [ ] 8.7 Hook expiry job payment ke booking expiry (jika payment expired → booking cancelled)
- [ ] 8.8 Tulis tests untuk payment lifecycle: create → success → e-ticket issued

## 9. Backend: Admin Dashboard Endpoints

- [ ] 9.1 Implementasi handler GET `/api/admin/stats` (total events, bookings, revenue, users)
- [ ] 9.2 Implementasi handler GET `/api/admin/bookings` dengan filter status, date range
- [ ] 9.3 Implementasi handler POST `/api/admin/users` untuk create admin user
- [ ] 9.4 Implementasi handler GET `/api/admin/users` untuk list users
- [ ] 9.5 Tulis seeder untuk admin user default + sample event

## 10. Frontend: Setup & Layout (Nuxt 3 + Tailwind)

- [x] 10.1 Konfigurasi Tailwind CSS di Nuxt 3 dengan custom theme (colors, fonts)
- [x] 10.2 Setup Pinia stores: `useAuthStore`, `useEventStore`, `useQueueStore`, `useBookingStore`
- [x] 10.3 Setup `$fetch` wrapper dengan auto-attach JWT token dan baseURL
- [x] 10.4 Buat layout default dengan navbar (logo, menu, user avatar) dan footer
- [x] 10.5 Buat middleware `auth` untuk protected pages dan `admin` untuk admin pages
- [x] 10.6 Buat halaman login & register dengan form validation (VeeValidate + Zod)

## 11. Frontend: Public Pages

- [ ] 11.1 Halaman Home (`/`) — list upcoming events dengan card grid + filter
- [ ] 11.2 Halaman Event Detail (`/events/[id]`) — banner, info, ticket categories dengan harga
- [ ] 11.3 Halaman Events List (`/events`) — searchable & filterable list
- [ ] 11.4 Komponen `EventCard`, `TicketCategoryCard`, `CountdownTimer` (reusable)

## 12. Frontend: War Tiket Flow

- [ ] 12.1 Halaman War (`/events/[id]/war`) — countdown, info event, tombol "Mulai War"
- [ ] 12.2 Halaman Antrian (`/events/[id]/queue`) — posisi, total antrian, ETA, polling tiap 2 detik
- [ ] 12.3 Komponen `QueueStatus` dengan progress bar dan animasi
- [ ] 12.4 Halaman Booking (`/events/[id]/booking`) — pilih quantity, ringkasan, tombol "Lanjut Bayar"
- [ ] 12.5 Integrasi flow: war → queue → ready → booking (dengan handling token expiry)

## 13. Frontend: Booking & Payment Pages

- [ ] 13.1 Halaman Payment (`/bookings/[id]/pay`) — instruksi pembayaran, timer countdown, simulasi button
- [ ] 13.2 Halaman My Bookings (`/my/bookings`) — list booking dengan status badges
- [ ] 13.3 Halaman E-Ticket (`/my/bookings/[id]`) — detail + kode e-ticket dengan QR placeholder
- [ ] 13.4 Halaman Profile (`/profile`) — view & edit profil user

## 14. Frontend: Admin Dashboard

- [ ] 14.1 Layout admin dengan sidebar navigation
- [ ] 14.2 Halaman Admin Dashboard (`/admin`) — statistik cards (revenue, bookings, events)
- [ ] 14.3 Halaman Admin Events (`/admin/events`) — CRUD events dengan form & table
- [ ] 14.4 Halaman Admin Event Detail (`/admin/events/[id]`) — manage ticket categories
- [ ] 14.5 Halaman Admin Bookings (`/admin/bookings`) — list & filter all bookings
- [ ] 14.6 Halaman Admin Payments (`/admin/payments`) — payment monitoring
- [ ] 14.7 Halaman Admin Users (`/admin/users`) — list & manage users

## 15. Integration, Testing & Polish

- [ ] 15.1 End-to-end test flow: register → login → browse event → war → queue → book → pay → e-ticket
- [ ] 15.2 Setup seed data (1 admin, 3 events dengan 2-3 kategori masing-masing)
- [ ] 15.3 Tulis README dengan setup instructions, env vars, dan deployment notes
- [ ] 15.4 Setup error handling & loading states konsisten di semua halaman
- [ ] 15.5 Tambahkan empty states & 404 page
- [ ] 15.6 Optimasi responsive design untuk mobile (Tailwind breakpoints)
- [ ] 15.7 Setup logging di Go (slog) dan request ID middleware
- [ ] 15.8 Performance check: load test 1000 concurrent users dengan `k6` atau `wrk`
- [ ] 15.9 Final review: cek ulang semua spec scenarios terpenuhi
