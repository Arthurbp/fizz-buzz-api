package main

import (
	"context"
	"encoding/json"
	"errors"
	"fizz-buzz-api/pkg/store"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func main() {
	l := log.New(os.Stdout, "", log.Lshortfile)

	// Connect to mongo
	s, err := store.NewClient("mongodb://mongo:27017", l)
	if err != nil {
		l.Fatal(err)
	}
	defer func() {
		if err = s.C.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Init HTTP router and setup routes
	h := handler(s, l)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		Handler:      h,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 120,
	}

	// Start the server
	go func() {
		l.Println("Starting server on port 8080")

		err := srv.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Println("Got signal:", sig)

	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func handler(s *store.Store, l *log.Logger) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/fizzbuzz", resultHandler(s, l)).Methods("GET")
	r.HandleFunc("/fizzbuzz/stats", statsHandler(s, l)).Methods("GET")
	return r
}

func resultHandler(s *store.Store, l *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q, err := parseFizzBuzzParams(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := FizzBuzzResponse(q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		s.InsertFizzBuzzQuery(ctx, q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func statsHandler(s *store.Store, l *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		result, err := s.AggregateFizzBuzzQueries(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(result) == 0 {
			http.Error(w, "No request have been sent", http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result[0])
	}
}

func parseFizzBuzzParams(r *http.Request) (store.FizzBuzzQuery, error) {
	var decoder = schema.NewDecoder()

	var q store.FizzBuzzQuery

	err := decoder.Decode(&q, r.URL.Query())
	if err != nil {
		return q, errors.New("Failed to parse parameters")
	}
	if q.Str1 == "" {
		return q, errors.New("string1 parameter is empty")
	}
	if q.Str2 == "" {
		return q, errors.New("string2 parameter is empty")
	}
	if q.Int1 < 1 {
		return q, errors.New("int1 must be greater than or equal to 1")
	}
	if q.Int2 < 1 {
		return q, errors.New("int2 must be greater than or equal to 1")
	}
	if q.Limit < 1 {
		return q, errors.New("limit must be greater than or equal to 1")
	}
	if q.Limit > 1000000 {
		return q, errors.New("limit must be less than 1000000")
	}
	return q, err
}
