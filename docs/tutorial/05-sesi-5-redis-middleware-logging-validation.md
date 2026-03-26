# Sesi 5 — Redis, Middleware, Logging, Validation

## Tujuan sesi

Di sesi ini kamu akan belajar:

- bagaimana Redis dipakai secara sederhana
- urutan middleware dan alasan urutannya
- bagaimana request logging dibuat
- bagaimana validation request bekerja
- kenapa concern lintas route dipisah ke middleware dan helper

Sesi ini membantu kamu memahami hal-hal yang membuat backend terasa “lebih nyata” di dunia kerja.

---

## File yang dipelajari

1. `internal/product/cache.go`
2. `internal/platform/cache/redis.go`
3. `internal/platform/http/router.go`
4. `internal/platform/http/middleware/tenant.go`
5. `internal/platform/http/middleware/auth.go`
6. `internal/platform/http/middleware/request_logger.go`
7. `internal/platform/logger/zap.go`
8. `internal/platform/http/response.go`
9. `internal/auth/handler.go`
10. `internal/product/handler.go`

---

## Bagian 1 — Redis

### File: `internal/platform/cache/redis.go`

File ini membuat koneksi Redis dari config.

Yang terjadi:
- membaca `RedisAddr`, `RedisPassword`, `RedisDB`
- membuat `redis.NewClient`
- melakukan `Ping` untuk memastikan koneksi hidup

### Kenapa Ping dilakukan di awal
Supaya kalau Redis salah konfigurasi, aplikasi gagal sejak awal.
Lebih baik gagal cepat daripada error diam-diam saat request pertama masuk.

---

### File: `internal/product/cache.go`

File ini adalah cache helper khusus untuk product.

Method penting:
- `GetProduct`
- `SetProduct`
- `GetList`
- `SetList`
- `Invalidate`

### Key Redis yang dipakai
- detail product: `product:<tenant_id>:<id>`
- list product: `products:<tenant_id>`

### Kenapa key membawa tenant
Karena cache juga harus tenant-aware.
Kalau tidak, tenant A bisa saja membaca cache milik tenant B.

### TTL cache
Cache diset dengan TTL 5 menit.

### Alur cache hit
Contoh `GET /products/1`:
1. service cek Redis
2. kalau ada → langsung kembalikan
3. database tidak dipanggil

### Alur cache miss
1. service cek Redis
2. tidak ada
3. query database
4. hasil disimpan ke Redis
5. response dikirim

### Kapan cache dihapus
Saat:
- create
- update
- delete

Karena data berubah dan cache lama menjadi basi.

---

## Bagian 2 — Middleware

### File: `internal/platform/http/router.go`

Urutan middleware di project ini:

1. `RequestID`
2. `RealIP`
3. `Recoverer`
4. `Timeout`
5. `Tenant`
6. `RequestLogger`

Untuk route private, ditambah:
7. `AuthJWT`

---

## Kenapa urutannya seperti itu

### 1. Request ID
Supaya setiap request punya ID unik.
Ini membantu tracing saat melihat log.

### 2. Real IP
Supaya aplikasi tahu IP asli client bila ada proxy.

### 3. Recoverer
Supaya panic tidak membuat server mati total.
Middleware ini menangkap panic dan mengembalikan response yang aman.

### 4. Timeout
Supaya request yang terlalu lama tidak menggantung selamanya.

### 5. Tenant
Supaya context tenant sudah tersedia untuk handler dan middleware lain.

### 6. Request logger
Agar log request bisa menyertakan:
- request_id
- path
- method
- status
- duration
- tenant_id

### 7. AuthJWT
Dipasang hanya untuk route protected.
Tujuannya supaya endpoint public seperti login tidak dipaksa punya token.

---

## Bagian 3 — Middleware tenant

### File: `internal/platform/http/middleware/tenant.go`

Tugasnya:
- membaca `X-Tenant-Id`
- menolak request jika header tidak ada
- menyimpan tenant ke context

### Kenapa tenant disimpan ke context
Karena tenant dipakai banyak bagian:
- handler
- service
- auth middleware
- request logger

Dengan context, semua komponen bisa membaca tenant dengan cara yang seragam.

---

## Bagian 4 — Middleware auth

### File: `internal/platform/http/middleware/auth.go`

Tugasnya:
- membaca header `Authorization`
- mengambil token Bearer
- parse JWT
- memastikan tenant token sama dengan tenant request
- menyimpan `AuthUser` ke context

### Kenapa auth middleware dipisah
Kalau setiap handler parse token sendiri:
- code berulang
- rawan salah
- sulit dirawat

Middleware membuat semua route protected punya aturan yang sama.

---

## Bagian 5 — Request logging

### File: `internal/platform/http/middleware/request_logger.go`

Logger request membuat log untuk setiap request setelah handler selesai.

Data yang dicatat:
- request_id
- method
- path
- status
- duration
- tenant_id

### Kenapa status direkam dengan `statusRecorder`
Karena `http.ResponseWriter` standar tidak otomatis memberi kita status code final.
Jadi dibuat wrapper kecil agar status bisa disimpan.

### Kenapa log request penting
Karena saat ada masalah, log request membantu menjawab:
- endpoint mana yang dipanggil?
- status berapa?
- tenant mana?
- request berapa lama?

---

## Bagian 6 — Setup zap logger

### File: `internal/platform/logger/zap.go`

Project ini membedakan logger berdasarkan `AppEnv`.

- `production` → `zap.NewProduction()`
- selain itu → `zap.NewDevelopment()`

### Kenapa dibedakan
Saat local development, log yang lebih mudah dibaca lebih nyaman.
Saat production, format log biasanya lebih terstruktur dan ringkas.

---

## Bagian 7 — Validation

Validation ada di handler auth dan product.

Contoh:
- `internal/auth/handler.go`
- `internal/product/handler.go`

Validator yang dipakai: `validator/v10`

### Bagaimana alurnya
1. request JSON didecode ke struct
2. handler menjalankan `validate.Struct(req)`
3. jika error, response error dikirim

### Kenapa validasi dilakukan dekat DTO request
Karena validasi ini berkaitan langsung dengan input HTTP.
Misalnya:
- field wajib
- minimal panjang string
- nilai minimal angka

Ini paling natural dilakukan dekat request masuk.

### Contoh validasi product
- `name` minimal 3 karakter
- `price` tidak boleh negatif
- `stock` tidak boleh negatif

---

## Bagian 8 — Helper response validasi

### File: `internal/platform/http/response.go`

Ada fungsi `ValidationDetails(err)`.

Tugasnya mengubah error validator menjadi bentuk yang lebih mudah dibaca, misalnya:
- field apa yang gagal
- rule apa yang gagal
- parameter rule berapa

Ini membuat error validation lebih informatif untuk client dan untuk kamu saat belajar.

---

## Step by step memahami concern lintas fitur

### Langkah 1
Buka `router.go` dan catat urutan middleware.

### Langkah 2
Buka `tenant.go` lalu pahami bagaimana tenant masuk ke context.

### Langkah 3
Buka `auth.go` middleware lalu pahami bagaimana token diambil dari header.

### Langkah 4
Buka `request_logger.go` lalu lihat field log apa saja yang dicatat.

### Langkah 5
Buka `product/cache.go` lalu lihat key Redis untuk list dan detail.

### Langkah 6
Buka `product/service.go` lalu cari kapan `GetList`, `GetProduct`, dan `Invalidate` dipanggil.

---

## Latihan kecil sesi 5

### Latihan 1 — Redis
Jawab:
- key cache list product bentuknya seperti apa?
- key cache detail product bentuknya seperti apa?
- kenapa tenant dimasukkan ke key?

### Latihan 2 — Middleware
Jawab:
- middleware mana yang wajib jalan untuk semua route?
- middleware mana yang hanya untuk route private?

### Latihan 3 — Validation
Jawab:
- validasi login ada di file mana?
- validasi create product ada di file mana?
- hasil error validation dirapikan oleh fungsi apa?

---

## Kapan pakai Redis dan kapan tidak

### Pakai Redis saat:
- ada data yang sering dibaca
- datanya tidak harus selalu real-time per milidetik
- kamu ingin mengurangi beban query database

### Jangan buru-buru pakai Redis saat:
- datanya jarang dibaca
- kompleksitas invalidasi cache lebih besar dari manfaatnya
- project masih sangat kecil dan performa belum menjadi masalah

Di starter ini Redis dipakai **secukupnya** sebagai sarana belajar pola cache yang masuk akal.

---

## Kesalahan umum pemula di sesi ini

### 1. Meletakkan semua logic ke middleware
Middleware hanya untuk concern lintas route.
Business logic tetap di service.

### 2. Lupa invalidasi cache
Akibatnya data lama tetap muncul.

### 3. Menaruh tenant di body request
Tenant di project ini dianggap konteks request, jadi lebih konsisten di header.

### 4. Logging terlalu banyak atau terlalu sedikit
Log yang baik cukup informatif, tapi tidak membanjiri noise.

---

## Checklist pemahaman sesi 5

Kamu harus bisa menjawab:

- koneksi Redis dibuat di file mana?
- cache product dibuat di file mana?
- urutan middleware di router bagaimana?
- request logging ada di file mana?
- validasi request dilakukan di file mana?

---

## Ringkasan sesi 5

Di sesi ini kamu belajar bahwa:

- Redis dipakai untuk cache list dan detail product
- middleware mengurus concern lintas route seperti tenant, auth, logging
- zap dipakai untuk log terstruktur
- validasi request dilakukan di handler dengan DTO request
- response validation dirapikan supaya mudah dipahami

Sesi berikutnya akan fokus ke:
**exception dan error handling**.
