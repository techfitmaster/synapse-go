package audit

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Entry represents a single audit log record.
// The Detail field (JSON string) allows business-specific extensions
// without changing the shared schema.
type Entry struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `gorm:"type:varchar(30);not null" json:"action"`
	Resource  string    `gorm:"type:varchar(200)" json:"resource"`
	Detail    string    `gorm:"type:text" json:"detail"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Store is the interface for persisting audit entries.
type Store interface {
	// Save persists an audit entry using a new database connection.
	Save(ctx context.Context, entry *Entry) error
	// SaveInTx persists an audit entry within an existing transaction.
	SaveInTx(tx *gorm.DB, entry *Entry) error
}
