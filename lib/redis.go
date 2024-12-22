package lib

import (
	"context"
	"strconv"

	"github.com/opium-bio/config"
	"github.com/opium-bio/utils"
	"github.com/redis/go-redis/v9"
)

func Redis() *redis.Client {
	cfg, err := config.LoadConfig("./config.toml")
	if err != nil {
		utils.Error("Error loading config", true)
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Hostname + ":" + strconv.Itoa(cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		utils.Error("Failed to connect to Redis: "+err.Error(), true)
		return nil
	}
	utils.Log("Successfully connected to Redis!")
	return client
}
