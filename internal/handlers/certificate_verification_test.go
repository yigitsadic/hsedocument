package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/yigitsadic/sertifikadogrula/internal/guard"
	"github.com/yigitsadic/sertifikadogrula/internal/store"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleCertificateVerificationNotFound(t *testing.T) {
	client := http.Client{}

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), guard.QRCodeCtxKey, "ABC")
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	})

	s := store.NewStore(nil)

	HandleCertificateVerification(r, s)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/certificate_verification", nil)
	if err != nil {
		t.Errorf("unexpected to see an error with this stage. error: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected to see an error with this stage. error: %s", err)
	}

	var responseObj store.QueryResult

	err = json.NewDecoder(res.Body).Decode(&responseObj)
	if err != nil {
		t.Errorf("unexpected to get error while decoding json. err=%s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code was %d", http.StatusOK)
	}

	if "ABC" != responseObj.QRCode {
		t.Errorf("expected to see %q as qr code but got=%q", "ABC", responseObj.QRCode)
	}

	if "not_verified" != responseObj.Status {
		t.Errorf("expected status was %q but got %q", "not_verified", responseObj.Status)
	}
}

func TestHandleCertificateVerificationNoCtxFound(t *testing.T) {
	client := http.Client{}

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			next.ServeHTTP(writer, request)
		})
	})

	s := store.NewStore(nil)

	HandleCertificateVerification(r, s)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/certificate_verification", nil)
	if err != nil {
		t.Errorf("unexpected to see an error with this stage. error: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected to see an error with this stage. error: %s", err)
	}

	var responseObj store.QueryResult

	err = json.NewDecoder(res.Body).Decode(&responseObj)
	if err != nil {
		t.Errorf("unexpected to get error while decoding json. err=%s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code was %d", http.StatusOK)
	}

	if defaultQRCode != responseObj.QRCode {
		t.Errorf("expected to see %q as qr code but got=%q", defaultQRCode, responseObj.QRCode)
	}

	if "not_verified" != responseObj.Status {
		t.Errorf("expected status was %q but got %q", "not_verified", responseObj.Status)
	}
}

func TestHandleCertificateVerificationFound(t *testing.T) {
	client := http.Client{}

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), guard.QRCodeCtxKey, "ABC")
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	})

	s := store.NewStore(nil)

	queryResult := store.QueryResult{
		Status:               "verified",
		MaskedFirstName:      "Yi***",
		MaskedLastName:       "***ıç",
		QRCode:               "ABC",
		CertificateName:      "İş Güvenliği",
		CertificateCreatedAt: "15-04-2021",
		LastUpdated:          time.Now().Format(store.DateTimeFormat),
	}
	s.QueryResults["ABC"] = &queryResult

	HandleCertificateVerification(r, s)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/certificate_verification", nil)
	if err != nil {
		t.Errorf("unexpected to see an error with this stage. error: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected to see an error with this stage. error: %s", err)
	}

	var responseObj store.QueryResult

	err = json.NewDecoder(res.Body).Decode(&responseObj)
	if err != nil {
		t.Errorf("unexpected to get error while decoding json. err=%s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code was %d", http.StatusOK)
	}

	if queryResult.QRCode != responseObj.QRCode {
		t.Errorf("expected to see %q as qr code but got=%q", queryResult.QRCode, responseObj.QRCode)
	}

	if queryResult.Status != responseObj.Status {
		t.Errorf("expected status was %q but got %q", queryResult.Status, responseObj.Status)
	}

	if queryResult.MaskedFirstName != responseObj.MaskedFirstName {
		t.Errorf("expected first name was=%q but got=%q", queryResult.MaskedFirstName, responseObj.MaskedFirstName)
	}

	if queryResult.MaskedLastName != responseObj.MaskedLastName {
		t.Errorf("expected first name was=%q but got=%q", queryResult.MaskedLastName, responseObj.MaskedLastName)
	}

	if queryResult.CertificateName != responseObj.CertificateName {
		t.Errorf("expected certificate name was=%q but got=%q", queryResult.CertificateName, responseObj.CertificateName)
	}
}
