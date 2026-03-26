# Sesi 7 — Swagger

## Tujuan sesi

Di sesi ini kamu akan belajar:

- fungsi Swagger/OpenAPI di project ini
- cara membaca file `swagger.yaml`
- bagaimana endpoint auth dan products didokumentasikan
- bagaimana bearer auth ditulis di spec
- bagaimana memakai Swagger untuk memahami API sebagai pemula

---

## File yang dipelajari

1. `docs/swagger.yaml`
2. `internal/platform/http/router.go`

---

## Apa itu Swagger di project ini

Di project ini, Swagger disimpan sebagai file YAML manual:

- `docs/swagger.yaml`

Artinya, spec API **ditulis manual**, bukan di-generate otomatis dari komentar code.

### Kenapa ini bagus untuk pemula
Karena kamu bisa melihat:
- endpoint apa saja yang ada
- method apa yang dipakai
- header apa yang wajib
- body request seperti apa
- response yang diharapkan

Dengan kata lain, Swagger di sini adalah “peta API”.

---

## Bagian 1 — Cara file swagger di-serve

### File: `internal/platform/http/router.go`

Ada route:

- `r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))`

Artinya folder `docs/` disajikan sebagai file statis.

Jadi kamu bisa membuka file:

- `http://localhost:8080/docs/swagger.yaml`

---

## Bagian 2 — Struktur dasar swagger

Di awal file ada:

- `openapi: 3.0.3`
- `info`
- `servers`
- `components`
- `paths`

### `info`
Menjelaskan nama API dan versinya.

### `servers`
Menjelaskan base URL.
Di project ini: `http://localhost:8080`

### `components`
Berisi bagian yang dipakai ulang, misalnya:
- security scheme
- schema request

### `paths`
Berisi daftar endpoint.

---

## Bagian 3 — Security scheme bearer auth

Di `components.securitySchemes`, project ini mendefinisikan:

- `bearerAuth`
- tipe `http`
- scheme `bearer`
- format `JWT`

Ini memberitahu pembaca API bahwa route tertentu butuh Bearer token.

---

## Bagian 4 — Schema request

Di `components.schemas`, ada contoh:

- `LoginRequest`
- `ProductRequest`

### Kenapa schema ini penting
Karena request body bisa dipakai ulang dan dibaca lebih mudah.

Contoh `ProductRequest` memberi tahu:
- `name`
- `description`
- `price`
- `stock`

beserta contoh nilainya.

---

## Bagian 5 — Paths auth

### `/api/v1/auth/login`
Method: `POST`

Yang didokumentasikan:
- perlu header `X-Tenant-Id`
- perlu body email + password
- response 200 jika berhasil
- response 401 jika gagal

### `/api/v1/auth/me`
Method: `GET`

Yang didokumentasikan:
- perlu Bearer token
- perlu header `X-Tenant-Id`

---

## Bagian 6 — Paths product

Endpoint yang didokumentasikan:
- `GET /api/v1/products`
- `POST /api/v1/products`
- `GET /api/v1/products/{id}`
- `PUT /api/v1/products/{id}`
- `DELETE /api/v1/products/{id}`

Semua endpoint products:
- butuh bearer auth
- butuh `X-Tenant-Id`

### Kenapa ini penting
Karena multi-tenant dan auth adalah bagian inti dari perilaku API.
Jadi dokumentasi harus menampilkan keduanya dengan jelas.

---

## Cara memakai Swagger sebagai pemula

Sebagai pemula, jangan melihat Swagger hanya sebagai formalitas.
Gunakan Swagger untuk menjawab pertanyaan berikut:

- URL endpoint-nya apa?
- method-nya apa?
- apakah perlu auth?
- apakah perlu header tenant?
- body request-nya apa?
- response suksesnya bagaimana?

Kalau bingung memakai suatu endpoint, buka `swagger.yaml` dulu sebelum membaca code.

---

## Step by step membaca satu endpoint

Mari ambil contoh `POST /api/v1/auth/login`.

### Langkah 1
Cari path `/api/v1/auth/login`.

### Langkah 2
Lihat method `post`.

### Langkah 3
Lihat `parameters`.
Kamu akan tahu bahwa `X-Tenant-Id` wajib dikirim.

### Langkah 4
Lihat `requestBody`.
Kamu akan tahu schema yang dipakai adalah `LoginRequest`.

### Langkah 5
Lihat `components.schemas.LoginRequest`.
Kamu akan tahu field:
- email
- password

### Langkah 6
Lihat `responses`.
Kamu akan tahu kemungkinan response utama.

---

## Kapan Swagger manual cocok

Swagger manual cocok saat:
- project masih kecil-menengah
- kamu ingin belajar API contract dengan jelas
- kamu belum ingin menambah generator/tooling lain
- kamu ingin menjaga dokumentasi tetap eksplisit

---

## Keterbatasan Swagger manual

Perlu jujur juga:
- kamu harus update file YAML sendiri saat endpoint berubah
- kalau lupa update, doc bisa tidak sinkron dengan code

Tapi untuk starter belajar, ini justru bagus karena kamu dipaksa memahami contract API.

---

## Latihan kecil sesi 7

### Latihan 1
Buka `docs/swagger.yaml`, lalu cari:
- schema `LoginRequest`
- schema `ProductRequest`

### Latihan 2
Cari endpoint `GET /api/v1/products/{id}` lalu jawab:
- parameter path apa yang diperlukan?
- header apa yang wajib?
- apakah endpoint ini butuh auth?

### Latihan 3
Bandingkan `swagger.yaml` dengan `router.go`.
Cocokkan apakah semua route di router sudah ada di Swagger.

---

## Checklist pemahaman sesi 7

Kamu harus bisa menjawab:

- file Swagger ada di mana?
- route untuk mengakses file Swagger di browser ada di mana?
- security scheme bearer ditulis di bagian mana?
- request schema product ada di bagian mana?
- kenapa Swagger manual cocok untuk starter belajar?

---

## Ringkasan sesi 7

Di sesi ini kamu belajar bahwa:

- API didokumentasikan lewat `docs/swagger.yaml`
- file Swagger disajikan sebagai file statis lewat router
- route auth dan products sudah didokumentasikan
- bearer auth dan header tenant dijelaskan di spec
- Swagger manual membantu kamu memahami contract API dengan lebih jelas

Sesi berikutnya akan membahas cara menjalankan service pendukung dan debug:
**Docker Compose dan debug**.
