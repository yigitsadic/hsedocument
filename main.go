package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yigitsadic/sertifikadogrula/internal/store"
	"os"
	"time"
)

func main() {
	// Read port from ENV variable. If not found, default to 5050
	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}

	app := fiber.New()

	app.Post("/api/certificate_verification", func(ctx *fiber.Ctx) error {
		return ctx.JSON(store.QueryResult{
			MaskedFirstName: "Yi***",
			MaskedLastName:  "Sa***",
			QRCode:          "Ae1epOlMn",
			CertificateName: "İş Sağlığı ve Ergonomi",
			LastUpdated:     time.Now(),
		})
	})

	app.Listen(fmt.Sprintf(":%s", port))
}
