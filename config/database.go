package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func ConnectDB() {
	var err error
	dsn := "postgres://postgres:sekretik123@localhost:5432/futbikSecond?sslmode=disable"
	log.Println("Подключение к базе данных:", dsn)

	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка подклчения к базе: ", err)
	}
	fmt.Println("Подключение к базе данных установлено!")
}
