package store

import (
	"errors"
	"github.com/yigitsadic/hsedocument/internal/sheet"
	"testing"
	"time"
)

var sampleError = errors.New("simple read error")

type errorClient struct {
}

func (m errorClient) ReadSheetsAPI() ([]sheet.RawQueryResult, error) {
	return nil, sampleError
}

type successfulClient struct {
}

func (m successfulClient) ReadSheetsAPI() ([]sheet.RawQueryResult, error) {
	return []sheet.RawQueryResult{
		{
			FullName:        "Aycan Çotoy",
			QRCode:          "ABC",
			EducationName:   "Lorem",
			CertificateDate: "2021-03-15",
		},
		{
			FullName:        "Yiğit Sadıç",
			QRCode:          "DEF",
			EducationName:   "Ipsum",
			CertificateDate: "2021-04-08",
		},
	}, nil
}

func TestStore_QueryInStore(t *testing.T) {
	s := NewStore(nil)

	r1 := QueryResult{
		MaskedFullName: "Yi*** ***ıç",
		QRCode:         "ygt",
		LastUpdated:    time.Now().Format(DateTimeFormat),
	}

	r2 := QueryResult{
		MaskedFullName: "Ay*** ***oy",
		QRCode:         "aycn",
		LastUpdated:    time.Now().Format(DateTimeFormat),
	}

	s.QueryResults[r1.QRCode] = &r1
	s.QueryResults[r2.QRCode] = &r2

	t.Run("it should trim whitespaces", func(t *testing.T) {
		got, err := s.QueryInStore("  " + r1.QRCode + "  ")

		if err != nil {
			t.Errorf("unexpected to get an error but got=%s", err)
		}

		if got.MaskedFullName != r1.MaskedFullName {
			t.Errorf("expected return not satisfied for masked full name")
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

		if got.MaskedFullName != r1.MaskedFullName {
			t.Errorf("expected return not satisfied for masked full name")
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

		if got.MaskedFullName != r1.MaskedFullName {
			t.Errorf("expected return not satisfied for masked full name")
		}

		if got.QRCode != r1.QRCode {
			t.Errorf("expected return not satisfied for QR code")
		}
	})
}

func TestStore_WriteToStore(t *testing.T) {
	input := []sheet.RawQueryResult{
		{
			FullName:        "Aycan Çotoy",
			QRCode:          "ABC",
			EducationName:   "Lorem",
			CertificateDate: "2021-03-15",
		},
		{
			FullName:        "Yiğit Sadıç",
			QRCode:          "DEF",
			EducationName:   "Ipsum",
			CertificateDate: "2021-04-08",
		},
	}

	expectedResult := make(map[string]*QueryResult)
	expectedResult["ABC"] = &QueryResult{
		MaskedFullName:  "Ay*** ***oy",
		QRCode:          "ABC",
		CertificateName: "Lorem",
	}
	expectedResult["DEF"] = &QueryResult{
		MaskedFullName:  "Yi*** ***ıç",
		QRCode:          "DEF",
		CertificateName: "Ipsum",
	}

	s := NewStore(nil)

	s.WriteToStore(input)

	for k, v := range s.QueryResults {
		if expectedResult[k].MaskedFullName != v.MaskedFullName {
			t.Errorf("expected to see %s but got %s", expectedResult[k].MaskedFullName, v.MaskedFullName)
		}

		if expectedResult[k].QRCode != v.QRCode {
			t.Errorf("expected to see %s but got %s", expectedResult[k].QRCode, v.QRCode)
		}

		if expectedResult[k].CertificateName != v.CertificateName {
			t.Errorf("expected to see %s but got %s", expectedResult[k].CertificateName, v.CertificateName)
		}
	}
}

func TestStore_FetchFromSheets(t *testing.T) {
	t.Run("it should handle error gracefully", func(t *testing.T) {
		s := NewStore(errorClient{})

		if err := s.FetchFromSheets(); err != sampleError {
			t.Errorf("expected to get an error while reading")
		}

		if len(s.QueryResults) != 0 {
			t.Errorf("unpexected to write any result to query results")
		}
	})

	t.Run("it should write to store", func(t *testing.T) {
		s := NewStore(successfulClient{})

		if err := s.FetchFromSheets(); err != nil {
			t.Errorf("unexpected to see an error but got %s", err)
		}

		// expected length 2

		if len(s.QueryResults) != 2 {
			t.Errorf("expected to see 2 records to be written in query results")
		}

		if a, ok := s.QueryResults["ABC"]; ok {
			if a.QRCode != "ABC" {
				t.Errorf("expected qr code reference not found")
			}
		} else {
			t.Errorf("expected to see ABC reference")
		}

		if a, ok := s.QueryResults["DEF"]; ok {
			if a.QRCode != "DEF" {
				t.Errorf("expected qr code reference not found")
			}
		} else {
			t.Errorf("expected to see ABC reference")
		}
	})
}
