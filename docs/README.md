# go-starter

> Catatan: dokumen ini menjelaskan struktur versi awal. Untuk struktur versi siap berkembang, lihat `docs/LEARNING_GUIDE_V2.md`.

Starter backend Go yang ramah untuk pemula, tetap sederhana, tetapi mengikuti praktik industri yang masuk akal untuk project kecil-menengah.

Module project ini adalah:

```go
github.com/niamfauzi/go-starter
```

> Catatan: `go.sum` belum saya sertakan di template ini karena file tersebut biasanya dihasilkan otomatis saat kamu menjalankan `go mod tidy` atau build pertama.

---

## Gambaran besar project

Project ini berisi contoh:

- CRUD sederhana untuk `products`
- JWT authentication
- middleware
- validation
- structured logging
- Redis cache
- Swagger / OpenAPI sederhana
- Docker Compose
- migration SQL manual
- seeder
- integration test

Tujuan utamanya adalah **belajar dengan struktur yang rapi tapi tidak overkill**.

---

## Daftar sesi belajar

### Sesi 1 — Gambaran besar dan struktur project
Fokus memahami kenapa folder dipisah.

### Sesi 2 — Config, database, migration, seeder
Fokus memahami startup aplikasi dan data awal.

### Sesi 3 — Auth JWT
Fokus login, token, dan route yang dilindungi.

### Sesi 4 — CRUD products
Fokus handler, service, repository.

### Sesi 5 — Redis, middleware, logging, validation
Fokus concern lintas fitur.

### Sesi 6 — Exception / error handling
Fokus error yang konsisten dan aman.

### Sesi 7 — Swagger
Fokus dokumentasi endpoint.

### Sesi 8 — Docker Compose dan debug
Fokus menjalankan project dengan nyaman.

### Sesi 9 — Integration test
Fokus test end-to-end dengan service nyata.

### Sesi 10 — Cara menambah fitur baru
Fokus pola berpikir agar kamu bisa menambah fitur lain sendiri.

---

## Struktur folder final

```text
.
├── .env.example
├── .vscode/
│   └── launch.json
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── migrate/
│   │   └── main.go
│   └── seeder/
│       └── main.go
├── configs/
│   └── config.go
├── docs/
│   ├── swagger.yaml
│   └── swagger.html
├── internal/
│   ├── auth/
│   │   ├── context.go
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── jwt.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── platform/
│   │   ├── cache/
│   │   │   └── redis.go
│   │   ├── db/
│   │   │   ├── migrator.go
│   │   │   └── mysql.go
│   │   ├── http/
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── logging.go
│   │   │   │   └── tenant.go
│   │   │   └── router.go
│   │   └── logger/
│   │       └── logger.go
│   ├── product/
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   └── shared/
│       ├── apperror/
│       │   └── error.go
│       ├── response/
│       │   └── response.go
│       └── tenant/
│           └── context.go
├── migrations/
│   ├── 000001_create_users.down.sql
│   ├── 000001_create_users.up.sql
│   ├── 000002_create_products.down.sql
│   └── 000002_create_products.up.sql
├── scripts/
│   └── demo-curl.sh
├── tests/
│   └── integration/
│       └── api_test.go
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

---

## Kenapa struktur ini dipilih?

Karena kita memakai **package-based structure**:

- `auth/` berisi semua hal tentang auth
- `product/` berisi semua hal tentang product
- `platform/` berisi hal teknis lintas fitur seperti DB, Redis, router, logger
- `shared/` berisi helper yang dipakai banyak fitur

Ini lebih mudah dipahami pemula daripada langsung memakai clean architecture yang sangat dalam.

---

## Kenapa handler / service / repository dipisah?

### Handler
Tugasnya menerima HTTP request dan mengirim HTTP response.

### Service
Tugasnya business logic. Contoh:
- login valid atau tidak
- create product lalu invalidasi cache
- cek tenant

### Repository
Tugasnya query ke database.

Dengan pemisahan ini, code:
- lebih rapi
- lebih mudah dites
- lebih mudah diubah tanpa merusak bagian lain

---

## Kenapa tenant diambil dari header `X-Tenant-Id`?

Karena untuk project multi-tenant, kita butuh cara yang konsisten agar setiap request tahu dia milik tenant siapa.

Kenapa header?
- mudah dibaca middleware
- tidak perlu diulang di body
- cocok untuk dipakai lintas endpoint

---

## Kenapa JWT logic tidak ditaruh semua di handler?

Kalau semua logic JWT ditaruh di handler:
- handler akan cepat menjadi terlalu panjang
- parsing token akan terulang di banyak tempat
- lebih susah di-maintain

Makanya kita pisahkan:
- `jwt.go` untuk generate/parse token
- middleware auth untuk validasi token
- handler cukup fokus ke request/response

---

## Kenapa migration SQL manual lebih baik untuk belajar?

Karena saat belajar backend, kamu perlu memahami schema database secara jelas.

Kalau pakai AutoMigrate:
- memang cepat
- tapi kamu kurang belajar bagaimana schema benar-benar dibentuk

Kalau pakai SQL manual:
- lebih jelas
- lebih terkontrol
- lebih mirip workflow banyak tim saat menjaga schema production

---

## Kenapa Redis dipakai sederhana?

Di project ini Redis hanya dipakai untuk:
- cache list products
- cache detail product

Tujuannya supaya kamu paham konsep:
- cache hit
- cache miss
- cache invalidation

Belum perlu dibuat rumit seperti distributed lock, pub/sub, queue, dan lain-lain.

---

## Setup awal

### 1. Clone / buka project
Masuk ke folder project ini.

### 2. Salin env file

```bash
cp .env.example .env
```

### 3. Install dependency Go

```bash
go mod tidy
```

Perintah ini akan mengunduh dependency dan otomatis membuat `go.sum`.

---

## Cara menjalankan dengan Docker Compose

### 1. Jalankan container

```bash
docker compose up -d
```

Ini akan menjalankan:
- MySQL
- Redis
- API

### 2. Jalankan migration

```bash
docker compose exec api go run ./cmd/migrate up
```

### 3. Jalankan seeder

```bash
docker compose exec api go run ./cmd/seeder
```

### 4. Cek health endpoint

```bash
curl -H "X-Tenant-Id: tenant-demo" http://localhost:8080/healthz
```

---

## Cara menjalankan tanpa Docker untuk API

Kalau MySQL dan Redis sudah berjalan secara lokal:

```bash
go run ./cmd/api
```

---

## Alur request dari HTTP sampai DB / cache

Contoh request create product:

1. Request masuk ke router
2. Middleware request ID memberi ID request
3. Middleware logger mencatat request
4. Middleware tenant mengambil `X-Tenant-Id`
5. Middleware auth memvalidasi JWT
6. Handler membaca body JSON
7. Handler validasi request
8. Handler memanggil service
9. Service memanggil repository
10. Repository menyimpan ke MySQL dengan GORM
11. Service menghapus cache Redis yang sudah usang
12. Response dikirim ke client

---

## Auth JWT

### Login demo

- tenant_id: `tenant-demo`
- email: `demo@example.com`
- password: `password123`

### Request login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login   -H "Content-Type: application/json"   -H "X-Tenant-Id: tenant-demo"   -d '{
    "email": "demo@example.com",
    "password": "password123"
  }'
```

### Ambil profile

```bash
curl http://localhost:8080/api/v1/auth/me   -H "Authorization: Bearer TOKEN_KAMU"   -H "X-Tenant-Id: tenant-demo"
```

### Kenapa semua endpoint product private?

Untuk starter multi-tenant, lebih aman bila dari awal semua endpoint product butuh token.
Ini membuat:
- alur tenant lebih konsisten
- peluang bocor data lebih kecil
- pola belajar lebih seragam

---

## CRUD products

### Create product

```bash
curl -X POST http://localhost:8080/api/v1/products   -H "Content-Type: application/json"   -H "X-Tenant-Id: tenant-demo"   -H "Authorization: Bearer TOKEN_KAMU"   -d '{
    "name": "Headset Belajar",
    "description": "Headset untuk meeting",
    "price": 250000,
    "stock": 10
  }'
```

### List products

```bash
curl http://localhost:8080/api/v1/products   -H "X-Tenant-Id: tenant-demo"   -H "Authorization: Bearer TOKEN_KAMU"
```

### Get detail product

```bash
curl http://localhost:8080/api/v1/products/1   -H "X-Tenant-Id: tenant-demo"   -H "Authorization: Bearer TOKEN_KAMU"
```

### Update product

```bash
curl -X PUT http://localhost:8080/api/v1/products/1   -H "Content-Type: application/json"   -H "X-Tenant-Id: tenant-demo"   -H "Authorization: Bearer TOKEN_KAMU"   -d '{
    "name": "Headset Belajar Updated",
    "description": "Deskripsi baru",
    "price": 275000,
    "stock": 15
  }'
```

### Delete product

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1   -H "X-Tenant-Id: tenant-demo"   -H "Authorization: Bearer TOKEN_KAMU"
```

---

## Validation

Validasi diletakkan dekat request DTO karena:
- mudah dilihat saat membaca request
- field yang wajib / format yang benar jadi jelas
- handler lebih mudah dipahami

Contoh:
- email wajib valid
- password minimal 8 karakter
- name product minimal 3 karakter
- price dan stock tidak boleh negatif

---

## Exception / error handling

Project ini memakai custom `AppError`.

Tujuannya:
- response error konsisten
- handler tidak perlu merakit error manual terus-menerus
- internal error tidak bocor ke client

Format error:

```json
{
  "success": false,
  "message": "validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "Email",
        "message": "format email tidak valid"
      }
    ]
  }
}
```

Jenis error yang dibedakan:
- validation error
- bad request
- unauthorized
- forbidden
- not found
- conflict
- internal server error

---

## Redis: cache hit dan cache miss

### Cache miss
Saat data belum ada di Redis:
1. service ambil data dari MySQL
2. data disimpan ke Redis
3. data dikirim ke client

### Cache hit
Saat data sudah ada di Redis:
1. service ambil data dari Redis
2. tidak perlu query MySQL
3. response lebih cepat

### Kapan cache dihapus?
Saat:
- create product
- update product
- delete product

Kenapa?
Karena data lama di cache sudah tidak valid.

---

## Logging

Project ini memakai zap.

Yang dicatat:
- request method
- path
- status code
- durasi
- request id
- error penting

Kenapa logging penting?
Karena saat terjadi bug, kita butuh jejak untuk mencari:
- request mana yang gagal
- berapa lama proses berjalan
- endpoint mana yang paling sering error

---

## Middleware dan urutannya

Urutan di router:

1. `RequestID`
2. `RealIP`
3. `Logger`
4. `Recoverer`
5. `Timeout`
6. `Tenant`
7. `AuthJWT` (hanya untuk route private)

Kenapa urutannya begitu?

- request ID sebaiknya dibuat di awal
- logger butuh request ID dan status akhir
- recoverer menjaga aplikasi tidak crash saat panic
- timeout mencegah request menggantung terlalu lama
- tenant diletakkan sebelum handler karena hampir semua logic membutuhkannya
- auth hanya dipasang di route yang memang harus login

---

## Swagger / OpenAPI

Project ini menyediakan:

- `docs/swagger.yaml`
- `docs/swagger.html`

Cara melihat:
1. Jalankan aplikasi
2. Buka:
   - `http://localhost:8080/docs/swagger.yaml`
   - `http://localhost:8080/docs/swagger.html`

> Catatan: `swagger.html` memakai asset Swagger UI dari CDN supaya setup tetap sederhana untuk pemula.

---

## Debug untuk pemula

File `.vscode/launch.json` sudah disiapkan untuk VS Code.

Langkah:
1. install extension Go di VS Code
2. buka folder project
3. tekan `F5`
4. pilih `Debug API`

Di balik layar, extension Go biasanya memakai Delve untuk debugging aplikasi Go.

---

## Integration test

Project ini menyediakan test integrasi di:

```text
tests/integration/api_test.go
```

Yang dites:
- login sukses
- login gagal
- create product dengan JWT valid
- create product tanpa JWT
- get product
- update product
- delete product
- tenant isolation

### Kenapa integration test penting?

Karena unit test hanya menguji potongan kecil,
sedangkan integration test menguji alur nyata:
- HTTP
- middleware
- JWT
- MySQL
- Redis

### Cara menjalankan

Pastikan Docker tersedia, lalu:

```bash
go test ./tests/integration -v
```

Test ini memakai Testcontainers untuk menyalakan container MySQL dan Redis secara otomatis.

---

## Seeder

Seeder dipisah dari migration karena:

- migration fokus ke schema
- seeder fokus ke data awal / data demo

Jalankan:

```bash
go run ./cmd/seeder
```

atau di Docker:

```bash
docker compose exec api go run ./cmd/seeder
```

---

## Cara menambah fitur baru

Misalnya kamu ingin menambah fitur `categories`.

Pola yang bisa kamu ikuti:

1. buat package `internal/category`
2. buat `model.go`
3. buat `dto.go`
4. buat `repository.go`
5. buat `service.go`
6. buat `handler.go`
7. tambahkan migration SQL manual
8. tambahkan route di router
9. jika perlu cache, tambahkan logic Redis di service
10. jika perlu auth, taruh route di group private

### Cara berpikirnya
Tanya ke diri sendiri:
- data apa yang disimpan?
- request apa yang diterima?
- business rule-nya apa?
- query database apa yang dibutuhkan?
- apakah perlu cache?
- apakah endpoint perlu login?

Kalau kamu bisa menjawab pertanyaan itu, biasanya struktur file akan mengikuti dengan alami.

---

## Best practice untuk pemula

- mulai dari struktur sederhana
- pisahkan handler / service / repository
- gunakan migration SQL manual agar schema jelas
- gunakan middleware untuk concern berulang
- simpan secret di environment
- hash password dengan bcrypt
- jangan campur query database langsung di handler
- gunakan tenant filter di semua query multi-tenant
- log error penting
- buat response konsisten

---

## Common mistakes untuk pemula

- menaruh semua code di `main.go`
- menaruh semua logic di handler
- lupa filter `tenant_id`
- menyimpan password plain text
- lupa invalidasi cache setelah update
- membuat Redis terlalu rumit sejak awal
- memakai AutoMigrate tanpa memahami schema
- mencampur error internal mentah ke response client
- membuat struktur terlalu “enterprise” untuk project yang masih kecil

---

## Ringkasan akhir

Starter ini sengaja dibuat:
- cukup rapi
- cukup realistis
- tidak terlalu berat
- nyaman dipelajari langkah demi langkah

Setelah kamu paham project ini, fitur berikutnya biasanya tinggal mengikuti pola yang sama.

Selamat belajar dan bereksperimen.
