# Sesi 10 — Cara Menambah Fitur Baru

## Tujuan sesi

Di sesi ini kamu akan belajar cara berpikir saat ingin menambah fitur baru ke starter project ini.

Targetnya bukan hanya “bisa menambah file”, tetapi tahu:
- file apa yang harus dibuat
- urutan membuat fitur
- kapan perlu repository, service, handler
- kapan perlu cache
- kapan perlu migration
- bagaimana menjaga project tetap rapi dan tidak overkill

---

## Prinsip utama saat menambah fitur

Project ini memakai pola yang sederhana:

- **model**
- **repository**
- **service**
- **handler**

Kalau fitur baru memang butuh HTTP endpoint dan data di database, biasanya kamu akan menambah empat bagian ini.

Untuk concern lintas fitur, kamu pakai:
- `shared/` bila utilitas dipakai banyak fitur
- `platform/` bila itu kebutuhan teknis lintas fitur

---

## Contoh fitur baru: `categories`

Misalnya kamu ingin membuat fitur `categories`.

### Hasil akhir yang kamu butuhkan kira-kira:
- tabel `categories`
- package `internal/category`
- endpoint CRUD category
- dokumentasi Swagger
- integration test

---

## Urutan paling aman untuk menambah fitur

### Langkah 1 — Tentukan kebutuhan fitur
Tulis dulu:
- data apa yang disimpan?
- endpoint apa yang diperlukan?
- apakah butuh auth?
- apakah harus tenant-aware?
- apakah butuh cache?

Contoh category:
- `id`
- `tenant_id`
- `name`
- `description`

---

### Langkah 2 — Buat migration SQL
Tambahkan file baru, misalnya:
- `migrations/000003_create_categories.up.sql`
- `migrations/000003_create_categories.down.sql`

Kenapa mulai dari migration?
Karena fitur yang menyimpan data butuh schema dulu.

---

### Langkah 3 — Buat model
Buat file:
- `internal/category/model.go`

Isi:
- struct `Category`
- `CreateRequest`
- `UpdateRequest`

### Kenapa mulai dari model setelah migration
Karena model Go harus mencerminkan bentuk data yang memang akan kamu simpan.

---

### Langkah 4 — Buat repository
Buat file:
- `internal/category/repository.go`

Method yang biasanya dibutuhkan:
- `Create`
- `List`
- `GetByID`
- `Update`
- `Delete`

### Aturan penting
Kalau fitur tenant-aware, query harus selalu memasukkan `tenant_id`.

---

### Langkah 5 — Buat service
Buat file:
- `internal/category/service.go`

Service bertugas:
- mengatur alur
- memanggil repository
- mengubah error GORM menjadi `AppError`
- mengelola cache jika ada

---

### Langkah 6 — Buat handler
Buat file:
- `internal/category/handler.go`

Handler bertugas:
- decode request
- validasi request
- ambil path param
- panggil service
- kirim response

---

### Langkah 7 — Pasang route
Ubah:
- `internal/platform/http/router.go`

Tambahkan route category ke group private atau public sesuai kebutuhan.

---

### Langkah 8 — Tambah seeder bila perlu
Kalau category perlu data demo, tambahkan ke:
- `cmd/seeder/main.go`

---

### Langkah 9 — Tambah Swagger
Ubah:
- `docs/swagger.yaml`

Tambahkan:
- schema request category
- path endpoint category
- contoh response

---

### Langkah 10 — Tambah integration test
Tambahkan test ke:
- `tests/integration/api_test.go`

Minimal test:
- create category
- list category
- get by id
- update
- delete
- tenant isolation

---

## Template mental sederhana saat menambah fitur

Setiap kali mau menambah fitur, tanyakan 7 hal ini:

1. data apa yang disimpan?
2. migration SQL-nya bagaimana?
3. model Go-nya bagaimana?
4. endpoint HTTP-nya apa?
5. service logic-nya apa?
6. query database-nya apa?
7. test dan Swagger-nya sudah ditambah belum?

Kalau kamu menjawab 7 pertanyaan ini, biasanya fitur baru akan lebih rapi.

---

## Kapan perlu cache untuk fitur baru

Tidak semua fitur butuh Redis.

### Pakai cache kalau:
- data sering dibaca
- query cukup sering
- perubahan data tidak terlalu sering

### Tidak perlu cache kalau:
- fitur masih sederhana
- data jarang dibaca
- kompleksitas invalidasi lebih besar dari manfaat

Sebagai pemula, lebih baik buat fitur berjalan dulu tanpa cache.
Kalau sudah stabil dan masuk akal, baru tambahkan cache.

---

## Kapan perlu shared package baru

Tambahkan sesuatu ke `internal/shared/` jika benar-benar dipakai banyak fitur.

Contoh:
- helper pagination umum
- helper parsing common query
- helper format error umum

Jangan buru-buru memindahkan sesuatu ke `shared` kalau baru dipakai satu fitur.

---

## Kapan perlu platform package baru

Tambahkan sesuatu ke `internal/platform/` jika itu concern teknis lintas fitur.

Contoh:
- mailer
- storage client
- queue client
- metrics

---

## Langkah belajar paling enak saat membuat fitur sendiri

Sebagai pemula, ikuti strategi ini:

### 1. Tiru fitur product
Jadikan `internal/product` sebagai pola utama.

### 2. Buat versi kecil dulu
Jangan langsung banyak field dan banyak endpoint.

### 3. Pastikan CRUD dasar jalan
Sebelum tambah cache, auth role, atau fitur lain.

### 4. Tambahkan test
Minimal satu flow sukses dan satu flow gagal.

### 5. Rapikan dokumentasi
Tambahkan Swagger dan README bila perlu.

---

## Contoh checklist saat menambah `categories`

- [ ] migration up/down dibuat
- [ ] model category dibuat
- [ ] request DTO dibuat
- [ ] repository dibuat
- [ ] service dibuat
- [ ] handler dibuat
- [ ] route dipasang
- [ ] Swagger ditambah
- [ ] seeder ditambah jika perlu
- [ ] integration test ditambah

Checklist seperti ini sangat membantu agar tidak ada langkah yang lupa.

---

## Best practice untuk pemula

### 1. Tambah fitur sedikit demi sedikit
Jangan langsung menambah banyak file tanpa urutan.

### 2. Pegang pola yang sama
Kalau `product` sudah punya pola yang jelas, ikuti pola itu.

### 3. Tenant-aware sejak awal
Kalau project multi-tenant, jangan menunda tenant filter.

### 4. Pisahkan input request dari model DB
Ini membuat validasi dan perubahan data lebih aman.

### 5. Tambahkan test untuk skenario penting
Terutama auth dan tenant isolation.

---

## Common mistakes untuk pemula

### 1. Menyalin fitur lama tanpa mengerti
Akibatnya kamu bingung saat bug muncul.

### 2. Lupa menambah migration
Model sudah ada, tapi tabel belum ada.

### 3. Lupa update router
Handler sudah dibuat, tapi route belum dipasang.

### 4. Lupa update Swagger
API jalan, dokumentasi tertinggal.

### 5. Lupa test tenant isolation
Ini sangat penting di project seperti ini.

### 6. Menambah abstraksi terlalu cepat
Misalnya membuat generic repository besar padahal belum perlu.
Di starter ini, lebih baik tetap sederhana.

---

## Latihan akhir sesi 10

Coba buat rencana fitur `categories` di catatanmu.

Tulis:
1. struktur tabel
2. file yang akan kamu buat
3. endpoint yang akan kamu tambahkan
4. validasi request yang dibutuhkan
5. test minimal yang akan kamu tulis

Jangan langsung coding. Rencanakan dulu.
Ini kebiasaan baik yang sangat membantu.

---

## Ringkasan sesi 10

Di sesi ini kamu belajar bahwa:

- menambah fitur baru sebaiknya mengikuti pola yang sudah ada
- urutan yang baik adalah migration → model → repository → service → handler → router → test → docs
- tidak semua fitur perlu cache
- tenant-aware harus dijaga sejak awal
- dokumentasi dan test harus ikut bertambah bersama fitur

Kalau kamu sudah sampai sesi ini, berarti kamu sudah punya peta belajar yang cukup lengkap untuk memahami dan mengembangkan starter project ini.
