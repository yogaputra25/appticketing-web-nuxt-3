## Why

Aplikasi web Nuxt 3 saat ini menggunakan Tailwind CSS dan sudah memakai utility classes responsive di beberapa tempat, namun belum ada standar mobile-first yang konsisten. Akibatnya tampilan di perangkat mobile (≤ 640px) sering terpotong, navbar overflow, tabel admin tidak bisa di-scroll horizontal, modal terlalu besar, dan touch target terlalu kecil. Mengingat target pasar adalah pengguna yang banyak mengakses lewat smartphone saat akan war tiket, tampilan mobile yang buruk menurunkan konversi dan merusak pengalaman "antrian adil" yang menjadi nilai jual utama produk.

## What Changes

- **Mobile-first design system**: Audit & refactor `tailwind.config.ts` dengan breakpoint default Tailwind (sm/md/lg/xl/2xl), tambah safe-area inset untuk device ber-notch, dan utility custom untuk touch target minimal 44px.
- **Layout & Navigasi**: Buat ulang `layouts/default.vue` dan `layouts/admin.vue` dengan navbar collapsible (hamburger menu) untuk mobile, drawer navigation untuk admin, dan sticky bottom action bar di halaman kritis (war, booking, payment).
- **Komponen reusable mobile-friendly**: Refactor `EventCard`, `TicketCategoryCard`, `QueueStatus`, `CountdownTimer` agar responsif (stack di mobile, grid di desktop).
- **Halaman publik**: Optimasi `pages/index.vue`, `pages/events/index.vue`, `pages/events/[id].vue`, `pages/events/[id]/war.vue`, `pages/events/[id]/queue.vue`, `pages/events/[id]/booking.vue` dengan padding/viewport yang sesuai mobile.
- **Halaman auth**: `pages/login.vue` & `pages/register.vue` full-width di mobile dengan input field yang touch-friendly.
- **Halaman admin**: `pages/admin/*` — tabel responsive dengan horizontal scroll, card view alternatif di mobile, filter collapse jadi accordion.
- **Halaman user**: `pages/my/bookings.vue`, `pages/my/bookings/[id].vue`, `pages/profile.vue` — list view vertical di mobile, modal/popup bottom-sheet style.
- **Utilitas bersama**: Tambah composable `useViewport()` untuk deteksi breakpoint, komponen `MobileBottomSheet.vue`, `MobileTabs.vue`, dan helper CSS untuk truncate multi-line.
- **Testing & verification**: Manual test di Chrome DevTools responsive mode untuk breakpoint 320px, 375px, 414px, 768px, 1024px, 1280px.

## Capabilities

### New Capabilities

- `responsive-web-ui`: Standar dan komponen untuk memastikan semua halaman web tampil optimal di mobile (≤ 640px) dan tetap bagus di tablet/desktop, dengan fokus pada touch target, layout adaptif, dan navigasi yang ramah mobile.

### Modified Capabilities

<!-- Tidak ada capability backend yang berubah — perubahan murni UI/frontend.
     Requirements frontend (web) sebelumnya tersebar tanpa standar; sekarang
     dikonsolidasikan ke capability baru `responsive-web-ui`. -->

## Impact

- **Frontend**: File-file di `apps/web/` (pages, layouts, components, composables, assets/css, tailwind.config.ts). Tidak ada perubahan backend.
- **Dependencies**: Tetap pada Tailwind CSS yang sudah ada; tidak menambah library baru. Memanfaatkan `@nuxtjs/tailwindcss` yang sudah terpasang.
- **No backend changes**: API contract tidak berubah. Ini murni peningkatan UI.
- **Browser support**: Modern browsers (Chrome/Safari/Firefox/Edge versi 2 tahun terakhir) dengan viewport meta `width=device-width, initial-scale=1`.
