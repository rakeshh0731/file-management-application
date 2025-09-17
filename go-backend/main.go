package main

import (
	"log"
	"net/http"

	"file-hub-go/api"
	"file-hub-go/config"
	"file-hub-go/database"
	"file-hub-go/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	// Initialize database connections
	database.InitMongoDB()
	database.InitUserDB()

	r := chi.NewRouter()

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.AppConfig.AllowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any browser
	})

	r.Use(chiMiddleware.Logger)
	r.Use(c.Handler)

	// --- Public Routes ---
	r.Group(func(r chi.Router) {
		r.Post("/api/auth/register", api.RegisterUser)
		r.Post("/api/auth/login", api.LoginUser)
	})

	// --- Protected Routes ---
	// These routes require a valid JWT
	r.Group(func(r chi.Router) {
		// Apply the JWT authentication middleware
		r.Use(middleware.JwtAuthentication)

		// File related routes
		r.Get("/api/files/", api.GetFiles)
		r.Post("/api/files/", api.UploadFile)
		r.Delete("/api/files/{id}/", api.DeleteFile)
	})

	// Serve static files from the 'uploads' directory
	fs := http.FileServer(http.Dir(config.AppConfig.UploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fs))

	log.Printf("Server is running on port %s", config.AppConfig.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.AppConfig.ServerPort, r))
}
