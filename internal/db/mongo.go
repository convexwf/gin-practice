package db

import (
	"context"
	"reflect"
	"time"

	. "github.com/convexwf/gin-practice/internal/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoModel struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var MongoDB *MongoModel

func ConnectDB() error {

	if MongoDB != nil {
		return nil
	}

	mongoURI := Config.GetString("mongo.uri")
	mongoDbName := Config.GetString("mongo.db_name")
	mongoConnectTimeout := Config.GetDuration("mongo.connect_timeout")

	MongoDB = &MongoModel{}
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout*time.Second)
	defer cancel()
	if MongoDB.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI)); err != nil {
		return errors.Wrap(err, "failed to connect to MongoDB")
	}
	if err = MongoDB.Client.Ping(ctx, readpref.Primary()); err != nil {
		return errors.Wrap(err, "failed to ping MongoDB")
	}
	MongoDB.Database = MongoDB.Client.Database(mongoDbName)

	Log.Info("successfully connected to MongoDB")
	return nil
}

func DisconnectDB() {
	if MongoDB == nil || MongoDB.Client == nil {
		Log.Warn("MongoDB is not connected")
	}
	if err := MongoDB.Client.Disconnect(context.Background()); err != nil {
		Log.Error("failed to disconnect from MongoDB")
	}
	MongoDB = nil
}

func InsertAWSRecords(records []AWSData) error {
	var dataToInsert []interface{}

	for _, record := range records {
		// Convert time.Time to primitive.DateTime
		receivedTime := primitive.NewDateTimeFromTime(record.ReceivedTime)

		// Create a new document to insert
		doc := bson.M{
			"station_id":    record.StationID,
			"received_time": receivedTime,
			"weather_info":  record.WeatherInfo,
		}

		dataToInsert = append(dataToInsert, doc)
	}

	_, err := db.MongoDB.Database.Collection("aws_data").InsertMany(context.Background(), dataToInsert)
	if err != nil {
		return err
	}

	return nil
}

func ConvertStructToBSONDoc(data interface{}) bson.D {
	doc := bson.D{}
	val := reflect.ValueOf(data)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		key := field.Tag.Get("bson")
		if key == "" {
			key = field.Name
		}

		value := val.Field(i).Interface()

		switch value.(type) {
		case []string:
			sliceVal := reflect.ValueOf(value)
			slice := make([]interface{}, sliceVal.Len())
			for j := 0; j < sliceVal.Len(); j++ {
				slice[j] = sliceVal.Index(j).Interface()
			}
			doc = append(doc, bson.E{Key: key, Value: slice})
		case map[string]interface{}:
			mapVal := reflect.ValueOf(value)
			doc = append(doc, bson.E{Key: key, Value: mapVal})
		default:
			doc = append(doc, bson.E{Key: key, Value: value})
		}
	}

	return doc
}
