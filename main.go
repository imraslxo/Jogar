package main

import (
	"futbikSecond/config"
	_ "futbikSecond/docs"
	"futbikSecond/routes"
	"github.com/joho/godotenv"
	"log"
)

// @title			My API
// @version		1.0
// @description	API for my Go project
// @host			localhost:8080
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла")
	}

	config.ConnectDB()
	defer config.DB.Close()

	config.RunMigrations()

	r := routes.Routes()
	r.Run(":8080")
}
