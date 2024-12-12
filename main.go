package main

import (
	"fmt"
	"log"

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
	}

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
		if !lib.ValidateTld(domain.Domain) {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Invalid domain format or TLD",
			})
		}
		return c.JSON(fiber.Map{
			"meow":   "🐱",
			"domain": domain.Domain,
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
