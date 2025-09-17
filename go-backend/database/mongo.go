package database

import (
	"context"
	"file-hub-go/config"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is the connection handle for the database client.
var DB *mongo.Client

// FileCollection is a handle to the "files" collection in MongoDB.
// We'll use this to perform operations on our file documents.
var FileCollection *mongo.Collection

func InitMongoDB() {
	// Ensure the uploads directory exists for storing physical files.
	if err := os.MkdirAll(config.AppConfig.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	// --- Database Connection ---
	// Get the connection string from environment variables with a fallback.
	// The `authSource=admin` is often needed when your user is defined in the `admin` database.
	connectionString := config.AppConfig.MongoURI
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB. We use a context with a timeout to prevent
	// the application from hanging indefinitely if it can't connect.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary to verify that the connection is alive.
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	DB = client

	// Get a handle for your "files" collection within the "filehub" database.
	FileCollection = client.Database("filehub").Collection("files")
}
