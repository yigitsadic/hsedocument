package guard

import (
	"context"
	"encoding/json"
	"net/http"
)

type jsonForm struct {
	QrCode string `json:"qr_code"`
}

var (
	unauthorizedResponse = map[string]string{
		"message": "unauthorized",
	}
)

const QRCodeCtxKey = "qrCode"

type Authentication struct {
}

func (a *Authentication) Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var f jsonForm
		json.NewDecoder(request.Body).Decode(&f)
		defer request.Body.Close()

		if f.QrCode != "" {
			ctx := context.WithValue(request.Context(), QRCodeCtxKey, f.QrCode)

			next.ServeHTTP(writer, request.WithContext(ctx))
			return
		}

		writer.WriteHeader(http.StatusUnauthorized)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(unauthorizedResponse)
	})
}
