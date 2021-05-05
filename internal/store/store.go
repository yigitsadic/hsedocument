package store

import (
	"errors"
	"sync"
	"time"
)

var (
	QRCodeNotFoundErr = errors.New("QR Code not found")
)

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
}

// Reads given QR code from store.
func (s *Store) QueryInStore(qr_code string) (*QueryResult, error) {
	result, ok := s.QueryResults[qr_code]

	if ok {
		return result, nil
	} else {
		return nil, QRCodeNotFoundErr
	}
}
