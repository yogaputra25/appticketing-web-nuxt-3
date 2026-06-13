## Context

Project ini adalah SaaS pemesanan tiket konser (War Tiket) dengan stack Nuxt 3 + Go (Chi) + PostgreSQL + Redis. Backend Go di `apps/api/` sudah selesai dan `router.go` mendeklarasikan banyak endpoint. Saat menjalankan flow admin dari frontend (login sebagai admin → klik "Tambah Event" → isi form → submit), request `POST /api/admin/events/{id}/categories` selalu kembali **404 page not found** (plain text, 19 bytes — chi default handler), padahal endpoint tersebut sudah dideklarasikan di `router.go` line 113 di dalam `r.Route("/admin", ...)`.

**Akar masalah:** Chi router menggunakan **radix trie** untuk route matching. Saat request `POST /api/admin/events/1/categories` masuk, chi memproses:
1. Cocok dengan prefix `/api/admin/events` (line 100) — MATCH
2. Masuk ke subroute `/admin/events`, cari handler untuk pattern `/{id}/categories` di antara handler yang dideklarasikan di dalam subroute tersebut
3. Pola yang dideklarasikan di dalam subroute `/admin/events` hanya: `/`, `/{id}`, `/{id}/publish` — `/{id}/categories` TIDAK ADA
4. Karena chi tidak melakukan **fallthrough ke sibling subroute** (`/admin`), ia langsung return 404 dengan body `404 page not found`

**Audit lengkap** (sudah dilakukan) menunjukkan:
- ✅ Endpoint lain (`POST /api/admin/events/{id}/publish`, `GET /api/admin/events/{id}`, `GET /api/admin/stats`, dll.) berfungsi dengan benar
- ✅ Beberapa 404 yang sebelumnya dikira "endpoint missing" sebenarnya adalah JSON `{"error":"not_found","message":"..."}` (52 bytes) — itu dari handler yang benar yang merespons "resource not found", bukan dari chi default
- ❌ Hanya `POST /api/admin/events/{eventId}/categories` yang benar-benar tidak ter-registrasi di trie chi

**Stack & constraints:**
- Go 1.22+, Chi v5
- File `apps/api/internal/router/router.go` adalah single source of truth untuk endpoint
- Build via Docker multi-stage (`Dockerfile`) dengan `go build` di stage `builder`
- Frontend tidak perlu diubah — saat backend fixed, flow otomatis bekerja

## Goals / Non-Goals

**Goals:**
- Memperbaiki router agar `POST /api/admin/events/{eventId}/categories` ter-registrasi dengan benar
- Mendokumentasikan pola router (event-scoped routes di dalam subroute) untuk mencegah regresi
- Membuat smoke test script untuk regresi detection — jalankan setelah perubahan router atau handler
- Verifikasi end-to-end flow admin (tambah event → tambah kategori → publish) bisa diselesaikan

**Non-Goals:**
- Refactor besar struktur router (mis. migrasi ke OpenAPI/Swagger generator)
- Migrasi ke chi v6 atau library router lain
- Menambah endpoint baru di luar yang sudah ada di `router.go`
- Mengubah API contract (path, method, request/response shape)
- Menambah CI/CD pipeline (smoke test cukup dijalankan manual dulu)
- Healthcheck improvement (low priority, di luar scope minimum)

## Decisions

### 1. Pindahkan `catH.Create` ke dalam subroute `/admin/events` (bukan deklarasi baru)

**Pilihan:** Edit `apps/api/internal/router/router.go` line 113 — pindahkan `r.Post("/events/{eventId}/categories", catH.Create)` dari blok `r.Route("/admin", ...)` ke dalam blok `r.Route("/admin/events", ...)` (sebelum tanda `})` penutup di line 108).

**Alasan:**
- **Fix minimal & deterministik** — 1 baris dipindah, tidak ada perubahan logika
- **Mengikuti pola router** — `/admin/events/{id}/publish` (POST) sudah ada di dalam subroute `/admin/events`; `/admin/events/{id}/categories` (POST) juga harus di sana
- **Zero risk** — tidak mengubah path, method, atau handler

**Alternatif:**
- Deklarasikan ulang di kedua tempat (DRY violation, susah di-maintain)
- Gunakan chi `With()` untuk inject handler di multiple subroutes — overkill
- Refactor pakai group terpisah untuk event-scoped operations — perubahan besar yang tidak perlu

### 2. Tulis komentar dokumentasi pola router di `router.go`

**Pilihan:** Tambah blok komentar di atas `r.Route("/admin/events", ...)` yang menjelaskan: "Event-scoped admin routes (paths starting with `/admin/events/{id}/...`) MUST be declared inside this subroute, not in the outer `/admin` subroute. Chi does not fall through to sibling subroutes on miss."

**Alasan:**
- Menjelaskan "why" — bukan hanya "what"
- Mencegah developer di masa depan untuk menambah route serupa di tempat yang salah
- Self-documenting code — lebih baik dari wiki terpisah yang bisa stale

**Alternatif:**
- Doc terpisah di `docs/router.md` — bisa outdated
- Tidak ada komentar — bug akan terulang

### 3. Smoke test script: `bash` + `curl` + `jq`

**Pilihan:** Buat `apps/api/scripts/smoke.sh` yang:
1. Login sebagai admin untuk dapat token
2. Iterasi list of (method, path, expected_status, body_predicate)
3. Untuk setiap endpoint, hit dan cek status code + content-type/body shape
4. Print PASS/FAIL summary
5. Exit non-zero jika ada FAIL

**Alasan:**
- **Standalone** — tidak butuh dependency runtime (cukup `bash`, `curl`, `jq`)
- **Idempotent** — bisa dijalankan berulang
- **Self-contained** — admin credentials di-hardcode (untuk dev/test only)
- **Bisa di-extend** — tambah endpoint baru cukup tambah 1 baris

**Alternatif:**
- Go test (`_test.go` dengan `httptest.NewServer`) — lebih robust tapi lebih banyak boilerplate, dan butuh Go di environment tester
- Postman/Newman — perlu export collection, susah di-version-control
- `k6` atau `wrk` — untuk load testing, bukan smoke test
- Python script — tambah dependency

### 4. Pembeda 404-page vs 404-handler di smoke test

**Pilihan:** Script membedakan dua jenis 404:
- **404 page not found** (plain text, 19 bytes) → **FAIL** (route missing di chi trie)
- **`{"error":"not_found","message":"..."}`** (JSON, 50+ bytes) → **PASS** (route exists, resource missing — itu behavior yang diharapkan)

**Alasan:**
- **Akurat** — sebelumnya salah diagnosa karena tidak membedakan dua jenis 404
- **Mencegah false positive** — endpoint yang return 404 JSON untuk resource yang belum ada tidak akan dianggap broken
- **Mencegah false negative** — endpoint yang sebenarnya missing (chi 404) akan terdeteksi

**Alternatif:**
- Hanya cek status code (kurang presisi)
- Hanya cek body (kurang reliable)

### 5. Tidak menambah healthcheck detail (db/redis ping)

**Pilihan:** Skip improvement `GET /healthz` untuk return status DB & Redis. Tetap return `{"status":"ok"}` saja.

**Alasan:**
- **Out of scope minimum** — bug yang harus diperbaiki adalah chi routing, bukan healthcheck
- **Bisa ditambah nanti** — change kecil, isolated, tidak mendesak
- **Jangan bengkakkan change** — keep-it-focused principle

**Alternatif:**
- Tambah sekarang — scope creep
- Tambah di change terpisah — overkill untuk 1 endpoint

## Risks / Trade-offs

- **Risiko:** Developer di masa depan menambah event-scoped route di outer `/admin` block, kembali 404
  → **Mitigasi:** Komentar dokumentasi di `router.go` + smoke test script yang akan fail jika regression terjadi

- **Risiko:** Smoke test script bisa outdated (route ditambah tapi smoke test tidak diupdate)
  → **Mitigasi:** Smoke test mengecek bahwa SEMUA path yang dideklarasikan di `router.go` ada di test list (atau generate otomatis via `grep` pattern dari router.go)

- **Risiko:** Smoke test butuh API running & database seeded
  → **Mitigasi:** Script assume `docker compose up -d` sudah jalan (default `http://localhost:8081` di .env); tampilkan pesan error yang jelas jika tidak bisa konek

- **Risiko:** Smoke test akan menciptakan data dummy (event, category) yang bisa menumpuk di DB dev
  → **Mitigasi:** Gunakan random suffix (timestamp) untuk nama event agar tidak bentrok; di production, script ini tidak dijalankan

- **Risiko:** Admin credentials di-hardcode di smoke test (`admin@example.com` / `admin123`)
  → **Mitigasi:** Bisa di-override via env var `ADMIN_EMAIL` dan `ADMIN_PASSWORD`; default hanya untuk dev

- **Trade-off:** Smoke test pakai bash + curl — kurang robust dibanding Go test untuk edge cases
  → **Mitigasi:** Cukup untuk verifikasi "route exists & respond as expected"; untuk business logic pakai Go unit/integration test (sudah ada di `repository/*_test.go`)

## Migration Plan

Tidak ada data migration (zero DB change). Rollout:

1. **Edit `apps/api/internal/router/router.go`** — pindahkan 1 baris + tambah komentar
2. **Rebuild API image** — `docker compose build --no-cache api` (force fresh)
3. **Restart container** — `docker compose up -d api`
4. **Verifikasi manual** — `curl POST /api/admin/events/{id}/categories` → 201 Created
5. **Buat `apps/api/scripts/smoke.sh`** — script bash untuk regression test
6. **Tambah `make smoke` di `Makefile`** — convenience target
7. **Jalankan smoke test** — `bash apps/api/scripts/smoke.sh` harus exit 0
8. **Commit** — 1 commit untuk router fix + 1 commit untuk smoke test (atau gabung)

**Rollback strategy:** Karena perubahan minimal (1 baris dipindah + 1 komentar), rollback cukup revert commit dan rebuild. Tidak ada data yang berubah.

## Open Questions

- Apakah perlu menambahkan audit log middleware (chi middleware) untuk mencatat setiap request ke console? — Bisa berguna untuk debugging, tapi di luar scope.
- Apakah smoke test harus dijalankan otomatis di CI? — Saat ini tidak ada CI setup; bisa ditambah nanti.
- Apakah perlu tambah endpoint `GET /api/healthz/db` dan `/api/healthz/redis` terpisah? — Bisa, tapi minor.
- Bagaimana cara handle chi route conflict jika ada route baru yang serupa? — Butuh ADR (Architecture Decision Record) atau convention doc untuk tim.
