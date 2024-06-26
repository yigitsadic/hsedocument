package guard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const qrCode = "6f06ce23-7442-418c-886d-43af1f656808"

func TestAuthenticationGuard(t *testing.T) {
	var readQrCode string

	client := http.Client{}
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, ok := r.Context().Value(QRCodeCtxKey).(string)
		if ok {
			readQrCode = v
		}

		w.WriteHeader(http.StatusOK)
	})

	a := Authentication{}

	ts := httptest.NewServer(a.Guard(hf))
	defer ts.Close()

	var response = struct {
		Message string `json:"message"`
	}{}

	t.Run("it should handle good token", func(t *testing.T) {
		req, err := buildRequest(t, ts.URL)
		if err != nil {
			t.Errorf("unexpected to see an error at this stage. err=%s", err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("unexpected to see an error at this stage. err=%s", err)
		}

		if res != nil && res.StatusCode != http.StatusOK {
			t.Errorf("expected to get 200 response")
		}

		if readQrCode != qrCode {
			t.Errorf("expected to qr code assigned to ctx")
		}
	})

	t.Run("it should handle empty body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
		if err != nil {
			t.Errorf("unexpected to see an error at this stage. err=%s", err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("unexpected to see an error at this stage. err=%s", err)
		}

		if res != nil && res.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected to get 401 response")
		}

		json.NewDecoder(res.Body).Decode(&response)

		if response.Message != "unauthorized" {
			t.Errorf("unexpected response. got=%+v", response)
		}
	})

}

func buildRequest(t *testing.T, url string) (*http.Request, error) {
	t.Helper()

	b := bytes.Buffer{}
	json.NewEncoder(&b).Encode(map[string]string{
		"qr_code": qrCode,
	})
	return http.NewRequest(http.MethodPost, url, &b)
}
