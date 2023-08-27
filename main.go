package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/4serviceSoftware/tech-task/handlers"
	"github.com/4serviceSoftware/tech-task/nodes"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	logger := log.New(os.Stdout, "tech-task-one ", log.LstdFlags)

	// getting main database connetcion
	// TODO: get db credentials from config
	dbUrl := "postgres://kbnq:root@localhost:5432/techtaskone"
	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		logger.Fatal("DB conn: " + err.Error())
	}
	defer db.Close()

	ctx := context.Background()

	// creating nodes repository
	nodesRepo := nodes.NewRepositoryPostgres(db, ctx)

	nh := handlers.NewNodes(nodesRepo, logger)

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/nodes", nh.Get).Methods("GET")
	r.HandleFunc("/nodes", nh.Post).Methods("POST")
	http.Handle("/", r)

	// TODO: get all this server settings from some store
	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}

	// start the server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	logger.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting for current operations to complete
	// TODO: get waiting time from some settings store
	shutdownCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(shutdownCtx)
}
