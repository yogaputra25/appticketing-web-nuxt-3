## Context

Project ini menggunakan Nuxt 3 (Vue 3) + Tailwind CSS untuk frontend di `apps/web/`. Saat ini halaman-halaman web sudah menggunakan utility class Tailwind di beberapa tempat, namun belum ada standar mobile-first yang konsisten. Pengujian manual di DevTools menunjukkan banyak masalah pada viewport ≤ 640px: navbar meluap, tabel admin tidak bisa di-scroll dengan baik, modal terlalu besar, form input terlalu kecil untuk tap, dan konten penting war tiket (countdown, posisi antrian, tombol bayar) tidak di permukaan mobile.

Stack & kondisi saat ini:
- **Nuxt 3** dengan `@nuxtjs/tailwindcss` (sudah terpasang)
- **Tailwind CSS** dengan konfigurasi default (breakpoint bawaan)
- **Pinia** untuk state management
- Halaman publik (`pages/index.vue`, `pages/events/*`, `pages/login.vue`, `pages/register.vue`) dan admin (`pages/admin/*`) sudah ada, sebagian di-mark `done` di tasks.md
- Backend (Go) dan API contract **tidak berubah**

Stakeholder: pengguna akhir yang banyak mengakses dari smartphone saat ingin war tiket, admin event yang mengelola event dari laptop/tablet, dan developer yang akan memaintain halaman.

## Goals / Non-Goals

**Goals:**
- Standar mobile-first untuk seluruh halaman web dengan breakpoint Tailwind default
- Touch target minimum 44×44px untuk semua interaksi di mobile
- Navigasi yang adaptif (hamburger di mobile, inline di desktop)
- Tabel admin yang usable di mobile (horizontal scroll dalam container atau card view)
- Sticky CTA di halaman kritis (war, queue, booking, payment)
- Safe area inset untuk device ber-notch
- Body text ≥ 16px (mencegah iOS auto-zoom)
- Komponen bersama: `MobileBottomSheet`, `MobileTabs`, composable `useViewport`
- Zero perubahan backend / API

**Non-Goals:**
- Mengubah API contract
- Menambah library CSS/UI baru (Headless UI, Naive UI, dsb.) — pakai Tailwind utility saja
- Mengimplementasikan PWA (installable, service worker, offline mode)
- Membuat mobile app native
- Redesain visual menyeluruh (warna, tipografi, branding) — fokus pada layout & interaksi
- Internationalization (i18n) — di luar scope
- Mengubah layout admin sidebar menjadi sesuatu yang sangat berbeda dari versi desktop

## Decisions

### 1. Pendekatan: Mobile-first dengan Tailwind default breakpoints

**Pilihan:** Tetap pakai Tailwind default breakpoints (`sm: 640px`, `md: 768px`, `lg: 1024px`, `xl: 1280px`, `2xl: 1536px`). Tulis utility dari `base` (mobile) naik, bukan turun.

**Alasan:**
- Tidak menambah learning curve
- Konsisten dengan ekosistem Tailwind
- Project ini audience Indonesia yang banyak pakai smartphone, jadi mobile-first = default yang tepat
- Breakpoint default sudah cukup untuk kasus kita (320px iPhone SE → 1920px desktop)

**Alternatif:**
- Container queries (`@container`) — terlalu baru, belum banyak dipakai, Tailwind support masih experimental
- Custom breakpoints — menambah konfigurasi tanpa kebutuhan nyata

### 2. Komponen `MobileBottomSheet` custom, bukan library eksternal

**Pilihan:** Buat komponen `MobileBottomSheet.vue` sendiri dengan `<Teleport to="body">` + transisi Tailwind.

**Alasan:**
- Use case terbatas (filter sheet, konfirmasi aksi)
- Menghindari dependency baru
- Bundle size lebih kecil
- Bisa disesuaikan dengan styling kita

**Alternatif:**
- Headless UI / Radix Vue — overkill untuk 1-2 use case
- Native `<dialog>` — support browser belum konsisten di iOS lama

### 3. Hamburger menu pakai `<details>` + state, bukan JS state manager

**Pilihan:** Drawer navigasi mobile di-handle dengan state Vue reaktif biasa (`ref<boolean>`), buka/tutup via tombol, dengan `v-on:keydown.esc` untuk aksesibilitas.

**Alasan:**
- Sederhana, tidak perlu Pinia
- State lokal cukup (cuma boolean)
- Esc-to-close gratis dengan keydown handler

**Alternatif:**
- Pinia store — overkill
- `<dialog>` element dengan `showModal()` — lumayan tapi Tailwind transitions lebih fleksibel

### 4. Sticky CTA pakai `fixed bottom-0` di mobile saja, non-sticky di desktop

**Pilihan:** Bottom action bar di-hide di `md:` ke atas (cukup ruang), dan hanya muncul di mobile.

**Alasan:**
- Di desktop konten lebih lebar, sticky bottom sering menutupi konten
- Di mobile, sticky bottom adalah pattern standar (iOS/Android)

**Alternatif:**
- Sticky di semua viewport — menutupi konten di desktop

### 5. Tabel admin: horizontal scroll di mobile, bukan card view terpisah

**Pilihan:** Bungkus `<table>` dengan container `overflow-x-auto` di mobile, dan tambahkan `min-w-full` ke table.

**Alasan:**
- Implementasi paling cepat, satu sumber data
- Mempertahankan struktur & styling tabel
- Card view alternatif butuh duplikasi data binding

**Alternatif:**
- Card view khusus mobile — banyak duplikasi, sync state lebih sulit
- DataTables / AG Grid — dependency besar, overkill

### 6. `useViewport()` composable pakai `matchMedia`, bukan window resize listener

**Pilihan:** `useViewport()` menggunakan `window.matchMedia` dengan reactive cleanup.

**Alasan:**
- Lebih performant (browser optimize matchMedia)
- Tidak fire di setiap pixel resize
- Reactive di Vue 3 dengan `ref` + `addEventListener('change', ...)` di `onMounted`

**Alternatif:**
- `@vueuse/core` useMediaQuery — menambah dependency (lihat decision 7)

### 7. Tidak menambah library baru

**Pilihan:** Tetap pada dependency yang ada: Nuxt 3, Tailwind, Pinia. Tidak menambah Headless UI, VueUse, dll.

**Alasan:**
- `package.json` saat ini minimalis
- Use case kita cukup sederhana untuk di-handle dengan Vue + Tailwind saja
- Mengurangi attack surface & maintenance

**Alternatif:**
- `@vueuse/core` untuk `useMediaQuery`, `useScrollLock`, dsb. — useful tapi tidak wajib
- Headless UI untuk accessible modal/menu — nice-to-have, bisa ditambah nanti

### 8. Safe area pakai CSS custom property & utility Tailwind

**Pilihan:** Tambah custom utility di `tailwind.config.ts`:
```ts
padding: { 'safe': 'env(safe-area-inset-top)' }
```

**Alasan:**
- Standar W3C CSS env()
- Tailwind bisa extend dengan plugin sederhana
- Single source of truth

**Alternatif:**
- Inline `env()` di setiap file — duplikasi

## Risks / Trade-offs

- **Risiko:** Per-page audit membutuhkan waktu dan ada risiko regresi visual di desktop
  → **Mitigasi:** Test di semua breakpoint (320, 375, 414, 768, 1024, 1280, 1536) untuk setiap halaman yang diubah; gunakan Chrome DevTools responsive mode + real device bila ada

- **Risiko:** Perubahan global (mis. font-size default) bisa mempengaruhi halaman yang sudah jadi
  → **Mitigasi:** Mulai dari halaman yang paling terdampak (war, queue, booking, payment, admin tables), kerjakan per-halaman, commit kecil

- **Risiko:** `useViewport()` bisa mismatch saat SSR vs client (hydration)
  → **Mitigasi:** Default ke `isMobile = false` di SSR, baru update di `onMounted` untuk hindari hydration mismatch; atau gunakan `ClientOnly` wrapper di tempat yang butuh

- **Risiko:** iOS Safari kadang punya quirk khusus (100vh, rubber band, input zoom)
  → **Mitigasi:** Set `font-size: 16px` di semua input, gunakan `min-h-screen` (bukan `h-screen`) untuk full-height layout, test di iOS Safari kalau memungkinkan

- **Risiko:** Hamburger menu drawer tidak bisa di-scroll jika konten panjang
  → **Mitigasi:** `overflow-y-auto` di drawer, `overscroll-contain` untuk mencegah scroll chaining

- **Trade-off:** Tidak pakai library komponen = lebih banyak kode custom. Bisa lebih lama untuk edge cases (focus trap, ARIA) di modal/drawer
  → **Mitigasi:** Implementasi minimum accessible: `role="dialog"`, `aria-modal="true"`, focus ke drawer saat buka, esc-to-close, restore focus saat tutup

- **Trade-off:** Sticky bottom bar menambah tinggi viewport pada konten
  → **Mitigasi:** Tambah bottom padding yang setara pada container konten (`pb-20 md:pb-0`)

## Migration Plan

Tidak ada migration data atau backend. Rollout frontend-only:

1. **Tahap 1 — Fondasi:** Update `tailwind.config.ts` (tambah custom utilities), tambah `assets/css/safe-area.css` (atau di main css), buat `composables/useViewport.ts`
2. **Tahap 2 — Komponen bersama:** `MobileBottomSheet.vue`, `MobileTabs.vue`
3. **Tahap 3 — Layout & navigasi:** Refactor `layouts/default.vue` & `layouts/admin.vue` dengan hamburger/drawer
4. **Tahap 4 — Halaman publik:** Audit & fix `pages/index.vue`, `pages/events/*`, `pages/login.vue`, `pages/register.vue`
5. **Tahap 5 — Halaman user:** `pages/my/bookings.vue`, `pages/my/bookings/[id].vue`, `pages/profile.vue`
6. **Tahap 6 — Halaman admin:** `pages/admin/*` (events, bookings, payments, users)
7. **Tahap 7 — Halaman kritis war:** `pages/events/[id]/war.vue`, `queue.vue`, `booking.vue`, `pages/bookings/[id]/pay.vue` — sticky CTA, countdown legibility
8. **Tahap 8 — Verification:** Manual test di DevTools, fix edge cases

**Rollback strategy:** Karena perubahan frontend-only dan tidak menyentuh backend/database, rollback cukup dengan revert commit dan redeploy. Tidak ada data migration.

## Open Questions

- Apakah admin diizinkan menggunakan mobile untuk operasional harian, atau admin selalu desktop? (Implikasi: seberapa agresif optimasi mobile di admin?)
- Apakah perlu dark mode? (Bisa berdampak pada kontras & mobile readability — out of scope untuk change ini tapi worth noting)
- Apakah tim punya akses ke real device iOS/Android untuk testing, atau hanya DevTools? (Mempengaruhi confidence level)
- Untuk sticky bottom bar di payment, apakah perlu integrasi dengan payment timer (countdown) — apakah timer ditampilkan di bar juga, atau hanya di konten?
