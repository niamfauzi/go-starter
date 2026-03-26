package main

import (
	"fmt"
	"os"

	"github.com/niamfauzi/go-starter/configs"
	"github.com/niamfauzi/go-starter/internal/platform/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("pakai: go run ./cmd/migrate up")
		fmt.Println("atau  : go run ./cmd/migrate down")
		os.Exit(1)
	}

	cfg := configs.MustLoad()

	sqlDB, err := db.NewSQL(cfg.MySQLDSN)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	migrator := db.NewMigrator(sqlDB)

	command := os.Args[1]
	switch command {
	case "up":
		if err := migrator.Up("./migrations"); err != nil {
			panic(err)
		}
		fmt.Println("migration up selesai")
	case "down":
		if err := migrator.DownLast("./migrations"); err != nil {
			panic(err)
		}
		fmt.Println("migration down selesai")
	default:
		fmt.Println("command tidak dikenal, gunakan: up atau down")
		os.Exit(1)
	}
}
