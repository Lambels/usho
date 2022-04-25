package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lambels/usho/handlers"
	"github.com/Lambels/usho/repo"
	"github.com/Lambels/usho/repo/file"
	"github.com/Lambels/usho/repo/mysql"
)

var (
	db   = flag.Bool("db", false, "indicates wether to use mysql store, if ignored will inmem store")
	dsn  = flag.String("dsn", "", "data source name for mysql database")
	path = flag.String("path", "./store", "indicates where the file storage should be located")
)

func main() {
	flag.Parse()
	var r repo.Repo
	var err error

	if *db {
		r, err = mysql.New(*dsn)
	} else {
		r, err = file.New(*path)
	}

	if err != nil {
		log.Fatalf("failed to open store: %v", err)
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
	log.Printf("Cought signal: %v\n", sig)
	log.Printf("Gracefully shutting down server\n")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
