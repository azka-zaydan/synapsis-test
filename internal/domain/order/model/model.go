package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type OrderStatus int

var (
	OrderPlacedStatus OrderStatus = 0
	OrderPaidStatus   OrderStatus = 1
)

type Order struct {
	ID            uuid.UUID     `db:"id"`
	UserID        uuid.UUID     `db:"user_id"`
	PaymentID     uuid.NullUUID `db:"payment_id"`
	TotalPrice    float64       `db:"total_price"`
	Status        int           `db:"status"`
	OrderAt       time.Time     `db:"order_at"`
	PaymentAt     null.Time     `db:"payment_at"`
	CompletedAt   null.Time     `db:"completed_at"`
	CreatedBy     uuid.UUID     `db:"created_by"`
	MetaCreatedAt time.Time     `db:"meta_created_at"`
	UpdatedBy     uuid.UUID     `db:"updated_by"`
	MetaUpdatedAt time.Time     `db:"meta_updated_at"`
	DeletedBy     uuid.NullUUID `db:"deleted_by"`
	MetaDeletedAt null.Time     `db:"meta_deleted_at"`
}
