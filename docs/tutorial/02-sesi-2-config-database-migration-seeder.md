# Sesi 2 — Config, Database, Migration, Seeder

## Tujuan sesi

Di sesi ini kamu akan belajar:

- bagaimana project membaca konfigurasi
- bagaimana koneksi MySQL dibuat
- kenapa project memakai GORM tapi migration tetap SQL manual
- bagaimana seeder bekerja
- langkah menjalankan database dari nol

Sesi ini penting karena sebelum memahami API, kamu harus tahu dulu dari mana aplikasi mendapat konfigurasi dan bagaimana database disiapkan.

---

## File yang dipelajari di sesi ini

Baca file berikut:

1. `configs/config.go`
2. `internal/platform/db/mysql.go`
3. `migrations/000001_create_users.up.sql`
4. `migrations/000002_create_products.up.sql`
5. `cmd/seeder/main.go`
6. `.env.example`
7. `scripts/migrate-up.sh`
8. `scripts/migrate-down.sh`

---

## Bagian 1 — Config

### File: `configs/config.go`

Tugas file ini:
- membaca environment variable
- memberi nilai default
- mengubah nilai string menjadi tipe yang dibutuhkan, misalnya `int` atau `time.Duration`

Struct config yang dipakai:

- `AppEnv`
- `HTTPPort`
- `MySQLDSN`
- `RedisAddr`
- `RedisPassword`
- `RedisDB`
- `JWTSecret`
- `JWTIssuer`
- `JWTAccessTokenTTL`

### Kenapa config dipisah dari `main.go`

Kalau semua `os.Getenv()` ditulis langsung di `main.go`, file `main.go` cepat menjadi panjang dan berisik.

Dengan memindahkan config ke file sendiri:
- code lebih rapi
- lebih mudah dicari
- lebih mudah diganti
- lebih gampang dipakai ulang di seeder atau test

---

## Bagian 2 — File `.env.example`

File ini adalah contoh environment variable.

Isinya antara lain:

- `APP_ENV=local`
- `HTTP_PORT=8080`
- `MYSQL_DSN=...`
- `REDIS_ADDR=localhost:6379`
- `JWT_SECRET=super-secret-change-me`

### Kenapa `.env.example` penting

Sebagai pemula, file ini membantu kamu tahu:
- variabel apa saja yang wajib ada
- nilai default lokal yang bisa dipakai
- format DSN yang benar

---

## Bagian 3 — Koneksi MySQL dengan GORM

### File: `internal/platform/db/mysql.go`

Di sini project membuka koneksi MySQL dengan GORM:

- `gorm.Open(mysql.Open(dsn), &gorm.Config{...})`

Setelah itu, project mengambil objek `sql.DB` di balik GORM untuk mengatur connection pool:

- `SetMaxOpenConns`
- `SetMaxIdleConns`
- `SetConnMaxLifetime`

### Kenapa ini penting

Koneksi database tidak boleh dibuka sembarangan per request.
Aplikasi biasanya membuka satu pool koneksi di awal, lalu memakainya ulang.

### Kenapa tetap pakai GORM

Karena GORM cukup ramah untuk belajar:
- query lebih mudah ditulis
- CRUD lebih cepat dibuat
- model Go bisa langsung dipakai

### Kenapa migration tetap SQL manual

Walau project memakai GORM, schema database **tidak** dibuat dengan `AutoMigrate`.

Alasannya:
- kamu belajar bentuk tabel dengan jelas
- kamu belajar index dan constraint dengan sadar
- perubahan schema lebih mudah direview
- lebih dekat ke praktik tim backend yang serius

---

## Bagian 4 — Migration SQL

### File: `migrations/000001_create_users.up.sql`

Migration ini membuat tabel `users`.

Kolom penting:
- `id`
- `tenant_id`
- `name`
- `email`
- `password_hash`
- `role`
- `is_active`
- `created_at`
- `updated_at`
- `deleted_at`

Index penting:
- unique `(tenant_id, email)`
- index `tenant_id`
- index `deleted_at`

### Kenapa email unik per tenant

Karena project ini multi-tenant.
Jadi tenant A dan tenant B secara konsep bisa saja memiliki email yang sama.

---

### File: `migrations/000002_create_products.up.sql`

Migration ini membuat tabel `products`.

Kolom penting:
- `id`
- `tenant_id`
- `name`
- `description`
- `price`
- `stock`
- `created_at`
- `updated_at`
- `deleted_at`

Index penting:
- `tenant_id`
- `deleted_at`
- `(tenant_id, name)`

### Kenapa `tenant_id` selalu ada

Karena semua data product harus terikat ke tenant tertentu.
Kalau tidak, data bisa tercampur antar tenant.

---

## Bagian 5 — Seeder

### File: `cmd/seeder/main.go`

Seeder bertugas mengisi data awal untuk belajar.

Yang di-seed:
- 1 user demo
- beberapa product demo

### User demo
- Tenant: `tenant-demo`
- Email: `demo@example.com`
- Password: `password123`

Password **tidak** disimpan dalam bentuk plain text.
Seeder memakai `bcrypt.GenerateFromPassword(...)` untuk membuat hash.

### Kenapa seeder dipisah dari migration

Karena migration dan seed punya tujuan berbeda.

#### Migration
Fokus ke **schema**
- buat tabel
- buat index
- ubah struktur database

#### Seeder
Fokus ke **data awal**
- user demo
- product demo

Dengan memisahkan keduanya, kamu bisa:
- menjalankan migration tanpa memaksa data demo masuk
- menjalankan seed berulang kali di local
- mengelola data awal lebih fleksibel

---

## Step by step menjalankan database dari nol

### Langkah 1 — salin env

```bash
cp .env.example .env
```

### Langkah 2 — jalankan MySQL dan Redis

```bash
docker compose up -d mysql redis
```

Tunggu sampai container sehat.

### Langkah 3 — jalankan migration

Project ini menyediakan script:

- `scripts/migrate-up.sh`
- `scripts/migrate-down.sh`

Kalau kamu sudah punya `migrate` CLI, jalankan sesuai isi script tersebut.

Contoh bentuk perintah:

```bash
migrate -path ./migrations -database "mysql://root:root@tcp(localhost:3306)/app_db?multiStatements=true" up
```

### Langkah 4 — jalankan seeder

```bash
go run ./cmd/seeder
```

### Langkah 5 — cek hasilnya

Masuk ke MySQL lalu cek tabel:

- `users`
- `products`

Kamu juga bisa cek isi awalnya.

---

## Cara berpikir pemula: kenapa urutan setup seperti ini

Urutannya bukan kebetulan.

1. **env dulu**  
   karena aplikasi perlu tahu koneksi ke mana

2. **container database dulu**  
   karena migration butuh MySQL aktif

3. **migration dulu**  
   karena tabel harus ada sebelum data diisi

4. **seeder setelah migration**  
   karena seed perlu tabel yang sudah jadi

5. **jalankan API terakhir**  
   karena API bergantung pada semuanya

---

## Latihan kecil sesi 2

Lakukan latihan ini:

1. buka `migrations/000001_create_users.up.sql`
2. sebutkan:
   - kolom mana yang menyimpan hash password?
   - kolom mana yang menyimpan tenant?
   - index mana yang mencegah email ganda per tenant?

Lalu:

3. buka `cmd/seeder/main.go`
4. cari bagian yang membuat hash password
5. cari bagian yang memakai `FirstOrCreate`

### Kenapa `FirstOrCreate` dipakai
Supaya seeder aman dijalankan berulang kali dan tidak selalu membuat duplikasi.

---

## Kesalahan umum pemula di sesi ini

### 1. Mengira GORM otomatis membuat tabel
Di project ini **tidak**. Tabel dibuat lewat migration SQL.

### 2. Menjalankan seeder sebelum migration
Ini akan gagal karena tabel belum ada.

### 3. Lupa `multiStatements=true` pada DSN saat migration
Beberapa tool migration membutuhkannya untuk skenario tertentu.

### 4. Menyimpan password plain text
Ini tidak aman. Gunakan hash.

---

## Checklist pemahaman sesi 2

Kamu harus bisa menjawab:

- config di-load dari file mana?
- koneksi MySQL dibuat di file mana?
- kenapa pakai migration SQL manual?
- kenapa seeder dipisah dari migration?
- data demo login disimpan di mana?

---

## Ringkasan sesi 2

Di sesi ini kamu belajar bahwa:

- config dibaca dari environment dan dikumpulkan dalam satu struct
- GORM dipakai untuk query, bukan untuk migrasi schema
- tabel `users` dan `products` dibuat lewat SQL manual
- seeder mengisi user demo dan products demo
- setup project harus dilakukan dengan urutan yang benar

Setelah pondasi siap, sesi berikutnya akan masuk ke:
**JWT authentication**.
