package lib

import (
	"context"
	"log"
	"sync"

	"github.com/opium-bio/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
)

func MongoDB() *mongo.Client {
	mongoOnce.Do(func() {
		cfg, err := config.LoadConfig("./config.toml")
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		clientOptions := options.Client().ApplyURI(cfg.MongoDB.String)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatalf("Error connecting to MongoDB: %v", err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatalf("Error pinging MongoDB: %v", err)
		}

		log.Printf("MongoDB connection established successfully")
		clientInstance = client
	})

	if clientInstanceErr != nil {
		log.Fatalf("MongoDB client initialization error: %v", clientInstanceErr)
	}

	return clientInstance
}
