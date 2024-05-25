package model

import (
	"context"
	"fmt"
	"time"

	"github.com/convexwf/gin-practice/app/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	UserID      string                 `bson:"user_id" json:"user_id"`
	UserName    string                 `bson:"user_name" json:"user_name"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	YearOfBirth int                    `bson:"year_of_birth" json:"year_of_birth"`
	MoneySpent  float64                `bson:"money_spent" json:"money_spent"`
	Tags        []string               `bson:"tags" json:"tags"`
	Info        map[string]interface{} `bson:"info" json:"info"`
}

// NewUser creates a new User object
func NewUser(userID, userName string, yearOfBirth int, moneySpent float64, tags []string, info map[string]interface{}) *User {
	return &User{
		UserID:      userID,
		UserName:    userName,
		CreatedAt:   time.Now(),
		YearOfBirth: yearOfBirth,
		MoneySpent:  moneySpent,
		Tags:        tags,
		Info:        info,
	}
}

// GetUserByID retrieves a user by user ID
func GetUserByID(userID string) (*User, error) {
	filter := bson.M{"user_id": userID}
	var user User
	err := db.MongoDB.Database.Collection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find user by ID %s", userID))
	}
	return &user, nil
}

// GetUserByName retrieves a user by user name
func GetUserByName(userName string) (*User, error) {
	filter := bson.M{"user_name": userName}
	var user User
	err := db.MongoDB.Database.Collection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find user by name %s", userName))
	}
	return &user, nil
}

// GetUserByYearOfBirth retrieves users by year of birth
func GetUserByYearOfBirth(yearOfBirth int) ([]User, error) {
	filter := bson.M{"year_of_birth": yearOfBirth}
	cursor, err := db.MongoDB.Database.Collection("users").Find(context.Background(), filter)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find users by year of birth %d", yearOfBirth))
	}
	var users []User
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decode users by year of birth %d", yearOfBirth))
	}
	return users, nil
}
