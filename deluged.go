// Steve Phillips / elimisteve
// 2013.04.28

package main

import (
	"./handlers"
	"./types"
	"github.com/bmizerany/pat"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	MONGO_URLS    = "localhost"
	DATABASE_NAME = "deluged"
)

var (
	session *mgo.Session
	db      *mgo.Database
	mux     = pat.New()
)

// Connect to DB
func init() {
	if DATABASE_NAME == "" {
		log.Fatal("Must set DATABASE_NAME")
	}

	var err error
	session, err = mgo.Dial(MONGO_URLS)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB '%s'", MONGO_URLS)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(DATABASE_NAME)

	// helpers.SetDB(db)
	types.SetDB(db)
}

// Define routes
func init() {
	mux.Get("/", http.HandlerFunc(handlers.GetIndex))
	mux.Post("/users/new", http.HandlerFunc(handlers.PostUser))
	mux.Post("/subdomains/new", http.HandlerFunc(handlers.PostSubdomain))
	// mux.Put("/subdomain", http.HandlerFunc(handlers.PostSubdomain))
	// mux.Post("/", http.HandlerFunc(handlers.PayPalCallback))
	// mux.Post("/token", http.HandlerFunc(handlers.PostToken))
	
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer session.Close()
	createUserIndex()

	// Start HTTP server
	server := SimpleHTTPServer(mux, ":9090")
	log.Printf("HTTP server trying to listen on %v...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP listen failed: %v\n", err)
	}
}

func createUserIndex() {
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	err := db.C("users").EnsureIndex(index)
	if err != nil {
		log.Printf("Tried to create users index: %v\n", err)
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
