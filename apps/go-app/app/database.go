package app

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var DBClient *mongo.Client

func StartMongo(uri string) {
	opts := options.Client()

	// mongo otel instrumentation
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(uri)
	DBClient, _ = mongo.Connect(context.Background(), opts)

	// seed database
	docs := []interface{}{
		bson.D{{"id", "1"}, {"title", "todo 1"}},
		bson.D{{"id", "2"}, {"title", "todo 2"}},
	}
	DBClient.Database("todo").Collection("todos").InsertMany(context.Background(), docs)
}
