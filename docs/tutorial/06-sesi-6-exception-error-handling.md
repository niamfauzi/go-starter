# Sesi 6 — Exception / Error Handling

## Tujuan sesi

Di sesi ini kamu akan belajar:

- bagaimana project membungkus error agar konsisten
- bagaimana error dipetakan menjadi HTTP status code
- bagaimana detail validation dikirim ke client
- bagaimana internal error tidak bocor ke client
- bagaimana error mengalir dari repository → service → handler → response

Ini sesi penting karena API yang rapi bukan hanya soal sukses, tetapi juga soal gagal dengan jelas dan aman.

---

## File yang dipelajari

1. `internal/shared/apperror/error.go`
2. `internal/platform/http/response.go`
3. `internal/auth/service.go`
4. `internal/product/service.go`
5. `internal/auth/handler.go`
6. `internal/product/handler.go`
7. `internal/platform/http/middleware/auth.go`
8. `internal/platform/http/middleware/tenant.go`

---

## Gambaran besar alur error

Mari lihat bentuk umumnya.

### Di repository
Repository biasanya mengembalikan error asli dari GORM atau database.

### Di service
Service mengubah error teknis menjadi error aplikasi yang lebih bermakna, misalnya:
- not found
- unauthorized
- internal server error

### Di handler
Handler tidak perlu tahu detail teknis database.
Handler cukup memanggil helper `WriteError`.

### Di response
Client mendapat format error yang konsisten.

---

## Bagian 1 — AppError

### File: `internal/shared/apperror/error.go`

Struct inti:

- `Code`
- `Message`
- `HTTPStatus`
- `Details`
- `Err`

### Fungsi field-field ini

#### `Code`
Kode mesin yang konsisten, misalnya:
- `VALIDATION_ERROR`
- `UNAUTHORIZED`
- `NOT_FOUND`

#### `Message`
Pesan singkat untuk client.

#### `HTTPStatus`
Status code HTTP yang dipakai.

#### `Details`
Dipakai untuk info tambahan, misalnya detail validation.

#### `Err`
Error asli internal. Ini berguna untuk logging, tapi tidak dikirim ke client.

---

## Constructor helper

Project menyediakan helper seperti:

- `Validation(...)`
- `BadRequest(...)`
- `Unauthorized(...)`
- `Forbidden(...)`
- `NotFound(...)`
- `Conflict(...)`
- `Internal(err)`

### Kenapa helper ini enak
Karena service/handler jadi lebih jelas saat membaca code.

Contoh:
- `return apperror.NotFound("product not found")`
- `return apperror.Unauthorized("invalid email or password")`

Ini jauh lebih mudah dibaca daripada merakit struct error panjang berulang kali.

---

## Bagian 2 — Helper response

### File: `internal/platform/http/response.go`

File ini punya dua helper utama:

- `WriteSuccess`
- `WriteError`

### Bentuk response sukses
Contoh:

```json
{
  "success": true,
  "message": "product fetched",
  "data": { ... }
}
```

### Bentuk response error
Contoh:

```json
{
  "success": false,
  "message": "validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [...]
  }
}
```

### Kenapa response dibuat seragam
Karena:
- client lebih mudah membaca
- frontend lebih mudah meng-handle
- debugging lebih mudah
- semua endpoint punya pola yang sama

---

## Bagian 3 — Validation error

Masih di file `response.go`, ada helper:

- `ValidationDetails(err)`

Fungsinya:
- mengambil error dari validator
- mengubahnya menjadi array field detail

Contohnya:
- field apa yang gagal
- tag validasi apa yang gagal
- nilai parameternya apa

Ini membuat pesan error lebih ramah untuk client.

---

## Bagian 4 — Error flow di auth

### File: `internal/auth/service.go`

Contoh pada `Login`:

- jika user tidak ditemukan → `Unauthorized("invalid email or password")`
- jika user tidak aktif → `Unauthorized("user is inactive")`
- jika password salah → `Unauthorized("invalid email or password")`
- jika generate token gagal → `Internal(err)`

### Kenapa “email tidak ditemukan” dan “password salah” memakai pesan yang sama
Ini best practice yang baik untuk login.
Tujuannya supaya attacker tidak mudah tahu apakah email tertentu benar-benar terdaftar.

---

## Bagian 5 — Error flow di product

### File: `internal/product/service.go`

Contoh pada `GetByID`:

- kalau GORM mengembalikan `ErrRecordNotFound` → `NotFound("product not found")`
- selain itu → `Internal(err)`

Contoh pada `Update` dan `Delete`:
- kalau data tidak ada → `NotFound`
- kalau query gagal → `Internal`

### Kenapa service yang melakukan mapping error
Karena service paling tahu arti bisnis dari error tersebut.

Contoh:
- error DB “record tidak ditemukan” di konteks bisnis berarti “product tidak ada”
- error DB “gagal konek” di konteks bisnis berarti “internal server error”

---

## Bagian 6 — Handler dan logging error

### File: `internal/auth/handler.go`
### File: `internal/product/handler.go`

Handler memakai:
- `platformhttp.WriteError(w, err)`

Handler tidak perlu membongkar semua error satu per satu, karena helper response sudah menanganinya.

### Kenapa beberapa error di-log dan beberapa tidak
Di handler product ada helper `logIfInternal(...)`.

Artinya:
- error 4xx biasanya tidak perlu dianggap error sistem
- error 5xx lebih penting untuk di-log

Ini membantu log tetap fokus ke masalah server yang perlu perhatian.

---

## Bagian 7 — Error dari middleware

Middleware juga bisa menghasilkan error.

Contoh:
- `tenant.go` → jika `X-Tenant-Id` hilang, kirim `BadRequest`
- `auth.go` → jika token hilang, kirim `Unauthorized`
- `auth.go` → jika tenant token mismatch, kirim `Forbidden`

### Kenapa middleware ikut memakai `WriteError`
Supaya format error tetap konsisten, walau error muncul sebelum handler dipanggil.

---

## Step by step memahami alur error nyata

Ambil contoh request ini:

```bash
curl http://localhost:8080/api/v1/products/9999 \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Tenant-Id: tenant-demo"
```

### Yang terjadi:
1. middleware tenant lolos
2. middleware auth lolos
3. handler `GetByID` dipanggil
4. service `GetByID` memanggil repository
5. repository tidak menemukan data
6. GORM mengembalikan `ErrRecordNotFound`
7. service mengubahnya menjadi `apperror.NotFound("product not found")`
8. handler memanggil `WriteError`
9. response 404 dikirim ke client

Ini contoh alur error yang bersih dan mudah dilacak.

---

## Latihan kecil sesi 6

### Latihan 1
Buka `internal/shared/apperror/error.go`, lalu jawab:
- field apa yang tidak dikirim ke client?
- field apa yang menyimpan status code?
- helper apa yang dipakai untuk internal server error?

### Latihan 2
Buka `internal/platform/http/response.go`, lalu jawab:
- fungsi apa yang mengirim response error?
- fungsi apa yang merapikan validation error?

### Latihan 3
Buka `internal/product/service.go`, lalu cari semua tempat yang mengembalikan:
- `NotFound`
- `Internal`

---

## Kenapa internal error tidak dibocorkan

Ini penting untuk keamanan dan kebersihan API.

Kalau error internal mentah dikirim ke client, bisa bocor:
- query database
- detail schema
- detail stack
- informasi yang tidak perlu dilihat user

Karena itu project ini:
- menyimpan error asli di field `Err`
- tetapi mengirim pesan aman seperti `internal server error`

---

## Kapan pakai BadRequest, Unauthorized, Forbidden, NotFound

### `BadRequest`
Dipakai saat request client salah bentuk.
Contoh:
- body JSON rusak
- product id bukan angka
- header tenant tidak ada

### `Unauthorized`
Dipakai saat client belum terautentikasi atau token salah.

### `Forbidden`
Dipakai saat client terautentikasi, tetapi tidak boleh melakukan aksi tertentu.
Contoh:
- token tenant mismatch

### `NotFound`
Dipakai saat resource tidak ditemukan.

### `Internal`
Dipakai saat ada kegagalan teknis di sisi server.

---

## Kesalahan umum pemula di error handling

### 1. Mengembalikan semua error sebagai 500
Ini membuat API sulit dipakai.

### 2. Mengirim error database mentah ke client
Ini tidak aman.

### 3. Menulis format error beda-beda di tiap handler
Ini membuat API tidak konsisten.

### 4. Menangani error database langsung di handler
Lebih baik service yang melakukan mapping makna bisnis.

---

## Checklist pemahaman sesi 6

Kamu harus bisa menjawab:

- custom app error ada di file mana?
- helper response error ada di file mana?
- siapa yang mengubah error GORM menjadi `NotFound`?
- kenapa error 5xx lebih penting untuk logging?
- kenapa error internal tidak dikirim mentah ke client?

---

## Ringkasan sesi 6

Di sesi ini kamu belajar bahwa:

- project memakai `AppError` agar error konsisten
- service bertugas memetakan error teknis menjadi error bisnis
- handler cukup memanggil `WriteError`
- response error selalu punya format yang sama
- detail validation dipisahkan agar mudah dibaca
- internal error tetap disimpan untuk log, tapi tidak dibocorkan ke client

Sesi berikutnya akan membahas dokumentasi API:
**Swagger**.
