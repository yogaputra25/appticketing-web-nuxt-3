# Sistem Pemesanan Tiket (War Tiket)

SaaS pemesanan tiket konser/event dengan fitur **war tiket** — antrian adil untuk tiket terbatas.

## Stack

- **Frontend**: Nuxt 3 (Vue 3) + Tailwind CSS + Pinia
- **Backend**: Go 1.22+ + Chi router + GORM
- **Database**: PostgreSQL 16 + Redis 7
- **Auth**: JWT

## Struktur Monorepo

```
appticketing/
├── apps/
│   ├── web/          # Nuxt 3 frontend
│   └── api/          # Go backend
├── docs/             # Dokumentasi tambahan
├── openspec/         # OpenSpec change artifacts
├── docker-compose.yml
├── .env.example
└── README.md
```

## Quick Start

```bash
# 1. Copy env
cp .env.example .env

# 2. Start infrastructure (postgres + redis)
docker compose up -d postgres redis

# 3. Run migrations
cd apps/api && make migrate-up

# 4. Start API
cd apps/api && go run cmd/server/main.go

# 5. Start web (terminal baru)
cd apps/web && npm install && npm run dev
```

Frontend: http://localhost:3000  
Backend: http://localhost:8080

## Smoke Test

```bash
cd apps/api && make smoke
```

Jalankan setelah `docker compose up -d` untuk verifikasi semua endpoint
ter-registrasi dengan benar. Membutuhkan API yang sedang berjalan.

## Dokumentasi

Detail desain & specs ada di `openspec/changes/sistem-pemesanan-tiket-war/`.

### Router Patterns

Lihat `docs/router.md` untuk penjelasan pola routing Chi (no-fallthrough rule).

### Troubleshooting

Jika mendapat response `404 page not found` (plain text, 19 bytes) padahal
route sudah ada di `apps/api/internal/router/router.go`, kemungkinan route
tersebut dideklarasikan di subroute yang salah. Pastikan event-scoped routes
(seperti `/admin/events/{id}/categories`) berada di dalam subroute
`/admin/events`, bukan di subroute `/admin` yang lebih luas.
Lihat `docs/router.md` untuk detail pola router.
