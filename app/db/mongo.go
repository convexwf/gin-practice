package db

import (
	"context"
	"time"

	. "github.com/convexwf/gin-practice/app/util"
	"github.com/pkg/errors"
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
	if MongoDB == nil {
		Log.Warn("MongoDB is not connected")
	}
	if err := MongoDB.Client.Disconnect(context.Background()); err != nil {
		Log.Error("failed to disconnect from MongoDB")
	}
	MongoDB = nil
}
