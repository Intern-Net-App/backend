package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Job struct {
	ID              primitive.ObjectID `bson:"_id"`
	JobTitle        string             `bson:"job_title"`
	JobDetailUrl    string             `bson:"job_detail_url"`
	JobListed       string             `bson:"job_listed"`
	CompanyName     string             `bson:"company_name"`
	CompanyLink     string             `bson:"company_link"`
	CompanyLocation string             `bson:"company_location"`
}
