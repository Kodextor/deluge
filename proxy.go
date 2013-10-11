// Steve Phillips / elimisteve
// 2013.09.29

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	HTTP_TARGET_URL = "http://69.196.150.25:8888"
)

var (
	httpProxy *httputil.ReverseProxy
	router = mux.NewRouter()
)

// HTTP
func init() {
	u, err := url.Parse(HTTP_TARGET_URL)
	if err != nil {
		log.Fatalf("Error parsing %s: %v\n", HTTP_TARGET_URL, err)
	}
	httpProxy = httputil.NewSingleHostReverseProxy(u)
}

func init() {
	router.HandleFunc("/{path:.*}", HTTPIndex).Methods("GET")

	http.Handle("/", router)
}

func main() {
	// Start HTTP server
	server := SimpleHTTPServer(router, ":8080")
	log.Printf("HTTP server trying to listen on %v...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP listen failed: %v\n", err)
	}
}

func HTTPIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request to /: `%+v`\n\n", r)
	httpProxy.ServeHTTP(w, r)
}

// Misc

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
