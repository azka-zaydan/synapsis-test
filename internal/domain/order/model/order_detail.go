package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type OrderDetail struct {
	ID                   uuid.UUID     `db:"id"`
	OrderID              uuid.UUID     `db:"order_id"`
	ProductID            uuid.UUID     `db:"product_id"`
	TotalItems           int           `db:"total_items"`
	SubtotalProductPrice float64       `db:"subtotal_product_price"`
	CreatedBy            uuid.UUID     `db:"created_by"`
	MetaCreatedAt        time.Time     `db:"meta_created_at"`
	UpdatedBy            uuid.UUID     `db:"updated_by"`
	MetaUpdatedAt        time.Time     `db:"meta_updated_at"`
	DeletedBy            uuid.NullUUID `db:"deleted_by"`
	MetaDeletedAt        null.Time     `db:"meta_deleted_at"`
}
