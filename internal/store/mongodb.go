package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	Client *mongo.Client
	DB     *mongo.Database
	Uri    string
}

func NewMongoStore(uri string) (*MongoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoStore{
		Client: client,
		DB:     client.Database("pbac"),
		Uri:    uri,
	}, nil
}

func (s *MongoStore) Policies() *mongo.Collection {
	return s.DB.Collection("policies")
}

func (s *MongoStore) Roles() *mongo.Collection {
	return s.DB.Collection("roles")
}

func (s *MongoStore) GetClient() *mongo.Client {
	return s.Client
}

func (s *MongoStore) GetURI() string {
	return s.Uri
}
