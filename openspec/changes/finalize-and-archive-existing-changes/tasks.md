## 1. Verifikasi — integrate-backend-api-to-frontend

- [x] 1.1 Docker containers sudah running, API health check: `GET /healthz` → 200 OK, `GET /api/events` → 200 OK
- [x] 1.2 Frontend dev server sudah running di `http://localhost:3000` — 200 OK
- [ ] 1.3 Verifikasi end-to-end user flow: register → login → browse event → war → queue → booking → pay (simulate success) → lihat e-ticket
- [ ] 1.4 Verifikasi admin flow: login admin → lihat dashboard stats → manage events (create + add categories + publish) → lihat bookings → lihat payments
- [ ] 1.5 Verifikasi my bookings: user bisa lihat list booking, lihat detail, batalkan booking pending
- [ ] 1.6 Cek console browser untuk error; cek network tab untuk status 200/201/4xx yang expected

## 2. Verifikasi — fix-chi-route-ordering-and-verify-all-apis

- [x] 2.1 Smoke test via PowerShell: 15 endpoints tested — PASS (termasuk war, booking, payment, admin CRUD)
- [x] 2.2 API container di-stop, curl return error code 000 — smoke script akan mendeteksi connection error

## 3. Archive

- [x] 3.1 Tidak ada perubahan kode yang belum di-commit (clean working tree, hanya file change artifact baru)
- [ ] 3.2 Archive change `integrate-backend-api-to-frontend`: `openspec archive change "integrate-backend-api-to-frontend"`
- [ ] 3.3 Archive change `fix-chi-route-ordering-and-verify-all-apis`: `openspec archive change "fix-chi-route-ordering-and-verify-all-apis"`
