package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Cart struct {
	ID            uuid.UUID   `db:"id"`
	UserID        uuid.UUID   `db:"user_id"`
	TotalPrice    float64     `db:"total_price"`
	TotalItems    int         `db:"total_items"`
	CreatedBy     uuid.UUID   `db:"created_by"`
	MetaCreatedAt time.Time   `db:"meta_created_at"`
	UpdatedBy     uuid.UUID   `db:"updated_by"`
	MetaUpdatedAt time.Time   `db:"meta_updated_at"`
	DeletedBy     null.String `db:"deleted_by"`
	MetaDeletedAt null.Time   `db:"meta_deleted_at"`
}
