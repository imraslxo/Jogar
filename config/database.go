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
	dsn := "postgres://postgres:sekretik123@localhost:5432/funbikSecond?sslmode=disable"

	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка подклчения к базе: ", err)
	}
	fmt.Println("Подключение к базе данных установлено!")
}
