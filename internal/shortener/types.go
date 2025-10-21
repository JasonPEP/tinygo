package shortener

import (
	"time"

	"gorm.io/gorm"
)

// Link represents a shortened URL record.
type Link struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Code         string    `gorm:"uniqueIndex;size:32;not null" json:"code"`
	LongURL      string    `gorm:"size:2048;not null" json:"long_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	HitCount     int64     `gorm:"default:0" json:"hit_count"`
	LastAccessAt time.Time `json:"last_access_at"`
}

// TableName returns the table name for the Link model
func (Link) TableName() string {
	return "links"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (l *Link) BeforeCreate(tx *gorm.DB) error {
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now()
	}
	if l.UpdatedAt.IsZero() {
		l.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (l *Link) BeforeUpdate(tx *gorm.DB) error {
	l.UpdatedAt = time.Now()
	return nil
}
