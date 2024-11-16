package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username          string             `bson:"username" json:"username"`
	Password          string             `bson:"password" json:"-"`
	Role              string             `bson:"role" json:"role"`
	Email             string             `bson:"email" json:"email"`
	EmailVerified     bool               `bson:"email_verified" json:"email_verified"`
	VerificationToken string             `bson:"verification_token,omitempty" json:"-"`
	ResetToken        string             `bson:"reset_token,omitempty" json:"-"`
	ResetTokenExpiry  time.Time          `bson:"reset_token_expiry,omitempty" json:"-"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
