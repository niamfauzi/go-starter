# Sesi 1 — Gambaran Besar dan Struktur Project

## Tujuan sesi

Di sesi ini kamu akan memahami:

- project ini dibuat untuk apa
- alur request dari HTTP sampai database
- fungsi tiap folder utama
- kenapa struktur project dibuat package-based
- file mana yang perlu dibaca dulu sebagai pemula

Sesi ini penting karena banyak pemula langsung lompat ke code detail, padahal belum punya peta besar. Akibatnya code terasa rumit. Jadi, di sesi 1 kita bikin dulu “peta project”-nya.

---

## Gambaran besar project

Starter project ini adalah backend Go yang sengaja dibuat:

- **cukup rapi**
- **tidak terlalu enterprise**
- **mudah dipelajari**
- **cukup dekat dengan praktik industri**

Fitur utamanya:

- login dengan JWT
- endpoint profile `/auth/me`
- CRUD products
- tenant isolation lewat header `X-Tenant-Id`
- MySQL sebagai database utama
- Redis untuk cache sederhana
- middleware umum seperti request id, timeout, recovery, auth, tenant
- logging terstruktur dengan zap
- migration manual dengan file SQL
- integration test

---

## Kenapa pakai struktur package-based

Project ini **tidak** memakai struktur berat seperti:

- `domain`
- `application`
- `infrastructure`
- `interface`
- `usecase`
- `port`

Sebaliknya, project ini memakai struktur berdasarkan **fitur/package**, contohnya:

- `internal/auth`
- `internal/product`

Pendekatan ini enak untuk pemula karena saat kamu ingin belajar auth, kamu tinggal masuk ke folder `internal/auth`. Saat kamu ingin belajar product, kamu tinggal masuk ke `internal/product`.

Artinya, code lebih mudah dicari dan tidak terlalu tersebar.

---

## Struktur folder utama

Berikut struktur utamanya:

```text
.
├── cmd/
│   ├── api/
│   └── seeder/
├── configs/
├── docs/
├── internal/
│   ├── auth/
│   ├── product/
│   ├── platform/
│   │   ├── cache/
│   │   ├── db/
│   │   ├── http/
│   │   │   └── middleware/
│   │   └── logger/
│   └── shared/
│       ├── apperror/
│       └── tenant/
├── migrations/
├── scripts/
├── tests/integration/
└── README.md
```

---

## Fungsi tiap folder

### `cmd/`
Folder ini berisi **entry point** aplikasi.

- `cmd/api/main.go` → titik awal aplikasi HTTP
- `cmd/seeder/main.go` → command untuk mengisi data awal

Sebagai pemula, anggap `cmd/` sebagai tombol “jalankan program”.

### `configs/`
Berisi loader environment variable.

Fungsinya:
- membaca `.env`
- menyusun config ke dalam struct
- membuat `main.go` tidak penuh dengan `os.Getenv`

### `docs/`
Berisi file Swagger/OpenAPI manual.

Di project ini isinya `swagger.yaml`.

### `internal/auth/`
Semua code yang berhubungan dengan auth:

- model user
- request login
- JWT manager
- repository user
- service auth
- handler auth
- auth context

### `internal/product/`
Semua code yang berhubungan dengan products:

- model product
- request create/update
- repository product
- service product
- handler product
- cache helper untuk product

### `internal/platform/`
Berisi hal-hal teknis yang dipakai lintas fitur.

- `cache/` → koneksi Redis
- `db/` → koneksi GORM MySQL
- `http/` → router, response helper, middleware
- `logger/` → setup zap

### `internal/shared/`
Berisi utilitas yang dipakai banyak fitur.

- `apperror/` → custom application error
- `tenant/` → context helper untuk tenant

### `migrations/`
Berisi file SQL manual untuk membuat tabel.

### `scripts/`
Berisi helper script migration.

### `tests/integration/`
Berisi integration test yang menjalankan dependency nyata.

---

## Alur request dari awal sampai akhir

Supaya lebih mudah, mari ikuti satu contoh request: `GET /api/v1/products`.

### Langkah 1 — request masuk ke router
File yang berperan: `internal/platform/http/router.go`

Router menentukan URL mana ditangani handler mana.

### Langkah 2 — middleware berjalan
Beberapa middleware yang berjalan:

- request ID
- real IP
- recoverer
- timeout
- tenant
- request logger

Kalau route butuh login, ada auth middleware juga.

### Langkah 3 — tenant dibaca dari header
Middleware tenant membaca `X-Tenant-Id`, lalu menyimpannya ke context.

### Langkah 4 — handler menerima request
Contoh: `internal/product/handler.go`

Handler akan:
- ambil tenant dari context
- baca parameter/body
- validasi input
- panggil service

### Langkah 5 — service menjalankan business logic
Contoh: `internal/product/service.go`

Service akan:
- cek cache dulu bila perlu
- panggil repository
- invalidasi cache saat data berubah
- mengubah error database menjadi error aplikasi

### Langkah 6 — repository akses MySQL
Contoh: `internal/product/repository.go`

Repository berbicara ke MySQL melalui GORM.

### Langkah 7 — response dikirim
File yang berperan: `internal/platform/http/response.go`

Semua response dibungkus dengan format yang konsisten.

---

## File pertama yang wajib dibaca

Kalau kamu benar-benar baru, baca file ini dulu secara berurutan:

1. `README.md`
2. `cmd/api/main.go`
3. `internal/platform/http/router.go`
4. `internal/auth/handler.go`
5. `internal/auth/service.go`
6. `internal/product/handler.go`
7. `internal/product/service.go`
8. `internal/product/repository.go`

Kenapa urutannya begitu?

Karena:
- `README` memberi gambaran
- `main.go` menunjukkan wiring
- `router.go` menunjukkan rute request
- handler/service/repository menunjukkan alur inti aplikasi

---

## Penjelasan sederhana file `cmd/api/main.go`

Di file ini semua dependency “dirakit”:

- load config
- buat logger
- buka koneksi MySQL
- buka koneksi Redis
- buat validator
- buat repository auth dan product
- buat service auth dan product
- buat handler auth dan product
- pasang router
- jalankan server HTTP

Sebagai pemula, pahami satu hal:
**`main.go` bukan tempat business logic.**
Fungsinya hanya merakit komponen.

---

## Kenapa handler, service, repository dipisah

Ini salah satu bagian paling penting.

### Handler
Tugasnya:
- menerima HTTP request
- decode JSON
- validasi input
- mengirim response

### Service
Tugasnya:
- aturan bisnis
- alur proses
- koordinasi cache/repository
- mapping error

### Repository
Tugasnya:
- query ke database

Kalau semua dicampur di satu file:
- code cepat berantakan
- susah dites
- susah dibaca
- susah diubah

---

## Latihan kecil sesi 1

Lakukan ini:

1. buka `README.md`
2. buka `cmd/api/main.go`
3. tulis di catatanmu:
   - config dibuat di mana
   - db dibuat di mana
   - router dibuat di mana
   - auth dibuat di mana
   - product dibuat di mana

Tujuan latihan ini supaya kamu hafal “pusat wiring” project.

---

## Checklist pemahaman sesi 1

Kalau kamu sudah paham sesi ini, kamu harus bisa menjawab:

- apa fungsi `cmd/api/main.go`?
- kenapa `auth` dan `product` dipisah ke package sendiri?
- apa bedanya `platform` dan `shared`?
- kenapa ada folder `migrations`?
- request HTTP lewat file mana dulu?

Kalau belum bisa menjawab, ulangi baca bagian struktur dan alur request.

---

## Ringkasan sesi 1

Di sesi ini kamu belajar bahwa:

- project disusun berdasarkan fitur, bukan layer berat
- `main.go` adalah tempat wiring, bukan tempat logic
- request mengalir dari router → middleware → handler → service → repository → response
- folder `auth` dan `product` adalah dua fitur utama
- folder `platform` menampung kebutuhan teknis lintas fitur

Setelah ini, sesi berikutnya akan fokus ke pondasi teknis:
**config, MySQL, migration, dan seeder**.
