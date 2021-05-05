package store

import (
	"errors"
	"strings"
	"sync"
	"time"
)

const (
	BaseUrl = "https://hsegroup.uz/kurumsal/certificate_verification?qr_code="
)

var (
	QRCodeNotFoundErr = errors.New("QR Code not found")
)

type RawQueryResult struct {
	FirstName       string
	LastName        string
	QRCode          string
	CertificateName string
}

type QueryResult struct {
	MaskedFirstName string    `json:"first_name"`
	MaskedLastName  string    `json:"last_name"`
	QRCode          string    `json:"qr_code"`
	CertificateName string    `json:"certificate_name"`
	LastUpdated     time.Time `json:"last_updated"`
}

type Store struct {
	QueryResults map[string]*QueryResult
	Mu           *sync.Mutex

	Ticker *time.Ticker
}

// Reads given QR code from store.
func (s *Store) QueryInStore(qrCode string) (res *QueryResult, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	qrCode = strings.TrimSpace(qrCode)

	// Search with qr code only
	if result, ok := s.QueryResults[qrCode]; ok {
		res = result

		return
	}

	// Search with full path
	if result, ok := s.QueryResults[strings.TrimPrefix(qrCode, BaseUrl)]; ok {
		res = result

		return
	}

	return nil, QRCodeNotFoundErr
}

// Writes raw query result to store with masking names.
func (s *Store) WriteToStore(results []RawQueryResult) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, result := range results {
		s.QueryResults[result.QRCode] = &QueryResult{
			MaskedFirstName: result.FirstName,
			MaskedLastName:  result.LastName,
			QRCode:          result.QRCode,
			CertificateName: result.CertificateName,
			LastUpdated:     time.Now(),
		}
	}
}
