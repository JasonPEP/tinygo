package test

import (
	"context"
	"path/filepath"
	"testing"

	"tinygo/internal/shortener"
	"tinygo/internal/storage"
)

func newTempStore(t *testing.T) *storageTestAdapter {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "links.json")
	fs, err := storage.NewFileStore(path)
	if err != nil {
		t.Fatalf("new filestore: %v", err)
	}
	return &storageTestAdapter{Store: fs, Path: path}
}

// Test basic shorten and resolve flow.
func TestService_SoftenAndResolve(t *testing.T) {
	st := newTempStore(t)
	svc := shortener.NewService(st.Store, "http://localhost:8080", 6)

	link, err := svc.Shorten(context.Background(), "https://golang.org", "")
	if err != nil {
		t.Fatalf("shorten: %v", err)
	}
	if link.Code == "" || link.LongURL != "https://golang.org" {
		t.Fatalf("unexpected link: %+v", link)
	}
	got, ok, err := svc.Resolve(context.Background(), link.Code)
	if err != nil || !ok {
		t.Fatalf("resolve: %v ok=%v", err, ok)
	}
	if got.LongURL != link.LongURL {
		t.Fatalf("resolve mismatch: %+v vs %+v", got, link)
	}
}

// --- local adapter (no cross-package export) ---
type storageTestAdapter struct {
	Store interface {
		Create(ctx context.Context, l shortener.Link) error
		Get(ctx context.Context, code string) (shortener.Link, bool, error)
		Delete(ctx context.Context, code string) error
		IncrementHit(ctx context.Context, code string) (shortener.Link, error)
		List(ctx context.Context) ([]shortener.Link, error)
	}
	Path string
}
