package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/models"
	"github.com/knakul853/accessmesh/internal/store"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestStore struct {
	*store.MongoStore
	client *mongo.Client
}

func setupTestStore(t *testing.T) *TestStore {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Use a test database
	db := client.Database("pbac_test")

	// Clean up any existing data
	err = db.Drop(ctx)
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}

	mongoStore := &store.MongoStore{
		Client: client,
		DB:     db,
	}

	return &TestStore{
		MongoStore: mongoStore,
		client:     client,
	}
}

func (ts *TestStore) Cleanup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ts.client.Disconnect(ctx); err != nil {
		t.Errorf("Failed to disconnect from MongoDB: %v", err)
	}
}

func TestPolicyHandler_Create(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.New()
	testStore := setupTestStore(t)
	defer testStore.Cleanup(t)

	handler := NewPolicyHandler(testStore.MongoStore)
	router.POST("/policies", handler.Create)

	policy := models.Policy{
		Role:     "manager",
		Resource: "/api/v1/orders",
		Action:   "read",
		Conditions: models.PolicyConditions{
			IPRange:   []string{"10.0.0.0/16"},
			TimeRange: []string{"08:00-20:00"},
		},
	}

	body, _ := json.Marshal(policy)
	req := httptest.NewRequest("POST", "/policies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Test
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Policy
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, policy.Role, response.Role)
	assert.NotEmpty(t, response.ID)
}
