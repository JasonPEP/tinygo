package storage

import (
	"context"
	"errors"

	"tinygo/internal/database"
	"tinygo/internal/shortener"

	"gorm.io/gorm"
)

// gormStore implements Store interface using GORM
type gormStore struct {
	db *gorm.DB
}

// NewGormStore creates a new GORM-based store
func NewGormStore() *gormStore {
	return &gormStore{
		db: database.DB,
	}
}

// Create saves a new link. Returns error if code exists.
func (s *gormStore) Create(ctx context.Context, l shortener.Link) error {
	result := s.db.WithContext(ctx).Create(&l)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("code already exists: " + l.Code)
		}
		return result.Error
	}
	return nil
}

// Get returns a link by code.
func (s *gormStore) Get(ctx context.Context, code string) (shortener.Link, bool, error) {
	var l shortener.Link
	result := s.db.WithContext(ctx).Where("code = ?", code).First(&l)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return shortener.Link{}, false, nil
		}
		return shortener.Link{}, false, result.Error
	}
	return l, true, nil
}

// Delete removes a link by code.
func (s *gormStore) Delete(ctx context.Context, code string) error {
	result := s.db.WithContext(ctx).Where("code = ?", code).Delete(&shortener.Link{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// IncrementHit increases hit counter and updates last access time.
func (s *gormStore) IncrementHit(ctx context.Context, code string) (shortener.Link, error) {
	var l shortener.Link

	// First, get the current record
	result := s.db.WithContext(ctx).Where("code = ?", code).First(&l)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return shortener.Link{}, ErrNotFound
		}
		return shortener.Link{}, result.Error
	}

	// Update hit count and last access time
	updates := map[string]interface{}{
		"hit_count":      gorm.Expr("hit_count + 1"),
		"last_access_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}

	result = s.db.WithContext(ctx).Model(&l).Updates(updates)
	if result.Error != nil {
		return shortener.Link{}, result.Error
	}

	// Get the updated record
	result = s.db.WithContext(ctx).Where("code = ?", code).First(&l)
	if result.Error != nil {
		return shortener.Link{}, result.Error
	}

	return l, nil
}

// List returns all links.
func (s *gormStore) List(ctx context.Context) ([]shortener.Link, error) {
	var links []shortener.Link
	result := s.db.WithContext(ctx).Find(&links)
	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}
