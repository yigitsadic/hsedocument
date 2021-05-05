package store

import (
	"errors"
	"github.com/yigitsadic/sertifikadogrula/internal/name_masker"
	"github.com/yigitsadic/sertifikadogrula/internal/sheet"
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

type QueryResult struct {
	Status               string    `json:"status"`
	MaskedFirstName      string    `json:"first_name"`
	MaskedLastName       string    `json:"last_name"`
	QRCode               string    `json:"qr_code"`
	CertificateName      string    `json:"certificate_name"`
	CertificateCreatedAt string    `json:"certificate_created_at"`
	LastUpdated          time.Time `json:"last_updated"`
}

type Store struct {
	QueryResults map[string]*QueryResult
	Mu           *sync.Mutex

	Client sheet.QueryClient
	Ticker *time.Ticker
}

func NewStore(client sheet.QueryClient) *Store {
	return &Store{
		QueryResults: make(map[string]*QueryResult),
		Mu:           &sync.Mutex{},
		Ticker:       time.NewTicker(1 * time.Hour),
		Client:       client,
	}
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
func (s *Store) WriteToStore(results []sheet.RawQueryResult) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, result := range results {
		s.QueryResults[result.QRCode] = &QueryResult{
			Status:               "verified",
			MaskedFirstName:      name_masker.MaskFirstName(result.FirstName),
			MaskedLastName:       name_masker.MaskLastName(result.LastName),
			QRCode:               result.QRCode,
			CertificateName:      result.CertificateName,
			CertificateCreatedAt: result.CertificateCreatedAt,
			LastUpdated:          time.Now(),
		}
	}
}

// Queries over given client and writes them to store.
func (s *Store) FetchFromSheets() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	got, err := s.Client.ReadSheetsAPI()
	if err != nil {
		return err
	}

	s.WriteToStore(got)

	return nil
}
