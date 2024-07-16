package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type ProductDBField string

const (
	Id            ProductDBField = "id"
	CategoryID    ProductDBField = "category_id"
	Name          ProductDBField = "name"
	Description   ProductDBField = "description"
	Price         ProductDBField = "price"
	Stock         ProductDBField = "stock"
	CreatedBy     ProductDBField = "created_by"
	MetaCreatedAt ProductDBField = "meta_created_at"
	UpdatedBy     ProductDBField = "updated_by"
	MetaUpdatedAt ProductDBField = "meta_updated_at"
	DeletedBy     ProductDBField = "deleted_by"
	MetaDeletedAt ProductDBField = "meta_deleted_at"
)

type Product struct {
	ID            uuid.UUID   `db:"id"`
	CategoryID    uuid.UUID   `db:"category_id"`
	Name          string      `db:"name"`
	Description   string      `db:"description"`
	Price         float64     `db:"price"`
	Stock         int         `db:"stock"`
	CreatedBy     string      `db:"created_by"`
	MetaCreatedAt time.Time   `db:"meta_created_at"`
	UpdatedBy     string      `db:"updated_by"`
	MetaUpdatedAt time.Time   `db:"meta_updated_at"`
	DeletedBy     null.String `db:"deleted_by"`
	MetaDeletedAt null.Time   `db:"meta_deleted_at"`
}
