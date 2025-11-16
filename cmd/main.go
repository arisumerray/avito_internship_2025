package main

import (
	"avito-test-2025/internal/api"
	"avito-test-2025/internal/repository/postgresql"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	repo, err := postgresql.NewRepository(postgresql.Config{
		Host:     "db",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "appdb",
	})
	if err != nil {
		log.Fatal(err)
	}
	app := api.New(repo)
	app.RegisterRoutes(r)
	r.Run(":8080")
}
