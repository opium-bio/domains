package lib

import (
	"context"
	"sync"

	"github.com/opium-bio/config"
	"github.com/opium-bio/utils"
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
			utils.Error("Error loading config", true)
		}
		clientOptions := options.Client().ApplyURI(cfg.MongoDB.String)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			utils.Error("Error connecting to MongoDB", true)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			utils.Error("Error pinging MongoDB", true)
		}

		utils.Log("Successfully connected to Mongodb!")
		clientInstance = client
	})

	if clientInstanceErr != nil {
		utils.Error("MongoDB client initialization error", true)
	}

	return clientInstance
}
