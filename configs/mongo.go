package configs

import (
	"context"
	log "github.com/sirupsen/logrus"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"person/internal/useful"
	"time"
)

func mongo() *driver.Collection {

	client, _ := driver.NewClient(options.Client().ApplyURI(properties.Mongo.Uri))
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	_ = client.Connect(ctx)
	err := client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(useful.ConnectDbError, err)
	}

	return client.Database(properties.Mongo.Database).Collection(properties.Mongo.Collection)
}