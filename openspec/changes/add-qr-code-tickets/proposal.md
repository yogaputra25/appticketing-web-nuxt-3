## Why

Tiket yang sudah dibeli belum memiliki QR code untuk verifikasi saat masuk event. Tanpa QR code, tidak ada cara mudah untuk memvalidasi keaslian tiket secara offline/on-site. Perlu sistem tiket dengan QR code unik per-tiket dan halaman scan untuk verifikasi.

## What Changes

- Buat model `Ticket` baru (satu baris per tiket, bukan per quantity) dengan unique ticket code dan QR code
- Generate unik tiket beserta QR code-nya saat pembayaran dikonfirmasi
- API endpoint baru untuk: daftar tiket user, detail tiket, verifikasi tiket via scan
- Halaman "My Tickets" di dashboard user yang menampilkan daftar tiket user
- Halaman detail tiket dengan QR code yang bisa discan
- Halaman scanner (scan QR code) untuk mencocokkan tiket — valid/invalid
- QR code berisi unique ticket identifier (bisa URL endpoint verifikasi)

## Capabilities

### New Capabilities
- `ticket-qr`: Tiket digital dengan QR code, user ticket list, QR scan verification

### Modified Capabilities
- `booking-management`: E-ticket generation berubah — tiket dibuat sebagai `Ticket` rows baru dengan QR code, bukan flat string list di booking

## Impact

- `apps/api/internal/model/`: New `Ticket` model + migrations
- `apps/api/internal/repository/`: New `TicketRepository`; update `BookingRepository` untuk create tickets
- `apps/api/internal/handler/`: New `TicketHandler` untuk list/detail/verify; update `PaymentHandler` untuk generate tickets
- `apps/api/internal/router/`: New routes for tickets
- `apps/web/`: New pages: My Tickets, Ticket Detail, Scanner
- Dependencies: library QR code generation (`github.com/skip2/go-qrcode`), library QR scanner untuk frontend (`vue-qrcode-reader`)
