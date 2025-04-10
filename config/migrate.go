package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib" // нужен для совместимости goose с pgx
	"github.com/pressly/goose/v3"
)

func RunMigrations() {
	dsn := os.Getenv("DB_URL")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД для миграции: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Ошибка установки диалекта goose: %v", err)
	}

	fmt.Println("🚀 Запуск миграций...")
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Ошибка при выполнении миграций: %v", err)
	}
	fmt.Println("Миграции применены успешно")
}
