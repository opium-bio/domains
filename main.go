package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/opium-bio/config"
)

func main() {
	cfg, err := config.LoadConfig("./config.toml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("Config loaded: %+v\n", cfg)
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"meow": "🐱",
		})
	})
	address := fmt.Sprintf(":%d", cfg.App.Port)
	log.Fatal(app.Listen(address))
}
