package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/opium-bio/config"
	"github.com/opium-bio/lib"
	"github.com/opium-bio/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Domain struct {
	Domain string `json:"domain"`
}

func main() {
	cfg, err := config.LoadConfig("./config.toml")
	if err != nil {
		utils.Error("Error loading config", true)
	}

	db := lib.MongoDB()
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			utils.Error("Error disconnecting MongoDB", true)
		}
	}()

	lib.CloudFlare()
	defer func() {
		utils.Error("Error", true)
	}()

	redis := lib.Redis()
	defer func() {
		if redis.Close() != nil {
			utils.Error("Error disconnecting Redis", true)
		}
	}()

	app := fiber.New()
	v1 := app.Group("/v1")

	v1.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"meow": "🐱",
		})
	})

	v1.Get("/domains", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{})
	})

	v1.Post("/add", func(c fiber.Ctx) error {
		domain := new(Domain)
		if err := c.Bind().Body(domain); err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body: " + err.Error(),
			})
		}

		if domain.Domain == "" {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain cannot be empty",
			})
		}

		err := redis.Set(context.Background(), "domain:"+domain.Domain, domain.Domain, 0).Err()
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to save to Redis",
			})
		}

		collection := db.Database("dev").Collection("domains")
		_, err = collection.InsertOne(context.Background(), bson.M{
			"domain": domain.Domain,
		})
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to save to MongoDB",
			})
		}

		err = lib.AddDomain(domain.Domain)
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to add domain to Cloudflare",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"domain":  domain.Domain,
		})
	})

	v1.Get("/checkdomain", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"meow": "🐱",
		})
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.App.Port), fiber.ListenConfig{
		DisableStartupMessage: true,
	}))
}
