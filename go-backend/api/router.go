package api

import (
    "github.com/go-chi/cors"
	// We import the `chi` package that we downloaded with `go get`.
	// This gives us the router functionality.
 	"github.com/go-chi/chi/v5"
 	"github.com/go-chi/chi/v5/middleware"
 	// We also import the `net/http` package, which is a built-in Go
 	"path/filepath"
 	// package for all things related to HTTP.
 	"net/http"
 )

// NewRouter creates and configures a new router. We will call this from main.go.
func NewRouter() http.Handler {
	// Create a new router instance.
	r := chi.NewRouter()

 	// --- CORS Middleware ---
 	// This allows our React frontend (running on localhost:3000) to make requests
 	// to our Go backend (running on localhost:8000).
 	r.Use(cors.Handler(cors.Options{
 		AllowedOrigins:   []string{"http://localhost:3000"},
 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
 		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
 		ExposedHeaders:   []string{"Link"},
 		AllowCredentials: true,
 		MaxAge:           300, // Maximum value not ignored by any major browsers
 	}))

	// --- Middleware ---
	// Middleware are functions that run on every request.
	// Logger prints a log line for each request to your console. Very useful for debugging!
	r.Use(middleware.Logger)
	// Recoverer gracefully handles panics and prevents the server from crashing.
	r.Use(middleware.Recoverer)
	// StripSlashes is a middleware that will match request paths with no trailing slashes.
 	r.Use(middleware.StripSlashes)

	// Define a route for the root path "/".
	// When a request comes to "/", the function provided will be executed.
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go backend is running!"))
	})

 	// --- API Routes ---
 	// This tells the router to use our new GetFiles function for requests to "/api/files/".
 	r.Get("/api/files", GetFiles)
 	r.Post("/api/files", UploadFile)
 	r.Delete("/api/files/{id}", DeleteFile)

 	// --- File Serving ---
 	// This creates a file server that serves static files from the "uploads" directory.
 	// The URL path "/uploads/" is stripped, so a request to "/uploads/foo.txt"
 	// will look for the file "uploads/foo.txt" on the disk.
 	uploadsDir := http.Dir(filepath.Join(".", "uploads"))
 	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(uploadsDir)))


	// Return the fully configured router.
	return r
}