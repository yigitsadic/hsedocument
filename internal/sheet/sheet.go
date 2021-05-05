package sheet

import "errors"

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
}
