package models

import "time"

// File represents the metadata for a file document in MongoDB.
type File struct {
	// The `bson` tag tells the MongoDB driver how to map this struct field
	// to a document field in the database. `_id` is the default primary key in MongoDB.
	// The `json` tag tells the `encoding/json` package how to serialize this field
	// for API responses.
	ID               string    `bson:"_id" json:"id"`
	File             string    `bson:"file" json:"file"` // Stores the path to the physical file
	OriginalFilename string    `bson:"original_filename" json:"original_filename"`
	FileType         string    `bson:"file_type" json:"file_type"`
	Size             int64     `bson:"size" json:"size"`
	Hash             string    `bson:"hash" json:"hash"`
	UploadedAt       time.Time `bson:"uploaded_at" json:"uploaded_at"`
}