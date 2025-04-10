package main

import (
	"futbikSecond/config"
	_ "futbikSecond/docs"
	"futbikSecond/routes"
)

// @title			My API
// @version		1.0
// @description	API for my Go project
// @host			localhost:8080
func main() {
	config.ConnectDB()
	defer config.DB.Close()

	config.RunMigrations()

	r := routes.Routes()
	r.Run(":8080")
}
