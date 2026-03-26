# Sesi 8 — Docker Compose dan Debug

## Tujuan sesi

Di sesi ini kamu akan belajar:

- bagaimana Docker Compose dipakai untuk menjalankan dependency
- apa fungsi service `mysql`, `redis`, dan `api`
- bagaimana environment variable dipasang di container
- bagaimana healthcheck membantu startup
- bagaimana melakukan debug lokal di VS Code

---

## File yang dipelajari

1. `docker-compose.yml`
2. `Dockerfile`
3. `.env.example`
4. `.vscode/launch.json`
5. `cmd/api/main.go`

---

## Bagian 1 — Docker Compose

### File: `docker-compose.yml`

Service yang ada:

- `mysql`
- `redis`
- `api`

---

## Service MySQL

Konfigurasi penting:

- image: `mysql:8.0`
- root password: `root`
- database: `app_db`
- port host: `3306:3306`
- volume: `mysql_data`

### Kenapa pakai volume
Supaya data database tidak hilang setiap container dimatikan.

### Healthcheck MySQL
Docker Compose juga menunggu MySQL sehat sebelum API mulai.

Ini membantu mengurangi masalah saat API start terlalu cepat sementara database belum siap.

---

## Service Redis

Konfigurasi penting:

- image: `redis:7-alpine`
- port host: `6379:6379`
- volume: `redis_data`

### Kenapa Redis juga diberi healthcheck
Supaya service lain tahu kapan Redis siap menerima koneksi.

---

## Service API

Service `api` dibuild dari project ini sendiri.

Environment yang dipasang:
- `APP_ENV`
- `HTTP_PORT`
- `MYSQL_DSN`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `JWT_SECRET`
- `JWT_ISSUER`
- `JWT_ACCESS_TOKEN_TTL`

### `depends_on`
API menunggu:
- MySQL sehat
- Redis sehat

---

## Bagian 2 — Cara menjalankan dengan Docker Compose

### Jalankan semua

```bash
docker compose up --build
```

### Jalankan hanya dependency

```bash
docker compose up -d mysql redis
```

Ini cocok kalau kamu ingin menjalankan API langsung lewat `go run ./cmd/api`.

---

## Urutan yang paling enak untuk pemula

### Opsi A — Jalankan semua dengan Docker Compose
Cocok jika kamu ingin cepat mencoba API.

Langkah:
1. `docker compose up --build`
2. jalankan migration
3. jalankan seeder jika belum otomatis kamu jalankan secara terpisah
4. coba endpoint

### Opsi B — Jalankan dependency dengan Docker, API dengan Go lokal
Cocok untuk belajar dan debug.

Langkah:
1. `docker compose up -d mysql redis`
2. jalankan migration
3. jalankan seeder
4. `go run ./cmd/api`

Opsi B biasanya lebih enak untuk belajar karena kamu bisa mengubah code dan menjalankan ulang API lebih mudah.

---

## Bagian 3 — Environment lokal

### File: `.env.example`

Untuk local development:
1. salin `.env.example` menjadi `.env`
2. sesuaikan bila perlu

Contoh:

```bash
cp .env.example .env
```

### Kenapa `.env` penting saat debug
Karena launch config VS Code akan membaca nilai dari file `.env`.

---

## Bagian 4 — Debug dengan VS Code

### File: `.vscode/launch.json`

Konfigurasi yang ada:
- name: `Debug API`
- program: `${workspaceFolder}/cmd/api`
- envFile: `${workspaceFolder}/.env`

### Artinya apa
Saat kamu menekan Run and Debug:
- VS Code menjalankan aplikasi dari `cmd/api`
- environment diambil dari `.env`

---

## Step by step debug untuk pemula

### 1. Install extension Go di VS Code
Pastikan extension Go terpasang.

### 2. Buka folder project
Buka root folder project, bukan hanya file tunggal.

### 3. Siapkan `.env`
Pastikan file `.env` ada.

### 4. Jalankan dependency
Kalau perlu:
```bash
docker compose up -d mysql redis
```

### 5. Jalankan migration dan seeder
Supaya aplikasi punya tabel dan data demo.

### 6. Buka menu Run and Debug
Pilih konfigurasi `Debug API`.

### 7. Pasang breakpoint
Contohnya di:
- `internal/auth/handler.go`
- `internal/auth/service.go`
- `internal/product/service.go`

### 8. Kirim request
Gunakan curl atau Postman.

### 9. Lihat alur program
Debugger akan berhenti di breakpoint dan kamu bisa:
- melihat isi request
- melihat tenant ID
- melihat request body
- melihat token claims
- melihat hasil query service

---

## Kapan debug paling berguna

Debug sangat membantu saat:
- login gagal dan kamu bingung kenapa
- token valid tapi request tetap ditolak
- create product berhasil tapi list tidak berubah
- cache terasa membingungkan
- tenant mismatch terjadi

---

## Bagian 5 — Hubungan Docker Compose dengan config

`docker-compose.yml` memasukkan environment ke container API.

Sementara saat kamu menjalankan API lokal, `configs/config.go` membaca nilai dari environment proses lokal.

Jadi sumber config sama-sama environment, hanya caranya berbeda:
- di container → dari `docker-compose.yml`
- di lokal → dari `.env` atau shell environment

---

## Latihan kecil sesi 8

### Latihan 1
Buka `docker-compose.yml`, lalu jawab:
- service apa saja yang ada?
- service mana yang memakai volume?
- service mana yang punya `depends_on`?

### Latihan 2
Buka `.vscode/launch.json`, lalu jawab:
- program yang dijalankan apa?
- env file yang dipakai apa?

### Latihan 3
Coba pasang breakpoint di `auth/service.go` method `Login`, lalu lakukan login.
Perhatikan:
- nilai email
- nilai tenantID
- kapan token dibuat

---

## Kesalahan umum pemula di sesi ini

### 1. Menjalankan API sebelum database siap
Gunakan healthcheck dan urutan setup yang benar.

### 2. Lupa migration
Container hidup tidak berarti tabel sudah ada.

### 3. Mengira `.env.example` otomatis dibaca
Yang biasa dipakai saat lokal adalah `.env`, bukan `.env.example`.

### 4. Debug tanpa dependency aktif
Kalau MySQL/Redis belum hidup, debug aplikasi akan gagal di startup.

---

## Checklist pemahaman sesi 8

Kamu harus bisa menjawab:

- fungsi masing-masing service di Docker Compose apa?
- kenapa healthcheck dipakai?
- bagaimana cara debug API di VS Code?
- kenapa `.env` penting untuk debug?
- kapan lebih enak menjalankan API lokal dibanding seluruhnya di container?

---

## Ringkasan sesi 8

Di sesi ini kamu belajar bahwa:

- Docker Compose dipakai untuk menyalakan MySQL, Redis, dan API
- healthcheck membantu urutan startup
- `.env` adalah sumber config lokal yang penting
- VS Code debug menjalankan `cmd/api`
- menjalankan dependency di Docker dan API di lokal sering menjadi cara belajar yang nyaman

Sesi berikutnya akan membahas:
**integration test**.
