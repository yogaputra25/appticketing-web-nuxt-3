## Why

Sistem pemesanan tiket konser/event dengan stok terbatas saat ini sering mengalami masalah: server crash saat traffic spike, pembelian ganda (double booking), dan tidak adanya mekanisme antrian yang adil. Sistem "war tiket" diperlukan untuk menangani lonjakan ribuan pengguna yang berebut tiket dalam hitungan detik secara adil dan transaksional. Ini adalah solusi SaaS yang bisa dipakai oleh berbagai event organizer.

## What Changes

- **Frontend**: Aplikasi Nuxt.js dengan Tailwind CSS untuk pengalaman pemesanan tiket yang responsif dan modern
- **Backend**: REST API Golang performa tinggi dengan concurrency model untuk menangani ribuan request simultan
- **Database**: PostgreSQL sebagai primary database (ACID compliance untuk transaksi tiket yang ketat)
- **Sistem War Tiket**: Mekanisme antrian (queue) dan locking untuk mencegah double booking
- **Auth System**: Registrasi, login, dan manajemen profil pengguna
- **Event Management**: CRUD event/konser oleh admin/event organizer
- **Ticket Inventory**: Manajemen stok tiket per kategori (VIP, Regular, dll)
- **Payment Integration**: Simulasi payment gateway dengan status tracking
- **Admin Dashboard**: Panel admin untuk monitoring penjualan dan stok tiket

## Capabilities

### New Capabilities

- `user-auth`: Autentikasi dan otorisasi pengguna — registrasi, login, JWT token, role-based access (user & admin)
- `event-management`: Pengelolaan event/konser oleh admin — membuat, mengedit, menjadwalkan event dengan detail seperti tanggal, venue, dan kategori tiket
- `ticket-inventory`: Manajemen inventori tiket dengan kuota per kategori, harga, dan status ketersediaan
- `ticket-war`: Sistem war tiket dengan antrian (queue), locking optimistik, dan timer reservation untuk menjamin keadilan — ini adalah core feature
- `payment`: Integrasi pembayaran — membuat order, memproses pembayaran, verifikasi status, dan expiry untuk pembayaran yang tidak diselesaikan
- `booking-management`: Riwayat dan manajemen pemesanan pengguna — melihat tiket yang sudah dibeli, status booking, dan e-ticket

### Modified Capabilities

<!-- Tidak ada existing capabilities yang dimodifikasi karena ini adalah project baru -->

## Impact

- **Frontend**: Project Nuxt.js baru dengan ekosistem Vue 3, Tailwind CSS, Pinia untuk state management
- **Backend**: Project Golang baru dengan framework HTTP (Chi/Gin), GORM atau SQLx untuk database, Redis untuk queue dan caching
- **Database**: PostgreSQL baru untuk data transaksional, Redis untuk antrian war tiket dan session
- **Infrastructure**: Docker Compose untuk development, potensi deploy ke cloud (VPS/AWS)
- **Third-party**: Midtrans/Xendit untuk payment gateway (simulasi di tahap awal)
