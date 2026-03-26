# Go Starter 1.25 — Dev Starter Guide

Dokumentasi ini menjelaskan cara menjalankan project Go starter untuk mode development menggunakan:

- Dockge
- Docker / Colima
- image dev `local/dev-go-air:1.25`
- Air untuk auto reload
- VS Code untuk coding
- `.env` project untuk konfigurasi aplikasi
- `.env` Dockge untuk konfigurasi global stack

Dokumentasi ini juga menjelaskan workflow harian development dan apa yang perlu dilakukan saat menambah package/library baru di Go.

---

## Tujuan Setup

Setup ini dibuat supaya development terasa:

- stabil
- cepat
- mudah dipahami
- nyaman dipakai harian
- tidak bercampur antara konfigurasi aplikasi dan konfigurasi stack Dockge

---

## Konsep Inti

Ada 2 jenis `.env` yang dipakai.

### 1. `.env` Dockge

Dipakai hanya untuk konfigurasi global stack, misalnya:

- `STACK_SOURCE_PATH`
- `STACK_SLUG`
- `PUBLISH_PORT`
- `PLATFORM`
- `CONTAINER_NAME`
- `DEV_IMAGE`

### 2. `.env` Project

Dipakai untuk konfigurasi aplikasi, misalnya:

- `APP_ENV`
- `HTTP_PORT`
- `MYSQL_DSN`
- `REDIS_ADDR`
- `JWT_SECRET`
- `JWT_ISSUER`
- `JWT_ACCESS_TOKEN_TTL`
- `GORM_DEBUG`
- `LOG_LEVEL`

Dengan pola ini:

- `compose.yaml` tetap rapi
- konfigurasi aplikasi tetap tinggal di folder project
- konfigurasi Dockge tidak tercampur dengan config app

---

## Struktur Project Minimal

```text
go-starter-1.25
├── .air.toml
├── .env
├── .env.example
├── go.mod
├── go.sum
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   └── server/
├── Dockerfile
├── Dockerfile.dev-go-air
└── README.md
```

Struktur minimal agar mode dev bisa jalan:

```text
go-starter-1.25
├── .air.toml
├── .env
├── go.mod
└── cmd/
    └── api/
        └── main.go
```

---

## Struktur Dockge

Source code asli bebas di mana saja, misalnya:

```text
/Users/endahfathonah/dev/projects/Research/go-starter-1.25
```

Lalu buat symlink ke folder stack Dockge:

```text
/Users/endahfathonah/dev/dockge/data/stacks/go-starter-1.25
  -> /Users/endahfathonah/dev/projects/Research/go-starter-1.25
```

Penting:

- Dockge mengelola stack dari folder stack miliknya
- source code tetap kamu edit dari folder asli
- symlink dipakai agar Dockge tetap bisa mengelola project itu

---

## File Penting

### `.air.toml`

Mengatur auto reload untuk Go.

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/api ./cmd/api"
  bin = "./tmp/api"
  full_bin = "./tmp/api"
  include_ext = ["go", "tpl", "tmpl", "html", "yaml", "yml", "json"]
  exclude_dir = ["tmp", "vendor", ".git", ".vscode"]
  delay = 300
  stop_on_error = true
  send_interrupt = true
  kill_delay = "500ms"

[log]
  time = true

[misc]
  clean_on_exit = true
```

### `.env` project

Contoh:

```env
APP_ENV=development
HTTP_PORT=8080

MYSQL_DSN=root:root@tcp(host.docker.internal:3306)/app_db?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=host.docker.internal:6379

JWT_SECRET=super-secret-change-me
JWT_ISSUER=go-starter
JWT_ACCESS_TOKEN_TTL=15m

GORM_DEBUG=true
LOG_LEVEL=debug
```

### `.env.example`

```env
APP_ENV=development
HTTP_PORT=8080
MYSQL_DSN=
REDIS_ADDR=
JWT_SECRET=
JWT_ISSUER=go-starter
JWT_ACCESS_TOKEN_TTL=15m
GORM_DEBUG=true
LOG_LEVEL=debug
```

### `Dockerfile.dev-go-air`

Image dev reusable untuk semua project Go.

```dockerfile
FROM golang:1.25-alpine

RUN apk add --no-cache git bash ca-certificates && update-ca-certificates

WORKDIR /workspace

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", "/workspace/.air.toml"]
```

---

## Compose Dev untuk Dockge

### `compose.yaml`

```yaml
services:
  api:
    image: ${DEV_IMAGE}
    pull_policy: never
    container_name: ${CONTAINER_NAME}
    working_dir: /workspace
    platform: ${PLATFORM}

    ports:
      - "${PUBLISH_PORT}:8080"

    env_file:
      - .env

    volumes:
      - ${STACK_SOURCE_PATH}:/workspace
      - gomodcache:/go/pkg/mod
      - gobuildcache:/root/.cache/go-build

    extra_hosts:
      - "host.docker.internal:host-gateway"

    command: ["air", "-c", "/workspace/.air.toml"]

    stdin_open: true
    tty: true

volumes:
  gomodcache:
    name: ${STACK_SLUG}_gomodcache
  gobuildcache:
    name: ${STACK_SLUG}_gobuildcache
```

### `.env` Dockge

```env
STACK_SLUG=go-starter
CONTAINER_NAME=go_starter_api
DEV_IMAGE=local/dev-go-air:1.25
PLATFORM=linux/arm64
PUBLISH_PORT=8080
STACK_SOURCE_PATH=/Users/endahfathonah/dev/dockge/data/stacks/go-starter-1.25
```

---

## Kenapa `env_file: .env`?

Karena `.env` aplikasi berada di root project yang sama dengan `compose.yaml`.

Jadi:

- `.env` Dockge dipakai untuk interpolasi `${...}`
- `.env` project dipakai sebagai `env_file` untuk aplikasi

Ini membuat:

- env global tetap di Dockge
- env aplikasi tetap di project

---

## Setup Awal Project

### 1. Buat project

```bash
mkdir -p /Users/endahfathonah/dev/projects/Research/go-starter-1.25/cmd/api
cd /Users/endahfathonah/dev/projects/Research/go-starter-1.25
```

### 2. Inisialisasi Go module

```bash
go mod init go-starter-1.25
```

### 3. Buat file minimal

Buat:

- `.air.toml`
- `.env`
- `cmd/api/main.go`

Contoh `cmd/api/main.go`:

```go
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

### 4. Build image dev

```bash
docker build -t local/dev-go-air:1.25 -f Dockerfile.dev-go-air .
```

### 5. Buat symlink Dockge

```bash
rm -rf /Users/endahfathonah/dev/dockge/data/stacks/go-starter-1.25

ln -s /Users/endahfathonah/dev/projects/Research/go-starter-1.25 \
/Users/endahfathonah/dev/dockge/data/stacks/go-starter-1.25
```

### 6. Buat stack di Dockge

- buka Dockge
- klik **Compose**
- isi nama stack: `go-starter-1.25`
- tempel `compose.yaml`
- isi `.env` Dockge
- save
- start

### 7. Cek aplikasi

```bash
curl http://localhost:8080/healthz
```

Kalau keluar `ok`, berarti setup berhasil.

---

## Daily Development Workflow

### Mulai kerja

1. Jalankan Colima.

```bash
colima start
```

2. Buka Dockge lalu klik **Start** pada stack project.

3. Buka project di VS Code dari source asli:

```text
/Users/endahfathonah/dev/projects/Research/go-starter-1.25
```

Jangan edit file dari folder Dockge stack. Edit selalu dari source asli.

### Saat coding

Kamu cukup edit file Go seperti biasa di VS Code.

Kalau ada perubahan code:

- Air mendeteksi perubahan
- app dibuild ulang
- process di-restart otomatis di container

Jadi kamu tidak perlu restart manual untuk perubahan `.go`.

### Melihat log

Gunakan Dockge UI:

- buka stack
- lihat tab **Logs**

Di situ kamu bisa lihat:

- error build
- panic
- output `log.Println`
- hasil reload Air

### Masuk ke container

Kalau perlu cek isi container:

- buka stack di Dockge
- klik **Terminal / Bash**

Contoh perintah yang sering dipakai:

```bash
pwd
ls -la
go env
cat .env
```

### Selesai kerja

1. Stop stack di Dockge.
2. Matikan Colima bila perlu.

```bash
colima stop
```

---

## Kapan Harus Restart?

### Tidak perlu restart

Kalau yang berubah:

- file `.go`
- logic handler
- service
- repository
- route
- middleware Go

Karena Air akan auto reload.

### Perlu restart container

Kalau yang berubah:

- `.air.toml`
- `.env`
- `compose.yaml`
- binary/image base
- dependency sistem di image

### Perlu build ulang image

Kalau yang berubah:

- `Dockerfile.dev-go-air`
- tool yang di-install di image
- versi Go dev image

---

## Workflow Update Harian

### Update code

Tinggal save file di VS Code. Air reload otomatis.

### Update `.env`

Restart stack dari Dockge.

### Update `compose.yaml`

Edit di Dockge lalu klik save atau redeploy.

### Update image dev

Build ulang:

```bash
docker build -t local/dev-go-air:1.25 -f Dockerfile.dev-go-air .
```

Lalu restart stack di Dockge.

---

## Menambah Package / Library Baru di Golang

Bagian ini penting supaya workflow development tetap rapi saat project mulai bertambah dependency.

### Kapan perlu menambah package?

Contoh saat kamu ingin menambah:

- router seperti `chi`
- ORM seperti `gorm`
- JWT library
- Redis client
- validator
- logger
- testing helper

### Cara menambah package Go biasa

Kalau package murni Go, jalankan dari root project:

```bash
go get github.com/go-chi/chi/v5
go mod tidy
```

Atau misalnya:

```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql
go mod tidy
```

Setelah itu:

- `go.mod` akan berubah
- `go.sum` akan berubah
- Air biasanya akan rebuild otomatis saat code berubah

Kalau container dev sedang jalan, bind mount akan membuat perubahan `go.mod` dan `go.sum` ikut terlihat di container.

### Kalau package baru butuh import di code

Misalnya kamu baru menambah package lalu menggunakannya di `main.go` atau file lain.

Yang terjadi:

- save file
- Air rebuild
- dependency diunduh jika belum ada
- app dijalankan ulang

Kalau download module pertama kali agak lama, itu normal.

### Kalau package butuh tool CLI tambahan

Contoh tool seperti:

- `air`
- migration CLI
- code generator
- linter tertentu

Kalau tool itu hanya dipakai di host, cukup install lokal.

Kalau tool itu harus tersedia di container dev, maka kamu perlu:

1. update `Dockerfile.dev-go-air`
2. build ulang image
3. restart stack di Dockge

Contoh menambah tool:

```dockerfile
RUN go install github.com/swaggo/swag/cmd/swag@latest
```

Lalu build ulang:

```bash
docker build -t local/dev-go-air:1.25 -f Dockerfile.dev-go-air .
```

### Kalau package butuh dependency OS / system library

Ini kasus yang sering dilupakan.

Contoh:

- package butuh `git`
- package butuh `gcc`
- package butuh header C
- package butuh `make`
- package butuh `sqlite-dev` atau library native lain

Kalau begitu, update juga `Dockerfile.dev-go-air`.

Contoh:

```dockerfile
RUN apk add --no-cache git bash ca-certificates build-base
```

Lalu build ulang image dan restart stack.

### Kalau package menambah env baru

Misalnya kamu menambah package S3 client dan butuh:

- `S3_ENDPOINT`
- `S3_ACCESS_KEY`
- `S3_SECRET_KEY`

Maka lakukan ini:

1. tambahkan env ke `.env`
2. tambahkan juga ke `.env.example`
3. dokumentasikan di README ini bila penting
4. restart stack agar env baru terbaca

### Checklist saat menambah package baru

Saat menambah package/library baru, cek hal ini:

- apakah cukup `go get` saja?
- apakah perlu `go mod tidy`?
- apakah perlu update code?
- apakah perlu update `.env`?
- apakah perlu update `Dockerfile.dev-go-air`?
- apakah perlu build ulang image?
- apakah perlu restart stack?

### Rule of thumb

Pakai patokan sederhana ini.

**Hanya package Go murni:**
- `go get`
- `go mod tidy`
- save code
- biasanya tidak perlu build ulang image

**Butuh tool baru di container:**
- update `Dockerfile.dev-go-air`
- build ulang image
- restart stack

**Butuh native/system dependency:**
- update `Dockerfile.dev-go-air`
- build ulang image
- restart stack

**Butuh env baru:**
- update `.env`
- update `.env.example`
- restart stack

---

## Checklist Cepat Saat Error

Jalankan ini di host:

```bash
ls -l /Users/endahfathonah/dev/dockge/data/stacks/go-starter-1.25
ls -la /Users/endahfathonah/dev/projects/Research/go-starter-1.25/.env
ls -la /Users/endahfathonah/dev/projects/Research/go-starter-1.25/.air.toml
docker ps
```

Yang harus benar:

- symlink ada
- `.env` ada
- `.air.toml` ada
- container running

---

## Error yang Paling Sering Terjadi

### 1. `.env not found`

Penyebab:

- file `.env` project belum ada
- `env_file` salah path
- salah paham antara `.env` Dockge dan `.env` aplikasi

Solusi:

- pastikan `.env` aplikasi ada di project root
- untuk setup ini gunakan `env_file: .env`

### 2. Air gagal build

Penyebab:

- `.air.toml` salah command
- `cmd/api` tidak ada
- `main.go` belum ada
- project belum punya `go.mod`

Solusi:

- cek `.air.toml`
- cek `cmd/api/main.go`
- jalankan `go mod init`

### 3. App tidak bisa connect MySQL / Redis

Penyebab:

- host salah
- service database belum jalan
- `MYSQL_DSN` atau `REDIS_ADDR` salah

Solusi:

- kalau DB di host machine, pakai `host.docker.internal`
- cek port DB
- cek credential

### 4. Perubahan code tidak reload

Penyebab:

- source tidak termount dengan benar
- `.air.toml` salah
- stack mount bukan ke folder project yang benar

Solusi:

- cek `STACK_SOURCE_PATH`
- cek symlink
- cek log Air

---

## Best Practice Harian

- edit source dari folder project asli
- start dan stop stack lewat Dockge
- jangan jalankan stack Dockge dari CLI
- pakai path absolut untuk source path
- pakai `.env` project untuk env aplikasi
- pakai `.env` Dockge hanya untuk global stack
- gunakan `pull_policy: never` untuk image lokal
- jangan mount `.:/workspace`
- setiap ada env baru, update juga `.env.example`
- setiap ada dependency baru, pikirkan apakah itu butuh update image atau cukup `go get`

---

## Larangan yang Sebaiknya Dihindari

Jangan lakukan ini:

```yaml
volumes:
  - .:/workspace
```

Jangan lakukan ini:

```yaml
command: ["air", "-c", ".air.toml"]
```

Jangan simpan env aplikasi di `.env` Dockge bila env itu milik app.

Jangan start dan stop stack Dockge campur antara UI dan CLI.

---

## Tips VS Code

Buka folder source asli, bukan folder stack Dockge.

Contoh:

```text
/Users/endahfathonah/dev/projects/Research/go-starter-1.25
```

Keuntungannya:

- indexing lebih normal
- extension Go lebih stabil
- file `.env`, `.air.toml`, dan source terlihat utuh
- tidak bingung dengan symlink stack

---

## Ringkasan Workflow Harian

### Pagi

- `colima start`
- Start stack di Dockge
- buka project di VS Code

### Saat kerja

- edit code
- Air auto reload
- lihat log di Dockge
- shell via Dockge bila perlu

### Saat ada perubahan config

- restart stack

### Saat update image

- build ulang image
- restart stack

### Sore / selesai kerja

- stop stack di Dockge
- `colima stop`

---

## Ringkasan Singkat

Pola development ini memisahkan:

- source code
- stack manager
- env aplikasi
- env global Dockge

Dengan pola ini:

- Dockge tetap rapi
- project tetap nyaman dikerjakan di VS Code
- Air auto reload berjalan
- config aplikasi tidak tercampur dengan config stack
- workflow harian jadi konsisten dan minim error
