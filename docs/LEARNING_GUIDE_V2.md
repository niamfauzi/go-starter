# Learning Guide V2

Panduan ini menjelaskan refactor project ke struktur **versi siap berkembang**.

Fokus utamanya:
- tetap ramah untuk pemula
- tetap package-based
- lebih siap saat fitur bertambah
- tidak menambah abstraction yang belum perlu

Module project tetap:

```go
github.com/niamfauzi/go-starter
```

## Sesi 1: Gambaran besar refactor ke versi siap berkembang

### Tujuan sesi
Memahami kenapa struktur lama perlu dirapikan, dan apa perubahan intinya.

### Gambaran besar perubahan
Sebelumnya, fitur bisnis seperti auth dan product langsung berada di:

```text
internal/auth
internal/product
```

Sekarang, fitur bisnis dipindah ke:

```text
internal/modules/auth
internal/modules/product
internal/modules/user
```

Lalu komponen teknis dipusatkan di:

```text
internal/platform
```

Dan helper lintas modul dipusatkan di:

```text
internal/shared
```

### Kenapa `auth` dan `product` dipindah ke `modules/`
Karena keduanya adalah fitur bisnis.

Cara berpikir sederhananya:
- kalau file itu membahas login, token, profile user, berarti itu business feature auth
- kalau file itu membahas create/update/delete/list product, berarti itu business feature product
- fitur bisnis sebaiknya hidup berdekatan agar mudah dicari saat project membesar

### Kenapa `platform/` dipisah dari `modules/`
Karena DB, Redis, HTTP server, logger, dan middleware bukan fitur bisnis.

Contoh:
- MySQL bukan fitur auth
- Redis bukan fitur product
- router bukan fitur user

Semua itu adalah infrastruktur teknis yang dipakai banyak modul.

### Kenapa `shared/` dibutuhkan
Karena ada helper yang dipakai lebih dari satu modul, misalnya:
- `apperror`
- `response`
- `validator`
- `tenant`

Kalau helper seperti ini ditaruh di module auth atau product, module lain akan terasa “menumpang”.

### Kapan sesuatu masuk `modules/`, `platform/`, atau `shared/`
- masuk `modules/` kalau itu logic bisnis per fitur
- masuk `platform/` kalau itu infrastruktur teknis
- masuk `shared/` kalau itu helper umum lintas modul

### Trade-off dibanding versi sederhana
Kelebihan:
- lebih rapi saat fitur bertambah
- import path lebih jelas
- folder business dan technical concern tidak bercampur

Konsekuensi:
- jumlah folder bertambah
- wiring dependency sedikit lebih panjang
- pemula perlu beberapa menit lebih lama untuk mengenali struktur

Trade-off ini masih masuk akal untuk project kecil-menengah yang mulai berkembang.

## Sesi 2: Menata folder modules, platform, dan shared

### Tujuan sesi
Memetakan struktur final yang dipakai setelah refactor.

### Struktur folder final

```text
.
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── migrate/
│   │   └── main.go
│   └── seeder/
│       └── main.go
├── configs/
│   ├── config.go
│   └── env.go
├── docs/
│   ├── LEARNING_GUIDE.md
│   ├── LEARNING_GUIDE_V2.md
│   ├── README.md
│   ├── swagger.html
│   ├── swagger.yaml
│   └── tutorial/
├── internal/
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── context.go
│   │   │   ├── delivery_http.go
│   │   │   ├── dto.go
│   │   │   ├── entity.go
│   │   │   ├── jwt.go
│   │   │   ├── password.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── product/
│   │   │   ├── cache.go
│   │   │   ├── delivery_http.go
│   │   │   ├── dto.go
│   │   │   ├── entity.go
│   │   │   ├── mapper.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── user/
│   │       ├── dto.go
│   │       ├── entity.go
│   │       ├── repository.go
│   │       └── service.go
│   ├── platform/
│   │   ├── cache/
│   │   │   └── redis.go
│   │   ├── db/
│   │   │   ├── gorm.go
│   │   │   ├── migrator.go
│   │   │   └── mysql.go
│   │   ├── http/
│   │   │   ├── router.go
│   │   │   ├── server.go
│   │   │   └── middleware/
│   │   │       ├── auth.go
│   │   │       ├── cors.go
│   │   │       ├── logging.go
│   │   │       ├── recovery.go
│   │   │       ├── requestid.go
│   │   │       ├── tenant.go
│   │   │       └── timeout.go
│   │   └── logger/
│   │       └── zap.go
│   └── shared/
│       ├── apperror/
│       │   ├── error.go
│       │   └── mapper.go
│       ├── response/
│       │   └── json.go
│       ├── tenant/
│       │   └── context.go
│       └── validator/
│           └── validator.go
├── migrations/
├── scripts/
├── tests/
│   └── integration/
│       └── api_test.go
├── docker-compose.yml
├── go.mod
└── README.md
```

### Catatan penting
Saya tetap mempertahankan `cmd/migrate/main.go`.

Alasannya sederhana:
- migration manual SQL masih bagian penting dari workflow
- command terpisah membuat pemula lebih mudah menjalankan `up` dan `down`
- ini tidak membuat struktur menjadi terlalu enterprise

Jadi ada sedikit penyesuaian praktis dibanding target awal, tetapi tetap konsisten dengan tujuan project.

## Sesi 3: Refactor auth module

### Tujuan sesi
Memisahkan tanggung jawab auth agar lebih jelas.

### File yang berubah
- `internal/modules/auth/delivery_http.go`
- `internal/modules/auth/service.go`
- `internal/modules/auth/repository.go`
- `internal/modules/auth/entity.go`
- `internal/modules/auth/dto.go`
- `internal/modules/auth/jwt.go`
- `internal/modules/auth/password.go`
- `internal/modules/auth/context.go`

### Kenapa file auth dipecah seperti ini
- `delivery_http.go`: fokus HTTP request/response
- `service.go`: fokus business logic login dan profile
- `repository.go`: akses data untuk kebutuhan auth
- `entity.go`: menyimpan `AuthUser` yang aman untuk context
- `dto.go`: request/response auth
- `jwt.go`: generate dan parse JWT
- `password.go`: hashing dan compare password
- `context.go`: simpan user login ke context

### Kenapa `password.go` dipisah
Supaya hashing password tidak tersebar di banyak file.

Kalau nanti ada kebutuhan:
- ganti cost bcrypt
- tambah helper verifikasi
- pindah ke argon2

maka perubahan terpusat.

### Contoh import path baru

```go
import "github.com/niamfauzi/go-starter/internal/modules/auth"
```

### Hasil akhir sederhana
Layer auth sekarang mudah dibaca:
- HTTP ada di delivery
- login logic ada di service
- token ada di jwt
- password ada di password helper

## Sesi 4: Refactor product module

### Tujuan sesi
Merapikan module product supaya lebih siap berkembang.

### File yang berubah
- `internal/modules/product/delivery_http.go`
- `internal/modules/product/service.go`
- `internal/modules/product/repository.go`
- `internal/modules/product/entity.go`
- `internal/modules/product/dto.go`
- `internal/modules/product/cache.go`
- `internal/modules/product/mapper.go`

### Kenapa product punya `cache.go`
Karena key Redis dan operasi cache product adalah detail teknis yang spesifik untuk product.

Kalau semua detail Redis ditaruh di `service.go`, file service cepat menjadi panjang.

### Kenapa product punya `mapper.go`
Karena entity database tidak selalu ideal langsung dikirim ke client.

Mapper dipakai saat:
- ingin memisahkan bentuk data DB dan response API
- ingin menyembunyikan field tertentu
- ingin menjaga handler tetap bersih

### Kenapa DTO dan entity dipisah
Karena:
- `entity` merepresentasikan tabel database
- `dto` merepresentasikan kontrak request/response API

Keduanya terlihat mirip sekarang, tetapi perannya berbeda.

### Kapan perlu projection struct
Projection cocok saat modul lain hanya butuh sebagian field.

Contoh:
- butuh `id`, `name`, `price` saja
- tidak butuh semua kolom `products`

Di kondisi itu, jangan duplikasi entity penuh. Buat DTO/projection yang memang sesuai kebutuhan.

## Sesi 5: Menambahkan user module bila diperlukan

### Tujuan sesi
Menentukan siapa pemilik utama tabel `users`.

### Kenapa `user` module dibuat
Karena tabel `users` dipakai lintas kebutuhan:
- auth perlu membaca user untuk login
- seeder perlu membuat user demo
- integration test perlu membuat user test

Kalau entity `users` tetap disimpan di auth, lama-lama auth akan menjadi “pemilik” tabel yang sebenarnya dipakai lintas fitur.

Itu kurang rapi.

### Aturan penting
**1 tabel database = 1 entity utama**

Artinya:
- tabel `users` punya entity utama di `internal/modules/user/entity.go`
- tabel `products` punya entity utama di `internal/modules/product/entity.go`

### Kenapa entity utama tidak boleh diduplikasi
Kalau diduplikasi:
- mapping GORM bisa tidak sinkron
- field bisa berbeda antar modul
- perubahan schema rawan lupa diupdate di semua tempat

Lebih aman punya satu entity utama, lalu modul lain memakai DTO/projection bila hanya butuh sebagian data.

### Kenapa `auth/repository.go` masih ada
Karena kebutuhan query auth tetap terasa sebagai concern auth.

Namun repository auth sekarang hanya menjadi lapisan auth-specific di atas module user. Ini kompromi yang cukup rapi dan tetap mudah dipahami pemula.

## Sesi 6: Menata middleware, HTTP server, logger, DB, Redis di platform

### Tujuan sesi
Memisahkan komponen teknis dari logic bisnis.

### File yang berubah
- `internal/platform/db/mysql.go`
- `internal/platform/db/gorm.go`
- `internal/platform/db/migrator.go`
- `internal/platform/cache/redis.go`
- `internal/platform/http/router.go`
- `internal/platform/http/server.go`
- `internal/platform/http/middleware/*.go`
- `internal/platform/logger/zap.go`

### Kenapa `mysql.go` dan `gorm.go` dipisah
Karena keduanya punya tujuan berbeda:
- `mysql.go` untuk `database/sql`, terutama migration manual
- `gorm.go` untuk ORM aplikasi sehari-hari

Pemisahan ini membuat niat file lebih jelas.

### Kenapa middleware dipecah
Supaya tiap concern punya file sendiri:
- `auth.go`
- `tenant.go`
- `logging.go`
- `cors.go`
- `requestid.go`
- `recovery.go`
- `timeout.go`

Ini masih masuk akal karena middleware memang concern lintas endpoint.

### Kenapa `server.go` ditambahkan
Supaya pembuatan `http.Server` tidak menumpuk di `main.go`.

`main.go` jadi lebih fokus pada wiring dependency.

## Sesi 7: Menata app error, response helper, validator helper, tenant helper di shared

### Tujuan sesi
Membuat helper lintas modul punya rumah yang jelas.

### File yang berubah
- `internal/shared/apperror/error.go`
- `internal/shared/apperror/mapper.go`
- `internal/shared/response/json.go`
- `internal/shared/validator/validator.go`
- `internal/shared/tenant/context.go`

### Kenapa `apperror` dipecah jadi `error.go` dan `mapper.go`
Karena ada dua tugas berbeda:
- mendefinisikan bentuk `AppError`
- memetakan error biasa menjadi `AppError`

Pemisahan ini masih sederhana dan membantu readability.

### Kenapa validator dibuat helper sendiri
Supaya nanti kalau ada kebutuhan:
- custom validation tag
- custom translation
- custom naming field

kita punya tempat yang jelas untuk menaruhnya.

### Kenapa response helper tetap sederhana
Karena kita belum butuh abstraction yang lebih berat.

`response.Success` dan `response.Error` sudah cukup untuk project ini.

## Sesi 8: Update router, wiring, import path, dan main.go

### Tujuan sesi
Menyambungkan semua package baru.

### File yang berubah
- `cmd/api/main.go`
- `internal/platform/http/router.go`

### Contoh import path baru

```go
import (
    "github.com/niamfauzi/go-starter/internal/modules/auth"
    "github.com/niamfauzi/go-starter/internal/modules/product"
    "github.com/niamfauzi/go-starter/internal/modules/user"
)
```

### Kenapa wiring di `main.go` masih manual
Karena project ini sengaja tidak memakai dependency injection framework.

Untuk starter belajar, wiring manual lebih baik karena:
- alur dependency terlihat jelas
- pemula mudah mengikuti arah data
- tidak ada magic tambahan

## Sesi 9: Update integration test, seeder, migration, dan Swagger bila perlu

### Tujuan sesi
Menjaga fitur lama tetap hidup setelah struktur diubah.

### File yang berubah
- `tests/integration/api_test.go`
- `cmd/seeder/main.go`
- `cmd/migrate/main.go`
- `docs/swagger.yaml`
- `docs/swagger.html`

### Penyesuaian penting
- seeder sekarang memakai `user.Entity` dan helper `auth.HashPassword`
- integration test memakai import path baru
- dokumentasi swagger pindah dari `openapi.yaml` ke `swagger.yaml`

### Kenapa migration tetap SQL manual
Karena ini memang aturan penting project:
- schema harus jelas
- perubahan database harus eksplisit
- belajar migration manual lebih berguna untuk project nyata

Jadi meskipun kita memakai GORM untuk query aplikasi, migration tetap tidak memakai AutoMigrate.

## Sesi 10: Panduan cara menambah modul baru dengan struktur ini

### Tujuan sesi
Memberi pola berpikir yang bisa diulang untuk fitur baru.

### Contoh menambah modul `category`

1. Buat folder `internal/modules/category`
2. Tambah `entity.go`
3. Tambah `dto.go`
4. Tambah `repository.go`
5. Tambah `service.go`
6. Tambah `delivery_http.go`
7. Tambah migration SQL manual
8. Daftarkan route di `internal/platform/http/router.go`
9. Tambah cache helper bila benar-benar perlu
10. Tambah test

### Kapan perlu membuat DTO
Buat DTO kalau:
- bentuk request berbeda dari entity
- response tidak ingin menampilkan semua field entity
- modul lain hanya butuh sebagian field

### Kapan perlu membuat mapper
Buat mapper kalau:
- entity dan response mulai berbeda
- ingin menjaga handler/service tetap rapi

Kalau entity dan response masih benar-benar sederhana, mapper kecil seperti di module product sudah cukup.

### Kapan jangan memecah file lagi
Jangan memecah file hanya demi terlihat rapi.

Contoh yang sengaja **tidak** saya pecah lebih jauh:
- service user tetap tipis
- response helper tetap satu file
- middleware wrapper tetap sederhana

Alasannya:
- kebutuhan belum kompleks
- refactor harus realistis
- readability lebih penting daripada jumlah file yang banyak

## Daftar file yang dipindah, diubah, dan ditambah

### Dipindah secara konsep
- `internal/auth/*` menjadi `internal/modules/auth/*`
- `internal/product/*` menjadi `internal/modules/product/*`

### Ditambah
- `internal/modules/user/*`
- `internal/platform/db/gorm.go`
- `internal/platform/http/server.go`
- `internal/platform/http/middleware/cors.go`
- `internal/platform/http/middleware/requestid.go`
- `internal/platform/http/middleware/recovery.go`
- `internal/platform/http/middleware/timeout.go`
- `internal/platform/logger/zap.go`
- `internal/shared/validator/validator.go`
- `internal/shared/apperror/mapper.go`
- `configs/env.go`
- `docs/swagger.yaml`
- `docs/LEARNING_GUIDE_V2.md`

### Diubah
- `cmd/api/main.go`
- `cmd/seeder/main.go`
- `cmd/migrate/main.go`
- `tests/integration/api_test.go`
- `internal/platform/http/router.go`
- `internal/platform/http/middleware/auth.go`
- `internal/platform/http/middleware/logging.go`
- `docs/swagger.html`

## Panduan refactor step by step dari struktur lama ke struktur baru

1. Buat folder `internal/modules`, `internal/platform`, dan `internal/shared` bila belum ada.
2. Pindahkan package bisnis ke `internal/modules`.
3. Tentukan pemilik entity utama untuk setiap tabel.
4. Pisahkan helper teknis ke `internal/platform`.
5. Pisahkan helper lintas modul ke `internal/shared`.
6. Ubah import path satu per satu.
7. Perbaiki wiring di `main.go`.
8. Perbarui seeder dan integration test.
9. Perbarui dokumentasi swagger.
10. Jalankan formatter dan test.

## File penting yang menjadi source of truth

Kalau ingin membaca hasil refactor langsung dari kode, mulai dari file ini:
- `cmd/api/main.go`
- `internal/platform/http/router.go`
- `internal/modules/auth/service.go`
- `internal/modules/product/service.go`
- `internal/modules/user/entity.go`
- `internal/shared/apperror/error.go`
- `tests/integration/api_test.go`

## Best practice untuk pemula

- simpan entity utama satu kali per tabel
- gunakan DTO saat kontrak API berbeda dari bentuk tabel
- buat mapper kecil saat response mulai berbeda dari entity
- pisahkan technical concern dari business concern
- jaga semua query tetap tenant-aware
- simpan hashing password di helper khusus
- jangan taruh query database langsung di handler
- pertahankan migration manual SQL
- jaga response error tetap konsisten

## Common mistakes untuk pemula

- memindahkan file hanya berdasarkan nama, bukan berdasarkan tanggung jawab
- menduplikasi entity `users` di auth dan user module sekaligus
- menaruh logic Redis detail di handler
- menaruh semua helper di `shared` walau sebenarnya hanya dipakai satu modul
- membuat abstraction terlalu cepat
- lupa update import path setelah pindah package
- lupa menjaga filter `tenant_id` di semua query
- lupa update test setelah refactor

## Ringkasan akhir

Struktur baru ini lebih siap berkembang karena:
- fitur bisnis terkumpul di `modules`
- komponen teknis terkumpul di `platform`
- helper lintas modul terkumpul di `shared`
- entity utama tidak diduplikasi
- wiring masih cukup sederhana untuk dipahami pemula

Ini bukan struktur paling “enterprise”, tetapi cukup realistis untuk project kecil-menengah yang mulai tumbuh.
