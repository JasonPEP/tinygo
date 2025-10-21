package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"tinygo/internal/shortener"
)

// ErrNotFound indicates code not exists.
var ErrNotFound = errors.New("link not found")

// fileStore stores links in memory with write-through JSON file.
type fileStore struct {
	mu    sync.RWMutex
	path  string
	links map[string]shortener.Link
}

type fileData struct {
	Links map[string]shortener.Link `json:"links"`
}

// NewFileStore creates or loads a file-backed store.
func NewFileStore(path string) (*fileStore, error) {
	fs := &fileStore{path: path, links: make(map[string]shortener.Link)}
	if err := fs.load(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *fileStore) load() error {
	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		// ensure dir exists
		if dir := filepath.Dir(s.path); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
		}
		return s.flush() // create empty file
	}
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return s.flush()
	}
	var fd fileData
	if err := json.Unmarshal(b, &fd); err != nil {
		return err
	}
	if fd.Links == nil {
		fd.Links = make(map[string]shortener.Link)
	}
	s.links = fd.Links
	return nil
}

func (s *fileStore) flush() error {
	s.mu.RLock()
	fd := fileData{Links: s.links}
	s.mu.RUnlock()

	tmp := s.path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(fd); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Create saves a new link. Returns error if code exists.
func (s *fileStore) Create(ctx context.Context, l shortener.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.links[l.Code]; ok {
		return fmt.Errorf("code already exists: %s", l.Code)
	}
	now := time.Now()
	if l.CreatedAt.IsZero() {
		l.CreatedAt = now
	}
	l.UpdatedAt = now
	s.links[l.Code] = l
	return s.flush()
}

// Get returns a link by code.
func (s *fileStore) Get(ctx context.Context, code string) (shortener.Link, bool, error) {
	s.mu.RLock()
	l, ok := s.links[code]
	s.mu.RUnlock()
	return l, ok, nil
}

// Delete removes a link by code.
func (s *fileStore) Delete(ctx context.Context, code string) error {
	s.mu.Lock()
	if _, ok := s.links[code]; !ok {
		s.mu.Unlock()
		return ErrNotFound
	}
	delete(s.links, code)
	s.mu.Unlock()
	return s.flush()
}

// IncrementHit increases hit counter and updates last access time.
func (s *fileStore) IncrementHit(ctx context.Context, code string) (shortener.Link, error) {
	s.mu.Lock()
	l, ok := s.links[code]
	if !ok {
		s.mu.Unlock()
		return shortener.Link{}, ErrNotFound
	}
	l.HitCount++
	l.LastAccessAt = time.Now()
	l.UpdatedAt = l.LastAccessAt
	s.links[code] = l
	s.mu.Unlock()
	if err := s.flush(); err != nil {
		return shortener.Link{}, err
	}
	return l, nil
}

// List returns all links.
func (s *fileStore) List(ctx context.Context) ([]shortener.Link, error) {
	s.mu.RLock()
	result := make([]shortener.Link, 0, len(s.links))
	for _, l := range s.links {
		result = append(result, l)
	}
	s.mu.RUnlock()
	return result, nil
}
