package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

func ptr(s string) *string {
	return &s
}

func sendDiscordEmbed(cfg *config.Config, title, description, color string) {
	webhookURL := cfg.Discord.Webhook
	if webhookURL == "" {
		log.Println("🚨 Discord webhook URL is not set in config")
		return
	}

	payload := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{
			{
				Title:       ptr(title),
				Description: ptr(description),
				Color:       ptr(fmt.Sprintf("%d", colorHex(color))),
				Footer: &discordwebhook.Footer{
					Text: ptr("🌐 Domains Notification"),
				},
			},
		},
	}
	err := discordwebhook.SendMessage(webhookURL, payload)
	if err != nil {
		log.Println("❌ Failed to send embed to Discord:", err)
	}
}

func sendDiscordError(cfg *config.Config, title, description, color string) {
	sendDiscordEmbed(cfg, fmt.Sprintf("❌ **ERROR**: %s", title), description, color)
}

func sendDiscordSuccess(cfg *config.Config, title, description, color string) {
	sendDiscordEmbed(cfg, fmt.Sprintf("✅ **SUCCESS**: %s", title), description, color)
}

func colorHex(color string) int {
	colors := map[string]int{
		"red":    0xFF0000,
		"green":  0x00FF00,
		"blue":   0x0000FF,
		"yellow": 0xFFFF00,
		"purple": 0x800080,
		"orange": 0xFFA500,
		"cyan":   0x00FFFF,
	}
	if hex, exists := colors[color]; exists {
		return hex
	}
	return 0x2F3136
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
			sendDiscordError(cfg, "Error connecting to Cloudflare", fmt.Sprintf("%v", err), "orange")
		}
	}()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://opium.bio",
	}))
	v1 := app.Group("/v1")

	v1.Post("/add", middleware.JWTAuth(db, cfg.JWT.Secret), func(c *fiber.Ctx) error {
		user := c.Locals("user").(middleware.User)
		domain := new(Domain)

		if err := c.BodyParser(domain); err != nil {
			log.Println("Error parsing request body:", err)
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body: " + err.Error(),
			})
		}
		if domain.Domain == "" {
			log.Println("Domain cannot be empty")
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain cannot be empty",
			})
		}
		var existingDomain FindDomains
		err := collection.FindOne(context.Background(), bson.M{"domain": domain.Domain}).Decode(&existingDomain)
		if err == nil {
			log.Println("Domain already exists in the database:", domain.Domain)
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
			log.Println("Failed to save to MongoDB:", err)
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to save to MongoDB",
			})
		}

		err = lib.AddDomain(domain.Domain)
		if err != nil {
			log.Println("Failed to add domain to Cloudflare:", err)
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to add domain to Cloudflare. Please try again later.",
			})
		}

		sendDiscordSuccess(cfg, "Domain Added", fmt.Sprintf("🎉 A new domain has been added: `%s` 🚀", domain.Domain), "green")
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
			log.Println("Domain cannot be empty")
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain cannot be empty",
			})
		}
		var domainthing FindDomains
		err := collection.FindOne(context.TODO(), bson.M{"domain": domain}).Decode(&domainthing)
		if err != nil {
			log.Println("Domain doesn't exist in the database:", domain)
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain doesn't exist in the database",
			})
		}
		if domainthing.AddedBy != user.ID.Hex() {
			log.Println("Not your domain:", domain)
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
			log.Println("Failed to fetch domains from the database:", err)
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Failed to fetch domains from the database",
			})
		}
		defer cursor.Close(context.TODO())
		for cursor.Next(context.TODO()) {
			var domain FindDomains
			if err := cursor.Decode(&domain); err != nil {
				log.Println("Failed to decode domain:", err)
				return c.JSON(fiber.Map{
					"success": false,
					"message": "Failed to decode domain",
				})
			}
			domains = append(domains, domain)
		}
		if err := cursor.Err(); err != nil {
			log.Println("Cursor error:", err)
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
	sendDiscordSuccess(cfg, "Domains API Started", "🚀 API has started successfully", "green")
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.App.Port)))
}
