package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
	uri    string
}

func NewMongoStore(uri string) (*MongoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoStore{
		client: client,
		db:     client.Database("pbac"),
		uri:    uri,
	}, nil
}

func (s *MongoStore) Policies() *mongo.Collection {
	return s.db.Collection("policies")
}

func (s *MongoStore) Roles() *mongo.Collection {
	return s.db.Collection("roles")
}

func (s *MongoStore) GetClient() *mongo.Client {
	return s.client
}

func (s *MongoStore) GetURI() string {
	return s.uri
}
