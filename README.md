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

## Dokumentasi

Detail desain & specs ada di `openspec/changes/sistem-pemesanan-tiket-war/`.
