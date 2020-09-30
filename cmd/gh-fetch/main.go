package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/karthikraobr/gh-fetch/internal/cache"
	"github.com/karthikraobr/gh-fetch/internal/gh"
	"github.com/karthikraobr/gh-fetch/internal/handlers"
	"github.com/karthikraobr/gh-fetch/internal/store"
)

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("PORT")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbName, password)
	log := log.New(os.Stdout, "", log.LstdFlags)
	store, err := store.New(connectionString, log)
	if err != nil {
		log.Fatal("could not initialize database")
		return
	}
	h := handlers.New(gh.New(nil, log), log, store, cache.New(100, 60))
	r := h.SetUpRouter()
	log.Fatal(r.Run(":8000"))
}
