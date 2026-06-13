## Why

Dua perubahan besar (`integrate-backend-api-to-frontend` dan `fix-chi-route-ordering-and-verify-all-apis`) sudah selesai diimplementasikan tetapi masih ada task manual verification yang belum dijalankan. Kita perlu memverifikasi end-to-end flow, menjalankan smoke test, dan mengarsipkan kedua perubahan agar status project bersih dan fokus ke perubahan selanjutnya.

## What Changes

- **Verifikasi manual** — Jalankan task manual verification yang tersisa di `integrate-backend-api-to-frontend` (10.1-10.5, 10.7, 10.8)
- **Smoke test verification** — Jalankan smoke test di `fix-chi-route-ordering-and-verify-all-apis` (4.5, 4.6)
- **Commit per area** — Commit hasil perbaikan backend payment, frontend path fix, pay page, admin payments
- **Archive** — Archive kedua perubahan yang sudah selesai

## Capabilities

### New Capabilities
<!-- Tidak ada capability baru — ini adalah perubahan proses finalisasi -->

### Modified Capabilities
<!-- Tidak ada requirement spec yang berubah — hanya verifikasi dan archive -->

## Impact

- Tidak ada perubahan kode — hanya verifikasi manual dan smoke test
- Status project: kedua perubahan akan diarsipkan setelah selesai
