## 1. Fix Chi Router

- [x] 1.1 Di `apps/api/internal/router/router.go`, pindahkan `r.Post("/events/{eventId}/categories", catH.Create)` dari blok `r.Route("/admin", ...)` (line ~113) ke dalam blok `r.Route("/admin/events", ...)` (sebelum `})` penutup)
- [x] 1.2 Tambah komentar dokumentasi di atas `r.Route("/admin/events", ...)` yang menjelaskan pola "event-scoped routes MUST be inside this subroute, not in outer /admin — chi does not fall through to sibling subroutes on miss"
- [x] 1.3 Verifikasi dengan `go build ./...` di `apps/api/` — harus compile tanpa error
- [x] 1.4 Rebuild Docker image: `docker compose build --no-cache api`
- [x] 1.5 Restart container: `docker compose up -d api`
- [x] 1.6 Smoke test: `curl -X POST http://localhost:8081/api/admin/events/1/categories -H "Authorization: Bearer $ADMIN_TOKEN" -d '{"name":"Test","price":100,"total_stock":10,"max_per_user":2}'` → harus return **201 Created** (sebelumnya 404)

## 2. Verify Other Endpoints (Re-audit)

- [x] 2.1 Test `GET /api/admin/events` → 200
- [x] 2.2 Test `GET /api/admin/bookings` → 200
- [x] 2.3 Test `GET /api/admin/bookings/1` → JSON `{"error":"not_found","message":"..."}` (PASS, bukan 404 page)
- [x] 2.4 Test `GET /api/admin/payments` → 200
- [x] 2.5 Test `GET /api/admin/users` → 200
- [x] 2.6 Test `POST /api/admin/events` → 201 (create new event)
- [x] 2.7 Test `POST /api/admin/events/{id}/publish` → 200 (publish event)
- [x] 2.8 Test `GET /api/events` (public) → 200
- [x] 2.9 Test `GET /api/events/1` (public) → 200
- [x] 2.10 Test `GET /api/events/1/categories` (public) → 200
- [x] 2.11 Test `POST /api/war/join?event_id=1` (auth) → 200
- [x] 2.12 Test `GET /api/war/status?event_id=1` (auth) → 200
- [x] 2.13 Test `POST /api/payments/create` (auth) → 200 atau 201 (atau JSON 404 jika booking_id tidak ada)
- [x] 2.14 Test `GET /api/payments/me` (auth) → 200
- [x] 2.15 Test `GET /api/admin/stats` (admin) → 200
- [x] 2.16 Test `GET /api/auth/me` (auth) → 200

## 3. End-to-End Flow Verification

- [x] 3.1 Login sebagai admin (`admin@example.com` / `admin123`) — dapat token
- [x] 3.2 Create event baru: `POST /api/admin/events` dengan `{title, venue, start_date, end_date}` — dapat event_id
- [x] 3.3 Add category: `POST /api/admin/events/{event_id}/categories` dengan `{name, price, total_stock, max_per_user}` — dapat category_id
- [x] 3.4 Publish event: `POST /api/admin/events/{event_id}/publish` — event status jadi `published`
- [x] 3.5 Verify public bisa lihat: `GET /api/events` — event baru muncul
- [x] 3.6 Register user baru
- [x] 3.7 User join war: `POST /api/war/join?event_id={event_id}` — dapat `{redirect_to_booking: true}` atau `{queued: true, position, token}`
- [x] 3.8 Lihat di frontend `http://localhost:3001/admin/events/new` dan `http://localhost:3001/admin/events/{id}` — flow tambah event & kategori bisa diselesaikan tanpa error

## 4. Create Smoke Test Script

- [x] 4.1 Buat `apps/api/scripts/smoke.sh` — bash script yang:
  - Define `API_BASE` (default `http://localhost:8080` dari .env)
  - Login admin → simpan token
  - Iterasi list of test cases (method, path, expected_status)
  - Untuk setiap test, hit endpoint dan cek response
  - Track PASS/FAIL count
  - Print summary di akhir
  - Exit 0 jika semua PASS, exit 1 jika ada FAIL
- [x] 4.2 Tambah minimal 12 test cases (termasuk yang diminta)
- [x] 4.3 Implementasi pembeda 404-page vs 404-handler:
  - Body `404 page not found` → FAIL
  - JSON body → PASS
- [x] 4.4 Tambah `make smoke` di `apps/api/Makefile` yang run script
- [ ] 4.5 Verifikasi: `bash apps/api/scripts/smoke.sh` → exit 0 dengan semua PASS (requires running API)
- [ ] 4.6 Verifikasi: matikan API container, jalankan script → exit non-zero dengan error jelas

## 5. Documentation

- [x] 5.1 Update `apps/api/internal/router/router.go` dengan komentar pola (sudah di task 1.2)
- [x] 5.2 Tambah `docs/router.md` — jelaskan chi no-fallthrough rule dengan contoh
- [x] 5.3 Update root `README.md`: tambah "Smoke Test" section dan referensi `docs/router.md`
- [x] 5.4 Tambah section "Troubleshooting" di root `README.md`

## 6. Commit & Cleanup

- [x] 6.1 Commit router fix: `fix(router): move catH.Create into /admin/events subroute to fix 404`
- [x] 6.2 Commit smoke test: `test(api): add smoke test script for endpoint regression detection`
- [x] 6.3 (Opsional) git tag: `v0.2.1-router-fix`
- [x] 6.4 Verifikasi tidak ada file temporary atau debug yang tertinggal
