# Sesi 3 — Auth JWT

## Tujuan sesi

Di sesi ini kamu akan belajar:

- bagaimana login bekerja
- bagaimana JWT dibuat
- bagaimana token divalidasi
- bagaimana user hasil auth disimpan ke context
- bagaimana endpoint `/auth/me` bekerja
- kenapa auth logic tidak ditaruh semua di handler

Ini salah satu sesi paling penting karena hampir semua backend modern punya kebutuhan auth.

---

## File yang dipelajari

Baca file ini berurutan:

1. `internal/auth/model.go`
2. `internal/auth/repository.go`
3. `internal/auth/jwt.go`
4. `internal/auth/service.go`
5. `internal/auth/context.go`
6. `internal/auth/handler.go`
7. `internal/platform/http/middleware/auth.go`
8. `internal/platform/http/router.go`

---

## Gambaran alur auth

Mari ikuti alur login dari awal sampai akhir.

### Langkah 1 — client kirim request login

Endpoint:
`POST /api/v1/auth/login`

Header:
- `X-Tenant-Id: tenant-demo`

Body:
```json
{
  "email": "demo@example.com",
  "password": "password123"
}
```

### Langkah 2 — middleware tenant berjalan
Sebelum handler login berjalan, middleware tenant akan membaca `X-Tenant-Id`.

Tenant ini lalu disimpan ke context request.

### Langkah 3 — handler login membaca body
File: `internal/auth/handler.go`

Handler:
- mengambil tenant dari context
- decode JSON body ke `LoginRequest`
- menjalankan validasi
- memanggil `service.Login(...)`

### Langkah 4 — service login menjalankan aturan bisnis
File: `internal/auth/service.go`

Service akan:
- mencari user berdasarkan `tenant_id + email`
- mengecek apakah user aktif
- membandingkan password plain text dengan hash bcrypt
- membuat JWT jika valid

### Langkah 5 — JWT dikembalikan ke client
Service mengembalikan `LoginResponse` yang berisi:
- `access_token`
- `token_type`
- `expires_in`
- data user

### Langkah 6 — client memakai token
Untuk endpoint yang butuh login, client mengirim:

```http
Authorization: Bearer <TOKEN>
X-Tenant-Id: tenant-demo
```

### Langkah 7 — middleware auth memvalidasi token
File: `internal/platform/http/middleware/auth.go`

Middleware:
- membaca `Authorization`
- mengambil token Bearer
- memanggil `jwtManager.ParseAccessToken(...)`
- membandingkan tenant token dengan tenant di header
- menyimpan user auth ke context

### Langkah 8 — handler protected tinggal membaca context
Contoh: endpoint `/auth/me`

Handler tidak perlu parse token lagi.
Ia cukup mengambil user dari context.

---

## Bagian 1 — Model auth

### File: `internal/auth/model.go`

Struct penting di sini:

- `User`
- `LoginRequest`
- `UserResponse`
- `LoginResponse`
- `Claims`
- `AuthUser`

### `User`
Ini model database untuk tabel `users`.

### `LoginRequest`
Ini DTO request login.

Validasinya:
- email wajib
- password wajib minimal 8 karakter

### `Claims`
Ini isi JWT.

Field penting:
- `UserID`
- `TenantID`
- `Email`
- `Role`

Ditambah `jwt.RegisteredClaims` untuk claim standar seperti:
- `sub`
- `iss`
- `iat`
- `exp`

### `AuthUser`
Ini bentuk user yang disimpan di context setelah token divalidasi.

---

## Bagian 2 — Repository auth

### File: `internal/auth/repository.go`

Repository auth berbicara ke MySQL lewat GORM.

Method penting:
- `FindByEmail`
- `FindByID`
- `Create`

### Kenapa query selalu tenant-aware
Lihat method `FindByEmail`:

- query berdasarkan `tenant_id`
- dan `email`

Ini penting agar user dari tenant lain tidak ikut terbaca.

---

## Bagian 3 — JWT manager

### File: `internal/auth/jwt.go`

Di file ini ada `JWTManager`.

Method penting:

- `GenerateAccessToken(user User)`
- `ParseAccessToken(tokenString string)`

### Saat membuat token
JWT manager:
- mengambil data user
- membuat `Claims`
- memberi waktu `IssuedAt`
- memberi waktu `ExpiresAt`
- menandatangani token dengan `JWT_SECRET`

### Saat membaca token
JWT manager:
- parse token
- memastikan signing method adalah `HS256`
- memvalidasi signature
- memvalidasi expiry
- mengembalikan `Claims`

### Kenapa cek algoritma penting
Supaya aplikasi tidak sembarangan menerima token dengan metode penandatanganan yang tidak diharapkan.

---

## Bagian 4 — Service auth

### File: `internal/auth/service.go`

Method penting:
- `Login`
- `Me`

### Alur `Login`
1. cari user berdasarkan tenant + email
2. kalau user tidak ditemukan → unauthorized
3. kalau user tidak aktif → unauthorized
4. bandingkan password input dengan hash di DB
5. kalau cocok, generate JWT
6. kembalikan response login

### Kenapa pakai `bcrypt.CompareHashAndPassword`
Karena database menyimpan hash, bukan password asli.
Jadi login harus membandingkan password input dengan hash tersebut.

### Alur `Me`
Service mengambil data user berdasarkan:
- tenant dari token
- user ID dari token

---

## Bagian 5 — Auth context

### File: `internal/auth/context.go`

File ini membuat helper:

- `NewContext(ctx, user)`
- `FromContext(ctx)`

Tujuannya:
- middleware auth menyimpan user ke context
- handler lain bisa membaca user dari context

### Kenapa ini penting
Supaya handler tidak perlu tahu detail JWT parsing.

---

## Bagian 6 — Handler auth

### File: `internal/auth/handler.go`

#### Method `Login`
Tugasnya:
- ambil tenant dari context
- decode JSON body
- validasi request
- panggil service
- kalau error, kirim response error
- kalau sukses, kirim response sukses

#### Method `Me`
Tugasnya:
- ambil auth user dari context
- panggil service
- kirim response

### Kenapa handler tidak membuat JWT sendiri
Karena handler seharusnya fokus pada HTTP saja.
Logika membuat token adalah logika auth/business, jadi lebih tepat ada di service/jwt manager.

---

## Bagian 7 — Middleware auth

### File: `internal/platform/http/middleware/auth.go`

Ini middleware yang menjaga endpoint private.

Alurnya:
1. ambil header `Authorization`
2. pecah format `Bearer <token>`
3. parse token
4. ambil tenant dari context
5. cek tenant di token sama dengan tenant di header
6. simpan `AuthUser` ke context

### Kenapa tenant token dicek lagi
Karena meskipun token valid, user tidak boleh memakai token tenant A untuk mengakses tenant B.

Itu sebabnya ada validasi:
- tenant di header
- harus sama dengan tenant di token

Ini bagian yang sangat penting untuk multi-tenant security.

---

## Bagian 8 — Router auth

### File: `internal/platform/http/router.go`

Bagian auth di router:

- `POST /api/v1/auth/login` → public
- `GET /api/v1/auth/me` → private

Untuk route private, router membungkus group dengan:
- `httpmw.AuthJWT(jwtManager)`

---

## Step by step mencoba login

### 1. Pastikan project sudah jalan
Kamu sudah harus:
- punya MySQL dan Redis aktif
- sudah menjalankan migration
- sudah menjalankan seeder
- sudah menjalankan API

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: tenant-demo" \
  -d '{
    "email": "demo@example.com",
    "password": "password123"
  }'
```

### 3. Ambil token dari response
Simpan `access_token`.

### 4. Pakai token untuk `/auth/me`

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Tenant-Id: tenant-demo"
```

### 5. Coba tenant mismatch
Ubah header jadi `X-Tenant-Id: tenant-lain` dengan token lama.

Harus gagal karena middleware auth memeriksa kecocokan tenant.

---

## Latihan kecil sesi 3

### Latihan 1
Buka `internal/auth/jwt.go`, lalu jawab:
- di mana token ditandatangani?
- di mana `exp` di-set?
- di mana `iss` di-set?

### Latihan 2
Buka `internal/platform/http/middleware/auth.go`, lalu jawab:
- token diambil dari header mana?
- format token yang diharapkan apa?
- apa yang terjadi jika tenant token tidak sama dengan tenant header?

### Latihan 3
Buka `internal/auth/service.go`, lalu jawab:
- apa yang terjadi jika email tidak ditemukan?
- apa yang terjadi jika password salah?
- apa yang terjadi jika user tidak aktif?

---

## Kesalahan umum pemula pada JWT

### 1. Menyimpan password asli
Jangan. Simpan hash.

### 2. Mengira JWT otomatis aman
JWT aman **jika**:
- secret dijaga
- token ada expiry
- signature divalidasi
- tenant/role diperiksa dengan benar

### 3. Parse JWT di semua handler
Itu membuat code berulang. Gunakan middleware.

### 4. Menganggap token valid cukup
Di project multi-tenant, token valid saja belum cukup.
Tenant pada token juga harus cocok dengan tenant request.

---

## Checklist pemahaman sesi 3

Kamu harus bisa menjawab:

- login memanggil service mana?
- JWT dibuat di file mana?
- token dibaca di file mana?
- auth user disimpan ke context di file mana?
- kenapa tenant mismatch harus ditolak?

---

## Ringkasan sesi 3

Di sesi ini kamu belajar bahwa:

- login dimulai dari handler, diproses di service, dan token dibuat oleh JWT manager
- password dibandingkan menggunakan bcrypt
- middleware auth membaca dan memvalidasi Bearer token
- tenant pada token harus sama dengan tenant pada request
- user hasil auth disimpan ke context untuk dipakai handler protected

Setelah auth selesai, sesi berikutnya akan fokus ke:
**CRUD products**.
