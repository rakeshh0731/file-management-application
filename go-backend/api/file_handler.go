package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"file-hub-go/config"
	"file-hub-go/database"
	"file-hub-go/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetFiles handles the logic for listing and filtering files.
func GetFiles(w http.ResponseWriter, r *http.Request) {
	// A filter document for our MongoDB query. bson.D preserves order.
	filter := bson.D{}
	params := r.URL.Query()

	// --- Filtering Logic (similar to your Django backend) ---

	// Search by filename (case-insensitive)
	if search := params.Get("search"); search != "" {
		// This creates a filter like: { "original_filename": { "$regex": "...", "$options": "i" } }
		filter = append(filter, bson.E{Key: "original_filename", Value: bson.D{{Key: "$regex", Value: search}, {Key: "$options", Value: "i"}}})
	}

	// Filter by file type (case-insensitive)
	if fileType := params.Get("file_type"); fileType != "" {
		filter = append(filter, bson.E{Key: "file_type", Value: bson.D{{Key: "$regex", Value: fileType}, {Key: "$options", Value: "i"}}})
	}

	// Filter by min size
	if sizeMinStr := params.Get("size_min"); sizeMinStr != "" {
		if sizeMin, err := strconv.ParseInt(sizeMinStr, 10, 64); err == nil {
			filter = append(filter, bson.E{Key: "size", Value: bson.D{{Key: "$gte", Value: sizeMin}}})
		}
	}

	// Filter by max size
	if sizeMaxStr := params.Get("size_max"); sizeMaxStr != "" {
		if sizeMax, err := strconv.ParseInt(sizeMaxStr, 10, 64); err == nil {
			filter = append(filter, bson.E{Key: "size", Value: bson.D{{Key: "$lte", Value: sizeMax}}})
		}
	}

	// Filter by uploaded after date
	if after := params.Get("uploaded_after"); after != "" {
		if t, err := time.Parse("2006-01-02", after); err == nil {
			filter = append(filter, bson.E{Key: "uploaded_at", Value: bson.D{{Key: "$gte", Value: t}}})
		}
	}

	// Filter by uploaded before date
	if before := params.Get("uploaded_before"); before != "" {
		if t, err := time.Parse("2006-01-02", before); err == nil {
			// Add 1 day to be inclusive of the 'before' date
			t = t.AddDate(0, 0, 1)
			filter = append(filter, bson.E{Key: "uploaded_at", Value: bson.D{{Key: "$lt", Value: t}}})
		}
	}

	// --- End Filtering ---

	var files []models.File
	// Set a timeout for the database operation.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find documents in the collection that match our filter.
	// We also sort by `uploaded_at` in descending order.
	opts := options.Find().SetSort(bson.D{{Key: "uploaded_at", Value: -1}})
	cursor, err := database.FileCollection.Find(ctx, filter, opts)
	if err != nil {
		http.Error(w, "Failed to fetch files from database", http.StatusInternalServerError)
		log.Printf("Error fetching files: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Decode all documents from the cursor into our `files` slice.
	if err = cursor.All(ctx, &files); err != nil {
		http.Error(w, "Failed to decode files", http.StatusInternalServerError)
		log.Printf("Error decoding files: %v", err)
		return
	}

	// If no files are found, return an empty JSON array `[]` instead of `null`.
	if files == nil {
		files = []models.File{}
	}

	// Set the response header and encode the files slice as JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// calculateFileHash calculates the SHA256 hash of a file.
func calculateFileHash(file io.ReadSeeker) (string, error) {
	// Rewind the file to the beginning before hashing
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Rewind the file again so it can be read for saving
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// UploadFile handles the logic for uploading a new file.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	var newFile models.File
	// Parse the multipart form, with a max file size from config
	if err := r.ParseMultipartForm(config.AppConfig.MaxUploadSize); err != nil {
		maxSizeMB := config.AppConfig.MaxUploadSize / 1024 / 1024
		http.Error(w, fmt.Sprintf("The uploaded file is too big. Please choose a file less than %dMB.", maxSizeMB), http.StatusBadRequest)
		return
	}
	// Get the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Calculate the file hash for deduplication
	fileHash, err := calculateFileHash(file)
	if err != nil {
		http.Error(w, "Could not calculate file hash", http.StatusInternalServerError)
		return
	}
	// --- Deduplication Logic ---
	var existingFile models.File
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if a file with the same hash already exists
	err = database.FileCollection.FindOne(ctx, bson.M{"hash": fileHash}).Decode(&existingFile)
	if err == nil {
		// DUPLICATE: A file with this hash exists. Create a new metadata entry
		// pointing to the same physical file.
		newFile = models.File{
			ID:               uuid.New().String(),
			File:             existingFile.File, // Point to the existing file path
			OriginalFilename: handler.Filename,
			FileType:         handler.Header.Get("Content-Type"),
			Size:             handler.Size,
			Hash:             fileHash,
			UploadedAt:       time.Now(),
		}
	} else {
		// NEW FILE: Save the physical file and create a new metadata entry.
		ext := filepath.Ext(handler.Filename)
		newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join(config.AppConfig.UploadDir, newFilename)

		// Create the destination file
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Could not save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file data to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Could not save file", http.StatusInternalServerError)
			return
		}
		newFile = models.File{
			ID:               uuid.New().String(),
			File:             "/" + filePath, // Store the URL path
			OriginalFilename: handler.Filename,
			FileType:         handler.Header.Get("Content-Type"),
			Size:             handler.Size,
			Hash:             fileHash,
			UploadedAt:       time.Now(),
		}
	}
	// Insert the new file metadata into the database
	_, err = database.FileCollection.InsertOne(context.Background(), newFile)
	if err != nil {
		http.Error(w, "Could not save file metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newFile)
}

// DeleteFile handles the logic for deleting a file.
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	// Get the file ID from the URL parameter
	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the file to get its hash and path before deleting
	var fileToDelete models.File
	err := database.FileCollection.FindOne(ctx, bson.M{"_id": fileID}).Decode(&fileToDelete)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	// Delete the metadata entry
	_, err = database.FileCollection.DeleteOne(ctx, bson.M{"_id": fileID})
	if err != nil {
		http.Error(w, "Failed to delete file metadata", http.StatusInternalServerError)
		return
	}
	// Check if any other files reference the same hash
	count, err := database.FileCollection.CountDocuments(ctx, bson.M{"hash": fileToDelete.Hash})
	if err != nil {
		log.Printf("Error checking for other file references: %v", err)
		// Continue without deleting the physical file to be safe
	} else if count == 0 {
		// No other files with the same hash, so delete the physical file
		// The path in the DB is like "/uploads/...", so we remove the leading "/"
		physicalPath := strings.TrimPrefix(fileToDelete.File, "/")
		if err := os.Remove(physicalPath); err != nil {
			log.Printf("Failed to delete physical file %s: %v", physicalPath, err)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
