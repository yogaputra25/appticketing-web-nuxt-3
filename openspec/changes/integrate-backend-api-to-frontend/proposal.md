## Why

Backend API Go (Chi + GORM) di `apps/api/` sudah selesai untuk modul auth, event, ticket category, war queue, dan booking (lihat `apps/api/internal/router/router.go` & `handlers/`). Frontend Nuxt 3 di `apps/web/` sudah membangun halaman-halaman UI dan kebanyakan sudah memakai `useApi()` / `useEventStore()` / `useQueueStore()` / `useBookingStore()` untuk memanggil endpoint, namun ada tiga masalah utama yang menghalangi demo end-to-end:

1. **Path & field name tidak konsisten** — beberapa halaman pakai prefix `/api/...` (mis. `/api/events`), sebagian lain pakai path tanpa prefix (mis. `/bookings`, `/auth/profile`). Sebagian page admin pakai `b.id?.slice(0, 8)` padahal `id` adalah `uint64`, dan ada halaman yang membaca `b.code` padahal model hanya punya `booking_code`. Akibatnya response backend tidak match dengan interface TypeScript frontend.
2. **Endpoint yang dipanggil belum ada di backend** — `pages/bookings/[id]/pay.vue` masih placeholder, dan `pages/my/bookings.vue`/`pages/my/bookings/[id].vue` memanggil endpoint yang tidak ada di router. Begitu juga `pages/admin/bookings/[id].vue` modal "Detail" yang memakai `GET /api/admin/bookings/{id}`.
3. **Store & composable shape drift** — `useQueueStore` & `useBookingStore` di `apps/web/stores/` mengembalikan field dengan nama berbeda dari response backend (`position` vs `is_ready`, `expires_at` vs `expires_at`, `category_id` vs `ticket_category_id`, dsb.).

Tanpa integrasi end-to-end, user tidak bisa melihat hasil nyala (event muncul, login bekerja, war flow berhasil, booking tercatat di DB). Change ini menyambungkan frontend dengan backend yang sudah ada, sekaligus melengkapi endpoint backend yang missing (payments, admin booking detail) sehingga demo end-to-end bisa dijalankan.

## What Changes

### Frontend fixes (path/shape alignment)

- **Normalisasi path API** — Semua panggilan di `pages/`, `components/`, dan `stores/` pakai prefix `/api/...` konsisten; ubah `pages/my/bookings.vue` dari `/bookings` ke `/api/bookings/me`, `pages/my/bookings/[id].vue` dari `/bookings/{id}` ke `/api/bookings/{id}`, `pages/profile.vue` dari `/auth/profile` ke `/api/auth/me`, dan `pages/admin/events/[id].vue` dll.
- **Perbaiki shape field di store & page** — Samakan nama field dengan response backend: `Booking` punya `total_amount` (bukan `total_price`), `booking_code` (bukan `code`), `e_ticket_codes: string[]` (bukan `ticket_count`); `QueueStatus` backend punya `is_ready`, `session_token`, `position`, `total_in_queue`, `estimated_wait` — update `useQueueStore` & `pages/events/[id]/queue.vue`.
- **Perbaiki field `id` yang dipotong** — `Booking.ID` adalah `number` (uint64) bukan string; hapus `b.id?.slice(0, 8)` dan pakai `b.booking_code` atau `String(b.id)`.
- **Perbaiki `pages/events/[id]/booking.vue`** — body `POST /api/bookings/reserve` harus berisi `event_id`, `session_token`, dan `items: [{category_id, quantity}]` (bukan `ticket_category_id`).
- **Perbaiki `pages/events/[id]/queue.vue`** — pakai `session_token` dari response `war/status` saat `is_ready=true`, dan `position`/`total_in_queue` saat menunggu.
- **Perbaiki `pages/admin/bookings/index.vue` & `pages/admin/users/index.vue`** — endpoint backend mengembalikan `{data, total, page, limit}` (bukan `last_page`); update pagination logic.

### Backend endpoint additions (yang kurang di frontend)

- **Tambah `GET /api/admin/bookings/{id}`** untuk modal Detail di admin (mengikuti pattern `DetailAdmin` di `event.go`).
- **Tambah modul Payment** — `POST /api/payments` untuk create payment, `POST /api/payments/{id}/simulate` (test mode), `GET /api/payments/me` (user), `GET /api/admin/payments` (admin) — agar `pages/bookings/[id]/pay.vue` & `pages/admin/payments/index.vue` benar-benar bekerja end-to-end. Model `Payment` sudah ada di DB (`payments` table) dan repository dasar sudah ada, hanya handler + route yang belum ditambahkan.
- **Tambah `GET /api/bookings/{id}/payment`** (opsional, untuk pay page) — mengembalikan instruksi pembayaran & status payment.

### Connection & verification

- **Pastikan `apiBase` di `nuxt.config.ts`** mengarah ke backend (`http://localhost:8080` untuk dev) dan `CORS` di backend allow origin Nuxt.
- **Tambah `.env.example` untuk web** (jika belum) — `NUXT_PUBLIC_API_BASE`.
- **Tambah `docker-compose.yml` entry** untuk menjalankan web + api + DB bersamaan (jika belum), atau dokumentasikan quick-start yang sudah ada.
- **Tambah script `npm run seed`** (jika belum) untuk populate data awal: 1 admin, 1 user, 1-2 event published dengan kategori tiket.

## Capabilities

### New Capabilities

- `frontend-api-integration`: Standar dan helper untuk memastikan seluruh halaman web (publik, user, admin) memakai endpoint backend yang benar dengan path, method, dan field name yang konsisten — berdasarkan kontrak `apps/api/internal/handler/*` dan `apps/api/internal/router/router.go`.
- `payment`: Integrasi pembayaran — create payment untuk booking, simulasi payment (success/fail), ambil payment user, monitoring payment oleh admin, dan payment expiry (mirip dengan `payment` capability yang sudah ada di change `sistem-pemesanan-tiket-war` tapi saat ini belum diimplementasikan di backend).

### Modified Capabilities

- `booking-management`: Halaman `/my/bookings` & `/my/bookings/[id]` harus memakai endpoint `GET /api/bookings/me` & `GET /api/bookings/{id}` (sudah ada) dengan shape `Booking` model. Backend perlu tambah `GET /api/admin/bookings/{id}` (untuk admin modal).
- `ticket-war`: Frontend `useQueueStore` & `pages/events/[id]/queue.vue` harus memakai response `{is_ready, session_token, position, total_in_queue, estimated_wait}` dari `GET /api/war/status`. Frontend `pages/events/[id]/booking.vue` harus mengirim `session_token` saat POST `/api/bookings/reserve`.

## Impact

- **Frontend** (`apps/web/`): file di `pages/`, `components/`, `stores/`, `composables/` — path & shape alignment, tanpa menambah fitur UI baru.
- **Backend** (`apps/api/`): tambah handler `payment.go` + 4 endpoint, tambah 1 admin endpoint `GET /api/admin/bookings/{id}`, register route di `router.go`.
- **Database**: tidak ada migration baru — tabel `payments` sudah ada. Model & repository `Payment` belum lengkap (perlu dicek `apps/api/internal/repository/payment.go`); jika belum ada, tambahkan repository mengikuti pattern `repository/booking.go`.
- **Config**: `nuxt.config.ts` runtimeConfig untuk `apiBase` (sudah ada `useRuntimeConfig().public.apiBase`).
- **CORS**: backend sudah allow `*` origins — tidak perlu diubah.
- **Backward compatibility**: tidak ada breaking change pada user-facing flow; perubahan shape difasilitasi dengan mengupdate store & page bersesuaian.
