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

// Retreive user by Email in Database
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil //User not Found
		}
		return nil, err //Error occured
	}

	return &user, nil
}

// Create User in Database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}
