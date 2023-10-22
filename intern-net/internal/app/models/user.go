package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents user data model
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

// Represents data needed for user Login
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
