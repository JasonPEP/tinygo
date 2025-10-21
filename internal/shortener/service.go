package shortener

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"tinygo/pkg/random"
)

var (
	codeRegexp     = regexp.MustCompile(`^[0-9A-Za-z_-]{3,32}$`)
	ErrInvalidURL  = errors.New("invalid url")
	ErrInvalidCode = errors.New("invalid code")
)

// Service contains business logic for creating and resolving short links.
type Service struct {
	store      Store
	codeLength int
	baseURL    string
	maxRetry   int
}

// NewService creates a shortener Service.
func NewService(store Store, baseURL string, codeLength int) *Service {
	return &Service{
		store:      store,
		codeLength: codeLength,
		baseURL:    baseURL,
		maxRetry:   5,
	}
}

// Shorten creates a short link optionally with a custom code.
func (s *Service) Shorten(ctx context.Context, longURL, customCode string) (Link, error) {
	if !isValidURL(longURL) {
		return Link{}, ErrInvalidURL
	}
	var code string
	if customCode != "" {
		if !codeRegexp.MatchString(customCode) {
			return Link{}, ErrInvalidCode
		}
		code = customCode
	} else {
		var err error
		code, err = random.Code(s.codeLength)
		if err != nil {
			return Link{}, fmt.Errorf("generate code: %w", err)
		}
	}

	l := Link{Code: code, LongURL: longURL}
	// If code exists, retry generate when not custom.
	for i := 0; i < s.maxRetry; i++ {
		if err := s.store.Create(ctx, l); err != nil {
			if customCode != "" {
				return Link{}, err
			}
			// re-generate and retry
			c, gerr := random.Code(s.codeLength)
			if gerr != nil {
				return Link{}, fmt.Errorf("regenerate code: %w", gerr)
			}
			l.Code = c
			continue
		}
		return l, nil
	}
	return Link{}, fmt.Errorf("exceeded retries to create short link")
}

// Resolve returns link by code without mutating stats.
func (s *Service) Resolve(ctx context.Context, code string) (Link, bool, error) {
	return s.store.Get(ctx, code)
}

// Hit increments hit counter and returns updated link.
func (s *Service) Hit(ctx context.Context, code string) (Link, error) {
	return s.store.IncrementHit(ctx, code)
}

// Delete removes a link.
func (s *Service) Delete(ctx context.Context, code string) error {
	return s.store.Delete(ctx, code)
}

// ShortURL builds absolute short URL.
func (s *Service) ShortURL(code string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, code)
}

// List returns all links.
func (s *Service) List(ctx context.Context) ([]Link, error) {
	return s.store.List(ctx)
}

func isValidURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

// Now is extracted for testing override when needed.
var Now = func() time.Time { return time.Now() }
