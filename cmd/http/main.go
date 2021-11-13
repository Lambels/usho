package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lambels/usho/handlers"
	"github.com/Lambels/usho/repo/file"
)

func main() {
	r, err := file.New("hello")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	ctx := context.Background()

	server := http.Server{
		Addr:    ":8080",
		Handler: handlers.NewService(r, ctx),
	}

	go func() {
		log.Println("Server Starting")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		} else {
			log.Println("Server Closed")
		}
	}()

	// Graceful shutdown
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigQuit
	log.Printf("cought signal: %v\n", sig)
	log.Printf("Gracefully shutting down server\n")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
