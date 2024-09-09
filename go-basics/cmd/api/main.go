package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-basics/internal/adapters/client"
	"github.com/osalomon89/go-basics/internal/adapters/handler"
	mysqlrepo "github.com/osalomon89/go-basics/internal/adapters/repository/mysql"
	"github.com/osalomon89/go-basics/internal/core/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// esCient, err := es8.NewDefaultClient()
	// if err != nil {
	// 	return err
	// }

	// repo := ds.NewEsRepository(esCient)
	// if err := repo.CreateIndex("items"); err != nil {
	// 	log.Fatalln(err)
	// }

	conn, err := mysqlrepo.GetConnectionDB()
	if err != nil {
		return err
	}

	repo := mysqlrepo.NewMySQLRepository(conn)

	providerClient := client.NewProviderClient()
	service := service.NewService(repo, providerClient) //itemServiceImpl

	h := handler.NewHandler(service) //ItemService

	r := mux.NewRouter()

	router := r.PathPrefix("/api-items/v1").Subrouter()

	router.HandleFunc("/hello", h.HelloHandler).Methods(http.MethodGet)

	router.HandleFunc("/items", h.GetAllItems).Methods(http.MethodGet)
	router.HandleFunc("/items", h.CreateItem).Methods(http.MethodPost)
	router.HandleFunc("/items/{id}", h.GetItemByID).Methods(http.MethodGet)
	router.HandleFunc("/items/{id}", h.UpdateItem).Methods(http.MethodPut)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

	return nil
}
