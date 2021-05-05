package guard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationGuard(t *testing.T) {
	client := http.Client{}
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	a := Authentication{Secret: "ABC"}

	ts := httptest.NewServer(a.Guard(hf))
	defer ts.Close()

	var response = struct {
		Message string `json:"message"`
	}{}

	t.Run("it should handle good token", func(t *testing.T) {
		req, err := buildRequest(t, ts.URL, "ABC")
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
	})

	t.Run("it should handle bad token", func(t *testing.T) {
		req, err := buildRequest(t, ts.URL, "DEF")
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

func buildRequest(t *testing.T, url, token string) (*http.Request, error) {
	t.Helper()

	b := bytes.Buffer{}
	json.NewEncoder(&b).Encode(map[string]string{
		"token":   token,
		"qr_code": "6f06ce23-7442-418c-886d-43af1f656808",
	})
	return http.NewRequest(http.MethodPost, url, &b)
}
