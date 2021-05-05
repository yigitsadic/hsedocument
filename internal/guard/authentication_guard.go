package guard

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

type token struct {
	Token string `json:"token"`
}

func AuthenticationGuard(c *fiber.Ctx) error {
	expectedToken := os.Getenv("TOKEN")
	if expectedToken == "" {
		return errors.New("unable to continue")
	}

	var t token
	err := c.BodyParser(&t)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if t.Token == expectedToken {
		return c.Next()
	}

	return errors.New("unable to continue")
}
