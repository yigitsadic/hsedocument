package store

import (
	"testing"
	"time"
)

func TestStore_QueryInStore(t *testing.T) {
	s := NewStore()

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

	t.Run("it should trim whitespaces", func(t *testing.T) {
		got, err := s.QueryInStore("  " + r1.QRCode + "  ")

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

	t.Run("it should find with full url", func(t *testing.T) {
		got, err := s.QueryInStore(BaseUrl + r1.QRCode)
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
}

func TestStore_WriteToStore(t *testing.T) {
	input := []RawQueryResult{
		{
			FirstName:       "Aycan",
			LastName:        "Çotoy",
			QRCode:          "ABC",
			CertificateName: "Lorem",
		},
		{
			FirstName:       "Yiğit",
			LastName:        "Sadıç",
			QRCode:          "DEF",
			CertificateName: "Ipsum",
		},
	}

	expectedResult := make(map[string]*QueryResult)
	expectedResult["ABC"] = &QueryResult{
		MaskedFirstName: "Ay***",
		MaskedLastName:  "***oy",
		QRCode:          "ABC",
		CertificateName: "Lorem",
	}
	expectedResult["DEF"] = &QueryResult{
		MaskedFirstName: "Yi***",
		MaskedLastName:  "***ıç",
		QRCode:          "DEF",
		CertificateName: "Ipsum",
	}

	s := NewStore()

	s.WriteToStore(input)

	for k, v := range s.QueryResults {
		if expectedResult[k].MaskedFirstName != v.MaskedFirstName {
			t.Errorf("expected to see %s but got %s", expectedResult[k].MaskedFirstName, v.MaskedFirstName)
		}

		if expectedResult[k].MaskedLastName != v.MaskedLastName {
			t.Errorf("expected to see %s but got %s", expectedResult[k].MaskedLastName, v.MaskedLastName)
		}

		if expectedResult[k].QRCode != v.QRCode {
			t.Errorf("expected to see %s but got %s", expectedResult[k].QRCode, v.QRCode)
		}

		if expectedResult[k].CertificateName != v.CertificateName {
			t.Errorf("expected to see %s but got %s", expectedResult[k].CertificateName, v.CertificateName)
		}
	}
}
