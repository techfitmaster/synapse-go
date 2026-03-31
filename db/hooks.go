package db

import (
	"time"

	"gorm.io/gorm"
)

// RegisterTimestampHooks registers GORM callbacks that auto-fill
// created_at on create and updated_at on create/update.
// Works with any model that has CreatedAt/UpdatedAt time.Time fields.
func RegisterTimestampHooks(db *gorm.DB) {
	_ = db.Callback().Create().Before("gorm:create").Register("synapse:timestamps:create", func(db *gorm.DB) {
		now := time.Now()
		if db.Statement.Schema != nil {
			if field := db.Statement.Schema.LookUpField("CreatedAt"); field != nil {
				if _, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
					_ = field.Set(db.Statement.Context, db.Statement.ReflectValue, now)
				}
			}
			if field := db.Statement.Schema.LookUpField("UpdatedAt"); field != nil {
				_ = field.Set(db.Statement.Context, db.Statement.ReflectValue, now)
			}
		}
	})

	_ = db.Callback().Update().Before("gorm:update").Register("synapse:timestamps:update", func(db *gorm.DB) {
		if db.Statement.Schema != nil {
			if field := db.Statement.Schema.LookUpField("UpdatedAt"); field != nil {
				db.Statement.SetColumn("UpdatedAt", time.Now())
			}
		}
	})
}
