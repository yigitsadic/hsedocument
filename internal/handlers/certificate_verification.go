package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/yigitsadic/sertifikadogrula/internal/guard"
	"github.com/yigitsadic/sertifikadogrula/internal/store"
	"net/http"
)

const (
	defaultQRCode = "c9492e5d-44eb-44f5-932b-0898ef52f48b"
)

func errorResponse(qrCode, lastUpdated string) store.QueryResult {
	if qrCode == "" {
		qrCode = defaultQRCode
	}

	return store.QueryResult{
		Status:               "not_verified",
		MaskedFirstName:      "",
		MaskedLastName:       "",
		QRCode:               qrCode,
		CertificateName:      "",
		CertificateCreatedAt: "",
		LastUpdated:          lastUpdated,
	}
}

func HandleCertificateVerification(r *chi.Mux, s *store.Store) {
	r.Post("/api/certificate_verification", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		v, ok := request.Context().Value(guard.QRCodeCtxKey).(string)
		if !ok {
			json.NewEncoder(writer).Encode(errorResponse("", s.LastUpdated))
			return
		}

		res, err := s.QueryInStore(v)
		if err != nil {
			json.NewEncoder(writer).Encode(errorResponse(v, s.LastUpdated))
			return
		}

		res.LastUpdated = s.LastUpdated
		json.NewEncoder(writer).Encode(res)
	})
}
