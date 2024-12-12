package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/opium-bio/config"
	"github.com/opium-bio/lib"
)

type Domain struct {
	Domain string `json:"domain"`
}

func main() {
	cfg, err := config.LoadConfig("./config.toml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(3)
	}
	db := lib.MongoDB()
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting MongoDB: %v", err)
			os.Exit(3)
		}
	}()

	fmt.Printf("Config loaded: %+v\n", cfg)
	app := fiber.New()
	v1 := app.Group("/v1")
	v1.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"meow": "🐱",
		})
	})
	v1.Post("/add", func(c fiber.Ctx) error {
		domain := new(Domain)
		if err := c.Bind().Body(domain); err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		if domain.Domain == "" {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Domain cannot be empty",
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
	address := fmt.Sprintf(":%d", cfg.App.Port)
	log.Fatal(app.Listen(address))
}
