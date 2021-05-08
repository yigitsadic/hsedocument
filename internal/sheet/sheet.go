package sheet

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	CannotReadFromGoogleErr = errors.New("cannot read from google")
)

type RawQueryResult struct {
	FirstName            string
	LastName             string
	QRCode               string
	CertificateName      string
	CertificateCreatedAt string
}

type QueryClient interface {
	ReadSheetsAPI() ([]RawQueryResult, error)
}

type Client struct {
	SheetId string
	APIKey  string
}

func (c Client) ReadSheetsAPI() ([]RawQueryResult, error) {
	srv, err := sheets.NewService(
		context.TODO(),
		option.WithAPIKey(c.APIKey),
		option.WithScopes(sheets.SpreadsheetsReadonlyScope),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	readRange := "Sertifika Veritabanı!A:H"
	resp, err := srv.Spreadsheets.Values.Get(c.SheetId, readRange).Do()
	if err != nil {
		return nil, err
	}

	var results []RawQueryResult

	for _, row := range resp.Values[1:] {
		if len(row) < 5 {
			continue
		}

		firstName, ok1 := row[0].(string)
		lastName, ok2 := row[1].(string)
		certificateName, ok3 := row[4].(string)
		certificateDate, ok4 := row[6].(string)
		referenceCode, ok5 := row[7].(string)

		if ok1 && ok2 && ok3 && ok4 && ok5 {
			results = append(results, RawQueryResult{
				FirstName:            firstName,
				LastName:             lastName,
				QRCode:               referenceCode,
				CertificateName:      certificateName,
				CertificateCreatedAt: certificateDate,
			})
		} else {
			continue
		}
	}

	return results, nil
}
