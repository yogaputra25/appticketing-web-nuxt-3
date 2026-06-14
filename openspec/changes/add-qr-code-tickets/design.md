## Context

Saat ini e-ticket code disimpan sebagai flat string list di `Booking.ETicketCodes`. Tidak ada entitas tiket individual, tidak ada QR code, dan tidak ada mekanisme verifikasi tiket secara on-site. User hanya bisa melihat booking yang berisi daftar kode e-ticket.

## Goals / Non-Goals

**Goals:**
- Buat model `Ticket` terpisah (satu row per tiket, bukan per quantity)
- Generate unique ticket code dan QR code untuk setiap tiket saat payment sukses
- API: daftar tiket user (`GET /api/tickets`), detail tiket (`GET /api/tickets/:id`), scan verifikasi (`POST /api/tickets/verify/:code`)
- Halaman My Tickets di dashboard user
- Halaman detail tiket dengan QR code
- Halaman scanner untuk verifikasi on-site
- QR code berisi URL verifikasi: `/tickets/v/<ticket_code>`

**Non-Goals:**
- Tidak ada PDF download tiket (cukup tampilan web)
- Tidak ada integrasi dengan hardware scanner fisik
- Tidak ada email notification tiket (belum)

## Decisions

### 1. Model: `Ticket` terpisah
- Current: `BookingItem` with `Quantity` + flat `ETicketCodes` list
- New: `Ticket` model — 1 row per tiket
- Fields: `ID`, `BookingID`, `TicketCode` (unique), `CategoryName`, `Status` (active/used/refunded), `ScannedAt`, `CreatedAt`
- Saat payment sukses, loop `BookingItem.Quantity` → create N `Ticket` rows
- **Rationale**: Relasi jelas per-tiket, memudahkan verifikasi individual dan tracking scan

### 2. QR code generation — frontend-side
- Backend hanya menyimpan `TicketCode` unik
- Frontend generate QR code image menggunakan library `qrcode` (npm)
- QR berisi URL: `{baseUrl}/tickets/v/{ticketCode}`
- **Rationale**: Tidak perlu tambahan dependency Go untuk generate image QR; frontend bisa generate sesuai kebutuhan UI

### 3. Verifikasi scan — endpoint backend
- `POST /api/tickets/verify/:code` — cek apakah `TicketCode` valid dan status `active`
- Jika valid → return 200 + data tiket, ubah status jadi `used` + set `ScannedAt`
- Jika invalid/used → return 404/400
- **Rationale**: Endpoint sederhana, bisa dipanggil dari halaman scanner atau langsung dari URL di QR

### 4. Halaman scanner — frontend
- Menggunakan `vue-qrcode-reader` (npm) untuk baca QR dari kamera
- Setelah scan, ekstrak `ticketCode` dari URL, panggil verify endpoint
- Tampilkan hasil valid/invalid di halaman
- **Rationale**: Library mature, Vue 3 compatible, works dengan Nuxt

## Risks / Trade-offs

- [Risk] Scanner perlu HTTPS untuk akses kamera → Mitigasi: pastikan dev server di localhost (allow), production wajib HTTPS
- [Risk] QR code berisi URL publik → Mitigasi: ticket code cukup unik dan tidak mengandung data sensitif; verify endpoint butuh auth (hanya petugas/scanner yang bisa akses)
- [Trade-off] Frontend generate QR lebih ringan untuk backend, tapi QR tidak bisa di-share sebagai file gambar langsung (harus dari halaman web). Untuk share, nanti bisa tambah generate PNG di sisi server jika diperlukan.
