package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type PaymentStatus int

var (
	Paid   PaymentStatus = 1
	Unpaid PaymentStatus = 0
)

type Payment struct {
	ID            uuid.UUID     `db:"id"`
	UserID        uuid.UUID     `db:"user_id"`
	OrderID       uuid.UUID     `db:"order_id"`
	PaymentMethod string        `db:"payment_method"`
	TotalPrice    float64       `db:"total_price"`
	Status        int           `db:"status"`
	PaymentAt     null.Time     `db:"payment_at"`
	CreatedBy     uuid.UUID     `db:"created_by"`
	MetaCreatedAt time.Time     `db:"meta_created_at"`
	UpdatedBy     uuid.UUID     `db:"updated_by"`
	MetaUpdatedAt time.Time     `db:"meta_updated_at"`
	DeletedBy     null.Time     `db:"deleted_by"`
	MetaDeletedAt uuid.NullUUID `db:"meta_deleted_at"`
}

func (m *Payment) Pay() {
	m.PaymentAt = null.TimeFrom(time.Now())
	m.Status = int(Paid)
}
