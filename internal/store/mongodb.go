package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Role     string `bson:"role" json:"role"`
}

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

func (s *MongoStore) Users() *mongo.Collection {
	return s.DB.Collection("users")
}

// GetAllUsers retrieves all users from the database
func (s *MongoStore) GetAllUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []User
	cursor, err := s.Users().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
