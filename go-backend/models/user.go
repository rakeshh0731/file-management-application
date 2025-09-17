package models

import "time"

// User defines the structure for a user in the database.
type User struct {
	ID           string    `bson:"_id" json:"id"`
	Username     string    `bson:"username" json:"username"`
	PasswordHash string    `json:"-"` // Do not expose password hash in JSON responses
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}
