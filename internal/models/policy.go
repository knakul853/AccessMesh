package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Policy struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Role       string             `bson:"role" json:"role"`
	Resource   string             `bson:"resource" json:"resource"`
	Action     string             `bson:"action" json:"action"`
	Conditions PolicyConditions   `bson:"conditions" json:"conditions"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type PolicyConditions struct {
	IPRange   []string `bson:"ip_range" json:"ip_range"`
	TimeRange []string `bson:"time_range" json:"time_range"`
}
