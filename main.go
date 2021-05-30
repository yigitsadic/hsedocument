package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/yigitsadic/hsedocument/internal/guard"
	"github.com/yigitsadic/hsedocument/internal/handlers"
	"github.com/yigitsadic/hsedocument/internal/sheet"
	"github.com/yigitsadic/hsedocument/internal/store"
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
		for _ = range s.Ticker.C {
			if err = s.FetchFromSheets(); err != nil {
				log.Println("Error occurred when fetching from Google Sheets. Err:", err)
			}
		}
	}()

	// Read port from ENV variable. If not found, default to 8080
	port, addr := os.Getenv("PORT"), os.Getenv("LISTEN_ADDR")
	if port == "" {
		port = "8080"
	}

	a := guard.Authentication{}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))
	r.Use(middleware.Logger)
	r.Use(a.Guard)
	handlers.HandleCertificateVerification(r, s)

	listenAddr := net.JoinHostPort(addr, port)

	log.Println("Server is starting on", listenAddr)
	// Start server.
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
