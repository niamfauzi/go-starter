# Sesi 4 — CRUD Products

## Tujuan sesi

Di sesi ini kamu akan belajar:

- bagaimana CRUD products dibangun
- bagaimana product selalu terikat ke tenant
- bagaimana handler, service, repository bekerja bersama
- kapan cache dibaca dan kapan cache dihapus
- bagaimana endpoint products dilindungi oleh auth

Fitur product adalah contoh terbaik untuk memahami alur backend end-to-end.

---

## File yang dipelajari

1. `internal/product/model.go`
2. `internal/product/repository.go`
3. `internal/product/service.go`
4. `internal/product/handler.go`
5. `internal/product/cache.go`
6. `internal/platform/http/router.go`

---

## Endpoint products

Di project ini semua endpoint products berada di route private:

- `GET /api/v1/products`
- `POST /api/v1/products`
- `GET /api/v1/products/{id}`
- `PUT /api/v1/products/{id}`
- `DELETE /api/v1/products/{id}`

### Kenapa semua route products dibuat private
Untuk starter multi-tenant seperti ini, keputusan paling aman adalah semua route products wajib auth.
Tujuannya:
- data tidak bocor
- alur security konsisten
- pemula lebih mudah memahami bahwa semua operasi product punya konteks user dan tenant

---

## Bagian 1 — Model product

### File: `internal/product/model.go`

Struct penting:

- `Product`
- `CreateRequest`
- `UpdateRequest`

### `Product`
Ini model utama untuk tabel `products`.

Field utamanya:
- `ID`
- `TenantID`
- `Name`
- `Description`
- `Price`
- `Stock`
- `CreatedAt`
- `UpdatedAt`
- `DeletedAt`

### `CreateRequest` dan `UpdateRequest`
Ini DTO untuk request body.

Validasinya:
- `name` wajib, min 3, max 150
- `description` max 500
- `price` wajib dan minimal 0
- `stock` wajib dan minimal 0

### Kenapa request dipisah dari model DB
Supaya:
- validasi request lebih jelas
- input HTTP tidak langsung dicampur dengan model database
- lebih aman dan mudah diubah

---

## Bagian 2 — Repository product

### File: `internal/product/repository.go`

Repository bertugas berbicara dengan MySQL melalui GORM.

Method yang ada:
- `Create`
- `List`
- `GetByID`
- `Update`
- `Delete`

### Contoh penting: `List`
Query selalu memakai:

- `Where("tenant_id = ?", tenantID)`

Artinya hanya data milik tenant aktif yang diambil.

### Contoh penting: `GetByID`
Query memakai:
- `tenant_id`
- `id`

Ini sangat penting untuk tenant isolation.

### Kenapa repository sengaja tipis
Karena repository sebaiknya fokus ke database saja.
Ia tidak perlu tahu validasi HTTP, JWT, atau format response.

---

## Bagian 3 — Service product

### File: `internal/product/service.go`

Service adalah pusat aturan bisnis product.

Method yang ada:
- `Create`
- `List`
- `GetByID`
- `Update`
- `Delete`

---

### Alur `Create`

1. service menerima `tenantID` dan `CreateRequest`
2. service membentuk `Product`
3. `TenantID` diisi dari context tenant
4. repository `Create` dipanggil
5. cache di-invalidate

### Kenapa cache dihapus setelah create
Karena list product lama bisa sudah tidak akurat.
Kalau tidak dihapus, client bisa mendapat data lama.

---

### Alur `List`

1. service cek Redis dulu
2. kalau cache ada → langsung kembalikan
3. kalau cache tidak ada → query ke database
4. simpan hasil ke Redis
5. kembalikan hasil ke client

Ini pola **cache aside** sederhana.

---

### Alur `GetByID`

1. service cek cache product per ID
2. kalau ada → langsung pakai
3. kalau tidak ada → ambil dari DB
4. kalau tidak ditemukan → `NotFound`
5. simpan ke cache
6. kembalikan hasil

---

### Alur `Update`

1. service cari product dulu berdasarkan tenant + id
2. kalau tidak ditemukan → `NotFound`
3. update field dari request
4. simpan ke DB
5. hapus cache list dan cache detail product

### Kenapa cari dulu sebelum update
Agar:
- tahu product benar-benar ada
- memastikan product milik tenant yang aktif
- menghindari update ke data yang salah

---

### Alur `Delete`

1. service cek dulu apakah product ada di tenant aktif
2. kalau tidak ada → `NotFound`
3. repository melakukan delete
4. cache dihapus

### Kenapa cek dulu sebelum delete
Supaya response lebih bermakna.
Kalau product tidak ada, API bisa mengembalikan `404` yang jelas.

---

## Bagian 4 — Handler product

### File: `internal/product/handler.go`

Handler product bertugas:
- ambil tenant dari context
- ambil ID dari path jika ada
- decode JSON body
- validasi request
- panggil service
- kirim response

### Method `Create`
- baca body
- validasi
- panggil `service.Create`

### Method `List`
- ambil tenant
- panggil `service.List`

### Method `GetByID`
- ambil `id` dari URL
- ubah ke angka
- panggil `service.GetByID`

### Method `Update`
- ambil `id`
- decode body
- validasi
- panggil `service.Update`

### Method `Delete`
- ambil `id`
- panggil `service.Delete`

### Kenapa parse ID dilakukan di handler
Karena parsing URL/path adalah urusan HTTP.
Service sebaiknya menerima data yang sudah bersih.

---

## Bagian 5 — Router products

### File: `internal/platform/http/router.go`

Bagian product route berada dalam group yang sudah memakai auth middleware.

Artinya:
- sebelum handler product berjalan
- token user harus valid dulu

Ini membuat handler product tidak perlu memikirkan login lagi.

---

## Step by step mencoba CRUD products

### 1. Login dulu
Ambil token lewat `/auth/login`.

### 2. Create product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Tenant-Id: tenant-demo" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Belajar",
    "description": "Laptop untuk belajar Go",
    "price": 7500000,
    "stock": 5
  }'
```

### 3. List products

```bash
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Tenant-Id: tenant-demo"
```

### 4. Get product by ID

```bash
curl http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Tenant-Id: tenant-demo"
```

### 5. Update product

```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Tenant-Id: tenant-demo" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Belajar Update",
    "description": "Deskripsi baru",
    "price": 7600000,
    "stock": 4
  }'
```

### 6. Delete product

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Tenant-Id: tenant-demo"
```

---

## Latihan kecil sesi 4

### Latihan 1
Buka `internal/product/service.go` lalu cari:
- di mana cache list dibaca?
- di mana cache detail dibaca?
- di mana cache dihapus setelah create/update/delete?

### Latihan 2
Buka `internal/product/repository.go` lalu jawab:
- query mana yang memastikan tenant isolation saat mengambil product by id?
- query mana yang memastikan delete hanya untuk tenant aktif?

### Latihan 3
Buka `internal/product/handler.go` lalu jawab:
- di method mana validasi create dilakukan?
- di method mana path param `id` di-parse?

---

## Cara berpikir saat membuat CRUD baru

Misalnya nanti kamu membuat fitur `categories`.

Polanya hampir sama:

1. buat `model.go`
2. buat DTO request
3. buat repository
4. buat service
5. buat handler
6. pasang route
7. kalau perlu cache, tambahkan pola yang sama

Pola ini sengaja sederhana supaya mudah diulang.

---

## Kesalahan umum pemula di CRUD

### 1. Lupa tenant filter di query
Ini sangat bahaya di project multi-tenant.

### 2. Menaruh validasi bisnis di handler semua
Handler cukup untuk validasi request. Aturan proses tetap di service.

### 3. Tidak menghapus cache setelah update/delete
Akibatnya data lama tetap muncul.

### 4. Update langsung tanpa cek data ada atau tidak
Akhirnya response error jadi tidak jelas.

---

## Checklist pemahaman sesi 4

Kamu harus bisa menjawab:

- product model ada di file mana?
- create/list/get/update/delete diproses di service mana?
- query GORM product ada di file mana?
- tenant filter dipakai di mana?
- cache dihapus kapan?

---

## Ringkasan sesi 4

Di sesi ini kamu belajar bahwa:

- CRUD dibagi menjadi handler, service, repository
- semua operasi product selalu memakai tenant
- list dan get by id memakai Redis sebagai cache
- create/update/delete harus menghapus cache
- route products berada di balik JWT auth

Setelah paham CRUD, sesi berikutnya akan membahas concern lintas fitur:
**Redis, middleware, logging, dan validation**.
