package store

import (
	"context"
	"time"

	"github.com/knakul853/accessmesh/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *MongoStore) Users() *mongo.Collection {
	return s.DB.Collection("users")
}

func (s *MongoStore) GetClient() *mongo.Client {
	return s.Client
}

func (s *MongoStore) GetURI() string {
	return s.Uri
}

// GetAllUsers retrieves all users from the database
func (s *MongoStore) GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
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

// UpdateUser updates a user in the database
func (s *MongoStore) UpdateUser(id string, user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":     user.Email,
			"role":      user.Role,
			"updatedAt": time.Now(),
		},
	}

	_, err = s.Users().UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// DeleteUser deletes a user from the database
func (s *MongoStore) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.Users().DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
