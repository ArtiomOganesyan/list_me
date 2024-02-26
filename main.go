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
)

func main() {

	var DbAddr string = "postgres://admin:admin@localhost:6000/list_me?sslmode=disable"
	var Port string = ":3000"

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