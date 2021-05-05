package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yigitsadic/sertifikadogrula/internal/guard"
	"github.com/yigitsadic/sertifikadogrula/internal/sheet"
	"github.com/yigitsadic/sertifikadogrula/internal/store"
	"log"
	"os"
	"sync"
	"time"
)

var (
	s              *store.Store
	initializeOnce sync.Once
	err            error
)

func initialGet() {
	err = s.FetchFromSheets()
	log.Println("Error occurred when fetching from Google Sheets. Err:", err)
}

func main() {
	// Read port from ENV variable. If not found, default to 5050
	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}

	// Read Sheet ID from ENV.
	sheetId := os.Getenv("SHEET_ID")
	if sheetId == "" {
		log.Panic("Sheet ID cannot be empty")
	}

	s = store.NewStore(sheet.Client{
		SheetId: sheetId,
	})
	initializeOnce.Do(initialGet)

	app := fiber.New()

	go func() {
		for {
			select {
			case <-s.Ticker.C:
				if err = s.FetchFromSheets(); err != nil {
					log.Println("Error occurred when fetching from Google Sheets. Err:", err)
				}
			}
		}
	}()

	app.Use(guard.AuthenticationGuard)
	app.Post("/api/certificate_verification", func(ctx *fiber.Ctx) error {
		b := struct {
			QRCode string `json:"qr_code"`
		}{}

		err = ctx.BodyParser(&b)
		if err != nil {
			return err
		}

		res, err := s.QueryInStore(b.QRCode)

		if err != nil {
			if err == store.QRCodeNotFoundErr {
				return ctx.JSON(store.QueryResult{
					Status:               "not_verified",
					MaskedFirstName:      "",
					MaskedLastName:       "",
					QRCode:               b.QRCode,
					CertificateName:      "",
					CertificateCreatedAt: "",
					LastUpdated:          time.Now(),
				})
			}

			return err
		}

		return ctx.JSON(res)
	})

	app.Listen(fmt.Sprintf(":%s", port))
}
