package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gtuk/discordwebhook"
	"github.com/opium-bio/config"
	"github.com/opium-bio/lib"
	"github.com/opium-bio/middleware"
	"go.mongodb.org/mongo-driver/bson"
)

type Domain struct {
	Domain string `json:"domain"`
}

type FindDomains struct {
	Domain    string `bson:"domain"`
	Donated   bool   `bson:"donated"`
	Status    string `bson:"status"`
	AddedBy   string `bson:"addedby"`
	DateAdded string `bson:"dateAdded"`
}

type User struct {
	Premium bool `json:"premium"`
}

func ptr(s string) *string {
	return &s
}

func sendDiscordError(cfg *config.Config, message string) {
	webhookURL := cfg.Discord.Webhook
	if webhookURL == "" {
		log.Println("Discord webhook URL is not set in config")
		return
	}
	payload := discordwebhook.Message{
		Content: ptr(fmt.Sprintf(":warning: **API Error**\n```%s```", message)),
	}

	err := discordwebhook.SendMessage(webhookURL, payload)
	if err != nil {
		log.Println("Failed to send error to Discord:", err)
	}
}

func main() {
	cfg, err := config.LoadConfig("./config.toml")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	db := lib.MongoDB()
	if db == nil {
		log.Fatal("Failed to connect to MongoDB")
	}
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		}
	}()
	collection := db.Database("prod").Collection("domains")

	lib.CloudFlare()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Unexpected error:", err)
			sendDiscordError(cfg, fmt.Sprintf("Unexpected error: %v", err))
		}
	}()

	app := fiber.New()
	v1 := app.Group("/v1")
	v1.Post("/add", middleware.JWTAuth(db, cfg.JWT.Secret), func(c *fiber.Ctx) error {
		user := c.Locals("user").(middleware.User)
		domain := new(Domain)
		
		var Usr User
		collection := db.Database("prod").Collection("users")
		err := collection.FindOne(context.Background(), bson.M{"_id": user.ID}).Decode(&Usr)
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to fetch user from the database",
			})
		}
		if !Usr.Premium {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "You need premium to add a domain",
			})
		}
		if err := c.BodyParser(domain); err != nil {
			sendDiscordError(cfg, fmt.Sprintf("Invalid request body: %v", err))
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
		var existingDomain FindDomains
		err = collection.FindOne(context.Background(), bson.M{"domain": domain.Domain}).Decode(&existingDomain)
		if err == nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain already exists in the database",
			})
		}

		newDomain := FindDomains{
			Domain:    domain.Domain,
			Donated:   false,
			Status:    "pending",
			AddedBy:   user.ID.Hex(),
			DateAdded: time.Now().Format(time.RFC3339),
		}
		_, err = collection.InsertOne(context.Background(), newDomain)
		if err != nil {
			sendDiscordError(cfg, fmt.Sprintf("MongoDB insert error: %v", err))
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to save to MongoDB",
			})
		}
		defer func() {
			if r := recover(); r != nil {
				sendDiscordError(cfg, fmt.Sprintf("Cloudflare panic: %v", r))
				c.JSON(fiber.Map{
					"success": false,
					"message": "Internal Cloudflare error occurred",
				})
			}
		}()

		err = lib.AddDomain(domain.Domain)
		if err != nil {
			sendDiscordError(cfg, fmt.Sprintf("Cloudflare API error: %v", err))
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to add domain to Cloudflare. Please try again later.",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Your domain has successfully been added and is pending approval.",
			"domain":  domain.Domain,
		})
	})
	v1.Get("/domain", middleware.JWTAuth(db, cfg.JWT.Secret), func(c *fiber.Ctx) error {
		user := c.Locals("user").(middleware.User)
		domain := c.Query("domain")
		if domain == "" {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain cannot be empty",
			})
		}
		var domainthing FindDomains
		err := collection.FindOne(context.TODO(), bson.M{"domain": domain}).Decode(&domainthing)
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain doesn't exist in the database",
			})
		}
		if domainthing.AddedBy != user.ID.Hex() {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Not your domain (make a ticket in the discord if this is wrong)",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": domainthing.Status,
		})
	})
	v1.Get("/my-domains", middleware.JWTAuth(db, cfg.JWT.Secret), func(c *fiber.Ctx) error {
		user := c.Locals("user").(middleware.User)
		var domains []FindDomains
		cursor, err := collection.Find(context.TODO(), bson.M{"addedby": user.ID.Hex()})
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to fetch domains from the database",
			})
		}
		defer cursor.Close(context.TODO())
		for cursor.Next(context.TODO()) {
			var domain FindDomains
			if err := cursor.Decode(&domain); err != nil {
				return c.JSON(fiber.Map{
					"success": false,
					"message": "Failed to decode domain",
				})
			}
			domains = append(domains, domain)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"domains": domains,
		})
	})
	v1.Get("/domains", func(c *fiber.Ctx) error {
		var domains []FindDomains
		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to fetch domains from the database",
			})
		}
		defer cursor.Close(context.TODO())
		for cursor.Next(context.TODO()) {
			var domain FindDomains
			if err := cursor.Decode(&domain); err != nil {
				return c.JSON(fiber.Map{
					"success": false,
					"message": "Failed to decode domain",
				})
			}
			domains = append(domains, domain)
		}
		if err := cursor.Err(); err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Cursor error",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(domains),
			"domains": domains,
		})
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.App.Port)))
}