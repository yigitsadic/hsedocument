package guard

import (
	"encoding/json"
	"net/http"
	"os"
)

type token struct {
	Token string `json:"token"`
}

var (
	unauthorizedResponse = map[string]string{
		"message": "unauthorized",
	}
)

func AuthenticationGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		expectedToken := os.Getenv("TOKEN")
		if expectedToken == "" {
			json.NewEncoder(writer).Encode(unauthorizedResponse)
			return
		}

		var t token

		if t.Token == expectedToken {
			next.ServeHTTP(writer, request.WithContext(request.Context()))
			return
		}

		json.NewEncoder(writer).Encode(unauthorizedResponse)
	})
}
