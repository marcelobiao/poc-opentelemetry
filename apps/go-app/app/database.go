package app

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var DBClient *mongo.Client

func StartMongo() {
	opts := options.Client()
	// mongo OTEL instrumentation
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI("mongodb://test:test@localhost:27017")
	DBClient, _ = mongo.Connect(context.Background(), opts)

	// seed the database with some todo's
	docs := []interface{}{
		bson.D{{"id", "1"}, {"title", "Buy groceries"}},
		bson.D{{"id", "2"}, {"title", "install Aspecto.io"}},
		bson.D{{"id", "3"}, {"title", "Buy dogz.io domain"}},
	}
	DBClient.Database("todo").Collection("todos").InsertMany(context.Background(), docs)
}
