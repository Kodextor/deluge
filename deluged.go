// Steve Phillips / elimisteve
// 2013.09.08

package main

import (
	"./handlers"
	// "./types"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	DATABASE_NAME = "deluge"
)

var (
	router = mux.NewRouter()
)

func init() {
	if DATABASE_NAME == "" {
		log.Fatal("Must set DATABASE_NAME")
	}

	// Define routes
	router.HandleFunc("/", handlers.GetIndex).Methods("GET")
	router.HandleFunc("/users/new", handlers.PostUser).Methods("POST")
	router.HandleFunc("/subdomains/new", handlers.PostSubdomain).Methods("POST")
	// router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	// router.HandleFunc("/subdomain", handlers.PostSubdomain).Methods("PUT")
	// router.HandleFunc("/", handlers.PayPalCallback).Methods("POST")
	// router.HandleFunc("/token", handlers.PostToken).Methods("POST")

	http.Handle("/", router)
}

func main() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Start HTTP server
	server := SimpleHTTPServer(router, ":9090")
	log.Printf("HTTP server trying to listen on %v...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP listen failed: %v\n", err)
	}
}

func SimpleHTTPServer(handler http.Handler, host string) *http.Server {
	server := http.Server{
		Addr:           host,
		Handler:        handler,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &server
}
