package paas

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// connect mongodb database
func ConnectMongoDB(dsn string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	clinetOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(ctx, clinetOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}
