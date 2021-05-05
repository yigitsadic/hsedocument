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
		context.Background(),
		option.WithScopes(sheets.SpreadsheetsReadonlyScope),
		option.WithAPIKey(c.APIKey),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	readRange := "B:F"
	resp, err := srv.Spreadsheets.Values.Get(c.SheetId, readRange).Do()
	if err != nil {
		return nil, err
	}

	var result []RawQueryResult

	for _, row := range resp.Values {
		if len(row) < 5 {
			continue
		}

		firstName, ok1 := row[0].(string)
		lastName, ok2 := row[1].(string)
		certificateName, ok3 := row[2].(string)
		certificateDate, ok4 := row[3].(string)
		referenceCode, ok5 := row[4].(string)

		if ok1 && ok2 && ok3 && ok4 && ok5 {
			result = append(result, RawQueryResult{
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

	return result, nil
}
