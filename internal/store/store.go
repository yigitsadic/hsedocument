package store

import (
	"errors"
	"sync"
	"time"
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
func (s *Store) QueryInStore(qr_code string) (*QueryResult, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	result, ok := s.QueryResults[qr_code]

	if ok {
		return result, nil
	} else {
		return nil, QRCodeNotFoundErr
	}
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
