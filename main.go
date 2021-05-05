package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yigitsadic/sertifikadogrula/internal/guard"
	"github.com/yigitsadic/sertifikadogrula/internal/sheet"
	"github.com/yigitsadic/sertifikadogrula/internal/store"
	"log"
	"net"
	"net/http"
	"time"

	"os"
	"sync"
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

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	r.With(a.Guard).Post("/api/certificate_verification", func(writer http.ResponseWriter, request *http.Request) {
		errResponse := store.QueryResult{
			Status:               "not_verified",
			MaskedFirstName:      "",
			MaskedLastName:       "",
			QRCode:               "6f06ce23-7442-418c-886d-43af1f656808",
			CertificateName:      "",
			CertificateCreatedAt: "",
			LastUpdated:          time.Now(),
		}
		writer.Header().Set("Content-Type", "application/json")
		b := struct {
			QRCode string `json:"qr_code"`
		}{}

		err = json.NewDecoder(request.Body).Decode(&b)
		if err != nil {
			json.NewEncoder(writer).Encode(errResponse)
			return
		}

		res, err := s.QueryInStore(b.QRCode)
		if err != nil {
			json.NewEncoder(writer).Encode(errResponse)
			return
		}

		json.NewEncoder(writer).Encode(res)
	})

	listenAddr := net.JoinHostPort(addr, port)

	log.Println("Server is starting on", listenAddr)
	// Start server.
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
