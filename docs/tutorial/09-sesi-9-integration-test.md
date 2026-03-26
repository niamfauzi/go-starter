# Sesi 9 — Integration Test

## Tujuan sesi

Di sesi ini kamu akan belajar:

- apa itu integration test
- bedanya dengan unit test
- kenapa integration test penting untuk API + DB + Redis
- bagaimana Testcontainers dipakai
- bagaimana test login, CRUD, dan tenant isolation dijalankan

---

## File yang dipelajari

1. `tests/integration/api_test.go`
2. `internal/platform/http/router.go`
3. `cmd/api/main.go`
4. `migrations/*.sql`
5. `cmd/seeder/main.go` (sebagai pembanding pola seed)

---

## Apa itu integration test

Integration test adalah test yang memeriksa apakah beberapa komponen benar-benar bekerja bersama.

Di project ini, integration test memeriksa kerja sama antara:
- HTTP router
- handler
- service
- repository
- MySQL
- Redis
- JWT auth
- tenant isolation

### Bedanya dengan unit test

#### Unit test
Menguji satu unit kecil secara terisolasi.
Biasanya dependency di-mock.

Contoh:
- test satu method service tanpa DB asli

#### Integration test
Menguji alur yang lebih nyata.
Biasanya memakai dependency sungguhan.

Contoh:
- login lewat HTTP
- token benar-benar dibuat
- product benar-benar tersimpan di MySQL
- Redis cache benar-benar aktif

---

## Kenapa integration test penting di project ini

Karena project ini punya banyak bagian yang harus saling cocok:

- router harus mengarah ke handler yang benar
- middleware tenant dan auth harus bekerja
- GORM harus benar konek ke MySQL
- cache harus bisa bicara ke Redis
- JWT harus bisa dipakai di route protected
- tenant isolation harus aman

Kalau hanya unit test, banyak masalah integrasi bisa lolos.

---

## Bagian 1 — Struktur umum `api_test.go`

### File: `tests/integration/api_test.go`

Di file ini ada beberapa bagian utama:

- `TestMain`
- helper container MySQL dan Redis
- helper migration
- helper seed
- helper login dan request JSON
- test-case utama

---

## Bagian 2 — `TestMain`

`TestMain` adalah setup besar yang dijalankan sebelum test-test lain.

Yang dilakukan:
1. start MySQL container
2. start Redis container
3. bangun config test
4. buka koneksi GORM
5. jalankan migration
6. isi seed data
7. buat logger, validator, Redis client
8. buat auth dan product dependencies
9. buat router
10. jalankan `httptest.NewServer(...)`

### Kenapa `httptest.NewServer` dipakai
Supaya test bisa memanggil API sungguhan melalui HTTP, tetapi tanpa perlu menjalankan server manual di terminal terpisah.

---

## Bagian 3 — Testcontainers

### `startMySQL`
Menjalankan container MySQL untuk test.

### `startRedis`
Menjalankan container Redis untuk test.

### Kenapa ini enak untuk belajar
Karena test jadi:
- lebih realistis
- tidak terlalu bergantung ke environment lokal yang sudah “kotor”
- lebih mendekati kondisi nyata

Tapi memang syaratnya:
- Docker harus aktif

---

## Bagian 4 — Migration dan seed dalam test

Test tidak mengandalkan database lokalmu yang sudah berisi data.
Sebaliknya, test:
- membuat environment sendiri
- menjalankan migration sendiri
- mengisi seed data sendiri

Ini bagus karena test menjadi lebih dapat diulang dan lebih stabil.

---

## Bagian 5 — Test case yang ada

### `TestLoginSuccess`
Memastikan login berhasil dengan credential yang benar.

### `TestLoginFail`
Memastikan login gagal jika password salah.

### `TestCreateProductWithoutJWT`
Memastikan route product ditolak jika tanpa token.

### `TestProductCRUDAndTenantIsolation`
Test ini memeriksa alur besar:
- login untuk ambil token
- create product
- get product
- update product
- coba akses dengan tenant lain
- delete product

Ini test yang sangat bagus untuk dipelajari karena memperlihatkan satu flow utuh.

---

## Bagian 6 — Helper `doJSON`

Helper ini membuat request HTTP dengan lebih singkat.

Yang dilakukan:
- marshal body ke JSON
- buat request
- pasang `Content-Type`
- pasang `X-Tenant-Id`
- pasang bearer token jika ada
- kirim request

### Kenapa helper ini penting
Supaya test tidak penuh dengan kode berulang.

---

## Bagian 7 — Helper `loginAndGetToken`

Helper ini:
1. memanggil endpoint login
2. membaca response JSON
3. mengambil `access_token`
4. mengembalikannya

Dengan helper ini, test CRUD menjadi lebih sederhana dibaca.

---

## Step by step menjalankan integration test

### 1. Pastikan Docker aktif
Karena testcontainers membutuhkan Docker.

### 2. Jalankan test
Dari root project:

```bash
go test ./tests/integration/... -v
```

Atau:

```bash
go test ./... -v
```

### 3. Perhatikan waktu startup
Pertama kali test bisa lebih lambat karena container image perlu diunduh.

### 4. Baca hasil test
Kalau gagal, lihat:
- test mana yang gagal
- status code yang didapat
- apakah MySQL/Redis container benar-benar start

---

## Cara membaca test dengan nyaman sebagai pemula

Baca urutannya seperti ini:

1. `TestMain`
2. `doJSON`
3. `loginAndGetToken`
4. `TestLoginSuccess`
5. `TestCreateProductWithoutJWT`
6. `TestProductCRUDAndTenantIsolation`

Urutan ini membantumu melihat:
- setup global
- helper umum
- test sederhana
- test skenario lengkap

---

## Latihan kecil sesi 9

### Latihan 1
Buka `TestMain` dan jawab:
- kapan MySQL container dijalankan?
- kapan Redis container dijalankan?
- kapan router dibuat?

### Latihan 2
Buka `TestProductCRUDAndTenantIsolation`, lalu jawab:
- token diambil dari helper mana?
- bagaimana create product dilakukan?
- bagaimana tenant mismatch diuji?

### Latihan 3
Tambahkan satu test kecil sendiri:
- panggil `GET /api/v1/auth/me` dengan token valid
- pastikan status code 200

Ini latihan bagus untuk mulai menulis test sendiri.

---

## Kapan integration test sebaiknya ditambah

Tambahkan integration test saat:
- menambah fitur baru yang menyentuh DB
- menambah auth/permission baru
- menambah cache atau middleware penting
- menemukan bug yang pernah terjadi, lalu ingin mencegah bug itu muncul lagi

---

## Kesalahan umum pemula di integration test

### 1. Mengira test ini sama dengan unit test
Tidak. Integration test lebih besar cakupannya.

### 2. Lupa Docker aktif
Testcontainers akan gagal.

### 3. Mengandalkan database lokal yang sudah ada
Integration test sebaiknya mandiri.

### 4. Membuat test terlalu saling bergantung
Usahakan satu test tidak terlalu bergantung pada urutan test lain.

---

## Checklist pemahaman sesi 9

Kamu harus bisa menjawab:

- apa itu integration test?
- kenapa project ini memakai Testcontainers?
- `TestMain` menyiapkan apa saja?
- helper `doJSON` dipakai untuk apa?
- test mana yang memeriksa tenant isolation?

---

## Ringkasan sesi 9

Di sesi ini kamu belajar bahwa:

- integration test menguji kerja sama banyak komponen sekaligus
- project ini membuat MySQL dan Redis test sendiri lewat Testcontainers
- migration dan seed juga dijalankan dalam test
- alur login dan CRUD diuji melalui HTTP sungguhan
- tenant isolation ikut diuji agar keamanan multi-tenant lebih terjamin

Sesi terakhir akan membahas:
**cara menambah fitur baru dengan pola yang sama**.
