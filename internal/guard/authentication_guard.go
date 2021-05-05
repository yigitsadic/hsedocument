package guard

import (
	"encoding/json"
	"net/http"
)

type token struct {
	Token string `json:"token"`
}

var (
	unauthorizedResponse = map[string]string{
		"message": "unauthorized",
	}
)

type Authentication struct {
	Secret string
}

func (a *Authentication) Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var t token
		json.NewDecoder(request.Body).Decode(&t)

		if t.Token == a.Secret {
			next.ServeHTTP(writer, request.WithContext(request.Context()))
			return
		}

		writer.WriteHeader(http.StatusUnauthorized)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(unauthorizedResponse)
	})
}
