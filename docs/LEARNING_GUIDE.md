# LEARNING GUIDE

Dokumen ini membantu kamu membaca project **bertahap per sesi**.

---

## Sesi 1 — Gambaran besar dan struktur project

### Tujuan sesi
Memahami bentuk project secara menyeluruh sebelum masuk ke code.

### File penting
- `README.md`
- `go.mod`
- `docker-compose.yml`
- `cmd/api/main.go`

### Yang perlu dipahami
- project ini memakai **package-based structure**
- `cmd/` berisi entry point
- `internal/auth` berisi semua code auth
- `internal/product` berisi semua code product
- `internal/platform` berisi hal teknis seperti DB, Redis, router, logger
- `internal/shared` berisi helper lintas fitur

### Ringkasan sesi
Kalau kamu sudah paham struktur folder, kamu akan lebih mudah membaca file lain tanpa bingung.

---

## Sesi 2 — Config, database, migration, seeder

### Tujuan sesi
Memahami bagaimana aplikasi membaca environment, terkoneksi ke MySQL, menjalankan migration SQL manual, lalu mengisi data awal.

### File penting
- `configs/config.go`
- `internal/platform/db/mysql.go`
- `internal/platform/db/migrator.go`
- `cmd/migrate/main.go`
- `cmd/seeder/main.go`
- `migrations/*.sql`

### Alur berpikir
1. `config.go` membaca env
2. `mysql.go` membuka koneksi database
3. `migrator.go` menjalankan SQL file manual
4. `cmd/migrate` hanya menjadi pintu masuk untuk runner migration
5. `cmd/seeder` mengisi data demo setelah schema siap

### Kenapa dipisah?
- config dipisah supaya `main.go` tidak penuh
- migration dipisah dari seeder supaya schema dan data demo tidak tercampur
- migration manual membuat kamu benar-benar paham bentuk tabel

### Ringkasan sesi
Setelah sesi ini, kamu sudah tahu bagaimana project “berdiri” dari nol.

---

## Sesi 3 — Auth JWT

### Tujuan sesi
Memahami login, pembuatan token, validasi token, dan penyimpanan user ke context.

### File penting
- `internal/auth/model.go`
- `internal/auth/dto.go`
- `internal/auth/repository.go`
- `internal/auth/service.go`
- `internal/auth/jwt.go`
- `internal/auth/context.go`
- `internal/auth/handler.go`
- `internal/platform/http/middleware/auth.go`

### Alur login
1. request login masuk ke handler
2. handler validasi body request
3. service mencari user lewat repository
4. password dibandingkan dengan bcrypt hash
5. JWT dibuat oleh `jwt.go`
6. token dikirim ke client

### Alur endpoint protected
1. client kirim `Authorization: Bearer <token>`
2. middleware auth membaca token
3. middleware parse token dan validasi claims
4. user disimpan ke context
5. handler tinggal membaca `auth.FromContext(...)`

### Kenapa auth logic tidak ditaruh di handler?
Agar handler tetap pendek dan fokus pada HTTP saja.

### Ringkasan sesi
Kalau sesi ini sudah paham, kamu sebenarnya sudah menguasai pondasi login API modern.

---

## Sesi 4 — CRUD products

### Tujuan sesi
Memahami alur CRUD yang rapi dari handler -> service -> repository.

### File penting
- `internal/product/model.go`
- `internal/product/dto.go`
- `internal/product/repository.go`
- `internal/product/service.go`
- `internal/product/handler.go`

### Alur create product
1. handler decode JSON
2. handler validasi request
3. service membentuk entity product
4. repository menyimpan ke MySQL dengan GORM
5. service invalidasi cache
6. handler kirim response

### Kenapa repository penting?
Karena query database sebaiknya tidak tersebar di handler.

### Ringkasan sesi
Kalau kamu bisa meniru pola product ini, kamu bisa membuat fitur lain seperti category, order, atau profile.

---

## Sesi 5 — Redis, middleware, logging, validation

### Tujuan sesi
Memahami concern yang dipakai lintas endpoint.

### File penting
- `internal/platform/cache/redis.go`
- `internal/platform/http/middleware/tenant.go`
- `internal/platform/http/middleware/auth.go`
- `internal/platform/http/middleware/logging.go`
- `internal/platform/http/router.go`
- `internal/platform/logger/logger.go`
- `internal/shared/tenant/context.go`

### Yang perlu dipahami
- Redis dipakai sebagai cache sederhana
- Tenant middleware mengambil `X-Tenant-Id`
- Auth middleware memvalidasi JWT
- Logging middleware mencatat request penting
- Router menyusun urutan middleware dan route

### Cache hit vs cache miss
- hit: data ditemukan di Redis
- miss: data diambil dari MySQL lalu disimpan ke Redis

### Ringkasan sesi
Setelah sesi ini, kamu paham kenapa project backend modern tidak hanya berisi handler dan query DB saja.

---

## Sesi 6 — Exception / error handling

### Tujuan sesi
Memahami cara membuat response error yang konsisten dan aman.

### File penting
- `internal/shared/apperror/error.go`
- `internal/shared/response/response.go`

### Alur error
1. repository / service bisa menghasilkan error
2. error dibungkus / diterjemahkan menjadi `AppError`
3. response helper mengubah `AppError` menjadi JSON response yang konsisten

### Kenapa penting?
Karena client butuh format error yang stabil, dan kita tidak boleh membocorkan detail internal server secara mentah.

### Ringkasan sesi
Kalau sesi ini dipahami, API kamu akan terasa jauh lebih rapi dan profesional.

---

## Sesi 7 — Swagger

### Tujuan sesi
Memahami dokumentasi API sederhana.

### File penting
- `docs/openapi.yaml`
- `docs/swagger.html`
- `internal/platform/http/router.go`

### Kenapa dibuat sederhana?
Supaya pemula bisa langsung melihat hasil dokumentasi tanpa harus bergantung ke generator yang lebih kompleks.

### Ringkasan sesi
Swagger membantu kamu dan orang lain mencoba endpoint lebih cepat.

---

## Sesi 8 — Docker Compose dan debug

### Tujuan sesi
Memahami cara menjalankan project secara konsisten dan cara debug lokal.

### File penting
- `docker-compose.yml`
- `Dockerfile`
- `.vscode/launch.json`

### Ringkasan sesi
Setelah sesi ini, kamu bisa menjalankan dan debug project dengan lebih nyaman.

---

## Sesi 9 — Integration test

### Tujuan sesi
Memahami pengujian end-to-end dengan service nyata.

### File penting
- `tests/integration/api_test.go`

### Yang dites
- login sukses
- login gagal
- create product dengan token
- create product tanpa token
- get product
- update product
- delete product
- tenant isolation

### Kenapa integration test penting?
Karena test ini tidak hanya memeriksa fungsi kecil, tapi memeriksa alur nyata dari HTTP sampai database.

### Ringkasan sesi
Ini sangat berguna saat project mulai bertambah besar.

---

## Sesi 10 — Cara menambah fitur baru

### Tujuan sesi
Membuat kamu mandiri setelah memahami starter ini.

### Pola dasar fitur baru
Misalnya ingin menambah `category`:
1. buat `internal/category/model.go`
2. buat `internal/category/dto.go`
3. buat `internal/category/repository.go`
4. buat `internal/category/service.go`
5. buat `internal/category/handler.go`
6. buat migration SQL manual
7. tambahkan route di router
8. tambahkan test

### Cara berpikir
Saat membuat fitur baru, pikirkan:
- data apa yang disimpan?
- siapa yang boleh akses?
- tenant filter perlu di mana?
- validasi request apa?
- perlu cache atau tidak?
- error apa yang mungkin muncul?

### Ringkasan sesi
Kalau pola ini kamu ikuti, fitur baru biasanya bisa dibuat dengan rapi tanpa harus mengubah keseluruhan project.
