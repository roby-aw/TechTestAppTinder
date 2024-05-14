package middlewares

import (
	"encoding/json"
	"os"
	"roby-backend-golang/utils"
	"strings"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/gofiber/fiber/v2"
)

func MiddleJWT(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")

	if !strings.Contains(authorizationHeader, "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "invalid token",
		})
	}

	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	secret := os.Getenv("JWT_SECRET")
	key, err := utils.Decode(secret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "invalid token",
		})
	}

	Header, _, err := jose.Decode(tokenString, key)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "invalid token",
		})
	}
	var str utils.JwtTokenClaimsUser
	err = json.Unmarshal([]byte(Header), &str)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "invalid token",
		})
	}
	timenow := time.Unix(str.Exp, 0)
	if timenow.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "invalid token",
		})

	}

	c.Locals("id", str.Sub)
	return c.Next()
}
