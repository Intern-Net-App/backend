package repositories

import (
	"context"

	"intern-net/internal/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handles MongoDB operations related to Users.
type UserRepository struct {
	collection *mongo.Collection
}

// Creates a new UserRepository instance
func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection}
}

// Find user by ID from MongoDB
func (ur *UserRepository) FindByID(id string) (*models.User, error) {
	filter := bson.M{"_id": id}

	var user models.User
	err := ur.collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
