package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zlyjoker102/simple-rest-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "logger", log.LstdFlags)

	uh := handlers.NewUser(l)
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", uh.GetUsers)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", uh.AddUser)
	postRouter.Use(uh.MiddlewareUserValid)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", uh.UpdateUser)
	putRouter.Use(uh.MiddlewareUserValid)


	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id:[0-9]+}", uh.DeleteUser)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        sm,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recived terminate, graceful shutdown", sig)
	s.ListenAndServe()

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
