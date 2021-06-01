package sheet

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type RawQueryResult struct {
	FullName           string
	Company            string
	QRCode             string
	EducationName      string
	EducationHours     string
	EducationDateStart string
	EducationDateEnd   string
	CertificateDate    string
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
		if len(row) < 8 {
			continue
		}

		fullName, ok1 := row[0].(string)
		company, ok2 := row[1].(string)
		educationName, ok3 := row[2].(string)
		educationHours, ok4 := row[3].(string)
		educationDateStart, ok5 := row[4].(string)
		educationDateEnd, ok6 := row[5].(string)
		referenceCode, ok7 := row[6].(string)
		certificateCreatedAt, ok8 := row[7].(string)

		if ok1 && ok2 && ok3 && ok4 && ok5 && ok6 && ok7 && ok8 {
			results = append(results, RawQueryResult{
				FullName:           fullName,
				QRCode:             referenceCode,
				EducationName:      educationName,
				Company:            company,
				EducationHours:     educationHours,
				EducationDateStart: educationDateStart,
				EducationDateEnd:   educationDateEnd,
				CertificateDate:    certificateCreatedAt,
			})
		} else {
			continue
		}
	}

	return results, nil
}
