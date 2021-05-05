package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yigitsadic/sertifikadogrula/internal/guard"
	"github.com/yigitsadic/sertifikadogrula/internal/handlers"
	"github.com/yigitsadic/sertifikadogrula/internal/sheet"
	"github.com/yigitsadic/sertifikadogrula/internal/store"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

var (
	s              *store.Store
	initializeOnce sync.Once
	err            error
)

func initialGet() {
	if err = s.FetchFromSheets(); err != nil {
		log.Println("Error occurred when fetching from Google Sheets. Err:", err)
	}
}

func main() {
	// Read Sheet ID from ENV.
	sheetId := os.Getenv("SHEET_ID")
	if sheetId == "" {
		log.Fatalln("Sheet ID cannot be empty")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("TOKEN cannot be empty")
	}

	apiKey := os.Getenv("SHEET_API_KEY")
	if apiKey == "" {
		log.Fatalln("SHEET_API_KEY cannot be empty")
	}

	s = store.NewStore(sheet.Client{
		SheetId: sheetId,
		APIKey:  apiKey,
	})
	initializeOnce.Do(initialGet)

	// Every 1 hour make request to Google Sheets.
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

	// Read port from ENV variable. If not found, default to 8080
	port, addr := os.Getenv("PORT"), os.Getenv("LISTEN_ADDR")
	if port == "" {
		port = "8080"
	}

	a := guard.Authentication{Secret: token}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(a.Guard)
	handlers.HandleCertificateVerification(r, s)

	listenAddr := net.JoinHostPort(addr, port)

	log.Println("Server is starting on", listenAddr)
	// Start server.
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
