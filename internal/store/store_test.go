package store

import (
	"sync"
	"testing"
	"time"
)

func TestStore_QueryInStore(t *testing.T) {
	s := &Store{
		QueryResults: make(map[string]*QueryResult),
		Mu:           &sync.Mutex{},
	}

	r1 := QueryResult{
		MaskedFirstName: "Yi***",
		MaskedLastName:  "Sa***",
		QRCode:          "ygt",
		LastUpdated:     time.Now(),
	}

	r2 := QueryResult{
		MaskedFirstName: "Ay***",
		MaskedLastName:  "Ço***",
		QRCode:          "aycn",
		LastUpdated:     time.Now(),
	}

	s.QueryResults[r1.QRCode] = &r1
	s.QueryResults[r2.QRCode] = &r2

	t.Run("it should return if finds", func(a *testing.T) {
		got, err := s.QueryInStore(r1.QRCode)
		if err != nil {
			t.Errorf("unexpected to get an error but got=%s", err)
		}

		if got.MaskedFirstName != r1.MaskedFirstName {
			t.Errorf("expected return not satisfied for masked first name")
		}

		if got.MaskedLastName != r1.MaskedLastName {
			t.Errorf("expected return not satisfied for masked last name")
		}

		if got.QRCode != r1.QRCode {
			t.Errorf("expected return not satisfied for QR code")
		}
	})

	t.Run("it should give error if cannot find", func(t *testing.T) {
		got, err := s.QueryInStore("ABCD")

		if got != nil {
			t.Errorf("unexpected to find")
		}

		if err != QRCodeNotFoundErr {
			t.Errorf("expected error was=%s but got=%s", QRCodeNotFoundErr, err)
		}
	})
}
