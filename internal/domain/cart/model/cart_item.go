package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type CartItem struct {
	ID            uuid.UUID   `db:"id"`
	CartID        uuid.UUID   `db:"cart_id"`
	ProductID     uuid.UUID   `db:"product_id"`
	Quantity      int         `db:"quantity"`
	TotalPrice    float64     `db:"total_price"`
	CreatedBy     uuid.UUID   `db:"created_by"`
	MetaCreatedAt time.Time   `db:"meta_created_at"`
	UpdatedBy     uuid.UUID   `db:"updated_by"`
	MetaUpdatedAt time.Time   `db:"meta_updated_at"`
	DeletedBy     null.String `db:"deleted_by"`
	MetaDeletedAt null.Time   `db:"meta_deleted_at"`
}
