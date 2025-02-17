package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
}

func JWTAuth(db *mongo.Client, jwtKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("token")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized.",
			})
		}

		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil || !parsedToken.Valid {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"message": "Token expired or invalid.",
			})
		}

		if claims["exp"] != nil {
			expTime := int64(claims["exp"].(float64))
			if expTime < time.Now().Unix() {
				return c.Status(401).JSON(fiber.Map{
					"success": false,
					"message": "Token has expired.",
				})
			}
		}

		userID, ok := claims["id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token claims.",
			})
		}

		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid user ID format.",
			})
		}
		collection := db.Database("prod").Collection("users")
		var user User
		err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
		if err != nil {
			fmt.Println("Error decoding user:", err)
			if err == mongo.ErrNoDocuments {
				return c.Status(401).JSON(fiber.Map{
					"success": false,
					"message": "Unauthorized.",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "An internal server error has occurred.",
			})
		}
		c.Locals("user", user)
		return c.Next()
	}
}
