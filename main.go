package main

import (
	"context"
	"list_me/db"
	"list_me/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// {th!2345
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file, assuming environment variables are set")
	}

	DbAddr := os.Getenv("DB_ADDR")
	if DbAddr == "" {
		log.Fatalf("DB_ADDR environment variable is not set. Application cannot start.")
	}

	Port := os.Getenv("PORT")
	if Port == "" {
		log.Println("PORT environment variable not set, defaulting to 8080")
		Port = "5999" // Default to port 8080 if not set
	}

	storage, err := db.NewStorage(DbAddr)
	if err != nil {
		log.Fatalf("db.NewStorage(): %v", err)
	}
	defer storage.Close()

	if err := storage.Init(); err != nil {
		log.Fatalf("storage.Init(): %v", err)
	}

	server := server.NewServer(Port, storage).Init()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	log.Println("Server started")

	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Server stopped")
}
