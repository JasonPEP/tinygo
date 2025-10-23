package test

import (
	"context"
	"testing"

	"tinygo/internal/shortener"
	"tinygo/internal/storage"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTempStore(t *testing.T) *storageTestAdapter {
	t.Helper()
	
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	
	// Auto migrate
	if err := db.AutoMigrate(&shortener.Link{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	
	store := storage.NewGormStore()
	// We need to set the database connection for testing
	// This is a bit of a hack, but works for testing
	store.SetDB(db)
	
	return &storageTestAdapter{Store: store}
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
		SetDB(db *gorm.DB)
	}
}
