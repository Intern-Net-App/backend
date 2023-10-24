package repositories

import (
	"context"
	"intern-net/internal/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobRepository struct {
	collection *mongo.Collection
}

// Create a new JobRepository Instance
func NewJobRepository(collection *mongo.Collection) *JobRepository {
	return &JobRepository{collection}
}

func (r *JobRepository) GetJobs(ctx context.Context, skip, limit int) ([]models.Job, error) {
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	filter := bson.M{}

	var jobs []models.Job

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
