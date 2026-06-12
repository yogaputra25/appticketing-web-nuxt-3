## 1. Foundation: Tailwind Config & Global CSS

- [x] 1.1 Update `apps/web/tailwind.config.ts` — pastikan breakpoints default aktif, tambah custom spacing untuk `safe` (`pt-safe`, `pb-safe`, `pl-safe`, `pr-safe`) berbasis `env(safe-area-inset-*)`
- [x] 1.2 Tambah global CSS di `apps/web/assets/css/main.css` — set `html { font-size: 16px }`, `body` base typography, dan `min-h-screen` helpers
- [x] 1.3 Tambah viewport meta di `apps/web/app.vue` atau `nuxt.config.ts` (jika belum): `width=device-width, initial-scale=1, viewport-fit=cover`
- [x] 1.4 Tambah plugin Tailwind atau utility kustom untuk `touch-target` (`min-h-[44px] min-w-[44px]`) dan `clamp-line-2` / `clamp-line-3` (truncate multi-line)
- [x] 1.5 Verifikasi tidak ada class `h-screen` yang harus diganti `min-h-screen` (iOS Safari quirk)

## 2. Shared Composables & Components

- [x] 2.1 Buat `apps/web/composables/useViewport.ts` — return reactive `{ isMobile, isTablet, isDesktop }` berbasis `matchMedia`; default SSR-safe
- [x] 2.2 Buat `apps/web/components/MobileBottomSheet.vue` — teleport ke body, slide-up di mobile, centered modal di `md:`, esc-to-close, focus trap sederhana
- [x] 2.3 Buat `apps/web/components/MobileTabs.vue` — horizontal scrollable tabs di mobile, centered di desktop
- [x] 2.4 Buat `apps/web/components/StickyBottomBar.vue` — slot wrapper dengan padding-bottom di konten saat visible di mobile
- [x] 2.5 Test komponen baru di Storybook-style page atau di DevTools dengan viewport switching

## 3. Layouts: Default & Admin

- [x] 3.1 Refactor `apps/web/layouts/default.vue` — navbar dengan logo + hamburger button (visible `md:hidden`), drawer navigasi di mobile, inline nav di `md:`, sticky user menu
- [x] 3.2 Pastikan footer full-width dengan padding yang sesuai mobile
- [x] 3.3 Refactor `apps/web/layouts/admin.vue` — sidebar collapse jadi drawer di mobile, hamburger toggle, top bar dengan user info
- [x] 3.4 Test layout di viewport 320px, 375px, 768px, 1280px

## 4. Public Pages: Home, Events, Auth

- [x] 4.1 Fix `apps/web/pages/index.vue` — hero section stack di mobile, event card grid 1 kolom (mobile) / 2 (sm) / 3 (md) / 4 (lg)
- [x] 4.2 Fix `apps/web/pages/events/index.vue` — search bar full-width mobile, filter sebagai bottom sheet atau accordion, card list view
- [x] 4.3 Fix `apps/web/pages/events/[id].vue` — banner full-width responsive, info grid stack di mobile, ticket categories card dengan sticky CTA
- [x] 4.4 Fix `apps/web/pages/login.vue` — form full-width mobile, input ≥ 44px tinggi, error message accessible
- [x] 4.5 Fix `apps/web/pages/register.vue` — sama seperti login, tambah password strength meter yang responsive

## 5. User Pages: Bookings & Profile

- [x] 5.1 Fix `apps/web/pages/my/bookings.vue` — list view vertical (cards) di mobile, table di desktop; status badge color-coded
- [x] 5.2 Fix `apps/web/pages/my/bookings/[id].vue` — e-ticket detail stack di mobile, QR code di tengah dengan caption, action buttons full-width
- [x] 5.3 Fix `apps/web/pages/profile.vue` — form full-width, field stack vertical di mobile, save button sticky bottom di mobile

## 6. Admin Pages

- [x] 6.1 Fix `apps/web/pages/admin/index.vue` — stat cards grid 1 (mobile) / 2 (sm) / 4 (lg), chart responsive
- [x] 6.2 Fix `apps/web/pages/admin/events/index.vue` — tabel dengan `overflow-x-auto` wrapper, filter jadi bottom sheet di mobile, "Add Event" FAB atau top button
- [x] 6.3 Fix `apps/web/pages/admin/events/new.vue` dan `pages/admin/events/[id].vue` — form layout responsive, ticket category manager sebagai accordion di mobile
- [x] 6.4 Fix `apps/web/pages/admin/bookings/index.vue` — sama pattern: scrollable table + filter sheet
- [x] 6.5 Fix `apps/web/pages/admin/payments/index.vue` — sama pattern
- [x] 6.6 Fix `apps/web/pages/admin/users/index.vue` — sama pattern

## 7. Critical War Tiket Flow

- [x] 7.1 Fix `apps/web/pages/events/[id]/war.vue` — countdown prominent (font ≥ 32px), info event ringkas, tombol "Mulai War" full-width di mobile dengan min-height 56px
- [x] 7.2 Fix `apps/web/pages/events/[id]/queue.vue` — posisi antrian & total di atas fold, progress bar visible, polling indicator, sticky bar dengan tombol "Keluar"
- [x] 7.3 Fix `apps/web/pages/events/[id]/booking.vue` — quantity selector dengan tombol +/- yang touch-friendly, summary card, sticky bottom "Lanjut Bayar" di mobile
- [x] 7.4 Fix `apps/web/pages/bookings/[id]/pay.vue` — payment instructions readable, countdown timer prominent, tombol "Bayar Sekarang" sticky bottom

## 8. Existing Component Audit

- [x] 8.1 Refactor `apps/web/components/EventCard.vue` — gambar `aspect-ratio`, judul clamp 2 baris, harga & tanggal stack di mobile
- [x] 8.2 Refactor `apps/web/components/TicketCategoryCard.vue` — quantity selector touch-friendly, harga prominent
- [x] 8.3 Refactor `apps/web/components/QueueStatus.vue` — angka posisi & total legible (font ≥ 24px), progress bar full-width
- [x] 8.4 Refactor `apps/web/components/CountdownTimer.vue` — font size scale, format DD:HH:MM:SS readable di mobile
- [x] 8.5 Audit komponen lain di `apps/web/components/` dan fix sesuai standar

## 9. Verification & Polish

- [ ] 9.1 Manual test semua halaman di Chrome DevTools responsive mode: 320, 375, 414, 768, 1024, 1280, 1536 px *(manual)*
- [ ] 9.2 Test di real device (jika tersedia) iOS Safari & Android Chrome *(manual)*
- [ ] 9.3 Lighthouse mobile audit — target Performance ≥ 90, Accessibility ≥ 95 *(manual)*
- [ ] 9.4 Test keyboard navigation & screen reader di komponen baru (drawer, bottom sheet) *(manual)*
- [ ] 9.5 Test dark mode / high contrast (jika relevan) *(manual)*
- [ ] 9.6 Final review: tidak ada horizontal scroll di mobile, semua CTA reachable, semua form usable *(manual)*
- [ ] 9.7 Update README dengan section "Mobile Support" jika perlu *(manual)*
