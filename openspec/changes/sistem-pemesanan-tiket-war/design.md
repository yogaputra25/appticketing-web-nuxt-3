## Context

Sistem pemesanan tiket konser/event dengan fitur "war tiket" menghadapi tantangan klasik: ribuan pengguna bersaing mendapatkan tiket terbatas dalam hitungan detik. Masalah utamanya: (1) double booking karena race condition, (2) server crash saat traffic spike, (3) pengalaman antrian yang tidak transparan. Solusi SaaS ini harus menjamin keadilan, integritas data, dan UX yang baik.

Stack teknologi:
- **Frontend**: Nuxt.js 3 (Vue 3) + Tailwind CSS + Pinia
- **Backend**: Go (Chi router) + GORM
- **Database**: PostgreSQL (transaksional) + Redis (antrian & cache)
- **Auth**: JWT

## Goals / Non-Goals

**Goals:**
- Sistem war tiket yang adil dengan antrian berbasis Redis
- Mencegah double booking dengan atomic operations (SELECT FOR UPDATE / optimistic locking)
- Soft reservation dengan TTL (5-10 menit) sebelum payment
- Performa tinggi: handle 1000+ concurrent users
- Real-time stock updates ke frontend
- Multi-tier pricing (VIP, Regular, dll)
- Dashboard admin untuk monitoring
- UI/UX modern dengan Tailwind CSS

**Non-Goals:**
- Real payment gateway integration (Xendit/Midtrans) di v1 — simulasi dulu
- Mobile app native (hanya web responsive)
- Social login (Google/FB) di v1
- E-ticket QR scanner untuk gate entry (future)

## Decisions

### 1. Database: PostgreSQL (bukan MongoDB)

**Pilihan: PostgreSQL**

**Alasan:**
- ACID compliance mutlak diperlukan untuk transaksi tiket (integritas stok & payment)
- `SELECT FOR UPDATE` (row-level locking) untuk mencegah double booking — fitur relational DB yang tidak dimiliki MongoDB
- Schema terstruktur (users, events, tickets, bookings) cocok untuk SQL
- Transaksi atomic di payment + reservation
- JSONB support tetap ada untuk fleksibilitas data fleksibel (event metadata)

**Alternatif yang dipertimbangkan:**
- **MongoDB**: Ditolak karena tidak punya transaksi multi-dokumen yang sama kuatnya, dan locking granular sulit. Stok tiket sangat relational.

### 2. War Tiket: Antrian Redis + Atomic Reservation

**Alasan:**
- Redis sorted set untuk queue (ZRANK untuk posisi antrian)
- Token-based queue (token sekali pakai, expire 5 menit)
- Atomic Lua script `DECRBY` di Redis untuk counter stok
- DB transaction finalisasi setelah payment

**Flow:**
1. User klik "War" → dapat queue token + posisi
2. User tunggu polling `/api/queue/status`
3. Saat giliran → redirect ke halaman booking 5 menit
4. POST `/api/booking/reserve` → atomic DB transaction: lock row ticket category, decrement stock, create booking dengan status `pending_payment` + expiry 10 menit
5. Payment → status `paid` atau expired (release stock)

### 3. Backend: Go + Chi + GORM

**Alasan:**
- Concurrency model Go (goroutines) ideal untuk ribuan request
- Chi router ringan, mudah di-middleware
- GORM untuk ORM dengan raw query support untuk locking
- Compile-time type safety

### 4. Frontend: Nuxt 3 + Tailwind + Pinia

**Alasan:**
- SSR Nuxt 3 untuk SEO & initial load cepat
- Tailwind untuk styling cepat dan konsisten
- Pinia untuk state management (auth, queue position, cart)
- API composition `useFetch` untuk komunikasi backend

### 5. Soft Reservation dengan TTL

**Alasan:**
- Mencegah monopolisasi stok saat payment
- Background job (cron) di Go release expired bookings → restore stock
- Redis TTL sebagai backup

### 6. Real-time Updates (Polling, bukan WebSocket di v1)

**Alasan:**
- Polling setiap 2-3 detik cukup untuk war tiket (queue & stock)
- Lebih sederhana, tidak perlu infrastructure WebSocket
- Bisa di-upgrade ke SSE/WebSocket di v2 jika perlu

## Risks / Trade-offs

- **High concurrent traffic pada saat war** → Mitigation: Redis rate limiting + Cloudflare/CDN + horizontal scaling Go
- **User refresh browser saat war** → Mitigation: token queue disimpan di httpOnly cookie + server-side session
- **DB bottleneck** → Mitigation: connection pooling, indexing pada ticket_categories.event_id, read replica di masa depan
- **Race condition di stok** → Mitigation: row-level lock + Redis atomic counter
- **Bot/auto-buy script** → Mitigation: CAPTCHA (Cloudflare Turnstile) di endpoint war + rate limit per IP/user
- **Stok tidak konsisten Redis vs Postgres** → Mitigation: Redis sebagai cache, DB sebagai source of truth; reconciliation job berkala

## Open Questions

- Apakah perlu integrasi payment gateway asli (Midtrans) di v1, atau simulasi dulu?
- Berapa TTL optimal untuk soft reservation? (default 10 menit, bisa dikonfigurasi)
- Apakah perlu multi-event organizer dalam 1 SaaS (multi-tenant), atau 1 organizer dulu?
