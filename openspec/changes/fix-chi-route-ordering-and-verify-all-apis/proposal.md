## Why

Saat admin mencoba menambah event baru atau menambahkan kategori tiket ke event yang ada dari frontend, request `POST /api/admin/events/{id}/categories` selalu kembali **404 Not Found** padahal endpoint tersebut sudah dideklarasikan di `apps/api/internal/router/router.go` (line 113). Investigasi menunjukkan akar masalahnya adalah **chi router tidak melakukan fallthrough** antara sibling sub-routes: ketika request cocok dengan prefix `/admin/events`, chi hanya mencari handler di dalam blok `r.Route("/admin/events", ...)`. Jika tidak ada yang match, ia mengembalikan 404 langsung — tidak pernah mencari di blok `r.Route("/admin", ...)` yang juga berisi `r.Post("/events/{eventId}/categories", catH.Create)`. Hal yang sama berpotensi terjadi untuk route lain di masa depan jika developer menambahkan handler event-scoped di outer `/admin` block.

Akibatnya, flow admin end-to-end (tambah event → tambah kategori → publish) tidak bisa diselesaikan — fitur inti produk tidak usable.

## What Changes

- **Fix chi router ordering** — Pindahkan `r.Post("/{eventId}/categories", catH.Create)` dari blok `r.Route("/admin", ...)` ke dalam blok `r.Route("/admin/events", ...)` agar path `/admin/events/{eventId}/categories` dapat di-resolve di trie yang sama. Tambahkan komentar inline yang menjelaskan pola ini untuk mencegah regresi.
- **Audit & verifikasi semua endpoint** — Jalankan smoke test (script curl) untuk semua endpoint yang dideklarasikan di `router.go` dan cocokkan dengan response yang diharapkan. Hasil audit menunjukkan semua endpoint lain SUDAH aktif dan berfungsi dengan benar — error 404 pada beberapa path sebelumnya (mis. `GET /api/admin/bookings/1`, `POST /api/payments/create` dengan `booking_id=1`) sebenarnya adalah response JSON `{"error":"not_found","message":"..."}` dari handler yang benar (52 bytes), bukan 404 page-not-found dari chi (19 bytes, plain text). Verifikasi ini mencegah re-open issue yang sebenarnya sudah selesai.
- **Dokumentasi pola router** — Tambahkan catatan singkat di `router.go` (komentar di atas blok admin routes) tentang pola "declare event-scoped routes inside the `/admin/events` subroute" untuk developer selanjutnya.
- **Opsional — tambah healthcheck route yang lebih informatif** — `GET /healthz` saat ini hanya return `{"status":"ok"}`. Bisa ditambah `{"status":"ok","db":"ok","redis":"ok"}` setelah ping dependency (low priority, di luar scope minimum).

## Capabilities

### New Capabilities

- `api-route-audit`: Kemampuan untuk melakukan audit end-to-end terhadap semua route yang dideklarasikan di `router.go` dan memastikan masing-masing merespons dengan benar. Disertai dengan script smoke test yang bisa dijalankan manual (atau di CI) untuk regresi.

### Modified Capabilities

<!-- Tidak ada requirement backend yang berubah — endpoint yang "baru" sebenarnya
     sudah ada di source code, hanya routing-nya yang perlu diperbaiki. -->

## Impact

- **Backend** (`apps/api/`): file `apps/api/internal/router/router.go` — pindahkan 1 route, tambah komentar. Tidak ada perubahan handler, tidak ada migration DB, tidak ada perubahan API contract (path, method, response shape tetap sama).
- **Frontend**: tidak ada perubahan kode. Flow admin yang sebelumnya 404 sekarang akan sukses otomatis.
- **Testing**: tambah script smoke test untuk 12-15 endpoint kritis (admin + public + authenticated) di folder `scripts/` atau `apps/api/scripts/smoke.sh` agar bisa dijalankan via `make smoke` atau `bash scripts/smoke.sh`.
- **Docs**: README opsional — tambahkan 1 paragraf tentang pola router chi yang dipakai (event-scoped routes di dalam sub-route, bukan di outer `/admin`).
- **No breaking change**: path & method endpoint tidak berubah.
