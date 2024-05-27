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
	Location    []float64              `bson:"location" json:"location"`
	YearOfBirth int                    `bson:"year_of_birth" json:"year_of_birth"`
	Tags        []string               `bson:"tags" json:"tags"`
	Info        map[string]interface{} `bson:"info" json:"info"`
}

// NewUser creates a new User object
func NewUser(userMap map[string]interface{}) (*User, error) {
	user := &User{}
	user.UserID = userMap["user_id"].(string)
	user.UserName = userMap["user_name"].(string)
	user.CreatedAt = userMap["created_at"].(time.Time)
	user.YearOfBirth = userMap["year_of_birth"].(int)
	user.Tags = userMap["tags"].([]string)
	user.Info = userMap["info"].(map[string]interface{})
	return user, nil
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

// GetAllYearsOfBirth retrieves all years of birth
func GetAllYearsOfBirth() ([]int, error) {
	resultSlice, err := db.MongoDB.Database.Collection("users").Distinct(context.Background(), "year_of_birth", bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find all years of birth")
	}
	var yearsOfBirth []int
	for _, result := range resultSlice {
		yearOfBirth := result.(int)
		yearsOfBirth = append(yearsOfBirth, yearOfBirth)
	}

	return yearsOfBirth, nil
}

// GetUsersByCreatedTime retrieves users by created time(start and end)
func GetUsersByCreatedTime(startTime time.Time, endTime time.Time) ([]User, error) {
	filter := bson.M{"created_at": bson.M{"$gte": startTime, "$lte": endTime}}
	cursor, err := db.MongoDB.Database.Collection("users").Find(context.Background(), filter)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find users by created time %s ~ %s", startTime, endTime))
	}
	var users []User
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decode users by created time %s ~ %s", startTime, endTime))
	}
	return users, nil
}
