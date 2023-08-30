package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/4serviceSoftware/tech-task/config"
	"github.com/4serviceSoftware/tech-task/db"
	"github.com/4serviceSoftware/tech-task/handlers"
	"github.com/4serviceSoftware/tech-task/internal/nodes"
	"github.com/4serviceSoftware/tech-task/internal/nodes/cachefile"
	"github.com/4serviceSoftware/tech-task/internal/repos"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "tech-task-one ", log.LstdFlags)

	ctx := context.Background()

	config, err := config.GetFromEnv(ctx)
	if err != nil {
		logger.Fatal("Config loading fail: " + err.Error())
	}

	// getting main database connetcion
	dbConn, err := db.GetPostgresConnection(ctx, config)
	if err != nil {
		logger.Fatal("DB conn: " + err.Error())
	}
	defer dbConn.Close()

	// creating nodes repository, nodes cashe and nodes service
	nodesRepo := repos.NewNodesRepositoryPostgres(ctx, dbConn)
	nodesCachefile := cachefile.NewCacheFile(config.NodesCacheFilename)
	nodesService := nodes.NewService(nodesRepo, nodesCachefile)

	nodesHandlers := handlers.NewNodes(nodesService, logger, config)

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/nodes", nodesHandlers.Get).Methods(http.MethodGet)
	r.HandleFunc("/nodes", nodesHandlers.Post).Methods(http.MethodPost)
	http.Handle("/", r)

	s := &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      r,
		IdleTimeout:  config.ServerIdleTimeout * time.Second,
		ReadTimeout:  config.ServerReadTimeout * time.Minute,
		WriteTimeout: config.ServerWriteTimeout * time.Second,
	}

	// start the server
	go func() {
		logger.Printf("Starting server at addr %s...\n", s.Addr)
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
	shutdownCtx, _ := context.WithTimeout(context.Background(), config.ServerShutdownTimeout*time.Second)
	s.Shutdown(shutdownCtx)
}
