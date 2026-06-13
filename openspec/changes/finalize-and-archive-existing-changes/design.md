## Context

Dua perubahan (`integrate-backend-api-to-frontend` dan `fix-chi-route-ordering-and-verify-all-apis`) sudah selesai implementasi. Tinggal task verifikasi manual dan archive yang belum dijalankan.

## Goals / Non-Goals

**Goals:**
- Verifikasi end-to-end flow user (register → login → war → booking → pay → e-ticket)
- Verifikasi admin flow (login → manage events → bookings → payments)
- Verifikasi my bookings (list, detail, cancel)
- Cek console browser & network tab
- Jalankan smoke test script
- Commit perubahan per area
- Archive kedua perubahan

**Non-Goals:**
- Tidak ada perubahan kode baru
- Tidak ada perubahan requirement spec

## Decisions

- Verifikasi dilakukan manual oleh user di browser (flow end-to-end tidak bisa diotomatisasi penuh)
- Smoke test API bisa dijalankan via script `apps/api/scripts/smoke.sh` terhadap Docker container yang sedang berjalan
- Commit akan dilakukan per area fungsional (backend payment, frontend path fix, pay page, admin payments) untuk memudahkan rollback

## Risks / Trade-offs

- Docker container mungkin mati — perlu `docker compose up -d postgres redis api` sebelum verifikasi
- Smoke test membutuhkan API running di `http://localhost:8081` — pastikan container aktif
