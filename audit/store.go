package audit

import (
	"context"

	"gorm.io/gorm"
)

// GORMStore implements Store using GORM.
type GORMStore struct {
	db *gorm.DB
}

// NewGORMStore creates a Store backed by the given GORM database.
func NewGORMStore(db *gorm.DB) *GORMStore {
	return &GORMStore{db: db}
}

// Save persists an audit entry using a new database connection.
func (s *GORMStore) Save(ctx context.Context, entry *Entry) error {
	return s.db.WithContext(ctx).Create(entry).Error
}

// SaveInTx persists an audit entry within an existing transaction.
func (s *GORMStore) SaveInTx(tx *gorm.DB, entry *Entry) error {
	return tx.Create(entry).Error
}
