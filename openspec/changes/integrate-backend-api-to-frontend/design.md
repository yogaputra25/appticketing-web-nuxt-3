## Context

Project ini adalah SaaS pemesanan tiket konser (War Tiket) dengan monorepo Nuxt 3 (frontend) + Go Chi (backend) + PostgreSQL + Redis. Backend Go untuk modul auth, event, ticket category, war queue, dan booking sudah selesai. Frontend sudah punya UI lengkap untuk publik, user, dan admin, dengan beberapa store Pinia (`useEventStore`, `useQueueStore`, `useBookingStore`) dan composable `useApi()`.

**Current state integrasi (hasil audit):**

| Area | Status | Catatan |
|---|---|---|
| Auth (login, register, me, update) | ✅ | Sudah pakai endpoint benar via `useAuthStore` |
| Home & Events list | ✅ | `useEventStore.fetchPublicEvents` → `/api/events` |
| Event detail | ✅ | `useEventStore.fetchEventDetail` → `/api/events/{id}` |
| War flow (join, queue, booking) | ⚠️ | Path benar tapi shape field frontend vs backend beda (mis. `token` vs `session_token`, `ticket_category_id` vs `category_id`) |
| My bookings | ❌ | Pakai `/bookings` (salah), harusnya `/api/bookings/me` |
| My booking detail | ❌ | Pakai `/bookings/{id}` (salah) + field `code`/`total_price`/`ticket_count` salah |
| Profile | ❌ | Pakai `/auth/profile` (salah), harusnya `PUT /api/auth/me` |
| Pay page | ❌ | Placeholder, belum ada endpoint |
| Admin dashboard | ✅ | Pakai `/api/admin/stats` & `/api/admin/bookings` |
| Admin events CRUD | ⚠️ | Endpoint benar tapi field shape perlu dicek (mis. `tickets_sold`) |
| Admin bookings modal detail | ❌ | `GET /api/admin/bookings/{id}` belum ada di router |
| Admin payments | ❌ | Endpoint `/api/admin/payments` belum ada di router |
| Admin users | ⚠️ | Endpoint ada tapi response shape pakai `last_page` (backend pakai `total`+`limit`) |

**Tambahan di backend yang dibutuhkan:** payment handler, payment repository, 4 payment routes, 1 admin booking detail route.

**Stack & constraints:**
- Go 1.22+, Chi v5, GORM, Redis (sorted set untuk queue, TTL untuk session)
- Nuxt 3, Vue 3, Pinia, Tailwind, `$fetch` dari Nuxt
- CORS sudah allow `*` di backend
- API contract = `apps/api/internal/router/router.go` adalah single source of truth

## Goals / Non-Goals

**Goals:**
- Sambungkan frontend ke backend end-to-end untuk flow: register → login → browse event → war → queue → booking → pay → e-ticket
- Normalisasi path & field name di frontend agar match dengan response backend
- Tambah payment backend (handler, repository, routes) sesuai spec `payment` yang sudah ada
- Tambah `GET /api/admin/bookings/{id}` untuk admin modal
- Verifikasi quick-start (docker compose) bisa menjalankan FE + BE + DB + seed data
- Zero data migration; tabel `payments` sudah ada

**Non-Goals:**
- Mengubah logika bisnis backend (sudah sesuai spec)
- Menambah fitur baru di UI (hanya perbaikan integrasi)
- Real payment gateway (tetap simulasi)
- Internationalization
- Optimasi performa tinggi (caching, CDN) — bisa di-iterasi selanjutnya
- Refactor besar di frontend (mis. migrasi ke Vue Query / TanStack)
- Real-time updates (WebSocket/SSE) — tetap polling

## Decisions

### 1. Single source of truth = `apps/api/internal/router/router.go` + handler signature

**Pilihan:** Frontend TypeScript types diselaraskan dengan struct JSON Go. Setiap endpoint yang dipanggil di FE harus bisa ditemukan di `router.go` dengan path & method yang persis sama.

**Alasan:**
- Menghilangkan inkonsistensi (`/api/...` vs tidak, `code` vs `booking_code`, dll.)
- Cukup scan `router.go` untuk verifikasi integrasi
- TypeScript types di `stores/*.ts` jadi single declaration point

**Alternatif:**
- Generate TypeScript types dari Go structs (via `go generate` + `oapi-codegen` atau `tygo`) — bagus tapi butuh tooling tambahan, di luar scope
- OpenAPI/Swagger spec — bisa ditambah nanti sebagai single source of truth

### 2. Path normalisasi: pakai grep + edit untuk fix FE

**Pilihan:** Edit file `pages/`, `stores/`, dan `composables/` di frontend yang masih pakai path/field salah, satu per satu, verifikasi dengan `grep`.

**Alasan:**
- Daftar file yang harus diperbaiki countable (sekitar 10 file)
- Lebih cepat dan deterministik daripada refactor global
- Bisa dilakukan paralel per-file

**Alternatif:**
- Tulis API client wrapper per-modul (`$api.events.list()`, `$api.bookings.reserve()`) — refactor besar, bisa di-iterasi

### 3. Payment backend: handler + repository terpisah, mengikuti pattern booking

**Pilihan:** Buat `apps/api/internal/handler/payment.go` dan `apps/api/internal/repository/payment.go` dengan API: `Create`, `GetByID`, `UpdateStatus`, `ListByUser`, `ListAll`, `MarkExpired`. Daftarkan di `router.go`.

**Alasan:**
- Konsisten dengan `booking.go` (yang sudah tested dan dipakai)
- Spec `payment` sudah jelas dari change sebelumnya
- Payment butuh transaksi (update booking status + generate e-ticket code) yang rapi di-handle di handler

**Alternatif:**
- Gabung payment ke `booking.go` handler — campur concern
- Pakai state machine library — overkill

### 4. Simulate endpoint tanpa autentikasi admin (test mode v1)

**Pilihan:** `POST /api/payments/{id}/simulate?status=success|fail` di-authenticate sebagai user (bukan admin), dan hanya pemilik payment yang boleh simulasi (mengikuti ownership check di `BookingHandler.Detail`).

**Alasan:**
- Spec v1 simulasi payment dipakai oleh user untuk demo
- Admin tidak perlu simulasi payment punya user lain
- Cukup validasi `payment.UserID == authenticated_user.ID`

**Alternatif:**
- Admin-only simulate — lebih aman tapi friction untuk demo
- Public simulate (siapa saja) — bahaya

### 5. E-ticket code generation: UUID per item, disimpan di `Booking.ETicketCodes []string`

**Pilihan:** Saat payment success, generate `len(booking.Items)` UUID, simpan ke `Booking.ETicketCodes` (JSONB column yang sudah ada). Frontend render list of codes.

**Alasan:**
- `Booking.ETicketCodes` sudah ada di model dengan scanner/valuer
- UUID v4 cukup unik dan pendek
- Konsisten dengan spec `payment` yang sudah ada (1 code per ticket item)

**Alternatif:**
- 1 code per booking (bukan per item) — lebih sederhana tapi kurang granular
- QR code image di server-side — out of scope v1 (placeholder di FE)

### 6. Frontend pay page: render instruksi sederhana + tombol "Bayar Sekarang" panggil simulate

**Pilihan:** `pages/bookings/[id]/pay.vue` panggil `POST /api/payments` saat mount untuk create payment (atau ambil existing pending), tampilkan amount, expiry countdown, dan tombol "Bayar Sekarang" yang panggil `POST /api/payments/{id}/simulate?status=success`. On success → redirect ke `/my/bookings/{id}`.

**Alasan:**
- Konsisten dengan spec v1 (simulasi payment)
- Countdown tetap pakai `CountdownTimer` existing dengan `expires_at` dari response
- Tidak butuh payment gateway UI yang kompleks

**Alternatif:**
- Integrasi Midtrans Snap — out of scope v1

### 7. CORS: keep `*` origins, no change

**Pilihan:** Tidak ubah CORS config di backend. Origin web `http://localhost:3000` sudah allowed via `*`.

**Alasan:**
- Sudah allow semua origins (cukup untuk dev)
- Untuk production bisa dipersempit, tapi di luar scope change ini

**Alternatif:**
- Whitelist origin spesifik — perubahan kecil tapi tidak wajib

### 8. Seed data: tambah `apps/api/cmd/seed/main.go` (jika belum ada) dengan 1 admin + 1 user + 2 events published + 3 kategori tiket

**Pilihan:** Tambah/cek `cmd/seed/main.go` agar bisa dijalankan via `make seed` atau `go run cmd/seed/main.go`. Buat data minimal untuk demo flow lengkap.

**Alasan:**
- Tanpa seed, demo perlu register + create event + add categories manual — friction
- Seed data menunjukkan expected shape response

**Alternatif:**
- SQL dump via migrations — lebih cepat tapi kurang fleksibel
- Factory-only (factory_girl) — Go tidak punya

### 9. Tidak tambah library frontend baru

**Pilihan:** Pakai `$fetch` + Pinia yang sudah ada. Tidak tambah TanStack Query, VueUse, dsb.

**Alasan:**
- Scope: integration, bukan architectural upgrade
- Kompleksitas yang ada cukup

**Alternatif:**
- TanStack Query untuk caching & retry otomatis — bagus tapi over-engineering untuk scope ini

## Risks / Trade-offs

- **Risiko:** Edit banyak file FE bisa introduce typo atau regresi visual kecil
  → **Mitigasi:** Test setiap flow end-to-end setelah edit; commit per file/area; jalankan `npm run build` untuk type-check

- **Risiko:** Backend payment handler baru bisa ada race condition (2 user bayar booking yang sama, atau user cancel saat admin simulasi)
  → **Mitigasi:** Pakai `SELECT ... FOR UPDATE` atau transaction di repository `payment.go` saat update status; double-check booking status sebelum simulasi

- **Risiko:** Generate e-ticket code di handler tanpa uniqueness check (UUID collision sangat kecil tapi mungkin)
  → **Mitigasi:** UUID v4 — probabilitas collision astronomis kecil; cukup untuk v1

- **Risiko:** Seed data bisa bentrok dengan data existing (unique constraint di email)
  → **Mitigasi:** Seed pakai email khusus (`admin@demo.test`, `user@demo.test`) dan `ON CONFLICT DO NOTHING` atau check dulu

- **Risiko:** Polling 2 detik di queue page bisa spam backend kalau 1000 user
  → **Mitigasi:** Out of scope — Redis sorted set di-handle efficient, dan rate limit per user sudah ada

- **Trade-off:** Path normalisasi "global edit" via grep mungkin skip edge case (path string concatenation, dsb.)
  → **Mitigasi:** Verifikasi manual setiap file yang diedit, dan grep ulang setelah edit

- **Trade-off:** Tambah payment backend menambah attack surface (simulate endpoint publicly accessible per-user)
  → **Mitigasi:** Ownership check (hanya pemilik payment boleh simulasi); batasi rate limit per payment ID; bisa tambah CAPTCHA di iterasi selanjutnya

- **Trade-off:** Tidak ada automatic retry di FE — kalau request gagal karena network blip, user harus refresh
  → **Mitigasi:** Tambah retry manual di critical flow (booking reserve) dengan konfirmasi; queue page sudah handle retry dengan `retryCount` & `maxRetries`

- **Risiko:** Format `date` (string ISO) vs `Date` object di FE bisa mismatch
  → **Mitigasi:** Selalu render lewat helper `formatDate()` yang handle `new Date(string)`; jangan asumsikan tipe

## Migration Plan

Tidak ada data migration (zero DB change selain menggunakan tabel `payments` yang sudah ada). Step rollout:

1. **Backend payment** — Tambah `repository/payment.go`, `handler/payment.go`, register routes di `router.go`. Verifikasi `go build` & `go test`.
2. **Backend admin booking detail** — Tambah `eventH.DetailAdmin` pattern: tambahkan `bookingH.GetByIDAdmin` di handler & register route.
3. **Frontend audit & fix** — Edit file FE satu per satu, mulai dari path/shape:
   - `pages/my/bookings.vue` → `/api/bookings/me`
   - `pages/my/bookings/[id].vue` → `/api/bookings/{id}` + field names fix
   - `pages/profile.vue` → `PUT /api/auth/me` + reload user
   - `pages/events/[id]/booking.vue` → kirim `session_token`, pakai `category_id`
   - `pages/events/[id]/queue.vue` → handle `is_ready` + `session_token`
   - `pages/events/[id]/war.vue` → pakai response shape `{ redirect_to_booking, queued, position, token }`
   - `pages/admin/events/[id].vue` → pakai field `booking_code` & `total_amount`
   - `pages/admin/bookings/index.vue` → pagination pakai `total/limit`, modal detail panggil `GET /api/admin/bookings/{id}`
   - `pages/admin/users/index.vue` → pagination fix
   - `pages/bookings/[id]/pay.vue` → implementasi payment flow
4. **Store updates** — `useQueueStore`, `useBookingStore` disesuaikan dengan backend shape.
5. **Seed** — Pastikan `cmd/seed/main.go` membuat data demo.
6. **Manual verification** — Run docker compose, register, browse event, war, queue, book, pay, lihat e-ticket; admin login, lihat stats, manage events.
7. **README** — Update quick-start di README jika ada perubahan command.

**Rollback strategy:** Karena tidak ada data migration, rollback cukup revert commit dan redeploy. Untuk backend, jika payment handler menyebabkan error, cukup disable route di `router.go` dan redeploy.

## Open Questions

- Apakah perlu refactor `useApi` agar support `baseURL` per-call (mis. untuk upload ke static file server berbeda)? — Kemungkinan tidak, semua lewat satu API.
- Apakah payment expiry background job perlu dibuat sekarang atau cukup pakai booking expiry job yang sudah ada? — Bisa pakai booking expiry job, payment otomatis expired via join.
- Apakah perlu rate limit spesifik untuk `/api/payments/{id}/simulate` (mencegah user simulasi berkali-kali)? — Bisa pakai middleware chi sederhana, prioritas rendah.
- Frontend perlu cara "refresh user" setelah update profile (saat ini `useAuthStore.updateMe` tidak ada) — tambah method atau cukup `fetchMe()` lagi?
- Apakah `pages/admin/payments/index.vue` perlu kolom "Action" (mis. simulasi dari sisi admin)? — Tabel read-only dulu sesuai spec.
