package Database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var client *mongo.Client
var Err error = nil
var Database *mongo.Database

func InitDB() {
	client, Err = mongo.Connect(
		context.TODO(),
		options.Client().
			ApplyURI(
				os.Getenv("MONGODB_URI")))
	if Err != nil {
		return
	}
	Database = client.Database("authServer")
}
