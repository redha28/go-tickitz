package main

import (
	"gotickitz/internal/routes"
	"gotickitz/pkg"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	pg, err := pkg.Connect()
	if err != nil {
		log.Printf("[ERROR] Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		log.Println("DB connection closed")
		pg.Close()
	}()
	router := routes.InitRouter(pg)

	router.Run("localhost:8080")
}
