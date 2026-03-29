# Learning Flow System

Dokumen ini fokus menjelaskan **alur kerja request dari awal sampai akhir** di project ini, dengan contoh utama:

```http
POST /api/v1/products
```

Tujuan dokumen ini adalah membantu kamu membaca project nyata dengan lebih tenang.

Jadi kita tidak hanya membahas teori, tetapi benar-benar mengikuti file yang ada di repo ini.

---

## 1. Gambaran besar alur request

Saat client memanggil `POST /api/v1/products`, alurnya secara sederhana adalah:

1. request masuk ke HTTP server
2. server meneruskan request ke router
3. router mencocokkan path dan method
4. middleware umum dijalankan
5. middleware tenant membaca `X-Tenant-Id`
6. middleware auth membaca `Authorization: Bearer ...`
7. handler product membaca body JSON
8. handler memvalidasi body request
9. handler memanggil service product
10. service membuat entity product
11. repository menyimpan entity ke MySQL lewat GORM
12. service menghapus cache Redis yang sudah stale
13. handler membentuk response sukses
14. middleware logging mencatat hasil request
15. response dikirim kembali ke client

Kalau ada error di tengah jalan, alurnya berhenti di titik error itu, lalu `response.Error(...)` akan membentuk response error JSON yang konsisten.

---

## 2. Daftar file yang terlibat

Berikut file yang paling penting untuk memahami flow `Create Product`.

### Entry point dan wiring
- `cmd/api/main.go`
  Peran: menyusun dependency aplikasi saat startup.
- `internal/platform/http/server.go`
  Peran: membuat `http.Server`.
- `internal/platform/http/router.go`
  Peran: menyusun route dan urutan middleware.

### Middleware
- `internal/platform/http/middleware/requestid.go`
  Peran: memberi request ID.
- `internal/platform/http/middleware/cors.go`
  Peran: menangani CORS sederhana.
- `internal/platform/http/middleware/logging.go`
  Peran: mencatat method, path, status, durasi, dan request ID.
- `internal/platform/http/middleware/recovery.go`
  Peran: menjaga panic tidak menjatuhkan server.
- `internal/platform/http/middleware/timeout.go`
  Peran: memberi batas waktu request.
- `internal/platform/http/middleware/tenant.go`
  Peran: membaca `X-Tenant-Id` lalu menyimpannya ke context.
- `internal/platform/http/middleware/auth.go`
  Peran: membaca dan memvalidasi JWT, lalu menyimpan user ke context.

### Product module
- `internal/modules/product/delivery_http.go`
  Peran: parse request, validasi, panggil service, kirim response.
- `internal/modules/product/service.go`
  Peran: business logic create product.
- `internal/modules/product/repository.go`
  Peran: simpan data ke database lewat GORM.
- `internal/modules/product/entity.go`
  Peran: entity utama tabel `products`.
- `internal/modules/product/dto.go`
  Peran: bentuk request dan response product.
- `internal/modules/product/mapper.go`
  Peran: mengubah entity menjadi response.
- `internal/modules/product/cache.go`
  Peran: operasi cache Redis untuk product.

### Auth dan tenant
- `internal/modules/auth/jwt.go`
  Peran: parse JWT.
- `internal/modules/auth/context.go`
  Peran: simpan user login ke context.
- `internal/shared/tenant/context.go`
  Peran: simpan tenant ke context.

### Infrastruktur data
- `internal/platform/db/gorm.go`
  Peran: koneksi GORM ke MySQL.
- `internal/platform/cache/redis.go`
  Peran: koneksi Redis.

### Shared helper
- `internal/shared/apperror/error.go`
  Peran: definisi `AppError`.
- `internal/shared/apperror/mapper.go`
  Peran: mengubah error biasa menjadi `AppError`.
- `internal/shared/response/json.go`
  Peran: membentuk response JSON sukses dan error.
- `internal/shared/validator/validator.go`
  Peran: membuat validator instance.

---

## 3. Urutan alur request Create Product

Sekarang kita ikuti alurnya satu per satu.

### Langkah 1: aplikasi startup lebih dulu

Sebelum request apa pun datang, [cmd/api/main.go](/Users/endahfathonah/dev/projects/Laundry/lab-go-transaksi/cmd/api/main.go) sudah:

- load config
- membuat logger zap
- membuka koneksi MySQL lewat GORM
- membuka koneksi Redis
- membuat validator
- membuat repository, service, dan delivery untuk auth dan product
- membuat router
- membuat HTTP server

Artinya saat request datang, semua dependency utama sudah siap.

### Langkah 2: request masuk ke HTTP server

Client mengirim request:

```http
POST /api/v1/products HTTP/1.1
Host: localhost:8080
Authorization: Bearer <token>
X-Tenant-Id: tenant-001
Content-Type: application/json
```

Request pertama kali diterima oleh `http.Server` yang dibuat di:

- `internal/platform/http/server.go`

Server lalu meneruskan request ke handler utama, yaitu router.

### Langkah 3: router mencocokkan endpoint

Di:

- `internal/platform/http/router.go`

router akan:

1. memasang middleware global
2. membuat group `/api/v1`
3. memasang route private `/products`
4. untuk route private, memasang middleware `AuthJWT`

Karena request ini adalah:

```text
POST /api/v1/products
```

maka router mengarahkan request ke:

- `productHandler.Create`

yang ada di:

- `internal/modules/product/delivery_http.go`

Tetapi sebelum sampai ke handler, request harus melewati middleware dulu.

### Langkah 4: middleware `RequestID`

Middleware ini memberi ID unik ke request.

Tujuannya:
- memudahkan tracing log
- memudahkan debugging saat banyak request masuk bersamaan

ID ini nanti dipakai middleware logging.

### Langkah 5: middleware `CORS`

Middleware ini menambahkan header CORS sederhana.

Untuk request normal seperti `POST`, middleware akan meneruskan request ke langkah berikutnya.

Kalau request `OPTIONS`, middleware bisa mengakhiri request lebih awal dengan status `204 No Content`.

### Langkah 6: middleware `Logger`

Middleware logging mulai menghitung waktu mulai request.

Penting:
- log final baru ditulis **setelah** request selesai diproses
- jadi logging middleware membungkus seluruh alur request

Artinya dia bisa melihat status akhir response.

### Langkah 7: middleware `Recovery`

Middleware ini berjaga-jaga kalau ada panic.

Kalau suatu saat ada bug yang menyebabkan panic:
- server tidak mati
- panic diubah menjadi response error

Untuk flow normal create product, middleware ini hanya lewat saja.

### Langkah 8: middleware `Timeout`

Middleware ini memberi batas waktu request.

Tujuannya:
- request tidak menggantung terlalu lama
- server lebih aman saat ada operasi lambat

Kalau waktu habis, request akan dihentikan.

### Langkah 9: middleware `Tenant`

File:

- `internal/platform/http/middleware/tenant.go`

Middleware ini membaca header:

```http
X-Tenant-Id: tenant-001
```

Lalu:
- trim spasi
- memastikan nilainya tidak kosong
- menyimpan tenant ke context dengan helper `internal/shared/tenant/context.go`

Kalau header tenant tidak ada:
- request berhenti di sini
- response error langsung dikirim
- handler product tidak dijalankan

### Langkah 10: middleware `AuthJWT`

File:

- `internal/platform/http/middleware/auth.go`

Karena route `/products` adalah private route, middleware auth dijalankan.

Yang dilakukan middleware ini:

1. membaca header `Authorization`
2. mengambil bearer token
3. memanggil `jwtManager.ParseAccessToken(...)`
4. memvalidasi token
5. membandingkan tenant di token dengan tenant di header
6. membuat `auth.AuthUser`
7. menyimpan `AuthUser` ke context

Kalau token kosong, invalid, expired, atau tenant token tidak cocok:
- request berhenti di sini
- response error langsung dikirim

### Langkah 11: handler `product.Create` mulai bekerja

File:

- `internal/modules/product/delivery_http.go`

Method:

```go
func (h *HTTPDelivery) Create(w http.ResponseWriter, r *http.Request)
```

Handler melakukan hal-hal yang berhubungan dengan HTTP:

1. ambil user dari context auth
2. decode body JSON ke `CreateRequest`
3. validasi request DTO
4. panggil service
5. kirim response sukses atau error

Handler **tidak** menyimpan data langsung ke DB.

Itu penting. Handler hanya menjadi penghubung antara HTTP dan business logic.

### Langkah 12: handler mengambil user dari context

Handler memanggil:

```go
auth.FromContext(r.Context())
```

Data ini berasal dari middleware JWT tadi.

Jadi data user berpindah seperti ini:

```text
JWT Claims -> AuthUser -> context -> handler
```

Kalau ternyata user tidak ada di context:
- handler mengembalikan `401 Unauthorized`

Secara normal, ini jarang terjadi karena route product sudah dibungkus middleware auth.

### Langkah 13: body JSON di-parse ke DTO

Handler menjalankan:

```go
var req CreateRequest
json.NewDecoder(r.Body).Decode(&req)
```

DTO yang dipakai berasal dari:

- `internal/modules/product/dto.go`

Isi DTO:
- `name`
- `description`
- `price`
- `stock`

Kalau JSON tidak valid:
- handler mengembalikan `400 Bad Request`

### Langkah 14: request divalidasi

Handler memanggil validator:

```go
h.validate.Struct(req)
```

Validator instance dibuat saat startup di:

- `internal/shared/validator/validator.go`

Rule validasi berasal dari tag struct di DTO, misalnya:
- `required`
- `min`
- `max`
- `gte`

Kalau validasi gagal:
- handler membentuk `apperror.Validation(...)`
- detail field error diambil dari `apperror.FromValidator(...)`
- response error JSON dikirim

### Langkah 15: handler memanggil service

Jika parsing dan validasi sukses, handler memanggil:

```go
h.service.Create(r.Context(), authUser.TenantID, req)
```

Data yang dikirim ke service:
- `context.Context`
- `tenantID`
- `CreateRequest`

Sekarang kita masuk ke business logic.

### Langkah 16: service membuat entity product

File:

- `internal/modules/product/service.go`

Service menerima DTO lalu membentuk entity:

```text
CreateRequest -> product.Entity
```

Service mengisi:
- `TenantID` dari tenant user/context
- `Name` dari request
- `Description` dari request
- `Price` dari request
- `Stock` dari request

Kenapa `TenantID` tidak diambil dari body?
Karena tenant adalah bagian dari identitas request, bukan data bebas yang boleh diisi client sesuka hati.

### Langkah 17: service memanggil repository

Masih di `service.go`, service memanggil:

```go
s.repo.Create(ctx, entity)
```

Service bertanggung jawab atas aturan bisnis, repository bertanggung jawab atas query database.

### Langkah 18: repository menyimpan ke MySQL lewat GORM

File:

- `internal/modules/product/repository.go`

Repository menjalankan:

```go
r.db.WithContext(ctx).Create(entity)
```

`r.db` adalah instance GORM yang dibuat di:

- `internal/platform/db/gorm.go`

Alurnya:

```text
service -> repository -> GORM -> MySQL
```

Jika insert berhasil:
- entity akan terisi `ID`, `CreatedAt`, `UpdatedAt`

Jika insert gagal:
- error dikembalikan ke service
- lalu naik ke handler
- lalu diubah menjadi response error

### Langkah 19: service invalidasi cache Redis

Setelah berhasil simpan ke MySQL, service memanggil:

```go
s.invalidateCache(ctx, tenantID, entity.ID)
```

Di balik layar, ini memakai:

- `internal/modules/product/cache.go`

Cache yang dihapus:
- `products:list:<tenant>`
- `products:get:<tenant>:<id>`

Kenapa cache dihapus, bukan langsung diisi ulang?
Karena untuk create product, cara paling sederhana dan aman adalah:
- hapus cache yang mungkin stale
- biarkan request list/detail berikutnya membangun cache baru

Kalau invalidasi cache gagal:
- request create product **tetap sukses**
- service hanya menulis warning log

Ini keputusan yang masuk akal karena sumber data utama tetap MySQL, bukan Redis.

### Langkah 20: service mengubah entity menjadi response

Service memanggil mapper:

```go
toResponse(*entity)
```

Mapper ada di:

- `internal/modules/product/mapper.go`

Transformasi ini membuat response API tidak terlalu menempel langsung ke entity DB.

### Langkah 21: handler mengirim response sukses

Setelah service mengembalikan data response, handler memanggil:

```go
response.Success(w, http.StatusCreated, "product berhasil dibuat", product)
```

Helper ini ada di:

- `internal/shared/response/json.go`

Response sukses dibentuk konsisten seperti:

```json
{
  "success": true,
  "message": "product berhasil dibuat",
  "data": {
    "id": 1,
    "tenant_id": "tenant-001",
    "name": "Laptop A",
    "description": "Laptop kerja",
    "price": 15000000,
    "stock": 5,
    "created_at": "2026-03-28T20:00:00Z",
    "updated_at": "2026-03-28T20:00:00Z"
  }
}
```

### Langkah 22: middleware logging mencatat hasil akhir

Karena middleware logging membungkus seluruh request, setelah handler selesai, middleware logging menulis log seperti:

- request ID
- method
- path
- status
- duration
- remote address

Jadi log request umum terjadi di akhir flow, bukan di awal.

---

## 4. Contoh request HTTP

Contoh request yang benar:

```http
POST /api/v1/products HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer <token_jwt>
X-Tenant-Id: tenant-001
```

Contoh body JSON:

```json
{
  "name": "Laptop A",
  "description": "Laptop kerja untuk admin",
  "price": 15000000,
  "stock": 5
}
```

Kalau memakai `curl`:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: tenant-001" \
  -H "Authorization: Bearer <token_jwt>" \
  -d '{
    "name": "Laptop A",
    "description": "Laptop kerja untuk admin",
    "price": 15000000,
    "stock": 5
  }'
```

---

## 5. Contoh data yang mengalir antar layer

Supaya lebih terasa, kita lihat bentuk datanya.

### A. Request DTO

Masuk ke handler sebagai:

```go
type CreateRequest struct {
    Name        string
    Description string
    Price       int64
    Stock       int
}
```

Contoh nilainya:

```json
{
  "name": "Laptop A",
  "description": "Laptop kerja untuk admin",
  "price": 15000000,
  "stock": 5
}
```

### B. Data auth dari context

Dari middleware JWT, handler mendapat:

```go
auth.AuthUser{
    UserID:   10,
    TenantID: "tenant-001",
    Email:    "admin@tenant-001.com",
    Role:     "admin",
}
```

### C. Entity yang dibentuk service

Service membentuk:

```go
Entity{
    TenantID:    "tenant-001",
    Name:        "Laptop A",
    Description: "Laptop kerja untuk admin",
    Price:       15000000,
    Stock:       5,
}
```

### D. Data yang disimpan ke DB

Repository mengirim entity itu ke MySQL.

Setelah insert sukses, entity biasanya menjadi:

```go
Entity{
    ID:          25,
    TenantID:    "tenant-001",
    Name:        "Laptop A",
    Description: "Laptop kerja untuk admin",
    Price:       15000000,
    Stock:       5,
    CreatedAt:   <timestamp>,
    UpdatedAt:   <timestamp>,
}
```

### E. Response DTO

Mapper mengubah entity menjadi response:

```go
Response{
    ID:          25,
    TenantID:    "tenant-001",
    Name:        "Laptop A",
    Description: "Laptop kerja untuk admin",
    Price:       15000000,
    Stock:       5,
    CreatedAt:   <timestamp>,
    UpdatedAt:   <timestamp>,
}
```

Jadi transformasinya bisa dibaca seperti ini:

```text
JSON request
-> CreateRequest DTO
-> Entity
-> insert ke MySQL
-> Response DTO
-> JSON response
```

---

## 6. Jelaskan peran tiap layer

### Router
Tugas router adalah memilih endpoint yang benar dan menyusun urutan middleware.

Router tidak boleh berisi business logic create product.

### Middleware
Tugas middleware adalah concern umum yang berulang:
- request ID
- CORS
- logging
- recovery
- timeout
- tenant
- auth

Middleware tidak boleh membuat product, karena itu bukan concern umum.

### Handler / delivery_http
Tugas handler adalah urusan HTTP:
- baca header/body
- validasi input
- panggil service
- kirim response

Handler tidak boleh menulis SQL langsung.

### Service
Tugas service adalah business logic:
- membentuk entity
- menentukan tenant yang dipakai
- menentukan kapan cache harus dihapus
- menentukan apa yang dianggap sukses/gagal secara bisnis

### Repository
Tugas repository adalah akses data:
- `Create`
- `FindByID`
- `Update`
- `Delete`

Repository tidak boleh peduli soal JSON request atau response client.

### Redis / cache helper
Tugas cache helper adalah detail Redis:
- nama key
- get/set cache
- invalidasi cache

Dengan ini, `service.go` tidak penuh oleh detail string key Redis.

### Shared error
Tugas `shared/apperror` adalah membuat error konsisten dan aman dikirim ke client.

### Shared response
Tugas `shared/response` adalah membuat bentuk response sukses/error selalu konsisten.

---

## 7. Alur error

Sekarang kita bahas beberapa error penting.

### Kasus 1: token JWT tidak ada

- Muncul di layer: middleware auth
- File: `internal/platform/http/middleware/auth.go`
- Alur:
  - middleware tidak menemukan bearer token
  - memanggil `response.Error(...)`
  - `apperror.Unauthorized(...)` diubah menjadi JSON response

Contoh response:

```json
{
  "success": false,
  "message": "bearer token wajib diisi",
  "error": {
    "code": "UNAUTHORIZED"
  }
}
```

### Kasus 2: token invalid atau expired

- Muncul di layer: middleware auth
- Penyebab: `jwtManager.ParseAccessToken(...)` gagal
- Hasil: `401 Unauthorized`

### Kasus 3: tenant header tidak ada

- Muncul di layer: middleware tenant
- File: `internal/platform/http/middleware/tenant.go`
- Hasil: `400 Bad Request`

Request berhenti sebelum masuk handler.

### Kasus 4: body JSON invalid

- Muncul di layer: product delivery
- File: `internal/modules/product/delivery_http.go`
- Penyebab: `json.NewDecoder(...).Decode(...)` gagal
- Hasil: `400 Bad Request`

### Kasus 5: validasi gagal

- Muncul di layer: product delivery
- Penyebab: validator menemukan field yang tidak valid
- Hasil: `400 Bad Request` dengan detail field error

Contoh:

```json
{
  "success": false,
  "message": "validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "Name",
        "message": "nilai terlalu pendek / kecil"
      }
    ]
  }
}
```

### Kasus 6: database gagal menyimpan

- Muncul awalnya di layer: repository
- Naik ke: service
- Naik lagi ke: handler
- Di handler:
  - error di-log dengan `h.logger.Error(...)`
  - lalu dipanggil `response.Error(...)`

Kalau error itu bukan `AppError`, `apperror.As(...)` akan mengubahnya menjadi:
- `500 Internal Server Error`

### Kasus 7: cache Redis gagal

- Muncul di layer: product service
- Saat: invalidasi cache
- Perlakuan:
  - ditulis sebagai warning log
  - request create product tetap dianggap sukses

Kenapa?
Karena create product sudah berhasil masuk ke MySQL.

Redis di sini hanya cache, bukan sumber data utama.

### Kasus 8: product tidak valid

Contoh:
- `name` terlalu pendek
- `price` negatif
- `stock` negatif

Error ini muncul di:
- layer handler/delivery saat validasi DTO

Bukan di repository.

Kenapa?
Karena data input sebaiknya ditolak sedini mungkin sebelum menyentuh database.

---

## 8. Alur logging

### Request di-log di mana

Request umum di-log di:

- `internal/platform/http/middleware/logging.go`

Log ini mencatat:
- `request_id`
- `method`
- `path`
- `status`
- `duration`
- `remote_addr`

### Error di-log di mana

Error spesifik create product di-log di:

- `internal/modules/product/delivery_http.go`

saat `service.Create(...)` mengembalikan error.

Warning invalidasi cache di-log di:

- `internal/modules/product/service.go`

### Data yang aman untuk di-log

Contoh yang aman:
- request ID
- method
- path
- status code
- durasi
- tenant ID bila perlu
- user ID bila perlu
- nama endpoint
- jenis error

### Data yang tidak boleh di-log

Contoh yang sebaiknya tidak di-log:
- password
- full bearer token
- secret JWT
- data sensitif yang tidak perlu

Prinsip sederhananya:
- log secukupnya untuk tracing
- jangan log rahasia

---

## 9. Diagram alur sederhana

```text
Client
  -> HTTP Server
  -> Router
  -> RequestID Middleware
  -> CORS Middleware
  -> Logging Middleware
  -> Recovery Middleware
  -> Timeout Middleware
  -> Tenant Middleware
  -> AuthJWT Middleware
  -> Product HTTP Delivery
  -> Product Service
  -> Product Repository
  -> GORM
  -> MySQL
  -> Product Cache Helper
  -> Redis (invalidate cache)
  -> Response Helper
  -> Client
```

Kalau terjadi error, alurnya bisa berhenti di middleware, handler, atau service, lalu langsung masuk ke:

```text
AppError / Error Mapper
  -> Response Error JSON
  -> Client
```

---

## 10. Contoh trace request Create Product

Sekarang kita simulasi konkret.

### Data request

- tenant header = `tenant-001`
- user ID dari JWT = `10`
- product name = `Laptop A`
- price = `15000000`
- stock = `5`

### Trace langkah demi langkah

1. User mengirim `POST /api/v1/products`.
2. Server menerima request.
3. Router menemukan route `/api/v1/products` dengan method `POST`.
4. Middleware `RequestID` membuat request ID.
5. Middleware `CORS` menambahkan header CORS.
6. Middleware `Logger` mulai hitung durasi.
7. Middleware `Recovery` siap menangkap panic.
8. Middleware `Timeout` memberi batas waktu request.
9. Middleware `Tenant` membaca `X-Tenant-Id: tenant-001`, lalu menyimpannya ke context.
10. Middleware `AuthJWT` membaca bearer token.
11. JWT berhasil di-parse, lalu didapat claim user ID `10` dan tenant `tenant-001`.
12. Middleware memastikan tenant header cocok dengan tenant di token.
13. Middleware menyimpan `AuthUser{UserID: 10, TenantID: "tenant-001", ...}` ke context.
14. Handler `product.Create` dipanggil.
15. Handler mengambil `AuthUser` dari context.
16. Handler membaca body JSON:
    - `name = "Laptop A"`
    - `price = 15000000`
    - `stock = 5`
17. Handler memvalidasi DTO. Semua valid.
18. Handler memanggil `productService.Create(ctx, "tenant-001", req)`.
19. Service membentuk entity:
    - `TenantID = "tenant-001"`
    - `Name = "Laptop A"`
    - `Price = 15000000`
    - `Stock = 5`
20. Service memanggil repository create.
21. Repository menjalankan insert ke MySQL lewat GORM.
22. MySQL menyimpan row baru dan mengembalikan ID, misalnya `25`.
23. Service menghapus cache:
    - `products:list:tenant-001`
    - `products:get:tenant-001:25`
24. Service mengubah entity menjadi response DTO.
25. Handler mengirim response `201 Created`.
26. Middleware logger mencatat request selesai dengan status `201`.
27. Client menerima JSON sukses.

---

## 11. Hubungkan dengan struktur project

Sekarang kita hubungkan flow tadi dengan struktur `modules / platform / shared`.

### Kenapa logic bisnis tidak ditaruh di middleware
Karena middleware dipakai lintas endpoint.

Kalau create product ditaruh di middleware:
- flow jadi membingungkan
- tidak semua endpoint butuh logic itu
- business logic bercampur dengan concern umum

Business logic create product cocok di:
- `internal/modules/product/service.go`

### Kenapa parsing request tidak ditaruh di repository
Karena repository tugasnya database.

Kalau repository ikut parse JSON:
- layer jadi bercampur
- query DB tidak lagi fokus
- susah dites

Parsing request cocok di:
- `internal/modules/product/delivery_http.go`

### Kenapa JWT logic tidak ditaruh semua di handler
Karena JWT dipakai berulang di banyak endpoint private.

Kalau semua logic JWT ada di handler:
- handler jadi panjang
- parsing token terulang
- rawan tidak konsisten

JWT cocok dipisah menjadi:
- parser/generator di `internal/modules/auth/jwt.go`
- penggunaan di middleware auth

### Kenapa tenant helper cocok di shared
Karena tenant dipakai lintas modul:
- auth
- product
- middleware
- mungkin nanti order, category, invoice

Jadi tenant context cocok di:
- `internal/shared/tenant/context.go`

### Kenapa Redis dan DB cocok di platform
Karena keduanya adalah infrastruktur teknis, bukan fitur bisnis.

MySQL dan Redis bisa dipakai banyak modul.

Jadi mereka cocok di:
- `internal/platform/db`
- `internal/platform/cache`

---

## 12. Tambahkan panduan belajar

### Cara membaca alur request untuk endpoint lain

Pola sederhananya:

1. mulai dari `router.go`
2. cari endpoint yang dimaksud
3. lihat middleware apa saja yang membungkus route itu
4. masuk ke `delivery_http.go`
5. lihat method handler yang dipanggil
6. lanjut ke `service.go`
7. lanjut ke `repository.go`
8. cek helper tambahan seperti cache, mapper, atau error helper

Kalau kamu mengikuti urutan itu, kamu akan lebih mudah memahami endpoint lain.

### Cara menambah endpoint baru dengan pola yang sama

Misalnya kamu ingin menambah:
- `POST /api/v1/categories`
- `POST /api/v1/orders`

Pola berpikirnya:

1. buat DTO request
2. buat handler/delivery untuk parse dan validasi
3. buat service untuk business logic
4. buat repository untuk query DB
5. buat entity bila ada tabel baru
6. buat route di router
7. tambahkan cache hanya jika memang perlu
8. gunakan `response.Success` dan `response.Error`

### Cara tracing bug jika request gagal

Kalau request gagal, baca dari urutan ini:

1. cek endpoint di router benar atau tidak
2. cek middleware tenant dan auth
3. cek apakah body JSON valid
4. cek validasi DTO
5. cek service apakah business logic benar
6. cek repository/GORM query
7. cek log middleware
8. cek log error spesifik di handler/service

Cara paling aman adalah tracing dari luar ke dalam:

```text
client -> router -> middleware -> handler -> service -> repository -> DB/cache
```

---

## Ringkasan pola untuk fitur lain

Kalau nanti kamu mau membuat `Create Category` atau `Create Order`, pola berpikirnya tetap sama:

1. tentukan request HTTP-nya
2. tentukan middleware apa yang wajib jalan
3. parse dan validasi request di delivery
4. letakkan aturan bisnis di service
5. simpan data lewat repository
6. gunakan cache hanya jika memang bermanfaat
7. bentuk response konsisten
8. pastikan error dan log tetap mudah ditelusuri

Jadi inti pola project ini adalah:

```text
HTTP masuk
-> middleware umum
-> delivery menangani HTTP
-> service menangani business logic
-> repository menangani database
-> cache membantu performa
-> response helper membentuk output
```

Kalau kamu sudah paham alur `Create Product`, biasanya kamu akan jauh lebih mudah memahami endpoint lain di project ini.
